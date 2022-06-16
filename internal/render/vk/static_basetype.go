package vk

// translatePublic_Bool32 is a type conversion function for special handling of Bool32 to bool. It is associated with the Bool32 type through exceptions.json
func translatePublic_Bool32(val Bool32) bool {
	return val != Bool32(FALSE)
}

// translateInternal_Bool32 is a type conversion function for special handling of bool to Bool32. It is associated with the Bool32 type through exceptions.json
func translateInternal_Bool32(val bool) Bool32 {
	if val {
		return Bool32(TRUE)
	} else {
		return Bool32(FALSE)
	}
}
