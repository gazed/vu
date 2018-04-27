// Copyright © 2013-2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// +build ignore
//
// Ignored because cgo attempts to compile it during normal builds.
// To build a native test application, compile this on git bash (mingw) using:
//     gcc -o winApp os_windows.c os_windows_test.c -lopengl32 -lgdi32 -Wall -m64

#include <stdio.h>
#include "os_windows.h"

// Needed to get console debug out from a windowed app.
// Is there a more recent and easier way to do this?
// From http://dslweb.nwnexus.com/~ast/dload/guicon.htm
#include <fcntl.h>
#include <io.h>
void RedirectIOToConsole();

// Tests windows native library.
// Example C program that ensures the graphic shell works.
// This tests the native layer implmentation without golang.
int main()
{
    RedirectIOToConsole();
    dev_run(); // Does not return. Calls prepRender() and renderFrame().
    return 0;
}

// prepRender is called one time after the application opens and
// the drawing context has been initialized.
void prepRender()
{
    dev_set_title("Test Window");
    dev_set_size(600, 200, 600, 400);
    long x, y, w, h;
    dev_size(&x, &y, &w, &h);
    printf("windows size %ld %ld %ld %ld\n", x, y, w, h);
}

// renderFrame is called for the application to update its state
// and render a frame.
void renderFrame()
{
    dev_swap();
}

// handleInput is called as user events occur.
void handleInput(long event, long data)
{
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
        } else {
            printf("press %ld\n", data);
        }
    } else if (event == devUp) {
        printf("release %ld\n", data);
    } else if (event == devScroll) {
        printf("scroll %ld\n", data);
    } else {
        printf("event %ld\n", event);
    }
}

// From http://dslweb.nwnexus.com/~ast/dload/guicon.htm
static const WORD MAX_CONSOLE_LINES = 5000;

// From http://dslweb.nwnexus.com/~ast/dload/guicon.htm
// Is there a more recent and easier way to do this?
void RedirectIOToConsole()
{
    int hConHandle;
    long lStdHandle;
    CONSOLE_SCREEN_BUFFER_INFO coninfo;
    FILE *fp;

    // allocate a console for this app
    AllocConsole();

    // set the screen buffer to be big enough to let us scroll text
    GetConsoleScreenBufferInfo(GetStdHandle(STD_OUTPUT_HANDLE), &coninfo);
    coninfo.dwSize.Y = MAX_CONSOLE_LINES;
    SetConsoleScreenBufferSize(GetStdHandle(STD_OUTPUT_HANDLE), coninfo.dwSize);

    // redirect unbuffered STDOUT to the console
    lStdHandle = (intptr_t)GetStdHandle(STD_OUTPUT_HANDLE);
    hConHandle = _open_osfhandle(lStdHandle, _O_TEXT);
    fp = _fdopen( hConHandle, "w" );
    *stdout = *fp;
    setvbuf( stdout, NULL, _IONBF, 0 );

    // redirect unbuffered STDIN to the console
    lStdHandle = (intptr_t)GetStdHandle(STD_INPUT_HANDLE);
    hConHandle = _open_osfhandle(lStdHandle, _O_TEXT);
    fp = _fdopen( hConHandle, "r" );
    *stdin = *fp;
    setvbuf( stdin, NULL, _IONBF, 0 );

    // redirect unbuffered STDERR to the console
    lStdHandle = (intptr_t)GetStdHandle(STD_ERROR_HANDLE);
    hConHandle = _open_osfhandle(lStdHandle, _O_TEXT);
    fp = _fdopen( hConHandle, "w" );
    *stderr = *fp;
    setvbuf( stderr, NULL, _IONBF, 0 );
}
