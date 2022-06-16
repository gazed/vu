// Copyright © 2024 Galvanized Logic Inc.

package vu

import (
	"fmt"
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

		passes := app.scenes.getFrame(app, app.frame)
		if len(passes) != 2 {
			t.Errorf("expected 2 render passes, got %d", len(passes))
		}
		lightCount := passes[render.Pass3D].Data[load.NLIGHTS][0]
		if lightCount != 2 {
			t.Errorf("expected 2-3D render lights, got %d", lightCount)
		}
	})

	t.Run("scene with three models", func(t *testing.T) {
		app := newApplication()
		app.ld.loadDefaultAssets(rc) // direct loads (no goroutine)
		scene := app.addScene(Scene3D)
		scene.AddModel("shd:tex3D", "msh:cube", "tex:color:test")
		scene.AddModel("shd:tex3D", "msh:icon", "tex:color:test")
		scene.AddModel("shd:tex3D", "msh:quad", "tex:color:test")

		if len(app.models.ready) != 3 {
			t.Errorf("expected 3 ready models got %d", len(app.models.ready))
		}

		passes := app.scenes.getFrame(app, app.frame)
		if len(passes) != 2 {
			t.Errorf("expected 2 render passes, got %d", len(passes))
		}

		packetCount := len(passes[render.Pass3D].Packets)
		if packetCount != 3 {
			t.Errorf("expected 3-3D render packets, got %d", packetCount)
		}
	})
}

// mock render context.
type mrc struct{}

func (rc *mrc) LoadTexture(img *load.ImageData) (uint32, error) { return 0, nil }
func (rc *mrc) LoadMesh(load.MeshData) (uint32, error)          { return 0, nil }
func (rc *mrc) LoadShader(config *load.Shader) (uint16, error)  { return 0, nil }

// go test -run Bucket
func TestBucket(t *testing.T) {

	// Check the bits placement
	t.Run("set bucket", func(t *testing.T) {
		b := newBucket(render.Pass2D).setShaderID(21).setDistance(2.4)

		// pull out the values to check placement.
		dist := math.Float32frombits(uint32(b & 0x00000000FFFFFFFF))
		shid := uint16((b & 0x0000FFFF00000000) >> 40)
		draw := uint16((b & 0x00FF000000000000) >> 48) // drawOpaque == 4
		pass := uint8((b & 0xFF00000000000000) >> 56)
		if pass != 254 || draw != 4 || shid != 21 || dist != 2.4 {
			t.Errorf("bad bucket %016x %d %d %d %f\n", b, pass, draw, shid, dist)
		}
	})

	// check that the draw packets are properly sorted.
	t.Run("check sort", func(t *testing.T) {
		packets := []render.Packet{
			// 2D
			{Tag: 11, Bucket: uint64(newBucket(render.Pass2D).setShaderID(10).setType(drawTransparent))},
			{Tag: 10, Bucket: uint64(newBucket(render.Pass2D).setShaderID(11).setType(drawOpaque))},
			{Tag: 8, Bucket: uint64(newBucket(render.Pass2D).setShaderID(12).setType(drawOpaque))},
			{Tag: 9, Bucket: uint64(newBucket(render.Pass2D).setShaderID(12).setType(drawOpaque))},
			{Tag: 12, Bucket: uint64(newBucket(render.Pass2D).setShaderID(10).setType(drawTransparent))},
			// 3D
			{Tag: 7, Bucket: uint64(newBucket(render.Pass3D).setShaderID(4).setType(drawTransparent).setDistance(11))},
			{Tag: 3, Bucket: uint64(newBucket(render.Pass3D).setShaderID(5).setDistance(12))},
			{Tag: 5, Bucket: uint64(newBucket(render.Pass3D).setShaderID(2).setDistance(13))},
			{Tag: 2, Bucket: uint64(newBucket(render.Pass3D).setShaderID(5).setDistance(22))},
			{Tag: 4, Bucket: uint64(newBucket(render.Pass3D).setShaderID(2).setDistance(23))},
			{Tag: 1, Bucket: uint64(newBucket(render.Pass3D).setShaderID(1).setType(drawSky))},
			{Tag: 6, Bucket: uint64(newBucket(render.Pass3D).setShaderID(4).setType(drawTransparent).setDistance(14))},
		}

		sort.SliceStable(packets, func(i, j int) bool { return packets[i].Bucket > packets[j].Bucket })

		order := ""
		for i := range packets {
			order += fmt.Sprintf("[%d:%d]", i+1, packets[i].Tag)
		}
		expect := "[1:1][2:2][3:3][4:4][5:5][6:6][7:7][8:8][9:9][10:10][11:11][12:12]"
		if order != expect {
			t.Errorf("expected %s got %s", expect, order)
		}
	})
}
