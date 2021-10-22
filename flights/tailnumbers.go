package flights

import (
	"context"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/iterator"
	"io"
	"sort"
	"sync"
)

// TailNumbersFromIterator will return the unique set of `swim:tail_number` values from one or more whosonfirst/go-whosonfirst-iterate sources for SFO Museum flight data.
func TailNumbersFromIterator(ctx context.Context, iterator_uri string, iterator_sources ...string) ([]string, error) {

	catalog := new(sync.Map)

	iter_cb := func(ctx context.Context, path string, r io.ReadSeeker, args ...interface{}) error {

		select {
		case <-ctx.Done():
			return nil
		default:
			// pass
		}
		body, err := io.ReadAll(r)

		if err != nil {
			return fmt.Errorf("Failed to read %s, %w", path, err)
		}

		rsp := gjson.GetBytes(body, "properties.swim:tail_number")

		if !rsp.Exists() {
			return nil
		}

		v := rsp.String()

		if v == "" {
			return nil
		}

		catalog.Store(v, true)
		return nil
	}

	iter, err := iterator.NewIterator(ctx, iterator_uri, iter_cb)

	if err != nil {
		return nil, fmt.Errorf("Failed to create iterator, %v", err)
	}

	err = iter.IterateURIs(ctx, iterator_sources...)

	if err != nil {
		return nil, fmt.Errorf("Failed to iterate sources, %v", err)
	}

	values := make([]string, 0)

	catalog.Range(func(v interface{}, ignore interface{}) bool {
		values = append(values, v.(string))
		return true
	})

	sort.Strings(values)
	return values, nil
}
