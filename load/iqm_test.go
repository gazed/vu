// Copyright © 2014-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package load

import (
	"testing"
)

// Uses vu/eg resource directories.
func TestLoadIqm(t *testing.T) {
	m := &ModData{}
	err := m.Load("rat", NewLocator().Dir("IQM", "../eg/models"))
	if err != nil || len(m.V) <= 0 {
		t.Error(err)
	}
}
