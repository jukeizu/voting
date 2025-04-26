TAG=$(shell git describe --tags --always)
VERSION=$(TAG:v%=%)
REPO=jukeizu/voting
GO=go
BUILD=GOARCH=amd64 $(GO) build -ldflags="-s -w -X main.Version=$(VERSION)" 
PROTOFILES=$(wildcard .protobuf/voting/*/*.proto)
PBFILES=$(patsubst %.proto,%.pb.go, $(PROTOFILES))
PROTOPBDEST="../../../api/protobuf-spec/$(patsubst %.proto,%pb, $(notdir $<))"

.PHONY: all deps test build build-linux run docker-build docker-save docker-deploy proto clean $(PROTOFILES)

all: deps test build 
deps:
	$(GO) mod download

test:
	$(GO) vet ./...
	$(GO) test -v -race ./...

build:
	$(BUILD) -o bin/voting-$(VERSION) ./cmd/...

build-linux:
	CGO_ENABLED=0 GOOS=linux $(BUILD) -a -installsuffix cgo -o bin/voting ./cmd/...

run: build
	./bin/voting-$(VERSION) -migrate

docker-build:
	docker build -t $(REPO):$(VERSION) .

docker-save:
	mkdir -p bin && docker save -o bin/image.tar $(REPO):$(VERSION)

docker-push:
	docker push $(REPO):$(VERSION)

proto: $(PBFILES)

%.pb.go: %.proto
	cd $(dir $<) && mkdir -p $(PROTOPBDEST) && protoc $(notdir $<) --go_out=$(PROTOPBDEST) --go-grpc_out=$(PROTOPBDEST) 

clean:
	@rm -f bin/*
