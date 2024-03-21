NAME=gymo
VERSION=0.0.1

.PHONY: build
## build: Compile the packages.
build:
	@CGO_ENABLED=0 go build -v -o $(NAME)

.PHONY: run
## run: Build and Run in development mode.
run:
	@nodemon --exec go run main.go --signal SIGTERM

.PHONY: run-prod
## run-prod: Build and Run in production mode.
run-prod:
	@nodemon --exec go run main.go --signal SIGTERM

.PHONY: clean
## clean: Clean project and previous builds.
clean:
	@rm -f $(NAME)

.PHONY: deps
## deps: Download modules
deps:
	@go mod download -x

.PHONY: test
## test: Run tests with verbose mode
test:
	@go test -v ./...

.PHONY: help
all: help
# help: show this help message
help: Makefile
	@echo
	@echo " Choose a command to run in "$(APP_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
