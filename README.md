# FFBC Backend


## Требования

- Go 1.26+
- Docker и Docker Compose
- [sqlc](https://sqlc.dev/) для регенерации PostgreSQL-кода

Установка CLI:

```bash
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
make create-admin \
  email=admin@example.com \
  password='change-me' \
  display_name='Administrator'
```

Параметры:

| Параметр | Обязательный | Назначение |
|---|---:|---|
| `--email` | Да | Email администратора |
| `--password` | Да | Пароль администратора |
| `--display-name` | Нет | Отображаемое имя |


Для запуска произвольной команды `adminctl` внутри контейнера:

```bash
make adminctl args='create-admin --email admin@example.com --password change-me'
```

После успешного выполнения CLI выводит ID и email созданного администратора.

Миграции также запускаются внутри отдельного Docker-контейнера:

```bash
make migrate-up
make migrate-down
make migrate-version
make migrate-force version=2
make migrate-create name=add_example_table
```

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
