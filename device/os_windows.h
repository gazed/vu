// Copyright © 2013-2017 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// os_windows.h exposes the native layer needed by os_windows.go.

#ifndef os_windows_h
#define os_windows_h

#include <stdio.h>
#include <windows.h>

// Initialize the underlying Windows layer, create the default application
// window and OpenGL context. This call will not return.
void dev_run();

// Customize the window and context by setting attributes after
// the display is initialized.
void dev_set_size(long x, long y, long w, long h);
void dev_set_title(char * label);

// Flip the front and back rendering buffers. This is expected to be called
// each pass through the event loop to display the most recent drawing.
void dev_swap();

// Cleans and releases all resources including the OpenGL context.
void dev_dispose();

// Copy and paste strings to and from the general clipboard.
// Strings returned by copy must be freed by the caller.
char* dev_clip_copy();
void dev_clip_paste(const char* string);

// Used to check if the application is full screen mode.
// Return 1 if the application is full screen, 0 otherwise.
unsigned char dev_fullscreen();

// Flip full screen mode. Must be called after starting processing
// of events with gs_read_dispatch().
void dev_toggle_fullscreen();

// Get the current main window drawing area size.
void dev_size(long *x, long *y, long *w, long *h);

// Show or hide cursor. Lock it if it is hidden.
void dev_show_cursor(unsigned char show);

// Get current cursor location.
void dev_cursor(long *x, long *y);

// Set the cursor location to the given screen coordinates.
void dev_set_cursor_location(long x, long y);

// Used to toggle between full screen and windowed mode.
typedef struct {
    unsigned char full;     // true when in full screen mode.
    unsigned char maxed;    // true if window is maximized.
    long          style;    // used to restore windowed mode style.
    long          ex_style; // used to restore windowed mode style.
    RECT          rect;     // used to restore windowed dimensions.
} screenInfo;


// device callback parameter values for user input events.
enum {
    devUp       = 1,
    devDown     = 2,
    devScroll   = 3,
    devResize   = 5,
    devFocusIn  = 6,
    devFocusOut = 7,

    // codes that do not conflict with other virtual key codes.
    devMouseL   = VK_LBUTTON, // 0x01 Left mouse button
    devMouseM   = VK_MBUTTON, // 0x04 Middle mouse button
    devMouseR   = VK_RBUTTON, // 0x02 Right mouse button
};

#endif
