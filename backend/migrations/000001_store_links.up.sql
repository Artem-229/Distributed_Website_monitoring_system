CREATE TABLE monitors (
    id UUID PRIMARY KEY,
    users_id UUID NOT NULL,
    urll TEXT NOT NULL,
    time_interval INTEGER NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW()
);