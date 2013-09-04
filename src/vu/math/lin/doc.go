// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

// Package math/lin provides a linear math library that includes vectors,
// matricies, and quaternions. Linear math operations are used extensively
// in 3D applications for describing virtual worlds, repositioning objects,
// and performing physics simulations of movement and collisions.
//
// While most heavy math is expected to be done on the GPU, the CPU 3D math
// is still likely to be done within rendering loops where performance is key
// Some general guidelines, verified with benchmarks, can be seen throughout
// the library.
//     - minimize function calls.
//     - use pointers to structures
//     - avoid instantiating new structures where possible.
//     - prefer multiply over divide
//
// The long term plan is that there will evenutally be a standard golang linear
// math library and this home grown one can be ditched.
package lin

// TODO look at incorporating GPGPU math calls into this package.
