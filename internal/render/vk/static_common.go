package vk

import (
	"bytes"
	"runtime"
	"unsafe"
)

// #include <stdlib.h>
// #include "dlload.h"
import "C"

// Vulkanizer allows conversion from go-vk style structs to Vulkan-native structs. This
// includes setting the structure type flag, converting slices to pointers, etc.
type Vulkanizer interface {
	Vulkanize() unsafe.Pointer
}

// Goifier converts Vulkan-native structs back into go-vk style structs
type Goifier interface {
	Goify() Vulkanizer
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

// Error implements the error interface
// TODO: A way for commands to indicate if the Result code is an error for that command, or an unexpected return value?
func (r Result) Error() string {
	return r.String()
}

type vkCommand struct {
	protoName string
	argCount  int
	hasReturn bool
	fnHandle  unsafe.Pointer
}

var dlHandle unsafe.Pointer

var overrideLibName string

// OverrideDefaultVulkanLibrary allows you to set a specific Vulkan library name to be used in your program. For
// example, if you want to enable the validation layers, those layers are only available in the Vulkan SDK libary. go-vk
// passes the name to the host operating system's library opening/search method, so you must provide a relative or
// absolute path if your Vulkan library is not in the default search path for the platform.
func OverrideDefaultVulkanLibrary(nameOrPath string) {
	overrideLibName = nameOrPath
}

func execTrampoline(cmd *vkCommand, args ...uintptr) uintptr {
	if dlHandle == nil {
		var libName string
		switch runtime.GOOS {
		case "windows":
			libName = "vulkan-1.dll"
		case "darwin":
			// TODO: Running on Mac/Darwin is tested only to the point of creating and
			// destroying a Vulkan instance.
			libName = "libMoltenVK.dylib"
		case "linux":
			// TODO: Running on Linux is tested only to the point of creating and
			// destroying a Vulkan instance.
			libName = "libvulkan.so"
		default:
			panic("Unsupported GOOS at OpenLibrary: " + runtime.GOOS)
		}

		if overrideLibName != "" {
			libName = overrideLibName
		}

		cstr := C.CString(libName)
		dlHandle = C.OpenLibrary(cstr)
		C.free(unsafe.Pointer(cstr))
	}

	// cmd := lazyCommands[commandKey]
	if cmd.fnHandle == nil {
		cmd.fnHandle = C.SymbolFromName(dlHandle, unsafe.Pointer(sys_stringToBytePointer(cmd.protoName)))
		// lazyCommands[commandKey] = cmd
	}

	if len(args) != cmd.argCount {
		panic("Wrong number of arguments passed for cmd " + cmd.protoName)
	}

	var result C.uintptr_t

	switch cmd.argCount {
	case 1:
		result = C.Trampoline3(cmd.fnHandle, C.uintptr_t(args[0]), 0, 0)
	case 2:
		result = C.Trampoline3(cmd.fnHandle, C.uintptr_t(args[0]), C.uintptr_t(args[1]), 0)
	case 3:
		result = C.Trampoline3(cmd.fnHandle, C.uintptr_t(args[0]), C.uintptr_t(args[1]), C.uintptr_t(args[2]))
	case 4:
		result = C.Trampoline6(cmd.fnHandle, C.uintptr_t(args[0]), C.uintptr_t(args[1]), C.uintptr_t(args[2]), C.uintptr_t(args[3]), 0, 0)
	case 5:
		result = C.Trampoline6(cmd.fnHandle, C.uintptr_t(args[0]), C.uintptr_t(args[1]), C.uintptr_t(args[2]), C.uintptr_t(args[3]), C.uintptr_t(args[4]), 0)
	case 6:
		result = C.Trampoline6(cmd.fnHandle, C.uintptr_t(args[0]), C.uintptr_t(args[1]), C.uintptr_t(args[2]), C.uintptr_t(args[3]), C.uintptr_t(args[4]), C.uintptr_t(args[5]))
	case 7:
		result = C.Trampoline9(cmd.fnHandle, C.uintptr_t(args[0]), C.uintptr_t(args[1]), C.uintptr_t(args[2]), C.uintptr_t(args[3]), C.uintptr_t(args[4]), C.uintptr_t(args[5]), C.uintptr_t(args[6]), 0, 0)
	case 8:
		result = C.Trampoline9(cmd.fnHandle, C.uintptr_t(args[0]), C.uintptr_t(args[1]), C.uintptr_t(args[2]), C.uintptr_t(args[3]), C.uintptr_t(args[4]), C.uintptr_t(args[5]), C.uintptr_t(args[6]), C.uintptr_t(args[7]), 0)
	case 9:
		result = C.Trampoline9(cmd.fnHandle, C.uintptr_t(args[0]), C.uintptr_t(args[1]), C.uintptr_t(args[2]), C.uintptr_t(args[3]), C.uintptr_t(args[4]), C.uintptr_t(args[5]), C.uintptr_t(args[6]), C.uintptr_t(args[7]), C.uintptr_t(args[8]))
	case 10:
		result = C.Trampoline12(cmd.fnHandle, C.uintptr_t(args[0]), C.uintptr_t(args[1]), C.uintptr_t(args[2]), C.uintptr_t(args[3]), C.uintptr_t(args[4]), C.uintptr_t(args[5]), C.uintptr_t(args[6]), C.uintptr_t(args[7]), C.uintptr_t(args[8]), C.uintptr_t(args[9]), 0, 0)
	case 11:
		result = C.Trampoline12(cmd.fnHandle, C.uintptr_t(args[0]), C.uintptr_t(args[1]), C.uintptr_t(args[2]), C.uintptr_t(args[3]), C.uintptr_t(args[4]), C.uintptr_t(args[5]), C.uintptr_t(args[6]), C.uintptr_t(args[7]), C.uintptr_t(args[8]), C.uintptr_t(args[9]), C.uintptr_t(args[10]), 0)
	default:
		// There are no commands with 0 or 12+ arguments as of Vulkan 1.3.204
		panic("Unhandled number of arguments passed for cmd " + cmd.protoName)
	}

	return uintptr(result)
}

func stringToNullTermBytes(s string) *byte {
	b := []byte(s)
	b = append(b, 0)
	return &b[0]
}

func nullTermBytesToString(b []byte) string {
	n := bytes.IndexByte(b, 0)
	return string(b[:n])
}
