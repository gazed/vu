// Copyright © 2013-2017 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"github.com/gazed/vu/audio/al"
)

// Dump the openal binding information. This is a basic audio package test that
// checks if the underlying OpenAL functions are available. Columns of function
// names marked [+]:available or [ ]:missing will be written the the console.
//
// CONTROLS: NA
func da() {
	al.Dump()
}
