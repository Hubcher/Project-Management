# Contracts

`contracts` — общий модуль с внутренними gRPC-контрактами системы.

Именно здесь хранятся:
- `.proto`-описания сервисов
- сгенерированный Go-код для protobuf и gRPC
- команды генерации и линтинга контрактов

## Зачем нужен этот модуль

`contracts` — источник истины для внутренних взаимодействий между сервисами.

Он нужен, чтобы:
- держать межсервисные интерфейсы отдельно от реализации
- переиспользовать один и тот же контракт в `gateway` и backend-сервисах
- централизованно управлять изменениями gRPC API для всех микросервисов
- не дублировать protobuf-схемы по нескольким репозиториям или сервисам

Важно:
- OpenAPI / Swagger здесь не генерируется
- документация по внешнему API находится только в `gateway`
- для внутренних интерфейсов источником истины являются `.proto`-файлы

## Структура

```text
contracts/
  Makefile
  README.md
  proto/
    auth/
    export/
    project/
    report/
    user/
    paymentCalendar/
  gen/
    go/ 
```

## Какие контракты есть сейчас

- `proto/auth/auth.proto` — `AuthService`
- `proto/user/user.proto` — `UserService`
- `proto/project/project.proto` — `ProjectService`
- `proto/report/report.proto` — `ReportService`
- `proto/export/export.proto` — контракт для будущего `export-service`
- `proto/patmentCalendar/paymentCalendar` - контракт для будущего `payment-calendar-service`

## Что генерируется

Команда генерации пишет код в `contracts/gen/go`.

Этот код используется:
- `gateway` для gRPC-клиентов
- внутренними сервисами для gRPC-серверов

## Генерация и проверки

Из каталога `contracts`:

```bash
make tools
make protobuf
make protolint
make golint
make lint
```

Из корня монорепозитория:

```bash
make proto
```

## Требования

Для работы с контрактами нужны:
- `protoc`
- `protoc-gen-go`
- `protoc-gen-go-grpc`
- `protolint`
- `golangci-lint`

Установка большинства инструментов автоматизирована через:

```bash
make tools
```

## Важное замечание по Windows

`contracts/Makefile` использует POSIX-команды вроде `find`, `mkdir -p` и `rm -rf`.

Поэтому на Windows удобнее запускать его через:
- Git Bash
- WSL
- MSYS2

В чистом PowerShell эти команды могут отсутствовать.

## Когда нужно менять contracts

Изменения в `contracts` и повторная генерация нужна, когда:
- добавляется новый RPC
- меняется структура сообщений
- меняется контракт между `gateway` и внутренним сервисом

Рабочий порядок такой:

1. меняешь `.proto`
2. запускаешь генерацию
3. обновляешь код сервисов и `gateway`
4. проверяешь, что новые сообщения и RPC используются согласованно
