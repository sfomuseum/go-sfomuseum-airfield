compile-airlines:
	go run -mod vendor cmd/compile-flysfo-airlines-data/main.go
	go run -mod vendor cmd/compile-sfomuseum-airlines-data/main.go
