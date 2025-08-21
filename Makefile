.Phony:build

build:
	@go build -o ./bin/bitcoin

.Phony:start

start: build
	./bin/bitcoin