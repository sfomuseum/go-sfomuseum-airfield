cli:
	@make cli-lookup

cli-lookup:
	go build -mod vendor -o bin/lookup cmd/lookup/main.go

compile:
	@make compile-airlines
	@make compile-airports
	@make compile-aircraft
	@make cli-lookup

compile-airlines:
	go run -mod vendor cmd/compile-flysfo-airlines-data/main.go
	go run -mod vendor cmd/compile-sfomuseum-airlines-data/main.go

compile-airports:
	go run -mod vendor cmd/compile-sfomuseum-airports-data/main.go

compile-aircraft:
	go run -mod vendor cmd/compile-sfomuseum-aircraft-data/main.go
