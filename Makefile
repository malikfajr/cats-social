add_db="postgres://postgres:secret@localhost:5432/cats-social?sslmode=disable"

postgres:
	docker run --name pg-sprint -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=secret -d postgres:15-alpine

createdb:
	docker exec -it pg-sprint createdb --username=postgres cats-social

dropdb:
	docker exec -it pg-sprint dropdb --username=postgres cats-social

migrateup:
	migrate -database ${add_db} -path db/migrations up

migratedown:
	migrate -database $(add_db) -path db/migrations down

.PHONY: postgres createdb dropdb


# migrate database "postgres://postgres:secretly@localhost:5432/cats-social?sslmode=disable" -path db/migrations up
