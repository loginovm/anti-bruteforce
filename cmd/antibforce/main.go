package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/loginovm/anti-bruteforce/internal/app"
	"github.com/loginovm/anti-bruteforce/internal/ratelimit"
	"github.com/loginovm/anti-bruteforce/internal/server/http"
	"github.com/loginovm/anti-bruteforce/internal/storage"
	_ "github.com/loginovm/anti-bruteforce/swagger/docs"
)

const (
	checkLimitPeriod    = 60 * time.Second // Period after which (ip/login/password)counter is reset
	recycleBucketPeriod = 3 * time.Minute
	maxBucketAge        = 10 * time.Minute
)

// Args command-line parameters.
type Args struct {
	ConfigPath string
}

func main() {
	var config Config
	args := ProcessArgs(&config)
	// read configuration from the file and environment variables
	if err := cleanenv.ReadConfig(args.ConfigPath, &config); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	store, err := CreateStorage(ctx, config)
	if err != nil {
		cancel()
		log.Fatal(err) //nolint:gocritic
	}
	defer store.Close()

	app := createApp(store)
	server := http.NewServer(config.App.URL, app, config.App.Swagger)

	go func() {
		<-ctx.Done()
		if err = server.Stop(); err != nil {
			log.Fatal("failed to stop server: " + err.Error())
		}
	}()

	log.Println("anti-bruteforce is running...")
	if err = server.Start(ctx); err != nil {
		log.Fatal("failed to start server: " + err.Error())
	}
}

func createApp(repo storage.Repo) *app.App {
	ipBucketStore := &ratelimit.MemStorage{}
	loginBucketStore := &ratelimit.MemStorage{}
	passBucketStore := &ratelimit.MemStorage{}
	ipCalc := ratelimit.NewCalc(checkLimitPeriod, ipBucketStore)
	loginCalc := ratelimit.NewCalc(checkLimitPeriod, loginBucketStore)
	passCalc := ratelimit.NewCalc(checkLimitPeriod, passBucketStore)
	app := app.New(loginCalc, passCalc, ipCalc, repo)

	go func() {
		for {
			ipBucketStore.Clean(time.Now(), maxBucketAge)
			loginBucketStore.Clean(time.Now(), maxBucketAge)
			passBucketStore.Clean(time.Now(), maxBucketAge)
			time.Sleep(recycleBucketPeriod)
		}
	}()

	return app
}

func CreateStorage(ctx context.Context, cfg Config) (*storage.Storage, error) {
	config := cfg.Datasource
	s := storage.New()
	err := s.Connect(ctx,
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			config.Host, config.Port, config.Username, config.Password, config.Name, config.Ssl))
	if err != nil {
		return nil, err
	}
	if err = s.RunMigration(config.MigrationsDir); err != nil {
		return nil, err
	}

	return s, nil
}

func ProcessArgs(cfg interface{}) Args {
	var a Args

	f := flag.NewFlagSet("Calendar app", 1)
	f.StringVar(&a.ConfigPath, "config", "config.toml", "Path to configuration file")

	// Embed config descriptions into command help
	fu := f.Usage
	f.Usage = func() {
		fu()
		envHelp, _ := cleanenv.GetDescription(cfg, nil)
		fmt.Fprintln(f.Output())
		fmt.Fprintln(f.Output(), envHelp)
	}

	f.Parse(os.Args[1:])
	return a
}
