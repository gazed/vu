// Copyright © 2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package grid

// hex.go is a start to capturing some hex grid related code.
//
// HUGE thank you to one of the best educational websites anywhere
// and the authoritative site for hex grids:
//    http://www.redblobgames.com/grids/hexagons/
//    http://www.redblobgames.com/grids/hexagons/implementation.html

import (
	"math"
)

// Hex represents the location of one hex within a grid of hexes.
// The locations are relative to the map origin is at 0, 0, 0.
// The coordinates are uint32 so a unique identifier can be created
// from the combination of Q and R.
type Hex struct {
	Q, R, S int32 // Cube coordinates: sum equals 0.
}

// NewHex returns a hex such that the coordinates Q+R-S=0.
func NewHex(q, r int32) *Hex {
	return &Hex{Q: q, R: r, S: -q - r}
}

// ZH is Zero hex. Used when a non-nil hex is needed.
// Value is always expected to be 0, 0, 0.
var ZH = &Hex{0, 0, 0}

// Eq checks for equality and returns true if Hex h has identical
// Q,R,S values to the given Hex a.
func (h *Hex) Eq(a *Hex) bool { return h.Q == a.Q && h.R == a.R && h.S == a.S }

// Add (+) adds hexes b and a storing the results of the addition in h.
// Hex h may be used as one or both of the parameters.
// For example (+=) is
//     h.Add(h, b)
// The updated Hex h is returned.
func (h *Hex) Add(a, b *Hex) *Hex {
	h.Q, h.R, h.S = a.Q+b.Q, a.R+b.R, a.S+b.S
	return h
}

// Sub (-) subtracts hex b from a storing the results of the subtraction in h.
// Hex h may be used as one or both of the parameters.
// For example (-=) is
//     h.Sub(h, b)
// The updated Hex h is returned.
func (h *Hex) Sub(a, b *Hex) *Hex {
	h.Q, h.R, h.S = a.Q-b.Q, a.R-b.R, a.S-b.S
	return h
}

// Mult (*) multiplies hex a by scale k, storing the results of the
// multiplication in h. Hex h may be used as the parameter.
// For example (*=) is
//     h.Mult(h, k)
// The updated Hex h is returned.
func (h *Hex) Mult(a *Hex, k int32) *Hex {
	h.Q, h.R, h.S = a.Q*k, a.R*k, a.S*k
	return h
}

// ID returns a single unique value for this hex tile.
// It is a bitwise combination of Q and R.
func (h *Hex) ID() uint64 {
	return (uint64(h.Q) << 32) | (0x00000000FFFFFFFF & uint64(h.R))
}

// Len determines the distance of Hex h to the origin 0,0,0.
func (h *Hex) Len() int {
	return int((math.Abs(float64(h.Q)) + math.Abs(float64(h.R)) + math.Abs(float64(h.S))) * 0.5)
}

// Dist returns the distance between hexes h and a.
func (h *Hex) Dist(a *Hex) int {
	dq, dr, ds := h.Q-a.Q, h.R-a.R, h.S-a.S
	return int((math.Abs(float64(dq)) + math.Abs(float64(dr)) + math.Abs(float64(ds))) * 0.5)
}

// ToPointy returns the 2D pointy grid location for this hex
// location using the given hex size.
func (h *Hex) ToPointy(size float64) (x, y float64) {
	sqrtOf3 := 1.732050807569
	x = size * sqrtOf3 * (float64(h.Q) + float64(h.R)*0.5)
	y = size * 1.5 * float64(h.R)
	return x, y
}

// ToFlat returns the 2D flat grid location for this hex
// location using the given hex size.
func (h *Hex) ToFlat(size float64) (x, y float64) {
	sqrtOf3 := 1.732050807569
	x = size * 1.5 * float64(h.Q)
	y = size * sqrtOf3 * (float64(h.R) + float64(h.Q)*0.5)
	return x, y
}

// The six hex grid direction constants and related coordinate differences.
// The movement directions reflect the movement angles for the two types
// of hex grids. Note the relative change to the QRS hex coordinate.
//            FLAT                   POINTY
//        -+0  0+-  +0-          0+-       +0-
//         UL   UP  UR            LU       RU
//            ↖ ⇧ ↗                  ↖  ↗
//                            -+0 LT ⇦  ⇨  RT +-0
//            ↙ ⇩ ↘                  ↙  ↘
//         DL   DN  DR            LD       RD
//        -0+  0-+  +-0          -0+       0-+
const (
	// Flat grids use Up/Down prefixes.
	DR = 0 // down and right
	UL = 1 // up and left
	UR = 2 // up and right
	DL = 3 // down and left
	UP = 4 // up
	DN = 5 // down

	// Pointy grids use Right/Left prefixes.
	RT = 0 // right
	LT = 1 // left
	RU = 2 // right and up
	LD = 3 // left and down
	LU = 4 // left and up
	RD = 5 // right and down
)

// offsets to move from a hex to one of the 6 possible neighbouring hexes.
// The indicies are different for flat/pointy hex grids.
var offsets = []Hex{ //                                   Flat   Pointy
	{Q: 1, R: -1, S: 0}, {Q: -1, R: 1, S: 0}, // S axis : DR/UL  RT/LT
	{Q: 1, R: 0, S: -1}, {Q: -1, R: 0, S: 1}, // R axis : UR/DL  RU/LD
	{Q: 0, R: 1, S: -1}, {Q: 0, R: -1, S: 1}, // Q axis : UP/DN  LU/RD
}

// Diff returns the hex that is added to the current hex to
// move in the given direction.
// Return the zero-hex if the movement is unrecognized.
func Diff(dir int) *Hex {
	if 0 <= dir && dir < len(offsets) {
		return &offsets[dir]
	}
	return ZH // zero hex for unrecognized directions.
}

// Move returns the next hex h when travelling in the given
// direction dir from hex a.
// Hex h may be used as one or both of the parameters, eg:
//     h.Next(h, N, moveNS)
// The updated vector h is returned.
func (h *Hex) Move(a *Hex, dir int) *Hex { return h.Add(a, Diff(dir)) }
