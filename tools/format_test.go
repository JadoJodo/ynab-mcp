package tools

import "testing"

func TestFormatMilliunits(t *testing.T) {
	tests := []struct {
		input int64
		want  string
	}{
		{0, "$0.00"},
		{1000, "$1.00"},
		{12500, "$12.50"},
		{-5000, "-$5.00"},
		{999, "$1.00"},
		{1, "$0.00"},
		{500, "$0.50"},
		{-123450, "-$123.45"},
	}
	for _, tt := range tests {
		got := FormatMilliunits(tt.input)
		if got != tt.want {
			t.Errorf("FormatMilliunits(%d) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestMilliunitsFromDollars(t *testing.T) {
	tests := []struct {
		input float64
		want  int64
	}{
		{0, 0},
		{1.00, 1000},
		{12.50, 12500},
		{-5.00, -5000},
		{0.999, 999},
		{123.456, 123456},
	}
	for _, tt := range tests {
		got := MilliunitsFromDollars(tt.input)
		if got != tt.want {
			t.Errorf("MilliunitsFromDollars(%f) = %d, want %d", tt.input, got, tt.want)
		}
	}
}
