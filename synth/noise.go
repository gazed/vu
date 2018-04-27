// Copyright © 2014-2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package synth

// Noise exposes noise generating algorithms.
type Noise interface {
	Gen2D(x, y float64) float64
	Gen3D(x, y, z float64) float64
}

// =============================================================================

// SimplexNoise exposes the simplex noise generator and its
// control parameters for easy manipulation.
type SimplexNoise struct {
	F float64 // Higher frequency results in finer features.
	G float64 // Gain limits final value amount.
	L float64 // Using 2.0 gives scales of 1, 1/2, 1/4 for the octaves.
	O int     // More octaves for sharper features.
	N Noise   // Simplex noise algorithm.
}

// NewSimplexNoise initializes a simplex noise generator using the given
// seed. The returned generator can be used to create data blocks of
// generated values.
func NewSimplexNoise(seed int64) *SimplexNoise {
	sn := &SimplexNoise{F: 0.5, G: 0.55, L: 2.0, O: 6}
	sn.N = newSimplex(seed)
	return sn
}

// Gen2D returns a generated noise value for the given x,y coordinate.
// Used to generate different 2D images based on the SimplexNoise parameters.
func (sn *SimplexNoise) Gen2D(x, y float64) float64 {
	total := 0.0
	nfreq := sn.F
	amplitude := sn.G
	for o := 0; o < sn.O; o++ {
		xval := float64(x) * nfreq
		yval := float64(y) * nfreq
		total += sn.N.Gen2D(xval, yval) * amplitude
		nfreq *= sn.L
		amplitude *= sn.G
	}
	return total
}
