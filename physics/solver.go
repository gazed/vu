// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.
//
// Solver is a un-optimized, scaled-down, golang version of the Bullet physics
//     bullet-2.81-rev2613/src/.../btSequentialImpulseConstraintSolver.(cpp/h)
// which has the following license:
//
//    Bullet Continuous Collision Detection and Physics Library
//    Copyright (c) 2003-2006 Erwin Coumans  http://continuousphysics.com/Bullet/
//
//    This software is provided 'as-is', without any express or implied warranty.
//    In no event will the authors be held liable for any damages arising from the use of this software.
//    Permission is granted to anyone to use this software for any purpose,
//    including commercial applications, and to alter it and redistribute it freely,
//    subject to the following restrictions:
//
//    1. The origin of this software must not be misrepresented; you must not claim that you wrote the original software.
//       If you use this software in a product, an acknowledgment in the product documentation would be appreciated but is not required.
//    2. Altered source versions must be plainly marked as such, and must not be misrepresented as being the original software.
//    3. This notice may not be removed or altered from any source distribution.

package physics

import (
	"log"
	"math"

	"github.com/gazed/vu/math/lin"
)

// solver calculates the linear and angular velocity for each colliding
// body that satisfies a given set of linear equations (constraints).
// The constraints are generated from the current list of contacting
// pairs and points. The solution technique used is PGS.
//
// Projected Gauss-Seidel (PGS) solves linear complementarity problems (LCP).
//    http://en.wikipedia.org/wiki/Linear_complementarity_problem
//    http://en.wikipedia.org/wiki/Gauss–Seidel_method
//    http://image.diku.dk/kenny/download/vriphys10_course/lcp.pdf
//    http://stackoverflow.com/questions/11719704/projected-gauss-seidel-for-lcp
//    http://www.cs.ubc.ca/labs/sensorimotor/projects/sp_sigasia08/
// From wikipedia:
//    "The Gauss–Seidel method is an iterative technique for solving a
//     square system of n linear equations"
type solver struct {
	info   *solverInfo         // Constants for the solver.
	constC []*solverConstraint // Contact related equations.
	constF []*solverConstraint // Friction related equations.

	// scratch variables are optimizations that avoid creating/destroying
	// temporary objects that are needed each timestep.
	v0, v1, v2 *lin.V3 // scratch vectors.
	ra, rb     *lin.V3 // scratch relative positions for converting contacts.
}

// newSolver creates the necessary space for the solver to work.
// This is expected to be called once on engine startup.
func newSolver() *solver {
	sol := &solver{}
	sol.info = newSolverInfo()
	sol.constC = []*solverConstraint{}
	sol.constF = []*solverConstraint{}
	sol.v0 = lin.NewV3()
	sol.v1 = lin.NewV3()
	sol.v2 = lin.NewV3()
	sol.ra = lin.NewV3()
	sol.rb = lin.NewV3()
	return sol
}

// solve is expected to be called each physics update. It creates constraints
// based on contact points and then solves the constraints by adjusting bodies
// velocities to satisfy the constraints.
func (sol *solver) solve(bodies map[uint32]*body, contactPairs map[uint64]*contactPair) {
	sol.setupConstraints(bodies, contactPairs)
	sol.solveIterations(sol.info)
	sol.finish(bodies, sol.info)
}

// solver top level definitions and kick-off.
// ============================================================================
// solver setup and initialization creates the system of equations (constraints)
// that need to be solved.

// setupConstraints ensures all data is properly initialized before the solver
// starts. It sets up the contact and friction constraints based on a list of
// bodies and the complete list of all contact information.
func (sol *solver) setupConstraints(bodies map[uint32]*body, contactPairs map[uint64]*contactPair) {

	// Create solver specific information for each movable body.
	// Static bodies do not have associated solver bodies.
	for _, b := range bodies {
		if sb := b.initSolverBody(); sb.oBody != nil {
			{ // scratch v0
				sb.linearVelocity.Add(sb.linearVelocity, sol.v0.Scale(b.lfor, b.imass*sol.info.timestep))
				sb.angularVelocity.Add(sb.angularVelocity, sol.v0.MultMv(b.iitw, b.afor).Scale(sol.v0, sol.info.timestep))
			} // scratch v0 free
		}
	}

	// Reset the solver constraint holders, keeping allocated memory.
	sol.constC = sol.constC[0:0]
	sol.constF = sol.constF[0:0]

	// Generate the solver constraints for each contact pair.
	for _, contactPair := range contactPairs {
		sol.convertContacts(contactPair, sol.info)
	}
}

// convertContacts generates solver constraints from the given contacting pair.
func (sol *solver) convertContacts(pair *contactPair, info *solverInfo) {
	bodyA, bodyB := pair.bodyA, pair.bodyB
	sbodA, sbodB := bodyA.sbod, bodyB.sbod
	if (sbodA == nil || sbodA.oBody == nil) && (sbodB == nil || sbodB.oBody == nil) {
		log.Printf("Dev error: ignoring collision between two static bodies.")
		return
	}

	// turn each of the contact points into two solver constraints:
	//   one solver constraint for the contact itself.
	//   one solver constraint for friction.
	for _, poc := range pair.pocs {
		if poc.sp.distance > pair.processingLimit {
			continue // don't create constraints for non-contacting points.
		}

		// Setup the contact constraint.
		ccon := poc.sp.constC0
		ccon.sbodA, ccon.sbodB = sbodA, sbodB
		ccon.oPoint = poc
		{ // scratch poc.sp.vel, ra, rb
			relPosA := sol.ra.Sub(poc.sp.worldA, sbodA.world.Loc)
			relPosB := sol.rb.Sub(poc.sp.worldB, sbodB.world.Loc)
			rvel := sol.setupContactConstraint(ccon, sbodA, sbodB, poc, info, relPosA, relPosB, poc.sp.vel)
			sol.constC = append(sol.constC, ccon)

			// Setup the friction constraint.
			fcon := poc.sp.constF0
			fcon.frictionIndex, ccon.frictionIndex = ccon, fcon
			{ // scratch v0
				poc.sp.lateralFrictionDir.Sub(poc.sp.vel, sol.v0.Scale(poc.sp.normalWorldB, rvel))
			} // scratch v0 free
			lateralRelativeVelocity := poc.sp.lateralFrictionDir.LenSqr()
			if lateralRelativeVelocity > lin.Epsilon {
				poc.sp.lateralFrictionDir.Scale(poc.sp.lateralFrictionDir, 1.0/math.Sqrt(lateralRelativeVelocity))
			} else {
				{ // scratch v0
					poc.sp.normalWorldB.Plane(poc.sp.lateralFrictionDir, sol.v0)
				} // scratch v0 free
			}
			sol.setupFrictionConstraint(fcon, poc.sp.lateralFrictionDir, sbodA, sbodB, poc.sp, relPosA, relPosB)
			sol.constF = append(sol.constF, fcon)
		} // scratch poc.sp.vel, ra, rb free
	}
}

// setupContactConstraint initializes contact based constraints.
// Expected to be called on solver setup for each contact point.
func (sol *solver) setupContactConstraint(sc *solverConstraint, sbodA, sbodB *solverBody,
	poc *pointOfContact, info *solverInfo, relPosA, relPosB, vel *lin.V3) (relativeVelocity float64) {
	bodyA, bodyB := sbodA.oBody, sbodB.oBody // either may be nil if body is static.
	{                                        // scratch v0, v1
		torqueAxis0 := sol.v0.Cross(relPosA, poc.sp.normalWorldB)
		sc.angularComponentA.SetS(0, 0, 0)
		if bodyA != nil {
			sc.angularComponentA.MultMv(bodyA.iitw, torqueAxis0)
		}
		torqueAxis1 := sol.v1.Cross(relPosB, poc.sp.normalWorldB)
		sc.angularComponentB.SetS(0, 0, 0)
		if bodyB != nil { // scratch v2
			sc.angularComponentB.MultMv(bodyB.iitw, sol.v2.Neg(torqueAxis1))
		} // scratch v2 free

		denom0, denom1 := 0.0, 0.0
		if bodyA != nil { // scratch v2
			vec := sol.v2.Cross(sc.angularComponentA, relPosA)
			denom0 = bodyA.imass + poc.sp.normalWorldB.Dot(vec)
		} // scratch v2 free
		if bodyB != nil { // scratch v2
			sol.v2.Neg(sc.angularComponentB).Cross(sol.v2, relPosB)
			denom1 = bodyB.imass + poc.sp.normalWorldB.Dot(sol.v2)
		} // scratch v2 free
		relaxation := 1.0
		sc.jacDiagABInv = relaxation / (denom0 + denom1)
		sc.normal.Set(poc.sp.normalWorldB)
		sc.relpos1CrossNormal.Set(torqueAxis0)
		sc.relpos2CrossNormal.Neg(torqueAxis1)
	} // scratch v0, v1 free

	// Calculate penetration, friction, and restitution.
	penetration := poc.sp.distance + info.linearSlop
	{ // scratch v0, v1
		v0, v1 := sol.v0.SetS(0, 0, 0), sol.v1.SetS(0, 0, 0)
		if bodyA != nil {
			bodyA.getVelocityInLocalPoint(relPosA, v0)
		}
		if bodyB != nil {
			bodyB.getVelocityInLocalPoint(relPosB, v1)
		}
		vel.Sub(v0, v1)
	} // scratch v0, v1 free
	sc.friction = poc.sp.combinedFriction
	relativeVelocity = poc.sp.normalWorldB.Dot(vel)
	restitution := poc.sp.combinedRestitution * -relativeVelocity
	if restitution <= 0.0 {
		restitution = 0.0
	}

	// Warm start uses the previously applied impulse as an initial guess.
	sc.appliedImpulse = poc.sp.warmImpulse * info.warmstartingFactor
	{ // scratch v0, v1
		linc, angc := sol.v0, sol.v1
		if bodyA != nil {
			sbodA.applyImpulse(linc.Scale(sc.normal, bodyA.imass), angc.Set(sc.angularComponentA), sc.appliedImpulse)
		}
		if bodyB != nil {
			sbodB.applyImpulse(linc.Scale(sc.normal, bodyB.imass), angc.Neg(sc.angularComponentB), -sc.appliedImpulse)
		}
	} // scratch v0, v1 free
	sc.appliedPushImpulse = 0.0

	velocityError := 0.0
	vel1Dotn, vel2Dotn := 0.0, 0.0
	if bodyA != nil {
		vel1Dotn = sc.normal.Dot(sbodA.linearVelocity) + sc.relpos1CrossNormal.Dot(sbodA.angularVelocity)
	}
	if bodyB != nil { // scratch v0
		vel2Dotn = sol.v0.Neg(sc.normal).Dot(sbodB.linearVelocity) + sc.relpos2CrossNormal.Dot(sbodB.angularVelocity)
	} // scratch v0 free
	velocityError = restitution - (vel1Dotn + vel2Dotn)
	erp := info.erp2
	if !info.splitImpulse || (penetration > info.splitImpulsePenetrationLimit) {
		erp = info.erp
	}
	positionalError := 0.0
	if penetration > 0 {
		velocityError -= penetration / info.timestep
	} else {
		positionalError = -penetration * erp / info.timestep
	}
	penetrationImpulse := positionalError * sc.jacDiagABInv
	velocityImpulse := velocityError * sc.jacDiagABInv
	if !info.splitImpulse || penetration > info.splitImpulsePenetrationLimit {

		// combine position and velocity into rhs
		sc.rhs = penetrationImpulse + velocityImpulse
		sc.rhsPenetration = 0.0
	} else {

		// split position and velocity into rhs and m_rhsPenetration
		sc.rhs = velocityImpulse
		sc.rhsPenetration = penetrationImpulse
	}
	sc.cfm = 0
	sc.lowerLimit = 0
	sc.upperLimit = 1e10
	return relativeVelocity
}

// setupFrictionConstraint initializes contact based constraints.
// Expected to be called on solver setup for each point of contact.
func (sol *solver) setupFrictionConstraint(sc *solverConstraint, normalAxis *lin.V3, sbodA, sbodB *solverBody,
	sp *solverPoint, relPosA, relPosB *lin.V3) {
	bodyA, bodyB := sbodA.oBody, sbodB.oBody // either may be nil if body is static.
	sc.sbodA, sc.sbodB = sbodA, sbodB
	sc.normal.Set(normalAxis)
	sc.friction = sp.combinedFriction
	sc.oPoint = nil
	sc.appliedImpulse = 0.0
	sc.appliedPushImpulse = 0.0

	// compute torque
	ftorqueAxis := sc.relpos1CrossNormal.Cross(relPosA, sc.normal)
	sc.angularComponentA.SetS(0, 0, 0)
	if bodyA != nil {
		sc.angularComponentA.MultMv(bodyA.iitw, ftorqueAxis)
	}
	{ // scratch v0
		ftorqueAxis = sc.relpos2CrossNormal.Cross(relPosB, sol.v0.Neg(sc.normal))
	} // scratch v0 free
	sc.angularComponentB.SetS(0, 0, 0)
	if bodyB != nil {
		sc.angularComponentB.MultMv(bodyB.iitw, ftorqueAxis)
	}

	// compute sc.jacDiagABInv
	denom0, denom1 := 0.0, 0.0
	if bodyA != nil { // scratch v0
		sol.v0.Cross(sc.angularComponentA, relPosA)
		denom0 = bodyA.imass + normalAxis.Dot(sol.v0)
	} // scratch v0 free
	if bodyB != nil { // scratch v0, v1
		sol.v0.Cross(sol.v1.Neg(sc.angularComponentB), relPosB)
		denom1 = bodyB.imass + normalAxis.Dot(sol.v0)
	} // scratch v0, v1 free
	relaxation := 1.0
	sc.jacDiagABInv = relaxation / (denom0 + denom1)

	// compute limits.
	vel1Dotn, vel2Dotn := 0.0, 0.0
	if bodyA != nil {
		vel1Dotn = sc.normal.Dot(sbodA.linearVelocity) + sc.relpos1CrossNormal.Dot(sbodA.angularVelocity)
	}
	if bodyB != nil { // scratch v0
		vel2Dotn = sol.v0.Neg(sc.normal).Dot(sbodB.linearVelocity) + sc.relpos2CrossNormal.Dot(sbodB.angularVelocity)
	} // scratch v0 free
	velocityError := -(vel1Dotn + vel2Dotn)  // negative relative velocity
	sc.rhs = velocityError * sc.jacDiagABInv // velocity impulse
	sc.cfm = 0
	sc.lowerLimit = 0
	sc.upperLimit = 1e10
	sc.rhsPenetration = 0
}

// solver setup and initialization
// =============================================================================
// solver solution methods are used iteratively once the system of equations
// (constraints) have been set up.

// solveIterations find the velocity solutions for the given set of constraints.
func (sol *solver) solveIterations(info *solverInfo) {

	// Special step to solve penetrations just for contact constraints.
	// This updates the push and turn velocities used to correct penetration.
	if info.splitImpulse {
		for iteration := 0; iteration < info.numIterations; iteration++ {
			for _, sc := range sol.constC {
				sol.resolveSplitPenetrationImpulse(sc.sbodA, sc.sbodB, sc)
			}
		}
	}

	// Solve all constraints. This updates the delta linear and angular velocities.
	maxIterations := info.numIterations
	for iteration := 0; iteration < maxIterations; iteration++ {
		sol.solveSingleIteration(iteration, info)
	}
}

// solveSingleIteration processes each constraint once. The end result will be updated
// solverBody deltaVelocity values that better match all the constraints.
func (sol *solver) solveSingleIteration(iteration int, info *solverInfo) {
	if iteration < info.numIterations {
		for _, sc := range sol.constC {
			sol.resolveSingleConstraint(sc.sbodA, sc.sbodB, sc, true)
		}
		for _, sc := range sol.constF {
			totalImpulse := sc.frictionIndex.appliedImpulse // contact contraint impulse
			if totalImpulse > 0 {
				sc.lowerLimit = -(sc.friction * totalImpulse)
				sc.upperLimit = sc.friction * totalImpulse
				sol.resolveSingleConstraint(sc.sbodA, sc.sbodB, sc, false)
			}

		}
	}
}

// resolveSingleConstraint uses Projected Gauss Seidel or the equivalent
// sequential impulse to find delta velocities that solve the given constraint.
func (sol *solver) resolveSingleConstraint(sbod1, sbod2 *solverBody, sc *solverConstraint, doUpper bool) {
	deltaImpulse := sc.rhs - sc.appliedImpulse*sc.cfm
	{ // scratch v0
		deltaVel1Dotn := sc.normal.Dot(sbod1.deltaLinearVelocity) + sc.relpos1CrossNormal.Dot(sbod1.deltaAngularVelocity)
		deltaVel2Dotn := sol.v0.Neg(sc.normal).Dot(sbod2.deltaLinearVelocity) + sc.relpos2CrossNormal.Dot(sbod2.deltaAngularVelocity)
		deltaImpulse -= deltaVel1Dotn * sc.jacDiagABInv
		deltaImpulse -= deltaVel2Dotn * sc.jacDiagABInv
	} // scratch v0 free
	sum := sc.appliedImpulse + deltaImpulse
	if sum < sc.lowerLimit {
		deltaImpulse = sc.lowerLimit - sc.appliedImpulse
		sc.appliedImpulse = sc.lowerLimit
	} else if doUpper && sum > sc.upperLimit {
		deltaImpulse = sc.upperLimit - sc.appliedImpulse
		sc.appliedImpulse = sc.upperLimit
	} else {
		sc.appliedImpulse = sum
	}
	{ // scratch v0, v1
		linc, angc := sol.v0, sol.v1
		sbod1.applyImpulse(linc.Mult(sc.normal, sbod1.invMass), angc.Set(sc.angularComponentA), deltaImpulse)
		sbod2.applyImpulse(linc.Mult(linc.Neg(sc.normal), sbod2.invMass), angc.Set(sc.angularComponentB), deltaImpulse)
	} // scratch v0, v1 free
}

// resolveSplitPenetrationImpulse uses push and turn impulses to separate inter-penetrating bodies.
func (sol *solver) resolveSplitPenetrationImpulse(sbod1, sbod2 *solverBody, sc *solverConstraint) {
	if sc.rhsPenetration != 0 {
		deltaImpulse := sc.rhsPenetration - sc.appliedPushImpulse*sc.cfm
		{ // scratch v0
			deltaVel1Dotn := sc.normal.Dot(sbod1.pushVelocity) + sc.relpos1CrossNormal.Dot(sbod1.turnVelocity)
			deltaVel2Dotn := sol.v0.Neg(sc.normal).Dot(sbod2.pushVelocity) + sc.relpos2CrossNormal.Dot(sbod2.turnVelocity)
			deltaImpulse -= deltaVel1Dotn * sc.jacDiagABInv
			deltaImpulse -= deltaVel2Dotn * sc.jacDiagABInv
		} // scratch v0 free
		sum := sc.appliedPushImpulse + deltaImpulse
		if sum < sc.lowerLimit {
			deltaImpulse = sc.lowerLimit - sc.appliedPushImpulse
			sc.appliedPushImpulse = sc.lowerLimit
		} else {
			sc.appliedPushImpulse = sum
		}
		{ // scratch v0, v1
			linc, angc := sol.v0, sol.v1
			sbod1.applyPushImpulse(linc.Mult(sc.normal, sbod1.invMass), angc.Set(sc.angularComponentA), deltaImpulse)
			sbod2.applyPushImpulse(linc.Mult(linc.Neg(sc.normal), sbod2.invMass), angc.Set(sc.angularComponentB), deltaImpulse)
		} // scratch v0, v1 free
	}
}

// finish incorporates the velocities calculated by the solver back into the original body.
func (sol *solver) finish(bodies map[uint32]*body, info *solverInfo) {

	// save the applied impulse for future contacts (warm start).
	for _, sc := range sol.constC {
		sc.oPoint.sp.warmImpulse = sc.appliedImpulse
	}

	// update the active bodies velocities from their corresponding solver bodies.
	for _, b := range bodies {
		if b.movable {
			// Update the solverBody velocities from the deltaVelocities.
			if info.splitImpulse {
				b.sbod.writebackVelocityAndTransform(info.timestep, info.splitImpulseTurnErp)
			} else {
				b.sbod.writebackVelocity()
			}

			// Update the original body velocities from the solverBody velocities.
			// Update the world transform if the body needs to be separated from other bodies.
			b.lvel.Set(b.sbod.linearVelocity)
			b.avel.Set(b.sbod.angularVelocity)
			if info.splitImpulse {
				b.world.Set(b.sbod.world)
			}
		}
	}
}

// solver solutions
// ============================================================================
// solverInfo

// solverInfo holds fixed value parameters that act as controls
// for solver results.
type solverInfo struct {
	numIterations                int
	damping                      float64
	friction                     float64
	timestep                     float64
	restitution                  float64
	maxErrorReduction            float64
	erp                          float64 // used as Baumgarte factor
	erp2                         float64 // used in split impulse
	splitImpulseTurnErp          float64
	linearSlop                   float64
	warmstartingFactor           float64 // damps previous applied impluses.
	splitImpulsePenetrationLimit float64
	splitImpulse                 bool
}

// newSolverInfo initializes the solver information.
func newSolverInfo() *solverInfo {
	si := &solverInfo{}
	si.damping = 1.0
	si.friction = 0.3
	si.timestep = 1.0 / 50.0
	si.restitution = 0.0
	si.maxErrorReduction = 20.0
	si.numIterations = 10
	si.erp = 0.2
	si.erp2 = 0.8
	si.splitImpulse = true
	si.splitImpulsePenetrationLimit = -0.04
	si.splitImpulseTurnErp = 0.1
	si.linearSlop = 0.0
	si.warmstartingFactor = 0.85
	return si
}

// solverInfo
// ============================================================================
// solverConstraint

// solverConstraint is a one-dimensional constraint along a normal axis between
// two bodies. It is used to solve contact and friction constraints.
type solverConstraint struct {
	normal             *lin.V3
	relpos1CrossNormal *lin.V3
	relpos2CrossNormal *lin.V3
	angularComponentA  *lin.V3
	angularComponentB  *lin.V3
	appliedPushImpulse float64
	appliedImpulse     float64
	friction           float64
	jacDiagABInv       float64
	rhs                float64
	cfm                float64
	lowerLimit         float64
	upperLimit         float64
	rhsPenetration     float64
	oPoint             *pointOfContact
	sbodA              *solverBody
	sbodB              *solverBody

	// if this is frictionConstraint then it points to contactConstraint.
	// if this is contactConstraint then it points to frictionConstraint.
	frictionIndex *solverConstraint
}

// newSolverConstraint allocates the memory needed for a solver constraint.
func newSolverConstraint() *solverConstraint {
	sc := &solverConstraint{}
	sc.normal = lin.NewV3()
	sc.relpos1CrossNormal = lin.NewV3()
	sc.relpos2CrossNormal = lin.NewV3()
	sc.angularComponentA = lin.NewV3()
	sc.angularComponentB = lin.NewV3()
	return sc
}

// solverConstraint
// ============================================================================
// solverBody

// solverBody is used to attach extra solver information to Body objects.
type solverBody struct {
	oBody                *body // reference to original body
	world                *lin.T
	linearVelocity       *lin.V3
	angularVelocity      *lin.V3
	deltaLinearVelocity  *lin.V3
	deltaAngularVelocity *lin.V3
	pushVelocity         *lin.V3
	turnVelocity         *lin.V3
	invMass              *lin.V3
	t0                   *lin.T  // scratch
	v0                   *lin.V3 // scratch
}

// Create a single fixed solver body since they are the same for
// all fixed bodies and nothing should ever update them.
var fsb *solverBody

// fixedSolverBody lazy initializes and returns the single fixed
// solver body that is used by all static solver bodies.
func fixedSolverBody() *solverBody {
	if fsb == nil {
		fsb = &solverBody{}
		fsb.oBody = nil
		fsb.world = lin.NewT().SetI()
		fsb.linearVelocity = lin.NewV3()
		fsb.angularVelocity = lin.NewV3()
		fsb.deltaLinearVelocity = lin.NewV3()
		fsb.deltaAngularVelocity = lin.NewV3()
		fsb.pushVelocity = lin.NewV3()
		fsb.turnVelocity = lin.NewV3()
		fsb.invMass = lin.NewV3()
		fsb.t0 = lin.NewT()
		fsb.v0 = lin.NewV3()
	}
	return fsb
}

// newSolverBody allocates space for body specific solver information.
// This is expected to be called for a movable body, ie. one that has mass
// and can have velocity.
func newSolverBody(bod *body) *solverBody {
	sb := &solverBody{}
	sb.oBody = bod // reference
	sb.world = lin.NewT().Set(bod.world)
	sb.linearVelocity = lin.NewV3().Set(bod.lvel)
	sb.angularVelocity = lin.NewV3().Set(bod.avel)
	sb.deltaLinearVelocity = lin.NewV3()
	sb.deltaAngularVelocity = lin.NewV3()
	sb.pushVelocity = lin.NewV3()
	sb.turnVelocity = lin.NewV3()
	sb.invMass = lin.NewV3().SetS(bod.imass, bod.imass, bod.imass)
	sb.t0 = lin.NewT()
	sb.v0 = lin.NewV3()
	return sb
}

// reset updates an existing solverBody with new body information.
func (sb *solverBody) reset(bod *body) {
	sb.oBody = bod
	sb.world.Set(bod.world)
	sb.linearVelocity.Set(bod.lvel)
	sb.angularVelocity.Set(bod.avel)
	sb.deltaLinearVelocity.SetS(0, 0, 0)
	sb.deltaAngularVelocity.SetS(0, 0, 0)
	sb.pushVelocity.SetS(0, 0, 0)
	sb.turnVelocity.SetS(0, 0, 0)
	sb.invMass.SetS(bod.imass, bod.imass, bod.imass)
}

// applyPushImpulse updates the push and turn velocity used to separate
// inter-penetrating bodies.
func (sb *solverBody) applyPushImpulse(linearComponent, angularComponent *lin.V3, impulseMagnitude float64) {
	if sb.oBody != nil {
		sb.pushVelocity.Add(sb.pushVelocity, linearComponent.Scale(linearComponent, impulseMagnitude))
		sb.turnVelocity.Add(sb.turnVelocity, angularComponent.Scale(angularComponent, impulseMagnitude))
	}
}

// applyImpulse updates the linear and angular velocity change needed to
// solve constraints.
func (sb *solverBody) applyImpulse(linearComponent, angularComponent *lin.V3, impulseMagnitude float64) {
	if sb.oBody != nil {
		sb.deltaLinearVelocity.Add(sb.deltaLinearVelocity, linearComponent.Scale(linearComponent, impulseMagnitude))
		sb.deltaAngularVelocity.Add(sb.deltaAngularVelocity, angularComponent.Scale(angularComponent, impulseMagnitude))
	}
}

// writebackVelocity uses the delta velocities calculated by the solver to
// update the bodies existing linear and angular velocities.
func (sb *solverBody) writebackVelocity() {
	if sb.oBody != nil {
		sb.linearVelocity.Add(sb.linearVelocity, sb.deltaLinearVelocity)
		sb.angularVelocity.Add(sb.angularVelocity, sb.deltaAngularVelocity)
	}
}

// writebackVelocityAndTransform is called when there are position changes needed to
// resolve body penentration. This is done by applying the push and turn velocities
// to the current world transform.
func (sb *solverBody) writebackVelocityAndTransform(timestep, splitImpulseTurnErp float64) {
	if sb.oBody != nil {
		sb.linearVelocity.Add(sb.linearVelocity, sb.deltaLinearVelocity)
		sb.angularVelocity.Add(sb.angularVelocity, sb.deltaAngularVelocity)

		// correct the position/orientation based on push/turn recovery
		if sb.pushVelocity.X != 0 || sb.pushVelocity.Y != 0 || sb.pushVelocity.Z != 0 ||
			sb.turnVelocity.X != 0 || sb.turnVelocity.Y != 0 || sb.turnVelocity.Z != 0 {
			{ // scratch v0, t0
				turnVelocity := sb.v0.Scale(sb.turnVelocity, splitImpulseTurnErp)
				sb.t0.Integrate(sb.world, sb.pushVelocity, turnVelocity, timestep)
				sb.world.Set(sb.t0)
			} // scratch v0, t0 free
		}
	}
}

// solverBody
// ============================================================================
// solverPoint

// solverPoint amalgamates information from contactPair and pointOfContact
// for easy access by the solver. Where necessary there is one solverPoint
// initialized for each pointOfContact.
type solverPoint struct {
	localA              *lin.V3 // Point of contact for A in A's local space.
	localB              *lin.V3 // Point of contact for B in B's local space.
	worldB              *lin.V3 // Point of contact for A in world space.
	worldA              *lin.V3 // Point of contact for B in world space.
	normalWorldB        *lin.V3 // Point of contact in world space.
	lateralFrictionDir  *lin.V3 // Normal axis in friction constraint.
	distance            float64 // Distance between A and B.
	combinedFriction    float64 // Total friction.
	combinedRestitution float64 // Total restitution.
	warmImpulse         float64 // Saved warm start impulse (previous impulse).

	// Each solver point allocates reusable solver constraints.
	// These act like a constraint pool.
	constC0 *solverConstraint // Contact contraint, one per point.
	constF0 *solverConstraint // Friction contraint, one per point.
	vel     *lin.V3           // scratch vector needed by solver setup.
}

// newSolverPoint allocates memory for a solverPoint.
func newSolverPoint() *solverPoint {
	sp := &solverPoint{}
	sp.localA = &lin.V3{}
	sp.localB = &lin.V3{}
	sp.worldA = &lin.V3{}
	sp.worldB = &lin.V3{}
	sp.normalWorldB = &lin.V3{}
	sp.lateralFrictionDir = &lin.V3{}
	sp.warmImpulse = 0

	// allocate scratch space for solver constraints.
	sp.constC0 = newSolverConstraint()
	sp.constF0 = newSolverConstraint()
	sp.vel = &lin.V3{}
	return sp
}

// reuse is expected to be used to transfer old solver point information
// to the current solver point. poc.prepForSolver has already updated
// all the other fields.
func (sp *solverPoint) reuse(oldp *solverPoint) {
	sp.warmImpulse = oldp.warmImpulse // set to 0 to disable warm starting.
}

// set updates sp to have a copy of the given solverPoint information.
func (sp *solverPoint) set(s0 *solverPoint) {
	sp.localA.Set(s0.localA)
	sp.localB.Set(s0.localB)
	sp.worldA.Set(s0.worldA)
	sp.worldB.Set(s0.worldB)
	sp.normalWorldB.Set(s0.normalWorldB)
	sp.lateralFrictionDir.Set(s0.lateralFrictionDir)
	sp.distance = s0.distance
	sp.combinedFriction = s0.combinedFriction
	sp.combinedRestitution = s0.combinedRestitution
	sp.warmImpulse = s0.warmImpulse
}
