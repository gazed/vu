// Copyright © 2025 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// os_darwin.h exposes definitions needed for the go wrapper.

// Initialize the underlying Cocoa layer, create the default application
// window and graphics context. Return a pointer to a CaMetalLayer.
long int dev_init(char * title, long x, long y, long w, long h);
void dev_run(); // run - does not return.

// Customize the window and context by setting attributes after
// the display is initialized.
void dev_set_size(long x, long y, long w, long h);
void dev_set_title(char * label);

// Flip the front and back rendering buffers. This is expected to be called
// each pass through the event loop to display the most recent drawing.
void dev_swap();

// Cleans and releases all resources.
void dev_dispose();

// Used to check if the application is full screen mode.
// Return 1 if the application is full screen, 0 otherwise.
unsigned char dev_fullscreen();

// Flip full screen mode. Must be called after starting processing
// of events with gs_read_dispatch().
void dev_toggle_fullscreen();

// Get the current main window location and dimensions in screen pixels.
void dev_size(long *x, long *y, long *w, long *h);

// Get current cursor location.
void dev_cursor(long *x, long *y);

// device callback parameter values for user input events.
enum {
    devUp       = 1,
    devDown     = 2,
    devScroll   = 3,
    devMod      = 4,    // modifier flags.
    devMoved    = 5,
    devResized  = 6,
    devFocusIn  = 7,
    devFocusOut = 8,
    devMouseL   = 0xA0, // Don't conflict with key code.
    devMouseM   = 0xA1, //   ""
    devMouseR   = 0xA2, //   ""
};
