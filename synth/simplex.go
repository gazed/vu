// Copyright © 2014-2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package synth

import (
	"math"
	"math/rand"
	"time"
)

// simplex is a gradient noise generator algorithm. Its purpose is to help
// create random world maps and images. The random seed can be saved to
// recreate a given map at a later date.
//
// From and thanks to:
//    http://staffwww.itn.liu.se/~stegu/simplexnoise/simplexnoise.pdf
//    http://www.itn.liu.se/~stegu/simplexnoise/SimplexNoise.java
// Modified with
//	  https://github.com/Ian-Parberry/Tobler/
// Be aware of patent on 3D simplex and higher although this code is based
// on an clean opensource implementation that should not be affected by the patent.
//    https://www.google.com/patents/US6867776
type simplex struct {
	f2, f3    float64     // skewing and unskewing factors...
	g2, g3    float64     // ...for 2 and 3 dimensions.
	seed      int64       // seed that is unique for a given map.
	pseudo    []byte      // pseudo randomly ordered numbers 0-255
	perm      []byte      // 512 pseudo random numbers from 0-255
	permMod12 []byte      // pseudo random numbers mod 12
	gradients []*gradient // slopes to adjacent points.
	rng       *rand.Rand  // random number generator.

	// Exponentially Distributed Noise by Ian Parberry.
	mag    []float64 // magnitudes.
	magExp float64   // gradient magnitude exponent.
}

// newSimplex generates random noise based on a seed. Use 0 for the
// seed to create a new random map. Use a previous seed to get
// a previously generated map.
func newSimplex(seed int64) *simplex {
	s := &simplex{}
	s.seed = seed
	if s.seed == 0 { // create random seed if one not specified.
		s.seed = time.Now().UnixNano()
	}
	s.rng = rand.New(rand.NewSource(s.seed))
	pbase := []byte{151, 160, 137, 91, 90, 15,
		131, 13, 201, 95, 96, 53, 194, 233, 7, 225, 140, 36, 103, 30, 69, 142, 8, 99, 37, 240, 21, 10, 23,
		190, 6, 148, 247, 120, 234, 75, 0, 26, 197, 62, 94, 252, 219, 203, 117, 35, 11, 32, 57, 177, 33,
		88, 237, 149, 56, 87, 174, 20, 125, 136, 171, 168, 68, 175, 74, 165, 71, 134, 139, 48, 27, 166,
		77, 146, 158, 231, 83, 111, 229, 122, 60, 211, 133, 230, 220, 105, 92, 41, 55, 46, 245, 40, 244,
		102, 143, 54, 65, 25, 63, 161, 1, 216, 80, 73, 209, 76, 132, 187, 208, 89, 18, 169, 200, 196,
		135, 130, 116, 188, 159, 86, 164, 100, 109, 198, 173, 186, 3, 64, 52, 217, 226, 250, 124, 123,
		5, 202, 38, 147, 118, 126, 255, 82, 85, 212, 207, 206, 59, 227, 47, 16, 58, 17, 182, 189, 28, 42,
		223, 183, 170, 213, 119, 248, 152, 2, 44, 154, 163, 70, 221, 153, 101, 155, 167, 43, 172, 9,
		129, 22, 39, 253, 19, 98, 108, 110, 79, 113, 224, 232, 178, 185, 112, 104, 218, 246, 97, 228,
		251, 34, 242, 193, 238, 210, 144, 12, 191, 179, 162, 241, 81, 51, 145, 235, 249, 14, 239, 107,
		49, 192, 214, 31, 181, 199, 106, 157, 184, 84, 204, 176, 115, 121, 50, 45, 127, 4, 150, 254,
		138, 236, 205, 93, 222, 114, 67, 29, 24, 72, 243, 141, 128, 195, 78, 66, 215, 61, 156, 180}

	// reorder the pseudo random numbers for each generator.
	for len(pbase) > 0 {
		index := s.rng.Intn(len(pbase))
		s.pseudo = append(s.pseudo, pbase[index])
		pbase = append(pbase[:index], pbase[index+1:]...)
	}
	s.gradients = []*gradient{
		{1, 1, 0}, {-1, 1, 0}, {1, -1, 0}, {-1, -1, 0},
		{1, 0, 1}, {-1, 0, 1}, {1, 0, -1}, {-1, 0, -1},
		{0, 1, 1}, {0, -1, 1}, {0, 1, -1}, {0, -1, -1}}
	s.perm = make([]byte, 512)
	s.permMod12 = make([]byte, 512)
	for cnt := 0; cnt < 512; cnt++ {
		s.perm[cnt] = s.pseudo[cnt&255]
		s.permMod12[cnt] = s.perm[cnt] % 12
	}

	// Skewing and unskewing factors for 2 and 3 dimensions
	s.f2, s.g2 = 0.5*(math.Sqrt(3.0)-1.0), (3.0-math.Sqrt(3.0))/6.0
	s.f3, s.g3 = 1.0/3.0, 1.0/6.0

	//gradient magnitude array. Base on: https://github.com/Ian-Parberry/Tobler/
	m := 1.0        //current magnitude
	s.magExp = 1.02 ///< Mu, the gradient magnitude exponent.
	s.mag = make([]float64, 512)
	for cnt := 0; cnt < 512; cnt++ {
		s.mag[cnt] = m
		m /= s.magExp
	}
	return s
}

// dot2D provides a dot product of the point with the gradient.
func (s *simplex) dot2D(g *gradient, x, y float64) float64 {
	return g.x*x + g.y*y
}

// dot3D provides a dot product of the point with the gradient.
func (s *simplex) dot3D(g *gradient, x, y, z float64) float64 {
	return g.x*x + g.y*y + g.z*z
}

// floor was benchmarked to be a *lot* faster than (int)Math.floor(x).
func (s *simplex) floor(x float64) int {
	xi := int(x)
	if x < float64(xi) {
		return xi - 1
	}
	return xi
}

// generate 2D simplex noise.
// From http://www.its.liu.se/~stegu/simplexnoise/SimplexNoise.java
func (s *simplex) Gen2D(xin, yin float64) float64 {
	var n0, n1, n2 float64 // Noise contributions from the three corners

	// Skew the input space to determine which simplex cell we're in
	h := (xin + yin) * s.f2 // Hairy factor for 2D
	i := int(s.floor(xin + h))
	j := int(s.floor(yin + h))
	t := float64(i+j) * s.g2
	xs := float64(i) - t // Unskew the cell origin back to (x,y) space
	ys := float64(j) - t
	x0 := xin - xs // The x,y distances from the cell origin
	y0 := yin - ys

	// For the 2D case, the simplex shape is an equilateral triangle.
	// Determine which simplex we are is.
	var i1, j1 int // Offsets for second (middle) corner of simplex in (i,j) coords
	if x0 > y0 {
		i1 = 1
		j1 = 0 // lower triangle, XY order: (0,0)->(1,0)->(1,1)
	} else {
		i1 = 0
		j1 = 1 // upper triangle, YX order: (0,0)->(0,1)->(1,1)
	}

	// A step of (1,0) in (i,j) means a step of (1-c,-c) in (x,y), and
	// a step of (0,1) in (i,j) means a step of (-c,1-c) in (x,y), where
	// c = (3-sqrt(3))/6
	x1 := x0 - float64(i1) + s.g2 // Offsets for middle corner in (x,y) unskewed coords
	y1 := y0 - float64(j1) + s.g2
	x2 := x0 - 1.0 + 2.0*s.g2 // Offsets for last corner in (x,y) unskewed coords
	y2 := y0 - 1.0 + 2.0*s.g2

	// Work out the hashed gradient indices of the three simplex corners
	ii := i & 255
	jj := j & 255
	gi0 := s.permMod12[ii+int(s.perm[jj])]
	gi1 := s.permMod12[ii+i1+int(s.perm[jj+j1])]
	gi2 := s.permMod12[ii+1+int(s.perm[jj+1])]

	// Calculate the contribution from the three corners
	t0 := 0.5 - x0*x0 - y0*y0
	if t0 < 0 {
		n0 = 0.0
	} else {
		t0 *= t0
		n0 = t0 * t0 * s.mag[gi0] * s.dot2D(s.gradients[gi0], x0, y0) // (x,y) of grad3 used for 2D gradient
	}
	t1 := 0.5 - x1*x1 - y1*y1
	if t1 < 0 {
		n1 = 0.0
	} else {
		t1 *= t1
		n1 = t1 * t1 * s.mag[gi1] * s.dot2D(s.gradients[gi1], x1, y1)
	}
	t2 := 0.5 - x2*x2 - y2*y2
	if t2 < 0 {
		n2 = 0.0
	} else {
		t2 *= t2
		n2 = t2 * t2 * s.mag[gi2] * s.dot2D(s.gradients[gi2], x2, y2)
	}

	// Add contributions from each corner to get the final noise value.
	// The result is scaled to return values in the interval [-1,1].
	return 70.0 * (n0 + n1 + n2)
}

// gradient is used by the noise generator.
// Each gradient is a direction vector to an adjacent height point.
type gradient struct {
	x, y, z float64
}

// Gen3D simplex noise
// From http://www.its.liu.se/~stegu/simplexnoise/SimplexNoise.java
func (s *simplex) Gen3D(xin, yin, zin float64) float64 {
	var n0, n1, n2, n3 float64 // Noise contributions from the four corners

	// Skew the input space to determine which simplex cell we're in
	h := (xin + yin + zin) * s.f3 // Very nice and simple skew factor for 3D
	i := int(s.floor(xin + h))
	j := int(s.floor(yin + h))
	k := int(s.floor(zin + h))
	t := float64(i+j+k) * s.g3
	xs := float64(i) - t // Unskew the cell origin back to (x,y,z) space
	ys := float64(j) - t
	zs := float64(k) - t
	x0 := xin - xs // The x,y,z distances from the cell origin
	y0 := yin - ys
	z0 := zin - zs

	// For the 3D case, the simplex shape is a slightly irregular tetrahedros.
	// Determine which simplex we are is.
	var i1, j1, k1 int // Offsets for second corner of simplex in (i,j,k) coords
	var i2, j2, k2 int // Offsets for third corner of simplex in (i,j,k) coords
	if x0 >= y0 {
		if y0 >= z0 {
			i1, j1, k1, i2, j2, k2 = 1, 0, 0, 1, 1, 0 // X Y Z order
		} else if x0 >= z0 {
			i1, j1, k1, i2, j2, k2 = 1, 0, 0, 1, 0, 1 // X Z Y order
		} else {
			i1, j1, k1, i2, j2, k2 = 0, 0, 1, 1, 0, 1 // Z X Y order
		}
	} else { // x0<y0
		if y0 < z0 {
			i1, j1, k1, i2, j2, k2 = 0, 0, 1, 0, 1, 1 // Z Y X order
		} else if x0 < z0 {
			i1, j1, k1, i2, j2, k2 = 0, 1, 0, 0, 1, 1 // Y Z X order
		} else {
			i1, j1, k1, i2, j2, k2 = 0, 1, 0, 1, 1, 0 // Y X Z order
		}
	}

	// A step of (1,0,0) in (i,j,k) means a step of (1-c,-c,-c) in (x,y,z),
	// a step of (0,1,0) in (i,j,k) means a step of (-c,1-c,-c) in (x,y,z), and
	// a step of (0,0,1) in (i,j,k) means a step of (-c,-c,1-c) in (x,y,z), where
	// c = 1/6.
	x1 := x0 - float64(i1) + s.g3 // Offsets for second corner in (x,y,z) coords
	y1 := y0 - float64(j1) + s.g3
	z1 := z0 - float64(k1) + s.g3
	x2 := x0 - float64(i2) + 2.0*s.g3 // Offsets for third corner in (x,y,z) coords
	y2 := y0 - float64(j2) + 2.0*s.g3
	z2 := z0 - float64(k2) + 2.0*s.g3
	x3 := x0 - 1.0 + 3.0*s.g3 // Offsets for last corner in (x,y,z) coords
	y3 := y0 - 1.0 + 3.0*s.g3
	z3 := z0 - 1.0 + 3.0*s.g3

	// Work out the hashed gradient indices of the four simplex corners
	ii := i & 255
	jj := j & 255
	kk := k & 255
	gi0 := s.permMod12[ii+int(s.perm[jj+int(s.perm[kk])])]
	gi1 := s.permMod12[ii+i1+int(s.perm[jj+j1+int(s.perm[kk+k1])])]
	gi2 := s.permMod12[ii+i2+int(s.perm[jj+j2+int(s.perm[kk+k2])])]
	gi3 := s.permMod12[ii+1+int(s.perm[jj+1+int(s.perm[kk+1])])]

	// Calculate the contribution from the four corners
	t0 := 0.5 - x0*x0 - y0*y0 - z0*z0
	if t0 < 0 {
		n0 = 0.0
	} else {
		t0 *= t0
		n0 = t0 * t0 * s.mag[gi0] * s.dot3D(s.gradients[gi0], x0, y0, z0)
	}
	t1 := 0.5 - x1*x1 - y1*y1 - z1*z1
	if t1 < 0 {
		n1 = 0.0
	} else {
		t1 *= t1
		n1 = t1 * t1 * s.mag[gi1] * s.dot3D(s.gradients[gi1], x1, y1, z1)
	}
	t2 := 0.5 - x2*x2 - y2*y2 - z2*z2
	if t2 < 0 {
		n2 = 0.0
	} else {
		t2 *= t2
		n2 = t2 * t2 * s.mag[gi2] * s.dot3D(s.gradients[gi2], x2, y2, z2)
	}
	t3 := 0.5 - x3*x3 - y3*y3 - z3*z3
	if t3 < 0 {
		n3 = 0.0
	} else {
		t3 *= t3
		n3 = t3 * t3 * s.mag[gi3] * s.dot3D(s.gradients[gi3], x3, y3, z3)
	}

	// Add contributions from each corner to get the final noise value.
	// The result is scaled to stay just inside [-1,1]
	return 32.0 * (n0 + n1 + n2 + n3)
}
