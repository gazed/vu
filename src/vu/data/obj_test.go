// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package data

import (
	"testing"
)

func TestInvalidLoadObj(t *testing.T) {
	load := &loader{}
	meshes, err := load.obj("../eg/models", "xxx.obj")
	if len(meshes) != 0 || err == nil {
		t.Error("Should not be able to load an invalid file.")
	}
}

func TestCorruptLoadObj(t *testing.T) {
	load := &loader{}
	meshes, err := load.obj("../eg/models", "corrupt.obj")
	if len(meshes) != 0 || err == nil {
		t.Error("Should ignore corrupt input.")
	}
}

func TestLoadObj1(t *testing.T) {
	load := &loader{}
	meshes, err := load.obj("../eg/models", "cube.obj")
	if len(meshes) != 1 || err != nil {
		t.Fatal("Could not load cube.obj")
	}
	m := meshes[0]
	if len(m.V) != 32 || len(m.N) != 24 || len(m.F) != 36 {
		t.Error("Improper sizes in cube.obj")
	}
}

func TestLoadObj2(t *testing.T) {
	load := &loader{}
	meshes, err := load.obj("../eg/models", "monkey.obj")
	if len(meshes) != 1 || err != nil {
		t.Fatal("Could not load monkey.obj")
	}
	m := meshes[0]
	if len(m.V) != 2028 || len(m.N) != 1521 || len(m.F) != 2904 {
		t.Error("Improper sizes in monkey.obj")
	}
}

func TestLoadLevel(t *testing.T) {
	load := &loader{}
	ms, err := load.obj("../eg/models", "level1.obj")
	if len(ms) != 3 || err != nil {
		t.Fatalf("Loaded %d meshes from level1.obj", len(ms))
	}
	if ms[0].Name != "Glow1" || ms[1].Name != "Block1" || ms[2].Name != "Floor1" {
		t.Error("Invalid name level1.obj")
	}
}

func TestLoadTile(t *testing.T) {
	load := &loader{}
	meshes, err := load.obj("../eg/models", "tile.obj")
	if len(meshes) != 1 || err != nil {
		t.Fatalf("Could not load %d tile.obj %s", len(meshes), err)
	}
	m := meshes[0]
	if len(m.V) != 56 || len(m.N) != 42 || len(m.T) != 28 || len(m.F) != 36 {
		t.Error("Improper sizes in tile.obj")
	}
}
