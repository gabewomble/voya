include .env
# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run: Run the application
run:
	@go run cmd/api/main.go

## test: Test the application
test:
	@echo "Testing..."
	@go test ./... -v


## itest: Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

## watch: Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

# ==================================================================================== #
# DATABASE
# ==================================================================================== #

## db/sql/fmt: Formats sql files
.PHONY: db/sql/fmt
db/sql/fmt:
	@echo 'Formatting sql files...'
	sleek ./queries/*.sql ./migrations/*.*.sql

## db/sql/gen: Generate repository package from sql queries
.PHONY: db/sql/gen
db/sql/gen:
	@echo 'Generating repository package from sql queries'
	@sqlc generate

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}/${DB_DATABASE} up

## db/migrations/down: apply all up database migrations
.PHONY: db/migrations/down
db/migrations/down: confirm
	@echo 'Running down migrations...'
	migrate -path ./migrations -database postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}/${DB_DATABASE} down


## db/seed: apply all up database migrations
.PHONY: db/seed
db/seed:
	@echo 'Seeding database...'
	psql postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}/${DB_DATABASE} -f ./seed/users.sql

# ==================================================================================== #
# BUILD
# ==================================================================================== #

## all: Build the application
.PHONY: all
all: build

.PHONY: build
build:
	@echo "Building...";
	@go build -o main cmd/api/main.go;

## clean: Clean the binary
.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -f main

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## fmt: format code
.PHONY: fmt
fmt: db/sql/fmt
	@echo 'Formatting code...'
	go fmt ./...


## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit: vendor fmt
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor

# ==================================================================================== #
# DOCKER
# ==================================================================================== #

## docker-run: Create DB container
.PHONY: docker-run
docker-run:
	@if docker compose up 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up; \
	fi

## docker-down: Shutdown DB container
.PHONY: docker-down
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi
