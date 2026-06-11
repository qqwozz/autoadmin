<p align="center">
  <h1 align="center">📋 AutoAdmin</h1>
  <p align="center">
    Backend API + Telegram Bot для управления записями, клиентами и расписанием
  </p>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.22-00ADD8?style=flat-square&logo=go" alt="Go">
  <img src="https://img.shields.io/badge/SQLite-3-003B57?style=flat-square&logo=sqlite" alt="SQLite">
  <img src="https://img.shields.io/badge/Docker-24.0-2496ED?style=flat-square&logo=docker" alt="Docker">
  <img src="https://img.shields.io/badge/Telegram-Bot-26A5E4?style=flat-square&logo=telegram" alt="Telegram">
</p>

---

## 🏗 Архитектура

```
app/
├── cmd/
│   ├── api/main.go              # REST API сервер
│   └── bot/main.go              # Telegram бот
├── internal/
│   ├── model/                   # Общие модели данных
│   ├── db/                      # Подключение к SQLite
│   ├── repository/              # Слой доступа к БД
│   │   ├── master.go            # CRUD мастеров + поиск по telegram_id
│   │   ├── client.go            # CRUD клиентов + логика no-show
│   │   ├── service.go           # CRUD услуг
│   │   ├── schedule.go          # Слоты расписания + занятые слоты
│   │   ├── working_hours.go     # Рабочие часы по дням
│   │   ├── tariff.go            # CRUD тарифов
│   │   └── misc.go              # Блокировка, рефералы, оплаты, уведомления
│   ├── service/                 # Бизнес-логика
│   │   └── booking.go           # Свободные слоты, подтверждение/отмена, no-show
│   ├── handler/                 # HTTP обработчики (REST API)
│   ├── middleware/              # Auth, CORS, логирование
│   ├── auth/                    # JWT генерация/валидация
│   ├── client/                  # HTTP клиент для вызова API
│   └── bot/                     # Telegram бот с командами
└── data.sql                     # Схема БД + тестовые данные
```

---

## 🚀 Быстрый старт

### Docker (рекомендуется)

```bash
# Запуск API + Bot
BOT_TOKEN=ваш_токен make env-up

# Просмотр логов
make env-logs

# Остановка
make env-down
```

### Локально

```bash
# API сервер
cd app && go run ./cmd/api

# Telegram бот (в отдельном терминале)
BOT_TOKEN=ваш_токен API_URL=http://localhost:8080 go run ./cmd/bot
```

---

## 🤖 Telegram Bot

Бот обращается к API и отвечает на команды в Telegram.

| Команда | Описание | Пример |
|---------|----------|--------|
| `/start` | Приветствие | `/start` |
| `/masters` | Список мастеров | `/masters` |
| `/master` | Детали мастера | `/master 1` |
| `/clients` | Список клиентов | `/clients` |
| `/client` | Детали клиента | `/client 1` |
| `/services` | Список услуг | `/services 1` |
| `/schedule` | Расписание записей | `/schedule 1 confirmed` |
| `/tariffs` | Тарифы | `/tariffs` |
| `/available` | Свободные слоты | `/available 1 1 2026-06-15` |
| `/dashboard` | Панель управления | `/dashboard 1` |

---

## 📡 API Endpoints

### Авторизация

| Метод | Путь | Описание |
|-------|------|----------|
| `POST` | `/api/auth/login` | Вход по `telegram_id`, возвращает JWT |

### Мастера

| Метод | Путь | Описание |
|-------|------|----------|
| `GET` | `/api/masters` | Список (фильтр: `?telegram_id=123`) |
| `GET` | `/api/masters/{id}` | Мастер по ID |
| `GET` | `/api/me/master` | Мастер по Telegram ID |
| `POST` | `/api/masters` | Создать |
| `PUT` | `/api/masters/{id}` | Обновить |
| `DELETE` | `/api/masters/{id}` | Удалить |

### Клиенты

| Метод | Путь | Описание |
|-------|------|----------|
| `GET` | `/api/clients` | Список клиентов |
| `GET` | `/api/clients/{id}` | Клиент по ID |
| `GET` | `/api/clients/by-telegram/{telegramId}` | Клиент по Telegram ID |
| `POST` | `/api/clients` | Создать |
| `PUT` | `/api/clients/{id}` | Обновить |
| `DELETE` | `/api/clients/{id}` | Удалить |

### Услуги

| Метод | Путь | Описание |
|-------|------|----------|
| `GET` | `/api/services` | Список (фильтр: `?master_id=1`) |
| `GET` | `/api/services/{id}` | Услуга по ID |
| `POST` | `/api/services` | Создать |
| `PUT` | `/api/services/{id}` | Обновить |
| `DELETE` | `/api/services/{id}` | Удалить |

### Записи

| Метод | Путь | Описание |
|-------|------|----------|
| `GET` | `/api/available-slots` | Свободные слоты (`?master_id=1&service_id=1&date=2026-06-15`) |
| `GET` | `/api/schedule-slots` | Список записей |
| `GET` | `/api/schedule-slots/{id}` | Запись по ID |
| `POST` | `/api/schedule-slots` | Создать запись |
| `PUT` | `/api/schedule-slots/{id}` | Обновить |
| `DELETE` | `/api/schedule-slots/{id}` | Удалить |
| `POST` | `/api/schedule-slots/{id}/confirm` | Подтвердить (с кодом) |
| `POST` | `/api/schedule-slots/{id}/cancel` | Отменить (`?by=user`) |
| `POST` | `/api/schedule-slots/{id}/no-show` | Пометить no-show |

### Рабочие часы

| Метод | Путь | Описание |
|-------|------|----------|
| `GET` | `/api/working-hours` | Список (фильтр: `?master_id=1`) |
| `GET` | `/api/working-hours/{id}` | По ID |
| `POST` | `/api/working-hours` | Создать |
| `PUT` | `/api/working-hours/{id}` | Обновить |
| `DELETE` | `/api/working-hours/{id}` | Удалить |

### Тарифы

| Метод | Путь | Описание |
|-------|------|----------|
| `GET` | `/api/tariffs` | Список тарифов |
| `GET` | `/api/tariffs/{id}` | По ID |
| `POST` | `/api/tariffs` | Создать |
| `PUT` | `/api/tariffs/{id}` | Обновить |
| `DELETE` | `/api/tariffs/{id}` | Удалить |

### Прочее

| Метод | Путь | Описание |
|-------|------|----------|
| `GET` | `/api/no-show-settings` | Настройки no-show |
| `GET` | `/api/blacklist` | Чёрный список |
| `GET` | `/api/blocked-slots` | Заблокированные слоты |
| `GET` | `/api/ref-codes` | Реферальные коды |
| `GET` | `/api/subscription-payments` | Платежи |
| `GET` | `/api/notifications-log` | Лог уведомлений |
| `GET` | `/api/master/{id}/dashboard` | Дашборд мастера |

---

## 🔐 Авторизация

**JWT (для фронтенда):**

```bash
# Получить токен
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"telegram_id": 12345}'

# Использовать токен
curl -H "Authorization: Bearer eyJ..." http://localhost:8080/api/masters
```

**Telegram ID (для бота):**

```bash
curl -H "X-Telegram-ID: 12345" http://localhost:8080/api/masters
```

---

## 🧪 Тесты

```bash
# Все тесты
make test

# Или напрямую
cd app && go test ./internal/... -v
```

---

## ⚙️ Переменные окружения

### API сервер

| Переменная | По умолчанию | Описание |
|------------|-------------|----------|
| `DB_PATH` | `/data/database.sqlite` | Путь к БД SQLite |
| `PORT` | `8080` | Порт сервера |
| `JWT_SECRET` | `autoadmin-dev-secret...` | Секрет для JWT |

### Telegram бот

| Переменная | Описание |
|------------|----------|
| `BOT_TOKEN` | Токен Telegram Bot API (обязательно) |
| `API_URL` | URL API сервера (по умолчанию `http://localhost:8080`) |

---

## 📦 Makefile

```bash
make env-up       # Запуск Docker (API + Bot)
make env-down     # Остановка
make env-rebuild  # Пересборка
make env-logs     # Логи
make test         # Тесты
make build        # Сборка бинарников
make run-api      # Запуск API локально
make run-bot      # Запуск бота локально
```

---

## 📁 Структура БД

```
masters              — Мастера
clients              — Клиенты
services             — Услуги
schedule_slots       — Записи
working_hours        — Рабочие часы
tariffs              — Тарифы
subscription_payments — Платежи
no_show_settings     — Настройки no-show
master_blacklist     — Чёрный список
blocked_slots        — Заблокированные слоты
master_ref_codes     — Реферальные коды
master_client_bindings — Привязки мастер-клиент
notifications_log    — Лог уведомлений
```
