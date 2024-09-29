package storage

import (
	"context"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/stdlib" // driver for sql db
	"github.com/jmoiron/sqlx"
	"github.com/loginovm/anti-bruteforce/internal/storage/models"
	"github.com/pressly/goose/v3"
)

var (
	_ BWListRepo   = (*Storage)(nil)
	_ SettingsRepo = (*Storage)(nil)
)

type Storage struct {
	db *sqlx.DB
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) GetWList(ctx context.Context) ([]models.BWListItem, error) {
	var items []models.BWListItem
	err := s.db.SelectContext(ctx, &items, `SELECT id, cidr, created_at FROM white_list`)

	return items, err
}

func (s *Storage) GetBList(ctx context.Context) ([]models.BWListItem, error) {
	var items []models.BWListItem
	err := s.db.SelectContext(ctx, &items, `SELECT id, cidr, created_at FROM black_list`)

	return items, err
}

func (s *Storage) AddWLItem(ctx context.Context, item models.BWListItem) error {
	query := `INSERT INTO white_list(cidr, created_at)
VALUES (:cidr, :created_at)`

	_, err := s.db.NamedExecContext(ctx, query, item)
	return err
}

func (s *Storage) DeleteWLItem(ctx context.Context, id int) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM white_list WHERE id = $1`, id)
	return err
}

func (s *Storage) AddBLItem(ctx context.Context, item models.BWListItem) error {
	query := `INSERT INTO black_list(cidr, created_at)
VALUES (:cidr, :created_at)`

	_, err := s.db.NamedExecContext(ctx, query, item)
	return err
}

func (s *Storage) DeleteBLItem(ctx context.Context, id int) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM black_list WHERE id = $1`, id)
	return err
}

func (s *Storage) GetSettings(ctx context.Context) (models.Setting, error) {
	var items []models.Setting
	err := s.db.SelectContext(ctx, &items, `
SELECT 
    id, 
    login_count, 
    password_count,
    ip_count 
FROM settings`)
	if err != nil {
		return models.Setting{}, err
	}
	if len(items) == 0 {
		return models.Setting{}, errors.New("settings is empty")
	}

	return items[0], err
}

func (s *Storage) UpdateSettings(ctx context.Context, setting models.Setting) error {
	query := `UPDATE settings
SET login_count = :login_count,
    password_count = :password_count,
    ip_count = :ip_count
WHERE id = :id
`
	_, err := s.db.NamedExecContext(ctx, query, setting)

	return err
}

func (s *Storage) RunMigration(migrationsDir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("cannot set dialect: %w", err)
	}

	if err := goose.Up(s.db.DB, migrationsDir); err != nil {
		return fmt.Errorf("cannot do up migration: %w", err)
	}

	return nil
}

func (s *Storage) Connect(ctx context.Context, dsn string) error {
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("cannot open pgx driver: %w", err)
	}

	s.db = db
	return db.PingContext(ctx)
}

func (s *Storage) Close() error {
	return s.db.Close()
}
