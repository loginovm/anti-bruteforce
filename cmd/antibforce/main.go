package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/loginovm/anti-bruteforce/internal/storage"
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

	Print(ctx, store)
}

func CreateStorage(ctx context.Context, cfg Config) (storage.Repo, error) {
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
	f.StringVar(&a.ConfigPath, "config", "/etc/calendar/config.toml", "Path to configuration file")

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

func Print(ctx context.Context, store storage.Repo) {
	wl, err := store.GetWList(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("White list: %+v\n", wl)
	bl, err := store.GetBList(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Black list: %+v\n", bl)
	s, err := store.GetSettings(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Settings: %+v\n", s)
}
