// fragment shader nshade uses the color passed from the vertex shader.

in  vec3 v_color; // color from vertex shader
out vec4 f_color; // fragment color.

void main(void) {
   f_color = vec4(v_color, 1.0);
}
