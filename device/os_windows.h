// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

#ifndef os_windows_h
#define os_windows_h

// os_windows.h defines the method calls needed by os_windows.go native layer.

#include <stdio.h>
#include <windows.h>

// Used to pass back user input each on each polling call.
typedef struct {
    long event;   // the user event. Zero if nothing is happening.
    long mousex;  // current mouse position is always filled in.
    long mousey;  // current mouse position is always filled in.
    long key;     // which key is currently pressed, if any.
    long mods;    // which modifier keys are currently pressed, if any.
    long scroll;  // the scroll amount if any.
} GSEvent;

// Used to toggle between full screen and windowed mode.
typedef struct {
    unsigned char full;     // true when in full screen mode.
    unsigned char maxed;    // true if window is maximized.
    long          style;    // used to restore windowed mode style.
    long          ex_style; // used to restore windowed mode style.
    RECT          rect;     // used to restore windowed dimensions.
} GSScreen;

// Initialize the underlying application window.
// Returns a reference to the application window (display).
long gs_display_init();

// Creates the device context (shell) on the given display.
// Returns a reference to the shell.
long gs_shell(long display);

// Open the application window on the given display.
void gs_shell_open(long display);

// Used to check for the user quitting the application.
// Return 1 as long as the user hasn't closed the window.
unsigned char gs_shell_alive(long shell);

// Used to check if the application is is full screen mode.
// Return 1 if the application is full screen, 0 otherwise.
unsigned char gs_fullscreen(long display);

// Flip full screen mode. Must be called after starting processing
// of events with gs_read_dispatch().
void gs_toggle_fullscreen(long display);

// Process a user event. This must be called inside an event loop in order
// for the application to work. The event is also processed to determine
// window events.
void gs_read_dispatch(long display, GSEvent *gs_urge);

// Get the current main window drawing area size.
void gs_size(long display, long *x, long *y, long *w, long *h);

// Show or hide cursor. Lock it if it is hidden.
void gs_show_cursor(long display, unsigned char show);

// Set the cursor location to the given screen coordinates.
void gs_set_cursor_location(long display, long x, long y);

// Create an OpenGL context using the given shell. Subsequent calls
// return the current context and ignoring the input parameters.
//
// Return 0 if a rendering context could not be created.
// This can happen if there are no renderers capable of handling
// the requested OpenGL attributes.
#ifdef __cplusplus
    extern "C" {
#endif
long gs_context(long long * display, long long * shell);
#ifdef __cplusplus
    }
#endif

// Flip the front and back rendering buffers. This is expected to be called
// each pass through the event loop to display the most recent drawing.
#ifdef __cplusplus
    extern "C" {
#endif
void gs_swap_buffers(long context);
#ifdef __cplusplus
    }
#endif

// Cleans and releases the application window.
#ifdef __cplusplus
    extern "C" {
#endif
void gs_display_dispose(long display);
#ifdef __cplusplus
    }
#endif

// Customize the window and context by setting attributes before the
// display or context is initialized.
void gs_set_attr_l(long attr, long value);
void gs_set_attr_s(long attr, char * value);

// Create the window, but don't open it.
long gs_create_window(HMODULE hInstance, LPSTR className);

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
    GS_SysKeyUp              = 0x0105, // WM_SYSKEYUP
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

// Provide key modifier bit masks. All currently pressed modifier
// keys come back combined into one bitmask value.
enum {
   GS_ShiftKeyMask      = 1 << 17,
   GS_ControlKeyMask    = 1 << 18,
   GS_CommandKeyMask    = 1 << 19,
   GS_FunctionKeyMask   = 1 << 20,
   GS_AlternateKeyMask  = 1 << 21,
};

#endif
