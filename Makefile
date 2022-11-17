SHELL=/bin/sh
IMAGE_TAG := $(shell git rev-parse HEAD)

GO_BUILD_VERSION := $(or $(GO_BUILD_VERSION), 'local')
GO_BUILD_BRANCH := $(or $(GO_BUILD_BRANCH), $(shell git branch | grep \* | cut -d ' ' -f2))
GO_BUILD_COMMIT := $(or $(GO_BUILD_COMMIT), $(shell git rev-parse HEAD))
GO_BUILD_TIME := $(or $(GO_BUILD_TIME), $(shell date '+%m-%d-%YT%H:%M:%S'))
GOOS := $(or $(GOOS), '')

PROTO_FILES += $(shell find proto -name "*.proto" -not -path '*google*' | sort -u)

.PHONY: test
test:
	go test -cover -race ./... -count=1

.PHONY: deps
deps:
	go mod download
	@$(MAKE) tidy

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=${GOOS} go build -o bin/app -ldflags="-X main.Version=${GO_BUILD_VERSION} -X main.Branch=${GO_BUILD_BRANCH} -X main.Commit=${GO_BUILD_COMMIT} -X main.BuildTime=${GO_BUILD_TIME}" cmd/*.go

.PHONY: lint
lint:
	golangci-lint run

.PHONY: ci-lint
ci-lint: tools lint tidy ## run linter and clean dependencies after it

.PHONY: run
run:
	go run cmd/*.go serve

.PHONY: proto
proto:
	rm -rf pkg/api/proto
	mkdir -p pkg/proto
	protoc -I$(shell pwd) --go_out=pkg/proto --go-grpc_out=pkg/proto $(PROTO_FILES)

.PHONY: docs
docs:
	protoc --doc_out=docs --doc_opt=html,index.html $(PROTO_FILES)
	protoc --doc_out=docs --doc_opt=markdown,api.md $(PROTO_FILES)

.PHONY: tools
tools:
	go install \
        github.com/golangci/golangci-lint/cmd/golangci-lint \
        github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc \
        google.golang.org/grpc \
        google.golang.org/grpc/cmd/protoc-gen-go-grpc \
        google.golang.org/protobuf/cmd/protoc-gen-go

.PHONY: gen
gen:
	go generate ./...
