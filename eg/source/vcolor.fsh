// vcolor fragment shader uses the color passed from the vertex shader.

in  vec4 v_color; // Interpoloated color from vertex shader
out vec4 f_color; // Fragment color.

void main(void) {
   f_color = v_color;
}
