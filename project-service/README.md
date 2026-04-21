# ProjectService

`project-service` отвечает за организационную структуру проектной деятельности компании.

Это внутренний gRPC-сервис. Публичный доступ к проектам идет только через `gateway`.

## Роль в системе

Сервис задает каркас всей проектной модели:
- какие проекты существуют
- в каком состоянии они находятся
- из каких этапов состоят
- кто участвует в проекте
- какие события происходят внутри проекта

В сервисе нет финансовой логики платежей, но он хранит проектные данные, на которые потом будут опираться другие подсистемы.

## Сущности

### `projects`

Хранит сам проект:
- `id`
- `project_code`
- `name`
- `description`
- `contract_number`
- `status`
- `customer_name`
- `manager_id`
- `planned_start_date`
- `planned_deadline`
- `actual_start_date`
- `actual_deadline`
- `planned_budget`

Поддерживаемые статусы проекта:
- `draft`
- `planned`
- `active`
- `paused`
- `completed`
- `cancelled`

Замечание:
- при создании `contract_number` может быть автоматически сгенерирован, если клиент его не передал

### `project_stages`

Хранит этапы проекта:
- `id`
- `project_id`
- `name`
- `description`
- `sequence_number`
- `status`
- `planned_start_date`
- `planned_end_date`
- `actual_start_date`
- `actual_end_date`
- `planned_income`
- `planned_expense`

Статусы этапов:
- `draft`
- `planned`
- `in_progress`
- `completed`
- `cancelled`

### `project_members`

Связывает проект и пользователей:
- `id`
- `project_id`
- `user_id`
- `role_in_project`
- `is_active`
- `joined_at`
- `left_at`

### `project_events`

Фиксирует ключевые события проекта:
- `id`
- `project_id`
- `stage_id`
- `name`
- `description`
- `planned_date`
- `actual_date`
- `status`

Статусы событий:
- `planned`
- `reached`
- `cancelled`

## Что умеет сервис

Полный CRUD по:
- проектам
- этапам проекта
- участникам проекта
- событиям проекта

## Внутренний интерфейс

Сервис публикует gRPC API:
- `Ping`
- `CreateProject`, `GetProject`, `ListProjects`, `UpdateProject`, `DeleteProject`
- `CreateStage`, `GetStage`, `ListStages`, `UpdateStage`, `DeleteStage`
- `CreateMember`, `GetMember`, `ListMembers`, `UpdateMember`, `DeleteMember`
- `CreateEvent`, `GetEvent`, `ListEvents`, `UpdateEvent`, `DeleteEvent`

Контракт находится в:
- [`contracts/proto/project/project.proto`](../contracts/proto/project/project.proto)

## Где находится контроль доступа

Сам `project-service` реализует модель данных и валидацию сущностей.

Проверка прав доступа выполняется в `gateway`:
- `admin` может управлять любыми проектами
- `manager` управляет своими проектами
- активные участники проекта могут читать доступные им данные

## Конфигурация

Базовый файл:
- `project-service/configs/config.yaml`

Поддерживаемые override:

| YAML | Env | Значение |
| --- | --- | --- |
| `log_level` | `LOG_LEVEL` | уровень логирования |
| `address` | `PROJECT_ADDRESS` | адрес gRPC-сервера |
| `db_address` | `DB_ADDRESS` | строка подключения к PostgreSQL |

## Запуск

Локально:

```bash
go run ./cmd -config configs/config.yaml
```

Через Docker Compose из корня репозитория:

```bash
docker compose up --build -d project-service
```

## Миграции

Миграции находятся в:
- `project-service/internal/adapters/db/postgres/migrations`

Текущая схема создается файлами:
- `000001_create_project.up.sql`
- `000002_create_project_stages.up.sql`
- `000003_create_project_members.up.sql`
- `000004_create_project_events.up.sql`

## Что важно помнить

- у сервиса нет внешнего REST API
- Swagger для него не ведется
- внешняя документация по проектам находится только в `gateway`
