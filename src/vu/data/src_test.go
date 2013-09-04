// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package data

import (
	"fmt"
	"testing"
)

// Dictate how errors get printed.
const format = "\ngot\n%s\nwanted\n%s"

// Check that fragment and vertex shaders can be read in.
func TestLoadSource(t *testing.T) {
	load := &loader{}
	source := load.src("../eg/shaders", "basic.fsh")
	expect := []string{
		"#version 150\n",
		"\n",
		"in vec4 ex_Color;\n",
		"out vec4 out_Color;\n",
		"\n",
		"void main(void)\n",
		"{\n",
		"out_Color = ex_Color;\n",
		"}\n"}
	got, want := fmt.Sprintf("%s", source), fmt.Sprintf("%s", expect)
	if got != want {
		t.Errorf(format, got, want)
	}
}
