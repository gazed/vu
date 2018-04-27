// Copyright © 2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

import (
	"testing"
)

func TestClone(t *testing.T) {
	m := newMesh("meshTest")
	m.vao = 2
	m.InitData(1, 2, StaticDraw, true)
	m.SetData(1, []float32{2, 4, 6, 8, 7, 5})
	m.InitFaces(StaticDraw).SetFaces([]uint16{1, 2, 3, 4, 5, 6})

	// check that the deep copy worked.
	c := m.clone()
	if c.vao != m.vao {
		t.Errorf("clone failed vao %d %d", c.vao, m.vao)
	}
	if c.faces.Len() != m.faces.Len() {
		t.Errorf("clone failed faces %d %d", c.faces.Len(), m.faces.Len())
	}
	if c.vdata[1].Len() != m.vdata[1].Len() {
		t.Errorf("clone failed vdata %d %d", c.vdata[1].Len(), m.vdata[1].Len())
	}
}
