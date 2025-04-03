package database

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"EffectiveMobile/config"
)

func RunMigrations(ctx context.Context, cfg *config.Config, conn *pgxpool.Pool) error {
	files := []string{
		"internal/database/migrations/create_person.sql",
	}
	for _, file := range files {
		sqlContent, err := ioutil.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read SQL file %s: %w", file, err)
		}

		sqlQuery := fmt.Sprintf(string(sqlContent), cfg.DB.Schema, cfg.DB.Schema)

		_, err = conn.Exec(ctx, sqlQuery)
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}

		log.Printf("Successfully executed migration: %s", file)
	}

	return nil
}

// CreateDataBase создаёт базу данных через системную БД `postgres`
func CreateDataBase(cfg *config.Config) error {
    sysDsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=%s",
        cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Sslmode)

    conn, err := pgx.Connect(context.Background(), sysDsn)
    if err != nil {
        return fmt.Errorf("failed to connect to system database: %w", err)
    }
    defer conn.Close(context.Background())

    // случай если ошибка была не по отсутствию бд
    var exists bool
    err = conn.QueryRow(context.Background(), "SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = $1)", cfg.DB.Dbname).Scan(&exists)
    if err != nil {
        return fmt.Errorf("failed to check database existence: %w", err)
    }

    if !exists {
        _, err = conn.Exec(context.Background(), fmt.Sprintf(`CREATE DATABASE %s;`, cfg.DB.Dbname))
        if err != nil {
            return fmt.Errorf("failed to create database: %w", err)
        }
        log.Printf("Database %s created successfully", cfg.DB.Dbname)
    } else {
        log.Printf("Database %s already exists", cfg.DB.Dbname)
    }

    return nil
}