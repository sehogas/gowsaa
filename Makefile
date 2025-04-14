########################################################
override TARGET=api-arca
VERSION=1.0
OS=linux
ARCH=amd64
FLAGS="-s -w"
CGO=0
########################################################

run:
	go run ./cmd/api/.

wsfe:
	go run ./cmd/gowsfe/.

swagger:
	swag init --parseDependency --dir=./cmd/api/ --output=./cmd/api/docs/

build:
	docker build -t $(TARGET):$(VERSION) .
	docker tag $(TARGET):$(VERSION) $(TARGET):latest

up:
	docker compose up -d --build

down:
	docker compose down

