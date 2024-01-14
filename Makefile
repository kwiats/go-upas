build:
	@go build -o bin/go-upas

run: build
	@./bin/go-upas

test: 
	@go test -v ./...

