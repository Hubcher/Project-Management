# ReportService

`report-service` отвечает за оперативную фиксацию фактически выполненной работы.

Это внутренний gRPC-сервис. Публичный доступ к отчетам идет только через `gateway`.

## Роль в системе

Сервис отвечает на вопрос:

Что именно делали сотрудники по проектам и этапам в конкретный день?

Это не финансовый сервис, но он критически важен для:
- прозрачности работ
- анализа трудозатрат
- контроля загрузки сотрудников
- управленческого цикла согласования отчетов

## Сущности

### `daily_reports`

Заголовок дневного отчета:
- `id`
- `user_id`
- `report_date`
- `status`
- `total_hours`
- `summary`
- `created_at`
- `updated_at`

Поддерживаемые статусы:
- `draft`
- `submitted`
- `approved`
- `rejected`

### `daily_report_entries`

Строки внутри отчета:
- `id`
- `report_id`
- `project_id`
- `stage_id`
- `work_type`
- `description`
- `hours_spent`
- `created_at`
- `updated_at`

Строки позволяют разносить время по проектам и этапам.

### `daily_report_comments`

Комментарии к отчету:
- `id`
- `report_id`
- `author_user_id`
- `comment`
- `created_at`

Комментарии нужны для замечаний, обратной связи и процесса ревью.

## Что умеет сервис

Полный CRUD по:
- дневным отчетам
- строкам отчета
- комментариям

Дополнительно сервис:
- валидирует UUID, даты, статусы и часы
- автоматически пересчитывает `total_hours` при изменении строк отчета
- проверяет существование родительского отчета при работе со строками и комментариями

## Внутренний интерфейс

Сервис публикует gRPC API:
- `Ping`
- `CreateReport`, `GetReport`, `ListReports`, `UpdateReport`, `DeleteReport`
- `CreateEntry`, `GetEntry`, `ListEntries`, `UpdateEntry`, `DeleteEntry`
- `CreateComment`, `GetComment`, `ListComments`, `UpdateComment`, `DeleteComment`

Контракт находится в:
- [`contracts/proto/report/report.proto`](../contracts/proto/report/report.proto)

## Где находится контроль доступа

`report-service` отвечает за хранение и валидацию данных.

Права доступа проверяет `gateway`:
- автор отчета работает со своими отчетами
- менеджер проекта может читать и ревьюить отчеты по своим проектам
- администратор может работать со всеми отчетами

Таким образом, внешний access control остается на уровне `gateway`, а не внутри самого сервиса.

## Конфигурация

Базовый файл:
- `report-service/configs/config.yaml`

Поддерживаемые override:

| YAML | Env | Значение |
| --- | --- | --- |
| `log_level` | `LOG_LEVEL` | уровень логирования |
| `address` | `REPORT_ADDRESS` | адрес gRPC-сервера |
| `db_address` | `DB_ADDRESS` | строка подключения к PostgreSQL |

## Запуск

Локально:

```bash
go run ./cmd -config configs/config.yaml
```

Через Docker Compose из корня репозитория:

```bash
docker compose up --build -d report-service
```

## Миграции

У этого сервиса каталог называется:
- `report-service/internal/adapters/db/postgres/migration`

Текущая схема создается файлами:
- `000001_create_daily_reports.up.sql`
- `000002_create_daily_entries.up.sql`
- `000003_create_daily_report_comments.up.sql`


## Что важно помнить

- у сервиса нет внешнего REST API
- Swagger для него не ведется
- внешняя документация по отчетам находится только в `gateway`
