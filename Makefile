NAME=gymo
VERSION=0.0.1

LDFLAGS = "-s -w -buildid="

build:
	@CGO_ENABLED=0 go build -ldflags=$(LDFLAGS) -trimpath -v -o $(NAME)

dev:
	air --build.cmd "CGO_ENABLED=0 go build -o $(NAME)" --build.bin "./$(NAME)"

clean:
	go clean -cache
	@rm -f $(NAME)

deps:
	@go mod download -x

test:
	@go test -v ./...

build-docker:
	docker build --progress=plain -t $(NAME) .

all: build

help: Makefile
	@echo
	@echo " Choose a command to run in "$(APP_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

.PHONY: all
