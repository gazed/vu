// SPDX-FileCopyrightText : © 2026 Galvanized Logic Inc.
// SPDX-License-Identifier: MIT

//go:build windows

package vk

import "golang.org/x/sys/windows"

// pointer to the vulkan DLL.
var dlHandle *windows.LazyDLL

// reference to a vulkan command in the DLL.
type vkCommand = uintptr

// Needed to find vulkan command references with vkGetInstanceProcAddr.
// Set externally to the vulkan instance from vkCreateInstance.
var VKInst uintptr

func loadLibrary(overrideLibName string) error {
	libName := "vulkan-1.dll"
	if overrideLibName != "" {
		libName = overrideLibName
	}
	dlHandle = windows.NewLazyDLL(libName)
	return nil
}

func sysStringToBytes(s string) *byte {
	p, _ := windows.BytePtrFromString(s)
	return p
}
