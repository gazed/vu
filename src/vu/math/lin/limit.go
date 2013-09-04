// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package lin

import "math"

// Deal with single floating point imprecision.
// Anything with less precision than this is considered zero.
const TOLERANCE float32 = 0.000001

// PI/180 is commonly needed to convert degrees to radians: radians = degrees * PI/180.
const PI_OVER_180 = math.Pi / 180

// IsEqual checks that two floating point numbers are essentially the same.
func IsEqual(f1, f2 float32) bool {
	diff := f1 - f2
	return diff < TOLERANCE && diff > -TOLERANCE
}

// IsZero checks if the floating point number is essentially zero.
func IsZero(value float32) bool {
	return math.Abs(float64(value)) < float64(TOLERANCE)
}

// IsOne checks if the floating point number is essentially one or
// negative one.
func IsOne(value float32) bool {
	return math.Abs(math.Abs(float64(value))-1) < float64(TOLERANCE)
}

// Lerp provides linear interpolation between floats f and fend where
// fraction is expected to be between 0 and 1.
func Lerp(f, fend, fraction float32) float32 {
	return (fend-f)*fraction + f
}
