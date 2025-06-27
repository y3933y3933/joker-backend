# Include variables from the .envrc file
include .envrc

# ==================================================================================== # 
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'


.PHONY: confirm
confirm:
	@echo 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]


# ==================================================================================== # 
# DEVELOPMENT
# ==================================================================================== #

## run: run the application
.PHONY: run
run:
	go run .

## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${DB_URL} 

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up:
	@echo 'Running up migrations...'
	goose -dir ./sql/migrations postgres $(DB_URL) up


## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new: confirm
	@echo 'Creating migration files for ${name}...'
	goose -dir ./sql/migrations -s create $(name) sql




# ==================================================================================== # 
# QUALITY CONTROL
# ==================================================================================== #
.PHONY: tidy
## tidy: tidy module dependencies and format all .go files
tidy:
	@echo 'Tidying module dependencies...'
	go mod tidy
	@echo 'Formatting .go files'
	go fmt ./...


## audit: run quality control checks
.PHONY: audit 
audit:
	@echo 'Checking module dependencies...'
	go mod tidy -diff
	@echo 'Vetting code...'
	go vet ./...
	go tool staticcheck ./...


# ==================================================================================== # 
# BUILD
# ==================================================================================== #

## build: build the application
.PHONY: build
build:
	@echo 'Building...'
	go build -o=./bin/app .