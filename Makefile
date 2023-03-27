# To initialise our database schema
# migrate create -ext sql -dir db/migration -seq init_schema

sqlc:
	sqlc generate
rmpostgres:
	docker rm some-postgres
postgres:
	docker run --name some-postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=ritik -d postgres:12-alpine
createdb:
	docker exec -it some-postgres createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it some-postgres dropdb simple_bank
migrateup:
	migrate -path ./db/migration -database "postgresql://root:ritik@localhost:5432/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path ./db/migration -database "postgresql://root:ritik@localhost:5432/simple_bank?sslmode=disable" -verbose down
test:
	go test -v -cover ./...
.PHONY:	createdb postgres dropdb rmpostgres migrateup migratedown sqlc