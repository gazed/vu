// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package render

import (
	"fmt"
	"strconv"
	"testing"
	"unsafe"
)

// Check that golang lays out the data structure as sequential floats.
// This is important as memory structures will be handed down to the
// c-language graphics layer.
func TestMemoryLayout(t *testing.T) {
	m4 := M4{
		11, 12, 13, 14,
		21, 22, 23, 24,
		31, 32, 33, 34,
		41, 42, 43, 44}
	oneFloat := uint64(unsafe.Sizeof(m4.X0))
	fourFloats := oneFloat * 4
	mema, _ := strconv.ParseUint(fmt.Sprintf("%d", &(m4.X0)), 0, 64)
	memb, _ := strconv.ParseUint(fmt.Sprintf("%d", &(m4.Y0)), 0, 64) // next value.
	if memb-mema != oneFloat {
		t.Errorf("Next value should be %d bytes. Was %d", oneFloat, memb-mema)
	}
	memc, _ := strconv.ParseUint(fmt.Sprintf("%d", &(m4.X1)), 0, 64) // next row.
	if memc-mema != fourFloats {
		t.Errorf("Next row should be %d bytes. Was %d", fourFloats, memc-mema)
	}
}
