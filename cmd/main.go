package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/m04kA/SMK-UserService/internal/config"
	"github.com/m04kA/SMK-UserService/internal/handlers/api/create_car"
	"github.com/m04kA/SMK-UserService/internal/handlers/api/create_user"
	"github.com/m04kA/SMK-UserService/internal/handlers/api/delete_car"
	"github.com/m04kA/SMK-UserService/internal/handlers/api/delete_current_user"
	"github.com/m04kA/SMK-UserService/internal/handlers/api/get_current_user"
	"github.com/m04kA/SMK-UserService/internal/handlers/api/update_car"
	"github.com/m04kA/SMK-UserService/internal/handlers/api/update_current_user"
	"github.com/m04kA/SMK-UserService/internal/handlers/middleware"
	carrepo "github.com/m04kA/SMK-UserService/internal/infra/storage/car"
	userrepo "github.com/m04kA/SMK-UserService/internal/infra/storage/user"
	userservice "github.com/m04kA/SMK-UserService/internal/service/user"
)

func main() {

	// Загружаем конфигурацию
	cfg, err := config.Load("./config.toml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Подключаемся к базе данных
	db, err := sqlx.Connect("postgres", cfg.Database.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Successfully connected to database")

	// Инициализируем репозитории
	userRepo := userrepo.NewRepository(db)
	carRepo := carrepo.NewRepository(db)

	// Инициализируем сервис
	service := userservice.NewUserService(userRepo, carRepo)

	// Инициализируем handlers
	createUserHandler := create_user.NewHandler(service)
	getCurrentUserHandler := get_current_user.NewHandler(service)
	updateCurrentUserHandler := update_current_user.NewHandler(service)
	deleteCurrentUserHandler := delete_current_user.NewHandler(service)
	createCarHandler := create_car.NewHandler(service)
	updateCarHandler := update_car.NewHandler(service)
	deleteCarHandler := delete_car.NewHandler(service)

	// Инициализируем middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWT.Secret)

	// Настраиваем роутер
	r := mux.NewRouter()

	// Public routes
	r.HandleFunc("/users", createUserHandler.Handle).Methods(http.MethodPost)

	// Protected routes
	protected := r.PathPrefix("").Subrouter()
	protected.Use(authMiddleware.JWTAuth)

	protected.HandleFunc("/users/me", getCurrentUserHandler.Handle).Methods(http.MethodGet)
	protected.HandleFunc("/users/me", updateCurrentUserHandler.Handle).Methods(http.MethodPut)
	protected.HandleFunc("/users/me", deleteCurrentUserHandler.Handle).Methods(http.MethodDelete)

	protected.HandleFunc("/users/me/cars", createCarHandler.Handle).Methods(http.MethodPost)
	protected.HandleFunc("/users/me/cars/{car_id}", updateCarHandler.Handle).Methods(http.MethodPatch)
	protected.HandleFunc("/users/me/cars/{car_id}", deleteCarHandler.Handle).Methods(http.MethodDelete)

	// Создаем HTTP сервер
	addr := fmt.Sprintf(":%d", cfg.Server.HTTPPort)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Starting server on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Ожидаем сигнал завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
