package main

import (
	"fmt"
	"testing"
)

func TestProcessKind(t *testing.T) {

	tests := []struct {
		path   string
		result processType
		isErr  bool
	}{
		{
			path:   "testdata/A/%^&*()(___and",
			result: FILE,
			isErr:  false,
		},
		{
			path:   "testdata/b 1&2",
			result: DIR,
			isErr:  false,
		},
		{
			path:   "testdata/b 1&2/",
			result: WALK,
			isErr:  false,
		},
		{
			path:   "testdata/doesntExist", // does not exist
			result: NONE,
			isErr:  true,
		},
		{
			path:   "testdata/A/%^&*()(___and/", // not a directory
			result: NONE,
			isErr:  true,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			result, err := processKind(tt.path)
			if err != nil {
				fmt.Println(err)
			}
			if got, want := (err != nil), tt.isErr; got != want {
				t.Fatalf("unexpected error %v", err)
			}
			if got, want := result, tt.result; got != want {
				t.Errorf("result got %d want %d", got, want)
			}
		})
	}
}
