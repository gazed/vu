// Copyright © 2013-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package load

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Mtl loads a Wavefront MTL file which is a text representation of one
// or more material descriptions. See the MTL file format specification at:
//    https://en.wikipedia.org/wiki/Wavefront_.obj_file#File_format
//    http://web.archive.org/web/20080813073052/
//    http://paulbourke.net/dataformats/mtl/
// The Reader r is expected to be opened and closed by the caller.
// A successful import overwrites the data in MtlData.
func Mtl(r io.Reader, d *MtlData) error {
	var f1, f2, f3 float32
	reader := bufio.NewReader(r)
	line, e1 := reader.ReadString('\n')
	for ; e1 == nil; line, e1 = reader.ReadString('\n') {
		tokens := strings.Split(line, " ")
		switch tokens[0] {
		case "Ka": // ambient
			if _, e := fmt.Sscanf(line, "Ka %f %f %f", &f1, &f2, &f3); e != nil {
				return fmt.Errorf("could not parse ambient values %s", e)
			}
			d.KaR, d.KaG, d.KaB = f1, f2, f3
		case "Kd": // diffuse
			if _, e := fmt.Sscanf(line, "Kd %f %f %f", &f1, &f2, &f3); e != nil {
				return fmt.Errorf("could not parse diffuse values %s", e)
			}
			d.KdR, d.KdG, d.KdB = f1, f2, f3
		case "Ks": // specular
			if _, e := fmt.Sscanf(line, "Ks %f %f %f", &f1, &f2, &f3); e != nil {
				return fmt.Errorf("could not parse specular values %s", e)
			}
			d.KsR, d.KsG, d.KsB = f1, f2, f3
		case "d": // transparency
			a, _ := strconv.ParseFloat(strings.TrimSpace(tokens[1]), 32)
			d.Alpha = float32(a)
		case "Ns": // specular exponent
			ns, _ := strconv.ParseFloat(strings.TrimSpace(tokens[1]), 32)
			d.Ns = float32(ns)
		case "newmtl": // material name
		case "Ni": // optical density - scaler. Ignored for now.
		case "illum": // illumination model - int. Ignored for now.
		}
	}
	return nil
}
