// Copyright © 2013-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// Package load fetches disk based 3D assets. Assets are loaded into
// one of the following intermediate data structures:
//    FntData.Load uses Fnt to load bitmapped characters.
//    ImgData.Load uses Png to load model textures.
//    ModData.Load uses Iqm to load animated models.
//    MshData.Load uses Obj to load static models.
//    MtlData.Load uses Mtl to load model lighting data.
//    ShdData.Load uses Src to load GPU shader programs.
//    SndData.Load uses Wav to load 3D audio.
// Each intermediate data format is currently associated with one file
// format. Asset loading is currently intended for smaller 3D applications
// where data is loaded directly from disk to memory, i.e. no database.
//
// A default file Locator is provided. It can be replaced with
// a different Locator that follows a different string based naming
// convention for finding disk based assets. Overall package load
// attempts to shield users from knowledge about:
//    File Formats  : how asset contents are stored on disk.
//    File Types    : how file types map to asset data structs.
//    File Locations: where asset files are stored on disk.
//
// Package load is provided as part of the vu (virtual universe) 3D engine.
package load

import (
	"fmt"
	"image"
	"io"

	"github.com/gazed/vu/math/lin"
)

// Design Notes: Balance between full functionality and maintainability,
//               ie: anything more than the absolute minimum implementation
//               needs to be justified by value added.
//             : Prefer documented public specifications for asset file types.
// FUTURE: Add full support for obj, iqm specs.
// FUTURE: wrap or develop more import formats. See possibilities at the
//         Open Asset Import Library: http://assimp.sourceforge.net
// FUTURE: Load data into formats that can be immediately transferred
//         to the GPU or audio card without further processing.
// FUTURE: Optional industrial strength database back end?

// =============================================================================

// FntData holds UV texture mapping information for a font.
// It is intended for populating rendered models of strings.
// This is an intermediate data format that needs further processing by
// something like vu/Model to bind the data to a GPU and associate
// it with a texture atlas containing the bitmapped font images.
type FntData struct {
	W, H  int       // Width and height
	Chars []ChrData // Character data.
}

// ChrData holds UV texture mapping information for one character.
// Expected to be used as part of FntData.
type ChrData struct {
	Char       rune // Character.
	X, Y, W, H int  // Character bit size.
	Xo, Yo, Xa int  // Character offset.
}

// Load font character mapping data. Existing FntData is
// overwritten with information found by the Locator.
func (d *FntData) Load(name string, l Locator) (err error) {
	fname := name + ".fnt" // FUTURE: other font file formats.
	var reader io.ReadCloser
	if reader, err = l.GetResource(fname); err != nil {
		return fmt.Errorf("Could not load glyphs from %s: %s\n", fname, err)
	}
	defer reader.Close()
	return Fnt(reader, d)
}

// FntData
// =============================================================================
// ImgData

// ImgData uses standard images as the underlying data format for textures.
// The image height, width, and bytes in (N)RGBA format are:
//    width  := img.Bounds().Max.X - img.Bounds().Min.X
//    height := img.Bounds().Max.Y - img.Bounds().Min.Y
//    rgba, _ := img.(*image.(N)RGBA)
// Note that golang NRGBA are images with an alpha channel, but without alpha
// pre-multiplication. RGBA are images originally without an alpha channel,
// but assigned an alpha of 1 when read in.
//
// This is an intermediate data format that needs further processing by
// something like vu/Model to bind the data to a GPU based texture.
type ImgData struct {
	Img image.Image
}

// Load image data. Existing ImgData is discarded and
// replaced with information found by the Locator.
func (d *ImgData) Load(name string, l Locator) (err error) {
	fname := name + ".png" // FUTURE: other image file formats.
	var reader io.ReadCloser
	if reader, err = l.GetResource(fname); err != nil {
		return fmt.Errorf("Could not load image from %s: %s\n", fname, err)
	}
	defer reader.Close()
	return Png(reader, d)
}

// FntData
// =============================================================================
// ModData

// ModData combines vertex data, animation data and some texture
// names for a complete animated model. It is an intermediate data
// format that needs further processing by something like vu/Model
// to bind the data to a GPU.
type ModData struct {
	MshData          // Vertex based data.
	AnmData          // Animation data.
	TMap    []TexMap // Texture name and vertex mapping data.
}

// AnmData holds the data necessary to a. It is an intermediate data
// format that needs further processing by something like vu/Model
// to bind the data to a GPU.
type AnmData struct {
	Movements []Movement // Indexes into the given frames.
	Blends    []byte     // Vertex blend indicies. Arranged as [][4]byte
	Weights   []byte     // Vertex blend weights.  Arranged as [][4]byte
	Joints    []int32    // Joint parent information for each joint.
	Frames    []*lin.M4  // Animation transforms: [NumFrames][NumJoints].
}

// Movement marks a number of frames as a particular animated move that
// affects frames from F0 to F0+FN. Expected to be used as part of AnmData.
type Movement struct {
	Name   string  // Name of the animation
	F0, Fn uint32  // First frame, number of frames.
	Rate   float32 // Frames per second.
}

// TexMap allows a model to have multiple textures. The named texture
// resource affects triangle faces from F0 to F0+FN. Expected to be used
// as part of ModData.
type TexMap struct {
	Name   string // Name of the texture resource.
	F0, Fn uint32 // First triangle face index and number of faces.
}

// Load model vertex and animation data. Existing ModData is
// overwritten with information found by the Locator.
func (d *ModData) Load(name string, l Locator) (err error) {
	fname := name + ".iqm" // FUTURE: other animated model file formats.
	var reader io.ReadCloser
	if reader, err = l.GetResource(fname); err != nil {
		return fmt.Errorf("Could not load animated model from %s: %s\n", fname, err)
	}
	defer reader.Close()
	return Iqm(reader, d)
}

// ModData
// =============================================================================
// MshData

// MshData stores vertex data from .obj files.
// It is intended for populating rendered models.
// The V,F buffers are expected to have data. The N,T,X buffers are optional.
//
// MshData is an intermediate data format that needs further processing
// by something like vu/Model to bind the data to a GPU.
type MshData struct {
	Name string    // Imported model name.
	V    []float32 // Vertex positions.    Arranged as [][3]float32
	N    []float32 // Vertex normals.      Arranged as [][3]float32
	T    []float32 // Texture coordinates. Arranged as [][2]float32
	X    []float32 // Vertex tangents.     Arranged as [][2]float32
	F    []uint16  // Triangle faces.      Arranged as [][3]uint16
}

// Load model mesh vertex data. Existing MshData is
// overwritten with information found by the Locator.
func (d *MshData) Load(name string, l Locator) (err error) {
	fname := name + ".obj" // FUTURE: other model mesh file formats.
	var reader io.ReadCloser
	if reader, err = l.GetResource(fname); err != nil {
		return fmt.Errorf("Could not load mesh data from %s: %s\n", fname, err)
	}
	defer reader.Close()
	return Obj(reader, d)
}

// MshData
// =============================================================================
// MtlData

// MtlData holds color and alpha information.
// It is intended for populating rendered models and is
// often needed for as attributes for shaders with lighting.
//
// MtlData is an intermediate data format that needs further processing
// by something like vu/Model to bind the data to uniforms in a GPU shader.
type MtlData struct {
	KaR, KaG, KaB float32 // Ambient color.
	KdR, KdG, KdB float32 // Diffuse color.
	KsR, KsG, KsB float32 // Specular color.
	Ns            float32 // Specular exponent.
	Alpha         float32 // Transparency
}

// Load model lighting material data. Existing MtlData is
// overwritten with information found by the Locator.
func (d *MtlData) Load(name string, l Locator) (err error) {
	fname := name + ".mtl" // FUTURE: other model material file formats.
	var reader io.ReadCloser
	if reader, err = l.GetResource(fname); err != nil {
		return fmt.Errorf("could not open %s %s", fname, err)
	}
	defer reader.Close()
	return Mtl(reader, d)
}

// MshData
// =============================================================================
// ShdData

// ShdData includes both vertex and fragment shader program source.
// Each line is terminated by a linefeed so that it will compile
// on the GPU.
//
// ShdData is an intermediate data format that needs further processing
// by something like vu/Model to compile the shader program and bind
// it to a GPU.
type ShdData struct {
	Vsh SrcData // Vertex shader.
	Fsh SrcData // Fragment (pixel) shader.
}

// Load vertex and fragment shader program source code. Shader source is
// appended to the existing ShdData source. Assumes both the vertex and
// fragment shader files have the same name prefix.
func (d *ShdData) Load(name string, l Locator) (err error) {
	fname := name + ".vsh"
	var vr io.ReadCloser
	if vr, err = l.GetResource(fname); err != nil {
		return fmt.Errorf("Load vertex shader error %s: %s\n", fname, err)
	}
	defer vr.Close()
	if d.Vsh, err = Src(vr); err != nil {
		return fmt.Errorf("Load vertex shader error %s: %s\n", fname, err)
	}

	fname = name + ".fsh"
	var fr io.ReadCloser
	if fr, err = l.GetResource(fname); err != nil {
		return fmt.Errorf("Load fragment shader error %s: %s\n", fname, err)
	}
	defer fr.Close()
	d.Fsh, err = Src(fr)
	return err
}

// ShdData
// =============================================================================
// SndData

// SndData consists of the actual audio data bytes along with sounds attributes
// that describe how the sound data is interpreted and played.
//
// SndData is an intermediate data format that needs further processing
// by something like vu/Pov to associate the noise with a 3D location
// and bind it to an audio card.
type SndData struct {
	Data  []byte         // the sound data bytes.
	Attrs *SndAttributes // Attributes describing the sound data.
}

// SndAttributes describe how sound data is interpreted and played.
type SndAttributes struct {
	Channels   uint16 // Number of audio channels.
	Frequency  uint32 // 8000, 44100, etc.
	DataSize   uint32 // Size of audio data.
	SampleBits uint16 // 8 bits = 8, 16 bits = 16, etc.
}

// Load sound data. Existing SoundData is
// overwritten with information found by the Locator.
func (d *SndData) Load(name string, l Locator) (err error) {
	fname := name + ".wav" // FUTURE: other 3D audio file formats.
	var reader io.ReadCloser
	if reader, err = l.GetResource(fname); err != nil {
		return fmt.Errorf("Load sound error %s: %s\n", fname, err)
	}
	defer reader.Close()
	return Wav(reader, d)
}
