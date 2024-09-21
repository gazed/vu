// Copyright © 2024 Galvanized Logic Inc.

package vu

import (
	"testing"
	"time"

	"github.com/gazed/vu/audio"
	"github.com/gazed/vu/load"
)

// go test -run Loader
// verify loader can create assets from file data.
func TestLoader(t *testing.T) {
	ld := newLoader()                // start goroutine.
	rc := &loaderTestRenderContext{} // mock render context.
	ac := &loaderTestAudioContext{}  // mock audio context.

	// queue the asset load requests
	files := []string{
		"core.png", "22:hack.ttf", "bloop.wav", "box0.glb", "tex3D.shd",
	}
	ld.importAssetData(files...)
	ld.importAssetData(files...) // check that duplicate requests are ignored

	// check if the loader goroutine hangs or panics.
	t.Run("load asset files", func(t *testing.T) {

		// loop to process the asset load responses from the loader goroutine.
		// monkey_solid generates 2 assets, so expect 7.
		expectedAssets, loadedAssets := 7, 0
		for i := 0; i < 200; i++ {
			loadedAssets += ld.loadAssets(rc, ac)
			if loadedAssets >= expectedAssets {
				ld.dispose() // close the worker input channel
				break
			}
			time.Sleep(20)
		}

		if len(ld.loaded) != 5 {
			t.Fatalf("expected loaded with 5 files got %v", ld.loaded)
		}
		for filename, v := range ld.loaded {
			if !v {
				t.Errorf("expected file %s to be marked as loaded", filename)
			}
		}
		if loaderTestTextureLoads != 2 {
			t.Errorf("expected 2 texture loads got %d", loaderTestTextureLoads)
		}
		if loaderTestMeshLoads != 1 {
			t.Errorf("expected 1 mesh loads got %d", loaderTestMeshLoads)
		}
		if loaderTestShaderLoads != 1 {
			t.Errorf("expected 1 shader loads got %d", loaderTestShaderLoads)
		}
		if loaderTestSoundLoads != 1 {
			t.Errorf("expected 1 audio loads got %d", loaderTestSoundLoads)
		}
	})

	// the loader workers are closed at this point
	// so any new file requests will panic.

	t.Run("request loaded assets", func(t *testing.T) {
		id1, id2 := eID(1), eID(2)
		ld.getAsset(assetID(tex, "core"), id1, loaderTestCallback)
		ld.getAsset(assetID(tex, "core"), id2, loaderTestCallback)
		ld.getAsset(assetID(msh, "box0"), id1, loaderTestCallback)
		ld.getAsset(assetID(shd, "tex3D"), id1, loaderTestCallback)
		ld.getAsset(assetID(aud, "bloop"), id1, loaderTestCallback)
		if loaderTestID1Callbacks != 4 {
			t.Errorf("entity 1 expected 4 callbacks, got %d", loaderTestID1Callbacks)
		}
		if loaderTestID2Callbacks != 1 {
			t.Errorf("entity 2 expected 1 callback, got %d", loaderTestID2Callbacks)
		}
	})

	t.Run("request unloaded assets", func(t *testing.T) {
		loaderTestID1Callbacks = 0
		loaderTestID2Callbacks = 0
		id1, id2 := eID(1), eID(2)
		ld.getAsset(assetID(msh, "core"), id1, loaderTestCallback) // invalid aid
		ld.getAsset(assetID(shd, "test"), id2, loaderTestCallback) // default assets not loaded
		if loaderTestID1Callbacks != 0 || loaderTestID2Callbacks != 0 {
			t.Error("expected 0 callbacks")
		}
		if len(ld.requests) != 2 {
			t.Errorf("expected 2 outstanding assets requests, got %d", len(ld.requests))
		}
	})
}

// track asset load callbacks to entities.
var loaderTestID1Callbacks = 0
var loaderTestID2Callbacks = 0

// mock the asset load complete callback
func loaderTestCallback(eid eID, a asset) {
	switch eid {
	case 1:
		loaderTestID1Callbacks += 1
	case 2:
		loaderTestID2Callbacks += 1
	}
}

// track asset loads to render and audio contexts.
var loaderTestTextureLoads = 0
var loaderTestTextureUpdates = 0
var loaderTestMeshLoads = 0
var loaderTestShaderLoads = 0
var loaderTestSoundLoads = 0

// mock the render.Load interface expected by the loader.
type loaderTestRenderContext struct{}

func (rc *loaderTestRenderContext) LoadTexture(img *load.ImageData) (tid uint32, err error) {
	loaderTestTextureLoads += 1
	return 0, nil
}
func (rc *loaderTestRenderContext) UpdateTexture(tid uint32, img *load.ImageData) (err error) {
	loaderTestTextureUpdates += 1
	return nil
}
func (rc *loaderTestRenderContext) LoadMesh(load.MeshData) (mid uint32, err error) {
	loaderTestMeshLoads += 1
	return 0, nil
}
func (rc *loaderTestRenderContext) LoadMeshes([]load.MeshData) (mids []uint32, err error) {
	loaderTestMeshLoads += 1
	return []uint32{}, nil
}

func (rc *loaderTestRenderContext) LoadShader(config *load.Shader) (tid uint16, err error) {
	loaderTestShaderLoads += 1
	return 0, nil
}

// mock the audio.Load interface expected by the loader.
type loaderTestAudioContext struct{}

func (ac *loaderTestAudioContext) LoadSound(sound, buff *uint64, d *audio.Data) error {
	loaderTestSoundLoads += 1
	return nil
}
