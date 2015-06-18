// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

#include "os_linux.h"

// FUTURE
// Example C program that ensures the graphic shell works.
// Compile this on git bash (mingw) using:
//     gcc -o App os_windows*.c -lopengl32 -lgdi32 -Wall -m64
// This is commented out because cgo attempts to compile it when building
// the go binary and will not like having duplicate main methods.
// To test: remove the comments to build using the gcc command above.
// int main( void )
// {
//     // initialize the window.
//     gs_set_attr_s(GS_AppName, "Test Window");
//     long long display = gs_display_init();
//     if (display == 0)
//     {
//         printf("Failed display init.\n");
//         exit( EXIT_FAILURE );
//     }
//
//     // create the window.
//     long long shell = gs_shell(display);
//     printf("display %ld shell %ld\n", (long)display, (long)shell);
//     long long context = gs_context(&display, &shell);
//     printf("display %ld shell %ld context %ld\n", (long)display, (long)shell, (long)context);
//     if (context == 0)
//     {
//         printf("Failed context init.\n");
//         exit( EXIT_FAILURE );
//     }
//     gs_shell_open(display);
//  long x, y, w, h;
//     gs_size(display, &x, &y, &w, &h);
//     printf("shell size %ld::%ld ", w, h);
//
//     // process user events.
//     GSEvent gsu = {-1, 0, 0, 0, 0, 0};
//     while (gs_shell_alive(display))
//     {
//         gsu.event = -1;
//         gsu.mousex = 0;
//         gsu.mousey = 0;
//         gsu.key = 0;
//         gsu.mods = 0;
//         gsu.scroll = 0;
//         gs_read_dispatch(display, &gsu);
//         if (gsu.event > 0)
//         {
//             printf("event %ld mouse %ld,%ld - ", gsu.event, gsu.mousex, gsu.mousey);
//             if (gsu.event == GS_KeyDown || gsu.event == GS_KeyUp)
//             {
//                 printf("keys 0x%.2lX:0x%.2lX - ", gsu.mods, gsu.key);
//             }
//             else if (gsu.event == GS_ScrollWheel)
//             {
//                 printf("wheel %ld - ", gsu.scroll);
//             }
//             printf("\n");
//         }
//         gs_swap_buffers(shell);
//     }
//     gs_display_dispose(display);
//  return 0;
// };
