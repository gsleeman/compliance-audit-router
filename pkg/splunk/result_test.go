package splunk

import (
	"reflect"
	"testing"
	"time"
)

func TestSearchResult_string(t *testing.T) {
	type args struct {
		field string
	}
	tests := []struct {
		name string
		a    SearchResult
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.string(tt.args.field); got != tt.want {
				t.Errorf("SearchResult.string() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchResult_slice(t *testing.T) {
	type args struct {
		field string
	}
	tests := []struct {
		name string
		a    SearchResult
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.slice(tt.args.field); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchResult.slice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchResult_time(t *testing.T) {
	type args struct {
		field string
	}
	tests := []struct {
		name string
		a    SearchResult
		args args
		want time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.time(tt.args.field); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchResult.time() = %v, want %v", got, tt.want)
			}
		})
	}
}
