// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

// Package load fetches disk based data that will be used for 3D assets.
// Data is loaded directly from disk for development builds and from a zip
// file attached to the binary for production builds.
//
// Data that can be loaded from disk is listed in the Loader interface.
// Data is returned in an intermediate format that is used to populate 3D
// render or audio based assets as follows:
//      Data                      File            Example Consumer
//     ------                    ------          -------------------
//    bitmapped fonts          : txtfile.fnt --> render.Model
//    colour and surface data  : txtfile.mtl --> render.Model
//    vertex data              : txtfile.obj --> render.Mesh
//    vertex shader program    : txtfile.vsh -┐
//    fragment shader program  : txtfile.fsh --> render.Shader
//    animated models          : txtfile.iqe --> render.Model
//    animated models          : binfile.iqm --> render.Model
//    images                   : binfile.png --> render.Texture
//    audio                    : binfile.wav --> audio.Sound
//
// Package load is currently intended for smaller 3D applications where data
// is loaded directly from files to memory, i.e. no database involved.
//
// Package load is provided as part of the vu (virtual universe) 3D engine.
package load

// Design Notes:
// FUTURE: wrap or develop more import formats. See the Open Asset Import Library:
//         http://assimp.sourceforge.net/
// FUTURE: Have more than one supported file type for a given resource.
// FUTURE: Industrial strength (database) back end?

import (
	"archive/zip"
	"image"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

// Loader provides methods for loading disk based data assets. Loader methods
// log development errors for unknown assets or unsupported data types.
// Loader files will return empty or nil data values when there are errors.
type Loader interface {

	// SetDir overrides the default directory location for the given asset type.
	// All directories are expected to be relative to the application location.
	SetDir(assetType int, dir string) Loader
	Dispose() // Properly terminate asset loading

	// Supported file formats.
	Png(name string) (img image.Image, err error)         // .png
	Mtl(name string) (mtl *MtlData, err error)            // .mtl
	Obj(name string) (obj []*ObjData, err error)          // .obj
	Fnt(name string) (fnt *FntData, err error)            // .fnt
	Vsh(name string) (src []string, err error)            // .vsh
	Fsh(name string) (src []string, err error)            // .fsh
	Wav(name string) (wh *WavHdr, data []byte, err error) // .wav
	Iqe(name string) (iqd *IqData, err error)             // .iqe
	Iqm(name string) (iqd *IqData, err error)             // .iqm

	// GetResource allows applications to include and find custom resources.
	GetResource(directory, name string) (file io.ReadCloser, err error)
}

// Asset directory identifiers for SetDir.
const (
	img = iota // Font and texture images.
	mod        // Model meshes and materials.
	snd        // Audio.
	src        // Font mapping files and shader source.
)

// NewLoader provides the default loader implmentation.
func NewLoader() Loader { return newLoader() }

// Loader interface
// ===========================================================================
// loader is the default Loader implementation.

// loader provides functions to assist getting asset data from disk into representations
// that are easily communicated to the audio and graphics layer. Loader is input only.
// Asset files are expected to be created by 3rdParty tools like Blender or Gimp.
type loader struct {
	reader *zip.ReadCloser // Used as the resource file if set. Otherwise use the file system.
	dir    map[int]string  // Data directory locations.
}

// newLoader creates the appropriate asset loader. Assets are in a zip
// file that is either included within the production binary or in a asset
// directory relative to the executable. Development builds have a nil
// loader.reader and will look locally on disk.
func newLoader() *loader {
	var resources *zip.ReadCloser // packaged resources.
	programName := os.Args[0]     // qualified path to executable
	resourceZip := path.Join(path.Dir(programName), "../Resources/resources.zip")
	if reader, err := zip.OpenReader(resourceZip); err == nil {
		resources = reader // the creator must call loader.dispose()
	} else if reader, err := zip.OpenReader(programName); err == nil {
		resources = reader // the creator must call loader.dispose()
	}
	l := &loader{reader: resources}
	l.dir = map[int]string{
		mod: "models",
		snd: "audio",
		src: "source",
		img: "images",
	}
	return l
}

// Comply with the Loader interface.
func (l *loader) Wav(name string) (wh *WavHdr, data []byte, err error) { return l.wav(name) }
func (l *loader) Png(name string) (img image.Image, err error)         { return l.png(name) }
func (l *loader) Fnt(name string) (fnt *FntData, err error)            { return l.fnt(name) }
func (l *loader) Vsh(name string) (src []string, err error)            { return l.txt(name + ".vsh") }
func (l *loader) Fsh(name string) (src []string, err error)            { return l.txt(name + ".fsh") }
func (l *loader) Mtl(name string) (mtl *MtlData, err error)            { return l.mtl(name) }
func (l *loader) Obj(name string) (obj []*ObjData, err error)          { return l.obj(name) }
func (l *loader) Iqe(name string) (iqd *IqData, err error)             { return l.iqe(name) }
func (l *loader) Iqm(name string) (iqd *IqData, err error)             { return l.iqm(name) }
func (l *loader) SetDir(dataType int, dir string) Loader               { return l.setDir(dataType, dir) }
func (l *loader) Dispose()                                             { l.dispose() }

// Expose the resource location ability in the Loader interface.
func (l *loader) GetResource(directory, name string) (file io.ReadCloser, err error) {
	return l.getResource(directory, name)
}

// dispose properly terminates the loader.
// This is only needed when the loader has been reading resources from a file.
func (l *loader) dispose() {
	if l.reader != nil {
		l.reader.Close()
	}
}

// setDir is used if a test case or application wishes to override
// one or more of the default resource directory locations.
func (l *loader) setDir(dataType int, dir string) *loader {
	switch dataType {
	case mod, src, snd, img:
		l.dir[dataType] = dir
	default:
		log.Printf("loader.setDir: unknown resource type")
	}
	return l
}

// getResource locates the named resource.  This is expected to be used either
// in production where the resources have been included with the application,
// or development where the resources are on disk in the local directory.
//
// The caller is responsible for closing the returned file.
func (l *loader) getResource(directory, name string) (file io.ReadCloser, err error) {
	filePath := strings.TrimSpace(path.Join(directory, name))
	if l.reader != nil {
		for _, resource := range l.reader.File {
			if filePath == resource.Name {
				rc, zerr := resource.Open()
				if zerr != nil {
					log.Printf("Could not open resource %s: %s", resource.Name, zerr)
					return nil, zerr
				}
				return rc, nil
			}
		}
	}
	return os.Open(filePath)
}
