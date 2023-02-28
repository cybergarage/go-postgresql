# Copyright (C) 2019 Satoshi Konno. All rights reserved.
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

GIT_ROOT=github.com/cybergarage/
PRODUCT_NAME=go-postgresql
PKG_NAME=postgresql

MODULE_ROOT=${PKG_NAME}
MODULE_PKG_ROOT=${GIT_ROOT}${PRODUCT_NAME}/${MODULE_ROOT}
MODULE_SRC_DIR=${PKG_NAME}
MODULE_SRCS=\
	${MODULE_SRC_DIR} \
	${MODULE_SRC_DIR}/protocol
MODULE_PKGS=\
	${MODULE_PKG_ROOT} \
	${MODULE_PKG_ROOT}/protocol

EXAMPLES_ROOT=examples
EXAMPLES_DEAMON_BIN=go-postgresqld
EXAMPLES_PKG_ROOT=${GIT_ROOT}${PRODUCT_NAME}/${EXAMPLES_ROOT}
EXAMPLES_SRC_DIR=${EXAMPLES_ROOT}
EXAMPLES_SRCS=\
	${EXAMPLES_SRC_DIR}/${EXAMPLES_DEAMON_BIN} \
	${EXAMPLES_SRC_DIR}/${EXAMPLES_DEAMON_BIN}/server \
	${EXAMPLES_SRC_DIR}/${EXAMPLES_DEAMON_BIN}/server/storage
EXAMPLES_DEAMON_ROOT=${EXAMPLES_PKG_ROOT}
EXAMPLES_PKGS=\
	${EXAMPLES_DEAMON_ROOT}/${EXAMPLES_DEAMON_BIN}/server \
	${EXAMPLES_DEAMON_ROOT}/${EXAMPLES_DEAMON_BIN}/server/storage
EXAMPLE_BINARIES=\
	${EXAMPLES_DEAMON_ROOT}/${EXAMPLES_DEAMON_BIN}

TEST_ROOT=${PKG_NAME}test
TEST_PKG_ROOT=${GIT_ROOT}${PRODUCT_NAME}/${TEST_ROOT}
TEST_SRC_DIR=${TEST_ROOT}
TEST_SRCS=\
	${TEST_SRC_DIR}/client \
	${TEST_SRC_DIR}/server \
	${TEST_SRC_DIR}/ycsb \
	${TEST_SRC_DIR}/sqltest
TEST_PKGS=\
	${TEST_PKG_ROOT}/client \
	${TEST_PKG_ROOT}/server \
	${TEST_PKG_ROOT}/ycsb \
	${TEST_PKG_ROOT}/sqltest

ALL_ROOTS=\
	${MODULE_ROOT} \
	${EXAMPLES_ROOT} \
	${TEST_ROOT}

ALL_SRCS=\
	${MODULE_SRCS} \
	${EXAMPLES_SRCS} \
	${TEST_SRCS}

ALL_PKGS=\
	${MODULE_PKGS} \
	${EXAMPLES_PKGS} \
	${TEST_PKGS}

BINARIES=${EXAMPLE_BINARIES}

.PHONY: clean test

all: test

format:
	gofmt -s -w ${ALL_ROOTS}

vet: format
	go vet ${ALL_PKGS}

lint: vet
	golangci-lint run ${ALL_SRCS}

build: vet
	go build -v ${MODULE_PKGS}

test: lint
	go test -v -cover -p=1 ${ALL_PKGS}

install: build
	go install -v -gcflags=${GCFLAGS} ${BINARIES}

clean:
	go clean -i ${ALL_PKGS}
