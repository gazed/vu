// Copyright © 2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package load

// IQE: Inter-Quake Export format
// A text based format for 3D models that includes skeletal animation. See:
//    http://sauerbraten.org/iqm and some iqe importer implementations at,
//    https://github.com/ccxvii/asstools/blob/master/iqe.c
//    https://raw.githubusercontent.com/ccxvii/asstools/master/iqe_import.py

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"vu/math/lin"
)

// TODO Finish this once model data is available to test against.

// iqe loads one text based inter-quake model format file. Where possible it
// shares data structures and import support methods with the IQM importer.
func (l *loader) iqe(filename string) (iqd *IqData, err error) {
	iqd = &IqData{Textures: []IqTexture{}}
	var file io.ReadCloser
	if file, err = l.getResource(l.dir[mod], filename+".iqe"); err == nil {
		defer file.Close()
		return l.loadIqe(file, iqd)
	}
	return
}

// loadIqe reads a valid IQE file into an IqData structure.
func (l *loader) loadIqe(file io.ReadCloser, iqd *IqData) (*IqData, error) {
	iqd = &IqData{}
	characters, _ := ioutil.ReadAll(bufio.NewReader(file))
	lines := strings.Split(string(characters), "\n")
	poses := []*transform{}
	frameCnt, jointCnt, animCnt := 0, 0, 0        // data structure counters.
	faceCnt, faceOff, texCnt, vpCnt := 0, 0, 0, 0 // data structure counters.
	var s0 string
	var i0, i1, i2, i3 int
	var f0, f1, f2, f3, f4, f5, f6, f7, f8, f9 float32
	for _, line := range lines {
		line = strings.TrimSpace(line)
		tokens := strings.Split(line, " ")
		if len(tokens) < 2 {
			continue
		}
		key := tokens[0]
		switch key {
		case "vp": // parse 3 floats and add a vertex position.
			if _, e := fmt.Sscanf(line, "vp %f %f %f", &f0, &f1, &f2); e == nil {
				iqd.V = append(iqd.V, f0, f1, f2)
				vpCnt++
			}
		case "vt": // parse 2 floats and add a texture coordinate.
			if _, e := fmt.Sscanf(line, "vt %f %f", &f0, &f1); e == nil {
				// TODO remove kludge with better data.
				if faceOff == 0 {
					iqd.T = append(iqd.T, f0, f1)
				} else {
					iqd.T = append(iqd.T, f0, f1-1) // should not be necessary.
				}
			}
		case "vn": // parse 3 floats and add a vertex normal.
			if _, e := fmt.Sscanf(line, "vn %f %f %f", &f0, &f1, &f2); e == nil {
				iqd.N = append(iqd.N, f0, f1, f2)
			}
		case "vb": // pairs of int indexes and float weights.
			switch len(tokens) {
			case 3:
				if _, e := fmt.Sscanf(line, "vb %d %f", &i0, &f0); e == nil {
					iqd.B = append(iqd.B, byte(i0), 0, 0, 0)
					iqd.W = append(iqd.W, byte(f0*255), 0, 0, 0)
				}
			case 5:
				if _, e := fmt.Sscanf(line, "vb %d %f %d %f",
					&i0, &f0, &i1, &f1); e == nil {
					iqd.B = append(iqd.B, byte(i0), byte(i1), 0, 0)
					iqd.W = append(iqd.W, byte(f0*255), byte(f1*255), 0, 0)
				}
			case 7:
				if _, e := fmt.Sscanf(line, "vb %d %f %d %f %d %f",
					&i0, &f0, &i1, &f1, &i2, &f2); e == nil {
					iqd.B = append(iqd.B, byte(i0), byte(i1), byte(i2), 0)
					iqd.W = append(iqd.W, byte(f0*255), byte(f1*255), byte(f1*255), 0)
				}
			case 9:
				if _, e := fmt.Sscanf(line, "vb %d %f %d %f %d %f %d %f",
					&i0, &f0, &i1, &f1, &i2, &f2, &i3, &f3); e == nil {
					iqd.B = append(iqd.B, byte(i0), byte(i1), byte(i2), byte(i3))
					iqd.W = append(iqd.W, byte(f0*255), byte(f1*255), byte(f2*255), byte(f3*255))
				}
			default:
				if len(tokens) > 9 {
					log.Printf("iqe: Exceeded limit of 4 joints per vertex")
				}
			}
		case "fm": // parse 3 ints and add a triangle face.
			if _, e := fmt.Sscanf(line, "fm %d %d %d", &i0, &i1, &i2); e == nil {
				fo := faceOff
				iqd.F = append(iqd.F, uint16(i0+fo), uint16(i1+fo), uint16(i2+fo))
				faceCnt++
				iqd.Textures[texCnt-1].Fn++ // update current texture face count.
			}
		case "mesh": // reset the face offset to the current number of verticies.
			faceOff = vpCnt
		case "material": // parse 1 string for material name.
			if _, e := fmt.Sscanf(line, "material %s", &s0); e == nil {
				s0 = strings.Trim(s0, "\"")
				iqd.Textures = append(iqd.Textures, IqTexture{s0, uint32(faceCnt), 0})
				texCnt++
			}
		case "joint": // parse 1 string for joint name, 1 int for parent.
			if _, e := fmt.Sscanf(line, "joint %s %d", &s0, &i0); e == nil {
				iqd.Joints = append(iqd.Joints, int32(i0))
				jointCnt++
			}
		case "pq": // parse pose transform floats, 3 translates, 4 rotates, 3 optional scales.
			// The first poses, equal to the number of joints, are base poses for joints.
			// The remaining groups are frame poses where each group is equal to the number of joints.
			switch len(tokens) {
			case 8:
				if _, e := fmt.Sscanf(line, "pq %f %f %f %f %f %f %f",
					&f0, &f1, &f2, &f3, &f4, &f5, &f6); e == nil {
					t := &lin.V3{float64(f0), float64(f1), float64(f2)}
					r := &lin.Q{float64(f3), float64(f4), float64(f5), float64(f6)}
					s := &lin.V3{1, 1, 1}
					poses = append(poses, &transform{t, r, s})
				}
			case 11:
				if _, e := fmt.Sscanf(line, "pq %f %f %f %f %f %f %f %f %f %f",
					&f0, &f1, &f2, &f3, &f4, &f5, &f6, &f7, &f8, &f9); e == nil {
					t := &lin.V3{float64(f0), float64(f1), float64(f2)}
					r := &lin.Q{float64(f3), float64(f4), float64(f5), float64(f6)}
					s := &lin.V3{float64(f7), float64(f8), float64(f9)}
					poses = append(poses, &transform{t, r, s})
				}
			}
		case "animation": // parse 1 string for animation name.
			if _, e := fmt.Sscanf(line, "animation %s", &s0); e == nil {
				s0 = strings.Trim(s0, "\"")
				iqd.Anims = append(iqd.Anims, IqAnim{s0, uint32(frameCnt), 0})
				animCnt++
			}
		case "frame": // keep track of the current frame number.
			frameCnt++
			iqd.Anims[animCnt-1].Fn++ // update current animation frame count.
		case "vc": // ignored for now. vertex colour of 4 floats.
		case "fa": // ignored for now.
		case "pm": // ignored for now.
		case "pa": // ignored for now.
		}
	}

	// Create the matrix transforms for the animation frames.
	scr := &scratch{}
	l.createBaseFrames(iqd, poses[:jointCnt], scr)
	iqd.Frames = l.createAnimationFrames(iqd.Joints, poses[jointCnt:], scr, jointCnt, frameCnt)
	return iqd, nil
}

// createAnimationFrames generates the frame data. Process each pose in each frame.
func (l *loader) createAnimationFrames(joints []int32, poses []*transform, scr *scratch, numJoints, numFrames int) []*lin.M4 {
	frames := make([]*lin.M4, numFrames*numJoints)
	for frame := 0; frame < numFrames; frame++ {
		for pose := 0; pose < numJoints; pose++ {
			p := poses[pose]
			cnt := frame*numJoints + pose
			parent := int(joints[pose])
			frames[cnt] = l.genFrame(scr, p, pose, numJoints, parent)
		}
	}
	return frames
}

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
