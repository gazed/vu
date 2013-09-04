// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package data

import (
	"testing"
)

// Check that a PNG image can be imported.
func TestLoadPng(t *testing.T) {
	load := &loader{}
	if img, _ := load.png("../eg/images", "xxx.png"); img != nil {
		t.Error("Image should be nil for bad files")
	}
	if img, _ := load.png("../eg/images", "image.png"); img == nil {
		t.Error("Could not load image file")
	}
}
