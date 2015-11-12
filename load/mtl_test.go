// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package load

import (
	"fmt"
	"testing"
)

// Uses vu/eg resource directories.
func TestLoadMtl(t *testing.T) {
	load := newLoader().setDir(mod, "../eg/models")
	m, err := load.mtl("red")
	if m == nil || err != nil {
		t.Fatalf("Should be able to load a valid material file %s", err)
	}
	got, want := fmt.Sprintf("%2.1f %2.1f %2.1f", m.KdR, m.KdG, m.KdB), "0.8 0.6 0.2"
	if got != want {
		t.Errorf(format, got, want)
	}
}
