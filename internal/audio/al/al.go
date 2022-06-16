// Copyright © 2013-2024 Galvanized Logic Inc.

//go:build windows

// Package al provides golang audio library bindings for OpenAL.
// Refer to the official OpenAL documentation for more information.
//
// Package al is provided as part of the vu (virtual universe) 3D engine.
package al

// OpenAL: https://openal.org
// Requires the 64-bit soft_oal.dll from:
// o https://openal-soft.org/openal-binaries/

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/gazed/vu/internal/device/win"
	"golang.org/x/sys/windows"
)

var (
	libopenal32 *windows.LazyDLL // Library

	// Functions AL/al.h
	alEnable               *windows.LazyProc
	alDisable              *windows.LazyProc
	alIsEnabled            *windows.LazyProc
	alGetString            *windows.LazyProc
	alGetBooleanv          *windows.LazyProc
	alGetIntegerv          *windows.LazyProc
	alGetFloatv            *windows.LazyProc
	alGetDoublev           *windows.LazyProc
	alGetBoolean           *windows.LazyProc
	alGetInteger           *windows.LazyProc
	alGetFloat             *windows.LazyProc
	alGetDouble            *windows.LazyProc
	alGetError             *windows.LazyProc
	alIsExtensionPresent   *windows.LazyProc
	alGetProcAddress       *windows.LazyProc
	alGetEnumValue         *windows.LazyProc
	alListenerf            *windows.LazyProc
	alListener3f           *windows.LazyProc
	alListenerfv           *windows.LazyProc
	alListeneri            *windows.LazyProc
	alListener3i           *windows.LazyProc
	alListeneriv           *windows.LazyProc
	alGetListenerf         *windows.LazyProc
	alGetListener3f        *windows.LazyProc
	alGetListenerfv        *windows.LazyProc
	alGetListeneri         *windows.LazyProc
	alGetListener3i        *windows.LazyProc
	alGetListeneriv        *windows.LazyProc
	alGenSources           *windows.LazyProc
	alDeleteSources        *windows.LazyProc
	alIsSource             *windows.LazyProc
	alSourcef              *windows.LazyProc
	alSource3f             *windows.LazyProc
	alSourcefv             *windows.LazyProc
	alSourcei              *windows.LazyProc
	alSource3i             *windows.LazyProc
	alSourceiv             *windows.LazyProc
	alGetSourcef           *windows.LazyProc
	alGetSource3f          *windows.LazyProc
	alGetSourcefv          *windows.LazyProc
	alGetSourcei           *windows.LazyProc
	alGetSource3i          *windows.LazyProc
	alGetSourceiv          *windows.LazyProc
	alSourcePlayv          *windows.LazyProc
	alSourceStopv          *windows.LazyProc
	alSourceRewindv        *windows.LazyProc
	alSourcePausev         *windows.LazyProc
	alSourcePlay           *windows.LazyProc
	alSourceStop           *windows.LazyProc
	alSourceRewind         *windows.LazyProc
	alSourcePause          *windows.LazyProc
	alSourceQueueBuffers   *windows.LazyProc
	alSourceUnqueueBuffers *windows.LazyProc
	alGenBuffers           *windows.LazyProc
	alDeleteBuffers        *windows.LazyProc
	alIsBuffer             *windows.LazyProc
	alBufferData           *windows.LazyProc
	alBufferf              *windows.LazyProc
	alBuffer3f             *windows.LazyProc
	alBufferfv             *windows.LazyProc
	alBufferi              *windows.LazyProc
	alBuffer3i             *windows.LazyProc
	alBufferiv             *windows.LazyProc
	alGetBufferf           *windows.LazyProc
	alGetBuffer3f          *windows.LazyProc
	alGetBufferfv          *windows.LazyProc
	alGetBufferi           *windows.LazyProc
	alGetBuffer3i          *windows.LazyProc
	alGetBufferiv          *windows.LazyProc
	alDopplerFactor        *windows.LazyProc
	alDopplerVelocity      *windows.LazyProc
	alSpeedOfSound         *windows.LazyProc
	alDistanceModel        *windows.LazyProc

	// Functions AL/ac.h
	alcCreateContext      *windows.LazyProc
	alcMakeContextCurrent *windows.LazyProc
	alcProcessContext     *windows.LazyProc
	alcSuspendContext     *windows.LazyProc
	alcDestroyContext     *windows.LazyProc
	alcGetCurrentContext  *windows.LazyProc
	alcGetContextsDevice  *windows.LazyProc
	alcOpenDevice         *windows.LazyProc
	alcCloseDevice        *windows.LazyProc
	alcGetError           *windows.LazyProc
	alcIsExtensionPresent *windows.LazyProc
	alcGetProcAddress     *windows.LazyProc
	alcGetEnumValue       *windows.LazyProc
	alcGetString          *windows.LazyProc
	alcGetIntegerv        *windows.LazyProc
	alcCaptureOpenDevice  *windows.LazyProc
	alcCaptureCloseDevice *windows.LazyProc
	alcCaptureStart       *windows.LazyProc
	alcCaptureStop        *windows.LazyProc
	alcCaptureSamples     *windows.LazyProc
)

// bind the methods to the function pointers
func Init() {
	// Library
	// libopenal32 = windows.NewLazyDLL("OpenAL32.dll")
	libopenal32 = windows.NewLazyDLL("soft_oal.dll")

	// Functions AL/al.h
	alEnable = libopenal32.NewProc("alEnable")
	alDisable = libopenal32.NewProc("alDisable")
	alIsEnabled = libopenal32.NewProc("alIsEnabled")
	alGetString = libopenal32.NewProc("alGetString")
	alGetBooleanv = libopenal32.NewProc("alGetBooleanv")
	alGetIntegerv = libopenal32.NewProc("alGetIntegerv")
	alGetFloatv = libopenal32.NewProc("alGetFloatv")
	alGetDoublev = libopenal32.NewProc("alGetDoublev")
	alGetBoolean = libopenal32.NewProc("alGetBoolean")
	alGetInteger = libopenal32.NewProc("alGetInteger")
	alGetFloat = libopenal32.NewProc("alGetFloat")
	alGetDouble = libopenal32.NewProc("alGetDouble")
	alGetError = libopenal32.NewProc("alGetError")
	alIsExtensionPresent = libopenal32.NewProc("alIsExtensionPresent")
	alGetProcAddress = libopenal32.NewProc("alGetProcAddress")
	alGetEnumValue = libopenal32.NewProc("alGetEnumValue")
	alListenerf = libopenal32.NewProc("alListenerf")
	alListener3f = libopenal32.NewProc("alListener3f")
	alListenerfv = libopenal32.NewProc("alListenerfv")
	alListeneri = libopenal32.NewProc("alListeneri")
	alListener3i = libopenal32.NewProc("alListener3i")
	alListeneriv = libopenal32.NewProc("alListeneriv")
	alGetListenerf = libopenal32.NewProc("alGetListenerf")
	alGetListener3f = libopenal32.NewProc("alGetListener3f")
	alGetListenerfv = libopenal32.NewProc("alGetListenerfv")
	alGetListeneri = libopenal32.NewProc("alGetListeneri")
	alGetListener3i = libopenal32.NewProc("alGetListener3i")
	alGetListeneriv = libopenal32.NewProc("alGetListeneriv")
	alGenSources = libopenal32.NewProc("alGenSources")
	alDeleteSources = libopenal32.NewProc("alDeleteSources")
	alIsSource = libopenal32.NewProc("alIsSource")
	alSourcef = libopenal32.NewProc("alSourcef")
	alSource3f = libopenal32.NewProc("alSource3f")
	alSourcefv = libopenal32.NewProc("alSourcefv")
	alSourcei = libopenal32.NewProc("alSourcei")
	alSource3i = libopenal32.NewProc("alSource3i")
	alSourceiv = libopenal32.NewProc("alSourceiv")
	alGetSourcef = libopenal32.NewProc("alGetSourcef")
	alGetSource3f = libopenal32.NewProc("alGetSource3f")
	alGetSourcefv = libopenal32.NewProc("alGetSourcefv")
	alGetSourcei = libopenal32.NewProc("alGetSourcei")
	alGetSource3i = libopenal32.NewProc("alGetSource3i")
	alGetSourceiv = libopenal32.NewProc("alGetSourceiv")
	alSourcePlayv = libopenal32.NewProc("alSourcePlayv")
	alSourceStopv = libopenal32.NewProc("alSourceStopv")
	alSourceRewindv = libopenal32.NewProc("alSourceRewindv")
	alSourcePausev = libopenal32.NewProc("alSourcePausev")
	alSourcePlay = libopenal32.NewProc("alSourcePlay")
	alSourceStop = libopenal32.NewProc("alSourceStop")
	alSourceRewind = libopenal32.NewProc("alSourceRewind")
	alSourcePause = libopenal32.NewProc("alSourcePause")
	alSourceQueueBuffers = libopenal32.NewProc("alSourceQueueBuffers")
	alSourceUnqueueBuffers = libopenal32.NewProc("alSourceUnqueueBuffers")
	alGenBuffers = libopenal32.NewProc("alGenBuffers")
	alDeleteBuffers = libopenal32.NewProc("alDeleteBuffers")
	alIsBuffer = libopenal32.NewProc("alIsBuffer")
	alBufferData = libopenal32.NewProc("alBufferData")
	alBufferf = libopenal32.NewProc("alBufferf")
	alBuffer3f = libopenal32.NewProc("alBuffer3f")
	alBufferfv = libopenal32.NewProc("alBufferfv")
	alBufferi = libopenal32.NewProc("alBufferi")
	alBuffer3i = libopenal32.NewProc("alBuffer3i")
	alBufferiv = libopenal32.NewProc("alBufferiv")
	alGetBufferf = libopenal32.NewProc("alGetBufferf")
	alGetBuffer3f = libopenal32.NewProc("alGetBuffer3f")
	alGetBufferfv = libopenal32.NewProc("alGetBufferfv")
	alGetBufferi = libopenal32.NewProc("alGetBufferi")
	alGetBuffer3i = libopenal32.NewProc("alGetBuffer3i")
	alGetBufferiv = libopenal32.NewProc("alGetBufferiv")
	alDopplerFactor = libopenal32.NewProc("alDopplerFactor")
	alDopplerVelocity = libopenal32.NewProc("alDopplerVelocity")
	alSpeedOfSound = libopenal32.NewProc("alSpeedOfSound")
	alDistanceModel = libopenal32.NewProc("alDistanceModel")

	// AL/alc.h
	alcCreateContext = libopenal32.NewProc("alcCreateContext")
	alcMakeContextCurrent = libopenal32.NewProc("alcMakeContextCurrent")
	alcProcessContext = libopenal32.NewProc("alcProcessContext")
	alcSuspendContext = libopenal32.NewProc("alcSuspendContext")
	alcDestroyContext = libopenal32.NewProc("alcDestroyContext")
	alcGetCurrentContext = libopenal32.NewProc("alcGetCurrentContext")
	alcGetContextsDevice = libopenal32.NewProc("alcGetContextsDevice")
	alcOpenDevice = libopenal32.NewProc("alcOpenDevice")
	alcCloseDevice = libopenal32.NewProc("alcCloseDevice")
	alcGetError = libopenal32.NewProc("alcGetError")
	alcIsExtensionPresent = libopenal32.NewProc("alcIsExtensionPresent")
	alcGetProcAddress = libopenal32.NewProc("alcGetProcAddress")
	alcGetEnumValue = libopenal32.NewProc("alcGetEnumValue")
	alcGetString = libopenal32.NewProc("alcGetString")
	alcGetIntegerv = libopenal32.NewProc("alcGetIntegerv")
	alcCaptureOpenDevice = libopenal32.NewProc("alcCaptureOpenDevice")
	alcCaptureCloseDevice = libopenal32.NewProc("alcCaptureCloseDevice")
	alcCaptureStart = libopenal32.NewProc("alcCaptureStart")
	alcCaptureStop = libopenal32.NewProc("alcCaptureStop")
	alcCaptureSamples = libopenal32.NewProc("alcCaptureSamples ")
}

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

func UTF16PtrToString(s *uint16) string {
	return windows.UTF16PtrToString(s)
}

// convert a uint boolean to a go bool
func cbool(albool uint) bool { return albool == TRUE }

// Special type mappings. Note that the context and device are pointers
// on Windows and Linux, but integers on OSX.
type (
	Context uintptr // C.struct_ALCcontext_struct
	Device  uintptr // C.struct_ALCdevice_struct
	Pointer unsafe.Pointer
)

// AL/al.h go bindings
func Enable(capability int32) {
	syscall.Syscall(alEnable.Addr(), 1,
		uintptr(capability),
		0,
		0)
}
func Disable(capability int32) {
	syscall.Syscall(alDisable.Addr(), 1,
		uintptr(capability),
		0,
		0)
}

func IsEnabled(capability int32) bool {
	ret, _, _ := syscall.Syscall(alIsEnabled.Addr(), 1,
		uintptr(capability),
		0,
		0)
	return ret == TRUE
}

func GetString(param int32) string {
	ret, _, _ := syscall.Syscall(alGetString.Addr(), 1,
		uintptr(param),
		0,
		0)
	return win.UTF16PtrToString((*uint16)(unsafe.Pointer(ret)))
}
func GetBooleanv(param int32, data *int8) {
	syscall.Syscall(alGetBooleanv.Addr(), 2,
		uintptr(param),
		uintptr(unsafe.Pointer(data)),
		0)
}
func GetIntegerv(param int32, data *int32) {
	syscall.Syscall(alGetIntegerv.Addr(), 2,
		uintptr(param),
		uintptr(unsafe.Pointer(data)),
		0)
}
func GetFloatv(param int32, data *float32) {
	syscall.Syscall(alGetFloatv.Addr(), 2,
		uintptr(param),
		uintptr(unsafe.Pointer(data)),
		0)
}
func GetDoublev(param int32, data *float64) {
	syscall.Syscall(alGetDoublev.Addr(), 2,
		uintptr(param),
		uintptr(unsafe.Pointer(data)),
		0)
}
func GetBoolean(param int32) bool {
	ret, _, _ := syscall.Syscall(alGetBoolean.Addr(), 1,
		uintptr(param),
		0,
		0)
	return ret == TRUE
}
func GetInteger(param int32) int32 {
	ret, _, _ := syscall.Syscall(alGetInteger.Addr(), 1,
		uintptr(param),
		0,
		0)
	return int32(ret)
}

func GetFloat(param int32) float32 {
	ret, _, _ := syscall.Syscall(alGetFloat.Addr(), 1,
		uintptr(param),
		0,
		0)
	return float32(ret)
}
func GetDouble(param int32) float64 {
	ret, _, _ := syscall.Syscall(alGetDouble.Addr(), 1,
		uintptr(param),
		0,
		0)
	return float64(ret)
}
func GetError() int32 {
	ret, _, _ := syscall.Syscall(alGetError.Addr(), 0,
		0,
		0,
		0)
	return int32(ret)
}

func IsExtensionPresent(extname string) bool {
	str16, err := syscall.UTF16PtrFromString(extname)
	if err != nil {
		return false
	}
	ret, _, _ := syscall.Syscall(alIsExtensionPresent.Addr(), 1,
		uintptr(unsafe.Pointer(str16)),
		0,
		0)
	return ret == TRUE
}

func GetProcAddress(fname string) Pointer {
	str16, err := syscall.UTF16PtrFromString(fname)
	if err != nil {
		return nil
	}
	ret, _, _ := syscall.Syscall(alGetProcAddress.Addr(), 1,
		uintptr(unsafe.Pointer(str16)),
		0,
		0)
	return Pointer(ret)
}

func GetEnumValue(ename string) int32 {
	str16, err := syscall.UTF16PtrFromString(ename)
	if err != nil {
		return -1
	}
	ret, _, _ := syscall.Syscall(alGetEnumValue.Addr(), 1,
		uintptr(unsafe.Pointer(str16)),
		0,
		0)
	return int32(ret)
}

func Listenerf(param int32, value float32) {
	syscall.Syscall(alListenerf.Addr(), 2,
		uintptr(param),
		uintptr(value),
		0)
}
func Listener3f(param int32, value1, value2, value3 float32) {
	syscall.Syscall6(alListener3f.Addr(), 4,
		uintptr(param),
		uintptr(value1),
		uintptr(value2),
		uintptr(value3),
		0,
		0)
}
func Listenerfv(param int32, values *float32) {
	syscall.Syscall(alListenerfv.Addr(), 2,
		uintptr(param),
		uintptr(unsafe.Pointer(values)),
		0)
}
func Listeneri(param int32, value int32) {
	syscall.Syscall(alListeneri.Addr(), 2,
		uintptr(param),
		uintptr(value),
		0)
}
func Listener3i(param int32, value1, value2, value3 int32) {
	syscall.Syscall6(alListener3i.Addr(), 4,
		uintptr(param),
		uintptr(value1),
		uintptr(value2),
		uintptr(value3),
		0,
		0)
}
func Listeneriv(param int32, values *int32) {
	syscall.Syscall(alListeneriv.Addr(), 2,
		uintptr(param),
		uintptr(unsafe.Pointer(values)),
		0)
}
func GetListenerf(param int32, value *float32) {
	syscall.Syscall(alGetListenerf.Addr(), 2,
		uintptr(param),
		uintptr(unsafe.Pointer(value)),
		0)
}
func GetListener3f(param int32, value1, value2, value3 *float32) {
	syscall.Syscall6(alGetListener3f.Addr(), 4,
		uintptr(param),
		uintptr(unsafe.Pointer(value1)),
		uintptr(unsafe.Pointer(value2)),
		uintptr(unsafe.Pointer(value3)),
		0,
		0)
}
func GetListenerfv(param int32, values *float32) {
	syscall.Syscall(alGetListenerfv.Addr(), 2,
		uintptr(param),
		uintptr(unsafe.Pointer(values)),
		0)
}
func GetListeneri(param int32, value *int32) {
	syscall.Syscall(alGetListeneri.Addr(), 2,
		uintptr(param),
		uintptr(unsafe.Pointer(value)),
		0)
}
func GetListener3i(param int32, value1, value2, value3 *int32) {
	syscall.Syscall6(alGetListener3i.Addr(), 4,
		uintptr(param),
		uintptr(unsafe.Pointer(value1)),
		uintptr(unsafe.Pointer(value2)),
		uintptr(unsafe.Pointer(value3)),
		0,
		0)
}
func GetListeneriv(param int32, values *int32) {
	syscall.Syscall(alGetListeneriv.Addr(), 2,
		uintptr(param),
		uintptr(unsafe.Pointer(values)),
		0)
}
func GenSources(n int32, sources *uint32) {
	syscall.Syscall(alGenSources.Addr(), 2,
		uintptr(n),
		uintptr(unsafe.Pointer(sources)),
		0)
}
func DeleteSources(n int32, sources *uint32) {
	syscall.Syscall(alDeleteSources.Addr(), 2,
		uintptr(n),
		uintptr(unsafe.Pointer(sources)),
		0)
}
func IsSource(sid uint32) bool {
	ret, _, _ := syscall.Syscall(alIsSource.Addr(), 1,
		uintptr(sid),
		0,
		0)
	return ret == TRUE
}
func Sourcef(sid uint32, param int32, value float32) {
	syscall.Syscall(alSourcef.Addr(), 3,
		uintptr(sid),
		uintptr(param),
		uintptr(value))
}
func Source3f(sid uint32, param int32, value1, value2, value3 float32) {
	syscall.Syscall6(alSource3f.Addr(), 5,
		uintptr(sid),
		uintptr(param),
		uintptr(value1),
		uintptr(value2),
		uintptr(value3),
		0)
}
func Sourcefv(sid uint32, param int32, values *float32) {
	syscall.Syscall(alSourcefv.Addr(), 3,
		uintptr(sid),
		uintptr(param),
		uintptr(unsafe.Pointer(values)))
}
func Sourcei(sid uint32, param int32, value int32) {
	syscall.Syscall(alSourcei.Addr(), 3,
		uintptr(sid),
		uintptr(param),
		uintptr(value))
}
func Source3i(sid uint32, param int32, value1, value2, value3 int32) {
	syscall.Syscall6(alSource3i.Addr(), 5,
		uintptr(sid),
		uintptr(param),
		uintptr(value1),
		uintptr(value2),
		uintptr(value3),
		0)
}
func Sourceiv(sid uint32, param int32, values *int32) {
	syscall.Syscall(alSourceiv.Addr(), 3,
		uintptr(sid),
		uintptr(param),
		uintptr(unsafe.Pointer(values)))
}
func GetSourcef(sid uint32, param int32, value *float32) {
	syscall.Syscall(alGetSourcef.Addr(), 3,
		uintptr(sid),
		uintptr(param),
		uintptr(unsafe.Pointer(value)))
}
func GetSource3f(sid uint32, param int32, value1, value2, value3 *float32) {
	syscall.Syscall6(alGetSource3f.Addr(), 5,
		uintptr(sid),
		uintptr(param),
		uintptr(unsafe.Pointer(value1)),
		uintptr(unsafe.Pointer(value2)),
		uintptr(unsafe.Pointer(value3)),
		0)
}
func GetSourcefv(sid uint32, param int32, values *float32) {
	syscall.Syscall(alGetSourcefv.Addr(), 3,
		uintptr(sid),
		uintptr(param),
		uintptr(unsafe.Pointer(values)))
}
func GetSourcei(sid uint32, param int32, value *int32) {
	syscall.Syscall(alGetSourcei.Addr(), 3,
		uintptr(sid),
		uintptr(param),
		uintptr(unsafe.Pointer(value)))
}
func GetSource3i(sid uint32, param int32, value1, value2, value3 *int32) {
	syscall.Syscall6(alGetSource3i.Addr(), 5,
		uintptr(sid),
		uintptr(param),
		uintptr(unsafe.Pointer(value1)),
		uintptr(unsafe.Pointer(value2)),
		uintptr(unsafe.Pointer(value3)),
		0)
}
func GetSourceiv(sid uint32, param int32, values *int32) {
	syscall.Syscall(alGetSourceiv.Addr(), 3,
		uintptr(sid),
		uintptr(param),
		uintptr(unsafe.Pointer(values)))
}
func SourcePlayv(ns int32, sids *uint32) {
	syscall.Syscall(alSourcePlayv.Addr(), 2,
		uintptr(ns),
		uintptr(unsafe.Pointer(sids)),
		0)
}
func SourceStopv(ns int32, sids *uint32) {
	syscall.Syscall(alSourceStopv.Addr(), 2,
		uintptr(ns),
		uintptr(unsafe.Pointer(sids)),
		0)
}
func SourceRewindv(ns int32, sids *uint32) {
	syscall.Syscall(alSourceRewindv.Addr(), 2,
		uintptr(ns),
		uintptr(unsafe.Pointer(sids)),
		0)
}
func SourcePausev(ns int32, sids *uint32) {
	syscall.Syscall(alSourcePausev.Addr(), 2,
		uintptr(ns),
		uintptr(unsafe.Pointer(sids)),
		0)
}
func SourcePlay(sid uint32) {
	syscall.Syscall(alSourcePlay.Addr(), 1,
		uintptr(sid),
		0,
		0)
}
func SourceStop(sid uint32) {
	syscall.Syscall(alSourceStop.Addr(), 1,
		uintptr(sid),
		0,
		0)
}
func SourceRewind(sid uint32) {
	syscall.Syscall(alSourceRewind.Addr(), 1,
		uintptr(sid),
		0,
		0)
}
func SourcePause(sid uint32) {
	syscall.Syscall(alSourcePause.Addr(), 1,
		uintptr(sid),
		0,
		0)
}
func SourceQueueBuffers(sid uint32, numEntries int32, bids *uint32) {
	syscall.Syscall(alSourceQueueBuffers.Addr(), 3,
		uintptr(sid),
		uintptr(numEntries),
		uintptr(unsafe.Pointer(bids)))
}
func SourceUnqueueBuffers(sid uint32, numEntries int32, bids *uint32) {
	syscall.Syscall(alSourceUnqueueBuffers.Addr(), 3,
		uintptr(sid),
		uintptr(numEntries),
		uintptr(unsafe.Pointer(bids)))
}
func GenBuffers(n int32, buffers *uint32) {
	syscall.Syscall(alGenBuffers.Addr(), 2,
		uintptr(n),
		uintptr(unsafe.Pointer(buffers)),
		0)
}
func DeleteBuffers(n int32, buffers *uint32) {
	syscall.Syscall(alDeleteBuffers.Addr(), 2,
		uintptr(n),
		uintptr(unsafe.Pointer(buffers)),
		0)
}
func IsBuffer(bid uint32) bool {
	ret, _, _ := syscall.Syscall(alIsBuffer.Addr(), 1,
		uintptr(bid),
		0,
		0)
	return ret == TRUE
}
func BufferData(bid uint32, format int32, data Pointer, size int32, freq int32) {
	syscall.Syscall6(alBufferData.Addr(), 5,
		uintptr(bid),
		uintptr(format),
		uintptr(data),
		uintptr(size),
		uintptr(freq),
		0)
}
func Bufferf(bid uint32, param int32, value float32) {
	syscall.Syscall(alBufferf.Addr(), 3,
		uintptr(bid),
		uintptr(param),
		uintptr(value))
}
func Buffer3f(bid uint32, param int32, value1, value2, value3 float32) {
	syscall.Syscall6(alBuffer3f.Addr(), 1,
		uintptr(bid),
		uintptr(param),
		uintptr(value1),
		uintptr(value2),
		uintptr(value3),
		0)
}
func Bufferfv(bid uint32, param int32, values *float32) {
	syscall.Syscall(alBufferfv.Addr(), 3,
		uintptr(bid),
		uintptr(param),
		uintptr(unsafe.Pointer(values)))
}
func Bufferi(bid uint32, param int32, value int32) {
	syscall.Syscall(alBufferi.Addr(), 3,
		uintptr(bid),
		uintptr(param),
		uintptr(value))
}
func Buffer3i(bid uint32, param int32, value1, value2, value3 int32) {
	syscall.Syscall6(alBuffer3i.Addr(), 5,
		uintptr(bid),
		uintptr(param),
		uintptr(value1),
		uintptr(value2),
		uintptr(value3),
		0)
}
func Bufferiv(bid uint32, param int32, values *int32) {
	syscall.Syscall(alBufferiv.Addr(), 3,
		uintptr(bid),
		uintptr(param),
		uintptr(unsafe.Pointer(values)))
}
func GetBufferf(bid uint32, param int32, value *float32) {
	syscall.Syscall(alGetBufferf.Addr(), 3,
		uintptr(bid),
		uintptr(param),
		uintptr(unsafe.Pointer(value)))
}
func GetBuffer3f(bid uint32, param int32, value1, value2, value3 *float32) {
	syscall.Syscall6(alGetBuffer3f.Addr(), 5,
		uintptr(bid),
		uintptr(param),
		uintptr(unsafe.Pointer(value1)),
		uintptr(unsafe.Pointer(value2)),
		uintptr(unsafe.Pointer(value3)),
		0)
}
func GetBufferfv(bid uint32, param int32, values *float32) {
	syscall.Syscall(alGetBufferfv.Addr(), 3,
		uintptr(bid),
		uintptr(param),
		uintptr(unsafe.Pointer(values)))
}
func GetBufferi(bid uint32, param int32, value *int32) {
	syscall.Syscall(alGetBufferi.Addr(), 3,
		uintptr(bid),
		uintptr(param),
		uintptr(unsafe.Pointer(value)))
}
func GetBuffer3i(bid uint32, param int32, value1, value2, value3 *int32) {
	syscall.Syscall6(alGetBuffer3i.Addr(), 5,
		uintptr(bid),
		uintptr(param),
		uintptr(unsafe.Pointer(value1)),
		uintptr(unsafe.Pointer(value2)),
		uintptr(unsafe.Pointer(value3)),
		0)
}
func GetBufferiv(bid uint32, param int32, values *int32) {
	syscall.Syscall(alGetBufferiv.Addr(), 3,
		uintptr(bid),
		uintptr(param),
		uintptr(unsafe.Pointer(values)))
}
func DopplerFactor(value float32) {
	syscall.Syscall(alDopplerFactor.Addr(), 1,
		uintptr(value),
		0,
		0)
}
func DopplerVelocity(value float32) {
	syscall.Syscall(alDopplerVelocity.Addr(), 1,
		uintptr(value),
		0,
		0)
}
func SpeedOfSound(value float32) {
	syscall.Syscall(alSpeedOfSound.Addr(), 1,
		uintptr(value),
		0,
		0)
}
func DistanceModel(distanceModel float32) {
	syscall.Syscall(alDistanceModel.Addr(), 1,
		uintptr(distanceModel),
		0,
		0)
}

// AL/alc.h go bindings
func CreateContext(device Device, attrlist *int32) Context {
	ret, _, _ := syscall.Syscall(alcCreateContext.Addr(), 2,
		uintptr(device),
		uintptr(unsafe.Pointer(attrlist)),
		0)
	return Context(ret)
}
func MakeContextCurrent(context Context) bool {
	ret, _, _ := syscall.Syscall(alcMakeContextCurrent.Addr(), 1,
		uintptr(context),
		0,
		0)
	return ret == TRUE
}
func ProcessContext(context Context) {
	syscall.Syscall(alcProcessContext.Addr(), 1,
		uintptr(context),
		0,
		0)
}
func SuspendContext(context Context) {
	syscall.Syscall(alcSuspendContext.Addr(), 1,
		uintptr(context),
		0,
		0)
}
func DestroyContext(context Context) {
	syscall.Syscall(alcDestroyContext.Addr(), 1,
		uintptr(context),
		0,
		0)
}
func GetCurrentContext() Context {
	ret, _, _ := syscall.Syscall(alcGetCurrentContext.Addr(), 0,
		0,
		0,
		0)
	return Context(ret)
}
func GetContextsDevice(context Context) Device {
	ret, _, _ := syscall.Syscall(alcGetContextsDevice.Addr(), 1,
		uintptr(context),
		0,
		0)
	return Device(ret)
}
func OpenDevice(devicename string) Device {
	if devicename == "" {
		// request default device
		ret, _, _ := syscall.Syscall(alcOpenDevice.Addr(), 1,
			uintptr(0),
			0,
			0)
		return Device(ret)
	}
	str16, err := syscall.UTF16PtrFromString(devicename)
	if err != nil {
		return 0
	}
	ret, _, _ := syscall.Syscall(alcOpenDevice.Addr(), 1,
		uintptr(unsafe.Pointer(str16)),
		0,
		0)
	return Device(ret)
}
func CloseDevice(device Device) bool {
	ret, _, _ := syscall.Syscall(alcCloseDevice.Addr(), 1,
		uintptr(device),
		0,
		0)
	return ret == TRUE
}
func GetDeviceError(device Device) int32 {
	ret, _, _ := syscall.Syscall(alcGetError.Addr(), 1,
		uintptr(device),
		0,
		0)
	return int32(ret)
}
func IsDeviceExtensionPresent(device Device, extname string) bool {
	str16, err := syscall.UTF16PtrFromString(extname)
	if err != nil {
		return false
	}
	ret, _, _ := syscall.Syscall(alcIsExtensionPresent.Addr(), 2,
		uintptr(device),
		uintptr(unsafe.Pointer(str16)),
		0)
	return ret == TRUE
}
func GetDeviceProcAddress(device Device, fname string) Pointer {
	str16, err := syscall.UTF16PtrFromString(fname)
	if err != nil {
		return nil
	}
	ret, _, _ := syscall.Syscall(alcGetProcAddress.Addr(), 2,
		uintptr(device),
		uintptr(unsafe.Pointer(str16)),
		0)
	return Pointer(ret)
}
func GetDeviceEnumValue(device Device, ename string) int32 {
	str16, err := syscall.UTF16PtrFromString(ename)
	if err != nil {
		return -1
	}
	ret, _, _ := syscall.Syscall(alcGetEnumValue.Addr(), 2,
		uintptr(device),
		uintptr(unsafe.Pointer(str16)),
		0)
	return int32(ret)
}

//	func GetDeviceString(device Device, param int32) string {
//		ret, _, _ := syscall.Syscall(alcGetString.Addr(), 2,
//			uintptr(device),
//			uintptr(param),
//			0)
//		return ""
//		// return C.GoString(alcGetString((C.uintptr_t)(device), C.int(param)))
//	}
func GetDeviceIntegerv(device Device, param int32, size int32, data *int32) {
	syscall.Syscall6(alcGetIntegerv.Addr(), 4,
		uintptr(device),
		uintptr(param),
		uintptr(size),
		uintptr(unsafe.Pointer(data)),
		0,
		0)
}
func CaptureOpenDevice(devicename string, frequency uint32, format int32, buffersize int32) Device {
	str16, err := syscall.UTF16PtrFromString(devicename)
	if err != nil {
		return 0
	}
	ret, _, _ := syscall.Syscall6(alcCaptureOpenDevice.Addr(), 4,
		uintptr(unsafe.Pointer(str16)),
		uintptr(frequency),
		uintptr(format),
		uintptr(buffersize),
		0,
		0)
	return Device(ret)
}
func CaptureCloseDevice(device Device) bool {
	ret, _, _ := syscall.Syscall(alcCaptureCloseDevice.Addr(), 1,
		uintptr(device),
		0,
		0)
	return ret == TRUE
}
func CaptureStart(device Device) {
	syscall.Syscall(alcCaptureStart.Addr(), 1,
		uintptr(device),
		0,
		0)
}
func CaptureStop(device Device) {
	syscall.Syscall(alcCaptureStop.Addr(), 1,
		uintptr(device),
		0,
		0)
}
func CaptureSamples(device Device, buffer Pointer, samples int) {
	syscall.Syscall(alcCaptureSamples.Addr(), 3,
		uintptr(device),
		uintptr(buffer),
		uintptr(samples))
}

// Show which function pointers are bound [+] or not bound [-].
// Expected to be used as a sanity check to see if the OpenAL libraries exist.
func BindingReport() (report []string) {
	report = []string{}

	// AL/al.h
	report = append(report, "AL")
	report = append(report, isBound(unsafe.Pointer(alEnable), "alEnable"))
	report = append(report, isBound(unsafe.Pointer(alDisable), "alDisable"))
	report = append(report, isBound(unsafe.Pointer(alIsEnabled), "alIsEnabled"))
	report = append(report, isBound(unsafe.Pointer(alGetString), "alGetString"))
	report = append(report, isBound(unsafe.Pointer(alGetBooleanv), "alGetBooleanv"))
	report = append(report, isBound(unsafe.Pointer(alGetIntegerv), "alGetIntegerv"))
	report = append(report, isBound(unsafe.Pointer(alGetFloatv), "alGetFloatv"))
	report = append(report, isBound(unsafe.Pointer(alGetDoublev), "alGetDoublev"))
	report = append(report, isBound(unsafe.Pointer(alGetBoolean), "alGetBoolean"))
	report = append(report, isBound(unsafe.Pointer(alGetInteger), "alGetInteger"))
	report = append(report, isBound(unsafe.Pointer(alGetFloat), "alGetFloat"))
	report = append(report, isBound(unsafe.Pointer(alGetDouble), "alGetDouble"))
	report = append(report, isBound(unsafe.Pointer(alGetError), "alGetError"))
	report = append(report, isBound(unsafe.Pointer(alIsExtensionPresent), "alIsExtensionPresent"))
	report = append(report, isBound(unsafe.Pointer(alGetProcAddress), "alGetProcAddress"))
	report = append(report, isBound(unsafe.Pointer(alGetEnumValue), "alGetEnumValue"))
	report = append(report, isBound(unsafe.Pointer(alListenerf), "alListenerf"))
	report = append(report, isBound(unsafe.Pointer(alListener3f), "alListener3f"))
	report = append(report, isBound(unsafe.Pointer(alListenerfv), "alListenerfv"))
	report = append(report, isBound(unsafe.Pointer(alListeneri), "alListeneri"))
	report = append(report, isBound(unsafe.Pointer(alListener3i), "alListener3i"))
	report = append(report, isBound(unsafe.Pointer(alListeneriv), "alListeneriv"))
	report = append(report, isBound(unsafe.Pointer(alGetListenerf), "alGetListenerf"))
	report = append(report, isBound(unsafe.Pointer(alGetListener3f), "alGetListener3f"))
	report = append(report, isBound(unsafe.Pointer(alGetListenerfv), "alGetListenerfv"))
	report = append(report, isBound(unsafe.Pointer(alGetListeneri), "alGetListeneri"))
	report = append(report, isBound(unsafe.Pointer(alGetListener3i), "alGetListener3i"))
	report = append(report, isBound(unsafe.Pointer(alGetListeneriv), "alGetListeneriv"))
	report = append(report, isBound(unsafe.Pointer(alGenSources), "alGenSources"))
	report = append(report, isBound(unsafe.Pointer(alDeleteSources), "alDeleteSources"))
	report = append(report, isBound(unsafe.Pointer(alIsSource), "alIsSource"))
	report = append(report, isBound(unsafe.Pointer(alSourcef), "alSourcef"))
	report = append(report, isBound(unsafe.Pointer(alSource3f), "alSource3f"))
	report = append(report, isBound(unsafe.Pointer(alSourcefv), "alSourcefv"))
	report = append(report, isBound(unsafe.Pointer(alSourcei), "alSourcei"))
	report = append(report, isBound(unsafe.Pointer(alSource3i), "alSource3i"))
	report = append(report, isBound(unsafe.Pointer(alSourceiv), "alSourceiv"))
	report = append(report, isBound(unsafe.Pointer(alGetSourcef), "alGetSourcef"))
	report = append(report, isBound(unsafe.Pointer(alGetSource3f), "alGetSource3f"))
	report = append(report, isBound(unsafe.Pointer(alGetSourcefv), "alGetSourcefv"))
	report = append(report, isBound(unsafe.Pointer(alGetSourcei), "alGetSourcei"))
	report = append(report, isBound(unsafe.Pointer(alGetSource3i), "alGetSource3i"))
	report = append(report, isBound(unsafe.Pointer(alGetSourceiv), "alGetSourceiv"))
	report = append(report, isBound(unsafe.Pointer(alSourcePlayv), "alSourcePlayv"))
	report = append(report, isBound(unsafe.Pointer(alSourceStopv), "alSourceStopv"))
	report = append(report, isBound(unsafe.Pointer(alSourceRewindv), "alSourceRewindv"))
	report = append(report, isBound(unsafe.Pointer(alSourcePausev), "alSourcePausev"))
	report = append(report, isBound(unsafe.Pointer(alSourcePlay), "alSourcePlay"))
	report = append(report, isBound(unsafe.Pointer(alSourceStop), "alSourceStop"))
	report = append(report, isBound(unsafe.Pointer(alSourceRewind), "alSourceRewind"))
	report = append(report, isBound(unsafe.Pointer(alSourcePause), "alSourcePause"))
	report = append(report, isBound(unsafe.Pointer(alSourceQueueBuffers), "alSourceQueueBuffers"))
	report = append(report, isBound(unsafe.Pointer(alSourceUnqueueBuffers), "alSourceUnqueueBuffers"))
	report = append(report, isBound(unsafe.Pointer(alGenBuffers), "alGenBuffers"))
	report = append(report, isBound(unsafe.Pointer(alDeleteBuffers), "alDeleteBuffers"))
	report = append(report, isBound(unsafe.Pointer(alIsBuffer), "alIsBuffer"))
	report = append(report, isBound(unsafe.Pointer(alBufferData), "alBufferData"))
	report = append(report, isBound(unsafe.Pointer(alBufferf), "alBufferf"))
	report = append(report, isBound(unsafe.Pointer(alBuffer3f), "alBuffer3f"))
	report = append(report, isBound(unsafe.Pointer(alBufferfv), "alBufferfv"))
	report = append(report, isBound(unsafe.Pointer(alBufferi), "alBufferi"))
	report = append(report, isBound(unsafe.Pointer(alBuffer3i), "alBuffer3i"))
	report = append(report, isBound(unsafe.Pointer(alBufferiv), "alBufferiv"))
	report = append(report, isBound(unsafe.Pointer(alGetBufferf), "alGetBufferf"))
	report = append(report, isBound(unsafe.Pointer(alGetBuffer3f), "alGetBuffer3f"))
	report = append(report, isBound(unsafe.Pointer(alGetBufferfv), "alGetBufferfv"))
	report = append(report, isBound(unsafe.Pointer(alGetBufferi), "alGetBufferi"))
	report = append(report, isBound(unsafe.Pointer(alGetBuffer3i), "alGetBuffer3i"))
	report = append(report, isBound(unsafe.Pointer(alGetBufferiv), "alGetBufferiv"))
	report = append(report, isBound(unsafe.Pointer(alDopplerFactor), "alDopplerFactor"))
	report = append(report, isBound(unsafe.Pointer(alDopplerVelocity), "alDopplerVelocity"))
	report = append(report, isBound(unsafe.Pointer(alSpeedOfSound), "alSpeedOfSound"))
	report = append(report, isBound(unsafe.Pointer(alDistanceModel), "alDistanceModel"))

	// AL/alc.h
	report = append(report, "ALC")
	report = append(report, isBound(unsafe.Pointer(alcCreateContext), "alcCreateContext"))
	report = append(report, isBound(unsafe.Pointer(alcMakeContextCurrent), "alcMakeContextCurrent"))
	report = append(report, isBound(unsafe.Pointer(alcProcessContext), "alcProcessContext"))
	report = append(report, isBound(unsafe.Pointer(alcSuspendContext), "alcSuspendContext"))
	report = append(report, isBound(unsafe.Pointer(alcDestroyContext), "alcDestroyContext"))
	report = append(report, isBound(unsafe.Pointer(alcGetCurrentContext), "alcGetCurrentContext"))
	report = append(report, isBound(unsafe.Pointer(alcGetContextsDevice), "alcGetContextsDevice"))
	report = append(report, isBound(unsafe.Pointer(alcOpenDevice), "alcOpenDevice"))
	report = append(report, isBound(unsafe.Pointer(alcCloseDevice), "alcCloseDevice"))
	report = append(report, isBound(unsafe.Pointer(alcGetError), "alcGetError"))
	report = append(report, isBound(unsafe.Pointer(alcIsExtensionPresent), "alcIsExtensionPresent"))
	report = append(report, isBound(unsafe.Pointer(alcGetProcAddress), "alcGetProcAddress"))
	report = append(report, isBound(unsafe.Pointer(alcGetEnumValue), "alcGetEnumValue"))
	report = append(report, isBound(unsafe.Pointer(alcGetString), "alcGetString"))
	report = append(report, isBound(unsafe.Pointer(alcGetIntegerv), "alcGetIntegerv"))
	report = append(report, isBound(unsafe.Pointer(alcCaptureOpenDevice), "alcCaptureOpenDevice"))
	report = append(report, isBound(unsafe.Pointer(alcCaptureCloseDevice), "alcCaptureCloseDevice"))
	report = append(report, isBound(unsafe.Pointer(alcCaptureStart), "alcCaptureStart"))
	report = append(report, isBound(unsafe.Pointer(alcCaptureStop), "alcCaptureStop"))
	report = append(report, isBound(unsafe.Pointer(alcCaptureSamples), "alcCaptureSamples"))
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
