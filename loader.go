// Copyright © 2015-2024 Galvanized Logic Inc.

package vu

// loader.go uses a goroutine to load asset data from disk.
// The asset data is then stored in an asset object and kept
// for reuse.

import (
	"fmt"
	"image"
	"image/color"
	"log/slog"
	"math"
	"path"
	"strings"
	"time"

	"github.com/gazed/vu/audio"
	"github.com/gazed/vu/load"
	"github.com/gazed/vu/render"
)

// assetLoader imports application assets data from disk.
// loader wraps the load package for the engine.
type assetLoader struct {

	// loaded tracks asset files.
	loaded map[string]bool // true: completed, false: still loading

	// requests tracks outstanding entities requests for assets.
	requests map[aid][]assetRequest

	// labelRequests match a font aid with an label Entity.
	// The label Entity needs a mesh generated that matches its string.
	labelRequests map[aid][]*Entity

	// loaded assets are remembered so future requests
	// can immediately return the same assets.
	assets map[aid]asset // loaded assets indexed by aid.

	// load asset files using a goroutine.
	loadAssetReq chan string           // request load asset file eg: "bloop.wav"
	loadedAssets chan []load.AssetData // file assets finished importing.
}

// newLoader is called once on startup by the engine.
func newLoader() *assetLoader {
	l := &assetLoader{}
	l.loaded = map[string]bool{}
	l.requests = map[aid][]assetRequest{}
	l.labelRequests = map[aid][]*Entity{}
	l.assets = map[aid]asset{}

	// allocate enough workers to avoid having to wait for a worker.
	numWorkers := 5
	l.loadAssetReq = make(chan string, 100)
	l.loadedAssets = make(chan []load.AssetData, 100)
	for wid := 1; wid <= numWorkers; wid++ {
		go importWorker(wid, l.loadAssetReq, l.loadedAssets)
	}
	return l
}

// getAsset returns the requested asset on the given callback.
// The asset is returned immediately if it is already loaded, or
// it is returned once it loads.
func (l *assetLoader) getAsset(aid aid, eid eID, callback assetReady) {
	if a, ok := l.assets[aid]; ok {
		callback(eid, a) // asset was already loaded.
		return
	}

	// register the asset request.
	if reqs, ok := l.requests[aid]; ok {
		l.requests[aid] = append(reqs, assetRequest{eid: eid, callback: callback})
	} else {
		l.requests[aid] = []assetRequest{{eid: eid, callback: callback}}
	}
}

// getLabelMesh remembers the request for a label mesh and handles
// it within loadAssets
func (l *assetLoader) getLabelMesh(aid aid, me *Entity) {
	if reqs, ok := l.labelRequests[aid]; ok {
		l.labelRequests[aid] = append(reqs, me)
	} else {
		l.labelRequests[aid] = []*Entity{me}
	}
}

// assetReady is the callback function needed to notify
// when requested asset is available.
type assetReady func(eid eID, a asset)

// requests tracks entities requests for assets.
type assetRequest struct {
	eid      eID        // entity requesting asset load.
	callback assetReady // callback on asset load complete.
}

// importAssetData gets the requested asset data for the given entity.
// The caller is responsible for adhering to the asset file name conventions.
func (l *assetLoader) importAssetData(assetFilenames ...string) {
	for _, filename := range assetFilenames {
		if _, ok := l.loaded[filename]; ok {
			continue // file is being loaded or has been loaded.
		}

		// first load request for this file.
		l.loaded[filename] = false // mark as loading
		l.loadAssetReq <- filename // put request on loader goroutine channel.
	}
}

// dispose is called when the engine is shutting down.
func (l *assetLoader) dispose() {
	// Close the worker queue since there are no more sends,
	// however keep the receiving channel open in case there
	// are workers trying to write to it.
	close(l.loadAssetReq) // shuts down idle workers.
}

// loadAssets checks for asset data from the goroutines and turns
// any data into loaded assets. As it is expected to run on the main thread
// each update tick, it will limit the amount of time spent processing assets
// so as to not stall the main loop.
func (l *assetLoader) loadAssets(rc render.Loader, ac audio.Loader) (assetsCreated int) {
	var timeUsed time.Duration
	start := time.Now()

	// track how much time is spent uploading assets
	// and return after a reasonable amount has elapsed.
	// FUTURE: possibly consume up to the amount of unused update time.
	timeLimit := 0.005 // max 5 milliseconds per update is reasonable.
	var err error
	assetsCreated = 0
	for timeUsed.Seconds() < timeLimit {
		select {
		case loaded := <-l.loadedAssets:

			// loaded should always contain one or more assets from a single file.
			if len(loaded) <= 0 {
				slog.Warn("investigate: no assets returned from worker")
				break
			}

			// filename helps uniquely identify the assets.
			filename := loaded[0].Filename
			ext := strings.ToLower(path.Ext(filename))
			name := strings.Replace(filename, ext, "", 1)
			assets := []asset{}
			for _, assetData := range loaded {
				if assetData.Err != nil {
					slog.Error("failed asset load", "filename", assetData.Filename, "error", assetData.Err)
					break // developer needs to debug why asset is missing.
				}
				switch data := assetData.Data.(type) {
				case load.MeshData:
					assetsCreated += 1
					msh := newMesh(name)
					msh.mid, err = rc.LoadMesh(data)
					if err != nil {
						slog.Error("LoadMesh failed", "error", err)
						break
					}
					assets = append(assets, msh)
					slog.Debug("new asset", "asset", "msh:"+msh.label(), "id", msh.mid, "filename", filename)
				case load.PBRMaterialData:
					assetsCreated += 1
					mat := newMaterial(name)
					mat.color = rgba{
						float32(data.ColorR),
						float32(data.ColorG),
						float32(data.ColorB),
						float32(data.ColorA),
					}
					mat.metallic = float32(data.Metallic)
					mat.roughness = float32(data.Roughness)
					assets = append(assets, mat)
					slog.Debug("new asset", "asset", "mat:"+mat.label(), "filename", filename)
				case *load.ImageData:
					assetsCreated += 1
					t := newTexture(name)
					t.opaque = data.Opaque
					t.tid, err = rc.LoadTexture(data)
					if err != nil {
						slog.Error("LoadTexture failed", "error", err)
						break
					}
					assets = append(assets, t)
					slog.Debug("new asset", "asset", "tex:"+t.label(), "id", t.tid, "opaque", t.opaque, "filename", filename)
				case *load.FontData:
					assetsCreated += 1
					f := newFont(name)
					f.setSize(data.W, data.H)
					for _, ch := range data.Chars {
						f.addChar(ch.Char, ch.X, ch.Y, ch.W, ch.H, ch.Xo, ch.Yo, ch.Xa)
					}
					assets = append(assets, f)
					slog.Debug("new asset", "asset", "fnt:"+f.label(), "filename", filename, "chars", len(f.chars))
				case *load.AudioData:
					assetsCreated += 1
					s := newSound(name)
					err = ac.LoadSound(&s.sid, &s.did, s.data) // upload audio data to audio device
					if err != nil {
						slog.Error("LoadSound failed", "error", err)
						break
					}
					assets = append(assets, s)
					slog.Debug("new asset", "asset", "snd:"+s.label(), "filename", filename)
				case *load.Shader:
					assetsCreated += 1
					s := newShader(name)
					s.setConfig(data)
					s.sid, err = rc.LoadShader(s.config)
					if err != nil {
						slog.Error("LoadShader failed", "error", err)
						break
					}
					assets = append(assets, s)
					slog.Debug("new asset", "asset", "shd:"+s.label(), "id", s.sid, "filename", filename)
				case load.ShaderData:
					// ignore since shader bytes are loaded directly from the render package.
					slog.Warn("load.ShaderBytes called") // unexpected. Testing?

				// case TODO animation asset from glb

				default:
					dtype := fmt.Sprintf("%T", data)
					slog.Error("unknown asset data", "datatype", dtype)
					break // developer needs sync code with the load package.
				}
			}
			l.loaded[filename] = true // mark file as as loaded

			// track loaded assets and notify requested asset listeners.
			for _, a := range assets {
				l.assets[a.aid()] = a // track loaded assets.

				// notify any outstanding requests for this asset.
				if reqs, ok := l.requests[a.aid()]; ok {
					for _, req := range reqs {
						req.callback(req.eid, a)
					}
				}
				delete(l.requests, a.aid())
			}

		default:
			l.loadLabels(rc)     // check outstanding label requests
			return assetsCreated // return if there are no loaded assets.
		}
		l.loadLabels(rc) // check outstanding label requests
		timeUsed = time.Since(start)
	}
	return assetsCreated
}

// loadLabels checks for outstanding label asset requests and
// creates the label mesh if all the assets are available.
func (l *assetLoader) loadLabels(rc render.Loader) {
	for aid, entities := range l.labelRequests {

		// get the matching font asset.
		if a, ok := l.assets[aid]; ok {
			for _, me := range entities {
				fnt, ok := a.(*font)
				if !ok {
					slog.Error("loadAssets expected font asset", "asset", a.label())
					break
				}
				str, wrap := me.labelData()
				sx, sy, md := fnt.setStr(str, wrap)

				// upload the mesh data
				mid, err := rc.LoadMesh(md)
				if err != nil {
					slog.Error("generateLabelMesh:LoadMesh failed", "error", err)
					break
				}

				// set the mesh asset on the label
				msh := newMesh(fmt.Sprintf("label%04d", mid))
				msh.mid = mid
				slog.Debug("new label mesh", "asset", "msh:"+msh.label(), "id", msh.mid)
				me.setLabelMesh(msh, sx, sy)
			}
			delete(l.labelRequests, aid)
		}
	}
}

// =============================================================================
// importWorker processes asset import requests on a goroutine.

// importWorker imports assets from persistent store. Currently persistent
// storage are just files on disk. This runs as a goroutine taking requests
// for assets from the needAsset channel and returning the loaded asset
// on the fetched channel.
func importWorker(wid int, loadAssetRequest <-chan string, loadedAssets chan<- []load.AssetData) {
	for filename := range loadAssetRequest {
		assetData := load.LoadAssetFile(filename)
		loadedAssets <- assetData
	}
}

// =============================================================================
// default assets

// loadDefaultAssets creates and pre-loads some basic and fallback assets.
func (l *assetLoader) loadDefaultAssets(rc render.Loader) (err error) {
	generateDefaultMeshes()

	// create default meshes
	m := newMesh("icon")
	m.mid, err = rc.LoadMesh(iconMeshData)
	if err != nil {
		return fmt.Errorf("LoadMesh icon: %w", err)
	}
	l.assets[m.aid()] = m
	slog.Debug("new asset", "asset", "msh:"+m.label(), "id", m.mid)

	m = newMesh("quad")
	m.mid, err = rc.LoadMesh(quadMeshData)
	if err != nil {
		return fmt.Errorf("LoadMesh quad: %w", err)
	}
	l.assets[m.aid()] = m
	slog.Debug("new asset", "asset", "msh:"+m.label(), "id", m.mid)

	m = newMesh("cube")
	m.mid, err = rc.LoadMesh(cubeMeshData)
	if err != nil {
		return fmt.Errorf("LoadMesh cube: %w", err)
	}
	l.assets[m.aid()] = m
	slog.Debug("new asset", "asset", "msh:"+m.label(), "id", m.mid)

	m = newMesh("circle")
	m.mid, err = rc.LoadMesh(circleMeshData)
	if err != nil {
		return fmt.Errorf("LoadMesh circle: %w", err)
	}
	l.assets[m.aid()] = m
	slog.Debug("new asset", "asset", "msh:"+m.label(), "id", m.mid)

	// create a basic texture to use for testing.
	t := newTexture("test")
	size, pixels := generateDefaultTexture()
	img := &load.ImageData{Width: size, Height: size, Pixels: pixels}
	t.tid, err = rc.LoadTexture(img)
	if err != nil {
		return fmt.Errorf("LoadTexture test: %w", err)
	}
	l.assets[t.aid()] = t
	slog.Debug("new asset", "asset", "tex:"+t.label(), "id", t.tid)

	// create a basic shader as the first shader.
	s := newShader("tex3D")
	s.config, err = load.ShaderConfig("tex3D.shd")
	if err != nil {
		return fmt.Errorf("load.ShaderConfig tex3D.shd: %w", err)
	}
	s.sid, err = rc.LoadShader(s.config)
	if err != nil {
		return fmt.Errorf("LoadShader tex3D: %w", err)
	}
	l.assets[s.aid()] = s
	slog.Debug("new asset", "asset", "shd:"+s.label(), "id", s.sid)
	return nil
}

func generateDefaultMeshes() {
	// 1x1 vertical 2D plane with tex-coordinates.
	iconMeshData[load.Vertexes] = load.F32Buffer(iconVerts, 2)
	iconMeshData[load.Texcoords] = load.F32Buffer(iconTexuv, 2)
	iconMeshData[load.Indexes] = load.U16Buffer(iconIndex)

	// 2x2 plane with normals and tex-coordinates.
	// Based on Blender default plane.
	quadMeshData[load.Vertexes] = load.F32Buffer(quadVerts, 3)
	quadMeshData[load.Texcoords] = load.F32Buffer(quadTexuv, 2)
	quadMeshData[load.Normals] = load.F32Buffer(quadNorms, 3)
	quadMeshData[load.Indexes] = load.U16Buffer(quadIndex)

	// 1x1x1 cube with normals and tex-coordinates for all 6 sides.
	// based on Blender default cube.
	cubeMeshData[load.Vertexes] = load.F32Buffer(cubeVerts, 3)
	cubeMeshData[load.Texcoords] = load.F32Buffer(cubeTexuv, 2)
	cubeMeshData[load.Normals] = load.F32Buffer(cubeNorms, 3)
	cubeMeshData[load.Indexes] = load.U16Buffer(cubeIndex)

	// unit circle by connecting 1000 points.
	circleVerts, circleIndex := unitCircle()
	circleMeshData[load.Vertexes] = load.F32Buffer(circleVerts, 3)
	circleMeshData[load.Indexes] = load.U16Buffer(circleIndex)
}

// 1x1 vertical 2D plane with tex-coordinates.
// Based on Blender default plane y-up.
var iconMeshData = make(load.MeshData, load.VertexTypes)
var iconVerts = []float32{
	-0.5, +0.5,
	+0.5, +0.5,
	-0.5, -0.5,
	+0.5, -0.5,
}
var iconTexuv = []float32{
	0.0, 1.0,
	1.0, 1.0,
	0.0, 0.0,
	1.0, 0.0,
}
var iconIndex = []uint16{
	0, 1, 3,
	0, 3, 2,
}

// 1x1 front facing plane for 3D labels and billboards.
var quadMeshData = make(load.MeshData, load.VertexTypes)
var quadVerts = []float32{
	-0.5, +0.5, 0.0, // top left
	-0.5, -0.5, 0.0, // bottom left
	+0.5, +0.5, 0.0, // top right
	+0.5, -0.5, 0.0, // bottom right
}
var quadNorms = []float32{
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
}
var quadTexuv = []float32{
	0.0, 1.0,
	0.0, 0.0,
	1.0, 1.0,
	1.0, 0.0,
}
var quadIndex = []uint16{
	2, 1, 3,
	2, 0, 1,
}

// 1x1x1 cube with normals and tex-coordinates for all 6 sides.
// based on Blender default cube.
var cubeMeshData = make(load.MeshData, load.VertexTypes)
var cubeVerts = []float32{
	+0.5, -0.5, +0.5,
	+0.5, -0.5, +0.5,
	+0.5, -0.5, +0.5,
	+0.5, +0.5, +0.5,
	+0.5, +0.5, +0.5,
	+0.5, +0.5, +0.5,
	-0.5, +0.5, +0.5,
	-0.5, +0.5, +0.5,
	-0.5, +0.5, +0.5,
	-0.5, -0.5, +0.5,
	-0.5, -0.5, +0.5,
	-0.5, -0.5, +0.5,
	+0.5, -0.5, -0.5,
	+0.5, -0.5, -0.5,
	+0.5, -0.5, -0.5,
	+0.5, +0.5, -0.5,
	+0.5, +0.5, -0.5,
	+0.5, +0.5, -0.5,
	-0.5, +0.5, -0.5,
	-0.5, +0.5, -0.5,
	-0.5, +0.5, -0.5,
	-0.5, -0.5, -0.5,
	-0.5, -0.5, -0.5,
	-0.5, -0.5, -0.5,
}
var cubeNorms = []float32{
	+0.0, -1.0, -0.0,
	+0.0, +0.0, +1.0,
	+1.0, +0.0, -0.0,
	+0.0, +0.0, +1.0,
	+0.0, +1.0, -0.0,
	+1.0, +0.0, -0.0,
	-1.0, +0.0, -0.0,
	+0.0, +0.0, +1.0,
	+0.0, +1.0, -0.0,
	-1.0, +0.0, -0.0,
	+0.0, -1.0, -0.0,
	+0.0, +0.0, +1.0,
	+0.0, -1.0, -0.0,
	+0.0, +0.0, -1.0,
	+1.0, +0.0, -0.0,
	+0.0, +0.0, -1.0,
	+0.0, +1.0, -0.0,
	+1.0, +0.0, -0.0,
	-1.0, +0.0, -0.0,
	+0.0, +0.0, -1.0,
	+0.0, +1.0, -0.0,
	-1.0, +0.0, -0.0,
	+0.0, -1.0, -0.0,
	+0.0, +0.0, -1.0,
}
var cubeTexuv = []float32{
	+0.625, +0.50,
	+0.625, +0.50,
	+0.625, +0.50,
	+0.375, +0.50,
	+0.375, +0.50,
	+0.375, +0.50,
	+0.625, +0.25,
	+0.625, +0.25,
	+0.625, +0.25,
	+0.375, +0.25,
	+0.375, +0.25,
	+0.375, +0.25,
	+0.625, +0.75,
	+0.625, +0.75,
	+0.875, +0.50,
	+0.375, +0.75,
	+0.125, +0.50,
	+0.375, +0.75,
	+0.625, +1.0,
	+0.625, +0.0,
	+0.875, +0.25,
	+0.375, +1.0,
	+0.125, +0.25,
	+0.375, +0.0,
}
var cubeIndex = []uint16{
	1, 3, 7, 1, 7, 11,
	13, 23, 19, 13, 19, 15,
	2, 14, 17, 2, 17, 5,
	4, 16, 20, 4, 20, 8,
	6, 18, 21, 6, 21, 9,
	12, 0, 10, 12, 10, 22,
}

// circleMeshData holds points generated by unitCircle.
var circleMeshData = make(load.MeshData, load.VertexTypes)

// unitCircle creates a set of points on the circumference of a unit circle
// of radius 1.0. Expected to be rendered as a line.
func unitCircle() (verts []float32, indexes []uint16) {
	count := 100
	fraction := (1.0 / float64(count)) * (math.Pi * 2.0)
	angle := 0.0
	for i := 0; i < count; i++ {
		x := float32(math.Cos(angle))
		y := float32(math.Sin(angle))
		angle += fraction

		// 0 to count-1 vertexes of 3 float32
		verts = append(verts, x, y, 0.0) // xyz per vertex

		// line indexes are 0:1, 1:2, ... count-2:count-1, count-1:0
		indexes = append(indexes, uint16(i), uint16(i+1))
		if i == count-1 {
			indexes[len(indexes)-1] = 0 // count-1:0
		}
	}
	return verts, indexes
}

// generateDefaultTexture creates a texture pattern that can be
// used as placeholders for other textures.
// Dump for testing using:
//   - f, _ := os.Create("image.png")
//   - png.Encode(f, img)
func generateDefaultTexture() (squareSize uint32, pixels []byte) {
	size := 256 // with and height in pixels

	// non-alpha-premultiplied 32-bit color
	img := image.NewNRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{size, size},
	})

	// rgba uint8
	gray := color.NRGBA{100, 100, 100, 0xff}

	// create a checkerboard pattern
	quad := size / 32 // 8x8 blocks
	for x := 0; x < size; x++ {
		val := (x / quad) % 2
		for y := 0; y < size; y++ {
			val2 := (y / quad) % 2
			if (val+val2)%2 == 0 {
				img.Set(x, y, gray)
			} else {
				img.Set(x, y, color.White)
			}
		}
	}
	return uint32(size), []byte(img.Pix)
}
