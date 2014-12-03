// Copyright © 2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package land

import (
	"math"
	"math/rand"
	"time"
)

// noise is a perlin noise generator algorithm. It's purpose is to help
// create random world maps. The random seed can be saved in order to
// recreate a given map at a later date.
//
// From and thanks to:
//   http://staffwww.itn.liu.se/~stegu/simplexnoise/simplexnoise.pdf
//   http://www.itn.liu.se/~stegu/simplexnoise/SimplexNoise.java
type noise struct {
	F2        float64     // skewing and unskewing factors...
	G2        float64     // ... for 2 dimensions.
	seed      int64       // seed that is unique for a given map.
	pseudo    []byte      // pseudo randomly ordered numbers 0-255
	perm512   []byte      // 512 pseudo random numbers from 0-255
	permMod12 []byte      // pseudo random numbers mod 12
	gradients []*gradient // slopes to adjacent points.
	random    *rand.Rand  // random number generator.
}

// newNoise generates random noise based on a seed. Use 0 for
// the seed to create a new random map. Use a previous seed to get a
// previously generated map.
func newNoise(seed int64) *noise {
	n := &noise{}
	n.seed = seed
	if n.seed == 0 {
		n.seed = time.Now().UnixNano()
	}
	n.random = rand.New(rand.NewSource(n.seed))
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
		index := n.random.Intn(len(pbase))
		n.pseudo = append(n.pseudo, pbase[index])
		pbase = append(pbase[:index], pbase[index+1:]...)
	}
	n.gradients = []*gradient{
		&gradient{1, 1, 0}, &gradient{-1, 1, 0}, &gradient{1, -1, 0}, &gradient{-1, -1, 0},
		&gradient{1, 0, 1}, &gradient{-1, 0, 1}, &gradient{1, 0, -1}, &gradient{-1, 0, -1},
		&gradient{0, 1, 1}, &gradient{0, -1, 1}, &gradient{0, 1, -1}, &gradient{0, -1, -1}}
	n.perm512 = make([]byte, 512)
	n.permMod12 = make([]byte, 512)
	for cnt := 0; cnt < 512; cnt++ {
		n.perm512[cnt] = n.pseudo[cnt&255]
		n.permMod12[cnt] = n.perm512[cnt] % 12
	}

	// Skewing and unskewing factors for 2 dimensions
	n.F2 = 0.5 * (math.Sqrt(3.0) - 1.0)
	n.G2 = (3.0 - math.Sqrt(3.0)) / 6.0
	return n
}

// dot2D provides a dot product of the point with the gradient.
func (n *noise) dot2D(g *gradient, x, y float64) float64 {
	return g.x*x + g.y*y
}

// generate 2D simplex noise.
func (n *noise) generate(xin, yin float64) float64 {
	var n0, n1, n2 float64 // Noise contributions from the three corners

	// Skew the input space to determine which simplex cell we're in
	s := (xin + yin) * n.F2 // Hairy factor for 2D
	i := int(math.Floor(xin + s))
	j := int(math.Floor(yin + s))
	t := float64(i+j) * n.G2
	X0 := float64(i) - t // Unskew the cell origin back to (x,y) space
	Y0 := float64(j) - t
	x0 := xin - X0 // The x,y distances from the cell origin
	y0 := yin - Y0

	// For the 2D case, the simplex shape is an equilateral triangle.
	// Determine which simplex we are in.
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
	x1 := x0 - float64(i1) + n.G2 // Offsets for middle corner in (x,y) unskewed coords
	y1 := y0 - float64(j1) + n.G2
	x2 := x0 - 1.0 + 2.0*n.G2 // Offsets for last corner in (x,y) unskewed coords
	y2 := y0 - 1.0 + 2.0*n.G2

	// Work out the hashed gradient indices of the three simplex corners
	ii := i & 255
	jj := j & 255
	gi0 := n.permMod12[ii+int(n.perm512[jj])]
	gi1 := n.permMod12[ii+i1+int(n.perm512[jj+j1])]
	gi2 := n.permMod12[ii+1+int(n.perm512[jj+1])]

	// Calculate the contribution from the three corners
	t0 := 0.5 - x0*x0 - y0*y0
	if t0 < 0 {
		n0 = 0.0
	} else {
		t0 *= t0
		n0 = t0 * t0 * n.dot2D(n.gradients[gi0], x0, y0) // (x,y) of grad3 used for 2D gradient
	}
	t1 := 0.5 - x1*x1 - y1*y1
	if t1 < 0 {
		n1 = 0.0
	} else {
		t1 *= t1
		n1 = t1 * t1 * n.dot2D(n.gradients[gi1], x1, y1)
	}
	t2 := 0.5 - x2*x2 - y2*y2
	if t2 < 0 {
		n2 = 0.0
	} else {
		t2 *= t2
		n2 = t2 * t2 * n.dot2D(n.gradients[gi2], x2, y2)
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
