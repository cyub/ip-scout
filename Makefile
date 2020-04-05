GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

BINARY_NAME=ip-scout
DOCKERCMD=docker
DOCKERBUILD=$(DOCKERCMD) build

all: test build
build: GeoLite2-City.mmdb
	$(GOBUILD) -o $(BINARY_NAME) -v

GeoLite2-City.mmdb:
	wget -c https://static.cyub.vip/GeoLite2-City.mmdb.tgz -O GeoLite2-City.mmdb.tgz
	tar -xzvf GeoLite2-City.mmdb.tgz && rm GeoLite2-City.mmdb.tgz

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

dev:
	air -c .air.conf

run: build
	./$(BINARY_NAME)

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME) -v

docker-build:
	docker run --rm -it -e GOOS=linux -e GOARCH=amd64 -v "$(GOPATH)":/go -w  /go/src/github.com/cyub/ip-scout golang:latest \
	$(GOBUILD)  -o "$(BINARY_NAME)" -v

build-image: build-linux
	$(DOCKERBUILD) --no-cache . -t cyub/ip-scout