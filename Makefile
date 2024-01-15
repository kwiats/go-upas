build:
	@go build -o bin/go-upas cmd/rest/main.go

run: build
	@./bin/go-upas 

test: 
	@go test -v ./...

