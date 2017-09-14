// Copyright © 2013-2017 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// The OSX (darwin) native layer implementation.
// This wraps the OSX API's (where the real work is done). See:
//    https://github.com/gamedevtech/CocoaOpenGLWindow
//    https://developer.apple.com/library/mac/qa/qa1385/_index.html
//    https://developer.apple.com/library/mac/documentation/GraphicsImaging/Conceptual/
//            OpenGL-MacProgGuide/opengl_designstrategies/opengl_designstrategies.html
//    https://lists.apple.com/archives/Mac-opengl/2010/Mar/msg00078.html
//
// Overall this is a nibless Cocoa window that expects its content to
// be rendered by OpenGL.

#import <Cocoa/Cocoa.h>
#import <QuartzCore/CVDisplayLink.h>
#import <OpenGL/OpenGL.h>
#import <OpenGL/gl3.h>
#include "os_darwin_amd64.h"

// The golang device callbacks that would normally be defined in _cgo_export.h
// are manually reproduced here and also implemented in the native test file
// os_darwin_amd64_test.m.
extern void prepRender();
extern void renderFrame();
extern void handleInput(long event, long data);

//==============================================================================

// VuView handles the OpenGL context and any window related events.
@interface VuView : NSOpenGLView <NSWindowDelegate> {
@public
    bool running;       // Engine is active and rendering frames.
    bool rendering;     // Currently generating one render frame.
    NSLock* appLock;    // Guards access from CVDisplayLink thread.
    CVDisplayLinkRef displayLink;
}
@end

@implementation VuView
-(BOOL)canBecomeKeyView { return YES; }
-(BOOL)acceptsFirstResponder { return YES; } // Window accepts input events

// Initialize creates the window frame and context pixel format.
-(id) initWithFrame: (NSRect) frame {
    NSOpenGLPixelFormatAttribute pixelAttrs[] = {
        NSOpenGLPFAOpenGLProfile, NSOpenGLProfileVersion3_2Core,
        NSOpenGLPFADepthSize,     24,
        NSOpenGLPFADoubleBuffer,
        0
    };
    NSOpenGLPixelFormat *pf = [[NSOpenGLPixelFormat alloc] initWithAttributes:pixelAttrs];
    if (pf == nil) {
        return nil;
    }
    self = [super initWithFrame:frame pixelFormat:[pf autorelease]];
    appLock = [[NSLock alloc] init];
    return self;
}

// prepareOpenGL creates the rendering context for the window.
-(void) prepareOpenGL {
   [super prepareOpenGL];
   [[self window] setLevel: NSNormalWindowLevel];
   [[self window] makeKeyAndOrderFront: self];

   // Make all the OpenGL calls to setup rendering and build the necessary
   // endering objects
   [[self openGLContext] makeCurrentContext];
   GLint swapInt = 1; // Synchronize buffer swaps with vertical refresh rate
   [[self openGLContext] setValues:&swapInt forParameter:NSOpenGLCPSwapInterval];

   // Create a display link capable of being used with all active displays
   CVDisplayLinkCreateWithActiveCGDisplays(&displayLink);

   // Set the renderer output callback function
   CVDisplayLinkSetOutputCallback(displayLink, &GlobalDisplayLinkCallback, self);
   CGLContextObj cglContext = (CGLContextObj)[ [self openGLContext] CGLContextObj];
   CGLPixelFormatObj cglPixelFormat = (CGLPixelFormatObj)[ [self pixelFormat] CGLPixelFormatObj];
   CVDisplayLinkSetCurrentCGDisplayFromOpenGLContext(displayLink, cglContext, cglPixelFormat);
}

// Continuous callback on main thread somewhat close to when a frame is
// requested by CVDisplayLink. Called from the main thread event
// after the last CVDisplayLink trigger.
-(void) renderFrame {
    NSAutoreleasePool *pool = [[NSAutoreleasePool alloc] init];
    [self setNeedsDisplay:YES];

    // Trigger application update to generate a new frame.
    // The renderFrame could result in calls to the dev_* functions below.
    renderFrame(); // Expected to call swap once done.
    [pool drain];
}

// getFrame is a CVDisplayLink callback triggered from a high priority thread
// based on the display refresh rate. Generally intended for video with the
// expectation that the display is ready for the next frame.
-(CVReturn) getFrame:(const CVTimeStamp*)outputTime {
    // Nothing is created on this thread so there is no autorelease pool.

    // Remember that a render request went out in order not to flood the main
    // thread with render requests. This is reset when a frame completes in
    // dev_swap(). Need a lock because this method is on a high performance
    // thread and dev_swap() is on the main thread.
    [appLock lock]; // guard rendering variable. See dev_swap.
    if (rendering) {
        // Don't wait around if the application is still busy and hasn't
        // finished rendering a frame. Just drop a frame and try again next
        // callback.
        // printf("dropping frame %llums\n", outputTime->hostTime/ 1000000);
    } else {
        rendering = true; // mark application as busy rendering.

        // Put the render request on the main thread to avoid having to lock
        // the OpenGL context. This also ensures the first renderFrame runs
        // after prepRender. It also prevents the application from running
        // its update in a high priority timer thread.
        dispatch_async(dispatch_get_main_queue(), ^{ [self renderFrame]; });
        // Equivalent call - need to profile if one is better than the other:
        // [self performSelectorOnMainThread:@selector(renderFrame)
        //                        withObject:nil waitUntilDone:NO ];
    }
    [appLock unlock];
    return kCVReturnSuccess;
}

// global handler for display refresh callbacks.
static CVReturn GlobalDisplayLinkCallback(CVDisplayLinkRef displayLink,
                const CVTimeStamp* now, const CVTimeStamp* outputTime,
                CVOptionFlags flagsIn, CVOptionFlags* flagsOut,
                void* displayLinkContext) {
    CVReturn result = [(VuView*)displayLinkContext getFrame:outputTime];
    return result;
}

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
- (void) keyUp: (NSEvent*) event          { handleInput(devUp, [event keyCode]); }
- (void) keyDown: (NSEvent*) event {
    if ([event isARepeat] == NO) {
        handleInput(devDown, [event keyCode]);
    }
}

// Get a single callback once resizing is finished.
// Master application is expected to update glViewport();
- (void)windowDidEndLiveResize:(NSNotification*)notification { handleInput(devResize, 0); }
- (void)windowDidMove:(NSNotification*)notification          { handleInput(devResize, 0); }


// focusEvent handles stopping the CVDisplay when the window
// loses focus.
-(void) focus: (int)focusEvent {
    [appLock lock];
    if (focusEvent == devFocusOut) {
        CVDisplayLinkStop(displayLink);
    }
    if (focusEvent == devFocusIn) {
        CVDisplayLinkStart(displayLink);
        rendering = false; // reset render flag.
    }
    [appLock unlock];
    handleInput(focusEvent, 0);
}

// Drawing updates are expected to stop when the window loses focus.
-(void)windowDidDeminiaturize:(NSNotification *)notification { [self focus:devFocusIn]; }
-(void)windowDidMiniaturize:(NSNotification *)notification   { [self focus:devFocusOut]; }
-(void)windowDidBecomeKey:(NSNotification *)notification     { [self focus:devFocusIn]; }
-(void)windowDidResignKey:(NSNotification *)notification     { [self focus:devFocusOut]; }

// Terminate window when the red X is pressed
-(void)windowWillClose:(NSNotification *)notification {
    if (running) {
        running = false;
        [appLock lock];
        CVDisplayLinkStop(displayLink);
        CVDisplayLinkRelease(displayLink);
        [appLock unlock];
        [NSApp terminate:self];
    }
}

// Cleanup
- (void) dealloc {
    [appLock release];
    [super dealloc];
}
@end

//==============================================================================

// VuAppDelegate ensures the application window opens and closes properly.
@interface VuAppDelegate : NSObject <NSApplicationDelegate> {}
@end

@implementation VuAppDelegate
-(void) applicationWillFinishLaunching:(NSNotification *)notification {
    // Needed to ensure window menu items work as expected.
    [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];
}
-(void) applicationDidFinishLaunching:(NSNotification *)notification {
    // Needed to pop the window to the front on launch.
    [NSApp activateIgnoringOtherApps:YES];
}
-(void) applicationDidBecomeActive:(NSNotification *)notification {
    NSView *view = [[NSApp mainWindow] contentView];
    if ( [view isKindOfClass:[VuView class]] ) {

        // Start CVDisplayLink after the window is ready to render.
        // Only do this once the first time the window is active.
        VuView * v = (VuView *)view;
        if (!v->running) {
            prepRender(); // intitial one time application callback.
            v->running = true;
            CVDisplayLinkStart(v->displayLink);
        }
    }
}
-(BOOL) applicationShouldTerminateAfterLastWindowClosed:(NSApplication *)sender {
    return YES; // Ensures application terminates when window is closed.
}
@end

//==============================================================================

// dev_run creates the the OSX objects needed to get an application window
// with an OpenGL context. It expects to be called by the main thread, and
// does not return. Execution is now run from the OSX event loop and will cause
// callbacks as events and render updates are required.
void dev_run() {
    // Autorelease Pool: Objects declared in this scope will be automatically
    //                   released at the end of it, when the pool is "drained".
    NSAutoreleasePool * pool = [[NSAutoreleasePool alloc] init];

    // Create a shared app instance.
    // This will initialize the global variable
    // 'NSApp' with the application instance.
    [NSApplication sharedApplication];
    [NSApp setDelegate: [[VuAppDelegate new] autorelease]];

    // Create a window:
    NSRect screenRect = [[NSScreen mainScreen] frame];
    NSRect bounds = NSMakeRect(0, 0, 640, 480);
    NSRect windowRect = NSMakeRect(NSMidX(screenRect) - NSMidX(bounds),
                                 NSMidY(screenRect) - NSMidY(bounds),
                                 bounds.size.width, bounds.size.height);
    NSUInteger windowStyle =  NSWindowStyleMaskTitled | NSWindowStyleMaskClosable |
                              NSWindowStyleMaskMiniaturizable | NSWindowStyleMaskResizable;
    NSWindow * window = [[NSWindow alloc] initWithContentRect:windowRect
                        styleMask:windowStyle
                        backing:NSBackingStoreBuffered
                        defer:NO];
    [NSWindow setAllowsAutomaticWindowTabbing: NO];
    [window autorelease];
    NSWindowController * windowController = [[NSWindowController alloc] initWithWindow:window];
    [windowController autorelease];

    // Create the application menu programatically.
    NSString *appName = @"App"; // startup title. Set later by application.
    NSMenu *menuBar = [[NSMenu new] autorelease];
    [NSApp setMainMenu: menuBar];

    // the main application menu assigns the default quit key Cmd-Q
    NSMenuItem *mi = [menuBar addItemWithTitle:@"Apple" action:NULL keyEquivalent:@""];
    NSMenu *m = [[NSMenu alloc] initWithTitle:@"Apple"];
    [NSApp performSelector:@selector(setAppleMenu:) withObject: m];
    [menuBar setSubmenu:m forItem:mi];
    [m addItemWithTitle:[NSString stringWithFormat:@"Quit %@", appName]
       action:@selector(terminate:) keyEquivalent:@"q"];

    // the view menu allows switching to full screen with Ctrl-Cmd-F
    mi = [menuBar addItemWithTitle:@"View" action:NULL keyEquivalent:@""];
    m = [[NSMenu alloc] initWithTitle:@"View"];
    [menuBar setSubmenu:m forItem:mi];
    [[m addItemWithTitle:@"Enter Full Screen" action:@selector(toggleFullScreen:) keyEquivalent:@"f"]
        setKeyEquivalentModifierMask:NSEventModifierFlagControl | NSEventModifierFlagCommand];

    // Create app delegate to handle system events
    VuView* view = [[[VuView alloc] initWithFrame:windowRect] autorelease];
    [window setContentView:view];
    [window setDelegate:view];
    [window setTitle:appName];
    [window orderFrontRegardless];
    [NSApp run]; // Run event loop. Does not return until application exits.
    [pool drain];
}

// Flip OpenGL render buffers. Context always double buffered.
// Marks frame as having been completed.
void dev_swap() {
    NSView *view = [[NSApp mainWindow] contentView];
    if ( [view isKindOfClass:[VuView class]] ) {
        VuView * v = (VuView *)view;
        [v->appLock lock];
        v->rendering = false;
        [v->appLock unlock];
        [[v openGLContext] flushBuffer];
    }
}

// Update window size. No validation on values.
void dev_set_size(long x, long y, long w, long h) {
    NSWindow *window = [NSApp mainWindow];
    NSRect frame = [window frame];
    NSRect content = [window contentRectForFrameRect:frame];
    CGFloat titleBarHeight = window.frame.size.height - content.size.height;
    CGSize windowSize = CGSizeMake(w, h + titleBarHeight);
    NSRect windowFrame = CGRectMake(x, y, windowSize.width, windowSize.height);
    [window setFrame:windowFrame display:YES animate:YES];
}

// Get current shell size.
void dev_size(long *x, long *y, long *w, long *h) {
    NSWindow *window = [NSApp mainWindow];
    NSRect frame = [window frame];
    NSRect content = [window contentRectForFrameRect:frame];
    *x = frame.origin.x;
    *y = frame.origin.y;
    *w = content.size.width;
    *h = content.size.height;
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

// Position the cursor at the given window location.
// Ensure screen coordinate origin is bottom left.
void dev_set_cursor_location(long x, long y) {
    NSWindow *window = [NSApp mainWindow];
    NSRect windowRect = [window frame];
    CGRect screenRect = CGDisplayBounds(CGMainDisplayID());

    // Flip y to get screen 0,0 at bottom left.
    CGPoint point = CGPointMake(windowRect.origin.x+x, screenRect.size.height-windowRect.origin.y-y);
    CGWarpMouseCursorPosition(point);
    CGAssociateMouseAndMouseCursorPosition(true);
}

// Show or hide cursor. Note that the cursor is not locked with
// CGAssociateMouseAndMouseCursorPosition(show) or the cursor position
// would not change.
void dev_show_cursor(unsigned char show) {
    if (show) {
        [NSCursor unhide];
    } else {
        [NSCursor hide];
    }
}

// Get the cursor position.
void dev_cursor(long *x, long *y) {
    NSWindow *window = [NSApp mainWindow];
    NSPoint mouse = [window mouseLocationOutsideOfEventStream];
    *x = mouse.x;
    *y = mouse.y;
}

// Return the current clipboard contents if the clipboard contains text.
// Otherwise return nil. Any returned strings must be freed by the caller.
char* dev_clip_copy() {
    NSPasteboard* pb = [NSPasteboard generalPasteboard];
    if (![[pb types] containsObject:NSStringPboardType]) {
        return NULL; // only deal with strings.
    }
    NSString* object = [pb stringForType:NSStringPboardType];
    if (!object) {
        return NULL; // only handle non-nil strings.
    }
    return strdup([object UTF8String]); // must be freed by caller.
}

// Paste the given string into the general clipboard.
void dev_clip_paste(const char* string) {
    NSArray* types = [NSArray arrayWithObjects:NSStringPboardType, nil];
    NSPasteboard* pb = [NSPasteboard generalPasteboard];
    [pb declareTypes:types owner:nil];
    [pb setString:[NSString stringWithUTF8String:string] forType:NSStringPboardType];
}

// Close down the application. Will cause call window terminate method.
void dev_dispose() { [NSApp terminate:nil]; }
