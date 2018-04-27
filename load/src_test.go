// Copyright © 2013-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package load

import (
	"fmt"
	"testing"
)

// Uses vu/eg resource directories.
func TestLoadSource(t *testing.T) {
	shd := &ShdData{}
	loc := NewLocator().Dir("VSH", "../eg/source").Dir("FSH", "../eg/source")
	err := shd.Load("tuv", loc)
	expect := []string{
		"in      vec2      v_t;     // interpolated texture coordinates.\n",
		"uniform sampler2D uv;      // texture sampler.\n",
		"out     vec4      f_color; // final fragment color.\n",
		"\n",
		"void main() {\n",
		"f_color = texture(uv, v_t);\n",
		"}\n",
	}
	got, want := fmt.Sprintf("%s", shd.Fsh), fmt.Sprintf("%s", expect)
	if err != nil || got != want {
		t.Errorf(format, got, want)
	}
}

// Dictate how errors get printed.
const format = "\ngot\n%s\nwanted\n%s"
