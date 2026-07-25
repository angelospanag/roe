package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/angelospanag/roe/db/migrations"
	"github.com/angelospanag/roe/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: migrate <up|down|status>")
		os.Exit(1)
	}

	if err := run(os.Args[1]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(cmd string) error {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config error: %w", err)
	}

	sqlDB, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("unable to open database: %w", err)
	}
	defer func() { _ = sqlDB.Close() }()

	goose.SetBaseFS(migrations.FS)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	switch cmd {
	case "up":
		return goose.Up(sqlDB, ".")
	case "down":
		return goose.Down(sqlDB, ".")
	case "status":
		return goose.Status(sqlDB, ".")
	default:
		return fmt.Errorf("unknown command: %s (want up, down, or status)", cmd)
	}
}
