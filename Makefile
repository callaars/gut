# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=gut

all: build
build: 
	$(GOBUILD) -o $(BINARY_NAME) -v
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
test:
	$(GOTEST)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)
install:
	cp $(BINARY_NAME) /usr/bin
deps:
	$(GOGET) github.com/fatih/color
	$(GOGET) github.com/phayes/permbits
	$(GOGET) github.com/urfave/cli

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME) -v
