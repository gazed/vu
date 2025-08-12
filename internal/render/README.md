
The `vk` directory is generated using `gazed/vk-gen` (based on `bbredesen/vk-gen`)

* get version 1.2 of the vk.xml specification (the latest spec breaks vk-gen)
```
curl https://raw.githubusercontent.com/KhronosGroup/Vulkan-Headers/v1.2.203/registry/vk.xml > vk.xml
```

* run the vk-gen command to get the windows bindings.
```
./vk-gen.exe -platform win32
```

* generate the enum strings.
```
cd vk;  go generate
```


# FUTURE wgpu

The `internal/wgpu` package would provide the webgpu bindings.
This could done by wrapping the javascript API exposed in the browser.
See: `https://github.com/mokiat/wasmgpu` (2 years with no changes)

Note that the webgpu spec is not released, and also the golang `syscall/js`,
used to wrap the javascript API, is not done either.

Could maybe generate the go code from the yml.
See `https://github.com/webgpu-native/webgpu-headers/blob/main/webgpu.yml`
The yml file is used to generate webgpu.h which is then used by dawn and webgpu-native.
Also see the pure rust webgpu implementation at `https://github.com/gfx-rs/wgpu`

Then render/wgpu.go would use the bindings to implement the engine render API.
This would only be for the `GOOS=js GOARCH=wasm`

## Working with webgpu

* build the wasm.
* GOOS=js GOARCH=wasm go build -o  ./web/main.wasm

* get the wasm shim from the go install
* setenv GOROOT="/c/Program Files/Go"
* cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" ./web

* run the server
* cd server; go run

* webgpu enabled by default in firefox on windows as of 2025jul.
* http://localhost:9090

