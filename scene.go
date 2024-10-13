// Copyright © 2017-2024 Galvanized Logic Inc.

package vu

// scene.go transforms application created models into render draw calls.

import (
	"log/slog"
	"math"
	"sort"

	"github.com/gazed/vu/load"
	"github.com/gazed/vu/render"
)

// SceneType corresponds to a render pass.
//   - Scene3D for the normal world rendering.
//   - Scene2D for the UI overlay rendering.
type SceneType uint32

const (
	//	Scene3D for the normal world rendering.
	Scene3D SceneType = SceneType(render.Pass3D)
	//	Scene2D for the UI overlay rendering.
	Scene2D SceneType = SceneType(render.Pass2D)
)

// Cam returns the camera instace for a scene, returning nil
// if the entity is not a scene.
//
// Depends on Eng.AddScene.
func (e *Entity) Cam() *Camera {
	if s := e.app.scenes.get(e.eid); s != nil {
		return s.cam
	}
	slog.Error("Cam needs AddScene", "eid", e.eid)
	return nil
}

// =============================================================================
// scene data

// scene contains application created resources used to render screen images.
// A scene groups one camera with a group of application created entities.
// Scene is created by the Application calling Eng.AddScene().
type scene struct {
	pid render.PassID // scene render pass
	eid eID           // Scene and top level scene graph node.
	fbo uint32        // Render target. Default 0: display buffer.

	// Cam is this scenes camera data. Guaranteed to be non-nil.
	cam *Camera // Created automatically with a new scene.
}

// newScene creates a new transform hiearchy branch with its own camera.
func newScene(eid eID, passID render.PassID) *scene {
	s := &scene{eid: eid, pid: passID, cam: newCamera()}
	switch s.pid {
	case render.Pass2D:
		s.cam.SetClip(0.0, 10.0) // orthographic default clip
	case render.Pass3D:
		s.cam.SetClip(0.1, 1000.0) // projection default clip
	}
	return s
}

// setProjection updates scenes cameras projection matrix to match
// the latest application window size.
func (s *scene) setProjection(ww, wh uint32) {
	w, h := float64(ww), float64(wh)
	c := s.cam
	switch {
	case s.pid == render.Pass2D:
		c.setOrthographic(0, w, 0, h, c.near, c.far)
	default:
		c.setPerspective(c.fov, w/h, c.near, c.far)
	}
	c.focus = true
}

// setPassUniformData sets the pass uniform data.
func (s *scene) setPassUniformData(app *application, pass *render.Pass) {
	pass.Uniforms[load.PROJ] = render.M4ToBytes(s.cam.pm, pass.Uniforms[load.PROJ])
	pass.Uniforms[load.VIEW] = render.M4ToBytes(s.cam.vm, pass.Uniforms[load.VIEW])
	cx, cy, cz := s.cam.At()
	pass.Uniforms[load.CAM] = render.V4SToBytes(cx, cy, cz, 0, pass.Uniforms[load.CAM])

	// adds any scene lights to the render pass.
	// lights are children of the scene.
	nlights := 0
	index := app.povs.index[s.eid] // scene eID
	n := app.povs.nodes[index]
	for _, kid := range n.kids {
		ki, ok := app.povs.index[kid]
		if !ok {
			continue
		}

		// ignore culled children
		kn := app.povs.nodes[ki]
		if culled := kn.cull; culled {
			continue
		}

		// get the light
		kp := app.povs.povs[ki]
		l := app.lights.get(kp.eid)
		if l == nil {
			continue
		}

		// fill in the light render information into pass.Lights.
		if nlights >= len(pass.Lights) {
			slog.Warn("scene:setPassUniformData to many lights")
			break
		}
		light := &pass.Lights[nlights]
		px, py, pz := kp.at()
		light.X, light.Y, light.Z = float32(px), float32(py), float32(pz)
		light.R, light.G, light.B = l.r, l.g, l.b
		light.Intensity = l.intensity
		nlights += 1

	}
	pass.Uniforms[load.LIGHTS] = render.LightsToBytes(pass.Lights, pass.Uniforms[load.LIGHTS])
	pass.Uniforms[load.NLIGHTS] = render.U8ToBytes(uint8(nlights), pass.Uniforms[load.NLIGHTS])
}

// =============================================================================
// scene component manager.

// scenes manages all the Scene instances.
// There's not many scenes so not much to optimize.
type scenes struct {
	all      map[eID]*scene // Scene instance data.
	released []asset        // Scene assets being disposed.

	// Scratch variables: reused each update.
	parts []uint32 // Flattened pov hiearchy.
}

// newScenes creates the scene component manager and is expected to
// be called once on startup.
func newScenes() *scenes {
	ss := &scenes{}
	ss.all = map[eID]*scene{}
	ss.parts = []uint32{} // updated each frame
	return ss
}

// create makes a new scene and associates it with the given entity.
// Nothing is created if there already is a scene for the given entity.
func (ss *scenes) create(eid eID, sceneType SceneType) *scene {
	scene, ok := ss.all[eid]
	if !ok {
		ss.all[eid] = newScene(eid, render.PassID(sceneType))
	}
	return scene // don't allow creating over existing scene.
}

// resize the scene cameras to the new window dimensions.
func (ss *scenes) resize(ww, wh uint32) {
	for _, scene := range ss.all {
		scene.cam.focus = true
	}
	ss.setViewMatrixes(ww, wh)
}

// setViewMatrixes calculates the current render frame camera locations and
// orientations. Called before rendering to adjust for app camera changes.
func (ss *scenes) setViewMatrixes(w, h uint32) {
	for _, scene := range ss.all {
		scene.setProjection(w, h)
		scene.cam.updateView()
	}
}

// get returns the Scene associated with the given entity.
func (ss *scenes) get(id eID) *scene { return ss.all[id] }

// getFrame converts the scene transform hierarchy to a frame of render packets.
//
// The provided frame memory is recycled in that the render packets are lazy
// allocated and reused each update. The updated frame is returned.
func (ss *scenes) getFrame(app *application, frame []render.Pass) []render.Pass {
	if len(ss.all) <= 0 {
		return frame // the app hasn't created scenes yet.
	}
	if len(ss.all) > 2 {
		slog.Error("one or two scenes supported: 3D or 2D or 3D+2D")
		return frame
	}

	// turn the scene models into a frame of render.Packets.
	for _, sc := range ss.all {
		pass := &frame[sc.pid]           // a scene is either a 3D or 2D render pass.
		pass.Reset()                     // reset and reuse previous pass.
		sc.setPassUniformData(app, pass) // set scene uniform data in the pass.
		if n := app.povs.getNode(sc.eid); n != nil && !n.cull {
			index := app.povs.index[sc.eid]
			ss.parts = ss.listParts(app, sc, index, ss.parts[:0])
			pass.Packets = ss.renderParts(app, sc, ss.parts, pass.Packets)

			// sort the render pass packets.
			sort.SliceStable(pass.Packets, func(i, j int) bool {
				return pass.Packets[i].Bucket < pass.Packets[j].Bucket
			})
		}
		frame[sc.pid] = *pass // save the updated pass.
	}
	return frame
}

// listParts recursively turns the Pov hierarchy into a flat list using a depth
// first traversal. Pov's not affecting the rendered scene are excluded.
//
// TODO: cull more Pov's - often based on camera distance.
func (ss *scenes) listParts(app *application, sc *scene, index uint32, parts []uint32) []uint32 {
	p := app.povs.povs[index]
	n := app.povs.nodes[index]
	if culled := n.cull; culled {
		return parts
	}

	// get the model for this part
	if m := app.models.get(p.eid); m != nil {
		w := p.tw.Loc
		parts = append(parts, index)
		if sc.pid == render.Pass3D {
			// save distance to camera for transparency sorting.
			// closer objects drawn last.
			m.tocam = sc.cam.distance(w.X, w.Y, w.Z)
		}
	}

	// recurse scene graph processing children of viable elements.
	for _, kid := range n.kids {
		if ki, ok := app.povs.index[kid]; ok {
			parts = ss.listParts(app, sc, ki, parts)
		}
	}
	return parts
}

// renderParts prepares for rendering by converting a sequenced list
// of pov's into render packets.
func (ss *scenes) renderParts(app *application, sc *scene, parts []uint32, packets render.Packets) render.Packets {
	modelsNotReady := []error{}

	// turn all the pov's, models, and cameras into render packets.
	var packet *render.Packet
	for _, index := range parts {
		p := &(app.povs.povs[index])

		// generate render packets for models with loaded assets.
		if m := app.models.get(p.eid); m != nil {
			if packets, packet = packets.GetPacket(); packet != nil {
				packet.Bucket = newBucket(sc.pid)

				// render model normally from scene camera.
				// This sets the shader uniforms in the render packet.
				if err := m.fillPacket(packet, p, sc.cam); err != nil {
					modelsNotReady = append(modelsNotReady, err)

					// exclude models that are not ready to render.
					packets = packets.DiscardLastPacket()
				}
			}
		}
	}
	if len(modelsNotReady) > 0 {
		slog.Debug("models not ready", "err#", len(modelsNotReady), "err0", modelsNotReady[0])
	}
	return packets
}

// dispose removes the scene data associated with the given entity.
// Nothing happens if there is no scene data. Returns a list of
// eids that need other components disposed.
func (ss *scenes) dispose(eid eID, dead []eID) []eID {
	delete(ss.all, eid)
	return dead
}

// =============================================================================

// newBucket produces a number that is used to order draw calls
// from lowest to highest bucket value, ie: sort a<b.
//   - Pass.... DrawType ShaderID ........ Distance to Camera.................
//     00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000
//     F   F    F   F    F   F    F   F    F   F    F   F    F   F    F   F
//   - Pass bits is the render pass.
//   - DrawType sorts within groups of objects.
//     8 Skydome     : draw first
//     4 Opaque      : draw next, generally drawn in order created.
//     1 Transparent : draw last, sorted back to front using distance to camera.
//   - ShaderID to reduce pipeline switches.
//
// Pass bits have their values reversed so that low
// values result in higher bucket numbers, ie:
//   - 0 (render.Pass3D) = 255 bucket Pass value - render first
//   - 1 (render.Pass2D) = 254 bucket Pass value - render next
func newBucket(pass render.PassID) uint64 {
	b := uint64(math.MaxUint8-pass) << 56 // render lower numbers before higher.
	return b | drawOpaque                 // opaque is default
}

// setBucketDistance to camera for sorting transparent objects.
// closer objects should be drawn last.
func setBucketDistance(b uint64, toCam float64) uint64 {
	return b | uint64(math.Float32bits(float32(toCam)))
}

// setBucketType marks the object as the given type.
// Expects one of the type values defined below.
func setBucketType(b uint64, t uint64) uint64 {
	return b&clearType | uint64(t) // mark as the given type.
}

// setBucketShader marks the object as the given type.
func setBucketShader(b uint64, sid uint16) uint64 {
	return b&clearShaderID | uint64(sid)<<40
}

// setBucketLayer
func setBucketLayer(b uint64, layer uint8) uint64 {
	if layer < 16 {
		return b&clearLayer | uint64(layer)<<52
	}
	return b
}

// Useful bits for setting or clearing the bucket.
const (
	clearDistance uint64 = 0xFFFFFFFF00000000
	clearShaderID uint64 = 0xFFFF0000FFFFFFFF

	// draw types.
	clearType       uint64 = 0xFFF0FFFFFFFFFFFF
	drawOpaque      uint64 = 0x0001000000000000 // opaque objects before transparent
	drawTransparent uint64 = 0x0008000000000000 // transparent objects last.

	// layers.
	clearLayer uint64 = 0xFF0FFFFFFFFFFFFF
)
