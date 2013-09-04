// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package data

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Load a Wavefront .mtl file which is a text representation of one
// or more material descriptions.  See the file format specification at:
//    https://en.wikipedia.org/wiki/Wavefront_.obj_file#File_format
//    http://web.archive.org/web/20080813073052/
//    http://local.wasp.uwa.edu.au/~pbourke/dataformats/mtl/
func (l loader) mtl(directory, filename string) (mat *Material, err error) {
	var file io.ReadCloser
	if file, err = l.getResource(directory, filename); err != nil {
		return mat, fmt.Errorf("could not open %s %s", filename, err)
	}
	defer file.Close()
	mat = &Material{}
	var f1, f2, f3 float32
	reader := bufio.NewReader(file)
	line, e1 := reader.ReadString('\n')
	for ; e1 == nil; line, e1 = reader.ReadString('\n') {
		tokens := strings.Split(line, " ")
		switch tokens[0] {
		case "Ka": // ambient
			if _, e := fmt.Sscanf(line, "Ka %f %f %f", &f1, &f2, &f3); e != nil {
				return mat, fmt.Errorf("could not parse ambient values %s", e)
			}
			mat.Ka.R, mat.Ka.G, mat.Ka.B = f1, f2, f3
		case "Kd": // diffuse
			if _, e := fmt.Sscanf(line, "Kd %f %f %f", &f1, &f2, &f3); e != nil {
				return mat, fmt.Errorf("could not parse diffuse values %s", e)
			}
			mat.Kd.R, mat.Kd.G, mat.Kd.B = f1, f2, f3
		case "Ks": // specular
			if _, e := fmt.Sscanf(line, "Ks %f %f %f", &f1, &f2, &f3); e != nil {
				return mat, fmt.Errorf("could not parse specular values %s", e)
			}
			mat.Ks.R, mat.Ks.G, mat.Ks.B = f1, f2, f3
		case "d": // transparency
			a, _ := strconv.ParseFloat(strings.TrimSpace(tokens[1]), 32)
			mat.Tr = float32(a)
		case "newmtl": // material name
		case "Ns": // specular exponent - scaler. Ignored for now.
		case "Ni": // optical density - scaler. Ignored for now.
		case "illum": // illumination model - int. Ignored for now.
		}
	}
	return
}
