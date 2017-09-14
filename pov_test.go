// Copyright © 2017 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package vu

import (
	"testing"
)

func TestDeleteRoot(t *testing.T) {
	eids := &eids{}
	povs := newPovs()
	p0 := povs.create(eids.create(), 0) // root  eid 1
	povs.create(eids.create(), p0.eid)  // child eid 2
	povs.create(eids.create(), p0.eid)  // child eid 3
	kids := povs.dispose(p0.eid, []eid{})
	if len(kids) != 2 || kids[0] != 2 || kids[1] != 3 {
		t.Errorf("%d %d %d", len(kids), kids[0], kids[1])
	}
}

func TestDeleteChild(t *testing.T) {
	eids := &eids{}
	povs := newPovs()
	p0 := povs.create(eids.create(), 0)      // root  eid 1
	k0 := povs.create(eids.create(), p0.eid) // child eid 2
	povs.create(eids.create(), p0.eid)       // child eid 3
	kids := povs.dispose(k0.eid, []eid{})
	if len(kids) != 0 || povs.eids[1] != 3 {
		t.Errorf("0:%d 3:%d", len(kids), povs.eids[1])
	}
	if povs.nodes[0].kids[0] != 3 || len(povs.nodes[0].kids) != 1 {
		t.Errorf("3:%d 1:%d", povs.nodes[0].kids[0], len(povs.nodes[0].kids))
	}
}

func TestDeleteMiddle(t *testing.T) {
	eids := &eids{}
	povs := newPovs()
	e1 := povs.create(eids.create(), 0)      // root0 eid 1          1
	e2 := povs.create(eids.create(), e1.eid) // child eid 2         / \
	povs.create(eids.create(), e1.eid)       // child eid 3        2   3
	e4 := povs.create(eids.create(), e2.eid) // root  eid 4       / \
	povs.create(eids.create(), e2.eid)       // child eid 5      4   5
	povs.create(eids.create(), e4.eid)       // child eid 6     / \
	povs.create(eids.create(), e4.eid)       // root1 eid 7    6   7
	kids := povs.dispose(e2.eid, []eid{})
	if len(kids) != 4 || povs.nodes[0].kids[0] != 3 {
		t.Errorf("0:%d 3:%d", len(kids), povs.nodes[0].kids[0])
	}
}
