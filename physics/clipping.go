// SPDX-FileCopyrightText : © 2024 Galvanized Logic Inc.
// SPDX-License-Identifier: BSD-2-Clause

package physics

import (
	"log/slog"
	"math"

	"github.com/gazed/vu/math/lin"
)

// cPlane
type cPlane struct {
	normal lin.V3
	point  lin.V3
}

// is_point_in_plane
func is_point_in_plane(plane *cPlane, position lin.V3) bool {
	distance := -plane.normal.Dot(&plane.point)
	if position.Dot(&plane.normal)+distance < 0.0 {
		return false
	}
	return true
}

// plane_edge_intersection
func plane_edge_intersection(plane *cPlane, start lin.V3, end lin.V3, out_point *lin.V3) bool {
	const EPSILON float64 = 0.000001
	ab := lin.NewV3().Sub(&end, &start)

	// Check that the edge and plane are not parallel and thus never intersect
	// We do this by projecting the line (start - A, End - B) ab along the plane
	ab_p := plane.normal.Dot(ab)
	if math.Abs(ab_p) > EPSILON {
		// Generate a random point on the plane (any point on the plane will suffice)
		distance := -plane.normal.Dot(&plane.point)
		p_co := lin.NewV3().Scale(&plane.normal, -distance)

		// Work out the edge factor to scale edge by
		// e.g. how far along the edge to traverse before it meets the plane.
		// This is computed by: -proj<plane_nrml>(edge_start - any_planar_point) / proj<plane_nrml>(edge_start - edge_end)
		fac := -plane.normal.Dot(lin.NewV3().Sub(&start, p_co)) / ab_p

		// Stop any large floating point divide issues with almost parallel planes
		fac = math.Min(math.Max(fac, 0.0), 1.0)

		// Return point on edge
		out_point.Add(&start, ab.Scale(ab, fac))
		return true
	}
	return false
}

// sutherland_hodgman clips the input polygon to the input clip planes
// If remove_instead_of_clipping is true, vertices that are lying outside the clipping planes will be removed instead of clipped
// Based on https://research.ncl.ac.uk/game/mastersdegree/gametechnologies/previousinformation/physics5collisionmanifolds/
func sutherland_hodgman(input_polygon []lin.V3, clip_planes []cPlane, remove_instead_of_clipping bool) (out_polygon []lin.V3) {
	if len(clip_planes) <= 0 {
		slog.Error("sutherland_hodgman called with no clip planes")
		return out_polygon
	}

	// Create temporary list of vertices
	// We will keep ping-pong'ing between the two lists updating them as we go.
	input := append([]lin.V3{}, input_polygon...)
	output := []lin.V3{}

	for i := 0; i < len(clip_planes); i++ {
		// If every single point has already been removed previously, just exit
		if len(input) == 0 {
			break
		}
		plane := &clip_planes[i]

		// Loop through each edge of the polygon and clip that edge against the current plane.
		temp_point, start_point := lin.NewV3(), input[len(input)-1]
		for j := 0; j < len(input); j++ {
			end_point := input[j]
			start_in_plane := is_point_in_plane(plane, start_point)
			end_in_plane := is_point_in_plane(plane, end_point)

			if remove_instead_of_clipping {
				if end_in_plane {
					output = append(output, end_point)
				}
			} else {
				// If the edge is entirely within the clipping plane, keep it as it is
				if start_in_plane && end_in_plane {
					output = append(output, end_point)
				} else if start_in_plane && !end_in_plane {
					// If the edge interesects the clipping plane, cut the edge along clip plane
					if plane_edge_intersection(plane, start_point, end_point, temp_point) {
						output = append(output, *temp_point)
					}
				} else if !start_in_plane && end_in_plane {
					if plane_edge_intersection(plane, start_point, end_point, temp_point) {
						output = append(output, *temp_point)
					}
					output = append(output, end_point)
				}
			}
			// ..otherwise the edge is entirely outside the clipping plane and should be removed/ignored
			start_point = end_point
		}
		// Swap input/output polygons, and clear output list for us to generate afresh
		tmp := input
		input = output
		output = tmp[:0] // clear array keeping allocated memory
	}
	return input
}

// get_closest_point_polygon
func get_closest_point_polygon(position lin.V3, reference_plane *cPlane) lin.V3 {
	d := lin.NewV3().Scale(&reference_plane.normal, -1.0).Dot(&reference_plane.point)
	t := lin.NewV3().Sub(&position, lin.NewV3().Scale(&reference_plane.normal, reference_plane.normal.Dot(&position)+d))
	return *t
}

// build_boundary_planes
func build_boundary_planes(convex_hull *collider_Convex_Hull, target_face_idx uint32) []cPlane {
	result := []cPlane{}
	face_neighbors := convex_hull.face_to_neighbors[target_face_idx]
	for i := 0; i < len(face_neighbors); i++ {
		neighbor_face := convex_hull.transformed_faces[face_neighbors[i]]
		p := cPlane{}
		p.point = convex_hull.transformed_vertices[neighbor_face.elements[0]]
		p.normal.Neg(&neighbor_face.normal)
		result = append(result, p)
	}
	return result
}

// get_face_with_most_fitting_normal
func get_face_with_most_fitting_normal(support_idx uint32, convex_hull *collider_Convex_Hull, normal lin.V3) uint32 {
	const EPSILON float64 = 0.000001
	support_faces := convex_hull.vertex_to_faces[support_idx]

	max_proj := -math.MaxFloat64
	var selected_face_idx uint32
	for i := 0; i < len(support_faces); i++ {
		face := convex_hull.transformed_faces[support_faces[i]]
		proj := face.normal.Dot(&normal)
		if proj > max_proj {
			max_proj = proj
			selected_face_idx = support_faces[i]
		}
	}
	return selected_face_idx
}

// get_edge_with_most_fitting_normal
func get_edge_with_most_fitting_normal(
	support1_idx, support2_idx uint32,
	convex_hull1, convex_hull2 *collider_Convex_Hull,
	normal lin.V3, edge_normal *lin.V3) v4Int {

	// inverted_normal := lin.NewV3().Neg(&normal) ... not used
	support1 := &convex_hull1.transformed_vertices[support1_idx]
	support2 := &convex_hull2.transformed_vertices[support2_idx]
	support1_neighbors := convex_hull1.vertex_to_neighbors[support1_idx]
	support2_neighbors := convex_hull2.vertex_to_neighbors[support2_idx]

	max_dot := -math.MaxFloat64
	selected_edges := v4Int{}

	for i := 0; i < len(support1_neighbors); i++ {
		neighbor1 := convex_hull1.transformed_vertices[support1_neighbors[i]]
		edge1 := lin.NewV3().Sub(support1, &neighbor1)
		for j := 0; j < len(support2_neighbors); j++ {
			neighbor2 := convex_hull2.transformed_vertices[support2_neighbors[j]]
			edge2 := lin.NewV3().Sub(support2, &neighbor2)

			current_normal := lin.NewV3().Cross(edge1, edge2).Unit()
			current_normal_inverted := lin.NewV3().Neg(current_normal)

			dot := current_normal.Dot(&normal)
			if dot > max_dot {
				max_dot = dot
				selected_edges.x = support1_idx
				selected_edges.y = support1_neighbors[i]
				selected_edges.z = support2_idx
				selected_edges.w = support2_neighbors[j]
				*edge_normal = *current_normal
			}
			dot = current_normal_inverted.Dot(&normal)
			if dot > max_dot {
				max_dot = dot
				selected_edges.x = support1_idx
				selected_edges.y = support1_neighbors[i]
				selected_edges.z = support2_idx
				selected_edges.w = support2_neighbors[j]
				*edge_normal = *current_normal_inverted
			}
		}
	}
	return selected_edges
}

// collision_distance_between_skew_lines calculates the distance between
// two indepedent skew lines in the 3D world.
// The first line is given by a known point P1 and a direction vector D1
// The second line is given by a known point P2 and a direction vector D2
// Outputs:
//
//	L1 is the closest POINT to the second line that belongs to the first line
//	L2 is the closest POINT to the first line that belongs to the second line
//	_N is the number that satisfies L1 = P1 + _N * D1
//	_M is the number that satisfies L2 = P2 + _M * D2
func collision_distance_between_skew_lines(p1, d1, p2, d2 lin.V3, l1, l2 *lin.V3, _n, _m *float64) bool {
	n1 := d1.X*d2.X + d1.Y*d2.Y + d1.Z*d2.Z
	n2 := d2.X*d2.X + d2.Y*d2.Y + d2.Z*d2.Z
	m1 := -d1.X*d1.X - d1.Y*d1.Y - d1.Z*d1.Z
	m2 := -d2.X*d1.X - d2.Y*d1.Y - d2.Z*d1.Z
	r1 := -d1.X*p2.X + d1.X*p1.X - d1.Y*p2.Y + d1.Y*p1.Y - d1.Z*p2.Z + d1.Z*p1.Z
	r2 := -d2.X*p2.X + d2.X*p1.X - d2.Y*p2.Y + d2.Y*p1.Y - d2.Z*p2.Z + d2.Z*p1.Z

	// Solve 2x2 linear system
	if (n1*m2)-(n2*m1) == 0 {
		return false
	}
	n := ((r1 * m2) - (r2 * m1)) / ((n1 * m2) - (n2 * m1))
	m := ((n1 * r2) - (n2 * r1)) / ((n1 * m2) - (n2 * m1))
	if l1 != nil {
		l1.Add(&p1, l1.Scale(&d1, m))
	}
	if l2 != nil {
		l2.Add(&p2, l2.Scale(&d2, n))
	}
	if _n != nil {
		*_n = n
	}
	if _m != nil {
		*_m = m
	}
	return true
}

// get_vertices_of_faces
func get_vertices_of_faces(hull *collider_Convex_Hull, face collider_Convex_Hull_Face) []lin.V3 {
	vertices := []lin.V3{}
	for i := 0; i < len(face.elements); i++ {
		vertices = append(vertices, hull.transformed_vertices[face.elements[i]])
	}
	return vertices
}

// convex_convex_contact_manifold
func convex_convex_contact_manifold(collider1, collider2 *collider,
	normal lin.V3, contacts []collider_Contact) []collider_Contact {
	if collider1.ctype != collider_TYPE_CONVEX_HULL || collider2.ctype != collider_TYPE_CONVEX_HULL {
		slog.Error("convex_convex_contact_manifold expects two COLLIDER_TYPE_CONVEX_HULL")
		return []collider_Contact{}
	}
	convex_hull1 := &collider1.convex_hull
	convex_hull2 := &collider2.convex_hull

	const EPSILON float64 = 0.0001

	inverted_normal := lin.NewV3().Neg(&normal)

	// vec3 edge_normal;
	support1_idx := support_point_get_index(convex_hull1, normal)
	support2_idx := support_point_get_index(convex_hull2, *inverted_normal)
	face1_idx := get_face_with_most_fitting_normal(support1_idx, convex_hull1, normal)
	face2_idx := get_face_with_most_fitting_normal(support2_idx, convex_hull2, *inverted_normal)
	face1 := convex_hull1.transformed_faces[face1_idx]
	face2 := convex_hull2.transformed_faces[face2_idx]
	edge_normal := lin.NewV3()
	edges := get_edge_with_most_fitting_normal(support1_idx, support2_idx, convex_hull1, convex_hull2, normal, edge_normal)

	chosen_normal1_dot := face1.normal.Dot(&normal)
	chosen_normal2_dot := face2.normal.Dot(inverted_normal)
	edge_normal_dot := edge_normal.Dot(&normal)

	if edge_normal_dot > chosen_normal1_dot+EPSILON && edge_normal_dot > chosen_normal2_dot+EPSILON {
		// Edge
		l1, l2 := lin.NewV3(), lin.NewV3()
		p1 := convex_hull1.transformed_vertices[edges.x]
		d1 := lin.NewV3().Sub(&convex_hull1.transformed_vertices[edges.y], &p1)
		p2 := convex_hull2.transformed_vertices[edges.z]
		d2 := lin.NewV3().Sub(&convex_hull2.transformed_vertices[edges.w], &p2)
		collision_distance_between_skew_lines(p1, *d1, p2, *d2, l1, l2, nil, nil)
		contact := collider_Contact{*l1, *l2, normal}
		contacts = append(contacts, contact)
	} else {
		// Face
		var reference_face_support_points []lin.V3
		var incident_face_support_points []lin.V3
		var boundary_planes []cPlane

		is_face1_the_reference_face := chosen_normal1_dot > chosen_normal2_dot
		if is_face1_the_reference_face {
			reference_face_support_points = get_vertices_of_faces(convex_hull1, face1)
			incident_face_support_points = get_vertices_of_faces(convex_hull2, face2)
			boundary_planes = build_boundary_planes(convex_hull1, face1_idx)
		} else {
			reference_face_support_points = get_vertices_of_faces(convex_hull2, face2)
			incident_face_support_points = get_vertices_of_faces(convex_hull1, face1)
			boundary_planes = build_boundary_planes(convex_hull2, face2_idx)
		}

		clipped_points := sutherland_hodgman(incident_face_support_points, boundary_planes, false)

		var reference_plane cPlane
		if is_face1_the_reference_face {
			reference_plane.normal.Neg(&face1.normal)
		} else {
			reference_plane.normal.Neg(&face2.normal)
		}
		reference_plane.point = reference_face_support_points[0]

		final_clipped_points := []lin.V3{}
		final_clipped_points = sutherland_hodgman(clipped_points, []cPlane{reference_plane}, true)

		for i := 0; i < len(final_clipped_points); i++ {
			point := final_clipped_points[i]
			closest_point := get_closest_point_polygon(point, &reference_plane)
			point_diff := lin.NewV3().Sub(&point, &closest_point)
			var contact_penetration float64

			// we are projecting the points that are in the incident face on the reference planes
			// so the points that we have are part of the incident object.
			var contact collider_Contact
			if is_face1_the_reference_face {
				contact_penetration = point_diff.Dot(&normal)
				contact.collision_point1.Sub(&point, lin.NewV3().Scale(&normal, contact_penetration))
				contact.collision_point2 = point
			} else {
				contact_penetration = -point_diff.Dot(&normal)
				contact.collision_point1 = point
				contact.collision_point2.Add(&point, lin.NewV3().Scale(&normal, contact_penetration))
			}
			contact.normal = normal
			if contact_penetration < 0.0 {
				contacts = append(contacts, contact)
			}
		}
	}
	if len(contacts) == 0 {
		slog.Debug("convex_convex_contact_manifold: no intersection was found")
	}
	return contacts
}

// clipping_get_contact_manifold
func clipping_get_contact_manifold(collider1, collider2 *collider,
	normal lin.V3, penetration float64, contacts []collider_Contact) []collider_Contact {
	// TODO: For now, we only consider CONVEX and SPHERE colliders.
	// If new colliders are added, we can think about making this more generic.

	switch {
	case collider1.ctype == collider_TYPE_SPHERE:
		sphere_collision_point := support_point(collider1, normal)

		var contact collider_Contact
		contact.collision_point1 = sphere_collision_point
		contact.collision_point2.Sub(&sphere_collision_point, lin.NewV3().Scale(&normal, penetration))
		contact.normal = normal
		contacts = append(contacts, contact)
	case collider2.ctype == collider_TYPE_SPHERE:
		inverse_normal := lin.NewV3().Neg(&normal)
		sphere_collision_point := support_point(collider2, *inverse_normal)

		var contact collider_Contact
		contact.collision_point1.Add(&sphere_collision_point, lin.NewV3().Scale(&normal, penetration))
		contact.collision_point2 = sphere_collision_point
		contact.normal = normal
		contacts = append(contacts, contact)
	case collider1.ctype == collider_TYPE_CONVEX_HULL && collider2.ctype == collider_TYPE_CONVEX_HULL:
		contacts = convex_convex_contact_manifold(collider1, collider2, normal, contacts)
	default:
		slog.Error("unsupported collider types", "c1", collider1.ctype, "c2", collider2.ctype)
	}
	return contacts
}
