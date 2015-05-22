#version 330

layout(location=0) in vec3 in_v; // particle locations.
uniform               mat4  mvm; // model view matrix
uniform               mat4  pm;  // projection matrix
 
// A particle effect shader. Consistent sizing based on:
// http://stackoverflow.com/questions/17397724/point-sprites-for-particle-system
void main(void)
{
   float spriteSize = 0.2;
   vec4 eyePos = mvm * vec4(in_v, 1);
   vec2 screenSize = vec2(800, 600);
   vec4 projVoxel = pm * vec4(spriteSize,spriteSize,eyePos.z,eyePos.w);
   vec2 projSize = screenSize * projVoxel.xy / projVoxel.w;
   gl_PointSize = 0.1 * (projSize.x+projSize.y);
   gl_Position = pm * eyePos; // particle location.
}
