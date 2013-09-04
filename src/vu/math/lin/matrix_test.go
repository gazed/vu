// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package lin

import (
	"fmt"
	"strconv"
	"testing"
)

// Check that golang lays out the data structure as sequential float32's.
func TestMemoryLayout(t *testing.T) {
	m4 := M4{
		11, 12, 13, 14,
		21, 22, 23, 24,
		31, 32, 33, 34,
		41, 42, 43, 44}
	mema, err := strconv.ParseUint(fmt.Sprintf("%d", &(m4.X0)), 0, 64)
	memb, _ := strconv.ParseUint(fmt.Sprintf("%d", &(m4.Y0)), 0, 64)
	if err != nil {
		fmt.Printf("error %s\n", err)
	}
	if memb-mema != 4 {
		t.Error("float32 should be 4 bytes")
	}
	memc, _ := strconv.ParseUint(fmt.Sprintf("%d", &(m4.X1)), 0, 64)
	if memc-mema != 16 {
		t.Error("vector xyzw:0 should be followed by vector xyzw:1")
	}
}

func TestCloneM3(t *testing.T) {
	m3 := M3{
		1, 2, 3,
		4, 5, 6,
		7, 8, 9}
	m := m3.Clone()
	expect := M3{
		1, 2, 3,
		4, 5, 6,
		7, 8, 9}
	got, want := m.Dump(), expect.Dump()
	if got != want {
		t.Errorf(format, got, want)
	}
}
func TestCloneM4(t *testing.T) {
	m4 := M4{
		11, 12, 13, 14,
		21, 22, 23, 24,
		31, 32, 33, 34,
		41, 42, 43, 44}
	m := m4.Clone()
	expect := M4{
		11, 12, 13, 14,
		21, 22, 23, 24,
		31, 32, 33, 34,
		41, 42, 43, 44}
	got, want := m.Dump(), expect.Dump()
	if got != want {
		t.Errorf(format, got, want)
	}
}

func TestM3(t *testing.T) {
	m4 := M4{
		11, 12, 13, 14,
		21, 22, 23, 24,
		31, 32, 33, 34,
		41, 42, 43, 44}
	m3 := m4.M3()
	expect := M3{
		11, 12, 13,
		21, 22, 23,
		31, 32, 33}
	got, want := m3.Dump(), expect.Dump()
	if got != want {
		t.Errorf(format, got, want)
	}
}

func TestIdentity3x3(t *testing.T) {
	m3 := M3Identity()
	expect := M3{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1}
	got, want := m3.Dump(), expect.Dump()
	if got != want {
		t.Errorf(format, got, want)
	}
}
func TestIdentity4x4(t *testing.T) {
	m4 := M4Identity()
	expect := M4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1}
	got, want := m4.Dump(), expect.Dump()
	if got != want {
		t.Errorf(format, got, want)
	}
}

func TestTranspose3x3(t *testing.T) {
	m3 := M3{
		1, 2, 3,
		4, 5, 6,
		7, 8, 9}
	m3.Transpose()
	expect := M3{
		1, 4, 7,
		2, 5, 8,
		3, 6, 9}
	got, want := m3.Dump(), expect.Dump()
	if got != want {
		t.Errorf(format, got, want)
	}
}
func TestTranspose4x4(t *testing.T) {
	m4 := M4{
		11, 12, 13, 14,
		21, 22, 23, 24,
		31, 32, 33, 34,
		41, 42, 43, 44}
	m4.Transpose()
	expect := M4{
		11, 21, 31, 41,
		12, 22, 32, 42,
		13, 23, 33, 43,
		14, 24, 34, 44}
	got, want := m4.Dump(), expect.Dump()
	if got != want {
		t.Errorf(format, got, want)
	}
}

func TestMultiply3x3(t *testing.T) {
	lm := &M3{
		1, 1, 1,
		1, 1, 1,
		1, 1, 1}
	rm := &M3{
		1, 2, 3,
		1, 2, 3,
		1, 2, 3}
	lm.Mult(rm)
	expect := M3{
		3, 6, 9,
		3, 6, 9,
		3, 6, 9}
	got, want := lm.Dump(), expect.Dump()
	if got != want {
		t.Errorf(format, got, want)
	}
}
func TestMultiply4x4(t *testing.T) {
	lm := &M4{
		1, 1, 1, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
		1, 1, 1, 1}
	rm := &M4{
		1, 2, 3, 4,
		1, 2, 3, 4,
		1, 2, 3, 4,
		1, 2, 3, 4}
	lm.Mult(rm)
	expect := M4{
		4, 8, 12, 16,
		4, 8, 12, 16,
		4, 8, 12, 16,
		4, 8, 12, 16}
	got, want := lm.Dump(), expect.Dump()
	if got != want {
		t.Errorf(format, got, want)
	}
}

func TestTranslate(t *testing.T) {
	m4 := &M4{
		1, 2, 3, 4,
		1, 2, 3, 4,
		1, 2, 3, 4,
		1, 2, 3, 4}
	m4.TranslateL(1, 2, 3)
	expect := M4{
		1, 2, 3, 4,
		1, 2, 3, 4,
		1, 2, 3, 4,
		7, 14, 21, 28}
	got, want := m4.Dump(), expect.Dump()
	if got != want {
		t.Errorf(format, got, want)
	}

	// do the same transform using matrix multiply.
	m4 = M4Translater(1, 2, 3)
	mw := &M4{
		1, 2, 3, 4,
		1, 2, 3, 4,
		1, 2, 3, 4,
		1, 2, 3, 4}
	m4.Mult(mw)
	got, want = m4.Dump(), expect.Dump()
	if got != want {
		t.Errorf(format, got, want)
	}

	// try the right side multiply translate.
	mw = &M4{
		1, 2, 3, 4,
		1, 2, 3, 4,
		1, 2, 3, 4,
		1, 2, 3, 4}
	mw.TranslateR(1, 2, 3)
	expect = M4{
		5, 10, 15, 4,
		5, 10, 15, 4,
		5, 10, 15, 4,
		5, 10, 15, 4}
	got, want = mw.Dump(), expect.Dump()
	if got != want {
		t.Errorf(format, got, want)
	}
}

func TestScale(t *testing.T) {
	m4 := &M4{
		1, 2, 3, 4,
		1, 2, 3, 4,
		1, 2, 3, 4,
		1, 2, 3, 4}
	m4.ScaleL(1, 2, 3)
	expect := M4{
		1, 2, 3, 4,
		2, 4, 6, 8,
		3, 6, 9, 12,
		1, 2, 3, 4}
	got, want := m4.Dump(), expect.Dump()
	if got != want {
		t.Errorf(format, got, want)
	}

	// do the same transform using matrix multiply.
	m4 = M4Scaler(1, 2, 3)
	mw := &M4{
		1, 2, 3, 4,
		1, 2, 3, 4,
		1, 2, 3, 4,
		1, 2, 3, 4}
	m4.Mult(mw)
	got, want = m4.Dump(), expect.Dump()
	if got != want {
		t.Errorf(format, got, want)
	}

	// try the right side multiply scale.
	m4 = &M4{
		1, 2, 3, 4,
		1, 2, 3, 4,
		1, 2, 3, 4,
		1, 2, 3, 4}
	m4.ScaleR(1, 2, 3)
	expect = M4{
		1, 4, 9, 4,
		1, 4, 9, 4,
		1, 4, 9, 4,
		1, 4, 9, 4}
	got, want = m4.Dump(), expect.Dump()
	if got != want {
		t.Errorf(format, got, want)
	}
}

func TestOrthographicProjection(t *testing.T) {
	m4 := M4Orthographic(2, 3, 4, 5, 6, 7)
	expect := M4{
		2, 0, +0, 0,
		0, 2, +0, 0,
		0, 0, -2, 0,
		-5, -9, -13, 1}
	got, want := m4.Dump(), expect.Dump()
	if got != want {
		t.Errorf(format, got, want)
	}
}

func TestPerspective(t *testing.T) {
	m4 := M4Perspective(1, 2, 3, 4)
	expect := M4{
		57.29, 0.0000, 0, 0,
		0.000, 114.59, 0, 0,
		0.000, 0.0000, -7, -1,
		0.000, 0.0000, -24, 0}
	got, want := m4.Dump(), expect.Dump()
	if got != want {
		t.Errorf(format, got, want)
	}
}

func TestIPerspective(t *testing.T) {
	p := M4Perspective(45, 800.0/600.0, 0.1, 50)
	ip := M4PerspectiveI(45, 800.0/600.0, 0.1, 50)
	p.Mult(ip)
	expect := M4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1}
	got, want := p.Dump(), expect.Dump()
	if got != want {
		t.Errorf(format, got, want)
	}
}

func TestIModelView(t *testing.T) {
	dir := &Q{0, 0, 0, 1}
	dir = QAxisAngle(&V3{0, 1, 0}, 45).Mult(dir)
	R := QAxisAngle(&V3{0, 0, 1}, 45).Mult(dir).M4()
	T := M4Translater(10, 0, -5)
	mv := R.Mult(T)
	imv := mv.Clone().IModelView()
	m := imv.Mult(mv)
	if IsZero(m.X0) {
		m.X0 = 0
	}
	if IsZero(m.Y0) {
		m.Y0 = 0
	}
	if IsZero(m.Z0) {
		m.Z0 = 0
	}
	if IsZero(m.X1) {
		m.X1 = 0
	}
	if IsZero(m.Y1) {
		m.Y1 = 0
	}
	if IsZero(m.Z1) {
		m.Z1 = 0
	}
	if IsZero(m.X2) {
		m.X2 = 0
	}
	if IsZero(m.Y2) {
		m.Y2 = 0
	}
	if IsZero(m.Z2) {
		m.Z2 = 0
	}
	if IsZero(m.X3) {
		m.X3 = 0
	}
	if IsZero(m.Y3) {
		m.Y3 = 0
	}
	if IsZero(m.Z3) {
		m.Z3 = 0
	}
	expect := M4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1}
	got, want := m.Dump(), expect.Dump()
	if got != want {
		t.Errorf(format, got, want)
	}
}

// Check that the standard scale-rotate-translate multiplication is quicker if
// new matricies are not created each time. Note that the S, and R matricies are
// altered each time through the loop of the "Same" test.
// Run 'go test -bench=".*"' to get something like:
//    BenchmarkMultSame	20000000	       109 ns/op
//    BenchmarkMultNew	10000000	       278 ns/op
func BenchmarkMultSame(b *testing.B) {
	S, R, T := M4Scaler(0.5, 0, 0), QAxisAngle(&V3{Y: 1}, 90).M4(), M4Translater(0, 0, 4)
	for cnt := 0; cnt < b.N; cnt++ {
		S.Mult(R.Mult(T))
	}
}
func BenchmarkMultNew(b *testing.B) {
	S, R, T := M4Scaler(0.5, 0, 0), QAxisAngle(&V3{Y: 1}, 90).M4(), M4Translater(0, 0, 4)
	for cnt := 0; cnt < b.N; cnt++ {
		S.MultNew(R.MultNew(T))
	}
}

// MultNew does the same as M4.Mult, but returns a new matrix each time.
func (m *M4) MultNew(rm *M4) *M4 {
	x0 := m.X0*rm.X0 + m.Y0*rm.X1 + m.Z0*rm.X2 + m.W0*rm.X3
	y0 := m.X0*rm.Y0 + m.Y0*rm.Y1 + m.Z0*rm.Y2 + m.W0*rm.Y3
	z0 := m.X0*rm.Z0 + m.Y0*rm.Z1 + m.Z0*rm.Z2 + m.W0*rm.Z3
	w0 := m.X0*rm.W0 + m.Y0*rm.W1 + m.Z0*rm.W2 + m.W0*rm.W3
	x1 := m.X1*rm.X0 + m.Y1*rm.X1 + m.Z1*rm.X2 + m.W1*rm.X3
	y1 := m.X1*rm.Y0 + m.Y1*rm.Y1 + m.Z1*rm.Y2 + m.W1*rm.Y3
	z1 := m.X1*rm.Z0 + m.Y1*rm.Z1 + m.Z1*rm.Z2 + m.W1*rm.Z3
	w1 := m.X1*rm.W0 + m.Y1*rm.W1 + m.Z1*rm.W2 + m.W1*rm.W3
	x2 := m.X2*rm.X0 + m.Y2*rm.X1 + m.Z2*rm.X2 + m.W2*rm.X3
	y2 := m.X2*rm.Y0 + m.Y2*rm.Y1 + m.Z2*rm.Y2 + m.W2*rm.Y3
	z2 := m.X2*rm.Z0 + m.Y2*rm.Z1 + m.Z2*rm.Z2 + m.W2*rm.Z3
	w2 := m.X2*rm.W0 + m.Y2*rm.W1 + m.Z2*rm.W2 + m.W2*rm.W3
	x3 := m.X3*rm.X0 + m.Y3*rm.X1 + m.Z3*rm.X2 + m.W3*rm.X3
	y3 := m.X3*rm.Y0 + m.Y3*rm.Y1 + m.Z3*rm.Y2 + m.W3*rm.Y3
	z3 := m.X3*rm.Z0 + m.Y3*rm.Z1 + m.Z3*rm.Z2 + m.W3*rm.Z3
	w3 := m.X3*rm.W0 + m.Y3*rm.W1 + m.Z3*rm.W2 + m.W3*rm.W3
	return &M4{
		x0, y0, z0, w0,
		x1, y1, z1, w1,
		x2, y2, z2, w2,
		x3, y3, z3, w3,
	}
}

// Check the time it takes to create an identity matrix with and
// without a function call is about the same.
// Run 'go test -bench=".*"' to get something like:
//    BenchmarkCreateCall	20000000   95.4 ns/op
//    BenchmarkCreateNoCall	20000000   95.4 ns/op
func BenchmarkCreateCall(b *testing.B) {
	var m *M4
	for cnt := 0; cnt < b.N; cnt++ {
		m = M4Identity()
	}
	m.X0 = 0 // make the compiler happy.
}
func BenchmarkCreateNoCall(b *testing.B) {
	var m *M4
	for cnt := 0; cnt < b.N; cnt++ {
		m = &M4{X0: 1, Y1: 1, Z2: 1, W3: 1}
	}
	m.X0 = 0 // make the compiler happy.
}
