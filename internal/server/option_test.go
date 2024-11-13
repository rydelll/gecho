package server

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestWithTimeout(t *testing.T) {
	cases := []struct {
		input   []byte
		want    []byte
		timeout time.Duration
	}{
		{timeout: time.Millisecond},
		{timeout: time.Second},
		{timeout: time.Minute},
	}

	for _, tc := range cases {
		t.Run(tc.timeout.String(), func(t *testing.T) {
			srv := new(Server)

			WithTimeout(tc.timeout)(srv)
			if diff := cmp.Diff(tc.timeout, srv.timeout); diff != "" {
				t.Errorf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}
