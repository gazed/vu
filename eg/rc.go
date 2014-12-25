// Copyright © 2014 Galvanized Logic. All rights reserved.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"fmt"
	"log"

	"github.com/gazed/vu"
	"github.com/gazed/vu/math/lin"
)

// rc demonstrates ray-casting and mouse-picking. In this example the plane
// is centered at the origin and the camera is tilted 45 degrees around X.
func rc() {
	rc := &rctag{}
	var err error
	if rc.eng, err = vu.New("Ray Cast", 400, 100, 800, 600); err != nil {
		log.Fatal("rc: error intitializing engine %s", err)
	}
	rc.eng.SetDirector(rc) // get user input through Director.Update()
	rc.create()            // create initial assests.
	defer rc.eng.Shutdown()
	defer catchErrors()
	rc.eng.Action()
}

// Globally unique "tag" for this example.
type rctag struct {
	eng    vu.Engine
	scene  vu.Scene
	cam    vu.Camera
	fsize  float64 // floor size in world space units.
	gsize  float64 // grid size: number visible/virtual tiles in floor image.
	ww, wh int     // window dimensions.
	floor  vu.Part // plane for raycast testing
	hilite vu.Part // tracks which tile is currently selected.
	s0, s1 vu.Part // spheres for raycast testing.
	s2, s3 vu.Part // spheres for raycast testing.
	banner vu.Part // shows selected grid locations.
}

// create is the startup asset creation.
func (rc *rctag) create() {
	rc.eng.Color(0.2, 0.2, 0.2, 1)
	rc.scene = rc.eng.AddScene(vu.VP)
	rc.cam = rc.scene.Cam()
	rc.cam.Spin(45, 0, 0)          // Tilt the camera and
	rc.cam.SetLocation(0, -14, 14) // ...point directly at 0, 0, 0
	rc.gsize = 32                  // 4x8 ie. image is 4x4 grid, tile.obj is oversampled by 8.

	// The ray cast target is a plane displaying the image of a 32x32 grid.
	rc.fsize = 10.0                                       // 2x2 plane to 20x20 plane.
	rc.floor = rc.scene.AddPart()                         // create the floor.
	rc.floor.SetForm(vu.NewPlane(0, 0, -1))               // the floors ray intersect shape.
	rc.floor.SetScale(rc.fsize, rc.fsize, 0)              // scale the model to fsize.
	rc.floor.SetRole("uv").SetMesh("tile").AddTex("tile") // put the image on the floor.
	rc.floor.Role().SetTexMode(0, vu.TEX_REPEAT)          // repeat for UV values > 1.

	// create a selected tile tracker.
	rc.hilite = rc.scene.AddPart()
	rc.hilite.SetScale(0.625, 0.625, 0.001) // scale to cover a single tile.
	rc.hilite.SetRole("uv").SetMesh("icon").AddTex("image")

	// Put spheres at the floor corners.
	rc.s0 = rc.makeSphere(10, 10, 0, 1, 0, 0)
	rc.s1 = rc.makeSphere(-10, 10, 0, 0, 1, 0)
	rc.s2 = rc.makeSphere(10, -10, 0, 0, 0, 1)
	rc.s3 = rc.makeSphere(-10, -10, 0, 1, 1, 0)

	// Add a banner to show the currently selected grid location.
	_, _, w, h := rc.eng.Size()
	over := rc.eng.AddScene(vu.VO)
	over.Set2D()
	over.Cam().SetOrthographic(0, float64(w), 0, float64(h), 0, 10)
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
	p.SetForm(vu.NewSphere(sz))
	p.SetLocation(x, y, z)
	p.SetScale(sz, sz, sz)
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
	rc.cam.SetPerspective(fov, ratio, near, far)
}

// raycast checks which grid tile is selected on a mouse click. It gets
// the picking ray direction and then intersect the ray against the
// geometry in world space.
func (rc *rctag) raycast(mx, my int) {
	rx, ry, rz := rc.cam.Ray(mx, my, rc.ww, rc.wh)
	ray := vu.NewRay(rx, ry, rz)
	ray.World().SetLoc(rc.cam.Location()) // camera is ray origin.

	// collide the ray with the plane and get the world-space contact point on hit.
	if hit, x, y, z := vu.Cast(ray, rc.floor.Form()); hit {
		bot := &lin.V3{-rc.fsize, -rc.fsize, 0}
		top := &lin.V3{rc.fsize, rc.fsize, 0}

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
		println("missed plane")
		rc.hilite.SetVisible(false) // missed the plane entirely.
		rc.banner.SetVisible(false)
	}
}

// hovercast checks the sphere each update and turns the spheres a different
// colour when the mouse is over them.
func (rc *rctag) hovercast(mx, my int) {
	rx, ry, rz := rc.cam.Ray(mx, my, rc.ww, rc.wh)
	ray := vu.NewRay(rx, ry, rz)
	ray.World().SetLoc(rc.cam.Location())
	parts := []vu.Part{rc.s0, rc.s1, rc.s2, rc.s3}
	colors := []rgb{rgb{1, 0, 0}, rgb{0, 1, 0}, rgb{0, 0, 1}, rgb{1, 1, 0}}
	for cnt, p := range parts {
		if hit, _, _, _ := vu.Cast(ray, p.Form()); hit {
			p.Role().SetKd(1, 1, 1)
		} else {
			p.Role().SetKd(colors[cnt].R, colors[cnt].G, colors[cnt].B)
		}
	}
}

type rgb struct{ R, G, B float64 } // rgb holds a colour.
