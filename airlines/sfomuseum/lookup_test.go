package sfomuseum

import (
	"context"
	"github.com/sfomuseum/go-sfomuseum-airfield"
	"testing"
)

func TestSFOMuseumLookup(t *testing.T) {

	wofid_tests := map[string]int64{
		"AC":         1159283597,
		"ACA":        1159283597,
		"AIR CANADA": 1159283597,
		"MOV":        1360700753,
		"NN":         1360700753,
		"77":         1159283643,
		"AHC":        1159283643,
	}

	schemes := []string{
		"airlines://sfomuseum",
		"airlines://sfomuseum/github",
	}

	ctx := context.Background()

	for _, s := range schemes {

		lu, err := airfield.NewLookup(ctx, s)

		if err != nil {
			t.Fatalf("Failed to create lookup for '%s', %v", s, err)
		}

		for code, wofid := range wofid_tests {

			results, err := lu.Find(ctx, code)

			if err != nil {
				t.Fatalf("Unable to find '%s' using scheme '%s', %v", code, s, err)
			}

			if len(results) != 1 {
				t.Fatalf("Invalid results for '%s' using scheme '%s'", code, s)
			}

			a := results[0].(*Airline)

			if a.WhosOnFirstId != wofid {
				t.Fatalf("Invalid match for '%s', expected %d but got %d using scheme '%s'", code, wofid, a.WhosOnFirstId, s)
			}
		}
	}

}
