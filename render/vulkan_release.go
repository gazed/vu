// Copyright © 2024 Galvanized Logic Inc.

//go:build !debug

package render

// vulkan_release.go ensures vulkan debug utilitlies are not shipped
// with the release build.

// addValidationLayer is excluded from release builds
func (vr *vulkanRenderer) addValidationLayer() error { return nil }
