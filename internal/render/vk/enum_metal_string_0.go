// Code generated by "stringer -output=enum_metal_string_0.go -type=ExportMetalObjectTypeFlagBitsEXT"; DO NOT EDIT.

package vk

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[EXPORT_METAL_OBJECT_TYPE_METAL_DEVICE_BIT_EXT-1]
	_ = x[EXPORT_METAL_OBJECT_TYPE_METAL_COMMAND_QUEUE_BIT_EXT-2]
	_ = x[EXPORT_METAL_OBJECT_TYPE_METAL_BUFFER_BIT_EXT-4]
	_ = x[EXPORT_METAL_OBJECT_TYPE_METAL_TEXTURE_BIT_EXT-8]
	_ = x[EXPORT_METAL_OBJECT_TYPE_METAL_IOSURFACE_BIT_EXT-16]
	_ = x[EXPORT_METAL_OBJECT_TYPE_METAL_SHARED_EVENT_BIT_EXT-32]
}

const (
	_ExportMetalObjectTypeFlagBitsEXT_name_0 = "EXPORT_METAL_OBJECT_TYPE_METAL_DEVICE_BIT_EXTEXPORT_METAL_OBJECT_TYPE_METAL_COMMAND_QUEUE_BIT_EXT"
	_ExportMetalObjectTypeFlagBitsEXT_name_1 = "EXPORT_METAL_OBJECT_TYPE_METAL_BUFFER_BIT_EXT"
	_ExportMetalObjectTypeFlagBitsEXT_name_2 = "EXPORT_METAL_OBJECT_TYPE_METAL_TEXTURE_BIT_EXT"
	_ExportMetalObjectTypeFlagBitsEXT_name_3 = "EXPORT_METAL_OBJECT_TYPE_METAL_IOSURFACE_BIT_EXT"
	_ExportMetalObjectTypeFlagBitsEXT_name_4 = "EXPORT_METAL_OBJECT_TYPE_METAL_SHARED_EVENT_BIT_EXT"
)

var (
	_ExportMetalObjectTypeFlagBitsEXT_index_0 = [...]uint8{0, 45, 97}
)

func (i ExportMetalObjectTypeFlagBitsEXT) String() string {
	switch {
	case 1 <= i && i <= 2:
		i -= 1
		return _ExportMetalObjectTypeFlagBitsEXT_name_0[_ExportMetalObjectTypeFlagBitsEXT_index_0[i]:_ExportMetalObjectTypeFlagBitsEXT_index_0[i+1]]
	case i == 4:
		return _ExportMetalObjectTypeFlagBitsEXT_name_1
	case i == 8:
		return _ExportMetalObjectTypeFlagBitsEXT_name_2
	case i == 16:
		return _ExportMetalObjectTypeFlagBitsEXT_name_3
	case i == 32:
		return _ExportMetalObjectTypeFlagBitsEXT_name_4
	default:
		return "ExportMetalObjectTypeFlagBitsEXT(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
