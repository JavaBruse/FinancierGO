-- Создание таблицы пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Таблица accounts
CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    number VARCHAR(255),
    balance DECIMAL(10,2) DEFAULT 0,
    currency VARCHAR(3),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Таблица transactions
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    account_id INTEGER REFERENCES accounts(id),
    amount DECIMAL(10,2),
    type VARCHAR(50),
    date TIMESTAMP DEFAULT NOW()
);

-- Таблица credits
CREATE TABLE credits (
    id SERIAL PRIMARY KEY,
    account_id INTEGER REFERENCES accounts(id),
    amount DECIMAL(10,2),
    rate DECIMAL(5,2),
    months INTEGER,
    start_date TIMESTAMP,
    next_payment TIMESTAMP,
    remaining DECIMAL(10,2)
);

-- Таблица cards
CREATE TABLE cards (
    id SERIAL PRIMARY KEY,
    account_id INTEGER REFERENCES accounts(id),
    encrypted TEXT,
    cvv_hash TEXT,
    hmac TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Таблица payment_schedules
CREATE TABLE payment_schedules (
    id SERIAL PRIMARY KEY,
    credit_id INTEGER REFERENCES credits(id),
    amount DECIMAL(10,2),
    due_date TIMESTAMP,
    paid BOOLEAN DEFAULT false
);