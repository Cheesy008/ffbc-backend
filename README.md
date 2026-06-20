# FFBC Backend


## Требования

- Go 1.26+
- Docker и Docker Compose
- [golang-migrate](https://github.com/golang-migrate/migrate) для работы с миграциями через Makefile
- [sqlc](https://sqlc.dev/) для регенерации PostgreSQL-кода

Установка CLI:

```bash
go install -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

## Быстрый запуск

Запустить PostgreSQL и приложение:

```bash
make docker-up
```

Применить миграции:

```bash
make migrate-up
```

Приложение будет доступно по адресу:

```text
http://localhost:8080
```

Документация:

```text
http://localhost:8080/api/docs
```

## Makefile и Docker
Для запуска Docker контейнеров используйте Makefile, все команды описаны внутри.


## adminctl

`adminctl` — CLI для административных операций. Сейчас поддерживается команда создания администратора.

Перед запуском примените миграции:

```bash
make postgres-up
make migrate-up
```

Создать администратора:

```bash
go run ./cmd/adminctl create-admin \
  --email admin@example.com \
  --password 'change-me' \
  --display-name 'Administrator'
```

Параметры:

| Параметр | Обязательный | Назначение |
|---|---:|---|
| `--email` | Да | Email администратора |
| `--password` | Да | Пароль администратора |
| `--display-name` | Нет | Отображаемое имя |


Можно собрать отдельный бинарный файл:

```bash
go build -o adminctl ./cmd/adminctl

./adminctl create-admin \
  --email admin@example.com \
  --password 'change-me'
```

После успешного выполнения CLI выводит ID и email созданного администратора.

## Работа с sqlc

После изменения SQL-запросов или схемы миграций перегенерируйте код:

```bash
sqlc generate
```

Сгенерированные пакеты находятся в:

```text
internal/admin/repository/postgres/sqlc/generated
internal/catalog/repository/postgres/sqlc/generated
```

## Проверка проекта

Форматирование:

```bash
gofmt -w $(find cmd internal -name '*.go')
```

Сборка и тесты:

```bash
go test ./...
```

Сборка приложения:

```bash
go build -o app ./cmd/app
```
