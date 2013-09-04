package gl

import "testing"

// The test passes if the binding layer can initialize without crashing.
func TestInit(t *testing.T) {
	Init()
}
