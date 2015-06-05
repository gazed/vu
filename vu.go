// Copyright © 2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// Package vu, virtual universe, provides 3D application support. Vu wraps
// subsystems like rendering, physics, data loading, audio, etc. to provide
// higher level functionality that includes:
//    • Transform graphs and composite objects.
//    • Timestepped update/render loop.
//    • Access to user input events.
//    • Cameras and transform manipulation.
//    • Coupling loaded assets to render and audio systems.
// Refer to the vu/eg package for examples of engine functionality.
//
// Vu dependencies are:
//    • OpenGL for graphics card access.        See package vu/render.
//    • OpenAL for sound card access.           See package vu/audio.
//    • Cocoa  for OSX windowing and input.     See package vu/device.
//    • WinAPI for Windows windowing and input. See package vu/device.
package vu

// Design note: Concurrency based on "Share memory by communicating".
//              http://golang.org/doc/codewalk/sharemem/
// For structures and concurrency the key is in passing a pointer
// to a structure passes ownership of that structure instance.

import (
	"fmt"
	"log"

	"github.com/gazed/vu/audio"
	"github.com/gazed/vu/device"
	"github.com/gazed/vu/render"
)

// New creates the Engine and initializes the underlying resources needed
// by the engine. It then starts application callbacks through the App
// interface. This is expected to be called once on application startup.
func New(app App, name string, wx, wy, ww, wh int) (err error) {
	m := &machine{} // main thread and device facing handler.
	if app == nil {
		return fmt.Errorf("No application. Shutting down.")
	}
	m.counts = map[uint32]*meshCount{}

	// initialize the os specific shell, graphics context, and input tracker.
	name, wx, wy, ww, wh = m.vet(name, wx, wy, ww, wh)
	m.dev = device.New(name, wx, wy, ww, wh)

	// initialize the audio layer.
	m.ac = audio.New()
	if err = m.ac.Init(); err != nil {
		m.shutdown()
		return // failed to initialize audio layer
	}

	// initialize the graphics layer.
	m.gc = render.New()
	if err = m.gc.Init(); err != nil {
		m.shutdown()
		return // failed to initialize graphics layer.
	}
	m.gc.Viewport(ww, wh)
	m.dev.Open()
	m.input = m.dev.Update()

	// Start the application facing loop for state updates and
	// enter the device facing loop for rendering and user input polling.
	m.reqs = make(chan msg)
	m.stop = make(chan bool)
	go runEngine(app, wx, wy, ww, wh, m.reqs, m.stop)
	m.startup()        // underlying device polling and rendering.
	defer m.shutdown() // ensure shutdown happens no matter what.
	return nil         // report successful termination.
}

// Engine constants used as input to various methods.
const (

	// Global graphic state constants. See Eng.State
	BLEND = render.BLEND // Alpha blending. Enabled by default.
	CULL  = render.CULL  // Backface culling. Enabled by default.
	DEPTH = render.DEPTH // Z-buffer awareness. Enabled by default.

	// Per-part rendering constants for Model.SetDrawMode.
	TRIANGLES = render.TRIANGLES // Triangles are the norm.
	POINTS    = render.POINTS    // Used for particle effects.
	LINES     = render.LINES     // Used for drawing squares and boxes.

	// User input key released indicator. Total time down, in update
	// ticks, is key down ticks minus RELEASED. See App.Update.
	RELEASED = device.KEY_RELEASED

	// Render buckets. Lower values drawn first.
	OPAQUE      = render.OPAQUE      // draw first
	TRANSPARENT = render.TRANSPARENT // draw after opaque
	OVERLAY     = render.OVERLAY     // draw last.

	// Texture rendering directives for Model.SetTexMode()
	TEX_REPEAT = iota // Repeat texture when UV greater than 1.
	TEX_CLAMP         // Clamp to texture edge.

	// 3D Direction constants. Primarily used for panning or
	// rotating a camera view. See Camera.Spin.
	XAxis // Affect only the X axis.
	YAxis // Affect only the Y axis.
	ZAxis // Affect only the Z axis.

	// Objects created by the application.
	// Note that POV is both a transform hierarchy node
	// and a particular location:orientation in 3D space.
	POV   // Transform hierarchy node and 3D location:orientation.
	MODEL // Rendered model attached to a Pov.
	BODY  // Physics body attached to a Pov.
	VIEW  // Camera and view transform attached to a Pov.
	NOISE // Sound attached to a Pov.
	LIGHT // Light attached to a Pov.
)

// vu
// =============================================================================
// machine  Defn: "engine is a device that drives a machine"
// This is the machine and it is driven by the application facing engine class.

// machine deals with initialization and handling of all underlying hardware;
// generally through the OS, GPU, and audio API's. Machine is expected to be
// run from the main thread -- this is expected and enforced by those
// aforementioned API's.
type machine struct {
	gc     render.Renderer // Graphics card interface layer.
	dev    device.Device   // Os specific window and rendering context.
	ac     audio.Audio     // Audio card interface layer.
	reqs   chan msg        // Requests from the application loop.
	stop   chan bool       // Used to shutdown the application loop.
	frame0 []render.Draw   // Previous render frame.
	frame1 []render.Draw   // Most recent render frame.
	input  *device.Pressed // Latest user keyboard and mouse input.

	// Counts keeps track of the number of faces and verticies for
	// each successfully bound mesh.
	counts map[uint32]*meshCount
}

// startup is called on the main thread. Only the main thread can interact
// with the device layer and the rendering context. This loop depends on
// frequent and regular calls from the application update both for polling
// user input and rendering.
func (m *machine) startup() {
	m.gc.Enable(BLEND, true) // match application startup state.
	m.gc.Enable(CULL, true)  // match application startup state.
	for m.dev != nil && m.dev.IsAlive() {
		req := <-m.reqs // guard against a closed channel.

		// handle all communication, blocking until there is a request
		// to process. Requests wait until the current request is finished.
		switch t := req.(type) {
		case *shutdown:
			return // exit immediately. The app loop is dead.
		case *appData:
			m.refreshAppData(t) // sync with engine that is blocking on reply.
		default:
			switch t := req.(type) {
			case *renderFrame:
				m.render(t)
			case *bindData:
				m.bind(t)
			case *setColour:
				m.gc.Color(t.r, t.g, t.b, t.a)
			case *enableAttr:
				m.gc.Enable(t.attr, t.enable)
			case *toggleScreen:
				m.dev.ToggleFullScreen()
			case *setVolume:
				m.ac.SetGain(t.gain)
			case *setCursor:
				m.dev.SetCursorAt(t.cx, t.cy)
			case *showCursor:
				m.dev.ShowCursor(t.enable)
			case *placeListener:
				m.ac.PlaceListener(t.x, t.y, t.z)
			case *playSound:
				m.ac.PlaySound(t.sid, t.x, t.y, t.z)
			case *releaseData:
				m.release(t)
			case nil:
				return // exit immediately: channel closed.
			default:
				log.Printf("machine: unknown msg %T", t)
			}
		}
	}
	close(m.stop) // The underlying device is gone, stop the app loop.
}

// shutdown stops and closes the engine and the applications shell.
func (m *machine) shutdown() {
	if m.ac != nil {
		m.ac.Shutdown()
		m.ac = nil
	}
	if m.dev != nil {
		m.dev.Dispose()
		m.dev = nil
	}
}

// render passes the frame draw data to the supporting render layer.
func (m *machine) render(r *renderFrame) {
	m.gc.Clear()
	if r.frame != nil { // update the previous and current render frames.
		m.frame0 = m.frame1
		m.frame1 = r.frame
	}

	// FUTURE: use the renderFrame interpolation and the previous frame
	//         for rendering between frame updates.
	for _, drawing := range m.frame1 {
		if drawing.Vao() > 0 {
			m.setCounts(drawing)
			m.gc.Render(drawing)
		} else {
			log.Printf("machine.render: bad mesh vao %d", drawing.Vao())
		}
	}
	m.dev.SwapBuffers()
}

// refreshAppData gathers user input and returns it on request.
// Only poll input when requested so input is not dropped.
// The underlying device layer collects input since last call.
// Expected to be called once per update tick.
func (m *machine) refreshAppData(data *appData) {
	data.input.convertInput(m.input, 0, 0) // refresh user data.
	if data.input.Resized {
		data.state.setScreen(m.dev.Size())
		m.gc.Viewport(data.state.W, data.state.H)
	}
	data.state.FullScreen = m.dev.IsFullScreen()
	data.reply <- data       // return refreshed app data.
	m.input = m.dev.Update() // get latest user input for next refresh.
}

func (m *machine) setCounts(d render.Draw) {
	if cnts, ok := m.counts[d.Vao()]; ok {
		d.SetCounts(cnts.faces, cnts.verticies)
	} else {
		log.Printf("machine.setCounts: must have mesh counts %d", d.Vao())
	}
}

// bind sends data to the graphics or audio card and replies when finished.
// Data needs to be bound once before it can be used for rendering or audio.
// Data needs rebinding if it is changed.
func (m *machine) bind(bd *bindData) {
	switch d := bd.data.(type) {
	case *mesh:
		err := m.gc.BindMesh(&d.vao, d.vdata, d.faces)
		if err != nil {
			bd.reply <- fmt.Errorf("Failed mesh bind %s: %s", d.name, err)
		} else {
			cnts, ok := m.counts[d.vao]
			if !ok {
				cnts = &meshCount{}
				m.counts[d.vao] = cnts
			}
			if d.faces != nil {
				cnts.faces = d.faces.Len()
			}
			if d.vdata != nil && len(d.vdata) > 0 {
				cnts.verticies = d.vdata[0].Len()
			}
			bd.reply <- nil
		}
	case *shader:
		var err error
		d.program, err = m.gc.BindShader(d.vsh, d.fsh, d.uniforms, d.layouts)
		if err != nil {
			bd.reply <- fmt.Errorf("Failed shader bind %s: %s", d.name, err)
		} else {
			bd.reply <- nil
		}
	case *texture:
		err := m.gc.BindTexture(&d.tid, d.img, d.repeat)
		if err != nil {
			bd.reply <- fmt.Errorf("Failed texture bind %s: %s", d.name, err)
		} else {
			bd.reply <- nil
		}
	case *sound:
		err := m.ac.BindSound(&d.sid, &d.did, d.data)
		if err != nil {
			bd.reply <- fmt.Errorf("Failed sound bind %s: %s", d.name, err)
		} else {
			bd.reply <- nil
		}
	default:
		bd.reply <- fmt.Errorf("No bindings for %T", d)
	}
}

// meshCount is used by bind to track the latest number
// of faces and verticies bound for a mesh.
type meshCount struct {
	faces     int // number of faces last bound.
	verticies int // number of verticies last bound.
}

// release figures out what data to release based on the releaseData type.
func (m *machine) release(rd *releaseData) {
	switch d := rd.data.(type) {
	case *mesh:
		m.gc.ReleaseMesh(d.vao)
	case *shader:
		m.gc.ReleaseShader(d.program)
	case *texture:
		m.gc.ReleaseTexture(d.tid)
	case *sound:
		m.ac.ReleaseSound(d.sid)
	default:
		log.Printf("machine.release: No bindings for %T", rd)
	}
}

// constants to ensure reasonable input values.
const (
	maxWindowTitle    = 40  // Max number of characters for a window title.
	minWindowSize     = 100 // Miniumum pixels for a window width or height.
	minWindowPosition = 0   // Bottom left corner of the screen.
)

// vet ensures that the startup parameters will result in a visible
// window with a window title.
func (m *machine) vet(name string, x0, y0, width, height int) (n string, x, y, w, h int) {
	if len(name) > maxWindowTitle {
		name = ""
	}
	if width < minWindowSize {
		width = minWindowSize
	}
	if height < minWindowSize {
		height = minWindowSize
	}
	if x0 < minWindowPosition {
		x0 = minWindowPosition
	}
	if y0 < minWindowPosition {
		y0 = minWindowPosition
	}
	return name, x0, y0, width, height
}

// machine
// =============================================================================
// msg

// msg: each structure below is a msg between concurrent goroutines.
//      Messages are pointers to one of the structures below.
type msg interface{}

// appData contains both user input and engine state passed to the
// application. A single copy, owned by the engine, is created on startup.
type appData struct {
	input *Input        // Refreshed each update.
	state *State        // Refreshed each update.
	reply chan *appData // For syncing updates between machine and operator.
}

// newAppData expects to be called on startup for
// updating and communicating user input and global state.
func newAppData() *appData {
	as := &appData{reply: make(chan *appData)}
	as.input = &Input{Down: map[string]int{}}
	as.state = &State{CullBacks: true, Blend: true}
	as.state.setColour(0, 0, 0, 1)
	return as
}

// shutdown is used to terminate a goroutine.
type shutdown struct{}

// bindData is a request to send data to the graphics or sound card.
type bindData struct {
	data  interface{} // msh, shd, tex, snd
	reply chan error
}

// renderFrame requests a render. New render data is supplied after
// an update and interpolation values are given for render requests
// between new frames. No reply is expected. Render data is expected
// to be created by the engine update loop and read/rendered by the
// vu machine.
type renderFrame struct {
	interp float64       // Fraction between 0 and 1.
	frame  []render.Draw // May be empty.
	ut     uint64        // Counter for debugging.
}

// placeListener locates the sounds listener in world space.
type placeListener struct {
	x, y, z float64
}

// playSound plays the given sound at the given world location.
type playSound struct {
	sid     uint32
	x, y, z float64
}

// state change messages. Operator to machine. Fire and forget.
type enableAttr struct {
	attr   uint32
	enable bool
}
type setColour struct{ r, g, b, a float32 }
type setVolume struct{ gain float64 }
type setCursor struct{ cx, cy int }
type showCursor struct{ enable bool }
type toggleScreen struct{}

// releaseData removes the underlying resources associated with one
// of the following:
//    bound and cached: *mesh, *shader, *texture, *sound, *noise,
//    cached only     : *material, *font, *animation
type releaseData struct {
	data interface{} // *mesh, *shader, *texture, *sound, *noise
}
