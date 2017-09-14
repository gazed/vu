// Copyright © 2017 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// The iOS (darwin) native layer implementation.
// This wraps the iOS API's (where the real work is done).

#include <stdio.h>

#import <UIKit/UIKit.h>
#import <GLKit/GLKit.h>
#import "os_darwin_arm64.h"

// The golang device callbacks that would normally be defined in _cgo_export.h
// are manually reproduced here and also implemented in the native test file
// os_darwin_arm64_test.m.
extern void prepRender();
extern void renderFrame();
extern void handleInput(long event, int d0, int d1);

// Declare VuController first because it is as a delegate helper in VuDelegate.
@interface VuController : GLKViewController<UIContentContainer, GLKViewDelegate>
@property (strong, nonatomic) EAGLContext *context;
@property (strong, nonatomic) GLKView *glview;
@property bool running;       // Engine is active and rendering frames.
@end

// =============================================================================
// VuDelegate

// VuDelegate is the primary handler for main loop events.
// It is handed to the main loop controller once on startup.
// See UIApplicationMain below.
@interface VuDelegate : UIResponder<UIApplicationDelegate>
@property (strong, nonatomic) UIWindow *window;
@property (strong, nonatomic) VuController *controller;
@end

// VuDelegate receives callbacks from the main event loop.
// It mainly cares about reporting when the app is active and visible.
// iOS apps can only make OpenGL calls when they are active and visible
// or the app is force terminated.
@implementation VuDelegate
- (BOOL)application:(UIApplication *)application didFinishLaunchingWithOptions:(NSDictionary *)launchOptions {
    self.window = [[UIWindow alloc] initWithFrame:[[UIScreen mainScreen] bounds]];
    self.controller = [[VuController alloc] initWithNibName:nil bundle:nil];
    self.window.rootViewController = self.controller;
    [self.window makeKeyAndVisible];
    return YES;
}

- (void)applicationDidBecomeActive:(UIApplication * )application {
    handleInput(devFocus, 0, 0);
}
- (void)applicationWillResignActive:(UIApplication *)application {
    handleInput(devFocus, 1, 0);
}
- (void)applicationDidEnterBackground:(UIApplication *)application {
    handleInput(devFocus, 1, 0);
}
- (void)applicationWillTerminate:(UIApplication *)application {
    handleInput(devFocus, 1, 0);
}
@end

// VuDelegate
// =============================================================================
// VuController

@implementation VuController
- (void)viewDidLoad
{
    [super viewDidLoad];
    self.context = [[EAGLContext alloc] initWithAPI:kEAGLRenderingAPIOpenGLES3];
    self.glview = (GLKView*)self.view;
    self.glview.drawableDepthFormat = GLKViewDrawableDepthFormat24;
    self.glview.multipleTouchEnabled = true;
    self.glview.context = self.context;
    self.glview.userInteractionEnabled = YES;
    self.glview.enableSetNeedsDisplay = YES; // only invoked once

    // FUTURE: profile the various devices for defaul values.
    //         Would like 60 minimum. Is the ipad pro is upto 120?
    // self.preferredFramesPerSecond = 30;
    //
    // FUTURE: Enable multisampling? Check performance.
    // self.glview.drawableMultisample = GLKViewDrawableMultisample4X;
    NSLog(@"%@", @"viewDidLoad");
}

- (void)viewWillTransitionToSize:(CGSize)size withTransitionCoordinator:(id<UIViewControllerTransitionCoordinator>)coordinator
{
    [coordinator animateAlongsideTransition:^(id<UIViewControllerTransitionCoordinatorContext> context) {
    } completion:^(id<UIViewControllerTransitionCoordinatorContext> context) {
        handleInput(devResize, 0, 0);
    }];
    NSLog(@"%@ %i %i", @"viewWillTransitionToSize", (int)size.width, (int)size.height);
}

// update is called prior to each drawInRect render frame.
//
// FUTURE: find a better spot to call the one-time prepRender.
//         The issue is that the system doesn't accept openGL calls
//         until this point in time.
- (void) update
{
    if (!self.running) {
        prepRender();        // intitial application callback.
        self.running = true; // only first time window is active.
    }
}

// drawInRect is requesting a new display frame. GLKit handles presenting the
// frame to the device when this returns.
//
// FUTURE: another option is to use CADisplayLink for timing if the
//         GLKViewController callbacks cause timing problems. Other issues
//         need to be dealt with: see:
// https://stackoverflow.com/questions/10080932/glkviewcontrollerdelegate-getting-blocked
// https://stackoverflow.com/questions/5944050/cadisplaylink-opengl-rendering-breaks-uiscrollview-behaviour
- (void)glkView:(GLKView *)view drawInRect:(CGRect)rect
{
    renderFrame();
}

// sendTouches inverts the y so 0,0 is at the bottom left.
static void sendTouches(int change, NSSet* touches)
{
    CGFloat scale = [UIScreen mainScreen].scale;
    CGSize size = [UIScreen mainScreen].bounds.size;
    for (UITouch* touch in touches) {
        CGPoint p = [touch locationInView:touch.view];
        NSLog(@"%@ touchId:%i x:%i y:%i scale:%i", @"sendTouches",
              change, (int)(p.x*scale), (int)(p.y*scale), (int)scale);
        handleInput(change, p.x*scale, (size.height-p.y)*scale);
    }
}
- (void)touchesBegan:(NSSet*)touches withEvent:(UIEvent*)event {
    sendTouches(devTouchBegin, touches);
}
- (void)touchesMoved:(NSSet*)touches withEvent:(UIEvent*)event {
    sendTouches(devTouchMove, touches);
}
- (void)touchesEnded:(NSSet*)touches withEvent:(UIEvent*)event {
    sendTouches(devTouchEnd, touches);
}
- (void)touchesCanceled:(NSSet*)touches withEvent:(UIEvent*)event {
    sendTouches(devTouchEnd, touches);
}
@end

// VuController
// =============================================================================
// methods for golang calling into iOS

// This is called on startup to hand control to the UIKit framework as per
// the standard iOS main. Control is returned using callbacks like drawInRect.
void dev_run(void)
{
    @autoreleasepool {
        UIApplicationMain(0, nil, nil, NSStringFromClass([VuDelegate class]));
    }
}

// dev_size gets display information in pixels.
// Orientation switches are handled by the OS.
void dev_size(int *w, int *h, int *scale)
{
    CGFloat s = [UIScreen mainScreen].scale; // 1.0, 2.0, or 3.0.
    *scale = (int)s;
    CGSize size = [UIScreen mainScreen].bounds.size;
    *w = *scale * (int)size.width;
    *h = *scale * (int)size.height;

    // FUTURE: check if application cares about orientation changes.
    // UIInterfaceOrientation orientation = [[UIApplication sharedApplication] statusBarOrientation];
}

// dev_dispose is called when the application terminates...
// Generally not done in iOS.
// FUTURE: remove this method.
void dev_dispose() {}

// dev_log outputs application log to the device console.
void dev_log(const char* log)
{
    NSLog(@"%s", log);
}
