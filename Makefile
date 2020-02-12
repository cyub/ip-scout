GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

BINARY_NAME=app
DOCKERCMD=docker
DOCKERBUILD=$(DOCKERCMD) build

all: test build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v
	./$(BINARY_NAME)

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME) -v

docker-build:
	docker run --rm -it -e GOOS=linux -e GOARCH=amd64 -v "$(GOPATH)":/go -w  /go/src/github.com/cyub/ip-scout golang:latest \
	$(GOBUILD)  -o "$(BINARY_NAME)" -v

build-image: build-linux
	$(DOCKERBUILD) --no-cache . -t cyub/ip-scout