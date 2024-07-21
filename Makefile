# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Main package file
MAIN_FILE=main.go

# Build output directory
BUILDDIR=build

# Binary output name
BINARY_NAME=myapp

all: clean build

build: 
	$(GOBUILD) -o $(BUILDDIR)/$(BINARY_NAME) -v $(MAIN_FILE)

clean: 
	$(GOCLEAN)
	rm -rf $(BUILDDIR)

run:
	$(GOBUILD) -o $(BUILDDIR)/$(BINARY_NAME) -v $(MAIN_FILE)
	./$(BUILDDIR)/$(BINARY_NAME)

test: 
	$(GOTEST) -v ./...

deps:
	$(GOGET) github.com/go-sql-driver/mysql

.PHONY: all build clean run test deps
