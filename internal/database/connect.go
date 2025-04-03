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
        log.Println("Database does not exist. Attempting to create...")
        err = CreateDataBase(cfg)
        if err != nil {
            return nil, fmt.Errorf("unable to create a database: %w", err)
        }

        // Вторая попытка
        pool, err = pgxpool.Connect(context.Background(), cfg.DB.ConnString())
        if err != nil {
            return nil, fmt.Errorf("unable to create a database: %w", err)
        }
    }

    log.Println("Connected to database")
    return pool, nil
}
