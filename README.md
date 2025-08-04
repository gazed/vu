<!-- Copyright © 2013-2024 Galvanized Logic Inc. -->

# Vu

Vu (Virtual Universe) is a 3D engine based on Go (Golang) programming language.
Vu is composed of packages, detailed in [GoDoc](http://godoc.org/github.com/gazed/vu),
and briefly summarized below. More getting started and background information is available
on the [Wiki](https://github.com/gazed/vu/wiki)

Vu is a small engine intended for simple games. It currently supports Vulkan on Windows.

* `vu/physics` handles spheres and convex hulls.
* `vu/render` uses Vulkan 1.3 without any extensions.
* `vu/device` supports a basic window, button presses, mouse clicks, and mouse movement. 

Vu started as a learning project for Go and 3D programming.
It is currently being used by the author to develop games and other 3D applications.

Sub packages
--------

* [audio](http://godoc.org/github.com/gazed/vu/audio) Positions and plays sounds in a 3D environment.
* [device](http://godoc.org/github.com/gazed/vu/device)  Links the application to native OS specific window and user events.
* [eg](http://godoc.org/github.com/gazed/vu/eg) Examples that both demonstrate and test the vu engine.
* [load](http://godoc.org/github.com/gazed/vu/load) Asset loaders including models, textures, audio, shaders, and bitmapped fonts.
* [math/lin](http://godoc.org/github.com/gazed/vu/math/lin) Linear math library for vectors, matricies, and quaternions.
* [physics](http://godoc.org/github.com/gazed/vu/physics) Repositions bodies based on simulated physics.
* [render](http://godoc.org/github.com/gazed/vu/render) 3D drawing and graphics interface.

Build and Test
------

* Clone the repo, ie: `git clone git@github.com:gazed/vu.git`
* Use `go generate` in `vu/assets/shaders` to build the shaders. Requires `glslc` executable.
* `go build` in `vu\eg` and run the examples.

**Build Dependencies**

* Go version 1.23 or later.
* Vulkan version 1.3 or later, and vulkan validation layer from the SDK at https://www.lunarg.com/vulkan-sdk/
* OpenAL (https://openal.org) latest 64-bit version `OpenAL32.dll`+`soft_oal.dll` from https://openal-soft.org/openal-binaries/
* `glslc` executable from https://github.com/google/shaderc is needed to build shaders in `vu/assets/shaders`.

Example Game
------

[Floworlds](https://www.floworlds.com) is a strategic puzzle terraforming game coded entirely in Go using the `vu` engine.

Credit Where Credit is Due
--------

Vu is built on the generous sharing of code and information by many talented, kind,
and lucid developers. Thank you in particular to:

* [https://www.youtube.com/@TravisVroman](https://www.youtube.com/@TravisVroman)
  Travis Vroman brilliantly explains how to use Vulkan in a game engine.
  His video devlog and code base provides an overall context for Vulkan that is
  difficult to get from the specification. Want a Vulkan C engine? Checkout 
  [https://github.com/travisvroman/kohi](https://github.com/travisvroman/kohi) Apache License.

* [https://github.com/bbredesen/go-vk](https://github.com/bbredesen/go-vk) MIT License.
  Go bindings for Vulkan. This is a huge amount of work due to the size and complexity of the Vulkan spec.
  Bindings were generated from a hacked version, `gazed/vk-gen`, that uses Syscall instead of Cgo.
  The generated bindings were dropped into `vu/internal/render/vk`

* [https://github.com/felipeek/raw-physics](https://github.com/felipeek/raw-physics) MIT License.
  Raw-physics provides the ability to collide spheres and convex hulls in less than 4000 lines of code.
  This was perfect for porting from C into Go `vu/physics` pretty much line for line.

* [https://github.com/lxn/win](https://github.com/lxn/win)
  Go bindings for the Windows API that do not require c-go and a C compiler.
  The parts needed were put directly into `vu/internal/device/win` along with the original License.

* [https://github.com/qmuntal/gltf](https://github.com/qmuntal/gltf) BSD 2-Clause License.
  GLTF helps tremendously for importing 3D model data and thanks to qmuntal it was possible
  to pull a few thousand lines of Go code into `vu/internal/load/gltf` and get a working importer.
