# AuthService

`auth-service` отвечает за аутентификацию, авторизацию и выпуск JWT.

Это внутренний gRPC-сервис. Внешние клиенты не должны обращаться к нему напрямую: публичный доступ идет только через `gateway`.

## Ответственность сервиса

`auth-service` хранит только данные, связанные с учетной записью:
- `email`
- `password_hash`
- `role`
- `is_active`
- timestamps и служебные поля безопасности

Личные данные пользователя не хранятся здесь. Они находятся в `user-service`.

## Что умеет сервис

- регистрировать учетные данные
- логинить пользователя по email и паролю
- валидировать JWT
- возвращать claims пользователя для `gateway`
- удалять учетные данные пользователя по orchestration-сценарию
- автоматически создавать bootstrap admin при старте

## Таблицы БД

### `auth_accounts`

Основные поля:
- `user_id`
- `email`
- `password_hash`
- `role`
- `is_active`
- `created_at`
- `updated_at`
- `password_changed_at`

Роль ограничена значениями:
- `admin`
- `user`

## Внутренний интерфейс

Сервис публикует gRPC API:
- `Ping`
- `Register`
- `Login`
- `Validate`
- `DeleteCredentials`

Контракт находится в:
- [`contracts/proto/auth/auth.proto`](../contracts/proto/auth/auth.proto)

## Как сервис используется в системе

Через `gateway`:
- `POST /api/auth/register`
- `POST /api/auth/login`
- `GET /api/auth/me`

При регистрации пользователя `gateway` оркестрирует сценарий:

1. создает учетные данные в `auth-service`
2. создает профиль в `user-service`
3. при ошибке на втором шаге откатывает создание credentials

## Bootstrap admin

При запуске сервис может автоматически создать администратора.

В текущем `docker-compose.yaml` для локальной разработки используются значения:
- email: `admin@example.com`
- password: `Admin123!`

Это dev-настройка. Для `production` такие значения `необходимо` задавать через секреты или env override.

## Конфигурация

Базовый файл:
- `auth-service/configs/config.yaml`

Поддерживаемые override:

| YAML | Env | Значение |
| --- | --- | --- |
| `log_level` | `LOG_LEVEL` | уровень логирования |
| `address` | `AUTH_ADDRESS` | адрес gRPC-сервера |
| `env` | `ENV` | окружение |
| `db_address` | `DB_ADDRESS` | строка подключения к PostgreSQL |
| `bootstrap_admin.enabled` | `BOOTSTRAP_ADMIN_ENABLED` | включить bootstrap admin |
| `bootstrap_admin.email` | `BOOTSTRAP_ADMIN_EMAIL` | email bootstrap admin |
| `bootstrap_admin.password` | `BOOTSTRAP_ADMIN_PASSWORD` | пароль bootstrap admin |

## Запуск

Локально:

```bash
go run ./cmd -config configs/config.yaml
```

Через Docker Compose из корня репозитория:

```bash
docker compose up --build -d auth-service
```

## Что важно помнить

- у сервиса нет внешнего REST API
- Swagger для него не ведется
- внешний клиент работает только через `gateway`
- источником истины для внутреннего API является `.proto`, а не OpenAPI
