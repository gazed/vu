Internal contains bindings needed by the corresponding engine packages. 
These are not suitable for sharing as they are often mangled versions from
other open source projects.

- `audio/al`   - OpenAL audio bindings. - created following the device/win binding pattern. 
- `device/win` - WinAPI bindings, see: https://github.com/lxn/win
- `load/gltf`  - GLTF bindings.   see: https://github.com/qmuntal/gltf
- `render/vk`  - Vulkan bindings, see: https://github.com/bbredesen/go-vk

These have been included from their original projects for the following reasons:

- to minimize dependencies in the vu project.
- to create a copy in case something happens to the original project. 
- to delete code that is not used in `vu`.
- 
