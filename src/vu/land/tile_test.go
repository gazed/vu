// Copyright © 2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package land

import (
	"testing"
)

// Unique map keys from tiles store and fetch map tiles.
func TestTileKey(t *testing.T) {
	if key := tileKey(1, 0, 0); key != "0" {
		t.Errorf("Expected key 0, got %s", key)
	}
	if key := tileKey(1, 1, 0); key != "1" {
		t.Errorf("Expected key 1, got %s", key)
	}
	if key := tileKey(1, 0, 1); key != "2" {
		t.Errorf("Expected key 2, got %s", key)
	}
	if key := tileKey(1, 1, 1); key != "3" {
		t.Errorf("Expected key 3, got %s", key)
	}
	if key := tileKey(3, 7, 7); key != "333" {
		t.Errorf("Expected key 333, got %s", key)
	}
	if key := tileKey(8, 255, 15); key != "11113333" {
		t.Errorf("Expected key 11113333, got %s", key)
	}
	if key := tileKey(8, 15, 255); key != "22223333" {
		t.Errorf("Expected key 22223333, got %s", key)
	}
}

// Reverse of the TileKey tests.
func TestKeyTile(t *testing.T) {
	if z, x, y := keyTile("0"); z != 1 || x != 0 || y != 0 {
		t.Errorf("Expected 1 0 0, got %d %d %d", z, x, y)
	}
	if z, x, y := keyTile("1"); z != 1 || x != 1 || y != 0 {
		t.Errorf("Expected 1 1 0, got %d %d %d", z, x, y)
	}
	if z, x, y := keyTile("2"); z != 1 || x != 0 || y != 1 {
		t.Errorf("Expected 1 0 1, got %d %d %d", z, x, y)
	}
	if z, x, y := keyTile("3"); z != 1 || x != 1 || y != 1 {
		t.Errorf("Expected 1 1 1, got %d %d %d", z, x, y)
	}
	if z, x, y := keyTile("333"); z != 3 || x != 7 || y != 7 {
		t.Errorf("Expected 3 7 7, got %d %d %d", z, x, y)
	}
	if z, x, y := keyTile("11113333"); z != 8 || x != 255 || y != 15 {
		t.Errorf("Expected 8 255 15, got %d %d %d", z, x, y)
	}
	if z, x, y := keyTile("22223333"); z != 8 || x != 15 || y != 255 {
		t.Errorf("Expected 8 15 255, got %d %d %d", z, x, y)
	}
}

// Uncomment to see a tile images. Recomment afterwards so as to not generate
// images during automated testing.
// func TestTileImageGeneration(t *testing.T) {
// 	l := newLand(1, 256, 124)
// 	writeImage("target/", "tile0.png", l.newTile(0, 0, 0).image(-0.25))
// 	writeImage("target/", "tile00.png", l.newTile(1, 0, 0).image(-0.25))
// 	writeImage("target/", "tile01.png", l.newTile(1, 0, 1).image(-0.25))
// 	writeImage("target/", "tile10.png", l.newTile(1, 1, 0).image(-0.25))
// 	writeImage("target/", "tile11.png", l.newTile(1, 1, 1).image(-0.25))
// }
