
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
