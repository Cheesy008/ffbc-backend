DATABASE_URL ?= postgres://ffbc:ffbc@localhost:5432/ffbc?sslmode=disable
MIGRATIONS_DIR ?= migrations

.PHONY: help docker-up docker-down postgres-up app-up migrate-create migrate-up migrate-down migrate-version migrate-force

docker-up:
	docker compose up -d

docker-down:
	docker compose down

postgres-up:
	docker compose up -d postgres

migrate-create:
	@test -n "$(name)" || (echo "Usage: make migrate-create name=create_table_name" && exit 1)
	mkdir -p $(MIGRATIONS_DIR)
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down 1

migrate-version:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" version

migrate-force:
	@test -n "$(version)" || (echo "Usage: make migrate-force version=1" && exit 1)
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" force $(version)
