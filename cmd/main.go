package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	_ "sub/cmd/docs"
	"sub/internal/config"
	database "sub/internal/db"
	"sub/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/swaggo/http-swagger"
)

// @title User Subscription API
// @version 1.0
// @description API для управления подписками пользователей
// @host localhost:8080
// @BasePath /
func main() {
	cfg := config.InitConf()

	db, err := database.ConnectDb(cfg.DbKey)
	if err != nil {
		log.Fatal(err)
	}

	r := setupRouter(db)

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.Port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	slog.Info("Server started on port " + strconv.Itoa(cfg.Port))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	slog.Info("Server exiting")
}

func setupRouter(db *sqlx.DB) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Post("/api/user", handlers.CreateUser(db))
	r.Get("/api/user/{id}", handlers.GetUser(db))
	r.Put("/api/user/{id}", handlers.UpdateUser(db))
	r.Get("/api/user", handlers.GetAllUsers(db))
	r.Delete("/api/user/{id}", handlers.DeleteUser(db))
	r.Get("/api/total-price", handlers.GetTotalPrice(db))

	return r
}
