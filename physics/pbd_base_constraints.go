// Copyright © 2024 Galvanized Logic Inc.

package physics

import (
	"log/slog"

	"github.com/gazed/vu/math/lin"
)

// position_Constraint_Preprocessed_Data;
type position_Constraint_Preprocessed_Data struct {
	b1                        *Body
	b2                        *Body
	r1_wc                     lin.V3
	r2_wc                     lin.V3
	b1_inverse_inertia_tensor lin.M3
	b2_inverse_inertia_tensor lin.M3
}

// angular_Constraint_Preprocessed_Data;
type angular_Constraint_Preprocessed_Data struct {
	b1                        *Body
	b2                        *Body
	b1_inverse_inertia_tensor lin.M3
	b2_inverse_inertia_tensor lin.M3
}

// #define USE_QUATERNIONS_LINEARIZED_FORMULAS

// calculate_positional_constraint_preprocessed_data
func calculate_positional_constraint_preprocessed_data(b1, b2 *Body, r1_lc, r2_lc lin.V3, pcpd *position_Constraint_Preprocessed_Data) {
	pcpd.b1 = b1
	pcpd.b2 = b2
	pcpd.r1_wc.MultQ(&r1_lc, &b1.world_rotation)
	pcpd.r2_wc.MultQ(&r2_lc, &b2.world_rotation)
	pcpd.b1_inverse_inertia_tensor = get_dynamic_inverse_inertia_tensor(b1)
	pcpd.b2_inverse_inertia_tensor = get_dynamic_inverse_inertia_tensor(b2)
}

// positional_constraint_get_delta_lambda
func positional_constraint_get_delta_lambda(pcpd *position_Constraint_Preprocessed_Data, h, compliance, lambda float64, delta_x lin.V3) float64 {
	c := delta_x.Len()

	// We need to avoid calculations when delta_x is zero or very very close to zero, otherwise we will might run into
	// big problems because of floating-point precision
	const EPSILON float64 = 1e-50
	if c <= EPSILON {
		return 0.0
	}
	b1 := pcpd.b1
	b2 := pcpd.b2
	r1_wc := pcpd.r1_wc
	r2_wc := pcpd.r2_wc
	b1_inverse_inertia_tensor := pcpd.b1_inverse_inertia_tensor
	b2_inverse_inertia_tensor := pcpd.b2_inverse_inertia_tensor

	n := lin.NewV3().SetS(delta_x.X/c, delta_x.Y/c, delta_x.Z/c)

	// calculate the inverse masses of both entities
	w1 := b1.inverse_mass + lin.NewV3().Cross(&r1_wc, n).Dot(lin.NewV3().MultMv(&b1_inverse_inertia_tensor, lin.NewV3().Cross(&r1_wc, n)))
	w2 := b2.inverse_mass + lin.NewV3().Cross(&r2_wc, n).Dot(lin.NewV3().MultMv(&b2_inverse_inertia_tensor, lin.NewV3().Cross(&r2_wc, n)))
	if w1+w2 == 0 {
		slog.Error("positional_constraint_get_delta_lambda: inverse mass is zero")
	}

	// calculate the delta_lambda (XPBD) and updates the constraint
	til_compliance := compliance / (h * h)
	delta_lambda := (-c - til_compliance*lambda) / (w1 + w2 + til_compliance)
	return delta_lambda
}

// positional_constraint_apply
// Apply the positional constraint, updating the position and orientation of the entities accordingly
func positional_constraint_apply(pcpd *position_Constraint_Preprocessed_Data, delta_lambda float64, delta_x lin.V3) {
	c := delta_x.Len()

	// We need to avoid calculations when delta_x is zero or very very close to zero, otherwise we will might run into
	// big problems because of floating-point precision
	const EPSILON float64 = 1e-50
	if c <= EPSILON {
		return
	}

	b1 := pcpd.b1
	b2 := pcpd.b2
	r1_wc := &pcpd.r1_wc
	r2_wc := &pcpd.r2_wc
	b1_inverse_inertia_tensor := &pcpd.b1_inverse_inertia_tensor
	b2_inverse_inertia_tensor := &pcpd.b2_inverse_inertia_tensor

	n := lin.NewV3().SetS(delta_x.X/c, delta_x.Y/c, delta_x.Z/c)

	// calculates the positional impulse
	positional_impulse := lin.NewV3().Scale(n, delta_lambda)

	// updates the position of the entities based on eq (6) and (7)
	if !b1.fixed {
		b1.world_position.Add(&b1.world_position, lin.NewV3().Scale(positional_impulse, b1.inverse_mass))
	}
	if !b2.fixed {
		b2.world_position.Add(&b2.world_position, lin.NewV3().Scale(positional_impulse, -b2.inverse_mass))
	}

	// updates the rotation of the entities based on eq (8) and (9)
	aux1 := lin.NewV3().MultMv(b1_inverse_inertia_tensor, lin.NewV3().Cross(r1_wc, positional_impulse))
	aux2 := lin.NewV3().MultMv(b2_inverse_inertia_tensor, lin.NewV3().Cross(r2_wc, positional_impulse))

	aux_q1 := lin.NewQ().SetS(aux1.X, aux1.Y, aux1.Z, 0.0)
	aux_q2 := lin.NewQ().SetS(aux2.X, aux2.Y, aux2.Z, 0.0)
	q1 := lin.NewQ().Mult(&b1.world_rotation, aux_q1) // apply aux_q1 to b1 world rot
	q2 := lin.NewQ().Mult(&b2.world_rotation, aux_q2) // apply aux_q2 to b2 world rot

	if !b1.fixed {
		b1.world_rotation.X = b1.world_rotation.X + 0.5*q1.X
		b1.world_rotation.Y = b1.world_rotation.Y + 0.5*q1.Y
		b1.world_rotation.Z = b1.world_rotation.Z + 0.5*q1.Z
		b1.world_rotation.W = b1.world_rotation.W + 0.5*q1.W
		// should we normalize?
		b1.world_rotation.Unit()
	}
	if !b2.fixed {
		b2.world_rotation.X = b2.world_rotation.X - 0.5*q2.X
		b2.world_rotation.Y = b2.world_rotation.Y - 0.5*q2.Y
		b2.world_rotation.Z = b2.world_rotation.Z - 0.5*q2.Z
		b2.world_rotation.W = b2.world_rotation.W - 0.5*q2.W
		// should we normalize?
		b2.world_rotation.Unit()
	}
}

// calculate_angular_constraint_preprocessed_data
func calculate_angular_constraint_preprocessed_data(b1, b2 *Body, acpd *angular_Constraint_Preprocessed_Data) {
	acpd.b1 = b1
	acpd.b2 = b2
	acpd.b1_inverse_inertia_tensor = get_dynamic_inverse_inertia_tensor(b1)
	acpd.b2_inverse_inertia_tensor = get_dynamic_inverse_inertia_tensor(b2)
}

// angular_constraint_get_delta_lambda
func angular_constraint_get_delta_lambda(acpd *angular_Constraint_Preprocessed_Data, h, compliance, lambda float64, delta_q lin.V3) float64 {
	theta := delta_q.Len()

	// We need to avoid calculations when delta_q is zero or very very close to zero, otherwise we will might run into
	// big problems because of floating-point precision
	const EPSILON float64 = 1e-50
	if theta <= EPSILON {
		return 0.0
	}

	// b1 := acpd.b1 ... not used
	// b2 := acpd.b2 ... not used
	b1_inverse_inertia_tensor := &acpd.b1_inverse_inertia_tensor
	b2_inverse_inertia_tensor := &acpd.b2_inverse_inertia_tensor

	n := lin.NewV3().SetS(delta_q.X/theta, delta_q.Y/theta, delta_q.Z/theta)

	// calculate the inverse masses of both entities
	w1 := n.Dot(lin.NewV3().MultMv(b1_inverse_inertia_tensor, n))
	w2 := n.Dot(lin.NewV3().MultMv(b2_inverse_inertia_tensor, n))
	// 	assert(w1 + w2 != 0.0);

	// calculate the delta_lambda (XPBD) and updates the constraint
	til_compliance := compliance / (h * h)
	delta_lambda := (-theta - til_compliance*lambda) / (w1 + w2 + til_compliance)
	return delta_lambda
}

// angular_constraint_apply
// Apply the angular constraint, updating the orientation of the entities accordingly
func angular_constraint_apply(acpd *angular_Constraint_Preprocessed_Data, delta_lambda float64, delta_q lin.V3) {
	theta := delta_q.Len()

	// We need to avoid calculations when delta_q is zero or very very close to zero, otherwise we will might run into
	// big problems because of floating-point precision
	const EPSILON float64 = 1e-50
	if theta <= EPSILON {
		return
	}

	b1 := acpd.b1
	b2 := acpd.b2
	b1_inverse_inertia_tensor := &acpd.b1_inverse_inertia_tensor
	b2_inverse_inertia_tensor := &acpd.b2_inverse_inertia_tensor

	n := lin.V3{delta_q.X / theta, delta_q.Y / theta, delta_q.Z / theta}

	// calculates the positional impulse
	positional_impulse := lin.NewV3().Scale(&n, -delta_lambda)

	// updates the rotation of the entities based on eq (8) and (9)
	aux1 := lin.NewV3().MultMv(b1_inverse_inertia_tensor, positional_impulse)
	aux2 := lin.NewV3().MultMv(b2_inverse_inertia_tensor, positional_impulse)

	aux_q1 := lin.NewQ().SetS(aux1.X, aux1.Y, aux1.Z, 0.0)
	aux_q2 := lin.NewQ().SetS(aux2.X, aux2.Y, aux2.Z, 0.0)
	q1 := lin.NewQ().Mult(&b1.world_rotation, aux_q1) // apply aux_q1 to b1 world rot
	q2 := lin.NewQ().Mult(&b2.world_rotation, aux_q2) // apply aux_q2 to b2 world rot
	if !b1.fixed {
		b1.world_rotation.X = b1.world_rotation.X + 0.5*q1.X
		b1.world_rotation.Y = b1.world_rotation.Y + 0.5*q1.Y
		b1.world_rotation.Z = b1.world_rotation.Z + 0.5*q1.Z
		b1.world_rotation.W = b1.world_rotation.W + 0.5*q1.W
		// should we normalize?
		b1.world_rotation.Unit()
	}
	if !b2.fixed {
		b2.world_rotation.X = b2.world_rotation.X - 0.5*q2.X
		b2.world_rotation.Y = b2.world_rotation.Y - 0.5*q2.Y
		b2.world_rotation.Z = b2.world_rotation.Z - 0.5*q2.Z
		b2.world_rotation.W = b2.world_rotation.W - 0.5*q2.W
		// should we normalize?
		b2.world_rotation.Unit()
	}
}
