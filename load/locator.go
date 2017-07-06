// Copyright © 2013-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package load

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Locator knows how to search disk based locations for files.
// Locator uses a built in knowledge of paths and file types.
// It uses a convention for locating file types in directories where
// the defaults can be overridden or added to using the Dir method.
type Locator interface {
	Dir(ext, dir string) Locator // Map a file extension to a directory.
	Dispose()                    // Properly terminate asset loading

	// GetResource allows applications to include and find custom resources.
	//   name: specific resource identifier, like a file or full file path.
	//   dir : prepended to the name path like a directory.
	GetResource(name string) (file io.ReadCloser, err error)
}

// NewLocator returns the default asset locator. The default Locator
// looks directly to disk for development builds and for a zip file for
// production builds. The default asset locator expects all locations
// are directories relative to the application location.
// The default Locator maps the following file types to the given directories.
//    PNG               : "images"
//    WAV               : "audio"
//    OBJ, IQM, MTL     : "models"
//    FNT, VSH, FSH, TXT: "source"
func NewLocator() Locator { return newLocator() }

// ===========================================================================
// locator implements Locator.

// locator knows where to find asset data on disk.
type locator struct {
	reader *zip.ReadCloser   // Used as the resource file if set.
	dirs   map[string]string //
}

// newLocator returns the default Locator implementation and asset
// directory locations. These are conventions for locating zipped assets
// in different situations.
func newLocator() *locator {
	var resources *zip.ReadCloser // packaged resources.
	programName := os.Args[0]     // qualified path to executable
	assetZip := path.Join(path.Dir(programName), "../Resources/assets.zip")
	if reader, err := zip.OpenReader(assetZip); err == nil {
		resources = reader // OSX packaged application.
	} else if reader, err := zip.OpenReader(programName); err == nil {
		resources = reader // windows non-store exe. Zip with Exe.
	} else {
		// windows store app.
		// use absolute path to executable since relative files
		// are not located when running as a properly installed appx.
		// Windows locates apps in
		//     /c/Program Files/WindowsApps/mangledAppName
		// data app data in
		//     ~/AppData/Local/Packages/mangledAppName/LocalCache/...
		programName = filepath.Dir(os.Args[0]) // executable directory
		absDir, err0 := filepath.Abs(programName)
		assetZip = path.Join(absDir, "Assets/assets.zip")
		if reader, err := zip.OpenReader(assetZip); err0 == nil && err == nil {
			resources = reader // Windows
		}
	}

	// if resources is still nil then this is likely a debug build
	// and GetResources below will attempt to read directly from disk.
	l := &locator{reader: resources}
	l.dirs = map[string]string{ // default directories for file locations.
		"OBJ":  "models",
		"IQM":  "models",
		"MTL":  "models",
		"WAV":  "audio",
		"TXT":  "source",
		"VSH":  "source",
		"FSH":  "source",
		"FNT":  "source",
		"JSON": "source",
		"PNG":  "images",
	}
	return l
}

// GetResource locates the named resource. This is expected to be used either
// in production where the resources have been included with the application,
// or development where the resources are on disk in the local directory.
//
// The caller is responsible for closing the returned file.
func (l *locator) GetResource(name string) (file io.ReadCloser, err error) {
	prefix, ext := "", ""
	if sep := strings.LastIndexAny(name, "."); sep != -1 {
		ext = strings.ToUpper(name[sep+1:])
	}
	if val, defined := l.dirs[ext]; defined { // optional group lookup.
		prefix = val
	}
	filePath := strings.TrimSpace(path.Join(prefix, name))
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

// Dir maps a file extention to a directory. Having a convention
// means that only the file name needs to be specified.
// The extension, by another convention, is the last 3 letters of
// the GetResource name parameter.
func (l *locator) Dir(ext, dir string) Locator {
	l.dirs[ext] = dir // add or overwrite where to find a file type.
	return l
}

// Dispose properly terminates the loader.
// This is only needed when the loader has been reading resources from a file.
func (l *locator) Dispose() {
	if l.reader != nil {
		l.reader.Close()
	}
}
