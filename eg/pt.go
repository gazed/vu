// Copyright © 2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"

	"github.com/gazed/vu"
	"github.com/gazed/vu/synth"
)

// pt is an EXPERIMENTAL example for exploring generated textures.
// A good explanation of procedurally generated textures can be found at:
//    http://www.upvector.com/?section=Tutorials&subsection=Intro%20to%20Procedural%20Textures
//
// CONTROLS:
//   Esc   : reset to default values
//   1-4   : precanned textures
//   DaUa  : select parameter       : up down
//   LaRa  : change value           : decrease increase
//   -=    : change value           : decrease increase
//   Sh    : change value by x10    : works with LaRa
//   Alt   : change value by x100   : works with LaRa
//   CD    : colors                 : enable disable
func pt() {
	pt := &pttag{}
	if err := vu.New(pt, "Procedural Textures", 400, 100, 800, 600); err != nil {
		log.Printf("pt: error starting engine %s", err)
	}
	defer catchErrors()
}

// pttag is the  unique "tag" to encapsulate the demo specific data.
type pttag struct {
	cam2 *vu.Camera          // 2D camera for textures and shape parameters.
	sn   *synth.SimplexNoise // main component of interesting textures.

	// controls for editing simplex noise parameters.
	row    int       // controls which parameter is being edited.
	rowbg  *vu.Pov   // hilight for the editing shape values.
	mag    float64   // shape parameter editing amount.
	mindex int       // controls amount of change when editing.
	labels []*vu.Pov // indexed by F, G, L, O
	amount *vu.Pov   // show parameter editing amount value.

	// memory for generating the images.
	tex   []*vu.Pov      // models to display the generated textures.
	imgs  []*image.NRGBA // separate texture images for each model.
	data  [][]float64    // placeholder for generated noise data.
	vdata [][]float64    // generated voronoi data.

	// texture generation functions and references.
	marble  []func(x, y float64) float64 // specific combination.
	cloud   []func(x, y float64) float64 // specific combination.
	strange []func(x, y float64) float64 // specific combination.
	center  []func(x, y float64) float64 // specific combination.
	texture []func(x, y float64) float64 // pointer to function list.
	colors  []color.NRGBA                // random color palette.
}

// Create is the startup asset creation.
func (pt *pttag) Create(eng vu.Eng, s *vu.State) {
	scene2 := eng.Root().NewPov()
	pt.cam2 = scene2.NewCam().SetUI()
	pt.resize(s.W, s.H)
	pt.mag = 0.01 // default edit value.

	// initialize the parts that create the image.
	imgSize, numImgs := 256, 6
	pt.sn = synth.NewSimplexNoise(124)
	for cnt := 0; cnt < numImgs; cnt++ {
		pt.imgs = append(pt.imgs, image.NewNRGBA(image.Rect(0, 0, imgSize, imgSize)))
	}
	pt.data = make([][]float64, imgSize)
	pt.vdata = make([][]float64, imgSize)
	for x := range pt.data {
		pt.data[x] = make([]float64, imgSize)
		pt.vdata[x] = make([]float64, imgSize)
	}

	// hiliting the selected shape parameters.
	pt.rowbg = scene2.NewPov().SetAt(130, 110, 0).SetScale(55, 20, 0)
	pt.rowbg.NewModel("alpha", "msh:icon", "mat:transparent_blue")

	// uv mapped icon to display generated texture.
	scale := 256.0
	tex := eng.Root().NewPov().SetAt(600, 200, 0).SetScale(scale, scale, 0)
	tex.NewModel("uv", "msh:icon").Make("tex:gen0")
	pt.tex = append(pt.tex, tex)

	// display the individual texture algorithms...
	tex = eng.Root().NewPov().SetAt(100, 500, 0).SetScale(128, 128, 0)
	tex.NewModel("uv", "msh:icon").Make("tex:fn1")
	pt.tex = append(pt.tex, tex)
	tex = eng.Root().NewPov().SetAt(250, 500, 0).SetScale(128, 128, 0)
	tex.NewModel("uv", "msh:icon").Make("tex:fn2")
	pt.tex = append(pt.tex, tex)
	tex = eng.Root().NewPov().SetAt(400, 500, 0).SetScale(128, 128, 0)
	tex.NewModel("uv", "msh:icon").Make("tex:fn3")
	pt.tex = append(pt.tex, tex)
	tex = eng.Root().NewPov().SetAt(550, 500, 0).SetScale(128, 128, 0)
	tex.NewModel("uv", "msh:icon").Make("tex:fn4")
	pt.tex = append(pt.tex, tex)
	tex = eng.Root().NewPov().SetAt(700, 500, 0).SetScale(128, 128, 0)
	tex.NewModel("uv", "msh:icon").Make("tex:fn5")
	pt.tex = append(pt.tex, tex)

	// group the majority of the generation functions.
	pt.marble = []func(x, y float64) float64{pt.land, pt.sinx}
	pt.cloud = []func(x, y float64) float64{pt.land, pt.radial}
	pt.strange = []func(x, y float64) float64{pt.smoke, pt.voronoi}
	pt.center = []func(x, y float64) float64{pt.smoke, pt.radial}
	pt.texture = pt.marble

	// Display the texture parameters.
	left := 50.0
	font := "lucidiaSu18"
	pt.labels = []*vu.Pov{nil, nil, nil, nil}
	pt.labels[F] = scene2.NewPov().SetAt(left, 120, 0)
	pt.labels[F].NewLabel("uv", font, font+"White")
	pt.labels[G] = scene2.NewPov().SetAt(left, 100, 0)
	pt.labels[G].NewLabel("uv", font, font+"White")
	pt.labels[L] = scene2.NewPov().SetAt(left, 80, 0)
	pt.labels[L].NewLabel("uv", font, font+"White")
	pt.labels[O] = scene2.NewPov().SetAt(left, 60, 0)
	pt.labels[O].NewLabel("uv", font, font+"White")
	pt.amount = scene2.NewPov().SetAt(left+120, 0, 0)
	pt.amount.NewLabel("uv", font, font+"White")
	pt.positionLabels()
	pt.updateLabels()

	// initialize the images: main composite image and...
	pt.genTex(pt.texture)                    // 1: noise image
	pt.tex[1].Model().Tex(0).Set(pt.imgs[1]) //
	pt.fillOne(pt.data, pt.sinx)             // 2: sin image
	pt.flatColor(pt.data, pt.imgs[2])        //
	pt.tex[2].Model().Tex(0).Set(pt.imgs[2]) //
	pt.fillOne(pt.data, pt.radial)           // 3: radial image
	pt.flatColor(pt.data, pt.imgs[3])        //
	pt.tex[3].Model().Tex(0).Set(pt.imgs[3]) //
	pt.fillVoronoi(pt.vdata)                 // 4: voronoi image
	pt.flatColor(pt.vdata, pt.imgs[4])       //
	pt.tex[4].Model().Tex(0).Set(pt.imgs[4]) //
	pt.fillOne(pt.data, pt.smoke)            // 5: smoke image
	pt.flatColor(pt.data, pt.imgs[5])        //
	pt.tex[5].Model().Tex(0).Set(pt.imgs[5]) //
}

// Update is the regular engine callback.
func (pt *pttag) Update(eng vu.Eng, in *vu.Input, s *vu.State) {
	if in.Resized {
		pt.resize(s.W, s.H)
	}
	for press, down := range in.Down {
		switch {

		// select simplex noise parameter.
		case press == vu.KDa && down == 1:
			// select the next parameter row.
			pt.row = (pt.row + 1) % len(pt.labels)
			pt.positionLabels()
		case press == vu.KUa && down == 1:
			// select the previous parameter row.
			pt.row = pt.row - 1
			if pt.row < 0 {
				pt.row = len(pt.labels) - 1
			}
			pt.positionLabels()

		// alter edit value amounts with control keys.
		// Hold one of these keys down at the same time as + or -.
		case press == vu.KShift:
			if down > 0 {
				pt.mag = 0.1
			} else {
				pt.mag = 0.01
			}
			pt.updateLabels()
		case press == vu.KAlt:
			if down > 0 {
				pt.mag = 1.0
			} else {
				pt.mag = 0.01
			}
			pt.updateLabels()

		// change the colors
		case press == vu.KC && down == 1: // enable colors
			pt.colors = pt.newColors(pt.colors)
			pt.genTex(pt.texture)
		case press == vu.KD && down == 1: // disable colors
			pt.colors = pt.colors[:0]
			pt.genTex(pt.texture)

		// change the current parameter value.
		case press == vu.KEsc && down == 1:
			pt.reset(pt.row)
			pt.updateLabels()
			pt.genTex(pt.texture)
		case (press == vu.KEqual || press == vu.KRa) && down == 1:
			// increase the current parameter by the current magnitude.
			pt.alter(pt.row, pt.mag)
			pt.updateLabels()
			pt.genTex(pt.texture)
		case (press == vu.KMinus || press == vu.KLa) && down == 1:
			// decrease the current parameter by the current magnitude.
			pt.alter(pt.row, -pt.mag)
			pt.updateLabels()
			pt.genTex(pt.texture)

		// pre-canned combined textures .
		case press == vu.K1 && down == 1:
			pt.texture = pt.marble
			pt.genTex(pt.texture)
		case press == vu.K2 && down == 1:
			pt.texture = pt.cloud
			pt.genTex(pt.texture)
		case press == vu.K3 && down == 1:
			pt.texture = pt.strange
			pt.genTex(pt.texture)
		case press == vu.K4 && down == 1:
			pt.texture = pt.center
			pt.genTex(pt.texture)
		}
	}
}

// resize handles user changes to the window size.
func (pt *pttag) resize(ww, wh int) {
	pt.cam2.SetOrthographic(0, float64(ww), 0, float64(wh), 0, 50)
}

// Indicies for simplex noise parameters.
const (
	F = iota // frequency.
	G        // gain
	L        // lacunarity
	O        // octaves.
)

// reset the indicated supershape parameter.
// Invalid parameter indicies are ignored.
func (pt *pttag) reset(index int) {
	switch index {
	case F:
		pt.sn.F = 0.5
	case G:
		pt.sn.G = 0.55
	case L:
		pt.sn.L = 2.0
	case O:
		pt.sn.O = 6
	}
}

// alter adds the given amount to the indicated parameter.
// Invalid parameter indicies are ignored.
func (pt *pttag) alter(index int, amount float64) {
	switch index {
	case F:
		pt.sn.F += amount
	case G:
		pt.sn.G += amount
	case L:
		pt.sn.L += amount
	case O:
		pt.sn.O += int(amount)
	}
}

// Algorithms for generating image data.
func (pt *pttag) land(x, y float64) float64 { return pt.sn.Gen2D(x, y) }
func (pt *pttag) smoke(x, y float64) float64 {
	if v := pt.sn.Gen2D(x, y) + 0.5; v > 0 {
		return v
	}
	return 0
}
func (pt *pttag) sinx(x, y float64) float64     { return (1.0 + math.Sin(x*50)) / 2 }
func (pt *pttag) gradient(x, y float64) float64 { return x }
func (pt *pttag) voronoi(x, y float64) float64 {
	ix := int(math.Max(math.Min(x*255, 255), 0)) // limit to 0-255.
	iy := int(math.Max(math.Min(y*255, 255), 0)) // limit to 0-255.
	return pt.vdata[ix][iy]
}

// radial is a 0-1 sized area with center at 0.5,0.5.
// Algorithm for generating image data.
func (pt *pttag) radial(x, y float64) float64 {
	dx, dy := x-0.5, y-0.5
	return math.Max(0.5-math.Sqrt(dx*dx+dy*dy), 0) // ensure 0 or greater
}

// fillCombo pipelines the given texture functions over each data point.
func (pt *pttag) fillCombo(data [][]float64, fns []func(x, y float64) float64) {
	flip := len(data) - 1 // Put 0,0 at bottom left.
	for x := range data {
		for y := range data[x] {
			x0, y0 := float64(x)/float64(flip), float64(y)/float64(flip)
			for _, fn := range fns {
				x0 = fn(x0, y0)
			}
			data[x][flip-y] = x0
		}
	}
}

// fillOne applies the given texture generation function to all data points.
func (pt *pttag) fillOne(data [][]float64, fn func(x, y float64) float64) {
	flip := len(data) - 1 // Put 0,0 at bottom left.
	for x := range data {
		for y := range data[x] {
			x0, y0 := float64(x)/float64(flip), float64(y)/float64(flip)
			data[x][flip-y] = fn(x0, y0)
		}
	}
}

// Brute force voronoi diagram. See https://en.wikipedia.org/wiki/Voronoi_diagram
// Currently drawn with with euclidean distance - could use manhatten distance.
// Points are scattered about for voronoi. Currently points are random within
// a grid cell within the surface - could just randomly scatter over the surface.
func (pt *pttag) fillVoronoi(data [][]float64) {
	var points []float64

	// Divide 0-1 into a grid based on size and scatter points
	// within each grid cell.
	size := 4
	div := 1.0 / float64(size)
	for xcell := 0; xcell < size; xcell++ {
		for ycell := 0; ycell < size; ycell++ {
			x, y := rand.Float64(), rand.Float64()           // 0-1 range.
			ox, oy := float64(xcell)*div, float64(ycell)*div // cell offset
			x, y = x*div+ox, y*div+oy                        // point in grid cell.
			points = append(points, x, y)
		}
	}

	// Use the distance to a point as the data value.
	flip := len(data) - 1 // Put 0,0 at bottom left.
	for x := range data {
		for y := range data[x] {
			fx, fy := float64(x)/float64(flip), float64(y)/float64(flip)
			distance := math.MaxFloat64
			for cnt := 0; cnt < len(points); cnt += 2 {
				px, py := points[cnt], points[cnt+1]
				dsqr := (fx-px)*(fx-px) + (fy-py)*(fy-py)
				if dsqr < distance {
					distance = dsqr
				}
			}
			data[x][flip-y] = math.Sqrt(distance) * 2
		}
	}
}

// genTex generates the texture and sets it on the icon model.
func (pt *pttag) genTex(fns []func(x, y float64) float64) {
	pt.fillCombo(pt.data, fns)
	if len(pt.colors) > 0 {
		pt.gradientColor(pt.data, pt.imgs[0], pt.colors)
	} else {
		pt.flatColor(pt.data, pt.imgs[0])
	}
	pt.tex[0].Model().Tex(0).Set(pt.imgs[0])

	// regenerate the noise images affected by a parameter change.
	pt.fillOne(pt.data, pt.land)
	pt.flatColor(pt.data, pt.imgs[1])
	pt.tex[1].Model().Tex(0).Set(pt.imgs[1])
	pt.fillOne(pt.data, pt.smoke)
	pt.flatColor(pt.data, pt.imgs[5])
	pt.tex[5].Model().Tex(0).Set(pt.imgs[5])
}

// =============================================================================
// Handle editing the noise parameters.

// positionLabels updates the row highlight and the
// edit amount to match the parameter row that has focus.
func (pt *pttag) positionLabels() {
	x, _, z := pt.rowbg.At()
	pt.rowbg.SetAt(x, float64(130-pt.row*20), z)
	x, _, z = pt.amount.At()
	pt.amount.SetAt(x, float64(120-pt.row*20), z)
}

// updateLabels regenerates the labels for the parameter labels.
func (pt *pttag) updateLabels() {
	s := fmt.Sprintf("Freq = %5.3f", pt.sn.F)
	pt.labels[F].Model().SetStr(s)
	s = fmt.Sprintf("Gain = %5.3f", pt.sn.G)
	pt.labels[G].Model().SetStr(s)
	s = fmt.Sprintf("Lac   = %5.3f", pt.sn.L)
	pt.labels[L].Model().SetStr(s)
	s = fmt.Sprintf("Oct   = %5d", pt.sn.O)
	pt.labels[O].Model().SetStr(s)

	// show the magnitude when changing param amounts.
	s = fmt.Sprintf("+/-   %3.2f", pt.mag)
	pt.amount.Model().SetStr(s)
}

// =============================================================================
// Procedurally generating colors... for some background see:
// http://devmag.org.za/2012/07/29/how-to-choose-colours-procedurally-algorithms/
//
// The following code adds color to generated images.

// flatColor translates the data into a gray scale color where the data values
// are expected to be in the range 0-1.
func (pt *pttag) flatColor(t [][]float64, img *image.NRGBA) *image.NRGBA {
	var c *color.NRGBA
	for x := range t {
		for y := range t[x] {
			val := uint8(t[x][y] * 255)
			c = &color.NRGBA{val, val, val, 255}
			img.SetNRGBA(x, y, *c)
		}
	}
	return img
}

// gradientColor colors the images by interpolating between the given colors.
func (pt *pttag) gradientColor(t [][]float64, img *image.NRGBA, colors []color.NRGBA) *image.NRGBA {
	var c *color.NRGBA
	for x := range t {
		for y := range t[x] {
			c = interpolateColor(colors, t[x][y])
			img.SetNRGBA(x, y, *c)
		}
	}
	return img
}

// newColors picks a completely random set of 2-5 colors.
// FUTURE: improve this so that the chosen colors work together.
func (pt *pttag) newColors(colors []color.NRGBA) []color.NRGBA {
	colors = colors[:0]       // reset keeping memory.
	num := rand.Int31n(4) + 2 // 2-5
	for cnt := int32(0); cnt < num; cnt++ {
		r, g, b := rand.Int31n(256), rand.Int31n(256), rand.Int31n(256)
		c := color.NRGBA{uint8(r), uint8(g), uint8(b), 255}
		colors = append(colors, c)
	}
	return colors
}

// interpolateColor returns a color based on the given value and
// a range of colours. The algorithm is based on answer 2 from:
// http://stackoverflow.com/questions/1236683/color-interpolation-between-3-colors-in-net
func interpolateColor(colors []color.NRGBA, x float64) (c *color.NRGBA) {
	sigma2 := 0.035
	r, g, b := 0.0, 0.0, 0.0
	total := 0.0
	step := 1.0 / float64(len(colors)-1)
	mu := 0.0
	for _ = range colors {
		total += math.Exp(-(x-mu)*(x-mu)/(2.0*sigma2)) / math.Sqrt(2.0*math.Pi*sigma2)
		mu += step
	}
	mu = 0.0
	for _, col := range colors {
		percent := math.Exp(-(x-mu)*(x-mu)/(2.0*sigma2)) / math.Sqrt(2.0*math.Pi*sigma2)
		mu += step
		r += float64(col.R) * percent / total
		g += float64(col.G) * percent / total
		b += float64(col.B) * percent / total
	}
	return &color.NRGBA{uint8(r), uint8(g), uint8(b), 255}
}
