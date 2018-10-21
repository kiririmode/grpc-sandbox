.DEFAULT_GOAL := help

DOCKER_TAG=latest

## Install dependencies to use this product
deps:
	$(if $(shell which dep), @echo "dep has already installed", go get -u github.com/golang/dep/cmd/dep)
	dep ensure

## Install dependencies to develop this product
devel-deps:
	go get -u github.com/Songmu/make2help/cmd/make2help
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get github.com/fullstorydev/grpcurl
	go install github.com/fullstorydev/grpcurl/cmd/grpcurl
	go get -u golang.org/x/lint/golint
	go get -u github.com/haya14busa/goverage
	go get github.com/ckaznocha/protoc-gen-lint

## Compile .proto to golang sources
pb:
	protoc -I. helloworld.proto --go_out=plugins=grpc:helloworld

## lint
lint:
	protoc -I. helloworld.proto --lint_out=.
	gometalinter ./...

## test
test:
	go test ./...

## get coverage
coverage:
	goverage -v -race -covermode=atomic -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

## gRPC サーバの Docker イメージを作成する。DOCKER_TAG に image の tag 名を設定できる (default: latest)。
docker-build:
	docker build --squash -t kiririmode/grpc-sandbox:$(DOCKER_TAG) .

## show help
help:
	@make2help $(MAKEFILE_LIST)

.PHONY: dev devel-deps pb help
