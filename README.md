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

# FinancierGO

## API Endpoints Documentation

### Authentication

#### Register
- **URL**: `/register`
- **Method**: `POST`
- **Body**:
```json
{
    "username": "string",
    "email": "string",
    "password": "string"
}
```

#### Login
- **URL**: `/login`
- **Method**: `POST`
- **Body**:
```json
{
    "email": "string",
    "password": "string"
}
```

### Accounts (Requires Authentication)

#### Create Account
- **URL**: `/api/accounts`
- **Method**: `POST`
- **Body**:
```json
{
    "currency": "string"
}
```

#### Transfer Money
- **URL**: `/api/transfer`
- **Method**: `POST`
- **Body**:
```json
{
    "from_account_id": "integer",
    "to_account_id": "integer",
    "amount": "float"
}
```

### Cards (Requires Authentication)

#### Create Card
- **URL**: `/api/cards`
- **Method**: `POST`
- **Body**:
```json
{
    "account_id": "integer",
    "cvv": "string"
}
```

### Credits (Requires Authentication)

#### Create Credit
- **URL**: `/api/credits`
- **Method**: `POST`
- **Body**:
```json
{
    "account_id": "integer",
    "amount": "float",
    "rate": "float",
    "months": "integer"
}
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
