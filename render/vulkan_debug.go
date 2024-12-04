// Copyright © 2024 Galvanized Logic Inc.

//go:build debug

package render

// vulkan_debug.go includes vulkan debug utilitlies when
// building with "-tags debug"

import (
	"fmt"
	"log/slog"

	"github.com/gazed/vu/internal/render/vk"
)

// init is called before main to override the addValidationLayer method.
func init() {
	addValidationLayer = func(layers []string) ([]string, error) {
		slog.Debug("vulkan validation added")
		props, err := vk.EnumerateInstanceLayerProperties()
		if err != nil {
			return layers, fmt.Errorf("vk.EnumerateInstanceLayerProperties: %w", err)
		}
		for _, p := range props {
			if p.LayerName == "VK_LAYER_KHRONOS_validation" {
				return append(layers, p.LayerName), nil
			}
		}

		// complain, but continue. This may be a debug load running on
		// a testers machine that does not have the vulkan SDK.
		// Install the LunarG Vulkan SDK where appropriate.
		slog.Error("khronos validation layer not found")
		return layers, nil
	}
}
