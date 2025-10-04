# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

SMK-UserService - сервис управления пользователями и их автомобилями для приложения автомойки. Использует Clean Architecture с разделением на domain, service, infrastructure и handlers.

## Development Commands

### Docker Environment
```bash
# Запуск PostgreSQL и миграций
docker-compose up -d

# Остановка
docker-compose down

# Пересоздание с удалением данных
docker-compose down -v && docker-compose up -d
```

### Build and Run
```bash
go run cmd/main.go
```

### Testing
```bash
go test ./...
```

### Database
- PostgreSQL работает на порту **5435** (не стандартный 5432)
- Connection string: `host=localhost port=5435 user=postgres password=postgres dbname=smk_userservice sslmode=disable`
- Миграции автоматически применяются при `docker-compose up`

## Architecture

### Clean Architecture Layers

1. **Domain Layer** (`internal/domain/`)
   - `user.go` - доменная модель User с db тегами для sqlx
   - `car.go` - доменная модель Car с db тегами для sqlx
   - Поля: TGLink (не TGLing), все nullable поля используют указатели

2. **Service Layer** (`internal/service/user/`)
   - `user_service.go` - полная бизнес-логика для User и Car
     - User: CreateUser, UpdateUser, DeleteUser, GetUserByID, GetUserWithCars
     - Car: CreateCar, UpdateCar (PATCH), DeleteCar
     - Обработка ошибок с wrapping через fmt.Errorf
   - `contract.go` - интерфейсы репозиториев и кастомные ошибки
     - ErrUserNotFound, ErrUserAlreadyExists, ErrCarNotFound, ErrCarAccessDenied
   - `models/models.go` - все DTO модели (User, Car, UserWithCars)

3. **Infrastructure Layer** (`internal/infra/storage/`)
   - `user/repository.go` - UserRepository с обработкой ошибок
   - `car/repository.go` - CarRepository с обработкой ошибок
   - Используют psqlbuilder для построения SQL запросов
   - Все методы имеют собственные кастомные ошибки (ErrCreateUser, ErrGetUser, etc.)

4. **Handlers Layer** (`internal/handlers/`)
   - `api/helpers.go` - утилиты (RespondJSON, RespondError, DecodeJSON и специфичные ошибки)
   - `api/create_user/handler.go` - POST /users
   - `api/get_current_user/handler.go` - GET /users/me (возвращает UserWithCars)
   - `api/update_current_user/handler.go` - PUT /users/me
   - `api/delete_current_user/handler.go` - DELETE /users/me
   - `api/create_car/handler.go` - POST /users/me/cars
   - `api/update_car/handler.go` - PATCH /users/me/cars/{car_id}
   - `api/delete_car/handler.go` - DELETE /users/me/cars/{car_id}
   - `middleware/auth.go` - JWT аутентификация middleware

5. **Configuration** (`internal/config/`)
   - `config.go` - загрузка config.toml с секциями logs, server, database, jwt

### Key Design Patterns

- **Repository Pattern**: Абстракция доступа к данным через интерфейсы
- **DTO Pattern**: Разделение доменных моделей и API контрактов
- **Dependency Injection**: Все зависимости через конструкторы
- **Error Wrapping**: Многоуровневая обработка ошибок (repository → service → handler)
- **Handler per Endpoint**: Каждый endpoint в отдельной папке для тестирования

### Database Schema

- **users table**:
  - tg_user_id BIGINT PRIMARY KEY
  - name, phone_number, tg_link, created_at
  - Index на phone_number

- **cars table**:
  - id UUID PRIMARY KEY (генерируется БД)
  - user_id BIGINT FK → users(tg_user_id) ON DELETE CASCADE
  - brand, model, license_plate (required)
  - color, size (nullable)
  - Indexes на user_id и license_plate

### Important Utilities

- **psqlbuilder** (`pkg/psqlbuilder/psqlbuilder.go`): обертки над squirrel
  - Автоматически использует Dollar placeholder ($1, $2) для PostgreSQL
  - Select(), Insert(), Update(), Delete()

## API Contract

API реализует OpenAPI спецификацию из `schemas/api/schema.yaml`:

### Public Endpoints
- `POST /users` - создание пользователя

### Protected Endpoints (JWT required)
- `GET /users/me` - получение пользователя с автомобилями
- `PUT /users/me` - обновление профиля
- `DELETE /users/me` - удаление профиля
- `POST /users/me/cars` - добавление автомобиля
- `PATCH /users/me/cars/{car_id}` - частичное обновление автомобиля
- `DELETE /users/me/cars/{car_id}` - удаление автомобиля

### Authentication
JWT Bearer token в заголовке Authorization:
- Формат: `Authorization: Bearer <token>`
- Token должен содержать claim `tg_user_id` (int64 или string)
- UserID извлекается из токена и используется для авторизации действий

## Configuration

`config.toml`:
```toml
[logs]
level = "info"

[server]
http_port = 8080

[database]
host = "localhost"
port = 5435  # Не стандартный порт!
user = "postgres"
password = "postgres"
dbname = "smk_userservice"
sslmode = "disable"

[jwt]
secret = "your-secret-key-change-in-production"
```

## Project Structure

```
SMK-UserService/
├── cmd/
│   └── main.go                 # Entry point с routing и DI
├── internal/
│   ├── config/
│   │   └── config.go          # Config loader
│   ├── domain/
│   │   ├── user.go
│   │   └── car.go
│   ├── service/user/
│   │   ├── contract.go        # Interfaces и errors
│   │   ├── user_service.go    # Business logic
│   │   └── models/models.go   # DTOs
│   ├── infra/storage/
│   │   ├── user/repository.go
│   │   └── car/repository.go
│   └── handlers/
│       ├── middleware/auth.go
│       └── api/
│           ├── helpers.go
│           ├── create_user/handler.go
│           ├── get_current_user/handler.go
│           ├── update_current_user/handler.go
│           ├── delete_current_user/handler.go
│           ├── create_car/handler.go
│           ├── update_car/handler.go
│           └── delete_car/handler.go
├── pkg/psqlbuilder/
├── migrations/
│   ├── 001_create_users_table.up.sql
│   ├── 001_create_users_table.down.sql
│   ├── 002_create_cars_table.up.sql
│   └── 002_create_cars_table.down.sql
├── docker-compose.yml
├── config.toml
└── .gitignore
```