<!-- Copyright © 2013-2017 Galvanized Logic Inc.                       -->
<!-- Use is governed by a BSD-style license found in the LICENSE file. -->

# Vu

Vu (Virtual Universe) is a 3D engine based on the modern programming language Go (Golang).
Vu is composed of packages, detailed in [GoDoc](http://godoc.org/github.com/gazed/vu),
and briefly summarized below.

Sub packages
--------

* [audio](http://godoc.org/github.com/gazed/vu/audio) Positions and plays sounds in a 3D environment.
* [audio/al](http://godoc.org/github.com/gazed/vu/audio/al) OpenAL bindings. Links the audio layer and the sound hardware.
* [device](http://godoc.org/github.com/gazed/vu/device)  Links the application to native OS specific window and user events.
* [load](http://godoc.org/github.com/gazed/vu/load) Asset loaders including models, textures, audio, shaders, and bitmapped fonts.
* [math/lin](http://godoc.org/github.com/gazed/vu/math/lin) Vector, matrix, quaternion, and transform linear math library.
* [physics](http://godoc.org/github.com/gazed/vu/physics) Repositions bodies based on simulated physics.
* [render](http://godoc.org/github.com/gazed/vu/render) 3D drawing and graphics interface.
* [render/gl](http://godoc.org/github.com/gazed/vu/render/gl) Generated OpenGL bindings. Links rendering system to graphics hardware.
* [render/gl/gen](http://godoc.org/github.com/gazed/vu/render/gl/gen)  OpenGL binding generator.

Less essential, but potentially more fun packages are:

* [eg](http://godoc.org/github.com/gazed/vu/eg) Examples that both demonstrate and validate the vu engine.
* [ai](http://godoc.org/github.com/gazed/vu/ai) Behaviour Tree for autonomous units.
* [grid](http://godoc.org/github.com/gazed/vu/grid) Grid based random level generators. A-star and flow field pathfinding.
* [synth](http://godoc.org/github.com/gazed/vu/synth) Procedural generation utilities.
* [tools/sdf](http://godoc.org/github.com/gazed/vu/tools/sdf) Signed distance field converstion utility.

Installation
-----

Ensure you have installed [Go](http://golang.org) 1.6+:

```bash
go get -u github.com/gazed/vu
```

Now you can build and run examples:

```bash
cd $GOPATH/src/github.com/gazed/vu/eg
go build .
./eg
```

**Build Dependencies**

* ``OS X``: Objective C and C compilers (clang) from Xcode command line tools.
* ``Windows``: C compiler (gcc) from mingw64-bit.

**Runtime Dependencies**

* OpenGL version 3.3 or later.
* OpenAL 64-bit version 2.1.

**Building on Windows**

* Vu has been built and tested on Windows using gcc from mingw64-bit.
  Mingw64 was installed to c:/mingw64.
  * Put OpenAL on the gcc library path by copying ``OpenAL32.dll`` to
    ``c:/mingw64/x86_64-w64-mingw32/lib/OpenAL32.dll``
* 64-bit OpenAL for Windows machines is available at http://openal.org/downloads.
  Running the OpenAL installer results in a ``c:/Windows/System32/OpenAL32.dll``.
* Building with Cygwin has not been attempted. It may have special needs.

Limitations
-----------

The engine and its packages include the essentials by design. In particular:

* There is no 3D editor.
* There is no networking package.
* Physics only handles boxes and spheres.
* The device layer interface provides only the absolute minimum from the underlying
  windowing system. Only OSX and Windows 7+ are currently supported.
* Rendering supports standard OpenGL 3.3 and later. OpenGL extensions are not used.
* The Windows platform is sometimes limited by the availability of OpenGL and OpenAL.
  Generally OpenGL issues are fixed by downloading manufacturer's graphic card drivers.
  However older laptops with Intel graphics don't always have OpenGL drivers.
