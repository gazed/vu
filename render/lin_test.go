// Copyright © 2013-2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package render

import (
	"fmt"
	"strconv"
	"testing"
	"unsafe"
)

// Check that golang lays out the data structure as sequential floats.
// Memory structures layout is important as the memory is handed down
// to the c-language graphics layer.
func TestMemoryLayout(t *testing.T) {
	x4 := m4{
		11, 12, 13, 14,
		21, 22, 23, 24,
		31, 32, 33, 34,
		41, 42, 43, 44}
	oneFloat := uint64(unsafe.Sizeof(x4.xx))
	fourFloats := oneFloat * 4
	mema, _ := strconv.ParseUint(fmt.Sprintf("%d", &(x4.xx)), 0, 64)
	memb, _ := strconv.ParseUint(fmt.Sprintf("%d", &(x4.xy)), 0, 64) // next value.
	if memb-mema != oneFloat {
		t.Errorf("Next value should be %d bytes. Was %d", oneFloat, memb-mema)
	}
	memc, _ := strconv.ParseUint(fmt.Sprintf("%d", &(x4.yx)), 0, 64) // next row.
	if memc-mema != fourFloats {
		t.Errorf("Next row should be %d bytes. Was %d", fourFloats, memc-mema)
	}
}

// Check that pointers are initialized to zero.
func TestNullPointer(t *testing.T) {
	var null0 unsafe.Pointer
	if ptr := uintptr(null0); ptr != 0 {
		t.Error("Default unsafe.Pointers should be zero.")
	}
}
