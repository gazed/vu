//go:build darwin || linux

package vk

import (
	"golang.org/x/sys/unix"
)

func sys_stringToBytePointer(s string) *byte {
	p, _ := unix.BytePtrFromString(s)
	return p
}
