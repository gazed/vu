// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package main

import (
	"log"
	"time"
	"vu"
	"vu/math/lin"
)

// cm, collision motion, shows how controlling a physics body works by applying
// velocity pushes instead of updating the body's location. It also demonstrates
// how collisions with another physics body can affect camera/player movement.
func cm() {
	cm := &cmtag{}
	var err error
	if cm.eng, err = vu.New("Collision Motion", 400, 100, 800, 600); err != nil {
		log.Printf("co: error intitializing engine %s", err)
		return
	}
	cm.eng.SetDirector(cm)
	defer cm.eng.Shutdown()
	cm.eng.Action()
}

// Globally unique "tag" that encapsulates example specific data.
type cmtag struct {
	eng       vu.Engine
	scene     vu.Scene
	bod       vu.Part
	holdoff   time.Duration // Time in milliseconds.
	last      time.Time     // Last time command was run.
	pmx, pmy  int           // Previous mouse position.
	switchCam bool          // Control camera viewpoint.
	turn      float64       // Amount of rotation about Y axis.
	look      float64       // Amount of rotation about X axis.
}

// Create is the engine intialization callback. Create the physics objects.
func (cm *cmtag) Create(eng vu.Engine) {
	cm.holdoff, _ = time.ParseDuration("1000ms")
	cm.last = time.Now()
	cm.scene = eng.AddScene(vu.VF)
	cm.scene.SetPerspective(60, float64(800)/float64(600), 0.1, 50)
	cm.scene.SetLightLocation(0, 5, 0)
	cm.scene.SetLightColour(0.4, 0.7, 0.9)
	cm.scene.SetViewLocation(0, 10, 25)

	// add a base for the other physics bodies to rest on.
	slab := cm.scene.AddPart()
	slab.SetLocation(0, 1, 0)
	slab.SetFacade("cube", "gouraud").SetMaterial("floor")
	slab.SetScale(100, 50, 100)
	slab.SetBody(vu.Box(50, 25, 50), 0, 0.4)
	slab.SetLocation(0, -25, 0)

	// create a fixed part.
	box := cm.scene.AddPart()
	box.SetLocation(0, 1, 0)
	box.SetFacade("cube", "gouraud").SetMaterial("cube")
	box.SetScale(5, 5, 5)
	box.SetBody(vu.Box(2.5, 2.5, 2.5), 0, 0) // sized to match scaling.

	// create a physics body that can be associated with a camera,
	cm.bod = cm.scene.AddPart()
	cm.bod.SetFacade("sphere", "gouraud").SetMaterial("sphere")
	cm.bod.SetBody(vu.Sphere(1), 1, 0)
	cm.bod.SetLocation(-6, 2, -2)

	// set some constant state.
	cm.eng.Enable(vu.BLEND, true)
	cm.eng.Enable(vu.CULL, true)
	cm.eng.Enable(vu.DEPTH, true)
	cm.eng.Color(0.1, 0.1, 0.1, 1.0)
	return
}

// Update is the regular engine callback.
func (cm *cmtag) Update(input *vu.Input) {
	if input.Resized {
		cm.resize()
	}

	// own the rotation.
	dir := lin.NewQ().SetAa(0, 1, 0, lin.Rad(cm.turn))
	cm.bod.SetRotation(dir.GetS())
	if cm.switchCam {

		// apply physics location to the camera.
		x, y, z := cm.bod.Location()
		cm.scene.SetViewLocation(x, y, z)
		cm.scene.SetViewRotation(dir.GetS())
		cm.scene.SetViewTilt(-cm.look)
	} else {

		// keep camera at a fixed location.
		cm.scene.SetViewLocation(0, 10, 25)
		cm.scene.SetViewRotation(0, 0, 0, 1)
		cm.scene.SetViewTilt(0)
	}

	// handle user requests.
	for press, rel := range input.Down {
		switch press {
		case "W":
			cm.move(zaxis, -1, rel < 0) // back and forth.
		case "S":
			cm.move(zaxis, 1, rel < 0)
		case "A":
			cm.move(xaxis, -1, rel < 0) // left and right
		case "D":
			cm.move(xaxis, 1, rel < 0)

		// switch camera between first and fixed third person.
		case "C":
			if time.Now().After(cm.last.Add(cm.holdoff)) {
				cm.switchCam = !cm.switchCam
				cm.last = time.Now()
			}
		}
	}

	// apply rotation based on mouse changes.
	xdiff := cm.pmx - input.Mx
	ydiff := cm.pmy - input.My
	cm.spin(yaxis, xdiff) // rotate left/right (around y-axis)
	cm.spin(xaxis, ydiff) // rotate up/down (around x-axis)
	cm.pmx, cm.pmy = input.Mx, input.My
}

// resize handles user screen/window changes.
func (cm *cmtag) resize() {
	x, y, width, height := cm.eng.Size()
	cm.eng.Resize(x, y, width, height)
	cm.scene.SetPerspective(60, float64(width)/float64(height), 0.1, 50)
}

// Used as parameters to indicate where movement is happening.
const (
	xaxis = iota
	yaxis
	zaxis
)

// move adjusts linear velocity.  This is how physics bodies should be repositioned.
func (cm *cmtag) move(axis, dir int, stop bool) {
	if stop {
		cm.bod.Stop() // remove linear velocities.
	} else {
		mov := 0.5 * float64(dir)
		switch axis {
		case xaxis:
			cm.bod.Move(mov, 0, 0)
		case zaxis:
			cm.bod.Move(0, 0, mov)
		}
	}
}

// spin tracks rotation/orientation/direction changes in separate variables
// They are applied during update.
func (cm *cmtag) spin(axis, diff int) {
	if diff != 0 {
		switch axis {
		case xaxis:
			cm.look += 0.1 * float64(diff)
			if cm.look > 90.0 {
				cm.look = 90.0
			}
			if cm.look < -90.0 {
				cm.look = -90.0
			}
		case yaxis:
			cm.turn += 0.5 * float64(diff)
			if cm.turn > 360 || cm.turn < -360 {
				cm.turn = 0
			}
		}
	}
}
