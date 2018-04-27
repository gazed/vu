// Copyright © 2017-2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// FUTURE: Skyboxes rendering with a single quad.
//   https://gamedev.stackexchange.com/questions/60313/implementing-a-skybox-with-glsl-version-330
//   http://www.rioki.org/2013/03/07/glsl-skybox.html

import (
	"github.com/gazed/vu/render"
)

// AddSky adds one sky dome to a 3D scene entity.
// A sky is an optional pov with static model that is rendered
// prior to and behind all other scene objects. It uses the scene camera
// rotation while ignoring the scene camera location.
//
// The returned entity expects to be populated with a sky dome model
// and sky texture. Nil is returned if the entity is not a 3D scene or
// if there is already a sky dome attached to the scene.
func (e *Ent) AddSky() *Ent {
	sky := e.app.scenes.createSky(e)
	if sky == nil {
		return nil
	}
	return &Ent{app: e.app, eid: sky.eid}
}

// newSky creates a pov entity outside the scene graph hierarchy.
// The application is responsible for adding the model and texture.
func newSky(app *application) *sky {
	s := &sky{}
	s.eid = app.eids.create()
	app.povs.create(s.eid, 0)            // Keep outside scene graph hierarchy.
	s.cam = newCamera().SetAt(0, 0.2, 0) // 0.2 avoids clipping the model bottom.
	return s
}

// sky holds the data needed for a sky dome. Each sky dome has its
// own eid and its data is tracked by the scene component manager.
type sky struct {
	eid eid     // holder for sky dome pov, model, and texture.
	cam *Camera // like scene.cam but ignores location.
}

// draw renders the skydome using the scene camera rotation so the
// dome rotates in-step with the scene.
func (s *sky) draw(app *application, scene *scene, f frame) frame {
	p := app.povs.get(s.eid)
	if m := app.models.getReady(s.eid); m != nil && p != nil {
		var draw **render.Draw
		if f, draw = f.getDraw(); draw != nil {
			scene.draw(*draw) // apply scene wide attributes.

			// sky needs to be drawn with sky cam.
			// Sky cam is updated by the scene component manager.
			p.draw(*draw, s.cam.pm, s.cam.vm)
			m.draw(*draw, nil) // draw dome model.
			(*draw).Depth = false
			(*draw).Bucket = setSky((*draw).Bucket) // specify render order.
		}
	}
	return f
}
