// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package al

import "testing"

// The test passes if the binding layer can initialize without crashing.
func TestInit(t *testing.T) {
	Init()
}
