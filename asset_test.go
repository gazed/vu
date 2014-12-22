// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

import (
	"testing"
)

// Check fetching and retrieving.
func TestDepot(t *testing.T) {

	// create and cache the asset.
	m := newMaterial("red")
	dep := newDepot()
	dep.cache(mat, m)

	// fetch a non-existing asset.
	m2 := newMaterial("blue")
	data := asset(m2)
	if err := dep.fetch(mat, &data); err == nil {
		t.Errorf("Should have failed.")
	}

	// fetch the cached asset.
	m3 := newMaterial("red")
	data = asset(m3)
	if err := dep.fetch(mat, &data); err != nil {
		t.Errorf("Should have worked %s", err)
	}
}
