// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package gl

import (
	"fmt"
)

// Dump shows which OpenGL functions have been bound to an underlying
// implementation.  This is not a guarantee that the bound functionality will
// work, but is an indication of what is supported on the current platform.
// Bindings can be dumped without an active OpenGL context.
//
// Bound functions are indicated with [+] and unbound with [ ].
func Dump() {
	Init()
	report := BindingReport()
	size := len(report)

	// print the report in columns.
	format := "%-35.35s %-35.35s\n"
	for cnt := 0; cnt < size/2+1; cnt++ {
		one := ""
		if cnt < size/2 {
			one = report[cnt]
		}
		two := ""
		if size/2+cnt < size-1 {
			two = report[size/2+cnt]
		}
		fmt.Printf(format, one, two)
	}
}
