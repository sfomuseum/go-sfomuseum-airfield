package sfomuseum

import (
	"context"
	"testing"
)

func TestFindCurrentAirline(t *testing.T) {

	tests := map[string]int64{
		"AC":         1159283597,
		"ACA":        1159283597,
		"AIR CANADA": 1159283597,
	}

	ctx := context.Background()

	for code, id := range tests {

		g, err := FindCurrentAirline(ctx, code)

		if err != nil {
			t.Fatalf("Failed to find current airline for %s, %v", code, err)
		}

		if g.WhosOnFirstId != id {
			t.Fatalf("Unexpected ID for airline %s. Got %d but expected %d", code, g.WhosOnFirstId, id)
		}
	}
}
