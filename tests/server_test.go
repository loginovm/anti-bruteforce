//go:build integration

package tests

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"os/exec"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/loginovm/anti-bruteforce/internal/app"
	"github.com/loginovm/anti-bruteforce/internal/ratelimit"
	api "github.com/loginovm/anti-bruteforce/internal/server/http"
	"github.com/loginovm/anti-bruteforce/internal/server/http/models"
	"github.com/loginovm/anti-bruteforce/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestLoginAttempts(t *testing.T) {
	testAttemptsCount := 3
	ctx, cancel := context.WithCancel(context.Background())
	defer func() { cleanUp(cancel) }()
	resetCounterAfter := 2 * time.Second // n,m,k attempts count are reset after 2 sec
	err := runServer(ctx, appCfg, resetCounterAfter)
	require.NoError(t, err)

	tests := []struct {
		name             string
		loginAttempts    int
		passwordAttempts int
		ipAttempts       int
		login            string
		password         string
		IP               string
	}{
		{name: "Login test", loginAttempts: testAttemptsCount, passwordAttempts: math.MaxInt32, ipAttempts: math.MaxInt32, login: "userA", IP: "1.1.1.1"},
		{name: "Password test", loginAttempts: math.MaxInt32, passwordAttempts: testAttemptsCount, ipAttempts: math.MaxInt32, password: "pass123", IP: "1.1.1.2"},
		{name: "IP test", loginAttempts: math.MaxInt32, passwordAttempts: math.MaxInt32, ipAttempts: testAttemptsCount, IP: "10.10.10.10"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err = setSettings(tc.loginAttempts, tc.passwordAttempts, tc.ipAttempts)
			require.NoError(t, err)

			req := models.CheckLoginAttemptRequest{Login: tc.login, Password: tc.password, IP: tc.IP}
			// First testAttemptsCount login attempts must be allowed
			testLoginAttempts(t, testAttemptsCount, true, req)
			// Next login attempt is forbidden
			testLoginAttempts(t, 1, false, req)

			// After configured time period elapsed login attempts count is reset
			time.Sleep(resetCounterAfter)
			// And first testAttemptsCount login attempts again are allowed
			testLoginAttempts(t, testAttemptsCount, true, req)
			// Next login attempt is forbidden
			testLoginAttempts(t, 1, false, req)
		})
	}
}

func testLoginAttempts(t *testing.T, count int, expectedResult bool, req models.CheckLoginAttemptRequest) {
	url := getServerURL(appCfg.AppURL, "check-login-attempt")
	var resp models.CheckAttemptResponse
	for i := 0; i < count; i++ {
		err := httpReq(http.MethodPut, url, req, &resp)
		require.NoError(t, err)
		require.Equal(t, expectedResult, resp.IsValid)
	}
}

func setSettings(loginAttempts, passwordAttempts, ipAttempts int) error {
	url := getServerURL(appCfg.AppURL, "settings")
	expected := models.Setting{
		LoginCount:    loginAttempts,
		PasswordCount: passwordAttempts,
		IPCount:       ipAttempts,
	}
	return httpReq(http.MethodPut, url, expected, nil)
}

func TestBlackList(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() { cleanUp(cancel) }()
	err := runServerDefault(ctx, appCfg)
	require.NoError(t, err)

	err = setSettings(100, 100, 100)
	require.NoError(t, err)

	login := "TestBlacklist"
	ip := "200.100.1.2"
	// Login attempts are allowed
	req := models.CheckLoginAttemptRequest{Login: login, Password: "123", IP: ip}
	testLoginAttempts(t, 5, true, req)

	// Add IP to blacklist
	url := getServerURL(appCfg.AppURL, "blacklist")
	err = httpReq(http.MethodPost, url, models.CidrRequest{Cidr: ip}, nil)
	require.NoError(t, err)

	// All login attempts are forbidden
	req = models.CheckLoginAttemptRequest{Login: login, Password: "123", IP: ip}
	testLoginAttempts(t, 5, false, req)
}

func TestWhiteList(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() { cleanUp(cancel) }()
	err := runServerDefault(ctx, appCfg)
	require.NoError(t, err)

	err = setSettings(1, 100, 100)
	require.NoError(t, err)

	login := "TestWhitelist"
	ip := "200.200.1.2"
	// LoginAttempts = 1. First attempt is allowed only
	req := models.CheckLoginAttemptRequest{Login: login, Password: "123", IP: ip}
	testLoginAttempts(t, 1, true, req)
	testLoginAttempts(t, 5, false, req)

	// Add IP to whitelist
	url := getServerURL(appCfg.AppURL, "whitelist")
	err = httpReq(http.MethodPost, url, models.CidrRequest{Cidr: ip}, nil)
	require.NoError(t, err)

	// All login attempts are allowed
	req = models.CheckLoginAttemptRequest{Login: login, Password: "123", IP: ip}
	testLoginAttempts(t, 5, true, req)
}

func TestUpdateSettings(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() { cleanUp(cancel) }()
	err := runServerDefault(ctx, appCfg)
	require.NoError(t, err)
	url := getServerURL(appCfg.AppURL, "settings")
	expected := models.Setting{
		LoginCount:    1000,
		PasswordCount: 5000,
		IPCount:       10000,
	}
	// update settings
	err = httpReq(http.MethodPut, url, expected, nil)
	require.NoError(t, err)

	var actual models.Setting
	// get settings
	err = httpGet(url, &actual)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func runServerDefault(ctx context.Context, cfg appConfig) error {
	return runServer(ctx, cfg, 60*time.Second)
}

func runServer(
	ctx context.Context,
	cfg appConfig,
	resetCounterAfter time.Duration,
) error {
	store, err := createStorage(ctx, cfg)
	if err != nil {
		return err
	}
	app := createApp(store, resetCounterAfter)
	server := api.NewServer(cfg.AppURL, app, "")
	var stopping atomic.Bool
	go func() {
		<-ctx.Done()
		stopping.Store(true)
		log.Println("Stopping server")
		if err = server.Stop(); err != nil {
			log.Fatal("failed to stop server: " + err.Error())
		}
	}()

	go func() {
		log.Println("anti-bruteforce is running...")
		if err = server.Start(ctx); err != nil {
			if !stopping.Load() {
				log.Fatal("failed to start server: \n" + err.Error())
			}
		}
		store.Close()
	}()

	return nil
}

func createApp(repo storage.Repo, resetCounterAfter time.Duration) *app.App {
	ipBucketStore := &ratelimit.MemStorage{}
	loginBucketStore := &ratelimit.MemStorage{}
	passBucketStore := &ratelimit.MemStorage{}
	ipCalc := ratelimit.NewCalc(resetCounterAfter, ipBucketStore)
	loginCalc := ratelimit.NewCalc(resetCounterAfter, loginBucketStore)
	passCalc := ratelimit.NewCalc(resetCounterAfter, passBucketStore)
	app := app.New(loginCalc, passCalc, ipCalc, repo)

	return app
}

func createStorage(ctx context.Context, cfg appConfig) (*storage.Storage, error) {
	if err := exec.Command("sh", "./create_test_db.sh").Run(); err != nil {
		return nil, err
	}
	time.Sleep(time.Millisecond * 1000)
	config := cfg.DB
	s := storage.New()
	err := s.Connect(ctx,
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			config.Host, config.Port, config.Username, config.Password, config.Name, "disable"))
	if err != nil {
		return nil, fmt.Errorf("connect to db: %w", err)
	}
	if err = s.RunMigration(config.MigrationsDir); err != nil {
		return nil, err
	}

	return s, nil
}

func cleanUp(cancel context.CancelFunc) {
	log.Println("Test cleanup start")
	cancel()
	if err := exec.Command("sh", "./remove_test_db.sh").Run(); err != nil {
		log.Println(err)
	}
}

func getServerURL(host string, path string) string {
	url := "http://" + strings.Trim(host, "/") + "/" + strings.Trim(path, "/")

	return url
}
