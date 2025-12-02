// SPDX-FileCopyrightText : © 2025 Galvanized Logic Inc.
// SPDX-License-Identifier: BSD-2-Clause

//go:build ignore
//
// Ignored because cgo attempts to compile it during normal builds.
//
// objective-c complile check can be run as follows:
// - export SDK=`xcrun --sdk iphoneos --show-sdk-path`
// - export CLANG=`xcrun --sdk iphoneos --find clang`
// - $CLANG -isysroot $SDK -framework Foundation -framework UIKit -framework Metal -framework MetalKit -framework QuartzCore -miphoneos-version-min=14.0 -o iosTest os_ios.m os_ios_test.m
//
// golang cross compile check can be run on macOS as follows:
// - export SDK=`xcrun --sdk iphoneos --show-sdk-path`
// - export CLANG=`xcrun --sdk iphoneos --find clang`
// - export IOS_FLAGS="-isysroot $SDK -arch arm64 -miphoneos-version-min=14.0"
// - env GOOS=ios GOARCH=arm64 CC=$CLANG CXX=$CLANG CGO_CFLAGS=$IOS_FLAGS CGO_LDFLAGS=$IOS_FLAGS CGO_ENABLED=1 go build
//
// The above only tests that everything compiles. To see console ouput
// create an XCode ios:objective-c project. Delete the existing project
// content and add the following files:
//
//   iosTest
//     iosTest
//        main.m      <- os_ios_test.m contents
//        os_ios.h
//        os_ios.m

#import <os/log.h>
#import "os_ios.h"

// Tests ios.
// Example Objective-C program that ensures the graphic shell works.
// This tests the native layer implmentation without golang.
int main(int argc, char *argv[]) {
    dev_log("Test ios device console output");
    dev_run(); // run the main event loop until it terminates.
    return 0;
}

// setView is called on startup to save a pointer to
// the Metal view layer.
void setView(long viewPointer) {
    os_log(OS_LOG_DEFAULT, "setView called %ld", viewPointer);
}

// renderFrame is called for the application to update its state
// and render a frame.
void renderFrame(void) {
    // dev_log("renderFrame called");
}

// handleInput reacts to user events.
void handleInput(long event, int d0, int d1) {
    if (event == devResized) {
        int w, h;
        dev_size(&w, &h);
        os_log(OS_LOG_DEFAULT, "handleInput resize called %d %d", w, h);
    } else if (event == devFocusIn) {
        dev_log("handleInput devFocus IN");
    } else if (event == devFocusOut) {
        dev_log("handleInput devFocus OUT");
    } else if (event == devTouchBegin) {
        os_log(OS_LOG_DEFAULT,  "handleInput devTouchBegin %d %d", d0, d1);
    } else if (event == devTouchMove) {
        os_log(OS_LOG_DEFAULT,  "handleInput devTouchMove %d %d", d0, d1);
    } else if (event == devTouchEnd) {
        os_log(OS_LOG_DEFAULT,  "handleInput devTouchEnd %d %d", d0, d1);
    } else {
        os_log(OS_LOG_DEFAULT,  "handleInput called %ld %d %d", event, d0, d1);
    }
}
