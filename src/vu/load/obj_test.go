// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package load

import (
	"testing"
)

func TestInvalidLoadObj(t *testing.T) {
	load := newLoader().setDir(mod, "../eg/models")
	meshes, err := load.obj("xxx")
	if len(meshes) != 0 || err == nil {
		t.Error("Should not be able to load an invalid file.")
	}
}

func TestCorruptLoadObj(t *testing.T) {
	load := newLoader().setDir(mod, "../eg/models")
	meshes, err := load.obj("corrupt")
	if len(meshes) != 0 || err == nil {
		t.Error("Should ignore corrupt input.")
	}
}

func TestLoadObj1(t *testing.T) {
	load := newLoader().setDir(mod, "../eg/models")
	meshes, err := load.obj("cube")
	if len(meshes) != 1 || err != nil {
		t.Fatal("Could not load cube.obj")
	}
	m := meshes[0]
	if len(m.V) != 24 || len(m.N) != 24 || len(m.F) != 36 {
		t.Error("Improper sizes in cube.obj")
	}
}

func TestLoadObj2(t *testing.T) {
	load := newLoader().setDir(mod, "../eg/models")
	meshes, err := load.obj("monkey")
	if len(meshes) != 1 || err != nil {
		t.Fatal("Could not load monkey.obj")
	}
	m := meshes[0]
	if len(m.V) != 1521 || len(m.N) != 1521 || len(m.F) != 2904 {
		t.Error("Improper sizes in monkey.obj")
	}
}

func TestLoadLevel(t *testing.T) {
	load := newLoader().setDir(mod, "../eg/models")
	ms, err := load.obj("level1")
	if len(ms) != 3 || err != nil {
		t.Fatalf("Loaded %d meshes from level1.obj", len(ms))
	}
	if ms[0].Name != "Glow1" || ms[1].Name != "Block1" || ms[2].Name != "Floor1" {
		t.Error("Invalid name level1.obj")
	}
}
