//go:build windows

package vk

import "golang.org/x/sys/windows"

func sys_stringToBytePointer(s string) *byte {
	p, _ := windows.BytePtrFromString(s)
	return p
}
