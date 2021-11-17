.PHONY: unit integration local lint migrate prune dev swag

BUILD_DIR = $(PWD)/build
MIGRATIONS= $(PWD)/migrations/
DATABASE = postgres://postgres:121073@localhost/compresso?sslmode=disable

SHELL := /bin/bash

unit:
	go clean -testcache
	go test -tags=unit -v -timeout 5s -cover ./...

integration:
	go clean -testcache
	go test -tags=integration -v -timeout 30s -cover ./...

lint:
	golangci-lint run ./...

swag:
	swag init -g cmd/auth/main.go -o docs/auth --exclude internal/video
	swag init -g cmd/videoapi/main.go -o docs/videoapi --exclude internal/auth

# ==============================================================================
# MIGRATIONS

migrate.up:
	migrate -path $(MIGRATIONS) -database "$(DATABASE)" up

migrate.down:
	migrate -path $(MIGRATIONS) -database "$(DATABASE)" down

# % - stands for version
migrate.force.%:
	migrate -path $(MIGRATIONS) -database "$(DATABASE)" force $*

# ==============================================================================
# Local usage of services, % in (auth, videoapi, videosvc)

%.clean:
	rm $(BUILD_DIR)/$*

%.build:
	go build -ldflags="-w -s" -o $(BUILD_DIR)/$* cmd/$*/main.go

%.run: %.clean %.build
	$(BUILD_DIR)/$*

#
# ==============================================================================
# Docker compose commands for local usage

local.up:
	echo "Starting local environment"
	docker-compose -f docker-compose.local.yml up --build -d

local.down:
	docker-compose -f docker-compose.local.yml down

#
# ==============================================================================
# Dev usage

dev.up:
	echo "Starting dev environment"
	docker-compose -f docker-compose.dev.yml up --build

dev.down:
	docker-compose -f docker-compose.dev.yml down

prune:
	docker system prune -f
