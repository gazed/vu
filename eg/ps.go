// SPDX-FileCopyrightText : © 2014-2025 Galvanized Logic Inc.
// SPDX-License-Identifier: BSD-2-Clause

package main

import (
	"log/slog"
	"time"

	"github.com/gazed/vu"
)

// ps primitive shapes explores creating geometric shapes and standard
// shape primitives using shaders. This example demonstrates:
//   - loading assets.
//   - creating a 3D scene.
//   - controlling scene camera movement.
//   - draw circles primitives, one from a shader, one from lines.
//   - draw 2D box line primitive.
//
// CONTROLS:
//   - W,S    : move forward, back
//   - A,D    : move left, right
//   - C,Z    : move up, down
//   - RMouse : look around
//   - Q      : quit and close window.
func ps() {
	defer catchErrors()
	ps := &pstag{}
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

	// Run will call Load once and then call Update each engine tick.
	eng.Run(ps, ps) // does not return while example is running.
}

// Globally unique "tag" that encapsulates example specific data.
type pstag struct {
	scene  *vu.Entity // 3D scene
	ui     *vu.Entity // 2D scene
	mx, my int32      // mouse position
	pitch  float64    // Up-down look direction.
	yaw    float64    // Left-right look direction.
}

// Load is the one time startup engine callback to create initial assets.
func (ps *pstag) Load(eng *vu.Engine) error {
	// import assets from asset files.
	// This creates the assets referenced by the models below.
	// Note that circle and quad meshes engine defaults.
	eng.ImportAssets("circle.shd", "lines.shd", "lines2D.shd", "pbr0.shd", "col3D.shd")

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

	// draw a 3D quad.
	q1 := ps.scene.AddModel("shd:col3D", "msh:quad")
	q1.SetAt(0, 0, -10).SetScale(5, 5, 1)
	q1.SetColor(0, 0, 1, 1) // blue

	// Add a 2D scene and draw a frame using a 2D line shader.
	ps.ui = eng.AddScene(vu.Scene2D)
	frame := ps.ui.AddModel("shd:lines2D", "msh:frame")
	frame.SetAt(800, 450, 0).SetScale(900, 500, 0).SetColor(1, 1, 1, 1)

	// draw a 2D circle
	c4 := ps.ui.AddModel("shd:lines2D", "msh:circle2D")
	c4.SetAt(800, 450, 0).SetScale(150, 150, 0)
	c4.SetColor(1, 0, 1, 1)
	return nil
}

// Update is the ongoing engine callback for rendering.
func (ps *pstag) Update(eng *vu.Engine, in *vu.Input, delta time.Duration) {
	// react to one time press events.
	for press := range in.Pressed {
		switch press {
		case vu.KQ:
			// quit if Q is pressed
			eng.Shutdown()
			return
		case vu.KF11:
			eng.ToggleFullscreen()
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
