// Copyright © 2024 Galvanized Logic Inc.

package load

// buffer.go provides a means of sharing vertex data and
// other data that can be uploaded to the GPU.

import (
	"fmt"
	"unsafe"
)

// Buffer holds byte data that can be uploaded to the GPU
type Buffer struct {
	Data   []byte // bytes in little-endian order.
	Count  uint32 // total number of elements, ie: how many vertex vec3's.
	Stride uint32 // bytes for one element, eg: 12 for float32:vec3
}

// F32Buffer converts a slice of float32s to a Buffer of bytes.
// Used to pass data to glsl vec float types
func F32Buffer(verts []float32, dimension uint32) Buffer {
	return Buffer{
		Stride: 4 * dimension, // 4 bytes for each float32 * (vec2 or vec3)
		Count:  uint32(len(verts)) / dimension,
		Data:   unsafe.Slice((*byte)(unsafe.Pointer(&verts[0])), len(verts)*4),
	}
}

// U32Buffer converts a slice of uint32 to a Buffer of bytes.
// Used to pass data to glsl uvec int types
func U32Buffer(data []uint32, dimension uint32) Buffer {
	return Buffer{
		Stride: 4 * dimension, // 4 bytes for each uint32 * (uvec3 or uvec4)
		Count:  uint32(len(data)) / dimension,
		Data:   unsafe.Slice((*byte)(unsafe.Pointer(&data[0])), len(data)*4),
	}
}

// U16Buffer converts a slice of uint16 to a Buffer of bytes.
// Used to pass data to glsl vertex indexes.
func U16Buffer(indexes []uint16) Buffer {
	return Buffer{
		Stride: 2, // 2 bytes for each uint16
		Count:  uint32(len(indexes)),
		Data:   unsafe.Slice((*byte)(unsafe.Pointer(&indexes[0])), len(indexes)*2),
	}
}

// PrintF32 dumps bytes as float32 vertices. Used to debug mesh data.
// eg: md[load.Vertexes].PrintF32()
func (buff Buffer) PrintF32(name string) {
	vsize := len(buff.Data) / 4 // number of vertexes.
	vx := unsafe.Slice((*float32)(unsafe.Pointer(&buff.Data[0])), vsize)
	fmt.Printf("%s:%d\n", name, len(vx))
	dim := vsize / int(buff.Count)
	for i := 0; i < len(vx); i += dim {
		switch dim {
		case 2:
			fmt.Printf("  %+f,%+f,\n", vx[i], vx[i+1])
		case 3:
			fmt.Printf("  %+f,%+f,%+f,\n", vx[i], vx[i+1], vx[i+2])
		}
	}
}

// PrintU16 dumps bytes as uint16 triangle indexes. Used to debug mesh data.
func (buff Buffer) PrintU16(name string) {
	ix := unsafe.Slice((*uint16)(unsafe.Pointer(&buff.Data[0])), buff.Count)
	fmt.Printf("%s:%d\n", name, len(ix))
	for i := 0; i < len(ix); i += 3 {
		fmt.Printf("  %d,%d,%d,\n", ix[i], ix[i+1], ix[i+2])
	}
}
