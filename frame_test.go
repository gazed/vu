// Copyright © 2017 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

import (
	"math"
	"testing"
)

// Check that the bits are at the right spots.
func TestSetBucket(t *testing.T) {
	bucket := setBucket(2, 200)
	bucket = setDist(bucket, 0.4)
	dist := math.Float32frombits(uint32(bucket & 0x00000000FFFFFFFF))
	over := uint16((bucket & 0x0000FFFF00000000) >> 40)
	pass := uint8((bucket & 0x00FF000000000000) >> 48)
	if dist != 0.4 || over != 55 || pass != 2 {
		t.Errorf("Bad bucket %016x %f %d %d\n", bucket, dist, over, pass)
	}
}

// Check setting transparency on a bucket.
func TestSetTransparent(t *testing.T) {
	b := setBucket(2, 200)
	b = setDist(b, 0.4)
	bt := setTransparent(b)
	dist := math.Float32frombits(uint32(bt & 0x00000000FFFFFFFF))
	over := uint16((bt & 0x0000FFFF00000000) >> 40)
	pass := uint8((bt & 0x00FF000000000000) >> 48)
	alpha := uint8((bt & alphaObj) >> 32)
	if dist != 0.4 || over != 55 || pass != 2 || alpha != 1 {
		t.Errorf("Bad bucket %016x %f %d %d %d\n", bt, dist, over, pass, alpha)
	}
}
