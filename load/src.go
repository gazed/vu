// Copyright © 2013-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package load

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

// SrcData is a slice of linefeed terminated strings used to load
// text based files.
type SrcData []string

// Src loads text based data into SrcData.
// Each line terminated with a single linefeed (needed for shader source).
// The Reader r is expected to be opened and closed by the caller.
// A successful import returns a new slice of strings in SrcData.
func Src(r io.Reader) (d SrcData, err error) {
	characters, berr := ioutil.ReadAll(bufio.NewReader(r))
	if berr != nil {
		return nil, fmt.Errorf("Load source %s\n", berr)
	}
	d = strings.Split(string(characters), "\n")

	// shader source lines must be terminated with a linefeed in order to compile.
	for cnt := range d {
		d[cnt] = strings.TrimSpace(d[cnt]) + "\n"
	}
	d = d[0 : len(d)-1] // remove extraneous last line from split.
	return d, nil
}
