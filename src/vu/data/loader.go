// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package data

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

// loader provides functions to assist getting data from disk into representations
// that are easily communicated to the graphics layer (resources). Examples of
// loading resources are:
//    Load material data from various formats (.mtl)
//    Load model (mesh) data from various formats (.obj)
//    Load images and texture data from files (.png)
//    Load sound from audio files (.wav)
//    Load bitmap fonts (glpyh) information from files (.fnt)
//    Load vertex and fragment shader programs from source files (.vsh, .fsh)
//
// loader is input only (no export). The resources files are expected to be created
// by 3rdParty tools like Blender or Gimp.
type loader struct {
	reader *zip.ReadCloser // Used as the resource file if set. Otherwise use the file system.
	dir    map[int]string  // Resource directory locations.
}

// Types of resources that are used both in loader and depot.
const (
	gly = iota // glyphs
	mat        // material
	msh        // mesh
	shd        // shader
	snd        // sound
	tex        // texture - checked as last const.
)

// newLoader creates the appropriate resource loader.  Resources are in a zip
// file that is either included within the production binary or in a Resource
// directory relative to the executable.  Development builds have a nil
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
		gly: "images",
		mat: "models",
		msh: "models",
		shd: "shaders",
		snd: "audio",
		tex: "images",
	}
	return l
}

// dispose properly terminates the loader.
// This is only needed when the loader has been reading resources from a file.
func (l loader) dispose() {
	if l.reader != nil {
		l.reader.Close()
	}
}

// load attempts to find the named resource on disk and turn it into the
// appropriate resource data type.
func (l *loader) load(name string, data interface{}) {
	switch data.(type) {
	case **Glyphs:
		glyphset, _ := data.(**Glyphs)
		*glyphset = l.glyphset(l.dir[gly], name)
	case **Material:
		material, _ := data.(**Material)
		*material = l.material(l.dir[mat], name)
	case **Mesh:
		mesh, _ := data.(**Mesh)
		*mesh = l.mesh(l.dir[msh], name)
	case **Shader:
		shader, _ := data.(**Shader)
		*shader = l.shader(l.dir[shd], name)
	case **Sound:
		sound, _ := data.(**Sound)
		*sound = l.sound(l.dir[snd], name)
	case **Texture:
		texture, _ := data.(**Texture)
		*texture = l.texture(l.dir[tex], name)
	default:
		log.Printf("loader.load: resource type is unknown")
	}
}

// setDir is used if a test case or application wishes to override
// one or more of the default resource directory locations.
func (l *loader) setDir(dir string, data interface{}) {
	switch data.(type) {
	case *Glyphs:
		l.dir[gly] = dir
	case *Material:
		l.dir[mat] = dir
	case *Mesh:
		l.dir[msh] = dir
	case *Shader:
		l.dir[shd] = dir
	case *Sound:
		l.dir[snd] = dir
	case *Texture:
		l.dir[tex] = dir
	default:
		log.Printf("loader.setDir: resource type is unknown")
	}
}

// glyphset loads glyph set resources from disk.  Nil is returned if there were
// loading problems.  Currently only the ".fnt" file format is supported.
func (l loader) glyphset(directory, name string) *Glyphs {
	if glyphs, err := l.fnt(directory, name+".fnt"); err == nil {
		glyphs.Name = name
		return glyphs
	}
	return nil
}

// material loads a material from the indicated file.  Material returns nil
// if the material cannot be loaded.  Currently only the Wavefront .mtl file
// format is supported.
func (l loader) material(directory, name string) *Material {
	mat, err := l.mtl(directory, name+".mtl")
	if err != nil {
		log.Printf("Could not load material %s %s\n", name, err)
		return nil
	}
	mat.Name = name
	return mat
}

// mesh loads one or more meshes from disk. Mesh can return an empty slice if
// nothing was found or if there were other loading errors.  Currently only
// the Wavefront .obj file format is supported.
func (l loader) mesh(directory, name string) *Mesh {
	meshes, err := l.obj(directory, name+".obj")
	if err != nil || len(meshes) <= 0 {
		log.Printf("Could not read meshes %s %s\n", name, err)
		return nil
	}
	return meshes[0] // ignore multiple meshes for now.
}

// shader loads the source for a shader program.  The shader still needs to be
// compiled, linked, and have its uniforms specified.  The shader name is
// expected to be the prefix for both vertex (.vsh) and fragment (.fsh) shader
// source files.
func (l loader) shader(directory, name string) *Shader {
	shader := &Shader{Name: name}
	shader.Vsh = l.src(directory, name+".vsh")
	shader.Fsh = l.src(directory, name+".fsh")
	return shader
}

// sound loads a sound from the indicated file.  Sound can return nil
// if the sound cannot be loaded.  Currently only the wave .wav file
// format is supported.
func (l loader) sound(directory, name string) *Sound {
	sound := &Sound{Name: name}
	if err := l.wav(sound, directory, name+".wav"); err != nil {
		log.Printf("Could not load sound %s %s\n", name, err)
		return nil
	}
	return sound
}

// texture loads a texture from disk and initializes it so that it is ready to
// use. Texture return nil if there were loading problems.  The name is
// expected to be the prefix for a .png file.
func (l loader) texture(directory, name string) *Texture {
	img, err := l.png(directory, name+".png")
	if err == nil {
		texture := &Texture{Name: name, Img: img}
		return texture
	} else {
		println("error ", err)
	}
	return nil
}

// getResource locates the named resource.  This is expected to be used either
// in production where the resources have been included with the application,
// or development where the resources are on disk in the local directory.
//
// The caller is responsible for closing the returned file.
func (l loader) getResource(directory, name string) (file io.ReadCloser, err error) {
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
