// Copyright © 2014-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package synth

import (
	"math"
)

// Form holds the 6 parameters needed to for the superformula to
// create supershapes. Form uses the the superformula to generate
// points in both 2D and 3D. See:
//    https://en.wikipedia.org/wiki/Superformula
//    http://paulbourke.net/geometry/supershape/
type Form struct {
	M  float64 // Angle multiplier
	N1 float64 // Overall exponent.
	N2 float64 // Cos exponent.
	N3 float64 // Sin exponent.
	A  float64 // Cos divisor
	B  float64 // Sin divisor
}

// CircleForm is the default state of a super formula whose
// shape is a 2D circle or 3D sphere with values 0,1,1,1,1,1.
// Not expected to be changed.
var CircleForm = &Form{M: 0, N1: 1, N2: 1, N3: 1, A: 1, B: 1}

// NewForm creates a supershape at its default circle/sphere shape.
func NewForm() *Form {
	return &Form{M: 0, N1: 1, N2: 1, N3: 1, A: 1, B: 1}
}

// SetValues sets the superform to the given values.
func (f *Form) SetValues(m, n1, n2, n3, a, b float64) {
	f.M, f.N1, f.N2, f.N3, f.A, f.B = m, n1, n2, n3, a, b
}

// Set the superform to the given Form fm.
func (f *Form) Set(fm *Form) {
	f.M, f.N1, f.N2, f.N3, f.A, f.B = fm.M, fm.N1, fm.N2, fm.N3, fm.A, fm.B
}

// Radius runs the superformula for the given angle in radians.
func (f *Form) Radius(angle float64) (r float64) {
	t1 := math.Pow(math.Abs(math.Cos(f.M*angle/4)/f.A), f.N2)
	t2 := math.Pow(math.Abs(math.Sin(f.M*angle/4)/f.B), f.N3)
	radius := math.Pow(t1+t2, 1/f.N1)
	if math.Abs(radius) == 0 {
		return 0
	}
	radius = 1 / radius
	return radius
}

// At2D returns the 2D point for a given angle in radians.
func (f *Form) At2D(angle float64) (x, y float64) {
	radius := f.Radius(angle)
	x = radius * math.Cos(angle)
	y = radius * math.Sin(angle)
	return x, y
}

// At3D returns the 3D point for the supplied lat and lon angles in radians.
//   lat is the angle in radians between Pi/2 and -Pi/2.
//   lon is the angle in radians between Pi and -Pi.
func (f *Form) At3D(lat, lon float64) (x, y, z float64) {
	r1 := f.Radius(lon)
	r2 := f.Radius(lat)
	x = r1 * math.Cos(lon) * r2 * math.Cos(lat)
	y = r1 * math.Sin(lon) * r2 * math.Cos(lat)
	z = r2 * math.Sin(lat)
	return x, y, z
}
