PROTOC_GEN_GO := $(GOPATH)/bin/protoc-gen-go
PROTOC := $(shell which protoc)
# A
# If protoc isn't on the path, set it to a target that's never up to date, so
# the install command always runs.
ifeq ($(PROTOC),)
    PROTOC = must-rebuild
endif

# Figure out which machine we're running on.
UNAME := $(shell uname)

$(PROTOC):
# Run the right installation command for the operating system.
ifeq ($(UNAME), Darwin)
	brew install protobuf
endif
ifeq ($(UNAME), Linux)
	sudo apt-get install protobuf-compiler
endif

# If $GOPATH/bin/protoc-gen-go does not exist, we'll run this command to install
# it.
$(PROTOC_GEN_GO):
	go get -u github.com/golang/protobuf/protoc-gen-go

OUTPUT_PATH := $(GOPATH)/src
IDL_DIR := idl/proto/
DATA_IDL_DIR := $(IDL_DIR)data/
MODELS_IDL_DIR := $(IDL_DIR)models/
PROTOS := $(DATA_IDL_DIR) $(MODELS_IDL_DIR) $(IDL_DIR)

.PHONY: proto build test clean

# Run the proto compiler on input
define compile_proto
	protoc -I. --proto_path=$(IDL_DIR) \
		--go_out=plugins=grpc:$(OUTPUT_PATH) $(1)*.proto
endef

# Generate the proto objects using the protoc-gen-go compiler
proto: $(PROTOC_GEN_GO) $(PROTOC)
	# TODO: Figure out how to successfully do this as a for loop
	$(call compile_proto,$(DATA_IDL_DIR))
	$(call compile_proto,$(MODELS_IDL_DIR))
	$(call compile_proto,$(IDL_DIR))

# Location of prometheus monitoring configuration
PROMETHEUS_CFG := "config/prometheus.yml"
# TODO: Run test before anything
build: clean test
	go build
	# TODO: Make sure docker container has prometheus
	prometheus --config.file=$(PROMETHEUS_CFG) &
	# TODO: Docker compose postgresql server
	./bloggo

clean:
	rm ./bloggo
	rm -rf $(GEN_DIR)

PKGS := $(shell go list ./ ... | grep -v vendor/)
# Generate proto objects and lint
test: proto $(GOMETALINTER)
	go test $(PKGS)

BIN_DIR := $(GOPATH)/bin
GOMETALINTER := $(BIN_DIR)/gometalinter

# Run go meta linter
$(GOMETALINTER):
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install &> /dev/null
