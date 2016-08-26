// Copyright © 2013-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package load

import (
	"testing"
)

// Uses vu/eg resource directories.
func TestLoadPng(t *testing.T) {
	img := &ImgData{}
	loc := NewLocator().Dir("PNG", "../eg/images")
	if err := img.Load("xxx", loc); err == nil {
		t.Error("Image should fail for bad files")
	}
	if err := img.Load("image", loc); err != nil {
		t.Error("Could not load image file")
	}
}
