// Copyright © 2024 Galvanized Logic Inc.

package render

import (
	"testing"
)

// check reusing slices of Packets
func TestPackets(t *testing.T) {
	pass := Pass{}
	var p *Packet

	// check that the initial number of packets is zero
	t.Run("check initial number of packets", func(t *testing.T) {
		c, l := cap(pass.Packets), len(pass.Packets)
		if c != 0 {
			t.Fatal("expected zero capacity got", c)
		}
		if l != 0 {
			t.Fatal("expected zero length got", l)
		}
	})

	t.Run("allocate 1 packet", func(t *testing.T) {
		if pass.Packets, p = pass.Packets.GetPacket(); p == nil {
			t.Fatal("expected a packet")
		}
		c, l := cap(pass.Packets), len(pass.Packets)
		if c != 1 {
			t.Fatal("expected 1 capacity got", c)
		}
		if l != 1 {
			t.Fatal("expected 1 length got", l)
		}
	})

	t.Run("allocate 100 more packets", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			if pass.Packets, p = pass.Packets.GetPacket(); p == nil {
				t.Fatal("expected a packet")
			}
		}
		c, l := cap(pass.Packets), len(pass.Packets)
		if c < 101 {
			t.Fatal("expected 101 capacity or greater got", c)
		}
		if l != 101 {
			t.Fatal("expected 1 length got", l)
		}
	})

	t.Run("reset", func(t *testing.T) {
		pass.Packets = pass.Packets[:0] // reset keeping underlying memory.
		c, l := cap(pass.Packets), len(pass.Packets)
		if c < 101 {
			t.Fatal("expected 101 capacity or greater got", c)
		}
		if l != 0 {
			t.Fatal("expected 0 length got", l)
		}
	})

	t.Run("allocate 10 more packets", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			if pass.Packets, p = pass.Packets.GetPacket(); p == nil {
				t.Fatal("expected a packet")
			}
		}
		c, l := cap(pass.Packets), len(pass.Packets)
		if c < 101 {
			t.Fatal("expected 101 capacity or greater got", c)
		}
		if l != 10 {
			t.Fatal("expected 10 length got", l)
		}
	})
}
