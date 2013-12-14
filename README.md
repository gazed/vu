vu
==

Vu (Virtual Universe) is a skeleton 3D engine written primarily in golang. 
Vu is composed of various packages which are detailed in the 
go docs and briefly summarized below.

Packages
--------

* ``vu`` The 3D application facing layer wraps and extends the other packages.
* ``vu/audio`` Positions and plays sounds in a 3D environment. 
* ``vu/audio/al`` OpenAL bindings. Links the audio layer and the sound hardware. 
* ``vu/data`` Resource data loaders including models, textures, audio, shaders, and bitmapped fonts.
* ``vu/device`` Links the application to native OS specific window and user events. 
* ``vu/math/lin`` Vector, matrix, quaternion, and transform math library.
* ``vu/move`` Repositions bodies based on simulated physics.
* ``vu/render`` 3D drawing interface.
* ``vu/render/gl`` Generated OpenGL bindings. Links the rendering system and the graphics hardware.
* ``vu/render/gl/gen`` OpenGL binding generator. 

Less essential, but potentially more fun packages are:

* ``vu/eg`` Examples that are used both to demonstrate and validate vu engine functionality.
* ``vu/grid`` Grid based random level generators.

Build
-----

* Ensure GOPATH contains the ``vu`` directory.
* Build from the ``vu`` directory using ``./build src`` or ``python build src``.
  All build output is located in the ``target`` directory.
* Build and run the examples from the ``vu/src/vu/eg`` directory using ``go build`` and ``./eg``. 

**Build Dependencies**

* go and standard go libraries.
* python for the build script.
* ``osx``: Objective C and C compilers (clang) from XCode command line tools.
* ``win``: C compiler (gcc) from mingw64-bit 
* git for product version numbering.

**Runtime Dependencies**

* OpenGL version 3.2 or later.
* OpenAL 64-bit version 2.1.

Limitations
-----------

The engine and its packages are bare bone by design - it's a skeleton 3D engine :). In particular:

* There is no game engine editor.
* Physics only handles boxes and sphere shapes. Only a few files of the bullet 
  physics engine were ported or wrapped to golang. Huge thank-you to bullet physics.
* Only one format is supported for each type of 3D resource data, e.g ``.obj`` for 3D models.
* The device layer interface provides only the absolute minimum from the underlying windowing system. 
  Only OSX and Windows 7 are currently supported.
* Rendering supports standard OpenGL 3.2 and later. OpenGL extensions are not used. 
* OpenGL may not be supported on all Windows graphics drivers (eg. laptop Intel based graphics)  
* 64-bit OpenAL may be difficult to locate for Windows machines. See 
  http://connect.creativelabs.com/openal/Downloads/oalinst.zip if/when their website is up.
* Building on Windows used golang with gcc from mingw64-bit. 
  Building with Cygwin may have special needs. 
