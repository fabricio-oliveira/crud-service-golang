
dependecies:
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# docker-build: @build docker image
docker-build: 
	@docker build -t test -f ./build/Dockerfile .

# build: @build app
build:
	@go build ./...

# debug: @debug start the app in debug mode
debug: dependecies
	@docker-compose up -d db db-provisioner
	@air

# lint: @lint check code quality
lint: dependecies
	golangci-lint run

# unit_test: @unit_test run the unit tests
unit_test: dependecies
	go test -ldflags="-s=false" -gcflags="-l" -v ./...