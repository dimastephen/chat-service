# Chat & Auth Service (WIP)

⚠️ **Проект находится в стадии доработки**  
Это pet-project на Go, включающий чат-сервис и сервис аутентификации (JWT).  
Основная цель — отработать архитектуру микросервисов, gRPC, gRPC-Gateway и работу с JWT.

---

## 📌 Архитектура

- **Сервисы**:
  - `chat` — сервис чата (создание, удаление чатов, отправка сообщений).
  - `auth` — сервис аутентификации (выдача access/refresh токенов, проверка прав, работа с ролями).

- **Основные технологии**:
  - Go (1.22+)
  - gRPC + gRPC-Gateway
  - Swagger (через `statik`)
  - PostgreSQL (миграции через `migrations/`)
  - JWT (access + refresh токены)
  - Docker / Docker Compose
  - zap (логирование)
  - interceptors (валидация, логирование)
  - собственная библиотека-утилита [utils](https://github.com/dimastephen/utils) (обёртка над DB и другие полезные функции)

---

## 🗂️ Структура проекта
├── chat/ # чат-сервис
│ ├── cmd/ # main.go
│ ├── internal/ # app, service, repository, api, interceptor и др.
│ ├── migrations/ # SQL миграции
│ └── api/ # protobuf/gRPC API
│
├── auth/ # сервис авторизации
│ ├── cmd/ # main.go
│ ├── internal/ # app, service, repository, jwt, crypto и др.
│ └── api/ # protobuf/gRPC API
│
├── docker-compose.yaml # общий запуск сервисов
└── Makefile # команды для запуска/сборки

## 🚀 Запуск
docker-compose up --build
