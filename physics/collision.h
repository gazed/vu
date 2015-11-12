// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// The basic types that are needed by box-box collision. 
// All kinds of nice allocation, memory alignment, and C++ class functionality 
// are lost from the original bullet physics types, but it gets things working.
typedef double btScalar;
typedef btScalar dMatrix3[4*3];
typedef btScalar btVector3[4];

// Consolidate the input box-box information into a single structure. 
typedef struct {
	btVector3 orgA, orgB; // Origin of boxes in world space. 
	dMatrix3  rotA, rotB; // 3x3 rotation transforms for boxes.
	btVector3 lenA, lenB; // Half-lengths of boxes.
} BoxBoxInput;

// Consolidate the output box-box collision information into a structure. 
typedef struct { 
	btVector3 n; // Normal of collision.
	btVector3 p; // Point of contact of collision.
	btScalar  d; // Depth of collision.
} BoxBoxContact; // One contact.
typedef struct {
	int           code;   // Collision face/edge indicator.
	int           ncp;    // Number of contact points.
	BoxBoxContact bbc[4]; // Points of contact.
} BoxBoxResults;          // All contacts (up to 4).

// Collide two boxes and generate points of contact. The number of
// contacts will be zero if the boxes did not actually collide.
void boxBoxClosestPoints(BoxBoxInput *in, BoxBoxResults *out);
