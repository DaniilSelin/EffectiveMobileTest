package main

import (                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 
	"log"
	"context"
	"os"
	"os/signal"
	"time"
	"net/http"
	"fmt"

	"EffectiveMobile/internal/transport/http/api"
	"EffectiveMobile/config"
	"EffectiveMobile/internal/database"
	"EffectiveMobile/internal/logger"
	"EffectiveMobile/internal/repository"
	"EffectiveMobile/internal/service"
)

func main() {
	ctx := context.Background()

	ctx, err := logger.New(ctx)
	if err != nil {
		log.Fatalf("Error create logger: %v", err)
	}

	cfg, err := config.LoadConfig("config/config.yml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// 1. Подключаемся к БД
	dbPool, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer dbPool.Close()

	// 2.. Запускаем миграции
	err = database.RunMigrations(ctx, cfg, dbPool)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// 3. Создаем репозитории
	personRepo := repository.NewPersonRepository(dbPool, cfg)

	// 4. Создаем сервисы
	personService := service.NewPersonService(personRepo)

	// 5. Создаем хэндлер
	handler := api.NewHandler(ctx, personService)

	// 6. Создаём роутер
	router := api.NewRouter(ctx, handler)

	// 7. Запускаем сервер
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting server on %s...", addr)

	srv := &http.Server{
		Addr: addr,
		Handler: router,
	}

	go func () {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	// 8. Завршаем работу сервер (Graceful Shutdown)
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server")
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	srv.Shutdown(ctx)
	log.Println("Server gracefull stopped")
}
