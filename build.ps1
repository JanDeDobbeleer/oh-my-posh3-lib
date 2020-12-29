cargo build --release
Copy-Item target/release/posh3.dll ./test
Copy-Item target/release/libposh3.dll.a ./test
