DOCKER_DATABASE_URL ?= postgres://ffbc:ffbc@postgres:5432/ffbc?sslmode=disable
MIGRATIONS_DIR ?= migrations
export MIGRATIONS_DIR

.PHONY: help docker-up docker-down postgres-up adminctl create-admin migrate-create migrate-up migrate-down migrate-version migrate-force

docker-up:
	docker compose up -d

docker-down:
	docker compose down

postgres-up:
	docker compose up -d postgres

adminctl:
	@test -n "$(args)" || (echo "Usage: make adminctl args='create-admin --email admin@example.com --password change-me'" && exit 1)
	docker compose run --rm adminctl $(args)

create-admin:
	@test -n "$(email)" || (echo "Usage: make create-admin email=admin@example.com password=change-me [display_name=Administrator]" && exit 1)
	@test -n "$(password)" || (echo "Usage: make create-admin email=admin@example.com password=change-me [display_name=Administrator]" && exit 1)
	docker compose run --rm adminctl create-admin --email "$(email)" --password "$(password)" $(if $(display_name),--display-name "$(display_name)")

migrate-create:
	@test -n "$(name)" || (echo "Usage: make migrate-create name=create_table_name" && exit 1)
	mkdir -p $(MIGRATIONS_DIR)
	docker compose run --rm --user "$(shell id -u):$(shell id -g)" migrate-create create -ext sql -dir /migrations -seq $(name)

migrate-up:
	docker compose run --rm --build migrate -path /migrations -database "$(DOCKER_DATABASE_URL)" up

migrate-down:
	docker compose run --rm --build migrate -path /migrations -database "$(DOCKER_DATABASE_URL)" down 1

migrate-version:
	docker compose run --rm --build migrate -path /migrations -database "$(DOCKER_DATABASE_URL)" version

migrate-force:
	@test -n "$(version)" || (echo "Usage: make migrate-force version=1" && exit 1)
	docker compose run --rm --build migrate -path /migrations -database "$(DOCKER_DATABASE_URL)" force $(version)
