package main

import (
	"reflect"
	"testing"
)

func TestParseDifficulties(t *testing.T) {
	tests := []struct {
		in      string
		want    []int
		wantErr bool
	}{
		{"4,5", []int{4, 5}, false},
		{"6", []int{6}, false},
		{" 4 , 5 ,6 ", []int{4, 5, 6}, false},
		{"64", []int{64}, false},
		{"", nil, true},
		{"4,", nil, true},
		{"4,a", nil, true},
		{"0", nil, true},
		{"65", nil, true},
		{"-1", nil, true},
	}
	for _, tt := range tests {
		got, err := parseDifficulties(tt.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("parseDifficulties(%q) error = %v, wantErr %v", tt.in, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("parseDifficulties(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}
