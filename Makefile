build:
	cargo build --release

macos: build
	cp target/release/libposh3.dylib ./test

unix: build
	cp target/release/libposh3.so ./test

