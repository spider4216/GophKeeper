APP_NAME=gophkeeper

DSN?=postgres://postgres:postgres@localhost:5432/gophkeeper
LOG_LEVEL?=debug

run-srv:
	SERVER_ADDRESS=:8080 DB_DSN=$(DSN) LOG_LEVEL=${LOG_LEVEL} go run ./cmd/server

migration-gen:
	migrate create -ext sql -dir ./migrations/${path} -seq $(name)

migrate-up:
	migrate -path ./migrations/${path} -database $(DSN) up $(ver)

migrate-down:
	migrate -path ./migrations/${path} -database $(DSN) down $(ver)

migrate-force:
	migrate -path ./migrations/${path} -database $(DSN) force $(ver)