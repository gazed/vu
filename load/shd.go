// SPDX-FileCopyrightText : © 2022-2024 Galvanized Logic Inc.
// SPDX-License-Identifier: MIT

package load

// shd.go reads generated shader reflection data from disk.
// The reflection data is used to map shader parameter data
// needed by the shader programs.

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

// =============================================================================
// Shader supplies the GPU render layer shader pipelines with uniform data
// from the CPU models.
type Shader struct {
	// Name identifies the shader modules. eg: Name_frag, Name_vert
	Name   string      // unique name for this shader.
	Pass   string      // shader render pass.
	Stages ShaderStage // bit flags for the shader stages.

	// user defined shader attribute, eg:
	//   struct LineAttribute { bool isLine; }; // line shader.
	//   [shader("vertex")] [Line(true)]
	DrawLines bool // true to render lines instead of triangles.

	// TODO user defined shader attribute
	CullModeNone bool // true disables backface culling.

	// Attrs layout location match the order of declaration in the shader.
	Attrs []ShaderAttribute

	// Uniforms must match the shader attributes in name and position
	// where scope also identifies a DescriptorSet.
	Uniforms []ShaderUniform
}

// ShaderAttribute identifies the data layouts for attributes.
type ShaderAttribute struct {
	Name   string         // unique name matching shader source code.
	Scope  AttributeScope // vertex or instanced
	AType  int            // attribute type
	DType  ShaderDataType // attribute data type.
	Stride uint32         // size in bytes
}

// ShaderUniform identifies the data layouts for uniforms.
type ShaderUniform struct {
	Name      string         // unique name matching shader source code.
	Scope     UniformScope   //
	UType     ShaderDataType //
	SceneUID  SceneUniform   // index to scene uniform data
	ModelUID  ModelUniform   // index to model uniform data
	Size      int            // size in bytes
	IsSampler bool           // true for sampler uniforms.
}

// GetSceneUniforms returns all the scene scoped uniforms from the
// Shader data.
func (s *Shader) GetSceneUniforms() (uniforms []*ShaderUniform) {
	for i, uni := range s.Uniforms {
		if uni.UType != Type_SAMPLER && uni.Scope == SceneScope {
			uniforms = append(uniforms, &s.Uniforms[i])
		}
	}
	return uniforms
}

// GetModelUniforms returns all the model scoped uniforms from the
// Shader data.
func (s *Shader) GetModelUniforms() (uniforms []*ShaderUniform) {
	for i, uni := range s.Uniforms {
		if uni.UType != Type_SAMPLER && uni.Scope == ModelScope {
			uniforms = append(uniforms, &s.Uniforms[i])
		}
	}
	return uniforms
}

// GetSceneUniforms returns the sampler uniforms from the Shader data.
func (s *Shader) GetSamplerUniforms() (uniforms []*ShaderUniform) {
	for i, uni := range s.Uniforms {
		if uni.UType == Type_SAMPLER && uni.Scope == MaterialScope {
			uniforms = append(uniforms, &s.Uniforms[i])
		}
	}
	return uniforms
}

// =============================================================================
// AttributeScope identifies the supported shader attribute types.
type AttributeScope uint8

const (
	VertexAttribute   AttributeScope = iota // per vertex data
	InstanceAttribute                       // per instance data.
)

// =============================================================================
// ShaderPass identifies a shader render pass.
type ShaderPass uint8

const (
	Renderpass_3D ShaderPass = iota // 3D is rendered before 2D
	Renderpass_2D                   //
)

// =============================================================================
// ShaderStage identifies the currently supported programmable
// shader stages that a shader pipeline can have.
type ShaderStage uint8

const (
	Stage_VERTEX   ShaderStage = 1 << iota // vertex processing.
	Stage_GEOMETRY                         // eg: turn points into quads.
	Stage_FRAGMENT                         // pixel processing.
)

// =============================================================================
// UniformScope identifies the three supported types of uniforms.
// The scope is directly related to the layout(set=x) value in the shader code.
type UniformScope uint8

const (
	SceneScope    UniformScope = iota // descriptor set=0
	MaterialScope                     // descriptor set=1  texture samplers
	ModelScope                        // descriptor set=2  FUTURE: model > 128b
	PushScope                         // push constants    model: 128 byte limit
)

// =============================================================================
// SceneUniform is scene scope uniform data.
type SceneUniform uint8

const (
	PROJ          SceneUniform = iota // scene
	VIEW                              // scene
	NMAT                              // scene
	CAM                               // scene
	LIGHTS                            // scene
	LIGHTCNT                          // scene
	TIME                              // scene
	SceneUniforms                     // must be last
)

// =============================================================================
// ModelUniform is model scope uniform data.
type ModelUniform uint8

const (
	MODEL         ModelUniform = iota // model transform
	SCALE                             // vec3
	PBRMR                             // vec4 x:metallic y:roughness
	COLOR                             // rgba
	F4                                // shader specific parameters.
	F16                               // shader specific parameters.
	ModelUniforms                     // must be last
)

// =============================================================================
// map the shader configuration to data structures.
// Some maps are public so external apps can add new shader attributes
// as needed by new shaders.
//
// FUTURE: use shader reflection to generate the shader descriptions
// and maps from the shader code.
var shaderStages = map[string]ShaderStage{
	"vertex":   Stage_VERTEX,
	"fragment": Stage_FRAGMENT,
}

// ShaderAttributes are used by the render layer to allocate buffers
// for attribute data. Matches the data type read from GLB files.
var ShaderAttributes = map[string]int{
	// names for model vertex data
	"position": Vertexes,
	"texcoord": Texcoords,
	"normal":   Normals,
	"tangent":  Tangents, // FUTURE
	"joint":    Joints,   // FUTURE animation
	"weight":   Weights,  // FUTURE animation

	// names for instanced model data
	"i_position": InstancePosition,
	"i_color":    InstanceColors,
	"i_scale":    InstanceScales,
}

// ShaderAttributeScope categorizes vertex attributes into
// per-vertex data or per-model-instance data.
// Expected use is for passing data from the engine to the render system.
var ShaderAttributeScope = map[string]AttributeScope{
	"vertex":   VertexAttribute,
	"instance": InstanceAttribute,
}

// ShaderSceneUniforms are shader uniforms that apply to
// a single scene (render pass).
// Used for passing data from the engine to the render system.
var ShaderSceneUniforms = map[string]SceneUniform{
	"proj":     PROJ,     //
	"view":     VIEW,     //
	"cam":      CAM,      //
	"lights":   LIGHTS,   //
	"lightCnt": LIGHTCNT, //
	"time":     TIME,     //
}

// ShaderModelUniforms are shader uniforms that apply to one model.
// Used for passing data from the engine to the render system.
var ShaderModelUniforms = map[string]ModelUniform{
	"model": MODEL, // 4x4 matrix
	"scale": SCALE, // 3 floats
	"color": COLOR, // rgba floats
	"pbrMR": PBRMR, // PBR metallic roughness.
	"f4":    F4,    // 4 floats
	"f16":   F16,   // 16 floats
}

// =============================================================================
// ShaderDataType helps describe shader attributes and uniforms
type ShaderDataType uint8

const (
	Type_INT ShaderDataType = iota
	Type_FLOAT
	Type_LIGHTARRAY // array of lights
	Type_MAT3
	Type_MAT4
	Type_MAT34
	Type_SAMPLER
	Type_VEC2
	Type_VEC3
	Type_VEC4
	Type_IVEC4
)

// ShaderUniformData are the supported uniform data types.
// Expected use is for passing data from the engine to the render system.
var ShaderUniformData = map[string]ShaderDataType{
	"int":     Type_INT,
	"lights":  Type_LIGHTARRAY,
	"sampler": Type_SAMPLER,

	// slang data types
	"float32":  Type_FLOAT,
	"matrix33": Type_MAT3,
	"matrix44": Type_MAT4,
	"matrix34": Type_MAT34,
	"vector2":  Type_VEC2,
	"vector3":  Type_VEC3,
	"vector4":  Type_VEC4,
}

// Shd loads json shader reflection data and returns it as
// a shader configuration struct. This is needed to provide
// a link between the models and the shader programs.
func Shd(name string, data []byte) (shader *Shader, err error) {
	shd := shaderJSON{}
	if err := json.Unmarshal(data, &shd); err != nil {
		return shader, fmt.Errorf("Shd: json unmarshal %w", err)
	}

	// map the shader pipeline stages to a bitmask and complain if
	// the engine doesn't yet support a pipeline stage.
	stages := ShaderStage(0)
	for _, ent := range shd.EntryPoints {
		stage, ok := shaderStages[ent.Stage]
		if !ok {
			return shader, fmt.Errorf("Shd:unsupported stage %s", ent.Stage)
			// ie: new FEATURE: add support for shader stage", "stage", ent.Stage
		}
		stages |= stage
	}

	// map the shader attributes using variable names to imply the attribute data type.
	pass := "3D"
	attrs := []ShaderAttribute{}
	for _, ent := range shd.EntryPoints {
		if ent.Stage != "vertex" {
			continue // only care about vertex shader attributes.
		}
		for _, parm := range ent.Parameters {
			for _, f := range parm.Type.Fields {

				// by convention the attribute variable name implies the data type.
				atype, ok1 := ShaderAttributes[f.Name]
				if !ok1 {
					return shader, fmt.Errorf("Shd:unsupported shader attribute %s", f.Name)
				}

				// instance data attributes names start with "i_"
				ascope := VertexAttribute
				if strings.HasPrefix(f.Name, "i_") {
					ascope = InstanceAttribute
				}

				// get the type of data for this field.
				dtype, stride := ShaderDataType(0), uint32(0)
				switch {
				case f.Type.Kind == "scalar" && f.Type.Element.ScalarType == "float32":
					dtype, stride = Type_FLOAT, 4
				case f.Type.Kind == "scalar" && f.Type.ScalarType == "float32":
					dtype, stride = Type_FLOAT, 4
				case f.Type.Kind == "vector" && f.Type.ElementCount == 2:
					dtype, stride = Type_VEC2, 8
				case f.Type.Kind == "vector" && f.Type.ElementCount == 3:
					dtype, stride = Type_VEC3, 12
				case f.Type.Kind == "vector" && f.Type.ElementCount == 4:
					dtype, stride = Type_VEC4, 16
				default:
					slog.Error("unknown attribute type", "field", f.Name, "kind", f.Type.Kind)
				}
				attrs = append(attrs, ShaderAttribute{
					Name:   f.Name,
					AType:  atype,
					Scope:  ascope,
					DType:  dtype,
					Stride: stride,
				})

				// use the number of vertex elements to determine the shader pass.
				if f.Name == "position" && f.Type.Kind == "vector" && f.Type.ElementCount == 2 {
					pass = "2D"
				}
			}
		}
	}

	// map the shader uniforms
	uniforms := []ShaderUniform{}
	for _, parm := range shd.Parameters {
		switch {
		case parm.Name == "samplers":
			// samplers are a material scope uniform. Color is the only
			// supported texture type right now, ie:
			// AddModel("shd:tint", "msh:icon", "tex:color:undo") <-- "color"
			u := ShaderUniform{
				Name:      "color",
				Scope:     MaterialScope,
				UType:     Type_SAMPLER,
				SceneUID:  PROJ,
				ModelUID:  COLOR,
				IsSampler: true,
			}
			uniforms = append(uniforms, u)
		default:
			for _, f := range parm.Type.Element.Fields {
				sceneUID, okScene := ShaderSceneUniforms[f.Name]
				modelUID, okModel := ShaderModelUniforms[f.Name]
				var scope UniformScope
				switch {
				case okScene && !okModel:
					scope = SceneScope
				case !okScene && okModel && parm.Binding.Kind == "pushConstantBuffer":
					scope = PushScope
				default:
					return shader, fmt.Errorf("Shd:unsupported uniform scope %s", f.Name)
				}

				dstr := f.Type.Element.ScalarType
				kind := f.Type.Kind
				switch {
				case kind == "matrix":
					dstr = kind + strconv.Itoa(f.Type.RowCount) + strconv.Itoa(f.Type.ColumnCount)
				case kind == "vector":
					dstr = kind + strconv.Itoa(f.Type.ElementCount)
				case kind == "array":
					dstr = "lights"
				default:
					slog.Error("unknown uniform type", "name", f.Name, "kind", kind)
				}
				utype, ok2 := ShaderUniformData[dstr]
				if !ok2 {
					return shader, fmt.Errorf("Shd:unsupported uniform data %s", dstr)
				}

				// uniform expected in the render pass or in the render packet.
				uniforms = append(uniforms, ShaderUniform{
					Name:     f.Name,
					Scope:    scope,
					UType:    utype,
					SceneUID: sceneUID,
					ModelUID: modelUID,
					Size:     f.Binding.Size,
				})
			}
		}
	}

	// create the shader configuration.
	shaderName, _ := strings.CutSuffix(name, ".shd")
	shader = &Shader{
		Name:     shaderName,
		Pass:     pass,
		Stages:   stages,
		Attrs:    attrs,
		Uniforms: uniforms,
	}

	// FUTURE: more render flags as needed.
	for _, ent := range shd.EntryPoints {
		if ent.Stage == "vertex" {
			for _, tag := range ent.Tags {
				if tag.Name == "Line" {
					shader.DrawLines = true
				}
				if tag.Name == "CullOff" {
					shader.CullModeNone = true
				}
			}
		}
	}
	// fmt.Printf("json shader %+v\n", shader)
	return shader, nil

}

// ==================================================================================================
// shaderJSON are structs used to parse the SPRIV json reflection files generated
// at the same time the shader *.spv files are created.
//
// Parsing structs are from a json-to-go converter on the reflection output.
// Only mapped the elements of interest.
type shaderJSON struct {
	EntryPoints []shaderEntryPoint `json:"entryPoints"`
	Parameters  []shaderParameter  `json:"parameters"`
}

// shader entry points (vertex attributes).
//
//	 "entryPoints": [
//	 {
//		"name": "vertMain",
//		"stage": "vertex",
//		"parameters": [
//		    {
//		        "name": "stage_input",
//		        "stage": "vertex",
//		        "binding": {"kind": "varyingInput", "index": 0, "count": 2},
//		        "type": {
//		            "kind": "struct",
//		            "name": "vertInput",
//		            "fields": [
//		    ...
//		"userAttribs": [
//		...
type shaderEntryPoint struct {
	Name       string            `json:"name"`
	Stage      string            `json:"stage"`
	Parameters []shaderEntryParm `json:"parameters"`
	Tags       []shaderTag       `json:"userAttribs"`
}
type shaderEntryParm struct {
	Name    string          `json:"name"`
	Stage   string          `json:"stage"`
	Binding shaderBinding   `json:"binding"`
	Type    shaderEntryType `json:"type"`
}
type shaderEntryType struct {
	Kind   string        `json:"kind"`
	Name   string        `json:"name"`
	Fields []shaderField `json:"fields"`
}

// shader parameters (uniforms).
//
//	 "parameters": [
//	 {
//		"name": "scene_uniforms",
//		"binding": {"kind": "descriptorTableSlot", "index": 0},
//		"type": {
//		    "kind": "constantBuffer",
//		    "elementType": {
//		        "kind": "struct",
//		        "fields": [
//		...
type shaderParameter struct {
	Name    string         `json:"name"`
	Binding shaderBinding  `json:"binding"`
	Type    shaderParmType `json:"type,omitempty"` //
}
type shaderParmType struct {
	Kind    string            `json:"kind"`
	Element shaderElementType `json:"elementType"`
}
type shaderElementType struct {
	Kind       string         `json:"kind"`
	Fields     []shaderField  `json:"fields"`
	BaseShape  string         `json:"baseShape"`
	ResultType shaderBaseType `json:"resultType"`
}

// shader fields are common to uniforms and attributes
//
//	 "fields": [
//	 {
//		"name": "position",
//		"type": {
//		    "kind": "vector",
//		    "elementCount": 2,
//		    "elementType": {
//		        "kind": "scalar",
//		        "scalarType": "float32"
//		    }
//		},
//		"stage": "vertex",
//		"binding": {"kind": "varyingInput", "index": 0}
//	 },
type shaderField struct {
	Name    string         `json:"name"`
	Type    shaderBaseType `json:"type"`
	Stage   string         `json:"stage"`
	Binding shaderBinding  `json:"binding"`
}
type shaderBaseType struct {
	Kind         string       `json:"kind"`
	RowCount     int          `json:"rowCount,omitempty"`
	ColumnCount  int          `json:"columnCount,omitempty"`
	ElementCount int          `json:"elementCount,omitempty"`
	Element      shaderScalar `json:"elementType"`
	ScalarType   string       `json:"scalarType"`
}
type shaderScalar struct {
	Kind       string `json:"kind"`
	ScalarType string `json:"scalarType"`
}
type shaderBinding struct {
	Kind   string `json:"kind"`   // all bindings
	Index  int    `json:"index"`  // kind: descriptorTableSlot, varyingInput, pushConstantBuffer,
	Offset int    `json:"offset"` // kind: uniform
	Size   int    `json:"size"`   // kind: uniform
	Count  int    `json:"count"`  // kind: varyingInput
	Space  int    `json:"space"`  // kind: descriptorTableSlot
}

// user attributes are shader pipeline configuration tags.
// The presence of a user attribute sets a pipeline config bool to true.
//
//	"userAttribs": [
//	{
//		"name": "Line", // ie: shader expects line vertex data instead of triangles
//		"arguments": [  //     args are ignored.
//		    1
//		]
//	}
type shaderTag struct {
	Name string `json:"name"` // attribute name. "value" can be any type.
}
