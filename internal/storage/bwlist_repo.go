package storage

import (
	"context"

	"github.com/loginovm/anti-bruteforce/internal/storage/models"
)

type BWListRepo interface {
	GetWList(ctx context.Context) ([]models.BWListItem, error)
	GetBList(ctx context.Context) ([]models.BWListItem, error)
	AddWLItem(ctx context.Context, item models.BWListItem) error
	DeleteWLItem(ctx context.Context, id int) error
	AddBLItem(ctx context.Context, item models.BWListItem) error
	DeleteBLItem(ctx context.Context, id int) error
}
