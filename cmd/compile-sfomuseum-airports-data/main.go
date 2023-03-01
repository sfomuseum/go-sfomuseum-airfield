package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/sfomuseum/go-sfomuseum-airfield/airports/sfomuseum"
	_ "github.com/whosonfirst/go-whosonfirst-iterate-git/v2"		
)

func main() {

	default_target := fmt.Sprintf("data/%s", sfomuseum.DATA_JSON)

	iterator_uri := flag.String("iterator-uri", "repo://?include=properties.sfomuseum:placetype=airport&exclude=properties.edtf:deprecated=.*", "...")
	iterator_source := flag.String("iterator-source", "/usr/local/data/sfomuseum-data-whosonfirst", "...")

	target := flag.String("target", default_target, "The path to write SFO Museum airport data.")
	stdout := flag.Bool("stdout", false, "Emit SFO Museum airport data to SDOUT.")

	flag.Parse()

	ctx := context.Background()

	writers := make([]io.Writer, 0)

	fh, err := os.OpenFile(*target, os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		log.Fatalf("Failed to open '%s', %v", *target, err)
	}

	writers = append(writers, fh)

	if *stdout {
		writers = append(writers, os.Stdout)
	}

	wr := io.MultiWriter(writers...)

	lookup, err := sfomuseum.CompileAirportsData(ctx, *iterator_uri, *iterator_source)

	if err != nil {
		log.Fatalf("Failed to compile data, %v", err)
	}

	enc := json.NewEncoder(wr)
	err = enc.Encode(lookup)

	if err != nil {
		log.Fatalf("Failed to marshal results, %v", err)
	}

}
