// SPDX-FileCopyrightText : © 2024 Galvanized Logic Inc.
// SPDX-License-Identifier: BSD-2-Clause

package physics

import (
	"github.com/gazed/vu/math/lin"
)

// calculate_external_force
// Calculate the sum of all external forces acting on an entity
func calculate_external_force(b *Body) lin.V3 {
	// center_of_mass := lin.NewV3() ... not used
	total_force := lin.NewV3()
	for i := 0; i < len(b.forces); i++ {
		total_force.Add(total_force, &b.forces[i].newtons)
	}
	return *total_force
}

// calculate_external_torque
// Calculate the sum of all external torques acting on an entity
func calculate_external_torque(b *Body) lin.V3 {
	center_of_mass := lin.NewV3()
	total_torque := lin.NewV3()
	distance := lin.NewV3()
	for i := 0; i < len(b.forces); i++ {
		distance.Sub(&b.forces[i].position, center_of_mass)
		total_torque.Add(total_torque, distance.Cross(distance, &b.forces[i].newtons))
	}
	return *total_torque
}

// get_dynamic_inertia_tensor
// Calculate the dynamic inertia tensor of an entity,
// i.e., the inertia tensor transformed considering entity's rotation
func get_dynamic_inertia_tensor(b *Body) lin.M3 {
	// Can only be used if the local->world matrix is orthogonal
	rotation_matrix := lin.NewM3().SetQ(&b.world_rotation)
	transposed_rotation_matrix := lin.NewM3().Transpose(rotation_matrix)
	aux := lin.NewM3().Mult(rotation_matrix, &b.inertia_tensor)
	aux.Mult(aux, transposed_rotation_matrix)
	return *aux
}

// get_dynamic_inverse_inertia_tensor
// Calculate the dynamic inverse inertia tensor of an entity,
// i.e., the inverse inertia tensor transformed considering entity's rotation
func get_dynamic_inverse_inertia_tensor(b *Body) lin.M3 {
	// Can only be used if the local->world matrix is orthogonal
	rotation_matrix := lin.NewM3().SetQ(&b.world_rotation)
	transposed_rotation_matrix := lin.NewM3().Transpose(rotation_matrix)
	aux := lin.NewM3().Mult(rotation_matrix, &b.inverse_inertia_tensor)
	aux.Mult(aux, transposed_rotation_matrix)
	return *aux
}
