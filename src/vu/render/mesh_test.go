// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package render

import (
	"testing"
)

// Check that meshes know how much data they contain.
func TestMeshSize(t *testing.T) {
	m := newMesh("test")
	m.InitData(0, 3, 0, false).SetData(0, []float32{0, 0, 0}) // 12 bytes
	m.InitFaces(0).SetFaces([]uint16{0, 0, 0, 0})             // 8 bytes
	if m.Size() != 20 {
		t.Errorf("Expected 20, got %d", m.Size())
	}

}

func TestValidMesh(t *testing.T) {
	m := newMesh("test")
	m.InitData(0, 3, 0, false).SetData(0, []float32{0, 0, 0, 0, 0, 0})       // 2 verticies
	m.InitData(1, 4, 0, false).SetData(1, []float32{0, 0, 0, 0, 0, 0, 0, 0}) // 2 verticies
	m.InitData(2, 2, 0, false).SetData(2, []float32{0, 0, 0, 0})             // 2 verticies
	if !m.valid() || m.numVerticies() != 2 {
		t.Errorf("Expected valid mesh with 2 vertices. Got %d", m.numVerticies())
	}
}

func TestInvalidMesh(t *testing.T) {
	m := newMesh("test")
	if m.valid() {
		t.Errorf("Got valid mesh when expecting invalid mesh with no data")
	}
	m.InitData(0, 3, 0, false).SetData(0, []float32{0, 0, 0, 0, 0, 0}) // 2 verticies
	m.InitData(1, 4, 0, false).SetData(1, []float32{0, 0, 0, 0})       // 1 vertex
	m.InitData(2, 2, 0, false).SetData(2, []float32{0, 0, 0, 0})       // 2 verticies
	if m.valid() {
		t.Errorf("Got valid mesh when expecting invalid mesh with mismatched counts")
	}
}
