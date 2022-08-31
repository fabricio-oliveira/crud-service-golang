
# docker-build: @build docker image
docker-build: 
	@docker build -t test -f ./build/Dockerfile .

# build: @build app
build:
	@go build ./...

# debug: @debug start the app in debug mode
debug:
	@go run ./...