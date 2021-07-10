all: build run

build:
	mkdir -p build
	go build -o build

clean:
	go clean
	rm -r build

run:
	./build/monitor --config=build/config.json --debug

test:
	mkdir -p build/
	go test -parallel 10 -v -coverprofile=build/c.out ./...
	go tool cover -html=build/c.out -o build/coverage.html

test_ci:
	go test -v -race -coverprofile=coverage.txt ./...

.PHONY: all test clean build
