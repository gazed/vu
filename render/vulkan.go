// Copyright © 2024 Galvanized Logic Inc.

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
	clear       [4]float32 // rgba clear color.
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

	// createRenderpasses
	render3D             vk.RenderPass    // world render pass.
	render2D             vk.RenderPass    // UI overlay pass.
	render3DFramebuffers []vk.Framebuffer // one framebuffer per swapchain image
	render2DFramebuffers []vk.Framebuffer // one framebuffer per swapchain image

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

	// render frame dynamic state.
	viewport vk.Viewport // same as frame size.
	scissor  vk.Rect2D   // same as frame size.

	// mesh vertex attribute buffers.
	vertexBuffers []vulkanBuffer // non-interleaved.
	maxVertexBuff []uint64       // max bytes allowed
	// instanced model data buffers.
	instanceBuffers []vulkanBuffer // non-interleaved.
	maxInstanceBuff []uint64       //max bytes allowed

	// application GPU resources.
	meshes    []vulkanMesh     // application GPU mesh data
	textures  []vulkanTexture  // application GPU texture data
	shaders   []vulkanShader   // shaders - one pipeline per shader.
	instances []vulkanInstance // application GPU instance data
}

// vkEnabledLayers can be modified by debug builds
// by overriding the addValidationLayer method.
var vkEnabledLayers []string = []string{} // enabled vulkan layers
var addValidationLayer func([]string) ([]string, error) = func(layers []string) ([]string, error) { return layers, nil }

// getVulkanRenderer acquires the vulkan resources needed to render scenes.
func getVulkanRenderer(dev *device.Device, title string) (vr *vulkanRenderer, err error) {
	vr = &vulkanRenderer{}
	vr.title = title
	vr.osdev = dev
	vr.frameWidth, vr.frameHeight = vr.osdev.SurfaceSize() // initial size
	vr.deviceExtensions = []string{vk.KHR_SWAPCHAIN_EXTENSION_NAME}

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
		vr.createSwapchainResources, // one swapchain with one render frame per image
		vr.createRenderFrames,       // two frames

		// two renderpass, each with their own framebuffers.
		// - first draw the 3D world
		// - then draw the 2D UI overlay.
		vr.createRenderpasses,
		vr.createFramebuffers,
		vr.createVertexBuffers,
		vr.createInstanceBuffers,
	}
	for _, create := range createFunctions {
		if err := create(); err != nil {
			vr.dispose()
			return nil, err
		}
	}

	// init is done... throw it back to the application
	// to create the shaders and upload the render data.
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
	vr.disposeInstanceBuffers()
	vr.disposeVertexBuffers()
	for i := range vr.textures {
		vr.dropTexture(uint32(i))
	}
	for sid := range vr.shaders {
		vr.disposeShader(&vr.shaders[sid])
	}

	// per renderpass..
	vr.disposeFramebuffers()
	if vr.render3D != 0 {
		vk.DestroyRenderPass(vr.device, vr.render3D, nil)
		vr.render3D = 0
	}
	if vr.render2D != 0 {
		vk.DestroyRenderPass(vr.device, vr.render2D, nil)
		vr.render2D = 0
	}

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
			ApiVersion:         vk.API_VERSION_1_2,
		},
		PpEnabledLayerNames:     vkEnabledLayers,
		PpEnabledExtensionNames: vr.instanceExtensions(), // vulkan_windows.go
	}
	vr.instance, err = vk.CreateInstance(&instanceInfo, nil)
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
		properties := vk.GetPhysicalDeviceProperties(d)

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
		features := vk.GetPhysicalDeviceFeatures(d)
		if !features.SamplerAnisotropy {
			slog.Warn("missing required features")
			break // device missing samplerAnisotropy
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

	// create the logical device with separate queues for each queue family
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
	deviceCreateInfo := vk.DeviceCreateInfo{
		PQueueCreateInfos: queueInfos,
		PEnabledFeatures: &vk.PhysicalDeviceFeatures{
			SamplerAnisotropy: true,
		},
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

	// TODO allow separate present and graphic queues
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

// =============================================================================
// renderpasses - each renderpass has attachments and requires a set of
//                framebuffers for each swapchain image.

// createRenderpasses creates the 3D and 2D renderpasses.
// The 3D world renderpass is run before the 2D overlay renderpass.
func (vr *vulkanRenderer) createRenderpasses() (err error) {
	renderpassInfo := vk.RenderPassCreateInfo{
		PAttachments: []vk.AttachmentDescription{
			{
				Format:         vr.surfaceFormat.Format,
				Samples:        vk.SAMPLE_COUNT_1_BIT,
				LoadOp:         vk.ATTACHMENT_LOAD_OP_CLEAR,
				StoreOp:        vk.ATTACHMENT_STORE_OP_STORE,
				StencilLoadOp:  vk.ATTACHMENT_LOAD_OP_DONT_CARE,
				StencilStoreOp: vk.ATTACHMENT_STORE_OP_DONT_CARE,
				InitialLayout:  vk.IMAGE_LAYOUT_UNDEFINED,
				FinalLayout:    vk.IMAGE_LAYOUT_COLOR_ATTACHMENT_OPTIMAL, // output to UI pass
			},
			{
				Format:         vr.depthFormat,
				Samples:        vk.SAMPLE_COUNT_1_BIT,
				LoadOp:         vk.ATTACHMENT_LOAD_OP_CLEAR,
				StoreOp:        vk.ATTACHMENT_STORE_OP_DONT_CARE,
				StencilLoadOp:  vk.ATTACHMENT_LOAD_OP_DONT_CARE,
				StencilStoreOp: vk.ATTACHMENT_STORE_OP_DONT_CARE,
				InitialLayout:  vk.IMAGE_LAYOUT_UNDEFINED,
				FinalLayout:    vk.IMAGE_LAYOUT_DEPTH_STENCIL_ATTACHMENT_OPTIMAL,
			},
		},
		PSubpasses: []vk.SubpassDescription{
			{
				PipelineBindPoint: vk.PIPELINE_BIND_POINT_GRAPHICS,
				PColorAttachments: []vk.AttachmentReference{
					{
						Attachment: 0,
						Layout:     vk.IMAGE_LAYOUT_COLOR_ATTACHMENT_OPTIMAL,
					},
				},
				PDepthStencilAttachment: &vk.AttachmentReference{
					Attachment: 1,
					Layout:     vk.IMAGE_LAYOUT_DEPTH_STENCIL_ATTACHMENT_OPTIMAL,
				},
			}},
		PDependencies: []vk.SubpassDependency{
			{
				SrcSubpass:    vk.SUBPASS_EXTERNAL,
				DstSubpass:    0,
				SrcStageMask:  vk.PIPELINE_STAGE_COLOR_ATTACHMENT_OUTPUT_BIT,
				SrcAccessMask: 0,
				DstStageMask:  vk.PIPELINE_STAGE_COLOR_ATTACHMENT_OUTPUT_BIT,
				DstAccessMask: vk.ACCESS_COLOR_ATTACHMENT_READ_BIT | vk.ACCESS_COLOR_ATTACHMENT_WRITE_BIT,
			}},
		Flags: 0,
	}
	vr.render3D, err = vk.CreateRenderPass(vr.device, &renderpassInfo, nil)

	// create the 2D UI overlay renderpass
	renderpassInfo = vk.RenderPassCreateInfo{
		PAttachments: []vk.AttachmentDescription{
			{
				Format:         vr.surfaceFormat.Format,
				Samples:        vk.SAMPLE_COUNT_1_BIT,
				LoadOp:         vk.ATTACHMENT_LOAD_OP_LOAD,
				StoreOp:        vk.ATTACHMENT_STORE_OP_STORE,
				StencilLoadOp:  vk.ATTACHMENT_LOAD_OP_DONT_CARE,
				StencilStoreOp: vk.ATTACHMENT_STORE_OP_DONT_CARE,
				InitialLayout:  vk.IMAGE_LAYOUT_COLOR_ATTACHMENT_OPTIMAL, // FinalLayout of 3D pass
				FinalLayout:    vk.IMAGE_LAYOUT_PRESENT_SRC_KHR,          // present to screen.
			},
		},
		PSubpasses: []vk.SubpassDescription{
			{
				PipelineBindPoint: vk.PIPELINE_BIND_POINT_GRAPHICS,
				PColorAttachments: []vk.AttachmentReference{
					{
						Attachment: 0,
						Layout:     vk.IMAGE_LAYOUT_COLOR_ATTACHMENT_OPTIMAL,
					},
				},
			}},
		PDependencies: []vk.SubpassDependency{
			{
				SrcSubpass:    vk.SUBPASS_EXTERNAL,
				DstSubpass:    0,
				SrcStageMask:  vk.PIPELINE_STAGE_COLOR_ATTACHMENT_OUTPUT_BIT,
				SrcAccessMask: 0,
				DstStageMask:  vk.PIPELINE_STAGE_COLOR_ATTACHMENT_OUTPUT_BIT,
				DstAccessMask: vk.ACCESS_COLOR_ATTACHMENT_READ_BIT | vk.ACCESS_COLOR_ATTACHMENT_WRITE_BIT,
			}},
		Flags: 0,
	}
	vr.render2D, err = vk.CreateRenderPass(vr.device, &renderpassInfo, nil)
	return nil
}

// createFramebuffers for each renderpass where each renderpass
// has a framebuffer for each swapchain image.
//
// Framebuffers are recreated on resizes.
func (vr *vulkanRenderer) createFramebuffers() (err error) {
	vr.render3DFramebuffers = make([]vk.Framebuffer, len(vr.images))
	for i := range vr.images {
		frameInfo := vk.FramebufferCreateInfo{
			RenderPass: vr.render3D,
			PAttachments: []vk.ImageView{
				vr.views[i],
				vr.depthImage.view,
			},
			Width:  vr.frameWidth,
			Height: vr.frameHeight,
			Layers: 1,
		}
		vr.render3DFramebuffers[i], err = vk.CreateFramebuffer(vr.device, &frameInfo, nil)
		if err != nil {
			return fmt.Errorf("vk.CreateFramebuffer: %w", err)
		}
	}
	vr.render2DFramebuffers = make([]vk.Framebuffer, len(vr.images))
	for i := range vr.images {
		frameInfo := vk.FramebufferCreateInfo{
			RenderPass: vr.render2D,
			PAttachments: []vk.ImageView{
				vr.views[i],
			},
			Width:  vr.frameWidth,
			Height: vr.frameHeight,
			Layers: 1,
		}
		vr.render2DFramebuffers[i], err = vk.CreateFramebuffer(vr.device, &frameInfo, nil)
		if err != nil {
			return fmt.Errorf("vk.CreateFramebuffer: %w", err)
		}
	}
	return nil
}

// disposeFramebuffers is needed for window resizes.
func (vr *vulkanRenderer) disposeFramebuffers() {
	for i := range vr.render2DFramebuffers {
		if vr.render2DFramebuffers[i] != 0 {
			vk.DestroyFramebuffer(vr.device, vr.render2DFramebuffers[i], nil)
			vr.render2DFramebuffers[i] = 0
		}
	}
	for i := range vr.render3DFramebuffers {
		if vr.render3DFramebuffers[i] != 0 {
			vk.DestroyFramebuffer(vr.device, vr.render3DFramebuffers[i], nil)
			vr.render3DFramebuffers[i] = 0
		}
	}
}

// createVertexBuffers creates separate buffers for the different possible
// mesh vertex data and triangle index data.
func (vr *vulkanRenderer) createVertexBuffers() (err error) {
	vr.vertexBuffers = make([]vulkanBuffer, load.VertexTypes)
	vr.maxVertexBuff = make([]uint64, load.VertexTypes)
	flags := vk.BUFFER_USAGE_VERTEX_BUFFER_BIT | vk.BUFFER_USAGE_TRANSFER_DST_BIT | vk.BUFFER_USAGE_TRANSFER_SRC_BIT
	props := vk.MEMORY_PROPERTY_DEVICE_LOCAL_BIT
	var size vk.DeviceSize
	var space vk.DeviceSize = vk.DeviceSize(GPUVertexBytes)

	// vertex positions.
	var buff *vulkanBuffer
	buff = &vr.vertexBuffers[load.Vertexes] // V3 or V2 float32
	size = 3 * 4 * space                    // 3-float32 * 4-bytes * lots of space.
	if err = vr.createBuffer(buff, size, flags, props); err != nil {
		return fmt.Errorf("createBuffers:position %w", err)
	}
	vr.maxVertexBuff[load.Vertexes] = uint64(size)

	// vertex texcoords
	buff = &vr.vertexBuffers[load.Texcoords] // V2 float32
	size = 2 * 4 * space                     // 2-float32 * 4-bytes * lots of space.
	if err = vr.createBuffer(buff, size, flags, props); err != nil {
		return fmt.Errorf("createBuffers:texcoord %w", err)
	}
	vr.maxVertexBuff[load.Texcoords] = uint64(size)

	// vertex colors
	buff = &vr.vertexBuffers[load.Colors] // V3 uint8
	size = 3 * 1 * space                  // 3-uint8 * 1-bytes * lots of space.
	if err = vr.createBuffer(buff, size, flags, props); err != nil {
		return fmt.Errorf("createBuffers:color %w", err)
	}
	vr.maxVertexBuff[load.Colors] = uint64(size)

	// vertex normals
	buff = &vr.vertexBuffers[load.Normals] // V3 float32
	size = 3 * 4 * space                   // 3-float32 * 4-bytes * lots of space.
	if err = vr.createBuffer(buff, size, flags, props); err != nil {
		return fmt.Errorf("createBuffers:normal %w", err)
	}
	vr.maxVertexBuff[load.Normals] = uint64(size)

	// FUTURE:
	// vertex load.Tangents V4 float32
	// vertex load.Weights  V4 uint8    animations
	// vertex load.Joints   V4 uint8    animations

	// triangle indexes
	buff = &vr.vertexBuffers[load.Indexes] // uint16
	size = 4 * space                       // 1-float32 * 4-bytes * lots of space.
	iflags := vk.BUFFER_USAGE_INDEX_BUFFER_BIT | vk.BUFFER_USAGE_TRANSFER_DST_BIT | vk.BUFFER_USAGE_TRANSFER_SRC_BIT
	if err = vr.createBuffer(buff, size, iflags, props); err != nil {
		return fmt.Errorf("createBuffers:index %w", err)
	}
	vr.maxVertexBuff[load.Indexes] = uint64(size)
	return nil
}

// disposeVertexbuffers
func (vr *vulkanRenderer) disposeVertexBuffers() {
	for i := range vr.vertexBuffers {
		vr.disposeBuffer(&vr.vertexBuffers[i])
	}
}

// createInstanceBuffers creates separate buffers for the different possible
// instance data buffers.
func (vr *vulkanRenderer) createInstanceBuffers() (err error) {
	vr.instanceBuffers = make([]vulkanBuffer, load.InstanceTypes)
	vr.maxInstanceBuff = make([]uint64, load.InstanceTypes)
	flags := vk.BUFFER_USAGE_VERTEX_BUFFER_BIT | vk.BUFFER_USAGE_TRANSFER_DST_BIT | vk.BUFFER_USAGE_TRANSFER_SRC_BIT
	props := vk.MEMORY_PROPERTY_DEVICE_LOCAL_BIT
	var size vk.DeviceSize
	var space vk.DeviceSize = vk.DeviceSize(GPUInstanceBytes)

	// instance positions.
	var buff *vulkanBuffer
	buff = &vr.instanceBuffers[load.InstancePosition] // V3 or V2 float32
	size = 3 * 4 * space                              // 3-float32 * 4-bytes * lots of space.
	if err = vr.createBuffer(buff, size, flags, props); err != nil {
		return fmt.Errorf("createBuffers:position %w", err)
	}
	vr.maxInstanceBuff[load.InstancePosition] = uint64(size)

	// instance colors
	buff = &vr.instanceBuffers[load.InstanceColors] // V3 uint8
	size = 3 * 1 * space                            // 3-uint8 * 1-bytes * lots of space.
	if err = vr.createBuffer(buff, size, flags, props); err != nil {
		return fmt.Errorf("createBuffers:color %w", err)
	}
	vr.maxInstanceBuff[load.InstanceColors] = uint64(size)

	// instance scale
	buff = &vr.instanceBuffers[load.InstanceScales] // float32
	size = 4 * space                                // float32 * 4-bytes * lots of space.
	if err = vr.createBuffer(buff, size, flags, props); err != nil {
		return fmt.Errorf("createBuffers:normal %w", err)
	}
	vr.maxInstanceBuff[load.InstanceScales] = uint64(size)
	return nil
}

// disposeInstancebuffers
func (vr *vulkanRenderer) disposeInstanceBuffers() {
	for i := range vr.instanceBuffers {
		vr.disposeBuffer(&vr.instanceBuffers[i])
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
	vr.disposeFramebuffers()       // delete the renderpass framebuffers.
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
	vr.createFramebuffers() // recreate the renderpass framebuffers.
	// ---
	vr.recreatingSwapchain = false // ok to render
	vr.resizesCompleted = vr.resizesRequested
	slog.Info("vulkan resize complete", "size", fmt.Sprintf("%d:%d", vr.frameWidth, vr.frameHeight))
	return nil
}

// ============================================================================
// buffers allocate GPU data memory and transport data from the CPU to the GPU.
// lock, load, unlock a buffer generally in that order.

type vulkanBuffer struct {
	handle vk.Buffer
	memory vk.DeviceMemory
}

// createBuffer allocates a buffer, memory, and binds the buffer to the memory
func (vr *vulkanRenderer) createBuffer(buff *vulkanBuffer, size vk.DeviceSize,
	usage vk.BufferUsageFlagBits, flags vk.MemoryPropertyFlags) (err error) {

	// create the buffer structure
	buffInfo := vk.BufferCreateInfo{
		Size:        size,
		Usage:       usage,
		SharingMode: vk.SHARING_MODE_EXCLUSIVE, // used in a single queue
	}
	if buff.handle, err = vk.CreateBuffer(vr.device, &buffInfo, nil); err != nil {
		return fmt.Errorf("vk.CreateBuffer: %w", err)
	}
	if buff.handle == 0 {
		// Check is here because there was a bug in the original vulkan bindings.
		return fmt.Errorf("vk.CreateBuffer: 0 handle: %+v", buffInfo)
	}

	// allocate the memory needed by the buffer.
	memRequirements := vk.GetBufferMemoryRequirements(vr.device, buff.handle)
	memType, err := vr.findMemoryType(memRequirements.MemoryTypeBits, flags)
	if err != nil {
		return err
	}
	allocateInfo := vk.MemoryAllocateInfo{
		AllocationSize:  memRequirements.Size,
		MemoryTypeIndex: memType,
	}
	if buff.memory, err = vk.AllocateMemory(vr.device, &allocateInfo, nil); err != nil {
		return fmt.Errorf("createBuff:vk.AllocateMemory: %w", err)
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
// Copying to device local memory is done through the host coherent staging buffer and
// this uses the transfer queue.
func (vr *vulkanRenderer) uploadData(pool vk.CommandPool, queue vk.Queue, buff *vulkanBuffer, offset uint64, data []byte) {

	// create a staging buffer and load data into the staging buffer.
	var staging vulkanBuffer
	usage := vk.BUFFER_USAGE_TRANSFER_SRC_BIT
	flags := vk.MEMORY_PROPERTY_HOST_VISIBLE_BIT | vk.MEMORY_PROPERTY_HOST_COHERENT_BIT
	size := vk.DeviceSize(len(data))
	err := vr.createBuffer(&staging, size, usage, flags)
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
	vr.copyGPUBuffer(pool, queue, staging.handle, 0, buff.handle, offset, size)
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
	cmd, err := vr.beginSingleUseCommand(pool)
	if err != nil {
		return fmt.Errorf("copyBuffer: get command: %w", err)
	}
	region := vk.BufferCopy{SrcOffset: 0, DstOffset: vk.DeviceSize(dstOffset), Size: size}
	vk.CmdCopyBuffer(cmd, src, dst, []vk.BufferCopy{region})

	// submit one-time command for execution and wait for it to complete.
	if err := vr.endSingleUseCommand(cmd, pool, vr.graphicsQ); err != nil {
		return fmt.Errorf("copyBuffer: end command: %w", err)
	}
	return nil
}

// =============================================================================
// mesh vertex and index data are a bunch of bytes that exist
// within larger vulkan buffers
//   mesh vertices are held in each frames vertex buffer
//   mesh indices are held in each frames index buffer

// vulkanMesh tracks where the mesh data is stored in the vertex buffers.
// Indexed by the load.MeshData types.
type vulkanMesh []vulkanBuffData
type vulkanBuffData struct {
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

	// Add the vertex data at the end of the current mesh data.
	// Get the current offsets for each vertex data type
	startingOffsets := make([]uint32, load.VertexTypes) // tracks first offset.
	for i := 0; i < load.VertexTypes; i++ {
		if len(vr.meshes) > 0 {
			prev := vr.meshes[len(vr.meshes)-1]
			for i := 0; i < load.VertexTypes; i++ {
				startingOffsets[i] = prev[i].offset + prev[i].stride*prev[i].count
			}
		}
	}

	// cosolidate the upload data into a buffer for each vertex data type.
	data := make([][]byte, load.VertexTypes)
	for _, msh := range meshes {

		// add the mesh data after the last mesh.
		meshOffsets := make([]uint32, load.VertexTypes)
		if len(vr.meshes) > 0 {
			prev := vr.meshes[len(vr.meshes)-1]
			for i := 0; i < load.VertexTypes; i++ {
				meshOffsets[i] = prev[i].offset + prev[i].stride*prev[i].count
			}
		}

		// track each mesh with a vulkan-mesh
		vmsh := make(vulkanMesh, load.VertexTypes)
		for i := 0; i < load.VertexTypes; i++ {
			if msh[i].Count > 0 {
				vmsh[i].count = msh[i].Count
				vmsh[i].stride = msh[i].Stride
				vmsh[i].offset = meshOffsets[i]

				// check if the data exceeds what was allocated.
				total := uint64(vmsh[i].offset + msh[i].Count*msh[i].Stride)
				if total > vr.maxVertexBuff[i] {
					return mids, fmt.Errorf("loadMeshes:insufficient memory %d", i)
				}

				// consolidate the upload data into temp buffers.
				data[i] = append(data[i], msh[i].Data...)

				// track the total amount of uploaded bytes
				GPUTotalMeshBytes += msh[i].Count * msh[i].Stride
			} else {
				// push forward the previous offset for the vertex data
				// types that were not used by this mesh.
				vmsh[i].offset = meshOffsets[i]
			}
		}
		mids = append(mids, uint32(len(vr.meshes))) // mesh ID for the new mesh.
		vr.meshes = append(vr.meshes, vmsh)
	}

	// upload the consolidated data once for each data type.
	for i := 0; i < load.VertexTypes; i++ {
		uploadData := data[i]
		if len(uploadData) > 0 {
			buff := &vr.vertexBuffers[i]
			offset := uint64(startingOffsets[i])
			vr.uploadData(vr.graphicsQCmdPool, vr.graphicsQ, buff, offset, uploadData)
		}
	}
	return mids, nil
}

// FUTURE: handle deallocates using linked lists.
// For now never deallocate so that the lastMesh is always valid.
func (vr *vulkanRenderer) dropMesh(mid uint32) {}

// track instanced data as a number of buffers.
type vulkanInstance []vulkanBuffData

// loadInstanceData stores instance data in GPU buffers.
// Immutable once uploaded. There is only one instance buffer for
// all frames. Updating instance data means adding a new data and
// refering to it in future draw calls.
//
// FUTURE: handle instance data (de/re)allocates using linked lists.
func (vr *vulkanRenderer) loadInstanceData(data []load.Buffer) (iid uint32, err error) {
	inst := make(vulkanInstance, load.InstanceTypes)

	// add the instance data at the end of the buffer
	offsets := make([]uint32, load.InstanceTypes)
	if len(vr.instances) > 0 {
		prev := vr.instances[len(vr.instances)-1]
		for i := 0; i < load.InstanceTypes; i++ {
			offsets[i] = prev[i].offset + prev[i].stride*prev[i].count
		}
	}
	for i := 0; i < load.InstanceTypes; i++ {
		if data[i].Count > 0 {
			inst[i].count = data[i].Count
			inst[i].stride = data[i].Stride
			inst[i].offset = offsets[i]

			// check if the data exceeds what was allocated.
			total := uint64(inst[i].offset + inst[i].count*inst[i].stride)
			if total > vr.maxInstanceBuff[i] {
				return iid, fmt.Errorf("loadInstanceData:insufficient memory %d", i)
			}

			// upload data
			buff := &vr.instanceBuffers[i]
			offset := uint64(inst[i].offset)
			vr.uploadData(vr.graphicsQCmdPool, vr.graphicsQ, buff, offset, data[i].Data)

			// track the total amount of uploaded bytes
			GPUTotalInstanceBytes += data[i].Count * data[i].Stride
		}
	}
	iid = uint32(len(vr.instances)) // instance ID for the new instance data.
	vr.instances = append(vr.instances, inst)
	return iid, nil
}

// updateInstanceData : see docs on render:UpdateInstanceData
func (vr *vulkanRenderer) updateInstanceData(iid uint32, data []load.Buffer) (err error) {
	if int(iid) >= len(vr.instances) {
		return fmt.Errorf("updateInstanceData invalid ID: %d", iid)
	}
	inst := vr.instances[iid]

	// check that the buffer data matches.
	for i := 0; i < load.InstanceTypes; i++ {
		if inst[i].count != data[i].Count {
			return fmt.Errorf("updateInstanceData count mismatch %d %d", inst[i].count, data[i].Count)
		}
		if inst[i].stride != data[i].Stride {
			return fmt.Errorf("updateInstanceData stride mismatch %d %d", inst[i].stride, data[i].Stride)
		}
	}

	// everything matches, so re-upload data.
	for i := 0; i < load.InstanceTypes; i++ {
		if data[i].Count > 0 {
			// upload data - existing instance data remains the same: count, stride, offset
			buff := &vr.instanceBuffers[i]
			offset := uint64(inst[i].offset)
			vr.uploadData(vr.graphicsQCmdPool, vr.graphicsQ, buff, offset, data[i].Data)
		}
	}
	return nil // everything ok.
}

// FUTURE: handle deallocates using linked lists.
// For now never deallocate so that the lastMesh is always valid.
func (vr *vulkanRenderer) dropInstanceData(iid uint32) {}

// drawMesh
func (vr *vulkanRenderer) drawMesh(frame *vulkanFrame, mid uint32, attrs []load.ShaderAttribute) {
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
		switch attr.AttrScope {
		case load.VertexAttribute:
			i := attr.AttrType
			if i < 0 || i >= load.Indexes {
				slog.Error("unsupported vertex attribute", "attribute_type", attr.AttrType)
				continue
			}
			buffs = append(buffs, vr.vertexBuffers[i].handle)
			offsets = append(offsets, vk.DeviceSize(vmsh[i].offset))
		}
	}
	vk.CmdBindVertexBuffers(frame.cmds, 0, buffs, offsets)

	// bind the vertex index data.
	ibuff := vr.vertexBuffers[load.Indexes].handle
	ioffset := vk.DeviceSize(vmsh[load.Indexes].offset)
	vk.CmdBindIndexBuffer(frame.cmds, ibuff, ioffset, vk.INDEX_TYPE_UINT16)

	// draw the mesh
	vk.CmdDrawIndexed(frame.cmds, vmsh[load.Indexes].count, 1, 0, 0, 0)
}

// drawInstancedMesh
func (vr *vulkanRenderer) drawInstancedMesh(frame *vulkanFrame, mid, instID, instCount uint32, attrs []load.ShaderAttribute) {
	if mid < 0 || mid >= uint32(len(vr.meshes)) {
		slog.Error("invalid mesh ID", "mid", mid)
		return
	}
	vmsh := vr.meshes[mid]
	if instID >= uint32(len(vr.instances)) {
		slog.Error("invalid instance ID", "inst_ID", instID)
		return
	}
	instData := vr.instances[instID]

	// bind the mesh vertex attribute data expected by the shader.
	// The attribute must match one of the supported vertex types.
	buffs := []vk.Buffer{}
	offsets := []vk.DeviceSize{}
	for _, attr := range attrs {
		switch attr.AttrScope {
		case load.VertexAttribute:
			i := attr.AttrType
			if i < 0 || i >= load.Indexes {
				slog.Error("unsupported vertex attribute", "attribute_type", attr.AttrType)
				continue
			}
			buffs = append(buffs, vr.vertexBuffers[i].handle)
			offsets = append(offsets, vk.DeviceSize(vmsh[i].offset))
		case load.InstanceAttribute:
			i := attr.AttrType
			if i < 0 || i >= load.InstanceTypes {
				slog.Error("unsupported instance attribute", "attribute_type", attr.AttrType)
				continue
			}
			buffs = append(buffs, vr.instanceBuffers[i].handle)
			offsets = append(offsets, vk.DeviceSize(instData[i].offset))
		}
	}
	vk.CmdBindVertexBuffers(frame.cmds, 0, buffs, offsets)

	// bind the triangle index data.
	ibuff := vr.vertexBuffers[load.Indexes].handle
	ioffset := vk.DeviceSize(vmsh[load.Indexes].offset)
	vk.CmdBindIndexBuffer(frame.cmds, ibuff, ioffset, vk.INDEX_TYPE_UINT16)

	// draw the instanced mesh
	vk.CmdDrawIndexed(frame.cmds, vmsh[load.Indexes].count, instCount, 0, 0, 0)
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
		MipLevels:     4,
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
	cmd, err := vr.beginSingleUseCommand(vr.graphicsQCmdPool)
	if err != nil {
		slog.Error("beginSingleUseCommand", "error", err)
		return
	}
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
	vr.endSingleUseCommand(cmd, vr.graphicsQCmdPool, vr.graphicsQ)
}

func (vr *vulkanRenderer) copyBufferToImage(buffer *vulkanBuffer, img *vulkanImage) {
	cmd, err := vr.beginSingleUseCommand(vr.graphicsQCmdPool)
	if err != nil {
		slog.Error("beginSingleUseCommand", "error", err)
		return
	}
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
	vr.endSingleUseCommand(cmd, vr.graphicsQCmdPool, vr.graphicsQ)
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
	usets uniformSets // shader uniform information

	// scene uniform data buffers for all scene uniform descriptors
	// indexed by vr.imageIndex
	sceneUniforms    vulkanBuffer // scene scope uniform data.
	sceneUniformsMap *byte        // unsafe pointer to uniform mapped memory

	// track the number of unique material instances for this shader.
	maxMaterials   uint32           // maximum materials supported by shader.
	materials      []vulkanMaterial // track unique materials (sets) with...
	nextMaterialID uint32           // ... a material ID.

	// material uniform data buffer for all of this shaders material uniform data.
	// indexed by vr.imageIndex and material ID.
	materialUniforms    vulkanBuffer // material scope uniform data.
	materialUniformsMap *byte        // unsafe pointer to uniform mapped memory

	// descriptors for scene and material uniform data.
	sceneLayout         vk.DescriptorSetLayout // scene uniforms per renderpass
	materialLayout      vk.DescriptorSetLayout // material uniforms per object
	descriptorPool      vk.DescriptorPool      // uniforms and samplers
	sceneDescriptorSets []vk.DescriptorSet     // one per image.
	sceneUpdated        []bool                 // true if descriptor set updated.
}

// vulkanMaterial tracks existing resources to help reuse descriptor sets.
type vulkanMaterial struct {
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
	load.DataType_FLOAT: vk.FORMAT_R32_SFLOAT,
	load.DataType_VEC2:  vk.FORMAT_R32G32_SFLOAT,
	load.DataType_VEC3:  vk.FORMAT_R32G32B32_SFLOAT,
	load.DataType_VEC4:  vk.FORMAT_R32G32B32A32_SFLOAT,
}

// loadShader creates a shader and corresponding pipeline based
// on the given shader configuration.
func (vr *vulkanRenderer) loadShader(config *load.Shader) (sid uint16, err error) {
	sid = uint16(len(vr.shaders)) // the shader ID if loadShader succeeds.
	shader := vulkanShader{name: config.Name}
	shader.attrs = append(shader.attrs, config.Attrs...)
	shader.usets = getUniformSets(config.Uniforms)
	shader.maxMaterials = 256 // FUTURE: get from shader config.
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
	// no longer need the modules once the shader has been created.
	for i := range stages {
		defer vk.DestroyShaderModule(vr.device, stages[i].Module, nil)
	}

	// non-interleaved attribute descriptions
	vertexAttrDescriptions := make([]vk.VertexInputAttributeDescription, len(config.Attrs))
	vertexBindingDescriptions := make([]vk.VertexInputBindingDescription, len(config.Attrs))
	if len(config.Attrs) > 0 {
		for i, attr := range config.Attrs {
			vertexAttrDescriptions[i].Location = uint32(i)
			vertexAttrDescriptions[i].Binding = uint32(i)
			vertexAttrDescriptions[i].Format = vulkanDataFormats[attr.DataType]
			vertexAttrDescriptions[i].Offset = 0

			vertexBindingDescriptions[i].Binding = uint32(i)
			vertexBindingDescriptions[i].Stride = load.DataTypeSizes[attr.DataType]
			switch attr.AttrScope {
			case load.VertexAttribute:
				vertexBindingDescriptions[i].InputRate = vk.VERTEX_INPUT_RATE_VERTEX
			case load.InstanceAttribute:
				vertexBindingDescriptions[i].InputRate = vk.VERTEX_INPUT_RATE_INSTANCE
			}
		}
	}
	vertexInputInfo := &vk.PipelineVertexInputStateCreateInfo{
		PVertexBindingDescriptions:   vertexBindingDescriptions,
		PVertexAttributeDescriptions: vertexAttrDescriptions,
	}

	// allocate the scene descriptor set layout if applicable.
	scenes := config.GetSceneUniforms()
	if len(scenes) > 0 {
		bindings := []vk.DescriptorSetLayoutBinding{}
		binding := vk.DescriptorSetLayoutBinding{
			Binding:         0,
			DescriptorType:  vk.DESCRIPTOR_TYPE_UNIFORM_BUFFER,
			DescriptorCount: 1,
			StageFlags:      vk.SHADER_STAGE_VERTEX_BIT | vk.SHADER_STAGE_FRAGMENT_BIT,
		}
		bindings = append(bindings, binding)
		shader.sceneLayout, err = vk.CreateDescriptorSetLayout(
			vr.device, &vk.DescriptorSetLayoutCreateInfo{PBindings: bindings}, nil)
	}

	// allocate the material scope descriptor set layout if applicable.
	samplers := config.GetSamplerUniforms()
	if len(samplers) > 0 {
		bindings := []vk.DescriptorSetLayoutBinding{}
		for range samplers {
			binding := vk.DescriptorSetLayoutBinding{
				Binding:         uint32(len(bindings)),
				DescriptorType:  vk.DESCRIPTOR_TYPE_COMBINED_IMAGE_SAMPLER,
				DescriptorCount: 1,
				StageFlags:      vk.SHADER_STAGE_FRAGMENT_BIT,
			}
			bindings = append(bindings, binding)
		}
		shader.materialLayout, err = vk.CreateDescriptorSetLayout(
			vr.device, &vk.DescriptorSetLayoutCreateInfo{PBindings: bindings}, nil)
	}

	// create descriptor pools. Each shader has its own pool that can allocate
	// a descriptor set for each render image, normally 3.
	shader.descriptorPool, err = vk.CreateDescriptorPool(vr.device,
		&vk.DescriptorPoolCreateInfo{
			MaxSets: 3 + 3*shader.maxMaterials + 3, // 3 scene sets + 3 per material.
			PPoolSizes: []vk.DescriptorPoolSize{
				{
					Typ:             vk.DESCRIPTOR_TYPE_UNIFORM_BUFFER,
					DescriptorCount: 1024,
				},
				{
					Typ:             vk.DESCRIPTOR_TYPE_COMBINED_IMAGE_SAMPLER,
					DescriptorCount: 4096,
				},
			},
			Flags: vk.DESCRIPTOR_POOL_CREATE_FREE_DESCRIPTOR_SET_BIT, // | vk.DESCRIPTOR_POOL_CREATE_UPDATE_AFTER_BIND_BIT;
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
		shader.sceneDescriptorSets, err = vk.AllocateDescriptorSets(vr.device,
			&vk.DescriptorSetAllocateInfo{
				DescriptorPool: shader.descriptorPool,
				PSetLayouts:    allocLayouts,
			})
		if err != nil {
			vr.disposeShader(&shader)
			return 0, err
		}
		shader.sceneUpdated = make([]bool, len(shader.sceneDescriptorSets))
	}

	// allocate material descriptor sets.
	if shader.materialLayout != 0 {

		// allocate three descriptor sets (one per surface image)...
		allocLayouts := []vk.DescriptorSetLayout{}
		for i := 0; i < int(vr.imageCount); i++ {
			allocLayouts = append(allocLayouts, shader.materialLayout)
		}
		// ...for each material.
		shader.materials = make([]vulkanMaterial, shader.maxMaterials)
		for i := uint32(0); i < shader.maxMaterials; i++ {
			shader.materials[i].updated = make([]bool, vr.imageCount)
			shader.materials[i].descriptorSets, err = vk.AllocateDescriptorSets(vr.device,
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
	if shader.materialLayout != 0 {
		pipelineLayouts = append(pipelineLayouts, shader.materialLayout)
	}
	layoutInfo := vk.PipelineLayoutCreateInfo{
		PSetLayouts: pipelineLayouts,
	}

	// setup pipeline push constants if they exist
	if shader.usets.modelSize > 0 {
		push_constant := vk.PushConstantRange{
			Offset:     0,
			Size:       maxModelUniformBytes, // allocate the max 128 bytes.
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
		DepthCompareOp:        vk.COMPARE_OP_LESS,
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

	// create the pipeline. Use shader naming convention to pick
	// a renderpass.
	pipelineInfo := vk.GraphicsPipelineCreateInfo{
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
		RenderPass:          vr.render3D,
		Subpass:             0,
		BasePipelineHandle:  0,
		BasePipelineIndex:   -1,
	}
	if config.Pass == "2D" {
		pipelineInfo.RenderPass = vr.render2D
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

		// load the shader module bytes
		filename := fmt.Sprintf("%s.%s.spv", name, stageName)
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
	if s.materialLayout != 0 {
		vk.DestroyDescriptorSetLayout(vr.device, s.materialLayout, nil)
		s.materialLayout = 0
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
	flags := vk.MEMORY_PROPERTY_HOST_VISIBLE_BIT | vk.MEMORY_PROPERTY_HOST_COHERENT_BIT | deviceLocalBits
	numImages := uint32(len(vr.images))

	// create enough scene uniform data buffer space for each surface image
	// map the uniform memory once for the lifetime of the app.
	bufferSize := vk.DeviceSize(maxSceneUniformBytes * numImages)
	err = vr.createBuffer(&s.sceneUniforms, bufferSize, vk.BUFFER_USAGE_UNIFORM_BUFFER_BIT, flags)
	if err != nil {
		return fmt.Errorf("sceneUniformsMap:vk.createBuffer: %w", err)
	}
	s.sceneUniformsMap, err = vk.MapMemory(vr.device, s.sceneUniforms.memory, 0, bufferSize, 0)
	if err != nil {
		return fmt.Errorf("sceneUniformsMap:vk.MapMemory: %w", err)
	}

	// create the material scope uniform buffer for this frame.
	// map the uniform memory once for the lifetime of the app.
	bufferSize = vk.DeviceSize(maxMaterialUniformBytes * s.maxMaterials * numImages)
	err = vr.createBuffer(&s.materialUniforms, bufferSize, vk.BUFFER_USAGE_UNIFORM_BUFFER_BIT, flags)
	if err != nil {
		return fmt.Errorf("materialUniformsMap:vk.createBuffer: %w", err)
	}
	s.materialUniformsMap, err = vk.MapMemory(vr.device, s.materialUniforms.memory, 0, bufferSize, 0)
	if err != nil {
		return fmt.Errorf("materialUniformsMap:vk.MapMemory: %w", err)
	}
	return nil
}

// disposeShaderUniformBuffers diposes resources allocated in createShaderUniformBuffers.
func (vr *vulkanRenderer) disposeShaderUniformBuffers(s *vulkanShader) {
	if s != nil {
		vr.disposeBuffer(&s.materialUniforms)
		vr.disposeBuffer(&s.sceneUniforms)
		s.materialUniformsMap = nil
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
	descriptorSet := shader.sceneDescriptorSets[vr.imageIndex]
	if !shader.sceneUpdated[vr.imageIndex] {
		offset := vk.DeviceSize(vr.imageIndex * maxSceneUniformBytes)
		descriptorSetWrites := []vk.WriteDescriptorSet{
			{
				DstSet:          descriptorSet,
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
		vk.UpdateDescriptorSets(vr.device, descriptorSetWrites, nil)
		shader.sceneUpdated[vr.imageIndex] = true
	}
	setNum := uint32(0) // scene is always set=0
	dsets := []vk.DescriptorSet{descriptorSet}
	frame := vr.frames[vr.frameIndex]
	vk.CmdBindDescriptorSets(frame.cmds, vk.PIPELINE_BIND_POINT_GRAPHICS, shader.pipeLayout, setNum, dsets, nil)
}

// setMaterialSamplers sets all the samplers expected by this shader.
// The samplers are in the order expected by the shader config.
func (vr *vulkanRenderer) setMaterialSamplers(shader *vulkanShader, tids []uint32) (matID uint32, err error) {
	// compare to the existing material samplers to see if there is a match
	for i := uint32(0); i < shader.nextMaterialID; i++ {
		if slices.Compare(tids, shader.materials[i].samplerSet) == 0 {
			return i, nil // reuse existing material.
		}
	}

	// create a new material for these textures on this shader.
	matID = shader.nextMaterialID
	if matID >= shader.maxMaterials {
		return 0, fmt.Errorf("setMaterialSamplers:max materials exceeded:%d", matID)
	}
	shader.materials[matID].samplerSet = append(shader.materials[matID].samplerSet, tids...)
	shader.nextMaterialID += 1
	return matID, nil
}

// applyMaterialUniforms updates the material scope descriptor sets
// to point to the proper material uniform data buffer
func (vr *vulkanRenderer) applyMaterialUniforms(shader *vulkanShader, matID uint32) {
	if shader.materialLayout == 0 {
		slog.Error("applyMaterialUniforms: no material uniforms", "shader", shader.name)
		return
	}
	if matID >= uint32(len(shader.materials)) {
		slog.Error("applyMaterialUniforms:invalid material ID", "material_id", matID)
		return
	}
	material := &shader.materials[matID]
	descriptorSet := material.descriptorSets[vr.imageIndex]
	setNum := uint32(1) // material uniforms are set=1 (scene uniforms at set=0)

	// check if the material descriptor set needs updating.
	if !material.updated[vr.imageIndex] {
		descriptorSetWrites := []vk.WriteDescriptorSet{}
		descriptorIndex := uint32(0)

		// FUTURE add support for material uniforms
		// if shader.usets.materialSize > 0 {
		// 	offset := vk.DeviceSize(vr.imageIndex*shader.maxMaterials*maxMaterialUniformBytes + matID*maxMaterialUniformBytes)
		// 	descriptorSetWrites = append(descriptorSetWrites, vk.WriteDescriptorSet{
		// 		DstSet:         descriptorSet,
		// 		DstBinding:     descriptorIndex,
		// 		DescriptorType: vk.DESCRIPTOR_TYPE_UNIFORM_BUFFER,
		// 		PBufferInfo: []vk.DescriptorBufferInfo{
		// 			{
		// 				Buffer: shader.materialUniforms.handle,
		// 				Offset: offset,
		// 				Rang:   vk.DeviceSize(maxMaterialUniformBytes),
		// 			},
		// 		},
		// 	})
		// 	descriptorIndex += 1
		// }

		// check for samplers
		if len(material.samplerSet) > 0 {
			samplerInfo := []vk.DescriptorImageInfo{}
			for _, tid := range material.samplerSet {
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
		material.updated[vr.imageIndex] = true
	}
	frame := &vr.frames[vr.frameIndex]
	dsets := []vk.DescriptorSet{descriptorSet}
	vk.CmdBindDescriptorSets(frame.cmds, vk.PIPELINE_BIND_POINT_GRAPHICS, shader.pipeLayout, setNum, dsets, nil)
}

// =============================================================================
// command utilities

// convenience create for command buffer that is used once.
func (vr *vulkanRenderer) beginSingleUseCommand(pool vk.CommandPool) (cmd vk.CommandBuffer, err error) {
	commands, err := vk.AllocateCommandBuffers(vr.device,
		&vk.CommandBufferAllocateInfo{
			CommandPool:        pool,
			Level:              vk.COMMAND_BUFFER_LEVEL_PRIMARY,
			CommandBufferCount: 1,
		})
	if err != nil {
		return cmd, err
	}

	cmd = commands[0]
	err = vk.BeginCommandBuffer(cmd, &vk.CommandBufferBeginInfo{Flags: vk.COMMAND_BUFFER_USAGE_ONE_TIME_SUBMIT_BIT})
	return cmd, err
}

// convenience end for command buffer that is used once.
func (vr *vulkanRenderer) endSingleUseCommand(cmd vk.CommandBuffer, pool vk.CommandPool, queue vk.Queue) error {
	if cmd == 0 {
		return fmt.Errorf("endSingleUseCommand: invalid command buffer")
	}
	vk.EndCommandBuffer(cmd)

	// submit the command buffer
	submitInfo := vk.SubmitInfo{PCommandBuffers: []vk.CommandBuffer{cmd}}
	if err := vk.QueueSubmit(queue, []vk.SubmitInfo{submitInfo}, 0); err != nil {
		return fmt.Errorf("vk.QueueSubmit: %w", err)
	}

	// wait for submit to finish....and then free the command buffer
	if err := vk.QueueWaitIdle(queue); err != nil {
		return fmt.Errorf("vk.QueueWaitIdle: %w", err)
	}
	vk.FreeCommandBuffers(vr.device, pool, []vk.CommandBuffer{cmd})
	return nil
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
	renderComplete vk.Semaphore // GPU completed, ready for presentation
	inFlightFence  vk.Fence     // render to frames not in use by GPU.
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
		if vr.frames[i].renderComplete != 0 {
			vk.DestroySemaphore(vr.device, vr.frames[i].renderComplete, nil)
			vr.frames[i].renderComplete = 0
		}
		if vr.frames[i].inFlightFence != 0 {
			vk.DestroyFence(vr.device, vr.frames[i].inFlightFence, nil)
			vr.frames[i].inFlightFence = 0
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
	fr.renderComplete, err = vk.CreateSemaphore(vr.device, &vk.SemaphoreCreateInfo{}, nil)
	if err != nil {
		return fmt.Errorf("vk.CreateSemaphore.2: %w", err)
	}

	// Create the fence in a signaled state, indicating that the first frame has
	// already been "rendered". This will prevent the application from waiting
	// indefinitely for the first frame to render since it cannot be rendered
	// until a frame is "rendered" before it.
	fenceInfo := vk.FenceCreateInfo{Flags: vk.FENCE_CREATE_SIGNALED_BIT}
	fr.inFlightFence, err = vk.CreateFence(vr.device, &fenceInfo, nil)
	if err != nil {
		return fmt.Errorf("vk.CreateFence: %w", err)
	}
	return nil
}

// =============================================================================
// beginFrame, render objects, endFrame

func (vr *vulkanRenderer) setClearColor(r, g, b, a float32) {
	vr.clear[0], vr.clear[1], vr.clear[2], vr.clear[3] = r, g, b, a
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
	err = vk.WaitForFences(vr.device, []vk.Fence{frame.inFlightFence}, true, waitFrame)
	if err != nil {
		return fmt.Errorf("beginFrame aborted: vk.WaitForFences: %w", err)
	}

	// acquire the next image from the swapchain.
	// Pass in the semaphore to be signalled when image is available again.
	vr.imageIndex, err = vk.AcquireNextImageKHR(vr.device, vr.swapchain, maxTimeout, frame.imageAvailable, 0)
	if err != nil {
		if err == vk.SUBOPTIMAL_KHR || err == vk.ERROR_OUT_OF_DATE_KHR {
			vr.resize(vr.frameWidth, vr.frameHeight)
			return nil // didn't quite work.
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

	// set pipeline dynamic state
	vr.setViewportAndScissor()
	vk.CmdSetViewport(frame.cmds, 0, []vk.Viewport{vr.viewport})
	vk.CmdSetScissor(frame.cmds, 0, []vk.Rect2D{vr.scissor})

	// color clear
	colorClear, ccv := vk.ClearValue{}, vk.ClearColorValue{}
	ccv.AsTypeFloat32(vr.clear)
	colorClear.AsColor(ccv)

	// depth buffer clear.
	depthClear := vk.ClearValue{}
	depthClear.AsDepthStencil(vk.ClearDepthStencilValue{
		Depth:   1.0,
		Stencil: 0.0,
	})

	// reset the scene descriptor set updates for all shaders each frame render.
	// OPTIMIZE: only reset the sceneUpdated if the view changed.
	for i := range vr.shaders {
		s := &vr.shaders[i]
		for j := range s.sceneUpdated {
			s.sceneUpdated[j] = false
		}
	}

	// first pass always 3D (can be empty if only 2D).
	// start the 3D world render pass
	render3DInfo := vk.RenderPassBeginInfo{
		RenderPass:  vr.render3D,
		Framebuffer: vr.render3DFramebuffers[vr.imageIndex],
		RenderArea: vk.Rect2D{
			Offset: vk.Offset2D{X: 0, Y: 0},
			Extent: vk.Extent2D{Width: vr.frameWidth, Height: vr.frameHeight},
		},
		PClearValues: []vk.ClearValue{colorClear, depthClear},
	}
	vk.CmdBeginRenderPass(frame.cmds, &render3DInfo, vk.SUBPASS_CONTENTS_INLINE)

	var shader *vulkanShader
	shaderID := uint16(math.MaxUint16) - 1
	if len(passes) > 0 && len(passes[Pass3D].Packets) > 0 {
		pass := passes[Pass3D]

		// draw 3D packets
		for _, packet := range pass.Packets {
			// TODO complain about packets without meshes.

			// change shader when necessary.
			if shaderID != packet.ShaderID {
				if packet.ShaderID >= uint16(len(vr.shaders)) {
					slog.Error("invalid shaderID", "shader_id", packet.ShaderID)
					continue
				}
				shaderID = packet.ShaderID // changing shaders.
				shader = &vr.shaders[shaderID]
				vk.CmdBindPipeline(frame.cmds, vk.PIPELINE_BIND_POINT_GRAPHICS, shader.pipe)

				// setting scene uniforms for this shader
				vr.setSceneUniforms(shader, pass)
				vr.applySceneUniforms(shader)
			}

			// update material samplers
			if len(packet.TextureIDs) > 0 {
				matID, _ := vr.setMaterialSamplers(shader, packet.TextureIDs)
				vr.applyMaterialUniforms(shader, matID)
			}

			// bind model scope uniforms for this shader.
			vr.setModelUniforms(shader, packet)
			if packet.IsInstanced && packet.InstanceCount > 0 {
				// draw multiple models.
				vr.drawInstancedMesh(frame, packet.MeshID, packet.InstanceID, packet.InstanceCount, shader.attrs)
			} else {
				// draw one model.
				vr.drawMesh(frame, packet.MeshID, shader.attrs)
			}
		}
	}
	vk.CmdEndRenderPass(frame.cmds)

	// second pass always 2D if present.
	// then the 2D UI overlay render pass
	render2DInfo := vk.RenderPassBeginInfo{
		RenderPass:  vr.render2D,
		Framebuffer: vr.render2DFramebuffers[vr.imageIndex],
		RenderArea: vk.Rect2D{
			Offset: vk.Offset2D{X: 0, Y: 0},
			Extent: vk.Extent2D{Width: vr.frameWidth, Height: vr.frameHeight},
		},
	}
	vk.CmdBeginRenderPass(frame.cmds, &render2DInfo, vk.SUBPASS_CONTENTS_INLINE)
	if len(passes) > 1 && len(passes[Pass2D].Packets) > 0 {
		pass := passes[Pass2D]

		// draw 2D packets
		for _, packet := range pass.Packets {

			// change shader when necessary.
			if shaderID != packet.ShaderID {
				if packet.ShaderID >= uint16(len(vr.shaders)) {
					slog.Error("invalid shaderID", "shader_id", packet.ShaderID)
					continue
				}
				shaderID = packet.ShaderID // changing shaders.
				shader = &vr.shaders[shaderID]
				vk.CmdBindPipeline(frame.cmds, vk.PIPELINE_BIND_POINT_GRAPHICS, shader.pipe)

				// setting scene uniforms for this shader
				vr.setSceneUniforms(shader, pass)
				vr.applySceneUniforms(shader)
			}

			// update material samplers
			if len(packet.TextureIDs) > 0 {
				matID, _ := vr.setMaterialSamplers(shader, packet.TextureIDs)
				if lastMatID != matID {
					lastMatID = matID
				}
				vr.applyMaterialUniforms(shader, matID)
			}

			// bind model scope uniforms and draw the model.
			vr.setModelUniforms(shader, packet)
			vr.drawMesh(frame, packet.MeshID, shader.attrs)
		}
	}
	vk.CmdEndRenderPass(frame.cmds)

	// end command recording
	if err = vk.EndCommandBuffer(frame.cmds); err != nil {
		return fmt.Errorf("vk.EndCommandBuffer: %w", err)
	}
	return nil
}

var lastMatID uint32 = 345234545

func (vr *vulkanRenderer) endFrame(dt time.Duration) (err error) {
	frame := &vr.frames[vr.frameIndex]

	// reset the frames fence to unsignalled
	if err := vk.ResetFences(vr.device, []vk.Fence{frame.inFlightFence}); err != nil {
		return fmt.Errorf("vk.ResetFences: %w", err)
	}

	// submit the frame for render, waits for imageAvailable, signals renderComplete.
	submitInfo := vk.SubmitInfo{
		PWaitSemaphores:   []vk.Semaphore{frame.imageAvailable},
		PWaitDstStageMask: []vk.PipelineStageFlags{vk.PIPELINE_STAGE_COLOR_ATTACHMENT_OUTPUT_BIT},
		PSignalSemaphores: []vk.Semaphore{frame.renderComplete}, // signal present frame
		PCommandBuffers:   []vk.CommandBuffer{frame.cmds},
	}
	if err = vk.QueueSubmit(vr.graphicsQ, []vk.SubmitInfo{submitInfo}, frame.inFlightFence); err != nil {
		return fmt.Errorf("vk.QueueSubmit %w", err)
	}

	// present the frame, waits for renderComplete.
	presentInfo := vk.PresentInfoKHR{
		PWaitSemaphores: []vk.Semaphore{frame.renderComplete}, // wait for GPU render
		PSwapchains:     []vk.SwapchainKHR{vr.swapchain},
		PImageIndices:   []uint32{vr.imageIndex},
	}
	err = vk.QueuePresentKHR(vr.presentQ, &presentInfo)
	if err != nil {
		if err == vk.SUBOPTIMAL_KHR || err == vk.ERROR_OUT_OF_DATE_KHR {
			vr.resize(vr.frameWidth, vr.frameHeight)
			return nil // didn't quite work.
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
			vr.setUniform(shader, u, 0, pass.Uniforms[u.passUID])
		}
	}
}

// setModelUniforms copies the render packet data to model scope uniforms.
func (vr *vulkanRenderer) setModelUniforms(shader *vulkanShader, packet Packet) {
	for _, u := range shader.usets.index {
		if u.scope == load.ModelScope {
			vr.setUniform(shader, u, packet.MeshID, packet.Uniforms[u.packetUID])
		}
	}
}

// setUniform overwrites the uniform buffer with the given data.
// The data size must match the size expected by the given shader.
func (vr *vulkanRenderer) setUniform(shader *vulkanShader, u *uniform, instance uint32, data []byte) {
	switch u.scope {
	case load.SceneScope:
		offset := uintptr(vr.imageIndex*maxSceneUniformBytes + u.offset)
		dst := (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(shader.sceneUniformsMap)) + offset))
		copy(unsafe.Slice(dst, len(data)), data)
	case load.MaterialScope:
		// FUTURE: add material scope uniforms.... not currently needed by any shader.
		slog.Warn("the FUTURE is now -> add material scope uniforms")
		// offset := uintptr(vr.imageIndex*shader.maxMaterials*maxMaterialUniformBytes + instance*maxMaterialUniformBytes + u.offset)
		// dst := (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(shader.materialUniformsMap)) + offset))
		// copy(unsafe.Slice(dst, len(data)), data)
	case load.ModelScope:
		frame := vr.frames[vr.frameIndex]
		vk.CmdPushConstants(frame.cmds, shader.pipeLayout, vk.SHADER_STAGE_VERTEX_BIT|vk.SHADER_STAGE_FRAGMENT_BIT, u.offset, data)
	}
}
