//go:build darwin || ios

// Copyright © 2025 Galvanized Logic Inc.

package render

import (
	"fmt"

	"github.com/gazed/vu/device"
	"github.com/gazed/vu/internal/render/vk"
)

// vulkan_apple.go contains the Vulkan apple (macos, ios) specific OS extensions

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
	ci := vk.MacOSSurfaceCreateInfoMVK{
		PView: viewPointer, // pointer to MTKView
	}
	vr.surface, err = vk.CreateMacOSSurfaceMVK(vr.instance, &ci, nil)
	return err
}
