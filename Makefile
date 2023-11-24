dc_build:
	go mod vendor
	docker build ./ -t webhook:latest

dc_run:
	docker run webhook:latest

build:
	go build cmd/server/main.go

run:
	go run cmd/server/main.go

test:
	go test ./...

lint:
	golangci-lint run