// Copyright © 2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// The microsoft (windows) native layer directx implementation.
// FUTURE This wraps the microsoft API's to create a directx context.
//        Wait for DirectX12 and Vulkan to settle down. 

// +build windows,dx
// FUTURE Use directx only when specifically asked for on windows.

#include "os_windows.h"
#include <d3d11_1.h>     // TDM-GCC-64/x86_64-w64-mingw32/include/d3d11_1.h
#include <dxgi1_2.h>     // TDM-GCC-64/x86_64-w64-mingw32/include/dxgi1_2.h

ID3D11Device *device;
ID3D11DeviceContext *context;
IDXGISwapChain1 *swapchain;
ID3D11RenderTargetView *renderTarget;
ID3D11Texture2D *backbuffer;

long gs_context(long long * display, long long * shell)
{
    // Create the device and device context objects
    // ID3D11Device *device;
    // ID3D11DeviceContext *context;
    D3D11CreateDevice(
        nullptr,                  // let directx choose the adapter.
        D3D_DRIVER_TYPE_HARDWARE, // hardware only, so...
        nullptr,                  // ... ignore software.
        0,
        nullptr,                  // ignore feature levels.
        0,                        // ... so feature level count is 0.
        D3D11_SDK_VERSION,        // compatibility Dx11 and up.
        &device,                  // device filled in.
        nullptr,                  // ignore feature levels.
        &context);                // device context filled in.
    IDXGIDevice2* dxgiDevice;
    device->QueryInterface(__uuidof(IDXGIDevice2), (void **)&dxgiDevice);
    dxgiDevice->SetMaximumFrameLatency(1);
    IDXGIAdapter* dxgiAdapter;
    dxgiDevice->GetAdapter(&dxgiAdapter);
    IDXGIFactory2* dxgiFactory;
    dxgiAdapter->GetParent(IID_PPV_ARGS(&dxgiFactory));

    // set up the swap chain description
    DXGI_SWAP_CHAIN_DESC1 scd = {0};
    scd.BufferUsage = DXGI_USAGE_RENDER_TARGET_OUTPUT; // how the swap chain should be used
    scd.BufferCount = 2;                               // a front buffer and a back buffer
    scd.Format = DXGI_FORMAT_B8G8R8A8_UNORM;           // the most common swap chain format
    scd.SwapEffect = DXGI_SWAP_EFFECT_FLIP_SEQUENTIAL; // the recommended flip mode
    scd.SampleDesc.Count = 1;                          // disable anti-aliasing

    // create the swap chain.
    HWND hwnd = HWND(LongToHandle(*display)); // window.
    dxgiFactory->CreateSwapChainForHwnd(
       device,                   // address of the device
       hwnd,                     // address of the window
       &scd,                     // address of the swap chain description
       nullptr,                  // ignore full screen mode.
       nullptr,                  // advanced
       &swapchain);              // address of the new swap chain pointer
    dxgiFactory->Release();

	// create a render target pointing to the back buffer
	swapchain->GetBuffer(0, __uuidof(ID3D11Texture2D), (void**)(&backbuffer));
	device->CreateRenderTargetView(backbuffer, nullptr, &renderTarget);
    return 1;
}

void gs_swap_buffers(long shell)
{
   context->OMSetRenderTargets(1, &renderTarget, nullptr);
   float color[4] = {0.0f, 0.2f, 0.4f, 1.0f};
   context->ClearRenderTargetView(renderTarget, color);
   swapchain->Present(1, 0);
}

void gs_display_dispose(long display)
{
   HWND hwnd = HWND(LongToHandle(display));
   HDC shell = GetDC(hwnd);
   renderTarget->Release();
   swapchain->Release();
   context->Release();
   device->Release();
   ReleaseDC(hwnd, shell);
   DestroyWindow(hwnd);
}
