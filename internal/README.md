Internal contains external code and/or bindings used in the corresponding engine packages.
These are not suitable for sharing as they are often mangled versions from
other open source projects and/or contain CGO based binding code.

- `audio/al`   - OpenAL audio bindings. - created following the device/win binding pattern.
- `device/win` - WinAPI bindings, see: https://github.com/lxn/win
- `load/gltf`  - GLTF bindings.   see: https://github.com/qmuntal/gltf
- `render/vk`  - Vulkan bindings, see: https://github.com/bbredesen/go-vk

Some or all of the above projects have been copied from their original projects in order to:

- minimize dependencies in the vu project.
- create a copy in case something happens to the original project.
- delete code that is not used in `vu`.
