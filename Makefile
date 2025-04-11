coem:
	go run ./cmd/gocoem/.

coemcons:
	go run ./cmd/gocoemcons/.

wsfe:
	go run ./cmd/gowsfe/.

coem-swagger:
	swag init --parseDependency --dir=./cmd/gocoem/ --output=./cmd/gocoem/docs/
