# Distributed Website Monitoring System
 
Система мониторинга доступности и скорости отклика веб-сайтов. Серверная часть написана на Go с Apache Kafka и PostgreSQL. Клиент - приложение для Android TV на Kotlin + Jetpack Compose с поддержкой навигации через пульт (D-pad) и отображением данных по 5 регионам.
 
## Stack
 
**Backend:** Go, PostgreSQL, Apache Kafka (KRaft), golang-migrate, Docker  
**Android TV client:** Kotlin, Jetpack Compose, Retrofit, OkHttp, Gson, StateFlow
 
## Структура репозитория
 
```
master
├── backend/     - Go-сервер, Kafka, PostgreSQL, Docker Compose
└── frontend/    - React (устаревшая версия, не используется)
 
branch: frontend - Android TV приложение на Kotlin
```
 
---
 
## Запуск backend (ветка master)
 
Требования: [Docker](https://docs.docker.com/get-docker/) и Docker Compose.
 
```bash
git clone
cd backend
docker compose up --build
```
 
Сервер поднимется на `http://localhost:8080`
 
Docker Compose автоматически запустит:
- PostgreSQL и применит миграции через golang-migrate
- Apache Kafka в режиме KRaft (без ZooKeeper)
- Go-приложение
> Первый запуск займёт чуть больше времени - Docker скачивает образы и собирает бинарник Go.
 
## Запуск Android TV клиента (ветка frontend)
  
```bash
git clone -b frontend
```

Требования: Android Studio (Ladybug или новее), Android SDK 36.

- Открыть папку проекта в Android Studio: File → Open
- Дождаться завершения Gradle sync
- В app/build.gradle указать IP-адрес машины с запущенным backend:  
groovy   buildConfigField("String", "BASE_URL", "\"http://<IP>:8080/\"")
- Подключить Android TV устройство или запустить AVD с образом Android TV: Device Manager → Create Device → TV
- Нажать Run

---
 
## API
 
| Метод | Путь | Описание | Auth |
|-------|------|----------|------|
| POST | `/registration` | Регистрация | — |
| POST | `/login` | Вход, возвращает JWT | — |
| GET | `/api/monitors` | Список мониторов | JWT |
| POST | `/api/addmonitor` | Добавить монитор | JWT |
| POST | `/api/deletemonitor` | Удалить монитор | JWT |
| GET | `/api/checks/{id}` | История проверок | JWT |
| GET | `/api/checks/{id}/regions` | Данные по 5 регионам | JWT |
 
JWT передаётся в заголовке: `Authorization: Bearer <token>`
 
---
 
## Как это работает
 
```
Android TV client
      │  HTTP + JWT
      ▼
  Go REST API
      │
      ├── PostgreSQL - хранит пользователей, мониторы, результаты, алерты
      │
      └── Worker (горутины)
              │  публикует событие после каждой проверки
              ▼
           Kafka topic: monitor.results
              │
              ▼
           Consumer - создаёт алерты при DOWN или latency > 700ms
```
 
Воркер опрашивает активные мониторы с заданным интервалом (30 сек / 1 мин / 5 мин / 15 мин / 1 час), измеряет время отклика и сохраняет результат. Клиент Android TV отображает историю проверок и сводку по 5 географическим регионам: Россия, США, Китай, Центральная Европа, Азиатско-Тихоокеанский регион.
