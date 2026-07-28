// SPDX-FileCopyrightText : © 2026 Galvanized Logic Inc.
// SPDX-License-Identifier: MIT

package grid

import (
	"testing"
)

// go test -run Offset
func TestOffset(t *testing.T) {

	// should look similar to grid_test:"dump grid" with an extra border.
	// Note that the offset grid is drawn using OffsetSymbols.
	t.Run("dump offset grid", func(t *testing.T) {
		g := Grid{
			{1, 1, 0, 1, 1},
			{1, 0, 0, 0, 1},
			{0, 0, 1, 0, 0},
			{1, 0, 0, 0, 1},
			{1, 1, 0, 1, 1},
		}
		ogrid := MakeOffsetGrid(g, true)
		ogrid.Dump(OffsetSymbols)
	})
}
