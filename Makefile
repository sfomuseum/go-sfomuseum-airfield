compile:
	@make compile-airlines
	@make compile-airports
	@make compile-aircraft

compile-airlines:
	go run -mod vendor cmd/compile-flysfo-airlines-data/main.go
	go run -mod vendor cmd/compile-sfomuseum-airlines-data/main.go

compile-airports:
	go run -mod vendor cmd/compile-sfomuseum-airports-data/main.go

compile-aircraft:
	go run -mod vendor cmd/compile-sfomuseum-aircraft-data/main.go
