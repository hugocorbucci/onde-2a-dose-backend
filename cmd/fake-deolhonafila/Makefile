NAME=fake-deolhonafila
REV=$(shell git rev-parse --short=10 HEAD)
REGISTRY=docker.io
IMAGE=${REGISTRY}/hugocorbucci/${NAME}
GOPATH:=$(shell echo ${GOPATH})
GOVERSION:=1.16
PROJECT_ROOT=$(shell cd ../.. && pwd)
SERVER_RELATIVE=$(shell pwd | sed -e "s|^${GOPATH}/||g" )
DOCKER_IP:=$(shell if (docker-machine env 2>/dev/null >/dev/null); then (docker-machine ip); else echo "127.0.0.1"; fi)
DB_PORT:=$(shell if (docker-compose ps mysql 2>/dev/null >/dev/null) && (docker-compose port mysql 3306 2>/dev/null >/dev/null); then (docker-compose port mysql 3306 | sed 's\#.*:\([0-9]*\)\#\1\#'); else echo "3306"; fi)
TARGET_URL:=$(shell if [ -z "${TARGET_URL}" ]; then echo "http://localhost:8080"; else echo "${TARGET_URL}"; fi)
.DEFAULT_GOAL := package

target:
	mkdir -p target

clean:
	-rm -Rf target
.PHONY: clean

target/$(NAME)-darwin: target main.go stub.json
	([ "$(shell uname)" = "Darwin" ] && GOOS='darwin' go build -o $@ .) || echo "Can't compile darwin executable"

target/$(NAME)-linux64: target main.go stub.json
	([ "$(shell uname)" = "Darwin" ] && docker run --rm -v "${GOPATH}":/home/guest -w "/home/guest/${SERVER_RELATIVE}" -e "CGO_ENABLED=0" -e "GOPATH=/home/guest" golang:${GOVERSION} go build -o $@ .) || GOOS='linux' go build -o $@ .

build: target/$(NAME)-linux64 target/$(NAME)-darwin
.PHONY: build

package: target/$(NAME)-linux64 Dockerfile
	docker build -t "${IMAGE}:latest" -t "${IMAGE}:${REV}" .
.PHONY: package

push: package
	docker push "${IMAGE}:${REV}"
	docker push "${IMAGE}:latest"
.PHONY: push
