// SPDX-FileCopyrightText : © 2024 Galvanized Logic Inc.
// SPDX-License-Identifier: BSD-2-Clause

package physics

// body.go ported from entity.ccp/entity.h.
// Used physics "body" to distinguish from vu/Entity

import (
	"log/slog"

	"github.com/gazed/vu/math/lin"
)

// Body is a one object within a physics simulation.
// Bodies are added for objects that need to participate in physics.
// Bodies transforms are set by the physics simulation and not the application.
type Body struct {
	world_position lin.V3 // set/get using accessor methods
	world_rotation lin.Q  // set/get using accessor methods
	world_scale    lin.V3 // set/get using accessor methods

	// Physics Related
	colliders                    []collider
	bounding_sphere_radius       float64
	forces                       []force
	inverse_mass                 float64
	inertia_tensor               lin.M3
	inverse_inertia_tensor       lin.M3
	angular_velocity             lin.V3
	linear_velocity              lin.V3
	fixed                        bool
	active                       bool
	deactivation_time            float64
	static_friction_coefficient  float64
	dynamic_friction_coefficient float64
	restitution_coefficient      float64

	// PBD Auxilar
	previous_world_position   lin.V3
	previous_world_rotation   lin.Q
	previous_linear_velocity  lin.V3
	previous_angular_velocity lin.V3
}

// force applies a force to a physics Body
type force struct {
	position     lin.V3 // local or world position.
	newtons      lin.V3 // amount of force.
	local_coords bool   // true if using local coordinates.
}

// body_create_ex
func body_create_ex(
	world_position lin.V3,
	world_rotation lin.Q,
	world_scale lin.V3,
	mass float64,
	colliders []collider,
	static_friction_coefficient float64,
	dynamic_friction_coefficient float64,
	restitution_coefficient float64,
	is_fixed bool) *Body {

	if static_friction_coefficient < 0.0 || static_friction_coefficient > 1.0 {
		slog.Error("invalid static_friction_coefficient", "value", static_friction_coefficient)
		return nil
	}
	if dynamic_friction_coefficient < 0.0 || dynamic_friction_coefficient > 1.0 {
		slog.Error("invalid dynamic_friction_coefficient", "value", dynamic_friction_coefficient)
		return nil
	}
	if restitution_coefficient < 0.0 || restitution_coefficient > 1.0 {
		slog.Error("invalid restitution_coefficient", "value", restitution_coefficient)
		return nil
	}

	body := &Body{}
	body.world_position = world_position
	body.world_rotation = world_rotation
	body.world_scale = world_scale

	// initial velocities all default to zero, V3{0,0,0}
	// body.angular_velocity
	// body.linear_velocity
	// body.previous_angular_velocity
	// body.previous_linear_velocity

	body.bounding_sphere_radius = colliders_get_bounding_sphere_radius(colliders)
	if is_fixed {
		body.inverse_mass = 0.0
		body.inertia_tensor = *lin.M3Z // this is not correct, but it shouldn't make a difference
		body.inverse_inertia_tensor = *lin.M3Z
	} else {
		body.inverse_mass = 1.0 / mass
		body.inertia_tensor = colliders_get_default_inertia_tensor(colliders, mass)
		body.inverse_inertia_tensor.Inv(&body.inertia_tensor)
	}
	body.forces = []force{}
	body.fixed = is_fixed
	body.active = true
	body.deactivation_time = 0.0
	body.colliders = colliders
	body.static_friction_coefficient = static_friction_coefficient
	body.dynamic_friction_coefficient = dynamic_friction_coefficient
	body.restitution_coefficient = restitution_coefficient
	if body.dynamic_friction_coefficient > body.static_friction_coefficient {
		slog.Warn("dynamic friction coefficient is greater than static friction coefficient")
	}
	return body
}

// body_get_by_id return the id for this body.
// bodies set each simulation step: see physics.go.
func body_get_by_id(id bid) *Body {
	return &bodies[id]
}

// SetPosition sets the bodies world position.
func (body *Body) SetPosition(world_position lin.V3) {
	body.world_position = world_position
}

// Position returns the bodies world position.
func (body *Body) Position() (world_position *lin.V3) {
	return &body.world_position
}
func (body *Body) Velocity() (velocity *lin.V3) {
	return &body.linear_velocity
}

// SetRotation sets the bodies rotation.
func (body *Body) SetRotation(world_rotation lin.Q) {
	body.world_rotation = world_rotation
}

// Rotation gets the bodies world rotation.
func (body *Body) Rotation() (world_rotation *lin.Q) {
	return &body.world_rotation
}

// SetScale sets the bodies scale
func (body *Body) SetScale(world_scale lin.V3) {
	body.world_scale = world_scale
}

// Activate set the body as active in the simulation.
func (body *Body) Activate() {
	body.active = true
	body.deactivation_time = 0.0
}

// Push adds to the bodies linear velocity.
func (body *Body) Push(x, y, z float64) {
	body.linear_velocity.X += x
	body.linear_velocity.Y += y
	body.linear_velocity.Z += z
}

// AddForce to this body.
// If local_coords is false, then the position and force are represented
// in world coordinates, assuming that the center of the
// world is the center of the entity. That is, the coordinate (0, 0, 0)
// corresponds to the center of the entity in world coords.
// If local_coords is true, then the position and force are
// represented in local coords.
func (body *Body) AddForce(position lin.V3, newtons lin.V3, local_coords bool) {
	if local_coords {
		// If the force and position are in local cords, we first convert them to world coords
		// (actually, we convert them to ~"world coords centered at entity"~)
		newtons.MultQ(&newtons, &body.world_rotation)

		// note that we don't need translation since we want to be centered at entity anyway
		model_matrix := body.get_model_matrix_no_translation()

		//  gm_mat4_multiply_vec3(&model_matrix, position, true); <-- as point
		mmv3 := lin.NewM3().SetM4(&model_matrix)
		position.MultMv(mmv3, &position)
		position.X += model_matrix.Xw // 03: translation X
		position.Y += model_matrix.Yw // 13: translation W
		position.Z += model_matrix.Zw // 23: translation Z
	}
	f := force{}
	f.newtons = newtons
	f.position = position
	body.forces = append(body.forces, f)
}

// ClearForces
//
//	void entity_clear_forces(Entity* entity) {
//	    array_clear(entity->forces);
//	}
func (body *Body) clear_forces() {
	body.forces = body.forces[:0] // clear keeping memory.
}

// get_model_matrix
func (body *Body) get_model_matrix() lin.M4 {
	scale_matrix := &lin.M4{
		body.world_scale.X, 0.0, 0.0, 0.0,
		0.0, body.world_scale.Y, 0.0, 0.0,
		0.0, 0.0, body.world_scale.Z, 0.0,
		0.0, 0.0, 0.0, 1.0,
	}
	rotation_matrix := lin.NewM4().SetQ(&body.world_rotation)
	translation_matrix := &lin.M4{
		1.0, 0.0, 0.0, body.world_position.X,
		0.0, 1.0, 0.0, body.world_position.Y,
		0.0, 0.0, 1.0, body.world_position.Z,
		0.0, 0.0, 0.0, 1.0,
	}
	model_matrix := rotation_matrix.Mult(rotation_matrix, scale_matrix)
	model_matrix.Mult(translation_matrix, model_matrix)
	return *model_matrix
}

// get_model_matrix_no_translation
func (body *Body) get_model_matrix_no_translation() lin.M4 {
	scale_matrix := &lin.M4{
		body.world_scale.X, 0.0, 0.0, 0.0,
		0.0, body.world_scale.Y, 0.0, 0.0,
		0.0, 0.0, body.world_scale.Z, 0.0,
		0.0, 0.0, 0.0, 1.0,
	}
	rotation_matrix := lin.NewM4().SetQ(&body.world_rotation)
	model_matrix := rotation_matrix.Mult(rotation_matrix, scale_matrix)
	return *model_matrix
}

// util_get_model_matrix_no_scale
func util_get_model_matrix_no_scale(rotation *lin.Q, translation lin.V3) lin.M4 {
	rotation_matrix := lin.NewM4().SetQ(rotation)
	translation_matrix := &lin.M4{
		Xx: 1.0, Xy: 0.0, Xz: 0.0, Xw: translation.X, // Row-major order.
		Yx: 0.0, Yy: 1.0, Yz: 0.0, Yw: translation.Y,
		Zx: 0.0, Zy: 0.0, Zz: 1.0, Zw: translation.Z,
		Wx: 0.0, Wy: 0.0, Wz: 0.0, Ww: 1.0,
	}
	model_matrix := lin.NewM4().Mult(translation_matrix, rotation_matrix)
	return *model_matrix
}
