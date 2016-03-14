// Copyright © 2014-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package main

import (
	"log"

	"github.com/gazed/vu"
	"github.com/gazed/vu/land"
)

// tm demonstrates creating, texturing, and rendering a dynamic terrain map
// from a generated height map. The intent is to mimic a surface/land map.
//
// The water is simulated by two planes with a higher transparent blue plane
// covering a lower opaque blue plane.
func tm() {
	tm := &tmtag{}
	if err := vu.New(tm, "Terrain Map", 400, 100, 800, 600); err != nil {
		log.Printf("tm: error starting engine %s", err)
	}
	defer catchErrors()
}

// Encapsulate example specific data with a unique "tag".
type tmtag struct {
	cam     vu.Camera
	ww, wh  int         // window size.
	gm      vu.Model    // visible surface model
	ground  vu.Pov      // visible surface.
	coast   vu.Pov      // shallow water plane.
	ocean   vu.Pov      // deep water plane.
	world   land.Land   // height map generation.
	surface vu.Surface  // data structure used to create land.
	evo     [][]float64 // used for evolution experiments.
}

// Create is the engine callback for initial asset creation.
func (tm *tmtag) Create(eng vu.Eng, s *vu.State) {
	tm.ww, tm.wh = s.W, s.H
	tm.cam = eng.Root().NewCam()
	tm.cam.SetOrthographic(0, float64(tm.ww), 0, float64(tm.wh), 0, 50)
	sun := eng.Root().NewPov().SetLocation(0, 5, 0)
	sun.NewLight().SetColor(0.4, 0.7, 0.9)

	// create the world surface.
	seed := int64(123)
	patchSize := 128
	tm.world = land.New(patchSize, seed)
	worldTile := tm.world.NewTile(1, 0, 0)
	textureRatio := 256.0 / 1024.0
	tm.surface = vu.NewSurface(patchSize, patchSize, 16, float32(textureRatio), 10)

	// create a separate surface for generating initial land textures.
	emap := land.New(patchSize, seed-1)
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
	tm.ground = eng.Root().NewPov().SetLocation(0, -300, -10).SetScale(scale, scale, 1)
	tm.gm = tm.ground.NewModel("land").AddTex("land")
	tm.gm.LoadMat("tint").SetUniform("ratio", textureRatio)
	tm.gm.NewMesh("land")
	tm.surface.Update(tm.gm, 0, 0)

	// Add water planes.
	tm.ocean = eng.Root().NewPov()
	tm.ocean.SetLocation(256, 0, -10.5)
	tm.ocean.SetScale(float64(tm.ww), float64(tm.wh), 1)
	tm.ocean.NewModel("alpha").LoadMesh("plane").LoadMat("blue")
	tm.coast = eng.Root().NewPov().SetLocation(256, 0, -10)
	tm.coast.SetScale(float64(tm.ww), float64(tm.wh), 1)
	tm.coast.NewModel("alpha").LoadMesh("plane").LoadMat("transparent_blue")
	return
}

// Update is the regular engine callback.
func (tm *tmtag) Update(eng vu.Eng, in *vu.Input, s *vu.State) {
	if in.Resized {
		tm.ww, tm.wh = s.W, s.H
		tm.cam.SetOrthographic(0, float64(s.W), 0, float64(s.H), 0, 50)
	}

	// process user presses.
	for press, _ := range in.Down {
		switch press {

		// Change the water level.
		case vu.K_LBkt:
			dir := tm.cam.Lookat()
			tm.ocean.Move(0, 0, 1*in.Dt, dir)
			tm.coast.Move(0, 0, 1*in.Dt, dir)
		case vu.K_RBkt:
			dir := tm.cam.Lookat()
			tm.ocean.Move(0, 0, -1*in.Dt, dir)
			tm.coast.Move(0, 0, -1*in.Dt, dir)

		// Demonstrate texture evolution using a texture atlas.
		case vu.K_Equal:
			tm.evolve(0.01)
		case vu.K_Minus:
			tm.evolve(-0.01)
		}
	}
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
	tm.surface.Update(tm.gm, 0, 0)
}
