// Copyright © 2015-2017 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

import (
	"testing"
)

func TestStringHash(t *testing.T) {
	expecting := uint32(3088560132)
	if hash := stringHash("hudbg"); hash != expecting {
		t.Errorf("Expecting %d, got %d", expecting, hash)
	}
	expecting = 2937429552
	if hash := stringHash("cloak"); hash != expecting {
		t.Errorf("Expecting %d, got %d", expecting, hash)
	}
}

func TestStringHashEmpty(t *testing.T) {
	if hash := stringHash(""); hash != 0 {
		t.Errorf("Hash of empty string should be zero, got %d", hash)
	}
}
