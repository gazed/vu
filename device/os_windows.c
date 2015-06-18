// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// The microsoft (windows) native layer implementation.
// This wraps the microsoft API's (where the real work is done).

#include "os_windows.h"

// Application defaults. Internal use only. Not really state per-se these are
// consulted at startup for initial values. These are updated using the
// gs_set_attr* functions.
struct AppDefaults {
    long gs_ShellX;
    long gs_ShellY;
    long gs_ShellWidth;
    long gs_ShellHeight;
    long gs_AlphaSize;
    long gs_DepthSize;
    char gs_AppName[40];
};
struct AppDefaults defaults = { 100, 100, 240, 280, 8, 24, TEXT("App") };

// 101 must be used in the resource file. For example:
//     101 ICON "application.ico"
// See windres for compiling resource files. Compiled resource files are
// included in a golang build using the .syso file type.
#define IDI_APPICON 101

// State used to track window closure. Needed to avoid accessing the
// external shell pointer after a window has closed. There is no sure way to
// check if an object pointer is valid once that object has been released.
long gs_win_alive = -1;

// Fifo queue of event (urge) structure used to pass back events of interest.
// Needed because read_dispatch only handles a single event, but one user action
// can produce multiple events. Only ever seen 2 events produced from one.
static GSEvent gs_events[] = {
    {0, -1, -1, 0, 0, 0},
    {0, -1, -1, 0, 0, 0},
    {0, -1, -1, 0, 0, 0},
    {0, -1, -1, 0, 0, 0},
    {0, -1, -1, 0, 0, 0},
};
static int gs_event_front = 0;
static int gs_event_rear = 0;
static int gs_event_size = sizeof(gs_events) / sizeof(gs_events[0]);

// Full screen toggle structure.
static GSScreen gs_screen = {0, 0, 0, 0, {0, 0, 0, 0}};

void gs_write_urge(long eid, long key, long scroll)
{
    GSEvent *eve = &(gs_events[gs_event_rear]);
    eve->event = eid;
    eve->key = key;
    eve->scroll = scroll;
    eve->mousex = -1;
    eve->mousey = -1;
    eve->mods = 0;
    gs_event_rear = (gs_event_rear + 1) % gs_event_size;
}

// Windows callback procedure. Handle a few events often returning 0 to mark
// them as handled. This method is mostly microsoft magic as each event may
// have its own behaviour and different return codes.
LRESULT CALLBACK gs_wnd_proc(HWND hwnd, UINT msg, WPARAM wParam, LPARAM lParam)
{
    switch( msg )
    {
        case WM_ACTIVATE:
        {
            long isActive = LOWORD(wParam) != WA_INACTIVE ? 1 : 0;
            long event = msg + isActive; // GS_WindowActive or GS_WindowInactive
            gs_write_urge(event, 0, 0);
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
            gs_win_alive = -2;
            PostQuitMessage( 0 );
            return 0;
        }
        case WM_KEYDOWN:
        case WM_KEYUP:
        {
            long key = wParam;

            // only care about modifiers with other keys.
            if (key == VK_SHIFT || key == VK_CONTROL || key == VK_MENU || key == VK_LWIN || key == VK_RWIN) {
                return 0;
            }
            gs_write_urge(msg, key, 0);
            return 0;
        }
        case WM_MBUTTONDOWN:
        case WM_RBUTTONDOWN:
        case WM_LBUTTONDOWN:
        {
            SetCapture(hwnd);
            gs_write_urge(msg, 0, 0);
            return 0;
        }
        case WM_LBUTTONUP:
        case WM_RBUTTONUP:
        case WM_MBUTTONUP:
        {
            ReleaseCapture();
            gs_write_urge(msg, 0, 0);
            return 0;
        }
        case WM_MOUSEWHEEL:
        {
            // flip scroll direction to match OSX.
            long scroll = -1 * (((int)wParam) >> 16) / WHEEL_DELTA;
            gs_write_urge(msg, 0, scroll);
            return 0;
        }
        case WM_SIZE:
        {
            // TODO detect when window is restored from maximized.
            if (wParam == SIZE_MAXIMIZED)
            {
                gs_write_urge(GS_WindowResized, 0, 0);
            }
            return 0;
        }
        case WM_EXITSIZEMOVE:
        {
            gs_write_urge(msg, 0, 0); // sends GS_WindowResized
            return 0;
        }
    }

    // Pass all unhandled messages to DefWindowProc
    return DefWindowProc( hwnd, msg, wParam, lParam );
}

// Get the current mouse position relative to the bottom left corner of the
// application window.
void gs_pos(long display, long *x, long *y)
{
    HWND hwnd = LongToHandle(display);
    POINT point;
    GetCursorPos(&point);
    ScreenToClient(hwnd, &point);
    RECT rect;
    GetClientRect(hwnd, &rect);
    *x = point.x;
    *y = rect.bottom - point.y;
}

// Process all queued up user events and send one of the processed events
// back to the application. Prefer PeekMessage (non-blocking) over
// GetMessage (blocking).
void gs_read_dispatch(long display, GSEvent *gs_urge)
{
    MSG msg;
    if ( PeekMessage( &msg, NULL, 0, 0, PM_REMOVE ) != 0 )
    {
        // handle quit immediately
        if (msg.message == WM_QUIT)
        {
            gs_win_alive = -2;
            return;
        }
        DispatchMessage( &msg ); // goes to wnd_proc

        // message queue has been processed, return interesting stuff.
        if ( gs_event_front != gs_event_rear )
        {
			GSEvent *eve = &(gs_events[gs_event_front]);
            gs_urge->event = eve->event;
            gs_urge->key = eve->key;
            gs_urge->scroll = eve->scroll;
	        gs_event_front = (gs_event_front + 1) % gs_event_size;
        }
    }

    // always send back the modifier keys.
    long mods = 0;
    if ( GetKeyState(VK_SHIFT) & 0x8000 )
    {
        mods |= GS_ShiftKeyMask;
    }
    if ( GetKeyState(VK_CONTROL) & 0x8000 )
    {
        mods |= GS_ControlKeyMask;
    }
    if ( GetKeyState(VK_MENU) & 0x8000 )
    {
        mods |= GS_AlternateKeyMask;
    }
    if ( GetKeyState(VK_LWIN) & 0x8000 )
    {
        mods |= GS_CommandKeyMask;
    }
    if ( GetKeyState(VK_RWIN) & 0x8000 )
    {
        mods |= GS_CommandKeyMask;
    }
    gs_urge->mods = mods;

    // update the mouse each time rather than dealing with mouse move events.
    gs_pos(display, &(gs_urge->mousex), &(gs_urge->mousey));
}

// Needed because the window will be destroyed and recreated in order
// to create the proper opengl context.
LPSTR gs_className = TEXT("GS_WIN");

// Create the window, but don't open it.
long gs_create_window(HMODULE hInstance, LPSTR className)
{
    DWORD style = WS_TILEDWINDOW | WS_CLIPCHILDREN | WS_CLIPSIBLINGS;
    DWORD exStyle = WS_EX_APPWINDOW;

    // calculate the real window size from the desired size.
    RECT desktop;
    GetWindowRect(GetDesktopWindow(), &desktop);
    RECT rect = {0, 0, defaults.gs_ShellWidth-1, defaults.gs_ShellHeight-1};
    AdjustWindowRectEx( &rect, style, FALSE, exStyle );
    long wWidth = rect.right - rect.left + 1;
    long wHeight = rect.bottom - rect.top + 1;
    long topy = desktop.bottom - defaults.gs_ShellY - wHeight;

    // create the window
    HWND display = CreateWindowEx(
        exStyle,                // Optional styles
        className,              // Window class
        defaults.gs_AppName,    // Window title
        style,                  // Window style
        defaults.gs_ShellX,     // Size and position.
        topy,                   // Size and position.
        wWidth,                 // Size and position.
        wHeight,                // Size and position.
        NULL,                   // Parent window
        NULL,                   // Menu
        hInstance,              // Module instance handle.
        NULL                    // Additional app data.
    );
    return HandleToLong(display);
}

// Initialize, register the application class and create the initial
// application window.
long gs_display_init()
{
    // Get the application instance.
    HMODULE hInstance = GetModuleHandle(NULL);

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
        return 0;
    }
    return gs_create_window(hInstance, gs_className);
}

// Destroy the application window. Attempt to remove the rendering context and
// the device context as well.
void gs_display_dispose(long display)
{
    HWND hwnd = LongToHandle(display);
    HDC shell = GetDC(hwnd);
    HGLRC context = wglGetCurrentContext();
    wglMakeCurrent(NULL, NULL);
    wglDeleteContext(context);
    ReleaseDC(hwnd, shell);
    DestroyWindow(hwnd);
}

// Get the device context. This must be called after creating the window and
// before creating the rendering context.
long gs_shell(long display)
{
    HWND hwnd = LongToHandle(display);
    HDC shell = GetDC(hwnd);
    if (shell == NULL)
    {
        printf("Failed to GetDC %ld %ld\n", display, GetLastError());
    }
    return HandleToLong(shell);
}

// Show the application window to the user. This is expected to be called after
// the rendering context has been created.
void gs_shell_open(long display)
{
    HWND hwnd = LongToHandle(display);
    ShowWindow(hwnd, SW_SHOW);
    SetForegroundWindow(hwnd);
    gs_win_alive = 1;
}

// Check if the application window is still active.
unsigned char gs_shell_alive(long display)
{
    return gs_win_alive == 1;
}


// Return 1 if the application is full screen, 0 otherwise.
unsigned char gs_fullscreen(long display)
{
    return gs_screen.full;
}

// Flip full screen mode. Expected to be called after starting processing
// of events with gs_read_dispatch(). Based on:
// http://src.chromium.org/viewvc/chrome/trunk/src/ui/views/win/
//        fullscreen_handler.cc?revision=HEAD&view=markup
void gs_toggle_fullscreen(long display)
{
    HWND hwnd = LongToHandle(display);
    if (!gs_screen.full)
    {
        gs_screen.maxed = IsZoomed(hwnd);
        if (gs_screen.maxed)
        {
            SendMessage(hwnd, WM_SYSCOMMAND, SC_RESTORE, 0);
        }
        gs_screen.style = GetWindowLong(hwnd, GWL_STYLE);
        gs_screen.ex_style = GetWindowLong(hwnd, GWL_EXSTYLE);
        GetWindowRect(hwnd, &gs_screen.rect);
    }
    gs_screen.full = !gs_screen.full;
    if (gs_screen.full)
    {
        SetWindowLong(hwnd, GWL_STYLE,
                   gs_screen.style & ~(WS_CAPTION | WS_THICKFRAME));
        SetWindowLong(hwnd, GWL_EXSTYLE,
                   gs_screen.ex_style & ~(WS_EX_DLGMODALFRAME |
                   WS_EX_WINDOWEDGE | WS_EX_CLIENTEDGE | WS_EX_STATICEDGE));
        MONITORINFO m_info;
        m_info.cbSize = sizeof(m_info);
        GetMonitorInfo(MonitorFromWindow(hwnd, MONITOR_DEFAULTTONEAREST), &m_info);
        RECT m_rect = m_info.rcMonitor;
        SetWindowPos(hwnd, NULL, m_rect.left, m_rect.top,
                     m_rect.right-m_rect.left, m_rect.bottom-m_rect.top,
                     SWP_NOZORDER | SWP_NOACTIVATE | SWP_FRAMECHANGED);
    }
    else
    {
        SetWindowLong(hwnd, GWL_STYLE, gs_screen.style);
        SetWindowLong(hwnd, GWL_EXSTYLE, gs_screen.ex_style);
        RECT m_rect = gs_screen.rect;
        SetWindowPos(hwnd, NULL, m_rect.left, m_rect.top,
                     m_rect.right-m_rect.left, m_rect.bottom-m_rect.top,
                     SWP_NOZORDER | SWP_NOACTIVATE | SWP_FRAMECHANGED);
        if (gs_screen.maxed)
        {
            SendMessage(hwnd, WM_SYSCOMMAND, SC_MAXIMIZE, 0);
        }
    }
    PostMessage(hwnd, WM_EXITSIZEMOVE, 0, 0); // Trigger window resize.
}

// Position the cursor at the given window location. The incoming coordinates
// are relative to the bottom left corner - switch that to be relative to the
// top left corner expected by windows.
void gs_set_cursor_location(long display, long x, long y)
{
    HWND hwnd = LongToHandle(display);
    RECT rect;
    if (GetClientRect(hwnd, &rect) != 0 )
    {
        POINT loc;
        loc.x = x;
        loc.y = rect.bottom - y;
        if (ClientToScreen(hwnd, &loc) != 0 )
        {
            SetCursorPos(loc.x, loc.y);
        }
    }
}

// Get the current application windows client area location and size.
void gs_size(long display, long *x, long *y, long *w, long *h)
{
    HWND hwnd = LongToHandle(display);
    RECT rect;
    GetClientRect(hwnd, &rect);
    *w = rect.right - rect.left;
    *h = rect.bottom - rect.top;
    RECT desktop;
    GetWindowRect(GetDesktopWindow(), &desktop);
    GetWindowRect(hwnd, &rect);
    *x = rect.left;
    *y = desktop.bottom - rect.bottom;
}

// Show or hide cursor. Lock it to the window if it is hidden.
void gs_show_cursor(long display, unsigned char show)
{
    if (show)
    {
        ReleaseCapture();
    }
    else
    {
        HWND hwnd = LongToHandle(display);
        SetCapture(hwnd);
    }
    ShowCursor( show );
}

// The initial pixel format is used to get an initial rendering context so
// that more rendering functions can be loaded. These new rendering functions
// allow a better, final, rendering context to be created.
int gs_get_initial_pixelformat(long shell)
{
    HDC hdc = LongToHandle(shell);
    int flags = PFD_DRAW_TO_WINDOW | PFD_SUPPORT_OPENGL | PFD_DOUBLEBUFFER;
    PIXELFORMATDESCRIPTOR pfd =
    {
        sizeof(PIXELFORMATDESCRIPTOR),
        1,                 // version
        flags,             // see above
        PFD_TYPE_RGBA,     // type of framebuffer
        32,                // colour depth
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
int gs_get_pixelformat(long shell)
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
    HDC hdc = LongToHandle(shell);
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

// gs_context creates an opengl context. Actually it creates two of them.
// The first context is used to find better functions to create the final
// context.  Note that the pixel format is done only once for a window so
// it must be correctly chosen.
long gs_context(long long * display, long long * shell)
{
    // create the initial context.
    HDC hdc = LongToHandle(*shell);
    HGLRC initialContext;
    int initial_pixelFormat = gs_get_initial_pixelformat(*shell);
    if (initial_pixelFormat != 0)
    {
        initialContext = wglCreateContext(hdc);
        if (initialContext != NULL)
        {
            if (!wglMakeCurrent(hdc, initialContext))
            {
                wglDeleteContext(initialContext);
                initialContext = NULL;
            }
        }
    }
    if (initialContext == NULL)
    {
        return 0; // failed to get even a simple context.
    }

    // now that there is a context, bind the opengl extensions and fail
    // if the supported extensions are too old or if they are not there.
    gs_wglGetExtensionsStringEXT = (PFNWGLGETEXTENSIONSSTRINGEXTPROC) wglGetProcAddress( "wglGetExtensionsStringEXT" );
    gs_wglSwapIntervalARB = (PFNWGLSWAPINTERVALEXTPROC) wglGetProcAddress( "wglSwapIntervalEXT" );
    gs_wglGetExtensionsStringARB = (PFNWGLGETEXTENSIONSSTRINGARBPROC) wglGetProcAddress( "wglGetExtensionsStringARB" );
    gs_wglCreateContextAttribsARB = (PFNWGLCREATECONTEXTATTRIBSARBPROC) wglGetProcAddress( "wglCreateContextAttribsARB" );
    gs_wglGetPixelFormatAttribivARB = (PFNWGLGETPIXELFORMATATTRIBIVARBPROC) wglGetProcAddress( "wglGetPixelFormatAttribivARB" );
    gs_wglChoosePixelFormatARB = (PFNWGLCHOOSEPIXELFORMATARBPROC) wglGetProcAddress( "wglChoosePixelFormatARB" );
    if (gs_wglGetExtensionsStringEXT == NULL ||
        gs_wglSwapIntervalARB ==  NULL ||
        gs_wglGetExtensionsStringARB ==  NULL ||
        gs_wglCreateContextAttribsARB ==  NULL ||
        gs_wglGetPixelFormatAttribivARB ==  NULL ||
        gs_wglChoosePixelFormatARB ==  NULL)
    {
        return 0;
    }

    // destroy and recreate the window and shell
    gs_display_dispose(*display);
    HMODULE hInstance = GetModuleHandle(NULL);
    *display = gs_create_window(hInstance, gs_className);
    *shell = gs_shell(*display);
    int pixelformat = gs_get_pixelformat(*shell);
    if (pixelformat == 0)
    {
        return 0;
    }

    // now create the context on the fresh window.
    int cnt = 0;
    int attribs[40];
    hdc = LongToHandle(*shell);

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
    HGLRC context = gs_wglCreateContextAttribsARB( hdc, NULL, attribs );
    if (context != NULL)
    {
        if (wglMakeCurrent(hdc, context))
        {
            return HandleToLong(context);
        }
    }
    return 0; // failed to get rendering context
}

// Flip the back and front buffers of the rendering context.
void gs_swap_buffers(long shell)
{
    HDC hdc = LongToHandle(shell);
    SwapBuffers(hdc);
}

// Set long attributes. Attributes only take effect if they are set before
// they are used to create the window or rendering context.
void gs_set_attr_l(long attr, long value)
{
   switch (attr) {
   case GS_ShellX:
       if (value > 0) { defaults.gs_ShellX = value; }
       break;
   case GS_ShellY:
       if (value > 0) { defaults.gs_ShellY = value; }
       break;
   case GS_ShellWidth:
       if (value > 0) { defaults.gs_ShellWidth = value; }
       break;
   case GS_ShellHeight:
       if (value > 0) { defaults.gs_ShellHeight = value; }
       break;
   case GS_AlphaSize:
       if (value >= 0) { defaults.gs_AlphaSize = value; }
       break;
   case GS_DepthSize:
       if (value >= 0) { defaults.gs_DepthSize = value; }
       break;
   }
}

// Set string attributes. Attributes only take effect if they are set before
// they are used to create the window or rendering context.
void gs_set_attr_s(long attr, char * value)
{
   switch (attr) {
   case GS_AppName:
       if (value != NULL && strlen(value) < 40) {
           strcpy( defaults.gs_AppName, TEXT(value) );
       }
       break;
   }
}
