// Copyright © 2013-2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// The microsoft (windows) native layer implementation.
// This wraps the microsoft windowing API's (where the real work is done).

#include <stdio.h>
#include "os_windows.h"

// App Icon related. 101 must be used in the resource file. For example:
//     101 ICON "application.ico"
// See windres for compiling resource files. Compiled resource files are
// included in a golang build using the .syso file type.
#define IDI_APPICON 101

// The golang callbacks that would normally be defined in _cgo_export.h
// are manually reproduced here and also implemented in the native test file
// os_windows_test.c.
extern void prepRender();
extern void renderFrame();
extern void handleInput(long event, long data);

// Globals to track the windows window and context handles.
long  dev_win_alive = -1; // Global used to track window closure.
HWND  display;            // Window handle.
HDC   shell;              // Device context handle.
HGLRC context;            // Rendering context handle.

// Global to help toggle full screen.
static screenInfo dev_screen = {0, 0, 0, 0, {0, 0, 0, 0}};

// Windows callback procedure. Handle a few events often returning 0 to mark
// them as handled. This method is mostly microsoft magic as each event has
// its own behaviour and different return codes.
//
// Called as frequently as possible to process user input and window changes.
LRESULT CALLBACK gs_wnd_proc(HWND hwnd, UINT msg, WPARAM wParam, LPARAM lParam)
{
    switch( msg )
    {
        case WM_ACTIVATE:
        {
            if (LOWORD(wParam) != WA_INACTIVE)
            {
                handleInput(devFocusIn, 0);
            }
            else
            {
                handleInput(devFocusOut, 0);
            }
            return 0;
        }
        case WM_SYSCOMMAND:
        {
            if ( (wParam & 0xfff0)  == SC_KEYMENU )
            {
                return 0;
            }
            break;
        }
        case WM_CLOSE:
        {
            dev_win_alive = -2;
            PostQuitMessage( 0 );
            return 0;
        }
        case WM_KEYDOWN:
        case WM_KEYUP:
        case WM_SYSKEYDOWN: // mod keys can mask regular keys.
        case WM_SYSKEYUP:   // track key releases like keyup for V in ALT-V.
        {
            long key = wParam;
            if (msg == WM_SYSKEYUP || msg == WM_KEYUP) {
                handleInput(devUp, key);
            }
            if (msg == WM_SYSKEYDOWN || msg == WM_KEYDOWN) {
                handleInput(devDown, key);
            }

            // send SYSKEY events to DefWindowProc so system stuff like
            // tabbing between windows still works.
            if (msg == WM_SYSKEYDOWN || msg == WM_SYSKEYUP) {
                return DefWindowProc(hwnd, msg, wParam, lParam);
            }
            return 0;
        }
        case WM_LBUTTONDOWN:
        {
            SetCapture(display);
            handleInput(devDown, devMouseL);
            return 0;
        }
        case WM_LBUTTONUP:
        {
            handleInput(devUp, devMouseL);
            ReleaseCapture();
            return 0;
        }
        case WM_MBUTTONDOWN:
        {
            SetCapture(display);
            handleInput(devDown, devMouseM);
            return 0;
        }
        case WM_MBUTTONUP:
        {
            handleInput(devUp, devMouseM);
            ReleaseCapture();
            return 0;
        }
        case WM_RBUTTONDOWN:
        {
            SetCapture(display);
            handleInput(devDown, devMouseR);
            return 0;
        }
        case WM_RBUTTONUP:
        {
            handleInput(devUp, devMouseR);
            ReleaseCapture();
            return 0;
        }
        case WM_MOUSEWHEEL:
        {
            // flip scroll direction to match OSX.
            long scroll = -1 * (((int)wParam) >> 16) / WHEEL_DELTA;
            handleInput(devScroll, scroll);
            return 0;
        }
        case WM_SIZE:
        {
            if (wParam == SIZE_MAXIMIZED || wParam == SIZE_RESTORED)
            {
                handleInput(devResize, 0);
            }
            return 0;
        }
        case WM_EXITSIZEMOVE:
        {
            handleInput(devResize, 0);
            return 0;
        }
    }

    // Pass all unhandled messages to DefWindowProc
    return DefWindowProc( hwnd, msg, wParam, lParam );
}

// Create the window. This is called twice on startup because
// windows are needed to get the initial and final rendering contexts.
HWND gs_create_window(HMODULE hInstance, LPSTR className)
{
    DWORD style = WS_TILEDWINDOW | WS_CLIPCHILDREN | WS_CLIPSIBLINGS;
    DWORD exStyle = WS_EX_APPWINDOW;

    // calculate the real window size from the desired size.
    RECT desktop;
    long xDefault = 600;
    long yDefault = 400;
    GetWindowRect(GetDesktopWindow(), &desktop);
    RECT rect = {0, 0, xDefault-1, yDefault-1};
    AdjustWindowRectEx( &rect, style, FALSE, exStyle );
    long wWidth = rect.right - rect.left + 1;
    long wHeight = rect.bottom - rect.top + 1;
    long yTop = desktop.bottom - yDefault - wHeight;

    // create the window
    HWND display = CreateWindowEx(
        exStyle,                // Optional styles
        className,              // Window class
        "WinTest",              // Window title
        style,                  // Window style
        600,                    // Size and position.
        yTop,                   // Size and position.
        wWidth,                 // Size and position.
        wHeight,                // Size and position.
        NULL,                   // Parent window
        NULL,                   // Menu
        hInstance,              // Module instance handle.
        NULL                    // Additional app data.
    );
    return display;
}

// =============================================================================
// What follows is code needed to get an OpenGL context.

// Bind an opengl extension to control the swap interval (and its dependency).
//     http://www.opengl.org/registry/specs/EXT/wgl_swap_control.txt
//     http://www.opengl.org/registry/specs/ARB/wgl_extensions_string.txt
typedef int (APIENTRY * PFNWGLSWAPINTERVALEXTPROC) (int);
typedef const char *(APIENTRY * PFNWGLGETEXTENSIONSSTRINGARBPROC)( HDC );

// Bind an opengl extension in order to create a proper context on windows.
// (and its dependency) See: wglCreateContextAttribsARB
//     http://www.opengl.org/wiki/Creating_an_OpenGL_Context_(WGL)
//     http://www.opengl.org/registry/specs/ARB/wgl_create_context.txt
//     http://www.opengl.org/registry/specs/EXT/wgl_extensions_string.txt
typedef HGLRC (WINAPI * PFNWGLCREATECONTEXTATTRIBSARBPROC) (HDC, HGLRC, const int *);
typedef const char *(APIENTRY * PFNWGLGETEXTENSIONSSTRINGEXTPROC)( void );
#define WGL_CONTEXT_MAJOR_VERSION_ARB             0x2091
#define WGL_CONTEXT_MINOR_VERSION_ARB             0x2092
#define WGL_CONTEXT_LAYER_PLANE_ARB               0x2093
#define WGL_CONTEXT_FLAGS_ARB                     0x2094
#define WGL_CONTEXT_PROFILE_MASK_ARB              0x9126
#define WGL_CONTEXT_DEBUG_BIT_ARB                 0x0001
#define WGL_CONTEXT_FORWARD_COMPATIBLE_BIT_ARB    0x0002
#define WGL_CONTEXT_CORE_PROFILE_BIT_ARB          0x00000001
#define WGL_CONTEXT_COMPATIBILITY_PROFILE_BIT_ARB 0x00000002

// Bind an opengl extension to read and choose a pixel format using attributes.
//     http://www.opengl.org/registry/specs/ARB/wgl_pixel_format.txt
typedef BOOL (WINAPI * PFNWGLGETPIXELFORMATATTRIBIVARBPROC) (HDC, int, int, UINT, const int *, int *);
typedef BOOL (WINAPI * PFNWGLCHOOSEPIXELFORMATARBPROC) (HDC, const int *, const FLOAT *, UINT, int *, UINT *);
#define WGL_NUMBER_PIXEL_FORMATS_ARB    0x2000
#define WGL_DRAW_TO_WINDOW_ARB          0x2001
#define WGL_SUPPORT_OPENGL_ARB          0x2010
#define WGL_ACCELERATION_ARB            0x2003
#define WGL_DOUBLE_BUFFER_ARB           0x2011
#define WGL_STEREO_ARB                  0x2012
#define WGL_PIXEL_TYPE_ARB              0x2013
#define WGL_COLOR_BITS_ARB              0x2014
#define WGL_RED_BITS_ARB                0x2015
#define WGL_GREEN_BITS_ARB              0x2017
#define WGL_BLUE_BITS_ARB               0x2019
#define WGL_ALPHA_BITS_ARB              0x201B
#define WGL_ACCUM_BITS_ARB              0x201D
#define WGL_ACCUM_RED_BITS_ARB          0x201E
#define WGL_ACCUM_GREEN_BITS_ARB        0x201F
#define WGL_ACCUM_BLUE_BITS_ARB         0x2020
#define WGL_ACCUM_ALPHA_BITS_ARB        0x2021
#define WGL_DEPTH_BITS_ARB              0x2022
#define WGL_STENCIL_BITS_ARB            0x2023
#define WGL_AUX_BUFFERS_ARB             0x2024
#define WGL_SAMPLE_BUFFERS_ARB          0x2041
#define WGL_SAMPLES_ARB                 0x2042
#define WGL_NO_ACCELERATION_ARB         0x2025
#define WGL_GENERIC_ACCELERATION_ARB    0x2026
#define WGL_FULL_ACCELERATION_ARB       0x2027
#define WGL_TYPE_RGBA_ARB               0x202B
#define WGL_TYPE_COLORINDEX_ARB         0x202C

// The initial pixel format is used to get an initial rendering context so
// that more rendering functions can be loaded. These new rendering functions
// allow a better, final, rendering context to be created.
int gs_get_initial_pixelformat(HDC hdc)
{
    // HDC hdc = LongToHandle(shell);
    int flags = PFD_DRAW_TO_WINDOW | PFD_SUPPORT_OPENGL | PFD_DOUBLEBUFFER;
    PIXELFORMATDESCRIPTOR pfd =
    {
        sizeof(PIXELFORMATDESCRIPTOR),
        1,                 // version
        flags,             // see above
        PFD_TYPE_RGBA,     // type of framebuffer
        32,                // color depth
        0,0,0,0,0,0,0,0,   // red, green, blue, alpha bits
        0,0,0,0,0,         // accum bits
        24,                // depth buffer bits.
        0,                 // stencil buffer bits.
        PFD_MAIN_PLANE,    // layer
        0,0,0,0,0          // unused
    };

    // create the temporary context using the proper pixel format.
    int pixelFormat = ChoosePixelFormat(hdc, &pfd);
    if (pixelFormat != 0)
    {
        SetPixelFormat(hdc, pixelFormat, &pfd);
        return pixelFormat;
    }
    return 0;
}

// OpenGL extensions that are bound after the first context is created.
PFNWGLGETEXTENSIONSSTRINGEXTPROC gs_wglGetExtensionsStringEXT = NULL;
PFNWGLSWAPINTERVALEXTPROC gs_wglSwapIntervalARB =  NULL;
PFNWGLGETEXTENSIONSSTRINGARBPROC gs_wglGetExtensionsStringARB =  NULL;
PFNWGLCREATECONTEXTATTRIBSARBPROC gs_wglCreateContextAttribsARB =  NULL;
PFNWGLGETPIXELFORMATATTRIBIVARBPROC gs_wglGetPixelFormatAttribivARB =  NULL;
PFNWGLCHOOSEPIXELFORMATARBPROC gs_wglChoosePixelFormatARB =  NULL;

// The final pixel format created using the bound rendering functions.
int gs_get_pixelformat(HDC hdc)
{
    const int attribList[] =
    {
        WGL_DRAW_TO_WINDOW_ARB, 1,
        WGL_SUPPORT_OPENGL_ARB, 1,
        WGL_DOUBLE_BUFFER_ARB, 1,
        WGL_PIXEL_TYPE_ARB, WGL_TYPE_RGBA_ARB,
        WGL_ACCELERATION_ARB, WGL_FULL_ACCELERATION_ARB,
        WGL_COLOR_BITS_ARB, 32,
        WGL_DEPTH_BITS_ARB, 24,
        WGL_STENCIL_BITS_ARB, 8,
        0, // end
    };
    int pixelFormat;
    UINT numFormats;
    gs_wglChoosePixelFormatARB(hdc, attribList, NULL, 1, &pixelFormat, &numFormats);
    if (pixelFormat != 0)
    {
        PIXELFORMATDESCRIPTOR pfd;
        DescribePixelFormat(hdc, pixelFormat, sizeof(PIXELFORMATDESCRIPTOR), &pfd);
        SetPixelFormat(hdc, pixelFormat, &pfd);
        return pixelFormat;
    }
    return 0;
}

// Destroy the application window. Attempt to remove the rendering context
// and the device context as well.
void gs_display_dispose(HWND hwnd) // (long display)
{
    // HWND hwnd = LongToHandle(display);
    HDC shell = GetDC(hwnd);
    HGLRC context = wglGetCurrentContext();
    wglMakeCurrent(NULL, NULL);
    wglDeleteContext(context);
    ReleaseDC(hwnd, shell);
    DestroyWindow(hwnd);
}

// gs_context creates an opengl context. Actually it creates two of them.
// The first context is used to find better functions to create the final
// context. Note that the pixel format is done only once for a window which
// is why the window is destroyed and recreated for the second context.
HGLRC gs_context(HWND *display, HDC *shell)
{
    // Get the application instance.
    HMODULE hInstance = GetModuleHandle(NULL);
    LPSTR gs_className = TEXT("GS_WIN");

    // Register the window class - once.
    WNDCLASSEX wc;
    wc.cbSize        = sizeof(WNDCLASSEX);
    wc.style         = CS_HREDRAW | CS_VREDRAW | CS_OWNDC;
    wc.lpfnWndProc   = (WNDPROC) gs_wnd_proc;
    wc.cbClsExtra    = 0;
    wc.cbWndExtra    = 0;
    wc.hInstance     = hInstance;
    wc.hIcon         = (HICON) LoadImage(hInstance, MAKEINTRESOURCE(IDI_APPICON), IMAGE_ICON, 0, 0, LR_SHARED);
    wc.hCursor       = LoadCursor(NULL, IDC_ARROW);
    wc.hbrBackground = NULL;
    wc.lpszMenuName  = NULL;
    wc.lpszClassName = gs_className;
    wc.hIconSm       = (HICON) LoadImage(hInstance, MAKEINTRESOURCE(IDI_APPICON), IMAGE_ICON, 0, 0, LR_SHARED);
    if(!RegisterClassEx(&wc))
    {
        printf("Failed RegisterClassEx %ld\n", GetLastError());
        return 0;
    }

    // Create the initial window and context in order to get a better
    // window and context.
    *display = gs_create_window(hInstance, gs_className);
    *shell = GetDC(*display);
    if (shell == NULL)
    {
        printf("Failed GetDC %ld\n", GetLastError());
        return 0;
    }

    // create the initial context.
    HGLRC initialContext;
    int initial_pixelFormat = gs_get_initial_pixelformat(*shell);
    if (initial_pixelFormat != 0)
    {
        initialContext = wglCreateContext(*shell);
        if (initialContext != NULL)
        {
            if (!wglMakeCurrent(*shell, initialContext))
            {
                wglDeleteContext(initialContext);
                initialContext = NULL;
            }
        }
    }
    if (initialContext == NULL)
    {
        printf("Failed initial context %ld\n", GetLastError());
        return 0; // failed to get even a simple context.
    }

    // now that there is a context, bind the opengl extensions and fail
    // if the supported extensions are too old or if they are not there.
    gs_wglGetExtensionsStringEXT    = (PFNWGLGETEXTENSIONSSTRINGEXTPROC)
                                       wglGetProcAddress( "wglGetExtensionsStringEXT" );
    gs_wglSwapIntervalARB           = (PFNWGLSWAPINTERVALEXTPROC)
                                       wglGetProcAddress( "wglSwapIntervalEXT" );
    gs_wglGetExtensionsStringARB    = (PFNWGLGETEXTENSIONSSTRINGARBPROC)
                                       wglGetProcAddress( "wglGetExtensionsStringARB" );
    gs_wglCreateContextAttribsARB   = (PFNWGLCREATECONTEXTATTRIBSARBPROC)
                                       wglGetProcAddress( "wglCreateContextAttribsARB" );
    gs_wglGetPixelFormatAttribivARB = (PFNWGLGETPIXELFORMATATTRIBIVARBPROC)
                                       wglGetProcAddress( "wglGetPixelFormatAttribivARB" );
    gs_wglChoosePixelFormatARB      = (PFNWGLCHOOSEPIXELFORMATARBPROC)
                                       wglGetProcAddress( "wglChoosePixelFormatARB" );
    if (gs_wglGetExtensionsStringEXT == NULL ||
        gs_wglSwapIntervalARB ==  NULL ||
        gs_wglGetExtensionsStringARB ==  NULL ||
        gs_wglCreateContextAttribsARB ==  NULL ||
        gs_wglGetPixelFormatAttribivARB ==  NULL ||
        gs_wglChoosePixelFormatARB ==  NULL)
    {
        printf("Failed binding render extentions %ld\n", GetLastError());
        return 0;
    }

    // destroy and recreate the initial window and context.
    gs_className = TEXT("GS_WIN");
    gs_display_dispose(*display);
    *display = gs_create_window(hInstance, gs_className);
    *shell = GetDC(*display);
    int pixelformat = gs_get_pixelformat(*shell);
    if (pixelformat == 0)
    {
        printf("Failed second pixel format %ld\n", GetLastError());
        return 0;
    }

    // now create the context on the fresh window.
    int cnt = 0;
    int attribs[40];

    // Use the expected baseline opengl 3.2
    attribs[cnt++] = WGL_CONTEXT_MAJOR_VERSION_ARB;
    attribs[cnt++] = 3;
    attribs[cnt++] = WGL_CONTEXT_MINOR_VERSION_ARB;
    attribs[cnt++] = 2;
    attribs[cnt++] = WGL_CONTEXT_FLAGS_ARB;
    attribs[cnt++] = WGL_CONTEXT_FORWARD_COMPATIBLE_BIT_ARB;
    attribs[cnt++] = WGL_CONTEXT_PROFILE_MASK_ARB;
    attribs[cnt++] = WGL_CONTEXT_CORE_PROFILE_BIT_ARB;
    attribs[cnt++] = 0;
    HGLRC context = gs_wglCreateContextAttribsARB( *shell, NULL, attribs );
    if (context == NULL)
    {
        printf("Failed second context %ld\n", GetLastError());
        return 0; // failed to get rendering context
    }
    wglMakeCurrent(*shell, context);
    ShowWindow(*display, SW_SHOW);
    SetForegroundWindow(*display);
    return context;
}

// The above code is responsible for getting a render context.
// =============================================================================
// Native layer wrappers.

// Process input and render frames. This is a simple game loop that expects
// more complex stuff like fixed time-step to be handled by renderFrame.
// Ensures user input is processed by routing events through gs_wnd_proc.
//
// It mimics macOS and iOS where the OS keeps the loop and calls the application
// to render based on the display refresh rate. Here rendering is called as fast
// as possible - which can be inefficient since frames are still only displayed
// as fast as the monitor refresh rate.
void dev_run()
{
    context = gs_context(&display, &shell);
    dev_win_alive = 1;
    prepRender();
    while (dev_win_alive == 1)
    {
        MSG msg;
        if ( PeekMessage( &msg, NULL, 0, 0, PM_REMOVE ) != 0 )
        {
            if (msg.message == WM_QUIT)
            {
                // handle quit immediately
                dev_win_alive = -2;
                return;
            }
            DispatchMessage( &msg ); // goes to gs_wnd_proc for processing.
        }
        renderFrame();
    }
    wglDeleteContext(context);
    gs_display_dispose(display);
}

// Swaps rendering buffer. Called after rendering a frame.
void dev_swap()
{
    SwapBuffers(shell);
}

// Cleans and releases all resources including the OpenGL context.
void dev_dispose() {
    dev_win_alive = -2;
    wglDeleteContext(context);
    gs_display_dispose(display);
}

// Used to check if the application is full screen mode.
// Return 1 if the application is full screen, 0 otherwise.
unsigned char dev_fullscreen()
{
    return dev_screen.full;
}

// Flip full screen mode. Expected to be called after starting
// processing of events. Based on:
// http://src.chromium.org/viewvc/chrome/trunk/src/ui/views/win/
//        fullscreen_handler.cc?revision=HEAD&view=markup
void dev_toggle_fullscreen()
{
    if (!dev_screen.full)
    {
        dev_screen.maxed = IsZoomed(display);
        if (dev_screen.maxed)
        {
            SendMessage(display, WM_SYSCOMMAND, SC_RESTORE, 0);
        }
        dev_screen.style = GetWindowLong(display, GWL_STYLE);
        dev_screen.ex_style = GetWindowLong(display, GWL_EXSTYLE);
        GetWindowRect(display, &dev_screen.rect);
    }
    dev_screen.full = !dev_screen.full;
    if (dev_screen.full)
    {
        SetWindowLong(display, GWL_STYLE,
                   dev_screen.style & ~(WS_CAPTION | WS_THICKFRAME));
        SetWindowLong(display, GWL_EXSTYLE,
                   dev_screen.ex_style & ~(WS_EX_DLGMODALFRAME |
                   WS_EX_WINDOWEDGE | WS_EX_CLIENTEDGE | WS_EX_STATICEDGE));
        MONITORINFO m_info;
        m_info.cbSize = sizeof(m_info);
        GetMonitorInfo(MonitorFromWindow(display, MONITOR_DEFAULTTONEAREST), &m_info);
        RECT m_rect = m_info.rcMonitor;
        SetWindowPos(display, NULL, m_rect.left, m_rect.top,
                     m_rect.right-m_rect.left, m_rect.bottom-m_rect.top,
                     SWP_NOZORDER | SWP_NOACTIVATE | SWP_FRAMECHANGED);
    }
    else
    {
        SetWindowLong(display, GWL_STYLE, dev_screen.style);
        SetWindowLong(display, GWL_EXSTYLE, dev_screen.ex_style);
        RECT m_rect = dev_screen.rect;
        SetWindowPos(display, NULL, m_rect.left, m_rect.top,
                     m_rect.right-m_rect.left, m_rect.bottom-m_rect.top,
                     SWP_NOZORDER | SWP_NOACTIVATE | SWP_FRAMECHANGED);
        if (dev_screen.maxed)
        {
            SendMessage(display, WM_SYSCOMMAND, SC_MAXIMIZE, 0);
        }
    }
    PostMessage(display, WM_EXITSIZEMOVE, 0, 0); // Trigger window resize.
}

// Show or hide cursor. Lock it if it is hidden.
void dev_show_cursor(unsigned char show) {
    if (show)
    {
        ReleaseCapture();
    }
    else
    {
        SetCapture(display);
    }
    ShowCursor( show );
}

// Get the current mouse position relative to the bottom left corner
// of the application window.
void dev_cursor(long *x, long *y)
{
    POINT point;
    GetCursorPos(&point);
    ScreenToClient(display, &point);
    RECT rect;
    GetClientRect(display, &rect);
    *x = point.x;
    *y = rect.bottom - point.y;
}

// Position the cursor at the given window location. The incoming coordinates
// are relative to the bottom left corner - switch that to be relative to the
// top left corner expected by windows.
void dev_set_cursor_location(long x, long y)
{
    RECT rect;
    if (GetClientRect(display, &rect) != 0 )
    {
        POINT loc;
        loc.x = x;
        loc.y = rect.bottom - y;
        if (ClientToScreen(display, &loc) != 0 )
        {
            SetCursorPos(loc.x, loc.y);
        }
    }
}

// Sets the windows size and location.
// The y value is reversed because the incoming coordinates are relative
// to the bottom left corner. Windows expects it to be the top left.
void dev_set_size(long x, long y, long w, long h)
{
    RECT desk;
    if (GetWindowRect(GetDesktopWindow(), &desk) != 0 )
    {
        RECT wind;
        if (GetWindowRect(display, &wind) != 0 ) {
            RECT disp;
            if (GetClientRect(display, &disp) != 0 )
            {
                int xExtra = wind.right - wind.left - disp.right;
                int yExtra = wind.bottom - wind.top - disp.bottom;
                y = desk.bottom - y - h;
                SetWindowPos(display, HWND_TOP, x, y, w+xExtra, h+yExtra, 0);
            }
        }
    }
}

// Get the current main window drawing area size.
// Reverse y so origin is bottom left.
void dev_size(long *x, long *y, long *w, long *h)
{
    RECT rect;
    GetClientRect(display, &rect);
    *w = rect.right - rect.left;
    *h = rect.bottom - rect.top;
    RECT desktop, window;
    GetWindowRect(GetDesktopWindow(), &desktop);
    GetWindowRect(display, &window);
    *x = window.left;
    int yExtra = window.bottom - window.top - rect.bottom;
    *y = desktop.bottom - window.bottom + yExtra;
}

// Sets the windows title.
void dev_set_title(char * label)
{
    SetWindowText(display, label);
}

// Returns a WCHAR string of the specified UTF-8 string
// The returned string must be freed by the caller.
// Needed by copy/paste to handle UTF8 strings.
WCHAR* utf8_wchar(const char* utf8) {
    int length = MultiByteToWideChar(CP_UTF8, 0, utf8, -1, NULL, 0);
    if (!length) {
        return NULL;
    }
    WCHAR* wide = calloc(length, sizeof(WCHAR));
    if (!MultiByteToWideChar(CP_UTF8, 0, utf8, -1, wide, length)) {
        free(wide);
        return NULL;
    }
    return wide; // needs to be freed by the caller.
}

// Returns a UTF-8 string version of the specified wide string
// The returned string must be freed by the caller.
// Needed by copy/paste to handle UTF8 strings.
char* wchar_utf8(const WCHAR* wide) {
    int length = WideCharToMultiByte(CP_UTF8, 0, wide, -1, NULL, 0, NULL, NULL);
    if (!length) {
        return NULL;
    }
    char* utf8 = calloc(length, sizeof(char));
    if (!WideCharToMultiByte(CP_UTF8, 0, wide, -1, utf8, length, NULL, NULL)) {
        free(utf8);
        return NULL;
    }
    return utf8; // needs to be freed by the caller.
}

// Return the current clipboard contents if the clipboard contains text.
// Otherwise return nil. Any returned strings must be freed by the caller.
char* dev_clip_copy()
{
    if (!IsClipboardFormatAvailable(CF_UNICODETEXT)) {
        return NULL;
    }
    if (!OpenClipboard(display)) {
        return NULL;
    }
    HANDLE stringHandle = GetClipboardData(CF_UNICODETEXT);
    if (!stringHandle) {
        CloseClipboard();
        return NULL;
    }
    char* clipboardString = wchar_utf8(GlobalLock(stringHandle));
    GlobalUnlock(stringHandle);
    CloseClipboard();
    if (!clipboardString) {
        return NULL;
    }
    return clipboardString;
}

// Paste the given string into the general clipboard.
void dev_clip_paste(const char* string)
{
    WCHAR* widestr = utf8_wchar(string);
    if (!widestr) {
        return;
    }
    size_t wstrlen = (wcslen(widestr) + 1) * sizeof(WCHAR);
    HANDLE stringHandle = GlobalAlloc(GMEM_MOVEABLE, wstrlen);
    if (!stringHandle) {
        free(widestr);
        return;
    }
    memcpy(GlobalLock(stringHandle), widestr, wstrlen);
    GlobalUnlock(stringHandle);
    if (!OpenClipboard(display)) {
        GlobalFree(stringHandle);
        free(widestr);
        return;
    }
    EmptyClipboard();
    SetClipboardData(CF_UNICODETEXT, stringHandle);
    CloseClipboard();
    free(widestr);
}
