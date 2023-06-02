PROJECT_NAME := "github.com/christophwitzko/flight-booking-service"
PKG := "$(PROJECT_NAME)"

.PHONY: dep build

all: build

dep: ## Get the dependencies
	@go mod download
	@go mod vendor

build: dep ## Build the binary file
	@go build -o build/ -v $(PKG)/cmd/flight-booking-service
