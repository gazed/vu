// Copyright © 2016 Galvanized Logic. All rights reserved.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"math"

	"github.com/gazed/vu"
	"github.com/gazed/vu/synth"
)

// ss is an EXPERIMENTAL example for exploring procedural models
// using the supershape formula.
// Lots of information and implementations on the web:
//    https://en.wikipedia.org/wiki/Superformula
//    http://paulbourke.net/geometry/supershape/
//
// CONTROLS:
//   Esc   : reset to default values
//   Tab   : toggle 2D 3D shapes
//   DaUa  : select parameter       : up down
//   LaRa  : change value           : decrease increase
//   -=    : change value           : decrease increase
//   Sh    : change value by x10    : works with LaRa
//   Alt   : change value by x100   : works with LaRa
//   WS    : move model             : forward back
//   AD    : rotate model           : left right
//   1-9,0 : precanned shapes
func ss() {
	ss := &superShape{}
	if err := vu.New(ss, "Super Shapes", 400, 100, 800, 600); err != nil {
		log.Printf("ss: error starting engine %s", err)
	}
	defer catchErrors()
}

// superShape is the  unique "tag" to encapsulate the demo specific data.
type superShape struct {
	cam3   *vu.Camera // 3D camera for 3D shapes.
	cam2   *vu.Camera // 2D camera for 2D shapes and shape parameters.
	m2D    *vu.Pov    // model to render the generated 2D shape.
	m3D    *vu.Pov    // model to render the generated 3D shape.
	vb     []float32  // reusable vertex buffer to create models.
	fb     []uint16   // reusable face buffer to create models.
	nb     []float32  // reusable normal buffer to create models.
	row    int        // controls which shape value is being edited.
	rowbg  *vu.Pov    // background row hilight for the editing shape values.
	orig   []float64  // original shape values.
	mag    float64    // shape parameter editing amount.
	mindex int        // which magnitude for editing parameters.
	show3D bool       // toggle between 2D and 3D visualization.

	// control the shape using different values for the parameters.
	shape  *synth.Form // superformula shape.
	labels []*vu.Pov   // indexed by M, N1, N2, N3, A, B
	amount *vu.Pov     // show parameter editing amount value.
}

// Create is the startup asset creation.
func (ss *superShape) Create(eng vu.Eng, s *vu.State) {
	scene3 := eng.Root().NewPov()
	ss.cam3 = scene3.NewCam()
	scene2 := eng.Root().NewPov()
	ss.cam2 = scene2.NewCam().SetUI()
	ss.resize(s.W, s.H)
	ss.orig = []float64{0, 1, 1, 1}
	ss.mag = 0.1 // default edit amount.

	// hiliting the selected shape parameters.
	ss.rowbg = scene2.NewPov().SetAt(115, 110, 0).SetScale(55, 20, 0)
	ss.rowbg.NewModel("alpha", "msh:icon", "mat:transparent_blue")

	// mesh for 2D shapes.
	scale := 100.0
	ss.shape = synth.NewForm()
	ss.m2D = eng.Root().NewPov().SetAt(400, 300, 0).SetScale(scale, scale, 0)
	ss.m2D.NewModel("solid").Make("msh:gen2d")
	ss.m2D.Model().Set(vu.DrawMode(vu.Lines)).SetUniform("kd", 1, 1, 1)
	ss.m2D.Cull = ss.show3D
	ss.genSquare(ss.m2D.Model()) // use hand generated model to start.

	// mesh for 3D shapes.
	// debug with .SetDrawMode(vu.Points) and "solid" shader.
	ss.m3D = scene3.NewPov().SetAt(0, 0, -8).SetScale(2.5, 2.5, 2.5)
	ss.m3D.NewModel("gouraud", "mat:blue").Make("msh:gen3d")
	ss.genSquare(ss.m3D.Model()) // use hand generated model to start.
	ss.m3D.Cull = !ss.show3D

	// Display the shape parameters.
	left := 50.0
	font := "lucidiaSu18"
	ss.labels = []*vu.Pov{nil, nil, nil, nil, nil, nil, nil}
	ss.labels[M] = scene2.NewPov().SetAt(left, 120, 0)
	ss.labels[M].NewLabel("uv", font, font+"White")
	ss.labels[N1] = scene2.NewPov().SetAt(left, 100, 0)
	ss.labels[N1].NewLabel("uv", font, font+"White")
	ss.labels[N2] = scene2.NewPov().SetAt(left, 80, 0)
	ss.labels[N2].NewLabel("uv", font, font+"White")
	ss.labels[N3] = scene2.NewPov().SetAt(left, 60, 0)
	ss.labels[N3].NewLabel("uv", font, font+"White")
	ss.labels[A] = scene2.NewPov().SetAt(left, 40, 0)
	ss.labels[A].NewLabel("uv", font, font+"White")
	ss.labels[B] = scene2.NewPov().SetAt(left, 20, 0)
	ss.labels[B].NewLabel("uv", font, font+"White")
	ss.amount = scene2.NewPov().SetAt(left+120, 0, 0)
	ss.amount.NewLabel("uv", font, font+"White")
	ss.positionLabels()
	ss.updateLabels()
}

// Update is the regular engine callback.
func (ss *superShape) Update(eng vu.Eng, in *vu.Input, s *vu.State) {
	if in.Resized {
		ss.resize(s.W, s.H)
	}
	for press, down := range in.Down {
		switch {

		// select shape parameter.
		case press == vu.KDa && down == 1:
			// select the next shape parameter row.
			ss.row = (ss.row + 1) % (B + 1)
			ss.positionLabels()
		case press == vu.KUa && down == 1:
			// select the previous shape parameter row.
			ss.row = ss.row - 1
			if ss.row < 0 {
				ss.row = B
			}
			ss.positionLabels()

		// alter edit value amounts with control keys.
		// Hold one of these keys down at the same time as + or -.
		case press == vu.KShift:
			if down > 0 {
				ss.mag = 1.0
			} else {
				ss.mag = 0.1
			}
			ss.updateLabels()
		case press == vu.KAlt:
			if down > 0 {
				ss.mag = 10.0
			} else {
				ss.mag = 0.1
			}
			ss.updateLabels()

		// change the current parameter value.
		case press == vu.KEsc && down == 1:
			// reset the shape parameter to its original value.
			ss.reset(ss.row)
			ss.updateLabels()
			ss.genShape()
		case (press == vu.KEqual || press == vu.KRa) && down == 1:
			// increase the current shape parameter by the current magnitude.
			ss.alter(ss.row, ss.mag)
			ss.updateLabels()
			ss.genShape()
		case (press == vu.KMinus || press == vu.KLa) && down == 1:
			// decrease the current shape parameter by the current magnitude.
			ss.alter(ss.row, -ss.mag)
			ss.updateLabels()
			ss.genShape()

		// move and rotate the 3D supershape model.
		case press == vu.KA:
			ss.m3D.Spin(0, -2, 0)
		case press == vu.KD:
			ss.m3D.Spin(0, 2, 0)
		case press == vu.KW:
			x, y, z := ss.m3D.At()
			ss.m3D.SetAt(x, y, z-0.2)
		case press == vu.KS:
			x, y, z := ss.m3D.At()
			ss.m3D.SetAt(x, y, z+0.2)

		// switch between 2D and 3D visualization.
		case press == vu.KTab && down == 1:
			ss.show3D = !ss.show3D
			ss.genShape()

		// pre-canned shape favorites.
		case press == vu.K1 && down == 1:
			ss.precanned(1)
		case press == vu.K2 && down == 1:
			ss.precanned(2)
		case press == vu.K3 && down == 1:
			ss.precanned(3)
		case press == vu.K4 && down == 1:
			ss.precanned(4)
		case press == vu.K5 && down == 1:
			ss.precanned(5)
		case press == vu.K6 && down == 1:
			ss.precanned(6)
		case press == vu.K7 && down == 1:
			ss.precanned(7)
		case press == vu.K8 && down == 1:
			ss.precanned(8)
		case press == vu.K9 && down == 1:
			ss.precanned(9)
		case press == vu.K0 && down == 1:
			ss.precanned(0)
		}
	}
}

// resize handles user changes to the window size.
func (ss *superShape) resize(ww, wh int) {
	ss.cam2.SetOrthographic(0, float64(ww), 0, float64(wh), 0, 50)
	ss.cam3.SetPerspective(60, float64(ww)/float64(wh), 0.1, 50)
}

// precanned for a quick demo of some cool reference shapes.
func (ss *superShape) precanned(index int) {
	switch index {
	case 0: // reset
		ss.shape.SetValues(0, 1, 1, 1, 1, 1)
	case 1:
		ss.shape.SetValues(4, 0.3, 0.3, 0.3, 1, 1)
	case 2:
		ss.shape.SetValues(19.0/6.0, 0.3, 0.3, 0.3, 1, 1)
	case 3:
		ss.shape.SetValues(5, 0.05, 1.7, 1.7, 1, 1)
	case 4:
		ss.shape.SetValues(6, 1000, 400, 400, 1, 1) // hexagon
	case 5:
		ss.shape.SetValues(4, 0.5, 0.5, 4, 1, 1)
	case 6:
		ss.shape.SetValues(5, 3, 3, 3, 0.9, 0.3)
	case 7:
		ss.shape.SetValues(30, 75, 1.5, 35, 1, 0.6)
	case 8:
		ss.shape.SetValues(2, 0.5, 0.5, 2, 1, 1.2)
	case 9:
		ss.shape.SetValues(20, 2, 6, 6, 0.9, 0.9)
	}
	ss.updateLabels()
	ss.genShape()
}

// genShape recreates the shape based on the current shape parameters.
// It uses 2D or 3D shape methods based on the current display mode.
func (ss *superShape) genShape() {
	ss.m3D.Cull = !ss.show3D
	ss.m2D.Cull = ss.show3D
	if ss.show3D {
		ss.gen3DShape(ss.m3D.Model())
	} else {
		ss.gen2DShape(ss.m2D.Model())
	}
}

// gen2DShape runs the superformula over a series of samples and
// then links the samples with a line.
func (ss *superShape) gen2DShape(model vu.Model) {
	vc := uint16(0)                // vertex counter.
	vb, fb := ss.vb[:0], ss.fb[:0] // keep any allocated memory.
	samples := 100
	sampleSize := 2 * math.Pi / float64(samples)
	for cnt := 0; cnt < samples; cnt++ {
		angle := float64(cnt) * sampleSize
		vx, vy := ss.shape.At2D(angle)
		vb = append(vb, float32(vx), float32(vy), 0)
		fb = append(fb, vc, vc+1)
		vc++
	}
	fb[len(fb)-1] = 0 // close loop by linking the last vertex to the first.
	ss.vb, ss.fb = vb, fb

	// reset the mesh data with the generated data.
	model.Mesh().InitData(0, 3, vu.DynamicDraw, false).SetData(0, vb)
	model.Mesh().InitFaces(vu.DynamicDraw).SetFaces(fb)
}

// gen3DShape runs the superformula using planar coordinate sampling and
// then generates the 3D model by creating triangles from the sampled points.
func (ss *superShape) gen3DShape(model vu.Model) {
	vb, fb, nb := ss.vb[:0], ss.fb[:0], ss.nb[:0] // reuse keeping allocated memory.
	halfPi := math.Pi / 2.0
	samples := 100

	// generate end point for first vertex.
	vx, vy, vz := ss.shape.At3D(-halfPi, -math.Pi)
	vb = append(vb, float32(vx), float32(vy), float32(vz))
	nb = append(nb, 0, 0, 0)

	// lat is the angle in radians -Pi/2 to Pi/2, ie: bottom to top. 19 samples.
	// lon is the angle in radians -Pi to Pi, ie: full circle sweep. 20 samples.
	latSample := math.Pi / float64(samples)
	lonSample := 2 * math.Pi / float64(samples)
	for lat := -halfPi + latSample; lat < halfPi; lat += latSample {
		for lon := -math.Pi; lon < math.Pi; lon += lonSample {
			vx, vy, vz = ss.shape.At3D(lat, lon)
			vb = append(vb, float32(vx), float32(vy), float32(vz))
			nb = append(nb, 0, 0, 0)
		}
	}

	// generate end point for last vertex.
	vx, vy, vz = ss.shape.At3D(halfPi, math.Pi)
	vb = append(vb, float32(vx), float32(vy), float32(vz))
	nb = append(nb, 0, 0, 0)

	// generate the triangle faces for the verticies. Create a quad,
	// two triangles at a time, that join adjacent vertices from the shape samples.
	var v1, v2, v3, v4 uint16
	for lat := 0; lat < samples-2; lat++ {
		for lon := 0; lon < samples; lon++ {
			v1 = uint16(1 + lat*samples + lon)
			v2 = uint16(1 + lat*samples + lon + 1)
			if lat == samples-1 {
				v3 = uint16(1 + lon + 1)
				v4 = uint16(1 + lon)
			} else {
				v3 = uint16(1 + (lat+1)*samples + lon + 1)
				v4 = uint16(1 + (lat+1)*samples + lon)
			}
			if lon == samples-1 {
				v2 = uint16(1 + (lat)*samples)
				v3 = uint16(1 + (lat+1)*samples)
			}
			fb = append(fb, v1, v2, v3)
			fb = append(fb, v1, v3, v4)
		}
	}

	// close the shape by joining the first langitude sweep indicies to the
	// first vertex and the last langitude sweep indicies to the last vertex.
	fv, lv := uint16(0), uint16(samples*(samples-1)+1) // first and last indicies.
	for cnt := 0; cnt < samples; cnt++ {
		v2 = uint16(cnt) + 2
		v3 = uint16(cnt) + 1
		if cnt == samples-1 {
			v2 = 1
		}
		fb = append(fb, fv, v2, v3)

		// last vertex.
		v1 = lv - uint16(cnt) - 2
		v2 = lv - uint16(cnt) - 1
		if cnt == samples-1 {
			v1 = lv - 1
		}
		fb = append(fb, v1, v2, lv)
	}
	ss.vb, ss.nb, ss.fb = vb, nb, fb // update references to new memory.
	ss.updateNormals()

	// reset the mesh data with the generated data.
	model.Mesh().InitData(0, 3, vu.DynamicDraw, false).SetData(0, vb)
	model.Mesh().InitData(1, 3, vu.DynamicDraw, false).SetData(1, nb)
	model.Mesh().InitFaces(vu.DynamicDraw).SetFaces(fb)
}

// updateNormals is called to process ss.vb and ss.nb.
// The normals at each vertex are the sums of all normals
// for the faces that share that vertex.
func (ss *superShape) updateNormals() {
	vb, nb, fb := ss.vb, ss.nb, ss.fb
	for fc := 0; fc < len(fb); fc += 3 {

		// calculate the normal for a triangle face.
		i1, i2, i3 := fb[fc]*3, fb[fc+1]*3, fb[fc+2]*3 // triangle indicies
		ax, ay, az := vb[i1], vb[i1+1], vb[i1+2]       // tri vertex 1
		bx, by, bz := vb[i2], vb[i2+1], vb[i2+2]       // tri vertex 2
		cx, cy, cz := vb[i3], vb[i3+1], vb[i3+2]       // tri vertex 3
		x2, y2, z2 := ax-bx, ay-by, az-bz
		x1, y1, z1 := cx-bx, cy-by, cz-bz
		nx, ny, nz := y1*z2-z1*y2, z1*x2-x1*z2, x1*y2-y1*x2 // cross product

		// add the face normal to each vertex of the triangle.
		nb[i1], nb[i1+1], nb[i1+2] = nb[i1]+nx, nb[i1+1]+ny, nb[i1+2]+nz
		nb[i2], nb[i2+1], nb[i2+2] = nb[i2]+nx, nb[i2+1]+ny, nb[i2+2]+nz
		nb[i3], nb[i3+1], nb[i3+2] = nb[i3]+nx, nb[i3+1]+ny, nb[i3+2]+nz
	}

	// normalize the vertex normals.
	for nc := 0; nc < len(nb); nc += 3 {
		nb[nc], nb[nc+1], nb[nc+2] = ss.normalize(nb[nc], nb[nc+1], nb[nc+2])
	}
}

// normalize the vector a,b,c returning a unit length vector na,nb,nc.
func (ss *superShape) normalize(a, b, c float32) (na, nb, nc float32) {
	length := float32(math.Sqrt(float64(a*a + b*b + c*c)))
	if length == 0 {
		return a, b, c
	}
	return a / length, b / length, c / length
}

// positionLabels updates the row highlight and the
// edit amount to match the parameter row that has focus.
func (ss *superShape) positionLabels() {
	x, _, z := ss.rowbg.At()
	ss.rowbg.SetAt(x, float64(130-ss.row*20), z)
	x, _, z = ss.amount.At()
	ss.amount.SetAt(x, float64(120-ss.row*20), z)
}

// updateLabels regenerates the labels for the parameter labels.
func (ss *superShape) updateLabels() {
	s := fmt.Sprintf("m  = %5.4f", ss.shape.M)
	ss.labels[M].Model().SetStr(s)
	s = fmt.Sprintf("n1 = %5.4f", ss.shape.N1)
	ss.labels[N1].Model().SetStr(s)
	s = fmt.Sprintf("n2 = %5.4f", ss.shape.N2)
	ss.labels[N2].Model().SetStr(s)
	s = fmt.Sprintf("n3 = %5.4f", ss.shape.N3)
	ss.labels[N3].Model().SetStr(s)
	s = fmt.Sprintf("a   = %5.4f", ss.shape.A)
	ss.labels[A].Model().SetStr(s)
	s = fmt.Sprintf("b   = %5.4f", ss.shape.B)
	ss.labels[B].Model().SetStr(s)

	// show the magnitude when changing param amounts.
	s = fmt.Sprintf("+/-   %3.1f", ss.mag)
	ss.amount.Model().SetStr(s)
}

// genSquare is an example of manually creating a square using lines.
// It shows how to generate the necessary verticies and faces that
// are algorithmically produced for generated shapes.
func (ss *superShape) genSquare(m vu.Model) {
	vb := ss.vb[:0] // reuse, keeping allocated memory.
	fb := ss.fb[:0] // reuse, keeping allocated memory.

	// mesh data for a square drawn with lines.
	vb = append(vb, -0.5, -0.5, 0.0) // 0
	vb = append(vb, 0.5, -0.5, 0.0)  // 1
	vb = append(vb, 0.5, 0.5, 0.0)   // 2
	vb = append(vb, -0.5, 0.5, 0.0)  // 3
	fb = append(fb, 0, 1, 1, 2, 2, 3, 3, 0)
	m.Mesh().InitData(0, 3, vu.DynamicDraw, false).SetData(0, vb)
	m.Mesh().InitFaces(vu.DynamicDraw).SetFaces(fb)
}

// Indicies for super forumula parameters.
const (
	M  = iota // angle multiplier
	N1        // overall exponent.
	N2        // cos exponent.
	N3        // sin exponent.
	A         // cos divisor
	B         // sin divisor
)

// reset the indicated supershape parameter.
// Invalid parameter indicies are ignored.
func (ss *superShape) reset(index int) {
	switch index {
	case M:
		ss.shape.M = synth.CircleForm.M
	case N1:
		ss.shape.N1 = synth.CircleForm.N1
	case N2:
		ss.shape.N2 = synth.CircleForm.N2
	case N3:
		ss.shape.N3 = synth.CircleForm.N3
	case A:
		ss.shape.A = synth.CircleForm.A
	case B:
		ss.shape.B = synth.CircleForm.B
	}
}

// alter adds the given amount to the indicated parameter.
// Invalid parameter indicies are ignored.
func (ss *superShape) alter(index int, amount float64) {
	switch index {
	case M:
		ss.shape.M += amount
	case N1:
		ss.shape.N1 += amount
	case N2:
		ss.shape.N2 += amount
	case N3:
		ss.shape.N3 += amount
	case A:
		ss.shape.A += amount
	case B:
		ss.shape.B += amount
	}
}
