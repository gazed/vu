// Code generated by go-vk from vk.xml at 2024-09-26 06:25:38.462104 -0400 EDT m=+0.680633101. DO NOT EDIT.

package vk

const (
	LOD_CLAMP_NONE float32 = 1000.0
)

const (
	FALSE                         uint32 = 0
	TRUE                          uint32 = 1
	LUID_SIZE                     uint32 = 8
	MAX_GLOBAL_PRIORITY_SIZE_EXT  uint32 = 16
	UUID_SIZE                     uint32 = 16
	MAX_MEMORY_HEAPS              uint32 = 16 // The maximum number of unique memory heaps, each of which supporting 1 or more memory types
	MAX_MEMORY_TYPES              uint32 = 32
	MAX_DEVICE_GROUP_SIZE         uint32 = 32
	MAX_DRIVER_INFO_SIZE          uint32 = 256
	MAX_DESCRIPTION_SIZE          uint32 = 256
	MAX_PHYSICAL_DEVICE_NAME_SIZE uint32 = 256
	MAX_DRIVER_NAME_SIZE          uint32 = 256
	MAX_EXTENSION_NAME_SIZE       uint32 = 256
	REMAINING_ARRAY_LAYERS        uint32 = ^uint32(0)
	REMAINING_MIP_LEVELS          uint32 = ^uint32(0)
	ATTACHMENT_UNUSED             uint32 = ^uint32(0)
	SHADER_UNUSED_KHR             uint32 = ^uint32(0)
	SUBPASS_EXTERNAL              uint32 = ^uint32(0)
	QUEUE_FAMILY_IGNORED          uint32 = ^uint32(0)
	QUEUE_FAMILY_EXTERNAL         uint32 = ^uint32(1)
	QUEUE_FAMILY_FOREIGN_EXT      uint32 = ^uint32(2)
)

const (
	WHOLE_SIZE uint64 = ^uint64(0)
)

// Extension names and versions
const (
	LUID_SIZE_KHR             uint32 = LUID_SIZE
	MAX_DEVICE_GROUP_SIZE_KHR uint32 = MAX_DEVICE_GROUP_SIZE
	MAX_DRIVER_INFO_SIZE_KHR  uint32 = MAX_DRIVER_INFO_SIZE
	MAX_DRIVER_NAME_SIZE_KHR  uint32 = MAX_DRIVER_NAME_SIZE
	QUEUE_FAMILY_EXTERNAL_KHR uint32 = QUEUE_FAMILY_EXTERNAL
	SHADER_UNUSED_NV          uint32 = SHADER_UNUSED_KHR
)
