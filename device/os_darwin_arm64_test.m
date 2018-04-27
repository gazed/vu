// Copyright © 2017-2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// +build ignore
//
// Ignored because cgo attempts to compile it during normal builds.
// Here is the command line to cross-compile an ios test application:
//     clang -isysroot /Applications/Xcode.app/Contents/Developer/Platforms/iPhoneOS.platform/Developer/SDKs/iPhoneOS11.2.sdk -framework Foundation -framework UIKit -framework GLKit -framework OpenGLES -framework AVFoundation -arch arm64 -miphoneos-version-min=8.0 -o iosApp os_darwin_arm64_test.m
// The above just creates the executable. In order to test the executable it has
// to be packaged with xcodebuild into an app with valid signing certificates.

#import <stdio.h>
#import <OpenGLES/ES3/gl.h>
#import <OpenGLES/ES3/glext.h>

#import "os_darwin_arm64.m"

// Tests darwin amm64 OSX.
// Example Objective-C program that ensures the graphic shell works.
// This tests the native layer implmentation without golang.
int main(int argc, char *argv[]) {
    NSLog(@"%@", @"Test darwin_arm64 ios device console output");
    dev_run(); // run the main event loop until it terminates.
    return 0;
}

// prepRender is called one time after the application opens and
// the drawing context has been initialized.
void prepRender() {
    int w, h, scale;
    dev_size(&w, &h, &scale);
    NSLog(@"%@ %i %i %i", @"prepRender called", w, h, scale);
}

// renderFrame is called for the application to update its state
// and render a frame.
void renderFrame() {
    NSLog(@"%@", @"   render called");
    glClearColor(0.0f, 0.8f, 0.1f, 1.0f);
    glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT);
}

// handleInput reacts to user events.
void handleInput(long event, int d0, int d1) {
    if (event == devResize) {
        int w, h, scale;
        dev_size(&w, &h, &scale);
        NSLog(@"%@ %i %i %i", @"handleInput resize called", w, h, scale);
    } else if (event == devFocus) {
        NSLog(@"%@ %i", @"handleInput devFocus", d0);
    } else if (event == devTouchBegin) {
        NSLog(@"%@ %i %i", @"handleInput devTouchBegin", d0, d1);
    } else if (event == devTouchMove) {
        NSLog(@"%@ %i %i", @"handleInput devTouchMove", d0, d1);
    } else if (event == devTouchEnd) {
        NSLog(@"%@ %i %i", @"handleInput devTouchEnd", d0, d1);
    } else {
        NSLog(@"%@ %li %i %i", @"handleInput called", event, d0, d1);
    }
}
