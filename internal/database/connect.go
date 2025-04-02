package database

import (
    "context"
    "fmt"
    "log"

    "github.com/jackc/pgx/v4/pgxpool"
    "EffectiveMobile/config"
)

// InitDB инициализирует подключение к базе данных
func InitDB(cfg *config.Config) (*pgxpool.Pool, error) {
    pool, err := pgxpool.Connect(context.Background(), cfg.DB.ConnString())
    if err != nil {
        return nil, fmt.Errorf("unable to connect to database: %w", err)
    }

    log.Println("Connected to database")
    return pool, nil
}
