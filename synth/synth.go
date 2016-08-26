// Copyright © 2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// Package synth is used to procedurally generate textures and models.
// Synthesized data is unique to a supplied random seed so the generated items
// can be different or identical based on need. They can also be as
// detailed or large as desired, mindfull that it is the clients responsibility
// to manage the memory required for large amounts of generated data.
//
// Package synth is provided as part of the vu (virtual universe) 3D engine.
package synth

// FUTURE: Ideas for possible future improvements:
// Procedural content generation:
//    http://unigine.com/articles/130605-procedural-content-generation/
//    http://unigine.com/articles/131016-procedural-content-generation2/
//    https://en.wikipedia.org/wiki/No_Man%27s_Sky  **NMSky rules!
//    https://en.wikipedia.org/wiki/Superformula
// Other methods of generating random landscapes:
//    http://www.lighthouse3d.com/opengl/terrain/index.php3?fault
//    http://www.lighthouse3d.com/opengl/terrain/index.php3?circles
