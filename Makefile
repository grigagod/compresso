.PHONY: clean test docker auth

BUILD_DIR = $(PWD)/build
AUTH_MIGRATIONS= $(PWD)/migrations/auth
AUTH_DATABASE = postgres://postgres:121073@localhost/auth?sslmode=disable

test: 
	go test -v -timeout 30s -cover ./...

auth.migrate.up:
	migrate -path $(AUTH_MIGRATIONS) -database "$(AUTH_DATABASE)" up

auth.migrate.down:
	migrate -path $(AUTH_MIGRATIONS) -database "$(AUTH_DATABASE)" down

auth.migrate.force:
	migrate -path $(AUTH_MIGRATIONS) -database "$(AUTH_DATABASE)" force $(version)

lint:
	golangci-lint run ./...

swag:
	echo "Starting swagger generating"
	swag init -g **/**/*.go

# ==============================================================================
# Docker compose commands

docker.build:
	echo "Starting docker environment"
	docker-compose -f docker-compose.yml up --build

docker.run:
	docker-compose -f docker-compose.yml up
# ==============================================================================
# Docker support

FILES := $(shell docker ps -aq)

down-local:
	docker stop $(FILES)
	docker rm $(FILES)

prune:
	docker system prune -f

logs-local:
	docker logs -f $(FILES)
