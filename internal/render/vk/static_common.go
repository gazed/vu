package vk

import (
	"bytes"
	"fmt"
	"unsafe"
)

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
	// hacked to avoid the stringer tool.
	return fmt.Sprintf("%d", r) // r.String()
}

var overrideLibName string

// OverrideDefaultVulkanLibrary allows you to set a specific Vulkan library name to be used in your program. For
// example, if you want to enable the validation layers, those layers are only available in the Vulkan SDK libary. go-vk
// passes the name to the host operating system's library opening/search method, so you must provide a relative or
// absolute path if your Vulkan library is not in the default search path for the platform.
func OverrideDefaultVulkanLibrary(nameOrPath string) {
	overrideLibName = nameOrPath
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
