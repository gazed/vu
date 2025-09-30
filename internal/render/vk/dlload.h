#ifndef __DLLOAD_H__
#define __DLLOAD_H__

#include <stdint.h>

void* OpenLibrary(const char *name);
void CloseLibrary(void *lib_handle);

size_t Trampoline3(void *symbol, uintptr_t p0, uintptr_t p1, uintptr_t p2);
size_t Trampoline6(void *symbol, uintptr_t p0, uintptr_t p1, uintptr_t p2, uintptr_t p3, uintptr_t p4, uintptr_t p5);
size_t Trampoline9(void *symbol, uintptr_t p0, uintptr_t p1, uintptr_t p2, uintptr_t p3, uintptr_t p4, uintptr_t p5, uintptr_t p6, uintptr_t p7, uintptr_t p8);
size_t Trampoline12(void *symbol, uintptr_t p0, uintptr_t p1, uintptr_t p2, uintptr_t p3, uintptr_t p4, uintptr_t p5, uintptr_t p6, uintptr_t p7, uintptr_t p8, uintptr_t p9, uintptr_t p10, uintptr_t p11);

void* SymbolFromName(void *lib_handle, const void *name);

#endif
