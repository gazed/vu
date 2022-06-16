// Copyright © 2024 Galvanized Logic Inc.

package vu

import (
	"testing"
)

// go test -run Asset
func TestAssetID(t *testing.T) {

	// Different asset types with the same name must still have unique IDs.
	t.Run("uniqueID", func(t *testing.T) {
		alist := []asset{
			newMesh("abc"),
			newShader("abc"),
			newTexture("abc"),
			newMaterial("abc"),
			newSound("abc"),
		}
		for i, a := range alist {
			for j, b := range alist {
				if i == j {
					continue
				}
				if a.aid() == b.aid() {
					t.Errorf("asset ID should have been unique")
				}
			}
		}
	})

	// A specific asset type with the same name has the same ID.
	t.Run("sameID", func(t *testing.T) {
		m1 := newMesh("abc")
		m2 := newMesh("abc")
		if m1.aid() != m2.aid() {
			t.Errorf("asset ID must be the same")
		}
	})

}
