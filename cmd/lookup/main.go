package main

import (
	_ "github.com/sfomuseum/go-sfomuseum-airfield/aircraft/icao"
	_ "github.com/sfomuseum/go-sfomuseum-airfield/aircraft/sfomuseum"
	_ "github.com/sfomuseum/go-sfomuseum-airfield/airlines/flysfo"
	_ "github.com/sfomuseum/go-sfomuseum-airfield/airlines/sfomuseum"
	_ "github.com/sfomuseum/go-sfomuseum-airfield/airports/sfomuseum"
)

import (
	"context"
	"flag"
	"fmt"
	"github.com/sfomuseum/go-sfomuseum-airfield"
	"log"
)

func main() {

	lookup_uri := flag.String("lookup-uri", "", "...")

	flag.Parse()

	ctx := context.Background()
	lookup, err := airfield.NewLookup(ctx, *lookup_uri)

	if err != nil {
		log.Fatal(err)
	}

	for _, code := range flag.Args() {

		results, err := lookup.Find(ctx, code)

		if err != nil {
			fmt.Printf("%s *** %s\n", code, err)
			continue
		}

		for _, a := range results {
			fmt.Println(a)
		}
	}

}
