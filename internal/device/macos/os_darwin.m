// Copyright © 2025 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.
//

#import <Cocoa/Cocoa.h>
#import <Metal/Metal.h>
#import <MetalKit/MetalKit.h>
#import "os_darwin.h"

// The golang device callbacks that would normally be defined in _cgo_export.h
// are manually reproduced here and also implemented in the native test file
// os_darwin_amd64_test.m.
// extern void prepRender();
extern void renderFrame();
extern void handleInput(long event, long data);

// Save a reference to the app window that is created on startup.
// This is used for resizing.
NSWindow *app_window; // guaranteed to init to 0.

// default clear color.
static MTLClearColor CLEAR_COLOR = { 0.0, 0.0, 0.0, 1.0 };

//==============================================================================
// VuRenderer is used to get render callbacks that are
// roughly related to the display refresh rate.
@interface VuRenderer : NSObject <MTKViewDelegate>
@end

@implementation VuRenderer { }
- (void)mtkView:(MTKView *)view drawableSizeWillChange:(CGSize) size { }
- (void)drawInMTKView:(MTKView *)view { renderFrame(); }
@end

//==============================================================================
// VuView handles the window content.
// NOTES:
// o MTKView is as UIView on ios and an NSView on macos.
// o UIView is used on iOS (Cocoa Touch): coords 0,0 in top left with positive values of Y going down
// o NSView on Mac (Cocoa): 0,0 in lower left with positive values of Y going up
@interface VuView : MTKView <NSWindowDelegate> {}
@end

@implementation VuView
-(BOOL)canBecomeKeyView { return YES; }
-(BOOL)acceptsFirstResponder { return YES; } // Window accepts input events

// Calls to super is discouraged except in the case of rightMouseDown
// where the apple documents of NSView say that super should be called.
- (void) rightMouseDown: (NSEvent*) event {
    [super rightMouseDown:event];
    handleInput(devDown, devMouseR);
}
- (void) rightMouseUp: (NSEvent*) event   { handleInput(devUp, devMouseR); }
- (void) mouseDown: (NSEvent*) event      { handleInput(devDown, devMouseL); }
- (void) mouseUp: (NSEvent*) event        { handleInput(devUp, devMouseL); }
- (void) otherMouseDown: (NSEvent*) event { handleInput(devDown, devMouseM); }
- (void) otherMouseUp: (NSEvent*) event   { handleInput(devUp, devMouseM); }
- (void) scrollWheel: (NSEvent*) event    { handleInput(devScroll, [event deltaY]); }
- (void) flagsChanged: (NSEvent*) event   { handleInput(devMod, [event modifierFlags]); }

// Get one callback once resizing is finished.
- (void)windowDidEndLiveResize:(NSNotification*)notification { handleInput(devResized, 0); }
- (void)windowDidMove:(NSNotification*)notification { handleInput(devMoved, 0); }

// register all key presses and releases.
- (void) keyUp: (NSEvent*) event { handleInput(devUp, [event keyCode]); }
- (void) keyDown: (NSEvent*) event {
    if ([event isARepeat] == NO) {
        handleInput(devDown, [event keyCode]);
    }
}

// focusEvent handles stopping the CVDisplay when the window
// loses focus.
-(void) focus: (int)focusEvent {
    if (focusEvent == devFocusOut) {
        handleInput(focusEvent, 0);
    }
    if (focusEvent == devFocusIn) {
        handleInput(focusEvent, 0);
    }
}

// Drawing updates are expected to stop when the window loses focus.
-(void)windowDidBecomeKey:(NSNotification *)notification {
    [self focus:devFocusIn]; // gained focus
}
-(void)windowDidResignKey:(NSNotification *)notification {
    [self focus:devFocusOut]; // lost focus
}

// Handle an orderly shutdown.
-(void)windowWillClose:(NSNotification *)notification {
    // The process that started NSApp will be killed.
    // FUTURE: look at adding some engine shutdown callbacks.
    //         to get similar behaviour with other platform shutdowns.
}

@end

//==============================================================================
// VuAppDelegate controls a windows lifecycle.
@interface VuAppDelegate : NSObject <NSApplicationDelegate> {}
@end

@implementation VuAppDelegate
-(void) applicationWillFinishLaunching:(NSNotification *)notification {
    // Needed to ensure window menu items work as expected.
    [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];
}
-(void) applicationDidFinishLaunching:(NSNotification *)notification {
    NSView *view = [[NSApp mainWindow] contentView];

    // Needed to pop the window to the front on launch.
    [NSApp activateIgnoringOtherApps:YES];
}

-(BOOL) applicationShouldTerminateAfterLastWindowClosed:(NSApplication *)sender {
    return YES;
}
@end

//==============================================================================

// create the display and view before calling dev_run().
// Returns a pointer to the CAMetalLayer.
long int dev_init(char * title, long x, long y, long w, long h) {
    // Create a shared app instance.
    // This will initialize the global variable
    // 'NSApp' with the application instance.
    [NSApplication sharedApplication];
    [NSApp setDelegate: [[VuAppDelegate new] autorelease]];

    // Create a window using the given location.
    NSRect windowRect = NSMakeRect(x, y, w, h);
    NSUInteger windowStyle =  NSWindowStyleMaskTitled | NSWindowStyleMaskClosable |
                              NSWindowStyleMaskMiniaturizable | NSWindowStyleMaskResizable;
    NSWindow * window = [[NSWindow alloc] initWithContentRect:windowRect
                        styleMask:windowStyle
                        backing:NSBackingStoreBuffered
                        defer:NO];
    app_window = window; // set global reference.
    [NSWindow setAllowsAutomaticWindowTabbing: NO];
    [window autorelease];
    NSWindowController * windowController = [[NSWindowController alloc] initWithWindow:window];
    [windowController autorelease];

    // Create the application menu programatically.
    // The app menu label will be the name of the executable.
    NSString *appName = [NSString stringWithUTF8String:title];
    NSMenu *menuBar = [[NSMenu new] autorelease];
    [NSApp setMainMenu: menuBar];

    // The main application menu assigns the default quit key Cmd-Q
    NSMenuItem *mi = [menuBar addItemWithTitle:@"Apple" action:NULL keyEquivalent:@""];
    NSMenu *m = [[NSMenu alloc] initWithTitle:@"Apple"];
    [NSApp performSelector:@selector(setAppleMenu:) withObject: m];
    [menuBar setSubmenu:m forItem:mi];
    [m addItemWithTitle:[NSString stringWithFormat:@"Quit %@", appName]
       action:@selector(terminate:) keyEquivalent:@"q"];

    // The view menu allows switching to full screen with Ctrl-Cmd-F
    mi = [menuBar addItemWithTitle:@"View" action:NULL keyEquivalent:@""];
    m = [[NSMenu alloc] initWithTitle:@"View"];
    [menuBar setSubmenu:m forItem:mi];
    [[m addItemWithTitle:@"Enter Full Screen" action:@selector(toggleFullScreen:) keyEquivalent:@"f"]
        setKeyEquivalentModifierMask:NSEventModifierFlagControl | NSEventModifierFlagCommand];

    // Add the view to the window.
    VuView* view = [[[VuView alloc] initWithFrame:windowRect] autorelease];
    view.clearColor = CLEAR_COLOR;
    [window setContentView:view];
    [window setDelegate:view];
    [window setTitle:appName];
    [window orderFrontRegardless];

    // Create the render delegate to get draw callbacks.
    VuRenderer *renderer = [VuRenderer alloc];
    view.delegate = renderer;
    return (long int)(view.layer); // return a pointer to the CAMetalLayer
}

// Expected to be called by the main thread. It does not return.
// Execution is now run from the OSX event loop and will cause
// callbacks as events and render updates are required.
void dev_run() {
    [NSApp run]; // Run event loop. Does not return until application exits.
}

// Update window size. No validation on values.
void dev_set_size(long x, long y, long w, long h) {
    if (app_window == 0) {
        return; // window not yet initialized.
    }
    NSRect frame = [app_window frame];
    NSRect content = [app_window contentRectForFrameRect:frame];
    CGFloat titleBarHeight = app_window.frame.size.height - content.size.height;
    CGSize windowSize = CGSizeMake(w, h + titleBarHeight);
    NSRect windowFrame = CGRectMake(x, y, windowSize.width, windowSize.height);
    [app_window setFrame:windowFrame display:YES animate:YES];
}

// Get current shell size.
void dev_size(long *x, long *y, long *w, long *h) {
    if (app_window == 0) {
        *x = 0;
        *y = 0;
        *w = 0;
        *h = 0;
        return; // window not yet initialized.
    }
    NSRect frame = [app_window frame];
    NSRect content = [app_window contentRectForFrameRect:frame];
    *x = (long)frame.origin.x;
    *y = (long)frame.origin.y;
    *w = (long)content.size.width;
    *h = (long)content.size.height;
}

// Update window title. Minimal effort is made to ensure a valid value.
void dev_set_title(char * value) {
    if (value != nil) {
        [[NSApp mainWindow] setTitle:[NSString stringWithUTF8String:value]];
    }
}

// Return 1 if the application is full screen, 0 otherwise.
// This needs to return the correct result right after a call
// to dev_toggle_fullscreen.
unsigned char dev_fullscreen() {
    NSWindow *window = [NSApp mainWindow];
    return (([window styleMask] & NSWindowStyleMaskFullScreen) == NSWindowStyleMaskFullScreen);
}

// Flip full screen mode.
void dev_toggle_fullscreen() {
    [[NSApp mainWindow] toggleFullScreen:nil];
}

// Get the cursor position.
void dev_cursor(long *x, long *y) {
    NSWindow *window = [NSApp mainWindow];
    NSPoint mouse = [window mouseLocationOutsideOfEventStream];
    *x = mouse.x;
    *y = mouse.y;
}

// Close down the application. Will cause call window terminate method.
void dev_dispose() { [NSApp terminate:nil]; }
