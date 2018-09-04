// Copyright © 2014-2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

// light.go holder for lighting related code. Lights currently only
//          have a color.
// FUTURE: handle multiple lights for one scene. Need a shader that
//         incorporates multiple lights into the final color.
//
// http://martindevans.me/game-development/2015/02/27/Drawing-Stuff-On-Other-Stuff-With-Deferred-Screenspace-Decals/
// https://mtnphil.wordpress.com/2014/05/24/decals-deferred-rendering/
//
// decals implies deferred rendering...
// https://gamedevelopment.tutsplus.com/articles/forward-rendering-vs-deferred-rendering--gamedev-12342
// https://www.gamedev.net/articles/programming/graphics/deferred-rendering-demystified-r2746/
// http://c0de517e.blogspot.ca/2011/01/mythbuster-deferred-rendering.html
// https://learnopengl.com/Advanced-Lighting/Deferred-Shading
// http://ogldev.atspace.co.uk/www/tutorial35/tutorial35.html
//
// however some alternatives are...
// https://turanszkij.wordpress.com/2017/10/12/forward-decal-rendering/
// https://owlcatgames.com/news/183.html
//
// Consider the possibility of light indexed deferred rendering...
// https://github.com/dtrebilco/lightindexed-deferredrender
//
// ... or screen space ambient occlusion.
// https://learnopengl.com/Advanced-Lighting/SSAO
//
// which leads to physically based rendering.
// https://learnopengl.com/PBR/Theory
// https://www.allegorithmic.com/pbr-guide

import (
	"log"
	"math"

	"github.com/gazed/vu/math/lin"
	"github.com/gazed/vu/render"
)

// Types of lights
const (
	DirectionalLight = iota
	PointLight
	SpotLight
)

// MakeLight adds a light to a scene, creating and returning a new Part
// on the given scene entity. Lighting aware shaders combine the lights
// a models material values. All objects in the scene with lighting
// capable shades will be influenced by the light.
//
// Internally light is has a position and orientation in world space.
// The default position is 0,0,0 and the default direction is 0,0,-1.
// A lights default color is white 1,1,1.
//
// Depends on Eng.AddScene.
func (e *Ent) MakeLight(kindOfLight int) *Ent {
	if s := e.app.scenes.get(e.eid); s != nil {
		light := e.AddPart()
		light.app.lights.create(e, light.eid, kindOfLight)
		return light
	}
	log.Printf("MakeLight needs AddScene %d", e.eid)
	return e
}

// SetLightColor assigns a color to the light. The lights color is combined
// with objects material color to produce a final color. The R,G,B light values
// are between 0 (no color), and 1 (full color).
//
// Depends on Ent.MakeLight.
func (e *Ent) SetLightColor(r, g, b float64) *Ent {
	if l := e.app.lights.get(e.eid); l != nil {
		l.r, l.g, l.b = r, g, b
		return e
	}
	log.Printf("SetLightColor needs MakeLight %d", e.eid)
	return e
}

// SetLightIntensity scales the amount of light affecting each component of
// the lighting model. The values are between 0 (dark) and 1 (full intensity).
//
// Depends on Ent.MakeLight.
func (e *Ent) SetLightIntensity(ambient, diffuse, specular float64) *Ent {
	if l := e.app.lights.get(e.eid); l != nil {
		l.ka, l.kd, l.ks = ambient, diffuse, specular
		return e
	}
	log.Printf("SetLightIntensity needs MakeLight %d", e.eid)
	return e
}

// SetLightAttenuation sets the lights constant, linear and quadratic
// attenuation factors: kc,kl,kq. Attenuation reduces the effect of light
// as distance increases. This affects PointLights and SpotLights.
//
// Depends on Ent.MakeLight(SpotLight or PointLight).
func (e *Ent) SetLightAttenuation(kc, kl, kq float64) *Ent {
	if l := e.app.lights.get(e.eid); l != nil &&
		(l.kind == SpotLight || l.kind == PointLight) {
		l.kc, l.kl, l.kq = kc, kl, kq
		return e
	}
	log.Printf("SetLightAttenuation needs MakeLight(SpotLight or PointLight) %d", e.eid)
	return e
}

// SetConeAngles sets the area affected by a spotlight in degrees.
// The cone angles affects the size of the spot light, where the outerAngle
// is expected to be slightly larger than the innerAngle in order to create a
// blur around the edges of the spotlight.
//
// Depends on Ent.MakeLight(SpotLight).
func (e *Ent) SetConeAngles(innerAngle, outerAngle float64) *Ent {
	if l := e.app.lights.get(e.eid); l != nil && l.kind == SpotLight {
		l.innerAngle = math.Cos(lin.Rad(innerAngle))
		l.outerAngle = math.Cos(lin.Rad(outerAngle))
		return e
	}
	log.Printf("SetCutoffs needs MakeLight(SpotLight) %d", e.eid)
	return e
}

// light entity methods
// =============================================================================
// light data.

// light holds light color for now. Anticipate other future light parameters.
// The lights world position are scatch values updated each render loop.
type light struct {
	kind       int     // DirectionalLight, PointLight, or SpotLight.
	r, g, b    float64 // Light color: values are 0 to 1.
	ka, kd, ks float64 // Light intensity: values are 0 to 1.

	// Light attenuation values. Used for point and spot lights.
	// attentuation = 1.0/(kc + kl*d + kq*d*d)
	kc, kl, kq float64 // constant, linear, quadratic attenuation values.

	// Cutoff angles limit the area of affect for spotlights.
	// 15 degree cutoff is set as: math.Cos(lin.Rad(15.0))
	innerAngle float64 // cosine of inner angle,
	outerAngle float64 // cosine of larger outer angle.
}

// newLight creates a white light.
func newLight() (l *light) {
	l = &light{}
	l.r, l.g, l.b = 1.0, 1.0, 1.0       // defaults
	l.ka, l.kd, l.ks = 0.1, 1.0, 1.0    //    "
	l.kc, l.kl, l.kq = 1.0, 0.09, 0.032 //    "
	return l
}

// draw turns a light into draw call data needed by the render shaders.
// This relies on a convention for shader uniform light structures.
// FUTURE: handle arrays of lights.
func (l *light) draw(d *render.Draw, p *pov) {
	if ref, ok := d.Uniforms["lightPosition"]; ok {
		wx, wy, wz := p.world()
		d.SetUniformData(ref, float32(wx), float32(wy), float32(wz))
	}
	if ref, ok := d.Uniforms["lightColor"]; ok {
		d.SetUniformData(ref, float32(l.r), float32(l.g), float32(l.b))
	}
	if l.kind == DirectionalLight {
		if ref, ok := d.Uniforms["dirLight.direction"]; ok {
			dx, dy, dz := lin.MultSQ(0, 0, -1, p.tn.Rot)
			d.SetUniformData(ref, float32(dx), float32(dy), float32(dz))
		}
		if ref, ok := d.Uniforms["dirLight.color"]; ok {
			d.SetUniformData(ref, float32(l.r), float32(l.g), float32(l.b))
		}
		if ref, ok := d.Uniforms["dirLight.intensity"]; ok {
			d.SetUniformData(ref, float32(l.ka), float32(l.kd), float32(l.ks))
		}
	}
	if l.kind == PointLight {
		if ref, ok := d.Uniforms["pointLight.position"]; ok {
			// light position
			wx, wy, wz := p.world()
			d.SetUniformData(ref, float32(wx), float32(wy), float32(wz))
		}
		if ref, ok := d.Uniforms["pointLight.color"]; ok {
			d.SetUniformData(ref, float32(l.r), float32(l.g), float32(l.b))
		}
		if ref, ok := d.Uniforms["pointLight.intensity"]; ok {
			d.SetUniformData(ref, float32(l.ka), float32(l.kd), float32(l.ks))
		}
		if ref, ok := d.Uniforms["pointLight.attenuation"]; ok {
			d.SetUniformData(ref, float32(l.kc), float32(l.kl), float32(l.kq))
		}
	}
	if l.kind == SpotLight {
		if ref, ok := d.Uniforms["spotLight.position"]; ok {
			wx, wy, wz := p.world()
			d.SetUniformData(ref, float32(wx), float32(wy), float32(wz))
		}
		if ref, ok := d.Uniforms["spotLight.direction"]; ok {
			dx, dy, dz := lin.MultSQ(0, 0, -1, p.tn.Rot)
			d.SetUniformData(ref, float32(dx), float32(dy), float32(dz))
		}
		if ref, ok := d.Uniforms["spotLight.color"]; ok {
			d.SetUniformData(ref, float32(l.r), float32(l.g), float32(l.b))
		}
		if ref, ok := d.Uniforms["spotLight.intensity"]; ok {
			d.SetUniformData(ref, float32(l.ka), float32(l.kd), float32(l.ks))
		}
		if ref, ok := d.Uniforms["spotLight.attenuation"]; ok {
			d.SetUniformData(ref, float32(l.kc), float32(l.kl), float32(l.kq))
		}
		if ref, ok := d.Uniforms["spotLight.cone"]; ok {
			d.SetUniformData(ref, float32(l.innerAngle), float32(l.outerAngle))
		}
	}
}

// light data.
// =============================================================================
// light component manager.

// lights manages all the active Light instances.
// There's not many lights so not much to optimize.
type lights struct {
	data  map[eid]*light // Light instance data.
	group map[eid][]eid  // Group lights by scene eid.
}

// newLights creates a light component manager. Expected to be called
// once on startup.
func newLights() *lights {
	return &lights{data: map[eid]*light{}, group: map[eid][]eid{}}
}

// get the light associated with the given entity, returning nil
// if there is no such light.
func (ls *lights) get(id eid) *light {
	if light, ok := ls.data[id]; ok {
		return light
	}
	return nil
}

// create and associate a light with the given entity. If there
// already is a light for the entity, dont create one and return
// the existing one instead.
func (ls *lights) create(scene *Ent, lightId eid, kindOfLight int) *light {
	l := newLight()
	l.kind = kindOfLight
	ls.data[lightId] = l // All lights.

	// Group lights by scene id.
	ls.group[scene.eid] = append(ls.group[scene.eid], lightId)
	return l
}

// draw populates the render.Draw data with any lighting data for
// the given scene.
func (ls *lights) draw(sceneId eid, d *render.Draw, app *application) {
	if lightIds, ok := ls.group[sceneId]; ok { // if scene has lights.
		for _, lightId := range lightIds {
			light := ls.data[lightId]
			pov := app.povs.get(lightId)
			light.draw(d, pov)
		}
	}
}

// position returns the world location of the directional light in the scene.
// is a temporary solution to explore casting a shadow.
// FUTURE: make this more generic and cast shadows for each light.
func (ls *lights) position(sceneId eid, app *application) (lx, ly, lz float64) {
	if lightIds, ok := ls.group[sceneId]; ok {
		for _, lightId := range lightIds {
			if light := ls.data[lightId]; light.kind == DirectionalLight {
				pov := app.povs.get(lightId)
				return pov.world()
			}
		}
	}
	return 0, 0, 0
}

// dispose the light associated for the given entity. Do nothing
// if no such light exists.
func (ls *lights) dispose(id eid) { delete(ls.data, id) }
