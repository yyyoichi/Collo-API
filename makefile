run:
	@export CGO_LDFLAGS="-L/usr/local/lib -lmecab -lstdc++" && \
	export CGO_CFLAGS="-I/path/to/include" && \
	go run .