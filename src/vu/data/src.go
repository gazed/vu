// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package data

import (
	"bufio"
	"io/ioutil"
	"strings"
)

// src loads the fragment or vertex shader source into a slice of strings.
// Each line must be properly terminated with a single linefeed. Nil is returned
// if the expected source file is not found on disk.
func (l loader) src(directory, name string) (source []string) {
	if file, err := l.getResource(directory, name); err == nil {
		defer file.Close()
		characters, _ := ioutil.ReadAll(bufio.NewReader(file))
		source = strings.Split(string(characters), "\n")

		// source lines must be terminated with a linefeed in order to compile.
		for cnt, _ := range source {
			source[cnt] = strings.TrimSpace(source[cnt]) + "\n"
		}

		// remove extraneous last line from split.
		source = source[0 : len(source)-1]
	}
	return
}
