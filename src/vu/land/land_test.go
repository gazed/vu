// Copyright © 2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package land

import (
	"testing"
)

func TestZoomLimit(t *testing.T) {
	if l := newLand(9, 256, 12345); l.lod != 9 {
		t.Error("Level of detail should be 9")
	}
}

// Uncomment to see a world image. Recomment afterwards so as to not generate
// images during automated testing.
// func TestWorldImageGeneration(t *testing.T) {
// 	l := newLand(0, 256, 124)
// 	writeImage("target/", "land.png", l.image(-0.25))
// }
