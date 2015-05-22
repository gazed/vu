#version 330

layout(location=0) in vec3  in_v;  // particle point.
layout(location=1) in vec2  in_d;  // particle data: x is index, y is lifespan 1->0.
uniform               mat4  mvm;   // model view matrix
uniform               mat4  pm;    // projection matrix
out                   float alpha;
out                   float index;
 
// A particle effect shader. Consistent sizing based on:
// http://stackoverflow.com/questions/17397724/point-sprites-for-particle-system
void main(void)
{
   float spriteSize = 2;
   vec4 eyePos = mvm * vec4(in_v, 1);
   vec2 screenSize = vec2(800, 600);
   vec4 projVoxel = pm * vec4(spriteSize,spriteSize,eyePos.z,eyePos.w);
   vec2 projSize = screenSize * projVoxel.xy / projVoxel.w;
   gl_PointSize = 0.1 * (projSize.x+projSize.y) * in_d.y;
   alpha = in_d.y;            // particle transparency.
   index = in_d.x;            // particle identifier.
   gl_Position = pm * eyePos; // particle location.
}
