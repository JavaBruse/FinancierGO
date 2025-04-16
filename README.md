# Проект FinancierGO

## Запуск приложения

### Зависимости

- Go 1.20+
- PostgreSQL (в Docker или локально)

### Сторонние библиотеки

- github.com/beevik/etree v1.5.0
- github.com/golang-jwt/jwt/v5 v5.2.2
- github.com/gorilla/mux v1.8.1
- github.com/lib/pq v1.10.9
- golang.org/x/crypto v0.37.0
- gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
- gopkg.in/yaml.v3 v3.0.1

Для загрузки всех зависимостей:

```shell
go mod download
```

### Запуск PostgreSQL в Docker

```bash
docker run --name FinancierGO   -e POSTGRES_USER=postgres   -e POSTGRES_PASSWORD=postgres   -e POSTGRES_DB=financier   -p 5432:5432   -d postgres:13
```

### Запуск проекта

```bash
go mod tidy
go run .
```

### Структура проекта:

```shell
FinancierGo/
├── config/
│   └── config.go              # Конфигурация приложения
├── internal/
│   ├── models/                # Структуры БД (Users, Accounts и т.д.)
│   ├── repositories/          # Работа с БД (SQL-запросы)
│   ├── services/              # Бизнес-логика
│   ├── handlers/              # HTTP-обработчики
│   ├── middleware/            # JWT и другие промежуточные обработчики
│   ├── utils/                 # Хэширование, PGP, SOAP, email
├── routes/
│   └── routes.go              # Регистрация маршрутов
├── pkg/
│   ├── scheduler/             # Планировщик задач для платежей
│   └── migrations/migrate.go  # Миграция БД
├── go.mod
├── go.sum
└── main.go                    # Точка входа
```

# Эндпоинты, варинты запросов и ответов.

## Аутентификация

### `POST /register`

Регистрация пользователя.

**Body:**

```json
{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "mysecret"
}
```

**Response:**

```json
{
  "id": 1,
  "username": "johndoe",
  "email": "john@example.com",
  "created_at": "2025-04-15T12:00:00Z"
}
```

---

### `POST /login`

Аутентификация и получение JWT-токена.

**Body:**

```json
{
  "email": "john@example.com",
  "password": "mysecret"
}
```

**Response:**

```json
{
  "token": "JWT-TOKEN"
}
```

---

## Счета

### `POST /api/accounts`

Создание банковского счета.

**Headers:**
`Authorization: Bearer JWT-TOKEN`

**Body:**

```json
{
  "currency": "RUB"
}
```

**Response:**

```json
{
  "id": 1,
  "user_id": 1,
  "number": "40817810000012345678",
  "balance": 0,
  "currency": "RUB"
}
```

---

### `POST /api/accounts/{id}/deposit`

Пополнение баланса банковского счета.

**Headers:**

`Authorization: Bearer JWT-TOKEN`

**Body:**

```json
{
  "amount": 1000.0
}
```

**Response:**

```json
{
  "status": "success"
}
```

---

### `POST /accounts/{id}/withdraw`

Списание средств со счета.

**Headers:**

`Authorization: Bearer JWT-TOKEN`

**Body:**

```json
{
  "amount": 1000.0
}
```

**Response:**

```json
{
  "status": "success"
}
```

---

### `POST /api/transfer`

Перевод между счетами.

**Body:**

```json
{
  "from_account_id": 1,
  "to_account_id": 2,
  "amount": 1500.75
}
```

**Response:**

```json
{
  "status": "ok"
}
```

---

## Карты

### `POST /api/cards`

Выпуск виртуальной карты.

**Body:**

```json
{
  "account_id": 1,
  "cvv": "123"
}
```

**Response:**

```json
{
  "card_id": 5
}
```

---

## Аналитика

### `GET /api/analytics`

Доходы и расходы за текущий месяц.

**Response:**

```json
{
  "income": 25000.0,
  "expense": 17200.5
}
```

---

### `GET /api/analytics/credit`

Кредитная нагрузка пользователя.

**Response:**

```json
{
  "debt": 82000.0
}
```

---

### `GET /api/accounts/{accountId}/predict?days=30`

Прогноз расходов по кредитам за N дней.

**Response:**

```json
{
  "planned_expense": 5000.0
}
```

---

## Кредиты

### `POST /api/credits`

Оформление кредита.

**Body:**

```json
{
  "account_id": 1,
  "amount": 50000.0,
  "rate": 10.0,
  "months": 12
}
```

**Response:**

```json
{
  "id": 3,
  "amount": 50000.0,
  "remaining": 50000.0,
  "rate": 10.0
}
```

---

### `GET /api/credits/{creditId}/schedule`

Получить график платежей по кредиту.

**Response:**

```json
[
  {
    "amount": 1500.0,
    "due_date": "2025-05-15T00:00:00Z",
    "paid": false
  }
]
```

---

## Интеграция с ЦБ

### `GET /api/cbr/key-rate`

Получение ключевой ставки ЦБ РФ (+5% маржи).

**Response:**

```json
{
  "key_rate": 16.0
}
```

---

## Email-уведомления

Email отправляется автоматически:

- при просрочке платежа,
- если на счете недостаточно средств,
- через SMTP (`gomail`).

---

## Защищенные эндпоинты

Все `/api/*` — требуют JWT в заголовке:

```
Authorization: Bearer JWT-TOKEN
```

### Соответсвие критериям оценки:

- Реализация слоя моделей (8 баллов):
  - Определение структур данных (соответствие таблицам БД) - есть в internal/models/
  - Сериализация/десериализация (теги JSON) - есть в моделях
  - Базовая валидация полей - есть в internal/services/user.go
  - Проверка уникальности - реализована в репозиториях
  - Полная валидация всех полей - реализована в сервисах
- Реализация слоя репозиториев (9 баллов):
  - Инкапсуляция SQL-запросов - реализовано в internal/repositories/
  - Параметризованные запросы - используются во всех запросах
  - Простейшая обработка ошибок БД - есть в репозиториях
  - Управление транзакциями - реализовано в internal/repositories/transaction.go
  - Обработка сложных ошибок БД - есть в репозиториях
- Реализация слоя сервисов (20 баллов):
  - Регистрация и аутентификация - реализовано в internal/services/user.go
  - Создание счетов, пополнение баланса - реализовано в internal/services/account.go
  - Переводы между счетами - реализовано в internal/services/transaction.go
  - Генерация карт (алгоритм Луна) - реализовано в internal/services/card.go
  - Кредиты: расчет аннуитетных платежей - реализовано в internal/services/credit.go
  - Интеграция с SMTP - реализовано в internal/services/notification.go
  - Интеграция с ЦБ РФ - реализовано в internal/services/currency.go
  - Шедулер для списания платежей - реализовано в internal/services/scheduler.go
  - Логирование через logrus - реализовано в internal/middleware/logging.go
- Реализация слоя обработчиков (12 баллов):
  - Валидация входных данных - реализовано в internal/handlers/
  - Формирование HTTP-ответов (JSON) - реализовано в internal/handlers/
  - Вызов методов сервисов - реализовано в обработчиках
  - Реализация всех эндпоинтов из ТЗ - реализовано в routes/routes.go
  - Проверка прав доступа к ресурсам - реализовано в internal/middleware/auth.go
- Реализация маршрутизации (5 баллов):
  - Публичные эндпоинты - реализовано в routes/routes.go
  - Защищенные эндпоинты - реализовано в routes/routes.go
- Реализация Middleware (6 баллов):
  - Проверка JWT-токенов - реализовано в internal/middleware/auth.go
  - Блокировка неавторизованных запросов - реализовано в internal/middleware/auth.go
  - Добавление ID пользователя в контекст - реализовано в internal/middleware/auth.go
- Безопасность (7 баллов):
  - Хеширование паролей (bcrypt) - реализовано в internal/services/user.go
  - Шифрование данных карт - реализовано в internal/services/card.go
  - Хеширование CVV - реализовано в internal/services/card.go
  - Проверка прав доступа к счетам - реализовано в internal/services/account.go
- База данных (2 балла):
  - Создание минимальных таблиц - реализовано в migrations/
