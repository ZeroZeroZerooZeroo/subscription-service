# Subscription Service

[![Go Version](https://img.shields.io/badge/Go-1.25.1%2B-blue.svg)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Database-informational.svg)](https://www.postgresql.org)
[![Docker](https://img.shields.io/badge/Docker-Containerization-2496ED.svg)](https://www.docker.com)
[![Swagger](https://img.shields.io/badge/Swagger-API%20Documentation-85EA2D.svg)](https://swagger.io)

Микросервис для управления подписками пользователей с подсчетом суммарной стоимости всех подписок за выбранный период с фильтрацией по id пользователя и названию подписки.

---

## 📖 Оглавление

- [Subscription Service](#subscription-service)
  - [📖 Оглавление](#-оглавление)
  - [🚀 О проекте](#-о-проекте)
  - [✨ Функциональность](#-функциональность)
  - [🛠 Технологии](#-технологии)
  - [🚀 Быстрый старт](#-быстрый-старт)
    - [Предварительные требования](#предварительные-требования)
    - [Установка и запуск](#установка-и-запуск)
  - [📚 API Документация](#-api-документация)
    - [Основные эндпоинты:](#основные-эндпоинты)
    - [Примеры запросов:](#примеры-запросов)

---

## 🚀 О проекте

Subscription Service — это микросервис на Go для управления подписками пользователей. Сервис предоставляет REST API для создания, обновления, удаления подписок и расчета их стоимости за указанный период.

## ✨ Функциональность

*   **Управление подписками:**
    *   Создание новых подписок с автоматическим расчетом даты окончания
    *   Просмотр информации о подписке по ID
    *   Обновление данных существующих подписок
    *   Удаление подписок
    *   Пагинированный список всех подписок

*   **Расчет стоимости:**
    *   Расчет общей стоимости подписок за указанный период с фильтрацей по пользователю и сервису

*   **Документация:**
    *   Полная Swagger документация API


## 🛠 Технологии

*   **Язык программирования:** [Go (Golang)](https://golang.org/)
*   **База данных:** [PostgreSQL](https://www.postgresql.org/) 
*   **Контейнеризация:** [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)
*   **Миграции БД:** [golang-migrate](https://github.com/golang-migrate/migrate)
*   **Документация API:** [Swagger](https://swagger.io/) с [swaggo](https://github.com/swaggo/swag)


## 🚀 Быстрый старт

### Предварительные требования

- Docker & Docker Compose
- Go 1.25+ 

### Установка и запуск

1. **Клонирование репозитория:**
```bash
git clone https://github.com/ZeroZeroZerooZeroo/subscription-service.git
cd subscription-service
```

2. Настройка переменных окружения:
Создайте файл .env в корне проекта со следующим содержимым:

```bash
SERVER_HOST=localhost
SERVER_PORT=8080

DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres 
DB_PASSWORD=password 
DB_NAME=subscription_service
DB_SSLMODE=disable
```

3. **Сборка и запуск приложения:**
```bash
docker-compose up -d --build
```

Сервис будет доступен по адресу: `http://localhost:8080`

## 📚 API Документация

После запуска сервиса документация API доступна по адресу:
- **Swagger UI:** http://localhost:8080/swagger/

### Основные эндпоинты:

| Метод | Путь | Описание | Параметры |
|-------|------|-----------|-----------|
| POST | `/subscriptions` | Создать подписку | - |
| GET | `/subscriptions?id={id}` | Получить подписку по ID | `id` (query) |
| PUT | `/subscriptions?id={id}` | Обновить подписку | `id` (query) |
| DELETE | `/subscriptions?id={id}` | Удалить подписку | `id` (query) |
| GET | `/subscriptions/list` | Список подписок | `limit`, `offset` (query) |
| POST | `/subscriptions/total-cost` | Расчет стоимости | - |

### Примеры запросов:

```bash
# Создание подписки
curl -X POST http://localhost:8080/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 1500,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "01-2025"
  }'

# Получить подписку по ID
curl "http://localhost:8080/subscriptions?id=1"

# Список подписок с пагинацией
curl "http://localhost:8080/subscriptions/list?limit=10&offset=0"

# Расчет стоимости
curl -X POST http://localhost:8080/subscriptions/total-cost \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "service_name": "Yandex Plus",
    "start_period": "01-2025",
    "end_period": "02-2025"
  }'

# Удалить подписку
curl -X DELETE "http://localhost:8080/subscriptions?id=1"

```
