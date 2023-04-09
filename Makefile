COVERAGE := TRUE
DISPLAY_COVERAGE := TRUE
HTML_FILE:= coverage.html
JUNIT_FILE:= rspec.xml
COBERTURA_FILE:= coverage.xml
REPORT_PATH:= ./test-reports

gen:
	## Go generate
	@go generate ./...

protoc_gen:
	protoc \
      -I . -I ${GOPATH}/src \
      --go_out=plugins=grpc:. \
      --go_opt=paths=source_relative \
	  --experimental_allow_proto3_optional \
      	proto/*/*.proto

run:
	cd cmd && go run main.go

##@ Development

PROJECT_NAME ?= chat-app
DOCKER_COMPOSE ?= cd dev && docker compose -p $(PROJECT_NAME)

dev-up: ## Run local environment for developing locally
	$(DOCKER_COMPOSE) up -d

dev-down: ## Shutdown the local environment
	$(DOCKER_COMPOSE) down

dev-build:
	@echo "==> Building image <=="
	@bash ./dev/build-image.sh