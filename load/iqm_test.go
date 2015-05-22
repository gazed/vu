// Copyright © 2014-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package load

import (
	"testing"
)

// Uses vu/eg resource directories.
func TestLoadIqm(t *testing.T) {
	load := newLoader().setDir(mod, "../eg/models")
	if iqm, err := load.iqm("mrfixit"); err != nil && len(iqm.V) > 0 {
		// if iqm, err := load.iqm("rat"); err != nil && len(iqm.V) > 0 {
		t.Error(err)
	}
}
