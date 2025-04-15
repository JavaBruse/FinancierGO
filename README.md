```shell
docker run --name FinancierGO -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=financier -p 5432:5432 -d postgres:13

```


### Конфигурация

```shell
SMTP_HOST=smtp.mail.ru
SMTP_PORT=465
SMTP_USER=bank@example.com
SMTP_PASS=supersecret
```

```shell
go test ./internal/handlers -v
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
│   └── scheduler/             # Планировщик задач для платежей
├── go.mod
├── go.sum
└── main.go                    # Точка входа
```

## 🔐 Аутентификация

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

## 💰 Счета

### `POST /accounts`
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

### `POST /transfer`
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

## 💳 Карты

### `POST /cards`
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

## 📊 Аналитика

### `GET /analytics`
Доходы и расходы за текущий месяц.

**Response:**
```json
{
  "income": 25000.0,
  "expense": 17200.5
}
```

---

### `GET /analytics/credit`
Кредитная нагрузка пользователя.

**Response:**
```json
{
  "debt": 82000.0
}
```

---

### `GET /accounts/{accountId}/predict?days=30`
Прогноз расходов по кредитам за N дней.

**Response:**
```json
{
  "planned_expense": 5000.0
}
```

---

## 🧾 Кредиты

### `POST /credits`
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

### `GET /credits/{creditId}/schedule`
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

## 📡 Интеграция с ЦБ

### `GET /cbr/key-rate`
Получение ключевой ставки ЦБ РФ (+5% маржи).

**Response:**
```json
{
  "key_rate": 16.0
}
```

---

## 📬 Email-уведомления

Email отправляется автоматически:
- при просрочке платежа,
- если на счете недостаточно средств,
- через SMTP (`gomail`).

---
## 🛡️ Защищенные эндпоинты

Все `/accounts`, `/transfer`, `/cards`, `/analytics`, `/credits`, `/predict` — требуют JWT в заголовке:

```
Authorization: Bearer JWT-TOKEN
```

#### Get Credit Schedule
- **URL**: `/api/credits/{creditId}/schedule`
- **Method**: `GET`
- **URL Parameters**: 
  - `creditId`: integer

### Notes:
1. All API endpoints except `/register` and `/login` require authentication
2. Authentication is done via JWT token in the Authorization header
3. The token should be included in the format: `Bearer <token>`
