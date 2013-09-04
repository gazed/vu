// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package data

import (
	"fmt"
	"testing"
)

// Check that a material file can be imported.
func TestLoadMtl(t *testing.T) {
	load := &loader{}
	m, err := load.mtl("../eg/models", "cube.mtl")
	if m == nil || err != nil {
		t.Fatalf("Should be able to load a valid material file %s", err)
	}
	got, want := fmt.Sprintf("%2.1f", m.Kd), "{0.6 0.6 0.6}"
	if got != want {
		t.Errorf(format, got, want)
	}
	got, want = fmt.Sprintf("%2.1f", m.Ka), "{0.2 0.2 0.2}"
	if got != want {
		t.Errorf(format, got, want)
	}
	got, want = fmt.Sprintf("%2.1f", m.Ks), "{0.5 0.5 0.5}"
	if got != want {
		t.Errorf(format, got, want)
	}
}
