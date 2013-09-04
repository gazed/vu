// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

#ifndef os_windows_h
#define os_windows_h

#include <stdio.h>
#include <windows.h>

// Bind an opengl extension to control the swap interval (and its dependency).
//     http://www.opengl.org/registry/specs/EXT/wgl_swap_control.txt
//     http://www.opengl.org/registry/specs/ARB/wgl_extensions_string.txt
typedef int (APIENTRY * PFNWGLSWAPINTERVALEXTPROC) (int);
typedef const char *(APIENTRY * PFNWGLGETEXTENSIONSSTRINGARBPROC)( HDC );

// Bind an opengl extension in order to create a proper context on windows.
// (and its dependency)
//     http://www.opengl.org/wiki/Creating_an_OpenGL_Context_(WGL)
//     http://www.opengl.org/registry/specs/ARB/wgl_create_context.txt
//     http://www.opengl.org/registry/specs/EXT/wgl_extensions_string.txt
// wglCreateContextAttribsARB
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

// Used to pass back user input each on each polling call.
typedef struct {
	long event;   // the user event. Zero if nothing is happening. 
	long mousex;  // current mouse position is always filled in.
	long mousey;  // current mouse position is always filled in.
	long key;     // which key is currently pressed, if any.
	long mods;	  // which modifier keys are currently pressed, if any.
	long scroll;  // the scroll amount if any. 
} GSEvent;

// Initialize the underlying application window.
// Returns a reference to the application window (display).
long gs_display_init();

// Cleans and releases the application window. 
void gs_display_dispose(long display);

// Creates the device context (shell) on the given display.
// Returns a reference to the shell. 
long gs_shell(long display);

// Open the applicaiton window on the given display.
void gs_shell_open(long display);

// Used to check for the user quitting the application.
// Return 1 as long as the user hasn't closed the window.
unsigned char gs_shell_alive(long shell);

// Process a user event.  This must be called inside an event loop in order
// for the application to work.  The event is also processed to determine
// window events.
void gs_read_dispatch(long display, GSEvent *gs_urge);

// Get the current main window drawing area size. 
void gs_size(long display, long *x, long *y, long *w, long *h);

// Show or hide cursor.  Lock it if it is hidden.
void gs_show_cursor(long display, unsigned char show);

// Set the cursor location to the given screen coordinates.
void gs_set_cursor_location(long display, long x, long y);

// Create an OpenGL context using the given shell. Subsequent calls will
// return the current context (ignoring the input parameter).
//
// This may return 0 if a rendering context could not be created.
// This can happen if there are no renderers capable of handling
// the requested OpenGL attributes.
long gs_context(long long * display, long long * shell);

// Flip the front and back rendering buffers. This is expected to be called
// each pass through the event loop to display the most recent drawing.
void gs_swap_buffers(long context);

// Customize the window and context by setting attributes before the
// display or context is initialized.
void gs_set_attr_l(long attr, long value);
void gs_set_attr_s(long attr, char * value);

// Used in the provided setter functions to set one or more of the
// following attributes.
enum AppAttributes 
{
	GS_AppName,     // Text("App")
	GS_ShellX,      // 100
	GS_ShellY,      // 100
	GS_ShellWidth,  // 640
	GS_ShellHeight, // 480
	GS_AlphaSize,   //  8
	GS_DepthSize    // 24
};

// Possible return values from gs_read_dispatch.
enum {
    GS_LeftMouseDown         = 0x0201, // WM_LBUTTONDOWN
    GS_LeftMouseUp           = 0x0202, // WM_LBUTTONUP
    GS_RightMouseDown        = 0x0204, // WM_RBUTTONDOWN
    GS_RightMouseUp          = 0x0205, // WM_RBUTTONUP
    GS_MouseMoved            = 0x0200, // WM_MOUSEMOVE 
    GS_MouseExited           = 0x02a3, // WM_MOUSELEAVE
    GS_KeyDown               = 0x0100, // WM_KEYDOWN
    GS_KeyUp                 = 0x0101, // WM_KEYUP
    GS_ScrollWheel           = 0x020a, // WM_MOUSEWHEEL
    GS_OtherMouseDown        = 0x0207, // WM_MBUTTONDOWN
    GS_OtherMouseUp          = 0x0208, // WM_MBUTTONUP
    GS_WindowResized         = 0x0232, // WM_EXITSIZEMOVE 
    GS_WindowMoved           = 0x0003, // WM_MOVE
    GS_WindowIconified       = 0x0019, // WM_SHOWWINDOW + true  (1)
    GS_WindowUniconified     = 0x0018, // WM_SHOWWINDOW + false (0)
    GS_WindowActive          = 0x0007, // WM_ACTIVATE + WA_ACTIVE (1)
    GS_WindowInactive        = 0x0006  // WM_ACTIVATE + WM_INACTIVE (0)
};

// Provide key modifier bit masks.  All currently pressed modifier keys come 
// back combined into one bitmask value.
enum {
   GS_ShiftKeyMask      = 1 << 1, 
   GS_ControlKeyMask    = 1 << 2,
   GS_CommandKeyMask    = 1 << 3,
   GS_FunctionKeyMask   = 1 << 4,
   GS_AlternateKeyMask  = 1 << 5,
};

#endif
