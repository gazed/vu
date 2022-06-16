// Copyright © 2013-2024 Galvanized Logic Inc.

package al

import "testing"

// The test passes if the binding layer can initialize without crashing.
func TestInit(t *testing.T) {
	Init()
	Dump()
}
