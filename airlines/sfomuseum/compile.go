package sfomuseum

import (
	"context"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
	"github.com/whosonfirst/go-whosonfirst-iterate/emitter"
	"github.com/whosonfirst/go-whosonfirst-iterate/iterator"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"io"
	_ "log"
	"sync"
)

func CompileAirlinesData(ctx context.Context, iterator_uri string, iterator_sources ...string) ([]*Airline, error) {

	lookup := make([]*Airline, 0)
	mu := new(sync.RWMutex)

	iter_cb := func(ctx context.Context, fh io.ReadSeeker, args ...interface{}) error {

		select {
		case <-ctx.Done():
			return nil
		default:
			// pass
		}

		path, err := emitter.PathForContext(ctx)

		if err != nil {
			return fmt.Errorf("Failed to derive path from context, %w", err)
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

		wof_id, err := properties.Id(body)

		if err != nil {
			return fmt.Errorf("Failed to derive ID for %s, %w", path, err)
		}

		wof_name, err := properties.Name(body)

		if err != nil {
			return fmt.Errorf("Failed to derive name for %s, %w", path, err)
		}

		fl, err := properties.IsCurrent(body)

		if err != nil {
			return fmt.Errorf("Failed to determine is current for %s, %v", path, err)
		}

		sfom_id := int64(-1)

		sfom_rsp := gjson.GetBytes(body, "properties.sfomuseum:airline_id")

		if sfom_rsp.Exists() {
			sfom_id = sfom_rsp.Int()
		}

		role_rsp := gjson.GetBytes(body, "properties.sfomuseum:airline_role")

		a := &Airline{
			WhosOnFirstId: wof_id,
			SFOMuseumId:   sfom_id,
			Name:          wof_name,
			Role:          role_rsp.String(),
			IsCurrent:     fl.Flag(),
		}

		concordances := properties.Concordances(body)

		if concordances != nil {

			iata_code, ok := concordances["iata:code"]

			if ok {
				a.IATACode = iata_code
			}

			icao_code, ok := concordances["icao:code"]

			if ok {
				a.ICAOCode = icao_code
			}

			callsign, ok := concordances["icao:callsign"]

			if ok {
				a.ICAOCallsign = callsign
			}

			id, ok := concordances["wd:id"]

			if ok {
				a.WikidataId = id
			}
		}

		mu.Lock()
		lookup = append(lookup, a)
		mu.Unlock()

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
