// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package render

// TODO finish (start?) lighting. There much more information to add.
// For example:
//     * more light attributes like the type of light.
//     * handling many lights.
//     * loading lights as data resources.

import (
	"vu/data"
)

// Light has a postion and a colour (and eventually lots more once lighting
// is properly designed and implemented).
type Light struct {
	X, Y, Z float32  // Light position.
	Ld      data.Rgb // Light colour.
}
