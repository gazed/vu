// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package load

import (
	"testing"
)

// Startup don't-crash-test.
func TestCreateLoader(t *testing.T) {
	if load := newLoader().setDir(mod, "../eg/models"); load == nil {
		t.Error("Can't create a new loader.")
	}
}
