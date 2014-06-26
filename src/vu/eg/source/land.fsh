#version 330

in      vec3      f_nm;   // normal
in      vec2      tuv0;   // texture coordinates
in      vec2      tuv1;   // texture coordinates
in      float     weight; // texture blend weighting 
uniform sampler2D uv;     // 
uniform vec3      ka;     // material ambient value
uniform vec4      l;      // untransformed light position
out     vec4      ffc;    // final fragment colour

vec4 surfaceColour() {
    vec4 tc = vec4(0.0, 0.0, 0.0, 1.0);
    tc += (1.0-weight) * texture(uv, tuv0);
    tc += weight * texture(uv, tuv1);
    return tc;
}

void main() {
   float diffuse = max(0.0, dot(normalize(f_nm), l.xyz));
   vec4 light = vec4(ka, 1.0) * diffuse;   
   vec4 surface = surfaceColour();
   ffc = light * surface;
}
