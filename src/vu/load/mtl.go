// Copyright © 2013-2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package load

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// MtlData holds colour and alpha information. It is an intermediate format
// that is intended for populating render assets.
type MtlData struct {
	KaR, KaG, KaB float32 // Ambient colour.
	KdR, KdG, KdB float32 // Diffuse colour.
	KsR, KsG, KsB float32 // Specular colour.
	Tr            float32 // Transparency
}

// Load a Wavefront .mtl file which is a text representation of one
// or more material descriptions.  See the file format specification at:
//    https://en.wikipedia.org/wiki/Wavefront_.obj_file#File_format
//    http://web.archive.org/web/20080813073052/
//    http://local.wasp.uwa.edu.au/~pbourke/dataformats/mtl/
func (l *loader) mtl(name string) (data *MtlData, err error) {
	mtl := &MtlData{}
	var file io.ReadCloser
	if file, err = l.getResource(l.dir[mod], name+".mtl"); err != nil {
		return mtl, fmt.Errorf("could not open %s %s", name+".mtl", err)
	}
	defer file.Close()
	var f1, f2, f3 float32
	reader := bufio.NewReader(file)
	line, e1 := reader.ReadString('\n')
	for ; e1 == nil; line, e1 = reader.ReadString('\n') {
		tokens := strings.Split(line, " ")
		switch tokens[0] {
		case "Ka": // ambient
			if _, e := fmt.Sscanf(line, "Ka %f %f %f", &f1, &f2, &f3); e != nil {
				return mtl, fmt.Errorf("could not parse ambient values %s", e)
			}
			mtl.KaR, mtl.KaG, mtl.KaB = f1, f2, f3
		case "Kd": // diffuse
			if _, e := fmt.Sscanf(line, "Kd %f %f %f", &f1, &f2, &f3); e != nil {
				return mtl, fmt.Errorf("could not parse diffuse values %s", e)
			}
			mtl.KdR, mtl.KdG, mtl.KdB = f1, f2, f3
		case "Ks": // specular
			if _, e := fmt.Sscanf(line, "Ks %f %f %f", &f1, &f2, &f3); e != nil {
				return mtl, fmt.Errorf("could not parse specular values %s", e)
			}
			mtl.KsR, mtl.KsG, mtl.KsB = f1, f2, f3
		case "d": // transparency
			a, _ := strconv.ParseFloat(strings.TrimSpace(tokens[1]), 32)
			mtl.Tr = float32(a)
		case "newmtl": // material name
		case "Ns": // specular exponent - scaler. Ignored for now.
		case "Ni": // optical density - scaler. Ignored for now.
		case "illum": // illumination model - int. Ignored for now.
		}
	}
	return mtl, nil
}
