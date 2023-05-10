APP_NAME=endor
UNAME := $(shell uname -s)

## ----------------------------------------------------------------------
## This is a help comment. The purpose of this Makefile is to demonstrate
## a simple help mechanism that uses comments defined alongside the rules
## ----------------------------------------------------------------------

help: ## Show help for each of the Makefile recipes.
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST)


clean: ## Removes go executables and  redis container
	go clean
	rm -rf ${APP_NAME}-darwin
	rm -rf ${APP_NAME}-linux
	docker rm -f ${APP_NAME}-redis

build: ## Runs clean and setup to build infra and dependencies and create go executables
	$(MAKE) clean
	$(MAKE) setup
	go mod tidy
	go mod download
	GOARCH=amd64 GOOS=darwin go build -o ${APP_NAME}-darwin .
	GOARCH=amd64 GOOS=linux go build -o ${APP_NAME}-linux .

setup: ## Pull latest Redis images and run as container @ port 6739
	docker pull redis:latest
	docker run --name ${APP_NAME}-redis -p 6379:6379 -d redis

run-mac: ## Run go executable for Mac
	./endor-darwin

run-linux: ## Run go executable for linux
	./endor-linux

run-tests: ## Run unit-tests cases
	go test