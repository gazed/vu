// Copyright © 2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package synth

// Z-order curve utility methods. A z-order number represents a single value
// of a global x,y positioning. The number of bits in the z-order number
// indicate the zoom level. Also known as Morton codes.
//     https://fgiesen.wordpress.com/2009/12/13/decoding-morton-codes
//     http://graphics.stanford.edu/~seander/bithacks.html#InterleaveBMN

// ZMerge two numbers by interleaving their bits and creating a new
// z-order encoded number. The b-bits will precede the a-bits which
// means the b value can be at most a 31 bit value.
func ZMerge(a, b uint32) (m uint64) { return expand(a) | (expand(b) << 1) }

// ZSplit a z-order encoded number by de-interleaving the bits
// into two numbers. The b value represents the higher order bits.
func ZSplit(m uint64) (a, b uint32) { return compact(m), compact(m >> 1) }

// expand prepares a number to be bit interleaved with another number
// by inserting zeros before each bit. Based on:
func expand(n uint32) uint64 {
	a := uint64(n) & 0x00000000ffffffff      // ---- ---- ---- ---- ---- ---- ---- ---- fedc ba98 7654 3210 fedc ba98 7654 3210
	a = (a ^ (a << 16)) & 0x0000ffff0000ffff // ---- ---- ---- ---- fedc ba98 7654 3210 ---- ---- ---- ---- fedc ba98 7654 3210
	a = (a ^ (a << 8)) & 0x00ff00ff00ff00ff  // ---- ---- fedc ba98 ---- ---- 7654 3210 ---- ---- fedc ba98 ---- ---- 7654 3210
	a = (a ^ (a << 4)) & 0x0f0f0f0f0f0f0f0f  // ---- fedc ---- ba98 ---- 7654 ---- 3210 ---- fedc ---- ba98 ---- 7654 ---- 3210
	a = (a ^ (a << 2)) & 0x3333333333333333  // --fe --dc --ba --98 --76 --54 --32 --10 --fe --dc --ba --98 --76 --54 --32 --10
	a = (a ^ (a << 1)) & 0x5555555555555555  // -f-e -d-c -b-a -9-8 -7-6 -5-4 -3-2 -1-0 -f-e -d-c -b-a -9-8 -7-6 -5-4 -3-2 -1-0
	return a
}

// compact reverses ZExpand by discarding every other bit
// and collapsing the remaining 32 bits together.
func compact(n uint64) uint32 {
	a := n & 0x5555555555555555              // -f-e -d-c -b-a -9-8 -7-6 -5-4 -3-2 -1-0 -f-e -d-c -b-a -9-8 -7-6 -5-4 -3-2 -1-0
	a = (a ^ (a >> 1)) & 0x3333333333333333  // --fe --dc --ba --98 --76 --54 --32 --10 --fe --dc --ba --98 --76 --54 --32 --10
	a = (a ^ (a >> 2)) & 0x0f0f0f0f0f0f0f0f  // ---- fedc ---- ba98 ---- 7654 ---- 3210 ---- fedc ---- ba98 ---- 7654 ---- 3210
	a = (a ^ (a >> 4)) & 0x00ff00ff00ff00ff  // ---- ---- fedc ba98 ---- ---- 7654 3210 ---- ---- fedc ba98 ---- ---- 7654 3210
	a = (a ^ (a >> 8)) & 0x0000ffff0000ffff  // ---- ---- ---- ---- fedc ba98 7654 3210 ---- ---- ---- ---- fedc ba98 7654 3210
	a = (a ^ (a >> 16)) & 0x00000000ffffffff // ---- ---- ---- ---- ---- ---- ---- ---- fedc ba98 7654 3210 fedc ba98 7654 3210
	return uint32(a)
}

// ZLabel returns a label for a zorder merge value.
func ZLabel(zoom uint, merge uint64) (key string) {
	mask := uint64(3)
	buff := make([]byte, zoom)
	for z := zoom; z > 0; z-- {
		part := byte('0')
		mask = 3 << ((z - 1) * 2)
		switch merge & mask >> ((z - 1) * 2) {
		case 3:
			part = '3'
		case 2:
			part = '2'
		case 1:
			part = '1'
		case 0:
		}
		buff[zoom-z] = part
	}
	return string(buff)
}
