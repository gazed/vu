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
	err := shd.Load("basic", loc)
	expect := []string{
		"in  vec4 ex_Color;\n",
		"out vec4 out_Color;\n",
		"\n",
		"void main(void)\n",
		"{\n",
		"out_Color = ex_Color;\n",
		"}\n"}
	got, want := fmt.Sprintf("%s", shd.Fsh), fmt.Sprintf("%s", expect)
	if err != nil || got != want {
		t.Errorf(format, got, want)
	}
}

// Dictate how errors get printed.
const format = "\ngot\n%s\nwanted\n%s"
