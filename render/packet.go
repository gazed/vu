// Copyright © 2024 Galvanized Logic Inc.

package render

import (
	"github.com/gazed/vu/load"
)

// Packet holds the GPU references, shader uniforms, and the
// model-view-projection transforms needed to draw a model in
// a single draw call.
type Packet struct {
	ShaderID   uint16   // GPU shader reference
	MeshID     uint32   // GPU mesh reference.
	TextureIDs []uint32 // GPU texture references.

	// packet (model) uniform data.
	Data map[load.PacketUniform][]byte

	// used to draw instanced meshes.
	IsInstanced   bool   // true for instanced models.
	InstanceID    uint32 // GPU instance data reference.
	InstanceCount uint32 // instance count for instanced models.

	// Rendering hints.
	Tag    uint32 // Application tag (entity ID) for debugging.
	Bucket uint64 // Used to sort packets. Lower buckets rendered first.
}

// Reset clears old draw data so the draw call can be reused.
func (p *Packet) Reset() {
	p.ShaderID = 0                  // default shader
	p.MeshID = 0                    // default mesh
	p.TextureIDs = p.TextureIDs[:0] // reset, keeping memory
	p.Tag = 0                       //
	p.Bucket = 0                    //
	p.IsInstanced = false           //
	p.InstanceID = 0                //
	p.InstanceCount = 0             //

	// reset the uniform data.
	for i := load.PacketUniform(0); i < load.PacketUniforms; i++ {
		d, ok := p.Data[i]
		if !ok {
			p.Data[i] = []byte{}
		} else {
			p.Data[i] = d[:0] // reset keeping memory
		}
	}
}

// Packets is a list of packets that is used to allocates render models.
// Packets are intended to be reused each render loop.
type Packets []Packet // variable number of packets.

// GetPacket returns a render.Packet from Packets. The list of packets
// is grown as needed and Packet instances are reused if available.
func (p Packets) GetPacket() (Packets, *Packet) {
	size := len(p)
	switch {
	case size == cap(p):
		p = append(p, Packet{})
		p[size].Data = map[load.PacketUniform][]byte{}
	case size < cap(p): // use previously allocated.
		p = p[:size+1]
		if p[size].Data == nil {
			p[size].Data = map[load.PacketUniform][]byte{}
		}
		p[size].Reset() // clear existing data.
	}
	return p, &p[size]
}
