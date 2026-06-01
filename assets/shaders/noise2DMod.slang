// SPDX-FileCopyrightText : © 2026 Galvanized Logic Inc.
// SPDX-License-Identifier: MIT

module noise2DMod;

// =============================================================================
// Simplex noise 2D based on https://www.shadertoy.com/view/Msf3WH
// Copyright © 2013 Inigo Quilez The MIT License

float2 hash2D(float2 p) { // replace this by something better
	p = float2( dot(p,float2(127.1,311.7)), dot(p,float2(269.5,183.3)) );
	return -1.0 + 2.0*fract(sin(p)*43758.5453123);
}

public float noise2D( in float2 p ) {
	const float K1 = 0.366025404; // (sqrt(3)-1)/2;
	const float K2 = 0.211324865; // (3-sqrt(3))/6;

	float2	i = floor( p + (p.x+p.y)*K1 );
	float2	a = p - i + (i.x+i.y)*K2;
	float m = step(a.y,a.x);
	float2	o = float2(m,1.0-m);
	float2	b = a - o + K2;
	float2	c = a - 1.0 + 2.0*K2;
	float3	h = max( 0.5-float3(dot(a,a), dot(b,b), dot(c,c) ), 0.0 );
	float3	n = h*h*h*h*float3( dot(a,hash2D(i+0.0)), dot(b,hash2D(i+o)), dot(c,hash2D(i+1.0)));
	return dot( n, float3(70.0) );
}

public float fbm2D(float2 uv) {
	float f = 0.0;
	uv *= 5.0;
	float2x2 m = float2x2( 1.6,  1.2, -1.2,  1.6 );
	f  = 0.5000*noise2D( uv ); uv = mul(uv,m);
	f += 0.2500*noise2D( uv ); uv = mul(uv,m);
	f += 0.1250*noise2D( uv ); uv = mul(uv,m);
	f += 0.0625*noise2D( uv ); uv = mul(uv,m);
	return f; // f = 0.5 + 0.5*f;
}
