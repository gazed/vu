// Copyright © 2014-2016 Galvanized Logic. All rights reserved.
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
//
// The key idea with is that mouse screen positioning can be interpreted as
// interacting with elements of a 3D scene. See vu/camera.go which holds
// the inverse matricies needed to transform screen space to world space.
//
// CONTROLS:
//   Lm    : show mouse hit
func rc() {
	rc := &rctag{}
	if err := vu.New(rc, "Ray Cast", 400, 100, 800, 600); err != nil {
		log.Printf("rc: error starting engine %s", err)
	}
	defer catchErrors()
}

// Globally unique "tag" for this example.
type rctag struct {
	cam    *vu.Camera // main 3D camera.
	ui     *vu.Camera // 2D overlay camera.
	fsize  float64    // floor size in world space units.
	gsize  float64    // grid size: number visible/virtual tiles in floor image.
	ww, wh int        // window dimensions.
	floor  *vu.Pov    // plane for raycast testing
	hilite *vu.Pov    // tracks which tile is currently selected.
	s0, s1 *vu.Pov    // spheres for raycast testing.
	s2, s3 *vu.Pov    // spheres for raycast testing.
	banner *vu.Pov    // shows selected grid locations.
}

// Create is the engine callback for initial asset creation.
func (rc *rctag) Create(eng vu.Eng, s *vu.State) {
	top := eng.Root().NewPov()
	rc.cam = top.NewCam()
	rc.cam.SetPitch(45)      // Tilt the camera and
	rc.cam.SetAt(0, -14, 14) // ...point directly at 0, 0, 0
	rc.gsize = 32            // 4x8 ie. image is 4x4 grid, tile.obj is oversampled by 8.

	// The ray cast target is a plane displaying the image of a 32x32 grid.
	rc.fsize = 10.0                                 // 2x2 plane to 20x20 plane.
	rc.floor = top.NewPov()                         // create the floor.
	rc.floor.NewBody(vu.NewPlane(0, 0, -1))         // the floors ray intersect shape.
	rc.floor.SetScale(rc.fsize, rc.fsize, 0)        // scale the model to fsize.
	rc.floor.NewModel("uv", "msh:tile", "tex:tile") // put the image on the floor.
	// TODO remove m.RepeatTex("tile")

	// create a selected tile tracker.
	rc.hilite = top.NewPov().SetScale(0.625, 0.625, 0.001) // scale to cover a single tile.
	rc.hilite.NewModel("uv", "msh:icon", "tex:image")

	// Put spheres at the floor corners.
	rc.s0 = rc.makeSphere(top, 10, 10, 0, 1, 0, 0)
	rc.s1 = rc.makeSphere(top, -10, 10, 0, 0, 1, 0)
	rc.s2 = rc.makeSphere(top, 10, -10, 0, 0, 0, 1)
	rc.s3 = rc.makeSphere(top, -10, -10, 0, 1, 1, 0)

	// Add a banner to show the currently selected grid location.
	top2D := eng.Root().NewPov()
	rc.ui = top2D.NewCam().SetUI()
	rc.banner = top2D.NewPov()
	rc.banner.SetAt(100, 100, 0)
	rc.banner.Cull = true
	rc.banner.NewLabel("uv", "lucidiaSu22", "lucidiaSu22White").SetStr("Overlay Text")

	// set non default engine state.
	eng.Set(vu.Color(0.2, 0.2, 0.2, 1.0))
	rc.resize(s.W, s.H)
}

// makeSphere creates a sphere at the given x, y, z location and with
// the given r, g, b color.
func (rc *rctag) makeSphere(parent *vu.Pov, x, y, z float64, r, g, b float32) *vu.Pov {
	sz := 0.5
	sp := parent.NewPov()
	sp.NewBody(vu.NewSphere(sz))
	sp.SetAt(x, y, z)
	sp.SetScale(sz, sz, sz)
	model := sp.NewModel("solid", "msh:sphere")
	model.SetUniform("kd", r, g, b)
	return sp
}

// Update is the engine frequent user-input/state-update callback.
func (rc *rctag) Update(eng vu.Eng, in *vu.Input, s *vu.State) {
	if in.Resized {
		rc.resize(s.W, s.H)
	}
	for press, down := range in.Down {
		switch press {
		case vu.KLm:
			if down == 1 {
				rc.raycast(in.Mx, in.My)
			}
		}
	}
	rc.hovercast(in.Mx, in.My)
}

// resize handles user screen/window changes.
func (rc *rctag) resize(ww, wh int) {
	rc.ww, rc.wh = ww, wh
	fov, ratio, near, far := 60.0, float64(ww)/float64(wh), 0.1, 500.0
	rc.cam.SetPerspective(fov, ratio, near, far)
	rc.ui.SetOrthographic(0, float64(ww), 0, float64(wh), 0, 10)
}

// raycast checks which grid tile is selected on a mouse click. It gets
// the picking ray direction and then intersect the ray against the
// geometry in world space.
func (rc *rctag) raycast(mx, my int) {
	rx, ry, rz := rc.cam.Ray(mx, my, rc.ww, rc.wh)
	ray := vu.NewRay(rx, ry, rz)
	ray.World().SetLoc(rc.cam.At()) // camera is ray origin.

	// collide the ray with the plane and get the world-space contact point on hit.
	if hit, x, y, z := vu.Cast(ray, rc.floor.Body()); hit {
		bot := &lin.V3{X: -rc.fsize, Y: -rc.fsize, Z: 0}
		top := &lin.V3{X: rc.fsize, Y: rc.fsize, Z: 0}

		// check if the plane hit was within the floor area.
		if x >= bot.X && x <= top.X && y >= bot.Y && y <= top.Y {

			// place a marker where the mouse hit.
			rc.hilite.SetAt(x, y, z)
			rc.hilite.Cull = false
			rc.banner.Cull = false
			if model := rc.banner.Model(); model != nil {
				// adjust and display grid coordinates. Map x, y to 0:31
				xsize, ysize := top.X-bot.X, top.Y-bot.Y
				gx := int(((x * 2 / xsize) + 1) / 2 * rc.gsize)
				gy := int(((y * 2 / ysize) + 1) / 2 * rc.gsize)
				model.SetStr(fmt.Sprintf("%d:%d", gx, gy))
			}
		} else {
			rc.hilite.Cull = true // missed the grid.
			rc.banner.Cull = true
		}
	} else {
		log.Printf("Missed plane.")
		rc.hilite.Cull = true // missed the plane entirely.
		rc.banner.Cull = true
	}
}

// hovercast checks the sphere each update and turns the spheres a different
// color when the mouse is over them.
func (rc *rctag) hovercast(mx, my int) {
	rx, ry, rz := rc.cam.Ray(mx, my, rc.ww, rc.wh)
	ray := vu.NewRay(rx, ry, rz)
	ray.World().SetLoc(rc.cam.At())
	parts := []*vu.Pov{rc.s0, rc.s1, rc.s2, rc.s3}
	colors := []rgb{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}, {1, 1, 0}}
	for cnt, p := range parts {
		if hit, _, _, _ := vu.Cast(ray, p.Body()); hit {
			p.Model().SetUniform("kd", 1, 1, 1)
		} else {
			rgb := colors[cnt]
			p.Model().SetUniform("kd", rgb.R, rgb.G, rgb.B)
		}
	}
}

type rgb struct{ R, G, B float32 } // rgb holds a color.
