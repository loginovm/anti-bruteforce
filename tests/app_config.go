//nolint:unused
package tests

type appConfig struct {
	AppURL string
	DB     datasource
}

type datasource struct {
	Host          string
	Port          string
	Username      string
	Password      string
	Name          string
	MigrationsDir string
}

var appCfg = appConfig{
	AppURL: "localhost:8088",
	DB: datasource{
		Host:          "localhost",
		Port:          "5532",
		Username:      "otus_user",
		Password:      "dev_pass",
		Name:          "antibforce",
		MigrationsDir: "migrations",
	},
}
