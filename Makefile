# Copyright (C) 2019 The go-postgresql Authors. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Need bash for pushd, popd
SHELL := bash

PREFIX?=$(shell pwd)

GITHUB_ROOT=
PACKAGE_NAME=go-postgresql
SOURCE_DIR=postgresql

PACKAGE_ID=${GITHUB_ROOT}${PACKAGE_NAME}/${SOURCE_DIR}
PACKAGES=\
	${PACKAGE_ID} \
	${PACKAGE_ID}/protocol

.PHONY: version clean

all: test

VERSION_GO=${SOURCE_DIR}/version.go

${VERSION_GO}: ${SOURCE_DIR}/version.gen
	$< > $@

version: ${VERSION_GO}

format:
	gofmt -w ${SOURCE_DIR}

vet: format
	go vet ${PACKAGES}

build: vet
	go build -v ${PACKAGES}

test: vet
	go test -v -cover -timeout 300s ${PACKAGES}

clean:
	go clean -i ${PACKAGES}
