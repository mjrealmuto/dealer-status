# Go params
GOCMD=go
GOIMAGE=golang:1.21
GOBUILD=$(GOCMD) build -v -o yomama

export SSH_FILE=

#Docker params
DOCKER_PATH=./docker/Dockerfile
DOCKER_COMPOSE_FILE=./docker-compose.yaml
DOCKER_BUILD=docker build . -t go-dealer-status-app:latest
DOCKER_RUN=docker compose -f '$(DOCKER_COMPOSE_FILE)' up --force-recreate --build
#DOCKER_RUN=docker run -e SSH_KEY="$(ID_RSA)" go-dealer-status-app:latest

default: run-dealer-status

build-dealer-status:
	@echo "Building dealer-status..."
	@echo $(DOCKER_BUILD)
	@$(DOCKER_BUILD)

run-dealer-status: build-dealer-status
	@echo "Running dealer-status..."
	@$(DOCKER_RUN)

	
