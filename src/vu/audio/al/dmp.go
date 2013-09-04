// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package al

import (
	"fmt"
	"strings"
)

// Dump shows which OpenAL functions have been bound to an
// underlying implementation and which haven't.  This is not a guarantee
// that the bound functionality will work, but is an indication of what
// is supported on the current platform. A OpenAL context is not needed
// in order to dump the bindings.
//
// The bound functions are listed on the left with [+] and the unbound
// functions are listed on the right with [-]
func Dump() {
	Init()
	report := BindingReport()

	// process the basic report for nicer output
	var bound = []string{}
	var unbound = []string{}
	var group = []string{}
	var groupAllBound = true
	for _, groupOrFn := range report {

		// found a new group.
		if groupOrFn[0:1] == "A" {

			// if there was a previous group, then...
			if len(group) > 0 {

				// ... put the previous group into one of bound/unbound lists.
				if groupAllBound {
					bound = append(bound, group...)
				} else {
					unbound = append(unbound, group...)
				}
			}

			// start the new group
			group = append([]string{}, groupOrFn)
			groupAllBound = true
		} else {

			// append a group function to the group.
			group = append(group, groupOrFn)
			if strings.Contains(groupOrFn, "[-]") {
				groupAllBound = false
			}
		}
	}

	// append the last group.
	if groupAllBound {
		bound = append(bound, group...)
	} else {
		unbound = append(unbound, group...)
	}

	// now print the report in columns.
	format := "%-35.35s %-35.35s\n"
	fmt.Printf(format, "Bound", "Unbound")

	// print based on the largest of the two columns.
	max := len(bound)
	if max < len(unbound) {
		max = len(unbound)
	}
	for cnt := 0; cnt < max; cnt++ {
		bnd := ""
		if cnt < len(bound) {
			bnd = bound[cnt]
		}
		ubnd := ""
		if cnt < len(unbound) {
			ubnd = unbound[cnt]
		}
		fmt.Printf(format, bnd, ubnd)
	}
}
