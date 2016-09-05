// Copyright © 2013-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package load

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/gazed/vu/math/lin"
)

// Obj loads a Wavefront OBJ file containing one or more mesh descriptions.
// A Wavefront OBJ file is a text representation of one or more 3D models.
// This loader supports a limited subset of the full specification. It is
// specifically looking for a single object triangle mesh with normals.
//    https://en.wikipedia.org/wiki/Wavefront_.obj_file#File_format
//    http://www.martinreddy.net/gfx/3d/OBJ.spec
// The Reader r is expected to be opened and closed by the caller.
// A successful import overwrites the data in ObjData.
func Obj(r io.Reader, d *MshData) error {
	objs := obj2Strings(r)
	if len(objs) <= 0 {
		return fmt.Errorf("No objects in .obj file")
	}
	odata := &objData{}
	faces, derr := obj2Data(objs[0].lines, odata)
	if derr != nil {
		return fmt.Errorf("obj2Data %s", derr)
	}
	return obj2MshData(objs[0].name, odata, faces, d)
}

// public inteface
// =============================================================================
// internal implementation for loading OBJ files.

// objStrings is an intermediate data structure used in parsing.
type objStrings struct {
	name  string
	lines []string
}

// objData is an intermediate data structure used in parsing.
// Each .obj file keeps a global count of the data below.  This is referenced
// from the face data.
type objData struct {
	v []dataPoint // vertices
	n []dataPoint // normals
	t []uvPoint   // texture coordinates
}

// dataPoint is an internal structure for passing vertices or normals.
type dataPoint struct {
	x, y, z float32
}

// uvPoint is an internal structure for passing texture coordinates.
type uvPoint struct {
	u, v float32
}

// face is an internal structure for passing face indexes.
type face struct {
	s []string // each point is a "x/y/z" value.
}

// obj2Strings reads in all the file data grouped by object name. This is needed
// because a single wavefront file can hold many objects. Separating the objects
// makes parsing easier.
func obj2Strings(file io.Reader) (objs []*objStrings) {
	objs = []*objStrings{}
	name := ""
	var curr *objStrings
	reader := bufio.NewReader(file)
	line, e1 := reader.ReadString('\n')
	for ; e1 == nil; line, e1 = reader.ReadString('\n') {
		line = strings.TrimSpace(line)
		tokens := strings.Split(line, " ")
		if len(tokens) == 2 && tokens[0] == "o" {
			name = strings.TrimSpace(tokens[1])
			curr = &objStrings{name, []string{}}
			objs = append(objs, curr)
		} else if len(name) > 0 {
			curr.lines = append(curr.lines, strings.TrimSpace(line))
		}
	}
	return
}

// obj2Data turns a wavefront object into numbers and temporary data structures.
//
// Note that the OBJ files refer to vertices and normals through a absolute
// count from the beginning of the file. OBJ files can be created from Blender.
func obj2Data(lines []string, odata *objData) (faces []face, err error) {
	for _, line := range lines {
		tokens := strings.Split(line, " ")
		var f1, f2, f3 float32
		var s1, s2, s3 string
		switch tokens[0] {
		case "v":
			if _, e := fmt.Sscanf(line, "v %f %f %f", &f1, &f2, &f3); e != nil {
				log.Printf("Bad vertex: %s\n", line)
				return faces, fmt.Errorf("could not parse vertex %s", e)
			}
			odata.v = append(odata.v, dataPoint{f1, f2, f3})
		case "vn":
			if _, e := fmt.Sscanf(line, "vn %f %f %f", &f1, &f2, &f3); e != nil {
				log.Printf("Bad normal: %s\n", line)
				return faces, fmt.Errorf("could not parse normal %s", e)
			}
			odata.n = append(odata.n, dataPoint{f1, f2, f3})
		case "vt":
			if _, e := fmt.Sscanf(line, "vt %f %f", &f1, &f2); e != nil {
				log.Printf("Bad texture coord: %s\n", line)
				return faces, fmt.Errorf("could not texture coordinate %s", e)
			}
			odata.t = append(odata.t, uvPoint{f1, 1 - f2})
		case "f":
			if _, e := fmt.Sscanf(line, "f %s %s %s", &s1, &s2, &s3); e != nil {
				log.Printf("Bad face: %s\n", line)
				return faces, fmt.Errorf("could not parse face %s", e)
			}
			faces = append(faces, face{[]string{s1, s2, s3}})
		case "o": // mesh name is processed before this method is called.
		case "s": // FUTURE: smoothing group - ignored for now.
		case "mtllib": // materials loaded separately and explicitly.
		case "usemtl": // material name - ignored, see above.
		}
	}
	return
}

// obj2MshData turns the data from .obj format into an internal OpenGL friendly
// format. The following information needs to be created for each mesh.
//
//    mesh.V = append(mesh.V, ...3-float32) - indexed from 0
//    mesh.N = append(mesh.N, ...3-float32) - indexed from 0
//    mesh.T = append(mesh.T, ...2-float32)	- indexed from 0
//    mesh.F = append(mesh.F, ...3-uint16)	- refers to above zero indexed values
//
// odata holds the global vertex, texture, and normal point information.
// faces are the indexes for this mesh.
//
// Additionally the normals at each vertex are generated as the sum of the
// normals for each face that shares that vertex.
func obj2MshData(name string, odata *objData, faces []face, data *MshData) (err error) {
	data.Name = name
	vmap := make(map[string]int) // the unique vertex data points for this face.
	vcnt := -1

	// process each vertex of each face.  Each one represents a combination vertex,
	// texture coordinate, and normal.
	v0 := &lin.V3{} // scratch
	v1 := &lin.V3{} // scratch
	for _, face := range faces {
		for pi := 0; pi < 3; pi++ {
			facei := face.s[pi]
			v, t, n := -1, -1, -1
			if v, t, n, err = parseFaceIndex(facei); err != nil {
				return fmt.Errorf("could not parse face data %s", err)
			}

			// cut down the amount of information passed around by reusing points
			// where the vertex and the texture coordinate information is the same.
			vertexIndex := fmt.Sprintf("%d/%d", v, t)
			if _, ok := vmap[vertexIndex]; !ok {

				// add a new data point.
				vcnt++
				vmap[vertexIndex] = vcnt
				data.V = append(data.V, odata.v[v].x, odata.v[v].y, odata.v[v].z)
				data.N = append(data.N, odata.n[n].x, odata.n[n].y, odata.n[n].z)
				if t != -1 {
					data.T = append(data.T, odata.t[t].u, odata.t[t].v)
				}
			} else {

				// update the normal at the vertex to be a combination of
				// all the normals of each face that shares the vertex.
				ni := vmap[vertexIndex] * 3
				n1 := v0.SetS(float64(data.N[ni]), float64(data.N[ni+1]), float64(data.N[ni+2]))
				n2 := v1.SetS(float64(odata.n[n].x), float64(odata.n[n].y), float64(odata.n[n].z))
				n2.Add(n2, n1).Unit()
				data.N[ni], data.N[ni+1], data.N[ni+2] = float32(n2.X), float32(n2.Y), float32(n2.Z)
			}
			data.F = append(data.F, uint16(vmap[vertexIndex]))
		}
	}
	if len(data.V) <= 0 || len(data.F) <= 0 {
		return fmt.Errorf("Minimally need vertex and face data for %s", name)
	}
	return err
}

// parseFace turns a face index point string (representing multiple indices)
// into 3 integer indices. The texture index is optional and is returned with
// a -1 value if it is not there.
func parseFaceIndex(findex string) (v, t, n int, err error) {
	v, t, n = -1, -1, -1
	if _, err = fmt.Sscanf(findex, "%d//%d", &v, &n); err != nil {
		if _, err = fmt.Sscanf(findex, "%d/%d/%d", &v, &t, &n); err != nil {
			return v, t, n, fmt.Errorf("Bad face (%s)\n", findex)
		}
	}
	v = int(v - 1)
	n = int(n - 1) // should all have the same value.
	if t != -1 {
		t = int(t - 1)
	}
	return
}
