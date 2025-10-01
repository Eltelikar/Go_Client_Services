MAINFILE=./cmd/client-services/main.go
SERVICE_NAME=app-builder

all: build
	docker-compose up

test:
	go test -count=5 ./internal/graph/

build:
	docker-compose build

clean-build: clear-build build
	docker-compose up

clear-build:
	docker-compose down --rmi all --volumes