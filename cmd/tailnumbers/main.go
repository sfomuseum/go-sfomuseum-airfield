// tailnumbers is a command line tool to emit the unique set of `swim:tail_number` values from one or more
// whosonfirst/go-whosonfirst-iterate sources for SFO Museum flight data. For example:
//
// 	bin/tailnumbers -iterator-uri /usr/local/data/sfomuseum-data-flights-2021-10/
//
// 	bin/tailnumbers -iterator-uri 'repo://?include=properties.flysfo:date=2021-10-07' /usr/local/data/sfomuseum-data-flights-2021-10/
//
// 	bin/tailnumbers -iterator-uri 'repo://?include=properties.flysfo:date=2021-10-07' /usr/local/data/sfomuseum-data-flights-2021-10/
package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/sfomuseum/go-sfomuseum-airfield/flights"	
)

func main() {

	iterator_uri := flag.String("iterator-uri", "repo://", "...")

	flag.Parse()

	iterator_sources := flag.Args()

	ctx := context.Background()

	tailnumbers, err := flights.TailNumbersFromIterator(ctx, *iterator_uri, iterator_sources...)

	if err != nil {
		log.Fatalf("Failed to derive tailnumbers, %v", err)
	}

	for _, v := range tailnumbers {
		fmt.Println(v)
	}
}
