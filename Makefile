.DEFAULT_GOAL := help

## Install dependencies to use this product
dev:
	$(if $(shell which dep), @echo "dep has already installed", go get -u github.com/golang/dep/cmd/dep)
	dep ensure

## Install dependencies to develop this product
devel-deps:
	brew install protobuf
	go get -u github.com/Songmu/make2help/cmd/make2help
	go get -u github.com/golang/protobuf/protoc-gen-go

## Compile .proto to golang sources
pb:
	protoc -I. helloworld.proto --go_out=plugins=grpc:helloworld

## show help
help:
	@make2help $(MAKEFILE_LIST)

.PHONY: dev devel-deps help
