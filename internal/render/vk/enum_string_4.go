// Code generated by "stringer -output=enum_string_4.go -type=PolygonMode,PresentModeKHR,PrimitiveTopology,ProvokingVertexModeEXT,QueryControlFlagBits,QueryPipelineStatisticFlagBits,QueryPoolSamplingModeINTEL,QueryResultFlagBits,QueryType,QueueFlagBits,QueueGlobalPriorityEXT,RasterizationOrderAMD,RayTracingShaderGroupTypeKHR,RenderPassCreateFlagBits,RenderingFlagBitsKHR,ResolveModeFlagBits,Result,SampleCountFlagBits,SamplerAddressMode,SamplerCreateFlagBits,SamplerMipmapMode,SamplerReductionMode,SamplerYcbcrModelConversion,SamplerYcbcrRange,ScopeNV,SemaphoreImportFlagBits,SemaphoreType,SemaphoreWaitFlagBits,ShaderFloatControlsIndependence,ShaderGroupShaderKHR,ShaderInfoTypeAMD"; DO NOT EDIT.

package vk

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[POLYGON_MODE_FILL-0]
	_ = x[POLYGON_MODE_LINE-1]
	_ = x[POLYGON_MODE_POINT-2]
	_ = x[POLYGON_MODE_FILL_RECTANGLE_NV-1000153000]
}

const (
	_PolygonMode_name_0 = "POLYGON_MODE_FILLPOLYGON_MODE_LINEPOLYGON_MODE_POINT"
	_PolygonMode_name_1 = "POLYGON_MODE_FILL_RECTANGLE_NV"
)

var (
	_PolygonMode_index_0 = [...]uint8{0, 17, 34, 52}
)

func (i PolygonMode) String() string {
	switch {
	case 0 <= i && i <= 2:
		return _PolygonMode_name_0[_PolygonMode_index_0[i]:_PolygonMode_index_0[i+1]]
	case i == 1000153000:
		return _PolygonMode_name_1
	default:
		return "PolygonMode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[PRESENT_MODE_IMMEDIATE_KHR-0]
	_ = x[PRESENT_MODE_MAILBOX_KHR-1]
	_ = x[PRESENT_MODE_FIFO_KHR-2]
	_ = x[PRESENT_MODE_FIFO_RELAXED_KHR-3]
	_ = x[PRESENT_MODE_SHARED_DEMAND_REFRESH_KHR-1000111000]
	_ = x[PRESENT_MODE_SHARED_CONTINUOUS_REFRESH_KHR-1000111001]
}

const (
	_PresentModeKHR_name_0 = "PRESENT_MODE_IMMEDIATE_KHRPRESENT_MODE_MAILBOX_KHRPRESENT_MODE_FIFO_KHRPRESENT_MODE_FIFO_RELAXED_KHR"
	_PresentModeKHR_name_1 = "PRESENT_MODE_SHARED_DEMAND_REFRESH_KHRPRESENT_MODE_SHARED_CONTINUOUS_REFRESH_KHR"
)

var (
	_PresentModeKHR_index_0 = [...]uint8{0, 26, 50, 71, 100}
	_PresentModeKHR_index_1 = [...]uint8{0, 38, 80}
)

func (i PresentModeKHR) String() string {
	switch {
	case 0 <= i && i <= 3:
		return _PresentModeKHR_name_0[_PresentModeKHR_index_0[i]:_PresentModeKHR_index_0[i+1]]
	case 1000111000 <= i && i <= 1000111001:
		i -= 1000111000
		return _PresentModeKHR_name_1[_PresentModeKHR_index_1[i]:_PresentModeKHR_index_1[i+1]]
	default:
		return "PresentModeKHR(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[PRIMITIVE_TOPOLOGY_POINT_LIST-0]
	_ = x[PRIMITIVE_TOPOLOGY_LINE_LIST-1]
	_ = x[PRIMITIVE_TOPOLOGY_LINE_STRIP-2]
	_ = x[PRIMITIVE_TOPOLOGY_TRIANGLE_LIST-3]
	_ = x[PRIMITIVE_TOPOLOGY_TRIANGLE_STRIP-4]
	_ = x[PRIMITIVE_TOPOLOGY_TRIANGLE_FAN-5]
	_ = x[PRIMITIVE_TOPOLOGY_LINE_LIST_WITH_ADJACENCY-6]
	_ = x[PRIMITIVE_TOPOLOGY_LINE_STRIP_WITH_ADJACENCY-7]
	_ = x[PRIMITIVE_TOPOLOGY_TRIANGLE_LIST_WITH_ADJACENCY-8]
	_ = x[PRIMITIVE_TOPOLOGY_TRIANGLE_STRIP_WITH_ADJACENCY-9]
	_ = x[PRIMITIVE_TOPOLOGY_PATCH_LIST-10]
}

const _PrimitiveTopology_name = "PRIMITIVE_TOPOLOGY_POINT_LISTPRIMITIVE_TOPOLOGY_LINE_LISTPRIMITIVE_TOPOLOGY_LINE_STRIPPRIMITIVE_TOPOLOGY_TRIANGLE_LISTPRIMITIVE_TOPOLOGY_TRIANGLE_STRIPPRIMITIVE_TOPOLOGY_TRIANGLE_FANPRIMITIVE_TOPOLOGY_LINE_LIST_WITH_ADJACENCYPRIMITIVE_TOPOLOGY_LINE_STRIP_WITH_ADJACENCYPRIMITIVE_TOPOLOGY_TRIANGLE_LIST_WITH_ADJACENCYPRIMITIVE_TOPOLOGY_TRIANGLE_STRIP_WITH_ADJACENCYPRIMITIVE_TOPOLOGY_PATCH_LIST"

var _PrimitiveTopology_index = [...]uint16{0, 29, 57, 86, 118, 151, 182, 225, 269, 316, 364, 393}

func (i PrimitiveTopology) String() string {
	if i < 0 || i >= PrimitiveTopology(len(_PrimitiveTopology_index)-1) {
		return "PrimitiveTopology(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _PrimitiveTopology_name[_PrimitiveTopology_index[i]:_PrimitiveTopology_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[PROVOKING_VERTEX_MODE_FIRST_VERTEX_EXT-0]
	_ = x[PROVOKING_VERTEX_MODE_LAST_VERTEX_EXT-1]
}

const _ProvokingVertexModeEXT_name = "PROVOKING_VERTEX_MODE_FIRST_VERTEX_EXTPROVOKING_VERTEX_MODE_LAST_VERTEX_EXT"

var _ProvokingVertexModeEXT_index = [...]uint8{0, 38, 75}

func (i ProvokingVertexModeEXT) String() string {
	if i < 0 || i >= ProvokingVertexModeEXT(len(_ProvokingVertexModeEXT_index)-1) {
		return "ProvokingVertexModeEXT(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ProvokingVertexModeEXT_name[_ProvokingVertexModeEXT_index[i]:_ProvokingVertexModeEXT_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[QUERY_CONTROL_PRECISE_BIT-1]
}

const _QueryControlFlagBits_name = "QUERY_CONTROL_PRECISE_BIT"

var _QueryControlFlagBits_index = [...]uint8{0, 25}

func (i QueryControlFlagBits) String() string {
	i -= 1
	if i >= QueryControlFlagBits(len(_QueryControlFlagBits_index)-1) {
		return "QueryControlFlagBits(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _QueryControlFlagBits_name[_QueryControlFlagBits_index[i]:_QueryControlFlagBits_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[QUERY_PIPELINE_STATISTIC_INPUT_ASSEMBLY_VERTICES_BIT-1]
	_ = x[QUERY_PIPELINE_STATISTIC_INPUT_ASSEMBLY_PRIMITIVES_BIT-2]
	_ = x[QUERY_PIPELINE_STATISTIC_COMPUTE_SHADER_INVOCATIONS_BIT-1024]
	_ = x[QUERY_PIPELINE_STATISTIC_VERTEX_SHADER_INVOCATIONS_BIT-4]
	_ = x[QUERY_PIPELINE_STATISTIC_GEOMETRY_SHADER_INVOCATIONS_BIT-8]
	_ = x[QUERY_PIPELINE_STATISTIC_GEOMETRY_SHADER_PRIMITIVES_BIT-16]
	_ = x[QUERY_PIPELINE_STATISTIC_CLIPPING_INVOCATIONS_BIT-32]
	_ = x[QUERY_PIPELINE_STATISTIC_CLIPPING_PRIMITIVES_BIT-64]
	_ = x[QUERY_PIPELINE_STATISTIC_FRAGMENT_SHADER_INVOCATIONS_BIT-128]
	_ = x[QUERY_PIPELINE_STATISTIC_TESSELLATION_CONTROL_SHADER_PATCHES_BIT-256]
	_ = x[QUERY_PIPELINE_STATISTIC_TESSELLATION_EVALUATION_SHADER_INVOCATIONS_BIT-512]
}

const (
	_QueryPipelineStatisticFlagBits_name_0 = "QUERY_PIPELINE_STATISTIC_INPUT_ASSEMBLY_VERTICES_BITQUERY_PIPELINE_STATISTIC_INPUT_ASSEMBLY_PRIMITIVES_BIT"
	_QueryPipelineStatisticFlagBits_name_1 = "QUERY_PIPELINE_STATISTIC_VERTEX_SHADER_INVOCATIONS_BIT"
	_QueryPipelineStatisticFlagBits_name_2 = "QUERY_PIPELINE_STATISTIC_GEOMETRY_SHADER_INVOCATIONS_BIT"
	_QueryPipelineStatisticFlagBits_name_3 = "QUERY_PIPELINE_STATISTIC_GEOMETRY_SHADER_PRIMITIVES_BIT"
	_QueryPipelineStatisticFlagBits_name_4 = "QUERY_PIPELINE_STATISTIC_CLIPPING_INVOCATIONS_BIT"
	_QueryPipelineStatisticFlagBits_name_5 = "QUERY_PIPELINE_STATISTIC_CLIPPING_PRIMITIVES_BIT"
	_QueryPipelineStatisticFlagBits_name_6 = "QUERY_PIPELINE_STATISTIC_FRAGMENT_SHADER_INVOCATIONS_BIT"
	_QueryPipelineStatisticFlagBits_name_7 = "QUERY_PIPELINE_STATISTIC_TESSELLATION_CONTROL_SHADER_PATCHES_BIT"
	_QueryPipelineStatisticFlagBits_name_8 = "QUERY_PIPELINE_STATISTIC_TESSELLATION_EVALUATION_SHADER_INVOCATIONS_BIT"
	_QueryPipelineStatisticFlagBits_name_9 = "QUERY_PIPELINE_STATISTIC_COMPUTE_SHADER_INVOCATIONS_BIT"
)

var (
	_QueryPipelineStatisticFlagBits_index_0 = [...]uint8{0, 52, 106}
)

func (i QueryPipelineStatisticFlagBits) String() string {
	switch {
	case 1 <= i && i <= 2:
		i -= 1
		return _QueryPipelineStatisticFlagBits_name_0[_QueryPipelineStatisticFlagBits_index_0[i]:_QueryPipelineStatisticFlagBits_index_0[i+1]]
	case i == 4:
		return _QueryPipelineStatisticFlagBits_name_1
	case i == 8:
		return _QueryPipelineStatisticFlagBits_name_2
	case i == 16:
		return _QueryPipelineStatisticFlagBits_name_3
	case i == 32:
		return _QueryPipelineStatisticFlagBits_name_4
	case i == 64:
		return _QueryPipelineStatisticFlagBits_name_5
	case i == 128:
		return _QueryPipelineStatisticFlagBits_name_6
	case i == 256:
		return _QueryPipelineStatisticFlagBits_name_7
	case i == 512:
		return _QueryPipelineStatisticFlagBits_name_8
	case i == 1024:
		return _QueryPipelineStatisticFlagBits_name_9
	default:
		return "QueryPipelineStatisticFlagBits(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[QUERY_POOL_SAMPLING_MODE_MANUAL_INTEL-0]
}

const _QueryPoolSamplingModeINTEL_name = "QUERY_POOL_SAMPLING_MODE_MANUAL_INTEL"

var _QueryPoolSamplingModeINTEL_index = [...]uint8{0, 37}

func (i QueryPoolSamplingModeINTEL) String() string {
	if i < 0 || i >= QueryPoolSamplingModeINTEL(len(_QueryPoolSamplingModeINTEL_index)-1) {
		return "QueryPoolSamplingModeINTEL(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _QueryPoolSamplingModeINTEL_name[_QueryPoolSamplingModeINTEL_index[i]:_QueryPoolSamplingModeINTEL_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[QUERY_RESULT_64_BIT-1]
	_ = x[QUERY_RESULT_WAIT_BIT-2]
	_ = x[QUERY_RESULT_WITH_AVAILABILITY_BIT-4]
	_ = x[QUERY_RESULT_PARTIAL_BIT-8]
}

const (
	_QueryResultFlagBits_name_0 = "QUERY_RESULT_64_BITQUERY_RESULT_WAIT_BIT"
	_QueryResultFlagBits_name_1 = "QUERY_RESULT_WITH_AVAILABILITY_BIT"
	_QueryResultFlagBits_name_2 = "QUERY_RESULT_PARTIAL_BIT"
)

var (
	_QueryResultFlagBits_index_0 = [...]uint8{0, 19, 40}
)

func (i QueryResultFlagBits) String() string {
	switch {
	case 1 <= i && i <= 2:
		i -= 1
		return _QueryResultFlagBits_name_0[_QueryResultFlagBits_index_0[i]:_QueryResultFlagBits_index_0[i+1]]
	case i == 4:
		return _QueryResultFlagBits_name_1
	case i == 8:
		return _QueryResultFlagBits_name_2
	default:
		return "QueryResultFlagBits(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[QUERY_TYPE_OCCLUSION-0]
	_ = x[QUERY_TYPE_PIPELINE_STATISTICS-1]
	_ = x[QUERY_TYPE_TIMESTAMP-2]
	_ = x[QUERY_TYPE_TRANSFORM_FEEDBACK_STREAM_EXT-1000028004]
	_ = x[QUERY_TYPE_PERFORMANCE_QUERY_KHR-1000116000]
	_ = x[QUERY_TYPE_ACCELERATION_STRUCTURE_COMPACTED_SIZE_KHR-1000150000]
	_ = x[QUERY_TYPE_ACCELERATION_STRUCTURE_SERIALIZATION_SIZE_KHR-1000150001]
	_ = x[QUERY_TYPE_ACCELERATION_STRUCTURE_COMPACTED_SIZE_NV-1000165000]
	_ = x[QUERY_TYPE_PERFORMANCE_QUERY_INTEL-1000210000]
}

const (
	_QueryType_name_0 = "QUERY_TYPE_OCCLUSIONQUERY_TYPE_PIPELINE_STATISTICSQUERY_TYPE_TIMESTAMP"
	_QueryType_name_1 = "QUERY_TYPE_TRANSFORM_FEEDBACK_STREAM_EXT"
	_QueryType_name_2 = "QUERY_TYPE_PERFORMANCE_QUERY_KHR"
	_QueryType_name_3 = "QUERY_TYPE_ACCELERATION_STRUCTURE_COMPACTED_SIZE_KHRQUERY_TYPE_ACCELERATION_STRUCTURE_SERIALIZATION_SIZE_KHR"
	_QueryType_name_4 = "QUERY_TYPE_ACCELERATION_STRUCTURE_COMPACTED_SIZE_NV"
	_QueryType_name_5 = "QUERY_TYPE_PERFORMANCE_QUERY_INTEL"
)

var (
	_QueryType_index_0 = [...]uint8{0, 20, 50, 70}
	_QueryType_index_3 = [...]uint8{0, 52, 108}
)

func (i QueryType) String() string {
	switch {
	case 0 <= i && i <= 2:
		return _QueryType_name_0[_QueryType_index_0[i]:_QueryType_index_0[i+1]]
	case i == 1000028004:
		return _QueryType_name_1
	case i == 1000116000:
		return _QueryType_name_2
	case 1000150000 <= i && i <= 1000150001:
		i -= 1000150000
		return _QueryType_name_3[_QueryType_index_3[i]:_QueryType_index_3[i+1]]
	case i == 1000165000:
		return _QueryType_name_4
	case i == 1000210000:
		return _QueryType_name_5
	default:
		return "QueryType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[QUEUE_GRAPHICS_BIT-1]
	_ = x[QUEUE_COMPUTE_BIT-2]
	_ = x[QUEUE_TRANSFER_BIT-4]
	_ = x[QUEUE_SPARSE_BINDING_BIT-8]
	_ = x[QUEUE_PROTECTED_BIT-16]
}

const (
	_QueueFlagBits_name_0 = "QUEUE_GRAPHICS_BITQUEUE_COMPUTE_BIT"
	_QueueFlagBits_name_1 = "QUEUE_TRANSFER_BIT"
	_QueueFlagBits_name_2 = "QUEUE_SPARSE_BINDING_BIT"
	_QueueFlagBits_name_3 = "QUEUE_PROTECTED_BIT"
)

var (
	_QueueFlagBits_index_0 = [...]uint8{0, 18, 35}
)

func (i QueueFlagBits) String() string {
	switch {
	case 1 <= i && i <= 2:
		i -= 1
		return _QueueFlagBits_name_0[_QueueFlagBits_index_0[i]:_QueueFlagBits_index_0[i+1]]
	case i == 4:
		return _QueueFlagBits_name_1
	case i == 8:
		return _QueueFlagBits_name_2
	case i == 16:
		return _QueueFlagBits_name_3
	default:
		return "QueueFlagBits(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[QUEUE_GLOBAL_PRIORITY_LOW_EXT-128]
	_ = x[QUEUE_GLOBAL_PRIORITY_MEDIUM_EXT-256]
	_ = x[QUEUE_GLOBAL_PRIORITY_HIGH_EXT-512]
	_ = x[QUEUE_GLOBAL_PRIORITY_REALTIME_EXT-1024]
}

const (
	_QueueGlobalPriorityEXT_name_0 = "QUEUE_GLOBAL_PRIORITY_LOW_EXT"
	_QueueGlobalPriorityEXT_name_1 = "QUEUE_GLOBAL_PRIORITY_MEDIUM_EXT"
	_QueueGlobalPriorityEXT_name_2 = "QUEUE_GLOBAL_PRIORITY_HIGH_EXT"
	_QueueGlobalPriorityEXT_name_3 = "QUEUE_GLOBAL_PRIORITY_REALTIME_EXT"
)

func (i QueueGlobalPriorityEXT) String() string {
	switch {
	case i == 128:
		return _QueueGlobalPriorityEXT_name_0
	case i == 256:
		return _QueueGlobalPriorityEXT_name_1
	case i == 512:
		return _QueueGlobalPriorityEXT_name_2
	case i == 1024:
		return _QueueGlobalPriorityEXT_name_3
	default:
		return "QueueGlobalPriorityEXT(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[RASTERIZATION_ORDER_STRICT_AMD-0]
	_ = x[RASTERIZATION_ORDER_RELAXED_AMD-1]
}

const _RasterizationOrderAMD_name = "RASTERIZATION_ORDER_STRICT_AMDRASTERIZATION_ORDER_RELAXED_AMD"

var _RasterizationOrderAMD_index = [...]uint8{0, 30, 61}

func (i RasterizationOrderAMD) String() string {
	if i < 0 || i >= RasterizationOrderAMD(len(_RasterizationOrderAMD_index)-1) {
		return "RasterizationOrderAMD(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _RasterizationOrderAMD_name[_RasterizationOrderAMD_index[i]:_RasterizationOrderAMD_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[RAY_TRACING_SHADER_GROUP_TYPE_GENERAL_KHR-0]
	_ = x[RAY_TRACING_SHADER_GROUP_TYPE_TRIANGLES_HIT_GROUP_KHR-1]
	_ = x[RAY_TRACING_SHADER_GROUP_TYPE_PROCEDURAL_HIT_GROUP_KHR-2]
	_ = x[RAY_TRACING_SHADER_GROUP_TYPE_GENERAL_NV-0]
	_ = x[RAY_TRACING_SHADER_GROUP_TYPE_PROCEDURAL_HIT_GROUP_NV-2]
	_ = x[RAY_TRACING_SHADER_GROUP_TYPE_TRIANGLES_HIT_GROUP_NV-1]
}

const _RayTracingShaderGroupTypeKHR_name = "RAY_TRACING_SHADER_GROUP_TYPE_GENERAL_KHRRAY_TRACING_SHADER_GROUP_TYPE_TRIANGLES_HIT_GROUP_KHRRAY_TRACING_SHADER_GROUP_TYPE_PROCEDURAL_HIT_GROUP_KHR"

var _RayTracingShaderGroupTypeKHR_index = [...]uint8{0, 41, 94, 148}

func (i RayTracingShaderGroupTypeKHR) String() string {
	if i < 0 || i >= RayTracingShaderGroupTypeKHR(len(_RayTracingShaderGroupTypeKHR_index)-1) {
		return "RayTracingShaderGroupTypeKHR(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _RayTracingShaderGroupTypeKHR_name[_RayTracingShaderGroupTypeKHR_index[i]:_RayTracingShaderGroupTypeKHR_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[RENDER_PASS_CREATE_TRANSFORM_BIT_QCOM-1000282000]
}

const _RenderPassCreateFlagBits_name = "RENDER_PASS_CREATE_TRANSFORM_BIT_QCOM"

var _RenderPassCreateFlagBits_index = [...]uint8{0, 37}

func (i RenderPassCreateFlagBits) String() string {
	i -= 1000282000
	if i >= RenderPassCreateFlagBits(len(_RenderPassCreateFlagBits_index)-1) {
		return "RenderPassCreateFlagBits(" + strconv.FormatInt(int64(i+1000282000), 10) + ")"
	}
	return _RenderPassCreateFlagBits_name[_RenderPassCreateFlagBits_index[i]:_RenderPassCreateFlagBits_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[RENDERING_CONTENTS_SECONDARY_COMMAND_BUFFERS_BIT_KHR-1]
	_ = x[RENDERING_SUSPENDING_BIT_KHR-2]
	_ = x[RENDERING_RESUMING_BIT_KHR-4]
}

const (
	_RenderingFlagBitsKHR_name_0 = "RENDERING_CONTENTS_SECONDARY_COMMAND_BUFFERS_BIT_KHRRENDERING_SUSPENDING_BIT_KHR"
	_RenderingFlagBitsKHR_name_1 = "RENDERING_RESUMING_BIT_KHR"
)

var (
	_RenderingFlagBitsKHR_index_0 = [...]uint8{0, 52, 80}
)

func (i RenderingFlagBitsKHR) String() string {
	switch {
	case 1 <= i && i <= 2:
		i -= 1
		return _RenderingFlagBitsKHR_name_0[_RenderingFlagBitsKHR_index_0[i]:_RenderingFlagBitsKHR_index_0[i+1]]
	case i == 4:
		return _RenderingFlagBitsKHR_name_1
	default:
		return "RenderingFlagBitsKHR(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[RESOLVE_MODE_NONE-0]
	_ = x[RESOLVE_MODE_SAMPLE_ZERO_BIT-1]
	_ = x[RESOLVE_MODE_AVERAGE_BIT-2]
	_ = x[RESOLVE_MODE_MIN_BIT-4]
	_ = x[RESOLVE_MODE_MAX_BIT-8]
	_ = x[RESOLVE_MODE_AVERAGE_BIT_KHR-2]
	_ = x[RESOLVE_MODE_MAX_BIT_KHR-8]
	_ = x[RESOLVE_MODE_MIN_BIT_KHR-4]
	_ = x[RESOLVE_MODE_NONE_KHR-0]
	_ = x[RESOLVE_MODE_SAMPLE_ZERO_BIT_KHR-1]
}

const (
	_ResolveModeFlagBits_name_0 = "RESOLVE_MODE_NONERESOLVE_MODE_SAMPLE_ZERO_BITRESOLVE_MODE_AVERAGE_BIT"
	_ResolveModeFlagBits_name_1 = "RESOLVE_MODE_MIN_BIT"
	_ResolveModeFlagBits_name_2 = "RESOLVE_MODE_MAX_BIT"
)

var (
	_ResolveModeFlagBits_index_0 = [...]uint8{0, 17, 45, 69}
)

func (i ResolveModeFlagBits) String() string {
	switch {
	case i <= 2:
		return _ResolveModeFlagBits_name_0[_ResolveModeFlagBits_index_0[i]:_ResolveModeFlagBits_index_0[i+1]]
	case i == 4:
		return _ResolveModeFlagBits_name_1
	case i == 8:
		return _ResolveModeFlagBits_name_2
	default:
		return "ResolveModeFlagBits(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ERROR_INVALID_OPAQUE_CAPTURE_ADDRESS - -1000257000]
	_ = x[ERROR_NOT_PERMITTED_EXT - -1000174001]
	_ = x[ERROR_FRAGMENTATION - -1000161000]
	_ = x[ERROR_INVALID_DRM_FORMAT_MODIFIER_PLANE_LAYOUT_EXT - -1000158000]
	_ = x[ERROR_INVALID_EXTERNAL_HANDLE - -1000072003]
	_ = x[ERROR_OUT_OF_POOL_MEMORY - -1000069000]
	_ = x[ERROR_INVALID_SHADER_NV - -1000012000]
	_ = x[ERROR_VALIDATION_FAILED_EXT - -1000011001]
	_ = x[ERROR_INCOMPATIBLE_DISPLAY_KHR - -1000003001]
	_ = x[ERROR_OUT_OF_DATE_KHR - -1000001004]
	_ = x[ERROR_NATIVE_WINDOW_IN_USE_KHR - -1000000001]
	_ = x[ERROR_SURFACE_LOST_KHR - -1000000000]
	_ = x[ERROR_UNKNOWN - -13]
	_ = x[ERROR_FRAGMENTED_POOL - -12]
	_ = x[ERROR_FORMAT_NOT_SUPPORTED - -11]
	_ = x[ERROR_TOO_MANY_OBJECTS - -10]
	_ = x[ERROR_INCOMPATIBLE_DRIVER - -9]
	_ = x[ERROR_FEATURE_NOT_PRESENT - -8]
	_ = x[ERROR_EXTENSION_NOT_PRESENT - -7]
	_ = x[ERROR_LAYER_NOT_PRESENT - -6]
	_ = x[ERROR_MEMORY_MAP_FAILED - -5]
	_ = x[ERROR_DEVICE_LOST - -4]
	_ = x[ERROR_INITIALIZATION_FAILED - -3]
	_ = x[ERROR_OUT_OF_DEVICE_MEMORY - -2]
	_ = x[ERROR_OUT_OF_HOST_MEMORY - -1]
	_ = x[NOT_READY-1]
	_ = x[TIMEOUT-2]
	_ = x[EVENT_SET-3]
	_ = x[EVENT_RESET-4]
	_ = x[INCOMPLETE-5]
	_ = x[SUBOPTIMAL_KHR-1000001003]
	_ = x[THREAD_IDLE_KHR-1000268000]
	_ = x[THREAD_DONE_KHR-1000268001]
	_ = x[OPERATION_DEFERRED_KHR-1000268002]
	_ = x[OPERATION_NOT_DEFERRED_KHR-1000268003]
	_ = x[PIPELINE_COMPILE_REQUIRED_EXT-1000297000]
	_ = x[ERROR_FRAGMENTATION_EXT - -1000161000]
	_ = x[ERROR_INVALID_EXTERNAL_HANDLE_KHR - -1000072003]
	_ = x[ERROR_INVALID_OPAQUE_CAPTURE_ADDRESS_KHR - -1000257000]
	_ = x[ERROR_INVALID_DEVICE_ADDRESS_EXT - -1000257000]
	_ = x[ERROR_OUT_OF_POOL_MEMORY_KHR - -1000069000]
	_ = x[ERROR_PIPELINE_COMPILE_REQUIRED_EXT-1000297000]
	_ = x[ERROR_FULL_SCREEN_EXCLUSIVE_MODE_LOST_EXT - -1000255000]
}

const _Result_name = "ERROR_INVALID_OPAQUE_CAPTURE_ADDRESSERROR_FULL_SCREEN_EXCLUSIVE_MODE_LOST_EXTERROR_NOT_PERMITTED_EXTERROR_FRAGMENTATIONERROR_INVALID_DRM_FORMAT_MODIFIER_PLANE_LAYOUT_EXTERROR_INVALID_EXTERNAL_HANDLEERROR_OUT_OF_POOL_MEMORYERROR_INVALID_SHADER_NVERROR_VALIDATION_FAILED_EXTERROR_INCOMPATIBLE_DISPLAY_KHRERROR_OUT_OF_DATE_KHRERROR_NATIVE_WINDOW_IN_USE_KHRERROR_SURFACE_LOST_KHRERROR_UNKNOWNERROR_FRAGMENTED_POOLERROR_FORMAT_NOT_SUPPORTEDERROR_TOO_MANY_OBJECTSERROR_INCOMPATIBLE_DRIVERERROR_FEATURE_NOT_PRESENTERROR_EXTENSION_NOT_PRESENTERROR_LAYER_NOT_PRESENTERROR_MEMORY_MAP_FAILEDERROR_DEVICE_LOSTERROR_INITIALIZATION_FAILEDERROR_OUT_OF_DEVICE_MEMORYERROR_OUT_OF_HOST_MEMORYNOT_READYTIMEOUTEVENT_SETEVENT_RESETINCOMPLETESUBOPTIMAL_KHRTHREAD_IDLE_KHRTHREAD_DONE_KHROPERATION_DEFERRED_KHROPERATION_NOT_DEFERRED_KHRPIPELINE_COMPILE_REQUIRED_EXT"

var _Result_map = map[Result]string{
	-1000257000: _Result_name[0:36],
	-1000255000: _Result_name[36:77],
	-1000174001: _Result_name[77:100],
	-1000161000: _Result_name[100:119],
	-1000158000: _Result_name[119:169],
	-1000072003: _Result_name[169:198],
	-1000069000: _Result_name[198:222],
	-1000012000: _Result_name[222:245],
	-1000011001: _Result_name[245:272],
	-1000003001: _Result_name[272:302],
	-1000001004: _Result_name[302:323],
	-1000000001: _Result_name[323:353],
	-1000000000: _Result_name[353:375],
	-13:         _Result_name[375:388],
	-12:         _Result_name[388:409],
	-11:         _Result_name[409:435],
	-10:         _Result_name[435:457],
	-9:          _Result_name[457:482],
	-8:          _Result_name[482:507],
	-7:          _Result_name[507:534],
	-6:          _Result_name[534:557],
	-5:          _Result_name[557:580],
	-4:          _Result_name[580:597],
	-3:          _Result_name[597:624],
	-2:          _Result_name[624:650],
	-1:          _Result_name[650:674],
	1:           _Result_name[674:683],
	2:           _Result_name[683:690],
	3:           _Result_name[690:699],
	4:           _Result_name[699:710],
	5:           _Result_name[710:720],
	1000001003:  _Result_name[720:734],
	1000268000:  _Result_name[734:749],
	1000268001:  _Result_name[749:764],
	1000268002:  _Result_name[764:786],
	1000268003:  _Result_name[786:812],
	1000297000:  _Result_name[812:841],
}

func (i Result) String() string {
	if str, ok := _Result_map[i]; ok {
		return str
	}
	return "Result(" + strconv.FormatInt(int64(i), 10) + ")"
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SAMPLE_COUNT_1_BIT-1]
	_ = x[SAMPLE_COUNT_2_BIT-2]
	_ = x[SAMPLE_COUNT_4_BIT-4]
	_ = x[SAMPLE_COUNT_8_BIT-8]
	_ = x[SAMPLE_COUNT_16_BIT-16]
	_ = x[SAMPLE_COUNT_32_BIT-32]
	_ = x[SAMPLE_COUNT_64_BIT-64]
}

const (
	_SampleCountFlagBits_name_0 = "SAMPLE_COUNT_1_BITSAMPLE_COUNT_2_BIT"
	_SampleCountFlagBits_name_1 = "SAMPLE_COUNT_4_BIT"
	_SampleCountFlagBits_name_2 = "SAMPLE_COUNT_8_BIT"
	_SampleCountFlagBits_name_3 = "SAMPLE_COUNT_16_BIT"
	_SampleCountFlagBits_name_4 = "SAMPLE_COUNT_32_BIT"
	_SampleCountFlagBits_name_5 = "SAMPLE_COUNT_64_BIT"
)

var (
	_SampleCountFlagBits_index_0 = [...]uint8{0, 18, 36}
)

func (i SampleCountFlagBits) String() string {
	switch {
	case 1 <= i && i <= 2:
		i -= 1
		return _SampleCountFlagBits_name_0[_SampleCountFlagBits_index_0[i]:_SampleCountFlagBits_index_0[i+1]]
	case i == 4:
		return _SampleCountFlagBits_name_1
	case i == 8:
		return _SampleCountFlagBits_name_2
	case i == 16:
		return _SampleCountFlagBits_name_3
	case i == 32:
		return _SampleCountFlagBits_name_4
	case i == 64:
		return _SampleCountFlagBits_name_5
	default:
		return "SampleCountFlagBits(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SAMPLER_ADDRESS_MODE_REPEAT-0]
	_ = x[SAMPLER_ADDRESS_MODE_MIRRORED_REPEAT-1]
	_ = x[SAMPLER_ADDRESS_MODE_CLAMP_TO_EDGE-2]
	_ = x[SAMPLER_ADDRESS_MODE_CLAMP_TO_BORDER-3]
	_ = x[SAMPLER_ADDRESS_MODE_MIRROR_CLAMP_TO_EDGE-1000014000]
	_ = x[SAMPLER_ADDRESS_MODE_MIRROR_CLAMP_TO_EDGE_KHR-1000014000]
}

const (
	_SamplerAddressMode_name_0 = "SAMPLER_ADDRESS_MODE_REPEATSAMPLER_ADDRESS_MODE_MIRRORED_REPEATSAMPLER_ADDRESS_MODE_CLAMP_TO_EDGESAMPLER_ADDRESS_MODE_CLAMP_TO_BORDER"
	_SamplerAddressMode_name_1 = "SAMPLER_ADDRESS_MODE_MIRROR_CLAMP_TO_EDGE"
)

var (
	_SamplerAddressMode_index_0 = [...]uint8{0, 27, 63, 97, 133}
)

func (i SamplerAddressMode) String() string {
	switch {
	case 0 <= i && i <= 3:
		return _SamplerAddressMode_name_0[_SamplerAddressMode_index_0[i]:_SamplerAddressMode_index_0[i+1]]
	case i == 1000014000:
		return _SamplerAddressMode_name_1
	default:
		return "SamplerAddressMode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SAMPLER_CREATE_SUBSAMPLED_BIT_EXT-1000218000]
	_ = x[SAMPLER_CREATE_SUBSAMPLED_COARSE_RECONSTRUCTION_BIT_EXT-1000218000]
}

const _SamplerCreateFlagBits_name = "SAMPLER_CREATE_SUBSAMPLED_BIT_EXT"

var _SamplerCreateFlagBits_index = [...]uint8{0, 33}

func (i SamplerCreateFlagBits) String() string {
	i -= 1000218000
	if i >= SamplerCreateFlagBits(len(_SamplerCreateFlagBits_index)-1) {
		return "SamplerCreateFlagBits(" + strconv.FormatInt(int64(i+1000218000), 10) + ")"
	}
	return _SamplerCreateFlagBits_name[_SamplerCreateFlagBits_index[i]:_SamplerCreateFlagBits_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SAMPLER_MIPMAP_MODE_NEAREST-0]
	_ = x[SAMPLER_MIPMAP_MODE_LINEAR-1]
}

const _SamplerMipmapMode_name = "SAMPLER_MIPMAP_MODE_NEARESTSAMPLER_MIPMAP_MODE_LINEAR"

var _SamplerMipmapMode_index = [...]uint8{0, 27, 53}

func (i SamplerMipmapMode) String() string {
	if i < 0 || i >= SamplerMipmapMode(len(_SamplerMipmapMode_index)-1) {
		return "SamplerMipmapMode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _SamplerMipmapMode_name[_SamplerMipmapMode_index[i]:_SamplerMipmapMode_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SAMPLER_REDUCTION_MODE_WEIGHTED_AVERAGE-0]
	_ = x[SAMPLER_REDUCTION_MODE_MIN-1]
	_ = x[SAMPLER_REDUCTION_MODE_MAX-2]
	_ = x[SAMPLER_REDUCTION_MODE_MAX_EXT-2]
	_ = x[SAMPLER_REDUCTION_MODE_MIN_EXT-1]
	_ = x[SAMPLER_REDUCTION_MODE_WEIGHTED_AVERAGE_EXT-0]
}

const _SamplerReductionMode_name = "SAMPLER_REDUCTION_MODE_WEIGHTED_AVERAGESAMPLER_REDUCTION_MODE_MINSAMPLER_REDUCTION_MODE_MAX"

var _SamplerReductionMode_index = [...]uint8{0, 39, 65, 91}

func (i SamplerReductionMode) String() string {
	if i < 0 || i >= SamplerReductionMode(len(_SamplerReductionMode_index)-1) {
		return "SamplerReductionMode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _SamplerReductionMode_name[_SamplerReductionMode_index[i]:_SamplerReductionMode_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SAMPLER_YCBCR_MODEL_CONVERSION_RGB_IDENTITY-0]
	_ = x[SAMPLER_YCBCR_MODEL_CONVERSION_YCBCR_IDENTITY-1]
	_ = x[SAMPLER_YCBCR_MODEL_CONVERSION_YCBCR_709-2]
	_ = x[SAMPLER_YCBCR_MODEL_CONVERSION_YCBCR_601-3]
	_ = x[SAMPLER_YCBCR_MODEL_CONVERSION_YCBCR_2020-4]
	_ = x[SAMPLER_YCBCR_MODEL_CONVERSION_RGB_IDENTITY_KHR-0]
	_ = x[SAMPLER_YCBCR_MODEL_CONVERSION_YCBCR_2020_KHR-4]
	_ = x[SAMPLER_YCBCR_MODEL_CONVERSION_YCBCR_601_KHR-3]
	_ = x[SAMPLER_YCBCR_MODEL_CONVERSION_YCBCR_709_KHR-2]
	_ = x[SAMPLER_YCBCR_MODEL_CONVERSION_YCBCR_IDENTITY_KHR-1]
}

const _SamplerYcbcrModelConversion_name = "SAMPLER_YCBCR_MODEL_CONVERSION_RGB_IDENTITYSAMPLER_YCBCR_MODEL_CONVERSION_YCBCR_IDENTITYSAMPLER_YCBCR_MODEL_CONVERSION_YCBCR_709SAMPLER_YCBCR_MODEL_CONVERSION_YCBCR_601SAMPLER_YCBCR_MODEL_CONVERSION_YCBCR_2020"

var _SamplerYcbcrModelConversion_index = [...]uint8{0, 43, 88, 128, 168, 209}

func (i SamplerYcbcrModelConversion) String() string {
	if i < 0 || i >= SamplerYcbcrModelConversion(len(_SamplerYcbcrModelConversion_index)-1) {
		return "SamplerYcbcrModelConversion(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _SamplerYcbcrModelConversion_name[_SamplerYcbcrModelConversion_index[i]:_SamplerYcbcrModelConversion_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SAMPLER_YCBCR_RANGE_ITU_FULL-0]
	_ = x[SAMPLER_YCBCR_RANGE_ITU_NARROW-1]
	_ = x[SAMPLER_YCBCR_RANGE_ITU_FULL_KHR-0]
	_ = x[SAMPLER_YCBCR_RANGE_ITU_NARROW_KHR-1]
}

const _SamplerYcbcrRange_name = "SAMPLER_YCBCR_RANGE_ITU_FULLSAMPLER_YCBCR_RANGE_ITU_NARROW"

var _SamplerYcbcrRange_index = [...]uint8{0, 28, 58}

func (i SamplerYcbcrRange) String() string {
	if i < 0 || i >= SamplerYcbcrRange(len(_SamplerYcbcrRange_index)-1) {
		return "SamplerYcbcrRange(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _SamplerYcbcrRange_name[_SamplerYcbcrRange_index[i]:_SamplerYcbcrRange_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SCOPE_DEVICE_NV-1]
	_ = x[SCOPE_WORKGROUP_NV-2]
	_ = x[SCOPE_SUBGROUP_NV-3]
	_ = x[SCOPE_QUEUE_FAMILY_NV-5]
}

const (
	_ScopeNV_name_0 = "SCOPE_DEVICE_NVSCOPE_WORKGROUP_NVSCOPE_SUBGROUP_NV"
	_ScopeNV_name_1 = "SCOPE_QUEUE_FAMILY_NV"
)

var (
	_ScopeNV_index_0 = [...]uint8{0, 15, 33, 50}
)

func (i ScopeNV) String() string {
	switch {
	case 1 <= i && i <= 3:
		i -= 1
		return _ScopeNV_name_0[_ScopeNV_index_0[i]:_ScopeNV_index_0[i+1]]
	case i == 5:
		return _ScopeNV_name_1
	default:
		return "ScopeNV(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SEMAPHORE_IMPORT_TEMPORARY_BIT-1]
	_ = x[SEMAPHORE_IMPORT_TEMPORARY_BIT_KHR-1]
}

const _SemaphoreImportFlagBits_name = "SEMAPHORE_IMPORT_TEMPORARY_BIT"

var _SemaphoreImportFlagBits_index = [...]uint8{0, 30}

func (i SemaphoreImportFlagBits) String() string {
	i -= 1
	if i >= SemaphoreImportFlagBits(len(_SemaphoreImportFlagBits_index)-1) {
		return "SemaphoreImportFlagBits(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _SemaphoreImportFlagBits_name[_SemaphoreImportFlagBits_index[i]:_SemaphoreImportFlagBits_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SEMAPHORE_TYPE_BINARY-0]
	_ = x[SEMAPHORE_TYPE_TIMELINE-1]
	_ = x[SEMAPHORE_TYPE_BINARY_KHR-0]
	_ = x[SEMAPHORE_TYPE_TIMELINE_KHR-1]
}

const _SemaphoreType_name = "SEMAPHORE_TYPE_BINARYSEMAPHORE_TYPE_TIMELINE"

var _SemaphoreType_index = [...]uint8{0, 21, 44}

func (i SemaphoreType) String() string {
	if i < 0 || i >= SemaphoreType(len(_SemaphoreType_index)-1) {
		return "SemaphoreType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _SemaphoreType_name[_SemaphoreType_index[i]:_SemaphoreType_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SEMAPHORE_WAIT_ANY_BIT-1]
	_ = x[SEMAPHORE_WAIT_ANY_BIT_KHR-1]
}

const _SemaphoreWaitFlagBits_name = "SEMAPHORE_WAIT_ANY_BIT"

var _SemaphoreWaitFlagBits_index = [...]uint8{0, 22}

func (i SemaphoreWaitFlagBits) String() string {
	i -= 1
	if i >= SemaphoreWaitFlagBits(len(_SemaphoreWaitFlagBits_index)-1) {
		return "SemaphoreWaitFlagBits(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _SemaphoreWaitFlagBits_name[_SemaphoreWaitFlagBits_index[i]:_SemaphoreWaitFlagBits_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SHADER_FLOAT_CONTROLS_INDEPENDENCE_32_BIT_ONLY-0]
	_ = x[SHADER_FLOAT_CONTROLS_INDEPENDENCE_ALL-1]
	_ = x[SHADER_FLOAT_CONTROLS_INDEPENDENCE_NONE-2]
	_ = x[SHADER_FLOAT_CONTROLS_INDEPENDENCE_32_BIT_ONLY_KHR-0]
	_ = x[SHADER_FLOAT_CONTROLS_INDEPENDENCE_ALL_KHR-1]
	_ = x[SHADER_FLOAT_CONTROLS_INDEPENDENCE_NONE_KHR-2]
}

const _ShaderFloatControlsIndependence_name = "SHADER_FLOAT_CONTROLS_INDEPENDENCE_32_BIT_ONLYSHADER_FLOAT_CONTROLS_INDEPENDENCE_ALLSHADER_FLOAT_CONTROLS_INDEPENDENCE_NONE"

var _ShaderFloatControlsIndependence_index = [...]uint8{0, 46, 84, 123}

func (i ShaderFloatControlsIndependence) String() string {
	if i < 0 || i >= ShaderFloatControlsIndependence(len(_ShaderFloatControlsIndependence_index)-1) {
		return "ShaderFloatControlsIndependence(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ShaderFloatControlsIndependence_name[_ShaderFloatControlsIndependence_index[i]:_ShaderFloatControlsIndependence_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SHADER_GROUP_SHADER_GENERAL_KHR-0]
	_ = x[SHADER_GROUP_SHADER_CLOSEST_HIT_KHR-1]
	_ = x[SHADER_GROUP_SHADER_ANY_HIT_KHR-2]
	_ = x[SHADER_GROUP_SHADER_INTERSECTION_KHR-3]
}

const _ShaderGroupShaderKHR_name = "SHADER_GROUP_SHADER_GENERAL_KHRSHADER_GROUP_SHADER_CLOSEST_HIT_KHRSHADER_GROUP_SHADER_ANY_HIT_KHRSHADER_GROUP_SHADER_INTERSECTION_KHR"

var _ShaderGroupShaderKHR_index = [...]uint8{0, 31, 66, 97, 133}

func (i ShaderGroupShaderKHR) String() string {
	if i < 0 || i >= ShaderGroupShaderKHR(len(_ShaderGroupShaderKHR_index)-1) {
		return "ShaderGroupShaderKHR(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ShaderGroupShaderKHR_name[_ShaderGroupShaderKHR_index[i]:_ShaderGroupShaderKHR_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SHADER_INFO_TYPE_STATISTICS_AMD-0]
	_ = x[SHADER_INFO_TYPE_BINARY_AMD-1]
	_ = x[SHADER_INFO_TYPE_DISASSEMBLY_AMD-2]
}

const _ShaderInfoTypeAMD_name = "SHADER_INFO_TYPE_STATISTICS_AMDSHADER_INFO_TYPE_BINARY_AMDSHADER_INFO_TYPE_DISASSEMBLY_AMD"

var _ShaderInfoTypeAMD_index = [...]uint8{0, 31, 58, 90}

func (i ShaderInfoTypeAMD) String() string {
	if i < 0 || i >= ShaderInfoTypeAMD(len(_ShaderInfoTypeAMD_index)-1) {
		return "ShaderInfoTypeAMD(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ShaderInfoTypeAMD_name[_ShaderInfoTypeAMD_index[i]:_ShaderInfoTypeAMD_index[i+1]]
}
