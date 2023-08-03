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

GIT_ROOT=github.com/cybergarage
PRODUCT_NAME=go-postgresql
MODULE_ROOT=${GIT_ROOT}/${PRODUCT_NAME}

PKG_NAME=postgresql
PKG_COVER=${PKG_NAME}-cover
PKG_SRC_ROOT=${PKG_NAME}
PKG=${MODULE_ROOT}/${PKG_SRC_ROOT}

TEST_SRC_ROOT=${PKG_NAME}test
TEST_PKG=${MODULE_ROOT}/${TEST_SRC_ROOT}

EXAMPLES_ROOT=examples
EXAMPLES_PKG_ROOT=${GIT_ROOT}/${PRODUCT_NAME}/${EXAMPLES_ROOT}
EXAMPLES_SRC_DIR=${EXAMPLES_ROOT}
EXAMPLE_BINARIES=\
	${EXAMPLES_PKG_ROOT}/go-postgresqld

BIN_ROOT=bin
BIN_PKG_ROOT=${GIT_ROOT}/${PRODUCT_NAME}/${BIN_ROOT}
BIN_SRC_DIR=${BIN_ROOT}
BIN_BINARIES=\
	${BIN_PKG_ROOT}/pgpcapdump

BINARIES=\
	${EXAMPLE_BINARIES} \
	${BIN_BINARIES}

.PHONY: clean test

all: test

format:
	gofmt -s -w ${PKG_SRC_ROOT} ${TEST_SRC_ROOT} ${EXAMPLES_SRC_DIR} ${BIN_SRC_DIR}

vet: format
	go vet ${PKG}

lint: vet
	golangci-lint run ${PKG_SRC_ROOT}/... ${TEST_SRC_ROOT}/... ${EXAMPLES_SRC_DIR}/... ${BIN_SRC_DIR}/...

build: vet
	go build -v ${BINARIES}

test: lint
	go test -v -p 1 -timeout 60s -cover -coverpkg=${PKG}/... -coverprofile=${PKG_COVER}.out ${PKG}/... ${TEST_PKG}/...
	go tool cover -html=${PKG_COVER}.out -o ${PKG_COVER}.html

test_only:
	go test -v -p 1 -timeout 60s -cover -coverpkg=${PKG} -coverprofile=${PKG_COVER}.out ${PKG}/... ${TEST_PKG}/...
	go tool cover -html=${PKG_COVER}.out -o ${PKG_COVER}.html

install:
	go install -v -gcflags=${GCFLAGS} ${BINARIES}

clean:
	go clean -i ${PKG}

watchtest:
	fswatch -o . -e ".*" -i "\\.go$$" | xargs -n1 -I{} make test

watchlint:
	fswatch -o . -e ".*" -i "\\.go$$" | xargs -n1 -I{} make lint