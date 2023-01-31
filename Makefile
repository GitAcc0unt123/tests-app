#include .env
#export

.PHONY:
.SILENT:
.DEFAULT_GOAL := docker-compose

# https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

test: ## Run unit tests
	go test --short -v -count=1 -timeout=30s -coverprofile=coverage.out ./...

test-integration: ## Run all integration tests
	docker run --rm -d --name=test_db -e POSTGRES_DB=test -e POSTGRES_USER=username -e POSTGRES_PASSWORD=qwerty123 -p 6001:5432 postgres:15
	sleep 2s
	migrate -path ./migrations -database "postgres://username:qwerty123@127.0.0.1:6001/test?sslmode=disable" up 1
	go test -v -count=1 ./integration-tests/
	docker stop test_db

cover: ## Run unit tests and open the coverage report
	go test --short -v -count=1 -timeout=30s -coverprofile=coverage.out ./...
	go tool cover -html=./coverage.out -o coverage.html
	open coverage.html

docker-compose:
	docker compose up --build --force-recreate

swagger: ## Generate api documentation
	$$(go env GOPATH)/bin/swag init
	$$(go env GOPATH)/bin/swag fmt

gomock: ## Generate service mocks
#	export PATH=$PATH:$(go env GOPATH)/bin
	go generate internal/service/service.go

install:
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golang/mock/mockgen@v1.6.0
#	go install golang.org/x/tools/gopls@latest

	curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz
	mv migrate.linux-amd64 /usr/bin/migrate

# migrate_params = -path ./migrations -database "postgres://<username>:<password>@ip-addr:port/postgres?sslmode=disable"
migrate-up:
	migrate $(migrate_params) up $(n)

migrate-down:
	migrate $(migrate_params) down $(n)

migrate-version:
	migrate $(migrate_params) version

migrate-goto:
	migrate $(migrate_params) goto $(v)

migrate-drop:
	migrate $(migrate_params) drop

migrate-create:
	migrate create -ext sql -dir ./migrations/ -seq $(NAME)