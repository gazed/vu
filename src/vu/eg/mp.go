// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package main

import (
	"log"
	"vu"
	"vu/math/lin"
	"vu/physics"
)

// mp demonstrates "mouse picking" which is the ability of selecting a 3D object
// using the mouse. Currently the vu engine does not fully support mouse picking, but the
// example is used as a placeholder for eventual inclusion.  The only part that is in
// the engine is the calculation of the pick ray (see vu/picking.go for more information):
//     r1 := vu.PickDir(mx, my, fov, w, h, n, f, mv)
// Test by clicking on and off of the target while moving around.  Mouse picking always
// knows when the object has been selected.
func mp() {
	mp := &mptag{}
	var err error
	if mp.eng, err = vu.New("Line", 400, 100, 800, 600); err != nil {
		log.Printf("mp: error intitializing engine %s", err)
		return
	}
	mp.run = 10            // move so many cubes worth in one second.
	mp.spin = 270          // spin so many degrees in one second.
	mp.eng.SetDirector(mp) // override user input handling.
	mp.stagePlay()
	defer mp.eng.Shutdown()
	mp.eng.Action()
}

// Globally unique "tag" for this example.
type mptag struct {
	eng    *vu.Eng
	scene  vu.Scene
	target vu.Part
	run    float32
	spin   float32
}

// Create the shooting range
func (mp *mptag) stagePlay() {
	mp.scene = mp.eng.AddScene(vu.VP)
	mp.scene.SetPerspective(60, float32(800)/float32(600), 0.1, 50)
	mp.scene.SetLightLocation(0, 10, 0)
	mp.scene.SetLightColour(0.4, 0.7, 0.9)
	mp.scene.SetViewLocation(0, 0, 3)

	// load the object for simulating a shot.
	mp.target = mp.scene.AddPart()
	mp.target.SetLocation(0, 0, -4)
	mp.target.SetFacade("sphere", "gouraud", "floor")

	// set some constant state.
	mp.eng.Enable(vu.BLEND, true)
	mp.eng.Enable(vu.CULL, true)
	mp.eng.Color(0.1, 0.1, 0.1, 1.0)
	return
}

// Handle engine callbacks.
func (mp *mptag) Focus(focus bool) {}
func (mp *mptag) Resize(x, y, width, height int) {
	mp.eng.ResizeViewport(x, y, width, height)
	mp.scene.SetPerspective(60, float32(width)/float32(height), 0.1, 50)
}
func (mp *mptag) Update(pressed []string, gt, dt float32) {
	moveDelta := mp.eng.Dt
	for _, p := range pressed {
		switch p {
		case "W":
			mp.scene.MoveView(0, 0, moveDelta*-mp.run)
		case "S":
			mp.scene.MoveView(0, 0, moveDelta*mp.run)
		case "A":
			mp.scene.PanView(vu.YAxis, mp.eng.Dt*mp.spin)
		case "D":
			mp.scene.PanView(vu.YAxis, mp.eng.Dt*-mp.spin)
		case "Lm":
			xm, ym := mp.eng.Xm, mp.eng.Ym
			if mp.picked(xm, ym, 45, 800, 600, 0.1, 50, mp.target) {
				println("hit", xm, ym)
			} else {
				println("miss", xm, ym)
			}
		}
	}
}

// TODO generalize and move into the engine. This code has not yet made it into the
// engine as nothing else needs picking at the momement.
func (mp *mptag) picked(mx, my int, fov, w, h, n, f float32, target vu.Part) bool {
	rx, ry, rz, rw := mp.scene.ViewRotation()
	rot := &lin.Q{rx, ry, rz, rw}
	view := rot.Inverse().M4()
	cx, cy, cz := mp.scene.ViewLocation()
	view.TranslateL(-cx, -cy, -cz)

	rx, ry, rz, rw = target.Rotation()
	rot = &lin.Q{rx, ry, rz, rw}
	model := rot.M4()
	model.ScaleL(target.Scale())
	model.TranslateR(target.Location())

	mv := model.Mult(view)
	r1 := vu.PickDir(mx, my, fov, w, h, n, f, mv)

	// now check for collision
	x, y, z := target.Location()
	sphere := physics.Sphere(x, y, z, 1.0)
	ray := physics.Ray(cx, cy, cz, r1.X, r1.Y, r1.Z)
	c := ray.Collide(sphere)
	return c != nil
}
