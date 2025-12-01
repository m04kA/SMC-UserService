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
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/m04kA/SMC-UserService/internal/config"
	"github.com/m04kA/SMC-UserService/internal/handlers/api/create_car"
	"github.com/m04kA/SMC-UserService/internal/handlers/api/create_user"
	"github.com/m04kA/SMC-UserService/internal/handlers/api/delete_car"
	"github.com/m04kA/SMC-UserService/internal/handlers/api/delete_current_user"
	"github.com/m04kA/SMC-UserService/internal/handlers/api/get_current_user"
	"github.com/m04kA/SMC-UserService/internal/handlers/api/get_selected_car"
	"github.com/m04kA/SMC-UserService/internal/handlers/api/get_user_by_id"
	"github.com/m04kA/SMC-UserService/internal/handlers/api/select_car"
	"github.com/m04kA/SMC-UserService/internal/handlers/api/update_car"
	"github.com/m04kA/SMC-UserService/internal/handlers/api/update_current_user"
	"github.com/m04kA/SMC-UserService/internal/handlers/middleware"
	carrepo "github.com/m04kA/SMC-UserService/internal/infra/storage/car"
	userrepo "github.com/m04kA/SMC-UserService/internal/infra/storage/user"
	userservice "github.com/m04kA/SMC-UserService/internal/service/user"
	"github.com/m04kA/SMC-UserService/pkg/logger"
)

func main() {
	// Инициализируем логгер
	if err := logger.Init("./logs/app.log"); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Close()

	logger.Info("Starting SMC-UserService...")

	// Загружаем конфигурацию
	cfg, err := config.Load("./config.toml")
	if err != nil {
		logger.Fatal("Failed to load config: %v", err)
	}
	logger.Info("Configuration loaded successfully")

	// Подключаемся к базе данных
	db, err := sqlx.Connect("postgres", cfg.Database.DSN())
	if err != nil {
		logger.Fatal("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		logger.Fatal("Failed to ping database: %v", err)
	}
	logger.Info("Successfully connected to database")

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
	getSelectedCarHandler := get_selected_car.NewHandler(service)
	selectCarHandler := select_car.NewHandler(service)
	getUserByIDHandler := get_user_by_id.NewHandler(service)

	// Настраиваем роутер
	r := mux.NewRouter()

	// Применяем metrics middleware ко всем роутам
	r.Use(middleware.Metrics)

	// Metrics endpoint
	r.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)

	// Public routes
	r.HandleFunc("/users", createUserHandler.Handle).Methods(http.MethodPost)

	// Internal routes (для межсервисного взаимодействия)
	r.HandleFunc("/internal/users/{tg_user_id}", getUserByIDHandler.Handle).Methods(http.MethodGet)
	r.HandleFunc("/internal/users/{tg_user_id}/cars/selected", getSelectedCarHandler.Handle).Methods(http.MethodGet)

	// Protected routes (требуют заголовок X-User-ID)
	protected := r.PathPrefix("").Subrouter()
	protected.Use(middleware.UserIDAuth)

	protected.HandleFunc("/users/me", getCurrentUserHandler.Handle).Methods(http.MethodGet)
	protected.HandleFunc("/users/me", updateCurrentUserHandler.Handle).Methods(http.MethodPut)
	protected.HandleFunc("/users/me", deleteCurrentUserHandler.Handle).Methods(http.MethodDelete)

	protected.HandleFunc("/users/me/cars", createCarHandler.Handle).Methods(http.MethodPost)
	protected.HandleFunc("/users/me/cars/{car_id}", updateCarHandler.Handle).Methods(http.MethodPatch)
	protected.HandleFunc("/users/me/cars/{car_id}", deleteCarHandler.Handle).Methods(http.MethodDelete)
	protected.HandleFunc("/users/me/cars/{car_id}/select", selectCarHandler.Handle).Methods(http.MethodPut)

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
		logger.Info("Starting server on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start: %v", err)
		}
	}()

	// Ожидаем сигнал завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown: %v", err)
	}

	logger.Info("Server stopped gracefully")
}
