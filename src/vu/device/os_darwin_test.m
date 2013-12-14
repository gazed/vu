// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

#include "os_darwin.m"
#import <stdio.h>
#import <OpenGL/gl3.h> // Needed for the OpenGL gl* functions.

// An native test application can be built using:
//     clang -framework Cocoa -framework OpenGL -o App os_darwin_test.m
// This gives an idea of the expected usage exposed gs_* functions.
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
