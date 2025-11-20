//go:build windows

package vk

import "golang.org/x/sys/windows"

func sys_stringToBytePointer(s string) *byte {
	p, _ := windows.BytePtrFromString(s)
	return p
}

var dlHandle *windows.LazyDLL

type vkCommand struct {
	protoName string
	argCount  int
	hasReturn bool
	fnHandle  *windows.LazyProc
}

func loadLibrary(overrideLibName string) error {
	libName := "vulkan-1.dll"
	if overrideLibName != "" {
		libName = overrideLibName
	}
	dlHandle = windows.NewLazyDLL(libName)
	return nil
}
