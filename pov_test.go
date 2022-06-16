// Copyright © 2017-2024 Galvanized Logic Inc.

package vu

import (
	"fmt"
	"testing"

	"github.com/gazed/vu/math/lin"
)

func TestDeleteRoot(t *testing.T) {
	ents := &entities{}
	povs := newPovs()
	p0 := povs.create(ents.create(), 0) // root  eid 1
	povs.create(ents.create(), p0.eid)  // child eid 2
	povs.create(ents.create(), p0.eid)  // child eid 3
	kids := povs.dispose(p0.eid, []eID{})
	if len(kids) != 2 || kids[0] != 2 || kids[1] != 3 {
		t.Errorf("%d %d %d", len(kids), kids[0], kids[1])
	}
}

func TestDeleteChild(t *testing.T) {
	ents := &entities{}
	povs := newPovs()
	p0 := povs.create(ents.create(), 0)      // root  eid 1
	k0 := povs.create(ents.create(), p0.eid) // child eid 2
	povs.create(ents.create(), p0.eid)       // child eid 3
	kids := povs.dispose(k0.eid, []eID{})
	if len(kids) != 0 || povs.eids[1] != 3 {
		t.Errorf("0:%d 3:%d", len(kids), povs.eids[1])
	}
	if povs.nodes[0].kids[0] != 3 || len(povs.nodes[0].kids) != 1 {
		t.Errorf("3:%d 1:%d", povs.nodes[0].kids[0], len(povs.nodes[0].kids))
	}
}

func TestDeleteMiddle(t *testing.T) {
	ents := &entities{}
	povs := newPovs()
	e1 := povs.create(ents.create(), 0)      // root0 eid 1          1
	e2 := povs.create(ents.create(), e1.eid) // child eid 2         / \
	povs.create(ents.create(), e1.eid)       // child eid 3        2   3
	e4 := povs.create(ents.create(), e2.eid) // root  eid 4       / \
	povs.create(ents.create(), e2.eid)       // child eid 5      4   5
	povs.create(ents.create(), e4.eid)       // child eid 6     / \
	povs.create(ents.create(), e4.eid)       // root1 eid 7    6   7
	kids := povs.dispose(e2.eid, []eID{})
	if len(kids) != 4 || povs.nodes[0].kids[0] != 3 {
		t.Errorf("0:%d 3:%d", len(kids), povs.nodes[0].kids[0])
	}
}

// Dump a matrix. Used to debug the pov transform methods.
func DumpM4(m *lin.M4) string {
	format := "[%+2.9f, %+2.9f, %+2.9f, %+2.9f]\n"
	str := fmt.Sprintf(format, m.Xx, m.Xy, m.Xz, m.Xw)
	str += fmt.Sprintf(format, m.Yx, m.Yy, m.Yz, m.Yw)
	str += fmt.Sprintf(format, m.Zx, m.Zy, m.Zz, m.Zw)
	str += fmt.Sprintf(format, m.Wx, m.Wy, m.Wz, m.Ww)
	return str
}
