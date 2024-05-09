// Copyright © 2024 Galvanized Logic Inc.

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
//   - circle primitive shader and vertex line circle
//   - generatede icospheres.
//
// CONTROLS:
//   - W,S    : move forward, back
//   - A,D    : move left, right
//   - C,Z    : move up, down
//   - RMouse : look around
//   - Q      : quit and close window.
func ps() {
	ps := &mhtag{}

	defer catchErrors()
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

	// import assets from asset files.
	// This creates the assets referenced by the models below.
	// Note that circle and quad meshes engine defaults.
	eng.ImportAssets("circle.shd", "lines.shd", "pbr0.shd")

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

	// create and draw an icosphere. At the lowest resolution this
	// looks bad because the normals are shared where vertexes are
	// part of multiple triangles.
	eng.GenIcosphere(0)
	s0 := ps.scene.AddModel("shd:pbr0", "msh:icosphere0")
	s0.SetAt(-3, 0, -10).SetColor(0, 0, 1, 1).SetMetallicRoughness(true, 0.2)

	// a higher resolution icosphere starts to look ok with lighting.
	eng.GenIcosphere(4)
	s2 := ps.scene.AddModel("shd:pbr0", "msh:icosphere4")
	s2.SetAt(+3, 0, -10).SetColor(0, 1, 0, 1).SetMetallicRoughness(true, 0.2)

	eng.Run(ps) // does not return while example is running.
}

// Globally unique "tag" that encapsulates example specific data.
type pstag struct {
	scene  *vu.Entity
	mx, my int32   // mouse position
	pitch  float64 // Up-down look direction.
	yaw    float64 // Left-right look direction.
}

// Update is the application engine callback.
func (ps *pstag) Update(eng *vu.Engine, in *vu.Input, delta time.Duration) {
	// react to one time press events.
	for press := range in.Pressed {
		switch press {
		case vu.KQ:
			// quit if Q is pressed
			eng.Shutdown()
			return
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
