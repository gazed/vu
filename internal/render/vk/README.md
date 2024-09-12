# go-vk - Vulkan 1.3 supporting Windows and Mac

***This package is in a beta state right now!*** It has been tested on Windows and Mac. Please report any bugs you find!

go-vk is a Go-langauge (and Go-style) binding around the Vulkan graphics API. Rather than just slapping a Cgo wrapper
around everything, Vulkan's functions, structures and other types have been translated to a Go-style API. For example, 
"native" Vulkan returns any resources you request in pointers your program passes into Vulkan. This allows
Vulkan to (generally) return a VkResult success or error code from the C function call. However, in Go, we have the
luxury of multiple return values, so this:

```C
VkInstance myInstance;
Result r = vkCreateInstance(&instanceCI, NULL, &myInstance);
if (r != VK_SUCCESS) {
    // Handle an error
}
// Use the instance handle
```

Becomes this:
```go
instance, err := vk.CreateInstance(&instanceCI, nil)
if err != nil {
    panic("Could not create a Vulkan instance!") // Don't panic
}
```

Likewise, the "Enumerate" group of functions returning an array of values in C require a call, an error check, an
allocation, another function call, and another error check:
```C
int deviceCount;
Result res = vkEnumeratePhysicalDevices(myInstance, &deviceCount, NULL);
if (res != VK_SUCCESS) { // Check the result, of course
    // handle the error
}
// ...and you really should also check that deviceCount > 0
if (deviceCount == 0) {
    // gracefully exit, since there are no GPU devices actually available on this machine
}

VkPhysicalDevice devices[deviceCount];

res = vkEnumeratePhysicalDevices(myInstance, &deviceCount, devices);
if (res != VK_SUCCESS) { // Check the result again
    // handle the error
}
// Now do something with the devices and make sure you hold on to deviceCount 
// so you don't go beyond the bounds of the array...
for (int i = 0; i < deviceCount; i++) {
    // Check device suitability, select a device, and hold on to that handle...
}
```

Yuck. Here's the same code in Go:
```go
if devices, err := vk.EnumeratePhysicalDevices(myInstance); err != nil {
    // handle the error
} else {
    // devices is a slice of vk.PhysicalDevice. Nice!
}
```

But there's more! Passing multiple values to a Vulkan command requires a pointer and count parameter, and sometimes
that count parameter is embedded in another struct. You can make life a little easier with C++'s `std::vector`. 
For example, specifying requested extensions at instance creation:

```C++
std::vector<const char*> requiredExtensions = {
    VK_KHR_SWAPCHAIN_EXTENSION_NAME, VK_KHR_SURFACE_EXTENSION_NAME
};

VkInstanceCreateInfo createInfo{};
createInfo.sType = VK_STRUCTURE_TYPE_INSTANCE_CREATE_INFO;
// Other create info props...

// set the size
createInfo.enabledExtensionCount = static_cast<uint32_t>(extensions.size());
// extract the data pointer from the vector
createInfo.ppEnabledExtensionNames = extensions.data();
```

versus:

```go
requiredExtensions := []string{vk.KHR_SWAPCHAIN_EXTENSION_NAME, vk.KHR_SURFACE_EXTENSION_NAME}

createInfo := vk.InstanceCreateInfo{
    // No structure type, no length member, and no pointer required.
    // Just assign the slice, or even instantiate it inline
    EnabledExtensionNames: requiredExtensions,
}
```

## Code Generation
This codebase is (almost) entirely generated from a `vk.xml` file by the [vk-gen](https://github.com/bbredesen/vk-gen)
tool. Updating go-vk for a new Vulkan version should be as easy as downloading the new vk.xml file from Khronos and
executing vk-gen. **This repository does not get direct modifications!** Any bug fixes or new features need to be made in
`vk-gen`, which will then be used re-generate this code base.

## Usage

Ensure that your GPU supports Vulkan and that a Vulkan library is installed in your system-default library location
(e.g., C:\windows\system32\vulkan-1.dll on Windows). This package uses Cgo to call Vulkan, so it needs to be enabled in 
your Go settings.

`$ go get github.com/bbredesen/go-vk@latest`

Builds for Vulkan API versions 1.1, 1.2, 1.3 (and future releases) will be tagged as releases of go-vk with matching
version numbers, if you want to use a specific version of the API. go-vk does not itself require the Vulkan SDK be installed,
as it reads symbols from the system-default Vulkan library at runtime. However, you will need the SDK installed to use
validation layers, shader compilers, etc. during development.

```go main.go
package main

import (
    "github.com/bbredesen/go-vk"
)
// Notice that you don't need to alias the import, it is already bound to the "vk" namespace

func main() {
    if encodedVersion, err := vk.EnumerateInstanceVersion(); err != nil {
        // Returned errors are vk.Results. You can directly compare err those 
        // predefined values to determine which error occured.
        // The string returned by Error() is the name of the code. For example,
        // vk.ERROR_OUT_OF_DATE_KHR.Error() == "ERROR_OUT_OF_DATE_KHR"
        fmt.Printf("EnumerateInstanceVersion failed! Error code was %s\n", err.Error())
        os.Exit(1)
    } else {
        fmt.Printf("Installed Vulkan version: %d.%d.%d\n", 
            vk.API_VERSION_MAJOR(encodedVersion), 
            vk.API_VERSION_MINOR(encodedVersion), 
            vk.API_VERSION_PATCH(encodedVersion),
        )
    }

    // Also notice that you don't need to set the StructureType field on your Go structs. 
    // In fact, the sType field doesn't even exist on the public side of the binding...it is automatically
    // added when you pass your struct through to a command.
    appInfo := vk.ApplicationInfo{
		ApplicationName:    "Example App",
		ApplicationVersion: vk.MAKE_VERSION(1, 0, 0),
		EngineVersion:      vk.MAKE_VERSION(1, 0, 0),
		ApiVersion:         vk.MAKE_VERSION(1, 3, 0),
	}

	icInfo := vk.InstanceCreateInfo{
		ApplicationInfo:       appInfo,
        // Extension names are built into the binding as const strings.
		EnabledExtensionNames: []string{vk.KHR_SURFACE_EXTENSION_NAME, vk.KHR_WIN32_SURFACE_EXTENSION_NAME},
        // Layer names are not built in, unfortunately...layers are not part of the core API spec and names are not present in vk.xml
		EnabledLayerNames:     []string{"VK_LAYER_KHRONOS_validation"},
	}

	instance, err := vk.CreateInstance(&icInfo, nil)
    // vk.SUCCESS is defined as nil, so you can also check for an error like this if preferred.
    if err != vk.SUCCESS {
        fmt.Printf("Failed to create Vulkan instance, error code was %s\n", err.Error())
        if err == vk.ERROR_INCOMPATIBLE_DRIVER { 
            /* ... */
        }
    }
    fmt.Printf("Vulkan instance created, handle value is 0x%x\n", instance)

    // Clean up after yourself before exiting!
    vk.DestroyInstance(instance)
}
```

`$ go run main.go`

A number of code samples and working demos, including an implementation of the excellent tutorial program from
[vulkan-tutorial.com](https://vulkan-tutorial.com), are available at [go-vk-samples](https://github.com/bbredesen/go-vk-samples)

## Library Structure

The Vulkan API is defined through a set of type categories, each of which has a corresponding source file in go-vk.
Thus, you will find all structs defined in struct.go, all commands defined in command.go, etc. Where
platform-specific types are neccessary, they are defined in separate files with appropriate go:build tags. The
`stringify` tool has also been run against enumerated types, so if `result == vk.NOT_READY` then `result.String() == "NOT_READY"`.

The underlying Vulkan implementation is actually accessed through a small Cgo wrapper, found in static_common.go; go-vk
opens the shared library and lazy-loads any requested symbols. All of the public-facing structs in Go are translated to
the appropriate memory layout before being passed through to the API, via each struct's Vulkanize() function.
Vulkanize()'s primary purpose is to convert slices to a length and pointer field in the internal struct, Go strings to
null-terminated byte pointers, and to recursively Vulkanize any non-primitive members.

The structs also have a Goify function to do the reverse: create slices
from a length and pointer field and create strings from null-terminated byte arrays. In practice, this is only used for
structs that are returned by the API, but Goify is implemented on all structs.

Note that you should never need to directly call Vulkanize() or Goify() (with one expection, noted below). Conversions are
automatically handled in the background when you call a Vulkan command. 

### Extended Structs
If you use pNext to extend any structures, you will need to manually build the chain by calling Vulkanize() and setting the returned pointer in the
base struct.

```go
instanceCI := vk.InstanceCreateInfo{
    // ...
}

validationFeatures := vk.ValidationFeaturesEXT{
    PEnabledValidationFeatures: []vk.ValidationFeatureEnableEXT{vk.VALIDATION_FEATURE_ENABLE_BEST_PRACTICES_EXT}
    // ...
}

instanceCI.PNext = unsafe.Pointer(validationFeatures.Vulkanize())
```

Leaving these as unsafe.Pointers was the simplest implementation to get the binding up and running. The next level of
implementation is to define pNext as a Vulkanizer interface type, and have Vulkanize build the chain. I've also
considered more specific interfaces flagged with empty functions,
since the spec does indicate for each struct with what other structs it extends (e.g., VkValidationFlagsEXT has a
structextends="VkInstanceCreateInfo" attribute).

## Mapped memory and copying data

Any practical Vulkan application will need to copy raw data between the CPU and GPU...loads to uniform buffers, texture
data, etc. are exposed through vkMapMemory. Unfortunately for us, Go is designed to avoid directly
managing and copying memory. To handle this, three specific utility functions are included with go-vk: MemCopySlice,
MemCopyObj, and MemCopy.

The first two two functions use generics to copy your data byte-for-byte to Vulkan in an abstract way, so Go 1.18 or higher is a
requirement.

The MemCopy function that accepts two unsafe.Pointers and a number of bytes to copy, but it is recommended
that you use MemCopyObj or MemCopySlice instead. It is really only offered in case you need to target a Go version less
than 1.18 (and hence do not have access to generics). In that case you could vendor a copy of go-vk in your project and
delete the two generic functions.

**There are no guardrails on any of these functions! You, the developer, are repsonbile for allocating enough memory
at the destination before calling them.** 

They do not (and cannot) check how much space is available behind the pointer you give them. Under the hood, they create "fake"
byte slices at the destination pointer and the source pointer or at the head of the input slice. It then uses Go's copy macro
to copy the data over.

Go version 1.20 includes some new functions in the unsafe package for copying slices to pointers, allowing you to "cast"
between pointers and slices and use the `copy()` macro. The MemCopy functions above were written before the 1.20 release and do
effectively the same thing. You are free to use whichever method you prefer.

In Go 1.20+, this:
```go
ptr, err := vk.MapMemory(/* ... */)

sl := unsafe.Slice((*VertexFormat)(ptr), len(vertices))
copy(sl, vertices)
```
...is functionally the same as this:
```go
ptr, err := vk.MapMemory(/* ... */)

vk.MemCopySlice(ptr, vertices)
```

## A note on unions

Vulkan includes a small number of C-union types, VkClearValue and VkClearColorValue probably being the most commonly used.
However, Go does not have any concept of unions in the language. In go-vk, those unions are implemented as a struct
containing all of the members of the union, which is resolved behind the scenes to the correct member. You will need to
set the field you intend to use by calling the `As<FieldName>` method on those structs. The struct's Vulkanize() method will
then extract the correct member for passing into the Vulkan API.

```go
var ccv vk.ClearColorValue
ccv.AsTypeFloat32(float32[4]{0.0, 0.0, 0.0, 1.0})
// The spec names this field float32, which is a reserved word in Go. vk-gen 
// renames these fields to TypeFloat32, TypeInt32, etc. to avoid any conflicts.
```

## Examples

See the [go-vk-samples](https://github.com/bbredesen/go-vk-samples) repo for a number of working Vulkan samples using
this library. The samples currently run on Windows and Mac.

## Known Issues

* VkAccelerationStructureMatrixMotionInstanceNV - embedded bit fields in uint32_t are not handled at all...this
  structure will not behave as intended and will likely cause a crash if used.
* H.264 and H.265 commands and types are almost certainly broken. Vulkan does provide a separate XML file in the vk.xml format for those
  types, but reading that file has not yet been implemented in vk-gen. As a placeholder, all of these types are defined
  as int32 through exceptions.json.
* The union type VkPipelineExecutableStatisticValueKHR is returned from Vulkan through VkPipelineExecutableStatisticKHR.
  Returned unions are not supported and there is no Goify() function associated. VkPipelineExecutableStatisticKHR is
  returned to the developer without the Value member populated.