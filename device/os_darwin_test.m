// Copyright © 2013-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// +build ignore
//
// Ignored because cgo attempts to compile it during normal builds.
// Here is the command line to build a native test application:
//     clang -framework Cocoa -framework OpenGL -o App os_darwin_test.m

#import <stdio.h>
#import <OpenGL/gl3.h> // Needed for the OpenGL gl* functions.

// Example Objective-C program that ensures the graphic shell works.
// This tests the native layer implmentation without golang.
#include "os_darwin.m"
int main(int argc, char *argv[]) {
    long display = gs_display_init();
    gs_set_attr_s(GS_AppName, "Gsh");
    gs_set_attr_l(GS_ShellX, 600);
    gs_set_attr_l(GS_ShellY, 400);
    gs_set_attr_l(GS_ShellWidth, 800);
    gs_set_attr_l(GS_ShellHeight, 600);
    long shell = gs_shell(display);

    // initialize the OpenGL context
    long context = gs_context(shell);
    printf("%s %s\n", glGetString(GL_RENDERER), glGetString(GL_VERSION));
    printf("%s\n", glGetString(GL_SHADING_LANGUAGE_VERSION));
    GLuint vao;
    glGenVertexArrays(1, &vao);
    printf("GenVertexArrays 0x%X vao=%d\n", glGetError(), vao);
    glClearColor(0.0, 0.3, 0.3, 1.0);

    // run the main event loop.
    gs_shell_open(display);
    GSEvent gsu = {-1, -1, -1, 0, 0, 0};
    do {
        // do some OpenGL drawing.
        glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT);

        // process user events.
        gsu.event = -1;
        gs_read_dispatch(display, &gsu);
        if (gsu.event >= 0) {

            // show current key code
            printf("mouse %ld,%ld - ", gsu.mousex, gsu.mousey);
            printf("[mods  0x%.2lX] - ", gsu.mods);
            if (gsu.event == GS_KeyUp) {
                printf("key up   0x%.2lX - ", gsu.key);
            } else if (gsu.event == GS_KeyDown) {
                printf("key down 0x%.2lX - ", gsu.key);

                // test copy and paste.
                if (gsu.key == 0x08) { // c key
                    char *s = gs_clip_copy();
                    printf(" \"%s\"", s);
                    free(s);
                } else if (gsu.key == 0x23) { // p key
                    gs_clip_paste("test paste string");
                } else if (gsu.key == 0x11) { // t key: Test fullscreen toggle.
                   printf(" toggle before %d", gs_fullscreen(display));
                   gs_toggle_fullscreen(display);
                   printf(" toggle after %d", gs_fullscreen(display));
                }
            } else if (gsu.event == GS_ScrollWheel) {
                printf("wheel %ld - ", gsu.scroll);
            }  else {
                printf("event %ld - ", gsu.event);
            }
            printf("\n");
        }
        gs_swap_buffers(context);
    } while (gs_shell_alive(shell));
    gs_display_dispose(display);
    return 0;
}
