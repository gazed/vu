// Copyright © 2024 Galvanized Logic Inc.

//go:build debug

package render

// vulkan_debug.go includes vulkan debug utilitlies when
// building with "-tags debug"

import (
	"fmt"

	"github.com/gazed/vu/internal/render/vk"
)

// addValidationLayers enables vulkan validation in debug builds.
// Adds the validation layer to the list of
//
// updates: vr.apiLayers
func (vr *vulkanRenderer) addValidationLayer() error {
	props, err := vk.EnumerateInstanceLayerProperties()
	if err != nil {
		return fmt.Errorf("vk.EnumerateInstanceLayerProperties: %w", err)
	}
	for _, p := range props {
		if p.LayerName == "VK_LAYER_KHRONOS_validation" {
			vr.apiLayers = append(vr.apiLayers, p.LayerName)
			return nil
		}
	}

	// validation layers are expected to be found in developer builds
	// Install the LunarG Vulkan SDK if they are missing
	return fmt.Errorf("khronos validation layer not found")
}
