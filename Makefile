APP_NAME = dmx_light_ctrl
VERSION = 1.0.0
HOST ?= 127.0.0.1
PORT ?= 8080
REMOTE ?= pi@raspberrypi.local
GOOS ?= linux
GOARCH ?= arm
GOARM ?= 5
IMAGE_URL ?= "https://downloads.raspberrypi.com/raspios_oldstable_armhf/images/raspios_oldstable_armhf-2024-03-12/2024-03-12-raspios-bullseye-armhf.img.xz"
SD_CARD ?= ""

.PHONY: help
help:
	@echo "make options\n\
		- all             clean, deps, docs, test, vet, fmt, lint & build\n\
		- deps            fetch all dependencies\n\
		- clean           clean build directory bin/\n\
		- build_local     build binary for local testing bin/local_${APP_NAME}\n\
		- build           build binary for raspberry deployment bin/${APP_NAME}\n\
		- test            run test cases\n\
		- run_local       run localy at HOST:PORT ${HOST}:${PORT}\n\
		- run             run binary on ${REMOTE}\n\
		- deploy          build & push latest binary to REMOTE:${REMOTE}\n\
		- help            display this message"

.PHONY: all
all: clean deps test vet fmt lint build

.PHONY: deps
deps:
	go mod tidy -compat 1.17

.PHONY: docs
docs:
	${GOPATH}/bin/go-swagger3 --module-path . --main-file-path ./internal/service/server/server.go --output ./internal/service/server/swagger/swagger_gen.yaml --schema-without-pkg --generate-yaml true

.PHONY: clean
clean:
	go clean
	rm -rf bin

.PHONY: build_local
build_local: docs
	CGO_ENABLED=0 go build -o bin/local_${APP_NAME} cmd/local/main.go

.PHONY: build
build: docs
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} GOARM=${GOARM} go build -o bin/${APP_NAME} cmd/raspberry/main.go

.PHONY: test
test: build
	go test ./...

vet:
	go vet ./...

fmt:
	go list -f '{{.Dir}}' ./... | grep -v /vendor/ | xargs -L1 gofmt -l

lint:
	go list ./... | grep -v /vendor/ | xargs -L1 golint -set_exit_status

.PHONY: run_local
run_local: build_local
	HOST=${HOST} PORT=${PORT} bin/local_${APP_NAME}

.PHONY: deploy
deploy: build
	rsync -Pa bin/${APP_NAME} ${REMOTE}:~/${APP_NAME}

.PHONY: run
run: deploy
	ssh ${REMOTE} ' \
		~/${APP_NAME} \
		'
