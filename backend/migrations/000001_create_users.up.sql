CREATE TABLE users (
    id UUID PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    login VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE monitors (
    id UUID PRIMARY KEY,
    users_id UUID NOT NULL,
    url TEXT NOT NULL,
    time_interval INTEGER NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE monitor_checks (
    id UUID PRIMARY KEY,
    time_interval INTEGER NOT NULL,
    responce_time FLOAT NOT NULL,
    checked_at TIMESTAMP NOT NULL DEFAULT NOW(),
    status_ok BOOLEAN NOT NULL
);