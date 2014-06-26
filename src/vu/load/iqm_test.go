// Copyright © 2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package load

import (
	"testing"
)

func TestLoadIqm(t *testing.T) {
	load := newLoader().setDir(mod, "../eg/models")
	if iqm, err := load.iqm("mrfixit"); err != nil && len(iqm.V) > 0 {
		// if iqm, err := load.iqm("rat"); err != nil && len(iqm.V) > 0 {
		t.Error(err)
	}
}
