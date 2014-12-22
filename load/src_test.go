// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package load

import (
	"fmt"
	"testing"
)

func TestLoadSource(t *testing.T) {
	load := newLoader().SetDir(src, "../eg/source")
	source, _ := load.Fsh("basic")
	expect := []string{
		"#version 330\n",
		"\n",
		"in  vec4 ex_Color;\n",
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

// Dictate how errors get printed.
const format = "\ngot\n%s\nwanted\n%s"
