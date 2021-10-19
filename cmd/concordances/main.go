package main

import (
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/iterator"
	"flag"
	"log"
	"io"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
	sfom_properties "github.com/sfomuseum/go-sfomuseum-feature/properties"	
	"context"
	"fmt"
	"os"
	"github.com/sfomuseum/go-csvdict"
	"strconv"
	"sync"
	"sort"
)

func main() {

	iterator_uri := flag.String("iterator-uri", "repo://", "")
	
	flag.Parse()

	iterator_sources := flag.Args()
	
	ctx := context.Background()

	writers := []io.Writer{
		os.Stdout,
	}

	mw := io.MultiWriter(writers...)

	mu := new(sync.RWMutex)
	
	var csv_wr *csvdict.Writer
	
	iter_cb := func(ctx context.Context, path string, fh io.ReadSeeker, args ...interface{}) error {

		body, err := io.ReadAll(fh)

		if err != nil {
			return fmt.Errorf("Failed to read %s, %v", path, err)
		}

		id, err := properties.Id(body)

		if err != nil {
			return fmt.Errorf("Failed to derive ID for %s, %v", path, err)
		}

		name, err := properties.Name(body)
		
		if err != nil {
			return fmt.Errorf("Failed to derive name for %s, %v", path, err)
		}

		pt, err := sfom_properties.Placetype(body)

		if err != nil {
			return fmt.Errorf("Failed to derive placetype for %s, %v", path, err)
		}

		sfom_id, _ := sfom_properties.Id(body)

		concordances := properties.Concordances(body)

		wd_id, ok := concordances["wd:id"]

		if !ok {
			return nil
		}

		mu.Lock()
		defer mu.Unlock()
		
		out := map[string]string{
			"placetype": pt,
			"sfom_id": strconv.FormatInt(sfom_id, 10),
			"wof_id": strconv.FormatInt(id, 10),
			"name": name,
			"wikidata_id": wd_id.(string),
		}

		if csv_wr == nil {

			fieldnames := make([]string, 0)

			for k, _ := range out {
				fieldnames = append(fieldnames, k)
			}

			sort.Strings(fieldnames)
			
			wr, err := csvdict.NewWriter(mw, fieldnames)

			if err != nil {
				return fmt.Errorf("Failed to create new CSV writer, %w", err)
			}

			wr.WriteHeader()
			csv_wr = wr

		}

		csv_wr.WriteRow(out)
		csv_wr.Flush()

		// log.Println(pt, id, sfom_id, name, wd_id)
		return nil
	}

	iter, err := iterator.NewIterator(ctx, *iterator_uri, iter_cb)

	if err != nil {
		log.Fatalf("Failed to create new iterator, %v", err)
	}

	err = iter.IterateURIs(ctx, iterator_sources...)

	if err != nil {
		log.Fatalf("Failed to iterate URIs, %v", err)
	}

	csv_wr.Flush()
}
