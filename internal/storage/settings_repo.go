package storage

import (
	"context"

	"github.com/loginovm/anti-bruteforce/internal/storage/models"
)

type SettingsRepo interface {
	GetSettings(ctx context.Context) (models.Setting, error)
	UpdateSettings(ctx context.Context, setting models.Setting) error
}
