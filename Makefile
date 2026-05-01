.PHONY: dev start build clean

dev:
	go run cmd/main.go

start: build
	./bin/api-gateway

build:
	go build -o bin/api-gateway cmd/main.go

clean:
	rm -rf bin/
