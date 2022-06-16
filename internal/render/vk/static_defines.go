package vk

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
