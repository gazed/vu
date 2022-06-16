// Copyright © 2024 Galvanized Logic Inc.

package physics

import (
	"log/slog"

	"github.com/gazed/vu/math/lin"
)

// broad_Collision_Pair
type broad_Collision_Pair struct {
	b1_id bid
	b2_id bid
}

// broad_get_collision_pairs
func broad_get_collision_pairs(bodies []Body) []broad_Collision_Pair {
	collision_pairs := []broad_Collision_Pair{}
	for i := 0; i < len(bodies); i++ {
		b1 := &bodies[i]
		for j := i + 1; j < len(bodies); j++ {
			b2 := &bodies[j]
			entities_distance := lin.NewV3().Sub(&b1.world_position, &b2.world_position).Len()

			// Increase the distance a little to account for moving objects.
			// @TODO: We should derivate this value from delta_time, forces, velocities, etc
			max_distance_for_collision := b1.bounding_sphere_radius + b2.bounding_sphere_radius + 0.1
			if entities_distance <= max_distance_for_collision {
				// body ID is the bodies array index.
				collision_pairs = append(collision_pairs, broad_Collision_Pair{bid(i), bid(j)})
			}
		}
	}
	return collision_pairs
}

// uf_find
func uf_find(body_to_parent_map map[bid]bid, x bid) bid {
	p, ok := body_to_parent_map[x]
	if !ok {
		slog.Error("missing body parent", "body_id", x)
	}
	if p == x {
		return x
	}
	return uf_find(body_to_parent_map, p)
}

// uf_union
func uf_union(body_to_parent_map map[bid]bid, x, y bid) {
	key := uf_find(body_to_parent_map, y)
	value := uf_find(body_to_parent_map, x)
	body_to_parent_map[key] = value
}

// uf_collect_all
func uf_collect_all(bodies []Body, collision_pairs []broad_Collision_Pair) map[bid]bid {
	body_to_parent_map := map[bid]bid{}
	for i := range bodies {
		// id := bodies[i].id
		body_to_parent_map[bid(i)] = bid(i)
	}
	for i := range collision_pairs {
		collision_pair := &collision_pairs[i]
		id1 := collision_pair.b1_id
		id2 := collision_pair.b2_id
		b1 := body_get_by_id(id1)
		b2 := body_get_by_id(id2)
		if !b1.fixed && !b2.fixed {
			uf_union(body_to_parent_map, id1, id2)
		}
	}
	return body_to_parent_map
}

// broad_collect_simulation_islands
func broad_collect_simulation_islands(bodies []Body, collision_pairs []broad_Collision_Pair, constraints []constraint) [][]bid {
	simulation_islands := [][]bid{}

	// Collect the simulation islands into an entity->parent map
	body_to_parent_map := uf_collect_all(bodies, collision_pairs)

	// Extra step: To avoid bugs, we need to make sure that entities that
	// are part of a same constraint are also part of the same island!
	for i := range constraints {
		c := &constraints[i]
		id1 := c.b1_id
		id2 := c.b2_id
		b1 := body_get_by_id(id1)
		b2 := body_get_by_id(id2)
		if !b1.fixed && !b2.fixed {
			uf_union(body_to_parent_map, id1, id2)
		}
	}

	// As a last step, transform the simulation islands into a nice structure
	simulation_islands_map := map[bid]uint32{}
	for i := range bodies {
		b := &bodies[i]
		if b.fixed {
			continue
		}
		parent := uf_find(body_to_parent_map, bid(i))
		simulation_island_idx, ok := simulation_islands_map[parent]
		if !ok {
			// Simulation Island not created yet.
			simulation_island_idx = uint32(len(simulation_islands))
			simulation_islands = append(simulation_islands, []bid{})
			simulation_islands_map[parent] = simulation_island_idx
		}
		simulation_islands[simulation_island_idx] = append(simulation_islands[simulation_island_idx], bid(i))
	}
	return simulation_islands
}

func broad_simulation_islands_destroy(simulation_islands [][]bid) {
	// unnecessary.
}
