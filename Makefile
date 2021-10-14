.PHONY: clean test auth build run local videoapi

BUILD_DIR = $(PWD)/build
MIGRATIONS= $(PWD)/migrations/
DATABASE = postgres://postgres:121073@localhost/compresso?sslmode=disable

test:
	go test -v -timeout 30s -cover ./...

lint:
	golangci-lint run ./...

# ==============================================================================
# MIGRATIONS


migrate.up:
	migrate -path $(MIGRATIONS) -database "$(DATABASE)" up

migrate.down:
	migrate -path $(MIGRATIONS) -database "$(DATABASE)" down

migrate.force:
	migrate -path $(MIGRATIONS) -database "$(DATABASE)" force $(version)

# ==============================================================================
# Auth service

auth.clean:
	rm $(BUILD_DIR)/auth

auth.build: clean
	go build -ldflags="-w -s" -o $(BUILD_DIR)/auth cmd/auth/main.go

auth.run: auth.clean auth.build
	$(BUILD_DIR)/auth

auth.swag:
	swag init -g cmd/auth/main.go -o docs/auth --exclude internal/image


# ===================================================
# Videoapi
videoapi.clean:
	rm $(BUILD_DIR)/videoapi

videoapi.build:
	go build -ldflags="-w -s" -o $(BUILD_DIR)/videoapi cmd/videoapi/main.go

videoapi.run: videoapi.clean videoapi.build
	$(BUILD_DIR)/videoapi



# ===================================================
# Videosvc
videosvc.clean:
	rm $(BUILD_DIR)/videosvc

videosvc.build:
	go build -ldflags="-w -s" -o $(BUILD_DIR)/videosvc cmd/videosvc/main.go

videosvc.run: videosvc.clean videosvc.build
	$(BUILD_DIR)/videosvc


# ==============================================================================
# Docker compose commands

local:
	echo "Starting local environment"
	docker-compose -f docker-compose.local.yml up --build -d

local.stop:
	docker-compose -f docker-compose.local.yml down

develop:
	echo "Starting docker environment"
	docker-compose -f docker-compose.dev.yml up --build

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
