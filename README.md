# BuildIn CRM Backend

**Project-management** - микросервисная CRM-система для управления проектами и отчётами. Стек: Go, PostgreSQL, Docker, Swagger (OpenAPI).

---

## 📦 Архитектура монорепозитория

```
buildin-crm-backend/
├── auth-service/
├── contracts/
├── user-service/
├── project-service/
├── report-service/
├── export-service/
├── gateway/ 
├── docker-compose.yaml
├── Makefile
└── README.md
```

---

## 🧑‍💻 Разработка

### 📌 Требования

- [Go 1.21+](https://go.dev/dl/)
- [Docker Desktop](https://www.docker.com/products/docker-desktop)
- Терминал: PowerShell, CMD или терминал из GoLand/VS Code

---

## 🚀 Быстрый старт

### 1. Необходимые зависимости

```powershell
make tools
```

### 2. Запуск инфраструктуры

Запуск всего кластера контейнеров

```powershell
make up
```

## ⚙️ Переменные окружения (`.env.example`)

Пример для всех сервисов:

```env
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=buildin_crm
```

---

## 📚 Swagger API

Каждый сервис содержит файл `api/swagger.yaml` — для описания API.  
Можно использовать Swagger UI, Postman или [Swagger Editor](https://editor.swagger.io/).

---
