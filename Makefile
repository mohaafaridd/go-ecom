build:
	@go build -o bin/ecom cmd/main.go
test:
	@go test -b ./...
run: build
	@./bin/ecom