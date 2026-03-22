MODULE := github.com/Hubcher/project-management/contracts
PROTO_DIR := proto
GEN_DIR := gen/go

PROTO_FILES := $(shell find $(PROTO_DIR) -type f -name "*.proto")

.PHONY: all proto protobuf protolint golint lint clean tools check-dirs

all: lint proto

proto: protobuf

check-dirs:
	@mkdir -p $(GEN_DIR)

protobuf: check-dirs
	protoc -I $(PROTO_DIR) \
		--go_out=$(GEN_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(GEN_DIR) \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)

protolint:
	protolint lint $(PROTO_DIR)

golint:
	golangci-lint run -E gocritic -v ./...

lint: protolint golint

clean:
	rm -rf $(GEN_DIR)

tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/yoheimuta/protolint/cmd/protolint@latest