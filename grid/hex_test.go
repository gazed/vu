// Copyright © 2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package grid

import (
	"fmt"
	"testing"

	"github.com/gazed/vu/math/lin"
)

func TestNew(t *testing.T) {
	h, want := NewHex(2, 2), &Hex{2, 2, -4}
	if !h.Eq(want) {
		t.Errorf(format, h.Dump(), want.Dump())
	}
	h, want = NewHex(2, -3), &Hex{2, -3, 1}
	if !h.Eq(want) {
		t.Errorf(format, h.Dump(), want.Dump())
	}
}

func TestID(t *testing.T) {
	h, want := NewHex(2, 3), uint64(8589934595) // 10 0000000000000000000000000000011
	if h.ID() != want {
		t.Errorf("Wanted %d got %d", want, h.ID())
	}
	h, want = NewHex(1, -1), uint64(8589934591)
	if h.ID() != want {
		t.Errorf("Wanted %d got %d", want, h.ID())
	}
	h, want = NewHex(0, -1), uint64(4294967295)
	if h.ID() != want {
		t.Errorf("Wanted %d got %d", want, h.ID())
	}
	h, want = NewHex(-1, -1), uint64(18446744073709551615) // 10 0000000000000000000000000000011
	if h.ID() != want {
		t.Errorf("Wanted %d got %d", want, h.ID())
	}
}

func TestAdd(t *testing.T) {
	h, want := &Hex{1, 1, -2}, &Hex{2, 2, -4}
	if !h.Add(h, h).Eq(want) {
		t.Errorf(format, h.Dump(), want.Dump())
	}
}
func TestSub(t *testing.T) {
	h, want := &Hex{1, 1, -2}, &Hex{0, 0, 0}
	if !h.Sub(h, h).Eq(want) {
		t.Errorf(format, h.Dump(), want.Dump())
	}
}
func TestLen(t *testing.T) {
	h := &Hex{2, 2, -4}
	if h.Len() != 4 {
		t.Error("Invalid length", h.Len())
	}
}
func TestDist(t *testing.T) {
	h, a := &Hex{1, 1, -2}, &Hex{4, 4, -8}
	if h.Dist(a) != 6 {
		t.Errorf("Invalid distance %d", h.Dist(a))
	}
	if h.Dist(h) != 0 {
		t.Error("Distance with self should be zero.")
	}
}
func TestScale(t *testing.T) {
	h, want := &Hex{1, 1, -2}, &Hex{4, 4, -8}
	if !h.Mult(h, 4).Eq(want) {
		t.Errorf(format, h.Dump(), want.Dump())
	}
}

func TestMoveFlat(t *testing.T) {
	h, want := &Hex{1, 1, -2}, &Hex{1, 2, -3}
	if !h.Move(h, UP).Eq(want) {
		t.Errorf(format, h.Dump(), want.Dump())
	}
	h, want = &Hex{1, 1, -2}, &Hex{1, 0, -1}
	if !h.Move(h, DN).Eq(want) {
		t.Errorf(format, h.Dump(), want.Dump())
	}
	h, want = &Hex{1, 1, -2}, &Hex{2, 1, -3}
	if !h.Move(h, UR).Eq(want) {
		t.Errorf(format, h.Dump(), want.Dump())
	}
	h, want = &Hex{1, 1, -2}, &Hex{0, 1, -1}
	if !h.Move(h, DL).Eq(want) {
		t.Errorf(format, h.Dump(), want.Dump())
	}
	h, want = &Hex{1, 1, -2}, &Hex{0, 2, -2}
	if !h.Move(h, UL).Eq(want) {
		t.Errorf(format, h.Dump(), want.Dump())
	}
	h, want = &Hex{1, 1, -2}, &Hex{2, 0, -2}
	if !h.Move(h, DR).Eq(want) {
		t.Errorf(format, h.Dump(), want.Dump())
	}
}

func TestMovePointy(t *testing.T) {
	h, want := &Hex{1, 1, -2}, &Hex{2, 0, -2}
	if !h.Move(h, RT).Eq(want) {
		t.Errorf(format, h.Dump(), want.Dump())
	}
	h, want = &Hex{1, 1, -2}, &Hex{0, 2, -2}
	if !h.Move(h, LT).Eq(want) {
		t.Errorf(format, h.Dump(), want.Dump())
	}
	h, want = &Hex{1, 1, -2}, &Hex{2, 1, -3}
	if !h.Move(h, RU).Eq(want) {
		t.Errorf(format, h.Dump(), want.Dump())
	}
	h, want = &Hex{1, 1, -2}, &Hex{0, 1, -1}
	if !h.Move(h, LD).Eq(want) {
		t.Errorf(format, h.Dump(), want.Dump())
	}
	h, want = &Hex{1, 1, -2}, &Hex{1, 2, -3}
	if !h.Move(h, LU).Eq(want) {
		t.Errorf(format, h.Dump(), want.Dump())
	}
	h, want = &Hex{1, 1, -2}, &Hex{1, 0, -1}
	if !h.Move(h, RD).Eq(want) {
		t.Errorf(format, h.Dump(), want.Dump())
	}
}

func TestToPointy(t *testing.T) {
	h, x, y := &Hex{1, 1, -2}, 2.598076, 1.5
	if px, py := h.ToPointy(1.0); !lin.Aeq(x, px) || !lin.Aeq(y, py) {
		t.Errorf("Wanted %f %f got %f %f", x, y, px, py)
	}
}

const format = "\ngot\n%s\nwanted\n%s"

func (h *Hex) Dump() string { return fmt.Sprintf("%2d", *h) }
