APP=ash
BRANCH_NAME=$(shell git rev-parse --abbrev-ref HEAD)
COMMIT=$(shell git log --pretty=format:"%s"  | head -n 1)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
BUILD_DIR = build
BUILD_TIME = $(shell date +%s)
BUILD_LCOMMIT =$(shell git log --pretty=format:"%s"  | head -n 1)


.PHONY: install build all test clean exec generate

clean:
	rm -rf ./${BUILD_DIR}/${APP}

test:
	go test -race ./...

buildold: clean test
	go build -o ./${BUILD_DIR}/${APP} ./cmd/main.go

build: clean test
	GO_ENABLED=0 go build \
		-ldflags "-s -w  -X 'version.BranchName=${BRANCH_NAME}' \
		-X 'version.Commit=${COMMIT}' -X 'version.BuildTime=${BUILD_TIME}'" \
       -o ./${BUILD_DIR}/${APP} ./cmd/main.go

install: build
	sudo cp ./${BUILD_DIR}/${APP} /usr/local/bin/${APP} && sudo chmod +x /usr/local/bin/${APP}

exec: build
	exec  ${BUILD_DIR}/${APP} -c=/Users/senya/Documents/go/my/ash/config/ash.yaml 

generate:
	go generate ./...

run: 
	exec  ${BUILD_DIR}/${APP} -c=/Users/senya/Documents/go/my/ash/config/ash.yaml 
