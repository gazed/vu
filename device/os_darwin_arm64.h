// Copyright © 2017 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// Wrappers for access to basic iOS OpenGLES graphic context.

// Initialize the underlying Cocoa layer, create the default application
// window and OpenGL context. This call will not return.
void dev_run();

// Get the display size in pixels - always full screen so bottom
// left is 0,0.
void dev_size(int *w, int *h, int *scale);

// Shutting down.
// FUTURE:  remove for iOS since apps don't quit in iOS, iOS quits them.
void dev_dispose();

// Wrap logging so that normal logs show up in the iOS device console.
void dev_log(const char* log);

// device events for user input.
enum {
	devTouchBegin = 0, // touch.TypeBegin
	devTouchMove  = 1, // touch.TypeMove
	devTouchEnd   = 2, // touch.TypeEnd
    devResize     = 5,
    devFocus      = 6,
};
