// SPDX-FileCopyrightText : © 2025 Galvanized Logic Inc.
// SPDX-License-Identifier: BSD-2-Clause

//go:build !windows

// Package al provides golang audio library bindings for OpenAL.
// Official OpenAL documentation can be found online. Prepend "AL_"
// to the function or constant names found in this package.
// Refer to the official OpenAL documentation for more information.
//
// Package al is provided as part of the vu (virtual universe) 3D engine.
package al

// Design Notes:
// These bindings were based on the OpenAL header files found at:
//   http://repo.or.cz/w/openal-soft.git/blob/6dab9d54d1719105e0183f941a2b3dd36e9ba902:/include/AL/al.h
//   http://repo.or.cz/w/openal-soft.git/blob/6dab9d54d1719105e0183f941a2b3dd36e9ba902:/include/AL/alc.h
// Check information available at openal.org.

// #cgo darwin  LDFLAGS: -framework OpenAL
// #cgo linux   LDFLAGS: -lopenal -ldl
// #cgo windows LDFLAGS: -lOpenAL32
//
// #include <stdlib.h>
// #if defined(__APPLE__)
// #include <dlfcn.h>
// #elif defined(_WIN32)
// #define WIN32_LEAN_AND_MEAN 1
// #include <windows.h>
// #else
// #include <dlfcn.h>
// #endif
//
// #ifdef _WIN32
// static HMODULE hmod = NULL;
// #elif !defined __APPLE__
// static void* plib = NULL;
// #endif
//
// // Helps bind function pointers to c functions.
// static void* bindMethod(const char* name) {
// #ifdef __APPLE__
// 	return dlsym(RTLD_DEFAULT, name);
// #elif _WIN32
// 	if(hmod == NULL) {
// 		hmod = LoadLibraryA("OpenAL32.dll");
// 	}
// 	return GetProcAddress(hmod, (LPCSTR)name);
// #else
// 	if(plib == NULL) {
// 		plib = dlopen("libopenal.so", RTLD_LAZY);
// 	}
// 	return dlsym(plib, name);
// #endif
// }
//
// #if defined(_WIN32)
//  #define AL_APIENTRY __cdecl
//  #define ALC_APIENTRY __cdecl
// #else
//  #define AL_APIENTRY
//  #define ALC_APIENTRY
// #endif
//
// // AL/al.h typedefs
// typedef char ALboolean;
// typedef char ALchar;
// typedef signed char ALbyte;
// typedef unsigned char ALubyte;
// typedef unsigned short ALushort;
// typedef int ALint;
// typedef unsigned int ALuint;
// typedef int ALsizei;
// typedef int ALenum;
// typedef float ALfloat;
// typedef double ALdouble;
// typedef void ALvoid;
//
// #ifndef AL_API
// #define AL_API extern
// #endif
//
// // AL/alc.h typedefs
// typedef struct ALCdevice_struct ALCdevice;
// typedef struct ALCcontext_struct ALCcontext;
// typedef char ALCboolean;
// typedef char ALCchar;
// typedef signed char ALCbyte;
// typedef unsigned char ALCubyte;
// typedef unsigned short ALCushort;
// typedef int ALCint;
// typedef unsigned int ALCuint;
// typedef int ALCsizei;
// typedef int ALCenum;
// typedef void ALCvoid;
//
// #ifndef ALC_API
// #define ALC_API extern
// #endif
//
// // AL/al.h pointers to functions bound to the OS specific library.
// void           (AL_APIENTRY *pfn_alEnable)( ALenum capability );
// void           (AL_APIENTRY *pfn_alDisable)( ALenum capability );
// ALboolean      (AL_APIENTRY *pfn_alIsEnabled)( ALenum capability );
// const ALchar*  (AL_APIENTRY *pfn_alGetString)( ALenum param );
// void           (AL_APIENTRY *pfn_alGetBooleanv)( ALenum param, ALboolean* data );
// void           (AL_APIENTRY *pfn_alGetIntegerv)( ALenum param, ALint* data );
// void           (AL_APIENTRY *pfn_alGetFloatv)( ALenum param, ALfloat* data );
// void           (AL_APIENTRY *pfn_alGetDoublev)( ALenum param, ALdouble* data );
// ALboolean      (AL_APIENTRY *pfn_alGetBoolean)( ALenum param );
// ALint          (AL_APIENTRY *pfn_alGetInteger)( ALenum param );
// ALfloat        (AL_APIENTRY *pfn_alGetFloat)( ALenum param );
// ALdouble       (AL_APIENTRY *pfn_alGetDouble)( ALenum param );
// ALenum         (AL_APIENTRY *pfn_alGetError)( void );
// ALboolean      (AL_APIENTRY *pfn_alIsExtensionPresent)(const ALchar* extname );
// void*          (AL_APIENTRY *pfn_alGetProcAddress)( const ALchar* fname );
// ALenum         (AL_APIENTRY *pfn_alGetEnumValue)( const ALchar* ename );
// void           (AL_APIENTRY *pfn_alListenerf)( ALenum param, ALfloat value );
// void           (AL_APIENTRY *pfn_alListener3f)( ALenum param, ALfloat value1, ALfloat value2, ALfloat value3 );
// void           (AL_APIENTRY *pfn_alListenerfv)( ALenum param, const ALfloat* values );
// void           (AL_APIENTRY *pfn_alListeneri)( ALenum param, ALint value );
// void           (AL_APIENTRY *pfn_alListener3i)( ALenum param, ALint value1, ALint value2, ALint value3 );
// void           (AL_APIENTRY *pfn_alListeneriv)( ALenum param, const ALint* values );
// void           (AL_APIENTRY *pfn_alGetListenerf)( ALenum param, ALfloat* value );
// void           (AL_APIENTRY *pfn_alGetListener3f)( ALenum param, ALfloat *value1, ALfloat *value2, ALfloat *value3 );
// void           (AL_APIENTRY *pfn_alGetListenerfv)( ALenum param, ALfloat* values );
// void           (AL_APIENTRY *pfn_alGetListeneri)( ALenum param, ALint* value );
// void           (AL_APIENTRY *pfn_alGetListener3i)( ALenum param, ALint *value1, ALint *value2, ALint *value3 );
// void           (AL_APIENTRY *pfn_alGetListeneriv)( ALenum param, ALint* values );
// void           (AL_APIENTRY *pfn_alGenSources)( ALsizei n, ALuint* sources );
// void           (AL_APIENTRY *pfn_alDeleteSources)( ALsizei n, const ALuint* sources );
// ALboolean      (AL_APIENTRY *pfn_alIsSource)( ALuint sid );
// void           (AL_APIENTRY *pfn_alSourcef)( ALuint sid, ALenum param, ALfloat value);
// void           (AL_APIENTRY *pfn_alSource3f)( ALuint sid, ALenum param, ALfloat value1, ALfloat value2, ALfloat value3 );
// void           (AL_APIENTRY *pfn_alSourcefv)( ALuint sid, ALenum param, const ALfloat* values );
// void           (AL_APIENTRY *pfn_alSourcei)( ALuint sid, ALenum param, ALint value);
// void           (AL_APIENTRY *pfn_alSource3i)( ALuint sid, ALenum param, ALint value1, ALint value2, ALint value3 );
// void           (AL_APIENTRY *pfn_alSourceiv)( ALuint sid, ALenum param, const ALint* values );
// void           (AL_APIENTRY *pfn_alGetSourcef)( ALuint sid, ALenum param, ALfloat* value );
// void           (AL_APIENTRY *pfn_alGetSource3f)( ALuint sid, ALenum param, ALfloat* value1, ALfloat* value2, ALfloat* value3);
// void           (AL_APIENTRY *pfn_alGetSourcefv)( ALuint sid, ALenum param, ALfloat* values );
// void           (AL_APIENTRY *pfn_alGetSourcei)( ALuint sid, ALenum param, ALint* value );
// void           (AL_APIENTRY *pfn_alGetSource3i)( ALuint sid, ALenum param, ALint* value1, ALint* value2, ALint* value3);
// void           (AL_APIENTRY *pfn_alGetSourceiv)( ALuint sid, ALenum param, ALint* values );
// void           (AL_APIENTRY *pfn_alSourcePlayv)( ALsizei ns, const ALuint *sids );
// void           (AL_APIENTRY *pfn_alSourceStopv)( ALsizei ns, const ALuint *sids );
// void           (AL_APIENTRY *pfn_alSourceRewindv)( ALsizei ns, const ALuint *sids );
// void           (AL_APIENTRY *pfn_alSourcePausev)( ALsizei ns, const ALuint *sids );
// void           (AL_APIENTRY *pfn_alSourcePlay)( ALuint sid );
// void           (AL_APIENTRY *pfn_alSourceStop)( ALuint sid );
// void           (AL_APIENTRY *pfn_alSourceRewind)( ALuint sid );
// void           (AL_APIENTRY *pfn_alSourcePause)( ALuint sid );
// void           (AL_APIENTRY *pfn_alSourceQueueBuffers)(ALuint sid, ALsizei numEntries, const ALuint *bids );
// void           (AL_APIENTRY *pfn_alSourceUnqueueBuffers)(ALuint sid, ALsizei numEntries, ALuint *bids );
// void           (AL_APIENTRY *pfn_alGenBuffers)( ALsizei n, ALuint* buffers );
// void           (AL_APIENTRY *pfn_alDeleteBuffers)( ALsizei n, const ALuint* buffers );
// ALboolean      (AL_APIENTRY *pfn_alIsBuffer)( ALuint bid );
// void           (AL_APIENTRY *pfn_alBufferData)( ALuint bid, ALenum format, const ALvoid* data, ALsizei size, ALsizei freq );
// void           (AL_APIENTRY *pfn_alBufferf)( ALuint bid, ALenum param, ALfloat value);
// void           (AL_APIENTRY *pfn_alBuffer3f)( ALuint bid, ALenum param, ALfloat value1, ALfloat value2, ALfloat value3 );
// void           (AL_APIENTRY *pfn_alBufferfv)( ALuint bid, ALenum param, const ALfloat* values );
// void           (AL_APIENTRY *pfn_alBufferi)( ALuint bid, ALenum param, ALint value);
// void           (AL_APIENTRY *pfn_alBuffer3i)( ALuint bid, ALenum param, ALint value1, ALint value2, ALint value3 );
// void           (AL_APIENTRY *pfn_alBufferiv)( ALuint bid, ALenum param, const ALint* values );
// void           (AL_APIENTRY *pfn_alGetBufferf)( ALuint bid, ALenum param, ALfloat* value );
// void           (AL_APIENTRY *pfn_alGetBuffer3f)( ALuint bid, ALenum param, ALfloat* value1, ALfloat* value2, ALfloat* value3);
// void           (AL_APIENTRY *pfn_alGetBufferfv)( ALuint bid, ALenum param, ALfloat* values );
// void           (AL_APIENTRY *pfn_alGetBufferi)( ALuint bid, ALenum param, ALint* value );
// void           (AL_APIENTRY *pfn_alGetBuffer3i)( ALuint bid, ALenum param, ALint* value1, ALint* value2, ALint* value3);
// void           (AL_APIENTRY *pfn_alGetBufferiv)( ALuint bid, ALenum param, ALint* values );
// void           (AL_APIENTRY *pfn_alDopplerFactor)( ALfloat value );
// void           (AL_APIENTRY *pfn_alDopplerVelocity)( ALfloat value );
// void           (AL_APIENTRY *pfn_alSpeedOfSound)( ALfloat value );
// void           (AL_APIENTRY *pfn_alDistanceModel)( ALenum distanceModel );
//
// // AL/al.h wrappers for the go bindings.
// AL_API void          AL_APIENTRY wrap_alEnable( int capability ) { (*pfn_alEnable)( capability ); }
// AL_API void          AL_APIENTRY wrap_alDisable( int capability ) { (*pfn_alDisable)( capability ); }
// AL_API unsigned int  AL_APIENTRY wrap_alIsEnabled( int capability ) { return (*pfn_alIsEnabled)( capability ); }
// AL_API const char*   AL_APIENTRY wrap_alGetString( int param ) { return (*pfn_alGetString)( param ); }
// AL_API void          AL_APIENTRY wrap_alGetBooleanv( int param, char* data ) { (*pfn_alGetBooleanv)( param, data ); }
// AL_API void          AL_APIENTRY wrap_alGetIntegerv( int param, int* data ) { (*pfn_alGetIntegerv)( param, data ); }
// AL_API void          AL_APIENTRY wrap_alGetFloatv( int param, float* data ) { (*pfn_alGetFloatv)( param, data ); }
// AL_API void          AL_APIENTRY wrap_alGetDoublev( int param, double* data ) { (*pfn_alGetDoublev)( param, data );}
// AL_API ALboolean     AL_APIENTRY wrap_alGetBoolean( int param ) { return (*pfn_alGetBoolean)( param ); }
// AL_API ALint         AL_APIENTRY wrap_alGetInteger( int param ) { return (*pfn_alGetInteger)( param ); }
// AL_API ALfloat       AL_APIENTRY wrap_alGetFloat( int param ) { return (*pfn_alGetFloat)( param ); }
// AL_API ALdouble      AL_APIENTRY wrap_alGetDouble( int param ) { return (*pfn_alGetDouble)( param ); }
// AL_API ALenum        AL_APIENTRY wrap_alGetError( void ) { return (*pfn_alGetError)(); }
// AL_API ALboolean     AL_APIENTRY wrap_alIsExtensionPresent( const char* extname ) { return (*pfn_alIsExtensionPresent)( extname ); }
// AL_API void*         AL_APIENTRY wrap_alGetProcAddress( const char* fname ) { return (*pfn_alGetProcAddress)( fname ); }
// AL_API ALenum        AL_APIENTRY wrap_alGetEnumValue( const char* ename ) { return (*pfn_alGetEnumValue)( ename ); }
// AL_API void          AL_APIENTRY wrap_alListenerf( int param, float value ) { (*pfn_alListenerf)( param, value ); }
// AL_API void          AL_APIENTRY wrap_alListener3f( int param, float value1, float value2, float value3 ) { (*pfn_alListener3f)( param, value1, value2, value3 ); }
// AL_API void          AL_APIENTRY wrap_alListenerfv( int param, const float* values ) { (*pfn_alListenerfv)( param, values ); }
// AL_API void          AL_APIENTRY wrap_alListeneri( int param, int value ) { (*pfn_alListeneri)( param, value ); }
// AL_API void          AL_APIENTRY wrap_alListener3i( int param, int value1, int value2, int value3 ) { (*pfn_alListener3i)( param, value1, value2, value3 ); }
// AL_API void          AL_APIENTRY wrap_alListeneriv( int param, const int* values ) { (*pfn_alListeneriv)( param, values ); }
// AL_API void          AL_APIENTRY wrap_alGetListenerf( int param, float* value ) { (*pfn_alGetListenerf)( param, value ); }
// AL_API void          AL_APIENTRY wrap_alGetListener3f( int param, float *value1, float *value2, float *value3 ) { (*pfn_alGetListener3f)( param, value1, value2, value3 ); }
// AL_API void          AL_APIENTRY wrap_alGetListenerfv( int param, float* values ) { (*pfn_alGetListenerfv)( param, values ); }
// AL_API void          AL_APIENTRY wrap_alGetListeneri( int param, int* value ) { (*pfn_alGetListeneri)( param, value ); }
// AL_API void          AL_APIENTRY wrap_alGetListener3i( int param, int *value1, int *value2, int *value3 ) { (*pfn_alGetListener3i)( param, value1, value2, value3 ); }
// AL_API void          AL_APIENTRY wrap_alGetListeneriv( int param, int* values ) { (*pfn_alGetListeneriv)( param, values ); }
// AL_API void          AL_APIENTRY wrap_alGenSources( int n, unsigned int* sources ) { (*pfn_alGenSources)( n, sources ); }
// AL_API void          AL_APIENTRY wrap_alDeleteSources( int n, const unsigned int* sources ) { (*pfn_alDeleteSources)( n, sources ); }
// AL_API ALboolean     AL_APIENTRY wrap_alIsSource( unsigned int sid ) { return (*pfn_alIsSource)( sid ); }
// AL_API void          AL_APIENTRY wrap_alSourcef( unsigned int sid, int param, float value ) { (*pfn_alSourcef)( sid, param, value ); }
// AL_API void          AL_APIENTRY wrap_alSource3f( unsigned int sid, int param, float value1, float value2, float value3 ) { (*pfn_alSource3f)( sid, param, value1, value2, value3 ); }
// AL_API void          AL_APIENTRY wrap_alSourcefv( unsigned int sid, int param, const float* values ) { (*pfn_alSourcefv)( sid, param, values ); }
// AL_API void          AL_APIENTRY wrap_alSourcei( unsigned int sid, int param, int value ) { (*pfn_alSourcei)( sid, param, value ); }
// AL_API void          AL_APIENTRY wrap_alSource3i( unsigned int sid, int param, int value1, int value2, int value3 ) { (*pfn_alSource3i)( sid, param, value1, value2, value3 ); }
// AL_API void          AL_APIENTRY wrap_alSourceiv( unsigned int sid, int param, const int* values ) { (*pfn_alSourceiv)( sid, param, values ); }
// AL_API void          AL_APIENTRY wrap_alGetSourcef( unsigned int sid, int param, float* value ) { (*pfn_alGetSourcef)( sid, param, value ); }
// AL_API void          AL_APIENTRY wrap_alGetSource3f( unsigned int sid, int param, float* value1, float* value2, float* value3) { (*pfn_alGetSource3f)( sid, param, value1, value2, value3); }
// AL_API void          AL_APIENTRY wrap_alGetSourcefv( unsigned int sid, int param, float* values ) { (*pfn_alGetSourcefv)( sid, param, values ); }
// AL_API void          AL_APIENTRY wrap_alGetSourcei( unsigned int sid,  int param, int* value ) { (*pfn_alGetSourcei)( sid,  param, value ); }
// AL_API void          AL_APIENTRY wrap_alGetSource3i( unsigned int sid, int param, int* value1, int* value2, int* value3) { (*pfn_alGetSource3i)( sid, param, value1, value2, value3); }
// AL_API void          AL_APIENTRY wrap_alGetSourceiv( unsigned int sid,  int param, int* values ) { (*pfn_alGetSourceiv)( sid, param, values ); }
// AL_API void          AL_APIENTRY wrap_alSourcePlayv( int ns, const unsigned int *sids ) { (*pfn_alSourcePlayv)( ns, sids ); }
// AL_API void          AL_APIENTRY wrap_alSourceStopv( int ns, const unsigned int *sids ) { (*pfn_alSourceStopv)( ns, sids ); }
// AL_API void          AL_APIENTRY wrap_alSourceRewindv( int ns, const unsigned int *sids ) { (*pfn_alSourceRewindv)( ns, sids ); }
// AL_API void          AL_APIENTRY wrap_alSourcePausev( int ns, const unsigned int *sids ) { (*pfn_alSourcePausev)( ns, sids ); }
// AL_API void          AL_APIENTRY wrap_alSourcePlay( unsigned int sid ) { (*pfn_alSourcePlay)( sid ); }
// AL_API void          AL_APIENTRY wrap_alSourceStop( unsigned int sid ) { (*pfn_alSourceStop)( sid ); }
// AL_API void          AL_APIENTRY wrap_alSourceRewind( unsigned int sid ) { (*pfn_alSourceRewind)( sid ); }
// AL_API void          AL_APIENTRY wrap_alSourcePause( unsigned int sid ) { (*pfn_alSourcePause)( sid ); }
// AL_API void          AL_APIENTRY wrap_alSourceQueueBuffers( unsigned int sid, int numEntries, const unsigned int *bids ) { (*pfn_alSourceQueueBuffers)( sid, numEntries, bids ); }
// AL_API void          AL_APIENTRY wrap_alSourceUnqueueBuffers( unsigned int sid, int numEntries, unsigned int *bids ) {(*pfn_alSourceUnqueueBuffers)( sid, numEntries, bids ); }
// AL_API void          AL_APIENTRY wrap_alGenBuffers( int n, unsigned int* buffers ) { (*pfn_alGenBuffers)( n, buffers ); }
// AL_API void          AL_APIENTRY wrap_alDeleteBuffers( int n, const unsigned int* buffers ) { (*pfn_alDeleteBuffers)( n, buffers ); }
// AL_API ALboolean     AL_APIENTRY wrap_alIsBuffer( unsigned int bid ) { return (*pfn_alIsBuffer)( bid ); }
// AL_API void          AL_APIENTRY wrap_alBufferData( unsigned int bid, int format, const ALvoid* data, int size, int freq ) { (*pfn_alBufferData)( bid, format, data, size, freq ); }
// AL_API void          AL_APIENTRY wrap_alBufferf( unsigned int bid, int param, float value ) { (*pfn_alBufferf)( bid, param, value ); }
// AL_API void          AL_APIENTRY wrap_alBuffer3f( unsigned int bid, int param, float value1, float value2, float value3 ) { (*pfn_alBuffer3f)( bid, param, value1, value2, value3 ); }
// AL_API void          AL_APIENTRY wrap_alBufferfv( unsigned int bid, int param, const float* values ) { (*pfn_alBufferfv)( bid, param, values ); }
// AL_API void          AL_APIENTRY wrap_alBufferi( unsigned int bid, int param, int value ) { (*pfn_alBufferi)( bid, param, value ); }
// AL_API void          AL_APIENTRY wrap_alBuffer3i( unsigned int bid, int param, int value1, int value2, int value3 ) { (*pfn_alBuffer3i)( bid, param, value1, value2, value3 ); }
// AL_API void          AL_APIENTRY wrap_alBufferiv( unsigned int bid, int param, const int* values ) { (*pfn_alBufferiv)( bid, param, values ); }
// AL_API void          AL_APIENTRY wrap_alGetBufferf( unsigned int bid, int param, float* value ) { (*pfn_alGetBufferf)( bid, param, value ); }
// AL_API void          AL_APIENTRY wrap_alGetBuffer3f( unsigned int bid, int param, float* value1, float* value2, float* value3) { (*pfn_alGetBuffer3f)( bid, param, value1, value2, value3); }
// AL_API void          AL_APIENTRY wrap_alGetBufferfv( unsigned int bid, int param, float* values ) { (*pfn_alGetBufferfv)( bid, param, values ); }
// AL_API void          AL_APIENTRY wrap_alGetBufferi( unsigned int bid, int param, int* value ) { (*pfn_alGetBufferi)( bid, param, value ); }
// AL_API void          AL_APIENTRY wrap_alGetBuffer3i( unsigned int bid, int param, int* value1, int* value2, int* value3) { (*pfn_alGetBuffer3i)( bid, param, value1, value2, value3); }
// AL_API void          AL_APIENTRY wrap_alGetBufferiv( unsigned int bid, int param, int* values ) { (*pfn_alGetBufferiv)( bid, param, values ); }
// AL_API void          AL_APIENTRY wrap_alDopplerFactor( float value ) { (*pfn_alDopplerFactor)( value ); }
// AL_API void          AL_APIENTRY wrap_alDopplerVelocity( float value ) { (*pfn_alDopplerVelocity)( value ); }
// AL_API void          AL_APIENTRY wrap_alSpeedOfSound( float value ) { (*pfn_alSpeedOfSound)( value ); }
// AL_API void          AL_APIENTRY wrap_alDistanceModel( int distanceModel ) { (*pfn_alDistanceModel)( distanceModel ); }
//
// // AL/alc.h pointers to functions bound to the OS specific library.
// ALCcontext *   (ALC_APIENTRY *pfn_alcCreateContext) (ALCdevice *device, const ALCint *attrlist);
// ALCboolean     (ALC_APIENTRY *pfn_alcMakeContextCurrent)( ALCcontext *context );
// void           (ALC_APIENTRY *pfn_alcProcessContext)( ALCcontext *context );
// void           (ALC_APIENTRY *pfn_alcSuspendContext)( ALCcontext *context );
// void           (ALC_APIENTRY *pfn_alcDestroyContext)( ALCcontext *context );
// ALCcontext *   (ALC_APIENTRY *pfn_alcGetCurrentContext)( void );
// ALCdevice *    (ALC_APIENTRY *pfn_alcGetContextsDevice)( ALCcontext *context );
// ALCdevice *    (ALC_APIENTRY *pfn_alcOpenDevice)( const ALCchar *devicename );
// ALCboolean     (ALC_APIENTRY *pfn_alcCloseDevice)( ALCdevice *device );
// ALCenum        (ALC_APIENTRY *pfn_alcGetError)( ALCdevice *device );
// ALCboolean     (ALC_APIENTRY *pfn_alcIsExtensionPresent)( ALCdevice *device, const ALCchar *extname );
// void *         (ALC_APIENTRY *pfn_alcGetProcAddress)(ALCdevice *device, const ALCchar *funcname );
// ALCenum        (ALC_APIENTRY *pfn_alcGetEnumValue)(ALCdevice *device, const ALCchar *enumname );
// const ALCchar* (ALC_APIENTRY *pfn_alcGetString)( ALCdevice *device, ALCenum param );
// void           (ALC_APIENTRY *pfn_alcGetIntegerv)( ALCdevice *device, ALCenum param, ALCsizei size, ALCint *data );
// ALCdevice *    (ALC_APIENTRY *pfn_alcCaptureOpenDevice)( const ALCchar *devicename, ALCuint frequency, ALCenum format, ALCsizei buffersize );
// ALCboolean     (ALC_APIENTRY *pfn_alcCaptureCloseDevice)( ALCdevice *device );
// void           (ALC_APIENTRY *pfn_alcCaptureStart)( ALCdevice *device );
// void           (ALC_APIENTRY *pfn_alcCaptureStop)( ALCdevice *device );
// void           (ALC_APIENTRY *pfn_alcCaptureSamples)( ALCdevice *device, ALCvoid *buffer, ALCsizei samples );
//
// // AL/alc.h wrappers for the go bindings.
// ALC_API uintptr_t    ALC_APIENTRY wrap_alcCreateContext( uintptr_t device, const int* attrlist ) { return (uintptr_t)(*pfn_alcCreateContext)((ALCdevice *)device, attrlist); }
// ALC_API ALCboolean   ALC_APIENTRY wrap_alcMakeContextCurrent( uintptr_t context ) { return (*pfn_alcMakeContextCurrent)( (ALCcontext *)context ); }
// ALC_API void         ALC_APIENTRY wrap_alcProcessContext( uintptr_t context ) { (*pfn_alcProcessContext)( (ALCcontext *)context ); }
// ALC_API void         ALC_APIENTRY wrap_alcSuspendContext( uintptr_t context ) { (*pfn_alcSuspendContext)( (ALCcontext *)context ); }
// ALC_API void         ALC_APIENTRY wrap_alcDestroyContext( uintptr_t context ) { (*pfn_alcDestroyContext)( (ALCcontext *)context ); }
// ALC_API uintptr_t    ALC_APIENTRY wrap_alcGetCurrentContext( void ) { return (uintptr_t)(*pfn_alcGetCurrentContext)(); }
// ALC_API uintptr_t    ALC_APIENTRY wrap_alcGetContextsDevice( uintptr_t context ) { return (uintptr_t)(*pfn_alcGetContextsDevice)( (ALCcontext *)context ); }
// ALC_API uintptr_t    ALC_APIENTRY wrap_alcOpenDevice( const char *devicename ) { return (uintptr_t)(*pfn_alcOpenDevice)( devicename ); }
// ALC_API ALCboolean   ALC_APIENTRY wrap_alcCloseDevice( uintptr_t device ) { return (*pfn_alcCloseDevice)( (ALCdevice *)device ); }
// ALC_API ALCenum      ALC_APIENTRY wrap_alcGetError( uintptr_t device ) { return (*pfn_alcGetError)( (ALCdevice *)device ); }
// ALC_API ALCboolean   ALC_APIENTRY wrap_alcIsExtensionPresent( uintptr_t device, const char *extname ) { return (*pfn_alcIsExtensionPresent)( (ALCdevice *)device, extname ); }
// ALC_API void  *      ALC_APIENTRY wrap_alcGetProcAddress( uintptr_t device, const char *funcname ) { return (*pfn_alcGetProcAddress)( (ALCdevice *)device, funcname ); }
// ALC_API ALCenum      ALC_APIENTRY wrap_alcGetEnumValue( uintptr_t device, const char *enumname ) { return (*pfn_alcGetEnumValue)( (ALCdevice *)device, enumname ); }
// ALC_API const char * ALC_APIENTRY wrap_alcGetString( uintptr_t device, int param ) { return (*pfn_alcGetString)( (ALCdevice *)device, param ); }
// ALC_API void         ALC_APIENTRY wrap_alcGetIntegerv( uintptr_t device, int param, int size, int *data ) { (*pfn_alcGetIntegerv)( (ALCdevice *)device, param, size, data ); }
// ALC_API uintptr_t    ALC_APIENTRY wrap_alcCaptureOpenDevice( const char *devicename, unsigned int frequency, int format, int buffersize ) { return (uintptr_t)(*pfn_alcCaptureOpenDevice)( devicename, frequency, format, buffersize ); }
// ALC_API ALCboolean   ALC_APIENTRY wrap_alcCaptureCloseDevice( uintptr_t device ) { return (*pfn_alcCaptureCloseDevice)( (ALCdevice *)device ); }
// ALC_API void         ALC_APIENTRY wrap_alcCaptureStart( uintptr_t device ) { (*pfn_alcCaptureStart)( (ALCdevice *)device ); }
// ALC_API void         ALC_APIENTRY wrap_alcCaptureStop( uintptr_t device ) { (*pfn_alcCaptureStop)( (ALCdevice *)device ); }
// ALC_API void         ALC_APIENTRY wrap_alcCaptureSamples( uintptr_t device, ALCvoid *buffer, int samples ) { (*pfn_alcCaptureSamples)( (ALCdevice *)device, buffer, samples ); }
//
// void al_init() {
//    // AL/al.h
//    pfn_alEnable                  = bindMethod("alEnable");
//    pfn_alDisable                 = bindMethod("alDisable");
//    pfn_alIsEnabled               = bindMethod("alIsEnabled");
//    pfn_alGetString               = bindMethod("alGetString");
//    pfn_alGetBooleanv             = bindMethod("alGetBooleanv");
//    pfn_alGetIntegerv             = bindMethod("alGetIntegerv");
//    pfn_alGetFloatv               = bindMethod("alGetFloatv");
//    pfn_alGetDoublev              = bindMethod("alGetDoublev");
//    pfn_alGetBoolean              = bindMethod("alGetBoolean");
//    pfn_alGetInteger              = bindMethod("alGetInteger");
//    pfn_alGetFloat                = bindMethod("alGetFloat");
//    pfn_alGetDouble               = bindMethod("alGetDouble");
//    pfn_alGetError                = bindMethod("alGetError");
//    pfn_alIsExtensionPresent      = bindMethod("alIsExtensionPresent");
//    pfn_alGetProcAddress          = bindMethod("alGetProcAddress");
//    pfn_alGetEnumValue            = bindMethod("alGetEnumValue");
//    pfn_alListenerf               = bindMethod("alListenerf");
//    pfn_alListener3f              = bindMethod("alListener3f");
//    pfn_alListenerfv              = bindMethod("alListenerfv");
//    pfn_alListeneri               = bindMethod("alListeneri");
//    pfn_alListener3i              = bindMethod("alListener3i");
//    pfn_alListeneriv              = bindMethod("alListeneriv");
//    pfn_alGetListenerf            = bindMethod("alGetListenerf");
//    pfn_alGetListener3f           = bindMethod("alGetListener3f");
//    pfn_alGetListenerfv           = bindMethod("alGetListenerfv");
//    pfn_alGetListeneri            = bindMethod("alGetListeneri");
//    pfn_alGetListener3i           = bindMethod("alGetListener3i");
//    pfn_alGetListeneriv           = bindMethod("alGetListeneriv");
//    pfn_alGenSources              = bindMethod("alGenSources");
//    pfn_alDeleteSources           = bindMethod("alDeleteSources");
//    pfn_alIsSource                = bindMethod("alIsSource");
//    pfn_alSourcef                 = bindMethod("alSourcef");
//    pfn_alSource3f                = bindMethod("alSource3f");
//    pfn_alSourcefv                = bindMethod("alSourcefv");
//    pfn_alSourcei                 = bindMethod("alSourcei");
//    pfn_alSource3i                = bindMethod("alSource3i");
//    pfn_alSourceiv                = bindMethod("alSourceiv");
//    pfn_alGetSourcef              = bindMethod("alGetSourcef");
//    pfn_alGetSource3f             = bindMethod("alGetSource3f");
//    pfn_alGetSourcefv             = bindMethod("alGetSourcefv");
//    pfn_alGetSourcei              = bindMethod("alGetSourcei");
//    pfn_alGetSource3i             = bindMethod("alGetSource3i");
//    pfn_alGetSourceiv             = bindMethod("alGetSourceiv");
//    pfn_alSourcePlayv             = bindMethod("alSourcePlayv");
//    pfn_alSourceStopv             = bindMethod("alSourceStopv");
//    pfn_alSourceRewindv           = bindMethod("alSourceRewindv");
//    pfn_alSourcePausev            = bindMethod("alSourcePausev");
//    pfn_alSourcePlay              = bindMethod("alSourcePlay");
//    pfn_alSourceStop              = bindMethod("alSourceStop");
//    pfn_alSourceRewind            = bindMethod("alSourceRewind");
//    pfn_alSourcePause             = bindMethod("alSourcePause");
//    pfn_alSourceQueueBuffers      = bindMethod("alSourceQueueBuffers");
//    pfn_alSourceUnqueueBuffers    = bindMethod("alSourceUnqueueBuffers");
//    pfn_alGenBuffers              = bindMethod("alGenBuffers");
//    pfn_alDeleteBuffers           = bindMethod("alDeleteBuffers");
//    pfn_alIsBuffer                = bindMethod("alIsBuffer");
//    pfn_alBufferData              = bindMethod("alBufferData");
//    pfn_alBufferf                 = bindMethod("alBufferf");
//    pfn_alBuffer3f                = bindMethod("alBuffer3f");
//    pfn_alBufferfv                = bindMethod("alBufferfv");
//    pfn_alBufferi                 = bindMethod("alBufferi");
//    pfn_alBuffer3i                = bindMethod("alBuffer3i");
//    pfn_alBufferiv                = bindMethod("alBufferiv");
//    pfn_alGetBufferf              = bindMethod("alGetBufferf");
//    pfn_alGetBuffer3f             = bindMethod("alGetBuffer3f");
//    pfn_alGetBufferfv             = bindMethod("alGetBufferfv");
//    pfn_alGetBufferi              = bindMethod("alGetBufferi");
//    pfn_alGetBuffer3i             = bindMethod("alGetBuffer3i");
//    pfn_alGetBufferiv             = bindMethod("alGetBufferiv");
//    pfn_alDopplerFactor           = bindMethod("alDopplerFactor");
//    pfn_alDopplerVelocity         = bindMethod("alDopplerVelocity");
//    pfn_alSpeedOfSound            = bindMethod("alSpeedOfSound");
//    pfn_alDistanceModel           = bindMethod("alDistanceModel");
//
//    // AL/alc.h
//    pfn_alcCreateContext          = bindMethod("alcCreateContext");
//    pfn_alcMakeContextCurrent     = bindMethod("alcMakeContextCurrent");
//    pfn_alcProcessContext         = bindMethod("alcProcessContext");
//    pfn_alcSuspendContext         = bindMethod("alcSuspendContext");
//    pfn_alcDestroyContext         = bindMethod("alcDestroyContext");
//    pfn_alcGetCurrentContext      = bindMethod("alcGetCurrentContext");
//    pfn_alcGetContextsDevice      = bindMethod("alcGetContextsDevice");
//    pfn_alcOpenDevice             = bindMethod("alcOpenDevice");
//    pfn_alcCloseDevice            = bindMethod("alcCloseDevice");
//    pfn_alcGetError               = bindMethod("alcGetError");
//    pfn_alcIsExtensionPresent     = bindMethod("alcIsExtensionPresent");
//    pfn_alcGetProcAddress         = bindMethod("alcGetProcAddress");
//    pfn_alcGetEnumValue           = bindMethod("alcGetEnumValue");
//    pfn_alcGetString              = bindMethod("alcGetString");
//    pfn_alcGetIntegerv            = bindMethod("alcGetIntegerv");
//    pfn_alcCaptureOpenDevice      = bindMethod("alcCaptureOpenDevice");
//    pfn_alcCaptureCloseDevice     = bindMethod("alcCaptureCloseDevice");
//    pfn_alcCaptureStart           = bindMethod("alcCaptureStart");
//    pfn_alcCaptureStop            = bindMethod("alcCaptureStop");
//    pfn_alcCaptureSamples         = bindMethod("alcCaptureSamples");
// }
//
import "C"
import "unsafe"
import "fmt"

// AL/al.h constants (with AL_ removed). Refer to the original header for constant documentation.
const (
	FALSE                     = 0
	TRUE                      = 1
	NONE                      = 0
	NO_ERROR                  = 0
	SOURCE_RELATIVE           = 0x202
	CONE_INNER_ANGLE          = 0x1001
	CONE_OUTER_ANGLE          = 0x1002
	PITCH                     = 0x1003
	POSITION                  = 0x1004
	DIRECTION                 = 0x1005
	VELOCITY                  = 0x1006
	LOOPING                   = 0x1007
	BUFFER                    = 0x1009
	GAIN                      = 0x100A
	MIN_GAIN                  = 0x100D
	MAX_GAIN                  = 0x100E
	ORIENTATION               = 0x100F
	SOURCE_STATE              = 0x1010
	INITIAL                   = 0x1011
	PLAYING                   = 0x1012
	PAUSED                    = 0x1013
	STOPPED                   = 0x1014
	BUFFERS_QUEUED            = 0x1015
	BUFFERS_PROCESSED         = 0x1016
	SEC_OFFSET                = 0x1024
	SAMPLE_OFFSET             = 0x1025
	BYTE_OFFSET               = 0x1026
	SOURCE_TYPE               = 0x1027
	STATIC                    = 0x1028
	STREAMING                 = 0x1029
	UNDETERMINED              = 0x1030
	FORMAT_MONO8              = 0x1100
	FORMAT_MONO16             = 0x1101
	FORMAT_STEREO8            = 0x1102
	FORMAT_STEREO16           = 0x1103
	REFERENCE_DISTANCE        = 0x1020
	ROLLOFF_FACTOR            = 0x1021
	CONE_OUTER_GAIN           = 0x1022
	MAX_DISTANCE              = 0x1023
	FREQUENCY                 = 0x2001
	BITS                      = 0x2002
	CHANNELS                  = 0x2003
	SIZE                      = 0x2004
	UNUSED                    = 0x2010
	PENDING                   = 0x2011
	PROCESSED                 = 0x2012
	INVALID_NAME              = 0xA001
	INVALID_ENUM              = 0xA002
	INVALID_VALUE             = 0xA003
	INVALID_OPERATION         = 0xA004
	OUT_OF_MEMORY             = 0xA005
	VENDOR                    = 0xB001
	VERSION                   = 0xB002
	RENDERER                  = 0xB003
	EXTENSIONS                = 0xB004
	DOPPLER_FACTOR            = 0xC000
	DOPPLER_VELOCITY          = 0xC001
	SPEED_OF_SOUND            = 0xC003
	DISTANCE_MODEL            = 0xD000
	INVERSE_DISTANCE          = 0xD001
	INVERSE_DISTANCE_CLAMPED  = 0xD002
	LINEAR_DISTANCE           = 0xD003
	LINEAR_DISTANCE_CLAMPED   = 0xD004
	EXPONENT_DISTANCE         = 0xD005
	EXPONENT_DISTANCE_CLAMPED = 0xD006
)

// AL/alc.h constants (with AL removed). Refer to the original header for constant documentation.
const (
	C_FALSE                            = 0
	C_TRUE                             = 1
	C_NO_ERROR                         = 0
	C_FREQUENCY                        = 0x1007
	C_REFRESH                          = 0x1008
	C_SYNC                             = 0x1009
	C_MONO_SOURCES                     = 0x1010
	C_STEREO_SOURCES                   = 0x1011
	C_INVALID_DEVICE                   = 0xA001
	C_INVALID_CONTEXT                  = 0xA002
	C_INVALID_ENUM                     = 0xA003
	C_INVALID_VALUE                    = 0xA004
	C_OUT_OF_MEMORY                    = 0xA005
	C_DEFAULT_DEVICE_SPECIFIER         = 0x1004
	C_DEVICE_SPECIFIER                 = 0x1005
	C_EXTENSIONS                       = 0x1006
	C_MAJOR_VERSION                    = 0x1000
	C_MINOR_VERSION                    = 0x1001
	C_ATTRIBUTES_SIZE                  = 0x1002
	C_ALL_ATTRIBUTES                   = 0x1003
	C_CAPTURE_DEVICE_SPECIFIER         = 0x310
	C_CAPTURE_DEFAULT_DEVICE_SPECIFIER = 0x311
	C_CAPTURE_SAMPLES                  = 0x312
)

// bind the methods to the function pointers
func Init() error {
	C.al_init()
	return nil
}

// convert a uint boolean to a go bool
func cbool(albool uint) bool {
	return albool == TRUE
}

// Special type mappings. Note that the context and device are pointers
// on Windows and Linux, but integers on OSX.
type (
	Context uintptr // C.struct_ALCcontext_struct
	Device  uintptr // C.struct_ALCdevice_struct
	Pointer unsafe.Pointer
)

// AL/al.h go bindings
func Enable(capability int32)               { C.wrap_alEnable(C.int(capability)) }
func Disable(capability int32)              { C.wrap_alDisable(C.int(capability)) }
func IsEnabled(capability int32) bool       { return cbool(uint(C.wrap_alIsEnabled(C.int(capability)))) }
func GetString(param int32) string          { return C.GoString(C.wrap_alGetString(C.int(param))) }
func GetBooleanv(param int32, data *int8)   { C.wrap_alGetBooleanv(C.int(param), (*C.char)(data)) }
func GetIntegerv(param int32, data *int32)  { C.wrap_alGetIntegerv(C.int(param), (*C.int)(data)) }
func GetFloatv(param int32, data *float32)  { C.wrap_alGetFloatv(C.int(param), (*C.float)(data)) }
func GetDoublev(param int32, data *float64) { C.wrap_alGetDoublev(C.int(param), (*C.double)(data)) }
func GetBoolean(param int32) bool           { return cbool(uint(C.wrap_alGetBoolean(C.int(param)))) }
func GetInteger(param int32) int32          { return int32(C.wrap_alGetInteger(C.int(param))) }
func GetFloat(param int32) float32          { return float32(C.wrap_alGetFloat(C.int(param))) }
func GetDouble(param int32) float64         { return float64(C.wrap_alGetDouble(C.int(param))) }
func GetError() int32                       { return int32(C.wrap_alGetError()) }
func IsExtensionPresent(extname string) bool {
	cstr := C.CString(extname)
	defer C.free(unsafe.Pointer(cstr))
	return cbool(uint(C.wrap_alIsExtensionPresent(cstr)))
}
func GetProcAddress(fname string) Pointer {
	cstr := C.CString(fname)
	defer C.free(unsafe.Pointer(cstr))
	return Pointer(C.wrap_alGetProcAddress(cstr))
}
func GetEnumValue(ename string) int32 {
	cstr := C.CString(ename)
	defer C.free(unsafe.Pointer(cstr))
	return int32(C.wrap_alGetEnumValue(cstr))
}
func Listenerf(param int32, value float32) { C.wrap_alListenerf(C.int(param), C.float(value)) }
func Listener3f(param int32, value1, value2, value3 float32) {
	C.wrap_alListener3f(C.int(param), C.float(value1), C.float(value2), C.float(value3))
}
func Listenerfv(param int32, values *float32) { C.wrap_alListenerfv(C.int(param), (*C.float)(values)) }
func Listeneri(param int32, value int32)      { C.wrap_alListeneri(C.int(param), C.int(value)) }
func Listener3i(param int32, value1, value2, value3 int32) {
	C.wrap_alListener3i(C.int(param), C.int(value1), C.int(value2), C.int(value3))
}
func Listeneriv(param int32, values *int32) { C.wrap_alListeneriv(C.int(param), (*C.int)(values)) }
func GetListenerf(param int32, value *float32) {
	C.wrap_alGetListenerf(C.int(param), (*C.float)(value))
}
func GetListener3f(param int32, value1, value2, value3 *float32) {
	C.wrap_alGetListener3f(C.int(param), (*C.float)(value1), (*C.float)(value2), (*C.float)(value3))
}
func GetListenerfv(param int32, values *float32) {
	C.wrap_alGetListenerfv(C.int(param), (*C.float)(values))
}
func GetListeneri(param int32, value *int32) { C.wrap_alGetListeneri(C.int(param), (*C.int)(value)) }
func GetListener3i(param int32, value1, value2, value3 *int32) {
	C.wrap_alGetListener3i(C.int(param), (*C.int)(value1), (*C.int)(value2), (*C.int)(value3))
}
func GetListeneriv(param int32, values *int32) {
	C.wrap_alGetListeneriv(C.int(param), (*C.int)(values))
}
func GenSources(n int32, sources *uint32)    { C.wrap_alGenSources(C.int(n), (*C.uint)(sources)) }
func DeleteSources(n int32, sources *uint32) { C.wrap_alDeleteSources(C.int(n), (*C.uint)(sources)) }
func IsSource(sid uint32) bool               { return cbool(uint(C.wrap_alIsSource(C.uint(sid)))) }
func Sourcef(sid uint32, param int32, value float32) {
	C.wrap_alSourcef(C.uint(sid), C.int(param), C.float(value))
}
func Source3f(sid uint32, param int32, value1, value2, value3 float32) {
	C.wrap_alSource3f(C.uint(sid), C.int(param), C.float(value1), C.float(value2), C.float(value3))
}
func Sourcefv(sid uint32, param int32, values *float32) {
	C.wrap_alSourcefv(C.uint(sid), C.int(param), (*C.float)(values))
}
func Sourcei(sid uint32, param int32, value int32) {
	C.wrap_alSourcei(C.uint(sid), C.int(param), C.int(value))
}
func Source3i(sid uint32, param int32, value1, value2, value3 int32) {
	C.wrap_alSource3i(C.uint(sid), C.int(param), C.int(value1), C.int(value2), C.int(value3))
}
func Sourceiv(sid uint32, param int32, values *int32) {
	C.wrap_alSourceiv(C.uint(sid), C.int(param), (*C.int)(values))
}
func GetSourcef(sid uint32, param int32, value *float32) {
	C.wrap_alGetSourcef(C.uint(sid), C.int(param), (*C.float)(value))
}
func GetSource3f(sid uint32, param int32, value1, value2, value3 *float32) {
	C.wrap_alGetSource3f(C.uint(sid), C.int(param), (*C.float)(value1), (*C.float)(value2), (*C.float)(value3))
}
func GetSourcefv(sid uint32, param int32, values *float32) {
	C.wrap_alGetSourcefv(C.uint(sid), C.int(param), (*C.float)(values))
}
func GetSourcei(sid uint32, param int32, value *int32) {
	C.wrap_alGetSourcei(C.uint(sid), C.int(param), (*C.int)(value))
}
func GetSource3i(sid uint32, param int32, value1, value2, value3 *int32) {
	C.wrap_alGetSource3i(C.uint(sid), C.int(param), (*C.int)(value1), (*C.int)(value2), (*C.int)(value3))
}
func GetSourceiv(sid uint32, param int32, values *int32) {
	C.wrap_alGetSourceiv(C.uint(sid), C.int(param), (*C.int)(values))
}
func SourcePlayv(ns int32, sids *uint32)   { C.wrap_alSourcePlayv(C.int(ns), (*C.uint)(sids)) }
func SourceStopv(ns int32, sids *uint32)   { C.wrap_alSourceStopv(C.int(ns), (*C.uint)(sids)) }
func SourceRewindv(ns int32, sids *uint32) { C.wrap_alSourceRewindv(C.int(ns), (*C.uint)(sids)) }
func SourcePausev(ns int32, sids *uint32)  { C.wrap_alSourcePausev(C.int(ns), (*C.uint)(sids)) }
func SourcePlay(sid uint32)                { C.wrap_alSourcePlay(C.uint(sid)) }
func SourceStop(sid uint32)                { C.wrap_alSourceStop(C.uint(sid)) }
func SourceRewind(sid uint32)              { C.wrap_alSourceRewind(C.uint(sid)) }
func SourcePause(sid uint32)               { C.wrap_alSourcePause(C.uint(sid)) }
func SourceQueueBuffers(sid uint32, numEntries int32, bids *uint32) {
	C.wrap_alSourceQueueBuffers(C.uint(sid), C.int(numEntries), (*C.uint)(bids))
}
func SourceUnqueueBuffers(sid uint32, numEntries int32, bids *uint32) {
	C.wrap_alSourceUnqueueBuffers(C.uint(sid), C.int(numEntries), (*C.uint)(bids))
}
func GenBuffers(n int32, buffers *uint32)    { C.wrap_alGenBuffers(C.int(n), (*C.uint)(buffers)) }
func DeleteBuffers(n int32, buffers *uint32) { C.wrap_alDeleteBuffers(C.int(n), (*C.uint)(buffers)) }
func IsBuffer(bid uint32) bool               { return cbool(uint(C.wrap_alIsBuffer(C.uint(bid)))) }
func BufferData(bid uint32, format int32, data Pointer, size int32, freq int32) {
	C.wrap_alBufferData(C.uint(bid), C.int(format), unsafe.Pointer(data), C.int(size), C.int(freq))
}
func Bufferf(bid uint32, param int32, value float32) {
	C.wrap_alBufferf(C.uint(bid), C.int(param), C.float(value))
}
func Buffer3f(bid uint32, param int32, value1, value2, value3 float32) {
	C.wrap_alBuffer3f(C.uint(bid), C.int(param), C.float(value1), C.float(value2), C.float(value3))
}
func Bufferfv(bid uint32, param int32, values *float32) {
	C.wrap_alBufferfv(C.uint(bid), C.int(param), (*C.float)(values))
}
func Bufferi(bid uint32, param int32, value int32) {
	C.wrap_alBufferi(C.uint(bid), C.int(param), C.int(value))
}
func Buffer3i(bid uint32, param int32, value1, value2, value3 int32) {
	C.wrap_alBuffer3i(C.uint(bid), C.int(param), C.int(value1), C.int(value2), C.int(value3))
}
func Bufferiv(bid uint32, param int32, values *int32) {
	C.wrap_alBufferiv(C.uint(bid), C.int(param), (*C.int)(values))
}
func GetBufferf(bid uint32, param int32, value *float32) {
	C.wrap_alGetBufferf(C.uint(bid), C.int(param), (*C.float)(value))
}
func GetBuffer3f(bid uint32, param int32, value1, value2, value3 *float32) {
	C.wrap_alGetBuffer3f(C.uint(bid), C.int(param), (*C.float)(value1), (*C.float)(value2), (*C.float)(value3))
}
func GetBufferfv(bid uint32, param int32, values *float32) {
	C.wrap_alGetBufferfv(C.uint(bid), C.int(param), (*C.float)(values))
}
func GetBufferi(bid uint32, param int32, value *int32) {
	C.wrap_alGetBufferi(C.uint(bid), C.int(param), (*C.int)(value))
}
func GetBuffer3i(bid uint32, param int32, value1, value2, value3 *int32) {
	C.wrap_alGetBuffer3i(C.uint(bid), C.int(param), (*C.int)(value1), (*C.int)(value2), (*C.int)(value3))
}
func GetBufferiv(bid uint32, param int32, values *int32) {
	C.wrap_alGetBufferiv(C.uint(bid), C.int(param), (*C.int)(values))
}
func DopplerFactor(value float32)         { C.wrap_alDopplerFactor(C.float(value)) }
func DopplerVelocity(value float32)       { C.wrap_alDopplerVelocity(C.float(value)) }
func SpeedOfSound(value float32)          { C.wrap_alSpeedOfSound(C.float(value)) }
func DistanceModel(distanceModel float32) { C.wrap_alDistanceModel(C.int(distanceModel)) }

// AL/alc.h go bindings
func CreateContext(device Device, attrlist *int32) Context {
	return (Context)(C.wrap_alcCreateContext((C.uintptr_t)(device), (*C.int)(attrlist)))
}
func MakeContextCurrent(context Context) bool {
	return cbool(uint(C.wrap_alcMakeContextCurrent((C.uintptr_t)(context))))
}
func ProcessContext(context Context) {
	C.wrap_alcProcessContext((C.uintptr_t)(context))
}
func SuspendContext(context Context) {
	C.wrap_alcSuspendContext((C.uintptr_t)(context))
}
func DestroyContext(context Context) {
	C.wrap_alcDestroyContext((C.uintptr_t)(context))
}
func GetCurrentContext() Context { return (Context)(C.wrap_alcGetCurrentContext()) }
func GetContextsDevice(context Context) Device {
	return (Device)(C.wrap_alcGetContextsDevice((C.uintptr_t)(context)))
}
func OpenDevice(devicename string) Device {
	if devicename != "" {
		return (Device)(C.wrap_alcOpenDevice(nil))
	}
	cstr := C.CString(devicename)
	defer C.free(unsafe.Pointer(cstr))
	return (Device)(C.wrap_alcOpenDevice(cstr))
}
func CloseDevice(device Device) bool {
	return cbool(uint(C.wrap_alcCloseDevice((C.uintptr_t)(device))))
}
func GetDeviceError(device Device) int32 {
	return int32(C.wrap_alcGetError((C.uintptr_t)(device)))
}
func IsDeviceExtensionPresent(device Device, extname string) bool {
	cstr := C.CString(extname)
	defer C.free(unsafe.Pointer(cstr))
	return cbool(uint(C.wrap_alcIsExtensionPresent((C.uintptr_t)(device), cstr)))
}
func GetDeviceProcAddress(device Device, fname string) Pointer {
	cstr := C.CString(fname)
	defer C.free(unsafe.Pointer(cstr))
	return Pointer(C.wrap_alcGetProcAddress((C.uintptr_t)(device), cstr))
}
func GetDeviceEnumValue(device Device, ename string) int32 {
	cstr := C.CString(ename)
	defer C.free(unsafe.Pointer(cstr))
	return int32(C.wrap_alcGetEnumValue((C.uintptr_t)(device), cstr))
}
func GetDeviceString(device Device, param int32) string {
	return C.GoString(C.wrap_alcGetString((C.uintptr_t)(device), C.int(param)))
}
func GetDeviceIntegerv(device Device, param int32, size int32, data *int32) {
	C.wrap_alcGetIntegerv((C.uintptr_t)(device), C.int(param), C.int(size), (*C.int)(data))
}
func CaptureOpenDevice(devicename string, frequency uint32, format int32, buffersize int32) Device {
	cstr := C.CString(devicename)
	defer C.free(unsafe.Pointer(cstr))
	return (Device)(C.wrap_alcCaptureOpenDevice(cstr, C.uint(frequency), C.int(format), C.int(buffersize)))
}
func CaptureCloseDevice(device Device) bool {
	return cbool(uint(C.wrap_alcCaptureCloseDevice((C.uintptr_t)(device))))
}
func CaptureStart(device Device) { C.wrap_alcCaptureStart((C.uintptr_t)(device)) }
func CaptureStop(device Device)  { C.wrap_alcCaptureStop((C.uintptr_t)(device)) }
func CaptureSamples(device Device, buffer Pointer, samples int) {
	C.wrap_alcCaptureSamples((C.uintptr_t)(device), unsafe.Pointer(buffer), C.int(samples))
}

// Show which function pointers are bound [+] or not bound [-].
// Expected to be used as a sanity check to see if the OpenAL libraries exist.
func BindingReport() (report []string) {
	report = []string{}

	// AL/al.h
	report = append(report, "AL")
	report = append(report, isBound(unsafe.Pointer(C.pfn_alEnable), "alEnable"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alDisable), "alDisable"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alIsEnabled), "alIsEnabled"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetString), "alGetString"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetBooleanv), "alGetBooleanv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetIntegerv), "alGetIntegerv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetFloatv), "alGetFloatv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetDoublev), "alGetDoublev"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetBoolean), "alGetBoolean"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetInteger), "alGetInteger"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetFloat), "alGetFloat"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetDouble), "alGetDouble"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetError), "alGetError"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alIsExtensionPresent), "alIsExtensionPresent"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetProcAddress), "alGetProcAddress"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetEnumValue), "alGetEnumValue"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alListenerf), "alListenerf"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alListener3f), "alListener3f"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alListenerfv), "alListenerfv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alListeneri), "alListeneri"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alListener3i), "alListener3i"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alListeneriv), "alListeneriv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetListenerf), "alGetListenerf"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetListener3f), "alGetListener3f"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetListenerfv), "alGetListenerfv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetListeneri), "alGetListeneri"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetListener3i), "alGetListener3i"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetListeneriv), "alGetListeneriv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGenSources), "alGenSources"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alDeleteSources), "alDeleteSources"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alIsSource), "alIsSource"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alSourcef), "alSourcef"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alSource3f), "alSource3f"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alSourcefv), "alSourcefv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alSourcei), "alSourcei"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alSource3i), "alSource3i"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alSourceiv), "alSourceiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetSourcef), "alGetSourcef"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetSource3f), "alGetSource3f"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetSourcefv), "alGetSourcefv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetSourcei), "alGetSourcei"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetSource3i), "alGetSource3i"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetSourceiv), "alGetSourceiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alSourcePlayv), "alSourcePlayv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alSourceStopv), "alSourceStopv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alSourceRewindv), "alSourceRewindv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alSourcePausev), "alSourcePausev"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alSourcePlay), "alSourcePlay"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alSourceStop), "alSourceStop"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alSourceRewind), "alSourceRewind"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alSourcePause), "alSourcePause"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alSourceQueueBuffers), "alSourceQueueBuffers"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alSourceUnqueueBuffers), "alSourceUnqueueBuffers"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGenBuffers), "alGenBuffers"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alDeleteBuffers), "alDeleteBuffers"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alIsBuffer), "alIsBuffer"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alBufferData), "alBufferData"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alBufferf), "alBufferf"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alBuffer3f), "alBuffer3f"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alBufferfv), "alBufferfv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alBufferi), "alBufferi"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alBuffer3i), "alBuffer3i"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alBufferiv), "alBufferiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetBufferf), "alGetBufferf"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetBuffer3f), "alGetBuffer3f"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetBufferfv), "alGetBufferfv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetBufferi), "alGetBufferi"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetBuffer3i), "alGetBuffer3i"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetBufferiv), "alGetBufferiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alDopplerFactor), "alDopplerFactor"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alDopplerVelocity), "alDopplerVelocity"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alSpeedOfSound), "alSpeedOfSound"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alDistanceModel), "alDistanceModel"))

	// AL/alc.h
	report = append(report, "ALC")
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcCreateContext), "alcCreateContext"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcMakeContextCurrent), "alcMakeContextCurrent"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcProcessContext), "alcProcessContext"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcSuspendContext), "alcSuspendContext"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcDestroyContext), "alcDestroyContext"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcGetCurrentContext), "alcGetCurrentContext"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcGetContextsDevice), "alcGetContextsDevice"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcOpenDevice), "alcOpenDevice"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcCloseDevice), "alcCloseDevice"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcGetError), "alcGetError"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcIsExtensionPresent), "alcIsExtensionPresent"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcGetProcAddress), "alcGetProcAddress"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcGetEnumValue), "alcGetEnumValue"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcGetString), "alcGetString"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcGetIntegerv), "alcGetIntegerv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcCaptureOpenDevice), "alcCaptureOpenDevice"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcCaptureCloseDevice), "alcCaptureCloseDevice"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcCaptureStart), "alcCaptureStart"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcCaptureStop), "alcCaptureStop"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcCaptureSamples), "alcCaptureSamples"))
	return
}

// isBound returns a string that indicates if the given function
// pointer is bound.
func isBound(pfn unsafe.Pointer, fn string) string {
	inc := " "
	if pfn != nil {
		inc = "+"
	}
	return fmt.Sprintf("   [%s] %s", inc, fn)
}
