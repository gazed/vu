// Copyright © 2024 Galvanized Logic Inc.

package main

import (
	"fmt"
	"log/slog"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/gazed/vu"
	"github.com/gazed/vu/load"
	"github.com/gazed/vu/math/lin"
)

// is draws instanced stars showing that it is possible to render
// many models (over 9000) in a single draw call.
// This example demonstrates:
//   - loading assets.
//   - creating a 3D scene with models.
//   - creating a 2D scene with labels.
//   - changing the scene camera direction.
//   - instanced models - drawing many models in a single draw call.
//   - billboard shader where the models always face the camera.
//   - ray picking - using the mouse to interact with a 3D model.
//
// Controls
//   - RMouse : look around
//   - Q      : quit and close window.
func is() {
	is := &istag{ww: 1600, wh: 900, names: map[int]*vu.Entity{}}

	defer catchErrors()
	eng, err := vu.NewEngine(
		vu.Windowed(),
		vu.Title("Instanced Stars"),
		vu.Size(200, 200, int32(is.ww), int32(is.wh)),
		vu.Background(0.0, 0.0, 0.0, 1.0),
	)
	if err != nil {
		slog.Error("is: engine start", "err", err)
		return
	}

	// load the bright star data.
	iData := []load.Buffer{}
	is.stars, iData, err = is.loadBrightStars()

	// import assets from asset files.
	// This creates the assets referenced by the models below.
	eng.ImportAssets("bbinst.shd", "bboard.shd", "star.png", "ring.png") // 3D assets
	eng.ImportAssets("label.shd", "lucidiaSu18.fnt", "lucidiaSu18.png")  // 2D assets

	// The scene holds the cameras and lighting information
	// and acts as the root for all models added to the scene.
	is.scene3D = eng.AddScene(vu.Scene3D)
	is.scene3D.Cam().SetClip(0.1, 100000)

	// show a picked star by drawing a ring around it.
	// Create before instance data so it is drawn after the instanced data.
	// This is because the normal transparency sorting doesn't work with instanced
	// data since instanced data does not have a distance to the camera..
	is.pick = is.scene3D.AddModel("shd:bboard", "msh:quad", "tex:color:ring")
	is.pick.Cull(true) // hide until mouse is over a star

	// one draw call to draw over 9000 stars.
	s1 := is.scene3D.AddInstancedModel("shd:bbinst", "msh:quad", "tex:color:star")
	s1.SetInstanceData(eng, uint32(len(is.stars)), iData)

	// create a 2D scene to show the star name when holding the mouse over a star.
	is.scene2D = eng.AddScene(vu.Scene2D)

	eng.SetResizeListener(is)
	eng.Run(is) // does not return while example is running.
}

// Globally unique "tag" that encapsulates example specific data.
type istag struct {
	scene3D *vu.Entity   // stars
	scene2D *vu.Entity   // star names
	pick    *vu.Entity   // highlight a picked star
	ww, wh  int          // window width, height
	mx, my  int32        // mouse position
	pitch   float64      // Up-down look direction.
	yaw     float64      // Left-right look direction.
	stars   []brightStar // bright star data.

	// star names are created as needed and then reused.
	names map[int]*vu.Entity
}

// Resize is called by the engine when the window size changes.
func (is *istag) Resize(windowWidth, windowHeight uint32) {
	// keep these updated for ray picking.
	is.ww, is.wh = int(windowWidth), int(windowHeight)
}

// Update is the application engine callback.
func (is *istag) Update(eng *vu.Engine, in *vu.Input, delta time.Duration) {
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
	xdiff, ydiff := in.Mx-is.mx, in.My-is.my // mouse move differences...
	is.mx, is.my = in.Mx, in.My              // ... from last mouse location.

	// react to continuous press events.
	lookSpeed := 15.0 * delta.Seconds()
	cam := is.scene3D.Cam()
	for press := range in.Down {
		switch press {
		case vu.KMR:
			if ydiff != 0 {
				is.pitch += float64(ydiff) * lookSpeed
			}
			if xdiff != 0 {
				is.yaw += float64(-xdiff) * lookSpeed
			}
			cam.SetPitch(is.pitch)
			cam.SetYaw(is.yaw)
		}
	}

	// raycast: show the label for a star if the mouse is over that star.
	for _, v := range is.names {
		v.Cull(true)
	}
	is.pick.Cull(true)
	if starID := is.hitStar(is.stars); starID > -1 {
		s := &is.stars[starID]
		if name, ok := is.names[starID]; !ok {
			name = is.scene2D.AddLabel(s.Name, 0, "shd:label", "fnt:lucidiaSu18", "tex:color:lucidiaSu18")
			name.SetAt(20, 20, 0).SetColor(1, 1, 0, 1) // yellow label
			is.names[starID] = name
		} else {
			name.Cull(false)
		}

		// move the pick ring closer to the camera (origin)
		// in order blend the star pixels with the transparent center ring pixels.
		x, y, z := s.Locus[0], s.Locus[1], s.Locus[2]
		dist2 := x*x + y*y + z*z // distance squared
		at := lin.NewV3().SetS(x, y, z)
		at.Sub(at, lin.NewV3().SetS(-s.Locus[0], -s.Locus[1], -s.Locus[2]).Unit())
		scale := 1.0
		if dist2 < 1000 {
			scale = 0.2 // scale down for really close stars
		}
		is.pick.SetAt(at.X, at.Y, at.Z).SetScale(scale, scale, 1).Cull(false)
	}
}

// hitStar returns a star that is under the mouse.
func (is *istag) hitStar(stars []brightStar) (starID int) {
	cam := is.scene3D.Cam()

	// get unit ray direction from camera.
	ray := lin.NewV3().SetS(cam.Ray(int(is.mx), int(is.my), is.ww, is.wh))
	if ray.X == 0 && ray.Y == 0 && ray.Z == 0 {
		return -1 // mouse is not not over window.
	}
	sphere := lin.NewV3()
	for i, s := range stars {
		sphere.SetS(s.Locus[0], s.Locus[1], s.Locus[2])
		hit, _, _, _ := cam.RayCastSphere(ray, sphere, 0.5)
		if hit {
			return i
		}
	}
	return -1
}

// loadBrightStars converts the yaml star data to instanced GPU buffer data.
func (is *istag) loadBrightStars() (stars []brightStar, data []load.Buffer, err error) {
	bytes, err := load.DataBytes("bright_star.yaml")
	if err != nil {
		return stars, data, fmt.Errorf("is:loadBrightStars %w", err)
	}
	stars = make([]brightStar, 9101)
	err = yaml.Unmarshal(bytes, &stars)
	if err != nil {
		return stars, data, fmt.Errorf("is:loadBrightStars %w", err)
	}
	positions := []float32{}
	colors := []float32{}
	scales := []float32{}
	for _, s := range stars {
		x, y, z := float32(s.Locus[0]), float32(s.Locus[1]), float32(s.Locus[2])
		if x == 0 && y == 0 && z == 0 {
			fmt.Printf("%d %+v\n", s.ID, s)
		}
		positions = append(positions, x, y, z)
		r, g, b := float32(s.Color[0]), float32(s.Color[1]), float32(s.Color[2])
		colors = append(colors, r, g, b)

		// tweak scale based on lumens and distance.
		scale := float32(1.0)
		dist2 := x*x + y*y + z*z // distance squared
		if dist2 < 1000 {
			scale = 0.1 // scale down really close stars
		}

		// scale up a bit based on lumens.
		lumens := float32(s.Lumen)
		if lumens > 1000 {
			lumens = 1000
		}
		scale += lumens / 1000.0
		scales = append(scales, scale)
	}

	// convert the bright star data to instanced data.
	data = make([]load.Buffer, load.InstanceTypes)
	data[load.InstanceLocus] = load.F32Buffer(positions, 3)
	data[load.InstanceColors] = load.F32Buffer(colors, 3)
	data[load.InstanceScales] = load.F32Buffer(scales, 1)
	return stars, data, nil
}

// brightStar is used to import the bright_star data.
type brightStar struct {
	ID    int       `yaml:"id"`    // bright star ID
	Name  string    `yaml:"name"`  // bright star name, proper name if available
	Locus []float64 `yaml:"locus"` // xyz in light years
	Color []float64 `yaml:"color"` // rgb in range 0-1
	Lumen float64   `yaml:"lumen"` // luminosity used for scaling
}
