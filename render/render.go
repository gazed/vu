// Copyright © 2024 Galvanized Logic Inc.

package render

// render.go provides API wrappers for the render specific APIs.

import (
	"fmt"
	"log/slog"
	"math"
	"time"
	"unsafe"

	"github.com/gazed/vu/device"
	"github.com/gazed/vu/load"
	"github.com/gazed/vu/math/lin"
)

// RenderAPI enumerates the possible render backends.
type RenderAPI int

// Vulkan is currently the only supported render API There are no plans to
// support OpenGL or any DirectX versions before DX12. Unlikely futures include:
//   - Nintendo    NVN - proprietary...unlikely to ship golang to this platform.
//   - Playstation GNM - proprietary...unlikely to ship golang to this platform.
const (
	VULKAN_RENDERER RenderAPI = iota // windows, linux, android
	DX12_RENDERER                    // FUTURE: xbox
	METAL_RENDERER                   // FUTURE: iOS, macOS, tvOS, watchOS, visionOS
)

// New creates an initialized renderer and returns a render context.
func New(api RenderAPI, dev *device.Device, appTitle string) (rc *Context, err error) {
	switch api {
	case VULKAN_RENDERER:
		vr, err := getVulkanRenderer(dev, appTitle)
		if err != nil {
			return nil, fmt.Errorf("render create failed %w", err)
		}
		return &Context{renderer: vr}, nil
	}
	return nil, fmt.Errorf("unsupported render API: %d", api)
}

// Context holds data for the rendering system and wraps the API
// specific renderers, ie: Vulkan, DX12, Metal.
type Context struct {
	renderer    renderAPI // Render API wrapper
	frameNumber int64     // frame counter
}

// Dispose releases renderer resources.
func (c *Context) Dispose() {
	if c.renderer != nil {
		c.renderer.dispose()
	}
	c.renderer = nil
}

// Draw renders the given render passes for one frame.
// Expected to be called many times per second.
func (c *Context) Draw(passes []Pass, dt time.Duration) (err error) {
	if c.renderer == nil {
		return fmt.Errorf("renderer not intiialized")
	}
	// FUTURE: do something with delta time which is currently ignored.

	// an error in beginFrame may not be a problem.
	if err = c.renderer.beginFrame(dt); err != nil {
		slog.Debug("beginFrame", "error", err)
		return nil // ignore this frame and keep going
	}

	// errors in drawFrame or endFrame are always a problem.
	if err = c.renderer.drawFrame(passes); err != nil {
		return fmt.Errorf("render.RecordFrame: %w", err)
	}
	if err = c.renderer.endFrame(dt); err != nil {
		return fmt.Errorf("render.EndFrame: %w", err)
	}
	c.frameNumber++
	return nil
}

// Resize updates the graphics resources to the given size.
// Expected to be called when the user resizes the app window.
func (c *Context) Resize(width, height uint32) { c.renderer.resize(width, height) }

// Size returns the current render surface size.
func (c *Context) Size() (width, height uint32) { return c.renderer.size() }

// LoadTexture creates GPU texture resources and uploads
// texture data to the GPU.
func (c *Context) LoadTexture(img *load.ImageData) (tid uint32, err error) {
	return c.renderer.loadTexture(img.Width, img.Height, img.Pixels)
}

// UpdateTexture updates the GPU texture data for the given texture ID.
// Do not update Textures that are being rendered. Double buffer the textures
// to update a texture and then swap for the rendered texture. UpdateTexture
// ignores textures that are not exactly the same size as the existing texture.
func (c *Context) UpdateTexture(tid uint32, img *load.ImageData) (err error) {
	return c.renderer.updateTexture(tid, img.Width, img.Height, img.Pixels)
}

// DropTexture removes the GPU texture resources
// for the given texture ID.
func (c *Context) DropTexture(tid uint32) { c.renderer.dropTexture(tid) }

// LoadMeshes allocates GPU resources for the mesh data.
func (c *Context) LoadMeshes(msh []load.MeshData) (mids []uint32, err error) {
	return c.renderer.loadMeshes(msh)
}

// LoadMesh allocates GPU resources for the mesh data.
func (c *Context) LoadMesh(msh load.MeshData) (mid uint32, err error) {
	mids, err := c.renderer.loadMeshes([]load.MeshData{msh})
	if err != nil || len(mids) != 1 {
		return 0, fmt.Errorf("LoadMesh %d %w", len(mids), err)
	}
	return mids[0], nil
}

// DropMesh discards the mesh resources.
func (c *Context) DropMesh(mid uint32) { c.renderer.dropMesh(mid) }

// LoadInstanceData allocates GPU resources for the instanced mesh data.
func (c *Context) LoadInstanceData(data []load.Buffer) (iid uint32, err error) {
	return c.renderer.loadInstanceData(data)
}

// UpdateInstanceData updates the GPU instance data for the given instance data ID.
// Do not update instance data that is being rendered. Double buffer the instance
// data and then swap with the rendered instance data. UpdateInstanceData ignores data
// buffers that are not exactly the same sizes and types as the existing data buffers.
func (c *Context) UpdateInstanceData(iid uint32, data []load.Buffer) (err error) {
	return c.renderer.updateInstanceData(iid, data)
}

// DropInstanced discards the instanced resources.
func (c *Context) DropInstanceData(iid uint32) { c.renderer.dropInstanceData(iid) }

// LoadShader prepare the GPU indicated GPU shader for rendering.
func (c *Context) LoadShader(config *load.Shader) (sid uint16, err error) {
	return c.renderer.loadShader(config)
}

// SetClearColor sets the color that is used to clear the display.
func (c *Context) SetClearColor(r, g, b, a float32) {
	c.renderer.setClearColor(r, g, b, a)
}

// The render context implements this interface.
// It allows engine tests to mock this part of the render context.
type Loader interface {
	LoadTexture(img *load.ImageData) (tid uint32, err error)
	UpdateTexture(tid uint32, img *load.ImageData) (err error)
	LoadMesh(msh load.MeshData) (mid uint32, err error)
	LoadMeshes(mdata []load.MeshData) (mids []uint32, err error)
	LoadShader(config *load.Shader) (mid uint16, err error)

	// FUTURE: LoadAnimation

}

// =============================================================================

// renderAPI is a generic set of render methods that must be implemented
// by the specific render APIs ie: vulkan, dx12, metal
type renderAPI interface {
	dispose() // called once on shutdown

	// set the default background clear color.
	setClearColor(r, g, b, a float32)

	// render a frame.
	beginFrame(deltaTime time.Duration) error
	drawFrame(passes []Pass) error
	endFrame(deltaTime time.Duration) error

	// render resize controls.
	size() (width, height uint32) // returns current size
	resize(width, height uint32)  // request size change
	isResizing() bool             // true when size is updating.

	// create a GPU texture and upload the mesh data.
	loadTexture(w, h uint32, pixels []byte) (tid uint32, err error)
	updateTexture(tid, w, h uint32, pixels []byte) (err error)
	dropTexture(tid uint32) // release texture resources

	// create a GPU shader using the given shader configuration
	loadShader(config *load.Shader) (sid uint16, err error)
	dropShader(sid uint16) // release shader resources

	// create GPU meshes by uploading the mesh vertex data.
	// return an identifier for each mesh.
	loadMeshes(msh []load.MeshData) (mid []uint32, err error)
	// FUTURE: updateMesh() replace mesh with new mesh data.
	dropMesh(mid uint32)

	// load instance data for an instanced mesh.
	// return an identifier for the instance data.
	loadInstanceData(data []load.Buffer) (iid uint32, err error)
	updateInstanceData(iid uint32, data []load.Buffer) (err error)
	dropInstanceData(iid uint32)
}

// =============================================================================

// uniformSets describes data passed to shader programs.
// It is generated from load.ShaderUniform configuration data.
type uniformSets struct {
	sceneSize    uint32    // set0: total scene uniforms byte size.
	materialSize uint32    // set1: total material uniforms byte size.
	modelSize    uint32    // set2: total model uniforms byte size.
	numSamplers  uint32    // number of uniform samplers.
	uniforms     []uniform // per-uniform data.

	// index maps the uniform names to the uniform data
	index map[string]*uniform // pointers to uniforms slice data.
}

// uniform describes a single uniform and is generated from
// a shader configuration.
type uniform struct {
	scope     load.UniformScope  // matches descriptor set, ie: scene is set:0
	offset    uint32             // uniform: offset is start of data bytes in buffer
	size      uint32             // uniform: data size in bytes.
	bind      uint32             // sampler: the sampler bind location
	passUID   load.PassUniform   // pass data index.
	packetUID load.PacketUniform // packet data index.
}

// hasUniform if the shader supports the given uniform.
func (us uniformSets) hasUniform(name string) bool {
	_, ok := us.index[name]
	return ok
}

// Current uniform data size limits. Increase limits as shader
// complexity increases. Note that model data is a hard limit for
// using push constants.
//
// Due to buffer alignment when setting data.
// The MinUniformBufferOffsetAlignment is at worst 256 bytes
// so create each uniform buffer of 256 bytes and complain if
// the uniform data exceeds this.
const (
	maxSceneUniformBytes    = 256 // scene data fits in 256 bytes
	maxMaterialUniformBytes = 256 // material data fits in 256 bytes
	maxModelUniformBytes    = 128 // model data fits in 128 bytes
)

// genUniforms creates shaderUniforms from shader the configuration.
func getUniformSets(configUniforms []load.ShaderUniform) (sets uniformSets) {
	sets.uniforms = make([]uniform, len(configUniforms))
	sets.index = map[string]*uniform{}
	for i, cu := range configUniforms {
		u := &sets.uniforms[i]
		u.scope = cu.Scope
		if cu.DataType == load.DataType_SAMPLER {
			u.bind = sets.numSamplers
			sets.numSamplers += 1
		} else {
			u.size = load.DataTypeSizes[cu.DataType]
			switch cu.Scope {
			case load.SceneScope:
				u.offset = sets.sceneSize
				sets.sceneSize += u.size
			case load.MaterialScope:
				u.offset = sets.materialSize
				sets.materialSize += u.size
			case load.ModelScope:
				u.offset = sets.modelSize
				sets.modelSize += u.size
			}
			u.passUID = cu.PassUID     // one of these two...
			u.packetUID = cu.PacketUID // ...will be valid.
		}
		sets.index[cu.Name] = u
	}

	// complain if a shader exceeds the amount of allocated uniform bytes.
	// Either increase available space or rework the shader.
	if sets.sceneSize > maxSceneUniformBytes {
		slog.Error("need to increase uniformBufferSize", "set0_scene", sets.sceneSize)
	}
	if sets.materialSize > maxMaterialUniformBytes {
		slog.Error("need to increase uniformBufferSize", "set1_material", sets.materialSize)
	}
	if sets.modelSize > maxModelUniformBytes {
		slog.Error("need to increase uniformBufferSize", "set2_model", sets.modelSize)
	}
	return sets
}

// =============================================================================
// the render data structures below are used to set shader uniform data.

// V4ToBytes returns a byte slice of float32 for the given float64 vector.
// The given byte slice is zeroed and returned filled with the vector bytes.
func V4ToBytes(v *lin.V4, bytes []byte) []byte {
	bytes = bytes[:0]
	return append(bytes, (&v4{}).set64(v).toBytes()...)
}

// V4SToBytes returns a byte slice of float32 for the given float64s
// The given byte slice is zeroed and returned filled with the vector bytes.
func V4SToBytes(x, y, z, w float64, bytes []byte) []byte {
	bytes = bytes[:0]
	return append(bytes, (&v4{}).set64S(x, y, z, w).toBytes()...)
}

// V4S32ToBytes returns a byte slice of float32 from the given values.
// The given byte slice is zeroed and returned filled with the float bytes.
func V4S32ToBytes(x, y, z, w float32, bytes []byte) []byte {
	bytes = bytes[:0]
	return append(bytes, (&v4{}).setS(x, y, z, w).toBytes()...)
}

// v4 is a vec4 of float32 that is used to set shader uniforms.
type v4 struct{ x, y, z, w float32 }

func (v *v4) set64(v64 *lin.V4) *v4 {
	v.x, v.y, v.z, v.w = float32(v64.X), float32(v64.Y), float32(v64.Z), float32(v64.W)
	return v
}
func (v *v4) set64S(x, y, z, w float64) *v4 {
	v.x, v.y, v.z, v.w = float32(x), float32(y), float32(z), float32(w)
	return v
}
func (v *v4) setS(x, y, z, w float32) *v4 {
	v.x, v.y, v.z, v.w = x, y, z, w
	return v
}

// toBytes returns the data as a byte array.
func (v *v4) toBytes() []byte {
	return (*[int(unsafe.Sizeof(*v))]byte)(unsafe.Pointer(v))[:]
}
func (v *v4) setInvalid() { v.x, v.y, v.z, v.w = -1, -1, -1, -1 }

// =============================================================================

// M4ToBytes returns a byte slice of float32 for the given float64 matrix.
// The given byte slice is zeroed and returned filled with the matrix bytes.
func M4ToBytes(m *lin.M4, bytes []byte) []byte {
	bytes = bytes[:0]
	return append(bytes, (&m4{}).set64(m).toBytes()...)
}

// m4 is a 4x4 matrix of float32 that is used to set shader uniforms.
type m4 struct {
	xx, xy, xz, xw float32 // indices 0, 1, 2, 3  [00, 01, 02, 03] X-Axis
	yx, yy, yz, yw float32 // indices 4, 5, 6, 7  [10, 11, 12, 13] Y-Axis
	zx, zy, zz, zw float32 // indices 8, 9, a, b  [20, 21, 22, 23] Z-Axis
	wx, wy, wz, ww float32 // indices c, d, e, f  [30, 31, 32, 33]
}

func (m *m4) set64(m64 *lin.M4) *m4 {
	m.xx, m.xy, m.xz, m.xw = float32(m64.Xx), float32(m64.Xy), float32(m64.Xz), float32(m64.Xw)
	m.yx, m.yy, m.yz, m.yw = float32(m64.Yx), float32(m64.Yy), float32(m64.Yz), float32(m64.Yw)
	m.zx, m.zy, m.zz, m.zw = float32(m64.Zx), float32(m64.Zy), float32(m64.Zz), float32(m64.Zw)
	m.wx, m.wy, m.wz, m.ww = float32(m64.Wx), float32(m64.Wy), float32(m64.Wz), float32(m64.Ww)
	return m
}

// V16ToBytes returns a byte slice of float32 from the given float64.
// The given byte slice is zeroed and returned filled with the float bytes.
func V16ToBytes(args []float64, bytes []byte) []byte {
	bytes = bytes[:0]
	m := &m4{}
	m.xx, m.xy, m.xz, m.xw = float32(args[0]), float32(args[1]), float32(args[2]), float32(args[3])
	m.yx, m.yy, m.yz, m.yw = float32(args[4]), float32(args[5]), float32(args[6]), float32(args[7])
	m.zx, m.zy, m.zz, m.zw = float32(args[8]), float32(args[9]), float32(args[10]), float32(args[11])
	m.wx, m.wy, m.wz, m.ww = float32(args[12]), float32(args[13]), float32(args[14]), float32(args[15])
	return append(bytes, m.toBytes()...)
}

// ToBytes returns the data as a byte array.
func (m *m4) toBytes() []byte {
	return (*[int(unsafe.Sizeof(*m))]byte)(unsafe.Pointer(m))[:]
}

// =============================================================================

// Light holds the location and color for a directional light.
// Effectively 2 float32 vec4 for 32 bytes.
type Light struct {
	X, Y, Z, W float32 // location - W is 1 for point light, 0 for directional.
	R, G, B    float32 // color
	Intensity  float32 // light intensity
}

// reset the light data before reusing the light struct.
// Called internally from pass.Reset()
func (l *Light) reset() {
	l.X, l.Y, l.Z, l.W = 0.0, 0.0, 0.0, 0.0
	l.R, l.G, l.B = 0.0, 0.0, 0.0
	l.Intensity = 0.0
}

// LightsToBytes converts a slice of lights to bytes.
// The given byte slice is zeroed and returned filled with the given light data.
func LightsToBytes(lights []Light, bytes []byte) []byte {
	bytes = bytes[:0]
	const maxLights = 3
	lbytes := (*[int(unsafe.Sizeof(Light{})) * maxLights]byte)(unsafe.Pointer(&lights[0]))[:]
	return append(bytes, lbytes...)
}

// =============================================================================

// U8ToBytes returns a byte slice containing the given uint8
// because the smallest uniform is an int.
// The given byte slice is zeroed and returned filled with the uint8.
func U8ToBytes(val uint8, bytes []byte) []byte {
	bytes = bytes[:0]
	return append(bytes, []byte{byte(val), 0, 0, 0}...)
}

// Int32ToBytes returns a byte slice containing the given int32.
// The given byte slice is zeroed and returned filled with the int32.
func Int32ToBytes(val int32, bytes []byte) []byte {
	bytes = bytes[:0]
	b0 := byte(val)
	b1 := byte(val >> 8)
	b2 := byte(val >> 16)
	b3 := byte(val >> 24)
	return append(bytes, b0, b1, b2, b3)
}

// Float32ToBytes returns a byte slice containing the given int32.
// The given byte slice is zeroed and returned filled with the int32.
func Float32ToBytes(val float32, bytes []byte) []byte {
	return Int32ToBytes(int32(math.Float32bits(val)), bytes)
}
