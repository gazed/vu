// SPDX-FileCopyrightText : © 2024 Galvanized Logic Inc.
// SPDX-License-Identifier: BSD-2-Clause

package physics

import (
	"fmt"
	"log/slog"
	"math"
	"slices"

	"github.com/gazed/vu/math/lin"
)

// polytope_from_gjk_simplex
func polytope_from_gjk_simplex(s *gjk_Simplex) (polytope []lin.V3, faces []v3Int) {
	if s.num != 4 {
		slog.Error("polytope_from_gjk_simplex: expecting simplex.num:4")
	}
	polytope = []lin.V3{s.a, s.b, s.c, s.d}
	faces = []v3Int{
		v3Int{0, 1, 2}, // ABC
		v3Int{0, 2, 3}, // ACD
		v3Int{0, 3, 1}, // ADB
		v3Int{1, 2, 3}, // BCD
	}
	return polytope, faces
}

// get_face_normal_and_distance_to_origin
func get_face_normal_and_distance_to_origin(face v3Int, polytope []lin.V3) (normal lin.V3, distance float64) {
	a := &polytope[face.x]
	b := &polytope[face.y]
	c := &polytope[face.z]

	ab := lin.NewV3().Sub(b, a)
	ac := lin.NewV3().Sub(c, a)
	n := lin.NewV3().Cross(ab, ac).Unit()
	if n.X == 0.0 && n.Y == 0.0 && n.Z == 0.0 {
		slog.Error("get_face_normal_and_distance_to_origin: normal is zero vector")
		return normal, distance
	}

	// When this value is not 0, it is possible that the normals
	// are not found even if the polytope is not degenerate
	const DISTANCE_TO_ORIGIN_TOLERANCE = 0.0000000000000

	// the distance from the face's *plane* to the origin (considering an infinite plane).
	distance = n.Dot(a)
	if distance < -DISTANCE_TO_ORIGIN_TOLERANCE {
		// if the distance is less than 0, it means that our normal is point inwards instead of outwards
		// in this case, we just invert both normal and distance
		// this way, we don't need to worry about face's winding
		n.Neg(n)
		distance = -distance
	} else if distance >= -DISTANCE_TO_ORIGIN_TOLERANCE && distance <= DISTANCE_TO_ORIGIN_TOLERANCE {
		// if the distance is exactly 0.0, then it means that the origin is lying exactly on the face.
		// in this case, we can't directly infer the orientation of the normal.
		// since our shape is convex, we analyze the other vertices of the hull to deduce the orientation
		was_able_to_calculate_normal := false
		for i := 0; i < len(polytope); i++ {
			current := polytope[i]
			auxiliar_distance := n.Dot(&current)
			if auxiliar_distance < -DISTANCE_TO_ORIGIN_TOLERANCE || auxiliar_distance > DISTANCE_TO_ORIGIN_TOLERANCE {
				// since the shape is convex, the other vertices should always be "behind" the normal plane
				if auxiliar_distance >= -DISTANCE_TO_ORIGIN_TOLERANCE {
					n.Neg(n)
				}
				was_able_to_calculate_normal = true
				break
			}
		}

		// If we were not able to calculate the normal, it means that ALL points of the polytope are in the same plane
		// Therefore, we either have a degenerate polytope or our tolerance is not big enough
		if !was_able_to_calculate_normal {
			panic(fmt.Errorf("all points on same plane or degenerate polytope"))
		}
	}
	return *n, distance
}

// add_edge
func add_edge(edges []v2Int, edge v2Int, polytope []lin.V3) []v2Int {
	// @TODO: we can use a hash table here
	for i := 0; i < len(edges); i++ {
		current := edges[i]
		if edge.x == current.x && edge.y == current.y {
			return slices.Delete(edges, i, i+1)
		}
		if edge.x == current.y && edge.y == current.x {
			return slices.Delete(edges, i, i+1)
		}

		// @TEMPORARY: Once indexes point to unique vertices, this won't be needed.
		current_v1 := polytope[current.x]
		current_v2 := polytope[current.y]
		edge_v1 := polytope[edge.x]
		edge_v2 := polytope[edge.y]
		if current_v1.Eq(&edge_v1) && current_v2.Eq(&edge_v2) {
			return slices.Delete(edges, i, i+1)
		}
		if current_v1.Eq(&edge_v2) && current_v2.Eq(&edge_v1) {
			return slices.Delete(edges, i, i+1)
		}
	}
	return append(edges, edge)
}

//	static vec3 triangle_centroid(vec3 p1, vec3 p2, vec3 p3) {
//		vec3 centroid = gm_vec3_add(gm_vec3_add(p2, p3), p1);
//		centroid = gm_vec3_scalar_product(1.0 / 3.0, centroid);
//		return centroid;
//	}
func triangle_centroid(p1, p2, p3 lin.V3) (centroid lin.V3) {
	centroid.Add(&p2, &p3).Add(&centroid, &p1)
	centroid.Scale(&centroid, 1.0/3.0)
	return centroid
}

// epa
func epa(collider1, collider2 *collider, simplex *gjk_Simplex) (normal lin.V3, penetration float64, success bool) {
	const epsilon float64 = 0.0001

	// build initial polytope from GJK simplex
	polytope, faces := polytope_from_gjk_simplex(simplex)

	normals := []lin.V3{}
	faces_distance_to_origin := []float64{}
	min_normal := lin.NewV3()
	min_distance := math.MaxFloat64
	for i := 0; i < len(faces); i++ {
		var normal lin.V3
		var distance float64
		face := faces[i]
		normal, distance = get_face_normal_and_distance_to_origin(face, polytope)
		normals = append(normals, normal)
		faces_distance_to_origin = append(faces_distance_to_origin, distance)
		if distance < min_distance {
			min_distance = distance
			*min_normal = normal
		}
	}

	edges := []v2Int{}
	converged := false
	for it := 0; it < 100; it++ {
		support_point := support_point_of_minkowski_difference(collider1, collider2, *min_normal)

		// If the support time lies on the face currently set as the closest to the origin, we are done.
		d := min_normal.Dot(&support_point)
		if math.Abs(d-min_distance) < epsilon {
			normal = *min_normal
			penetration = min_distance
			converged = true
			// slog.Debug("epa converged", "iterations", it)
			break
		}

		// add new point to polytope
		new_point_index := uint32(len(polytope))
		polytope = append(polytope, support_point)

		// Expand Polytope
		loop_counter := 0
		for i := 0; i < len(normals); i++ {
			normal := normals[i]
			face := faces[i]

			// If the face normal points towards the support point, we need to reconstruct it.
			centroid := triangle_centroid(polytope[face.x], polytope[face.y], polytope[face.z])

			// If the face normal points towards the support point, we need to reconstruct it.
			if normal.Dot(lin.NewV3().Sub(&support_point, &centroid)) > 0.0 {
				face := faces[i]

				edge1 := v2Int{face.x, face.y}
				edge2 := v2Int{face.y, face.z}
				edge3 := v2Int{face.z, face.x}

				edges = add_edge(edges, edge1, polytope)
				edges = add_edge(edges, edge2, polytope)
				edges = add_edge(edges, edge3, polytope)

				// Relative order between the two arrays should be kept.
				faces = slices.Delete(faces, i, i+1)
				faces_distance_to_origin = slices.Delete(faces_distance_to_origin, i, i+1)
				normals = slices.Delete(normals, i, i+1)

				i -= 1

				loop_counter += 1
				if loop_counter > 1000 {
					panic(fmt.Errorf("epa: infinite loop"))
				}
			}
		}

		for i := 0; i < len(edges); i++ {
			edge := edges[i]
			new_face := v3Int{x: edge.x, y: edge.y, z: new_point_index}
			faces = append(faces, new_face)

			new_face_normal, new_face_distance := get_face_normal_and_distance_to_origin(new_face, polytope)
			normals = append(normals, new_face_normal)
			faces_distance_to_origin = append(faces_distance_to_origin, new_face_distance)
		}

		min_distance = math.MaxFloat64
		for i := 0; i < len(faces_distance_to_origin); i++ {
			distance := faces_distance_to_origin[i]
			if distance < min_distance {
				min_distance = distance
				min_normal = &normals[i]
			}
		}
		edges = edges[:0]
	}
	if !converged {
		slog.Warn("EPA did not converge.")
	}
	return normal, penetration, converged
}
