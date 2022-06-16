// Copyright © 2024 Galvanized Logic Inc.

package physics

import (
	"log/slog"
	"math"

	"github.com/gazed/vu/math/lin"
)

// support_point_get_index
func support_point_get_index(convex_hull *collider_Convex_Hull, direction lin.V3) uint32 {
	var selected_index uint32
	max_dot := -math.MaxFloat64
	for i := 0; i < len(convex_hull.transformed_vertices); i++ {
		dot := convex_hull.transformed_vertices[i].Dot(&direction)
		if dot > max_dot {
			selected_index = uint32(i)
			max_dot = dot
		}
	}
	return selected_index
}

// support_point
func support_point(collider *collider, direction lin.V3) lin.V3 {
	v3 := lin.NewV3()
	switch collider.ctype {
	case collider_TYPE_CONVEX_HULL:
		selected_index := support_point_get_index(&collider.convex_hull, direction)
		return collider.convex_hull.transformed_vertices[selected_index]
	case collider_TYPE_SPHERE:
		v3.Add(&collider.sphere.center, lin.NewV3().Scale(lin.NewV3().Set(&direction).Unit(), float64(collider.sphere.radius)))
		return *v3
	}
	slog.Error("unsupported collider type", "collider_type", collider.ctype)
	return *v3
}

// support_point_of_minkowski_difference
func support_point_of_minkowski_difference(collider1, collider2 *collider, direction lin.V3) lin.V3 {
	support1 := support_point(collider1, direction)
	support2 := support_point(collider2, *(lin.NewV3().Scale(&direction, -1)))
	return *(lin.NewV3().Sub(&support1, &support2))
}
