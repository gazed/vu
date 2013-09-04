// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

// The OSX (darwin) native layer implementation.
// This wraps the OSX API's (where the real work is done).

#import <Cocoa/Cocoa.h>
#import "os_darwin.h"

// Application defaults. Internal use only.
// Not really state per-se these are consulted at startup for initial values.
// These are set before they are used with the gs_set_attr* functions.
typedef struct {
    long gs_ShellX;
    long gs_ShellY;
    long gs_ShellWidth;
    long gs_ShellHeight;
    long gs_AlphaSize;
    long gs_DepthSize;
    NSString *gs_AppName;
} GSDefaults;
GSDefaults defaults = { 100, 100, 640, 480, 8, 24, @"App" };

// Get current mouse position independent of the event.
void gs_pos(long display, float *x, float *y) {
	NSWindow *window = [(id)display mainWindow];
    NSPoint origin = [window mouseLocationOutsideOfEventStream];
    *x = origin.x;
    *y = origin.y;
}

// Position the cursor at the given window location.  
void gs_set_cursor_location(long display, long x, long y) {
    NSWindow *window = [(id)display mainWindow];
    NSRect windowRect = [window frame]; // origin at bottom left in screen coordinates.
    CGRect screenRect = CGDisplayBounds(CGMainDisplayID()); // origin at top left.
    CGPoint point = CGPointMake(windowRect.origin.x+x, screenRect.size.height-windowRect.origin.y-y);
	CGWarpMouseCursorPosition(point);
	CGAssociateMouseAndMouseCursorPosition(true);
}

// Get any key modifiers independent of the current event.
void gs_mod(long display, long *key_mod) {
    *key_mod = (long) [NSEvent modifierFlags];
}

// Get the key code associated with the current event.  This will be 0 if the
// last event was not a key event.
void gs_key(long display, long *key_code) {
    *key_code = 0;
    NSEvent *event = [(id)display currentEvent];
    NSEventType etype = [event type];
    if ( NSKeyDown == etype || NSKeyUp == etype) {
        *key_code = [event keyCode];
	}
}

// Get the current scroll wheel value.  This will be 0 if the last event was
// not a scroll event.
void gs_scroll(long display, float *x_delta, float *y_delta) {
    *x_delta = 0;
    *y_delta = 0;
    NSEvent *event = [(id)display currentEvent];
    if (NSScrollWheel == [event type]) {
        *x_delta = [event deltaX];
        *y_delta = [event deltaY];
    }
}

// Show or hide cursor.  Note that the cursor is not locked with
// CGAssociateMouseAndMouseCursorPosition(show) or the cursor position
// would not change.
void gs_show_cursor(unsigned char show) {
	if (show) {
        [NSCursor unhide];
	} else {
        [NSCursor hide];
	}
}

// Called before running the application to create a few menu items.
//
// This uses a hidden API (setAppleMenu) so that the application menu can
// be set programatically.  This avoids the use of a NIB (binary file
// created by Interface Builder). The benefits are not having to store
// binary files in source control, and not forcing an opaque layer
// on users and maintainers of this simple library.
// All indications point to this hidden API staying around until it can
// be replaced with a supported one.  There are some hints of Apple
// eventually moving to NIB'less applications.
static void createMenus(NSApplication *display)
{
    // the top level menu bar.
    NSString *appName = defaults.gs_AppName;
    NSMenu *menuBar = [[NSMenu alloc] initWithTitle: @"MainMenu"];
    [display setMainMenu: menuBar];

    // the main application menu
    NSMenuItem *mi = [menuBar addItemWithTitle:@"Apple" action:NULL keyEquivalent:@""];
    NSMenu *m = [[NSMenu alloc] initWithTitle:@"Apple"];
    [display performSelector:@selector(setAppleMenu:) withObject: m];
    [menuBar setSubmenu:m forItem:mi];
    [m addItemWithTitle:[NSString stringWithFormat:@"Quit %@", appName]
                 action:@selector(orderOut:)
          keyEquivalent:@"q"];

    // the view menu
    mi = [menuBar addItemWithTitle:@"View" action:NULL keyEquivalent:@""];
    m = [[NSMenu alloc] initWithTitle:@"View"];
    [menuBar setSubmenu:m forItem:mi];
    [[m addItemWithTitle:@"Enter Full Screen"
                  action:@selector(toggleFullScreen:)
           keyEquivalent:@"f"]
        setKeyEquivalentModifierMask:NSControlKeyMask | NSCommandKeyMask];
}

// State used to track window events.  Set when a window event occurs, it is
// expected to be queried each time through the event loop (gs_read_dispatch)
// and then reset back to zero.
long wEvent = 0;

// State used to track window closure. This is needed to avoid accessing the
// external shell pointer after a window has closes.  There is no sure way to
// check if an object pointer is valid once that object has been released.
long gs_win_alive = -1;

// Used to get window notifications since it is far easier to let the window code
// figure out what particular mouse clicks and drags mean.
// These will be triggered as the underlying window processes the mouse moves and
// clicks sent during the gs_read_dispatch calls.
@interface EventDelegate : NSView <NSWindowDelegate> { }
@end
@implementation EventDelegate
-(void)windowWillClose:(NSNotification *)notification { gs_win_alive = -2; }
-(void)windowDidResize:(NSNotification *)notification { wEvent = GS_WindowResized; }
-(void)windowDidMove:(NSNotification *)notification { wEvent = GS_WindowMoved; }
-(void)windowDidMiniaturize:(NSNotification *)notification { wEvent = GS_WindowIconified; }
-(void)windowDidDeminiaturize:(NSNotification *)notification { wEvent = GS_WindowUniconified; }
-(void)windowDidBecomeKey:(NSNotification *)notification { wEvent = GS_WindowActive; }
-(void)windowDidResignKey:(NSNotification *)notification { wEvent = GS_WindowInactive; }

// let OS know that this app handles keys in order to prevent beeping.
-(BOOL)canBecomeKeyView { return YES; }
-(BOOL)acceptsFirstResponder { return YES; }
-(void)keyUp:(NSNotification *)notification {  }
-(void)keyDown:(NSNotification *)notification { }
@end

// Create the top level application (display).
long gs_display_init() {
    NSAutoreleasePool *pool = [[NSAutoreleasePool alloc] init];
    NSApplication *display = [NSApplication sharedApplication];
    [pool drain];
    return (long) display;
}

// Cleanup and quit the application.
void gs_display_dispose(long display) {
    NSOpenGLContext *context = [NSOpenGLContext currentContext];
    if (context != nil) {
        [NSOpenGLContext clearCurrentContext];
        [(id) context release];
    }
    [(id)display terminate: nil];
}

// Create the window.
long createShell(long display) {
    NSRect frame = NSMakeRect( defaults.gs_ShellX, defaults.gs_ShellY, defaults.gs_ShellWidth, defaults.gs_ShellHeight );
    unsigned int styleMask = NSTitledWindowMask | NSClosableWindowMask | NSMiniaturizableWindowMask | NSResizableWindowMask;
    NSWindow *window = [[NSWindow alloc]
        initWithContentRect:frame
                  styleMask:styleMask
                    backing:NSBackingStoreBuffered
                      defer:false];
    [window setTitle:defaults.gs_AppName];

    // Hook in the delegate.
    EventDelegate *delegate = [[[EventDelegate alloc] initWithFrame:frame] autorelease];
    [window setContentView:delegate];
    [window setDelegate:delegate];
    [window makeKeyWindow];
    [window orderBack:nil];
    [window setCollectionBehavior:[window collectionBehavior] | NSWindowCollectionBehaviorFullScreenPrimary];
	gs_win_alive = 1;

    // add in the menus
    createMenus((id) display);
    return (long) window;
}

// Create the window.  Note that the "run" loop is driven externally by calling
// (and it must be called) the "gs_read_dispatch" function.
long gs_shell(long display) {
    NSAutoreleasePool *pool = [[NSAutoreleasePool alloc] init];
    long shell = createShell(display);
    [pool drain];
    return shell;
}

// Open the window.  Note that the "run" loop is driven externally by calling
// (and it must be called) the "gs_read_dispatch" function.
void gs_shell_open(long display) {
    NSAutoreleasePool *pool = [[NSAutoreleasePool alloc] init];
    // bring the window to the front.
    ProcessSerialNumber psn = { 0, kCurrentProcess };
    TransformProcessType( &psn, kProcessTransformToForegroundApplication );
    SetFrontProcess( &psn );

    // this is needed to finialize hooking in the menu bar.
    [(id) display finishLaunching];
    [pool drain];
}

// Called from gs_read_dispatch to fill out an event type to return.
// Only return one event per call so choose a window event over a basic event
// (if there's a choice).
void gs_urge_info(long display, NSEvent *event, GSEvent *urge) {
	if (nil != event) {
        urge->event = (long) [event type];
        if (wEvent != 0) {
            urge->event = wEvent;
            wEvent = 0;

			// update the opengl context each window resize and move.
			if (urge->event == GS_WindowResized || urge->event == GS_WindowMoved) {
    			NSOpenGLContext *context = [NSOpenGLContext currentContext];
    			[(id)context update];
			}
        }

        // add event applicable information. 
        if (urge->event == GS_KeyDown || urge->event == GS_KeyUp) {
            long mods, kcode;
            gs_mod(display, &(urge->mods));
            gs_key(display, &(urge->key));
        } else if (urge->event == GS_ScrollWheel) {
            float dx, dy;
            gs_scroll(display, &dx, &dy);
			urge->scroll = (long)dy;
		} 
	}

	// always update the mouse.
	float mx, my;
	gs_pos(display, &mx, &my);
	urge->mousex = (long)mx;
	urge->mousey = (long)my;
}

// Get the next event and return the type. This mimics [NSApplication run] and allows
// external control over processing user events.  The expectation is that this is called
// in a tight loop like a gaming control loop.
void gs_read_dispatch(long display, GSEvent *gs_urge) {
    NSAutoreleasePool *pool = [[NSAutoreleasePool alloc] init];
    NSDate *date = [NSDate distantPast];
    NSEvent *event =
        [(id)display
        nextEventMatchingMask:NSAnyEventMask
                    untilDate:date
                       inMode:NSDefaultRunLoopMode
                      dequeue:YES];

    // have the view and window process the basic events
	if (nil != event) {
        [(id)display sendEvent:event];
	}
    [pool release];
    gs_urge_info(display, event, gs_urge);
}

// The window is hidden in the windowWillClose event and the main loop is expected
// to call this method to check if the shell should be terminated.  This allows the
// application a chance to do any final cleanup before everything stops.
//
// While this does not kill an iconified window, it should be replaced if a
// more appropriate window state can be found for communicating an application
// shutdown.
unsigned char gs_shell_alive(long shell) {
	id win = (id)shell;
	return ((gs_win_alive == 1) && ([win isMiniaturized] == YES || [win isVisible] == YES));
}

// Create the OpenGL context.  This must be called after the shell has
// been created.
long gs_context(long shell) {
    NSOpenGLContext *context = [NSOpenGLContext currentContext];
    if (context == nil) {
        NSOpenGLPixelFormatAttribute pixelAttrs[] = {
            NSOpenGLPFAOpenGLProfile, NSOpenGLProfileVersion3_2Core,
            NSOpenGLPFADoubleBuffer,
            NSOpenGLPFADepthSize,     defaults.gs_DepthSize,
            0
        };
        NSOpenGLPixelFormat *pixelFormat = [[NSOpenGLPixelFormat alloc] initWithAttributes:pixelAttrs];
        if (pixelFormat != nil) {
            context = [[NSOpenGLContext alloc] initWithFormat:pixelFormat
                                                 shareContext:nil];
            if (context != nil) {
                [context setView:[(id)shell contentView]];
                [context makeCurrentContext];
            }
        }
    }
    return (long) context;
}

// Flip the front and back buffers.
void gs_swap_buffers(long context) {
    [(id)context flushBuffer];
}

// Get current shell size. 
void gs_size(long shell, float *x, float *y, float *w, float *h) { 
	NSRect frame = [(id)shell frame];
	NSRect content = [(id)shell contentRectForFrameRect:frame];
	*x = frame.origin.x;
	*y = frame.origin.y;
    *w = content.size.width;
    *h = content.size.height;
}

// Update startup numeric defaults. Minimal effort to ensure a good value.
void gs_set_attr_l(long attr, long value) {
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

// Update startup string defaults. Minimal effort to ensure a good value.
void gs_set_attr_s(long attr, char * value) {
    switch (attr) {
    case GS_AppName:
        if (value != nil) {
            defaults.gs_AppName = [NSString stringWithUTF8String:value];
        }
        break;
    }
}

