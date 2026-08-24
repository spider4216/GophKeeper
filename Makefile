APP_NAME=gophkeeper

DSN?=postgres://postgres:postgres@localhost:5432/gophkeeper
DSN_CLIENT?=postgres://postgres:postgres@localhost:5432/gophkeeper_client_2

LOG_LEVEL?=debug
JWT_KEY?=qwerty
ENCRYPT_KEY?=12345678901234567890123456789012
SERVER_HOST?=127.0.0.1:8080

run-srv:
	JWT_KEY=${JWT_KEY} SERVER_ADDRESS=:8080 DB_DSN=$(DSN) LOG_LEVEL=${LOG_LEVEL} go run ./cmd/server

client-reg:
	SERVER_HOST=${SERVER_HOST} LOG_LEVEL=${LOG_LEVEL} DB_DSN=$(DSN_CLIENT) go run ./cmd/client register --email=${email} --password=${pass}

client-login:
	SERVER_HOST=${SERVER_HOST} JWT_KEY=${JWT_KEY} LOG_LEVEL=${LOG_LEVEL} DB_DSN=$(DSN_CLIENT) go run ./cmd/client login --email=${email} --password=${pass}

client-insert-loginpass:
	SERVER_HOST=${SERVER_HOST} JWT_KEY=${JWT_KEY} LOG_LEVEL=${LOG_LEVEL} ENCRYPT_KEY=${ENCRYPT_KEY} DB_DSN=$(DSN_CLIENT) go run ./cmd/client insert-loginpass --login=${login} --password=${pass} --title=${title} --token=${jwt}

client-update-loginpass:
	SERVER_HOST=${SERVER_HOST} JWT_KEY=${JWT_KEY} LOG_LEVEL=${LOG_LEVEL} ENCRYPT_KEY=${ENCRYPT_KEY} DB_DSN=$(DSN_CLIENT) go run ./cmd/client update-loginpass --login=${login} --password=${pass} --id=${id} --token=${jwt}

client-delete-item:
	SERVER_HOST=${SERVER_HOST} JWT_KEY=${JWT_KEY} LOG_LEVEL=${LOG_LEVEL} ENCRYPT_KEY=${ENCRYPT_KEY} DB_DSN=$(DSN_CLIENT) go run ./cmd/client delete-item --id=${id} --token=${jwt}

client-sync-send:
	SERVER_HOST=${SERVER_HOST} JWT_KEY=${JWT_KEY} LOG_LEVEL=${LOG_LEVEL} ENCRYPT_KEY=${ENCRYPT_KEY} DB_DSN=$(DSN_CLIENT) go run ./cmd/client sync-send --token=${jwt}

client-sync-get:
	SERVER_HOST=${SERVER_HOST} JWT_KEY=${JWT_KEY} LOG_LEVEL=${LOG_LEVEL} ENCRYPT_KEY=${ENCRYPT_KEY} DB_DSN=$(DSN_CLIENT) go run ./cmd/client sync-get --token=${jwt}

client-list:
	SERVER_HOST=${SERVER_HOST} JWT_KEY=${JWT_KEY} LOG_LEVEL=${LOG_LEVEL} ENCRYPT_KEY=${ENCRYPT_KEY} DB_DSN=$(DSN_CLIENT) go run ./cmd/client list --token=${jwt}

client-view:
	SERVER_HOST=${SERVER_HOST} JWT_KEY=${JWT_KEY} LOG_LEVEL=${LOG_LEVEL} ENCRYPT_KEY=${ENCRYPT_KEY} DB_DSN=$(DSN_CLIENT) go run ./cmd/client view --item_id=${id} --token=${jwt}

migration-gen:
	migrate create -ext sql -dir ./migrations/${path} -seq $(name)

migrate-up:
	migrate -path ./migrations/${path} -database $(DSN) up $(ver)

migrate-down:
	migrate -path ./migrations/${path} -database $(DSN) down $(ver)

migrate-force:
	migrate -path ./migrations/${path} -database $(DSN) force $(ver)