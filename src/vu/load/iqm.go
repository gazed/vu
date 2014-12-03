// Copyright © 2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package load

// IQM: Inter-Quake Model format.
// A binary format for 3D models that includes skeletal animation. See:
//    http://sauerbraten.org/iqm
//    http://www.opengl.org/wiki/Skeletal_Animation
//    http://content.gpwiki.org/index.php?title=OpenGL:Tutorials:Basic_Bones_System

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"vu/math/lin"
)

// IqData is model data from IQM or IQE files. It is an intermediate format
// that is intended for populating render.Model instances.
type IqData struct {
	Name string // Data name from IQM or IQE file.

	// Mesh and Texture data create the static model.
	V        []float32   // Vertex positions.  Arranged as [][3]float32
	N        []float32   // Vertex normals.    Arranged as [][3]float32
	X        []float32   // Vertex tangents.   Arranged as [][2]float32
	T        []float32   // Vertex tex coords. Arranged as [][2]float32
	F        []uint16    // Triangle Faces.    Arranged as [][3]uint16
	Textures []IqTexture // One or more model textures

	// Optional animation data. Blend indicies indicate which vertex
	// is influenced by which joint, up to 4 joints per vertex. Blend
	// weights gives the amount of influence of a joint on a vertex.
	Anims  []IqAnim  // One or more animations.
	B      []byte    // Vertex blend indicies. Arranged as [][4]byte
	W      []byte    // Vertex blend weights.  Arranged as [][4]byte
	Joints []int32   // Joint parent information for each joint.
	Frames []*lin.M4 // Animation transforms: [NumFrames][NumJoints].
}

// IqTexture allows a model to have multiple textures. The named texture
// resource affects triangle faces from F0 to F0+FN. Expected to be used
// as part of IqData.
type IqTexture struct {
	Name   string // Name of the texture resource.
	F0, Fn uint32 // First triangle face index and number of triangle faces.
}

// IqAnim allows a model to have multiple animations. The named animation
// affects frames from F0 to F0+FN. Expected to be used as part of IqData.
type IqAnim struct {
	Name   string // Name of the animation
	F0, Fn uint32 // First frame, number of frames.
}

// =============================================================================

// iqm loads binary inter-quake model format files.
func (l *loader) iqm(filename string) (iqd *IqData, err error) {
	iqd = &IqData{Textures: []IqTexture{}}
	var file io.ReadCloser
	if file, err = l.getResource(l.dir[mod], filename+".iqm"); err == nil {
		defer file.Close()
		return l.loadIqm(file, iqd)
	}
	return
}

// loadIqm reads a valid IQM file into an IqData structure.
func (l *loader) loadIqm(file io.ReadCloser, iqd *IqData) (*IqData, error) {
	hdr := &iqmheader{}
	if err := binary.Read(file, binary.LittleEndian, hdr); err != nil {
		return iqd, fmt.Errorf("Invalid .iqm file: %s", err)
	}

	// sanity check the data.
	if !bytes.Equal(iqmMagic, hdr.Magic[:]) {
		return iqd, fmt.Errorf("Invalid .iqm header magic: %s", string(hdr.Magic[:]))
	}
	if hdr.Version != 2 {
		return iqd, fmt.Errorf("Expecting .iqm version 2, got : %d", hdr.Version)
	}
	if hdr.Filesize > (16 << 20) {
		return iqd, fmt.Errorf("Not loading .iqm files bigger than 16MB")
	}

	// Get the data into memory. Not all readers return all the data in one shot.
	bytesRead := uint32(0)
	dataSize := hdr.Filesize - iqmheaderSize
	data := make([]byte, dataSize)
	inbuff := make([]byte, dataSize)
	for bytesRead < dataSize {
		inbytes, readErr := file.Read(inbuff)
		if readErr != nil {
			return iqd, fmt.Errorf("Corrupt .iqm file")
		}
		for cnt := 0; cnt < inbytes; cnt++ {
			data[bytesRead] = inbuff[cnt]
			bytesRead += 1
		}
	}
	if bytesRead != dataSize {
		return iqd, fmt.Errorf("Invalid .iqm file")
	}
	scratch := &scratch{}
	if hdr.Num_meshes > 0 {
		if err := l.loadIqmMeshes(hdr, data, iqd, scratch); err != nil {
			return iqd, err
		}
	}
	if hdr.Num_anims > 0 {
		if err := l.loadIqmAnims(hdr, data, iqd, scratch); err != nil {
			return iqd, err
		}
	}
	return iqd, nil
}

// scratch is temporary memory for loading a single model.
// The data is initialized in loadIqmMeshes, used in loadIqmAnims
type scratch struct {
	labels           map[uint32]string // resource string identifiers.
	baseframe        []*lin.M4         // joint transform.
	inversebaseframe []*lin.M4         // inverse joint transform.
}

// loadIqmMeshes parses the vertex data from the file data into the
// IqData structure.
func (l *loader) loadIqmMeshes(hdr *iqmheader, data []byte, iqd *IqData, scr *scratch) (err error) {
	buff := bytes.NewReader(data)

	// Get all the text labels referenced by other structures. Index the labels by
	// their byte position which is how they are referenced
	buff.Seek(int64(hdr.Ofs_text-iqmheaderSize), 0)
	text := make([]byte, hdr.Num_text)
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
	buff.Seek(int64(hdr.Ofs_vertexarrays-iqmheaderSize), 0)
	for cnt := 0; cnt < int(hdr.Num_vertexarrays); cnt++ {
		if err = binary.Read(buff, binary.LittleEndian, va); err != nil {
			return fmt.Errorf("Invalid .iqm file: %s", err)
		}
		switch va.Type {
		case iQM_POSITION:
			iqd.V = make([]float32, va.Size*hdr.Num_vertexes)
			err = l.readVertexData(data, va, iQM_FLOAT, 3, iqd.V)
		case iQM_NORMAL:
			iqd.N = make([]float32, va.Size*hdr.Num_vertexes)
			err = l.readVertexData(data, va, iQM_FLOAT, 3, iqd.N)
		case iQM_TANGENT:
			iqd.X = make([]float32, va.Size*hdr.Num_vertexes)
			err = l.readVertexData(data, va, iQM_FLOAT, 4, iqd.X)
		case iQM_TEXCOORD:
			iqd.T = make([]float32, va.Size*hdr.Num_vertexes)
			err = l.readVertexData(data, va, iQM_FLOAT, 2, iqd.T)

		// Indexes and weights are sent to the GPU as bytes in order
		// to reduce the amount of data transferred.
		case iQM_BLENDINDEXES:
			iqd.B = make([]byte, va.Size*hdr.Num_vertexes)
			err = l.readVertexData(data, va, iQM_UBYTE, 4, iqd.B)
		case iQM_BLENDWEIGHTS:
			iqd.W = make([]byte, va.Size*hdr.Num_vertexes)
			err = l.readVertexData(data, va, iQM_UBYTE, 4, iqd.W)
			// Note: blend weights are normalized to 0-1 floats on transfer to the GPU.
		}
		if err != nil {
			return err
		}
	}

	// Get the triangle face data.
	buff.Seek(int64(hdr.Ofs_triangles-iqmheaderSize), 0)
	faces := make([]uint32, 3*hdr.Num_triangles)
	if err = binary.Read(buff, binary.LittleEndian, faces); err != nil {
		return fmt.Errorf("Invalid .iqm triangles %s", err)
	}
	iqd.F = make([]uint16, 3*hdr.Num_triangles)
	for cnt := 0; cnt < len(faces); cnt++ {
		iqd.F[cnt] = uint16(faces[cnt])
	}

	// Multiple meshes mean means that multiple textures are used for this model.
	msh := &iqmmesh{}
	buff.Seek(int64(hdr.Ofs_meshes-iqmheaderSize), 0)
	for cnt := 0; cnt < int(hdr.Num_meshes); cnt++ {
		if err = binary.Read(buff, binary.LittleEndian, msh); err != nil {
			return fmt.Errorf("Invalid .iqm file: %s", err)
		}
		itex := IqTexture{}
		itex.Name = scr.labels[msh.Material] // Name of the mesh resource.
		itex.F0, itex.Fn = msh.First_triangle, msh.Num_triangles
		iqd.Textures = append(iqd.Textures, itex)
	}
	return nil
}

// readVertexData reads and validates a set of vertex data from an IQM file.
func (l *loader) readVertexData(data []byte, va *iqmvertexarray, dtype, dspan uint32, outData interface{}) (err error) {
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

// loadIqmAnims parses the animation data from the file data into the
// IqData structure.
func (l *loader) loadIqmAnims(hdr *iqmheader, data []byte, iqd *IqData, scr *scratch) (err error) {
	if hdr.Num_poses != hdr.Num_joints {
		return fmt.Errorf("Invalid .iqm joints must equal poses", hdr.Num_poses, hdr.Num_joints)
	}
	buff := bytes.NewReader(data)

	// Read base pose joint data.
	jnts := make([]iqmjoint, hdr.Num_joints)
	buff.Seek(int64(hdr.Ofs_joints-iqmheaderSize), 0)
	if err = binary.Read(buff, binary.LittleEndian, jnts); err != nil {
		return fmt.Errorf("Invalid .iqm file: %s", err)
	}
	iqd.Joints = make([]int32, hdr.Num_joints)

	// process the joint base transforms using an intermediate form.
	basePoses := []*transform{}
	for cnt, j := range jnts {
		iqd.Joints[cnt] = j.Parent // save the joint parent data

		// put the pose data into a transform ready structure.
		t := &lin.V3{float64(j.Translate[0]), float64(j.Translate[1]), float64(j.Translate[2])}
		r := &lin.Q{float64(j.Rotate[0]), float64(j.Rotate[1]), float64(j.Rotate[2]), float64(j.Rotate[3])}
		s := &lin.V3{float64(j.Scale[0]), float64(j.Scale[1]), float64(j.Scale[2])}
		basePoses = append(basePoses, &transform{t, r, s})
	}
	l.createBaseFrames(iqd, basePoses, scr)

	// Get the per frame pose data.
	buff.Seek(int64(hdr.Ofs_poses-iqmheaderSize), 0)
	poses := make([]iqmpose, hdr.Num_poses)
	if err = binary.Read(buff, binary.LittleEndian, poses); err != nil {
		return fmt.Errorf("Invalid .iqm poses %s", err)
	}

	// Get the animation data.
	buff.Seek(int64(hdr.Ofs_anims-iqmheaderSize), 0)
	animData := make([]iqmanim, hdr.Num_anims)
	if err = binary.Read(buff, binary.LittleEndian, animData); err != nil {
		return fmt.Errorf("Invalid .iqm animations %s", err)
	}
	for _, adata := range animData {
		anim := IqAnim{}
		anim.Name = scr.labels[adata.Name]
		anim.F0 = adata.First_frame
		anim.Fn = adata.Num_frames
		iqd.Anims = append(iqd.Anims, anim)
	}

	// Get the animation frames.
	buff.Seek(int64(hdr.Ofs_frames-iqmheaderSize), 0)
	frameData := make([]uint16, hdr.Num_frames*hdr.Num_framechannels)
	if err = binary.Read(buff, binary.LittleEndian, frameData); err != nil {
		return fmt.Errorf("Invalid .iqm frames %s", err)
	}

	// Generate the final animation frames from base poses and frame poses.
	pt := &transform{&lin.V3{}, &lin.Q{}, &lin.V3{}} // pose transform
	iqd.Frames = make([]*lin.M4, hdr.Num_frames*hdr.Num_poses)
	fcnt := 0
	for frame := 0; frame < int(hdr.Num_frames); frame++ {
		for pose := 0; pose < int(hdr.Num_poses); pose++ {
			p := poses[pose]
			pt.t.X, fcnt = l.getPoseChannel(&p, frameData, fcnt, 0, 0x01)
			pt.t.Y, fcnt = l.getPoseChannel(&p, frameData, fcnt, 1, 0x02)
			pt.t.Z, fcnt = l.getPoseChannel(&p, frameData, fcnt, 2, 0x04)
			pt.q.X, fcnt = l.getPoseChannel(&p, frameData, fcnt, 3, 0x08)
			pt.q.Y, fcnt = l.getPoseChannel(&p, frameData, fcnt, 4, 0x10)
			pt.q.Z, fcnt = l.getPoseChannel(&p, frameData, fcnt, 5, 0x20)
			pt.q.W, fcnt = l.getPoseChannel(&p, frameData, fcnt, 6, 0x40)
			pt.s.X, fcnt = l.getPoseChannel(&p, frameData, fcnt, 7, 0x80)
			pt.s.Y, fcnt = l.getPoseChannel(&p, frameData, fcnt, 8, 0x100)
			pt.s.Z, fcnt = l.getPoseChannel(&p, frameData, fcnt, 9, 0x200)

			// Combine all the data into a animation ready frame transform matrix.
			cnt := frame*int(hdr.Num_poses) + pose
			iqd.Frames[cnt] = l.genFrame(scr, pt, pose, int(hdr.Num_poses), int(p.Parent))
		}
	}
	return nil
}

// getPoseChannel is a helper method that builds per-frame pose animation
// transform from a compressed/sparse format.
func (l *loader) getPoseChannel(p *iqmpose, fdata []uint16, fcnt, index int, mask uint32) (float64, int) {
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
func (l *loader) createBaseFrames(iqd *IqData, poses []*transform, scr *scratch) {
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
		parent := iqd.Joints[cnt]
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
func (l *loader) genFrame(scr *scratch, pt *transform, pcnt, numPoses, parent int) *lin.M4 {
	pt.q.Unit()
	m4 := lin.NewM4().SetQ(pt.q)
	m4.Transpose(m4).ScaleSM(pt.s.X, pt.s.Y, pt.s.Z)       // apply scale before rotation.
	m4.Wx, m4.Wy, m4.Wz, m4.Ww = pt.t.X, pt.t.Y, pt.t.Z, 1 // translation added in, not multiplied.
	if parent >= 0 {
		// parentBasePose * childPose * childInverseBasePose
		return m4.Mult(scr.inversebaseframe[pcnt], m4).Mult(m4, scr.baseframe[parent])
	} else {
		// childPose * childInverseBasePose
		return m4.Mult(scr.inversebaseframe[pcnt], m4)
	}
}

// =============================================================================
// The binary structures for an IQM file is from sauerbraten.org/iqm/iqm.txt

// iqmheader provides the indexes to the remaining data and is at the begining
// of the iqm file.
type iqmheader struct {
	Magic                                                 [16]byte // the string "INTERQUAKEMODEL\0".
	Version                                               uint32   // Must be version 2.
	Filesize                                              uint32   // Total bytes in the file.
	Flags                                                 uint32
	Num_text, Ofs_text                                    uint32
	Num_meshes, Ofs_meshes                                uint32 // Number and data offset of meshes.
	Num_vertexarrays, Num_vertexes, Ofs_vertexarrays      uint32
	Num_triangles, Ofs_triangles, Ofs_adjacency           uint32
	Num_joints, Ofs_joints                                uint32
	Num_poses, Ofs_poses                                  uint32
	Num_anims, Ofs_anims                                  uint32
	Num_frames, Num_framechannels, Ofs_frames, Ofs_bounds uint32
	Num_comment, Ofs_comment                              uint32
	Num_extensions, Ofs_extensions                        uint32 // A linked list, not a contiguous array.
}

// iqmheaderSize is (16 bytes)+(27 fields)*(4 bytes) = 124 bytes.
var iqmheaderSize = uint32(124)

// iqmMagic is the first 16 bytes of a valid IQM file.
var iqmMagic = []byte{'I', 'N', 'T', 'E', 'R', 'Q', 'U', 'A', 'K', 'E', 'M', 'O', 'D', 'E', 'L', 0}

type iqmmesh struct {
	Name                          uint32 // unique name for the mesh, if desired
	Material                      uint32 // set to a name of a non-unique material or texture
	First_vertex, Num_vertexes    uint32
	First_triangle, Num_triangles uint32
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
	Name                    uint32
	First_frame, Num_frames uint32
	Framerate               float32
	Flags                   uint32
}

// all vertex array entries must ordered as defined below, if present
// i.e. position comes before normal comes before ... comes before custom
// where a format and size is given, this means models intended for portable use should use these
// an IQM implementation is not required to honor any other format/size than those recommended
// however, it may support other format/size combinations for these types if it desires
const ( // vertex array type
	iQM_POSITION     = 0 // float, 3
	iQM_TEXCOORD     = 1 // float, 2
	iQM_NORMAL       = 2 // float, 3
	iQM_TANGENT      = 3 // float, 4
	iQM_BLENDINDEXES = 4 // ubyte, 4
	iQM_BLENDWEIGHTS = 5 // ubyte, 4
	iQM_COLOR        = 6 // ubyte, 4

	// all values up to IQM_CUSTOM are reserved for future use
	// any value >= IQM_CUSTOM is interpreted as CUSTOM type
	// the value then defines an offset into the string table, where offset = value - IQM_CUSTOM
	// this must be a valid string naming the type
	iQM_CUSTOM = 0x10
)

const ( // vertex array format
	iQM_BYTE   = 0
	iQM_UBYTE  = 1
	iQM_SHORT  = 2
	iQM_USHORT = 3
	iQM_INT    = 4
	iQM_UINT   = 5
	iQM_HALF   = 6
	iQM_FLOAT  = 7
	iQM_DOUBLE = 8
)
