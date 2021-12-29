cli:
	@make cli-lookup
	# @make cli-create
	@make cli-stats

cli-create:
	go build -mod vendor -o bin/import-airport cmd/import-airport/main.go

cli-lookup:
	go build -mod vendor -o bin/lookup cmd/lookup/main.go

cli-stats:
	go build -mod vendor -o bin/tailnumbers cmd/tailnumbers/main.go

compile:
	@make compile-airlines
	@make compile-airports
	@make compile-aircraft
	@make cli-lookup

compile-airlines:
	go run -mod vendor cmd/compile-flysfo-airlines-data/main.go -iterator-uri 'git:///tmp?exclude=properties.edtf:deprecated=.*' https://github.com/sfomuseum-data/sfomuseum-data-enterprise.git
	go run -mod vendor cmd/compile-sfomuseum-airlines-data/main.go  -iterator-uri 'git:///tmp?exclude=properties.edtf:deprecated=.*' https://github.com/sfomuseum-data/sfomuseum-data-enterprise.git

compile-airports:
	go run -mod vendor cmd/compile-sfomuseum-airports-data/main.go  -iterator-uri 'git:///tmp?include=properties.sfomuseum:placetype=airport&exclude=properties.edtf:deprecated=.*' https://github.com/sfomuseum-data/sfomuseum-data-whosonfirst.git

compile-aircraft:
	go run -mod vendor cmd/compile-sfomuseum-aircraft-data/main.go -iterator-uri 'git:///tmp?exclude=properties.edtf:deprecated=.*' https://github.com/sfomuseum-data/sfomuseum-data-aircraft.git
