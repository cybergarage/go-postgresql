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

# Usage Guide BenchBase
# https://github.com/cmu-db/benchbase?tab=readme-ov-file#usage-guide

# Variables
#BENCHBASE_GITURL := https://github.com/cmu-db
BENCHBASE_GITURL := https://github.com/cybergarage

BENCHBASE_DIR := benchbase
BENCHBASE_TARGET := postgres
BENCHBASE_JAR := benchbase.jar
BENCHBASE_CONFIG := tpcc_config.xml
BENCHBASE_TARGET_DIR := $(BENCHBASE_DIR)/target/benchbase-$(BENCHBASE_TARGET)
BENCHBASE_TARGET_JAR := $(BENCHBASE_TARGET_DIR)/$(BENCHBASE_JAR)

# Default target
.PHONY: all setup clean help

all: setup

# Setup BenchBase
setup: $(BENCHBASE_JAR) $(BENCHBASE_CONFIG)
	@echo "BenchBase setup completed."
	java -jar $(BENCHBASE_JAR) -h

# Download and build BenchBase if JAR doesn't exist
$(BENCHBASE_JAR):
	@echo "Setting up BenchBase..."
	rm -rf $(BENCHBASE_DIR) lib data config scripts
	git clone --depth 1 $(BENCHBASE_GITURL)/benchbase.git
	cd $(BENCHBASE_DIR) && ./mvnw package -P $(BENCHBASE_TARGET)
	cd $(BENCHBASE_DIR)/target && tar xvzf benchbase-$(BENCHBASE_TARGET).tgz
	cp $(BENCHBASE_TARGET_JAR) .
	mv $(BENCHBASE_TARGET_DIR)/lib .
	mv $(BENCHBASE_TARGET_DIR)/data .
	mv $(BENCHBASE_TARGET_DIR)/config .
	mv $(BENCHBASE_TARGET_DIR)/scripts .

$(BENCHBASE_CONFIG):
	@echo "Creating default BenchBase configuration..."
	cp config/postgres/sample_tpcc_config.xml tpcc_config.xml;

# Clean up generated files
clean:
	rm -rf $(BENCHBASE_DIR) lib data config scripts $(BENCHBASE_JAR)
	@echo "Cleaned up BenchBase files."

# Run BenchBase with specified configuration
run: $(BENCHBASE_JAR) $(BENCHBASE_CONFIG)
	@echo "Running BenchBase with configuration $(BENCHBASE_CONFIG)..."
	java -jar $(BENCHBASE_JAR) -b tpcc -c $(BENCHBASE_CONFIG) --create=true --load=true --execute=true

# Show help
help:
	@echo "BenchBase Makefile"
	@echo ""
	@echo "Targets:"
	@echo "  all     - Setup BenchBase (default)"
	@echo "  setup   - Setup BenchBase and show help"
	@echo "  clean   - Remove all generated files"
	@echo "  help    - Show this help message"
	@echo ""
	@echo "Variables:"
	@echo "  BENCHBASE_GITURL=$(BENCHBASE_GITURL)"
	@echo "  BENCHBASE_TARGET=$(BENCHBASE_TARGET)"