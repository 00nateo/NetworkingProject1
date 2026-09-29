# Makefile for simpleclient (CS4254 Fall 2026, Project 1)
# Name of the output executable
BINARY := simpleclient

# Group number used in the submission tarball name (override on the command line).
GROUP ?= 1

# Files that go into the final submission tarball.
SUBMIT_FILES := main.go go.mod Makefile README secret_flags

.PHONY: all build run submit clean

all: build

build:
	@test -f go.mod || go mod init simpleclient
	CGO_ENABLED=0 GOTOOLCHAIN=local go build -o $(BINARY) .

# Convenience target for testing: builds, then runs with $(ARGS).
run: build
	./$(BINARY) $(ARGS)

# Package the source files (not the binary) for submission.
submit:
	tar czf group$(GROUP)_project1.tar.gz $(SUBMIT_FILES)

clean:
	rm -f $(BINARY) group*_project1.tar.gz
