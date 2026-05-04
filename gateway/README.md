# Gateway

`gateway` — единственная внешняя точка входа в систему.

Он принимает REST-запросы от клиента, валидирует JWT, проверяет роли и права доступа, а затем оркестрирует вызовы во внутренние gRPC-сервисы:
- `auth-service`
- `user-service`
- `project-service`
- `report-service`
- `export-service`
- `payment-calendar-service`

## Роль в системе

`gateway` отвечает за:
- внешний REST API
- авторизацию по Bearer JWT
- преобразование REST DTO в gRPC DTO и обратно
- orchestration-сценарии между несколькими сервисами
- Swagger / OpenAPI документацию для внешнего API

`gateway` не хранит собственную БД.

## Что важно понимать

- внешние клиенты общаются только с `gateway`
- внутренние сервисы не обязаны публиковать REST
- вся OpenAPI-документация относится только к `gateway`
- если меняется внешний REST-контракт, обновлять нужно handler-аннотации и Swagger в `gateway`

## Основные группы доступных endpoints

### Auth

- `POST /api/auth/register`
- `POST /api/auth/login`
- `GET /api/auth/me`

### Users

- `POST /api/users`
- `GET /api/users`
- `GET /api/users/{id}`
- `PUT /api/users/{id}`
- `DELETE /api/users/{id}`

### Projects

- `POST /api/projects`
- `GET /api/projects`
- `GET /api/projects/{id}`
- `PUT /api/projects/{id}`
- `DELETE /api/projects/{id}`

### Project stages

- `POST /api/projects/{projectId}/stages`
- `GET /api/projects/{projectId}/stages`
- `GET /api/project-stages/{id}`
- `PUT /api/project-stages/{id}`
- `DELETE /api/project-stages/{id}`

### Project members

- `POST /api/projects/{projectId}/members`
- `GET /api/projects/{projectId}/members`
- `GET /api/project-members/{id}`
- `PUT /api/project-members/{id}`
- `DELETE /api/project-members/{id}`

### Project events

- `POST /api/projects/{projectId}/events`
- `GET /api/projects/{projectId}/events`
- `GET /api/project-events/{id}`
- `PUT /api/project-events/{id}`
- `DELETE /api/project-events/{id}`

### Reports

- `POST /api/reports`
- `GET /api/reports`
- `GET /api/reports/{id}`
- `PUT /api/reports/{id}`
- `DELETE /api/reports/{id}`

### Report entries

- `POST /api/reports/{reportId}/entries`
- `GET /api/reports/{reportId}/entries`
- `GET /api/report-entries/{id}`
- `PUT /api/report-entries/{id}`
- `DELETE /api/report-entries/{id}`

### Report comments

- `POST /api/reports/{reportId}/comments`
- `GET /api/reports/{reportId}/comments`
- `GET /api/report-comments/{id}`
- `PUT /api/report-comments/{id}`
- `DELETE /api/report-comments/{id}`

### System

- `GET /api/ping`
- `GET /swagger/index.html`
- `GET /openapi/swagger.json`
- `GET /openapi/swagger.yaml`

## Авторизация и права

`gateway` является `middleware`, а именно проверяет:
- наличие и валидность JWT не самостоятельно, а через `auth-service`
- роль пользователя
- принадлежность ресурса пользователю
- проектные права, например manager или active member

Примеры:
- список всех пользователей доступен только администратору
- обычный пользователь видит только свой профиль
- только `admin` и `manager` могут создавать проекты; обычные сотрудники (`user`) работают только в назначенных им проектах
- менеджер проекта может работать с проектом и отчетами по связанным проектам
- администратор имеет глобальный доступ

## Swagger / OpenAPI

Swagger поддерживается только здесь.

Спецификация генерируется в:
- `gateway/api/swagger.json`
- `gateway/api/swagger.yaml`

Swagger UI доступен по адресу:
- [http://localhost:28080/swagger/index.html](http://localhost:28080/swagger/index.html)

Сырой OpenAPI:
- [http://localhost:28080/openapi/swagger.json](http://localhost:28080/openapi/swagger.json)
- [http://localhost:28080/openapi/swagger.yaml](http://localhost:28080/openapi/swagger.yaml)

## Конфигурация

Базовый файл:
- `gateway/configs/config.yaml`

Переопределение через env поддерживается.

| YAML                       | Env                        | Значение                         |
|----------------------------|----------------------------|----------------------------------|
| `log_level`                | `LOG_LEVEL`                | уровень логирования              |
| `api-server.address`       | `ADDRESS`                  | адрес HTTP-сервера               |
| `api-server.timeout`       | `API_TIMEOUT`              | timeout для API                  |
| `auth_address`             | `AUTH_ADDRESS`             | адрес `auth-service`             |
| `user_address`             | `USER_ADDRESS`             | адрес `user-service`             |
| `project_address`          | `PROJECT_ADDRESS`          | адрес `project-service`          |
| `report_address`           | `REPORT_ADDRESS`           | адрес `report-service`           |
| `export_address`           | `EXPORT_ADDRESS`           | адрес `export-service`           |
| `payment_Calendar_address` | `PAYMENT_CALENDAR_ADDRESS` | адрес `payment-calendar-service` |
  

## Запуск

Локально:

```bash
go run ./cmd -config configs/config.yaml
```

Через Docker Compose из корня репозитория:

```bash
make swagger
docker compose up --build -d gateway
```

`make swagger` важен, потому что контейнер `gateway` раздает сгенерированные файлы из `gateway/api`.

## Зависимости

Для полноценной работы `gateway` нужны:
- `auth-service`
- `user-service`
- `project-service`
- `report-service`
- `export-service`
- `patment-calendar-service`

Также именно `gateway` зависит от актуальных `.proto`-контрактов из `contracts`.
