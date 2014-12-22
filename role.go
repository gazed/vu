// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

import (
	"image"
	"log"

	"github.com/gazed/vu/math/lin"
	"github.com/gazed/vu/render"
)

// Role manages rendered 3D objects.
// Role is the link between assets and the data needed by the rendering
// system. Assets, such as shaders, mesh, textures, etc. are specified
// as unique strings in the Role methods. Assets are loaded if necessary,
// and converted to internal data formats that are sent to the rendering
// system. Specifically, Role is the mediator between lazy loaded vu/load
// assets used to provide the shader data needed by a vu/render/Model.
type Role interface {
	Shader() (shader string)               // The shader for this role.
	SetMaterial(name string) Role          // Optional surface colour.
	SetKd(r, g, b float64)                 // Kd is the diffuse colour.
	SetAlpha(a float64)                    // Set or get the shaders
	Alpha() (a float64)                    // ...specific alpha value.
	Uniform(id string) (val []float32)     // Shader uniforms are named
	SetUniform(id string, val interface{}) // ...floats or float slices.

	// Mesh combined with a Texture form the basis of most visible objects.
	Mesh() render.Mesh                // Mesh groups per-vertex data where
	NewMesh(mesh string) render.Mesh  //  ...meshes can be initially empty,
	SetMesh(mesh string) Role         //  ...or created from assets.
	Tex(index int) (tex string)       // Texture is an image drawn on a mesh
	AddTex(texture string) Role       // ...There can be multiple in a shader,
	UseTex(texture string, index int) // ...and they can be replaced,
	RemTex(index int)                 // ...deleted,
	SetImg(i image.Image, index int)  // ...set from generated data,
	SetTexMode(index int, mode int)   // ...or rendered differently.

	// Fonts are used to display small text phrases using a default mesh.
	// SetPhrase generates a billboard mesh based on existing font data and
	// returns the mesh width in pixels. Fonts assume that a "uv" type shader
	// and one texture have been specified.
	SetFont(font string) Role            // Font is the character mapping.
	SetPhrase(phrase string) (int, Role) // Display the phrase.

	// Particle effects are either Application driven or shader based.
	// Set the effect to match the shader for this role. Particle shaders
	// treat mesh vertex data as points.
	SetEffect(e Effect) // A CPU or GPU based particle system.

	// Animations are loaded directly from disk. In addition to setting
	// the mesh and texture data from the animation file, extra animation
	// data like joints and frames are sent to the animation specific
	// shader. Animated models often have multiple named animations.
	// The done function is called each time the animation loop completes.
	Animation() render.Animation              // Nil if no animation.
	SetAnimation(id string) Role              // IQM/IQE animated model.
	PlayMovement(index int, done func()) bool // Return true if available.
	Movements() []string                      // Animation names.

	// Light: each role can have one light with a location and a colour.
	// Expected for shaders, like gouraud, that handle lighting.
	LightLocation() (x, y, z float64) // Get or
	SetLightLocation(x, y, z float64) // ...set light location.
	SetLightColour(r, g, b float64)   // Set light colour.

	// GPU pipeline controls for this role.
	SetDrawMode(mode int) // One of TRIANGLES (default), POINTS, LINES
	Set2D()               // Turn off Depth testing.
	SetCullOff()          // Turn off back-face culling for this role.
}

// FUTURE: there are possible render optimizations to be had by sorting models
//         based on shader/vao/textures, etc. The idea is to reduce expensive
//         GPU state switches. This sorting may be better implemented in the
//         vu/render package.
// FUTURE: Have the API make it clear that the following are generally just
//         different, and mutually exclusive, ways of specifying a
//         shader + mesh + texture combo:
//            • Mesh/Texture/Material combo.
//            • Font (mesh, font texture, font data, phrase)
//            • Effect (mesh, particle effect)
//            • Animation (mesh, texture[s], animation data)
//         Can doing this as components make the API cleaner and easier to use?

// Role interface
// ===========================================================================
// role - Role implementation

// role ensures that the graphic assets specified in the Role methods are
// loaded initialized and passed on to the rendering layer. This is the
// link between the asset loader, vu/load, and rendering, vu/render.
type role struct {
	model   render.Model // render shader and shader data.
	assets  *assets      // Asset manager.
	fnt     *font        // Font texture uv mappings.
	effect  *effect      // Optional particle effect.
	lloc    []float32    // light location
	lcolour []float32    // light colour
	tm      *lin.M4      // Scratch transform matrix.
}

// newRole allocates the necessary data structures.
func newRole(shaderName string, a *assets) *role {
	r := &role{}
	r.assets = a
	shader := r.assets.getShader(shaderName)
	r.model = r.assets.newModel(shader)
	r.tm = &lin.M4{}

	// Set default scene light data. This is updated later.
	r.lloc = []float32{0, 0, 0, 1}
	r.lcolour = []float32{0, 0, 0}
	r.SetUniform("l", r.lloc)
	r.SetUniform("ld", r.lcolour)
	return r
}

// dispose ensures that graphics objects are removed from the GPU and cache
// once they are no longer needed. Non-GPU objects can remain in the cache,
// but GPU objects need to be removed and, possibly,re-added to be properly
// initialized.
func (r *role) dispose() {
	msh := r.model.Mesh()
	shd := r.model.Shader()
	textures := r.model.Textures()
	r.assets.remModel(r.model)
	r.model.Dispose()
	if msh != nil && !msh.Bound() {
		r.assets.remMesh(msh)
	}
	if shd != nil && !shd.Bound() {
		r.assets.remShader(shd)
	}
	for _, t := range textures {
		if t != nil && !t.Bound() {
			r.assets.remTexture(t)
		}
	}
	r.model = nil
}

// Role interface implementation.
func (r *role) SetDrawMode(mode int)                { r.model.SetDrawMode(mode) }
func (r *role) Set2D()                              { r.model.Set2D() }
func (r *role) SetCullOff()                         { r.model.SetCullOff() }
func (r *role) SetImg(image image.Image, index int) { r.model.SetImage(image, index) }
func (r *role) SetTexMode(index int, mode int)      { r.model.TexMode(index, mode) }
func (r *role) RemTex(index int)                    { r.model.RemTexture(index) }
func (r *role) SetAlpha(a float64)                  { r.model.SetAlpha(a) }
func (r *role) Alpha() (a float64)                  { return r.model.Alpha() }
func (r *role) Mesh() render.Mesh                   { return r.model.Mesh() }
func (r *role) NewMesh(mesh string) render.Mesh {
	msh := r.assets.newMesh(mesh)
	r.model.SetMesh(msh)
	return msh
}
func (r *role) SetMesh(mesh string) Role {
	msh := r.assets.getMesh(mesh)
	r.model.SetMesh(msh)
	return r
}
func (r *role) SetAnimation(id string) Role {
	r.model = r.assets.getModel(id, r.model)
	return r
}
func (r *role) Shader() (shader string) {
	if shd := r.model.Shader(); shd != nil {
		return shd.Name()
	}
	return ""
}
func (r *role) Tex(index int) (texture string) {
	if t := r.model.Texture(index); t != nil {
		return t.Name()
	}
	return ""
}
func (r *role) AddTex(texture string) Role {
	if t := r.assets.getTexture(texture); t != nil {
		r.model.AddTexture(t)
	}
	return r
}
func (r *role) UseTex(texture string, index int) {
	if t := r.assets.getTexture(texture); t != nil {
		r.model.UseTexture(t, index)
	}
}
func (r *role) SetFont(font string) Role {
	r.fnt = r.assets.getFont(font)
	return r
}
func (r *role) SetPhrase(phrase string) (int, Role) {
	if r.model != nil && r.fnt != nil {
		if r.model.Mesh() == nil {
			r.model.SetMesh(r.model.Gc().NewMesh("phrase"))
		}
		width := r.fnt.Panel(r.model.Mesh(), phrase)
		return width, r
	}
	return 0, r
}
func (r *role) SetEffect(e Effect) {
	if eff, ok := e.(*effect); ok {
		r.effect = eff
	}
}
func (r *role) SetMaterial(name string) Role {
	mat := r.assets.getMaterial(name)
	r.model.SetUniform("kd", []float32{mat.kd.R, mat.kd.G, mat.kd.B})
	r.model.SetUniform("ks", []float32{mat.ks.R, mat.ks.G, mat.ks.B})
	r.model.SetUniform("ka", []float32{mat.ka.R, mat.ka.G, mat.ka.B})
	r.model.SetAlpha(float64(mat.tr))
	return r
}

// SetKd casts down to the float32's expected by the GPU.
func (rl *role) SetKd(r, g, b float64) {
	rl.model.SetUniform("kd", []float32{float32(r), float32(g), float32(b)})
}

// Lighting.
func (r *role) LightLocation() (x, y, z float64) {
	return float64(r.lloc[0]), float64(r.lloc[1]), float64(r.lloc[2])
}
func (r *role) SetLightLocation(x, y, z float64) {
	r.lloc[0], r.lloc[1], r.lloc[2] = float32(x), float32(y), float32(z)
}
func (rl *role) SetLightColour(r, g, b float64) {
	rl.lcolour[0], rl.lcolour[1], rl.lcolour[2] = float32(r), float32(g), float32(b)
}

// SetUniform looks for single floats or slices of floats.
func (r *role) SetUniform(name string, value interface{}) {
	if r.model != nil {
		switch v := value.(type) {
		case float32:
			r.model.SetUniform(name, []float32{v})
		case float64:
			r.model.SetUniform(name, []float32{float32(v)})
		case int:
			r.model.SetUniform(name, []float32{float32(v)})
		case []float32:
			r.model.SetUniform(name, v)
		case []float64:
			f32s := []float32{}
			for _, f64 := range v {
				f32s = append(f32s, float32(f64))
			}
			r.model.SetUniform(name, f32s)
		default:
			log.Print("Part.SetUniform: unknown type", value)
		}
	}
}
func (r *role) Uniform(name string) (values []float32) { return r.model.Uniform(name) }

// update the models transform and any ongoing particle effects or animations.
func (r *role) update(m, v, p *lin.M4, dt float64) {
	r.model.SetMvTransform(r.tm.Mult(m, v))     // model-view
	r.model.SetMvpTransform(r.tm.Mult(r.tm, p)) // model-view-projection
	r.model.Animate(dt)                         // nil animations ignored.
	if r.effect != nil {
		r.effect.Update(r.Mesh(), dt)
	}
}

// Animation
func (r *role) Animation() render.Animation {
	if r.model != nil {
		return r.model.Animation()
	}
	return nil
}

// PlayMovement plays the indicated movement.
func (r *role) PlayMovement(index int, done func()) bool {
	if r.model != nil {
		return r.model.PlayMovement(index, done)
	}
	return false
}
func (r *role) Movements() []string {
	if r.model != nil {
		return r.model.Movements()
	}
	return []string{}
}
