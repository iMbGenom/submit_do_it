package domain

import (
	"testing"
)

func TestSpotID(t *testing.T) {
	tests := []struct {
		name     string
		floor    int
		row      int
		col      int
		expected string
	}{
		{"AllZeros", 0, 0, 0, "0-0-0"},
		{"PositiveCase", 2, 3, 4, "2-3-4"},
		{"NegativeCase", -1, -2, -3, "-1--2--3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spot := &Spot{
				Floor: tt.floor,
				Row:   tt.row,
				Col:   tt.col,
			}
			got := spot.ID()
			if got != tt.expected {
				t.Errorf("Spot.ID() = %q, want %q", got, tt.expected)
			}
		})
	}
}
