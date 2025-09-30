package vk

import "unsafe"

// sl is a Go slice header (or rather, it matches a slice header byte-for-byte).
type sl struct {
	addr uintptr
	len  int
	cap  int
}

/* MemCopySlice provides an abstracted memory copy function, intended for use with mapped memory ranges. Note that the
destination is passed as an unsafe.Pointer and this function cannot determine if you have allocated enough
memory at that location. You are responsible for requesting enough memory from Vulkan! Unexpected behavior or crashes
are very possible if this function is misused or abused.
*/
func MemCopySlice[T any](dest unsafe.Pointer, src []T) {
	if len(src) == 0 {
		return
	}

	bytes := len(src) * int(unsafe.Sizeof(src[0]))

	sl_src := *(*[]byte)(unsafe.Pointer(&sl{uintptr(unsafe.Pointer(&src[0])), bytes, bytes}))
	sl_dest := *(*[]byte)(unsafe.Pointer(&sl{uintptr(dest), bytes, bytes}))

	copy(sl_dest, sl_src)
}

/* MemCopyObj provides an abstracted memory copy function for a single piece of data (a struct or primitive type),
intended for use with mapped memory ranges. Note that the destination is passed as an unsafe.Pointer and this function
cannot determine if you have allocated enough memory at that location. You are responsible for requesting enough memory
from Vulkan! Unexpected behavior or crashes are very possible if this function is misused or abused.

NOTE: If you pass a slice to this function, the slice header will be copied, not the contents! Use [MemCopySlice] instead.
*/
func MemCopyObj[T any](dest unsafe.Pointer, src *T) {
	bytes := int(unsafe.Sizeof(*src))

	sl_src := *(*[]byte)(unsafe.Pointer(&sl{uintptr(unsafe.Pointer(src)), bytes, bytes}))
	sl_dest := *(*[]byte)(unsafe.Pointer(&(sl{uintptr(dest), bytes, bytes})))

	copy(sl_dest, sl_src)
}

/* MemCopy is the closest to C's memcpy...you provide two pointers and a number of bytes to copy, and it will move the
data around. Using [MemCopySlice] or [MemCopyObj] instead is highly recommended. There are no guardrails on this function!
*/
func MemCopy(dest, src unsafe.Pointer, len int) {
	sl_src := *(*[]byte)(unsafe.Pointer(&sl{uintptr(src), len, len}))
	sl_dest := *(*[]byte)(unsafe.Pointer(&sl{uintptr(dest), len, len}))

	copy(sl_dest, sl_src)
}
