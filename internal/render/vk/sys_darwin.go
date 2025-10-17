//go:build darwin && !ios

package vk

import "fmt"
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

func loadLibrary(overrideLibName string) error {
	libName := "libMoltenVK.dylib"
	if overrideLibName != "" {
		libName = overrideLibName
	}
	cstr := C.CString(libName)
	dlHandle = C.OpenLibrary(cstr)
	C.free(unsafe.Pointer(cstr))

	// report any errors.
	if dlHandle == nil {
		cStr := C.OpenLibraryError()
		str := C.GoString(cStr)
		return fmt.Errorf("LoadVulkan failed: %s", str)
	}
	return nil
}

func sys_stringToBytePointer(s string) *byte {
	p, _ := unix.BytePtrFromString(s)
	return p
}
