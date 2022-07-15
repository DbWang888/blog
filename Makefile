migration_init:
	migrate create -ext sql -dir db/migration -seq init_schema

migration_up:
	migrate -path db/migration -database "mysql://root:4524@tcp(localhost:3306)/blog" -verbose up

migration_down:
	migrate -path db/migration -database "mysql://root:4524@tcp(localhost:3306)/blog" -verbose down

sqlc-init:
	docker run --rm -v "d:\ProgramPro\workspace\go\blog:/src" -w /src kjconroy/sqlc init

sqlc-generate:
	docker run --rm -v "d:\ProgramPro\workspace\go\blog:/src" -w /src kjconroy/sqlc generate

.PHONY: migration migration_up migration_down sqlc-init sqlc-generate