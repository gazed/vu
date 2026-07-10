// SPDX-FileCopyrightText : © 2026 Galvanized Logic Inc.
// SPDX-License-Identifier: MIT

module rayMod;

// ray queries are used for hard shadows.

// acceleration structure for ray queries.
layout(set=0, binding=1) public RaytracingAccelerationStructure tlas;

// traceShadowOccluded returns true if the ray hits anything.
public bool traceShadowOccluded(float3 origin, float3 direction, float3 normal, float tMin, float tMax) {
    if (dot(normal, direction) <= 0) {
        return true; // skip fragments that are not facing the light
    }
    RayDesc ray;
    ray.Origin = origin;       // start at fragment - in world space.
    ray.Direction = direction; // direction to the light - vector in world space.
    ray.TMin = tMin;           // starting offset so as to not hit origin surface.
    ray.TMax = tMax;           // maximum distance to light.

    RayQuery q;
    // OPTIMIZATION: Accept first hit and end search is sufficient for hard shadows
    q.TraceRayInline(tlas, RAY_FLAG_ACCEPT_FIRST_HIT_AND_END_SEARCH, 0xFF, ray);
    while (q.Proceed()) {
        // no additional logic needed for hard shadows.
		// FUTURE: add more queries for soft shadows.
    }
    return (q.CommittedStatus() == COMMITTED_TRIANGLE_HIT);
}
