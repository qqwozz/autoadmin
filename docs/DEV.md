# AutoAdmin

Telegram-платформа управления расписанием для мастеров (репетиторы, фотографы, специалисты).

[Идея проекта](docs/IDEA.md)

---

## Быстрый старт

### Запуск в Docker

```bash
make env-up
```

Сервер доступен на `http://localhost:8080`.

### Запуск локально

```bash
# Требуется: Go 1.22+, GCC (для CGO/sqlite3)
make run
```

---

## Структура проекта

```
autoadmin/
├── app/
│   ├── main.go          # Точка входа, HTTP сервер, роутинг
│   ├── models.go        # Структуры данных (Go structs)
│   ├── handlers.go      # HTTP обработчики (CRUD)
│   ├── data.sql         # Схема SQLite + начальные данные
│   ├── go.mod
│   └── go.sum
├── docs/
│   └── IDEA.md          # Описание идеи проекта
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── README.md
```

---

## Технологии

| Компонент | Технология |
|-----------|------------|
| Бэкенд | Go 1.22 |
| БД | SQLite (WAL mode) |
| Контейнеры | Docker, Docker Compose |
| Роутинг | `net/http` (Go 1.22 ServeMux) |

---

## API Reference

Base URL: `http://localhost:8080`

### Masters

| Метод | Путь | Описание |
|-------|------|----------|
| `GET` | `/api/masters` | Список всех мастеров |
| `GET` | `/api/masters/{id}` | Получить мастера по ID |
| `POST` | `/api/masters` | Создать мастера |
| `PUT` | `/api/masters/{id}` | Обновить мастера |
| `DELETE` | `/api/masters/{id}` | Удалить мастера |

### Clients

| Метод | Путь | Описание |
|-------|------|----------|
| `GET` | `/api/clients` | Список всех клиентов |
| `GET` | `/api/clients/{id}` | Получить клиента по ID |
| `POST` | `/api/clients` | Создать клиента |
| `PUT` | `/api/clients/{id}` | Обновить клиента |
| `DELETE` | `/api/clients/{id}` | Удалить клиента |

### Services

| Метод | Путь | Описание |
|-------|------|----------|
| `GET` | `/api/services` | Список услуг |
| `GET` | `/api/services?master_id={id}` | Услуги конкретного мастера |
| `GET` | `/api/services/{id}` | Получить услугу |
| `POST` | `/api/services` | Создать услугу |
| `PUT` | `/api/services/{id}` | Обновить услугу |
| `DELETE` | `/api/services/{id}` | Удалить услугу |

### Schedule Slots

| Метод | Путь | Описание |
|-------|------|----------|
| `GET` | `/api/schedule-slots` | Все слоты |
| `GET` | `/api/schedule-slots?master_id={id}&status={status}` | С фильтрами |
| `GET` | `/api/schedule-slots/{id}` | Получить слот |
| `POST` | `/api/schedule-slots` | Создать слот |
| `PUT` | `/api/schedule-slots/{id}` | Обновить слот |
| `DELETE` | `/api/schedule-slots/{id}` | Удалить слот |

### Working Hours

| Метод | Путь | Описание |
|-------|------|----------|
| `GET` | `/api/working-hours?master_id={id}` | Рабочие часы мастера |
| `GET` | `/api/working-hours/{id}` | Получить запись |
| `POST` | `/api/working-hours` | Создать |
| `PUT` | `/api/working-hours/{id}` | Обновить |
| `DELETE` | `/api/working-hours/{id}` | Удалить |

### Tariffs

| Метод | Путь | Описание |
|-------|------|----------|
| `GET` | `/api/tariffs` | Список тарифов |
| `GET` | `/api/tariffs/{id}` | Получить тариф |
| `POST` | `/api/tariffs` | Создать тариф |
| `PUT` | `/api/tariffs/{id}` | Обновить тариф |
| `DELETE` | `/api/tariffs/{id}` | Удалить тариф |

### No-Show Settings

| Метод | Путь | Описание |
|-------|------|----------|
| `GET` | `/api/no-show-settings` | Все настройки |
| `GET` | `/api/no-show-settings/{masterId}` | Настройки мастера |
| `PUT` | `/api/no-show-settings/{masterId}` | Обновить настройки |

### Blacklist

| Метод | Путь | Описание |
|-------|------|----------|
| `GET` | `/api/blacklist?master_id={id}` | Чёрный список мастера |
| `POST` | `/api/blacklist` | Добавить в чёрный список |
| `DELETE` | `/api/blacklist/{id}` | Удалить из чёрного списка |

### Blocked Slots

| Метод | Путь | Описание |
|-------|------|----------|
| `GET` | `/api/blocked-slots?master_id={id}` | Заблокированные слоты |
| `POST` | `/api/blocked-slots` | Заблокировать слот |
| `DELETE` | `/api/blocked-slots/{id}` | Разблокировать |

### Ref Codes

| Метод | Путь | Описание |
|-------|------|----------|
| `GET` | `/api/ref-codes?master_id={id}` | Реферальные коды |
| `POST` | `/api/ref-codes` | Создать код |
| `DELETE` | `/api/ref-codes/{id}` | Удалить код |

### Subscription Payments

| Метод | Путь | Описание |
|-------|------|----------|
| `GET` | `/api/subscription-payments?master_id={id}` | Платежи мастера |
| `POST` | `/api/subscription-payments` | Создать платёж |

### Notifications Log

| Метод | Путь | Описание |
|-------|------|----------|
| `GET` | `/api/notifications-log?user_type={type}&user_id={id}` | Лог уведомлений (лимит 100) |

### Dashboard

| Метод | Путь | Описание |
|-------|------|----------|
| `GET` | `/api/master/{id}/dashboard` | Дашборд мастера (мастер + услуги + часы + статистика) |

---

## Примеры запросов

### Создать мастера

```bash
curl -X POST http://localhost:8080/api/masters \
  -H "Content-Type: application/json" \
  -d '{
    "telegram_id": 123456789,
    "name": "Иван Петров",
    "phone": "+79001234567"
  }'
```

### Создать услугу

```bash
curl -X POST http://localhost:8080/api/services \
  -H "Content-Type: application/json" \
  -d '{
    "master_id": 1,
    "name": "Консультация",
    "duration_minutes": 60,
    "price": 1500.00
  }'
```

### Записать клиента

```bash
curl -X POST http://localhost:8080/api/schedule-slots \
  -H "Content-Type: application/json" \
  -d '{
    "master_id": 1,
    "client_id": 1,
    "service_id": 1,
    "start_time": "2026-06-15 14:00:00",
    "end_time": "2026-06-15 15:00:00"
  }'
```

### Получить дашборд мастера

```bash
curl http://localhost:8080/api/master/1/dashboard
```

---

## Статусы записей

| Статус | Описание |
|--------|----------|
| `pending_confirmation` | Ожидает подтверждения по коду |
| `confirmed` | Подтверждена |
| `completed` | Клиент пришёл |
| `no_show` | Неявка |
| `cancelled` | Отменена клиентом |
| `cancelled_by_master` | Отменена мастером |

---

## Makefile команды

| Команда | Описание |
|---------|----------|
| `make env-up` | Запуск в Docker |
| `make env-down` | Остановка |
| `make env-rebuild` | Пересборка |
| `make env-logs` | Логи |
| `make env-shell` | Shell в контейнере |
| `make db-shell` | SQLite shell |
| `make db-schema` | Показать схему БД |
| `make db-tables` | Список таблиц |
| `make db-counts` | Количество записей |
| `make db-reset` | Сброс и пересоздание БД |
| `make run` | Запуск локально |

---

## Переменные окружения

| Переменная | По умолчанию | Описание |
|------------|--------------|----------|
| `DB_PATH` | `/data/database.sqlite` | Путь к файлу SQLite |
| `PORT` | `8080` | Порт сервера |
