// SPDX-FileCopyrightText : © 2026 Galvanized Logic Inc.
// SPDX-License-Identifier: MIT

package vk

import (
	"bytes"
	"unsafe"
)

// ============================================================================
// The functions in this file are Go versions of C preprocessor macros from vk.xml. Each function here is associated with its
// Vulkan name through exceptions.json

func deprecatedMakeVersion(major, minor, patch uint32) uint32 {
	return major<<22 | minor<<12 | patch
}
func deprecatedVersionMajor(version uint32) uint32 {
	return version >> 22
}
func deprecatedVersionMinor(version uint32) uint32 {
	return version >> 12 & 0x3FF
}
func deprecatedVersionPatch(version uint32) uint32 {
	return version & 0xFFF
}

func makeVersion(major, minor, patch uint32) uint32 {
	return major<<22 | minor<<12 | patch
}

func makeApiVersion(variant, major, minor, patch uint32) uint32 {
	return variant<<29 | major<<22 | minor<<12 | patch
}

func apiVersionVariant(version uint32) uint32 {
	return version >> 29
}
func apiVersionMajor(version uint32) uint32 {
	return version >> 22 & 0x7F
}
func apiVersionMinor(version uint32) uint32 {
	return version >> 12 & 0x3FF
}
func apiVersionPatch(version uint32) uint32 {
	return version & 0xFFF
}

// max is an internal utility function, used in processing struct member slice/array lengths
func max(nums ...int) int {
	rval := 0
	for _, v := range nums {
		if rval < v {
			rval = v
		}
	}
	return rval
}

// nullTermBytesToString converts c string arrays to golang strings.
func nullTermBytesToString(b []byte) string {
	n := bytes.IndexByte(b, 0)
	return string(b[:n])
}

// LoadVulkan loads the vulkan library. This must be called once on startup.
// Override the default Vulkan library name by setting libName to a non-empty string.
// For example, if you want to enable the validation layers, those layers are only available in the
// Vulkan SDK libary. go-vk passes the name to the host operating system's library opening/search method,
// so you must provide a relative or absolute path if your Vulkan library is not in the default search
// path for the platform.
func LoadVulkan(overrideLibName string) error {

	// non-empty strings override the default name.
	return loadLibrary(overrideLibName)
}

// ============================================================================
// sl is a Go slice header (or rather, it matches a slice header byte-for-byte).
type sl struct {
	addr uintptr
	len  int
	cap  int
}

// MemCopySlice provides an abstracted memory copy function, intended for use with mapped memory ranges. Note that the
// destination is passed as an unsafe.Pointer and this function cannot determine if you have allocated enough
// memory at that location. You are responsible for requesting enough memory from Vulkan! Unexpected behavior or crashes
// are very possible if this function is misused or abused.
func MemCopySlice[T any](dest unsafe.Pointer, src []T) {
	if len(src) == 0 {
		return
	}
	bytes := len(src) * int(unsafe.Sizeof(src[0]))
	sl_src := *(*[]byte)(unsafe.Pointer(&sl{uintptr(unsafe.Pointer(&src[0])), bytes, bytes}))
	sl_dest := *(*[]byte)(unsafe.Pointer(&sl{uintptr(dest), bytes, bytes}))
	copy(sl_dest, sl_src)
}
