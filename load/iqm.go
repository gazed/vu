// Copyright © 2014-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package load

// IQM: Inter-Quake Model format.
// A binary format for 3D models that includes skeletal animation:
//    http://www.opengl.org/wiki/Skeletal_Animation
//    http://content.gpwiki.org/index.php?title=OpenGL:Tutorials:Basic_Bones_System

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/gazed/vu/math/lin"
)

// Iqm loads an Inter-Quake model IQM file into ModData.
// IQM is A binary format for 3D models that includes skeletal animation.
// See: http://sauerbraten.org/iqm. This loader has been tested against
// a subset of the full specification.
// The Reader r is expected to be opened and closed by the caller.
func Iqm(r io.Reader, d *ModData) error {
	hdr := &iqmheader{}
	if err := binary.Read(r, binary.LittleEndian, hdr); err != nil {
		return fmt.Errorf("Invalid .iqm file: %s", err)
	}

	// sanity check the data.
	if !bytes.Equal(iqmMagic, hdr.Magic[:]) {
		return fmt.Errorf("Invalid .iqm header magic: %s", string(hdr.Magic[:]))
	}
	if hdr.Version != 2 {
		return fmt.Errorf("Expecting .iqm version 2, got : %d", hdr.Version)
	}
	if hdr.Filesize > (16 << 20) {
		return fmt.Errorf("Not loading .iqm files bigger than 16MB")
	}

	// Get the data into memory. Not all readers return all the data in one shot.
	bytesRead := uint32(0)
	dataSize := hdr.Filesize - iqmheaderSize
	data := make([]byte, dataSize)
	inbuff := make([]byte, dataSize)
	for bytesRead < dataSize {
		inbytes, readErr := r.Read(inbuff)
		if readErr != nil {
			return fmt.Errorf("Corrupt .iqm file")
		}
		for cnt := 0; cnt < inbytes; cnt++ {
			data[bytesRead] = inbuff[cnt]
			bytesRead++
		}
	}
	if bytesRead != dataSize {
		return fmt.Errorf("Invalid .iqm file")
	}
	scratch := &scratch{}
	if hdr.NumMeshes > 0 {
		if err := loadIqmMeshes(hdr, data, d, scratch); err != nil {
			return err
		}
	}
	if hdr.NumAnims > 0 {
		if err := loadIqmAnims(hdr, data, d, scratch); err != nil {
			return err
		}
	}
	return nil
}

// public inteface
// =============================================================================
// internal implementation for loading IQM files.

// scratch is temporary memory for loading a single model.
// The data is initialized in loadIqmMeshes, used in loadIqmAnims
type scratch struct {
	labels           map[uint32]string // resource string identifiers.
	baseframe        []*lin.M4         // joint transform.
	inversebaseframe []*lin.M4         // inverse joint transform.
}

// loadIqmMeshes parses the vertex data from the file data
// into the IqData structure.
func loadIqmMeshes(hdr *iqmheader, data []byte, mod *ModData, scr *scratch) (err error) {
	buff := bytes.NewReader(data)

	// Get all the text labels referenced by other structures. Index the
	// labels by their byte position which is how they are referenced
	buff.Seek(int64(hdr.OfsText-iqmheaderSize), 0)
	text := make([]byte, hdr.NumText)
	if err = binary.Read(buff, binary.LittleEndian, text); err != nil {
		return fmt.Errorf("Invalid .iqm text block %s", err)
	}
	scr.labels = map[uint32]string{}
	last := uint32(0)
	for cnt := 0; cnt < len(text); cnt++ {
		if text[cnt] == 0 {
			if cnt == int(last) {
				scr.labels[last] = ""
			} else {
				scr.labels[last] = string(text[last:cnt])
			}
			last = uint32(cnt) + 1
		}
	}

	// Get the vertex data.
	va := &iqmvertexarray{}
	buff.Seek(int64(hdr.OfsVertexArrays-iqmheaderSize), 0)
	for cnt := 0; cnt < int(hdr.NumVertexArrays); cnt++ {
		if err = binary.Read(buff, binary.LittleEndian, va); err != nil {
			return fmt.Errorf("Invalid .iqm file: %s", err)
		}
		switch va.Type {
		case iQMPOSITION:
			mod.V = make([]float32, va.Size*hdr.NumVertexes)
			err = readVertexData(data, va, iQMFLOAT, 3, mod.V)
		case iQMNORMAL:
			mod.N = make([]float32, va.Size*hdr.NumVertexes)
			err = readVertexData(data, va, iQMFLOAT, 3, mod.N)
		case iQMTANGENT:
			mod.X = make([]float32, va.Size*hdr.NumVertexes)
			err = readVertexData(data, va, iQMFLOAT, 4, mod.X)
		case iQMTEXCOORD:
			mod.T = make([]float32, va.Size*hdr.NumVertexes)
			err = readVertexData(data, va, iQMFLOAT, 2, mod.T)

		// Indexes and weights are sent to the GPU as bytes in order
		// to reduce the amount of data transferred.
		case iQMBLENDINDEXES:
			mod.Blends = make([]byte, va.Size*hdr.NumVertexes)
			err = readVertexData(data, va, iQMUBYTE, 4, mod.Blends)
		case iQMBLENDWEIGHTS:
			mod.Weights = make([]byte, va.Size*hdr.NumVertexes)
			err = readVertexData(data, va, iQMUBYTE, 4, mod.Weights)
			// Note: blend weights are normalized to 0-1 floats on GPU transfer.
		}
		if err != nil {
			return err
		}
	}

	// Get the triangle face data.
	buff.Seek(int64(hdr.OfsTriangles-iqmheaderSize), 0)
	faces := make([]uint32, 3*hdr.NumTriangles)
	if err = binary.Read(buff, binary.LittleEndian, faces); err != nil {
		return fmt.Errorf("Invalid .iqm triangles %s", err)
	}
	mod.F = make([]uint16, 3*hdr.NumTriangles)
	for cnt := 0; cnt < len(faces); cnt++ {
		mod.F[cnt] = uint16(faces[cnt])
	}

	// Multiple meshes mean means that multiple textures are used for this model.
	msh := &iqmmesh{}
	buff.Seek(int64(hdr.OfsMeshes-iqmheaderSize), 0)
	for cnt := 0; cnt < int(hdr.NumMeshes); cnt++ {
		if err = binary.Read(buff, binary.LittleEndian, msh); err != nil {
			return fmt.Errorf("Invalid .iqm file: %s", err)
		}
		itex := TexMap{}
		itex.Name = scr.labels[msh.Material] // Name of the mesh resource.
		itex.F0, itex.Fn = msh.FirstTriangle, msh.NumTriangles
		mod.TMap = append(mod.TMap, itex)
	}
	return nil
}

// readVertexData reads and validates a set of vertex data from an IQM file.
func readVertexData(data []byte, va *iqmvertexarray, dtype, dspan uint32, outData interface{}) (err error) {
	if va.Format != dtype || va.Size != dspan {
		return fmt.Errorf("Invalid .iqm vertex data array")
	}
	vbuff := bytes.NewReader(data)
	vbuff.Seek(int64(va.Offset-iqmheaderSize), 0)
	if err = binary.Read(vbuff, binary.LittleEndian, outData); err != nil {
		return fmt.Errorf("Invalid .iqm vertex data %s", err)
	}
	return nil
}

// loadIqmAnims parses the animation data from the file data.
func loadIqmAnims(hdr *iqmheader, data []byte, mod *ModData, scr *scratch) (err error) {
	if hdr.NumPoses != hdr.NumJoints {
		return fmt.Errorf("Invalid .iqm joints %d must equal poses %d", hdr.NumJoints, hdr.NumPoses)
	}
	buff := bytes.NewReader(data)

	// Read base pose joint data.
	jnts := make([]iqmjoint, hdr.NumJoints)
	buff.Seek(int64(hdr.OfsJoints-iqmheaderSize), 0)
	if err = binary.Read(buff, binary.LittleEndian, jnts); err != nil {
		return fmt.Errorf("Invalid .iqm file: %s", err)
	}
	mod.Joints = make([]int32, hdr.NumJoints)

	// process the joint base transforms using an intermediate form.
	basePoses := []*transform{}
	for cnt, j := range jnts {
		mod.Joints[cnt] = j.Parent // save the joint parent data

		// FUTURE: use joint names as attachment points where the convention
		//         is to have "attachment" in the joint name, and not have
		//         the joint affect any verticies.

		// put the pose data into a transform ready structure.
		t := &lin.V3{X: float64(j.Translate[0]), Y: float64(j.Translate[1]), Z: float64(j.Translate[2])}
		r := &lin.Q{X: float64(j.Rotate[0]), Y: float64(j.Rotate[1]), Z: float64(j.Rotate[2]), W: float64(j.Rotate[3])}
		s := &lin.V3{X: float64(j.Scale[0]), Y: float64(j.Scale[1]), Z: float64(j.Scale[2])}
		basePoses = append(basePoses, &transform{t, r, s})
	}
	createBaseFrames(mod, basePoses, scr)

	// Get the per frame pose data.
	buff.Seek(int64(hdr.OfsPoses-iqmheaderSize), 0)
	poses := make([]iqmpose, hdr.NumPoses)
	if err = binary.Read(buff, binary.LittleEndian, poses); err != nil {
		return fmt.Errorf("Invalid .iqm poses %s", err)
	}

	// Get the animation data.
	buff.Seek(int64(hdr.OfsAnims-iqmheaderSize), 0)
	animData := make([]iqmanim, hdr.NumAnims)
	if err = binary.Read(buff, binary.LittleEndian, animData); err != nil {
		return fmt.Errorf("Invalid .iqm animations %s", err)
	}
	for _, adata := range animData {
		anim := Movement{}
		anim.Name = scr.labels[adata.Name]
		anim.F0 = adata.FirstFrame
		anim.Fn = adata.NumFrames
		anim.Rate = adata.Framerate
		mod.Movements = append(mod.Movements, anim)
	}

	// Get the animation frames.
	buff.Seek(int64(hdr.OfsFrames-iqmheaderSize), 0)
	frameData := make([]uint16, hdr.NumFrames*hdr.NumFrameChannels)
	if err = binary.Read(buff, binary.LittleEndian, frameData); err != nil {
		return fmt.Errorf("Invalid .iqm frames %s", err)
	}

	// Generate the final animation frames from base poses and frame poses.
	pt := &transform{&lin.V3{}, &lin.Q{}, &lin.V3{}} // pose transform
	mod.Frames = make([]*lin.M4, hdr.NumFrames*hdr.NumPoses)
	fcnt := 0
	for frame := 0; frame < int(hdr.NumFrames); frame++ {
		for pose := 0; pose < int(hdr.NumPoses); pose++ {
			p := poses[pose]
			pt.t.X, fcnt = getPoseChannel(&p, frameData, fcnt, 0, 0x01)
			pt.t.Y, fcnt = getPoseChannel(&p, frameData, fcnt, 1, 0x02)
			pt.t.Z, fcnt = getPoseChannel(&p, frameData, fcnt, 2, 0x04)
			pt.q.X, fcnt = getPoseChannel(&p, frameData, fcnt, 3, 0x08)
			pt.q.Y, fcnt = getPoseChannel(&p, frameData, fcnt, 4, 0x10)
			pt.q.Z, fcnt = getPoseChannel(&p, frameData, fcnt, 5, 0x20)
			pt.q.W, fcnt = getPoseChannel(&p, frameData, fcnt, 6, 0x40)
			pt.s.X, fcnt = getPoseChannel(&p, frameData, fcnt, 7, 0x80)
			pt.s.Y, fcnt = getPoseChannel(&p, frameData, fcnt, 8, 0x100)
			pt.s.Z, fcnt = getPoseChannel(&p, frameData, fcnt, 9, 0x200)

			// Combine all the data into a animation ready frame transform matrix.
			cnt := frame*int(hdr.NumPoses) + pose
			mod.Frames[cnt] = genFrame(scr, pt, pose, int(hdr.NumPoses), int(p.Parent))
		}
	}
	return nil
}

// getPoseChannel is a helper method that builds per-frame pose animation
// transform from a compressed/sparse format.
func getPoseChannel(p *iqmpose, fdata []uint16, fcnt, index int, mask uint32) (float64, int) {
	channel := float64(p.Channeloffset[index])
	if p.Channelmask&mask == mask {
		channel += float64(fdata[fcnt]) * float64(p.Channelscale[index])
		fcnt++
	}
	return channel, fcnt
}

// createBaseFrames constructs the joint transform base-pose matricies.
// These are temporary structures used later in genFrame to prepare the
// per-frame animation data.
func createBaseFrames(mod *ModData, poses []*transform, scr *scratch) {
	i3 := lin.NewM3()                                   // scratch
	vx, vy, vz := lin.NewV3(), lin.NewV3(), lin.NewV3() // scratch
	numJoints := len(poses)
	scr.baseframe = make([]*lin.M4, numJoints)        // joint transforms.
	scr.inversebaseframe = make([]*lin.M4, numJoints) // inverse transforms.
	for cnt := 0; cnt < numJoints; cnt++ {
		j := poses[cnt]

		// Get the joint transform.
		j.q.Unit() // ensure unit quaternion.
		m4 := lin.NewM4().SetQ(j.q)
		m4.Transpose(m4).ScaleSM(j.s.X, j.s.Y, j.s.Z)       // apply scale before rotation.
		m4.Wx, m4.Wy, m4.Wz, m4.Ww = j.t.X, j.t.Y, j.t.Z, 1 // translation added in, not multiplied.
		scr.baseframe[cnt] = m4

		// invert the joint transform for frame generation later on.
		i3.Inv(i3.SetM4(m4))
		itx := -vx.SetS(i3.Xx, i3.Yx, i3.Zx).Dot(j.t)
		ity := -vy.SetS(i3.Xy, i3.Yy, i3.Zy).Dot(j.t)
		itz := -vz.SetS(i3.Xz, i3.Yz, i3.Zz).Dot(j.t)
		i4 := lin.NewM4()
		i4.Xx, i4.Xy, i4.Xz, i4.Xw = i3.Xx, i3.Xy, i3.Xz, 0
		i4.Yx, i4.Yy, i4.Yz, i4.Yw = i3.Yx, i3.Yy, i3.Yz, 0
		i4.Zx, i4.Zy, i4.Zz, i4.Zw = i3.Zx, i3.Zy, i3.Zz, 0
		i4.Wx, i4.Wy, i4.Wz, i4.Ww = itx, ity, itz, 1
		scr.inversebaseframe[cnt] = i4

		// Combine the joint transforms and inverse transform with the parent transform.
		parent := mod.Joints[cnt]
		if parent >= 0 {
			// childBasePose * parentBasePose
			scr.baseframe[cnt].Mult(scr.baseframe[cnt], scr.baseframe[parent])

			// childInverseBasePose * parentInverseBasePose
			scr.inversebaseframe[cnt].Mult(scr.inversebaseframe[parent], scr.inversebaseframe[cnt])
		}
	}
}

// Concatenate each pose with the inverse base pose to avoid doing this at animation time.
// If the joint has a parent, then it needs to be pre-concatenated with its parent's base pose.
// Thus it all negates at animation time like so:
//    (parentPose * parentInverseBasePose) * (parentBasePose * childPose * childInverseBasePose) =>
//    parentPose * (parentInverseBasePose * parentBasePose) * childPose * childInverseBasePose =>
//    parentPose * childPose * childInverseBasePose
func genFrame(scr *scratch, pt *transform, pcnt, numPoses, parent int) *lin.M4 {
	pt.q.Unit()
	m4 := lin.NewM4().SetQ(pt.q)
	m4.Transpose(m4).ScaleSM(pt.s.X, pt.s.Y, pt.s.Z)       // apply scale before rotation.
	m4.Wx, m4.Wy, m4.Wz, m4.Ww = pt.t.X, pt.t.Y, pt.t.Z, 1 // translation added in, not multiplied.
	if parent >= 0 {
		// parentBasePose * childPose * childInverseBasePose
		return m4.Mult(scr.inversebaseframe[pcnt], m4).Mult(m4, scr.baseframe[parent])
	}
	// childPose * childInverseBasePose
	return m4.Mult(scr.inversebaseframe[pcnt], m4)
}

// =============================================================================
// The binary structures for an IQM file is from sauerbraten.org/iqm/iqm.txt

// iqmheader provides the indexes to the remaining data and is at the beginning
// of the iqm file.
type iqmheader struct {
	Magic                                             [16]byte // the string "INTERQUAKEMODEL\0".
	Version                                           uint32   // Must be version 2.
	Filesize                                          uint32   // Total bytes in the file.
	Flags                                             uint32
	NumText, OfsText                                  uint32
	NumMeshes, OfsMeshes                              uint32 // Number and data offset of meshes.
	NumVertexArrays, NumVertexes, OfsVertexArrays     uint32
	NumTriangles, OfsTriangles, OfsAdjacency          uint32
	NumJoints, OfsJoints                              uint32
	NumPoses, OfsPoses                                uint32
	NumAnims, OfsAnims                                uint32
	NumFrames, NumFrameChannels, OfsFrames, OfsBounds uint32
	NumComment, OfsComment                            uint32
	NumExtensions, OfsExtensions                      uint32 // A linked list, not a contiguous array.
}

// iqmheaderSize is (16 bytes)+(27 fields)*(4 bytes) = 124 bytes.
var iqmheaderSize = uint32(124)

// iqmMagic is the first 16 bytes of a valid IQM file.
var iqmMagic = []byte{'I', 'N', 'T', 'E', 'R', 'Q', 'U', 'A', 'K', 'E', 'M', 'O', 'D', 'E', 'L', 0}

type iqmmesh struct {
	Name                        uint32 // unique name for the mesh, if desired
	Material                    uint32 // set to a name of a non-unique material or texture
	FirstVertex, NumVertexes    uint32
	FirstTriangle, NumTriangles uint32
}

type iqmvertexarray struct {
	Type   uint32 // type or custom name
	Flags  uint32
	Format uint32 // component format
	Size   uint32 // number of components
	Offset uint32 // offset to array of tightly packed components, with num_vertexes * size total entries
	// offset must be aligned to max(sizeof(format), 4)
}

type iqmtriangle struct {
	Vertex [3]uint32
}

// translate is translation <Tx, Ty, Tz>, and
// rotate is quaternion rotation <Qx, Qy, Qz, Qw> (in relative/parent local space)
// scale is pre-scaling <Sx, Sy, Sz>
// output = (input*scale)*rotation + translation
type iqmjoint struct {
	Name      uint32
	Parent    int32 // parent < 0 means this is a root bone
	Translate [3]float32
	Rotate    [4]float32
	Scale     [3]float32
}

// channels 0..2 are translation <Tx, Ty, Tz> and
// channels 3..6 are quaternion rotation <Qx, Qy, Qz, Qw> (in relative/parent local space)
// channels 7..9 are scale <Sx, Sy, Sz>
// output = (input*scale)*rotation + translation
type iqmpose struct {
	Parent        int32  // parent < 0 means this is a root bone
	Channelmask   uint32 // mask of which 10 channels are present for this joint pose
	Channeloffset [10]float32
	Channelscale  [10]float32
}

type iqmanim struct {
	Name                  uint32
	FirstFrame, NumFrames uint32
	Framerate             float32
	Flags                 uint32
}

// all vertex array entries must ordered as defined below, if present
// i.e. position comes before normal comes before ... comes before custom
// where a format and size is given, this means models intended for portable use should use these
// an IQM implementation is not required to honor any other format/size than those recommended
// however, it may support other format/size combinations for these types if it desires
const ( // vertex array type
	iQMPOSITION     = 0 // float, 3
	iQMTEXCOORD     = 1 // float, 2
	iQMNORMAL       = 2 // float, 3
	iQMTANGENT      = 3 // float, 4
	iQMBLENDINDEXES = 4 // ubyte, 4
	iQMBLENDWEIGHTS = 5 // ubyte, 4
	iQMCOLOR        = 6 // ubyte, 4

	// all values up to IQM_CUSTOM are reserved for future use
	// any value >= IQM_CUSTOM is interpreted as CUSTOM type
	// the value then defines an offset into the string table, where offset = value - IQM_CUSTOM
	// this must be a valid string naming the type
	iQMCUSTOM = 0x10
)

const ( // vertex array format
	iQMBYTE   = 0
	iQMUBYTE  = 1
	iQMSHORT  = 2
	iQMUSHORT = 3
	iQMINT    = 4
	iQMUINT   = 5
	iQMHALF   = 6
	iQMFLOAT  = 7
	iQMDOUBLE = 8
)

// transform is a temporary structure to help read pose information from
// IQE/IQM files. It consists of:
//    translation <Tx, Ty, Tz>,
//    quaternion rotation <Qx, Qy, Qz, Qw> (in relative/parent local space)
//    scale is pre-scaling <Sx, Sy, Sz>
// It is combined as follows: (input*scale)*rotation + translation
type transform struct {
	t *lin.V3 // translate.
	q *lin.Q  // rotation (orientation) quaternion.
	s *lin.V3 // scale.
}
