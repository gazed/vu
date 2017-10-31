// Copyright © 2014-2017 Galvanized Logic. All rights reserved.
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
	defer catchErrors()
	if err := vu.Run(&rctag{}); err != nil {
		log.Printf("rc: error starting engine %s", err)
	}
}

// Globally unique "tag" for this example.
type rctag struct {
	scene  *vu.Ent // main 3D scene.
	ui     *vu.Ent // 2D overlay scene.
	fsize  float64 // floor size in world space units.
	gsize  float64 // grid size: number visible/virtual tiles in floor image.
	ww, wh int     // window dimensions.
	floor  *vu.Ent // plane for raycast testing
	hilite *vu.Ent // tracks which tile is currently selected.
	s0, s1 *vu.Ent // spheres for raycast testing.
	s2, s3 *vu.Ent // spheres for raycast testing.
	banner *vu.Ent // shows selected grid locations.
}

// Create is the engine callback for initial asset creation.
func (rc *rctag) Create(eng vu.Eng, s *vu.State) {
	eng.Set(vu.Title("Ray Cast"), vu.Size(400, 100, 800, 600))
	eng.Set(vu.Color(0.2, 0.2, 0.2, 1.0))
	rc.scene = eng.AddScene()
	rc.scene.Cam().SetClip(0.1, 50).SetFov(60).SetAt(0, -14, 14)
	rc.scene.Cam().SetPitch(45) // Tilt the camera
	rc.gsize = 32               // 4x8 ie. image is 4x4 grid, tile.obj is oversampled by 8.

	// The ray cast target is a plane displaying the image of a 32x32 grid.
	rc.fsize = 10.0                                  // 2x2 plane to 20x20 plane.
	rc.floor = rc.scene.AddPart()                    // create the floor.
	rc.floor.MakeBody(vu.Plane(0, 0, -1))            // the floors ray intersect shape.
	rc.floor.SetScale(rc.fsize, rc.fsize, 0)         // scale the model to fsize.
	rc.floor.MakeModel("uv", "msh:tile", "tex:tile") // put the image on the floor.

	// create a selected tile tracker.
	rc.hilite = rc.scene.AddPart().SetScale(0.625, 0.625, 0.001) // scale to cover a single tile.
	rc.hilite.MakeModel("uv", "msh:icon", "tex:image")

	// Put spheres at the floor corners.
	rc.s0 = rc.makeSphere(rc.scene.AddPart(), 10, 10, 0, 1, 0, 0)
	rc.s1 = rc.makeSphere(rc.scene.AddPart(), -10, 10, 0, 0, 1, 0)
	rc.s2 = rc.makeSphere(rc.scene.AddPart(), 10, -10, 0, 0, 0, 1)
	rc.s3 = rc.makeSphere(rc.scene.AddPart(), -10, -10, 0, 1, 1, 0)

	// // Add a banner to show the currently selected grid location.
	rc.ui = eng.AddScene().SetUI()
	rc.ui.Cam().SetClip(0, 10)
	rc.banner = rc.ui.AddPart()
	rc.banner.SetAt(100, 100, 0)
	rc.banner.Cull(true)
	rc.banner.MakeLabel("txt", "lucidiaSu22").Typeset("Overlay Text")
}

// makeSphere creates a sphere at the given x, y, z location and with
// the given r, g, b color.
func (rc *rctag) makeSphere(sp *vu.Ent, x, y, z float64, r, g, b float32) *vu.Ent {
	sz := 0.5
	sp.MakeBody(vu.Sphere(sz))
	sp.SetAt(x, y, z).SetScale(sz, sz, sz)
	model := sp.MakeModel("solid", "msh:sphere")
	model.SetUniform("kd", r, g, b)
	return sp
}

// Update is the engine frequent user-input/state-update callback.
func (rc *rctag) Update(eng vu.Eng, in *vu.Input, s *vu.State) {
	rc.ww, rc.wh = s.W, s.H
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

// raycast checks which grid tile is selected on a mouse click. It gets
// the picking ray direction and then intersect the ray against the
// geometry in world space.
func (rc *rctag) raycast(mx, my int) {
	cam := rc.scene.Cam()
	rx, ry, rz := cam.Ray(mx, my, rc.ww, rc.wh)
	ray := vu.Ray(rx, ry, rz)
	ray.World().SetLoc(cam.At()) // camera is ray origin.

	// collide the ray with the plane and get the world-space contact point on hit.
	if hit, x, y, z := rc.floor.Cast(ray); hit {
		bot := &lin.V3{X: -rc.fsize, Y: -rc.fsize, Z: 0}
		top := &lin.V3{X: rc.fsize, Y: rc.fsize, Z: 0}

		// check if the plane hit was within the floor area.
		if x >= bot.X && x <= top.X && y >= bot.Y && y <= top.Y {

			// place a marker where the mouse hit.
			rc.hilite.SetAt(x, y, z)
			rc.hilite.Cull(false)
			rc.banner.Cull(false)

			// adjust and display grid coordinates. Map x, y to 0:31
			xsize, ysize := top.X-bot.X, top.Y-bot.Y
			gx := int(((x * 2 / xsize) + 1) / 2 * rc.gsize)
			gy := int(((y * 2 / ysize) + 1) / 2 * rc.gsize)
			rc.banner.Typeset(fmt.Sprintf("%d:%d", gx, gy))
		} else {
			rc.hilite.Cull(true) // missed the grid.
			rc.banner.Cull(true)
		}
	} else {
		log.Printf("Missed plane.")
		rc.hilite.Cull(true) // missed the plane entirely.
		rc.banner.Cull(true)
	}
}

// hovercast checks the sphere each update and turns the spheres a different
// color when the mouse is over them.
func (rc *rctag) hovercast(mx, my int) {
	cam := rc.scene.Cam()
	rx, ry, rz := cam.Ray(mx, my, rc.ww, rc.wh)
	ray := vu.Ray(rx, ry, rz)
	ray.World().SetLoc(cam.At())
	parts := []*vu.Ent{rc.s0, rc.s1, rc.s2, rc.s3}
	colors := []rgb{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}, {1, 1, 0}}
	for cnt, p := range parts {
		if hit, _, _, _ := p.Cast(ray); hit {
			p.SetUniform("kd", 1, 1, 1)
		} else {
			rgb := colors[cnt]
			p.SetUniform("kd", rgb.R, rgb.G, rgb.B)
		}
	}
}

type rgb struct{ R, G, B float32 } // rgb holds a color.
