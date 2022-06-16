// Copyright © 2024 Galvanized Logic Inc.

package load

import (
	"io/ioutil"
	"strings"
	"testing"
)

// go test -run Fnt
func TestFnt(t *testing.T) {
	SetAssetDir(".fnt", "../assets/fonts")
	fnt, err := Font("lucidiaSu18.fnt")
	if err != nil {
		t.Fatalf("glyph load failed %s", err)
	}
	if fnt.W != 256 || fnt.H != 256 || len(fnt.Chars) != 247 {
		t.Errorf("invalid font data: %d %d %d", fnt.W, fnt.H, len(fnt.Chars))
	}
}

// go test -run Glb
func TestGlb(t *testing.T) {
	SetAssetDir(".glb", "../assets/models")
	t.Run("no-glb", func(t *testing.T) {
		assets := Model("does_not_exist.glb")
		if len(assets) != 1 {
			t.Fatalf("glb load failed")
		}
		if assets[0].Err == nil {
			t.Errorf("expected file not found error %s", assets[0].Err)
		}
	})

	// should return
	t.Run("monkey0", func(t *testing.T) {
		assets := Model("monkey0.glb")
		if len(assets) != 2 {
			t.Fatalf("glb load failed")
		}

		// 0. mesh data
		md, ok := assets[0].Data.(MeshData)
		if !ok {
			t.Fatalf("exepcted mesh data")
		}
		if len(md[Vertexes].Data) != 23592 {
			t.Errorf("exepcted 23592 bytes got %d bytes", len(md[Vertexes].Data))
		}

		// 1. material data
		mat, ok := assets[1].Data.(PBRMaterialData)
		if !ok {
			t.Fatalf("exepcted PBR material data")
		}
		if mat.ColorB <= 0 {
			t.Errorf("exepcted a color")
		}
		if mat.Metallic <= 0 {
			t.Errorf("exepcted a metallic value")
		}
		if mat.Roughness <= 0 {
			t.Errorf("exepcted a roughness value")
		}
	})
}

func TestWav(t *testing.T) {
	SetAssetDir(".wav", "../assets/audio")
	snd, err := Audio("bloop.wav")
	if err != nil || int(snd.Attrs.DataSize) != len(snd.Data) {
		t.Errorf("wave load failed %s", err)
	}
}

func TestImage(t *testing.T) {
	SetAssetDir(".png", "../assets/images")
	img, err := Image("keyboard.png")
	if err != nil && len(img.Pixels) > 0 {
		t.Errorf("image load failed %s", err)
	}
}

// go test -run Shader
func TestShader(t *testing.T) {
	SetAssetDir(".shd", "../assets/shaders")

	// check if all the shaders can be loaded.
	t.Run("validate", func(t *testing.T) {
		files, _ := ioutil.ReadDir("../assets/shaders")
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".shd") {
				_, err := ShaderConfig(file.Name())
				if err != nil {
					t.Fatalf("fix shader %s %s", file.Name(), err)
				}
			}
		}
	})

	t.Run("tex3D", func(t *testing.T) {
		shd, err := ShaderConfig("tex3D.shd")
		if err != nil || shd.Name != "tex3D" || shd.Pass != "3D" {
			t.Errorf("shader configuration load failed %s", err)
		}
		if len(shd.Attrs) != 2 {
			t.Fatalf("expected 2 attributes got %d", len(shd.Attrs))
		}
		if shd.Attrs[1].AttrType != Texcoords {
			t.Fatalf("expected texcoord as second attribute")
		}
		if len(shd.Uniforms) != 4 {
			t.Fatalf("expected 4 uniforms got %d", len(shd.Uniforms))
		}
	})

	t.Run("pbr0", func(t *testing.T) {
		shd, err := ShaderConfig("pbr0.shd")
		if err != nil || shd.Name != "pbr0" || shd.Pass != "3D" {
			t.Errorf("shader configuration load failed %s", err)
		}
	})

	t.Run("bbinst", func(t *testing.T) {
		shd, err := ShaderConfig("bbinst.shd")
		if err != nil || shd.Name != "bbinst" || shd.Pass != "3D" {
			t.Errorf("shader configuration load failed %s", err)
		}
		if len(shd.Attrs) != 5 {
			t.Fatalf("expected 5 attributes got %d", len(shd.Attrs))
		}
		if shd.Attrs[2].AttrType != InstanceLocus ||
			shd.Attrs[3].AttrType != InstanceColors ||
			shd.Attrs[4].AttrType != InstanceScales {
			t.Fatalf("invalid instance attribute type")
		}
	})
}

// dumpMesh can be used to print mesh data for debugging.
func dumpMesh() {
	assets := Model("box0.glb")
	md, ok := assets[0].Data.(MeshData)
	if !ok {
		return
	}
	md[Vertexes].PrintF32("Vertexes")
	md[Normals].PrintF32("Normals")
	md[Texcoords].PrintF32("Texcoords")
	md[Indexes].PrintU16("Indexes")
}
