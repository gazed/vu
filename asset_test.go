// Copyright © 2018 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

import (
	"testing"
)

// Test (almost) unique ids from strings.
func TestAid(t *testing.T) {
	a1 := assetID(shd, "coloredInstanced")
	a2 := assetID(shd, "texturedInstanced")
	if a1 == a2 {
		t.Errorf("Need unique asset ids")
	}
}

func TestStringHash(t *testing.T) {
	expecting := uint64(12628427582448915685)
	if hash := stringHash("hudbg"); hash != expecting {
		t.Errorf("Expecting %d, got %d", expecting, hash)
	}
	expecting = 14906461943560666909
	if hash := stringHash("cloak"); hash != expecting {
		t.Errorf("Expecting %d, got %d", expecting, hash)
	}
}

func TestStringHashEmpty(t *testing.T) {
	if hash := stringHash(""); hash != 0 {
		t.Errorf("Hash of empty string should be zero, got %d", hash)
	}
}
