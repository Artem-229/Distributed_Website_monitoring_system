# Distributed Website Monitoring System

Web application for distributed website monitoring. Built with Go + React(for now).

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
