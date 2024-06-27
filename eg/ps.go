// Copyright © 2024 Galvanized Logic Inc.

package main

import (
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/gazed/vu"
	"github.com/gazed/vu/load"
	"github.com/gazed/vu/math/lin"
)

// ps primitive shapes explores creating geometric shapes and standard
// shape primitives using shaders. This example demonstrates:
//   - loading assets.
//   - creating a 3D scene.
//   - controlling scene camera movement.
//   - draw circles primitives, one from a shader, one from lines.
//   - generate icosphere meshes using eng.MakeMesh.
//
// CONTROLS:
//   - W,S    : move forward, back
//   - A,D    : move left, right
//   - C,Z    : move up, down
//   - RMouse : look around
//   - Q      : quit and close window.
func ps() {
	ps := &mhtag{}

	defer catchErrors()
	eng, err := vu.NewEngine(
		vu.Windowed(),
		vu.Title("Primitive Shapes"),
		vu.Size(200, 200, 1600, 900),
		vu.Background(0.01, 0.01, 0.01, 1.0),
	)
	if err != nil {
		slog.Error("ps: engine start", "err", err)
		return
	}

	// import assets from asset files.
	// This creates the assets referenced by the models below.
	// Note that circle and quad meshes engine defaults.
	eng.ImportAssets("circle.shd", "lines.shd", "pbr0.shd")

	// The scene holds the cameras and lighting information
	// and acts as the root for all models added to the scene.
	ps.scene = eng.AddScene(vu.Scene3D)

	// add one directional light. SetAt sets the direction.
	ps.scene.AddLight(vu.DirectionalLight).SetAt(-1, -2, -2)

	// Draw a 3D line circle using a shader and a quad.
	scale := 3.0
	c1 := ps.scene.AddModel("shd:circle", "msh:quad")
	c1.SetAt(-1.5, 0, -5).SetScale(scale, scale, scale)

	// Draw a 3D line circle using a circle model and lines.
	c2 := ps.scene.AddModel("shd:lines", "msh:circle")
	c2.SetAt(+1.5, 0, -5).SetScale(scale, scale, scale)
	c2.SetColor(0, 1, 0, 1) // green
	// draw a half size line circle.
	c3 := ps.scene.AddModel("shd:lines", "msh:circle")
	c3.SetAt(+3.0, 0, -5).SetScale(scale/2, scale/2, scale/2)
	c3.SetColor(1, 0, 0, 1) // red

	// create and draw an icosphere. At the lowest resolution this
	// looks bad because the normals are shared where vertexes are
	// part of multiple triangles.
	genIcosphereMesh(eng, 0)
	s0 := ps.scene.AddModel("shd:pbr0", "msh:icosphere0")
	s0.SetAt(-3, 0, -10).SetColor(0, 0, 1, 1).SetMetallicRoughness(true, 0.2)

	// a higher resolution icosphere starts to look ok with lighting.
	genIcosphereMesh(eng, 4)
	s2 := ps.scene.AddModel("shd:pbr0", "msh:icosphere4")
	s2.SetAt(+3, 0, -10).SetColor(0, 1, 0, 1).SetMetallicRoughness(true, 0.2)

	eng.Run(ps) // does not return while example is running.
}

// Globally unique "tag" that encapsulates example specific data.
type pstag struct {
	scene  *vu.Entity
	mx, my int32   // mouse position
	pitch  float64 // Up-down look direction.
	yaw    float64 // Left-right look direction.
}

// Update is the application engine callback.
func (ps *pstag) Update(eng *vu.Engine, in *vu.Input, delta time.Duration) {
	// react to one time press events.
	for press := range in.Pressed {
		switch press {
		case vu.KQ:
			// quit if Q is pressed
			eng.Shutdown()
			return
		}
	}

	// get mouse position difference from last update.
	xdiff, ydiff := in.Mx-ps.mx, in.My-ps.my // mouse move differences...
	ps.mx, ps.my = in.Mx, in.My              // ... from last mouse location.

	// react to continuous press events.
	lookSpeed := 15.0 * delta.Seconds()
	move := 10.0 // move so many units worth in one second.
	speed := move * delta.Seconds()
	cam := ps.scene.Cam()
	for press := range in.Down {
		switch press {
		case vu.KW:
			cam.Move(0, 0, -speed, cam.Lookat()) // -Z forward (into screen)
		case vu.KS:
			cam.Move(0, 0, speed, cam.Lookat()) // +Z back (away from screen)
		case vu.KA:
			cam.Move(-speed, 0, 0, cam.Lookat()) // left
		case vu.KD:
			cam.Move(speed, 0, 0, cam.Lookat()) // right
		case vu.KC:
			cam.Move(0, speed, 0, cam.Lookat()) // up
		case vu.KZ:
			cam.Move(0, -speed, 0, cam.Lookat()) // down
		case vu.KMR:
			if ydiff != 0 {
				ps.pitch = ps.limitPitch(ps.pitch + float64(-ydiff)*lookSpeed)
				cam.SetPitch(ps.pitch)
			}
			if xdiff != 0 {
				ps.yaw += float64(-xdiff) * lookSpeed
				cam.SetYaw(ps.yaw)
			}
		}
	}
}

// limitPitch ensures that look up/down is limited to 90 degrees.
// This helps reduce confusion when looking around.
func (ps *pstag) limitPitch(pitch float64) float64 {
	switch {
	case pitch > 90:
		return 90
	case pitch < -90:
		return -90
	}
	return pitch
}

// genIcosphereMesh creates a unit sphere made of triangles.
// Higher subdivisions create more triangles. Supported values are 0-7:
//   - 0:       20 triangles
//   - 1:       80 triangles
//   - 2:      320 triangles
//   - 3:     1280 triangles
//   - 4:     5120 triangles
//   - 5:   20_480 triangles
//   - 6:   81_920 triangles
func genIcosphereMesh(eng *vu.Engine, subdivisions int) (err error) {
	if subdivisions < 0 || subdivisions > 6 {
		return fmt.Errorf("genIcosphereMesh: unsupported subdivision %d", subdivisions)
	}

	// create the initial icosphere mesh data.
	verts, indexes := genIcosphere(subdivisions)

	// generate triangle normals. This produces the same number of indexes
	// and triangles but more vertexes since the vertexes are not shared
	// between triangles - they each must have their own normal.
	newVerts := []float32{}
	normals := []float32{}
	newIndexes := []uint16{}
	for i := 0; i < len(indexes); i += 3 {
		v1, v2, v3 := indexes[i], indexes[i+1], indexes[i+2]
		p1x, p1y, p1z := verts[v1*3], verts[v1*3+1], verts[v1*3+2]
		p2x, p2y, p2z := verts[v2*3], verts[v2*3+1], verts[v2*3+2]
		p3x, p3y, p3z := verts[v3*3], verts[v3*3+1], verts[v3*3+2]
		newVerts = append(newVerts, p1x, p1y, p1z)
		newVerts = append(newVerts, p2x, p2y, p2z)
		newVerts = append(newVerts, p3x, p3y, p3z)

		// use midpoint of triangle as normal for the triangle vertexes.
		mx := (p1x + p2x + p3x) / 3.0
		my := (p1y + p2y + p3y) / 3.0
		mz := (p1z + p2z + p3z) / 3.0
		normal := lin.NewV3().SetS(float64(mx), float64(my), float64(mz)).Unit()
		nx := float32(normal.X)
		ny := float32(normal.Y)
		nz := float32(normal.Z)
		normals = append(normals, nx, ny, nz)
		normals = append(normals, nx, ny, nz)
		normals = append(normals, nx, ny, nz)

		// same number of triangles... but now pointing to unique vertexes/normals.
		newIndexes = append(newIndexes, uint16(i), uint16(i+1), uint16(i+2))
	}

	// load the generated data into a mesh using eng.MakeMesh.
	meshData := make(load.MeshData, load.VertexTypes)
	meshData[load.Vertexes] = load.F32Buffer(newVerts, 3)
	meshData[load.Normals] = load.F32Buffer(normals, 3)
	meshData[load.Indexes] = load.U16Buffer(newIndexes)
	meshTag := fmt.Sprintf("icosphere%d", subdivisions)
	return eng.MakeMesh(meshTag, meshData)
}

// genIcosphere creates mesh data for a unit sphere based on triangles.
// The number of vertexes increases with each subdivision.
// Based on:
//   - http://blog.andreaskahler.com/2009/06/creating-icosphere-mesh-in-code.html
//
// The normals on a unit sphere nothing more than the direction from the
// center of the sphere to each vertex.
//
// Using uint16 for indexes limits the number of vertices to 65535.
//
// FUTURE: look at a slower icosphere subdivision that avoids exponential growth, see:
// https://devforum.roblox.com/t/hex-planets-dev-blog-i-generating-the-hex-sphere/769805
func genIcosphere(subdivisions int) (vertexes []float32, indexes []uint16) {
	midPointCache := map[int64]uint16{} // stores new midpoint vertex indexes.

	// addVertex is a closure that adds a vertex, ensuring that the
	// vertex is on a unit sphere. Note the vertex is also the normal.
	// Return the index of the vertex
	addVertex := func(x, y, z float32) uint16 {
		length := float32(math.Sqrt(float64(x*x + y*y + z*z)))
		vertexes = append(vertexes, x/length, y/length, z/length)
		return uint16(len(vertexes)/3) - 1 // indexes start at 0.
	}

	// getMidPoint is a closure that fetches or creates the
	// midpoint index between indexes p1 and p2.
	getMidPoint := func(p1, p2 uint16) (index uint16) {

		// first check if the middle point has already been added as a vertex.
		smallerIndex, greaterIndex := p1, p2
		if p2 < p1 {
			smallerIndex, greaterIndex = p2, p1
		}
		key := int64(smallerIndex)<<32 + int64(greaterIndex)
		if val, ok := midPointCache[key]; ok {
			return val
		}

		// not cached, then add a new vertex
		p1X, p1Y, p1Z := vertexes[p1*3], vertexes[p1*3+1], vertexes[p1*3+2]
		p2X, p2Y, p2Z := vertexes[p2*3], vertexes[p2*3+1], vertexes[p2*3+2]
		midx := (p1X + p2X) / 2.0
		midy := (p1Y + p2Y) / 2.0
		midz := (p1Z + p2Z) / 2.0

		// add vertex makes sure point is on unit sphere
		index = addVertex(midx, midy, midz)

		// cache the new midpoint and return index
		midPointCache[key] = index
		return index
	}

	// create initial 12 vertices of a icosahedron
	// from the corners of 3 orthogonal planes.
	t := float32((1.0 + math.Sqrt(5.0)) / 2.0)
	addVertex(-1, +t, 0) // corners of XY-plane
	addVertex(+1, +t, 0)
	addVertex(-1, -t, 0)
	addVertex(+1, -t, 0)
	addVertex(0, -1, +t) // corners of YZ-plane
	addVertex(0, +1, +t)
	addVertex(0, -1, -t)
	addVertex(0, +1, -t)
	addVertex(+t, 0, -1) // corners of XZ-plane
	addVertex(+t, 0, +1)
	addVertex(-t, 0, -1)
	addVertex(-t, 0, +1)

	// create 20 triangles of the icosahedron
	// 5 faces around point 0
	indexes = append(indexes, 0, 11, 5)
	indexes = append(indexes, 0, 5, 1)
	indexes = append(indexes, 0, 1, 7)
	indexes = append(indexes, 0, 7, 10)
	indexes = append(indexes, 0, 10, 11)

	// 5 adjacent faces
	indexes = append(indexes, 1, 5, 9)
	indexes = append(indexes, 5, 11, 4)
	indexes = append(indexes, 11, 10, 2)
	indexes = append(indexes, 10, 7, 6)
	indexes = append(indexes, 7, 1, 8)

	// 5 faces around point 3
	indexes = append(indexes, 3, 9, 4)
	indexes = append(indexes, 3, 4, 2)
	indexes = append(indexes, 3, 2, 6)
	indexes = append(indexes, 3, 6, 8)
	indexes = append(indexes, 3, 8, 9)

	// 5 adjacent faces
	indexes = append(indexes, 4, 9, 5)
	indexes = append(indexes, 2, 4, 11)
	indexes = append(indexes, 6, 2, 10)
	indexes = append(indexes, 8, 6, 7)
	indexes = append(indexes, 9, 8, 1)

	// create new triangles for each level of subdivision
	for i := 0; i < subdivisions; i++ {

		// create 4 new triangles to replace each existing triangle.
		newIndexes := []uint16{}
		for i := 0; i < len(indexes); i += 3 {
			v1, v2, v3 := indexes[i], indexes[i+1], indexes[i+2]
			a := getMidPoint(v1, v2) // create or fetch mid-point vertex.
			b := getMidPoint(v2, v3) //   ""
			c := getMidPoint(v3, v1) //   ""

			newIndexes = append(newIndexes, v1, a, c)
			newIndexes = append(newIndexes, v2, b, a)
			newIndexes = append(newIndexes, v3, c, b)
			newIndexes = append(newIndexes, a, b, c)
		}

		// replace the old indexes with the new ones.
		indexes = newIndexes
	}
	return vertexes, indexes
}
