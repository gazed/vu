// SPDX-FileCopyrightText : © 2026 Galvanized Logic Inc.
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"image"
	"image/draw"
	"log/slog"
	"time"

	"github.com/gazed/vu"
	"github.com/gazed/vu/math/lin"
)

// sn draws a sphere using shader noise.
// This example demonstrates:
//   - 3D shader noise.
//   - loading assets.
//   - creating a 3D scene.
//   - binary GLTF (GLB) imports for mesh.
//
// CONTROLS:
//   - A,D    : spin left, right
//   - {,}    : adjust noise scale
//   - -,+    : adjust noise frequency
//   - Q      : quit and close window.
func sn() {
	defer catchErrors()
	sn := &sntag{scale: 2.0, frequency: 1.5}
	eng, err := vu.NewEngine(
		vu.Windowed(),
		vu.Title("Shader Noise"),
		vu.Size(200, 200, 1600, 900),
		vu.Background(0.01, 0.01, 0.01, 1.0),
	)
	if err != nil {
		slog.Error("mh: engine start", "err", err)
		return
	}

	// Run will call Load once and then call Update each engine tick.
	eng.Run(sn, sn) // does not return while example is running.
}

// Globally unique "tag" that encapsulates example specific data.
type sntag struct {
	scene *vu.Entity // 3D
	ball  *vu.Entity // sphere model

	// updatable label.
	ui         *vu.Entity   // 2D scene
	stats      *vu.Entity   // noise parameter label
	statsReady bool         // true when stat assets are loaded.
	text       *image.NRGBA // the stats image texture.

	// rotation control
	pos *lin.V3 // initial location.
	rot float64 // rotation around origin.

	// noise control
	scale     float32 // noise scaling.
	frequency float32 // noise frequency.
}

// Load is the one time startup engine callback to create initial assets.
func (sn *sntag) Load(eng *vu.Engine) error {

	// import assets from asset files.
	// This creates the assets referenced by the models below.
	eng.ImportAssets("noise.shd", "sphere.glb")    // 3D assets
	eng.ImportAssets("image2D.shd", "22:hack.ttf") // 2D assets

	// The scene holds the cameras and lighting information
	// and acts as the root for all models added to the scene.
	sn.scene = eng.AddScene(vu.Scene3D)
	sn.pos = lin.NewV3().SetS(0, 0, 4)
	sn.scene.Cam().SetAt(sn.pos.X, sn.pos.Y, sn.pos.Z)

	// add a sphere at the origin.
	const sphere_radius = 1.2849 // from blender
	sn.ball = sn.scene.AddModel("shd:noise", "msh:sphere")
	sn.ball.SetScale(2, 2, 2).SetAt(0, 0, 0)
	sn.ball.SetModelUniform("f4", []float32{sn.scale, sn.frequency, 0.0, 0.0})

	// view the noise parameter changes.
	sn.ui = eng.AddScene(vu.Scene2D)

	// double buffer the stats texture, one for updating, one for rendering.
	imgWidth, imgHeight := 512, 192
	sn.text = image.NewNRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	sn.stats = sn.ui.AddModel("shd:image2D", "msh:icon", "fnt:hack22")
	sn.stats.SetAt(280, 120, 0).SetScale(float64(imgWidth), float64(imgHeight), 0.0)
	sn.stats.AddUpdatableTexture(eng, "stats", sn.text)
	return nil
}

// Update is the ongoing engine callback.
func (sn *sntag) Update(eng *vu.Engine, in *vu.Input, delta time.Duration) {

	// wait for the font to load before the initial text update.
	// Afterwards only need to update if it changes.
	if !sn.statsReady {
		sn.statsReady = sn.updateStats(eng) // initialize stats label.
	}

	// react to one time press events.
	for press := range in.Pressed {
		switch press {
		case vu.KQ:
			// quit if Q is pressed
			eng.Shutdown()
			return
		}
	}

	// react to continuous press events.
	cam := sn.scene.Cam()
	lookSpeed := 250 * delta.Seconds()
	for press := range in.Down {
		switch press {
		case vu.KA:
			// rotate camera left around center
			sn.rot -= lookSpeed
			transformAroundOrigin := lin.NewT().SetLoc(0, 0, 0).SetAa(0, 1, 0, lin.Rad(sn.rot))
			at := transformAroundOrigin.App(lin.NewV3().Set(sn.pos))
			cam.SetAt(at.X, at.Y, at.Z)
			cam.SetYaw(sn.rot)
		case vu.KD:
			// rotate camera right around center
			sn.rot += lookSpeed
			transformAroundOrigin := lin.NewT().SetLoc(0, 0, 0).SetAa(0, 1, 0, lin.Rad(sn.rot))
			at := transformAroundOrigin.App(lin.NewV3().Set(sn.pos))
			cam.SetAt(at.X, at.Y, at.Z)
			cam.SetYaw(sn.rot)

		// play with the noise parameters.
		case vu.KRBkt: // { frequency up
			sn.scale = min(5.0, sn.scale+0.025)
			sn.ball.SetModelUniform("f4", []float32{sn.scale, sn.frequency, 0.0, 0.0})
			sn.updateStats(eng)
		case vu.KLBkt: // } frequency down
			sn.scale = max(1.0, sn.scale-0.025)
			sn.ball.SetModelUniform("f4", []float32{sn.scale, sn.frequency, 0.0, 0.0})
			sn.updateStats(eng)
		case vu.KEqual: // + scale up
			sn.frequency = min(4.0, sn.frequency+0.015)
			sn.ball.SetModelUniform("f4", []float32{sn.scale, sn.frequency, 0.0, 0.0})
			sn.updateStats(eng)
		case vu.KMinus: // - scale down
			sn.frequency = max(1.5, sn.frequency-0.015)
			sn.ball.SetModelUniform("f4", []float32{sn.scale, sn.frequency, 0.0, 0.0})
			sn.updateStats(eng)
		}
	}
}

// update the statistics label.
func (sn *sntag) updateStats(eng *vu.Engine) bool {
	draw.Draw(sn.text, sn.text.Bounds(), image.Transparent, image.Point{}, draw.Src)
	label := fmt.Sprintf("scale([]):%2.2f frequency(-+):%2.3f", sn.scale, sn.frequency)
	ready := sn.stats.WriteImageText("hack22", label, 0, 0, sn.text)
	sn.stats.UpdateTexture(eng, sn.text)
	return ready == nil
}
