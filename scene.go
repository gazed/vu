// Copyright © 2017-2024 Galvanized Logic Inc.

package vu

// scene.go transforms application created models into render draw calls.

import (
	"log/slog"
	"math"
	"sort"

	"github.com/gazed/vu/load"
	"github.com/gazed/vu/math/lin"
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
	pass.Data[load.PROJ] = render.M4ToBytes(s.cam.pm, pass.Data[load.PROJ])
	pass.Data[load.VIEW] = render.M4ToBytes(s.cam.vm, pass.Data[load.VIEW])
	cx, cy, cz := s.cam.At()
	pass.Data[load.CAM] = render.V4SToBytes(cx, cy, cz, 0, pass.Data[load.CAM])

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
	pass.Data[load.LIGHTS] = render.LightsToBytes(pass.Lights, pass.Data[load.LIGHTS])
	pass.Data[load.NLIGHTS] = render.U8ToBytes(uint8(nlights), pass.Data[load.NLIGHTS])
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
	ss.setRenderCamera(ww, wh)
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
				return pass.Packets[i].Bucket > pass.Packets[j].Bucket
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

	if ready := app.models.getReady(p.eid); ready != nil {
		w := p.tw.Loc
		parts = append(parts, index)
		if sc.pid == render.Pass3D {
			// save distance to camera for transparency sorting.
			// closer objects drawn last.
			ready.tocam = sc.cam.distance(w.X, w.Y, w.Z)
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

	// turn all the pov's, models, and cameras into render packets.
	var packet *render.Packet
	for _, index := range parts {
		p := &(app.povs.povs[index])

		// generate render packets for models with loaded assets.
		if m := app.models.getReady(p.eid); m != nil && m.mesh != nil {
			if packets, packet = packets.GetPacket(); packet != nil {
				packet.Bucket = uint64(newBucket(sc.pid))

				// render model normally from scene camera.
				// This sets the expected shader uniforms into the draw call.
				m.fillPacket(packet, p, sc.cam)
			}

		}
	}
	return packets
}

// setPrev saves the previous locations and orientations.
// Called before application update and used in setRenderCamera
func (ss *scenes) setPrev() {
	for _, scene := range ss.all {
		c := scene.cam
		c.prev.Set(c.at)
	}
}

// setRenderCamera calculates the current render frame camera locations and
// orientations. Called before rendering to adjust for app camera changes.
func (ss *scenes) setRenderCamera(w, h uint32) {
	for _, scene := range ss.all {
		c := scene.cam

		// Cameras that haven't changed or moved already
		// have the correct transform matricies.
		if !c.focus && c.prev.Eq(c.at) {
			continue
		}
		if c.focus { // camera changed.
			scene.setProjection(w, h)
		}
		c.focus = false
		t0 := lin.NewT()
		t0.Loc.Set(c.at.Loc)
		t0.Rot.Set(c.at.Rot)

		// Set the view transform. Updates c.vm.
		c.vt(t0, c.vm) // updates view transform c.vm

		// Inverse only matters for perspective view transforms.
		c.it(t0, c.ivm) // updates inverse view transform c.ivm.
	}
}

// dispose removes the scene data associated with the given entity.
// Nothing happens if there is no scene data. Returns a list of
// eids that need other components disposed.
func (ss *scenes) dispose(eid eID, dead []eID) []eID {
	delete(ss.all, eid)
	return dead
}

// =============================================================================

// bucket is used to sort render packets.
type bucket uint64

// setBucket produces a number that is used to order draw calls.
// Higher values are rendered before lower values.
// setBucket creates an opaque object for the given render pass.
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
func newBucket(pass render.PassID) bucket {
	b := bucket(math.MaxUint8-pass) << 56 // render lower numbers before higher.
	return b | drawOpaque                 // opaque is default
}

// setDist to camera for sorting transparent objects.
func (b bucket) setDistance(toCam float64) bucket {
	return b | bucket(math.Float32bits(float32(toCam)))
}

// setType marks the object as the given type.
// Expects one of the type values defined below.
func (b bucket) setType(t bucket) bucket {
	return b&clearType | t // mark as the given type.
}

// setShaderID marks the object as the given type.
func (b bucket) setShaderID(sid uint16) bucket {
	return b&clearShaderID | bucket(sid)<<40
}

// Useful bits for setting or clearing the bucket.
const (
	clearDistance bucket = 0xFFFFFFFF00000000
	clearShaderID bucket = 0xFFFF0000FFFFFFFF
	clearType     bucket = 0xFF00FFFFFFFFFFFF

	// draw types.
	drawSky         bucket = 0x0008000000000000 // sky before other objects
	drawOpaque      bucket = 0x0004000000000000 // opaque objects before transparent
	drawTransparent bucket = 0x0001000000000000 // transparent objects last.
)
