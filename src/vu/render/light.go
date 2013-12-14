// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package render

// FUTURE finish (start?) lighting. There much more information to add. For example:
//     * more light attributes like the type of light.
//     * handling many lights.
//     * loading lights as data resources.
// See: http://gamedev.tutsplus.com/articles/glossary/forward-rendering-vs-deferred-rendering/

import (
	"vu/data"
)

// Light has a position and a colour.
type Light struct {
	X, Y, Z float32  // Light position.
	Ld      data.Rgb // Light colour.
}
