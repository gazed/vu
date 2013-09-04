// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package main

import (
	"vu/audio/al"
)

// Dump the openal binding information.  This is a basic audio package test that checks
// if the underlying OpenAL functions are available. A column of function names with
// [+] signs beside will be written the the console if everything is ok.
func da() {
	al.Dump()
}
