package flysfo

import (
	"context"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/iterator"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"io"
	"log"
	"sync"
)

func CompileAirlinesData(ctx context.Context, iterator_uri string, iterator_sources ...string) ([]Airline, error) {

	lookup := make([]Airline, 0)
	mu := new(sync.RWMutex)

	seen := new(sync.Map)

	iter_cb := func(ctx context.Context, path string, fh io.ReadSeeker, args ...interface{}) error {

		select {
		case <-ctx.Done():
			return nil
		default:
			// pass
		}

		_, uri_args, err := uri.ParseURI(path)

		if err != nil {
			return fmt.Errorf("Failed to parse %s, %w", path, err)
		}

		if uri_args.IsAlternate {
			return nil
		}

		body, err := io.ReadAll(fh)

		if err != nil {
			return fmt.Errorf("Failed to read %s, %w", path, err)
		}

		pt_rsp := gjson.GetBytes(body, "properties.sfomuseum:placetype")

		if pt_rsp.String() != "airline" {
			return nil
		}

		concordances := properties.Concordances(body)

		if concordances == nil {
			return nil
		}

		iata_code, ok := concordances["flysfo:code"]

		if !ok {
			return nil
		}

		fl, err := properties.IsCurrent(body)

		if err != nil {
			return fmt.Errorf("Failed to determine whether %s is current, %v", path, err)
		}

		if !fl.IsTrue() || !fl.IsKnown() {
			return nil
		}

		wof_id, err := properties.Id(body)

		if err != nil {
			return fmt.Errorf("Failed to derive ID, %w", err)
		}

		wof_name, err := properties.Name(body)

		if err != nil {
			return fmt.Errorf("Failed to derive name, %w", err)
		}

		a := Airline{
			WhosOnFirstId: wof_id,
			FlysfoID:      fmt.Sprintf("%s", iata_code),
			Name:          wof_name,
			IATACode:      fmt.Sprintf("%s", iata_code),
		}

		icao_code, ok := concordances["icao:code"]

		if ok {
			w, ok := seen.Load(icao_code)

			if ok {
				log.Println("WARNING", icao_code, w, wof_id)

			} else {
				seen.Store(icao_code, wof_id)
			}

			a.ICAOCode = fmt.Sprintf("%s", icao_code)

		} else {
			log.Println("WARNING", "Missing ICAO code", wof_id)
		}

		mu.Lock()
		defer mu.Unlock()

		lookup = append(lookup, a)

		return nil
	}

	iter, err := iterator.NewIterator(ctx, iterator_uri, iter_cb)

	if err != nil {
		return nil, fmt.Errorf("Failed to create iterator, %w", err)
	}

	err = iter.IterateURIs(ctx, iterator_sources...)

	if err != nil {
		return nil, fmt.Errorf("Failed to iterate sources, %w", err)
	}

	return lookup, nil
}
