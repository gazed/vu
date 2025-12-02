// SPDX-FileCopyrightText : © 2022-2024 Galvanized Logic Inc.
// SPDX-License-Identifier: BSD-2-Clause

package vu

// asset.go describe application resources that can be loaded from files.
// The assets are used by game entities to create game objects.

import (
	"hash/crc64"
	"math"
	"math/rand"

	"github.com/gazed/vu/audio"
	"github.com/gazed/vu/load"
)

// ============================================================================
// asset

// asset uniquely identfies a loaded resource (whereas an entity identifies
// a game object which may be comprised of one or more assets)

// Assets implement the asset interface.
// An asset is unique based on a name and an asset type.
type asset interface {
	aid() aid      // Unique ID based on combination of asset type and name.
	label() string // Asset name.
}

// aid is a unique asset identifier.
// Asset identifiers are hashes generated from an asset name and type.
type aid uint64

// kind returns the type of asset data for this aid.
func (a aid) kind() uint32 { return uint32(a & math.MaxUint8) }

// Asset types used in aid and aid.kind.
type assetType uint32

const (
	fnt        assetType = iota // font
	shd                         // shader configuration
	mat                         // material
	msh                         // mesh
	tex                         // texture
	aud                         // audio
	anm                         // animation rig and skin
	assetTypes                  // end of asset types - must be last.
)

// assetID produces a unique asset identifier using for the given
// asset type t, and asset name. Keep as many stringHash bits as possible
// to avoid collisions.
func assetID(t assetType, name string) aid { return aid(t) + aid(stringHash(name))<<8 }

// stringHash turns a string into a number.
func stringHash(s string) (hash uint64) {
	return crc64.Checksum([]byte(s), crcTable)
}

var crcTable = crc64.MakeTable(rand.Uint64())

// ============================================================================
// mesh references vertex data.
type mesh struct {
	name string // Unique mesh name.
	tag  aid    // name and type as a number.
	mid  uint32 // GPU vertex data reference.
}

// newMesh allocates space for a mesh structure,
// including space to store buffer data.
func newMesh(name string) *mesh {
	m := &mesh{name: name, tag: assetID(msh, name)}
	return m
}

// implement assset interface
func (m *mesh) aid() aid      { return m.tag }  // hashed type and name.
func (m *mesh) label() string { return m.name } // asset name

// ============================================================================
// shader references a loaded shader and its configuration.
type shader struct {
	name string // Unique shader name.
	tag  aid    // name and type as a number.
	sid  uint16 // GPU shader reference.

	// shader configuration and list of uniforms.
	config *load.Shader // shader configuration.
}

// newShader allocates space for a mesh structure,
// including space to store buffer data.
func newShader(name string) *shader {
	s := &shader{name: name, tag: assetID(shd, name)}
	return s
}

// setConfig is called when the shader configuration has been loaded.
func (s *shader) setConfig(cfg *load.Shader) {
	s.config = cfg // the full shader config.
}

// implement assset interface
func (s *shader) aid() aid      { return s.tag }  // hashed type and name.
func (s *shader) label() string { return s.name } // asset name

// ============================================================================
// texture references 2D images
//
// texture is an an optional, but very common, part of a rendered model.
// texture data is copied to the graphics card. One or more textures
// can be associated with a model entity and consumed by a shader.
type texture struct {
	name   string // Unique name of the texture.
	tag    aid    // Name and type as a number.
	tid    uint32 // GPU texture reference.
	opaque bool   // All pixels have alpha==1.0.
}

// newTexture allocates space for a texture asset.
func newTexture(name string) *texture {
	return &texture{name: name, tag: assetID(tex, name)}
}

// implement assset interface
func (t *texture) aid() aid      { return t.tag }  // hashed type and name.
func (t *texture) label() string { return t.name } // asset name

// =============================================================================
// material

// material is used to color a mesh. It specifies the surface color and
// how the surface is lit. Materials are applied to a rendered model by
// a shader.
type material struct {
	name      string  // Unique matrial name.
	tag       aid     // name and type as a number.
	color     rgba    // material base color for a solid PBR.
	metallic  float32 // metallic value if no m-r texture
	roughness float32 // roughness value if no m-r texture
}

// newMaterial allocates space for material values.
// default values are all zero except color alpha set to 1.0.
func newMaterial(name string) *material {
	mat := &material{name: name, tag: assetID(mat, name)}
	mat.color.a = 1.0 // opaque.
	return mat
}

// implement assset interface
func (m *material) aid() aid      { return m.tag }  // hashed type and name.
func (m *material) label() string { return m.name } // asset name

// rgb holds color information where each field is expected to contain
// a value from 0.0 to 1.0. A value of 0 means none of that color while a value
// of 1 means as much as possible of that color.
type rgba struct {
	r float32 // Red.
	g float32 // Green.
	b float32 // Blue.
	a float32 // Alpha.
}

// =============================================================================
// sound

// sound is an engine sound asset. Expected to be accessed through
// the sounds component.
type sound struct {
	name       string      // Unique name of the sound.
	tag        aid         // name and type as a number.
	sid        uint64      // Audio card identifier related to sound location.
	did        uint64      // Audio data reference identifier.
	lx, ly, lz float64     // noise location.
	data       *audio.Data // noise data.
}

// newSound allocates space for a texture object.
func newSound(name string) *sound {
	return &sound{name: name, tag: assetID(aud, name), data: &audio.Data{}}
}

// implement assset interface
func (s *sound) aid() aid      { return s.tag }  // hashed type and name.
func (s *sound) label() string { return s.name } // asset name
