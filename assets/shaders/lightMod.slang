// SPDX-FileCopyrightText : © 2026 Galvanized Logic Inc.
// SPDX-License-Identifier: MIT

module lightMod;

// Light is uniform data.
// Light is either a sun light, point light, or spot light.
// - a sun light uses: color, intensity,		 , direction.
// - point lights use: color, intensity, position,			, attenuation.
// - spot lights use : color, intensity, position, direction, attenuation, cutoff.
public struct Light {
	public float4 color;       // RGB are 0-1, A:unused
	public float4 pos;	       // XYZ is the world space light position, W:unused
	public float4 dir;	       // XYZ is world space light direction, W:unused
	public float  intensity;   //
    public float  attenuation; //
    public float  cutoff;      // cos(cutoff angle in radians)
	public int    type;        // 0:none 1:sun 2:point 3:spot
};

// Light.type matches the vu light types.
public static const int NoLight    = 0;
public static const int SunLight   = 1;
public static const int PointLight = 2;
public static const int SpotLight  = 3;

// LightInfo holds calculations based on light type.
public struct LightInfo {
	public float4 Color;  // fragment rgba
	public float3 World;  // vertex world position.
	public float3 L;      // normalized light direction
	public float3 I;      // light color * intensity
};

// common fragment shader code.
public bool updateLightInfo(Light l, inout LightInfo li) {
	float3 ldir;
	float3 ldist;
	switch (l.type) {
	case SunLight:
		li.L = normalize(l.pos.xyz);         // sun light direction
		li.I = l.color.xyz * l.intensity;    // light color * light intensity.
		return true;
	case PointLight:
		ldir = l.pos.xyz - li.World;         // light direction.
		ldist = length(ldir);                // get distance before normalization.
		li.L = normalize(ldir);		         // normalize light direction.
		li.I = l.color.xyz * l.intensity;    // intensity
		li.I /= (ldist * ldist);             // intensity drops off by distance squared.
		return true;
	case SpotLight:
		ldir = l.pos.xyz - li.World;         // light direction.
		ldist = length(ldir);                // get distance before normalization.
		li.L = normalize(ldir);		         // normalize light direction.
		float beamEdge = dot(li.L, normalize(-l.dir.xyz));
		li.I = float3(0); // outside spot light beam
		if (beamEdge > l.cutoff) {
			// inside spot light beam
			li.I = l.color.xyz * l.intensity; // intensity
			li.I /= (ldist * ldist);          // intensity drops off by distance squared.

			// add blur to beam edge.
			float falloff =  1.0 - ((1.0 - beamEdge)/(1.0 - l.cutoff));
			li.I = li.I*falloff;
		}
		return true;
	default:
		return false; // NoLight
	}
}
