package storage

type Repo interface {
	BWListRepo
	SettingsRepo
	Close() error
}
