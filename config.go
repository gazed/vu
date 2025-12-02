// SPDX-FileCopyrightText : © 2022 Galvanized Logic Inc.
// SPDX-License-Identifier: BSD-2-Clause

package vu

// config.go reduces the NewEngine API footprint using functional options.
// See: http://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
//      https://commandcenter.blogspot.ca/2014/01/self-referential-functions-and-design.html

// Config contains configuration attributes that can be set by the game
// before running the engine game loop.
type Config struct {

	// attributes for windowed games
	title    string // window title
	windowed bool   // true to run in windowed mode.
	x, y     int32  // display top left corner in pixels
	w, h     int32  // display width and height in pixels

	// display default background color
	r, g, b, a float32 // red, green, blue, alpha: range 0-1
}

// configDefaults provides reasonable defaults so the game
// runs even if no configuration attributes are set.
var configDefaults = Config{
	title:    "VU",  // default title
	windowed: false, // default full screen.
	x:        0,     // top left corner
	y:        0,     // top left corner
	w:        800,   // default 16:9 ratio
	h:        450,   // default 16:9 ratio
	r:        0.0,   // default black
	g:        0.0,   // default black
	b:        0.0,   // default black
	a:        1.0,   // default opaque
}

// Attr defines optional application attributes that can be used to
// configure the engine.
//
//	eng, err := vu.NewEngine(
//	   vu.Title("Keyboard Controller"),
//	   vu.Size(200, 200, 900, 400),
//	   vu.Background(0.45, 0.45, 0.45, 1.0),
//	)
type Attr func(*Config) // type for attribute overrides

// Title sets the window title when using windowed mode.
// For use in NewEngine().
func Title(t string) Attr {
	return func(c *Config) { c.title = t }
}

// Size sets the window top left corner location
// and size in pixels when using windowed mode.
func Size(x, y, w, h int32) Attr {
	// FUTURE: revisit the upper bounds.
	return func(c *Config) {
		// limit to reasonable locations.
		if x >= 0 && x < 10_000 {
			c.x = x
		}
		if y >= 0 && y < 10_000 {
			c.y = y
		}

		// limit to resonable sizes.
		if w > 10 && w < 10_000 {
			c.w = w
		}
		if h > 10 && h < 10_000 {
			c.h = h
		}
	}
}

// Windowed mode instead of fullscreen.
func Windowed() Attr {
	return func(c *Config) { c.windowed = true }
}

// Background display clear color.
func Background(r, g, b, a float32) Attr {
	return func(c *Config) { c.r = r; c.g = g; c.b = b; c.a = a }
}
