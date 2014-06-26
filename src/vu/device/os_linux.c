// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

// FUTURE The ubuntu (linux) native layer implementation.
// This wraps the linux windows API's (where the real work is done).

#include "os_linux.h"

// Application defaults. Internal use only.  Not really state per-se these are 
// consulted at startup for initial values. These are updated using the 
// gs_set_attr* functions.
struct AppDefaults {
    long gs_ShellX;
    long gs_ShellY;
    long gs_ShellWidth;
    long gs_ShellHeight;
    long gs_AlphaSize;
    long gs_DepthSize;
    char gs_AppName[40];
};
struct AppDefaults defaults = { 100, 100, 240, 280, 8, 24, "App" };

// Initialize, register the application class and create the initial
// application window.
long gs_display_init()
{
	return 0;
} 

// Destroy the application window. Attempt to remove the rendering context and
// the device context as well.
void gs_display_dispose(long display) 
{
}

// Get the device context. This must be called after creating the window and
// before creating the rendering context.
long gs_shell(long display)
{
	return 0;
}

// Show the application window to the user. This is expected to be called after
// the rendering context has been created.
void gs_shell_open(long display)
{
}

// Check if the application window is still active.
unsigned char gs_shell_alive(long display) 
{ 
	return 0;
}

// Get the current mouse position relative to the bottom left corner of the
// application window.
void gs_pos(long display, long *x, long *y) 
{
}

// Position the cursor at the given window location.  The incoming coordinates
// are relative to the bottom left corner - switch that to be relative to the
// top left corner.
void gs_set_cursor_location(long display, long x, long y) 
{
}

// Process all queued up user events and send one of the processed events back
// to the application.
void gs_read_dispatch(long display, GSEvent *gs_urge)
{
}

// Get the current application windows client area location and size.
void gs_size(long display, long *x, long *y, long *w, long *h)
{
}

// Show or hide cursor.  Lock it to the window if it is hidden.
void gs_show_cursor(long display, unsigned char show) 
{
}

// gs_context creates an opengl context.  Actually it creates two of them.
// The first context is used to find better functions to create the final
// context.  Note that the pixel format is done only once for a window so
// it must be correctly chosen.
long gs_context(long long * display, long long * shell) 
{
	return 0;
}

// Flip the back and front buffers of the rendering context.
void gs_swap_buffers(long shell) 
{
}

// Set long attributes. Attributes only take effect if they are set before 
// they are used to create the window or rendering context. 
void gs_set_attr_l(long attr, long value) 
{
   switch (attr) {
   case GS_ShellX:
       if (value > 0) { defaults.gs_ShellX = value; }
       break;
   case GS_ShellY:
       if (value > 0) { defaults.gs_ShellY = value; }
       break;
   case GS_ShellWidth:
       if (value > 0) { defaults.gs_ShellWidth = value; }
       break;
   case GS_ShellHeight:
       if (value > 0) { defaults.gs_ShellHeight = value; }
       break;
   case GS_AlphaSize:
       if (value >= 0) { defaults.gs_AlphaSize = value; }
       break;
   case GS_DepthSize:
       if (value >= 0) { defaults.gs_DepthSize = value; }
       break;
   }
}

// Set string attributes.  Attributes only take effect if they are set before 
// they are used to create the window or rendering context. 
void gs_set_attr_s(long attr, char * value)
{
   switch (attr) {
   case GS_AppName:
       if (value != NULL && strlen(value) < 40) {
           strcpy( defaults.gs_AppName, value ); 
       }
       break;
   }
}


