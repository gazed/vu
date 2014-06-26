// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package vu

// FUTURE finish (start?) lighting. There much more information to add. For example:
//     * more light attributes like the type of light.
//     * handling many lights.
//     * loading lights as data resources.
// See: http://gamedev.tutsplus.com/articles/glossary/forward-rendering-vs-deferred-rendering/
//      http://antongerdelan.net/opengl/deferredshading.html

// More light resources:
//    http://tomdalling.com/blog/modern-opengl/06-diffuse-point-lighting/
//    http://tomdalling.com/blog/modern-opengl/07-more-lighting-ambient-specular-attenuation-gamma/
//    http://www.learnopengles.com/android-lesson-two-ambient-and-diffuse-lighting/

// Lighting is almost entirely based on shaders...
//    http://www.gamedev.net/page/resources/_/technical/opengl/the-basics-of-glsl-40-shaders-r2861
//    http://www.swiftless.com/glsltuts.html

// light has a position and a colour.
type light struct {
	x, y, z float32 // Light position.
	ld      rgb     // Light colour.
}

// light
// ===========================================================================
// rgb

// rgb holds colour information where each field is expected to contain
// a value from 0.0 to 1.0. A value of 0 means none of that colour while a value
// of 1.0 means as much as possible of that colour. For example:
//     black := &Rgb{0, 0, 0}     white := &Rgb{1, 1, 1}
//     red   := &Rgb{1, 0, 0}     gray  := &Rgb{0.5, 0.5, 0.5}
type rgb struct {
	R float32 // Red.
	G float32 // Green.
	B float32 // Blue.
}
