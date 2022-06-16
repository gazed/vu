// Copyright © 2024 Galvanized Logic Inc.

package load

// glb.go imports a subset of the GLTF specification.

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"log/slog"

	"github.com/gazed/vu/internal/load/gltf"
)

// Glb accepts a subset of a binary gltf.Document that describes
// a single mesh model. The single models limitiations is enforced as follows:
//   - one Scene
//   - one Node
//   - one Mesh
//   - one Mesh.Primitive
//
// These conform to a single model exported from Blender.
func Glb(name string, r io.Reader) (data []AssetData) {
	doc := gltf.Document{}
	gltf.NewDecoder(r).Decode(&doc)

	// reduce complexity by focusing on single mesh models.
	if len(doc.Scenes) != 1 {
		return []AssetData{{Filename: name, Err: fmt.Errorf("expecting one gltf Scene")}}
	}
	if len(doc.Nodes) != 1 {
		return []AssetData{{Filename: name, Err: fmt.Errorf("expecting one gltf Node")}}
	}
	if len(doc.Meshes) != 1 {
		return []AssetData{{Filename: name, Err: fmt.Errorf("expecting one gltf Mesh")}}
	}

	// primitive contains the mesh data.
	if len(doc.Meshes[0].Primitives) != 1 {
		return []AssetData{{Filename: name, Err: fmt.Errorf("expecting one gltf Mesh.Primitive")}}
	}
	primitive := doc.Meshes[0].Primitives[0]
	if len(primitive.Attributes) < 1 {
		return []AssetData{{Filename: name, Err: fmt.Errorf("expecting one or more vertex attributes")}}
	}

	// Blender exports non-interleaved vertex data so use that.
	md := make(MeshData, VertexTypes)
	for k, v := range primitive.Attributes {
		accessor := doc.Accessors[v]
		buffview := doc.BufferViews[*accessor.BufferView]
		buff := doc.Buffers[buffview.Buffer].Data
		offset := buffview.ByteOffset
		byteLen := buffview.ByteLength

		// sanity check: all the vertex data counts must match.
		if accessor.Count != doc.Accessors[0].Count {
			return []AssetData{{Filename: name, Err: fmt.Errorf("invalid vertex attributes")}}
		}

		// make copies of the gltf doc data.
		switch k {
		case gltf.POSITION:
			if accessor.Type != gltf.AccessorVec3 || accessor.ComponentType != gltf.ComponentFloat {
				return []AssetData{{Filename: name, Err: fmt.Errorf("expecting vec3:float32 vertexes")}}
			}
			md[Vertexes].Stride = 12 // bytes for vec3 of float32
			md[Vertexes].Count = accessor.Count
			md[Vertexes].Data = make([]byte, byteLen)
			copy(md[Vertexes].Data, buff[offset:offset+byteLen])
		case gltf.NORMAL:
			if accessor.Type != gltf.AccessorVec3 || accessor.ComponentType != gltf.ComponentFloat {
				return []AssetData{{Filename: name, Err: fmt.Errorf("expecting vec3:float32 normals")}}
			}
			md[Normals].Stride = 12 // bytes for vec3 of float32
			md[Normals].Count = accessor.Count
			md[Normals].Data = make([]byte, byteLen)
			copy(md[Normals].Data, buff[offset:offset+byteLen])
		case gltf.TEXCOORD_0:
			if accessor.Type != gltf.AccessorVec2 || accessor.ComponentType != gltf.ComponentFloat {
				return []AssetData{{Filename: name, Err: fmt.Errorf("expecting vec2:float32 texcoords")}}
			}
			md[Texcoords].Stride = 8 // 4x2 bytes for vec2 of float32
			md[Texcoords].Count = accessor.Count
			md[Texcoords].Data = make([]byte, byteLen)
			copy(md[Texcoords].Data, buff[offset:offset+byteLen])
		default:
			slog.Warn("unprocessed data", "data", k)
		}
	}

	// get mesh indices if they exist
	if primitive.Indices != nil {
		accessor := doc.Accessors[*primitive.Indices]
		if accessor.Type != gltf.AccessorScalar || accessor.ComponentType != gltf.ComponentUshort {
			return []AssetData{{Filename: name, Err: fmt.Errorf("expecting uint16 indexes")}}
		}
		md[Indexes].Stride = 2 // bytes for uint16
		md[Indexes].Count = accessor.Count
		viewIndex := accessor.BufferView
		buffview := doc.BufferViews[*viewIndex]
		offset := buffview.ByteOffset
		byteLen := buffview.ByteLength
		buff := doc.Buffers[buffview.Buffer].Data
		md[Indexes].Data = make([]byte, byteLen)
		copy(md[Indexes].Data, buff[offset:offset+byteLen])
	}
	data = append(data, AssetData{Filename: name, Data: md, Err: nil})

	// load material, expecting 1 PRB material.
	if len(doc.Materials) != 1 || doc.Materials[0].PBRMetallicRoughness == nil {
		return []AssetData{{Filename: name, Err: fmt.Errorf("expecting one PBR material")}}
	}
	mat := doc.Materials[0].PBRMetallicRoughness
	switch {
	case mat.BaseColorTexture != nil && mat.MetallicRoughnessTexture != nil:
		// support PRB materials with base color and metallic-roughness textures
		// corresponds to the pbr_texture shader - create two textures
		bc, err := getPNGTexture(&doc, mat.BaseColorTexture.Index)
		if err != nil {
			return []AssetData{{Filename: name, Err: fmt.Errorf("PBR base color:%w", err)}}
		}
		mr, err := getPNGTexture(&doc, mat.MetallicRoughnessTexture.Index)
		if err != nil {
			return []AssetData{{Filename: name, Err: fmt.Errorf("PBR material-roughness:%w", err)}}
		}
		data = append(data, AssetData{Filename: name, Data: bc, Err: nil})
		data = append(data, AssetData{Filename: name, Data: mr, Err: nil})
		slog.Debug("load.Glb: texture based PBR")
	case mat.BaseColorTexture != nil:
		// support PRB materials with just a base color texture.
		// corresponds to the pbr_base_color shader - create one texture and one material.
		img, err := getPNGTexture(&doc, mat.BaseColorTexture.Index)
		if err != nil {
			return []AssetData{{Filename: name, Err: err}}
		}
		data = append(data, AssetData{Filename: name, Data: img, Err: nil})

		// the other material values are constant.
		pbr := PBRMaterialData{
			Metallic:  mat.MetallicFactorOrDefault(),
			Roughness: mat.RoughnessFactorOrDefault(),
		}
		data = append(data, AssetData{Filename: name, Data: pbr, Err: nil})
		slog.Debug("load.Glb: color texture PBR")
	case mat.BaseColorFactor != nil:
		// support solid shaded PRB materials.
		// corresponds to the pbr_solid shader - create one material.
		color := mat.BaseColorFactorOrDefault()
		pbr := PBRMaterialData{
			Metallic:  mat.MetallicFactorOrDefault(),
			Roughness: mat.RoughnessFactorOrDefault(),
		}
		pbr.ColorR = color[0]
		pbr.ColorG = color[1]
		pbr.ColorB = color[2]
		pbr.ColorA = color[3]
		data = append(data, AssetData{Filename: name, Data: pbr, Err: nil})
		slog.Debug("load.Glb: solid color PBR")
	default:
		return []AssetData{{Filename: name, Err: fmt.Errorf("missing PBR material parameters")}}
	}

	// FUTURE load animation (joints, and weights) vertex data

	return data
}

// getPNGTexture attemps to return a png NRGBA image for the given texture.
func getPNGTexture(doc *gltf.Document, textureIndex uint32) (idata *ImageData, err error) {
	tex := doc.Textures[textureIndex]
	if tex.Source == nil {
		return nil, fmt.Errorf("expecting texture image index")
	}
	img := doc.Images[*tex.Source]
	if img.MimeType != "image/png" {
		return nil, fmt.Errorf("expecting png texture image")
	}
	if img.BufferView == nil {
		return nil, fmt.Errorf("expecting texture buffer view")
	}

	// get the buffer view
	view := doc.BufferViews[*img.BufferView]
	offset := view.ByteOffset
	byteLen := view.ByteLength

	// get the bytes from the buffer data
	buff := doc.Buffers[view.Buffer].Data
	ibytes := buff[offset : offset+byteLen]

	// get a png image from the bytes
	pngImg, err := png.Decode(bytes.NewReader(ibytes))
	if err != nil {
		return nil, fmt.Errorf("expecting png image data")
	}

	// return the image data for the png image.
	idata = &ImageData{}
	switch t := pngImg.(type) {
	case *image.NRGBA:
		idata.Pixels = []byte(t.Pix)
	default:
		return nil, fmt.Errorf("expecting png NRGBA, got %T", t)
	}
	idata.Width = uint32(pngImg.Bounds().Size().X)
	idata.Height = uint32(pngImg.Bounds().Size().Y)
	return idata, err
}
