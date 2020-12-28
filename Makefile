ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))

build:
	cd lib/command && cargo build --release
	cp lib/command/target/release/libcommand.dylib lib/
	go build -ldflags="-r $(ROOT_DIR)lib" main.go

run: build
	./main
