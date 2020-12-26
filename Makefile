#!/usr/bin/make -f

.PHONY: all help build clean distclean build-demos

all: help

help:
	@echo "usage: make {help|test|clean|demos}"
	@echo
	@echo "  test: perform all available tests"
	@echo "  clean: cleans package  and built files"
	@echo "  demos: builds the boxes, mouse and unicode demos"
	@echo

test:
	@echo -n "vetting cdk ..."
	@go vet && echo " done"
	@echo "testing cdk ..."
	@go test -cover -coverprofile=coverage.out ./...
	@echo "test coverage ..."
	@go tool cover -html=coverage.out

clean:
	@echo "cleaning"
	@go clean ./...      || true
	@rm -fv beep         || true
	@rm -fv boxes        || true
	@rm -fv colors       || true
	@rm -fv hello_world  || true
	@rm -fv mouse        || true
	@rm -fv unicode      || true
	@rm -fv coverage.out || true

demos: clean build
	@echo "building beep"
	@go build -v _demos/beep.go
	@echo "building boxes"
	@go build -v _demos/boxes.go
	@echo "building colors"
	@go build -v _demos/colors.go
	@echo "building hello_world"
	@go build -v _demos/hello_world.go
	@echo "building mouse"
	@go build -v _demos/mouse.go
	@echo "building unicode"
	@go build -v _demos/unicode.go
