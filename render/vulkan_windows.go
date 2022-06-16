// Copyright © 2024 Galvanized Logic Inc.

package render

import (
	"fmt"

	"github.com/gazed/vu/device"
	"github.com/gazed/vu/internal/render/vk"
	"golang.org/x/sys/windows"
)

// vulkan_windows.go contains the Vulkan windows specific OS extensions

// instanceExtensions return extensions needed for the VkInstance.
func (vr *vulkanRenderer) instanceExtensions() []string {
	return []string{
		vk.KHR_SURFACE_EXTENSION_NAME,
		vk.KHR_WIN32_SURFACE_EXTENSION_NAME,
	}
}

// createSurface associates a vulkan instance with a winOS window.
// It gets display surface information from the windows specific method
// in the device layer.
func (vr *vulkanRenderer) createSurface() (err error) {
	hinstance, hwnd, err := device.GetRenderSurfaceInfo(vr.osdev)
	if hinstance == 0 || hwnd == 0 || err != nil {
		return fmt.Errorf("device.GetWindowsSurfaceInfo failed %w", err)
	}
	ci := vk.Win32SurfaceCreateInfoKHR{
		Hinstance: windows.Handle(hinstance),
		Hwnd:      windows.HWND(hwnd),
	}
	vr.surface, err = vk.CreateWin32SurfaceKHR(vr.instance, &ci, nil)
	return err
}
