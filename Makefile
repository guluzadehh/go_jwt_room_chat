build:
	@go build -o bin/gochat cmd/gochat/main.go

run: build
	@./bin/gochat