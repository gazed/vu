// Copyright © 2016-2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// aid.go defines asset identfiers. Assets are data used by shaders.
// See eid.go for entity identifiers.

import (
	"hash/crc64"
	"math"
	"math/rand"
)

// asset describes any data asset that can uniquely identify itself.
// Based on unique names within a given asset group.
// Assets are game data like meshes and textures that are
// either read from disk or created algorithmically.
type asset interface {
	aid() aid      // Data type and name combined.
	label() string // Identifier unique with data type.
}

// ============================================================================

// aid is a unique asset identifier.
// Asset identifiers are hashes generated from an asset name and type.
type aid uint64

// kind returns the type of asset data for this aid.
func (a aid) kind() uint32 { return uint32(a & math.MaxUint8) }

// Asset types used in aid and aid.kind.
const (
	fnt        = iota // font
	shd               // shader
	mat               // material
	msh               // mesh
	tex               // texture
	snd               // sound
	anm               // animation
	assetTypes        // end of asset types.
)

// =============================================================================
// asset utility methods.

// assetID produces a unique asset identifier using for the given
// asset type t, and asset name. Keep as many stringHash bits as possible
// to avoid collisions.
func assetID(t int, name string) aid { return aid(t) + aid(stringHash(name))<<8 }

var crcTable = crc64.MakeTable(rand.Uint64())

// stringHash turns a string into a number.
func stringHash(s string) (hash uint64) {
	return crc64.Checksum([]byte(s), crcTable)
}
