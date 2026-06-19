.PHONY: run build test swagger

run:
	go run cmd/server/main.go

build:
	go build -o bin/server cmd/server/main.go

test:
	go test ./...

swagger:
	swag init -g cmd/server/main.go --parseInternal

swagger-fmt:
	swag fmt
