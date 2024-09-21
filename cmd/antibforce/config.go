package main

type Config struct {
	App struct {
		URL string `toml:"url"`
	} `toml:"app"`
	Datasource struct {
		Host          string `toml:"host"`
		Port          string `toml:"port"`
		Username      string `toml:"user"`
		Password      string `toml:"password"`
		Name          string `toml:"db-name"`
		Ssl           string `toml:"ssl"`
		MigrationsDir string `toml:"migrations-dir"`
	} `toml:"datasource"`
}
