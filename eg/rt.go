// Copyright © 2015-2016 Galvanized Logic. All rights reserved.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/gazed/vu/device"
	"github.com/gazed/vu/load"
	"github.com/gazed/vu/math/lin"
	"github.com/gazed/vu/render"
	"github.com/gazed/vu/render/gl"
)

// rt helps to understand ray tracing basics. Real time hardware supported ray
// tracing is a possible (far) future rendering alternative. Ray tracing is
// sufficiently different from standard rasterization it would likely need
// its own 3D engine.
// The code in this example is broken into two sections:
//    1. OpenGL based code that displays a single texture on a quad mesh.
//    2. Ray trace code that generates a ray trace image.
// Some general ray tracing reading ...
//   http://www.ics.uci.edu/~gopi/CS211B/RayTracing%20tutorial.pdf
//   http://www.gamasutra.com/blogs/AlexandruVoica/20140318/213148/Practical_techniques_for_ray_tracing_in_games.php?print=1
//   http://www.igorsevo.com/Article.aspx?article=A+simple+real-time+raytracer+in+CUDA+and+OpenGL
//   http://www.researchgate.net/publication/220183679_OptiX_A_General_Purpose_Ray_Tracing_Engine
//
// CONTROLS: NA
func rt() {
	rt := new(rtrace)
	dev := device.New("Ray Trace", 400, 400, 512, 512)
	rt.scene = rt.createScene() // create the scene for the ray tracer.
	rt.img = rt.rayTrace()      // create the ray traced image.
	rt.initRender()             // initialize opengl.
	dev.Open()
	for dev.IsAlive() {
		rt.update(dev)
		rt.drawScene()
		dev.SwapBuffers()
	}
	dev.Dispose()
}

// Encapsulate this examples methods using this structure.
type rtrace struct {
	vao     uint32     // vertex array object reference.
	mvp     render.Mvp // transform matrix for rendering.
	mvpID   int32      // mvp uniform id
	shaders uint32     // shader program reference.

	// texture information.
	img   *image.NRGBA // Texture data.
	texID uint32       // Graphics card texture identifier.
	tex2D int32        // texture sampler uniform id

	// quad mesh information
	verts []float32 // triangle verticies.
	faces []uint8   // connect verticies into two triangles (one quad)
	uvs   []float32 // uv texture coordinates.

	// ray trace information.
	iw, ih int      // image width and height.
	scene  []lin.V3 // scene comprised of sphere locations.
	procs  int      // number of runtime processors.

	// statistics.
	sampleCalls int // number of calls to the sampler.
	traceCalls  int // number of calls to the ray collision tracer.
}

// initRender is one time initialization that creates a single VAO
// to display a single ray trace generated texture.
func (rt *rtrace) initRender() {
	rt.verts = []float32{ // four verticies for a quad.
		0, 0, 0,
		4, 0, 0,
		0, 4, 0,
		4, 4, 0,
	}
	rt.faces = []uint8{ // create quad from 2 triangles.
		0, 2, 1,
		1, 2, 3,
	}
	rt.uvs = []float32{ // texture coordinates to sample the image.
		0, 0,
		1, 0,
		0, 1,
		1, 1,
	}

	// Start up OpenGL and create a single vertex array object.
	gl.Init()
	gl.GenVertexArrays(1, &rt.vao)
	gl.BindVertexArray(rt.vao)

	// vertex data.
	var vbuff uint32
	gl.GenBuffers(1, &vbuff)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbuff)
	gl.BufferData(gl.ARRAY_BUFFER, int64(len(rt.verts)*4), gl.Pointer(&(rt.verts[0])), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, 0)
	gl.EnableVertexAttribArray(0)

	// faces data.
	var ebuff uint32
	gl.GenBuffers(1, &ebuff)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebuff)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int64(len(rt.faces)), gl.Pointer(&(rt.faces[0])), gl.STATIC_DRAW)

	// texture coordatinates
	var tbuff uint32
	gl.GenBuffers(1, &tbuff)
	gl.BindBuffer(gl.ARRAY_BUFFER, tbuff)
	gl.BufferData(gl.ARRAY_BUFFER, int64(len(rt.uvs)*4), gl.Pointer(&(rt.uvs[0])), gl.STATIC_DRAW)
	var tattr uint32 = 2
	gl.VertexAttribPointer(tattr, 2, gl.FLOAT, false, 0, 0)
	gl.EnableVertexAttribArray(tattr)

	// use ray trace generated texture image.
	bounds := rt.img.Bounds()
	width, height := int32(bounds.Dx()), int32(bounds.Dy())
	ptr := gl.Pointer(&(rt.img.Pix[0]))
	gl.GenTextures(1, &rt.texID)
	gl.BindTexture(gl.TEXTURE_2D, rt.texID)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, width, height, 0, gl.RGBA, gl.UNSIGNED_BYTE, ptr)

	// texture sampling shader.
	shader := &load.ShdData{}
	err := shader.Load("tuv", load.NewLocator())
	if err != nil {
		log.Fatalf("Failed to load shaders %s\n", err)
	}
	rt.shaders = gl.CreateProgram()
	if err := gl.BindProgram(rt.shaders, shader.Vsh, shader.Fsh); err != nil {
		log.Fatalf("Failed to create program: %s\n", err)
	}
	rt.mvpID = gl.GetUniformLocation(rt.shaders, "mvpm")
	rt.tex2D = gl.GetUniformLocation(rt.shaders, "sampler2D")
	rt.mvp = render.NewMvp().Set(lin.NewM4().Ortho(0, 4, 0, 4, 0, 10))

	// set some state that doesn't need to change during drawing.
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
}

// update handles user input.
func (rt *rtrace) update(dev device.Device) {
	pressed := dev.Update()
	if pressed.Resized {
		_, _, ww, wh := dev.Size()
		gl.Viewport(0, 0, int32(ww), int32(wh))
	}
}

// drawScene renders the single texture on the quad.
func (rt *rtrace) drawScene() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(rt.shaders)
	gl.Uniform1i(rt.tex2D, 0)
	gl.ActiveTexture(gl.TEXTURE0 + 0)
	gl.BindVertexArray(rt.vao)
	gl.UniformMatrix4fv(rt.mvpID, 1, false, rt.mvp.Pointer())
	gl.DrawElements(gl.TRIANGLES, int32(len(rt.faces)), gl.UNSIGNED_BYTE, 0)

	// cleanup
	gl.ActiveTexture(0)
	gl.UseProgram(0)
	gl.BindVertexArray(0)
}

// opengl
// =============================================================================
// ray tracer.

// Credits:
// The ray trace example originates from Paul S. Heckbert "A Minimal Ray Tracer"
// and is based on Andrew Kensler's business card version of it.
//   http://www.cs.utah.edu/%7Eaek/code/card.cpp
// which was analyzed and commented on:
//   http://fabiensanglard.net/rayTracing_back_of_business_card/index.php
// and inspired a go code version:
//   https://github.com/abelsson/rays/blob/master/gorays/main.go
// The ray trace code below is based on the above links and
// uses the vu/math/lin linear math libaries.
//
// The example is brute force Whitted Ray tracing. It is simply optimized
// in that it uses only spheres for quicker ray collision calculations.
// Other simple optimizations could include
//   Use fewer spheres in the scene.
//   Use fewer rays per pixel - currently 64.
//
// Program flow:
//     rayTrace() - creates image using a render worker per processor.
//     render()   - one call to render for each image row. Render generates
//                  64 rays per image row pixel.
//     sample()   - calculates the color for one ray.
//     trace()    - collides a ray with the scene.

// The scene to be rendered. Each non-space is a sphere.
// The scene background, sky and floor, are generated.
// Some tests were run to capture the effect of adding spheres:
//    23 spheres. Sampler:10864441 Tracer:15557132 Time 2.353326s
//     1 sphere.  Sampler: 9614633 Tracer:11353321 Time 0.883293s
var art = []string{
	"                   ",
	"                   ",
	"                   ",
	"  *     *          ",
	"  *     *   *    * ",
	"  *     *   *    * ",
	"   *   *    *    * ",
	"    * *     *    * ",
	"     *       ****  ",
}

// createScene initializes the image size and positions spheres in 3D space.
// This must be called once before beginning raytracing.
func (rt *rtrace) createScene() []lin.V3 {
	rt.iw, rt.ih = 512, 512
	numRows := len(art)
	numCols := len(art[0])
	spheres := make([]lin.V3, 0, numRows*numCols)
	for k := numCols - 1; k >= 0; k-- {
		for j := numRows - 1; j >= 0; j-- {
			if art[j][numCols-1-k] != ' ' {
				location := lin.V3{X: -float64(k), Y: 3, Z: -float64(numRows-1-j) - 4}
				spheres = append(spheres, location)
			}
		}
	}
	return spheres
}

// rayTrace creates a single ray traced image.
// Its job is to create the workers that will render each row of the image.
// createScene() needs to have been called first.
func (rt *rtrace) rayTrace() *image.NRGBA {
	start := time.Now()         // track total raytrace time.
	rt.procs = runtime.NumCPU() // equals GOMAXPROCS since Go 1.5.
	img := image.NewNRGBA(image.Rect(0, 0, rt.iw, rt.ih))

	// we're nominally tracing rays for pixel (x,y) in the direction of ax+by+c.
	// At the image midpoint, this should be `g`
	g, a, b, c := lin.NewV3(), lin.NewV3(), lin.NewV3(), lin.NewV3()
	g.SetS(-5.5, -16, 0).Unit()                        // camera direction.
	a.SetS(0, 0, 1).Cross(a, g).Unit().Scale(a, 0.002) // camera up vector.
	b.Cross(g, a).Unit().Scale(b, 0.002)               // right vector.

	// Comment from Aeg:
	// "offset from the eye point (ignoring lens perturbation `t`)
	//  to the corner of the focal plane."
	c.Add(a, b).Scale(c, -256).Add(c, g)

	// create one worker goroutine per processor.
	rows := make(chan row, rt.ih)
	var wg sync.WaitGroup
	wg.Add(rt.procs)
	for i := 0; i < rt.procs; i++ {
		go rt.worker(*a, *b, *c, img, rows, &wg) // pass vectors by value.
	}

	// start assigning image rows to the workers.
	for y := (rt.ih - 1); y >= 0; y-- {
		rows <- row(y)
	}
	close(rows) // closing the worker comm channel causes workers to terminate...
	wg.Wait()   // ... once they finish their current assignment.

	// dump some render statistics.
	used := time.Since(start)
	log.Printf("Sample:%d Trace:%d Time:%fs ", rt.sampleCalls, rt.traceCalls, used.Seconds())
	return img
}

// worker reads from the 'rows' channel.
// A render is started for each row read from the channel.
func (rt *rtrace) worker(a, b, c lin.V3, img *image.NRGBA, rows <-chan row, wg *sync.WaitGroup) {
	defer wg.Done() // signal completion once all rows are processed.

	// render one row at a time, blocking if there are no more rows,
	// terminating once the rows channel is closed.
	seed := rand.Uint32()
	for r := range rows {
		r.render(rt, a, b, c, img, &seed)
	}
}

type row int // row is for rendering the colors for one row of the image.

// render one row of pixels by calculating a color for each pixel.
// The image pixel row number is r. Fill the pixel color into the
// image after the color has been calculated.
func (r row) render(rt *rtrace, a, b, c lin.V3, img *image.NRGBA, seed *uint32) {
	rgba := color.NRGBA{0, 0, 0, 255}
	t, v1, v2 := lin.NewV3(), lin.NewV3(), lin.NewV3() // temp vectors.
	color, orig, dir := lin.NewV3(), lin.NewV3(), lin.NewV3()
	for x := (rt.iw - 1); x >= 0; x-- {
		color.SetS(13, 13, 13) // Use a very dark default color.

		// Cast 64 rays per pixel for blur (stochastic sampling) and soft-shadows.
		for cnt := 0; cnt < 64; cnt++ {

			// Add randomness to the camera origin 17,16,8
			t.Scale(&a, rnd(seed)-0.5).Scale(t, 99).Add(t, v1.Scale(&b, rnd(seed)-0.5).Scale(v1, 99))
			orig.SetS(17, 16, 8).Add(orig, t)

			// Add randomness to the camera direction.
			rnda := rnd(seed) + float64(x)
			rndb := float64(r) + rnd(seed)
			dir.Scale(t, -1)
			dir.Add(dir, v1.Scale(&a, rnda).Add(v1, v2.Scale(&b, rndb)).Add(v1, &c).Scale(v1, 16))
			dir.Unit()

			// accumulate the color from each of the 64 rays.
			sample := rt.sample(*orig, *dir, seed)
			color = sample.Scale(&sample, 3.5).Add(&sample, color)
		}

		// set the final pixel color in the image.
		rgba.R = byte(color.X) // red
		rgba.G = byte(color.Y) // green
		rgba.B = byte(color.Z) // blue
		img.SetNRGBA(rt.iw-x, int(r), rgba)
	}
}

// sample calculates the color value for a given ray (origin and direction)
// shot into the scene. It relies on the trace() method to determine what the
// ray hit and recursively calls itself to add in the color values of any
// child rays. This method performs the job of a rasterization pipeline shader
// by calculating color based on light and normals wherever rays hit scene
// objects.
func (rt *rtrace) sample(orig, dir lin.V3, seed *uint32) (color lin.V3) {
	rt.sampleCalls++                        // track number of times called
	st, dist, bounce := rt.trace(orig, dir) // check ray scene collision.
	obounce := bounce

	// generate a sky color if the ray is going up.
	if st == missHigh {
		p := 1 - dir.Z // make the sky color lighter closer to the horizon.
		p = p * p
		p = p * p
		color.SetS(0.7, 0.6, 1)
		return *color.Scale(&color, p)
	}

	// add randomness to light for soft shadows.
	hitAt, lightDir, tmpv := lin.NewV3(), lin.NewV3(), lin.NewV3()
	hitAt.Add(&orig, tmpv.Scale(&dir, dist))
	lightDir.SetS(9+rnd(seed), 9+rnd(seed), 16).Add(lightDir, tmpv.Scale(hitAt, -1)).Unit()
	lightIntensity := lightDir.Dot(&bounce) // lambertian factor based on angle of light source.

	// check if the spot is in shadow by tracing a ray from the
	// intersection point to the light. Its in shadow if there is a hit.
	shadowFactor := 1.0
	if lightIntensity < 0 {
		lightIntensity = 0
		shadowFactor = 0
	} else {
		var hitStatus int
		if hitStatus, _, bounce = rt.trace(*hitAt, *lightDir); hitStatus != missHigh {
			lightIntensity = 0
			shadowFactor = 0
		}
	}

	// generate a floor color if the ray was going down.
	if st == missLow {
		hitAt.Scale(hitAt, 0.2)
		color.SetS(3, 3, 3) // gray floor squares.
		if int(math.Ceil(hitAt.X)+math.Ceil(hitAt.Y))&1 == 1 {
			color.SetS(1, 1, 3) // blue floor squares.
		}
		return *(color.Scale(&color, lightIntensity*0.2+0.1))
	}

	// calculate the color 'rgb' with diffuse and specular component.
	// r is the reflection vector.
	reflectDir := lin.NewV3()
	reflectDir.Add(&dir, tmpv.Scale(&obounce, obounce.Dot(tmpv.Scale(&dir, -2))))
	rgb := lightDir.Dot(reflectDir.Scale(reflectDir, shadowFactor))
	rgb = math.Pow(rgb, 99)

	// cast a child ray from where the parent ray hit.
	// Add in the result of the color from the child ray.
	color.SetS(rgb, rgb, rgb)
	addColor := rt.sample(*hitAt, *reflectDir, seed)
	addColor.Scale(&addColor, 0.5)
	return *(color.Add(&color, &addColor))
}

// trace casts a ray (origin, direction) to see if the ray hits any
// of the spheres in the scene. The possible return values are:
//   missHigh : no hit and ray goes up.
//   missLow  : no hit and ray goes down.
//   hit      : hit so return distance and reflection ray.
func (rt *rtrace) trace(orig, dir lin.V3) (hitStatus int, minHitDistance float64, bounce lin.V3) {
	rt.traceCalls++ // track number of times called.
	minHitDistance = 1e9
	hitStatus = missHigh
	s := -orig.Z / dir.Z
	if 0.01 < s {
		minHitDistance = s
		bounce.SetS(0, 0, 1)
		hitStatus = missLow
	}

	// cast the ray against each sphere in the scene.
	// http://www.lighthouse3d.com/tutorials/maths/ray-sphere-intersection/
	// http://kylehalladay.com/blog/tutorial/math/2013/12/24/Ray-Sphere-Intersection.html
	tempv := lin.NewV3()
	for i := range rt.scene {
		tempv.Add(&orig, &(rt.scene[i])) // ray origin + sphere center.
		b := tempv.Dot(&dir)             // represent the intersection ray...
		c := tempv.Dot(tempv) - 1        // ... and the sphere radius
		b2 := b * b                      // ... using squared magnitudes.

		// if the ray intersected the sphere.
		if b2 > c {
			q := b2 - c                      // convert the squared length values....
			hitDistance := -b - math.Sqrt(q) // ... to the actual collision distance.

			// remember the minimum hit distance.
			if hitDistance < minHitDistance && hitDistance > 0.01 {
				minHitDistance = hitDistance
				bounce = *tempv
				hitStatus = hit
			}
		}
	}
	if hitStatus == hit {
		bounce.Add(&bounce, tempv.Scale(&dir, minHitDistance)).Unit()
	}
	return
}

// ray hit status is one of the following.
const (
	missHigh = iota // ray will hit the sky.
	missLow         // ray will hit the floor.
	hit             // ray has hit a sphere.
)

// rnd returns a pseudo random value between 0 and 1 based on a given seed.
// This method provides about a 10x speedup over rand.Float64().
func rnd(s *uint32) float64 {
	ss := *s
	ss += ss
	ss ^= 1
	if int32(ss) < 0 {
		ss ^= 0x88888eef
	}
	*s = ss
	return float64(*s%95) / float64(95)
}
