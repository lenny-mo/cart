
GOPATH:=$(shell go env GOPATH)
MODIFY= proto/

.PHONY: proto
proto:
    
	protoc  --micro_out=${MODIFY} --go_out=${MODIFY} proto/cart/cart.proto

.PHONY: build
build: proto

	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cart-service *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t cart-service:latest
