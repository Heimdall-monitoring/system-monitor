all: build

build:
	go build -o build

clean:
	go clean

test:
	mkdir -p build/
	# go test -v -coverprofile=build/c.out .
	# go tool cover -html=build/c.out -o build/coverage.html

.PHONY: all test clean build
