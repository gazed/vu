// Copyright © 2013-2024 Galvanized Logic Inc.

// Package load fetches disk based 3D asset data. It's main purpose is to
// find the asset file and load its data into intermediate data structs.
//   - ".spv"  spir-v shader module byte code
//   - ".png"  image data
//   - ".shd"  shader configuration description
//   - ".glb"  vertex data, image data, animation data, material data
//   - ".wav"  audio data
//   - ".fnt"  font mapping data
//   - ".yaml" data file
//
// This package is primary used internally for getting data from disk
// that is then upload to the render and audio systems.
//
// Package load is provided as part of the vu (virtual universe) 3D engine.
package load

// FUTURE move load to internal/load and expose any useful application
// facing APIs through Engine wrapper methods.

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"path"
	"strings"
)

// =============================================================================
// asset location conventions.

// assetDirs expects asset files to be in directories based on asset type.
// These are the default directories and can be overridden using SetAssetDir.
var assetDirs = map[string]string{
	".spv":  "assets/shaders", // spir-v compiled shader byte files
	".shd":  "assets/shaders", // yaml shader configuration files.
	".png":  "assets/images",  // png images, often textures.
	".glb":  "assets/models",  // glb scenes, meshes, materials, animations, textures,...
	".fnt":  "assets/fonts",   // glyph files exspect corresponding image.
	".wav":  "assets/audio",   // sound data.
	".yaml": "assets/data",    // data files
}

// SetAssetDir can be used to add or change the default directory
// conventions for finding assets.
func SetAssetDir(ext, dir string) {
	ext = strings.ToLower(ext)
	assetDirs[ext] = dir
}

// LoadAssetFile returns one or more AssetData structs from the
// given asset filename.
func LoadAssetFile(fname string) []AssetData {
	switch getFileExtension(fname) {
	case ".spv":
		data, err := ShaderBytes(fname)
		return []AssetData{{Filename: fname, Data: ShaderData(data), Err: err}}
	case ".shd":
		shader, err := ShaderConfig(fname)
		return []AssetData{{Filename: fname, Data: shader, Err: err}}
	case ".png":
		img, err := Image(fname)
		return []AssetData{{Filename: fname, Data: img, Err: err}}
	case ".glb":
		return Model(fname) // possible to have multiple assets
	case ".wav":
		aud, err := Audio(fname)
		return []AssetData{{Filename: fname, Data: aud, Err: err}}
	case ".fnt":
		fnt, err := Font(fname)
		return []AssetData{{Filename: fname, Data: fnt, Err: err}}
	}
	err := fmt.Errorf("unsupported asset file %s", fname)
	return []AssetData{{Filename: fname, Err: err}}
}

// AssetData is used to return loaded data.
// The caller is expected to switch on the Data type.
type AssetData struct {
	Filename string      // request filename with extension
	Err      error       // nil if load was successful
	Data     interface{} // struct of loaded data.
}

// getData returns the raw bytes in the requested file.
//
//	filename: name of the file including the file extension.
func getData(filename string) (data []byte, err error) {
	assetDir := "" // default to local directory.
	extension := getFileExtension(filename)
	if dir, defined := assetDirs[extension]; defined {
		assetDir = dir
	}
	filepath := strings.TrimSpace(path.Join(assetDir, filename))
	return os.ReadFile(filepath)
}

// getFileExtension returns the given filename extension
// including the leading dot ".".
func getFileExtension(filename string) (extension string) {
	return strings.ToLower(path.Ext(filename))
}

// =============================================================================
// data files ie: ".yaml"

// ByteData loads the bytes for the given file.
type ByteData []byte

// DataBytes loads bytes from generic data files.
func DataBytes(name string) (data ByteData, err error) {
	data, err = getData(name)
	if err != nil {
		return data, fmt.Errorf("data byte load %s: %w", name, err)
	}
	return data, nil
}

// =============================================================================
// ".spv" glsl shader byte code

// ShaderData differentiates shader data from other byte data.
type ShaderData []byte

// ShaderBytes loads compiled shader (spir-v) byte data.
func ShaderBytes(name string) (data ShaderData, err error) {
	data, err = getData(name)
	if err != nil {
		return data, fmt.Errorf("shader byte data load %s: %w", name, err)
	}
	return data, nil
}

// =============================================================================
// ".shd" custom yaml shader configuration data.

// ShaderConfig loads compiled shader (spir-v) byte data.
func ShaderConfig(name string) (cfg *Shader, err error) {
	data, err := getData(name)
	if err != nil {
		return cfg, fmt.Errorf("shader config load %s: %w", name, err)
	}
	return Shd(name, data)
}

// =============================================================================
// ".png" image loading.

// ImageData contains image data for uploading to the GPU.
type ImageData struct {
	Width  uint32
	Height uint32
	Pixels []byte
	Opaque bool
}

// Image loads .png images as the underlying data format for textures.
func Image(name string) (idata *ImageData, err error) {
	data, err := getData(name)
	if err != nil {
		return idata, fmt.Errorf("image load %s: %w", name, err)
	}
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return idata, fmt.Errorf("image decode %s: %w", name, err)
	}
	idata = &ImageData{}
	switch t := img.(type) {
	case *image.NRGBA:
		idata.Pixels = []byte(t.Pix)
		idata.Opaque = t.Opaque()
	default:
		return idata, fmt.Errorf("image expecting NRGBA: %s: %T", name, t)
	}
	idata.Width = uint32(img.Bounds().Size().X)
	idata.Height = uint32(img.Bounds().Size().Y)
	return idata, nil
}

// =============================================================================
// ".glb" gltf binary files
// https://github.com/KhronosGroup/glTF-Tutorials/blob/master/gltfTutorial/README.md

// Vertex MeshData attribute types.
const (
	Vertexes    = iota // 0 required:V3 float32
	Texcoords          // 1 optional:V2 float32
	Normals            // 2 optional:V3 float32
	Tangents           // 3 optional:V4 float32
	Colors             // 4 optional:V3 uint8
	Joints             // 5 FUTURE:V4 uint8    animations
	Weights            // 6 FUTURE:V4 uint8    animations
	Indexes            // 7 required:uint16             - must be second last.
	VertexTypes        // 8 number of vertex data types - must be last.
)

// Instance Data attribute types describe per-instance model data.
// These are defined here because they are similar to vertex data types.
const (
	InstancePosition = iota // position 0 V3 float32
	InstanceColors          // 1 V3 uint8
	InstanceScales          // 2 float32
	InstanceTypes           // number of instance data types - must be last.
)

// MeshData contains per-vertex data. Data will be index the GLTF buffer.
type MeshData []Buffer

// PBRMaterialData describes a PBR solid material.
type PBRMaterialData struct {
	ColorR    float64 // material base color for solid PBR.
	ColorG    float64 //  ""
	ColorB    float64 //  ""
	ColorA    float64 //  ""
	Metallic  float64 // metallic value if no m-r texture
	Roughness float64 // roughness value if no m-r texture
}

// Model loads 3D data including mesh, texture, material, animation,
// and transform data. Errors are returned inside AssetData which
// will always contain at least one element.
func Model(name string) (data []AssetData) {
	dbytes, err := getData(name)
	if err != nil {
		return []AssetData{{Filename: name, Err: fmt.Errorf("model load %s: %w", name, err)}}
	}
	return Glb(name, bytes.NewReader(dbytes))
}

// =============================================================================
// ".fnt" - font glyph mappings

// FontData holds UV texture mapping information for a font.
// It is intended for populating rendered models of strings.
// This is an intermediate data format that needs further processing by
// something like a vu.Ent.MakeModel to bind the data to a GPU and associate
// it with a texture atlas containing the bitmapped font images.
type FontData struct {
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

// Font loads font character mapping data. Existing FontData is
// overwritten with information found by the Locator.
func Font(name string) (fnt *FontData, err error) {
	data, err := getData(name)
	if err != nil {
		return fnt, fmt.Errorf("font load %s: %w", name, err)
	}
	return Fnt(bytes.NewReader(data))
}

// =============================================================================
// ".wav" - audio data.

// AudioData consists of the actual audio data bytes along with sounds attributes
// that describe how the sound data is interpreted and played.
//
// AudioData is an intermediate data format that needs further processing
// to associate the sound with a 3D location and bind it to an audio device.
type AudioData struct {
	Data  []byte           // the sound data bytes.
	Attrs *AudioAttributes // Attributes describing the sound data.
}

// AudioAttributes describe how sound data is interpreted and played.
type AudioAttributes struct {
	Channels   uint16 // Number of audio channels.
	Frequency  uint32 // 8000, 44100, etc.
	DataSize   uint32 // Size of audio data.
	SampleBits uint16 // 8 bits = 8, 16 bits = 16, etc.
}

// Audio loads audio data.
func Audio(name string) (aud *AudioData, err error) {
	data, err := getData(name)
	if err != nil {
		return aud, fmt.Errorf("audio load %s: %w", name, err)
	}
	return Wav(bytes.NewReader(data))
}
