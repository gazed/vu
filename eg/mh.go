// SPDX-FileCopyrightText : © 2022-2025 Galvanized Logic Inc.
// SPDX-License-Identifier: BSD-2-Clause

package main

import (
	"log/slog"
	"time"

	"github.com/gazed/vu"
)

// mh draws blender monkeys using PBR shaders.
// This example demonstrates:
//   - loading assets.
//   - creating a 3D scene.
//   - adding a light to a scene.
//   - controlling scene camera movement.
//   - shaders for Physically Based Rendering (PBR).
//   - binary GLTF (GLB) imports for mesh, texture, and material assets.
//
// CONTROLS:
//   - W,S    : move forward, back
//   - A,D    : move left, right
//   - C,Z    : move up, down
//   - RMouse : look around
//   - Q      : quit and close window.
func mh() {
	defer catchErrors()
	mh := &mhtag{}
	eng, err := vu.NewEngine(
		vu.Windowed(),
		vu.Title("Monkey Heads"),
		vu.Size(200, 200, 1600, 900),
		vu.Background(0.01, 0.01, 0.01, 1.0),
	)
	if err != nil {
		slog.Error("mh: engine start", "err", err)
		return
	}

	// Run will call Load once and then call Update each engine tick.
	eng.Run(mh, mh) // does not return while example is running.
}

// Globally unique "tag" that encapsulates example specific data.
type mhtag struct {
	scene  *vu.Entity
	mx, my int32   // mouse position
	pitch  float64 // Up-down look direction.
	yaw    float64 // Left-right look direction.
}

// Load is the one time startup engine callback to create initial assets.
func (mh *mhtag) Load(eng *vu.Engine) error {

	// import assets from asset files.
	// This creates the assets referenced by the models below.
	eng.ImportAssets("pbr0.shd", "monkey0.glb", "pbr1.shd", "monkey1.glb")

	// The scene holds the cameras and lighting information
	// and acts as the root for all models added to the scene.
	mh.scene = eng.AddScene(vu.Scene3D)

	// add one directional light. SetAt sets the direction.
	mh.scene.AddLight(vu.DirectionalLight).SetAt(-1, -2, -2)

	// add monkey heads: facing +Z by default.
	// Request model assets and render once those assets have been loaded.
	//
	// pbr0 uses fixed color:metallic:roughness values.
	mh1 := mh.scene.AddModel("shd:pbr0", "msh:monkey0", "mat:monkey0")
	mh1.SetAt(-1.5, 0, -5)
	//
	// pbr1 uses fixed a color texture and fixed metallic:roughness values.
	// The texture matches the shader "color" sampler.
	mh2 := mh.scene.AddModel("shd:pbr1", "msh:monkey1", "tex:color:monkey1", "mat:monkey1")
	mh2.SetAt(+1.5, 0, -5)
	return nil
}

// Update is the ongoing engine callback.
func (mh *mhtag) Update(eng *vu.Engine, in *vu.Input, delta time.Duration) {
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
	xdiff, ydiff := in.Mx-mh.mx, in.My-mh.my // mouse move differences...
	mh.mx, mh.my = in.Mx, in.My              // ... from last mouse location.

	// react to continuous press events.
	lookSpeed := 15.0 * delta.Seconds()
	move := 10.0 // move so many units worth in one second.
	speed := move * delta.Seconds()
	cam := mh.scene.Cam()
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
				mh.pitch = mh.limitPitch(mh.pitch + float64(-ydiff)*lookSpeed)
				cam.SetPitch(mh.pitch)
			}
			if xdiff != 0 {
				mh.yaw += float64(-xdiff) * lookSpeed
				cam.SetYaw(mh.yaw)
			}
		}
	}
}

// limitPitch ensures that look up/down is limited to 90 degrees.
// This helps reduce confusion when looking around.
func (mh *mhtag) limitPitch(pitch float64) float64 {
	switch {
	case pitch > 90:
		return 90
	case pitch < -90:
		return -90
	}
	return pitch
}
