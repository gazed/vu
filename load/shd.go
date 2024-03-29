// Copyright © 2024 Galvanized Logic Inc.

package load

// shd.go reads a shader description from disk.
// Shader descriptions are used by the engine to map
// the shader data to the shader programs.

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// map the shader configuration to data structures.
// Some maps are public so external apps can add new shader attributes
// as needed by new shaders.
//
// FUTURE: use shader reflection to generate the shader descriptions
// and maps from the shader code.
var shaderStages = map[string]ShaderStage{
	"vert": Stage_VERTEX,
	"geom": Stage_GEOMETRY,
	"frag": Stage_FRAGMENT,
}

// ShaderAttributes map glsl vertex attributes from shader
// programs to model data provided by code.
// Expected use is for passing data from the engine to the render system.
var ShaderAttributes = map[string]int{
	// names for model vertex data
	"position": Vertexes,
	"texcoord": Texcoords,
	"normal":   Normals,
	"tangent":  Tangents,
	"v_color":  Colors,
	"joint":    Joints,
	"weight":   Weights,

	// names for instanced model data
	"i_locus": InstanceLocus,
	"i_color": InstanceColors,
	"i_scale": InstanceScales,
}

// ShaderAttributeScope categorizes vertex attributes into
// per-vertex data or per-model-instance data.
// Expected use is for passing data from the engine to the render system.
var ShaderAttributeScope = map[string]AttributeScope{
	"vertex":   VertexAttribute,
	"instance": InstanceAttribute,
}

// ShaderAttributeData are the currently supported data types
// for vertex and model-instance data.
// Expected use is for passing data from the engine to the render system.
var ShaderAttributeData = map[string]ShaderDataType{
	"float": DataType_FLOAT,
	"vec2":  DataType_VEC2,
	"vec3":  DataType_VEC3,
	"vec4":  DataType_VEC4,
}

// ShaderUniformScope is used to bind data in shaders.
// Expected use is for passing data from the engine to the render system.
var ShaderUniformScope = map[string]UniformScope{
	"scene":    SceneScope,    // global
	"material": MaterialScope, // per shader
	"model":    ModelScope,    // per model
}

// ShaderPassUniforms are shader uniforms that apply to the
// a single scene (render pass).
// Expected use is for passing data from the engine to the render system.
var ShaderPassUniforms = map[string]PassUniform{
	"proj":    PROJ,
	"view":    VIEW,
	"cam":     CAM,
	"lights":  LIGHTS,
	"nlights": NLIGHTS,
}

// ShaderPacketUniforms are shader uniforms that apply to one model.
// Render data for a model is put into a render.Packet.
// Expected use is for passing data from the engine to the render system.
var ShaderPacketUniforms = map[string]PacketUniform{
	"model":    MODEL,
	"scale":    SCALE,
	"color":    COLOR,
	"material": MATERIAL,
}

// ShaderUniformData are the supported uniform data types.
// Expected use is for passing data from the engine to the render system.
var ShaderUniformData = map[string]ShaderDataType{
	"int":     DataType_INT,
	"light3":  DataType_LIGHT3,
	"mat3":    DataType_MAT3,
	"mat4":    DataType_MAT4,
	"sampler": DataType_SAMPLER,
	"vec2":    DataType_VEC2,
	"vec3":    DataType_VEC3,
	"vec4":    DataType_VEC4,
}

// Shd loads a yaml shader configuration and returns it as
// a shader configuration struct. This is needed to provide
// a link between the render models and the shader programs.
func Shd(name string, data []byte) (shader *Shader, err error) {
	var cfg shaderConfig
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return shader, fmt.Errorf("Shd: yaml %w", err)
	}

	// convert the shader stage strings to a bit-mask.
	stages := ShaderStage(0)
	for _, stg := range cfg.Stages {
		stage, ok := shaderStages[stg]
		if !ok {
			return shader, fmt.Errorf("Shd:unsupported stage %s", stg)
		}
		stages |= stage
	}

	// convert the attribute descriptions
	attrs := []ShaderAttribute{}
	for _, a := range cfg.Attrs {
		atype, ok1 := ShaderAttributes[a.Name]
		if !ok1 {
			return shader, fmt.Errorf("Shd:unsupported shader attribute %s", a.Name)
		}
		ascope, ok2 := ShaderAttributeScope[a.Scope]
		if !ok2 {
			return shader, fmt.Errorf("Shd:unsupported attribute scope %s", a.Scope)
		}
		dtype, ok3 := ShaderAttributeData[a.Data]
		if !ok3 {
			return shader, fmt.Errorf("Shd:unsupported attribute data %s", a.Data)
		}
		attrs = append(attrs, ShaderAttribute{
			Name:      a.Name,
			AttrType:  atype,
			AttrScope: ascope,
			DataType:  dtype,
		})
	}

	// convert the uniform descriptions
	uniforms := []ShaderUniform{}
	for _, u := range cfg.Uniforms {
		scope, ok1 := ShaderUniformScope[u.Scope]
		if !ok1 {
			return shader, fmt.Errorf("Shd:unsupported uniform scope %s", u.Scope)
		}
		dtype, ok2 := ShaderUniformData[u.Data]
		if !ok2 {
			return shader, fmt.Errorf("Shd:unsupported uniform data %s", u.Data)
		}

		// uniform must be in the render pass or in the render packet.
		passID, ok3 := ShaderPassUniforms[u.Name]
		packetID, ok4 := ShaderPacketUniforms[u.Name]
		if !ok3 && !ok4 {
			return shader, fmt.Errorf("Shd:unsupported uniform %s", u.Name)
		}
		uniforms = append(uniforms, ShaderUniform{
			Name:      u.Name,
			Scope:     scope,
			DataType:  dtype,
			PassUID:   passID,
			PacketUID: packetID,
		})
	}

	// create the shader configuration.
	shader = &Shader{
		Name:     cfg.Name,
		Pass:     cfg.Pass,
		Stages:   stages,
		Attrs:    attrs,
		Uniforms: uniforms,
	}

	// add the render flags to configure the shader pipeline,
	// FUTURE: more render flags as needed.
	if cfg.Render != "" {
		shader.CullModeNone = strings.Contains(cfg.Render, "cullOff")
		shader.DrawLines = strings.Contains(cfg.Render, "drawLines")
	}

	// return the shader
	return shader, nil
}

// shaderConfig is used to load string based shader configuration.
// The yaml is string based so that it is easier to read.
type shaderConfig struct {
	Name   string   `yaml:"name"`
	Pass   string   `yaml:"pass"`
	Stages []string `yaml:"stages"`
	Render string   `yaml:"render"` // render flags.
	Attrs  []struct {
		Name  string `yaml:"name"`
		Data  string `yaml:"data"`
		Scope string `yaml:"scope"` // vertex or instanced
	} `yaml:"attrs"`
	Uniforms []struct {
		Name  string `yaml:"name"`
		Data  string `yaml:"data"`
		Scope string `yaml:"scope"`
	} `yaml:"uniforms"`
}

// =============================================================================
// shader structs needed to create shader and to bridge the render object
// data with the uniforms needed by the shader.

// Shader describes shader data and is used to generate
// render shaders on startup.
//
// FUTURE: replace the yaml files with data gathered from shader code reflection.
// Not really a simple way to do Spirv reflection at the moment as the current
// tools generate lots of data which would be a pain to parse.
// See: spirv-cross, spirv-reflect command line tools.
type Shader struct {
	// Name identifies the shader modules. eg: Name.frag, Name.vert
	Name   string      // unique name for this shader.
	Pass   string      // renderpass name for this shader.
	Stages ShaderStage // bit flags for the shader stages.

	// Set from shaderConfig.Render flags.
	CullModeNone bool // true disables backface culling.
	DrawLines    bool // true to render lines instead of triangles.

	// Attrs must match the shader attributes in name and position, ie:
	//   Attr[0].Name == position   ... which matches
	//	 "layout(location=0) in vec3 position;"
	Attrs []ShaderAttribute

	// Uniforms must match the shader attributes in name and position
	// where scope also identifies a DescriptorSet.
	Uniforms []ShaderUniform
}

// ShaderAttribute identifies the data layouts for attributes.
type ShaderAttribute struct {
	Name      string         // unique name matching shader source code.
	AttrType  int            // matching mesh vertex attribute type
	AttrScope AttributeScope // vertex or instanced
	DataType  ShaderDataType // describes the attribute data.
}

// ShaderUniform identifies the data layouts for uniforms.
type ShaderUniform struct {
	Name      string         // unique name matching shader source code.
	Scope     UniformScope   //
	DataType  ShaderDataType //
	PassUID   PassUniform    // index to pass data
	PacketUID PacketUniform  // index to packet data
}

func (s *Shader) GetSceneUniforms() (uniforms []*ShaderUniform) {
	for i, uni := range s.Uniforms {
		if uni.DataType != DataType_SAMPLER && uni.Scope == SceneScope {
			uniforms = append(uniforms, &s.Uniforms[i])
		}
	}
	return uniforms
}
func (s *Shader) GetMaterialUniforms() (uniforms []*ShaderUniform) {
	for i, uni := range s.Uniforms {
		if uni.DataType != DataType_SAMPLER && uni.Scope == MaterialScope {
			uniforms = append(uniforms, &s.Uniforms[i])
		}
	}
	return uniforms
}
func (s *Shader) GetSamplerUniforms() (uniforms []*ShaderUniform) {
	for i, uni := range s.Uniforms {
		if uni.DataType == DataType_SAMPLER && uni.Scope == MaterialScope {
			uniforms = append(uniforms, &s.Uniforms[i])
		}
	}
	return uniforms
}

// =============================================================================

// ShaderPass identifies the render pass that the shader uses.
type ShaderPass uint8

const (
	Renderpass_3D ShaderPass = iota // 3D is rendered before 2D
	Renderpass_2D                   //
)

// PassUniform is scene level data.
type PassUniform uint8

const (
	PROJ         PassUniform = iota // scene
	VIEW                            // scene
	CAM                             // scene
	LIGHTS                          // scene
	NLIGHTS                         // scene
	PassUniforms                    // must be last
)

// =============================================================================

// ShaderStage identifies the currently supported programmable
// shader stages that a shader can have.
type ShaderStage uint8

const (
	Stage_VERTEX   ShaderStage = 1 << iota // vertex processing.
	Stage_GEOMETRY                         // eg: turn points into quads.
	Stage_FRAGMENT                         // pixel processing.
)

// =============================================================================

// AttributeScope identifies the two types of attributes.
type AttributeScope uint8

const (
	VertexAttribute AttributeScope = iota
	InstanceAttribute
)

// =============================================================================

// UniformScope identifies the three supported types of uniforms.
// The scope is directly related to the layout(set=x) value in the shader code.
type UniformScope uint8

const (
	SceneScope    UniformScope = iota // set per render pass: set=0
	MaterialScope                     // set per material     : set=1
	ModelScope                        // set per model      : push constants
)

// PacketUniform is model level data.
type PacketUniform uint8

const (
	MODEL          PacketUniform = iota // model
	SCALE                               // model
	COLOR                               // model
	MATERIAL                            // model
	PacketUniforms                      // must be last
)

// =============================================================================

// ShaderDataType helps describe shader attributes and uniforms
type ShaderDataType uint8

const (
	DataType_INT ShaderDataType = iota
	DataType_FLOAT
	DataType_LIGHT3 // array of 3 lights
	DataType_MAT3
	DataType_MAT4
	DataType_SAMPLER
	DataType_VEC2
	DataType_VEC3
	DataType_VEC4
)

var DataTypeSizes = map[ShaderDataType]uint32{
	DataType_INT:     4,  // int32
	DataType_FLOAT:   4,  // float32
	DataType_LIGHT3:  96, // 3 light struct of 2 vec4 float
	DataType_MAT3:    36, // 9 float32
	DataType_MAT4:    64, // 16 float32
	DataType_SAMPLER: 0,  // samplers don't have size.
	DataType_VEC2:    8,  // float32 vec2
	DataType_VEC3:    12, // float32 vec3
	DataType_VEC4:    16, // float32 vec4
}
