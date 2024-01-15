APP = ash
BUILD_DIR = build
BUILD_TIME = $(shell date +%s)
BUILD_LCOMMIT =$(shell git log --pretty=format:"%s"  | head -n 1)


.PHONY: install build all test clean exec

clean:
	rm -rf ./${BUILD_DIR}/${APP}

test:
	go test ./...

build: clean test
	go build -o ./${BUILD_DIR}/${APP} ./cmd/main.go

install: build
	sudo cp ./${BUILD_DIR}/${APP} /usr/local/bin/${APP} && sudo chmod +x /usr/local/bin/${APP}

exec: build
	exec  ${BUILD_DIR}/${APP}

generate:
	go generate ./...
.PHONY: generate
