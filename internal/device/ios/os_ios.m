// Copyright © 2025 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// The iOS (darwin) native layer implementation.
// This wraps the iOS API's (where the real work is done).

#import <os/log.h>
#import <MetalKit/MetalKit.h>
#import <QuartzCore/CAMetalLayer.h>
#import <UIKit/UIKit.h>

#import "os_ios.h"

// The golang device callbacks that would normally be defined in _cgo_export.h
// are manually reproduced here and also implemented in the native test file
// os_ios_test.m.
extern void renderFrame(void);
extern void handleInput(long event, int d0, int d1);
extern void setView(long viewPointer);

// =============================================================================
// VuView is a Metal compatible view.
@interface VuView : MTKView
@end

@implementation VuView

// Returns a Metal-compatible layer.
+(Class) layerClass { return [CAMetalLayer class]; }

// sendTouches forwards the touch input.
static void sendTouches(int change, NSSet* touches)
{
    CGFloat scale = [UIScreen mainScreen].nativeScale;
    for (UITouch* touch in touches) {
        CGPoint p = [touch locationInView:touch.view];
        handleInput(change, p.x*scale, p.y*scale);
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

// =============================================================================
// VuController
@interface VuController : UIViewController
@property (nonatomic, strong) VuView *vuView;
@property (nonatomic, strong) CADisplayLink *displayLink;
@end

@implementation VuController

- (void)viewDidLoad {
    [super viewDidLoad];
    self.view.userInteractionEnabled = YES;
    // FUTURE handle multi-touch using
    // self.view.multipleTouchEnabled = YES; // default is NO

    // create the metal MTKView view.
    self.vuView = [[VuView alloc] initWithFrame:[UIScreen mainScreen].bounds];
    // self.vuView.delegate = self; TODO needed if using CADisplayLink?
    [self.view addSubview: self.vuView];

    // set contentsScale for CAMetalLayer
    CGFloat scale = [UIScreen mainScreen].nativeScale;
    self.vuView.layer.contentsScale = scale;

    // initialize the render timer.
    self.displayLink = [CADisplayLink displayLinkWithTarget:self selector:@selector(render:)];
    [self.displayLink addToRunLoop:[NSRunLoop mainRunLoop] forMode:NSDefaultRunLoopMode];

    // return a pointer to the CAMetalLayer
    setView((long)self.vuView.layer);
}

// renders a frame by calling back to the application engine.
- (void)render:(CADisplayLink *)sender {
    renderFrame();
}

// Called for screen rotations.
- (void) viewWillTransitionToSize:(CGSize) size
        withTransitionCoordinator:(id<UIViewControllerTransitionCoordinator>) coordinator {
    [super viewWillTransitionToSize:size withTransitionCoordinator:coordinator];
    handleInput(devResized, 0, 0);
}

-(void) viewDidDisappear: (BOOL) animated {
	[self.displayLink invalidate];
	[self.displayLink release];
	[super viewDidDisappear: animated];
}

@end

// =============================================================================
// SceneDelegate
@interface SceneDelegate : UIResponder <UIWindowSceneDelegate>
@property (strong, nonatomic) UIWindow *window;
@end

@implementation SceneDelegate

@synthesize window = _window;

- (void)scene:(UIScene *)scene willConnectToSession:(UISceneSession *)session options:(UISceneConnectionOptions *)connectionOptions {
    self.window = [[UIWindow alloc] initWithWindowScene:(UIWindowScene *)scene];
    self.window.rootViewController = [[VuController alloc] init];
    self.window.hidden = NO;
    self.window.userInteractionEnabled = YES;
    self.window.rootViewController.view.backgroundColor = [UIColor blackColor];
    [self.window makeKeyAndVisible];
}

- (void)sceneDidBecomeActive:(UIScene *)scene {
    handleInput(devFocusIn, 0, 0);
}
- (void)sceneWillResignActive:(UIScene *)scene {
    handleInput(devFocusOut, 0, 0);
}

@end

// =============================================================================
// VuDelegate
@interface VuDelegate : UIResponder<UIApplicationDelegate>
@end

// VuDelegate receives callbacks from the main event loop.
@implementation VuDelegate

- (UISceneConfiguration *)application:(UIApplication *)application configurationForConnectingSceneSession:(UISceneSession *)connectingSceneSession options:(UISceneConnectionOptions *)options {
    UISceneConfiguration *configuration = [[UISceneConfiguration alloc] initWithName:nil sessionRole:connectingSceneSession.role];
    configuration.delegateClass = SceneDelegate.class;
    return configuration;
}

@end

// =============================================================================
// golang binding methods.

// This is called on startup to hand control to the UIKit framework as per
// the standard iOS main. Control is returned using callbacks.
void dev_run(void) {
    int argc = 1;
    char *argv[] = {""};
    @autoreleasepool {
        UIApplicationMain(argc, argv, nil, NSStringFromClass([VuDelegate class]));
    }
}

// dev_size gets display information in pixels.
// Orientation switches are handled by the OS.
void dev_size(int *w, int *h) {
    CGSize size = [[UIScreen mainScreen] nativeBounds].size;
    *w = (int)size.width;
    *h = (int)size.height;
}

// dev_log outputs application log to the device console.
void dev_log(const char* log) {
    os_log(OS_LOG_DEFAULT, "%{public}s", log);
}
