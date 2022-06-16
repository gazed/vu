// Copyright © 2024 Galvanized Logic Inc.

package physics

import (
	"log/slog"
	"math"

	"github.com/gazed/vu/math/lin"
)

// pbd_Axis_Type;
type pbd_Axis_Type uint8

const (
	pbd_POSITIVE_X_AXIS pbd_Axis_Type = iota
	pbd_NEGATIVE_X_AXIS
	pbd_POSITIVE_Y_AXIS
	pbd_NEGATIVE_Y_AXIS
	pbd_POSITIVE_Z_AXIS
	pbd_NEGATIVE_Z_AXIS
)

// constraint_Type;
type constraint_Type uint8

const (
	positional_CONSTRAINT constraint_Type = iota
	collision_CONSTRAINT
	mutual_ORIENTATION_CONSTRAINT
	hinge_JOINT_CONSTRAINT
	spherical_JOINT_CONSTRAINT
)

// positional_Constraint;
type positional_Constraint struct {
	r1_lc      lin.V3
	r2_lc      lin.V3
	compliance float64
	lambda     float64
	distance   lin.V3
}

// collision_Constraint;
type collision_Constraint struct {
	r1_lc    lin.V3
	r2_lc    lin.V3
	normal   lin.V3
	lambda_t float64
	lambda_n float64
}

// mutual_Orientation_Constraint;
type mutual_Orientation_Constraint struct {
	compliance float64
	lambda     float64
}

// hinge_Joint_Constraint;
type hinge_Joint_Constraint struct {
	r1_lc      lin.V3
	r2_lc      lin.V3
	compliance float64
	lambda_pos float64

	b1_aligned_axis     pbd_Axis_Type
	b2_aligned_axis     pbd_Axis_Type
	lambda_aligned_axes float64

	limited           bool
	upper_limit       float64
	lower_limit       float64
	b1_limit_axis     pbd_Axis_Type
	b2_limit_axis     pbd_Axis_Type
	lambda_limit_axes float64
}

// spherical_Joint_Constraint;
type spherical_Joint_Constraint struct {
	r1_lc      lin.V3
	r2_lc      lin.V3
	lambda_pos float64

	lambda_swing      float64
	swing_upper_limit float64
	swing_lower_limit float64
	b1_swing_axis     pbd_Axis_Type
	b2_swing_axis     pbd_Axis_Type

	lambda_twist      float64
	twist_upper_limit float64
	twist_lower_limit float64
	b1_twist_axis     pbd_Axis_Type
	b2_twist_axis     pbd_Axis_Type
}

// constraint;
type constraint struct {
	ctype constraint_Type
	b1_id bid
	b2_id bid

	// one of the following based on ctype
	positional_constraint         positional_Constraint
	collision_constraint          collision_Constraint
	mutual_orientation_constraint mutual_Orientation_Constraint
	hinge_joint_constraint        hinge_Joint_Constraint
	spherical_joint_constraint    spherical_Joint_Constraint
}

// #define LINEAR_SLEEPING_THRESHOLD 0.10
// #define ANGULAR_SLEEPING_THRESHOLD 0.10
// #define DEACTIVATION_TIME_TO_BE_INACTIVE 1.0
const (
	linear_SLEEPING_THRESHOLD        = 0.10
	angular_SLEEPING_THRESHOLD       = 0.10
	deactivation_TIME_TO_BE_INACTIVE = 1.0
)

// pbd_positional_constraint_init
func pbd_positional_constraint_init(constraint *constraint, b1_id, b2_id bid, r1_lc, r2_lc lin.V3, compliance float64, distance lin.V3) {
	constraint.ctype = positional_CONSTRAINT
	constraint.b1_id = b1_id
	constraint.b2_id = b2_id
	constraint.positional_constraint.r1_lc = r1_lc
	constraint.positional_constraint.r2_lc = r2_lc
	constraint.positional_constraint.compliance = compliance
	constraint.positional_constraint.distance = distance
}

// pbd_mutual_orientation_constraint_init
func pbd_mutual_orientation_constraint_init(constraint *constraint, b1_id, b2_id bid, compliance float64) {
	constraint.ctype = mutual_ORIENTATION_CONSTRAINT
	constraint.b1_id = b1_id
	constraint.b2_id = b2_id
	constraint.mutual_orientation_constraint.compliance = compliance
}

// pbd_hinge_joint_constraint_unlimited_init
func pbd_hinge_joint_constraint_unlimited_init(constraint *constraint, b1_id, b2_id bid,
	r1_lc, r2_lc lin.V3, compliance float64, b1_aligned_axis, b2_aligned_axis pbd_Axis_Type) {
	constraint.ctype = hinge_JOINT_CONSTRAINT
	constraint.b1_id = b1_id
	constraint.b2_id = b2_id
	constraint.hinge_joint_constraint.r1_lc = r1_lc
	constraint.hinge_joint_constraint.r2_lc = r2_lc
	constraint.hinge_joint_constraint.compliance = compliance
	constraint.hinge_joint_constraint.b1_aligned_axis = b1_aligned_axis
	constraint.hinge_joint_constraint.b2_aligned_axis = b2_aligned_axis
	constraint.hinge_joint_constraint.limited = false
}

// pbd_hinge_joint_constraint_limited_init
func pbd_hinge_joint_constraint_limited_init(constraint *constraint, b1_id, b2_id bid,
	r1_lc, r2_lc lin.V3, compliance float64,
	b1_aligned_axis, b2_aligned_axis pbd_Axis_Type, b1_limit_axis, b2_limit_axis pbd_Axis_Type,
	lower_limit, upper_limit float64) {
	constraint.ctype = hinge_JOINT_CONSTRAINT
	constraint.b1_id = b1_id
	constraint.b2_id = b2_id
	constraint.hinge_joint_constraint.r1_lc = r1_lc
	constraint.hinge_joint_constraint.r2_lc = r2_lc
	constraint.hinge_joint_constraint.compliance = compliance
	constraint.hinge_joint_constraint.b1_aligned_axis = b1_aligned_axis
	constraint.hinge_joint_constraint.b2_aligned_axis = b2_aligned_axis
	constraint.hinge_joint_constraint.limited = true
	constraint.hinge_joint_constraint.b1_limit_axis = b1_limit_axis
	constraint.hinge_joint_constraint.b2_limit_axis = b2_limit_axis
	constraint.hinge_joint_constraint.lower_limit = lower_limit
	constraint.hinge_joint_constraint.upper_limit = upper_limit
}

// pbd_spherical_joint_constraint_init
func pbd_spherical_joint_constraint_init(constraint *constraint, b1_id, b2_id bid,
	r1_lc, r2_lc lin.V3, b1_swing_axis, b2_swing_axis, b1_twist_axis, b2_twist_axis pbd_Axis_Type,
	swing_lower_limit, swing_upper_limit, twist_lower_limit, twist_upper_limit float64) {
	constraint.ctype = spherical_JOINT_CONSTRAINT
	constraint.b1_id = b1_id
	constraint.b2_id = b2_id
	constraint.spherical_joint_constraint.r1_lc = r1_lc
	constraint.spherical_joint_constraint.r2_lc = r2_lc
	constraint.spherical_joint_constraint.b1_swing_axis = b1_swing_axis
	constraint.spherical_joint_constraint.b2_swing_axis = b2_swing_axis
	constraint.spherical_joint_constraint.b1_twist_axis = b1_twist_axis
	constraint.spherical_joint_constraint.b2_twist_axis = b2_twist_axis
	constraint.spherical_joint_constraint.swing_lower_limit = swing_lower_limit
	constraint.spherical_joint_constraint.swing_upper_limit = swing_upper_limit
	constraint.spherical_joint_constraint.twist_lower_limit = twist_lower_limit
	constraint.spherical_joint_constraint.twist_upper_limit = twist_upper_limit
}

// positional_constraint_solve
func positional_constraint_solve(constraint *constraint, h float64) {
	if constraint.ctype != positional_CONSTRAINT {
		slog.Error("positional_constraint_solve: invalid constraint")
		return
	}
	b1 := body_get_by_id(constraint.b1_id)
	b2 := body_get_by_id(constraint.b2_id)

	attachment_distance := lin.NewV3().Sub(&b1.world_position, &b2.world_position)
	delta_x := lin.NewV3().Sub(attachment_distance, &constraint.positional_constraint.distance)

	var pcpd position_Constraint_Preprocessed_Data
	calculate_positional_constraint_preprocessed_data(b1, b2, constraint.positional_constraint.r1_lc,
		constraint.positional_constraint.r2_lc, &pcpd)
	delta_lambda := positional_constraint_get_delta_lambda(&pcpd, h, constraint.positional_constraint.compliance,
		constraint.positional_constraint.lambda, *delta_x)
	positional_constraint_apply(&pcpd, delta_lambda, *delta_x)
	constraint.positional_constraint.lambda += delta_lambda
}

// calculate_p_til
func calculate_p_til(b *Body, r_lc lin.V3) (out lin.V3) {
	o := lin.NewV3().Add(&b.previous_world_position, lin.NewV3().MultQ(&r_lc, &b.previous_world_rotation))
	return *o
}

// calculate_p
func calculate_p(b *Body, r_lc lin.V3) (out lin.V3) {
	o := lin.NewV3().Add(&b.world_position, lin.NewV3().MultQ(&r_lc, &b.world_rotation))
	return *o
}

// collision_constraint_solve
func collision_constraint_solve(constraint *constraint, h float64) {
	if constraint.ctype != collision_CONSTRAINT {
		slog.Error("collision_constraint_solve: expecting collision constraint")
		return
	}
	b1 := body_get_by_id(constraint.b1_id)
	b2 := body_get_by_id(constraint.b2_id)

	var pcpd position_Constraint_Preprocessed_Data
	calculate_positional_constraint_preprocessed_data(b1, b2, constraint.collision_constraint.r1_lc, constraint.collision_constraint.r2_lc, &pcpd)

	// here we calculate 'p1' and 'p2' in order to calculate 'd', as stated in sec (3.5)
	p1 := lin.NewV3().Add(&b1.world_position, &pcpd.r1_wc)
	p2 := lin.NewV3().Add(&b2.world_position, &pcpd.r2_wc)
	d := lin.NewV3().Sub(p1, p2).Dot(&constraint.collision_constraint.normal)

	if d > 0.0 {
		delta_x := lin.NewV3().Scale(&constraint.collision_constraint.normal, d)
		delta_lambda := positional_constraint_get_delta_lambda(&pcpd, h, 0.0, constraint.collision_constraint.lambda_n, *delta_x)
		positional_constraint_apply(&pcpd, delta_lambda, *delta_x)
		constraint.collision_constraint.lambda_n += delta_lambda

		// Recalculate entity pair preprocessed data and p1/p2
		calculate_positional_constraint_preprocessed_data(b1, b2, constraint.collision_constraint.r1_lc, constraint.collision_constraint.r2_lc, &pcpd)

		p1.Add(&b1.world_position, &pcpd.r1_wc)
		p2.Add(&b2.world_position, &pcpd.r2_wc)
		delta_lambda = positional_constraint_get_delta_lambda(&pcpd, h, 0.0, constraint.collision_constraint.lambda_t, *delta_x)

		// We should also add a constraint for static friction, but only if lambda_t < u_s * lambda_n
		static_friction_coefficient := (b1.static_friction_coefficient + b2.static_friction_coefficient) / 2.0

		lambda_n := constraint.collision_constraint.lambda_n
		lambda_t := constraint.collision_constraint.lambda_t + delta_lambda
		// @NOTE(fek): This inequation shown in 3.5 was changed because the lambdas will always be negative!
		if lambda_t > static_friction_coefficient*lambda_n {
			p1_til := lin.NewV3().Add(&b1.previous_world_position,
				lin.NewV3().MultQ(&constraint.collision_constraint.r1_lc, &b1.previous_world_rotation))
			p2_til := lin.NewV3().Add(&b2.previous_world_position,
				lin.NewV3().MultQ(&constraint.collision_constraint.r2_lc, &b2.previous_world_rotation))
			delta_p := lin.NewV3().Sub(lin.NewV3().Sub(p1, p1_til), lin.NewV3().Sub(p2, p2_til))
			delta_p_t := lin.NewV3().Sub(delta_p, lin.NewV3().Scale(&constraint.collision_constraint.normal,
				delta_p.Dot(&constraint.collision_constraint.normal)))

			positional_constraint_apply(&pcpd, delta_lambda, *delta_p_t)
			constraint.collision_constraint.lambda_t += delta_lambda
		}
	}
}

// mutual_orientation_constraint_solve
func mutual_orientation_constraint_solve(constraint *constraint, h float64) {
	if constraint.ctype != mutual_ORIENTATION_CONSTRAINT {
		slog.Error("collision_constraint_solve: expecting mutual orientation constraint")
		return
	}
	b1 := body_get_by_id(constraint.b1_id)
	b2 := body_get_by_id(constraint.b2_id)

	var acpd angular_Constraint_Preprocessed_Data
	calculate_angular_constraint_preprocessed_data(b1, b2, &acpd)

	q2_inv := lin.NewQ().Inv(lin.NewQ().Set(&b2.world_rotation))
	aux := lin.NewQ().Mult(q2_inv, &b1.world_rotation) // apply world rot to inverse.
	delta_q := lin.NewV3().SetS(2.0*aux.X, 2.0*aux.Y, 2.0*aux.Z)

	delta_lambda := angular_constraint_get_delta_lambda(&acpd, h, constraint.mutual_orientation_constraint.compliance,
		constraint.mutual_orientation_constraint.lambda, *delta_q)
	angular_constraint_apply(&acpd, delta_lambda, *delta_q)
	constraint.mutual_orientation_constraint.lambda += delta_lambda
}

// limit_angle
func limit_angle(n, n1, n2 lin.V3, alpha, beta float64) (delta_q lin.V3, ok bool) {
	// Calculate phi, which is the angle between n1 and n2 with respect to the rotation vector n
	phi := math.Asin(n.Dot(lin.NewV3().Cross(&n1, &n2)))
	// Asin returns the angle in the interval [-pi/2,+pi/2], which is already correct if the angle between n1 and n2 is acute.
	// however, if n1 and n2 forms an obtuse angle, we need to manually differentiate. In this case, n1 dot n2 is less than 0.
	// For example, if the angle between n1 and n2 is 30 degrees, then sin(30)=0.5, but if the angle is 150, sin(150)=0.5 as well,
	// thus in both cases Asin will return 30 degrees (pi/6)

	if n1.Dot(&n2) < 0.0 {
		phi = math.Pi - phi // this will do the trick and fix the angle
	}
	// now our angle is between [-pi/2, 3pi/2].

	// maps the inner range [pi, 3pi/2] to [-pi, -pi/2]
	if phi > math.Pi {
		phi = phi - 2.0*math.Pi
	}
	// now our angle is between [-pi, pi]

	// this is useless?
	if phi < -math.Pi {
		phi = phi + 2.0*math.Pi
	}

	if phi < alpha || phi > beta {
		// at this point, phi represents the angle between n1 and n2

		// clamp phi to get the limit angle, i.e., the angle that we wanna 'be at'
		phi = lin.Clamp(phi, alpha, beta)
		// create a quaternion that represents this rotation
		rot := lin.NewQ().SetAa(n.X, n.Y, n.Z, phi)

		// rotate n1 by the limit angle, so n1 will get very close to n2, except for the extra rotation that we wanna get rid of
		n1.MultQ(&n1, rot)

		// calculate delta_q based on this extra rotation
		delta_q.Cross(&n1, &n2)
		return delta_q, true
	}
	return delta_q, false
}

// get_axis_in_world_coords
func get_axis_in_world_coords(entity_rotation *lin.Q, axis pbd_Axis_Type) (av lin.V3) {
	switch axis {
	case pbd_POSITIVE_X_AXIS:
		return *av.Right(entity_rotation)
	case pbd_NEGATIVE_X_AXIS:
		return *av.RightInverted(entity_rotation)
	case pbd_POSITIVE_Y_AXIS:
		return *av.Up(entity_rotation)
	case pbd_NEGATIVE_Y_AXIS:
		return *av.UpInverted(entity_rotation)
	case pbd_POSITIVE_Z_AXIS:
		return *av.Forward(entity_rotation)
	case pbd_NEGATIVE_Z_AXIS:
		return *av.ForwardInverted(entity_rotation)
	}
	slog.Error("get_axis_in_world_coords: impossible axis")
	return lin.V3{}
}

// hinge_joint_constraint_solve
func hinge_joint_constraint_solve(constraint *constraint, h float64) {
	if constraint.ctype != hinge_JOINT_CONSTRAINT {
		slog.Error("hinge_joint_constraint_solve: expecting hinge joint constraint")
		return
	}
	b1 := body_get_by_id(constraint.b1_id)
	b2 := body_get_by_id(constraint.b2_id)

	// Angular Constraint to make sure the aligned axis are kept aligned
	var acpd angular_Constraint_Preprocessed_Data
	calculate_angular_constraint_preprocessed_data(b1, b2, &acpd)

	b1_a_wc := get_axis_in_world_coords(&b1.world_rotation, constraint.hinge_joint_constraint.b1_aligned_axis)
	b2_a_wc := get_axis_in_world_coords(&b2.world_rotation, constraint.hinge_joint_constraint.b2_aligned_axis)
	delta_q := lin.V3{}
	delta_q.Cross(&b1_a_wc, &b2_a_wc)

	delta_lambda := angular_constraint_get_delta_lambda(&acpd, h, constraint.hinge_joint_constraint.compliance,
		constraint.hinge_joint_constraint.lambda_aligned_axes, delta_q)

	angular_constraint_apply(&acpd, delta_lambda, delta_q)
	constraint.hinge_joint_constraint.lambda_aligned_axes += delta_lambda

	// Positional constraint to ensure that the distance between both entities are correct
	// @TODO: optmize preprocessed datas
	var pcpd position_Constraint_Preprocessed_Data
	calculate_positional_constraint_preprocessed_data(b1, b2, constraint.hinge_joint_constraint.r1_lc,
		constraint.hinge_joint_constraint.r2_lc, &pcpd)

	p1 := lin.NewV3().Add(&b1.world_position, &pcpd.r1_wc)
	p2 := lin.NewV3().Add(&b2.world_position, &pcpd.r2_wc)
	delta_r := lin.NewV3().Sub(p1, p2)
	delta_x := delta_r

	delta_lambda = positional_constraint_get_delta_lambda(&pcpd, h, 0.0, constraint.hinge_joint_constraint.lambda_pos, *delta_x)
	positional_constraint_apply(&pcpd, delta_lambda, *delta_x)
	constraint.hinge_joint_constraint.lambda_pos += delta_lambda

	// Finally, angular constraint to ensure the joint angle limit is respected
	if constraint.hinge_joint_constraint.limited {
		n1 := get_axis_in_world_coords(&b1.world_rotation, constraint.hinge_joint_constraint.b1_limit_axis)
		n2 := get_axis_in_world_coords(&b2.world_rotation, constraint.hinge_joint_constraint.b2_limit_axis)
		n := get_axis_in_world_coords(&b1.world_rotation, constraint.hinge_joint_constraint.b1_aligned_axis)
		alpha := constraint.hinge_joint_constraint.lower_limit
		beta := constraint.hinge_joint_constraint.upper_limit

		delta_q, ok := limit_angle(n, n1, n2, alpha, beta)
		if ok {
			// Angular Constraint
			var acpd angular_Constraint_Preprocessed_Data
			calculate_angular_constraint_preprocessed_data(b1, b2, &acpd)

			delta_lambda := angular_constraint_get_delta_lambda(&acpd, h, 0.0, constraint.hinge_joint_constraint.lambda_limit_axes, delta_q)
			angular_constraint_apply(&acpd, delta_lambda, delta_q)
			constraint.hinge_joint_constraint.lambda_limit_axes += delta_lambda
		}
	}
}

// spherical_joint_constraint_solve
func spherical_joint_constraint_solve(constraint *constraint, h float64) {
	if constraint.ctype != spherical_JOINT_CONSTRAINT {
		slog.Error("spherical_joint_constraint_solve: expecting spherical joint constraint")
		return
	}
	const EPSILON float64 = 1e-50

	b1 := body_get_by_id(constraint.b1_id)
	b2 := body_get_by_id(constraint.b2_id)

	// Positional constraint to ensure that the distance between both entities are correct
	var pcpd position_Constraint_Preprocessed_Data
	calculate_positional_constraint_preprocessed_data(b1, b2, constraint.spherical_joint_constraint.r1_lc,
		constraint.spherical_joint_constraint.r2_lc, &pcpd)

	p1 := lin.NewV3().Add(&b1.world_position, &pcpd.r1_wc)
	p2 := lin.NewV3().Add(&b2.world_position, &pcpd.r2_wc)
	delta_r := lin.NewV3().Sub(p1, p2)
	delta_x := delta_r

	delta_lambda := positional_constraint_get_delta_lambda(&pcpd, h, 0.0, constraint.spherical_joint_constraint.lambda_pos, *delta_x)
	positional_constraint_apply(&pcpd, delta_lambda, *delta_x)
	constraint.spherical_joint_constraint.lambda_pos += delta_lambda

	// Angular constraint to ensure the swing angle limit is respected
	n1 := get_axis_in_world_coords(&b1.world_rotation, constraint.spherical_joint_constraint.b1_swing_axis)
	n2 := get_axis_in_world_coords(&b2.world_rotation, constraint.spherical_joint_constraint.b2_swing_axis)
	n := lin.NewV3().Cross(&n1, &n2)
	n_len := n.Len()

	if n_len > EPSILON {
		n.SetS(n.X/n_len, n.Y/n_len, n.Z/n_len)

		alpha := constraint.spherical_joint_constraint.swing_lower_limit
		beta := constraint.spherical_joint_constraint.swing_upper_limit

		delta_q, ok := limit_angle(*n, n1, n2, alpha, beta)
		if ok {
			// Angular Constraint
			var acpd angular_Constraint_Preprocessed_Data
			calculate_angular_constraint_preprocessed_data(b1, b2, &acpd)

			delta_lambda = angular_constraint_get_delta_lambda(&acpd, h, 0.0, constraint.spherical_joint_constraint.lambda_swing, delta_q)
			angular_constraint_apply(&acpd, delta_lambda, delta_q)
			constraint.spherical_joint_constraint.lambda_swing += delta_lambda
		}
	}

	// Angular constraint to ensure the twist angle limit is respected
	a1 := get_axis_in_world_coords(&b1.world_rotation, constraint.spherical_joint_constraint.b1_swing_axis)
	z1 := get_axis_in_world_coords(&b1.world_rotation, constraint.spherical_joint_constraint.b1_twist_axis)
	a2 := get_axis_in_world_coords(&b2.world_rotation, constraint.spherical_joint_constraint.b2_swing_axis)
	z2 := get_axis_in_world_coords(&b2.world_rotation, constraint.spherical_joint_constraint.b2_twist_axis)
	n.Add(&a1, &a2)
	n_len = n.Len()

	if n_len > EPSILON {
		n.SetS(n.X/n_len, n.Y/n_len, n.Z/n_len)

		n1.Sub(&z1, lin.NewV3().Scale(n, n.Dot(&z1)))
		n2.Sub(&z2, lin.NewV3().Scale(n, n.Dot(&z2)))
		n1_len := n1.Len()
		n2_len := n2.Len()
		if n1_len > EPSILON && n2_len > EPSILON {
			n1.SetS(n1.X/n1_len, n1.Y/n1_len, n1.Z/n1_len)
			n2.SetS(n2.X/n2_len, n2.Y/n2_len, n2.Z/n2_len)

			alpha := constraint.spherical_joint_constraint.twist_lower_limit
			beta := constraint.spherical_joint_constraint.twist_upper_limit

			delta_q, ok := limit_angle(*n, n1, n2, alpha, beta)
			if ok {
				// Angular Constraint
				var acpd angular_Constraint_Preprocessed_Data
				calculate_angular_constraint_preprocessed_data(b1, b2, &acpd)

				delta_lambda := angular_constraint_get_delta_lambda(&acpd, h, 0.0, constraint.spherical_joint_constraint.lambda_twist, delta_q)
				angular_constraint_apply(&acpd, delta_lambda, delta_q)
				constraint.spherical_joint_constraint.lambda_twist += delta_lambda
			}
		}
	}
}

// solve_constraint
func solve_constraint(constraint *constraint, h float64) {
	switch constraint.ctype {
	case positional_CONSTRAINT:
		positional_constraint_solve(constraint, h)
		return
	case collision_CONSTRAINT:
		collision_constraint_solve(constraint, h)
		return
	case mutual_ORIENTATION_CONSTRAINT:
		mutual_orientation_constraint_solve(constraint, h)
		return
	case hinge_JOINT_CONSTRAINT:
		// hinge_joint_constraint_solve(constraint, h)
		return
	case spherical_JOINT_CONSTRAINT:
		// spherical_joint_constraint_solve(constraint, h)
		return
	default:
		slog.Error("solve_constraint: unsupported constraint", "constraint_type", constraint.ctype)
	}
}

// clipping_contact_to_collision_constraint
func clipping_contact_to_collision_constraint(b1, b2 *Body, bid1, bid2 bid, contact *collider_Contact, constraint *constraint) {
	constraint.ctype = collision_CONSTRAINT
	constraint.b1_id = bid1
	constraint.b2_id = bid2
	constraint.collision_constraint.normal = contact.normal
	constraint.collision_constraint.lambda_n = 0.0
	constraint.collision_constraint.lambda_t = 0.0

	r1_wc := lin.NewV3().Sub(&contact.collision_point1, &b1.world_position)
	r2_wc := lin.NewV3().Sub(&contact.collision_point2, &b2.world_position)

	q1_inv := lin.NewQ().Inv(lin.NewQ().Set(&b1.world_rotation))
	constraint.collision_constraint.r1_lc.MultQ(r1_wc, q1_inv)

	q2_inv := lin.NewQ().Inv(lin.NewQ().Set(&b2.world_rotation))
	constraint.collision_constraint.r2_lc.MultQ(r2_wc, q2_inv)
}

// copy_constraints
func copy_constraints(constraints []constraint) []constraint {
	if constraints == nil {
		return []constraint{}
	}
	copied_constraints := make([]constraint, len(constraints))
	copy(copied_constraints, constraints)
	for i := 0; i < len(copied_constraints); i++ {
		constraint := &copied_constraints[i]

		// Reset lambda
		switch constraint.ctype {
		case positional_CONSTRAINT:
			constraint.positional_constraint.lambda = 0.0
		case collision_CONSTRAINT:
			constraint.collision_constraint.lambda_t = 0.0
			constraint.collision_constraint.lambda_n = 0.0
		case mutual_ORIENTATION_CONSTRAINT:
			constraint.mutual_orientation_constraint.lambda = 0.0
		case hinge_JOINT_CONSTRAINT:
			constraint.hinge_joint_constraint.lambda_pos = 0.0
			constraint.hinge_joint_constraint.lambda_aligned_axes = 0.0
			constraint.hinge_joint_constraint.lambda_limit_axes = 0.0
		case spherical_JOINT_CONSTRAINT:
			constraint.spherical_joint_constraint.lambda_pos = 0.0
			constraint.spherical_joint_constraint.lambda_swing = 0.0
			constraint.spherical_joint_constraint.lambda_twist = 0.0
		}
	}
	return copied_constraints
}

// pbd_simulate
func pbd_simulate(dt float64, bodies []Body, num_substeps, num_pos_iters uint32, enable_collisions bool) {
	pbd_simulate_with_constraints(dt, bodies, []constraint{}, num_substeps, num_pos_iters, enable_collisions)
}

// pbd_simulate_with_constraints
func pbd_simulate_with_constraints(dt float64, bodies []Body, external_constraints []constraint,
	num_substeps, num_pos_iters uint32, enable_collisions bool) {
	if dt <= 0.0 {
		return
	}
	h := dt / float64(num_substeps)

	broad_collision_pairs := broad_get_collision_pairs(bodies)
	simulation_islands := broad_collect_simulation_islands(bodies, broad_collision_pairs, external_constraints)

	// All entities will be contained in the simulation islands.
	// Update deactivation time and also, at the same time, its active status
	for j := 0; j < len(simulation_islands); j++ {
		simulation_island := simulation_islands[j]
		all_inactive := true
		for k := 0; k < len(simulation_island); k++ {
			b := body_get_by_id(simulation_island[k])

			linear_velocity_len := b.linear_velocity.Len()
			angular_velocity_len := b.angular_velocity.Len()
			if linear_velocity_len < linear_SLEEPING_THRESHOLD && angular_velocity_len < angular_SLEEPING_THRESHOLD {
				b.deactivation_time += dt // we should use 'dt' if doing once per frame
			} else {
				b.deactivation_time = 0.0
			}
			if b.deactivation_time < deactivation_TIME_TO_BE_INACTIVE {
				all_inactive = false
			}
		}

		// We only set entities to inactive if the whole island is inactive!
		for k := 0; k < len(simulation_island); k++ {
			b := body_get_by_id(simulation_island[k])
			b.active = !all_inactive
		}
	}
	broad_simulation_islands_destroy(simulation_islands)

	// The main loop of the PBD simulation
	for i := 0; i < int(num_substeps); i++ {
		for j := 0; j < len(bodies); j++ {
			b := &bodies[j]

			// Stores the previous position and orientation of the entity
			b.previous_world_position = b.world_position
			b.previous_world_rotation = b.world_rotation
			if b.fixed || !b.active {
				continue
			}

			// Calculate the external force and torque of the entity
			external_force := calculate_external_force(b)
			external_torque := calculate_external_torque(b)

			// Update the entity position and linear velocity based on the current velocity and applied forces
			b.linear_velocity.Add(&b.linear_velocity, external_force.Scale(&external_force, h*b.inverse_mass))
			b.world_position.Add(&b.world_position, lin.NewV3().Scale(&b.linear_velocity, h))

			// Update the entity orientation and angular velocity based on the current velocity and applied forces
			b_inverse_inertia_tensor := get_dynamic_inverse_inertia_tensor(b)
			b_inertia_tensor := get_dynamic_inertia_tensor(b)
			b.angular_velocity.Add(
				&b.angular_velocity,
				lin.NewV3().Scale(
					lin.NewV3().MultMv(
						&b_inverse_inertia_tensor,
						lin.NewV3().Sub(
							&external_torque,
							lin.NewV3().Cross(
								&b.angular_velocity,
								lin.NewV3().MultMv(&b_inertia_tensor, &b.angular_velocity)))),
					h))

			aux := lin.NewQ().SetS(b.angular_velocity.X, b.angular_velocity.Y, b.angular_velocity.Z, 0.0)
			q := lin.NewQ().Mult(&b.world_rotation, aux) // apply aux to world rotation.
			b.world_rotation.X = b.world_rotation.X + h*0.5*q.X
			b.world_rotation.Y = b.world_rotation.Y + h*0.5*q.Y
			b.world_rotation.Z = b.world_rotation.Z + h*0.5*q.Z
			b.world_rotation.W = b.world_rotation.W + h*0.5*q.W
			// should we normalize?
			b.world_rotation.Unit()
		}

		// Create the constraints array
		constraints := copy_constraints(external_constraints)

		// As explained in sec 3.5, in each substep we need to check for collisions
		if enable_collisions {
			for j := 0; j < len(broad_collision_pairs); j++ {
				bid1 := broad_collision_pairs[j].b1_id
				bid2 := broad_collision_pairs[j].b2_id
				b1 := body_get_by_id(bid1)
				b2 := body_get_by_id(bid2)

				// If b1 is "colliding" with b2, they must be either both active or both inactive
				if !b1.fixed && !b2.fixed {
					if b1.active != b2.active {
						slog.Error("pbd_simulate_with_constraints: both bodies must be active or inactive")
						continue
					}
				}

				// No need to solve the collision if both entities are either inactive or fixed
				if (b1.fixed || !b1.active) && (b2.fixed || !b2.active) {
					continue
				}

				colliders_update(b1.colliders, b1.world_position, &b1.world_rotation)
				colliders_update(b2.colliders, b2.world_position, &b2.world_rotation)
				contacts := colliders_get_contacts(b1.colliders, b2.colliders)
				for l := 0; l < len(contacts); l++ {
					contact := &contacts[l]
					var constraint constraint
					clipping_contact_to_collision_constraint(b1, b2, bid1, bid2, contact, &constraint)
					constraints = append(constraints, constraint)
				}
			}
		}

		// Now we run the PBD solver with NUM_POS_ITERS iterations
		for j := 0; j < int(num_pos_iters); j++ {
			for k := 0; k < len(constraints); k++ {
				constraint := &constraints[k]
				solve_constraint(constraint, h)
			}
		}

		// The PBD velocity update
		for j := 0; j < len(bodies); j++ {
			b := &bodies[j]
			if b.fixed || !b.active {
				continue
			}

			// We start by storing the current velocities (this is needed for the velocity solver that comes at the end of the loop)
			b.previous_linear_velocity = b.linear_velocity
			b.previous_angular_velocity = b.angular_velocity

			// Update the linear velocity based on the position difference
			b.linear_velocity.Scale(lin.NewV3().Sub(&b.world_position, &b.previous_world_position), 1.0/h)

			// Update the angular velocity based on the orientation difference
			inv := lin.NewQ().Inv(lin.NewQ().Set(&b.previous_world_rotation))
			delta_q := lin.NewQ().Mult(inv, &b.world_rotation) // apply world rot to inv
			if delta_q.W >= 0.0 {
				b.angular_velocity.Scale(lin.NewV3().SetS(delta_q.X, delta_q.Y, delta_q.Z), 2.0/h)
			} else {
				b.angular_velocity.Scale(lin.NewV3().SetS(delta_q.X, delta_q.Y, delta_q.Z), -2.0/h)
			}
		}

		// The velocity solver - we run this additional solver for every collision that we found
		for j := 0; j < len(constraints); j++ {
			constraint := &constraints[j]
			if constraint.ctype == collision_CONSTRAINT {
				b1 := body_get_by_id(constraint.b1_id)
				b2 := body_get_by_id(constraint.b2_id)
				n := constraint.collision_constraint.normal
				lambda_n := constraint.collision_constraint.lambda_n
				// lambda_t := constraint.collision_constraint.lambda_t .. not used

				var pcpd position_Constraint_Preprocessed_Data
				calculate_positional_constraint_preprocessed_data(b1, b2,
					constraint.collision_constraint.r1_lc,
					constraint.collision_constraint.r2_lc, &pcpd)
				v1 := b1.linear_velocity
				w1 := b1.angular_velocity
				v2 := b2.linear_velocity
				w2 := b2.angular_velocity
				// We start by calculating the relative normal and tangential velocities at the contact point, as described in (3.6)
				// @NOTE: equation (29) was modified here
				v := lin.NewV3().Sub(
					lin.NewV3().Add(
						&v1,
						lin.NewV3().Cross(&w1, &pcpd.r1_wc)),
					lin.NewV3().Add(
						&v2,
						lin.NewV3().Cross(&w2, &pcpd.r2_wc)))
				vn := n.Dot(v)
				vt := lin.NewV3().Sub(v, lin.NewV3().Scale(&n, vn))

				// delta_v stores the velocity change that we need to perform at the end of the solver
				delta_v := lin.NewV3()

				// we start by applying Coloumb's dynamic friction force
				dynamic_friction_coefficient := (b1.dynamic_friction_coefficient + b2.dynamic_friction_coefficient) / 2.0
				fn := lambda_n / h // simplifly h^2 by ommiting h in the next calculation
				// @NOTE: equation (30) was modified here
				fact := math.Min(dynamic_friction_coefficient*math.Abs(fn), vt.Len())
				// update delta_v
				delta_v.Add(delta_v, lin.NewV3().Scale(lin.NewV3().Set(vt).Unit(), -fact))

				// Now we handle restitution
				old_v1 := b1.previous_linear_velocity
				old_w1 := b1.previous_angular_velocity
				old_v2 := b2.previous_linear_velocity
				old_w2 := b2.previous_angular_velocity
				v_til := lin.NewV3().Sub(
					lin.NewV3().Add(
						&old_v1,
						lin.NewV3().Cross(&old_w1, &pcpd.r1_wc)),
					lin.NewV3().Add(&old_v2,
						lin.NewV3().Cross(&old_w2, &pcpd.r2_wc)))
				vn_til := n.Dot(v_til)
				e := b1.restitution_coefficient * b2.restitution_coefficient
				// @NOTE: equation (34) was modified here
				fact = -vn + math.Min(-e*vn_til, 0.0)
				// update delta_v
				delta_v.Add(delta_v, lin.NewV3().Scale(&n, fact))

				// Finally, we end the solver by applying delta_v, considering the inverse masses of both entities
				_w1 := b1.inverse_mass + lin.NewV3().Cross(&pcpd.r1_wc, &n).Dot(
					lin.NewV3().MultMv(
						&pcpd.b1_inverse_inertia_tensor,
						lin.NewV3().Cross(&pcpd.r1_wc, &n)))
				_w2 := b2.inverse_mass + lin.NewV3().Cross(&pcpd.r2_wc, &n).Dot(
					lin.NewV3().MultMv(
						&pcpd.b2_inverse_inertia_tensor,
						lin.NewV3().Cross(&pcpd.r2_wc, &n)))
				p := lin.NewV3().Scale(delta_v, 1.0/(_w1+_w2))

				if !b1.fixed {
					b1.linear_velocity.Add(&b1.linear_velocity, lin.NewV3().Scale(p, b1.inverse_mass))
					b1.angular_velocity.Add(&b1.angular_velocity,
						lin.NewV3().MultMv(&pcpd.b1_inverse_inertia_tensor, lin.NewV3().Cross(&pcpd.r1_wc, p)))
				}
				if !b2.fixed {
					b2.linear_velocity.Add(&b2.linear_velocity, lin.NewV3().Neg(lin.NewV3().Scale(p, b2.inverse_mass)))
					b2.angular_velocity.Add(&b2.angular_velocity,
						lin.NewV3().MultMv(&pcpd.b2_inverse_inertia_tensor, lin.NewV3().Neg(lin.NewV3().Cross(&pcpd.r2_wc, p))))
				}
			} else if constraint.ctype == hinge_JOINT_CONSTRAINT {
				// TODO: Joint damping
			}
		}
	}
}
