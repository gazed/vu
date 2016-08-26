#version 400 core

// shader to update a particle system based on a simple kinematics function
// http://antongerdelan.net/opengl/particles.html
// Consistent sizing based on:
// http://stackoverflow.com/questions/17397724/point-sprites-for-particle-system

layout (location = 0) in vec3 in_v; // initial velocity
layout (location = 1) in float start_time;

uniform float time; // system time in seconds
uniform mat4  mvm;  // model view matrix
uniform mat4  pm;   // projection matrix

// the fragment shader can use this for it's output colour's alpha component 
out float opacity;

void main()
{
	// work out the elapsed time for _this particle_ after its start time
	float t = time - start_time;

	// allow time to loop around so particle emitter keeps going
	t = mod (t, 3.0);
	opacity = 0.0;
	vec3 p = vec3(0, 0, 0);          // emitter location
	vec3 a = vec3 (0.0, -1.0, 0.0);  // gravity

	// standard kinematics equation of motion with velocity and
	// acceleration (gravity)
	p += in_v * t + 0.5 * a * t * t;

	// gradually make particle fade to invisible over 3 seconds
	opacity = 1.0 - (t / 3.0);

    // ensure the point size is adjusted for camera distance.
    float spriteSize = 0.2;
    vec4 eyePos = mvm * vec4(p, 1);
    vec2 screenSize = vec2(800, 600);
    vec4 projVoxel = pm * vec4(spriteSize,spriteSize,eyePos.z,eyePos.w);
    vec2 projSize = screenSize * projVoxel.xy / projVoxel.w;
    gl_PointSize = 0.1 * (projSize.x+projSize.y);
    gl_Position = pm * eyePos; // particle location.
}
