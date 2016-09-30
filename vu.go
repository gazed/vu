// Copyright © 2015-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// Package vu - virtual universe, provides 3D application support. Vu wraps
// subsystems like rendering, physics, asset loading, audio, etc. to provide
// higher level functionality that includes:
//    • Transform graphs and composite objects.
//    • Timestepped update/render loop.
//    • Access to user input events.
//    • Cameras and transform manipulation.
//    • Delivering loaded assets to render and audio devices.
// Refer to the vu/eg package for examples of engine functionality.
//
// Vu dependencies are:
//    • OpenGL for graphics card access.        See package vu/render.
//    • OpenAL for sound card access.           See package vu/audio.
//    • Cocoa  for OSX windowing and input.     See package vu/device.
//    • WinAPI for Windows windowing and input. See package vu/device.
package vu

// vu.go contains the main thread, machine, which is controlled by the engine
// state updater. The vu machine communicates with the hardware devices while
// the engine communicates with the application:
//    User -> devices <-> machine <-> engine <-> application -> User
//
// Concurrency design is based on "Share memory by communicating"
//     http://golang.org/doc/codewalk/sharemem
// in which ownership of structs is transferred when passing struct pointers
// between goroutines

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/gazed/vu/audio"
	"github.com/gazed/vu/device"
	"github.com/gazed/vu/render"
)

// New creates the Engine and initializes the underlying resources needed
// by the engine. It then starts application callbacks through the engine
// App interface. New is expected to be called one time on application startup.
//    app  : application callback handler.
//    name : window title.
//    wx,wy: bottom left window position in screen pixels.
//    ww,wh: window width and height in screen pixels.
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
	m.frame1 = frame{} // Previous render frame.
	m.frame0 = frame{} // Most recent render frame.

	// Start the application facing loop for state updates.
	// Run the device facing loop for rendering and user input polling.
	m.reqs = make(chan msg)
	m.stop = make(chan bool)
	m.draw = make(chan frame)
	go runEngine(app, wx, wy, ww, wh, m.reqs, m.draw, m.stop)
	m.run()            // underlying device polling and rendering.
	defer m.shutdown() // ensure shutdown happens no matter what.
	return nil         // report successful termination.
}

// Engine constants needed as input to methods as noted.
const (

	// Global graphic state constants. See Eng.State
	Blend       = render.Blend       // Alpha blending. Enabled by default.
	CullFace    = render.CullFace    // Backface culling. Enabled by default.
	DepthTest   = render.DepthTest   // Z-buffer awareness. Enabled by default.
	StaticDraw  = render.StaticDraw  // Data created once and rendered many times.
	DynamicDraw = render.DynamicDraw // Data is continually being updated.

	// Per-model rendering constants for Model DrawMode option.
	Triangles = render.Triangles // Triangles are the norm.
	Points    = render.Points    // Used for particle effects.
	Lines     = render.Lines     // Used for drawing squares and boxes.

	// KeyReleased indicator. Total time down, in update ticks,
	// is key down ticks minus KeyReleased. See App.Update.
	KeyReleased = device.KeyReleased

	// Application created and controlled objects associated with
	// the transform hierarchy. See Pov.Dispose.
	PovNode  = iota // Transform hierarchy node, 3D location:orientation.
	PovModel        // Rendered model attached to a Pov.
	PovBody         // Physics body attached to a Pov.
	PovCam          // Camera attached to a Pov.
	PovSound        // Sound attached to a Pov.
	PovLight        // Light attached to a Pov.
	PovLayer        // Render pass layer attached to a Pov.
)

// vu
// =============================================================================
// This is the machine and it is driven by the application facing engine class.

// machine deals with initialization and handling of all underlying hardware;
// generally through the OS, GPU, and audio API's. Machine is expected to be
// run from the main thread -- this is enforced by those aforementioned API's.
// Machine process requests from the application facing engine class.
type machine struct {
	gc     render.Renderer // Graphics card interface layer.
	dev    device.Device   // Os specific window and rendering context.
	ac     audio.Audio     // Audio card interface layer.
	input  *device.Pressed // Latest user keyboard and mouse input.
	frame1 frame           // Previous render frame.
	frame0 frame           // Most recent render frame.
	draw   chan frame      // return frame for updating.
	reqs   chan msg        // Requests from the application loop.
	stop   chan bool       // Used to shutdown the engine.

	// Counts keeps track of the number of faces and verticies
	// for each successfully bound mesh.
	counts map[uint32]*meshCount
}

// run is the main thread. Only the main thread can interact with the
// device layer and the rendering context. This loop depends on frequent
// and regular calls from the application update both for polling
// user input and rendering.
func (m *machine) run() {
	m.gc.Enable(Blend, true)    // expected application startup state.
	m.gc.Enable(CullFace, true) // expected application startup state.
	for m.dev != nil && m.dev.IsAlive() {
		req := <-m.reqs // req is nil for closed channel.

		// handle all communication, blocking until there is a request
		// to process. Requests wait until the current request is finished.
		switch t := req.(type) {
		case *shutdown:
			return // exit immediately. User shutdown engine.
		case *appData:
			m.refreshAppData(t) // poll to refresh device input.
		default:
			switch t := req.(type) {
			case *renderFrame:
				m.render(t)
			case *bindData:
				m.bind(t)
			case *setColor:
				m.gc.Color(t.r, t.g, t.b, t.a)
			case *enableAttr:
				m.gc.Enable(t.attr, t.enable)
			case *clampTex:
				m.gc.SetTextureMode(t.tid, true)
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
	close(m.stop) // The underlying device is gone, stop the engine.
}

// shutdown properly cleans up and closes the device layers.
func (m *machine) shutdown() {
	if m.ac != nil {
		m.ac.Dispose()
		m.ac = nil
	}
	if m.dev != nil {
		m.dev.Dispose()
		m.dev = nil
	}
}

// render passes the frame draw data to the supporting render layer.
// If there is a new frame, then the unused frame is sent back for updating.
func (m *machine) render(r *renderFrame) {

	// update the previous and current render frames with a new frame.
	if r.fr != nil && len(r.fr) > 0 {
		drawFrame := m.frame1 // return this frame to be updated.
		m.frame1 = m.frame0   // previous frame.
		m.frame0 = r.fr       // new frame.
		m.draw <- drawFrame   // return frame for updating.
	}

	// FUTURE: use interpolation between current and previous frames
	//         for render requests between frame updates.
	m.gc.Clear()
	for _, drawing := range m.frame0 {
		if drawing.Vao > 0 {
			m.setCounts(drawing)
			m.gc.Render(drawing)
		} else {
			log.Printf("machine.render: bad mesh vao %d", drawing.Vao)
		}
	}
	m.dev.SwapBuffers()
}

// refreshAppData gathers user input and returns it on request.
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

// setCounts ensures the draw frame has the most recent number of
// verticies and faces. The objects may have been updated after the
// draw object was created by the engine.
func (m *machine) setCounts(d *render.Draw) {
	if cnts, ok := m.counts[d.Vao]; ok {
		d.SetCounts(cnts.faces, cnts.verticies)
	} else {
		log.Printf("machine.setCounts: must have mesh counts %d", d.Vao)
	}
}

// bind sends data to the graphics or audio card and replies on the supplied
// channel when finished. Data needs to be bound once before it can be used
// for rendering or audio. Data needs rebinding if it is changed.
func (m *machine) bind(bd *bindData) {
	switch d := bd.data.(type) {
	case *mesh:
		bd.reply <- m.bindOne(d)
	case *shader:
		bd.reply <- m.bindOne(d)
	case *texture:
		bd.reply <- m.bindOne(d)
	case *sound:
		bd.reply <- m.bindOne(d)
	case *layer:
		bd.reply <- m.bindOne(d)
	case []asset:
		var err error
		for _, a := range d {
			if err = m.bindOne(a); err != nil {
				break
			}
		}
		bd.reply <- err
	default:
		bd.reply <- fmt.Errorf("No bindings for %T", d)
	}
}

// bindOne handles a single bind request. Expected to be called
// from bind().
func (m *machine) bindOne(a asset) error {
	switch d := a.(type) {
	case *mesh:
		err := m.gc.BindMesh(&d.vao, d.vdata, d.faces)
		if err == nil {
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
		}
		return err
	case *shader:
		var err error
		d.program, err = m.gc.BindShader(d.vsh, d.fsh, d.uniforms, d.layouts)
		return err
	case *texture:
		return m.gc.BindTexture(&d.tid, d.img)
	case *sound:
		return m.ac.BindSound(&d.sid, &d.did, d.data)
	case *layer:
		return m.gc.BindFrame(d.attr, &d.bid, &d.tex.tid, &d.db)
	}
	return fmt.Errorf("machine:bindOne. unhandled bind request")
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
	case *layer:
		m.gc.ReleaseFrame(d.bid, d.tex.tid, d.db)
		d.bid, d.tex.tid, d.db = 0, 0, 0
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

// msg are requests and data handled by the engine goroutine.
// Messages are pointers to one of the structures below.
// The struct type is the message and the struct fields carry the data.
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
	as.input = &Input{Down: map[int]int{}}
	as.state = &State{CullBacks: true, Blend: true}
	as.state.setColor(0, 0, 0, 1)
	return as
}

// shutdown is used to terminate a goroutine.
type shutdown struct{}

// bindData is a request to send data to the graphics or sound card.
type bindData struct {
	data  interface{} // 1 or more assets to be bound.
	reply chan error  // reply with any bind errors on this channel.
}

// renderFrame requests a render. New render data is supplied after
// an update. Interpolation values are given for render requests
// between new frames. No reply is expected. Render data is expected
// to be created by the engine update loop and processed by the
// vu machine.
type renderFrame struct {
	fr     frame   // May be empty.
	interp float64 // Fraction between 0 and 1.
	ut     uint64  // Counter for debugging.
}

// placeListener locates the sounds listener in world space.
type placeListener struct {
	x, y, z float64
}

// playSound plays the given sound at the given world location.
type playSound struct {
	sid     uint64
	x, y, z float64
}

// state change messages. Engine to machine. Fire and forget.
type enableAttr struct {
	attr   uint32
	enable bool
}
type setColor struct{ r, g, b, a float32 }
type setVolume struct{ gain float64 }
type setCursor struct{ cx, cy int }
type showCursor struct{ enable bool }
type toggleScreen struct{}
type clampTex struct{ tid uint32 }

// releaseData is used to request the removal a resources associated
// with one of the following:
//    bound and cached: *mesh, *shader, *texture, *sound, *noise,
//    cached only     : *material, *font, *animation
//    bound only      : *view
type releaseData struct {
	data interface{}
}

// =============================================================================

// catchErrors should be defered at the top of each goroutine so that
// errors can be logged in production loads as required by the application.
func catchErrors() {
	if r := recover(); r != nil {
		log.Printf("Panic %s: %s Shutting down.", r, debug.Stack())
		os.Exit(-1)
	}
}

// =============================================================================

// meshCount is used by bind to track the latest number
// of faces and verticies bound for a mesh.
type meshCount struct {
	faces     int // number of faces last bound.
	verticies int // number of verticies last bound.
}
