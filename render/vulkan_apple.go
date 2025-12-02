// SPDX-FileCopyrightText : © 2025 Galvanized Logic Inc.
// SPDX-License-Identifier: BSD-2-Clause

//go:build darwin || ios

package render

import (
	"fmt"

	"github.com/gazed/vu/device"
	"github.com/gazed/vu/internal/render/vk"
)

// vulkan_apple.go contains the Vulkan apple (macos, ios) specific OS extensions
// ie: VK_EXT_metal_surface supports both macos and ios.

// instanceExtensions return extensions needed for the VkInstance.
func (vr *vulkanRenderer) instanceExtensions() []string {
	return []string{
		vk.KHR_SURFACE_EXTENSION_NAME,
		vk.EXT_METAL_SURFACE_EXTENSION_NAME, // "VK_EXT_metal_surface"
	}
}

// createSurface associates a vulkan instance with an apple device display.
// It gets display surface information from the device layer.
func (vr *vulkanRenderer) createSurface() (err error) {
	viewPointer, err := device.GetRenderSurfaceInfo(vr.osdev)
	if viewPointer == nil || err != nil {
		return fmt.Errorf("device.GetRenderSurfaceInfo failed %w", err)
	}
	ci := vk.MetalSurfaceCreateInfoEXT{
		PLayer: (*vk.CAMetalLayer)(viewPointer),
	}
	vr.surface, err = vk.CreateMetalSurfaceEXT(vr.instance, &ci, nil)
	return err
}
