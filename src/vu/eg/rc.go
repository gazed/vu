// Copyright © 2014 Galvanized Logic. All rights reserved.
// Use is governed by a FreeBSD license found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"vu"
	"vu/math/lin"
)

// rc demonstrates ray-casting and mouse-picking. What this demo
// also shows is that ray casting needs better integration into
// the engine. Specifically:
//    • move/Solid needs closer integration with move/Body.
//    • pickDir should be moved into the engine.
//    • the inverse projection, and inverse view matricies should be
//      kept up to date as part of vu.Scene.
func rc() {
	rc := &rctag{}
	var err error
	if rc.eng, err = vu.New("Ray Cast", 400, 100, 800, 600); err != nil {
		log.Fatal("rc: error intitializing engine %s", err)
	}
	rc.eng.SetDirector(rc) // override user input handling.
	defer rc.eng.Shutdown()
	defer catchErrors()
	rc.eng.Action()
}

// Globally unique "tag" for this example.
type rctag struct {
	eng    vu.Engine
	scene  vu.Scene
	fsize  float64 // floor size in world space units.
	gsize  float64 // grid size: number visible/virtual tiles in floor image.
	ray    *lin.V4 // ray direction.
	ipm    *lin.M4 // inverse projection matrix. Updated if projection changes.
	ivm    *lin.M4 // inverse view matrix. Updated if the camera moves.
	ww, wh int     // window dimensions.
	floor  vu.Part // plane for raycast testing
	hilite vu.Part // tracks which tile is currently selected.
	s0, s1 vu.Part // spheres for raycast testing.
	s2, s3 vu.Part // spheres for raycast testing.
	banner vu.Part // shows selected grid locations.
}

// Create is the engine one-time initialization callback.
func (rc *rctag) Create(eng vu.Engine) {
	eng.Color(0.2, 0.2, 0.2, 1)
	rc.scene = eng.AddScene(vu.VP)
	rc.ray = lin.NewV4()
	rc.ipm = lin.NewM4()
	rc.ivm = lin.NewM4I()
	rc.gsize = 32 // 4x8 ie. tile image is 4x4 grid, tile.obj is oversampled by 8.

	// The ray cast target is a plane displaying the image of a 32x32 grid.
	rc.fsize = 10.0                                       // 2x2 plane to 20x20 plane.
	rc.floor = rc.scene.AddPart()                         // create the floor.
	rc.floor.SetSolid(vu.Plane(0, 0, -1))                 // the floors ray intersect shape.
	rc.floor.SetScale(rc.fsize, rc.fsize, 0)              // scale the model to fsize.
	rc.floor.SetLocation(0, 0, -20)                       // move away from camera, along -Z.
	rc.floor.SetRole("uv").SetMesh("tile").AddTex("tile") // put the image on the floor.
	rc.floor.Role().SetTexMode(0, vu.TEX_REPEAT)          // repeat for UV values > 1.
	rc.floor.Spin(-45, 0, 0)

	// create a selected tile tracker.
	rc.hilite = rc.scene.AddPart()
	rc.hilite.SetScale(0.3125, 0.3125, 0.001) // scale to cover a single tile.
	rc.hilite.Spin(-45, 0, 0)
	rc.hilite.SetRole("uv").SetMesh("icon").AddTex("image")

	// Put spheres at the floor corners when the floor is rotated 45 deg about X.
	rc.s0 = rc.makeSphere(10, 7.07, -20-7.07, 1, 0, 0)
	rc.s1 = rc.makeSphere(-10, 7.07, -20-7.07, 0, 1, 0)
	rc.s2 = rc.makeSphere(10, -7.07, -20+7.07, 0, 0, 1)
	rc.s3 = rc.makeSphere(-10, -7.07, -20+7.07, 1, 1, 0)

	// Add a banner to show the currently selected grid location.
	_, _, w, h := eng.Size()
	over := eng.AddScene(vu.VO)
	over.Set2D()
	over.SetOrthographic(0, float64(w), 0, float64(h), 0, 10)
	rc.banner = over.AddPart()
	rc.banner.SetRole("uv").AddTex("weblySleek22White").SetFont("weblySleek22")
	rc.banner.SetLocation(100, 100, 0)
	rc.banner.SetVisible(false)
	rc.resize()
}

// makeSphere creates a sphere at the given x, y, z location and with
// the given r, g, b colour.
func (rc *rctag) makeSphere(x, y, z, r, g, b float64) vu.Part {
	sz := 0.5
	p := rc.scene.AddPart()
	p.SetSolid(vu.Ball(sz))
	p.SetLocation(x, y, z).SetScale(sz, sz, sz)
	p.SetRole("flat").SetMesh("sphere").SetKd(r, g, b)
	return p
}

// Update is the engine frequent user-input/state-update callback.
func (rc *rctag) Update(in *vu.Input) {
	if in.Resized {
		rc.resize()
	}
	for press, down := range in.Down {
		switch press {
		case "Lm":
			if down == 1 {
				rc.raycast(in.Mx, in.My)
			}
		case "Q":
			rc.floor.Spin(-1, 0, 0)
			rc.hilite.Spin(-1, 0, 0)
			rc.hilite.SetVisible(false)
		case "E":
			rc.floor.Spin(1, 0, 0)
			rc.hilite.Spin(1, 0, 0)
			rc.hilite.SetVisible(false)
		}
	}
	rc.hovercast(in.Mx, in.My)
}

// resize handles user screen/window changes.
func (rc *rctag) resize() {
	var x, y int
	x, y, rc.ww, rc.wh = rc.eng.Size()
	rc.eng.Resize(x, y, rc.ww, rc.wh)
	fov, ratio, near, far := 60.0, float64(rc.ww)/float64(rc.wh), 0.1, 500.0
	rc.scene.SetPerspective(fov, ratio, near, far)
	rc.ipm.PerspInv(fov, ratio, near, far)
}

// raycast checks which grid tile is selected on a mouse click. It gets
// the picking ray direction and then intersect the ray against the
// geometry in world space.
func (rc *rctag) raycast(mx, my int) {
	rd := rc.pickDir(mx, my, rc.ww, rc.wh, rc.ipm, rc.ivm, rc.ray)
	ray := vu.Ray(rd.X, rd.Y, rd.Z)

	// collide the ray with the plane and get the world-space contact point on hit.
	if hit, x, y, z := vu.Cast(ray, rc.floor.Solid()); hit {

		// find the floor corners in world-space.
		rot := lin.NewQ().SetS(rc.floor.Rotation()).Unit() // current floor spin.
		bot := &lin.V3{-rc.fsize, -rc.fsize, 0}
		top := &lin.V3{rc.fsize, rc.fsize, 0}
		bot.MultQ(bot, rot)
		top.MultQ(top, rot)

		// check if the plane hit was within the floor area.
		if x >= bot.X && x <= top.X && y >= bot.Y && y <= top.Y {

			// place a marker where the mouse hit.
			rc.hilite.SetLocation(x, y, z)
			rc.hilite.SetVisible(true)

			// adjust and display grid coordinates. Map x, y to 0:31
			xsize, ysize := top.X-bot.X, top.Y-bot.Y
			gx := int(((x * 2 / xsize) + 1) / 2 * rc.gsize)
			gy := int(((y * 2 / ysize) + 1) / 2 * rc.gsize)
			rc.banner.Role().SetPhrase(fmt.Sprintf("%d:%d", gx, gy))
			rc.banner.SetVisible(true)
		} else {
			rc.hilite.SetVisible(false) // missed the grid.
			rc.banner.SetVisible(false)
		}
	} else {
		rc.hilite.SetVisible(false) // missed the plane entirely.
		rc.banner.SetVisible(false)
	}
}

// hovercast checks the sphere each update and turns the spheres a different
// colour when the mouse is over them.
func (rc *rctag) hovercast(mx, my int) {
	rd := rc.pickDir(mx, my, rc.ww, rc.wh, rc.ipm, rc.ivm, rc.ray)
	ray := vu.Ray(rd.X, rd.Y, rd.Z)
	parts := []vu.Part{rc.s0, rc.s1, rc.s2, rc.s3}
	colors := []rgb{rgb{1, 0, 0}, rgb{0, 1, 0}, rgb{0, 0, 1}, rgb{1, 1, 0}}
	for cnt, p := range parts {
		if hit, _, _, _ := vu.Cast(ray, p.Solid()); hit {
			p.Role().SetKd(1, 1, 1)
		} else {
			p.Role().SetKd(colors[cnt].R, colors[cnt].G, colors[cnt].B)
		}
	}
}

type rgb struct{ R, G, B float64 } // rgb holds a colour.

// pickDir applies inverse transforms to derive world space coordinates for
// a ray projected from the camera through the mouse's screen position. See:
//     http://bookofhook.com/mousepick.pdf
//     http://antongerdelan.net/opengl/raycasting.html
//     http://schabby.de/picking-opengl-ray-tracing/
//     (opengl FAQ Picking 20.0.010)
//     http://www.opengl.org/archives/resources/faq/technical/selection.htm
//     http://www.codeproject.com/Articles/625787/Pick-Selection-with-OpenGL-and-OpenCL
// FUTURE: incorporate this into the engine. Possibly Engine, Stage, or Scene.
func (rc *rctag) pickDir(mx, my, ww, wh int, ipm, ivm *lin.M4, ray *lin.V4) *lin.V4 {
	if mx >= 0 && mx <= ww && my >= 0 && my <= wh {
		clipx := float64(2*mx)/float64(ww) - 1 // mx to range -1:1
		clipy := float64(2*my)/float64(wh) - 1 // my to range -1:1
		clip := ray.SetS(clipx, clipy, -1, 1)

		// use the inverse perspective to go from clip to eye (view) coordinates
		eye := clip.MultMv(ipm, clip)
		eye.Z = -1 // into the screen
		eye.W = 0  // want a vector, not a point

		// use the inverse view to go from eye (view) coordinates to world coordinates.
		ray = eye.MultMv(ivm, eye)
	}
	return ray.Unit() // return the normalized direction vector.
}
