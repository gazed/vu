// Copyright © 2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package main

import (
	"log"
	"vu"
	"vu/land"
)

// tm demonstrates creating and rendering a dynamic terrain map from a
// generated height map. The intent is to mimic a landscape seen from above.
// The main engine pieces of the demo are are vu/Surface and package vu/land.
func tm() {
	tm := &tmtag{}
	tm.ww, tm.wh = 1024, 768
	var err error
	if tm.eng, err = vu.New("Terrain Map", 1200, 100, tm.ww, tm.wh); err != nil {
		log.Printf("tm: error intitializing engine %s", err)
		return
	}
	tm.eng.SetDirector(tm)  // override user input handling.
	defer tm.eng.Shutdown() // shut down the engine.
	defer catchErrors()
	tm.eng.Action()
}

// Encapsulate example specific data with a unique "tag".
type tmtag struct {
	eng     vu.Engine   // 3D engine.
	ww, wh  int         // window size.
	scene   vu.Scene    // overall scene with lighting.
	ground  vu.Part     // visible surface.
	coast   vu.Part     // shallow water plane.
	ocean   vu.Part     // deep water plane.
	world   land.Land   // height map generation.
	surface vu.Surface  // data structure used to create land.
	evo     [][]float64 // used for evolution experiments.
}

// Create is the engine intialization callback.
func (tm *tmtag) Create(eng vu.Engine) {
	tm.scene = eng.AddScene(vu.VP)
	tm.scene.SetOrthographic(0, float64(tm.ww), 0, float64(tm.wh), 0, 200)
	tm.scene.SetLightColour(0.1, 0, 0)
	tm.scene.SetLightLocation(5, 5, -5)

	// create the world surface.
	seed := int64(123)
	patchSize := 128
	tm.world = land.New(1, patchSize, seed)
	worldTile := tm.world.NewTile(1, 0, 0)
	textureRatio := 256.0 / 1024.0
	tm.surface = vu.NewSurface(patchSize, patchSize, 16, float32(textureRatio), 10)

	// create a separate surface for generating initial land textures.
	emap := land.New(1, patchSize, seed-1)
	etile := emap.NewTile(1, 0, 0)
	etopo := etile.Topo()

	// merge the land height and land texture information into a single surface.
	tm.evo = make([][]float64, patchSize)
	for x := range tm.evo {
		tm.evo[x] = make([]float64, patchSize)
	}
	numTextures := 3.0
	pts := tm.surface.Pts()
	topo := worldTile.Topo()
	for x := range topo {
		for y := range topo[x] {
			pts[x][y].Height = float32(topo[x][y])
			evolution := (etopo[x][y] + 1) * 0.5 * numTextures // (-1,1 map to 0-2), map to 0-3
			pts[x][y].Tindex = int(evolution)
			pts[x][y].Blend = float32(evolution) - float32(int(evolution))
			tm.evo[x][y] = evolution // remember for later.
		}
	}

	// Add a rendering component for the surface data.
	scale := 10.0
	tm.ground = tm.scene.AddPart().SetLocation(0, 0, -10).SetScale(scale, scale, 1)
	tm.ground.SetRole("land").AddTex("land")
	tm.ground.Role().SetMaterial("land").SetUniform("ratio", textureRatio)
	tm.surface.Update(tm.ground.Role().Mesh(), 0, 0)

	// Add water planes.
	tm.ocean = tm.scene.AddPart().SetLocation(256, 0, -10.5)
	tm.ocean.SetScale(float64(tm.ww), float64(tm.wh), 1)
	tm.ocean.SetRole("flat").SetMesh("plane").SetMaterial("blue2")
	tm.coast = tm.scene.AddPart().SetLocation(256, 0, -10)
	tm.coast.SetScale(float64(tm.ww), float64(tm.wh), 1)
	tm.coast.SetRole("flat").SetMesh("plane").SetMaterial("blue")
	return
}

// Update is the regular engine callback.
func (tm *tmtag) Update(in *vu.Input) {
	if in.Resized {
		tm.resize()
	}

	// process user presses.
	rate := 4.0
	for press, _ := range in.Down {
		switch press {

		// Change the water level.
		case "[":
			tm.ocean.Move(0, 0, 1*in.Dt)
			tm.coast.Move(0, 0, 1*in.Dt)
		case "]":
			tm.ocean.Move(0, 0, -1*in.Dt)
			tm.coast.Move(0, 0, -1*in.Dt)

		// Demonstrate evolution using a texture atlas.
		case "KP+":
			tm.evolve(0.005)
		case "KP-":
			tm.evolve(-0.005)

		// Move the map around.
		case "Ua":
			x, y, z := tm.ground.Location()
			tm.ground.SetLocation(x, y+rate, z)
		case "Da":
			x, y, z := tm.ground.Location()
			tm.ground.SetLocation(x, y-rate, z)
		case "La":
			x, y, z := tm.ground.Location()
			tm.ground.SetLocation(x-rate, y, z)
		case "Ra":
			x, y, z := tm.ground.Location()
			tm.ground.SetLocation(x+rate, y, z)
		}
	}
}

// resize handles user screen/window changes.
func (tm *tmtag) resize() {
	var x, y int
	x, y, tm.ww, tm.wh = tm.eng.Size()
	tm.eng.Resize(x, y, tm.ww, tm.wh)
}

// evolve slowly transitions from one texture to the next. This depends
// on seqentially ordering the similar textures in the texture atlas.
func (tm *tmtag) evolve(rate float64) {
	for x := range tm.evo {
		for y := range tm.evo[x] {
			eveo := tm.evo[x][y]
			even := tm.evo[x][y] + float64(rate)
			switch {
			case even > 2.99:
				even = 2.99
			case even < 0:
				even = 0
			}
			if eveo != even {
				tm.evo[x][y] = even
				tm.surface.Pts()[x][y].Tindex = int(even)
				tm.surface.Pts()[x][y].Blend = float32(even) - float32(int(even))
			}
		}
	}
	tm.surface.Update(tm.ground.Role().Mesh(), 0, 0)
}
