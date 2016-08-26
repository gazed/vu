// Copyright © 2013-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package load

import (
	"testing"
)

var modDir = "../eg/models"

// Uses vu/eg resource directories.
func TestInvalidLoadObj(t *testing.T) {
	msh := &MshData{}
	if err := msh.Load("xxx", NewLocator().Dir("OBJ", modDir)); err == nil {
		t.Error("Should not be able to load an invalid file.")
	}
}

func TestCorruptLoadObj(t *testing.T) {
	msh := &MshData{}
	if err := msh.Load("corrupt", NewLocator().Dir("OBJ", modDir)); err == nil {
		t.Error("Should ignore corrupt input.")
	}
}

func TestLoadObj1(t *testing.T) {
	msh := &MshData{}
	if err := msh.Load("cube", NewLocator().Dir("OBJ", modDir)); err != nil {
		t.Fatal("Could not load cube.obj")
	}
	if len(msh.V) != 24 || len(msh.N) != 24 || len(msh.F) != 36 {
		t.Error("Improper sizes in cube.obj")
	}
}

func TestLoadObj2(t *testing.T) {
	msh := &MshData{}
	if err := msh.Load("monkey", NewLocator().Dir("OBJ", modDir)); err != nil {
		t.Fatal("Could not load monkey.obj")
	}
	if len(msh.V) != 1521 || len(msh.N) != 1521 || len(msh.F) != 2904 {
		t.Error("Improper sizes in monkey.obj")
	}
}

// A cube with uv maps needs duplicated verticies to get the proper
// mapping for each cube face.
func TestLoadObj3(t *testing.T) {
	msh := &MshData{}
	if err := msh.Load("block", NewLocator().Dir("OBJ", modDir)); err != nil {
		t.Fatal("Could not load block.obj")
	}
	if len(msh.V) != 60 || len(msh.N) != 60 || len(msh.F) != 36 {
		t.Error("Improper sizes in block.obj", len(msh.V), len(msh.N), len(msh.F))
	}
}
