// SPDX-FileCopyrightText : © 2026 Galvanized Logic Inc.
// SPDX-License-Identifier: MIT

module modNoise3D;

// =============================================================================
// Noise - value - 3D based on https://www.shadertoy.com/view/4sfGzS
// Copyright © 2013 Inigo Quilez The MIT License

float hash3D(float3 p) {  // replace this by something better
	p  = fract( p*0.3183099+.1 );
	p *= 17.0;
	return frac( p.x*p.y*p.z*(p.x+p.y+p.z) );
}

public float noise3D(float3 x) {
	float3 i = floor(x);
	float3 f = fract(x);
	f = f*f*(3.0-2.0*f);
	return lerp(lerp(lerp(hash3D(i+float3(0,0,0)),
						  hash3D(i+float3(1,0,0)),f.x),
					 lerp(hash3D(i+float3(0,1,0)),
						  hash3D(i+float3(1,1,0)),f.x),f.y),
				lerp(lerp(hash3D(i+float3(0,0,1)),
						  hash3D(i+float3(1,0,1)),f.x),
					 lerp(hash3D(i+float3(0,1,1)),
						  hash3D(i+float3(1,1,1)),f.x),f.y),f.z);
}

static const float3x3 m = float3x3( 0.00,  0.80,  0.60,
								   -0.80,  0.36, -0.48,
								   -0.60, -0.48,  0.64 );

// fbm3D layers noise using Fractal Brownian Motion.
// Use 4 octaves and rotate the results each octave.
public float fbm3D(float3 q, float freq) {
	float amp = 0.5; // amplitude per octave.
	float f = 0.0;	 // final noise value.
	f  = amp*noise3D( q ); q = mul(q, m)*(freq+0.1); amp *= 0.5;
	f += amp*noise3D( q ); q = mul(q, m)*(freq+0.2); amp *= 0.5;
	f += amp*noise3D( q ); q = mul(q, m)*(freq+0.3); amp *= 0.5;
	f += amp*noise3D( q ); q = mul(q, m)*(freq+0.1); amp *= 0.5;
	return f;
}
