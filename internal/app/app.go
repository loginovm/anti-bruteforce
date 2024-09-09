package app

import (
	"github.com/loginovm/anti-bruteforce/internal/ratelimit"
)

type App struct {
	limit        int
	loginCalc    ratelimit.Calc
	passwordCalc ratelimit.Calc
	ipCalc       ratelimit.Calc
}

func New(limit int, loginCalc, passwordCalc, ipCalc ratelimit.Calc) *App {
	return &App{
		limit:        limit,
		loginCalc:    loginCalc,
		passwordCalc: passwordCalc,
		ipCalc:       ipCalc,
	}
}

func (app *App) CheckAttempt(login, password, ip string) bool {
	isAttemptValid := app.loginCalc.TryIncrement(login, app.limit) &&
		app.passwordCalc.TryIncrement(password, app.limit) &&
		app.ipCalc.TryIncrement(ip, app.limit)

	return isAttemptValid
}
