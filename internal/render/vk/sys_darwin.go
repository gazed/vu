//go:build darwin || ios

package vk

import "golang.org/x/sys/unix"
import "unsafe"

// #include <stdlib.h>
// #include "dlload.h"
import "C"

var dlHandle unsafe.Pointer

type vkCommand struct {
	protoName string
	argCount  int
	hasReturn bool
	fnHandle  unsafe.Pointer
}

// called once automatically on package init.
func init() {
	libName := "libMoltenVK.dylib"
	if overrideLibName != "" {
		libName = overrideLibName
	}
	cstr := C.CString(libName)
	dlHandle = C.OpenLibrary(cstr)
	C.free(unsafe.Pointer(cstr))
}

func sys_stringToBytePointer(s string) *byte {
	p, _ := unix.BytePtrFromString(s)
	return p
}
