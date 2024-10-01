package main

type Config struct {
	App struct {
		URL     string `toml:"url"`
		Swagger string `toml:"swagger"`
	} `toml:"app"`
	Datasource struct {
		Host          string `toml:"host" env:"DB_HOST"`
		Port          string `toml:"port" env:"DB_PORT"`
		Username      string `toml:"user" env:"DB_USER"`
		Password      string `toml:"password" env:"DB_PASS"`
		Name          string `toml:"db-name" env:"DB_NAME"`
		Ssl           string `toml:"ssl" env:"DB_SSL"`
		MigrationsDir string `toml:"migrations-dir" env:"DB_MIGRATIONS"`
	} `toml:"datasource"`
}
