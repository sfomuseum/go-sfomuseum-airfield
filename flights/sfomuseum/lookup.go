package sfomuseum

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sfomuseum/go-sfomuseum-airfield/flights"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-iterate/emitter"
	"github.com/whosonfirst/go-whosonfirst-iterate/iterator"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"io"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

var lookup_table *sync.Map
var lookup_idx int64

var lookup_init sync.Once
var lookup_init_err error

type FlightsLookupFunc func(context.Context)

// FlightsLookup provides an implementation of the `Lookup` interface for locating existing (WOF) flight records.
type FlightsLookup struct {
	flights.FlightsLookup
}

func init() {
	ctx := context.Background()
	flights.RegisterFlightsLookup(ctx, "sfomuseum", NewSFOMuseumLookup)
	lookup_idx = int64(0)
}

func NewSFOMuseumLookup(ctx context.Context, uri string) (flights.FlightsLookup, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	// Account for both:
	// airfield.NewLookup(ctx, "flights://sfomuseum/iterator")
	// airlines.NewAirLinesLookup(ctx, "sfomuseum://iterator")

	var source string

	switch u.Host {
	case "sfomuseum":
		source = u.Path
	default:
		source = u.Host
	}

	switch source {
	case "iterator":

		q := u.Query()

		iterator_uri := q.Get("uri")
		iterator_sources := q["source"]

		return NewSFOMuseumLookupFromIterator(ctx, iterator_uri, iterator_sources...)

	default:

		return nil, fmt.Errorf("Invalid or unsupported constructor (URI)")
	}
}

// NewSFOMuseumLookup will return an `FlightsLookupFunc` function instance that, when invoked, will populate an `architecture.Lookup` instance with data stored in `r`.
// `r` will be closed when the `FlightsLookupFunc` function instance is invoked.
// It is assumed that the data in `r` will be formatted in the same way as the procompiled (embedded) data stored in `data/sfomuseum.json`.
func NewSFOMuseumLookupFuncWithReader(ctx context.Context, r io.ReadCloser) FlightsLookupFunc {

	defer r.Close()

	var flights_list []*Flight

	dec := json.NewDecoder(r)
	err := dec.Decode(&flights_list)

	if err != nil {

		lookup_func := func(ctx context.Context) {
			lookup_init_err = err
		}

		return lookup_func
	}

	return NewSFOMuseumLookupFuncWithFlights(ctx, flights_list)
}

func NewSFOMuseumLookupFuncWithFlights(ctx context.Context, flights_list []*Flight) FlightsLookupFunc {

	lookup_func := func(ctx context.Context) {

		table := new(sync.Map)

		for _, data := range flights_list {

			select {
			case <-ctx.Done():
				return
			default:
				// pass
			}

			appendData(ctx, table, data)
		}

		lookup_table = table
	}

	return lookup_func
}

// NewSFOMuseumLookupWithLookupFunc will return an `flights.Lookup` instance derived by data compiled using `lookup_func`.
func NewSFOMuseumLookupWithLookupFunc(ctx context.Context, lookup_func FlightsLookupFunc) (flights.FlightsLookup, error) {

	fn := func() {
		lookup_func(ctx)
	}

	lookup_init.Do(fn)

	if lookup_init_err != nil {
		return nil, lookup_init_err
	}

	l := FlightsLookup{}
	return &l, nil
}

func NewSFOMuseumLookupFromIterator(ctx context.Context, iterator_uri string, iterator_sources ...string) (flights.FlightsLookup, error) {

	flights_list, err := CompileFlightsData(ctx, iterator_uri, iterator_sources...)

	if err != nil {
		return nil, fmt.Errorf("Failed to compile flight data, %w", err)
	}

	lookup_func := NewSFOMuseumLookupFuncWithFlights(ctx, flights_list)
	return NewSFOMuseumLookupWithLookupFunc(ctx, lookup_func)
}

func (l *FlightsLookup) Find(ctx context.Context, code string) ([]interface{}, error) {

	pointers, ok := lookup_table.Load(code)

	if !ok {
		return nil, fmt.Errorf("Code '%s' not found", code)
	}

	flights_list := make([]interface{}, 0)

	for _, p := range pointers.([]string) {

		if !strings.HasPrefix(p, "pointer:") {
			return nil, fmt.Errorf("Invalid pointer '%s'", p)
		}

		row, ok := lookup_table.Load(p)

		if !ok {
			return nil, fmt.Errorf("Invalid pointer '%s'", p)
		}

		flights_list = append(flights_list, row.(*Flight))
	}

	return flights_list, nil
}

func (l *FlightsLookup) Append(ctx context.Context, data interface{}) error {
	return appendData(ctx, lookup_table, data.(*Flight))
}

func appendData(ctx context.Context, table *sync.Map, data *Flight) error {

	idx := atomic.AddInt64(&lookup_idx, 1)

	pointer := fmt.Sprintf("pointer:%d", idx)
	table.Store(pointer, data)

	str_wofid := strconv.FormatInt(data.WhosOnFirstId, 10)
	str_sfomid := data.SFOMuseumFlightId

	possible_codes := []string{
		str_wofid,
		str_sfomid,
	}

	for _, code := range possible_codes {

		if code == "" {
			continue
		}

		pointers := make([]string, 0)
		has_pointer := false

		others, ok := table.Load(code)

		if ok {

			pointers = others.([]string)
		}

		for _, dupe := range pointers {

			if dupe == pointer {
				has_pointer = true
				break
			}
		}

		if has_pointer {
			continue
		}

		pointers = append(pointers, pointer)
		table.Store(code, pointers)
	}

	return nil
}

// CompileFlightsData will generate a list of `Flight` struct to be used as the source data for an `SFOMuseumLookup` instance.
// The list of gate are compiled by iterating over one or more source. `iterator_uri` is a valid `whosonfirst/go-whosonfirst-iterate` URI
// and `iterator_sources` are one more (iterator) URIs to process.
func CompileFlightsData(ctx context.Context, iterator_uri string, iterator_sources ...string) ([]*Flight, error) {

	lookup := make([]*Flight, 0)
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
			return fmt.Errorf("Failed load feature from %s, %w", path, err)
		}

		wofid_rsp := gjson.GetBytes(body, "properties.wof:id")
		sfomid_rsp := gjson.GetBytes(body, "properties.sfomuseum:flight_id")

		if !wofid_rsp.Exists() {
			return fmt.Errorf("Missing wof:id property (%s)", path)
		}

		if !sfomid_rsp.Exists() {
			return fmt.Errorf("Missing sfomuseum:flight_id property (%s)", path)
		}

		fl := &Flight{
			WhosOnFirstId:     wofid_rsp.Int(),
			SFOMuseumFlightId: sfomid_rsp.String(),
		}

		mu.Lock()
		lookup = append(lookup, fl)
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
