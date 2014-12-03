// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// OSX support to get a window with an OpenGL graphic context.
// All the window is expected to do is to be able to run in full screen mode,
// quit, show up in the dock, and participate in command-tab application switching.
//
// Maintenance Notes:
// The design is to wrap cocoa functionality in C syntax method so a simple
// binding layer can be created. Design goals include:
//    - minimize state, passing in needed information where possible.
//    - keep in/out parameter types to the basic C types.
//    - minimize the number of calls.
//    - use reasonable defaults where possible.
//    - duplicate enums were necessary so that no extra includes are needed.

// Used to pass back user input each polling call.
typedef struct {
	long event;   // the user event. Zero if nothing is happening. 
	long mousex;  // current mouse position is always filled in.
	long mousey;  // current mouse position is always filled in.
	long key;     // which key, or mouse button was affected, if any.
	long mods;	  // which modifier keys are currently pressed, if any.
	long scroll;  // the scroll amount if any. 
} GSEvent;

// Initialize the underlying Cocoa layer and create the default application.
// Returns a reference to the shared NSApplication instance (display).
long gs_display_init();

// Cleans and releases all resources including the OpenGL context.
void gs_display_dispose(long display);

// Creates the window (shell) on the given display.
// Returns a reference to the shell. 
long gs_shell(long display);

// Creates the window (shell) on the given display.
// Returns a reference to a NSWindow instance.
void gs_shell_open(long display);

// Used to check for the user quitting the application.
// Return 1 as long as the user hasn't closed the window.
unsigned char gs_shell_alive(long shell);

// Process a user event. This must be called inside an event loop in order
// for the application to work. The event is also processed to determine
// window events.
void gs_read_dispatch(long display, GSEvent *gs_urge);

// Get the current main window drawing area size. 
void gs_size(long shell, float *x, float*y, float *w, float *h);

// Show or hide cursor. Lock it if it is hidden.
void gs_show_cursor(unsigned char show);

// Set the cursor location to the given screen coordinates.
void gs_set_cursor_location(long display, long x, long y);

// Create an OpenGL context using the given shell. Subsequent calls will
// return the current context (ignoring the input parameter).
//
// This may return 0 if a rendering context could not be created.
// This can happen if there are no renderers capable of handling
// the requested OpenGL attributes.
long gs_context(long shell);

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

// Possible return values from gs_read_dispatch. The NSEventType
// values are defined in NSEvent.h
enum {
    GS_LeftMouseDown         = 1,   // NSLeftMouseDown
    GS_LeftMouseUp           = 2,   // NSLeftMouseUp
    GS_RightMouseDown        = 3,   // NSRightMouseDown
    GS_RightMouseUp          = 4,   // NSRightMouseUp
    GS_MouseMoved            = 5,   // NSMouseMoved
    GS_MouseEntered          = 8,   // NSMouseEntered
    GS_MouseExited           = 9,   // NSMouseExited
    GS_KeyDown               = 10,  // NSKeyDown
    GS_KeyUp                 = 11,  // NSKeyUp
    GS_ModKeysChanged        = 12,  // NSFlagsChanged
    GS_ScrollWheel           = 22,  // NSScrollWheel
    GS_OtherMouseDown        = 25,  // NSOtherMouseDown
    GS_OtherMouseUp          = 26,  // NSOtherMouseUp

    // Extra event types that don't conflict with NSEventType.
    GS_WindowResized         = 50,
    GS_WindowMoved           = 51,
    GS_WindowIconified       = 52,
    GS_WindowUniconified     = 53,
    GS_WindowActive          = 54,
    GS_WindowInactive        = 55
};

// Wrap the underlying key modifier definitions.
// All currently pressed modifier keys come back combined into one bitmask value.
enum {
   GS_ShiftKeyMask      = 1 << 17,  // NSShiftKeyMask
   GS_ControlKeyMask    = 1 << 18,  // NSControlKeyMask
   GS_AlternateKeyMask  = 1 << 19,  // NSAlternateKeyMask
   GS_CommandKeyMask    = 1 << 20,  // NSCommandKeyMask
   GS_FunctionKeyMask   = 1 << 23,  // NSFunctionKeyMask
};
