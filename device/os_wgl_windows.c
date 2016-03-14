// Copyright © 2015-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// The microsoft (windows) native layer opengl implementation.
// This wraps the microsoft API's to create a opengl context.

// +build windows,!dx
// Use opengl by default on windows.

#include "os_windows.h"

// Bind an opengl extension to control the swap interval (and its dependency).
//     http://www.opengl.org/registry/specs/EXT/wgl_swap_control.txt
//     http://www.opengl.org/registry/specs/ARB/wgl_extensions_string.txt
typedef int (APIENTRY * PFNWGLSWAPINTERVALEXTPROC) (int);
typedef const char *(APIENTRY * PFNWGLGETEXTENSIONSSTRINGARBPROC)( HDC );

// Bind an opengl extension in order to create a proper context on windows.
// (and its dependency) See: wglCreateContextAttribsARB
//     http://www.opengl.org/wiki/Creating_an_OpenGL_Context_(WGL)
//     http://www.opengl.org/registry/specs/ARB/wgl_create_context.txt
//     http://www.opengl.org/registry/specs/EXT/wgl_extensions_string.txt
typedef HGLRC (WINAPI * PFNWGLCREATECONTEXTATTRIBSARBPROC) (HDC, HGLRC, const int *);
typedef const char *(APIENTRY * PFNWGLGETEXTENSIONSSTRINGEXTPROC)( void );
#define WGL_CONTEXT_MAJOR_VERSION_ARB             0x2091
#define WGL_CONTEXT_MINOR_VERSION_ARB             0x2092
#define WGL_CONTEXT_LAYER_PLANE_ARB               0x2093
#define WGL_CONTEXT_FLAGS_ARB                     0x2094
#define WGL_CONTEXT_PROFILE_MASK_ARB              0x9126
#define WGL_CONTEXT_DEBUG_BIT_ARB                 0x0001
#define WGL_CONTEXT_FORWARD_COMPATIBLE_BIT_ARB    0x0002
#define WGL_CONTEXT_CORE_PROFILE_BIT_ARB          0x00000001
#define WGL_CONTEXT_COMPATIBILITY_PROFILE_BIT_ARB 0x00000002

// Bind an opengl extension to read and choose a pixel format using attributes.
//     http://www.opengl.org/registry/specs/ARB/wgl_pixel_format.txt
typedef BOOL (WINAPI * PFNWGLGETPIXELFORMATATTRIBIVARBPROC) (HDC, int, int, UINT, const int *, int *);
typedef BOOL (WINAPI * PFNWGLCHOOSEPIXELFORMATARBPROC) (HDC, const int *, const FLOAT *, UINT, int *, UINT *);
#define WGL_NUMBER_PIXEL_FORMATS_ARB    0x2000
#define WGL_DRAW_TO_WINDOW_ARB          0x2001
#define WGL_SUPPORT_OPENGL_ARB          0x2010
#define WGL_ACCELERATION_ARB            0x2003
#define WGL_DOUBLE_BUFFER_ARB           0x2011
#define WGL_STEREO_ARB                  0x2012
#define WGL_PIXEL_TYPE_ARB              0x2013
#define WGL_COLOR_BITS_ARB              0x2014
#define WGL_RED_BITS_ARB                0x2015
#define WGL_GREEN_BITS_ARB              0x2017
#define WGL_BLUE_BITS_ARB               0x2019
#define WGL_ALPHA_BITS_ARB              0x201B
#define WGL_ACCUM_BITS_ARB              0x201D
#define WGL_ACCUM_RED_BITS_ARB          0x201E
#define WGL_ACCUM_GREEN_BITS_ARB        0x201F
#define WGL_ACCUM_BLUE_BITS_ARB         0x2020
#define WGL_ACCUM_ALPHA_BITS_ARB        0x2021
#define WGL_DEPTH_BITS_ARB              0x2022
#define WGL_STENCIL_BITS_ARB            0x2023
#define WGL_AUX_BUFFERS_ARB             0x2024
#define WGL_SAMPLE_BUFFERS_ARB          0x2041
#define WGL_SAMPLES_ARB                 0x2042
#define WGL_NO_ACCELERATION_ARB         0x2025
#define WGL_GENERIC_ACCELERATION_ARB    0x2026
#define WGL_FULL_ACCELERATION_ARB       0x2027
#define WGL_TYPE_RGBA_ARB               0x202B
#define WGL_TYPE_COLORINDEX_ARB         0x202C

// The initial pixel format is used to get an initial rendering context so
// that more rendering functions can be loaded. These new rendering functions
// allow a better, final, rendering context to be created.
int gs_get_initial_pixelformat(long shell)
{
    HDC hdc = LongToHandle(shell);
    int flags = PFD_DRAW_TO_WINDOW | PFD_SUPPORT_OPENGL | PFD_DOUBLEBUFFER;
    PIXELFORMATDESCRIPTOR pfd =
    {
        sizeof(PIXELFORMATDESCRIPTOR),
        1,                 // version
        flags,             // see above
        PFD_TYPE_RGBA,     // type of framebuffer
        32,                // color depth
        0,0,0,0,0,0,0,0,   // red, green, blue, alpha bits
        0,0,0,0,0,         // accum bits
        24,                // depth buffer bits.
        0,                 // stencil buffer bits.
        PFD_MAIN_PLANE,    // layer
        0,0,0,0,0          // unused
    };

    // create the temporary context using the proper pixel format.
    int pixelFormat = ChoosePixelFormat(hdc, &pfd);
    if (pixelFormat != 0)
    {
        SetPixelFormat(hdc, pixelFormat, &pfd);
        return pixelFormat;
    }
    return 0;
}

// OpenGL extensions that are bound after the first context is created.
PFNWGLGETEXTENSIONSSTRINGEXTPROC gs_wglGetExtensionsStringEXT = NULL;
PFNWGLSWAPINTERVALEXTPROC gs_wglSwapIntervalARB =  NULL;
PFNWGLGETEXTENSIONSSTRINGARBPROC gs_wglGetExtensionsStringARB =  NULL;
PFNWGLCREATECONTEXTATTRIBSARBPROC gs_wglCreateContextAttribsARB =  NULL;
PFNWGLGETPIXELFORMATATTRIBIVARBPROC gs_wglGetPixelFormatAttribivARB =  NULL;
PFNWGLCHOOSEPIXELFORMATARBPROC gs_wglChoosePixelFormatARB =  NULL;

// The final pixel format created using the bound rendering functions.
int gs_get_pixelformat(long shell)
{
    const int attribList[] =
    {
        WGL_DRAW_TO_WINDOW_ARB, 1,
        WGL_SUPPORT_OPENGL_ARB, 1,
        WGL_DOUBLE_BUFFER_ARB, 1,
        WGL_PIXEL_TYPE_ARB, WGL_TYPE_RGBA_ARB,
        WGL_ACCELERATION_ARB, WGL_FULL_ACCELERATION_ARB,
        WGL_COLOR_BITS_ARB, 32,
        WGL_DEPTH_BITS_ARB, 24,
        WGL_STENCIL_BITS_ARB, 8,
        0, // end
    };
    int pixelFormat;
    UINT numFormats;
    HDC hdc = LongToHandle(shell);
    gs_wglChoosePixelFormatARB(hdc, attribList, NULL, 1, &pixelFormat, &numFormats);
    if (pixelFormat != 0)
    {
        PIXELFORMATDESCRIPTOR pfd;
        DescribePixelFormat(hdc, pixelFormat, sizeof(PIXELFORMATDESCRIPTOR), &pfd);
        SetPixelFormat(hdc, pixelFormat, &pfd);
        return pixelFormat;
    }
    return 0;
}

// gs_context creates an opengl context. Actually it creates two of them.
// The first context is used to find better functions to create the final
// context. Note that the pixel format is done only once for a window so
// it must be correctly chosen.
long gs_context(long long * display, long long * shell)
{
    // create the initial context.
    HDC hdc = LongToHandle(*shell);
    HGLRC initialContext;
    int initial_pixelFormat = gs_get_initial_pixelformat(*shell);
    if (initial_pixelFormat != 0)
    {
        initialContext = wglCreateContext(hdc);
        if (initialContext != NULL)
        {
            if (!wglMakeCurrent(hdc, initialContext))
            {
                wglDeleteContext(initialContext);
                initialContext = NULL;
            }
        }
    }
    if (initialContext == NULL)
    {
        return 0; // failed to get even a simple context.
    }

    // now that there is a context, bind the opengl extensions and fail
    // if the supported extensions are too old or if they are not there.
    gs_wglGetExtensionsStringEXT = (PFNWGLGETEXTENSIONSSTRINGEXTPROC) wglGetProcAddress( "wglGetExtensionsStringEXT" );
    gs_wglSwapIntervalARB = (PFNWGLSWAPINTERVALEXTPROC) wglGetProcAddress( "wglSwapIntervalEXT" );
    gs_wglGetExtensionsStringARB = (PFNWGLGETEXTENSIONSSTRINGARBPROC) wglGetProcAddress( "wglGetExtensionsStringARB" );
    gs_wglCreateContextAttribsARB = (PFNWGLCREATECONTEXTATTRIBSARBPROC) wglGetProcAddress( "wglCreateContextAttribsARB" );
    gs_wglGetPixelFormatAttribivARB = (PFNWGLGETPIXELFORMATATTRIBIVARBPROC) wglGetProcAddress( "wglGetPixelFormatAttribivARB" );
    gs_wglChoosePixelFormatARB = (PFNWGLCHOOSEPIXELFORMATARBPROC) wglGetProcAddress( "wglChoosePixelFormatARB" );
    if (gs_wglGetExtensionsStringEXT == NULL ||
        gs_wglSwapIntervalARB ==  NULL ||
        gs_wglGetExtensionsStringARB ==  NULL ||
        gs_wglCreateContextAttribsARB ==  NULL ||
        gs_wglGetPixelFormatAttribivARB ==  NULL ||
        gs_wglChoosePixelFormatARB ==  NULL)
    {
        return 0;
    }

    // destroy and recreate the window and shell
    LPSTR gs_className = TEXT("GS_WIN");
    gs_display_dispose(*display);
    HMODULE hInstance = GetModuleHandle(NULL);
    *display = gs_create_window(hInstance, gs_className);
    *shell = gs_shell(*display);
    int pixelformat = gs_get_pixelformat(*shell);
    if (pixelformat == 0)
    {
        return 0;
    }

    // now create the context on the fresh window.
    int cnt = 0;
    int attribs[40];
    hdc = LongToHandle(*shell);

    // Use the expected baseline opengl 3.2
    attribs[cnt++] = WGL_CONTEXT_MAJOR_VERSION_ARB;
    attribs[cnt++] = 3;
    attribs[cnt++] = WGL_CONTEXT_MINOR_VERSION_ARB;
    attribs[cnt++] = 2;
    attribs[cnt++] = WGL_CONTEXT_FLAGS_ARB;
    attribs[cnt++] = WGL_CONTEXT_FORWARD_COMPATIBLE_BIT_ARB;
    attribs[cnt++] = WGL_CONTEXT_PROFILE_MASK_ARB;
    attribs[cnt++] = WGL_CONTEXT_CORE_PROFILE_BIT_ARB;
    attribs[cnt++] = 0;
    HGLRC context = gs_wglCreateContextAttribsARB( hdc, NULL, attribs );
    if (context != NULL)
    {
        if (wglMakeCurrent(hdc, context))
        {
            return HandleToLong(context);
        }
    }
    return 0; // failed to get rendering context
}

// Flip the back and front buffers of the rendering context.
void gs_swap_buffers(long shell)
{
    HDC hdc = LongToHandle(shell);
    SwapBuffers(hdc);
}

// Destroy the application window. Attempt to remove the rendering context and
// the device context as well.
void gs_display_dispose(long display)
{
    HWND hwnd = LongToHandle(display);
    HDC shell = GetDC(hwnd);
    HGLRC context = wglGetCurrentContext();
    wglMakeCurrent(NULL, NULL);
    wglDeleteContext(context);
    ReleaseDC(hwnd, shell);
    DestroyWindow(hwnd);
}
