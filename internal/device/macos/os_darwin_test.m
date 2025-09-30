// Copyright © 2025 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// +build ignore
//
// Ignored because cgo attempts to compile it during normal builds.
// Use the following to build a native test application:
//     clang -o OSXApp -framework Cocoa -framework Metal -framework MetalKit -framework QuartzCore os_darwin_test.m

#import <stdio.h>
#import "os_darwin.m"

// Test the native layer implmentation without golang.
// If this works then a metal view and window has been created.
int main(int argc, const char *argv[]) {
    dev_init("OSXApp", 100, 100, 1280, 720); // Create display and view
    dev_run();                               // Does not return - the callbacks below will start getting called.
    return 0;
}

// renderFrame callbacks are very frequent and ignored.
void renderFrame() { }

// handleInput dumps user input.
void handleInput(long event, long data) {
    long x, y, w, h;
    if (event == devDown) {
        if (data == 0x11) { // t key or cmd-ctrl-F
            printf("toggle fullscreen\n");
            dev_toggle_fullscreen();
        } else if (data == 0x01) { // s key
            dev_size(&x, &y, &w, &h);
            printf("surface size x:%ld y:%ld w:%ld h:%ld\n", x, y, w, h);
        } else if (data == devMouseL) { // left click
            printf("left mouse click\n");
        } else {
            printf("press %ld\n", data);
        }
    } else if (event == devUp) {
        printf("release %ld\n", data);
    } else if (event == devScroll) {
        printf("scroll %ld\n", data);
    } else if (event == devMod) {
        printf("modifiers %ld\n", data);
    } else if (event == devFocusIn) {
        printf("focus gained %ld\n", data);
    } else if (event == devFocusOut) {
        printf("focus lost %ld\n", data);
    } else if (event == devResized) {
        long x, y, w, h;
        dev_size(&x, &y, &w, &h);
        printf("resized x:%ld y:%ld w:%ld h:%ld\n", x, y, w, h);
    } else if (event == devMoved) {
        long x, y, w, h;
        dev_size(&x, &y, &w, &h);
        printf("moved x:%ld y:%ld w:%ld h:%ld\n", x, y, w, h);
    } else {
        printf("event %ld\n", event);
    }
}
