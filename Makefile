container_runtime := $(shell which podman || which docker)

$(info using ${container_runtime})

.PHONY: up down swagger swagger-full clean lint proto tools

up: down
	${container_runtime} compose up --build -d

down:
	${container_runtime} compose down

swagger:
	cd gateway && mkdir -p api && $$(go env GOPATH)/bin/swag init -g cmd/main.go -o api --ot json,yaml

swagger-full:
	cd gateway && mkdir -p api && $$(go env GOPATH)/bin/swag init -g cmd/main.go -o api --ot go,json,yaml

clean:
	${container_runtime} compose down -v

lint:
	make -C contracts lint

proto:
	make -C contracts protobuf

tools:
	go install github.com/yoheimuta/protolint/cmd/protolint@latest
	go install github.com/swaggo/swag/cmd/swag@v1.16.6
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $$(go env GOPATH)/bin v2.4.0
	@echo "checking protobuf compiler, if it fails follow guide at https://protobuf.dev/installation/"
	@command -v protoc >/dev/null 2>&1 && echo OK || exit 1
