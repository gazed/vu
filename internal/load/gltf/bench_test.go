package gltf

import "testing"

func BenchmarkOpenBinary(b *testing.B) {
	benchs := []struct {
		name string
	}{
		{"../../assets/models/BoxVertexColors.glb"},
		{"../../assets/models/OrientationTest.glb"},
	}
	for _, bb := range benchs {
		b.Run(bb.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := Open(bb.name)
				if err != nil {
					b.Errorf("Open() error = %v", err)
					return
				}
			}
		})
	}
}
