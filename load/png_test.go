// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package load

import (
	"testing"
)

// Uses vu/eg resource directories.
func TestLoadPng(t *testing.T) {
	load := newLoader().setDir(img, "../eg/images")
	if img, _ := load.png("xxx"); img != nil {
		t.Error("Image should be nil for bad files")
	}
	if img, _ := load.png("image"); img == nil {
		t.Error("Could not load image file")
	}
}
