// SPDX-FileCopyrightText : © 2022-2025 Galvanized Logic Inc.
// SPDX-License-Identifier: MIT

package render

// vulkan.go is the wrapper for the Vulkan API.
// It is organized with initialization near the top of the file and
// rendering at the bottom with rough groupings along the way.
// It is one big file because the author has not found a nice way to break
// it up that helps with comprehension, ie: Vulkan makes sense if one
// knows it, and no amount of file reorg seems to help if one does not.

import (
	"fmt"
	"log/slog"
	"math"
	"slices"
	"time"
	"unsafe"

	"github.com/gazed/vu/device"
	"github.com/gazed/vu/internal/render/vk"
	"github.com/gazed/vu/load"
	"github.com/gazed/vu/math/lin"
)

// vulkanRenderer contains vulkan specific rendering information.
// Variables are grouped by the function that initializes them.
type vulkanRenderer struct {
	title       string     // application name.
	clearColor  [4]float32 // rgba clear color.
	frameWidth  uint32     // current app size
	frameHeight uint32     //  ""

	// display resizing
	resizeWidth         uint32 // new size when resizing, 0 afterwards
	resizeHeight        uint32 //  ""
	resizesRequested    uint64 // track outstanding resize requests.
	resizesCompleted    uint64 //  ""
	recreatingSwapchain bool   // true when updating size.

	// createInstance initializes the root of the vulkan hierarchy.
	instance vk.Instance // vulkan root

	// createSurface links the vulkan instance to an OS display
	osdev   *device.Device // injected in activate()
	surface vk.SurfaceKHR  // device specific.

	// selectPhysicalDevice selects a GPU
	physicalDevice         vk.PhysicalDevice
	deviceExtensions       []string
	graphicsQIndex         uint32 // index chosen from graphics queue family.
	transferQIndex         uint32 // index chosen from transfer queue family.
	presentQIndex          uint32 // index chosen from present queue family.
	deviceLocalHostVisible bool   // true for device local and  host visible buffers.

	// createLogicalDevice initializes vulkan GPU resources
	device    vk.Device // logical device
	graphicsQ vk.Queue  // queue created from queue index
	transferQ vk.Queue  //  ""
	presentQ  vk.Queue  //  ""

	// setRenderProperties
	surfaceFormat      vk.SurfaceFormatKHR            // chosen surface format
	surfacePresentMode vk.PresentModeKHR              // chosen present mode
	surfaceTransform   vk.SurfaceTransformFlagBitsKHR // surface transform flags
	depthFormat        vk.Format                      // 3D requires depth
	frameCount         uint32                         // two frames
	imageCount         uint32                         // three swapchain images.

	// createCommandPools
	graphicsQCmdPool vk.CommandPool // graphics queue command pool

	// createSwapchainResources
	// imageIndex tracks the swapchain image acquired by vkAcquireNextImageKHR.
	// This can be any of the swapchain images.
	// frameIndex tracks frame information in order and will always increment
	// each frame and loop using mod frameCount: 0, 1, 0, 1, 0, 1..
	swapchain  vk.SwapchainKHR //
	depthImage vulkanImage     // 3D requires depth
	images     []vk.Image      // images owned by swapchain
	imageIndex uint32          // index for images - set by vkAcquireNextImageKHR
	views      []vk.ImageView  // one view per swapchain image
	frames     []vulkanFrame   // frame resources for maxFrames
	frameIndex uint32          // index for frames - loop using mod maxFrames

	// create a semaphore per image. These are signaled when the GPU has
	// finished rendering an image and the image is ready for presentation.
	imageRendered []vk.Semaphore // wait for images renders prior to presentation.

	// render frame dynamic state.
	viewport vk.Viewport // same as frame size.
	scissor  vk.Rect2D   // same as frame size.

	// mesh vertex attribute data buffers.
	vertexMem []vulkanBuffer // 0:vData, 1:vIndex
	vertexPtr []uint32       // 0:vData, 1:vIndex - first empty data index.

	// ray acceleration data buffers.
	accelMem    []vulkanBuffer         // acceleration storage buffers BLAS, INST, TLAS, SCRATCH
	accelPtr    []vk.DeviceSize        // current end of each accel buffer
	accelMeshes map[uint32]accelStruct // BLAS accel data indexed by mesh ID
	accelScenes map[uint32]accelStruct // TLAS accel data indexed by scene ID

	// application GPU resources.
	meshes   []vulkanMesh    // application GPU mesh data
	textures []vulkanTexture // application GPU texture data
	shaders  []vulkanShader  // shaders - one pipeline per shader.
}

// buffer indexing.
const (
	// vertex data buffer types.
	vData  = 0 // all (non-index) vertex data
	vIndex = 1 // all vertex index data.

	// acceleration data buffer types.
	aBLAS    = 0 // all bottom level acceleration structs
	aINST0   = 1 // frame 0 TLAS instance acceleration structs
	aINST1   = 2 // frame 1 TLAS instance acceleration structs
	aTLAS    = 3 // the top level acceleration struct
	aSCRATCH = 4 // scratch buffer for acceleration create/update
)

// vkEnabledLayers can be modified by debug builds
// by overriding the addValidationLayer method.
var vkEnabledLayers []string = []string{} // enabled vulkan layers
var addValidationLayer func([]string) ([]string, error) = func(layers []string) ([]string, error) { return layers, nil }

// getVulkanRenderer acquires the vulkan resources needed to render scenes.
func getVulkanRenderer(dev *device.Device, title string) (vr *vulkanRenderer, err error) {
	vr = &vulkanRenderer{}
	vr.accelMeshes = map[uint32]accelStruct{}
	vr.accelScenes = map[uint32]accelStruct{}
	vr.title = title

	// load the vulkan library
	if err := vk.LoadVulkan(""); err != nil {
		return nil, fmt.Errorf("failed to load vulkan library: %w", err)
	}

	// create the vulkan stack.
	vr.osdev = dev
	vr.frameWidth, vr.frameHeight = vr.osdev.SurfaceSize() // initial size
	vr.deviceExtensions = []string{
		vk.KHR_SWAPCHAIN_EXTENSION_NAME,                // VK_KHR_swapchain is window system specific.
		vk.KHR_RAY_QUERY_EXTENSION_NAME,                // VK_KHR_ray_query
		vk.KHR_ACCELERATION_STRUCTURE_EXTENSION_NAME,   // VK_KHR_acceleration_structure
		vk.KHR_BUFFER_DEVICE_ADDRESS_EXTENSION_NAME,    // VK_KHR_buffer_device_address
		vk.KHR_DEFERRED_HOST_OPERATIONS_EXTENSION_NAME, // VK_KHR_deferred_host_operations
		vk.EXT_DESCRIPTOR_INDEXING_EXTENSION_NAME,      // VK_EXT_descriptor_indexing
		vk.KHR_SPIRV_1_4_EXTENSION_NAME,                // VK_KHR_spirv_1_4
	}

	// acquire resources from the top down.
	createFunctions := []func() error{
		// one time on startup
		vr.createInstance,       // one instance
		vr.createSurface,        // one surface - vulkan_windows.go
		vr.selectPhysicalDevice, // one physical device selected
		vr.createLogicalDevice,  // one logical device with queues
		vr.createCommandPools,   // one pool per queue family

		// render properties and swapchain are set on initialization
		// and updated and/or recreated on a window resize.
		vr.setRenderProperties,      // surface format, depth format, etc.
		vr.createImageSemaphores,    // one render semaphore per image.
		vr.createSwapchainResources, // one swapchain with one render frame per image
		vr.createRenderFrames,       // two frames

		// create storage for mesh data.
		vr.createVertexBuffers,

		// create storage for ray query acceleration data.
		vr.createAccelerationBuffers,
	}
	for _, create := range createFunctions {
		if err := create(); err != nil {
			vr.dispose()
			return nil, err
		}
	}

	// init is done... throw it back to the application
	// to load the shaders and render data.
	slog.Info("vulkan initialized")
	return vr, err
}

// dispose releases vulkan resources in the opposite order
// they were allocated..
func (vr *vulkanRenderer) dispose() {
	if vr.device != 0 {
		vk.DeviceWaitIdle(vr.device)
	}

	// remove application allocated resources.
	vr.disposeAccelerationStructs()
	vr.disposeAccelerationBuffers()
	vr.disposeVertexBuffers()
	for i := range vr.textures {
		vr.dropTexture(uint32(i))
	}
	for sid := range vr.shaders {
		vr.disposeShader(&vr.shaders[sid])
	}

	// dispose the per-image render semaphores.
	vr.disposeImageSemaphores()

	// swapchain and related resources: image views, depthbuffer
	// Also scraps all application mesh data.
	vr.disposeRenderFrames()
	vr.disposeSwapchainResources()

	// instance, surface, physical device, logical device, command pools.
	vr.disposeCommandPools()
	if vr.device != 0 {
		vk.DestroyDevice(vr.device, nil)
		vr.device = 0
	}
	if vr.surface != 0 {
		vk.DestroySurfaceKHR(vr.instance, vr.surface, nil)
		vr.surface = 0
	}
	if vr.instance != 0 {
		vk.DestroyInstance(vr.instance, nil)
		vr.instance = 0
	}
}

// =============================================================================
// one-time initialization at startup.
// vulkan instance, physical and logical device selection.

// createInstance initializes the root of the vulkan hierarchy.
func (vr *vulkanRenderer) createInstance() (err error) {
	vkEnabledLayers, err := addValidationLayer(vkEnabledLayers) // vulkan_debug.go
	if err != nil {
		return err
	}

	// create the vulkan instance.
	instanceInfo := vk.InstanceCreateInfo{
		PApplicationInfo: &vk.ApplicationInfo{
			PApplicationName:   vr.title,
			ApplicationVersion: vk.MAKE_VERSION(1, 0, 0),
			PEngineName:        "vu",
			EngineVersion:      vk.HEADER_VERSION_COMPLETE,
			ApiVersion:         vk.API_VERSION_1_3,
		},
		PpEnabledLayerNames:     vkEnabledLayers,
		PpEnabledExtensionNames: vr.instanceExtensions(), // vulkan_windows.go
	}
	vr.instance, err = vk.CreateInstance(&instanceInfo, nil)

	// set the vulkan instance in the vulkan bindings.
	// It is used to find function pointers to other vulkan methods.
	vk.VKInst = uintptr(vr.instance)
	return err
}

// deviceCandidate is a physical device that meets the requirements.
type deviceCandidate struct {
	physicalDevice vk.PhysicalDevice // the physical device
	graphicsQIndex uint32            // queue index
	transferQIndex uint32            // queue index
	presentQIndex  uint32            // queue index
}

// selectPhysicalDevice finds a physical device for 3D rendering.
func (vr *vulkanRenderer) selectPhysicalDevice() error {
	devices, err := vk.EnumeratePhysicalDevices(vr.instance)
	if err != nil {
		fmt.Errorf("vk.EnumeratePhysicalDevices: %w", err)
	}
	surface := surfaceProperties{}

	GPU = INTEGRATED_GPU // set global default unless discrete GPU found.
	candidates := []deviceCandidate{}
	for _, d := range devices {
		props2 := vk.GetPhysicalDeviceProperties2(d)
		properties := props2.Properties

		// require vulkan 1.4
		if properties.ApiVersion < vk.API_VERSION_1_3 {
			slog.Warn("vulkan version to low",
				"name", properties.DeviceName,
				"version", vr.version(properties.ApiVersion))
			break
		}

		// ensure that the device has the required queues.
		graphicsQIndex, transferQIndex, presentQIndex := -1, -1, -1
		minTransferScore := 255
		queueProps := vk.GetPhysicalDeviceQueueFamilyProperties(d)
		for i, qp := range queueProps {
			transferScore := 0
			if (qp.QueueFlags & vk.QUEUE_GRAPHICS_BIT) == vk.QUEUE_GRAPHICS_BIT {
				transferScore += 1
				graphicsQIndex = i
			}
			if (qp.QueueFlags & vk.QUEUE_COMPUTE_BIT) == vk.QUEUE_COMPUTE_BIT {
				transferScore += 1
			}

			// look for a queue that is dedicated to transfers,
			// ideally not shared with computer and/or graphics queues
			if (qp.QueueFlags & vk.QUEUE_TRANSFER_BIT) == vk.QUEUE_TRANSFER_BIT {
				if transferScore < minTransferScore {
					minTransferScore = transferScore
					transferQIndex = i
				}
			}

			// ensure the device can present an image
			canPresent, err := vk.GetPhysicalDeviceSurfaceSupportKHR(d, uint32(i), vr.surface)
			if err != nil {
				return fmt.Errorf("vk.GetPhysicalDeviceSurfaceSupportKHR:")
			}
			if canPresent {
				presentQIndex = i
			}
		}
		if graphicsQIndex < 0 || transferQIndex < 0 || presentQIndex < 0 {
			slog.Warn("missing required Queues")
			break // missing required queue
		}

		// update the surface swapchain information.
		if err := vr.getSurfaceProperties(&surface, d); err != nil {
			return fmt.Errorf("querySwapchainSupport %w", err)
		}
		if len(surface.formats) < 1 || len(surface.presentModes) < 1 {
			slog.Warn("missing swapchain support")
			break // required swapchain suipport not present on this device.
		}

		// query device extensions
		extensions, err := vk.EnumerateDeviceExtensionProperties(d, "")
		if err != nil {
			return fmt.Errorf("vk.EnumerateDeviceExtensionProperties:")
		}
		availableExtensions := map[string]bool{}
		for _, ext := range extensions {
			availableExtensions[ext.ExtensionName] = true
		}

		// check that the required device extensions are available
		for _, req := range vr.deviceExtensions {
			if _, ok := availableExtensions[req]; !ok {
				slog.Warn("missing required device extension", "name", req)
				break // try the next device.
			}
		}

		vr.deviceLocalHostVisible = false
		memProps := vk.GetPhysicalDeviceMemoryProperties(d)
		for i := uint32(0); i < memProps.MemoryTypeCount; i++ {
			if (memProps.MemoryTypes[i].PropertyFlags&vk.MEMORY_PROPERTY_HOST_VISIBLE_BIT) != 0 &&
				(memProps.MemoryTypes[i].PropertyFlags&vk.MEMORY_PROPERTY_DEVICE_LOCAL_BIT) != 0 {
				vr.deviceLocalHostVisible = true
				break
			}
		}

		// check for required features
		features := vk.GetPhysicalDeviceFeatures(d, &vk.PhysicalDeviceFeatures{})
		if !features.SamplerAnisotropy {
			slog.Warn("missing features.SamplerAnisotropy", "device", vr.physicalDevice)
			continue // device missing samplerAnisotropy
		}

		// reaching here means that the device meets all requirements.
		slog.Info("vulkan device found", "device", d,
			"name", properties.DeviceName,
			"driver", vr.version(properties.DriverVersion),
			"api", vr.version(properties.ApiVersion))

		// save devices that meet the requirements
		dc := deviceCandidate{
			physicalDevice: d,                      // the physical device
			graphicsQIndex: uint32(graphicsQIndex), // queue index
			transferQIndex: uint32(transferQIndex), // queue index
			presentQIndex:  uint32(presentQIndex),  // queue index
		}

		// prefer discrete GPU if available, ie: put the discrete GPUs at the front.
		if properties.DeviceType == vk.PHYSICAL_DEVICE_TYPE_DISCRETE_GPU {
			candidates = append([]deviceCandidate{dc}, candidates...)
			GPU = DISCRETE_GPU // discrete GPU used when found.
		} else {
			candidates = append(candidates, dc)
		}
	}

	// save the preferred physical device in the vulkan context
	if len(candidates) > 0 {
		dc := candidates[0] // discrete GPU if available.
		vr.physicalDevice = dc.physicalDevice
		vr.graphicsQIndex = dc.graphicsQIndex
		vr.transferQIndex = dc.transferQIndex
		vr.presentQIndex = dc.presentQIndex
		slog.Info("vulkan device created", "device", vr.physicalDevice, "discrete", GPU == DISCRETE_GPU)
		return nil // found a physical device.
	}
	return fmt.Errorf("no physical device found")
}

// createLogicalDevice
func (vr *vulkanRenderer) createLogicalDevice() (err error) {

	// create the queues for each queue family
	qIndicies := []uint32{vr.graphicsQIndex}
	if vr.graphicsQIndex != vr.presentQIndex {
		qIndicies = append(qIndicies, vr.presentQIndex) // separate present queue
	}
	if vr.graphicsQIndex != vr.transferQIndex {
		qIndicies = append(qIndicies, vr.transferQIndex) // separate transfer queue
	}
	queueInfos := []vk.DeviceQueueCreateInfo{}
	for _, qIndex := range qIndicies {
		qi := vk.DeviceQueueCreateInfo{}
		qi.QueueFamilyIndex = qIndex
		qi.PQueuePriorities = []float32{1.0}
		qi.Flags = 0
		qi.PNext = nil
		queueInfos = append(queueInfos, qi)
	}

	// chain device extensions that are not yet core features but are needed now.
	enableAcceleration := vk.PhysicalDeviceAccelerationStructureFeaturesKHR{}
	enableAcceleration.AccelerationStructure = true
	enableRayQuery := vk.PhysicalDeviceRayQueryFeaturesKHR{}
	enableRayQuery.RayQuery = true
	enableRayQuery.PNext = unsafe.Pointer(enableAcceleration.ToVK())
	// chain the 1.2 features
	features12 := vk.PhysicalDeviceVulkan12Features{}
	features12.DescriptorBindingVariableDescriptorCount = true
	features12.ShaderSampledImageArrayNonUniformIndexing = true
	features12.BufferDeviceAddress = true
	features12.DescriptorIndexing = true
	features12.PNext = unsafe.Pointer(enableRayQuery.ToVK())
	// chain the 1.3 features
	features13 := vk.PhysicalDeviceVulkan13Features{}
	features13.Synchronization2 = true
	features13.DynamicRendering = true
	features13.PNext = unsafe.Pointer(features12.ToVK())
	// chain the features
	features2 := vk.PhysicalDeviceFeatures2{}
	features2.Features.SamplerAnisotropy = true
	features2.PNext = unsafe.Pointer(features13.ToVK())

	// create the device with the desired features,
	// fail unless all features are supported.
	deviceCreateInfo := vk.DeviceCreateInfo{
		PNext:                   unsafe.Pointer(features2.ToVK()),
		PQueueCreateInfos:       queueInfos,
		PEnabledFeatures:        nil, // pNext is using VkPhysicalDeviceFeatures2
		PpEnabledExtensionNames: vr.deviceExtensions,
	}
	vr.device, err = vk.CreateDevice(vr.physicalDevice, &deviceCreateInfo, nil)
	if err != nil {
		return fmt.Errorf("vk.CreateDevice:: %w", err)
	}

	// create the queues from the queue indicies
	vr.graphicsQ = vk.GetDeviceQueue(vr.device, vr.graphicsQIndex, 0)
	vr.presentQ = vk.GetDeviceQueue(vr.device, vr.presentQIndex, 0)
	vr.transferQ = vk.GetDeviceQueue(vr.device, vr.transferQIndex, 0)

	// success
	slog.Debug("vulkan logical device created")
	return nil
}

// surfaceProperties
type surfaceProperties struct {
	formats      []vk.SurfaceFormatKHR
	capabilities vk.SurfaceCapabilitiesKHR
	presentModes []vk.PresentModeKHR
}

// getSurfaceProperites gets swapchain surface information from the physicalDevice.
func (vr *vulkanRenderer) getSurfaceProperties(surf *surfaceProperties, pd vk.PhysicalDevice) (err error) {
	surf.capabilities, err = vk.GetPhysicalDeviceSurfaceCapabilitiesKHR(pd, vr.surface)
	if err != nil {
		return fmt.Errorf("vk.GetPhysicalDeviceSurfaceCapabilitiesKHR: %w", err)
	}
	surf.formats, err = vk.GetPhysicalDeviceSurfaceFormatsKHR(pd, vr.surface)
	if err != nil {
		return fmt.Errorf("vk.GetPhysicalDeviceSurfaceFormatsKHR: %w", err)
	}
	surf.presentModes, err = vk.GetPhysicalDeviceSurfacePresentModesKHR(pd, vr.surface)
	if err != nil {
		return fmt.Errorf("vk.GetPhysicalDeviceSurfacePresentModesKHR: %w", err)
	}
	return nil
}

// createCommandPools - currently just the graphics queue command pool.
func (vr *vulkanRenderer) createCommandPools() (err error) {
	poolInfo := vk.CommandPoolCreateInfo{
		QueueFamilyIndex: vr.graphicsQIndex,
		Flags:            vk.COMMAND_POOL_CREATE_RESET_COMMAND_BUFFER_BIT,
	}
	vr.graphicsQCmdPool, err = vk.CreateCommandPool(vr.device, &poolInfo, nil)
	return err
}
func (vr *vulkanRenderer) disposeCommandPools() {
	if vr.graphicsQCmdPool != 0 {
		vk.DestroyCommandPool(vr.device, vr.graphicsQCmdPool, nil)
		vr.graphicsQCmdPool = 0
	}
}

// ============================================================================
// render properties set on initialization and updated on window resize.

// setRenderProperties sets or updates the following properties needed for
// the renderpass and swapchain
// : vr.surfaceFormat
// : vr.surfacePresentMode
// : vr.surfaceTransform
// : vr.frameWidth
// : vr.frameHeight
// : vr.depthFormat
func (vr *vulkanRenderer) setRenderProperties() (err error) {
	surface := surfaceProperties{}
	if err = vr.getSurfaceProperties(&surface, vr.physicalDevice); err != nil {
		return err
	}

	// choose a swap surface format
	foundFormat := false
	for _, sf := range surface.formats {

		// preferred format supported by all graphics cards
		// if sf.Format == vk.FORMAT_B8G8R8A8_UNORM && sf.ColorSpace == vk.COLOR_SPACE_SRGB_NONLINEAR_KHR {
		if sf.Format == vk.FORMAT_B8G8R8A8_SRGB && sf.ColorSpace == vk.COLOR_SPACE_SRGB_NONLINEAR_KHR {
			vr.surfaceFormat = sf
			foundFormat = true
			break
		}
	}
	if !foundFormat {
		// otherwise just use the first one.
		vr.surfaceFormat = surface.formats[0]
		slog.Warn("vulkan using default surface format")
	}

	// find the best present mode.
	vr.surfacePresentMode = vk.PRESENT_MODE_FIFO_KHR // always exists.
	for _, mode := range surface.presentModes {
		if mode == vk.PRESENT_MODE_MAILBOX_KHR {
			vr.surfacePresentMode = mode
			break
		}
	}

	// remember the surface transform
	vr.surfaceTransform = surface.capabilities.CurrentTransform // default
	if (surface.capabilities.SupportedTransforms & vk.SURFACE_TRANSFORM_IDENTITY_BIT_KHR) == vk.SURFACE_TRANSFORM_IDENTITY_BIT_KHR {
		vr.surfaceTransform = vk.SURFACE_TRANSFORM_IDENTITY_BIT_KHR
	} else {
		slog.Warn("VK_SURFACE_TRANSFORM_IDENTITY_BIT_KHR not available")
	}

	// update the extent and clamp to the GPU supported values
	min := surface.capabilities.MinImageExtent
	max := surface.capabilities.MaxImageExtent
	vr.frameWidth = lin.Clamp(vr.frameWidth, min.Width, max.Width)
	vr.frameHeight = lin.Clamp(vr.frameHeight, min.Height, max.Height)
	if surface.capabilities.CurrentExtent.Width == ^uint32(0) {
		slog.Warn("surface extent is max") // warn surface extent is unset.
	}

	// set vr.depthFormat
	if err = vr.setDepthFormat(); err != nil {
		return err
	}

	// ideally two frames so the CPU can record commands in one frame while the
	// GPU renders an image using the other.
	// Ask for 3 swapchain images in case one image is presented a bit slow.
	// Note: surface.capabilities.MinImageCount is at least 1.
	//       surface.capabilities.MaxImageCount of 0 means unlimited.
	vr.frameCount, vr.imageCount = 2, 3
	if surface.capabilities.MaxImageCount != 0 && vr.imageCount > surface.capabilities.MaxImageCount {
		vr.imageCount = surface.capabilities.MaxImageCount
	}
	slog.Info("vulkan surface",
		"width", vr.frameWidth, "height", vr.frameHeight,
		"frames", vr.frameCount, "images", vr.imageCount)
	return nil
}

// setDepthFormat updates vr.depthFormat from supported depth formats.
func (vr *vulkanRenderer) setDepthFormat() (err error) {
	candidates := []vk.Format{ // ordered by preference
		vk.FORMAT_D32_SFLOAT,         // 32 bit signed depth
		vk.FORMAT_D32_SFLOAT_S8_UINT, // 32 bit signed depth, 8 bit unsigned stencil
		vk.FORMAT_D24_UNORM_S8_UINT,  // 24 bit unsized normalized depth, 8 bit unsigned stencil
	}
	flags := uint32(vk.FORMAT_FEATURE_DEPTH_STENCIL_ATTACHMENT_BIT)
	for _, format := range candidates {
		props := vk.GetPhysicalDeviceFormatProperties(vr.physicalDevice, format)
		if uint32(props.OptimalTilingFeatures)&flags == flags {
			vr.depthFormat = format
			return nil
		}
		if uint32(props.LinearTilingFeatures)&flags == flags {
			vr.depthFormat = format
			return nil
		}
	}
	vr.depthFormat = vk.FORMAT_UNDEFINED
	return fmt.Errorf("setDepthFormat failed")
}

// createImageSemaphores tracks the render complete semaphore for each image.
// Called once on startup after vr.imageCount is initialized. See:
// - https://docs.vulkan.org/guide/latest/swapchain_semaphore_reuse.html
func (vr *vulkanRenderer) createImageSemaphores() (err error) {
	vr.imageRendered = make([]vk.Semaphore, vr.imageCount)
	for i := range vr.imageRendered {
		vr.imageRendered[i], err = vk.CreateSemaphore(vr.device, &vk.SemaphoreCreateInfo{}, nil)
		if err != nil {
			return fmt.Errorf("vk.CreateSemaphore.2: %w", err)
		}
	}
	return nil
}

// disposeImageSemaphores is called once on shutdown.
func (vr *vulkanRenderer) disposeImageSemaphores() {
	for i := range vr.imageRendered {
		if vr.imageRendered[i] != 0 {
			vk.DestroySemaphore(vr.device, vr.imageRendered[i], nil)
			vr.imageRendered[i] = 0
		}
	}
}

// =============================================================================
// swapchain recreated each window resize.

// createSwapchainResources creates all the swapchain resources.
// It is also used to recreate the swapchain after a resize.
func (vr *vulkanRenderer) createSwapchainResources() (err error) {
	creates := []func() error{
		vr.createSwapchain,   // swapchain
		vr.createDepthBuffer, // swapchain has one depth buffer
		vr.createImageViews,  // one view per swapchain image
	}
	for _, create := range creates {
		if err := create(); err != nil {
			vr.dispose()
			return err
		}
	}
	slog.Debug("vulkan swapchain created", "images", len(vr.images))
	return nil
}

// disposeSwapchainResources called on shutdown or display resize
func (vr *vulkanRenderer) disposeSwapchainResources() {
	for i := range vr.images {
		if vr.views[i] != 0 {
			vk.DestroyImageView(vr.device, vr.views[i], nil)
			vr.views[i] = 0
		}
	}
	vr.disposeImage(&vr.depthImage)
	if vr.swapchain != 0 {
		vk.DestroySwapchainKHR(vr.device, vr.swapchain, nil)
		vr.swapchain = 0
	}
}

// createSwapchain initializes the underlying render image frames.
func (vr *vulkanRenderer) createSwapchain() (err error) {
	extent := vk.Extent2D{Width: vr.frameWidth, Height: vr.frameHeight}

	// create the swapchain using a shared queue or
	// separate queues for graphics and presentation
	swapchainInfo := vk.SwapchainCreateInfoKHR{
		Flags:            0,
		Surface:          vr.surface,
		MinImageCount:    vr.imageCount,
		ImageFormat:      vr.surfaceFormat.Format,
		ImageColorSpace:  vr.surfaceFormat.ColorSpace,
		ImageExtent:      extent,
		ImageArrayLayers: 1,
		ImageUsage:       vk.IMAGE_USAGE_COLOR_ATTACHMENT_BIT,
		ImageSharingMode: vk.SHARING_MODE_EXCLUSIVE,
		PreTransform:     vr.surfaceTransform,
		CompositeAlpha:   vk.COMPOSITE_ALPHA_OPAQUE_BIT_KHR,
		PresentMode:      vr.surfacePresentMode,
		Clipped:          true,
		OldSwapchain:     0,
	}

	// FUTURE: allow separate present and graphic queues
	// if vr.graphicsQIndex != vr.presentQIndex {
	// 	swapchainInfo.ImageSharingMode = vk.SHARING_MODE_CONCURRENT
	// 	swapchainInfo.PQueueFamilyIndices = []uint32{vr.graphicsQIndex, vr.presentQIndex}
	// } else {
	// 	swapchainInfo.ImageSharingMode = vk.SHARING_MODE_EXCLUSIVE
	// }
	vr.swapchain, err = vk.CreateSwapchainKHR(vr.device, &swapchainInfo, nil)
	if err != nil {
		return fmt.Errorf("vk.CreateSwapchainKHR:: %w", err)
	}
	vr.images, err = vk.GetSwapchainImagesKHR(vr.device, vr.swapchain)
	if err != nil {
		return fmt.Errorf("vk.GetSwapchainImagesKHR: %w", err)
	}
	return nil
}

// createImageViews - one per swapchain image
func (vr *vulkanRenderer) createImageViews() (err error) {
	vr.views = make([]vk.ImageView, len(vr.images))
	for i := range vr.images {
		vr.views[i], err = vr.createImageView(vr.images[i], vr.surfaceFormat.Format, vk.IMAGE_ASPECT_COLOR_BIT)
		if err != nil {
			return fmt.Errorf("createFrameImageView: %w", err)
		}
	}
	return nil
}

// createDepthBuffer creates a depth image resource separate from the swapchain.
func (vr *vulkanRenderer) createDepthBuffer() (err error) {
	vr.depthImage = vulkanImage{}
	vr.depthImage.width = vr.frameWidth
	vr.depthImage.height = vr.frameHeight
	err = vr.createImage(&vr.depthImage, vr.depthFormat,
		vk.IMAGE_USAGE_DEPTH_STENCIL_ATTACHMENT_BIT,
		vk.MEMORY_PROPERTY_DEVICE_LOCAL_BIT)
	if err != nil {
		return err
	}
	vr.depthImage.view, err = vr.createImageView(vr.depthImage.handle, vr.depthFormat, vk.IMAGE_ASPECT_DEPTH_BIT)
	if err != nil {
		return err
	}
	return nil
}

// createVertexBuffers creates one buffer for all
// mesh vertex data (instanced and non-instanced)
// and a separate buffer for triangle index data.
func (vr *vulkanRenderer) createVertexBuffers() (err error) {
	vr.vertexMem = make([]vulkanBuffer, 2)
	vr.vertexPtr = []uint32{0, 0}
	memProps := vk.MEMORY_PROPERTY_DEVICE_LOCAL_BIT

	// vertex data.
	size := vk.DeviceSize(GPUMaxVertexBytes)
	usage := vk.BUFFER_USAGE_VERTEX_BUFFER_BIT |
		vk.BUFFER_USAGE_TRANSFER_DST_BIT |
		vk.BUFFER_USAGE_STORAGE_BUFFER_BIT |
		vk.BUFFER_USAGE_SHADER_DEVICE_ADDRESS_BIT |
		vk.BUFFER_USAGE_ACCELERATION_STRUCTURE_BUILD_INPUT_READ_ONLY_BIT_KHR
	if err = vr.createBuffer(&vr.vertexMem[vData], size, usage, memProps); err != nil {
		return fmt.Errorf("createVertexBuffers:data %w", err)
	}

	// need vertex data pointer when building acceleration structs
	i0 := vk.BufferDeviceAddressInfo{Buffer: vr.vertexMem[vData].handle}
	vr.vertexMem[vData].address = vk.GetBufferDeviceAddress(vr.device, &i0)

	// triangle indexes
	isize := vk.DeviceSize(GPUMaxVIndexBytes)
	iusage := vk.BUFFER_USAGE_INDEX_BUFFER_BIT | vk.BUFFER_USAGE_STORAGE_BUFFER_BIT |
		vk.BUFFER_USAGE_TRANSFER_DST_BIT |
		vk.BUFFER_USAGE_STORAGE_BUFFER_BIT |
		vk.BUFFER_USAGE_SHADER_DEVICE_ADDRESS_BIT |
		vk.BUFFER_USAGE_ACCELERATION_STRUCTURE_BUILD_INPUT_READ_ONLY_BIT_KHR
	if err = vr.createBuffer(&vr.vertexMem[vIndex], isize, iusage, memProps); err != nil {
		return fmt.Errorf("createVertexBuffers:indexes %w", err)
	}

	// need vertex triangle index pointer when building acceleration structs
	i1 := vk.BufferDeviceAddressInfo{Buffer: vr.vertexMem[vIndex].handle}
	vr.vertexMem[vIndex].address = vk.GetBufferDeviceAddress(vr.device, &i1)
	return nil
}
func (vr *vulkanRenderer) disposeVertexBuffers() {
	if len(vr.vertexMem) == 2 {
		vr.disposeBuffer(&vr.vertexMem[vIndex])
		vr.disposeBuffer(&vr.vertexMem[vData])
		vr.vertexMem = nil // trigger garbage collect
	}
}

// createAccelerationBuffers reserves device storage space for
// the acceleration structures needed for ray query and ray trace.
//
// FUTURE: optimizations recommended by:
// https://nvpro-samples.github.io/vk_raytracing_tutorial_KHR/concepts/acceleration-structures/
//   - Optimization Strategy: While you can allocate a separate scratch buffer
//     for each build, a more efficient approach is to determine the size of the
//     largest scratch buffer needed and reuse a single buffer of that size,
//     provided you handle the synchronization correctly between build calls.
//   - TLAS Updates: For TLAS updates (e.g., for animations), you can often keep
//     reusing the same scratch buffer as long as you ensure the previous frame's
//     acceleration structure writes are complete before starting the new build.
func (vr *vulkanRenderer) createAccelerationBuffers() (err error) {
	vr.accelMem = make([]vulkanBuffer, aSCRATCH+1)
	vr.accelPtr = make([]vk.DeviceSize, aSCRATCH+1)
	memProps := vk.MEMORY_PROPERTY_DEVICE_LOCAL_BIT

	// BLAS buffers
	size := vk.DeviceSize(GPUMaxAccelerationBLASBytes)
	usage := vk.BUFFER_USAGE_ACCELERATION_STRUCTURE_STORAGE_BIT_KHR |
		vk.BUFFER_USAGE_SHADER_DEVICE_ADDRESS_BIT |
		vk.BUFFER_USAGE_ACCELERATION_STRUCTURE_BUILD_INPUT_READ_ONLY_BIT_KHR
	if err = vr.createBuffer(&vr.accelMem[aBLAS], size, usage, memProps); err != nil {
		return fmt.Errorf("createAccelerationBuffers:BLAS %w", err)
	}

	// TLAS instance buffer for frame0
	size = vk.DeviceSize(GPUMaxAccelerationINSTBytes)
	usage = vk.BUFFER_USAGE_SHADER_DEVICE_ADDRESS_BIT |
		vk.BUFFER_USAGE_TRANSFER_DST_BIT |
		vk.BUFFER_USAGE_ACCELERATION_STRUCTURE_BUILD_INPUT_READ_ONLY_BIT_KHR
	if err = vr.createBuffer(&vr.accelMem[aINST0], size, usage, memProps); err != nil {
		return fmt.Errorf("createAccelerationBuffers:INST0 %w", err)
	}
	info := vk.BufferDeviceAddressInfo{Buffer: vr.accelMem[aINST0].handle}
	vr.accelMem[aINST0].address = vk.GetBufferDeviceAddress(vr.device, &info)
	// TLAS instance buffer for frame1
	size = vk.DeviceSize(GPUMaxAccelerationINSTBytes)
	usage = vk.BUFFER_USAGE_SHADER_DEVICE_ADDRESS_BIT |
		vk.BUFFER_USAGE_TRANSFER_DST_BIT |
		vk.BUFFER_USAGE_ACCELERATION_STRUCTURE_BUILD_INPUT_READ_ONLY_BIT_KHR
	if err = vr.createBuffer(&vr.accelMem[aINST1], size, usage, memProps); err != nil {
		return fmt.Errorf("createAccelerationBuffers:INST1 %w", err)
	}
	info = vk.BufferDeviceAddressInfo{Buffer: vr.accelMem[aINST1].handle}
	vr.accelMem[aINST1].address = vk.GetBufferDeviceAddress(vr.device, &info)

	// TLAS buffer
	size = vk.DeviceSize(GPUMaxAccelerationTLASBytes)
	usage = vk.BUFFER_USAGE_ACCELERATION_STRUCTURE_STORAGE_BIT_KHR |
		vk.BUFFER_USAGE_SHADER_DEVICE_ADDRESS_BIT
	if err = vr.createBuffer(&vr.accelMem[aTLAS], size, usage, memProps); err != nil {
		return fmt.Errorf("createAccelerationBuffers:TLAS %w", err)
	}

	// Scratch buffer
	size = vk.DeviceSize(GPUMaxAccelerationScratchBytes)
	usage = vk.BUFFER_USAGE_STORAGE_BUFFER_BIT |
		vk.BUFFER_USAGE_SHADER_DEVICE_ADDRESS_BIT
	if err = vr.createBuffer(&vr.accelMem[aSCRATCH], size, usage, memProps); err != nil {
		return fmt.Errorf("createAccelerationBuffers:SCRATCH %w", err)
	}
	info = vk.BufferDeviceAddressInfo{Buffer: vr.accelMem[aSCRATCH].handle}
	vr.accelMem[aSCRATCH].address = vk.GetBufferDeviceAddress(vr.device, &info)
	return nil
}
func (vr *vulkanRenderer) disposeAccelerationBuffers() {
	if len(vr.accelMem) == aSCRATCH+1 {
		vr.disposeBuffer(&vr.accelMem[aBLAS])
		vr.disposeBuffer(&vr.accelMem[aINST0])
		vr.disposeBuffer(&vr.accelMem[aINST1])
		vr.disposeBuffer(&vr.accelMem[aTLAS])
		vr.disposeBuffer(&vr.accelMem[aSCRATCH])
		vr.accelMem = nil // trigger garbage collect
		vr.accelPtr = nil //   ""
	}
}

// =============================================================================
// resize handling - mostly affects the swapchain and framebuffers.

// resize implements render.renderer
func (vr *vulkanRenderer) resize(width, height uint32) {
	vr.resizeWidth = width   // new width
	vr.resizeHeight = height // new height
	vr.resizesRequested += 1 // request resize
}
func (vr *vulkanRenderer) size() (width, height uint32) {
	return vr.frameWidth, vr.frameHeight
}

// isResizing returns true if there is an outstanding resize.
func (vr *vulkanRenderer) isResizing() bool {
	return vr.resizesRequested != vr.resizesCompleted
}

// resizeSwapchain recreates everything that is affected by a size change.
func (vr *vulkanRenderer) resizeSwapchain() (err error) {
	if vr.recreatingSwapchain {
		return fmt.Errorf("resizeFrame aborted: recreating swapchain")
	}
	if vr.frameWidth == 0 || vr.frameHeight == 0 {
		// require at least 1 pixel in each dimension.
		return fmt.Errorf("resizeFrame aborted: w:%d h:%d", vr.frameWidth, vr.frameHeight)
	}
	vr.recreatingSwapchain = true // stops rendering while recreating swapchain
	// ---
	vk.DeviceWaitIdle(vr.device)   // wait for idle before destroying swapchain
	vr.disposeSwapchainResources() // delete the swapchain.
	surface := surfaceProperties{} // requery swapchain support
	if err = vr.getSurfaceProperties(&surface, vr.physicalDevice); err != nil {
		return err
	}
	vr.frameWidth = uint32(vr.resizeWidth)               // update to new size
	vr.frameHeight = uint32(vr.resizeHeight)             //
	vr.resizeWidth = 0                                   // mark resize as complete
	vr.resizeHeight = 0                                  //
	if err = vr.createSwapchainResources(); err != nil { // recreate swapchain.
		return err
	}
	vr.recreatingSwapchain = false // ok to render
	vr.resizesCompleted = vr.resizesRequested
	slog.Info("vulkan resize complete", "size", fmt.Sprintf("%d:%d", vr.frameWidth, vr.frameHeight))
	return nil
}

// ============================================================================
// buffers allocate GPU data memory and transport data from the CPU to the GPU.
// lock, load, unlock a buffer generally in that order.

type vulkanBuffer struct {
	handle  vk.Buffer
	memory  vk.DeviceMemory
	address vk.DeviceAddress
}

// createBuffer allocates a buffer, memory, and binds the buffer to the memory
func (vr *vulkanRenderer) createBuffer(buff *vulkanBuffer, size vk.DeviceSize,
	usage vk.BufferUsageFlagBits, properties vk.MemoryPropertyFlags) (err error) {

	// create the buffer structure
	buffInfo := vk.BufferCreateInfo{
		Size:        size,
		Usage:       usage,
		SharingMode: vk.SHARING_MODE_EXCLUSIVE, // used in a single queue
	}
	if buff.handle, err = vk.CreateBuffer(vr.device, &buffInfo, nil); err != nil {
		return fmt.Errorf("vk.CreateBuffer: %w", err)
	}

	// allocate the memory needed by the buffer.
	memRequirements := vk.GetBufferMemoryRequirements(vr.device, buff.handle)
	memType, err := vr.findMemoryType(memRequirements.MemoryTypeBits, properties)
	if err != nil {
		return err
	}
	allocateInfo := vk.MemoryAllocateInfo{
		AllocationSize:  memRequirements.Size,
		MemoryTypeIndex: memType,
	}
	if (usage & vk.BUFFER_USAGE_SHADER_DEVICE_ADDRESS_BIT) != 0 {
		allocFlagsInfo := vk.MemoryAllocateFlagsInfo{Flags: vk.MEMORY_ALLOCATE_DEVICE_ADDRESS_BIT}
		allocateInfo.PNext = unsafe.Pointer(allocFlagsInfo.ToVK())
	}
	if buff.memory, err = vk.AllocateMemory(vr.device, &allocateInfo, nil); err != nil {
		return fmt.Errorf("createBuffer:vk.AllocateMemory: %w", err)
	}

	// bind the buffer to the memory
	if err = vk.BindBufferMemory(vr.device, buff.handle, buff.memory, 0); err != nil {
		return fmt.Errorf("vk.BindBufferMemory: %w", err)
	}
	return nil
}

// disposeBuffer returns buffer resources.
func (vr *vulkanRenderer) disposeBuffer(buff *vulkanBuffer) {
	if buff.memory != 0 {
		vk.FreeMemory(vr.device, buff.memory, nil)
		buff.memory = 0
	}
	if buff.handle != 0 {
		vk.DestroyBuffer(vr.device, buff.handle, nil)
		buff.handle = 0
	}
}

// Device local memory is faster than host coherent memory.
// Copying to device local memory is done through the host coherent
// staging buffer and this uses the transfer queue.
func (vr *vulkanRenderer) uploadData(pool vk.CommandPool, queue vk.Queue, handle vk.Buffer, offset uint64, data []byte) {

	// create a staging buffer and load data into the staging buffer.
	var staging vulkanBuffer
	usage := vk.BUFFER_USAGE_TRANSFER_SRC_BIT
	props := vk.MEMORY_PROPERTY_HOST_VISIBLE_BIT | vk.MEMORY_PROPERTY_HOST_COHERENT_BIT
	size := vk.DeviceSize(len(data))
	err := vr.createBuffer(&staging, size, usage, props)
	if err != nil {
		slog.Debug("uploadData:createBuffer", "error", err)
		return
	}
	err = vr.loadCPUBuffer(&staging, 0, data)
	if err != nil {
		slog.Debug("uploadData:loadCPUBuffer", "error", err)
		return
	}

	// copy the data from staging into the given GPU buffer
	vr.copyGPUBuffer(pool, queue, staging.handle, 0, handle, offset, size)
	vr.disposeBuffer(&staging) // clean up staging.
}

// loadCPUBuffer copies data to the CPU visible buffer,
// normally a HOST_VISIBLE, HOST_COHERENT staging buffer.
func (vr *vulkanRenderer) loadCPUBuffer(buff *vulkanBuffer, offset vk.DeviceSize, data []byte) error {
	size := vk.DeviceSize(len(data))
	ptr, err := vk.MapMemory(vr.device, buff.memory, offset, size, 0)
	if err != nil {
		return fmt.Errorf("vk.MapMemory: %w", err)
	}
	vk.MemCopySlice(unsafe.Pointer(ptr), data)
	vk.UnmapMemory(vr.device, buff.memory)
	return nil
}

// copyGPUBuffer copies data from one data buffer to another.
func (vr *vulkanRenderer) copyGPUBuffer(pool vk.CommandPool, queue vk.Queue,
	src vk.Buffer, srcOffset uint64,
	dst vk.Buffer, dstOffset uint64, size vk.DeviceSize) error {

	// create a one-time copy command.
	cmd := vr.beginSingleCommand("copyGPUBuffer", pool)
	region := vk.BufferCopy{SrcOffset: 0, DstOffset: vk.DeviceSize(dstOffset), Size: size}
	vk.CmdCopyBuffer(cmd, src, dst, []vk.BufferCopy{region})

	// submit one-time command for execution and wait for it to complete.
	vr.endSingleCommand("copyGPUBuffer", cmd, pool, vr.graphicsQ)
	return nil
}

// =============================================================================
// mesh vertex and index data are a bunch of bytes that exist
// within larger vulkan buffers
//   mesh vertices are held in each frames vertex buffer
//   mesh indices are held in each frames index buffer

// vulkanMesh tracks mesh data vertex buffer locations.
// Indexed by the load.MeshData types.
// - vertex data is stored in the vr.vertexMem[vData] buffer
// - vertex index data is stored in the vr.vertexMem[vIndex] buffer
type vulkanMesh map[int]vulkanAttrData
type vulkanAttrData struct {
	count  uint32 // number of elements, ie: vertexes, instances
	stride uint32 // number of bytes per element.
	offset uint32 // start location of data in the buffer.
}

// loadMeshes stores mesh data in GPU buffers.
// Buffers match the GLTF import data format.
// Immutable once uploaded. There is only one vertex buffer for
// all frames. Updating a mesh means adding a new one and refering to it
// in future draw calls.
//
// FUTURE: handle buffer (de/re)allocates using linked lists.
func (vr *vulkanRenderer) loadMeshes(meshes []load.MeshData) (mids []uint32, err error) {

	// cosolidate the upload data into a buffer for each vertex data type.
	data := [][]byte{[]byte{}, []byte{}} // vData, vIndex temp upload buffers.
	maxData := []uint32{GPUMaxVertexBytes, GPUMaxVIndexBytes}
	dataOffsets := []uint32{vr.vertexPtr[vData], vr.vertexPtr[vIndex]}
	for _, msh := range meshes {

		// track each mesh with a corresponding vulkan-mesh
		vmsh := map[int]vulkanAttrData{}
		for i := range load.VertexTypes {
			dataType := vData
			if i == load.Indexes {
				dataType = vIndex
			}
			if msh[i].Count > 0 {
				v := vulkanAttrData{
					count:  msh[i].Count,
					stride: msh[i].Stride,
					offset: dataOffsets[dataType],
				}
				vmsh[i] = v

				// check if the data exceeds what was allocated.
				if uint32(v.offset+v.count*v.stride) > maxData[dataType] {
					return mids, fmt.Errorf("loadMeshes:insufficient memory %d", dataType)
				}

				// consolidate the upload data into temp buffers.
				data[dataType] = append(data[dataType], msh[i].Data...)

				// update the last mesh data offset.
				dataOffsets[dataType] += v.stride * v.count
			}
		}
		mids = append(mids, uint32(len(vr.meshes))) // mesh ID for the new mesh.
		vr.meshes = append(vr.meshes, vmsh)
	}

	// upload the consolidated data once for each data type.
	for dataType := range []int{vData, vIndex} {
		uploadData := data[dataType]
		if len(uploadData) <= 0 {
			return mids, fmt.Errorf("loadMeshes:missing position and/or index data %d", dataType)
		}
		buff := &vr.vertexMem[dataType]
		offset := uint64(vr.vertexPtr[dataType])
		vr.uploadData(vr.graphicsQCmdPool, vr.graphicsQ, buff.handle, offset, uploadData)
		vr.vertexPtr[dataType] += uint32(len(uploadData))

		// track the total amount of uploaded bytes
		// to a global that the app can access.
		GPUTotalMeshBytes += uint32(len(uploadData))
	}
	if RaytraceEnabled {
		// create the per mesh BLAS acceleration data
		// so the mesh can be used in a ray trace/query shader.
		for _, mid := range mids {
			vr.createBLAS(mid)
		}
	}
	return mids, nil
}

// loadInstances adds instance data to an existing mesh.
// Immutable once uploaded.
func (vr *vulkanRenderer) loadInstances(mid uint32, idata []load.Buffer) (err error) {
	if int(mid) >= len(vr.meshes) {
		return fmt.Errorf("loadInstances invalid mesh ID: %d", mid)
	}

	// use the existing vulkan mesh for this model
	uploadData := []byte{} // temp upload buffers.
	vmsh := vr.meshes[mid]
	newDataIndex := vr.vertexPtr[vData]
	for i := range load.VertexTypes {
		if i == load.Indexes || idata[i].Count <= 0 {
			continue
		}
		v := vulkanAttrData{
			count:  idata[i].Count,
			stride: idata[i].Stride,
			offset: newDataIndex,
		}
		vmsh[i] = v

		// check if the data exceeds what was allocated.
		if uint32(v.offset+v.count*v.stride) > GPUMaxVertexBytes {
			return fmt.Errorf("loadInstances:insufficient memory %d", vData)
		}

		// consolidate the upload data into temp buffers.
		uploadData = append(uploadData, idata[i].Data...)

		// update the last mesh data offset.
		newDataIndex += v.stride * v.count
	}

	// upload the index data once
	buff := &vr.vertexMem[vData]
	offset := uint64(vr.vertexPtr[vData])
	vr.uploadData(vr.graphicsQCmdPool, vr.graphicsQ, buff.handle, offset, uploadData)
	vr.vertexPtr[vData] += uint32(len(uploadData))

	// track the total amount of uploaded bytes
	// to a global that the app can access.
	GPUTotalMeshBytes += uint32(len(uploadData))
	return nil
}

// FUTURE: handle deallocates using linked lists.
// For now never deallocate so that the lastMesh is always valid.
func (vr *vulkanRenderer) dropMesh(mid uint32) {}

// updateMesh : see docs on render:UpdateMesh
func (vr *vulkanRenderer) updateMesh(mid uint32, data []load.Buffer) (err error) {
	if int(mid) >= len(vr.meshes) {
		return fmt.Errorf("updateMesh invalid ID: %d", mid)
	}
	msh := vr.meshes[mid]

	// check that the buffer data matches.
	for i := range load.VertexTypes {
		if msh[i].count != data[i].Count {
			return fmt.Errorf("updateMesh count mismatch %d %d", msh[i].count, data[i].Count)
		}
		if msh[i].stride != data[i].Stride {
			return fmt.Errorf("updateMesh stride mismatch %d %d", msh[i].stride, data[i].Stride)
		}
	}

	// everything matches, so re-upload data.
	for i := range load.VertexTypes {
		if data[i].Count > 0 {
			dataType := vData
			if i == load.Indexes {
				dataType = vIndex
			}

			// upload data - existing data size is the same: count, stride, offset
			buff := &vr.vertexMem[dataType]
			offset := uint64(msh[i].offset)
			vr.uploadData(vr.graphicsQCmdPool, vr.graphicsQ, buff.handle, offset, data[i].Data)
		}
	}
	return nil // everything ok.
}

// drawMesh
func (vr *vulkanRenderer) drawMesh(frame *vulkanFrame, mid uint32, attrs []load.ShaderAttribute, isInstance bool, instanceCount uint32) {
	if mid < 0 || mid >= uint32(len(vr.meshes)) {
		slog.Error("invalid mesh ID", "mid", mid)
		return
	}
	vmsh := vr.meshes[mid]

	// bind the mesh vertex attribute data expected by the shader.
	// The attribute must match one of the supported vertex types.
	buffs := []vk.Buffer{}
	offsets := []vk.DeviceSize{}
	for _, attr := range attrs {
		i := attr.AType
		if i < 0 || i >= load.Indexes {
			slog.Error("unsupported vertex attribute", "attribute_type", attr.AType)
			continue
		}
		if _, ok := vmsh[i]; !ok {
			slog.Error("vertex attribute mismatch", "attribute_type", attr.AType)
			continue
		}
		buffs = append(buffs, vr.vertexMem[vData].handle)
		offsets = append(offsets, vk.DeviceSize(vmsh[i].offset))
	}
	vk.CmdBindVertexBuffers(frame.cmds, 0, buffs, offsets)

	// bind the vertex index data.
	ibuff := vr.vertexMem[vIndex].handle
	ioffset := vk.DeviceSize(vmsh[load.Indexes].offset)
	vk.CmdBindIndexBuffer(frame.cmds, ibuff, ioffset, vk.INDEX_TYPE_UINT16)

	// draw the mesh
	if isInstance {
		vk.CmdDrawIndexed(frame.cmds, vmsh[load.Indexes].count, instanceCount, 0, 0, 0)
	} else {
		vk.CmdDrawIndexed(frame.cmds, vmsh[load.Indexes].count, 1, 0, 0, 0)
	}
}

// =============================================================================
// image utilities

// vulkanImage holds data needed for creating and disposing of vulkan images.
// Can be a texture image or a depth buffer, etc.
type vulkanImage struct {
	handle vk.Image
	view   vk.ImageView
	memory vk.DeviceMemory
	width  uint32
	height uint32
}

// create a vkCreateImage
func (vr *vulkanRenderer) createImage(img *vulkanImage,
	format vk.Format,
	usage vk.ImageUsageFlags,
	memoryFlags vk.MemoryPropertyFlags) (err error) {

	// create the requested image
	imgInfo := vk.ImageCreateInfo{
		ImageType: vk.IMAGE_TYPE_2D,
		Extent: vk.Extent3D{
			Width:  img.width,
			Height: img.height,
			Depth:  1,
		},
		MipLevels:     1,
		ArrayLayers:   1,
		Format:        format,
		Tiling:        vk.IMAGE_TILING_OPTIMAL,
		InitialLayout: vk.IMAGE_LAYOUT_UNDEFINED,
		Usage:         usage,
		Samples:       vk.SAMPLE_COUNT_1_BIT,
		SharingMode:   vk.SHARING_MODE_EXCLUSIVE,
	}
	img.handle, err = vk.CreateImage(vr.device, &imgInfo, nil)
	if err != nil {
		return fmt.Errorf("vk.CreateImage: %w", err)
	}

	// check if the required image memory exists.
	memReqs := vk.GetImageMemoryRequirements(vr.device, img.handle)
	memType, err := vr.findMemoryType(memReqs.MemoryTypeBits, memoryFlags)
	if err != nil {
		return err
	}

	// allocate the image memory
	img.memory, err = vk.AllocateMemory(vr.device,
		&vk.MemoryAllocateInfo{
			AllocationSize:  memReqs.Size,
			MemoryTypeIndex: uint32(memType),
		}, nil)
	if err != nil {
		return fmt.Errorf("vk.AllocateMemory: %w", err)
	}
	err = vk.BindImageMemory(vr.device, img.handle, img.memory, 0)
	if err != nil {
		return fmt.Errorf("vk.BindImageMemory: %w", err)
	}
	return nil
}

// create a vk.CreateImageView for the given image
func (vr *vulkanRenderer) createImageView(img vk.Image, format vk.Format, aspectFlags vk.ImageAspectFlags) (view vk.ImageView, err error) {
	createInfo := vk.ImageViewCreateInfo{
		Image:    img,
		ViewType: vk.IMAGE_VIEW_TYPE_2D,
		Format:   format,
		SubresourceRange: vk.ImageSubresourceRange{
			AspectMask:     aspectFlags,
			BaseMipLevel:   0,
			LevelCount:     1,
			BaseArrayLayer: 0,
			LayerCount:     1,
		},
	}
	return vk.CreateImageView(vr.device, &createInfo, nil)
}

// transitionImageLayout switches image layout,
// see use in loadTexture.
func (vr *vulkanRenderer) transitionImageLayout(img *vulkanImage, format vk.Format, oldLayout vk.ImageLayout, newLayout vk.ImageLayout) {
	cmd := vr.beginSingleCommand("transitionImageLayout", vr.graphicsQCmdPool)
	barrier := vk.ImageMemoryBarrier{
		SrcAccessMask:       0,
		DstAccessMask:       0,
		OldLayout:           oldLayout,
		NewLayout:           newLayout,
		SrcQueueFamilyIndex: vk.QUEUE_FAMILY_IGNORED,
		DstQueueFamilyIndex: vk.QUEUE_FAMILY_IGNORED,
		Image:               img.handle,
		SubresourceRange: vk.ImageSubresourceRange{
			AspectMask:     vk.IMAGE_ASPECT_COLOR_BIT,
			BaseMipLevel:   0,
			LevelCount:     1,
			BaseArrayLayer: 0,
			LayerCount:     1,
		},
	}
	var sourceStage vk.PipelineStageFlags
	var destinationStage vk.PipelineStageFlags
	switch {
	case oldLayout == vk.IMAGE_LAYOUT_UNDEFINED && newLayout == vk.IMAGE_LAYOUT_TRANSFER_DST_OPTIMAL:
		barrier.SrcAccessMask = 0
		barrier.DstAccessMask = vk.ACCESS_TRANSFER_WRITE_BIT
		sourceStage = vk.PIPELINE_STAGE_TOP_OF_PIPE_BIT
		destinationStage = vk.PIPELINE_STAGE_TRANSFER_BIT
	case oldLayout == vk.IMAGE_LAYOUT_TRANSFER_DST_OPTIMAL && newLayout == vk.IMAGE_LAYOUT_SHADER_READ_ONLY_OPTIMAL:
		barrier.SrcAccessMask = vk.ACCESS_TRANSFER_WRITE_BIT
		barrier.DstAccessMask = vk.ACCESS_SHADER_READ_BIT
		sourceStage = vk.PIPELINE_STAGE_TRANSFER_BIT
		destinationStage = vk.PIPELINE_STAGE_FRAGMENT_SHADER_BIT
	default:
		slog.Error("unsupported layout transition!")
	}
	vk.CmdPipelineBarrier(cmd, sourceStage, destinationStage, 0, nil, nil, []vk.ImageMemoryBarrier{barrier})
	vr.endSingleCommand("transitionImageLayout", cmd, vr.graphicsQCmdPool, vr.graphicsQ)
}

func (vr *vulkanRenderer) copyBufferToImage(buffer *vulkanBuffer, img *vulkanImage) {
	cmd := vr.beginSingleCommand("copyBufferToImage", vr.graphicsQCmdPool)
	region := vk.BufferImageCopy{
		BufferOffset:      0,
		BufferRowLength:   0,
		BufferImageHeight: 0,
		ImageSubresource: vk.ImageSubresourceLayers{
			AspectMask:     vk.IMAGE_ASPECT_COLOR_BIT,
			MipLevel:       0,
			BaseArrayLayer: 0,
			LayerCount:     1,
		},
		ImageOffset: vk.Offset3D{X: 0, Y: 0, Z: 0},
		ImageExtent: vk.Extent3D{Width: img.width, Height: img.height, Depth: 1},
	}
	vk.CmdCopyBufferToImage(cmd, buffer.handle, img.handle, vk.IMAGE_LAYOUT_TRANSFER_DST_OPTIMAL, []vk.BufferImageCopy{region})
	vr.endSingleCommand("copyBufferToImage", cmd, vr.graphicsQCmdPool, vr.graphicsQ)
}

func (vr *vulkanRenderer) disposeImage(img *vulkanImage) {
	if img.view != 0 {
		vk.DestroyImageView(vr.device, img.view, nil)
		img.view = 0
	}
	if img.memory != 0 {
		vk.FreeMemory(vr.device, img.memory, nil)
		img.memory = 0
	}
	if img.handle != 0 {
		vk.DestroyImage(vr.device, img.handle, nil)
		img.handle = 0
	}
}

func (vr *vulkanRenderer) findMemoryType(filter uint32, memFlags vk.MemoryPropertyFlags) (memoryType uint32, err error) {
	memProps := vk.GetPhysicalDeviceMemoryProperties(vr.physicalDevice)
	for i := uint32(0); i < memProps.MemoryTypeCount; i++ {
		flagsMatch := memProps.MemoryTypes[i].PropertyFlags&memFlags == memFlags
		if (filter&1<<i != 0) && flagsMatch {
			return i, nil
		}
	}
	return 0, fmt.Errorf("failed to find appropriate memory type.")
}

// =============================================================================
// texture combines an image and a sampler.

// vulkanTexture is a vulkanImage plus an image sampler.
type vulkanTexture struct {
	image   vulkanImage
	sampler vk.Sampler
}

// loadTexture stores image data in a GPU buffer
// and creates a corresponding texture sampler.
//
// CURRENT: immutable once uploaded. Updating a texture means adding a new
// texture and refering to it in future draw calls.
func (vr *vulkanRenderer) loadTexture(w, h uint32, pixels []byte) (tid uint32, err error) {
	vr.textures = append(vr.textures, vulkanTexture{})
	tid = uint32(len(vr.textures) - 1)
	tex := &vr.textures[tid]

	// put image data into staging buffer
	imageSize := vk.DeviceSize(len(pixels))     //
	GPUTotalTextureBytes += uint32(len(pixels)) // track the total amount of uploaded bytes
	stagingBuffer := vulkanBuffer{}
	err = vr.createBuffer(&stagingBuffer, imageSize, vk.BUFFER_USAGE_TRANSFER_SRC_BIT,
		vk.MEMORY_PROPERTY_HOST_VISIBLE_BIT|vk.MEMORY_PROPERTY_HOST_COHERENT_BIT)
	if err != nil {
		vr.disposeBuffer(&stagingBuffer)
		return 0, err
	}
	if err = vr.loadCPUBuffer(&stagingBuffer, 0, pixels); err != nil {
		vr.disposeBuffer(&stagingBuffer)
		return 0, err
	}

	// upload GPU image
	format := vk.FORMAT_R8G8B8A8_SRGB
	tex.image.width = w
	tex.image.height = h
	err = vr.createImage(&tex.image, format,
		vk.IMAGE_USAGE_TRANSFER_DST_BIT|vk.IMAGE_USAGE_SAMPLED_BIT,
		vk.MEMORY_PROPERTY_DEVICE_LOCAL_BIT)
	if err != nil {
		vr.disposeBuffer(&stagingBuffer)
		return 0, err
	}
	vr.transitionImageLayout(&tex.image, format, vk.IMAGE_LAYOUT_UNDEFINED, vk.IMAGE_LAYOUT_TRANSFER_DST_OPTIMAL)
	vr.copyBufferToImage(&stagingBuffer, &tex.image)
	vr.transitionImageLayout(&tex.image, format, vk.IMAGE_LAYOUT_TRANSFER_DST_OPTIMAL, vk.IMAGE_LAYOUT_SHADER_READ_ONLY_OPTIMAL)
	vr.disposeBuffer(&stagingBuffer)

	// create the texture view
	tex.image.view, err = vr.createImageView(tex.image.handle, format, vk.IMAGE_ASPECT_COLOR_BIT)
	if err != nil {
		return 0, err
	}

	// create a sampler. These are immutable and can be shared
	// by different shaders and pipelines.
	devProps := vk.GetPhysicalDeviceProperties(vr.physicalDevice)
	samplerInfo := vk.SamplerCreateInfo{
		MagFilter:               vk.FILTER_LINEAR,
		MinFilter:               vk.FILTER_LINEAR,
		AddressModeU:            vk.SAMPLER_ADDRESS_MODE_REPEAT,
		AddressModeV:            vk.SAMPLER_ADDRESS_MODE_REPEAT,
		AddressModeW:            vk.SAMPLER_ADDRESS_MODE_REPEAT,
		AnisotropyEnable:        true,
		MaxAnisotropy:           devProps.Limits.MaxSamplerAnisotropy,
		BorderColor:             vk.BORDER_COLOR_INT_OPAQUE_BLACK,
		UnnormalizedCoordinates: false,
		CompareEnable:           false,
		CompareOp:               vk.COMPARE_OP_ALWAYS,
		MipmapMode:              vk.SAMPLER_MIPMAP_MODE_LINEAR,
	}
	tex.sampler, err = vk.CreateSampler(vr.device, &samplerInfo, nil)
	if err != nil {
		return 0, fmt.Errorf("vk.CreateSampler: %w", err)
	}
	return tid, nil
}
func (vr *vulkanRenderer) dropTexture(tid uint32) {
	if tid >= uint32(len(vr.textures)) {
		slog.Error("invalid texture ID", "tid", tid)
	}
	tex := &vr.textures[tid]
	vr.disposeImage(&tex.image)
	if tex.sampler != 0 {
		vk.DestroySampler(vr.device, tex.sampler, nil)
		tex.sampler = 0
	}
}

// updateTexture : see docs on render:UpdateTexture
func (vr *vulkanRenderer) updateTexture(tid, width, height uint32, pixels []byte) (err error) {
	if tid >= uint32(len(vr.textures)) {
		return fmt.Errorf("updateTexture invalid texture ID %d", tid)
	}
	tex := vr.textures[tid]
	if tex.image.width != width || tex.image.height != height {
		return fmt.Errorf("updateTexture expected image size %d:%d got %d:%d",
			tex.image.width, tex.image.height, width, height)
	}

	// put image data into staging buffer
	imageSize := vk.DeviceSize(len(pixels))
	stagingBuffer := vulkanBuffer{}
	err = vr.createBuffer(&stagingBuffer, imageSize, vk.BUFFER_USAGE_TRANSFER_SRC_BIT,
		vk.MEMORY_PROPERTY_HOST_VISIBLE_BIT|vk.MEMORY_PROPERTY_HOST_COHERENT_BIT)
	if err != nil {
		vr.disposeBuffer(&stagingBuffer)
		return fmt.Errorf("updateTexture create staging %w", err)
	}
	if err = vr.loadCPUBuffer(&stagingBuffer, 0, pixels); err != nil {
		vr.disposeBuffer(&stagingBuffer)
		return fmt.Errorf("updateTexture upload staging %w", err)
	}

	// upload GPU image
	format := vk.FORMAT_R8G8B8A8_SRGB
	vr.transitionImageLayout(&tex.image, format, vk.IMAGE_LAYOUT_UNDEFINED, vk.IMAGE_LAYOUT_TRANSFER_DST_OPTIMAL)
	vr.copyBufferToImage(&stagingBuffer, &tex.image)
	vr.transitionImageLayout(&tex.image, format, vk.IMAGE_LAYOUT_TRANSFER_DST_OPTIMAL, vk.IMAGE_LAYOUT_SHADER_READ_ONLY_OPTIMAL)
	vr.disposeBuffer(&stagingBuffer)
	return nil
}

// =============================================================================
// acceleration structures for ray query and ray trace.

// accelStruct used for both BLAS and TLAS.
type accelStruct struct {
	handle  vk.AccelerationStructureKHR // handle for the structure.
	address vk.DeviceAddress            // address to the storage buffer.
}

// createAccelerationStruct creates the struct and stores it in the
// appropriate pre-allocated storage buffer.
func (vr *vulkanRenderer) createAccelerationStruct(accel *accelStruct,
	accelType int, structSize vk.DeviceSize) (err error) {

	// map to the vulkan accel struct type.
	var structType vk.AccelerationStructureTypeKHR
	switch accelType {
	case aBLAS:
		structType = vk.ACCELERATION_STRUCTURE_TYPE_BOTTOM_LEVEL_KHR
	case aTLAS:
		structType = vk.ACCELERATION_STRUCTURE_TYPE_TOP_LEVEL_KHR
	default:
		slog.Error("createAccelerationStruct: invalid type", "type", accelType)
		return // developer error - should have been caught during testing.
	}

	// create and store the acceleration struct.
	createInfo := vk.AccelerationStructureCreateInfoKHR{
		Buffer: vr.accelMem[accelType].handle, // one buffer for BLAS, one for TLAS
		Offset: vr.accelPtr[accelType],        // append at end of storage buffer.
		Size:   structSize,
		Typ:    structType,
	}
	accel.handle, err = vk.CreateAccelerationStructureKHR(vr.device, &createInfo, nil)
	if err != nil {
		return fmt.Errorf("createAccelerationStruct: %w", err)
	}
	addressInfo := vk.AccelerationStructureDeviceAddressInfoKHR{AccelerationStructure: accel.handle}
	accel.address = vk.GetAccelerationStructureDeviceAddressKHR(vr.device, &addressInfo)

	// track where the next acceleration structure needs to be placed.
	// vulkan requires the struct offset to be on a 256 byte boundary.
	alignment := uint64(256)                                 // offset boundary requirement
	endOffset := uint64(vr.accelPtr[accelType] + structSize) // offset to end of data...
	offset := (endOffset + alignment - 1) & ^(alignment - 1) // ...plus alignment padding...
	vr.accelPtr[accelType] = vk.DeviceSize(offset)           // ...to get required offset.
	return nil
}

// disposeAccelerationStructs destroys all the acceleration structs
// create by the above createAccelerationStruct. Expected to be
// called once on shutdown.
func (vr *vulkanRenderer) disposeAccelerationStructs() {
	for _, tlas := range vr.accelScenes {
		vk.DestroyAccelerationStructureKHR(vr.device, tlas.handle, nil)
	}
	for _, blas := range vr.accelMeshes {
		vk.DestroyAccelerationStructureKHR(vr.device, blas.handle, nil)
	}
	clear(vr.accelScenes)
	clear(vr.accelMeshes)
}

// createBottomLevelAccelerationStruct for each mesh.
// FUTURE: build for all vr.meshes, ie: wait until all meshes are ready.
func (vr *vulkanRenderer) createBLAS(mid uint32) (err error) {
	mesh := vr.meshes[mid]
	vpos := mesh[load.Vertexes] // vertex buffer position
	ipos := mesh[load.Indexes]  // index buffer position
	numTriangles := ipos.count / 3
	if vpos.stride != 12 {
		return // ignore 2D meshes (vertexes are either 2D or 3D)
	}

	// create the triangle data for the acceleration structure
	vertexBufferDeviceAddress := vk.DeviceOrHostAddressConstKHR{}
	vertexBufferDeviceAddress.AsDeviceAddress(vr.vertexMem[vData].address + vk.DeviceAddress(vpos.offset))
	indexBufferDeviceAddress := vk.DeviceOrHostAddressConstKHR{}
	indexBufferDeviceAddress.AsDeviceAddress(vr.vertexMem[vIndex].address + vk.DeviceAddress(ipos.offset))
	triangles := vk.AccelerationStructureGeometryTrianglesDataKHR{
		VertexFormat: vk.FORMAT_R32G32B32_SFLOAT,
		VertexData:   vertexBufferDeviceAddress,
		MaxVertex:    vpos.count - 1,
		VertexStride: vk.DeviceSize(vpos.stride),
		IndexType:    vk.INDEX_TYPE_UINT16,
		IndexData:    indexBufferDeviceAddress,
	}
	geometry := vk.AccelerationStructureGeometryDataKHR{}
	geometry.AsTriangles(triangles)
	accelerationGeometry := vk.AccelerationStructureGeometryKHR{
		Flags:        vk.GEOMETRY_OPAQUE_BIT_KHR,
		GeometryType: vk.GEOMETRY_TYPE_TRIANGLES_KHR,
		Geometry:     geometry,
	}
	geomInfo := vk.AccelerationStructureBuildGeometryInfoKHR{
		Typ:         vk.ACCELERATION_STRUCTURE_TYPE_BOTTOM_LEVEL_KHR,
		Mode:        vk.BUILD_ACCELERATION_STRUCTURE_MODE_BUILD_KHR,
		Flags:       vk.BUILD_ACCELERATION_STRUCTURE_PREFER_FAST_TRACE_BIT_KHR,
		PGeometries: []vk.AccelerationStructureGeometryKHR{accelerationGeometry},
	}

	// Get sizes needed for building the acceleration structure.
	buildSizes := vk.GetAccelerationStructureBuildSizesKHR(
		vr.device,
		vk.ACCELERATION_STRUCTURE_BUILD_TYPE_DEVICE_KHR,
		&geomInfo,
		[]uint32{numTriangles},
	)
	if uint32(buildSizes.BuildScratchSize) > GPUMaxAccelerationScratchBytes {
		slog.Error("Need bigger BLAS scratch", "scratch_size", buildSizes.BuildScratchSize)
	}

	// create the BLAS acceleration structure
	blas := accelStruct{}
	structSize := buildSizes.AccelerationStructureSize
	err = vr.createAccelerationStruct(&blas, aBLAS, structSize)
	if err != nil {
		return fmt.Errorf("create BLAS struct: %w", err)
	}
	vr.accelMeshes[mid] = blas // keep reference for TLAS and shutdown.

	// add the information needed to build the acceleration struct
	scratchDeviceAddress := vk.DeviceOrHostAddressKHR{}
	scratchDeviceAddress.AsDeviceAddress(vr.accelMem[aSCRATCH].address)
	geomInfo.DstAccelerationStructure = blas.handle
	geomInfo.ScratchData = scratchDeviceAddress

	// add number of triangles to the build range.
	rangeInfo := vk.AccelerationStructureBuildRangeInfoKHR{
		PrimitiveCount:  numTriangles,
		PrimitiveOffset: 0, // 0 offset from triangles.IndexData defined above.
		FirstVertex:     0, // ignored as indicies are used in triangles above.
		TransformOffset: 0, // not used for BLAS.
	}

	// build the acceleration structure on the device via a one-time command buffer submission
	// The structure is built into the allocate memory.
	cmd := vr.beginSingleCommand("createBLAS", vr.graphicsQCmdPool)
	vk.CmdBuildAccelerationStructuresKHR(
		cmd,
		[]vk.AccelerationStructureBuildGeometryInfoKHR{geomInfo},
		[]*vk.AccelerationStructureBuildRangeInfoKHR{&rangeInfo},
	)
	vr.endSingleCommand("createBLAS", cmd, vr.graphicsQCmdPool, vr.graphicsQ)
	return nil
}

// (acceleration) ainstances are lazy created and reused each frame.
var ainstances = []vk.AccelerationStructureInstanceKHR{}

// check each frame-in-flight for any TLAS changes.
// Compare using the EID's of what was rendered this frame vs last frame.
var instanceEIDs = [][]uint32{[]uint32{}, []uint32{}}

// getAccelInstance returns an acceleration instances data structure,
// creating a new one if necessary, otherwise reusing an existing.
func (vr *vulkanRenderer) getAccelInst() (ai *vk.AccelerationStructureInstanceKHR) {
	size := len(ainstances)
	switch {
	case size == cap(ainstances):
		ainstances = append(ainstances, vk.AccelerationStructureInstanceKHR{})
	case size < cap(ainstances):
		ainstances = ainstances[:size+1] // reuse previously allocated.
	}
	return &ainstances[size]
}

// updateTopLevelAccelerationStruct
func (vr *vulkanRenderer) updateTLAS(pass Pass) (err error) {
	ainstances = ainstances[:0] // reset keeping memory
	previousEIDs := append([]uint32{}, instanceEIDs[vr.frameIndex]...)
	instanceEIDs[vr.frameIndex] = instanceEIDs[vr.frameIndex][:0]
	for _, packet := range pass.Packets {
		if packet.IsInstanced {
			continue // ignore instanced meshes.
		}
		eid := packet.EID    // identifies a mesh instance.
		mid := packet.MeshID // identifies a BLAS
		blas, ok := vr.accelMeshes[mid]
		if !ok {
			// should be a BLAS since the BLAS is created with the mesh.
			slog.Error("updateTLAS missing BLAS", "mid", mid, "eid", eid)
			return
		}
		instance := vr.getAccelInst() // get the next accel instance.
		instanceEIDs[vr.frameIndex] = append(instanceEIDs[vr.frameIndex], eid)

		// get the model transform data and...
		mat4x4Bytes := packet.Uniforms[load.MODEL]
		m := (*m4)(unsafe.Pointer(&mat4x4Bytes[0]))
		// ...set the acceleration instance transform.
		tm := &instance.Transform.Matrix
		tm[0][0], tm[0][1], tm[0][2], tm[0][3] = m.xx, m.yx, m.zx, m.wx
		tm[1][0], tm[1][1], tm[1][2], tm[1][3] = m.xy, m.yy, m.zy, m.wy
		tm[2][0], tm[2][1], tm[2][2], tm[2][3] = m.xz, m.yz, m.zz, m.wz

		// create the instance acceleration struct
		// manually handle the bitfield masking (not part of vulkan bindings yet)
		maskBitfield := uint32(0xFF << 24)                                                      // | InstanceCustomIndex
		flagBitfield := uint32(vk.GEOMETRY_INSTANCE_TRIANGLE_FACING_CULL_DISABLE_BIT_KHR << 24) // | InstanceShaderBindingTableRecordOffset
		instance.Mask_InstanceCustomIndex = maskBitfield
		instance.Flags_InstanceShaderBindingTableRecordOffset = flagBitfield
		instance.AccelerationStructureReference = uint64(blas.address)
	}
	if len(ainstances) <= 0 {
		return // no model instances, nothing to do.
	}

	// update TLAS if there are no changes from last time.
	// Any instance changes require a rebuild.
	updateTLAS := slices.Equal(instanceEIDs[vr.frameIndex], previousEIDs)

	// upload all the instance structs to the TLAS Instance storage buffer.
	// overwriting the contents of the previous instance data buffer.
	// Use per frame buffers for instance data so that it is safe to update.
	// NOTE: instance is byte compatible with instance.ToVK()
	size := len(ainstances) * int(unsafe.Sizeof(ainstances[0]))
	data := unsafe.Slice((*byte)(unsafe.Pointer(&ainstances[0])), size)
	var buff vulkanBuffer
	switch vr.frameIndex {
	case 0:
		buff = vr.accelMem[aINST0]                // acceleration instance buffer.
		vr.accelPtr[aINST0] = vk.DeviceSize(size) // end of data
	default:
		buff = vr.accelMem[aINST1]
		vr.accelPtr[aINST1] = vk.DeviceSize(size) // end of data
	}
	vr.uploadData(vr.graphicsQCmdPool, vr.graphicsQ, buff.handle, 0, data)

	// reference the current accleration instance buffer.
	instanceDataAddress := vk.DeviceOrHostAddressConstKHR{}
	instanceDataAddress.AsDeviceAddress(buff.address)
	instanceData := vk.AccelerationStructureGeometryInstancesDataKHR{
		ArrayOfPointers: false,               // device address to array of instance structs.
		Data:            instanceDataAddress, // device address of instance buffer
	}

	// create geometry referencing the acceleration instance data.
	geometry := vk.AccelerationStructureGeometryDataKHR{}
	geometry.AsInstances(instanceData)
	accelerationGeometry := vk.AccelerationStructureGeometryKHR{
		Flags:        vk.GEOMETRY_OPAQUE_BIT_KHR,
		GeometryType: vk.GEOMETRY_TYPE_INSTANCES_KHR,
		Geometry:     geometry,
	}
	buildInfo := vk.AccelerationStructureBuildGeometryInfoKHR{
		Typ: vk.ACCELERATION_STRUCTURE_TYPE_TOP_LEVEL_KHR,
		Flags: vk.BUILD_ACCELERATION_STRUCTURE_PREFER_FAST_BUILD_BIT_KHR |
			vk.BUILD_ACCELERATION_STRUCTURE_ALLOW_UPDATE_BIT_KHR,
		Mode:        vk.BUILD_ACCELERATION_STRUCTURE_MODE_BUILD_KHR,
		PGeometries: []vk.AccelerationStructureGeometryKHR{accelerationGeometry},
	}

	// check if this the first time or an update.
	tlas, tlasExists := vr.accelScenes[0]
	if tlasExists && updateTLAS {
		buildInfo.Mode = vk.BUILD_ACCELERATION_STRUCTURE_MODE_UPDATE_KHR
	}

	// get the TLAS build sizes.
	buildSizes := vk.GetAccelerationStructureBuildSizesKHR(
		vr.device,
		vk.ACCELERATION_STRUCTURE_BUILD_TYPE_DEVICE_KHR,
		&buildInfo,
		[]uint32{1})
	if uint32(buildSizes.BuildScratchSize) > GPUMaxAccelerationScratchBytes {
		slog.Error("Need bigger TLAS scratch", "scratch_size", buildSizes.BuildScratchSize)
	}
	if uint32(buildSizes.AccelerationStructureSize) > GPUMaxAccelerationTLASBytes {
		slog.Error("Need bigger TLAS buffer", "buffer_size", buildSizes.AccelerationStructureSize)
	}

	// if this is a first time build then create the first TLAS structure
	// with extra space for updates.
	if !tlasExists {
		tlas = accelStruct{}
		vr.createAccelerationStruct(&tlas, aTLAS, vk.DeviceSize(GPUMaxAccelerationTLASBytes))
		vr.accelScenes[0] = tlas // save tlas reference to destroy it on shutdown.
	}

	// create the TLAS build information.
	scratchDeviceAddress := vk.DeviceOrHostAddressKHR{}
	scratchDeviceAddress.AsDeviceAddress(vr.accelMem[aSCRATCH].address)
	if tlasExists && updateTLAS {
		buildInfo.SrcAccelerationStructure = tlas.handle
	}
	buildInfo.DstAccelerationStructure = tlas.handle
	buildInfo.ScratchData = scratchDeviceAddress
	rangeInfo := vk.AccelerationStructureBuildRangeInfoKHR{
		PrimitiveCount:  uint32(len(ainstances)), // number of instance structs.
		PrimitiveOffset: 0,
		FirstVertex:     0,
		TransformOffset: 0,
	}

	// Build the acceleration structure on the device via a one-time command buffer submission
	cmd := vr.beginSingleCommand("updateTLAS", vr.graphicsQCmdPool)
	if tlasExists {
		preBarrier := vk.MemoryBarrier{
			SrcAccessMask: vk.ACCESS_ACCELERATION_STRUCTURE_WRITE_BIT_KHR |
				vk.ACCESS_TRANSFER_WRITE_BIT | vk.ACCESS_SHADER_READ_BIT,
			DstAccessMask: vk.ACCESS_ACCELERATION_STRUCTURE_READ_BIT_KHR |
				vk.ACCESS_ACCELERATION_STRUCTURE_WRITE_BIT_KHR,
		}
		vk.CmdPipelineBarrier(cmd,
			vk.PIPELINE_STAGE_ACCELERATION_STRUCTURE_BUILD_BIT_KHR|
				vk.PIPELINE_STAGE_TRANSFER_BIT|
				vk.PIPELINE_STAGE_FRAGMENT_SHADER_BIT, // srcStageMask
			vk.PIPELINE_STAGE_ACCELERATION_STRUCTURE_BUILD_BIT_KHR, // dstStageMask
			0, // dependencyFlags
			[]vk.MemoryBarrier{preBarrier},
			nil,
			nil,
		)
	}
	// update or rebuild depending on instance changes above.
	vk.CmdBuildAccelerationStructuresKHR(
		cmd,
		[]vk.AccelerationStructureBuildGeometryInfoKHR{buildInfo},
		[]*vk.AccelerationStructureBuildRangeInfoKHR{&rangeInfo},
	)
	if tlasExists {
		postBarrier := vk.MemoryBarrier{
			SrcAccessMask: vk.ACCESS_ACCELERATION_STRUCTURE_WRITE_BIT_KHR,
			DstAccessMask: vk.ACCESS_ACCELERATION_STRUCTURE_READ_BIT_KHR |
				vk.ACCESS_SHADER_READ_BIT,
		}
		vk.CmdPipelineBarrier(cmd,
			vk.PIPELINE_STAGE_ACCELERATION_STRUCTURE_BUILD_BIT_KHR, // srcStageMask
			vk.PIPELINE_STAGE_ACCELERATION_STRUCTURE_BUILD_BIT_KHR|
				vk.PIPELINE_STAGE_FRAGMENT_SHADER_BIT, // dstStageMask
			0, // dependencyFlags
			[]vk.MemoryBarrier{postBarrier},
			nil,
			nil,
		)
	}
	vr.endSingleCommand("updateTLAS", cmd, vr.graphicsQCmdPool, vr.graphicsQ)
	return nil
}

// getBufferDeviceAddres returns the device address for the given buffer.
func (vr *vulkanRenderer) getBufferDeviceAddress(buff vk.Buffer) (addr vk.DeviceAddress) {
	addrInfo := vk.BufferDeviceAddressInfo{Buffer: buff}
	return vk.GetBufferDeviceAddress(vr.device, &addrInfo)
}

// =============================================================================
// shaders are GPU programs.

// vulkanShader groups shader related information including the
// different shader stage modules, pipeline, and descriptor layouts.
type vulkanShader struct {
	name       string            // for debug
	pipe       vk.Pipeline       //
	pipeLayout vk.PipelineLayout //

	// shader vertex attributes
	attrs []load.ShaderAttribute

	// uniform information to help create and update descriptor sets.
	usets    uniformSets // shader uniform information
	usesTLAS bool        // true when shader uses a top level accel struct.

	// scene uniform data buffers for all scene uniform descriptors
	// indexed by vr.imageIndex
	sceneUniforms    vulkanBuffer // scene scope uniform data.
	sceneUniformsMap *byte        // unsafe pointer to uniform mapped memory

	// track the number of unique material instances for this shader.
	maxSamplers   uint32          // maximum samplers supported by shader.
	samplers      []vulkanSampler // track unique sampler (sets) with...
	nextSamplerID uint32          // ... a sampler ID.

	// descriptors for scene and material uniform data.
	sceneLayout    vk.DescriptorSetLayout // scene uniforms per renderpass
	samplLayout    vk.DescriptorSetLayout // sample uniforms per shader.
	descriptorPool vk.DescriptorPool      // uniforms and samplers
	sceneDescSets  []vk.DescriptorSet     // one per image.
	sceneUpdated   []bool                 // true if descriptor set updated.
}

// vulkanSampler tracks existing resources to help reuse descriptor sets.
type vulkanSampler struct {
	samplerSet     []uint32           // unique set of samplers.
	descriptorSets []vk.DescriptorSet // one per surface image
	updated        []bool             // one per surface image - set when updated.
}

// map the shader stage bit flags to vulkan stage bit flags.
var vulkanStages = map[load.ShaderStage]vk.ShaderStageFlagBits{
	load.Stage_VERTEX:   vk.SHADER_STAGE_VERTEX_BIT,
	load.Stage_GEOMETRY: vk.SHADER_STAGE_GEOMETRY_BIT,
	load.Stage_FRAGMENT: vk.SHADER_STAGE_FRAGMENT_BIT,
}

// map the load data types to VK_FORMAT
var vulkanDataFormats = map[load.ShaderDataType]vk.Format{
	load.Type_FLOAT: vk.FORMAT_R32_SFLOAT,
	load.Type_VEC2:  vk.FORMAT_R32G32_SFLOAT,
	load.Type_VEC3:  vk.FORMAT_R32G32B32_SFLOAT,
	load.Type_VEC4:  vk.FORMAT_R32G32B32A32_SFLOAT,
}

// loadShader creates a shader and corresponding pipeline based
// on the given shader configuration.
func (vr *vulkanRenderer) loadShader(config *load.Shader) (sid uint16, err error) {
	sid = uint16(len(vr.shaders)) // the shader ID if loadShader succeeds.
	shader := vulkanShader{name: config.Name}
	shader.attrs = append(shader.attrs, config.Attrs...)
	shader.usets = getUniformSets(config.Uniforms)
	shader.maxSamplers = 256 // FUTURE: make configurable.
	vr.createShaderUniformBuffers(&shader)

	// stages for this shader. Module field is set by loadShaderModules
	stages := []vk.PipelineShaderStageCreateInfo{}
	switch config.Stages {
	case load.Stage_VERTEX | load.Stage_FRAGMENT:
		stages = append(stages, vk.PipelineShaderStageCreateInfo{Stage: vk.SHADER_STAGE_VERTEX_BIT, PName: "main"})
		stages = append(stages, vk.PipelineShaderStageCreateInfo{Stage: vk.SHADER_STAGE_FRAGMENT_BIT, PName: "main"})
	default:
		vr.disposeShader(&shader)
		return 0, fmt.Errorf("unsupported shader stages %d", config.Stages)
	}

	// load the modules for the given shader stages
	// shader modules can be released once the pipeline is created.
	err = vr.loadShaderModules(&shader, config.Name, stages)
	if err != nil {
		vr.disposeShader(&shader)
		return 0, err
	}
	// release the modules now that the shader has been created.
	for i := range stages {
		defer vk.DestroyShaderModule(vr.device, stages[i].Module, nil)
	}

	// non-interleaved vertex attribute descriptions
	vertexAttrDescriptions := make([]vk.VertexInputAttributeDescription, len(config.Attrs))
	vertexBindingDescriptions := make([]vk.VertexInputBindingDescription, len(config.Attrs))
	if len(config.Attrs) > 0 {
		for i, attr := range config.Attrs {
			// attribute descriptions
			vertexAttrDescriptions[i].Location = uint32(i)
			vertexAttrDescriptions[i].Binding = uint32(i)
			vertexAttrDescriptions[i].Format = vulkanDataFormats[attr.DType]
			vertexAttrDescriptions[i].Offset = 0

			// binding descriptions
			vertexBindingDescriptions[i].Binding = uint32(i)
			vertexBindingDescriptions[i].Stride = attr.Stride
			if attr.IsInstanced {
				vertexBindingDescriptions[i].InputRate = vk.VERTEX_INPUT_RATE_INSTANCE
			} else {
				vertexBindingDescriptions[i].InputRate = vk.VERTEX_INPUT_RATE_VERTEX
			}
		}
	}
	vertexInputInfo := &vk.PipelineVertexInputStateCreateInfo{
		PVertexBindingDescriptions:   vertexBindingDescriptions,
		PVertexAttributeDescriptions: vertexAttrDescriptions,
	}

	// every shader expects scene uniforms.
	scenes := config.GetSceneUniforms()
	if len(scenes) <= 0 {
		vr.disposeShader(&shader)
		return 0, fmt.Errorf("expecting scene uniforms")
	}

	// create the scene descriptor layout.
	accelUniforms := config.GetAccelerationUniforms()
	accelCount := uint32(len(accelUniforms)) // expecting 0 or 1
	shader.usesTLAS = accelCount > 0         // shader requires TLAS
	if accelCount > 1 {
		// developer to check why more than one would be needed.
		slog.Error("support for one TLAS per scene", "count", accelCount)
	}
	bindings := []vk.DescriptorSetLayoutBinding{
		vk.DescriptorSetLayoutBinding{
			Binding:         0,
			DescriptorType:  vk.DESCRIPTOR_TYPE_UNIFORM_BUFFER,
			DescriptorCount: 1,
			StageFlags:      vk.SHADER_STAGE_VERTEX_BIT | vk.SHADER_STAGE_FRAGMENT_BIT,
		},
		vk.DescriptorSetLayoutBinding{
			Binding:         1,
			DescriptorType:  vk.DESCRIPTOR_TYPE_ACCELERATION_STRUCTURE_KHR,
			DescriptorCount: accelCount,
			StageFlags:      vk.SHADER_STAGE_FRAGMENT_BIT,
		},
	}
	setLayout := vk.DescriptorSetLayoutCreateInfo{PBindings: bindings}
	shader.sceneLayout, err = vk.CreateDescriptorSetLayout(vr.device, &setLayout, nil)
	if err != nil {
		slog.Error("invalid scene layout", "err", err)
	}

	// allocate the sampler descriptor set layout if applicable.
	samplers := config.GetSamplerUniforms()
	if len(samplers) > 0 {
		bindings := []vk.DescriptorSetLayoutBinding{}
		binding := vk.DescriptorSetLayoutBinding{
			Binding:         0,
			DescriptorType:  vk.DESCRIPTOR_TYPE_COMBINED_IMAGE_SAMPLER,
			DescriptorCount: uint32(len(samplers)),
			StageFlags:      vk.SHADER_STAGE_FRAGMENT_BIT,
		}
		bindings = append(bindings, binding)
		shader.samplLayout, err = vk.CreateDescriptorSetLayout(
			vr.device, &vk.DescriptorSetLayoutCreateInfo{PBindings: bindings}, nil)
		if err != nil {
			slog.Error("invalid sampler layout", "err", err)
		}
	}

	// create one descriptor pool for descriptor sets in this shader.
	// FUTURE: create descriptor pools for groups of similar shaders,
	//         or one large pool for all shaders.
	// FUTURE: rationalize the descriptorCount sizes.
	shader.descriptorPool, err = vk.CreateDescriptorPool(vr.device,
		&vk.DescriptorPoolCreateInfo{
			MaxSets: 3 + 3*shader.maxSamplers + 3,
			PPoolSizes: []vk.DescriptorPoolSize{
				{Typ: vk.DESCRIPTOR_TYPE_UNIFORM_BUFFER, DescriptorCount: 1024},
				{Typ: vk.DESCRIPTOR_TYPE_ACCELERATION_STRUCTURE_KHR, DescriptorCount: 1024},
				{Typ: vk.DESCRIPTOR_TYPE_COMBINED_IMAGE_SAMPLER, DescriptorCount: 4096},
			},
			Flags: vk.DESCRIPTOR_POOL_CREATE_FREE_DESCRIPTOR_SET_BIT,
			// FUTURE: | vk.DESCRIPTOR_POOL_CREATE_UPDATE_AFTER_BIND_BIT;
		}, nil)
	if err != nil {
		vr.disposeShader(&shader)
		return 0, err
	}

	// allocate scene descriptor sets, one per image
	if shader.sceneLayout != 0 {
		allocLayouts := []vk.DescriptorSetLayout{}
		for i := 0; i < int(vr.imageCount); i++ {
			allocLayouts = append(allocLayouts, shader.sceneLayout)
		}
		shader.sceneDescSets, err = vk.AllocateDescriptorSets(vr.device,
			&vk.DescriptorSetAllocateInfo{
				DescriptorPool: shader.descriptorPool,
				PSetLayouts:    allocLayouts,
			})
		if err != nil {
			vr.disposeShader(&shader)
			return 0, err
		}
		shader.sceneUpdated = make([]bool, len(shader.sceneDescSets))
	}

	// allocate sampler descriptor sets.
	if shader.samplLayout != 0 && len(samplers) > 0 {

		// allocate three descriptor sets (one per surface image)...
		allocLayouts := []vk.DescriptorSetLayout{}
		for i := 0; i < int(vr.imageCount); i++ {
			allocLayouts = append(allocLayouts, shader.samplLayout)
		}
		// ...for each sampler.
		shader.samplers = make([]vulkanSampler, shader.maxSamplers)
		for i := uint32(0); i < shader.maxSamplers; i++ {
			shader.samplers[i].updated = make([]bool, vr.imageCount)
			shader.samplers[i].descriptorSets, err = vk.AllocateDescriptorSets(vr.device,
				&vk.DescriptorSetAllocateInfo{
					DescriptorPool: shader.descriptorPool,
					PSetLayouts:    allocLayouts,
				})
			if err != nil {
				vr.disposeShader(&shader)
				return 0, err
			}
		}
	}

	// create the pipeline layout
	pipelineLayouts := []vk.DescriptorSetLayout{}
	if shader.sceneLayout != 0 {
		pipelineLayouts = append(pipelineLayouts, shader.sceneLayout)
	}
	if shader.samplLayout != 0 {
		pipelineLayouts = append(pipelineLayouts, shader.samplLayout)
	}
	layoutInfo := vk.PipelineLayoutCreateInfo{
		PSetLayouts: pipelineLayouts,
	}

	// setup pipeline push constants if they exist
	if shader.usets.pushSize > 0 {
		push_constant := vk.PushConstantRange{
			Offset:     0,
			Size:       maxPushConstantBytes, // allocate the max 128 bytes.
			StageFlags: vk.SHADER_STAGE_VERTEX_BIT | vk.SHADER_STAGE_FRAGMENT_BIT,
		}
		layoutInfo.PPushConstantRanges = []vk.PushConstantRange{push_constant}
	}
	shader.pipeLayout, err = vk.CreatePipelineLayout(vr.device, &layoutInfo, nil)
	if err != nil {
		vr.disposeShader(&shader)
		return 0, fmt.Errorf("vk.CreatePipelineLayout: %w", err)
	}

	// all engine shaders expect triangles
	inputAssembly := vk.PipelineInputAssemblyStateCreateInfo{
		Topology:               vk.PRIMITIVE_TOPOLOGY_TRIANGLE_LIST,
		PrimitiveRestartEnable: false,
	}
	if config.DrawLines {
		inputAssembly.Topology = vk.PRIMITIVE_TOPOLOGY_LINE_LIST
	}

	// viewport and scissor are set as dynamic state later.
	vr.setViewportAndScissor()
	viewportState := vk.PipelineViewportStateCreateInfo{
		PViewports: []vk.Viewport{vr.viewport},
		PScissors:  []vk.Rect2D{vr.scissor},
	}

	// rasterization state can be changed by the shader config.
	rasterizerState := vk.PipelineRasterizationStateCreateInfo{
		DepthClampEnable:        false,
		RasterizerDiscardEnable: false,
		PolygonMode:             vk.POLYGON_MODE_FILL,
		LineWidth:               1.0,
		CullMode:                vk.CULL_MODE_BACK_BIT, // default
		FrontFace:               vk.FRONT_FACE_COUNTER_CLOCKWISE,
		DepthBiasEnable:         false,
	}
	if config.CullModeNone {
		rasterizerState.CullMode = vk.CULL_MODE_NONE
	}

	// multisampling
	multisampleState := vk.PipelineMultisampleStateCreateInfo{
		SampleShadingEnable:  false,
		RasterizationSamples: vk.SAMPLE_COUNT_1_BIT,
		MinSampleShading:     1.0,
	}

	// depth and stencil testing
	depthStencil := vk.PipelineDepthStencilStateCreateInfo{
		DepthTestEnable:       true,
		DepthWriteEnable:      true,
		DepthCompareOp:        vk.COMPARE_OP_LESS_OR_EQUAL,
		DepthBoundsTestEnable: false,
		StencilTestEnable:     false,
		Front:                 vk.StencilOpState{},
		Back:                  vk.StencilOpState{},
		MinDepthBounds:        0,
		MaxDepthBounds:        1.0,
	}

	// describe how colors are written to the render image.
	colorBlend := vk.PipelineColorBlendStateCreateInfo{
		LogicOpEnable: false,
		LogicOp:       vk.LOGIC_OP_COPY,
		PAttachments: []vk.PipelineColorBlendAttachmentState{
			vk.PipelineColorBlendAttachmentState{
				BlendEnable:         true, // transparent geometry
				ColorWriteMask:      vk.COLOR_COMPONENT_R_BIT | vk.COLOR_COMPONENT_G_BIT | vk.COLOR_COMPONENT_B_BIT | vk.COLOR_COMPONENT_A_BIT,
				SrcColorBlendFactor: vk.BLEND_FACTOR_SRC_ALPHA,
				DstColorBlendFactor: vk.BLEND_FACTOR_ONE_MINUS_SRC_ALPHA,
				ColorBlendOp:        vk.BLEND_OP_ADD,

				// blend the alpha values
				SrcAlphaBlendFactor: vk.BLEND_FACTOR_SRC_ALPHA,
				DstAlphaBlendFactor: vk.BLEND_FACTOR_ONE_MINUS_SRC_ALPHA,
				AlphaBlendOp:        vk.BLEND_OP_ADD,

				// ... or... set the final alpha to the new source alpha
				// SrcAlphaBlendFactor: vk.BLEND_FACTOR_ONE,
				// DstAlphaBlendFactor: vk.BLEND_FACTOR_ZERO,
				// AlphaBlendOp:        vk.BLEND_OP_ADD,
			},
		},
	}

	// pipeline dynamic state can be changed without re-creating the pipeline
	// Helps when resizing windows. Viewport and Scissor are expected to be set
	// in each command buffer.
	dynamicStateInfo := vk.PipelineDynamicStateCreateInfo{
		PDynamicStates: []vk.DynamicState{
			vk.DYNAMIC_STATE_VIEWPORT,
			vk.DYNAMIC_STATE_SCISSOR,
		},
	}

	// use dynamic rendering (replaces render passes and frame buffers)
	pipelineRenderingInfo := vk.PipelineRenderingCreateInfo{
		PColorAttachmentFormats: []vk.Format{vr.surfaceFormat.Format},
		DepthAttachmentFormat:   vr.depthFormat,
	}

	// create the pipeline.
	pipelineInfo := vk.GraphicsPipelineCreateInfo{
		PNext:               unsafe.Pointer(pipelineRenderingInfo.ToVK()),
		PStages:             stages,
		PVertexInputState:   vertexInputInfo,
		PInputAssemblyState: &inputAssembly,
		PViewportState:      &viewportState,
		PRasterizationState: &rasterizerState,
		PMultisampleState:   &multisampleState,
		PDepthStencilState:  &depthStencil,
		PColorBlendState:    &colorBlend,
		PDynamicState:       &dynamicStateInfo,
		PTessellationState:  nil,
		Layout:              shader.pipeLayout,
		Subpass:             0,
		BasePipelineHandle:  0,
		BasePipelineIndex:   -1,
	}
	pipelines, err := vk.CreateGraphicsPipelines(vr.device, 0, []vk.GraphicsPipelineCreateInfo{pipelineInfo}, nil)
	if err != nil {
		vr.disposeShader(&shader)
		return 0, fmt.Errorf("vk.CreateGraphicsPipeline: %w", err)
	}
	shader.pipe = pipelines[0]

	// success... add the shader to the list of loaded shaders.
	vr.shaders = append(vr.shaders, shader)
	return sid, nil
}

// loadShader loads the .spv and creates a shader module for each shader stage
// using a file name convention for locating shader .spv files.
func (vr *vulkanRenderer) loadShaderModules(shader *vulkanShader, name string, stages []vk.PipelineShaderStageCreateInfo) (err error) {
	supportedStages := map[vk.ShaderStageFlagBits]string{
		vk.SHADER_STAGE_VERTEX_BIT:   "vert",
		vk.SHADER_STAGE_FRAGMENT_BIT: "frag",
	}
	for i, stage := range stages {
		stageName, ok := supportedStages[stage.Stage]
		if !ok {
			return fmt.Errorf("unsupported shader stage %d", stage.Stage)
		}

		// load the shader module bytes using the
		// shader generation naming convention.
		filename := fmt.Sprintf("%s_%s.spv", name, stageName)
		shaderCode, err := load.ShaderBytes(filename)
		if err != nil {
			return fmt.Errorf("load shader: %w", err)
		}

		// create the shader module.
		stages[i].Module, err = vk.CreateShaderModule(vr.device,
			&vk.ShaderModuleCreateInfo{
				CodeSize: uintptr(len(shaderCode)),
				PCode:    (*uint32)(unsafe.Pointer(&shaderCode[0])),
			}, nil)
		if err != nil {
			return fmt.Errorf("vk.CreateShaderModule: %w", err)
		}
	}
	return nil
}

// dropShader drops the shader for the given shader ID.
func (vr *vulkanRenderer) dropShader(sid uint16) {
	if sid < 0 || sid >= uint16(len(vr.shaders)) {
		slog.Error("dropShader:invalid shader ID", "sid", sid)
		// also not allowed to drop the default shader.
	}
	vr.disposeShader(&vr.shaders[sid])
}

// disposeShader releases shader resources.
func (vr *vulkanRenderer) disposeShader(s *vulkanShader) {
	if s.descriptorPool != 0 {
		// destroying the pool also destroys the descriptorSets
		vk.DestroyDescriptorPool(vr.device, s.descriptorPool, nil)
		s.descriptorPool = 0
	}
	vr.disposeShaderUniformBuffers(s)
	if s.samplLayout != 0 {
		vk.DestroyDescriptorSetLayout(vr.device, s.samplLayout, nil)
		s.samplLayout = 0
	}
	if s.sceneLayout != 0 {
		vk.DestroyDescriptorSetLayout(vr.device, s.sceneLayout, nil)
		s.sceneLayout = 0
	}
	if s.pipe != 0 {
		vk.DestroyPipeline(vr.device, s.pipe, nil)
		s.pipe = 0
	}
	if s.pipeLayout != 0 {
		vk.DestroyPipelineLayout(vr.device, s.pipeLayout, nil)
		s.pipeLayout = 0
	}
}

// createShaderUniformBuffers creates the uniform data buffers for this shader
func (vr *vulkanRenderer) createShaderUniformBuffers(s *vulkanShader) (err error) {

	// check if the device supports device local buffers.
	deviceLocalBits := vk.MemoryPropertyFlags(0)
	if vr.deviceLocalHostVisible {
		// add device local if supported by the physical device.
		deviceLocalBits = vk.MEMORY_PROPERTY_DEVICE_LOCAL_BIT
	}
	props := vk.MEMORY_PROPERTY_HOST_VISIBLE_BIT | vk.MEMORY_PROPERTY_HOST_COHERENT_BIT | deviceLocalBits
	numImages := uint32(len(vr.images))

	// create enough scene uniform data buffer space for each surface image
	// map the uniform memory once for the lifetime of the app.
	bufferSize := vk.DeviceSize(maxSceneUniformBytes * numImages)
	usage := vk.BUFFER_USAGE_UNIFORM_BUFFER_BIT
	err = vr.createBuffer(&s.sceneUniforms, bufferSize, usage, props)
	if err != nil {
		return fmt.Errorf("sceneUniforms:vk.createBuffer: %w", err)
	}
	s.sceneUniformsMap, err = vk.MapMemory(vr.device, s.sceneUniforms.memory, 0, bufferSize, 0)
	if err != nil {
		return fmt.Errorf("sceneUniformsMap:vk.MapMemory: %w", err)
	}
	return nil
}

// disposeShaderUniformBuffers diposes resources allocated in createShaderUniformBuffers.
func (vr *vulkanRenderer) disposeShaderUniformBuffers(s *vulkanShader) {
	if s != nil {
		vr.disposeBuffer(&s.sceneUniforms)
		s.sceneUniformsMap = nil
	}
}

// applySceneUniforms updates the scene descriptor sets to point
// to the scene uniform data buffer
func (vr *vulkanRenderer) applySceneUniforms(shader *vulkanShader) {
	if shader.sceneLayout == 0 {
		slog.Error("applySceneUniforms: no scene uniforms", "shader", shader.name)
		return
	}
	sceneSet := shader.sceneDescSets[vr.imageIndex]
	if !shader.sceneUpdated[vr.imageIndex] {
		offset := vk.DeviceSize(vr.imageIndex * maxSceneUniformBytes)
		descWrites := []vk.WriteDescriptorSet{
			{
				DstSet:          sceneSet,
				DstBinding:      0,
				DstArrayElement: 0,
				DescriptorType:  vk.DESCRIPTOR_TYPE_UNIFORM_BUFFER,
				PBufferInfo: []vk.DescriptorBufferInfo{
					{
						Buffer: shader.sceneUniforms.handle,
						Offset: offset,
						Rang:   vk.DeviceSize(maxSceneUniformBytes),
					},
				},
			},
		}

		// update the TLAS when the shader expects an acceleration structure
		if shader.usesTLAS {
			tlas, tlasExists := vr.accelScenes[0]
			if tlasExists {

				// As this isn't part of Vulkan's core, information is passed via pNext chaining
				accelWriteSet := vk.WriteDescriptorSetAccelerationStructureKHR{}
				accelWriteSet.PAccelerationStructures = []vk.AccelerationStructureKHR{tlas.handle}
				descWrites = append(descWrites, vk.WriteDescriptorSet{
					PNext:           unsafe.Pointer(accelWriteSet.ToVK()),
					DstSet:          sceneSet,
					DstBinding:      1,
					DstArrayElement: 0,
					DescriptorType:  vk.DESCRIPTOR_TYPE_ACCELERATION_STRUCTURE_KHR,
				})
			}
		}
		vk.UpdateDescriptorSets(vr.device, descWrites, nil)
		shader.sceneUpdated[vr.imageIndex] = true
	}
	firstSet := uint32(0) // scene is always set=0
	dsets := []vk.DescriptorSet{sceneSet}
	frame := vr.frames[vr.frameIndex]
	vk.CmdBindDescriptorSets(frame.cmds, vk.PIPELINE_BIND_POINT_GRAPHICS, shader.pipeLayout, firstSet, dsets, nil)
}

// setSamplers sets all the samplers expected by this shader.
// The samplers are in the order expected by the shader config.
func (vr *vulkanRenderer) setSamplers(shader *vulkanShader, tids []uint32) (matID uint32, err error) {
	// compare texture IDs to the existing samplers
	for i := uint32(0); i < shader.nextSamplerID; i++ {
		if slices.Compare(tids, shader.samplers[i].samplerSet) == 0 {
			return i, nil // match: reuse existing sampler.
		}
	}

	// create a new sampler for these textures on this shader.
	matID = shader.nextSamplerID
	if matID >= shader.maxSamplers {
		return 0, fmt.Errorf("setSamplers:max samplers exceeded:%d", matID)
	}
	shader.samplers[matID].samplerSet = append(shader.samplers[matID].samplerSet, tids...)
	shader.nextSamplerID += 1
	return matID, nil
}

// applySamplerUniforms updates the material scope descriptor sets
// to point to the proper material uniform data buffer
func (vr *vulkanRenderer) applySamplerUniforms(shader *vulkanShader, sampID uint32) {
	if shader.samplLayout == 0 {
		slog.Error("applySamplersUniforms: no sampler uniforms", "shader", shader.name)
		return
	}
	if sampID >= uint32(len(shader.samplers)) {
		slog.Error("applySamplersUniforms:invalid sampler ID", "ID", sampID)
		return
	}
	sampler := &shader.samplers[sampID]
	descriptorSet := sampler.descriptorSets[vr.imageIndex]
	setNum := uint32(1) // sampler uniforms are set=1

	// check if the sampler descriptor set needs updating.
	if !sampler.updated[vr.imageIndex] {
		descriptorSetWrites := []vk.WriteDescriptorSet{}
		descriptorIndex := uint32(0)

		// check for samplers
		if len(sampler.samplerSet) > 0 {
			samplerInfo := []vk.DescriptorImageInfo{}
			for _, tid := range sampler.samplerSet {
				t := vr.textures[tid]
				samplerInfo = append(samplerInfo, vk.DescriptorImageInfo{
					ImageLayout: vk.IMAGE_LAYOUT_SHADER_READ_ONLY_OPTIMAL,
					ImageView:   t.image.view,
					Sampler:     t.sampler,
				})
			}
			descriptorSetWrites = append(descriptorSetWrites, vk.WriteDescriptorSet{
				DstSet:         descriptorSet,
				DstBinding:     descriptorIndex,
				DescriptorType: vk.DESCRIPTOR_TYPE_COMBINED_IMAGE_SAMPLER,
				PImageInfo:     samplerInfo,
			})
		}
		vk.UpdateDescriptorSets(vr.device, descriptorSetWrites, nil)
		sampler.updated[vr.imageIndex] = true
	}
	frame := &vr.frames[vr.frameIndex]
	dsets := []vk.DescriptorSet{descriptorSet}
	vk.CmdBindDescriptorSets(frame.cmds, vk.PIPELINE_BIND_POINT_GRAPHICS, shader.pipeLayout, setNum, dsets, nil)
}

// =============================================================================
// command utilities

// convenience create for command buffer that is used once.
func (vr *vulkanRenderer) beginSingleCommand(fn string, pool vk.CommandPool) (cmd vk.CommandBuffer) {
	commands, err := vk.AllocateCommandBuffers(vr.device,
		&vk.CommandBufferAllocateInfo{
			CommandPool:        pool,
			Level:              vk.COMMAND_BUFFER_LEVEL_PRIMARY,
			CommandBufferCount: 1,
		})
	if err != nil {
		slog.Error("beginSingleCommand:vk.AllocateCommandBuffers", "fn", fn, "error", err)
		return cmd
	}
	cmd = commands[0]
	info := vk.CommandBufferBeginInfo{Flags: vk.COMMAND_BUFFER_USAGE_ONE_TIME_SUBMIT_BIT}
	err = vk.BeginCommandBuffer(cmd, &info)
	if err != nil {
		slog.Error("beginSingleCommand:vk.CommandBufferBeginInfo", "fn", fn, "error", err)
	}
	return cmd
}

// convenience end for command buffer that is used once.
func (vr *vulkanRenderer) endSingleCommand(fn string, cmd vk.CommandBuffer, pool vk.CommandPool, queue vk.Queue) {
	if cmd == 0 {
		slog.Error("endSingleCommand: invalid command buffer", "fn", fn)
		return
	}
	vk.EndCommandBuffer(cmd)

	// submit the command buffer
	submitInfo := vk.SubmitInfo{PCommandBuffers: []vk.CommandBuffer{cmd}}
	if err := vk.QueueSubmit(queue, []vk.SubmitInfo{submitInfo}, 0); err != nil {
		slog.Error("endSingleCommand:vk.QueueSubmit", "fn", fn, "error", err)
		return
	}

	// wait for submit to finish....and then free the command buffer
	if err := vk.QueueWaitIdle(queue); err != nil {
		slog.Error("endSingleCommand:vk.QueueWaitIdle", "fn", fn, "error", err)
		return
	}
	vk.FreeCommandBuffers(vr.device, pool, []vk.CommandBuffer{cmd})
}

// =============================================================================
// Rendering a frame

// vulkanFrame holds the per-frame data.
// Frames rotate through rendering as follows:
// : frame record  - CPU records commands
// : frame render  - GPU processes commands to create image
// : frame present - display image to user
//
// Multiple frames (at least 2) allow CPU to be generating a frame while
// the GPU is processing another. Synchonization ensures each frames data
// set is accessed by one process (cpu, gpu, display) at a time.
type vulkanFrame struct {
	cmds vk.CommandBuffer // graphics queue commands.

	// frame render synchonization.
	imageAvailable vk.Semaphore // done presenting, ready for rendering.
	fence          vk.Fence     // render to frames not in use by GPU.
}

func (vr *vulkanRenderer) createRenderFrames() (err error) {
	vr.frameIndex = 0 // start rendering at frame zero
	vr.frames = make([]vulkanFrame, vr.frameCount)
	for i := range vr.frames {
		if err = vr.createFrameCommandBuffer(&vr.frames[i]); err != nil {
			return err
		}
		if err = vr.createFrameSyncronization(&vr.frames[i]); err != nil {
			return err
		}
	}
	return nil
}
func (vr *vulkanRenderer) disposeRenderFrames() {
	for i := range vr.frames {
		if vr.frames[i].imageAvailable != 0 {
			vk.DestroySemaphore(vr.device, vr.frames[i].imageAvailable, nil)
			vr.frames[i].imageAvailable = 0
		}
		if vr.frames[i].fence != 0 {
			vk.DestroyFence(vr.device, vr.frames[i].fence, nil)
			vr.frames[i].fence = 0
		}
		if vr.frames[i].cmds != 0 {
			vk.FreeCommandBuffers(vr.device, vr.graphicsQCmdPool, []vk.CommandBuffer{vr.frames[i].cmds})
			vr.frames[i].cmds = 0
		}
	}
}

// createFrameCommandBuffer ensures there is one graphics command buffer
// for each swapchain image.
func (vr *vulkanRenderer) createFrameCommandBuffer(fr *vulkanFrame) (err error) {
	buffInfo := vk.CommandBufferAllocateInfo{
		CommandPool:        vr.graphicsQCmdPool,
		Level:              vk.COMMAND_BUFFER_LEVEL_PRIMARY,
		CommandBufferCount: 1,
	}
	commands, err := vk.AllocateCommandBuffers(vr.device, &buffInfo)
	if err != nil || len(commands) != 1 {
		return fmt.Errorf("vk.AllocateCommandBuffers: %w", err)
	}
	fr.cmds = commands[0]
	return nil
}

// createFrameSyncronization creates the semaphores and fences needed to coordinate rendering.
func (vr *vulkanRenderer) createFrameSyncronization(fr *vulkanFrame) (err error) {
	fr.imageAvailable, err = vk.CreateSemaphore(vr.device, &vk.SemaphoreCreateInfo{}, nil)
	if err != nil {
		return fmt.Errorf("vk.CreateSemaphore.1: %w", err)
	}

	// Create the fence in a signaled state, indicating that the first frame has
	// already been "rendered". This will prevent the application from waiting
	// indefinitely for the first frame to render since it cannot be rendered
	// until a frame is "rendered" before it.
	fenceInfo := vk.FenceCreateInfo{Flags: vk.FENCE_CREATE_SIGNALED_BIT}
	fr.fence, err = vk.CreateFence(vr.device, &fenceInfo, nil)
	if err != nil {
		return fmt.Errorf("vk.CreateFence: %w", err)
	}
	return nil
}

// =============================================================================
// beginFrame, render objects, endFrame

func (vr *vulkanRenderer) setClearColor(r, g, b, a float32) {
	vr.clearColor[0], vr.clearColor[1], vr.clearColor[2], vr.clearColor[3] = r, g, b, a
}

func (vr *vulkanRenderer) beginFrame(dt time.Duration) (err error) {
	if vr.recreatingSwapchain {
		if err := vk.DeviceWaitIdle(vr.device); err != nil {
			return fmt.Errorf("vk.DeviceWaitIdle1: %w", err)
		}
		return fmt.Errorf("beginFrame aborted: recreating swapchain")
	}

	// handle screen resizes.
	if vr.isResizing() {
		if err := vk.DeviceWaitIdle(vr.device); err != nil {
			return fmt.Errorf("vk.DeviceWaitIdle.2: %w", err)
		}
		if err := vr.resizeSwapchain(); err != nil {
			return fmt.Errorf("recreateSwapchain: %w", err)
		}
		return fmt.Errorf("beginFrame aborted: resized swapchain")
	}

	// wait for current frame to complete - fence starts in signalled state for first frame,
	// and is later signalled when the frame is finished.
	frame := &vr.frames[vr.frameIndex]
	err = vk.WaitForFences(vr.device, []vk.Fence{frame.fence}, true, waitFrame)
	if err != nil {
		return fmt.Errorf("beginFrame aborted: vk.WaitForFences: %w", err)
	}

	// reset the frames fence to unsignalled
	if err := vk.ResetFences(vr.device, []vk.Fence{frame.fence}); err != nil {
		return fmt.Errorf("vk.ResetFences: %w", err)
	}

	// acquire the next image from the swapchain.
	// Pass in the semaphore to be signalled when image is available again.
	vr.imageIndex, err = vk.AcquireNextImageKHR(vr.device, vr.swapchain, maxTimeout, frame.imageAvailable, 0)
	if err != nil {
		if err == vk.SUBOPTIMAL_KHR || err == vk.ERROR_OUT_OF_DATE_KHR {
			return nil // wait for resize
		}
		return fmt.Errorf("beginFrame aborted: vk.AcquireNextImageKHR: %w", err)
	}
	return nil // next is recordFrame
}

// drawFrame issues all the draw commands needed to render all the
// render passes for this frame.
func (vr *vulkanRenderer) drawFrame(passes []Pass) (err error) {
	frame := &vr.frames[vr.frameIndex]

	// clear the command buffer and start recording commands.
	vk.ResetCommandBuffer(frame.cmds, 0) // reset
	cmdInfo := &vk.CommandBufferBeginInfo{Flags: vk.COMMAND_BUFFER_USAGE_ONE_TIME_SUBMIT_BIT}
	err = vk.BeginCommandBuffer(frame.cmds, cmdInfo)
	if err != nil {
		return fmt.Errorf("vk.BeginCommandBuffer %w", err)
	}

	// transition the image layouts to what is needed for drawing a frame.
	drawBarriers := []vk.ImageMemoryBarrier2{
		vk.ImageMemoryBarrier2{
			SrcStageMask:  vk.PIPELINE_STAGE_2_COLOR_ATTACHMENT_OUTPUT_BIT,
			SrcAccessMask: 0,
			DstStageMask:  vk.PIPELINE_STAGE_2_COLOR_ATTACHMENT_OUTPUT_BIT,
			DstAccessMask: vk.ACCESS_2_COLOR_ATTACHMENT_READ_BIT | vk.ACCESS_2_COLOR_ATTACHMENT_WRITE_BIT,
			OldLayout:     vk.IMAGE_LAYOUT_UNDEFINED,
			NewLayout:     vk.IMAGE_LAYOUT_ATTACHMENT_OPTIMAL,
			Image:         vr.images[vr.imageIndex],
			SubresourceRange: vk.ImageSubresourceRange{
				AspectMask: vk.IMAGE_ASPECT_COLOR_BIT,
				LevelCount: 1, LayerCount: 1},
		},
		vk.ImageMemoryBarrier2{
			SrcStageMask:  vk.PIPELINE_STAGE_2_LATE_FRAGMENT_TESTS_BIT,
			SrcAccessMask: vk.ACCESS_2_DEPTH_STENCIL_ATTACHMENT_WRITE_BIT,
			DstStageMask:  vk.PIPELINE_STAGE_2_EARLY_FRAGMENT_TESTS_BIT,
			DstAccessMask: vk.ACCESS_2_DEPTH_STENCIL_ATTACHMENT_WRITE_BIT,
			OldLayout:     vk.IMAGE_LAYOUT_UNDEFINED,
			NewLayout:     vk.IMAGE_LAYOUT_ATTACHMENT_OPTIMAL,
			Image:         vr.depthImage.handle,
			SubresourceRange: vk.ImageSubresourceRange{
				AspectMask: vk.IMAGE_ASPECT_DEPTH_BIT,
				LevelCount: 1, LayerCount: 1},
		},
	}
	drawBarrierInfo := vk.DependencyInfo{PImageMemoryBarriers: drawBarriers}
	vk.CmdPipelineBarrier2(frame.cmds, &drawBarrierInfo)

	// set pipeline dynamic state
	vr.setViewportAndScissor()
	vk.CmdSetViewport(frame.cmds, 0, []vk.Viewport{vr.viewport})
	vk.CmdSetScissor(frame.cmds, 0, []vk.Rect2D{vr.scissor})

	// color clear
	colorClear, ccv := vk.ClearValue{}, vk.ClearColorValue{}
	ccv.AsTypeFloat32(vr.clearColor)
	colorClear.AsColor(ccv)

	// depth buffer clear.
	depthClear := vk.ClearValue{}
	depthClear.AsDepthStencil(vk.ClearDepthStencilValue{
		Depth:   1.0,
		Stencil: 0.0,
	})

	// reset the scene descriptor set updates for all shaders each frame render.
	// OPTIMIZE: only reset sceneUpdated if the view changed.
	for i := range vr.shaders {
		s := &vr.shaders[i]
		for j := range s.sceneUpdated {
			s.sceneUpdated[j] = false
		}
	}

	// define how the rendering attachment will be used.
	colorInfo := vk.RenderingAttachmentInfo{
		ImageView:   vr.views[vr.imageIndex],
		ImageLayout: vk.IMAGE_LAYOUT_ATTACHMENT_OPTIMAL,
		LoadOp:      vk.ATTACHMENT_LOAD_OP_CLEAR,
		StoreOp:     vk.ATTACHMENT_STORE_OP_STORE,
		ClearValue:  colorClear,
	}
	depthInfo := vk.RenderingAttachmentInfo{
		ImageView:   vr.depthImage.view,
		ImageLayout: vk.IMAGE_LAYOUT_ATTACHMENT_OPTIMAL,
		LoadOp:      vk.ATTACHMENT_LOAD_OP_CLEAR,
		StoreOp:     vk.ATTACHMENT_STORE_OP_DONT_CARE,
		ClearValue:  depthClear,
	}
	renderInfo := vk.RenderingInfo{
		RenderArea: vk.Rect2D{
			Offset: vk.Offset2D{X: 0, Y: 0},
			Extent: vk.Extent2D{Width: vr.frameWidth, Height: vr.frameHeight},
		},
		LayerCount:        1,
		PColorAttachments: []vk.RenderingAttachmentInfo{colorInfo},
		PDepthAttachment:  &depthInfo,
	}

	// render a frame.
	vk.CmdBeginRendering(frame.cmds, &renderInfo)
	var shader *vulkanShader
	shaderID := uint16(math.MaxUint16) - 1
	for passID, pass := range passes {
		if RaytraceEnabled {
			// update/rebuild the TLAS acceleration structure
			// based on the passed in objects.
			if passID == 0 { // only for the 3D pass
				vr.updateTLAS(pass)
			}
		}

		// render each object, swapping shaders as needed.
		for _, packet := range pass.Packets {

			// change shader when necessary.
			if shaderID != packet.ShaderID {
				if packet.ShaderID >= uint16(len(vr.shaders)) {
					slog.Error("invalid shaderID", "shader_id", packet.ShaderID)
					continue
				}
				shaderID = packet.ShaderID     // changing shaders.
				shader = &vr.shaders[shaderID] //   ""
				vk.CmdBindPipeline(frame.cmds, vk.PIPELINE_BIND_POINT_GRAPHICS, shader.pipe)

				// set the scene uniforms for this shader
				vr.setSceneUniforms(shader, pass)
				vr.applySceneUniforms(shader)
			}

			// update texture samplers as needed.
			if len(packet.TextureIDs) > 0 {
				samplerID, _ := vr.setSamplers(shader, packet.TextureIDs)
				vr.applySamplerUniforms(shader, samplerID)
			}

			// bind model scope uniforms for this shader.
			vr.setModelUniforms(shader, packet)
			vr.drawMesh(frame, packet.MeshID, shader.attrs, packet.IsInstanced, packet.InstanceCount)
		}
	}
	vk.CmdEndRendering(frame.cmds)

	// transition the image layout to what is needed for presenting a frame.
	presentBarriers := []vk.ImageMemoryBarrier2{
		vk.ImageMemoryBarrier2{
			SrcStageMask:  vk.PIPELINE_STAGE_2_COLOR_ATTACHMENT_OUTPUT_BIT,
			SrcAccessMask: vk.ACCESS_2_COLOR_ATTACHMENT_WRITE_BIT,
			DstStageMask:  vk.PIPELINE_STAGE_2_COLOR_ATTACHMENT_OUTPUT_BIT,
			DstAccessMask: 0,
			OldLayout:     vk.IMAGE_LAYOUT_ATTACHMENT_OPTIMAL,
			NewLayout:     vk.IMAGE_LAYOUT_PRESENT_SRC_KHR,
			Image:         vr.images[vr.imageIndex],
			SubresourceRange: vk.ImageSubresourceRange{
				AspectMask: vk.IMAGE_ASPECT_COLOR_BIT,
				LevelCount: 1, LayerCount: 1},
		},
	}
	presentBarrierInfo := vk.DependencyInfo{PImageMemoryBarriers: presentBarriers}
	vk.CmdPipelineBarrier2(frame.cmds, &presentBarrierInfo)

	// end command recording
	if err = vk.EndCommandBuffer(frame.cmds); err != nil {
		return fmt.Errorf("vk.EndCommandBuffer: %w", err)
	}
	return nil
}

func (vr *vulkanRenderer) endFrame(dt time.Duration) (err error) {
	frame := &vr.frames[vr.frameIndex]

	// submit the frame for render, waits for imageAvailable, signals renderComplete.
	submitInfo := vk.SubmitInfo{
		PWaitSemaphores:   []vk.Semaphore{frame.imageAvailable},
		PWaitDstStageMask: []vk.PipelineStageFlags{vk.PIPELINE_STAGE_COLOR_ATTACHMENT_OUTPUT_BIT},
		PSignalSemaphores: []vk.Semaphore{vr.imageRendered[vr.imageIndex]},
		PCommandBuffers:   []vk.CommandBuffer{frame.cmds},
	}
	if err = vk.QueueSubmit(vr.graphicsQ, []vk.SubmitInfo{submitInfo}, frame.fence); err != nil {
		return fmt.Errorf("vk.QueueSubmit %w", err)
	}

	// present the frame, waits for renderComplete.
	presentInfo := vk.PresentInfoKHR{
		PWaitSemaphores: []vk.Semaphore{vr.imageRendered[vr.imageIndex]}, // wait for GPU render
		PSwapchains:     []vk.SwapchainKHR{vr.swapchain},
		PImageIndices:   []uint32{vr.imageIndex},
	}
	err = vk.QueuePresentKHR(vr.presentQ, &presentInfo)
	if err != nil {
		if err == vk.SUBOPTIMAL_KHR || err == vk.ERROR_OUT_OF_DATE_KHR {
			return nil // wait for resize.
		}
		return fmt.Errorf("endFrame aborted: vkQueuePresentKHR: %w", err)
	}
	vr.frameIndex = (vr.frameIndex + 1) % vr.frameCount
	return nil
}

// setViewportAndScissor updates viewport and scissor to match frame size.
// https://www.saschawillems.de/blog/2019/03/29/flipping-the-vulkan-viewport/
func (vr *vulkanRenderer) setViewportAndScissor() {
	vr.viewport.X = 0 // upper left corner
	vr.viewport.Y = 0 //  ""
	vr.viewport.Width = float32(vr.frameWidth)
	vr.viewport.Height = float32(vr.frameHeight)
	vr.viewport.MinDepth = 0.0
	vr.viewport.MaxDepth = 1.0
	vr.scissor.Offset.X = 0
	vr.scissor.Offset.Y = 0
	vr.scissor.Extent.Width = vr.frameWidth
	vr.scissor.Extent.Height = vr.frameHeight
}

// used for waits.
var waitFrame uint64 = uint64(time.Duration(100 * time.Millisecond)) // 10 frames/sec
var maxTimeout uint64 = math.MaxUint64                               // no limit.

// version returns a vulkan integer version as a string.
func (vr *vulkanRenderer) version(v uint32) string {
	return fmt.Sprintf("%d.%d.%d", vk.VERSION_MAJOR(v), vk.VERSION_MINOR(v), vk.VERSION_PATCH(v))
}

// setSceneUniforms copies the render pass data to scene scope uniforms.
func (vr *vulkanRenderer) setSceneUniforms(shader *vulkanShader, pass Pass) {
	for _, u := range shader.usets.index {
		if u.scope == load.SceneScope {
			vr.setUniform(shader, u, pass.Uniforms[u.sceneUID])
		}
	}
}

// setModelUniforms copies the render packet data to model scope uniforms.
func (vr *vulkanRenderer) setModelUniforms(shader *vulkanShader, packet Packet) {
	for _, u := range shader.usets.index {
		if u.scope == load.PushScope {
			vr.setUniform(shader, u, packet.Uniforms[u.modelUID])
		}
	}
}

// setUniform overwrites the uniform buffer with the given data.
// The data size must match the size expected by the given shader.
func (vr *vulkanRenderer) setUniform(shader *vulkanShader, u *uniform, data []byte) {
	switch u.scope {
	case load.SceneScope:
		offset := uintptr(vr.imageIndex*maxSceneUniformBytes + u.offset)
		dst := (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(shader.sceneUniformsMap)) + offset))
		copy(unsafe.Slice(dst, len(data)), data)
	case load.PushScope:
		frame := vr.frames[vr.frameIndex]
		vk.CmdPushConstants(frame.cmds, shader.pipeLayout, vk.SHADER_STAGE_VERTEX_BIT|vk.SHADER_STAGE_FRAGMENT_BIT, u.offset, data)
	}
}
