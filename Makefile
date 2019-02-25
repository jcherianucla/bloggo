PROTOC_GEN_GO := $(GOPATH)/bin/protoc-gen-go
PROTOC := $(shell which protoc)
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

IDL_DIR := idl/proto/
GEN_DIR := .gen/

.PHONY: proto build test clean
# Generate the proto objects using the protoc-gen-go compiler
proto: $(PROTOC_GEN_GO) $(PROTOC) clean
	mkdir -p $(GEN_DIR)
	protoc -I/usr/local/include -I. -I$(GOPATH)/src \
			--go_out=plugins=grpc:$(GEN_DIR) $(IDL_DIR)*.proto

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
