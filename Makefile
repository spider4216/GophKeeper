APP_NAME=gophkeeper

DSN?=postgres://postgres:postgres@localhost:5432/gophkeeper
DSN_CLIENT?=postgres://postgres:postgres@localhost:5432/gophkeeper_client_2

LOG_LEVEL?=debug
JWT_KEY?=qwerty
ENCRYPT_KEY?=12345678901234567890123456789012
SERVER_HOST?=http://127.0.0.1:8080
CLI_VERSION?=1.0.0

run-srv:
	JWT_KEY=${JWT_KEY} SERVER_ADDRESS=:8080 DB_DSN=$(DSN) LOG_LEVEL=${LOG_LEVEL} go run ./cmd/server

run-srv-ssl:
	ENABLE_HTTPS=true JWT_KEY=${JWT_KEY} SERVER_ADDRESS=:8080 DB_DSN=$(DSN) LOG_LEVEL=${LOG_LEVEL} go run ./cmd/server

client-reg:
	SERVER_HOST=${SERVER_HOST} LOG_LEVEL=${LOG_LEVEL} DB_DSN=$(DSN_CLIENT) go run ./cmd/client register --email=${email} --password=${pass}

client-login:
	SERVER_HOST=${SERVER_HOST} JWT_KEY=${JWT_KEY} LOG_LEVEL=${LOG_LEVEL} DB_DSN=$(DSN_CLIENT) go run ./cmd/client login --email=${email} --password=${pass}

client-insert-loginpass:
	SERVER_HOST=${SERVER_HOST} JWT_KEY=${JWT_KEY} LOG_LEVEL=${LOG_LEVEL} ENCRYPT_KEY=${ENCRYPT_KEY} DB_DSN=$(DSN_CLIENT) go run ./cmd/client insert-loginpass --login=${login} --password=${pass} --title=${title} --user-id=${userid}

client-update-loginpass:
	SERVER_HOST=${SERVER_HOST} JWT_KEY=${JWT_KEY} LOG_LEVEL=${LOG_LEVEL} ENCRYPT_KEY=${ENCRYPT_KEY} DB_DSN=$(DSN_CLIENT) go run ./cmd/client update-loginpass --login=${login} --password=${pass} --id=${id} --meta-id=${metaid} --title=${title} --user-id=${userid}

client-delete-item:
	SERVER_HOST=${SERVER_HOST} JWT_KEY=${JWT_KEY} LOG_LEVEL=${LOG_LEVEL} ENCRYPT_KEY=${ENCRYPT_KEY} DB_DSN=$(DSN_CLIENT) go run ./cmd/client delete-item --id=${id} --user-id=${userid}

client-sync-send:
	SERVER_HOST=${SERVER_HOST} JWT_KEY=${JWT_KEY} LOG_LEVEL=${LOG_LEVEL} ENCRYPT_KEY=${ENCRYPT_KEY} DB_DSN=$(DSN_CLIENT) go run ./cmd/client sync-send --user-id=${userid}

client-sync-get:
	SERVER_HOST=${SERVER_HOST} JWT_KEY=${JWT_KEY} LOG_LEVEL=${LOG_LEVEL} ENCRYPT_KEY=${ENCRYPT_KEY} DB_DSN=$(DSN_CLIENT) go run ./cmd/client sync-get --user-id=${userid}

client-list:
	SERVER_HOST=${SERVER_HOST} JWT_KEY=${JWT_KEY} LOG_LEVEL=${LOG_LEVEL} ENCRYPT_KEY=${ENCRYPT_KEY} DB_DSN=$(DSN_CLIENT) go run ./cmd/client list --user-id=${userid}

client-view:
	SERVER_HOST=${SERVER_HOST} JWT_KEY=${JWT_KEY} LOG_LEVEL=${LOG_LEVEL} ENCRYPT_KEY=${ENCRYPT_KEY} DB_DSN=$(DSN_CLIENT) go run ./cmd/client view --item_id=${id} --user-id=${userid}

client-version:
	DB_DSN=$(DSN_CLIENT) go run ./cmd/client version

client-create-text:
	SERVER_HOST=${SERVER_HOST} JWT_KEY=${JWT_KEY} LOG_LEVEL=${LOG_LEVEL} ENCRYPT_KEY=${ENCRYPT_KEY} DB_DSN=$(DSN_CLIENT) go run ./cmd/client create-text --path=${path} --user-id=${userid} --title=${title}

migration-gen:
	migrate create -ext sql -dir ./migrations/${path} -seq $(name)

migrate-up:
	migrate -path ./migrations/${path} -database $(DSN) up $(ver)

migrate-down:
	migrate -path ./migrations/${path} -database $(DSN) down $(ver)

migrate-force:
	migrate -path ./migrations/${path} -database $(DSN) force $(ver)

lint:
	golangci-lint run

fmt:
	golangci-lint fmt

crt:
	go run ./cmd/crtcert

build-client-linux-amd64:
	GOOS=linux GOARCH=amd64 go build \
		-ldflags "-X github.com/spider4216/GophKeeper/internal/client/version.Version=$(CLI_VERSION) -X 'github.com/spider4216/GophKeeper/internal/client/version.BuildDate=$(shell date +'%Y/%m/%d %H:%M:%S')'" \
		-o bin/client-linux-amd64 \
		./cmd/client

build-client-linux-arm64:
	GOOS=linux GOARCH=arm64 go build \
		-ldflags "-X github.com/spider4216/GophKeeper/internal/client/version.Version=$(CLI_VERSION) -X 'github.com/spider4216/GophKeeper/internal/client/version.BuildDate=$(shell date +'%Y/%m/%d %H:%M:%S')'" \
		-o bin/client-linux-arm64 \
		./cmd/client

build-client-linux-arm32:
	GOOS=linux GOARCH=arm GOARM=7 go build \
		-ldflags "-X github.com/spider4216/GophKeeper/internal/client/version.Version=$(CLI_VERSION) -X 'github.com/spider4216/GophKeeper/internal/client/version.BuildDate=$(shell date +'%Y/%m/%d %H:%M:%S')'" \
		-o bin/client-linux-arm32 \
		./cmd/client

build-client-win-amd64:
	GOOS=windows GOARCH=amd64 go build \
		-ldflags "-X github.com/spider4216/GophKeeper/internal/client/version.Version=$(CLI_VERSION) -X 'github.com/spider4216/GophKeeper/internal/client/version.BuildDate=$(shell date +'%Y/%m/%d %H:%M:%S')'" \
		-o bin/client-win-amd64 \
		./cmd/client

build-client-win-arm64:
	GOOS=windows GOARCH=arm64 go build \
		-ldflags "-X github.com/spider4216/GophKeeper/internal/client/version.Version=$(CLI_VERSION) -X 'github.com/spider4216/GophKeeper/internal/client/version.BuildDate=$(shell date +'%Y/%m/%d %H:%M:%S')'" \
		-o bin/client-win-arm64 \
		./cmd/client

build-client-win-386:
	GOOS=windows GOARCH=386 go build \
		-ldflags "-X github.com/spider4216/GophKeeper/internal/client/version.Version=$(CLI_VERSION) -X 'github.com/spider4216/GophKeeper/internal/client/version.BuildDate=$(shell date +'%Y/%m/%d %H:%M:%S')'" \
		-o bin/client-win-386 \
		./cmd/client

build-client-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build \
		-ldflags "-X github.com/spider4216/GophKeeper/internal/client/version.Version=$(CLI_VERSION) -X 'github.com/spider4216/GophKeeper/internal/client/version.BuildDate=$(shell date +'%Y/%m/%d %H:%M:%S')'" \
		-o bin/client-darwin-amd64 \
		./cmd/client

build-client-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build \
		-ldflags "-X github.com/spider4216/GophKeeper/internal/client/version.Version=$(CLI_VERSION) -X 'github.com/spider4216/GophKeeper/internal/client/version.BuildDate=$(shell date +'%Y/%m/%d %H:%M:%S')'" \
		-o bin/client-darwin-arm64 \
		./cmd/client