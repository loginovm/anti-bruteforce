package app

import (
	"context"
	"fmt"
	"net/netip"
	"strings"
	"time"

	"github.com/loginovm/anti-bruteforce/internal/ratelimit"
	"github.com/loginovm/anti-bruteforce/internal/storage"
	"github.com/loginovm/anti-bruteforce/internal/storage/models"
)

var (
	errIPInvalid   = NewAppError("value is not a valid IP address")
	errCidrInvalid = NewAppError("value is not a valid IP/CIDR")
)

type App struct {
	loginCalc    *ratelimit.Calc
	passwordCalc *ratelimit.Calc
	ipCalc       *ratelimit.Calc
	storage      storage.Repo
}

func New(
	loginCalc *ratelimit.Calc,
	passwordCalc *ratelimit.Calc,
	ipCalc *ratelimit.Calc,
	storage storage.Repo,
) *App {
	return &App{
		loginCalc:    loginCalc,
		passwordCalc: passwordCalc,
		ipCalc:       ipCalc,
		storage:      storage,
	}
}

func (app *App) CheckLoginAttempt(ctx context.Context, login, password, ip string) (bool, error) {
	if !isValidIP(ip) {
		return false, fmt.Errorf("%w : %s", errIPInvalid, ip)
	}
	blist, err := app.GetBList(ctx)
	if err != nil {
		return false, err
	}
	isInBList, err := isIPInList(ip, blist)
	if err != nil {
		return false, err
	}
	if isInBList {
		return false, nil
	}
	wlist, err := app.GetWList(ctx)
	if err != nil {
		return false, err
	}
	isInWList, err := isIPInList(ip, wlist)
	if err != nil {
		return false, err
	}
	if isInWList {
		return true, nil
	}
	settings, err := app.storage.GetSettings(ctx)
	if err != nil {
		return false, err
	}
	isAttemptValid := app.loginCalc.TryIncrement(login, settings.LoginCount) &&
		app.passwordCalc.TryIncrement(password, settings.PasswordCount) &&
		app.ipCalc.TryIncrement(ip, settings.IPCount)

	return isAttemptValid, nil
}

func (app *App) ResetIPBucket(ip string) {
	app.ipCalc.ResetBucket(ip)
}

func (app *App) ResetLoginBucket(login string) {
	app.loginCalc.ResetBucket(login)
}

func (app *App) AddIPToBlackList(ctx context.Context, ip string) error {
	if !isValidIP(ip) {
		return fmt.Errorf("%w : %s", errCidrInvalid, ip)
	}
	ip = normalizeIP(ip)
	blist, err := app.GetBList(ctx)
	if err != nil {
		return err
	}
	isIPAdded := app.isIPAlreadyAdded(ip, blist)
	if !isIPAdded {
		if err = app.storage.AddBLItem(ctx, models.BWListItem{Cidr: ip, CreatedAt: time.Now()}); err != nil {
			return err
		}
	}

	return nil
}

func (app *App) AddIPToWhiteList(ctx context.Context, ip string) error {
	if !isValidIP(ip) {
		return fmt.Errorf("%w : %s", errCidrInvalid, ip)
	}
	ip = normalizeIP(ip)
	wlist, err := app.GetWList(ctx)
	if err != nil {
		return err
	}
	isIPAdded := app.isIPAlreadyAdded(ip, wlist)
	if !isIPAdded {
		if err = app.storage.AddWLItem(ctx, models.BWListItem{Cidr: ip, CreatedAt: time.Now()}); err != nil {
			return err
		}
	}

	return nil
}

func (app *App) DeleteIPFromBlackList(ctx context.Context, ip string) error {
	ip = normalizeIP(ip)
	blist, err := app.GetBList(ctx)
	if err != nil {
		return err
	}
	for _, v := range blist {
		if ip == v.Cidr {
			if err = app.storage.DeleteBLItem(ctx, v.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func (app *App) DeleteIPFromWhiteList(ctx context.Context, ip string) error {
	ip = normalizeIP(ip)
	wlist, err := app.GetWList(ctx)
	if err != nil {
		return err
	}
	for _, v := range wlist {
		if ip == v.Cidr {
			if err = app.storage.DeleteWLItem(ctx, v.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func (app *App) GetWList(ctx context.Context) ([]models.BWListItem, error) {
	whiteList, err := app.storage.GetWList(ctx)
	if err != nil {
		return nil, err
	}

	return whiteList, nil
}

func (app *App) GetBList(ctx context.Context) ([]models.BWListItem, error) {
	blackList, err := app.storage.GetBList(ctx)
	if err != nil {
		return nil, err
	}

	return blackList, nil
}

func (app *App) GetSettings(ctx context.Context) (models.Setting, error) {
	s, err := app.storage.GetSettings(ctx)
	if err != nil {
		return models.Setting{}, err
	}

	return s, nil
}

func (app *App) UpdateSettings(ctx context.Context, s models.Setting) error {
	if err := app.storage.UpdateSettings(ctx, s); err != nil {
		return err
	}

	return nil
}

func isIPInList(ip string, list []models.BWListItem) (bool, error) {
	ipAddr, err := netip.ParseAddr(ip)
	if err != nil {
		return false, err
	}
	for _, v := range list {
		network, err := netip.ParsePrefix(v.Cidr)
		if err != nil {
			return false, err
		}
		if network.Contains(ipAddr) {
			return true, nil
		}
	}

	return false, nil
}

func (app *App) isIPAlreadyAdded(ip string, list []models.BWListItem) bool {
	for _, v := range list {
		if ip == v.Cidr {
			return true
		}
	}

	return false
}

func isValidIP(ip string) bool {
	if strings.ContainsRune(ip, '/') {
		_, err := netip.ParsePrefix(ip)
		return err == nil
	}

	_, err := netip.ParseAddr(ip)
	return err == nil
}

func normalizeIP(ip string) string {
	if !strings.ContainsRune(ip, '/') {
		ip += "/32"
	}

	return ip
}
