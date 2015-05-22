// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// +build ignore
//
// Ignored because cgo attempts to compile it during normal builds.
// To build a native test application, compile this on git bash (mingw) using:
//     gcc -o App os_windows*.c -lopengl32 -lgdi32 -Wall -m64

#include "os_windows.h"

// Example C program that ensures the graphic shell works.
// This tests the native layer implmentation without golang.
int main( void )
{
    // initialize the window.
    gs_set_attr_s(GS_AppName, "Test Window");
    long long display = gs_display_init();
    if (display == 0)
    {
        printf("Failed display init.\n");
        exit( EXIT_FAILURE );
    }

    // create the window.
    long long shell = gs_shell(display);
    printf("display %ld shell %ld\n", (long)display, (long)shell);
    long long context = gs_context(&display, &shell);
    printf("display %ld shell %ld context %ld\n", (long)display, (long)shell, (long)context);
    if (context == 0)
    {
        printf("Failed context init.\n");
        exit( EXIT_FAILURE );
    }
    gs_shell_open(display);
    long x, y, w, h;
    gs_size(display, &x, &y, &w, &h);
    printf("shell size %ld::%ld ", w, h);

    // process user events.
    GSEvent gsu = {-1, 0, 0, 0, 0, 0};
    while (gs_shell_alive(display))
    {
        gsu.event = -1;
        gsu.mousex = 0;
        gsu.mousey = 0;
        gsu.key = 0;
        gsu.mods = 0;
        gsu.scroll = 0;
        gs_read_dispatch(display, &gsu);
        if (gsu.event >= 0) {

            // show current key code
            printf("mouse %ld,%ld - ", gsu.mousex, gsu.mousey);
            printf("[mods  0x%.2lX] - ", gsu.mods);
            if (gsu.event == GS_KeyUp) {
                printf("key up   0x%.2lX - ", gsu.key);
            } else if (gsu.event == GS_KeyDown) {
                printf("key down 0x%.2lX - ", gsu.key);
            } else if (gsu.event == GS_ScrollWheel) {
                printf("wheel %ld - ", gsu.scroll);
            }  else {
                printf("event %ld - ", gsu.event);
            }
            printf("\n");
        }
        gs_swap_buffers(shell);
    }
    gs_display_dispose(display);
    return 0;
};
