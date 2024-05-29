.PHONY: all \
	install \
	update \
	update-all \
	format \
	vet \
	serve \
	serve-race \
	logs \
	build \
	test \
	bench \
	clean \
	help \
	test-cover-count \
	cover-count \
	test-cover-atomic \
	cover-atomic \
	html-cover-count \
	html-cover-atomic \
	run-cover-count \
	run-cover-atomic \
	view-cover-count \
	view-cover-atomic \
	docker \
	docker-up \
	docker-up-no-daemon \
	docker-down \
	docker-down-rm \
	docker-cli-register \
	docker-cli-build

.DEFAULT_GOAL=help

include .env

# Read: https://kodfabrik.com/journal/a-good-makefile-for-go

# Go parameters
CURRENT_PATH=$(shell pwd)
MAIN_PATH=$(CURRENT_PATH)/cmd/main.go
GO_CMD=go
GO_INSTALL=$(GO_CMD) install
GO_RUN=$(GO_CMD) run
GO_BUILD=$(GO_CMD) build
GO_CLEAN=$(GO_CMD) clean
GO_TEST=$(GO_CMD) test
GO_GET=$(GO_CMD) get
GO_MOD=$(GO_CMD) mod
GO_TOOL=$(GO_CMD) tool
GO_VET=$(GO_CMD) vet
GO_FMT=$(GO_CMD) fmt
BINARY_NAME=fiber-boilerplate
BINARY_UNIX=$(BINARY_NAME)_unix
DOCKER_COMPOSE=docker-compose
DOCKER=docker

## all: Test and build application
all: test build

## install: Run go install
install:
	$(GO_INSTALL) ./...

## update: Update modules
update:
	$(GO_GET) -u ./... && $(GO_MOD) tidy

## update-all: Update all modules
update-all:
	$(GO_GET) -u ./... all && $(GO_MOD) tidy

## format: Run go fmt
format:
	$(GO_FMT) ./...

## vet: Run go vet
vet: format
	$(GO_VET) ./...

## serve: Serve API
serve:
	$(GO_RUN) $(MAIN_PATH) run

## serve-race: Serve API with -race option
serve-race:
	$(GO_RUN) run -race $(MAIN_PATH)

## logs: Display server logs
logs:
	$(GO_RUN) $(MAIN_PATH) logs --server

build: format
	$(GO_BUILD) -ldflags "-s -w" -o $(BINARY_NAME) -v $(MAIN_PATH)

## test: Run test
test:
	$(GO_TEST) -cover -v ./...

test-cover-count: 
	$(GO_TEST) -covermode=count -coverprofile=cover-count.out ./...

test-cover-atomic: 
	$(GO_TEST) -covermode=atomic -coverprofile=cover-atomic.out ./...

cover-count:
	$(GO_TOOL) cover -func=cover-count.out

cover-atomic:
	$(GO_TOOL) cover -func=cover-atomic.out

html-cover-count:
	$(GO_TOOL) cover -html=cover-count.out

html-cover-atomic:
	$(GO_TOOL) cover -html=cover-atomic.out

run-cover-count: test-cover-count cover-count
run-cover-atomic: test-cover-atomic cover-atomic
view-cover-count: test-cover-count html-cover-count
view-cover-atomic: test-cover-atomic html-cover-atomic

## bench: Run benchmarks
bench: 
	$(GO_TEST) -benchmem -bench=. ./...

## docker: Stop running containers, build docker-compose.yml file and run containers
docker: docker-down docker-up

## docker-up: Build docker-compose.yml file and run containers
docker-up:
	$(DOCKER_COMPOSE) up --build --force-recreate -d

## docker-up-no-daemon: Build docker-compose.yml file and run containers in non daemon mode
docker-up-no-daemon:
	$(DOCKER_COMPOSE) up --build --force-recreate

## docker-down: Stop running containers
docker-down:
	$(DOCKER_COMPOSE) down --remove-orphans

## docker-down-rm: Stop running containers and remove linked volumes
docker-down-rm:
	$(DOCKER_COMPOSE) down --remove-orphans --volumes

## docker-cli-build: Build project for CLI
docker-cli-build:
	$(DOCKER) build -f Dockerfile -t fiber-boilerplate-cli .

## docker-cli-register: Run CLI container to register an admin user
docker-cli-register: docker-cli-build
	$(DOCKER) run -i --rm --net fiber-boilerplate_backend --link fiber-boilerplate-mysql fiber-boilerplate-cli register -l Admin -f Admin -e admin@gmail.com -p 'K-qy,Kg{<AB*XX;V3}_/x19u>1BBl!d'

## clean: Clean files
clean: 
	$(GO_CLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

help: Makefile
	@echo
	@echo "Choose a command run in "$(APP_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
	@echo
