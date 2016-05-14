// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package load

import (
	"bufio"
	"io"
	"io/ioutil"
	"strings"
)

// txt loads the fragment or vertex shader source into a slice of strings.
// Each line must be terminated with a single linefeed. Nil is returned
// if the expected source file is not found on disk.
func (l *loader) txt(fileName string) (lines []string, err error) {
	var file io.ReadCloser
	if file, err = l.getResource(l.dir[src], fileName); err == nil {
		defer file.Close()
		characters, _ := ioutil.ReadAll(bufio.NewReader(file))
		lines = strings.Split(string(characters), "\n")

		// shader source lines must be terminated with a linefeed in order to compile.
		for cnt := range lines {
			lines[cnt] = strings.TrimSpace(lines[cnt]) + "\n"
		}

		// remove extraneous last line from split.
		lines = lines[0 : len(lines)-1]
	}
	return
}
