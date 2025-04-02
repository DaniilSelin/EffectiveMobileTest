package database

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"EffectiveMobile/config"
)

func RunMigrations(ctx context.Context, cfg config.Config, conn *pgxpool.Pool) error {
	files := []string{
		"internal/database/migrations/create_user.sql",
	}
	for _, file := range files {
		sqlContent, err := ioutil.ReadFile(file)
		sqlContent := fmt.Sprintf(string(sqlContent), cfg.DB.Schema, cfg.DB.Schema)
		_, err = conn.Exec(ctx, sqlContent)

		if err != nil {
			return fmt.Errorf("failed to read SQL file %s: %w", file, err)
		}

		_, err = conn.Exec(ctx, string(sqlContent))
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}

		log.Printf("Successfully executed migration: %s", file)
	}

	return nil
}
