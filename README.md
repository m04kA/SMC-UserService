# SMK-UserService

Сервис управления пользователями и их автомобилями для приложения автомойки.

## 🏗️ Архитектура

Проект построен на **Clean Architecture** с четким разделением слоев:
- **Domain** - доменные модели (User, Car)
- **Service** - бизнес-логика
- **Repository** - работа с БД (PostgreSQL + sqlx)
- **Handlers** - HTTP API
- **Middleware** - JWT аутентификация, Prometheus метрики
- **Logging** - многоуровневое логирование (INFO, WARN, ERROR)

## 🚀 Быстрый старт

### 1. Запуск инфраструктуры (PostgreSQL, Prometheus, Grafana)

```bash
docker-compose up -d
```

Сервисы:
- **PostgreSQL**: порт **5435** (не стандартный 5432)
- **Prometheus**: http://localhost:9091
- **Grafana**: http://localhost:3001 (admin/admin)

### 2. Запуск сервера

```bash
go run cmd/main.go
```

Сервер запустится на `http://localhost:8080`

Логи записываются в консоль и `logs/app.log` (WARN и ERROR)

### 3. Тестирование API

#### Создание пользователя (публичный endpoint)
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "tg_user_id": 123456789,
    "name": "Иван",
    "phone_number": "+79991234567",
    "tg_link": "@ivan"
  }'
```

#### Получение текущего пользователя (требует JWT)
```bash
curl -X GET http://localhost:8080/users/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 📋 API Endpoints

### Public
- `POST /users` - создание пользователя

### Protected (требуют JWT)
- `GET /users/me` - получение пользователя с автомобилями
- `PUT /users/me` - обновление профиля
- `DELETE /users/me` - удаление профиля
- `POST /users/me/cars` - добавление автомобиля
- `PATCH /users/me/cars/{car_id}` - обновление автомобиля (car_id: int64)
- `DELETE /users/me/cars/{car_id}` - удаление автомобиля (car_id: int64)

### Monitoring
- `GET /metrics` - Prometheus метрики в формате OpenMetrics

## 🔧 Разработка

### Установка зависимостей
```bash
go mod tidy
```

### Сборка проекта
```bash
go build -o bin/server ./cmd/main.go
```

### Запуск тестов
```bash
go test ./...
```

### Остановка Docker
```bash
docker-compose down
```

### Полная очистка (с данными)
```bash
docker-compose down -v
```

## 📁 Структура проекта

```
SMK-UserService/
├── cmd/
│   └── main.go                           # Entry point
├── internal/
│   ├── config/                           # Конфигурация
│   ├── domain/                           # Доменные модели (User, Car)
│   ├── service/user/                     # Бизнес-логика + DTOs
│   ├── infra/storage/                    # Репозитории (PostgreSQL)
│   └── handlers/
│       ├── api/                          # HTTP handlers (handler per endpoint)
│       └── middleware/                   # Auth + Metrics middleware
├── pkg/
│   ├── logger/                           # Кастомный логгер
│   └── psqlbuilder/                      # Утилиты для SQL (squirrel wrapper)
├── monitoring/
│   ├── prometheus/prometheus.yml         # Prometheus конфигурация
│   └── grafana/                          # Grafana dashboards + datasources
├── migrations/                           # SQL миграции (golang-migrate)
├── schemas/api/schema.yaml               # OpenAPI 3.1.0 спецификация
├── config.toml                           # Конфигурация приложения
└── docker-compose.yml                    # Docker окружение (PostgreSQL, Prometheus, Grafana)
```

## ⚙️ Конфигурация

Файл `config.toml`:
- `[logs]` - уровень логирования
- `[server]` - порт HTTP сервера (по умолчанию 8080)
- `[database]` - настройки подключения к PostgreSQL (порт 5435)
- `[jwt]` - секретный ключ для JWT токенов

## 🔐 Аутентификация

JWT токен должен содержать claim `tg_user_id` (int64 или string).

Формат заголовка:
```
Authorization: Bearer <your-jwt-token>
```

Для генерации тестового токена используйте утилиту в `pkg/gentoken/main.go`.

## 📊 Мониторинг и логирование

### Prometheus метрики
- `http_requests_total` - счётчик HTTP запросов
- `http_request_duration_seconds` - длительность запросов
- `http_requests_in_flight` - активные запросы

Доступны на http://localhost:8080/metrics

### Grafana Dashboard
Автоматически загружается при старте:
- HTTP Request Rate (req/m)
- Requests In Flight
- Request Duration (p95, p99)
- Status Codes Distribution
- Requests by Endpoint

Доступен на http://localhost:3001 (admin/admin)

### Логирование
Все логи пишутся в консоль. Логи уровня WARN и ERROR дополнительно сохраняются в `logs/app.log`.

Уровни:
- **INFO** - успешные операции
- **WARN** - ошибки валидации и авторизации (4xx)
- **ERROR** - внутренние ошибки сервера (5xx)

## 📚 Документация

Полная документация архитектуры и паттернов проектирования находится в [CLAUDE.md](CLAUDE.md).

API Contract описан в [schemas/api/schema.yaml](schemas/api/schema.yaml) (OpenAPI 3.1.0).
