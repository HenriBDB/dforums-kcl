package communication_test

import (
	"bytes"
	"testing"
)

// https://gist.github.com/samalba/6059502
func assertEqualBytes(t *testing.T, a []byte, b []byte) {
	if !bytes.Equal(a, b) {
		t.Fatalf("%d != %d", a, b)
	}
}
