// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package al

import (
	"fmt"
)

// Dump shows which OpenAL functions have been bound to an
// underlying implementation and which haven't.  This is not a guarantee
// that the bound functionality will work, but is an indication of what
// is supported on the current platform.
// Bindings can be dumped even without an active OpenAL context.
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
