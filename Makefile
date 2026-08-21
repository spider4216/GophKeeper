APP_NAME=gophkeeper

DSN?=postgres://postgres:postgres@localhost:5432/gophkeeper
DSN_CLIENT?=postgres://postgres:postgres@localhost:5432/gophkeeper_client

LOG_LEVEL?=debug
JWT_KEY?=qwerty
ENCRYPT_KEY?=qwerty

run-srv:
	JWT_KEY=${JWT_KEY} SERVER_ADDRESS=:8080 DB_DSN=$(DSN) LOG_LEVEL=${LOG_LEVEL} go run ./cmd/server

client-reg:
	LOG_LEVEL=${LOG_LEVEL} DB_DSN=$(DSN_CLIENT) go run ./cmd/client register --email=${email} --password=${pass}

client-login:
	LOG_LEVEL=${LOG_LEVEL} DB_DSN=$(DSN_CLIENT) go run ./cmd/client login --email=${email} --password=${pass}

client-insert-loginpass:
	LOG_LEVEL=${LOG_LEVEL} ENCRYPT_KEY=${ENCRYPT_KEY} DB_DSN=$(DSN_CLIENT) go run ./cmd/client  insert-loginpass --login=${login} --password=${pass} --title=${title} --token=${jwt}

client-sync-send:
	LOG_LEVEL=${LOG_LEVEL} ENCRYPT_KEY=${ENCRYPT_KEY} DB_DSN=$(DSN_CLIENT) go run ./cmd/client  sync-send --token=${jwt}

client-list:
	LOG_LEVEL=${LOG_LEVEL} ENCRYPT_KEY=${ENCRYPT_KEY} DB_DSN=$(DSN_CLIENT) go run ./cmd/client  list --token=${jwt}

migration-gen:
	migrate create -ext sql -dir ./migrations/${path} -seq $(name)

migrate-up:
	migrate -path ./migrations/${path} -database $(DSN) up $(ver)

migrate-down:
	migrate -path ./migrations/${path} -database $(DSN) down $(ver)

migrate-force:
	migrate -path ./migrations/${path} -database $(DSN) force $(ver)