# Copyright (C) 2019 The go-postgresql Authors. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http:#www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

SHELL := bash

PREFIX?=$(shell pwd)

GOBIN := $(shell go env GOPATH)/bin
PATH := $(GOBIN):$(PATH)

GIT_ROOT=github.com/cybergarage
PRODUCT_NAME=go-postgresql
MODULE_ROOT=${GIT_ROOT}/${PRODUCT_NAME}

PKG_NAME=postgresql
PKG_VER=$(shell git describe --abbrev=0 --tags)
PKG_COVER=${PKG_NAME}-cover
PKG_SRC_ROOT=${PKG_NAME}
PKG=${MODULE_ROOT}/${PKG_SRC_ROOT}

TEST_SRC_ROOT=${PKG_NAME}test
TEST_PKG=${MODULE_ROOT}/${TEST_SRC_ROOT}

EXAMPLES_ROOT=examples
EXAMPLES_SRC_ROOT=${EXAMPLES_ROOT}
EXAMPLES_DEAMON_BIN=go-postgresqld
EXAMPLES_PKG_ROOT=${GIT_ROOT}/${PRODUCT_NAME}/${EXAMPLES_ROOT}
EXAMPLE_BINARIES=\
	${EXAMPLES_PKG_ROOT}/${EXAMPLES_DEAMON_BIN}

EXAMPLES_DOCKER_TAG=cybergarage/${EXAMPLES_DEAMON_BIN}:${PKG_VER}
EXAMPLES_DOCKER_TAG_LATEST=cybergarage/${EXAMPLES_DEAMON_BIN}:latest

BIN_ROOT=bin
BIN_PKG_ROOT=${GIT_ROOT}/${PRODUCT_NAME}/${BIN_ROOT}
BIN_SRC_DIR=${BIN_ROOT}
BIN_BINARIES=\
	${BIN_PKG_ROOT}/pgpcapdump

BINARIES=\
	${EXAMPLE_BINARIES}

.PHONY: clean test version sysbench certs
.IGNORE: lint

all: test

version:
	@pushd ${PKG_SRC_ROOT} && ./version.gen > version.go && popd
	-git commit ${PKG_SRC_ROOT}/version.go -m "Update version"

format: version
	gofmt -s -w ${PKG_SRC_ROOT} ${TEST_SRC_ROOT} ${EXAMPLES_SRC_ROOT}

vet: format
	go vet ${PKG}

lint: vet
	golangci-lint run ${PKG_SRC_ROOT}/... ${TEST_SRC_ROOT}/... ${EXAMPLES_SRC_ROOT}/...

test: lint
	chmod og-rwx  ${TEST_SRC_ROOT}/certs/client-key.pem
	go test -v -p 1 -timeout 10m -cover -coverpkg=${PKG}/... -coverprofile=${PKG_COVER}.out ${PKG}/... ${TEST_PKG}/...
	go tool cover -html=${PKG_COVER}.out -o ${PKG_COVER}.html

test_only:
	chmod og-rwx  ${TEST_SRC_ROOT}/certs/client-key.pem
	go test -v -p 1 -timeout 10m -cover -coverpkg=${PKG}/... -coverprofile=${PKG_COVER}.out ${PKG}/... ${TEST_PKG}/...
	go tool cover -html=${PKG_COVER}.out -o ${PKG_COVER}.html

certs:
	@pushd ${TEST_SRC_ROOT}/certs && make && popd

sysbench:
	go test -v -p 1 -run ^TestSysbench ${TEST_PKG}/sysbench

build: vet
	go build -v ${BINARIES}

install:
	go install -v -gcflags=${GCFLAGS} ${BINARIES}

run: install
	${GOBIN}/${EXAMPLES_DEAMON_BIN} --debug

image: test
	docker image build -t${EXAMPLES_DOCKER_TAG} -t${EXAMPLES_DOCKER_TAG_LATEST} .
	docker push ${EXAMPLES_DOCKER_TAG_LATEST}

image-push: image
	docker push ${EXAMPLES_DOCKER_TAG}

rund: image
	docker container run -it --rm -p 5432:5432 ${EXAMPLES_DOCKER_TAG_LATEST}

clean:
	go clean -i ${PKG}

watchtest:
	fswatch -o . -e ".*" -i "\\.go$$" | xargs -n1 -I{} make test

watchlint:
	fswatch -o . -e ".*" -i "\\.go$$" | xargs -n1 -I{} make lint
