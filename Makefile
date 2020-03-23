.PHONY: default test build

OS := $(shell uname)
VERSION ?= 1.0.0

# target #

default: unit-test integration-test build-api build-cron

build-api:
	@echo "Setup covid-api"
ifeq ($(OS),Linux)
	@echo "Build covid-api..."
	GOOS=linux  go build -ldflags "-s -w -X main.Version=$(VERSION)" -o api cmd/api/main.go
endif
ifeq ($(OS) ,Darwin)
	@echo "Build covid-api..."
	GOOS=darwin go build -ldflags "-X main.Version=$(VERSION)" -o api cmd/api/main.go
endif
	@echo "Succesfully Build for ${OS} version:= ${VERSION}"

build-cron:
	@echo "Setup covid-cron"
ifeq ($(OS),Linux)
	@echo "Build covid-cron..."
	GOOS=linux  go build -ldflags "-s -w -X main.Version=$(VERSION)" -o cron cmd/cron/main.go
endif
ifeq ($(OS) ,Darwin)
	@echo "Build covid-cron..."
	GOOS=darwin go build -ldflags "-X main.Version=$(VERSION)" -o cron cmd/cron/main.go
endif
	@echo "Succesfully Build for ${OS} version:= ${VERSION}"


# Test Packages

unit-test:
	@go test -count=1 -v --cover ./... -tags="unit"

integration-test:
	@go test -count=1 -v --cover -tags="integration" -p 1 ./... --env-path=`pwd`/.env --config-yaml=`pwd`/config.yaml
