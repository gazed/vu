// Copyright © 2024 Galvanized Logic Inc.

package vu

import (
	"math"
	"sort"
	"testing"

	"github.com/gazed/vu/load"
	"github.com/gazed/vu/render"
)

// go test -run Scene
func TestScene(t *testing.T) {
	rc := &mrc{} // mock render context.

	t.Run("scene with lights", func(t *testing.T) {
		app := newApplication()
		app.ld.loadDefaultAssets(rc) // direct loads (no goroutine)
		scene := app.addScene(Scene3D)
		scene.AddLight(DirectionalLight).SetLight(0.8, 0.8, 0.8, 10).SetAt(1, 2, 3)
		scene.AddLight(PointLight).SetLight(1.0, 0.1, 0.1, 10).SetAt(5, 4, 3)

		// check passes
		passes := app.scenes.getFrame(app, app.frame)
		if len(passes) != 2 {
			t.Errorf("expected 2 render passes, got %d", len(passes))
		}

		// check light counts.
		lightCount := passes[render.Pass3D].Uniforms[load.NLIGHTS][0]
		if lightCount != 2 {
			t.Errorf("expected 2-3D render lights, got %d", lightCount)
		}
	})

	t.Run("scene with three models", func(t *testing.T) {
		app := newApplication()
		app.ld.loadDefaultAssets(rc) // direct loads (no goroutine)
		scene := app.addScene(Scene3D)
		scene.AddModel("shd:icon", "msh:cube", "tex:color:test")
		scene.AddModel("shd:icon", "msh:icon", "tex:color:test")
		scene.AddModel("shd:icon", "msh:quad", "tex:color:test")

		// check passes
		passes := app.scenes.getFrame(app, app.frame)
		if len(passes) != 2 {
			t.Errorf("expected 2 render passes, got %d", len(passes))
		}

		// check packets
		packetCount := len(passes[render.Pass3D].Packets)
		if packetCount != 3 {
			t.Errorf("expected 3-3D render packets, got %d", packetCount)
		}
	})
}

// mock render context.
type mrc struct{}

func (rc *mrc) LoadTexture(img *load.ImageData) (uint32, error)           { return 0, nil }
func (rc *mrc) UpdateTexture(tid uint32, img *load.ImageData) (err error) { return nil }
func (rc *mrc) LoadMesh(load.MeshData) (uint32, error)                    { return 0, nil }
func (rc *mrc) LoadMeshes([]load.MeshData) ([]uint32, error)              { return []uint32{0}, nil }
func (rc *mrc) LoadShader(config *load.Shader) (uint16, error)            { return 0, nil }

// go test -run Bucket
func TestBucket(t *testing.T) {

	// Check the bits placement
	t.Run("set bucket", func(t *testing.T) {
		b := setBucketShader(setBucketDistance(newBucket(render.Pass2D), 2.4), 21)
		b = setBucketLayer(b, 7)

		// pull out the values to check placement.
		dist := math.Float32frombits(uint32(b & 0x00000000FFFFFFFF))
		draw := uint16((b & 0x000F000000000000) >> 48) // drawOpaque == 4
		shid := uint16((b & 0x0000FFFF00000000) >> 40)
		pass := uint8((b & 0xFF00000000000000) >> 56)
		layer := uint16((b & 0x00F0000000000000) >> 52)
		if pass != 254 || layer != 7 || draw != 4 || shid != 21 || dist != 2.4 {
			t.Errorf("bad bucket %016x %d %d %d %d %f\n", b, pass, draw, layer, shid, dist)
		}
	})

	// check sorting
	t.Run("bucket sort", func(t *testing.T) {
		packets := render.Packets{
			render.Packet{Bucket: 256},
			render.Packet{Bucket: 10},
			render.Packet{Bucket: 0},
			render.Packet{Bucket: 1025},
			render.Packet{Bucket: 20},
		}

		// sort the render packets.
		sort.SliceStable(packets, func(i, j int) bool {
			return packets[i].Bucket < packets[j].Bucket
		})

		if packets[0].Bucket != 0 || packets[4].Bucket != 1025 {
			t.Errorf("expected low to high %+v\n", packets)
		}
	})
}
