// Copyright © 2014-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// Package ai provides support for application unit behaviour.
// This is an experimental package that currently provides a behaviour
// tree implementation.
//
// Package ai is provided as part of the vu (virtual universe) 3D engine.
package ai

// Design Notes:
// Other AI support may be provided if:
//    • It is generic, ie. most of AI seems to be application
//      specific so it is unclear how much an engine can help.
//    • It is usefull, ie. simple AI's can be done with if-else statements,
//      so no engine help is necessary.
// Possible areas to investigate are:
//    • Helping smooth unit movement (assuming unit has a path from
//      vu/grid path or flow).

// More information at:
// https://web.cs.ship.edu/~djmoon/gaming/gaming-notes/ai-movement.pdf
// http://www.raywenderlich.com/24824/introduction-to-ai-programming-for-games
// http://www.hobbygamedev.com/articles/vol8/real-time-videogame-ai/
// http://www.ramalila.net/Adventures/AI/RealTime.html
