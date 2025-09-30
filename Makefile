MAINFILE=./cmd/client-services/main.go

all: build
	docker-compose up

build:
	docker-compose build

clean-build: clear-build build
	docker-compose up

clear-build:
	docker-compose down --rmi all --volumes