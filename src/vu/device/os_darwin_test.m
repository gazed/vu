// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

#import <stdio.h>
#import "os_darwin.h"
#import <OpenGL/gl3.h> // Needed for the OpenGL gl* functions.

// Used to test the dylib generated from gs_darwin.[hc]
// This also gives an idea of the expected usage exposed gs_* functions.
//
// A local version of the dynamic library can be built using:
//     clang -dynamiclib -fno-common -framework Cocoa -o libvudev.1.dylib os_darwin.m
// An test application (expected usage example) can be built using:
//     clang -I. libvudev.1.dylib -framework Cocoa -framework OpenGL -o App os_darwin_test.m
//
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
            if (gsu.event == GS_KeyDown || gsu.event == GS_KeyUp) {
                printf("keys 0x%.2lX:0x%.2lX - ", gsu.mods, gsu.key);
            } else if (gsu.event == GS_ScrollWheel) {
		    	printf("wheel %ld - ", gsu.scroll);
		    } 
		    printf("\n");
		}
        gs_swap_buffers(context);
    } while (gs_shell_alive(shell));
    gs_display_dispose(display);
    return 0;
}
