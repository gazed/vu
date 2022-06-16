// Copyright © 2024 Galvanized Logic Inc.

package render

import (
	"testing"

	"github.com/gazed/vu/math/lin"
)

// ============================================================================
// Benchmarking

// go test -bench=64
// Check the cost of converting 4x4 matrix of float64 to float32
// The last few runs showed around 3ns
func Benchmark64to32(b *testing.B) {
	m32 := &m4{}
	m64 := &lin.M4{11, 12, 13, 14, 21, 22, 23, 24, 31, 32, 33, 34, 41, 42, 43, 44}
	for cnt := 0; cnt < b.N; cnt++ {
		m32.set64(m64)
	}
}
