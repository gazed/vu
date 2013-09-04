// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package physics

// Motion is simulated using the Runge Kutta order 4 integrator. See:
//      http://stackoverflow.com/questions/1668098/runge-kutta-rk4-integration-for-game-physics
//      http://gafferongames.com/game-physics/integration-basics
//      http://www.cs.cmu.edu/~baraff/pbm/diffyq.pdf for RK4 formula.
//      http://www.cs.cmu.edu/~baraff/pbm/rigid1.pdf for physics.
//
// Using the above information linear movement is calculated as:
//    p = vm       linear momentum (p) is mass times velocity.
//    F = dp/dt    force is linear momentum (p) over time
//    v = p/m      velocity is linear momentum divided by mass
//    v = dx/dt    velocity is location (x) over time (meters per second)
// There are corresponding variables that can be used for angular momentum.
//
// Moving a rigid body is done by applying force at a point.  The force
// affects linear momentum and angular momentum separately as follows:
//    Flinear = F                 linear component of the force
// 	  Ftorque = F cross (p - x)   rotational component of the force
// For linear motion the entire force is applied directly to angular momentum.
// For angular motion the torque is calculated based on based on the cross
// product of the force vector and the point (p) on the object relative to the
// center of mass (x) of the object.
//
// Forces like gravity and friction are added through the Forces structure and
// associated methods.  In general:
//    Force applied: object moves if the force can overcome gravity and friction.
//    Force removed: object moves until gravity and friction slows it to a stop.
//    Torque applied: object rotates if the torque can overcome gravity and friction.
//    Torque removed: object rotates until gravity and friction slows it to a stop.
//
// Gravity applies a uniform force down (-Y direction).  Because it is uniform
// it has no effect on torque.  Of course there had better be a floor to collide
// with or the object will drop forever.
//     Fg = m * a
//     gravity = mass * m_gravity_acceleration;
import (
	"vu/math/lin"
)

// motion holds attributes used to simulate real-world motion on a rigid body.
// Rigid body motion allows the linear motion and rotational motion to be
// calculated separately.
type motion struct {

	// loc is the location of the center of mass in world coordinates (meters).
	// Loc is updated by the linear momentum of the object.  It can be changed
	// directly to teleport an object to a location.
	// The allociation/creation is expected to be provided through NewMotion.
	loc *lin.V3

	// dir is the current direction and amount of rotation around the axis
	// indicated by the direction. Dir (orientation: angle/axis) is represented
	// as a unit quaternion.  The allociation/creation is expected to be provided
	// through NewMotion.
	dir *lin.Q

	// linm is linear momentum.  This is the primary motion attribute in that
	// applying force to this object directly changes linear momentum.
	// i.e. the change in linear momentum is equivalent to the total force
	// acting on a body.  Linear momentum is mass*velocity so its units are
	// in kilogram meters per second.
	linm *lin.V3

	// angm is angular momentum. This is the primary rotation attribute in that
	// applying torque to this object directly changes angular momentum.
	// i.e. the change in angular momentum is equivalent to the total torque
	// acting on a body.  Angular momentum is used because it is constant over
	// time while angular velocity is not.  In other words, a body with no torque
	// will have constant angular momentum over time, but its angular velocity
	// may change.
	angm *lin.V3

	// linv is linear velocity. It is calculated from linear momentum and has
	// units in meters per second.
	linv *lin.V3

	// angv is angular velocity. It is calculated from angular momentum.
	angv *lin.V3

	// mass is the mass of the body in kilograms.  It affects both linear
	// and angular momentum.
	mass float32

	// imass, inverse mass is used to convert momentum to velocity.
	// This is calcuated once on object creation.
	imass float32

	// spin is used to track the rate of change of the rotation. Spin is updated
	// whenever rotation or angular momentum is changed.  It is used to seed
	// the derivative calculations which result in updates in the body's
	// current rotation.
	spin *lin.Q

	// angular mass (inertia tensor) describes how the mass in a body is
	// distributed relative to the bodys center of mass.  It is a scaling
	// factor between angular momentum and angular velocity and is
	// constant over a simulation of movement. Normally it is a matrix,
	// but has beeen simplified it to a single value as we're currently
	// only dealing with cubes.
	amass float32

	// iamass (inverse angular mass) is used to convert angular momentum to
	// angular velocity.  This is calcuated once on object creation.
	iamass float32

	// f applys external forces that adjust the motion each update. The most
	// obvious external forces are gravity and linear and angular damping.
	// An object will continue its current linear and angular motion unless
	// changed directly or the external forces are applied.
	f *forces
}

// newMotion returns a new resting "un-moving" motion.
// Any of the primary motion attributes can be set afterwards.
// Note that location and rotation are injected, not created.
func newMotion(size, mass float32, location *lin.V3, direction *lin.Q) (m *motion) {
	m = &motion{}
	m.mass = mass
	m.imass = 1.0 / m.mass
	m.amass = m.mass * size * size * 1.0 / 6.0
	m.iamass = 1.0 / m.amass

	// set to default resting state.
	m.loc = location
	m.linm = &lin.V3{}
	m.linv = &lin.V3{}
	m.angm = &lin.V3{}
	m.angv = &lin.V3{}
	m.dir = direction
	m.spin = lin.QIdentity()
	m.dir.Unit()
	m.updateSpin()
	m.setLinearMomentum(m.linm)
	m.setAngularMomentum(m.angm)
	m.f = &forces{100, 100, 9810} // 9.81 m/sec*sec
	return
}

// clone returns a duplicate of the calling motion.
func (mo *motion) clone() *motion {
	m := &motion{}
	m.loc = mo.loc.Clone()
	m.linm = mo.linm.Clone()
	m.angm = mo.angm.Clone()
	m.dir = mo.dir.Clone()
	m.linv = mo.linv.Clone()
	m.spin = mo.spin.Clone()
	m.angv = mo.angv.Clone()
	m.mass = mo.mass
	m.imass = mo.imass
	m.amass = mo.amass
	m.iamass = mo.iamass
	m.f = mo.f
	return m
}

// location returns a safe duplicate of the current location coordinates.
func (mo *motion) location() *lin.V3 { return mo.loc.Clone() }

// setLocation updates the bodies location.
func (mo *motion) setLocation(location *lin.V3) {
	mo.loc.X = location.X
	mo.loc.Y = location.Y
	mo.loc.Z = location.Z
}

// linearMomentum returns a copy of the current linear momentum values.
func (mo *motion) linearMomentum() *lin.V3 {
	return &lin.V3{mo.linm.X, mo.linm.Y, mo.linm.Z}
}

// setLinearMomentum updates momentum and velocity.
func (mo *motion) setLinearMomentum(momentum *lin.V3) {
	mo.linm.X = momentum.X
	mo.linm.Y = momentum.Y
	mo.linm.Z = momentum.Z
	mo.linv.X = momentum.X * mo.imass
	mo.linv.Y = momentum.Y * mo.imass
	mo.linv.Z = momentum.Z * mo.imass
}

// linearVelocity returns a copy of the current linear velocity.
// This is a convenience method so that the velocity doesn't have
// to be recalculated by dividing the linear momentum by the mass.
func (mo *motion) linearVelocity() *lin.V3 {
	return &lin.V3{mo.linv.X, mo.linv.Y, mo.linv.Z}
}

// rotation returns a safe duplicate of the current direction values.
func (mo *motion) rotation() *lin.Q { return mo.dir.Clone() }

// setRotation updates the axis of rotation and the amount of rotation..
func (mo *motion) setRotation(direction *lin.Q) {
	mo.dir.X = direction.X
	mo.dir.Y = direction.Y
	mo.dir.Z = direction.Z
	mo.dir.W = direction.W
	mo.dir.Unit()
	mo.updateSpin()
}

// angularMomentum returns a copy of the current angular momentum values.
func (mo *motion) angularMomentum() *lin.V3 {
	return &lin.V3{mo.angm.X, mo.angm.Y, mo.angm.Z}
}

// setAngularMomentum updates angular momentum, angular velocity and spin.
func (mo *motion) setAngularMomentum(momentum *lin.V3) {
	mo.angm.X = momentum.X
	mo.angm.Y = momentum.Y
	mo.angm.Z = momentum.Z
	mo.angv.X = momentum.X * mo.iamass
	mo.angv.Y = momentum.Y * mo.iamass
	mo.angv.Z = momentum.Z * mo.iamass
	mo.updateSpin()
}

// updateSpin adjust the motions spin. The updated spin is
//    spin = angularVelocity*direction*0.5.
// from http://www.cs.cmu.edu/~baraff/pbm/rigid1.pdf (4-2, page G24)
func (mo *motion) updateSpin() {
	if mo.angv != nil && mo.dir != nil {
		mo.spin = &lin.Q{mo.angv.X, mo.angv.Y, mo.angv.Z, 0}
		mo.spin.Scale(0.5).Mult(mo.dir)
		mo.spin.Unit()
	}
}

// derivative stores temporary values used to update primary motion attributes.
// It is used to hold derivative values during the recalulation of motion.
// The RK4 integrator needs to calculate four derivative values each timestep.
type derivative struct {
	// Linear velocity and force relate to linear momentum.
	velocity *lin.V3 // Velocity is the derivative of location.
	force    *lin.V3 // Force in the derivative of linear momentum.

	// Spin and torque relate to angular momentum.
	spin   *lin.Q  // Spin is the derivative of the direction.
	torque *lin.V3 // Torque is the derivative of angular momentum.
}

// isMoving is called to check if the body is moving.  A body is moving if it has any
// linear or angular momentum. This can be used to avoid calling Integrate
// on stationary bodies.
func (mo *motion) isMoving() bool {
	return !(lin.IsZero(mo.linm.X) && lin.IsZero(mo.linm.Y) && lin.IsZero(mo.linm.Z) &&
		lin.IsZero(mo.angm.X) && lin.IsZero(mo.angm.Y) && lin.IsZero(mo.angm.Z))
}

// Integrate motion forward by dt seconds.
// Uses an RK4 integrator to numerically integrate with error O(5).
func (mo *motion) integrate(t, dt float32) {
	d1 := mo.clone().derive(t, 0, nil)
	d2 := mo.clone().derive(t, dt*0.5, d1)
	d3 := mo.clone().derive(t, dt*0.5, d2)
	d4 := mo.clone().derive(t, dt, d3)
	amount := 1.0 / 6.0 * dt

	// combine derivatives to update location using velocity.
	mo.loc.X += amount * (d1.velocity.X + 2*(d2.velocity.X+d3.velocity.X) + d4.velocity.X)
	mo.loc.Y += amount * (d1.velocity.Y + 2*(d2.velocity.Y+d3.velocity.Y) + d4.velocity.Y)
	mo.loc.Z += amount * (d1.velocity.Z + 2*(d2.velocity.Z+d3.velocity.Z) + d4.velocity.Z)
	mo.setLocation(mo.loc)

	// combine derivatives to update linear momentum using force
	mo.linm.X += amount * (d1.force.X + 2*(d2.force.X+d3.force.X) + d4.force.X)
	mo.linm.Y += amount * (d1.force.Y + 2*(d2.force.Y+d3.force.Y) + d4.force.Y)
	mo.linm.Z += amount * (d1.force.Z + 2*(d2.force.Z+d3.force.Z) + d4.force.Z)
	mo.setLinearMomentum(mo.linm)

	// combine derivatives to update axis of rotation using spin.
	mo.dir.X += amount * (d1.spin.X + 2*(d2.spin.X+d3.spin.X) + d4.spin.X)
	mo.dir.Y += amount * (d1.spin.Y + 2*(d2.spin.Y+d3.spin.Y) + d4.spin.Y)
	mo.dir.Z += amount * (d1.spin.Z + 2*(d2.spin.Z+d3.spin.Z) + d4.spin.Z)
	mo.dir.W += amount * (d1.spin.W + 2*(d2.spin.W+d3.spin.W) + d4.spin.W)
	mo.setRotation(mo.dir)

	// combine derivatives to update angular momentum using torque.
	mo.angm.X += amount * (d1.torque.X + 2*(d2.torque.X+d3.torque.X) + d4.torque.X)
	mo.angm.Y += amount * (d1.torque.Y + 2*(d2.torque.Y+d3.torque.Y) + d4.torque.Y)
	mo.angm.Z += amount * (d1.torque.Z + 2*(d2.torque.Z+d3.torque.Z) + d4.torque.Z)
	mo.setAngularMomentum(mo.angm)
}

// derive motion values at a future time (time+deltaTime) using the
// calling motion and passed in derivative to calculate a new derivative.
func (mo *motion) derive(t, dt float32, dp *derivative) *derivative {
	// calculate a new motion if there was a previous derivative
	if dp != nil && dt != 0 {

		// calculate a future location based on a previous derivative.
		mo.loc.X += dp.velocity.X * dt
		mo.loc.Y += dp.velocity.Y * dt
		mo.loc.Z += dp.velocity.Z * dt
		mo.setLocation(mo.loc)

		// calculate a future linear momentum based on a previous derivative.
		mo.linm.X += dp.force.X * dt
		mo.linm.Y += dp.force.Y * dt
		mo.linm.Z += dp.force.Z * dt
		mo.setLinearMomentum(mo.linm)

		// calculate a future axis of rotation based on a previous derivative.
		mo.dir.X += dp.spin.X * dt
		mo.dir.Y += dp.spin.Y * dt
		mo.dir.Z += dp.spin.Z * dt
		mo.dir.W += dp.spin.W * dt
		mo.setRotation(mo.dir)

		// calculate a future angular momentumbased on a previous derivative.
		mo.angm.X += dp.torque.X * dt
		mo.angm.Y += dp.torque.Y * dt
		mo.angm.Z += dp.torque.Z * dt
		mo.setAngularMomentum(mo.angm)
	}

	d := &derivative{}
	d.velocity = &lin.V3{mo.linv.X, mo.linv.Y, mo.linv.Z}
	d.spin = &lin.Q{mo.spin.X, mo.spin.Y, mo.spin.Z, mo.spin.W}
	if mo.f != nil {
		d.force, d.torque = mo.f.apply(t, dt, mo.linv, mo.angv)
	}
	return d
}

// interpolate between two motion states, mo and mend, to create a new
// intermediate motion im. this smooths out movement between ticks where not
// all the motion for a full tick has been applied.
func (mo *motion) interpolate(mend *motion, alpha float32) (im *motion) {
	im = mend.clone()
	im.setLocation(mo.loc.Scale(1 - alpha).Add(mend.loc.Scale(alpha)))
	im.setLinearMomentum(mo.linm.Scale(1 - alpha).Add(mend.linm.Scale(alpha)))
	im.setRotation(mo.dir.Nlerp(mend.dir, alpha))
	im.setAngularMomentum(mo.angm.Scale(1 - alpha).Add(mend.angm.Scale(alpha)))
	return
}

// forces tracks external forces that adjust the motion each update. The most
// obvious external forces are gravity and linear and angular damping.
// An object will continue its current linear and angular motion unless
// changed directly or the external forces are applied.
type forces struct {
	// ldamp damps linear momentum. It does this by opposing the existing
	// movement at a fraction of the existing speed.  It is expected to be a fraction
	// between 0 and 1 that is applied as a negative amount to current velocity.
	// A damping fraction of 1 would stop existing motion, where 0 would have the
	// object continue its current movement forever.
	//
	// Right now all directions are affected equally by damping.
	ldamp float32

	// adamp damps angular momentum.  It behaves the same as
	// ldamp but with respect to angular velocity.
	adamp float32

	// g:gravity drops the Y value (straight down).  The force of gravity on
	// earth is 9.81m/s*s.  Thus the Y value is decreased by this amount each
	// second.
	g float32
}

// damp changes the default damping values.  Use, generally small fractions,
// between 0 (no damping) and 1 (object stops moving immediately).
// The damping values are expected to externally calculated and given as
// the final sum of forces acting on an object.
func (f *forces) damp(linearDamping, angularDamping float32) {
	f.ldamp = linearDamping
	f.adamp = angularDamping
}

// gravity changes the force of gravity in the simulated motion.
func (f *forces) gravity(gravity float32) {
	f.g = gravity
}

// calculate force and torque for motion at time t.  These will be used to
// udpate the linear momentum (using force) and angular momentum (using torque).
// The linear and angular velocity inputs are not changed and two new values,
// force and torque, are created and returned.
func (f *forces) apply(t, dt float32, linv *lin.V3, angv *lin.V3) (force *lin.V3, torque *lin.V3) {
	ldamp := dt * f.ldamp
	adamp := dt * f.adamp
	gravity := dt * f.g
	force = &lin.V3{-linv.X * ldamp, -linv.Y*ldamp - gravity, -linv.Z * ldamp}
	torque = &lin.V3{-angv.X * adamp, -angv.Y * adamp, -angv.Z * adamp}
	return
}
