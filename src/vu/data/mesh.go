// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package data

// Mesh holds 3D model data in a format that is easily consumed by a rendering
// layer. A mesh is expected to be reused and thus does not contain any instance
// information like location or scale.
//
// An approximation is to keep an individual mesh to less than 65,000 verticies
// so that it works on lower capability systems like mobile devices.
type Mesh struct {
	Name string // unique mesh name.

	// Verticies are points in 3D space.  Each vertex is specified with
	// (x, y, z, w) where
	//    x, y are  mandatory 2D minimum coordinates,
	//    z can be defaulted to 0 for a two dimensional vertex, and
	//    w can be defaulted to 1 for a constant depth.
	V []float32 // arranged as [][4]float32

	// Normals for each vertex. Normals are specified as (x, y, z) where each
	// value is between 0.0 and 1.0. The slice length is expected to be the same
	// length as the vertex slice, and each normal is calculated to be the
	// normalized sum of the normals for each face that shares the vertex.
	N []float32 // arranged as [][3]float32

	// Texture coordinates. Specifies how the texture data is aligned relative
	// to the fragment. There are two floats for each texture (u, v)  where the
	// values are expected to be between 0.0 and 1.0. The texture data is loaded
	// separately from an image file.
	//
	// The texture corresponding to these texture coordinates is expected to be
	// available to the rendering layer at the same time as this mesh.
	T []float32 // arranged as [][2]float32

	// Faces are used to index the vertex data into shapes. These are
	// expected to be triangles so 3 indicies form one face.  Each face value
	// refers to a single vertex.
	F []uint16 // arranged as [][3]uint16

	// Vao is an vertex array object that is a reference to the above
	// vertex and normal data.  The reference is to data is stored on the graphics
	// card. The "*buf" variables are graphic card references for
	// the above data.
	Vao  uint32 // vertex array references all buffers.
	Vbuf uint32 // vertex buffer
	Tbuf uint32 // texture buffer
	Nbuf uint32 // normals buffer
	Fbuf uint32 // face index buffer
}
