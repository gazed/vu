// Copyright © 2024 Galvanized Logic Inc.

package physics

import (
	"log/slog"
	"math"
	"slices"

	"github.com/gazed/vu/math/lin"
)

// collider_Contact;
type collider_Contact struct {
	collision_point1 lin.V3
	collision_point2 lin.V3
	normal           lin.V3
}

// collider_Convex_Hull_Face;
type collider_Convex_Hull_Face struct {
	elements []uint32
	normal   lin.V3
}

// collider_Convex_Hull;
type collider_Convex_Hull struct {
	vertices             []lin.V3
	transformed_vertices []lin.V3
	faces                []collider_Convex_Hull_Face
	transformed_faces    []collider_Convex_Hull_Face
	vertex_to_faces      [][]uint32
	vertex_to_neighbors  [][]uint32
	face_to_neighbors    [][]uint32
}

// collider_Sphere;
type collider_Sphere struct {
	radius float32 // typedef float r32;
	center lin.V3  // default 0,0,0
}

// collider_Type;
type collider_Type uint8

const (
	collider_TYPE_SPHERE collider_Type = iota
	collider_TYPE_CONVEX_HULL
)

// collider;
type collider struct {
	ctype       collider_Type
	convex_hull collider_Convex_Hull // hull...
	sphere      collider_Sphere      // ...or sphere based on type
}

// @NOTE: for simplicity (and speed), we don't deal with scaling in the colliders.
// therefore, if the object is scaled, the collider needs to be recreated
// (and the vertices should be already scaled when creating it)
func collider_sphere_create(radius float32) collider {
	var collider collider
	collider.ctype = collider_TYPE_SPHERE
	collider.sphere.radius = radius
	collider.sphere.center = lin.V3{}
	return collider
}

// collider_sphere_destroy
func collider_sphere_destroy(collider *collider) {}

// get_sphere_collider_bounding_sphere_radius
func get_sphere_collider_bounding_sphere_radius(collider *collider) float64 {
	return float64(collider.sphere.radius)
}

// do_triangles_share_same_vertex
func do_triangles_share_same_vertex(t1, t2 v3Int) bool {
	return t1.x == t2.x || t1.x == t2.y || t1.x == t2.z ||
		t1.y == t2.x || t1.y == t2.y || t1.y == t2.z ||
		t1.z == t2.x || t1.z == t2.y || t1.z == t2.z
}

// do_faces_share_same_vertex
func do_faces_share_same_vertex(e1, e2 []uint32) bool {
	for i := 0; i < len(e1); i++ {
		i1 := e1[i]
		for j := 0; j < len(e2); j++ {
			i2 := e2[j]
			if i1 == i2 {
				return true
			}
		}
	}
	return false
}

// collect_faces_planar_to
func collect_faces_planar_to(hull []lin.V3, hull_triangle_faces []v3Int, triangle_faces_to_neighbor_faces_map [][]uint32,
	is_triangle_face_already_processed_arr []bool, face_to_test_idx uint32, target_normal lin.V3, out []v3Int) []v3Int {

	EPSILON := 0.000001
	face_to_test := hull_triangle_faces[face_to_test_idx]
	v1 := hull[face_to_test.x]
	v2 := hull[face_to_test.y]
	v3 := hull[face_to_test.z]

	v12 := lin.NewV3().Sub(&v2, &v1)
	v13 := lin.NewV3().Sub(&v3, &v1)
	face_normal := lin.NewV3().Cross(v12, v13).Unit()

	if is_triangle_face_already_processed_arr[face_to_test_idx] {
		return out
	}
	projection := face_normal.Dot(&target_normal)
	if (projection-1.0) > -EPSILON && (projection-1.0) < EPSILON {
		out = append(out, face_to_test)
		is_triangle_face_already_processed_arr[face_to_test_idx] = true

		neighbor_faces := triangle_faces_to_neighbor_faces_map[face_to_test_idx]
		for i := 0; i < len(neighbor_faces); i++ {
			neighbor_face_idx := neighbor_faces[i]
			out = collect_faces_planar_to(hull, hull_triangle_faces, triangle_faces_to_neighbor_faces_map,
				is_triangle_face_already_processed_arr, neighbor_face_idx, target_normal, out)
		}
	}
	return out
}

// get_edge_index
func get_edge_index(edges []v2Int, edge v2Int) int32 {
	for i := 0; i < len(edges); i++ {
		current_edge := edges[i]
		if current_edge.x == edge.x && current_edge.y == edge.y {
			return int32(i)
		}
		if current_edge.x == edge.y && current_edge.y == edge.x {
			return int32(i)
		}
	}
	return -1
}

// create_convex_hull_face
func create_convex_hull_face(triangles []v3Int, face_normal lin.V3) collider_Convex_Hull_Face {

	// Collect the edges that form the border of the face
	edges := []v2Int{}
	for i := 0; i < len(triangles); i++ {
		triangle := triangles[i]
		edge1 := v2Int{triangle.x, triangle.y}
		edge2 := v2Int{triangle.y, triangle.z}
		edge3 := v2Int{triangle.z, triangle.x}
		edge1_idx := int(get_edge_index(edges, edge1))
		if edge1_idx >= 0 {
			last := len(edges) - 1
			edges[edge1_idx] = edges[last]             // replace deleted with last.
			edges = slices.Delete(edges, last, last+1) // delete last
		} else {
			edges = append(edges, edge1)
		}
		edge2_idx := int(get_edge_index(edges, edge2))
		if edge2_idx >= 0 {
			last := len(edges) - 1
			edges[edge2_idx] = edges[last]             // replace deleted with last.
			edges = slices.Delete(edges, last, last+1) // delete last
		} else {
			edges = append(edges, edge2)
		}
		edge3_idx := int(get_edge_index(edges, edge3))
		if edge3_idx >= 0 {
			last := len(edges) - 1
			edges[edge3_idx] = edges[last]             // replace deleted with last.
			edges = slices.Delete(edges, last, last+1) // delete last
		} else {
			edges = append(edges, edge3)
		}
	}

	// Nicely order the edges
	for i := 0; i < len(edges); i++ {
		current_edge := edges[i]
		for j := i + 1; j < len(edges); j++ {
			candidate_edge := edges[j]
			if current_edge.y != candidate_edge.x && current_edge.y != candidate_edge.y {
				continue
			}
			if current_edge.y == candidate_edge.y {
				tmp := candidate_edge.x
				candidate_edge.x = candidate_edge.y
				candidate_edge.y = tmp
			}

			tmp := edges[i+1]
			edges[i+1] = candidate_edge
			edges[j] = tmp
		}
	}
	// assert(edges[0].x == edges[array_length(edges) - 1].y);

	// Simply create the face elements based on the edges
	face_elements := []uint32{}
	for i := 0; i < len(edges); i++ {
		current_edge := edges[i]
		face_elements = append(face_elements, uint32(current_edge.x))
	}

	face := collider_Convex_Hull_Face{}
	face.elements = face_elements
	face.normal = face_normal
	return face
}

// is_neighbor_already_in_vertex_to_neighbors_map
func is_neighbor_already_in_vertex_to_neighbors_map(vertex_to_neighbors []uint32, neighbor uint32) bool {
	for i := 0; i < len(vertex_to_neighbors); i++ {
		if vertex_to_neighbors[i] == neighbor {
			return true
		}
	}
	return false
}

// get_convex_hull_collider_bounding_sphere_radius
func get_convex_hull_collider_bounding_sphere_radius(collider *collider) float64 {
	max_distance := 0.0
	for i := 0; i < len(collider.convex_hull.vertices); i++ {
		v := collider.convex_hull.vertices[i]
		distance := v.Len()
		if distance > max_distance {
			max_distance = distance
		}
	}
	return max_distance
}

// Create a convex hull from the vertices+indices
// For now, we assume that the mesh is already a convex hull
// This function only makes sure that vertices are unique - duplicated vertices will be merged.
func collider_convex_hull_create(vertices []lin.V3, indices []uint32) collider {
	vertex_to_idx_map := map[lin.V3]uint32{}

	// Build hull, eliminating duplicated vertex
	hull := []lin.V3{}
	for i := 0; i < len(vertices); i++ {
		current_vertex := vertices[i]
		current_index := uint32(0)
		if _, ok := vertex_to_idx_map[current_vertex]; !ok {
			current_index = uint32(len(hull))
			hull = append(hull, current_vertex)
			vertex_to_idx_map[current_vertex] = current_index
		}
	}

	// Collect all triangle faces that compose the just-built hull
	hull_triangle_faces := []v3Int{}
	for i := 0; i < len(indices); i += 3 {
		i1 := indices[i]
		i2 := indices[i+1]
		i3 := indices[i+2]
		v1 := vertices[i1]
		v2 := vertices[i2]
		v3 := vertices[i3]

		new_i1, ok1 := vertex_to_idx_map[v1]
		new_i2, ok2 := vertex_to_idx_map[v2]
		new_i3, ok3 := vertex_to_idx_map[v3]
		if !ok1 || !ok2 || !ok3 {
			slog.Error("collider_convex_hull_create: dev error: check your code")
		}

		triangle := v3Int{new_i1, new_i2, new_i3}
		hull_triangle_faces = append(hull_triangle_faces, triangle)
	}

	// Prepare vertex to faces map
	vertex_to_faces_map := make([][]uint32, len(hull))
	for i := 0; i < len(hull); i++ {
		vertex_to_faces_map[i] = []uint32{}
	}

	// Prepare vertex to neighbors map
	vertex_to_neighbors_map := make([][]uint32, len(hull))
	for i := 0; i < len(hull); i++ {
		vertex_to_neighbors_map[i] = []uint32{}
	}

	// Prepare triangle faces to neighbors map
	triangle_faces_to_neighbor_faces_map := make([][]uint32, len(hull_triangle_faces))
	for i := 0; i < len(hull_triangle_faces); i++ {
		triangle_faces_to_neighbor_faces_map[i] = []uint32{}
	}

	// Create the vertex to neighbors map
	for i := 0; i < len(hull_triangle_faces); i++ {
		triangle_face := hull_triangle_faces[i]
		for j := 0; j < len(hull_triangle_faces); j++ {
			if i == j {
				continue
			}
			face_to_test := hull_triangle_faces[j]
			if do_triangles_share_same_vertex(triangle_face, face_to_test) {
				triangle_faces_to_neighbor_faces_map[i] = append(triangle_faces_to_neighbor_faces_map[i], uint32(j))
			}
		}

		// Fill vertex to edges map
		if !is_neighbor_already_in_vertex_to_neighbors_map(vertex_to_neighbors_map[triangle_face.x], triangle_face.y) {
			vertex_to_neighbors_map[triangle_face.x] = append(vertex_to_neighbors_map[triangle_face.x], triangle_face.y)
		}
		if !is_neighbor_already_in_vertex_to_neighbors_map(vertex_to_neighbors_map[triangle_face.x], triangle_face.z) {
			vertex_to_neighbors_map[triangle_face.x] = append(vertex_to_neighbors_map[triangle_face.x], triangle_face.z)
		}
		if !is_neighbor_already_in_vertex_to_neighbors_map(vertex_to_neighbors_map[triangle_face.y], triangle_face.x) {
			vertex_to_neighbors_map[triangle_face.y] = append(vertex_to_neighbors_map[triangle_face.y], triangle_face.x)
		}
		if !is_neighbor_already_in_vertex_to_neighbors_map(vertex_to_neighbors_map[triangle_face.y], triangle_face.z) {
			vertex_to_neighbors_map[triangle_face.y] = append(vertex_to_neighbors_map[triangle_face.y], triangle_face.z)
		}
		if !is_neighbor_already_in_vertex_to_neighbors_map(vertex_to_neighbors_map[triangle_face.z], triangle_face.x) {
			vertex_to_neighbors_map[triangle_face.z] = append(vertex_to_neighbors_map[triangle_face.z], triangle_face.x)
		}
		if !is_neighbor_already_in_vertex_to_neighbors_map(vertex_to_neighbors_map[triangle_face.z], triangle_face.y) {
			vertex_to_neighbors_map[triangle_face.z] = append(vertex_to_neighbors_map[triangle_face.z], triangle_face.y)
		}
	}

	// Collect all 'de facto' faces of the convex hull
	faces := []collider_Convex_Hull_Face{}
	is_triangle_face_already_processed_arr := make([]bool, len(hull_triangle_faces))
	for i := 0; i < len(hull_triangle_faces); i++ {
		if is_triangle_face_already_processed_arr[i] {
			continue
		}

		triangle_face := hull_triangle_faces[i]
		v1 := hull[triangle_face.x]
		v2 := hull[triangle_face.y]
		v3 := hull[triangle_face.z]
		v12 := lin.NewV3().Sub(&v2, &v1)
		v13 := lin.NewV3().Sub(&v3, &v1)
		normal := lin.NewV3().Cross(v12, v13).Unit()
		planar_faces := []v3Int{}
		planar_faces = collect_faces_planar_to(hull, hull_triangle_faces, triangle_faces_to_neighbor_faces_map,
			is_triangle_face_already_processed_arr, uint32(i), *normal, planar_faces)

		new_face := create_convex_hull_face(planar_faces, *normal)
		new_face_index := uint32(len(faces))
		faces = append(faces, new_face)

		// Fill vertex to faces map accordingly
		for j := 0; j < len(planar_faces); j++ {
			planar_face := planar_faces[j]
			vertex_to_faces_map[planar_face.x] = append(vertex_to_faces_map[planar_face.x], new_face_index)
			vertex_to_faces_map[planar_face.y] = append(vertex_to_faces_map[planar_face.y], new_face_index)
			vertex_to_faces_map[planar_face.z] = append(vertex_to_faces_map[planar_face.z], new_face_index)
		}
	}

	// Prepare face to neighbors map
	face_to_neighbor_faces_map := make([][]uint32, len(faces))
	for i := 0; i < len(faces); i++ {
		face_to_neighbor_faces_map[i] = []uint32{}
	}

	// Fill faces to neighbor faces map
	for i := 0; i < len(faces); i++ {
		face := faces[i] // collider_Convex_Hull_Face
		for j := 0; j < len(faces); j++ {
			if i == j {
				continue
			}
			candidate_face := faces[j] // collider_Convex_Hull_Face
			if do_faces_share_same_vertex(face.elements, candidate_face.elements) {
				face_to_neighbor_faces_map[i] = append(face_to_neighbor_faces_map[i], uint32(j))
			}
		}
	}

	// create the convex hull collider
	collider := collider{}
	collider.ctype = collider_TYPE_CONVEX_HULL
	convex_hull := &collider.convex_hull
	convex_hull.faces = faces
	convex_hull.transformed_faces = append(convex_hull.transformed_faces, faces...)
	convex_hull.vertices = hull
	convex_hull.transformed_vertices = append(convex_hull.transformed_vertices, hull...)
	convex_hull.vertex_to_faces = vertex_to_faces_map
	convex_hull.vertex_to_neighbors = vertex_to_neighbors_map
	convex_hull.face_to_neighbors = face_to_neighbor_faces_map
	return collider
}

// collider_convex_hull_destroy - unnecessary
func collider_convex_hull_destroy(collider *collider) {}

// collider_destroy - unnecessary
func collider_destroy(collider *collider) {}

// colliders_destroy - unnecessary
func colliders_destroy(colliders []collider) {}

// collider_update
func collider_update(collider *collider, translation lin.V3, rotation *lin.Q) {
	switch collider.ctype {
	case collider_TYPE_CONVEX_HULL:
		model_matrix_no_scale := util_get_model_matrix_no_scale(rotation, translation)
		for i := 0; i < len(collider.convex_hull.transformed_vertices); i++ {
			vertex := lin.NewV4().SetS(
				collider.convex_hull.vertices[i].X,
				collider.convex_hull.vertices[i].Y,
				collider.convex_hull.vertices[i].Z,
				1.0,
			)
			transformed_vertex := lin.NewV4().MultMv(&model_matrix_no_scale, vertex)
			transformed_vertex.Scale(transformed_vertex, 1.0/transformed_vertex.W)
			collider.convex_hull.transformed_vertices[i].SetS(transformed_vertex.X, transformed_vertex.Y, transformed_vertex.Z)
		}
		for i := 0; i < len(collider.convex_hull.transformed_faces); i++ {
			normal := collider.convex_hull.faces[i].normal
			mm := lin.NewM3().SetM4(&model_matrix_no_scale)
			transformed_normal := lin.NewV3().MultMv(mm, &normal)
			collider.convex_hull.transformed_faces[i].normal = *(transformed_normal.Unit())
		}
	case collider_TYPE_SPHERE:
		collider.sphere.center = translation
	default:
		slog.Error("collider_update: unsupported collider type", "collider_type", collider.ctype)
	}
}

// colliders_update
func colliders_update(colliders []collider, translation lin.V3, rotation *lin.Q) {
	for i := 0; i < len(colliders); i++ {
		collider_update(&colliders[i], translation, rotation)
	}
}

// colliders_get_default_inertia_tensor
func colliders_get_default_inertia_tensor(colliders []collider, mass float64) lin.M3 {
	// For now, the center of mass is always assumed to be at 0,0,0
	if len(colliders) == 1 {
		collider := &colliders[0]
		if collider.ctype == collider_TYPE_SPHERE {
			// for now we assume the sphere is centered at its center of mass (because then the inertia tensor is simple)
			// assert(gm_vec3_is_zero(collider->sphere.center));
			I := (2.0 / 5.0) * mass * float64(collider.sphere.radius*collider.sphere.radius)
			result := lin.NewM3()
			result.Xx = I
			result.Yy = I
			result.Zz = I
			return *result
		}
	}

	total_num_vertices := 0
	for i := 0; i < len(colliders); i++ {
		collider := &colliders[i]
		total_num_vertices += len(collider.convex_hull.vertices)
	}
	mass_per_vertex := mass / float64(total_num_vertices)

	result := lin.NewM3()
	for i := 0; i < len(colliders); i++ {
		collider := &colliders[i]
		if collider.ctype != collider_TYPE_CONVEX_HULL {
			slog.Error("colliders_get_default_inertia_tensor expects collider_TYPE_CONVEX_HULL")
		}
		for j := 0; j < len(collider.convex_hull.vertices); j++ {
			v := collider.convex_hull.vertices[j]
			result.Xx += mass_per_vertex * (v.Y*v.Y + v.Z*v.Z)
			result.Xy += mass_per_vertex * v.X * v.Y
			result.Xz += mass_per_vertex * v.X * v.Z
			result.Yx += mass_per_vertex * v.X * v.Y
			result.Yy += mass_per_vertex * (v.X*v.X + v.Z*v.Z)
			result.Yz += mass_per_vertex * v.Y * v.Z
			result.Zx += mass_per_vertex * v.X * v.Z
			result.Zy += mass_per_vertex * v.Y * v.Z
			result.Zz += mass_per_vertex * (v.X*v.X + v.Y*v.Y)
		}
	}
	return *result
}

// collider_get_bounding_sphere_radius
func collider_get_bounding_sphere_radius(collider *collider) float64 {
	switch collider.ctype {
	case collider_TYPE_CONVEX_HULL:
		return get_convex_hull_collider_bounding_sphere_radius(collider)
	case collider_TYPE_SPHERE:
		return get_sphere_collider_bounding_sphere_radius(collider)
	}
	slog.Error("collider_get_bounding_sphere_radius: no bounding radius")
	return 0.0
}

// colliders_get_bounding_sphere_radius
func colliders_get_bounding_sphere_radius(colliders []collider) float64 {
	max_bounding_sphere_radius := -math.MaxFloat64
	for i := 0; i < len(colliders); i++ {
		collider := &colliders[i]
		bounding_sphere_radius := collider_get_bounding_sphere_radius(collider)
		if bounding_sphere_radius > max_bounding_sphere_radius {
			max_bounding_sphere_radius = bounding_sphere_radius
		}
	}
	return max_bounding_sphere_radius
}

// collider_get_contacts
func collider_get_contacts(collider1, collider2 *collider, contacts []collider_Contact) []collider_Contact {
	var simplex gjk_Simplex
	normal := lin.NewV3()

	// If both colliders are spheres, calling EPA is not only extremely slow, but also provide bad results.
	// GJK is also not necessary. In this case, just calculate everything analytically.
	if collider1.ctype == collider_TYPE_SPHERE && collider2.ctype == collider_TYPE_SPHERE {
		distance_vector := lin.NewV3().Sub(&collider1.sphere.center, &collider2.sphere.center)
		distance_sqd := distance_vector.Dot(distance_vector)
		min_distance := collider1.sphere.radius + collider2.sphere.radius
		if distance_sqd < float64(min_distance*min_distance) {
			// Spheres are colliding
			normal.Sub(&collider2.sphere.center, &collider1.sphere.center).Unit()
			penetration := float64(min_distance) - math.Sqrt(distance_sqd)
			contacts = clipping_get_contact_manifold(collider1, collider2, *normal, penetration, contacts)
		}
		return contacts
	}

	// Call GJK to check if there is a collision
	if gjk_collides(collider1, collider2, &simplex) {
		// There is a collision.  Get the collision normal using EPA
		normal, penetration, ok := epa(collider1, collider2, &simplex)
		if !ok {
			return contacts
		}

		// Finally, clip the results to get the result manifold
		contacts = clipping_get_contact_manifold(collider1, collider2, normal, penetration, contacts)
	}
	return contacts
}

// colliders_get_contacts
func colliders_get_contacts(colliders1, colliders2 []collider) []collider_Contact {
	contacts := []collider_Contact{}
	for i := 0; i < len(colliders1); i++ {
		collider1 := &colliders1[i]
		for j := 0; j < len(colliders2); j++ {
			collider2 := &colliders2[j]
			contacts = collider_get_contacts(collider1, collider2, contacts)
		}
	}
	return contacts
}
