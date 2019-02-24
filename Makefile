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

.PHONY: proto
proto: $(PROTOC_GEN_GO) $(PROTOC)
	rm -rf $(GEN_DIR)
	mkdir -p $(GEN_DIR)
	protoc -I/usr/local/include -I. -I$(GOPATH)/src \
			--go_out=plugins=grpc:$(GEN_DIR) $(IDL_DIR)*.proto

PKGS := $(shell go list ./ ... | grep -v vendor/)

.PHONY: test
test: lint
	go test $(PKGS)

BIN_DIR := $(GOPATH)/bin
GOMETALINTER := $(BIN_DIR)/gometalinter

$(GOMETALINTER):
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install &> /dev/null
