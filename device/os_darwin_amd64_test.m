// Copyright © 2013-2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// +build ignore
//
// Ignored because cgo attempts to compile it during normal builds.
// Use the following to build a native test application:
//     clang -o osxApp -framework Cocoa -framework Quartz -framework OpenGL os_darwin_amd64_test.m

#import <stdio.h>
#import "os_darwin_amd64.m"

// Tests darwin amd64 OSX.
// Example Objective-C program that ensures the graphic shell works.
// This tests the native layer implmentation without golang.
int main(int argc, char *argv[]) {
    dev_run(); // Does not return. Calls initialize() and update().
    return 0;
}

// prepRender is called one time after the application opens and
// the drawing context has been initialized.
void prepRender() {
    dev_set_title("osxApp");
    dev_set_size(0, 0, 600, 400);
    printf("%s %s\n", glGetString(GL_RENDERER), glGetString(GL_VERSION));
    printf("%s\n", glGetString(GL_SHADING_LANGUAGE_VERSION));
    glClearColor(0.0, 0.3, 0.3, 1.0);
}

// renderFrame is called for the application to update its state
// and render a frame.
void renderFrame() {
    glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT);
    dev_swap();
}

void handleInput(long event, long data) {
    if (event == devDown) {
        if (data == 0x08) { // c key
            char *s = dev_clip_copy();
            printf(" \"%s\"\n", s);
            free(s);
        } else if (data == 0x23) { // p key
            dev_clip_paste("test paste string");
        } else if (data == 0x11) { // t key
            dev_toggle_fullscreen();
        } else if (data == devMouseL) { // left click
            printf("left mouse click\n");
        }
    } else if (event == devUp) {
        printf("release %ld\n", data);
    } else if (event == devScroll) {
        printf("scroll %ld\n", data);
    } else if (event == devMod) {
        printf("modifiers %ld\n", data);
    } else {
        printf("event %ld\n", event);
    }
}
