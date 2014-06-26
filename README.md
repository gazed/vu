vu
==

Vu (Virtual Universe) is a minimalist 3D engine written primarily in golang. 
Vu is composed of packages, detailed in the go docs, and briefly summarized below.

Packages
--------

* ``vu`` The 3D application facing layer wraps and extends the other packages.
* ``vu/audio`` Positions and plays sounds in a 3D environment. 
* ``vu/audio/al`` OpenAL bindings. Links the audio layer and the sound hardware. 
* ``vu/load`` Asset loaders including models, textures, audio, shaders, and bitmapped fonts.
* ``vu/device`` Links the application to native OS specific window and user events. 
* ``vu/math/lin`` Vector, matrix, quaternion, and transform linear math library.
* ``vu/move`` Repositions bodies based on simulated physics.
* ``vu/render`` 3D drawing interface.
* ``vu/render/gl`` Generated OpenGL bindings. Links rendering system to graphics hardware.
* ``vu/render/gl/gen`` OpenGL binding generator. 

Less essential, but potentially more fun packages are:

* ``vu/eg`` Examples that both demonstrate and validate the vu engine.
* ``vu/grid`` Grid based random level generators.
* ``vu/land`` Height map and land surface generator. 
* ``vu/panel`` 2D overlay layer for 3D applications.

Build
-----

* Ensure GOPATH contains the ``vu`` directory.
* Build using ``build`` from the ``vu`` directory. Eg:
    * OSX:  
        * ``cd vu``
        * ``./build src``
    * WIN:
        * ``cd vu``
        * ``python build src``.
* Build and run the examples from the ``vu/src/vu/eg`` directory:
        * ``cd vu/src/vu/eg``
        * ``go build``
        * ``./eg``

**Build Dependencies**

* go1.3
* python for the build script.
* ``OSX``: Objective C and C compilers (clang) from XCode command line tools.
* ``WIN``: C compiler (gcc) from mingw64-bit 

**Runtime Dependencies**

* OpenGL version 3.3 or later.
* OpenAL 64-bit version 2.1.

Limitations
-----------

The engine and its packages are minamalist by design. In particular:

* There is no 3D editor.
* There is no networking support.
* Physics only handles boxes and sphere shapes. 
* The device layer interface provides only the absolute minimum from the underlying
  windowing system. Only OSX, Windows 7 and 8 are currently supported.
* Rendering supports standard OpenGL 3.3 and later. OpenGL extensions are not used.
* Windows is limited by the availability of OpenGL and OpenAL. Generally
  OpenGL issues are fixed by downloading manufacturer's graphic card drivers.
  However older laptops with Intel graphics don't always have OpenGL drivers.

**Building on Windows**

* 64-bit OpenAL may be difficult to locate for Windows machines.
  Try ``http://kcat.strangesoft.net/openal.html/openal-soft-1.15.1-bin.zip``
* Bampf has been built and tested on Windows using gcc from mingw64-bit.
  Mingw64 was installed to c:/mingw64.
* Put OpenAL on the gcc library path by copying
  ``openal-soft-1.15.1-bin/Win64/soft_oal.dll`` to
  ``c:/mingw64/x86_64-w64-mingw32/lib/OpenAL32.dll``
* Building with Cygwin has not been attempted. It may have special needs.

