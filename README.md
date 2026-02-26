# Distributed Website Monitoring System

Приложение для мониторинга доступности сайтов. Написано на Go + React(пока что).

## Stack

**Backend:** Go, Gin, PostgreSQL, golang-migrate  
**Frontend:** React, Vite

## Запуск проекта:

### Клонировать репозиторий

```bash
git clone 
cd 
```

### Frontend

```bash
cd frontend
npm install
npm run dev
```

Запускается на `http://localhost:5173`

### Backend

```bash
cd backend
docker compose up --build
```

Запускается на `http://localhost:8080`

> Перед запуском backend убедитесь что у вас запущен докер
