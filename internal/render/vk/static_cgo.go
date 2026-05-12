// SPDX-FileCopyrightText : © 2026 Galvanized Logic Inc.
// SPDX-License-Identifier: MIT

//go:build !windows

package vk

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/unix"
)

/*
#include <stdlib.h>
#include <stdio.h>
#include <stdint.h>

#ifndef __DLLOAD_H__
#define __DLLOAD_H__

#include <stdint.h>

void* OpenLibrary(const char *name);
char* OpenLibraryError();
void CloseLibrary(void *lib_handle);

size_t Trampoline3(void *symbol, uintptr_t p0, uintptr_t p1, uintptr_t p2);
size_t Trampoline6(void *symbol, uintptr_t p0, uintptr_t p1, uintptr_t p2, uintptr_t p3, uintptr_t p4, uintptr_t p5);
size_t Trampoline9(void *symbol, uintptr_t p0, uintptr_t p1, uintptr_t p2, uintptr_t p3, uintptr_t p4, uintptr_t p5, uintptr_t p6, uintptr_t p7, uintptr_t p8);
size_t Trampoline12(void *symbol, uintptr_t p0, uintptr_t p1, uintptr_t p2, uintptr_t p3, uintptr_t p4, uintptr_t p5, uintptr_t p6, uintptr_t p7, uintptr_t p8, uintptr_t p9, uintptr_t p10, uintptr_t p11);
size_t Trampoline15(void *symbol, uintptr_t p0, uintptr_t p1, uintptr_t p2, uintptr_t p3, uintptr_t p4, uintptr_t p5, uintptr_t p6, uintptr_t p7, uintptr_t p8, uintptr_t p9, uintptr_t p10, uintptr_t p11, uintptr_t p12, uintptr_t p13, uintptr_t p14);

void* SymbolFromName(void *lib_handle, const void *name);

#endif

#include <dlfcn.h>

void* OpenLibrary(const char *name) {
	void* lib_handle = dlopen(name, RTLD_LOCAL|RTLD_LAZY);
	if (!lib_handle) {
		return NULL;
	}
	return lib_handle;
}

// Expected to be called if OpenLibrary returns nil.
char* OpenLibraryError() {
    return dlerror();
}

void CloseLibrary(void *lib_handle) {
	if (dlclose(lib_handle) != 0) {
		printf("Problem closing library: %s", dlerror());
	}
}

void* SymbolFromName(void *lib_handle, const void *name) {
	return dlsym(lib_handle, (const char*) name);
}

typedef size_t (*vkGeneric_func3)(uintptr_t, uintptr_t, uintptr_t);
typedef size_t (*vkGeneric_func6)(uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t);
typedef size_t (*vkGeneric_func9)(uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t);
typedef size_t (*vkGeneric_func12)(uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t);
typedef size_t (*vkGeneric_func15)(uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t);

size_t Trampoline3(void *symbol, size_t p0, size_t p1, size_t p2) {
	return ((vkGeneric_func3) symbol)(p0, p1, p2);
}
size_t Trampoline6(void *symbol, size_t p0, size_t p1, size_t p2, size_t p3, size_t p4, size_t p5) {
	return ((vkGeneric_func6) symbol)(p0, p1, p2, p3, p4, p5);
}
size_t Trampoline9(void *symbol, size_t p0, size_t p1, size_t p2, size_t p3, size_t p4, size_t p5, size_t p6, size_t p7, size_t p8) {
	return ((vkGeneric_func9) symbol)(p0, p1, p2, p3, p4, p5, p6, p7, p8);
}
size_t Trampoline12(void *symbol, uintptr_t p0, uintptr_t p1, uintptr_t p2, uintptr_t p3, uintptr_t p4, uintptr_t p5, uintptr_t p6, uintptr_t p7, uintptr_t p8, uintptr_t p9, uintptr_t p10, uintptr_t p11) {
		return ((vkGeneric_func12) symbol)(p0, p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, p11);
}
size_t Trampoline15(void *symbol, uintptr_t p0, uintptr_t p1, uintptr_t p2, uintptr_t p3, uintptr_t p4, uintptr_t p5, uintptr_t p6, uintptr_t p7, uintptr_t p8, uintptr_t p9, uintptr_t p10, uintptr_t p11, uintptr_t p12, uintptr_t p13, uintptr_t p14) {
		return ((vkGeneric_func15) symbol)(p0, p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, p11, p12, p13, p14);
}
*/
import "C"

var dlHandle unsafe.Pointer

// c-api function pointer that is lazy initialized.
type vkCommand unsafe.Pointer

func libName() string {
	switch runtime.GOOS {
	case "darwin":
		return "libMoltenVK.dylib"
	case "ios":
		return "MoltenVK.framework/MoltenVK"
	}
	return "libvulkan.so"
}

func loadLibrary(overrideLibName string) error {
	libName := libName()
	if overrideLibName != "" {
		libName = overrideLibName
	}
	cstr := C.CString(libName)
	dlHandle = C.OpenLibrary(cstr)
	C.free(unsafe.Pointer(cstr))

	// report any errors.
	if dlHandle == nil {
		cStr := C.OpenLibraryError()
		str := C.GoString(cStr)
		return fmt.Errorf("LoadVulkan failed: %s", str)
	}
	return nil
}

func sysStringToBytes(s string) *byte {
	p, _ := unix.BytePtrFromString(s)
	return p
}
