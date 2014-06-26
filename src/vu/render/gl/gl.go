// Package gl provides a golang 3D graphics library
// that was auto-generated from open gl spec src/vu/render/gl/gen/glcorearb.h-v4.3.
// The official OpenGL documentation for any of the constants
// or methods can be found online. Just prepend "GL_"
// to the function or constants names in this package.
//
// Package gen is provided as part of the vu (virtual universe) 3D engine.
package gl

// #cgo darwin  LDFLAGS: -framework OpenGL
// #cgo linux   LDFLAGS: -lGL -ldl
// #cgo windows LDFLAGS: -lopengl32
//
// #include <stdlib.h>
// #if defined(__APPLE__)
// #include <dlfcn.h>
// #elif defined(_WIN32)
// #define WIN32_LEAN_AND_MEAN 1
// #include <windows.h>
// #else
// #include <dlfcn.h>
// #endif
//
// #ifdef _WIN32
// static HMODULE hmod = NULL;
// #elif !defined __APPLE__
// static void* plib = NULL;
// #endif
//
// /* Helps bind function pointers to c functions. */
// static void* bindMethod(const char* name) {
// #ifdef __APPLE__
// 	return dlsym(RTLD_DEFAULT, name);
// #elif _WIN32
// 	void* pf = wglGetProcAddress((LPCSTR)name);
// 	if(pf) {
// 		return pf;
// 	}
// 	if(hmod == NULL) {
// 		hmod = LoadLibraryA("opengl32.dll");
// 	}
// 	return GetProcAddress(hmod, (LPCSTR)name);
// #else
// 	if(plib == NULL) {
// 		plib = dlopen("libGL.so", RTLD_LAZY);
// 	}
// 	return dlsym(plib, name);
// #endif
// }
//
// /* Helper method for string arrays */
// static char**newStringArray(int numberOfStrings) {
//     return calloc(sizeof(char*), numberOfStrings);
// }
//
// /* Helper method for string arrays */
// static void assignString(char **stringArray, char *string, int index) {
//     stringArray[index] = string;
// }
//
// /* Helper method for string arrays */
// static void freeStringArray(char **stringArray, int numberOfStrings) {
//     int cnt;
//     for (cnt = 0; cnt < numberOfStrings; cnt++)
//     {
//         free(stringArray[cnt]);
//     }
//     free(stringArray);
// }
//
// #ifndef APIENTRY
// #define APIENTRY
// #endif
// #ifndef APIENTRYP
// #define APIENTRYP APIENTRY *
// #endif
// #ifndef GLAPI
// #define GLAPI extern
// #endif
//
// /* Base GL types */
//
// typedef unsigned int GLenum;
// typedef unsigned char GLboolean;
// typedef unsigned int GLbitfield;
// typedef signed char GLbyte;
// typedef short GLshort;
// typedef int GLint;
// typedef int GLsizei;
// typedef unsigned char GLubyte;
// typedef unsigned short GLushort;
// typedef unsigned int GLuint;
// typedef unsigned short GLhalf;
// typedef float GLfloat;
// typedef float GLclampf;
// typedef double GLdouble;
// typedef double GLclampd;
// typedef void GLvoid;
//
// #include <stddef.h>
// #ifndef GL_VERSION_2_0
// /* GL type for program/shader text */
// typedef char GLchar;
// #endif
//
// #ifndef GL_VERSION_1_5
// /* GL types for handling large vertex buffer objects */
// typedef ptrdiff_t GLintptr;
// typedef ptrdiff_t GLsizeiptr;
// #endif
//
// #ifndef GL_ARB_vertex_buffer_object
// /* GL types for handling large vertex buffer objects */
// typedef ptrdiff_t GLintptrARB;
// typedef ptrdiff_t GLsizeiptrARB;
// #endif
//
// #ifndef GL_ARB_shader_objects
// /* GL types for program/shader text and shader object handles */
// typedef char GLcharARB;
// typedef unsigned int GLhandleARB;
// #endif
//
// /* GL type for "half" precision (s10e5) float data in host memory */
// #ifndef GL_ARB_half_float_pixel
// typedef unsigned short GLhalfARB;
// #endif
//
// #ifndef GL_NV_half_float
// typedef unsigned short GLhalfNV;
// #endif
//
// #ifndef GLEXT_64_TYPES_DEFINED
// /* This code block is duplicated in glxext.h, so must be protected */
// #define GLEXT_64_TYPES_DEFINED
// /* Define int32_t, int64_t, and uint64_t types for UST/MSC */
// /* (as used in the GL_EXT_timer_query extension). */
// #if defined(__STDC_VERSION__) && __STDC_VERSION__ >= 199901L
// #include <inttypes.h>
// #elif defined(__sun__) || defined(__digital__)
// #include <inttypes.h>
// #if defined(__STDC__)
// #if defined(__arch64__) || defined(_LP64)
// typedef long int int64_t;
// typedef unsigned long int uint64_t;
// #else
// typedef long long int int64_t;
// typedef unsigned long long int uint64_t;
// #endif /* __arch64__ */
// #endif /* __STDC__ */
// #elif defined( __VMS ) || defined(__sgi)
// #include <inttypes.h>
// #elif defined(__SCO__) || defined(__USLC__)
// #include <stdint.h>
// #elif defined(__UNIXOS2__) || defined(__SOL64__)
// typedef long int int32_t;
// typedef long long int int64_t;
// typedef unsigned long long int uint64_t;
// #elif defined(_WIN32) && defined(__GNUC__)
// #include <stdint.h>
// #elif defined(_WIN32)
// typedef __int32 int32_t;
// typedef __int64 int64_t;
// typedef unsigned __int64 uint64_t;
// #else
// /* Fallback if nothing above works */
// #include <inttypes.h>
// #endif
// #endif
//
// #ifndef GL_EXT_timer_query
// typedef int64_t GLint64EXT;
// typedef uint64_t GLuint64EXT;
// #endif
//
// #ifndef GL_ARB_sync
// typedef int64_t GLint64;
// typedef uint64_t GLuint64;
// typedef struct __GLsync *GLsync;
// #endif
//
// #ifndef GL_ARB_cl_event
// /* These incomplete types let us declare types compatible with OpenCL's cl_context and cl_event */
// struct _cl_context;
// struct _cl_event;
// #endif
//
// #ifndef GL_ARB_debug_output
// typedef void (APIENTRY *GLDEBUGPROCARB)(GLenum source,GLenum type,GLuint id,GLenum severity,GLsizei length,const GLchar *message,GLvoid *userParam);
// #endif
//
// #ifndef GL_AMD_debug_output
// typedef void (APIENTRY *GLDEBUGPROCAMD)(GLuint id,GLenum category,GLenum severity,GLsizei length,const GLchar *message,GLvoid *userParam);
// #endif
//
// #ifndef GL_KHR_debug
// typedef void (APIENTRY *GLDEBUGPROC)(GLenum source,GLenum type,GLuint id,GLenum severity,GLsizei length,const GLchar *message,GLvoid *userParam);
// #endif
//
// #ifndef GL_NV_vdpau_interop
// typedef GLintptr GLvdpauSurfaceNV;
// #endif
// GLboolean (APIENTRYP pfn_glIsBuffer)(GLuint buffer);
// GLAPI unsigned char APIENTRY wrap_glIsBuffer(unsigned int buffer) { return (*pfn_glIsBuffer)((GLuint)buffer); }
//
// GLboolean (APIENTRYP pfn_glIsEnabled)(GLenum cap);
// GLAPI unsigned char APIENTRY wrap_glIsEnabled(unsigned int cap) { return (*pfn_glIsEnabled)((GLenum)cap); }
//
// GLboolean (APIENTRYP pfn_glIsEnabledi)(GLenum target, GLuint index);
// GLAPI unsigned char APIENTRY wrap_glIsEnabledi(unsigned int target, unsigned int index) { return (*pfn_glIsEnabledi)((GLenum)target, (GLuint)index); }
//
// GLboolean (APIENTRYP pfn_glIsFramebuffer)(GLuint framebuffer);
// GLAPI unsigned char APIENTRY wrap_glIsFramebuffer(unsigned int framebuffer) { return (*pfn_glIsFramebuffer)((GLuint)framebuffer); }
//
// GLboolean (APIENTRYP pfn_glIsNamedStringARB)(GLint namelen, const GLchar *name);
// GLAPI unsigned char APIENTRY wrap_glIsNamedStringARB(int namelen, const char* name) { return (*pfn_glIsNamedStringARB)((GLint)namelen, (const GLchar *)name); }
//
// GLboolean (APIENTRYP pfn_glIsProgram)(GLuint program);
// GLAPI unsigned char APIENTRY wrap_glIsProgram(unsigned int program) { return (*pfn_glIsProgram)((GLuint)program); }
//
// GLboolean (APIENTRYP pfn_glIsProgramPipeline)(GLuint pipeline);
// GLAPI unsigned char APIENTRY wrap_glIsProgramPipeline(unsigned int pipeline) { return (*pfn_glIsProgramPipeline)((GLuint)pipeline); }
//
// GLboolean (APIENTRYP pfn_glIsQuery)(GLuint id);
// GLAPI unsigned char APIENTRY wrap_glIsQuery(unsigned int id) { return (*pfn_glIsQuery)((GLuint)id); }
//
// GLboolean (APIENTRYP pfn_glIsRenderbuffer)(GLuint renderbuffer);
// GLAPI unsigned char APIENTRY wrap_glIsRenderbuffer(unsigned int renderbuffer) { return (*pfn_glIsRenderbuffer)((GLuint)renderbuffer); }
//
// GLboolean (APIENTRYP pfn_glIsSampler)(GLuint sampler);
// GLAPI unsigned char APIENTRY wrap_glIsSampler(unsigned int sampler) { return (*pfn_glIsSampler)((GLuint)sampler); }
//
// GLboolean (APIENTRYP pfn_glIsShader)(GLuint shader);
// GLAPI unsigned char APIENTRY wrap_glIsShader(unsigned int shader) { return (*pfn_glIsShader)((GLuint)shader); }
//
// GLboolean (APIENTRYP pfn_glIsSync)(GLsync sync);
// GLAPI unsigned char APIENTRY wrap_glIsSync(GLsync sync) { return (*pfn_glIsSync)((GLsync)sync); }
//
// GLboolean (APIENTRYP pfn_glIsTexture)(GLuint texture);
// GLAPI unsigned char APIENTRY wrap_glIsTexture(unsigned int texture) { return (*pfn_glIsTexture)((GLuint)texture); }
//
// GLboolean (APIENTRYP pfn_glIsTransformFeedback)(GLuint id);
// GLAPI unsigned char APIENTRY wrap_glIsTransformFeedback(unsigned int id) { return (*pfn_glIsTransformFeedback)((GLuint)id); }
//
// GLboolean (APIENTRYP pfn_glIsVertexArray)(GLuint array);
// GLAPI unsigned char APIENTRY wrap_glIsVertexArray(unsigned int array) { return (*pfn_glIsVertexArray)((GLuint)array); }
//
// GLboolean (APIENTRYP pfn_glUnmapBuffer)(GLenum target);
// GLAPI unsigned char APIENTRY wrap_glUnmapBuffer(unsigned int target) { return (*pfn_glUnmapBuffer)((GLenum)target); }
//
// GLenum (APIENTRYP pfn_glCheckFramebufferStatus)(GLenum target);
// GLAPI unsigned int APIENTRY wrap_glCheckFramebufferStatus(unsigned int target) { return (*pfn_glCheckFramebufferStatus)((GLenum)target); }
//
// GLenum (APIENTRYP pfn_glClientWaitSync)(GLsync sync, GLbitfield flags, GLuint64 timeout);
// GLAPI unsigned int APIENTRY wrap_glClientWaitSync(GLsync sync, unsigned int flags, unsigned long long timeout) { return (*pfn_glClientWaitSync)((GLsync)sync, (GLbitfield)flags, (GLuint64)timeout); }
//
// GLenum (APIENTRYP pfn_glGetError)(void);
// GLAPI unsigned int APIENTRY wrap_glGetError() { return (*pfn_glGetError)(); }
//
// GLenum (APIENTRYP pfn_glGetGraphicsResetStatusARB)(void);
// GLAPI unsigned int APIENTRY wrap_glGetGraphicsResetStatusARB() { return (*pfn_glGetGraphicsResetStatusARB)(); }
//
// GLint (APIENTRYP pfn_glGetAttribLocation)(GLuint program, const GLchar *name);
// GLAPI int APIENTRY wrap_glGetAttribLocation(unsigned int program, const char* name) { return (*pfn_glGetAttribLocation)((GLuint)program, (const GLchar *)name); }
//
// GLint (APIENTRYP pfn_glGetFragDataIndex)(GLuint program, const GLchar *name);
// GLAPI int APIENTRY wrap_glGetFragDataIndex(unsigned int program, const char* name) { return (*pfn_glGetFragDataIndex)((GLuint)program, (const GLchar *)name); }
//
// GLint (APIENTRYP pfn_glGetFragDataLocation)(GLuint program, const GLchar *name);
// GLAPI int APIENTRY wrap_glGetFragDataLocation(unsigned int program, const char* name) { return (*pfn_glGetFragDataLocation)((GLuint)program, (const GLchar *)name); }
//
// GLint (APIENTRYP pfn_glGetProgramResourceLocation)(GLuint program, GLenum programInterface, const GLchar *name);
// GLAPI int APIENTRY wrap_glGetProgramResourceLocation(unsigned int program, unsigned int programInterface, const char* name) { return (*pfn_glGetProgramResourceLocation)((GLuint)program, (GLenum)programInterface, (const GLchar *)name); }
//
// GLint (APIENTRYP pfn_glGetProgramResourceLocationIndex)(GLuint program, GLenum programInterface, const GLchar *name);
// GLAPI int APIENTRY wrap_glGetProgramResourceLocationIndex(unsigned int program, unsigned int programInterface, const char* name) { return (*pfn_glGetProgramResourceLocationIndex)((GLuint)program, (GLenum)programInterface, (const GLchar *)name); }
//
// GLint (APIENTRYP pfn_glGetSubroutineUniformLocation)(GLuint program, GLenum shadertype, const GLchar *name);
// GLAPI int APIENTRY wrap_glGetSubroutineUniformLocation(unsigned int program, unsigned int shadertype, const char* name) { return (*pfn_glGetSubroutineUniformLocation)((GLuint)program, (GLenum)shadertype, (const GLchar *)name); }
//
// GLint (APIENTRYP pfn_glGetUniformLocation)(GLuint program, const GLchar *name);
// GLAPI int APIENTRY wrap_glGetUniformLocation(unsigned int program, const char* name) { return (*pfn_glGetUniformLocation)((GLuint)program, (const GLchar *)name); }
//
// GLsync (APIENTRYP pfn_glCreateSyncFromCLeventARB)(struct _cl_context * context, struct _cl_event * event, GLbitfield flags);
// GLAPI GLsync APIENTRY wrap_glCreateSyncFromCLeventARB(struct _cl_context* context, struct _cl_event* event, unsigned int flags) { return (*pfn_glCreateSyncFromCLeventARB)((struct _cl_context *)context, (struct _cl_event *)event, (GLbitfield)flags); }
//
// GLsync (APIENTRYP pfn_glFenceSync)(GLenum condition, GLbitfield flags);
// GLAPI GLsync APIENTRY wrap_glFenceSync(unsigned int condition, unsigned int flags) { return (*pfn_glFenceSync)((GLenum)condition, (GLbitfield)flags); }
//
// GLuint (APIENTRYP pfn_glCreateProgram)(void);
// GLAPI unsigned int APIENTRY wrap_glCreateProgram() { return (*pfn_glCreateProgram)(); }
//
// GLuint (APIENTRYP pfn_glCreateShader)(GLenum type);
// GLAPI unsigned int APIENTRY wrap_glCreateShader(unsigned int t_ype) { return (*pfn_glCreateShader)((GLenum)t_ype); }
//
// GLuint (APIENTRYP pfn_glCreateShaderProgramv)(GLenum type, GLsizei count, const GLchar* const *strings);
// GLAPI unsigned int APIENTRY wrap_glCreateShaderProgramv(unsigned int t_ype, int count, const char* const* strings) { return (*pfn_glCreateShaderProgramv)((GLenum)t_ype, (GLsizei)count, (const GLchar* const *)strings); }
//
// GLuint (APIENTRYP pfn_glGetDebugMessageLog)(GLuint count, GLsizei bufsize, GLenum *sources, GLenum *types, GLuint *ids, GLenum *severities, GLsizei *lengths, GLchar *messageLog);
// GLAPI unsigned int APIENTRY wrap_glGetDebugMessageLog(unsigned int count, int bufsize, unsigned int* sources, unsigned int* types, unsigned int* ids, unsigned int* severities, int* lengths, unsigned char* messageLog) { return (*pfn_glGetDebugMessageLog)((GLuint)count, (GLsizei)bufsize, (GLenum *)sources, (GLenum *)types, (GLuint *)ids, (GLenum *)severities, (GLsizei *)lengths, (GLchar *)messageLog); }
//
// GLuint (APIENTRYP pfn_glGetDebugMessageLogARB)(GLuint count, GLsizei bufsize, GLenum *sources, GLenum *types, GLuint *ids, GLenum *severities, GLsizei *lengths, GLchar *messageLog);
// GLAPI unsigned int APIENTRY wrap_glGetDebugMessageLogARB(unsigned int count, int bufsize, unsigned int* sources, unsigned int* types, unsigned int* ids, unsigned int* severities, int* lengths, unsigned char* messageLog) { return (*pfn_glGetDebugMessageLogARB)((GLuint)count, (GLsizei)bufsize, (GLenum *)sources, (GLenum *)types, (GLuint *)ids, (GLenum *)severities, (GLsizei *)lengths, (GLchar *)messageLog); }
//
// GLuint (APIENTRYP pfn_glGetProgramResourceIndex)(GLuint program, GLenum programInterface, const GLchar *name);
// GLAPI unsigned int APIENTRY wrap_glGetProgramResourceIndex(unsigned int program, unsigned int programInterface, const char* name) { return (*pfn_glGetProgramResourceIndex)((GLuint)program, (GLenum)programInterface, (const GLchar *)name); }
//
// GLuint (APIENTRYP pfn_glGetSubroutineIndex)(GLuint program, GLenum shadertype, const GLchar *name);
// GLAPI unsigned int APIENTRY wrap_glGetSubroutineIndex(unsigned int program, unsigned int shadertype, const char* name) { return (*pfn_glGetSubroutineIndex)((GLuint)program, (GLenum)shadertype, (const GLchar *)name); }
//
// GLuint (APIENTRYP pfn_glGetUniformBlockIndex)(GLuint program, const GLchar *uniformBlockName);
// GLAPI unsigned int APIENTRY wrap_glGetUniformBlockIndex(unsigned int program, const char* uniformBlockName) { return (*pfn_glGetUniformBlockIndex)((GLuint)program, (const GLchar *)uniformBlockName); }
//
// GLvoid* (APIENTRYP pfn_glMapBuffer)(GLenum target, GLenum access);
// GLAPI void* APIENTRY wrap_glMapBuffer(unsigned int target, unsigned int access) { return (*pfn_glMapBuffer)((GLenum)target, (GLenum)access); }
//
// GLvoid* (APIENTRYP pfn_glMapBufferRange)(GLenum target, GLintptr offset, GLsizeiptr length, GLbitfield access);
// GLAPI void* APIENTRY wrap_glMapBufferRange(unsigned int target, long long offset, long long length, unsigned int access) { return (*pfn_glMapBufferRange)((GLenum)target, (GLintptr)offset, (GLsizeiptr)length, (GLbitfield)access); }
//
// const GLubyte * (APIENTRYP pfn_glGetString)(GLenum name);
// GLAPI const char* APIENTRY wrap_glGetString(unsigned int name) { return (char *) (*pfn_glGetString)((GLenum)name); }
//
// const GLubyte * (APIENTRYP pfn_glGetStringi)(GLenum name, GLuint index);
// GLAPI const char* APIENTRY wrap_glGetStringi(unsigned int name, unsigned int index) { return (char *) (*pfn_glGetStringi)((GLenum)name, (GLuint)index); }
//
// void (APIENTRYP pfn_glActiveShaderProgram)(GLuint pipeline, GLuint program);
// GLAPI void APIENTRY wrap_glActiveShaderProgram(unsigned int pipeline, unsigned int program) {  (*pfn_glActiveShaderProgram)((GLuint)pipeline, (GLuint)program); }
//
// void (APIENTRYP pfn_glActiveTexture)(GLenum texture);
// GLAPI void APIENTRY wrap_glActiveTexture(unsigned int texture) {  (*pfn_glActiveTexture)((GLenum)texture); }
//
// void (APIENTRYP pfn_glAttachShader)(GLuint program, GLuint shader);
// GLAPI void APIENTRY wrap_glAttachShader(unsigned int program, unsigned int shader) {  (*pfn_glAttachShader)((GLuint)program, (GLuint)shader); }
//
// void (APIENTRYP pfn_glBeginConditionalRender)(GLuint id, GLenum mode);
// GLAPI void APIENTRY wrap_glBeginConditionalRender(unsigned int id, unsigned int mode) {  (*pfn_glBeginConditionalRender)((GLuint)id, (GLenum)mode); }
//
// void (APIENTRYP pfn_glBeginQuery)(GLenum target, GLuint id);
// GLAPI void APIENTRY wrap_glBeginQuery(unsigned int target, unsigned int id) {  (*pfn_glBeginQuery)((GLenum)target, (GLuint)id); }
//
// void (APIENTRYP pfn_glBeginQueryIndexed)(GLenum target, GLuint index, GLuint id);
// GLAPI void APIENTRY wrap_glBeginQueryIndexed(unsigned int target, unsigned int index, unsigned int id) {  (*pfn_glBeginQueryIndexed)((GLenum)target, (GLuint)index, (GLuint)id); }
//
// void (APIENTRYP pfn_glBeginTransformFeedback)(GLenum primitiveMode);
// GLAPI void APIENTRY wrap_glBeginTransformFeedback(unsigned int primitiveMode) {  (*pfn_glBeginTransformFeedback)((GLenum)primitiveMode); }
//
// void (APIENTRYP pfn_glBindAttribLocation)(GLuint program, GLuint index, const GLchar *name);
// GLAPI void APIENTRY wrap_glBindAttribLocation(unsigned int program, unsigned int index, const char* name) {  (*pfn_glBindAttribLocation)((GLuint)program, (GLuint)index, (const GLchar *)name); }
//
// void (APIENTRYP pfn_glBindBuffer)(GLenum target, GLuint buffer);
// GLAPI void APIENTRY wrap_glBindBuffer(unsigned int target, unsigned int buffer) {  (*pfn_glBindBuffer)((GLenum)target, (GLuint)buffer); }
//
// void (APIENTRYP pfn_glBindBufferBase)(GLenum target, GLuint index, GLuint buffer);
// GLAPI void APIENTRY wrap_glBindBufferBase(unsigned int target, unsigned int index, unsigned int buffer) {  (*pfn_glBindBufferBase)((GLenum)target, (GLuint)index, (GLuint)buffer); }
//
// void (APIENTRYP pfn_glBindBufferRange)(GLenum target, GLuint index, GLuint buffer, GLintptr offset, GLsizeiptr size);
// GLAPI void APIENTRY wrap_glBindBufferRange(unsigned int target, unsigned int index, unsigned int buffer, long long offset, long long size) {  (*pfn_glBindBufferRange)((GLenum)target, (GLuint)index, (GLuint)buffer, (GLintptr)offset, (GLsizeiptr)size); }
//
// void (APIENTRYP pfn_glBindFragDataLocation)(GLuint program, GLuint color, const GLchar *name);
// GLAPI void APIENTRY wrap_glBindFragDataLocation(unsigned int program, unsigned int color, const char* name) {  (*pfn_glBindFragDataLocation)((GLuint)program, (GLuint)color, (const GLchar *)name); }
//
// void (APIENTRYP pfn_glBindFragDataLocationIndexed)(GLuint program, GLuint colorNumber, GLuint index, const GLchar *name);
// GLAPI void APIENTRY wrap_glBindFragDataLocationIndexed(unsigned int program, unsigned int colorNumber, unsigned int index, const char* name) {  (*pfn_glBindFragDataLocationIndexed)((GLuint)program, (GLuint)colorNumber, (GLuint)index, (const GLchar *)name); }
//
// void (APIENTRYP pfn_glBindFramebuffer)(GLenum target, GLuint framebuffer);
// GLAPI void APIENTRY wrap_glBindFramebuffer(unsigned int target, unsigned int framebuffer) {  (*pfn_glBindFramebuffer)((GLenum)target, (GLuint)framebuffer); }
//
// void (APIENTRYP pfn_glBindImageTexture)(GLuint unit, GLuint texture, GLint level, GLboolean layered, GLint layer, GLenum access, GLenum format);
// GLAPI void APIENTRY wrap_glBindImageTexture(unsigned int unit, unsigned int texture, int level, unsigned char layered, int layer, unsigned int access, unsigned int format) {  (*pfn_glBindImageTexture)((GLuint)unit, (GLuint)texture, (GLint)level, (GLboolean)layered, (GLint)layer, (GLenum)access, (GLenum)format); }
//
// void (APIENTRYP pfn_glBindProgramPipeline)(GLuint pipeline);
// GLAPI void APIENTRY wrap_glBindProgramPipeline(unsigned int pipeline) {  (*pfn_glBindProgramPipeline)((GLuint)pipeline); }
//
// void (APIENTRYP pfn_glBindRenderbuffer)(GLenum target, GLuint renderbuffer);
// GLAPI void APIENTRY wrap_glBindRenderbuffer(unsigned int target, unsigned int renderbuffer) {  (*pfn_glBindRenderbuffer)((GLenum)target, (GLuint)renderbuffer); }
//
// void (APIENTRYP pfn_glBindSampler)(GLuint unit, GLuint sampler);
// GLAPI void APIENTRY wrap_glBindSampler(unsigned int unit, unsigned int sampler) {  (*pfn_glBindSampler)((GLuint)unit, (GLuint)sampler); }
//
// void (APIENTRYP pfn_glBindTexture)(GLenum target, GLuint texture);
// GLAPI void APIENTRY wrap_glBindTexture(unsigned int target, unsigned int texture) {  (*pfn_glBindTexture)((GLenum)target, (GLuint)texture); }
//
// void (APIENTRYP pfn_glBindTransformFeedback)(GLenum target, GLuint id);
// GLAPI void APIENTRY wrap_glBindTransformFeedback(unsigned int target, unsigned int id) {  (*pfn_glBindTransformFeedback)((GLenum)target, (GLuint)id); }
//
// void (APIENTRYP pfn_glBindVertexArray)(GLuint array);
// GLAPI void APIENTRY wrap_glBindVertexArray(unsigned int array) {  (*pfn_glBindVertexArray)((GLuint)array); }
//
// void (APIENTRYP pfn_glBindVertexBuffer)(GLuint bindingindex, GLuint buffer, GLintptr offset, GLsizei stride);
// GLAPI void APIENTRY wrap_glBindVertexBuffer(unsigned int bindingindex, unsigned int buffer, long long offset, int stride) {  (*pfn_glBindVertexBuffer)((GLuint)bindingindex, (GLuint)buffer, (GLintptr)offset, (GLsizei)stride); }
//
// void (APIENTRYP pfn_glBlendColor)(GLfloat red, GLfloat green, GLfloat blue, GLfloat alpha);
// GLAPI void APIENTRY wrap_glBlendColor(float red, float green, float blue, float alpha) {  (*pfn_glBlendColor)((GLfloat)red, (GLfloat)green, (GLfloat)blue, (GLfloat)alpha); }
//
// void (APIENTRYP pfn_glBlendEquation)(GLenum mode);
// GLAPI void APIENTRY wrap_glBlendEquation(unsigned int mode) {  (*pfn_glBlendEquation)((GLenum)mode); }
//
// void (APIENTRYP pfn_glBlendEquationSeparate)(GLenum modeRGB, GLenum modeAlpha);
// GLAPI void APIENTRY wrap_glBlendEquationSeparate(unsigned int modeRGB, unsigned int modeAlpha) {  (*pfn_glBlendEquationSeparate)((GLenum)modeRGB, (GLenum)modeAlpha); }
//
// void (APIENTRYP pfn_glBlendEquationSeparatei)(GLuint buf, GLenum modeRGB, GLenum modeAlpha);
// GLAPI void APIENTRY wrap_glBlendEquationSeparatei(unsigned int buf, unsigned int modeRGB, unsigned int modeAlpha) {  (*pfn_glBlendEquationSeparatei)((GLuint)buf, (GLenum)modeRGB, (GLenum)modeAlpha); }
//
// void (APIENTRYP pfn_glBlendEquationSeparateiARB)(GLuint buf, GLenum modeRGB, GLenum modeAlpha);
// GLAPI void APIENTRY wrap_glBlendEquationSeparateiARB(unsigned int buf, unsigned int modeRGB, unsigned int modeAlpha) {  (*pfn_glBlendEquationSeparateiARB)((GLuint)buf, (GLenum)modeRGB, (GLenum)modeAlpha); }
//
// void (APIENTRYP pfn_glBlendEquationi)(GLuint buf, GLenum mode);
// GLAPI void APIENTRY wrap_glBlendEquationi(unsigned int buf, unsigned int mode) {  (*pfn_glBlendEquationi)((GLuint)buf, (GLenum)mode); }
//
// void (APIENTRYP pfn_glBlendEquationiARB)(GLuint buf, GLenum mode);
// GLAPI void APIENTRY wrap_glBlendEquationiARB(unsigned int buf, unsigned int mode) {  (*pfn_glBlendEquationiARB)((GLuint)buf, (GLenum)mode); }
//
// void (APIENTRYP pfn_glBlendFunc)(GLenum sfactor, GLenum dfactor);
// GLAPI void APIENTRY wrap_glBlendFunc(unsigned int sfactor, unsigned int dfactor) {  (*pfn_glBlendFunc)((GLenum)sfactor, (GLenum)dfactor); }
//
// void (APIENTRYP pfn_glBlendFuncSeparate)(GLenum sfactorRGB, GLenum dfactorRGB, GLenum sfactorAlpha, GLenum dfactorAlpha);
// GLAPI void APIENTRY wrap_glBlendFuncSeparate(unsigned int sfactorRGB, unsigned int dfactorRGB, unsigned int sfactorAlpha, unsigned int dfactorAlpha) {  (*pfn_glBlendFuncSeparate)((GLenum)sfactorRGB, (GLenum)dfactorRGB, (GLenum)sfactorAlpha, (GLenum)dfactorAlpha); }
//
// void (APIENTRYP pfn_glBlendFuncSeparatei)(GLuint buf, GLenum srcRGB, GLenum dstRGB, GLenum srcAlpha, GLenum dstAlpha);
// GLAPI void APIENTRY wrap_glBlendFuncSeparatei(unsigned int buf, unsigned int srcRGB, unsigned int dstRGB, unsigned int srcAlpha, unsigned int dstAlpha) {  (*pfn_glBlendFuncSeparatei)((GLuint)buf, (GLenum)srcRGB, (GLenum)dstRGB, (GLenum)srcAlpha, (GLenum)dstAlpha); }
//
// void (APIENTRYP pfn_glBlendFuncSeparateiARB)(GLuint buf, GLenum srcRGB, GLenum dstRGB, GLenum srcAlpha, GLenum dstAlpha);
// GLAPI void APIENTRY wrap_glBlendFuncSeparateiARB(unsigned int buf, unsigned int srcRGB, unsigned int dstRGB, unsigned int srcAlpha, unsigned int dstAlpha) {  (*pfn_glBlendFuncSeparateiARB)((GLuint)buf, (GLenum)srcRGB, (GLenum)dstRGB, (GLenum)srcAlpha, (GLenum)dstAlpha); }
//
// void (APIENTRYP pfn_glBlendFunci)(GLuint buf, GLenum src, GLenum dst);
// GLAPI void APIENTRY wrap_glBlendFunci(unsigned int buf, unsigned int src, unsigned int dst) {  (*pfn_glBlendFunci)((GLuint)buf, (GLenum)src, (GLenum)dst); }
//
// void (APIENTRYP pfn_glBlendFunciARB)(GLuint buf, GLenum src, GLenum dst);
// GLAPI void APIENTRY wrap_glBlendFunciARB(unsigned int buf, unsigned int src, unsigned int dst) {  (*pfn_glBlendFunciARB)((GLuint)buf, (GLenum)src, (GLenum)dst); }
//
// void (APIENTRYP pfn_glBlitFramebuffer)(GLint srcX0, GLint srcY0, GLint srcX1, GLint srcY1, GLint dstX0, GLint dstY0, GLint dstX1, GLint dstY1, GLbitfield mask, GLenum filter);
// GLAPI void APIENTRY wrap_glBlitFramebuffer(int srcX0, int srcY0, int srcX1, int srcY1, int dstX0, int dstY0, int dstX1, int dstY1, unsigned int mask, unsigned int filter) {  (*pfn_glBlitFramebuffer)((GLint)srcX0, (GLint)srcY0, (GLint)srcX1, (GLint)srcY1, (GLint)dstX0, (GLint)dstY0, (GLint)dstX1, (GLint)dstY1, (GLbitfield)mask, (GLenum)filter); }
//
// void (APIENTRYP pfn_glBufferData)(GLenum target, GLsizeiptr size, const GLvoid *data, GLenum usage);
// GLAPI void APIENTRY wrap_glBufferData(unsigned int target, long long size, const void* data, unsigned int usage) {  (*pfn_glBufferData)((GLenum)target, (GLsizeiptr)size, (const GLvoid *)data, (GLenum)usage); }
//
// void (APIENTRYP pfn_glBufferSubData)(GLenum target, GLintptr offset, GLsizeiptr size, const GLvoid *data);
// GLAPI void APIENTRY wrap_glBufferSubData(unsigned int target, long long offset, long long size, const void* data) {  (*pfn_glBufferSubData)((GLenum)target, (GLintptr)offset, (GLsizeiptr)size, (const GLvoid *)data); }
//
// void (APIENTRYP pfn_glClampColor)(GLenum target, GLenum clamp);
// GLAPI void APIENTRY wrap_glClampColor(unsigned int target, unsigned int clamp) {  (*pfn_glClampColor)((GLenum)target, (GLenum)clamp); }
//
// void (APIENTRYP pfn_glClear)(GLbitfield mask);
// GLAPI void APIENTRY wrap_glClear(unsigned int mask) {  (*pfn_glClear)((GLbitfield)mask); }
//
// void (APIENTRYP pfn_glClearBufferData)(GLenum target, GLenum internalformat, GLenum format, GLenum type, const void *data);
// GLAPI void APIENTRY wrap_glClearBufferData(unsigned int target, unsigned int internalformat, unsigned int format, unsigned int t_ype, const void* data) {  (*pfn_glClearBufferData)((GLenum)target, (GLenum)internalformat, (GLenum)format, (GLenum)t_ype, (const void *)data); }
//
// void (APIENTRYP pfn_glClearBufferSubData)(GLenum target, GLenum internalformat, GLintptr offset, GLsizeiptr size, GLenum format, GLenum type, const void *data);
// GLAPI void APIENTRY wrap_glClearBufferSubData(unsigned int target, unsigned int internalformat, long long offset, long long size, unsigned int format, unsigned int t_ype, const void* data) {  (*pfn_glClearBufferSubData)((GLenum)target, (GLenum)internalformat, (GLintptr)offset, (GLsizeiptr)size, (GLenum)format, (GLenum)t_ype, (const void *)data); }
//
// void (APIENTRYP pfn_glClearBufferfi)(GLenum buffer, GLint drawbuffer, GLfloat depth, GLint stencil);
// GLAPI void APIENTRY wrap_glClearBufferfi(unsigned int buffer, int drawbuffer, float depth, int stencil) {  (*pfn_glClearBufferfi)((GLenum)buffer, (GLint)drawbuffer, (GLfloat)depth, (GLint)stencil); }
//
// void (APIENTRYP pfn_glClearBufferfv)(GLenum buffer, GLint drawbuffer, const GLfloat *value);
// GLAPI void APIENTRY wrap_glClearBufferfv(unsigned int buffer, int drawbuffer, const float* value) {  (*pfn_glClearBufferfv)((GLenum)buffer, (GLint)drawbuffer, (const GLfloat *)value); }
//
// void (APIENTRYP pfn_glClearBufferiv)(GLenum buffer, GLint drawbuffer, const GLint *value);
// GLAPI void APIENTRY wrap_glClearBufferiv(unsigned int buffer, int drawbuffer, const int* value) {  (*pfn_glClearBufferiv)((GLenum)buffer, (GLint)drawbuffer, (const GLint *)value); }
//
// void (APIENTRYP pfn_glClearBufferuiv)(GLenum buffer, GLint drawbuffer, const GLuint *value);
// GLAPI void APIENTRY wrap_glClearBufferuiv(unsigned int buffer, int drawbuffer, const unsigned int* value) {  (*pfn_glClearBufferuiv)((GLenum)buffer, (GLint)drawbuffer, (const GLuint *)value); }
//
// void (APIENTRYP pfn_glClearColor)(GLfloat red, GLfloat green, GLfloat blue, GLfloat alpha);
// GLAPI void APIENTRY wrap_glClearColor(float red, float green, float blue, float alpha) {  (*pfn_glClearColor)((GLfloat)red, (GLfloat)green, (GLfloat)blue, (GLfloat)alpha); }
//
// void (APIENTRYP pfn_glClearDepth)(GLdouble depth);
// GLAPI void APIENTRY wrap_glClearDepth(double depth) {  (*pfn_glClearDepth)((GLdouble)depth); }
//
// void (APIENTRYP pfn_glClearDepthf)(GLfloat d);
// GLAPI void APIENTRY wrap_glClearDepthf(float d) {  (*pfn_glClearDepthf)((GLfloat)d); }
//
// void (APIENTRYP pfn_glClearNamedBufferDataEXT)(GLuint buffer, GLenum internalformat, GLenum format, GLenum type, const void *data);
// GLAPI void APIENTRY wrap_glClearNamedBufferDataEXT(unsigned int buffer, unsigned int internalformat, unsigned int format, unsigned int t_ype, const void* data) {  (*pfn_glClearNamedBufferDataEXT)((GLuint)buffer, (GLenum)internalformat, (GLenum)format, (GLenum)t_ype, (const void *)data); }
//
// void (APIENTRYP pfn_glClearNamedBufferSubDataEXT)(GLuint buffer, GLenum internalformat, GLenum format, GLenum type, GLsizeiptr offset, GLsizeiptr size, const void *data);
// GLAPI void APIENTRY wrap_glClearNamedBufferSubDataEXT(unsigned int buffer, unsigned int internalformat, unsigned int format, unsigned int t_ype, long long offset, long long size, const void* data) {  (*pfn_glClearNamedBufferSubDataEXT)((GLuint)buffer, (GLenum)internalformat, (GLenum)format, (GLenum)t_ype, (GLsizeiptr)offset, (GLsizeiptr)size, (const void *)data); }
//
// void (APIENTRYP pfn_glClearStencil)(GLint s);
// GLAPI void APIENTRY wrap_glClearStencil(int s) {  (*pfn_glClearStencil)((GLint)s); }
//
// void (APIENTRYP pfn_glColorMask)(GLboolean red, GLboolean green, GLboolean blue, GLboolean alpha);
// GLAPI void APIENTRY wrap_glColorMask(unsigned char red, unsigned char green, unsigned char blue, unsigned char alpha) {  (*pfn_glColorMask)((GLboolean)red, (GLboolean)green, (GLboolean)blue, (GLboolean)alpha); }
//
// void (APIENTRYP pfn_glColorMaski)(GLuint index, GLboolean r, GLboolean g, GLboolean b, GLboolean a);
// GLAPI void APIENTRY wrap_glColorMaski(unsigned int index, unsigned char r, unsigned char g, unsigned char b, unsigned char a) {  (*pfn_glColorMaski)((GLuint)index, (GLboolean)r, (GLboolean)g, (GLboolean)b, (GLboolean)a); }
//
// void (APIENTRYP pfn_glColorP3ui)(GLenum type, GLuint color);
// GLAPI void APIENTRY wrap_glColorP3ui(unsigned int t_ype, unsigned int color) {  (*pfn_glColorP3ui)((GLenum)t_ype, (GLuint)color); }
//
// void (APIENTRYP pfn_glColorP3uiv)(GLenum type, const GLuint *color);
// GLAPI void APIENTRY wrap_glColorP3uiv(unsigned int t_ype, const unsigned int* color) {  (*pfn_glColorP3uiv)((GLenum)t_ype, (const GLuint *)color); }
//
// void (APIENTRYP pfn_glColorP4ui)(GLenum type, GLuint color);
// GLAPI void APIENTRY wrap_glColorP4ui(unsigned int t_ype, unsigned int color) {  (*pfn_glColorP4ui)((GLenum)t_ype, (GLuint)color); }
//
// void (APIENTRYP pfn_glColorP4uiv)(GLenum type, const GLuint *color);
// GLAPI void APIENTRY wrap_glColorP4uiv(unsigned int t_ype, const unsigned int* color) {  (*pfn_glColorP4uiv)((GLenum)t_ype, (const GLuint *)color); }
//
// void (APIENTRYP pfn_glCompileShader)(GLuint shader);
// GLAPI void APIENTRY wrap_glCompileShader(unsigned int shader) {  (*pfn_glCompileShader)((GLuint)shader); }
//
// void (APIENTRYP pfn_glCompileShaderIncludeARB)(GLuint shader, GLsizei count, const GLchar* *path, const GLint *length);
// GLAPI void APIENTRY wrap_glCompileShaderIncludeARB(unsigned int shader, int count, const char** path, const int* length) {  (*pfn_glCompileShaderIncludeARB)((GLuint)shader, (GLsizei)count, (const GLchar* *)path, (const GLint *)length); }
//
// void (APIENTRYP pfn_glCompressedTexImage1D)(GLenum target, GLint level, GLenum internalformat, GLsizei width, GLint border, GLsizei imageSize, const GLvoid *data);
// GLAPI void APIENTRY wrap_glCompressedTexImage1D(unsigned int target, int level, unsigned int internalformat, int width, int border, int imageSize, const void* data) {  (*pfn_glCompressedTexImage1D)((GLenum)target, (GLint)level, (GLenum)internalformat, (GLsizei)width, (GLint)border, (GLsizei)imageSize, (const GLvoid *)data); }
//
// void (APIENTRYP pfn_glCompressedTexImage2D)(GLenum target, GLint level, GLenum internalformat, GLsizei width, GLsizei height, GLint border, GLsizei imageSize, const GLvoid *data);
// GLAPI void APIENTRY wrap_glCompressedTexImage2D(unsigned int target, int level, unsigned int internalformat, int width, int height, int border, int imageSize, const void* data) {  (*pfn_glCompressedTexImage2D)((GLenum)target, (GLint)level, (GLenum)internalformat, (GLsizei)width, (GLsizei)height, (GLint)border, (GLsizei)imageSize, (const GLvoid *)data); }
//
// void (APIENTRYP pfn_glCompressedTexImage3D)(GLenum target, GLint level, GLenum internalformat, GLsizei width, GLsizei height, GLsizei depth, GLint border, GLsizei imageSize, const GLvoid *data);
// GLAPI void APIENTRY wrap_glCompressedTexImage3D(unsigned int target, int level, unsigned int internalformat, int width, int height, int depth, int border, int imageSize, const void* data) {  (*pfn_glCompressedTexImage3D)((GLenum)target, (GLint)level, (GLenum)internalformat, (GLsizei)width, (GLsizei)height, (GLsizei)depth, (GLint)border, (GLsizei)imageSize, (const GLvoid *)data); }
//
// void (APIENTRYP pfn_glCompressedTexSubImage1D)(GLenum target, GLint level, GLint xoffset, GLsizei width, GLenum format, GLsizei imageSize, const GLvoid *data);
// GLAPI void APIENTRY wrap_glCompressedTexSubImage1D(unsigned int target, int level, int xoffset, int width, unsigned int format, int imageSize, const void* data) {  (*pfn_glCompressedTexSubImage1D)((GLenum)target, (GLint)level, (GLint)xoffset, (GLsizei)width, (GLenum)format, (GLsizei)imageSize, (const GLvoid *)data); }
//
// void (APIENTRYP pfn_glCompressedTexSubImage2D)(GLenum target, GLint level, GLint xoffset, GLint yoffset, GLsizei width, GLsizei height, GLenum format, GLsizei imageSize, const GLvoid *data);
// GLAPI void APIENTRY wrap_glCompressedTexSubImage2D(unsigned int target, int level, int xoffset, int yoffset, int width, int height, unsigned int format, int imageSize, const void* data) {  (*pfn_glCompressedTexSubImage2D)((GLenum)target, (GLint)level, (GLint)xoffset, (GLint)yoffset, (GLsizei)width, (GLsizei)height, (GLenum)format, (GLsizei)imageSize, (const GLvoid *)data); }
//
// void (APIENTRYP pfn_glCompressedTexSubImage3D)(GLenum target, GLint level, GLint xoffset, GLint yoffset, GLint zoffset, GLsizei width, GLsizei height, GLsizei depth, GLenum format, GLsizei imageSize, const GLvoid *data);
// GLAPI void APIENTRY wrap_glCompressedTexSubImage3D(unsigned int target, int level, int xoffset, int yoffset, int zoffset, int width, int height, int depth, unsigned int format, int imageSize, const void* data) {  (*pfn_glCompressedTexSubImage3D)((GLenum)target, (GLint)level, (GLint)xoffset, (GLint)yoffset, (GLint)zoffset, (GLsizei)width, (GLsizei)height, (GLsizei)depth, (GLenum)format, (GLsizei)imageSize, (const GLvoid *)data); }
//
// void (APIENTRYP pfn_glCopyBufferSubData)(GLenum readTarget, GLenum writeTarget, GLintptr readOffset, GLintptr writeOffset, GLsizeiptr size);
// GLAPI void APIENTRY wrap_glCopyBufferSubData(unsigned int readTarget, unsigned int writeTarget, long long readOffset, long long writeOffset, long long size) {  (*pfn_glCopyBufferSubData)((GLenum)readTarget, (GLenum)writeTarget, (GLintptr)readOffset, (GLintptr)writeOffset, (GLsizeiptr)size); }
//
// void (APIENTRYP pfn_glCopyImageSubData)(GLuint srcName, GLenum srcTarget, GLint srcLevel, GLint srcX, GLint srcY, GLint srcZ, GLuint dstName, GLenum dstTarget, GLint dstLevel, GLint dstX, GLint dstY, GLint dstZ, GLsizei srcWidth, GLsizei srcHeight, GLsizei srcDepth);
// GLAPI void APIENTRY wrap_glCopyImageSubData(unsigned int srcName, unsigned int srcTarget, int srcLevel, int srcX, int srcY, int srcZ, unsigned int dstName, unsigned int dstTarget, int dstLevel, int dstX, int dstY, int dstZ, int srcWidth, int srcHeight, int srcDepth) {  (*pfn_glCopyImageSubData)((GLuint)srcName, (GLenum)srcTarget, (GLint)srcLevel, (GLint)srcX, (GLint)srcY, (GLint)srcZ, (GLuint)dstName, (GLenum)dstTarget, (GLint)dstLevel, (GLint)dstX, (GLint)dstY, (GLint)dstZ, (GLsizei)srcWidth, (GLsizei)srcHeight, (GLsizei)srcDepth); }
//
// void (APIENTRYP pfn_glCopyTexImage1D)(GLenum target, GLint level, GLenum internalformat, GLint x, GLint y, GLsizei width, GLint border);
// GLAPI void APIENTRY wrap_glCopyTexImage1D(unsigned int target, int level, unsigned int internalformat, int x, int y, int width, int border) {  (*pfn_glCopyTexImage1D)((GLenum)target, (GLint)level, (GLenum)internalformat, (GLint)x, (GLint)y, (GLsizei)width, (GLint)border); }
//
// void (APIENTRYP pfn_glCopyTexImage2D)(GLenum target, GLint level, GLenum internalformat, GLint x, GLint y, GLsizei width, GLsizei height, GLint border);
// GLAPI void APIENTRY wrap_glCopyTexImage2D(unsigned int target, int level, unsigned int internalformat, int x, int y, int width, int height, int border) {  (*pfn_glCopyTexImage2D)((GLenum)target, (GLint)level, (GLenum)internalformat, (GLint)x, (GLint)y, (GLsizei)width, (GLsizei)height, (GLint)border); }
//
// void (APIENTRYP pfn_glCopyTexSubImage1D)(GLenum target, GLint level, GLint xoffset, GLint x, GLint y, GLsizei width);
// GLAPI void APIENTRY wrap_glCopyTexSubImage1D(unsigned int target, int level, int xoffset, int x, int y, int width) {  (*pfn_glCopyTexSubImage1D)((GLenum)target, (GLint)level, (GLint)xoffset, (GLint)x, (GLint)y, (GLsizei)width); }
//
// void (APIENTRYP pfn_glCopyTexSubImage2D)(GLenum target, GLint level, GLint xoffset, GLint yoffset, GLint x, GLint y, GLsizei width, GLsizei height);
// GLAPI void APIENTRY wrap_glCopyTexSubImage2D(unsigned int target, int level, int xoffset, int yoffset, int x, int y, int width, int height) {  (*pfn_glCopyTexSubImage2D)((GLenum)target, (GLint)level, (GLint)xoffset, (GLint)yoffset, (GLint)x, (GLint)y, (GLsizei)width, (GLsizei)height); }
//
// void (APIENTRYP pfn_glCopyTexSubImage3D)(GLenum target, GLint level, GLint xoffset, GLint yoffset, GLint zoffset, GLint x, GLint y, GLsizei width, GLsizei height);
// GLAPI void APIENTRY wrap_glCopyTexSubImage3D(unsigned int target, int level, int xoffset, int yoffset, int zoffset, int x, int y, int width, int height) {  (*pfn_glCopyTexSubImage3D)((GLenum)target, (GLint)level, (GLint)xoffset, (GLint)yoffset, (GLint)zoffset, (GLint)x, (GLint)y, (GLsizei)width, (GLsizei)height); }
//
// void (APIENTRYP pfn_glCullFace)(GLenum mode);
// GLAPI void APIENTRY wrap_glCullFace(unsigned int mode) {  (*pfn_glCullFace)((GLenum)mode); }
//
// void (APIENTRYP pfn_glDebugMessageCallback)(GLDEBUGPROC callback, const void *userParam);
// GLAPI void APIENTRY wrap_glDebugMessageCallback(GLDEBUGPROC callback, const void* userParam) {  (*pfn_glDebugMessageCallback)((GLDEBUGPROC)callback, (const void *)userParam); }
//
// void (APIENTRYP pfn_glDebugMessageCallbackARB)(GLDEBUGPROCARB callback, const GLvoid *userParam);
// GLAPI void APIENTRY wrap_glDebugMessageCallbackARB(GLDEBUGPROCARB callback, const void* userParam) {  (*pfn_glDebugMessageCallbackARB)((GLDEBUGPROCARB)callback, (const GLvoid *)userParam); }
//
// void (APIENTRYP pfn_glDebugMessageControl)(GLenum source, GLenum type, GLenum severity, GLsizei count, const GLuint *ids, GLboolean enabled);
// GLAPI void APIENTRY wrap_glDebugMessageControl(unsigned int source, unsigned int t_ype, unsigned int severity, int count, const unsigned int* ids, unsigned char enabled) {  (*pfn_glDebugMessageControl)((GLenum)source, (GLenum)t_ype, (GLenum)severity, (GLsizei)count, (const GLuint *)ids, (GLboolean)enabled); }
//
// void (APIENTRYP pfn_glDebugMessageControlARB)(GLenum source, GLenum type, GLenum severity, GLsizei count, const GLuint *ids, GLboolean enabled);
// GLAPI void APIENTRY wrap_glDebugMessageControlARB(unsigned int source, unsigned int t_ype, unsigned int severity, int count, const unsigned int* ids, unsigned char enabled) {  (*pfn_glDebugMessageControlARB)((GLenum)source, (GLenum)t_ype, (GLenum)severity, (GLsizei)count, (const GLuint *)ids, (GLboolean)enabled); }
//
// void (APIENTRYP pfn_glDebugMessageInsert)(GLenum source, GLenum type, GLuint id, GLenum severity, GLsizei length, const GLchar *buf);
// GLAPI void APIENTRY wrap_glDebugMessageInsert(unsigned int source, unsigned int t_ype, unsigned int id, unsigned int severity, int length, const char* buf) {  (*pfn_glDebugMessageInsert)((GLenum)source, (GLenum)t_ype, (GLuint)id, (GLenum)severity, (GLsizei)length, (const GLchar *)buf); }
//
// void (APIENTRYP pfn_glDebugMessageInsertARB)(GLenum source, GLenum type, GLuint id, GLenum severity, GLsizei length, const GLchar *buf);
// GLAPI void APIENTRY wrap_glDebugMessageInsertARB(unsigned int source, unsigned int t_ype, unsigned int id, unsigned int severity, int length, const char* buf) {  (*pfn_glDebugMessageInsertARB)((GLenum)source, (GLenum)t_ype, (GLuint)id, (GLenum)severity, (GLsizei)length, (const GLchar *)buf); }
//
// void (APIENTRYP pfn_glDeleteBuffers)(GLsizei n, const GLuint *buffers);
// GLAPI void APIENTRY wrap_glDeleteBuffers(int n, const unsigned int* buffers) {  (*pfn_glDeleteBuffers)((GLsizei)n, (const GLuint *)buffers); }
//
// void (APIENTRYP pfn_glDeleteFramebuffers)(GLsizei n, const GLuint *framebuffers);
// GLAPI void APIENTRY wrap_glDeleteFramebuffers(int n, const unsigned int* framebuffers) {  (*pfn_glDeleteFramebuffers)((GLsizei)n, (const GLuint *)framebuffers); }
//
// void (APIENTRYP pfn_glDeleteNamedStringARB)(GLint namelen, const GLchar *name);
// GLAPI void APIENTRY wrap_glDeleteNamedStringARB(int namelen, const char* name) {  (*pfn_glDeleteNamedStringARB)((GLint)namelen, (const GLchar *)name); }
//
// void (APIENTRYP pfn_glDeleteProgram)(GLuint program);
// GLAPI void APIENTRY wrap_glDeleteProgram(unsigned int program) {  (*pfn_glDeleteProgram)((GLuint)program); }
//
// void (APIENTRYP pfn_glDeleteProgramPipelines)(GLsizei n, const GLuint *pipelines);
// GLAPI void APIENTRY wrap_glDeleteProgramPipelines(int n, const unsigned int* pipelines) {  (*pfn_glDeleteProgramPipelines)((GLsizei)n, (const GLuint *)pipelines); }
//
// void (APIENTRYP pfn_glDeleteQueries)(GLsizei n, const GLuint *ids);
// GLAPI void APIENTRY wrap_glDeleteQueries(int n, const unsigned int* ids) {  (*pfn_glDeleteQueries)((GLsizei)n, (const GLuint *)ids); }
//
// void (APIENTRYP pfn_glDeleteRenderbuffers)(GLsizei n, const GLuint *renderbuffers);
// GLAPI void APIENTRY wrap_glDeleteRenderbuffers(int n, const unsigned int* renderbuffers) {  (*pfn_glDeleteRenderbuffers)((GLsizei)n, (const GLuint *)renderbuffers); }
//
// void (APIENTRYP pfn_glDeleteSamplers)(GLsizei count, const GLuint *samplers);
// GLAPI void APIENTRY wrap_glDeleteSamplers(int count, const unsigned int* samplers) {  (*pfn_glDeleteSamplers)((GLsizei)count, (const GLuint *)samplers); }
//
// void (APIENTRYP pfn_glDeleteShader)(GLuint shader);
// GLAPI void APIENTRY wrap_glDeleteShader(unsigned int shader) {  (*pfn_glDeleteShader)((GLuint)shader); }
//
// void (APIENTRYP pfn_glDeleteSync)(GLsync sync);
// GLAPI void APIENTRY wrap_glDeleteSync(GLsync sync) {  (*pfn_glDeleteSync)((GLsync)sync); }
//
// void (APIENTRYP pfn_glDeleteTextures)(GLsizei n, const GLuint *textures);
// GLAPI void APIENTRY wrap_glDeleteTextures(int n, const unsigned int* textures) {  (*pfn_glDeleteTextures)((GLsizei)n, (const GLuint *)textures); }
//
// void (APIENTRYP pfn_glDeleteTransformFeedbacks)(GLsizei n, const GLuint *ids);
// GLAPI void APIENTRY wrap_glDeleteTransformFeedbacks(int n, const unsigned int* ids) {  (*pfn_glDeleteTransformFeedbacks)((GLsizei)n, (const GLuint *)ids); }
//
// void (APIENTRYP pfn_glDeleteVertexArrays)(GLsizei n, const GLuint *arrays);
// GLAPI void APIENTRY wrap_glDeleteVertexArrays(int n, const unsigned int* arrays) {  (*pfn_glDeleteVertexArrays)((GLsizei)n, (const GLuint *)arrays); }
//
// void (APIENTRYP pfn_glDepthFunc)(GLenum func);
// GLAPI void APIENTRY wrap_glDepthFunc(unsigned int f_unc) {  (*pfn_glDepthFunc)((GLenum)f_unc); }
//
// void (APIENTRYP pfn_glDepthMask)(GLboolean flag);
// GLAPI void APIENTRY wrap_glDepthMask(unsigned char flag) {  (*pfn_glDepthMask)((GLboolean)flag); }
//
// void (APIENTRYP pfn_glDepthRange)(GLdouble near, GLdouble far);
// GLAPI void APIENTRY wrap_glDepthRange(double n_ear, double f_ar) {  (*pfn_glDepthRange)((GLdouble)n_ear, (GLdouble)f_ar); }
//
// void (APIENTRYP pfn_glDepthRangeArrayv)(GLuint first, GLsizei count, const GLdouble *v);
// GLAPI void APIENTRY wrap_glDepthRangeArrayv(unsigned int first, int count, const double* v) {  (*pfn_glDepthRangeArrayv)((GLuint)first, (GLsizei)count, (const GLdouble *)v); }
//
// void (APIENTRYP pfn_glDepthRangeIndexed)(GLuint index, GLdouble n, GLdouble f);
// GLAPI void APIENTRY wrap_glDepthRangeIndexed(unsigned int index, double n, double f) {  (*pfn_glDepthRangeIndexed)((GLuint)index, (GLdouble)n, (GLdouble)f); }
//
// void (APIENTRYP pfn_glDepthRangef)(GLfloat n, GLfloat f);
// GLAPI void APIENTRY wrap_glDepthRangef(float n, float f) {  (*pfn_glDepthRangef)((GLfloat)n, (GLfloat)f); }
//
// void (APIENTRYP pfn_glDetachShader)(GLuint program, GLuint shader);
// GLAPI void APIENTRY wrap_glDetachShader(unsigned int program, unsigned int shader) {  (*pfn_glDetachShader)((GLuint)program, (GLuint)shader); }
//
// void (APIENTRYP pfn_glDisable)(GLenum cap);
// GLAPI void APIENTRY wrap_glDisable(unsigned int cap) {  (*pfn_glDisable)((GLenum)cap); }
//
// void (APIENTRYP pfn_glDisableVertexAttribArray)(GLuint index);
// GLAPI void APIENTRY wrap_glDisableVertexAttribArray(unsigned int index) {  (*pfn_glDisableVertexAttribArray)((GLuint)index); }
//
// void (APIENTRYP pfn_glDisablei)(GLenum target, GLuint index);
// GLAPI void APIENTRY wrap_glDisablei(unsigned int target, unsigned int index) {  (*pfn_glDisablei)((GLenum)target, (GLuint)index); }
//
// void (APIENTRYP pfn_glDispatchCompute)(GLuint num_groups_x, GLuint num_groups_y, GLuint num_groups_z);
// GLAPI void APIENTRY wrap_glDispatchCompute(unsigned int num_groups_x, unsigned int num_groups_y, unsigned int num_groups_z) {  (*pfn_glDispatchCompute)((GLuint)num_groups_x, (GLuint)num_groups_y, (GLuint)num_groups_z); }
//
// void (APIENTRYP pfn_glDispatchComputeIndirect)(GLintptr indirect);
// GLAPI void APIENTRY wrap_glDispatchComputeIndirect(long long indirect) {  (*pfn_glDispatchComputeIndirect)((GLintptr)indirect); }
//
// void (APIENTRYP pfn_glDrawArrays)(GLenum mode, GLint first, GLsizei count);
// GLAPI void APIENTRY wrap_glDrawArrays(unsigned int mode, int first, int count) {  (*pfn_glDrawArrays)((GLenum)mode, (GLint)first, (GLsizei)count); }
//
// void (APIENTRYP pfn_glDrawArraysIndirect)(GLenum mode, const GLvoid *indirect);
// GLAPI void APIENTRY wrap_glDrawArraysIndirect(unsigned int mode, const void* indirect) {  (*pfn_glDrawArraysIndirect)((GLenum)mode, (const GLvoid *)indirect); }
//
// void (APIENTRYP pfn_glDrawArraysInstanced)(GLenum mode, GLint first, GLsizei count, GLsizei instancecount);
// GLAPI void APIENTRY wrap_glDrawArraysInstanced(unsigned int mode, int first, int count, int instancecount) {  (*pfn_glDrawArraysInstanced)((GLenum)mode, (GLint)first, (GLsizei)count, (GLsizei)instancecount); }
//
// void (APIENTRYP pfn_glDrawArraysInstancedBaseInstance)(GLenum mode, GLint first, GLsizei count, GLsizei instancecount, GLuint baseinstance);
// GLAPI void APIENTRY wrap_glDrawArraysInstancedBaseInstance(unsigned int mode, int first, int count, int instancecount, unsigned int baseinstance) {  (*pfn_glDrawArraysInstancedBaseInstance)((GLenum)mode, (GLint)first, (GLsizei)count, (GLsizei)instancecount, (GLuint)baseinstance); }
//
// void (APIENTRYP pfn_glDrawBuffer)(GLenum mode);
// GLAPI void APIENTRY wrap_glDrawBuffer(unsigned int mode) {  (*pfn_glDrawBuffer)((GLenum)mode); }
//
// void (APIENTRYP pfn_glDrawBuffers)(GLsizei n, const GLenum *bufs);
// GLAPI void APIENTRY wrap_glDrawBuffers(int n, const unsigned int* bufs) {  (*pfn_glDrawBuffers)((GLsizei)n, (const GLenum *)bufs); }
//
// void (APIENTRYP pfn_glDrawElements)(GLenum mode, GLsizei count, GLenum type, const GLvoid *indices);
// GLAPI void APIENTRY wrap_glDrawElements(unsigned int mode, int count, unsigned int t_ype, long long indicies) {  (*pfn_glDrawElements)((GLenum)mode, (GLsizei)count, (GLenum)t_ype, (const GLvoid *)indicies); }
//
// void (APIENTRYP pfn_glDrawElementsBaseVertex)(GLenum mode, GLsizei count, GLenum type, const GLvoid *indices, GLint basevertex);
// GLAPI void APIENTRY wrap_glDrawElementsBaseVertex(unsigned int mode, int count, unsigned int t_ype, long long indicies, int basevertex) {  (*pfn_glDrawElementsBaseVertex)((GLenum)mode, (GLsizei)count, (GLenum)t_ype, (const GLvoid *)indicies, (GLint)basevertex); }
//
// void (APIENTRYP pfn_glDrawElementsIndirect)(GLenum mode, GLenum type, const GLvoid *indirect);
// GLAPI void APIENTRY wrap_glDrawElementsIndirect(unsigned int mode, unsigned int t_ype, const void* indirect) {  (*pfn_glDrawElementsIndirect)((GLenum)mode, (GLenum)t_ype, (const GLvoid *)indirect); }
//
// void (APIENTRYP pfn_glDrawElementsInstanced)(GLenum mode, GLsizei count, GLenum type, const GLvoid *indices, GLsizei instancecount);
// GLAPI void APIENTRY wrap_glDrawElementsInstanced(unsigned int mode, int count, unsigned int t_ype, long long indicies, int instancecount) {  (*pfn_glDrawElementsInstanced)((GLenum)mode, (GLsizei)count, (GLenum)t_ype, (const GLvoid *)indicies, (GLsizei)instancecount); }
//
// void (APIENTRYP pfn_glDrawElementsInstancedBaseInstance)(GLenum mode, GLsizei count, GLenum type, const void *indices, GLsizei instancecount, GLuint baseinstance);
// GLAPI void APIENTRY wrap_glDrawElementsInstancedBaseInstance(unsigned int mode, int count, unsigned int t_ype, const void* indices, int instancecount, unsigned int baseinstance) {  (*pfn_glDrawElementsInstancedBaseInstance)((GLenum)mode, (GLsizei)count, (GLenum)t_ype, (const void *)indices, (GLsizei)instancecount, (GLuint)baseinstance); }
//
// void (APIENTRYP pfn_glDrawElementsInstancedBaseVertex)(GLenum mode, GLsizei count, GLenum type, const GLvoid *indices, GLsizei instancecount, GLint basevertex);
// GLAPI void APIENTRY wrap_glDrawElementsInstancedBaseVertex(unsigned int mode, int count, unsigned int t_ype, long long indicies, int instancecount, int basevertex) {  (*pfn_glDrawElementsInstancedBaseVertex)((GLenum)mode, (GLsizei)count, (GLenum)t_ype, (const GLvoid *)indicies, (GLsizei)instancecount, (GLint)basevertex); }
//
// void (APIENTRYP pfn_glDrawElementsInstancedBaseVertexBaseInstance)(GLenum mode, GLsizei count, GLenum type, const void *indices, GLsizei instancecount, GLint basevertex, GLuint baseinstance);
// GLAPI void APIENTRY wrap_glDrawElementsInstancedBaseVertexBaseInstance(unsigned int mode, int count, unsigned int t_ype, const void* indices, int instancecount, int basevertex, unsigned int baseinstance) {  (*pfn_glDrawElementsInstancedBaseVertexBaseInstance)((GLenum)mode, (GLsizei)count, (GLenum)t_ype, (const void *)indices, (GLsizei)instancecount, (GLint)basevertex, (GLuint)baseinstance); }
//
// void (APIENTRYP pfn_glDrawRangeElements)(GLenum mode, GLuint start, GLuint end, GLsizei count, GLenum type, const GLvoid *indices);
// GLAPI void APIENTRY wrap_glDrawRangeElements(unsigned int mode, unsigned int start, unsigned int end, int count, unsigned int t_ype, long long indicies) {  (*pfn_glDrawRangeElements)((GLenum)mode, (GLuint)start, (GLuint)end, (GLsizei)count, (GLenum)t_ype, (const GLvoid *)indicies); }
//
// void (APIENTRYP pfn_glDrawRangeElementsBaseVertex)(GLenum mode, GLuint start, GLuint end, GLsizei count, GLenum type, const GLvoid *indices, GLint basevertex);
// GLAPI void APIENTRY wrap_glDrawRangeElementsBaseVertex(unsigned int mode, unsigned int start, unsigned int end, int count, unsigned int t_ype, long long indicies, int basevertex) {  (*pfn_glDrawRangeElementsBaseVertex)((GLenum)mode, (GLuint)start, (GLuint)end, (GLsizei)count, (GLenum)t_ype, (const GLvoid *)indicies, (GLint)basevertex); }
//
// void (APIENTRYP pfn_glDrawTransformFeedback)(GLenum mode, GLuint id);
// GLAPI void APIENTRY wrap_glDrawTransformFeedback(unsigned int mode, unsigned int id) {  (*pfn_glDrawTransformFeedback)((GLenum)mode, (GLuint)id); }
//
// void (APIENTRYP pfn_glDrawTransformFeedbackInstanced)(GLenum mode, GLuint id, GLsizei instancecount);
// GLAPI void APIENTRY wrap_glDrawTransformFeedbackInstanced(unsigned int mode, unsigned int id, int instancecount) {  (*pfn_glDrawTransformFeedbackInstanced)((GLenum)mode, (GLuint)id, (GLsizei)instancecount); }
//
// void (APIENTRYP pfn_glDrawTransformFeedbackStream)(GLenum mode, GLuint id, GLuint stream);
// GLAPI void APIENTRY wrap_glDrawTransformFeedbackStream(unsigned int mode, unsigned int id, unsigned int stream) {  (*pfn_glDrawTransformFeedbackStream)((GLenum)mode, (GLuint)id, (GLuint)stream); }
//
// void (APIENTRYP pfn_glDrawTransformFeedbackStreamInstanced)(GLenum mode, GLuint id, GLuint stream, GLsizei instancecount);
// GLAPI void APIENTRY wrap_glDrawTransformFeedbackStreamInstanced(unsigned int mode, unsigned int id, unsigned int stream, int instancecount) {  (*pfn_glDrawTransformFeedbackStreamInstanced)((GLenum)mode, (GLuint)id, (GLuint)stream, (GLsizei)instancecount); }
//
// void (APIENTRYP pfn_glEnable)(GLenum cap);
// GLAPI void APIENTRY wrap_glEnable(unsigned int cap) {  (*pfn_glEnable)((GLenum)cap); }
//
// void (APIENTRYP pfn_glEnableVertexAttribArray)(GLuint index);
// GLAPI void APIENTRY wrap_glEnableVertexAttribArray(unsigned int index) {  (*pfn_glEnableVertexAttribArray)((GLuint)index); }
//
// void (APIENTRYP pfn_glEnablei)(GLenum target, GLuint index);
// GLAPI void APIENTRY wrap_glEnablei(unsigned int target, unsigned int index) {  (*pfn_glEnablei)((GLenum)target, (GLuint)index); }
//
// void (APIENTRYP pfn_glEndConditionalRender)(void);
// GLAPI void APIENTRY wrap_glEndConditionalRender() {  (*pfn_glEndConditionalRender)(); }
//
// void (APIENTRYP pfn_glEndQuery)(GLenum target);
// GLAPI void APIENTRY wrap_glEndQuery(unsigned int target) {  (*pfn_glEndQuery)((GLenum)target); }
//
// void (APIENTRYP pfn_glEndQueryIndexed)(GLenum target, GLuint index);
// GLAPI void APIENTRY wrap_glEndQueryIndexed(unsigned int target, unsigned int index) {  (*pfn_glEndQueryIndexed)((GLenum)target, (GLuint)index); }
//
// void (APIENTRYP pfn_glEndTransformFeedback)(void);
// GLAPI void APIENTRY wrap_glEndTransformFeedback() {  (*pfn_glEndTransformFeedback)(); }
//
// void (APIENTRYP pfn_glFinish)(void);
// GLAPI void APIENTRY wrap_glFinish() {  (*pfn_glFinish)(); }
//
// void (APIENTRYP pfn_glFlush)(void);
// GLAPI void APIENTRY wrap_glFlush() {  (*pfn_glFlush)(); }
//
// void (APIENTRYP pfn_glFlushMappedBufferRange)(GLenum target, GLintptr offset, GLsizeiptr length);
// GLAPI void APIENTRY wrap_glFlushMappedBufferRange(unsigned int target, long long offset, long long length) {  (*pfn_glFlushMappedBufferRange)((GLenum)target, (GLintptr)offset, (GLsizeiptr)length); }
//
// void (APIENTRYP pfn_glFramebufferParameteri)(GLenum target, GLenum pname, GLint param);
// GLAPI void APIENTRY wrap_glFramebufferParameteri(unsigned int target, unsigned int pname, int param) {  (*pfn_glFramebufferParameteri)((GLenum)target, (GLenum)pname, (GLint)param); }
//
// void (APIENTRYP pfn_glFramebufferRenderbuffer)(GLenum target, GLenum attachment, GLenum renderbuffertarget, GLuint renderbuffer);
// GLAPI void APIENTRY wrap_glFramebufferRenderbuffer(unsigned int target, unsigned int attachment, unsigned int renderbuffertarget, unsigned int renderbuffer) {  (*pfn_glFramebufferRenderbuffer)((GLenum)target, (GLenum)attachment, (GLenum)renderbuffertarget, (GLuint)renderbuffer); }
//
// void (APIENTRYP pfn_glFramebufferTexture)(GLenum target, GLenum attachment, GLuint texture, GLint level);
// GLAPI void APIENTRY wrap_glFramebufferTexture(unsigned int target, unsigned int attachment, unsigned int texture, int level) {  (*pfn_glFramebufferTexture)((GLenum)target, (GLenum)attachment, (GLuint)texture, (GLint)level); }
//
// void (APIENTRYP pfn_glFramebufferTexture1D)(GLenum target, GLenum attachment, GLenum textarget, GLuint texture, GLint level);
// GLAPI void APIENTRY wrap_glFramebufferTexture1D(unsigned int target, unsigned int attachment, unsigned int textarget, unsigned int texture, int level) {  (*pfn_glFramebufferTexture1D)((GLenum)target, (GLenum)attachment, (GLenum)textarget, (GLuint)texture, (GLint)level); }
//
// void (APIENTRYP pfn_glFramebufferTexture2D)(GLenum target, GLenum attachment, GLenum textarget, GLuint texture, GLint level);
// GLAPI void APIENTRY wrap_glFramebufferTexture2D(unsigned int target, unsigned int attachment, unsigned int textarget, unsigned int texture, int level) {  (*pfn_glFramebufferTexture2D)((GLenum)target, (GLenum)attachment, (GLenum)textarget, (GLuint)texture, (GLint)level); }
//
// void (APIENTRYP pfn_glFramebufferTexture3D)(GLenum target, GLenum attachment, GLenum textarget, GLuint texture, GLint level, GLint zoffset);
// GLAPI void APIENTRY wrap_glFramebufferTexture3D(unsigned int target, unsigned int attachment, unsigned int textarget, unsigned int texture, int level, int zoffset) {  (*pfn_glFramebufferTexture3D)((GLenum)target, (GLenum)attachment, (GLenum)textarget, (GLuint)texture, (GLint)level, (GLint)zoffset); }
//
// void (APIENTRYP pfn_glFramebufferTextureLayer)(GLenum target, GLenum attachment, GLuint texture, GLint level, GLint layer);
// GLAPI void APIENTRY wrap_glFramebufferTextureLayer(unsigned int target, unsigned int attachment, unsigned int texture, int level, int layer) {  (*pfn_glFramebufferTextureLayer)((GLenum)target, (GLenum)attachment, (GLuint)texture, (GLint)level, (GLint)layer); }
//
// void (APIENTRYP pfn_glFrontFace)(GLenum mode);
// GLAPI void APIENTRY wrap_glFrontFace(unsigned int mode) {  (*pfn_glFrontFace)((GLenum)mode); }
//
// void (APIENTRYP pfn_glGenBuffers)(GLsizei n, GLuint *buffers);
// GLAPI void APIENTRY wrap_glGenBuffers(int n, unsigned int* buffers) {  (*pfn_glGenBuffers)((GLsizei)n, (GLuint *)buffers); }
//
// void (APIENTRYP pfn_glGenFramebuffers)(GLsizei n, GLuint *framebuffers);
// GLAPI void APIENTRY wrap_glGenFramebuffers(int n, unsigned int* framebuffers) {  (*pfn_glGenFramebuffers)((GLsizei)n, (GLuint *)framebuffers); }
//
// void (APIENTRYP pfn_glGenProgramPipelines)(GLsizei n, GLuint *pipelines);
// GLAPI void APIENTRY wrap_glGenProgramPipelines(int n, unsigned int* pipelines) {  (*pfn_glGenProgramPipelines)((GLsizei)n, (GLuint *)pipelines); }
//
// void (APIENTRYP pfn_glGenQueries)(GLsizei n, GLuint *ids);
// GLAPI void APIENTRY wrap_glGenQueries(int n, unsigned int* ids) {  (*pfn_glGenQueries)((GLsizei)n, (GLuint *)ids); }
//
// void (APIENTRYP pfn_glGenRenderbuffers)(GLsizei n, GLuint *renderbuffers);
// GLAPI void APIENTRY wrap_glGenRenderbuffers(int n, unsigned int* renderbuffers) {  (*pfn_glGenRenderbuffers)((GLsizei)n, (GLuint *)renderbuffers); }
//
// void (APIENTRYP pfn_glGenSamplers)(GLsizei count, GLuint *samplers);
// GLAPI void APIENTRY wrap_glGenSamplers(int count, unsigned int* samplers) {  (*pfn_glGenSamplers)((GLsizei)count, (GLuint *)samplers); }
//
// void (APIENTRYP pfn_glGenTextures)(GLsizei n, GLuint *textures);
// GLAPI void APIENTRY wrap_glGenTextures(int n, unsigned int* textures) {  (*pfn_glGenTextures)((GLsizei)n, (GLuint *)textures); }
//
// void (APIENTRYP pfn_glGenTransformFeedbacks)(GLsizei n, GLuint *ids);
// GLAPI void APIENTRY wrap_glGenTransformFeedbacks(int n, unsigned int* ids) {  (*pfn_glGenTransformFeedbacks)((GLsizei)n, (GLuint *)ids); }
//
// void (APIENTRYP pfn_glGenVertexArrays)(GLsizei n, GLuint *arrays);
// GLAPI void APIENTRY wrap_glGenVertexArrays(int n, unsigned int* arrays) {  (*pfn_glGenVertexArrays)((GLsizei)n, (GLuint *)arrays); }
//
// void (APIENTRYP pfn_glGenerateMipmap)(GLenum target);
// GLAPI void APIENTRY wrap_glGenerateMipmap(unsigned int target) {  (*pfn_glGenerateMipmap)((GLenum)target); }
//
// void (APIENTRYP pfn_glGetActiveAtomicCounterBufferiv)(GLuint program, GLuint bufferIndex, GLenum pname, GLint *params);
// GLAPI void APIENTRY wrap_glGetActiveAtomicCounterBufferiv(unsigned int program, unsigned int bufferIndex, unsigned int pname, int* params) {  (*pfn_glGetActiveAtomicCounterBufferiv)((GLuint)program, (GLuint)bufferIndex, (GLenum)pname, (GLint *)params); }
//
// void (APIENTRYP pfn_glGetActiveAttrib)(GLuint program, GLuint index, GLsizei bufSize, GLsizei *length, GLint *size, GLenum *type, GLchar *name);
// GLAPI void APIENTRY wrap_glGetActiveAttrib(unsigned int program, unsigned int index, int bufSize, int* length, int* size, unsigned int* t_ype, char* name) {  (*pfn_glGetActiveAttrib)((GLuint)program, (GLuint)index, (GLsizei)bufSize, (GLsizei *)length, (GLint *)size, (GLenum *)t_ype, (GLchar *)name); }
//
// void (APIENTRYP pfn_glGetActiveSubroutineName)(GLuint program, GLenum shadertype, GLuint index, GLsizei bufsize, GLsizei *length, GLchar *name);
// GLAPI void APIENTRY wrap_glGetActiveSubroutineName(unsigned int program, unsigned int shadertype, unsigned int index, int bufsize, int* length, unsigned char* name) {  (*pfn_glGetActiveSubroutineName)((GLuint)program, (GLenum)shadertype, (GLuint)index, (GLsizei)bufsize, (GLsizei *)length, (GLchar *)name); }
//
// void (APIENTRYP pfn_glGetActiveSubroutineUniformName)(GLuint program, GLenum shadertype, GLuint index, GLsizei bufsize, GLsizei *length, GLchar *name);
// GLAPI void APIENTRY wrap_glGetActiveSubroutineUniformName(unsigned int program, unsigned int shadertype, unsigned int index, int bufsize, int* length, char* name) {  (*pfn_glGetActiveSubroutineUniformName)((GLuint)program, (GLenum)shadertype, (GLuint)index, (GLsizei)bufsize, (GLsizei *)length, (GLchar *)name); }
//
// void (APIENTRYP pfn_glGetActiveSubroutineUniformiv)(GLuint program, GLenum shadertype, GLuint index, GLenum pname, GLint *values);
// GLAPI void APIENTRY wrap_glGetActiveSubroutineUniformiv(unsigned int program, unsigned int shadertype, unsigned int index, unsigned int pname, int* values) {  (*pfn_glGetActiveSubroutineUniformiv)((GLuint)program, (GLenum)shadertype, (GLuint)index, (GLenum)pname, (GLint *)values); }
//
// void (APIENTRYP pfn_glGetActiveUniform)(GLuint program, GLuint index, GLsizei bufSize, GLsizei *length, GLint *size, GLenum *type, GLchar *name);
// GLAPI void APIENTRY wrap_glGetActiveUniform(unsigned int program, unsigned int index, int bufSize, int* length, int* size, unsigned int* t_ype, unsigned char* name) {  (*pfn_glGetActiveUniform)((GLuint)program, (GLuint)index, (GLsizei)bufSize, (GLsizei *)length, (GLint *)size, (GLenum *)t_ype, (GLchar *)name); }
//
// void (APIENTRYP pfn_glGetActiveUniformBlockName)(GLuint program, GLuint uniformBlockIndex, GLsizei bufSize, GLsizei *length, GLchar *uniformBlockName);
// GLAPI void APIENTRY wrap_glGetActiveUniformBlockName(unsigned int program, unsigned int uniformBlockIndex, int bufSize, int* length, unsigned char* uniformBlockName) {  (*pfn_glGetActiveUniformBlockName)((GLuint)program, (GLuint)uniformBlockIndex, (GLsizei)bufSize, (GLsizei *)length, (GLchar *)uniformBlockName); }
//
// void (APIENTRYP pfn_glGetActiveUniformBlockiv)(GLuint program, GLuint uniformBlockIndex, GLenum pname, GLint *params);
// GLAPI void APIENTRY wrap_glGetActiveUniformBlockiv(unsigned int program, unsigned int uniformBlockIndex, unsigned int pname, int* params) {  (*pfn_glGetActiveUniformBlockiv)((GLuint)program, (GLuint)uniformBlockIndex, (GLenum)pname, (GLint *)params); }
//
// void (APIENTRYP pfn_glGetActiveUniformName)(GLuint program, GLuint uniformIndex, GLsizei bufSize, GLsizei *length, GLchar *uniformName);
// GLAPI void APIENTRY wrap_glGetActiveUniformName(unsigned int program, unsigned int uniformIndex, int bufSize, int* length, unsigned char* uniformName) {  (*pfn_glGetActiveUniformName)((GLuint)program, (GLuint)uniformIndex, (GLsizei)bufSize, (GLsizei *)length, (GLchar *)uniformName); }
//
// void (APIENTRYP pfn_glGetActiveUniformsiv)(GLuint program, GLsizei uniformCount, const GLuint *uniformIndices, GLenum pname, GLint *params);
// GLAPI void APIENTRY wrap_glGetActiveUniformsiv(unsigned int program, int uniformCount, const unsigned int* uniformIndices, unsigned int pname, int* params) {  (*pfn_glGetActiveUniformsiv)((GLuint)program, (GLsizei)uniformCount, (const GLuint *)uniformIndices, (GLenum)pname, (GLint *)params); }
//
// void (APIENTRYP pfn_glGetAttachedShaders)(GLuint program, GLsizei maxCount, GLsizei *count, GLuint *obj);
// GLAPI void APIENTRY wrap_glGetAttachedShaders(unsigned int program, int maxCount, int* count, unsigned int* obj) {  (*pfn_glGetAttachedShaders)((GLuint)program, (GLsizei)maxCount, (GLsizei *)count, (GLuint *)obj); }
//
// void (APIENTRYP pfn_glGetBooleani_v)(GLenum target, GLuint index, GLboolean *data);
// GLAPI void APIENTRY wrap_glGetBooleani_v(unsigned int target, unsigned int index, unsigned char* data) {  (*pfn_glGetBooleani_v)((GLenum)target, (GLuint)index, (GLboolean *)data); }
//
// void (APIENTRYP pfn_glGetBooleanv)(GLenum pname, GLboolean *params);
// GLAPI void APIENTRY wrap_glGetBooleanv(unsigned int pname, unsigned char* params) {  (*pfn_glGetBooleanv)((GLenum)pname, (GLboolean *)params); }
//
// void (APIENTRYP pfn_glGetBufferParameteri64v)(GLenum target, GLenum pname, GLint64 *params);
// GLAPI void APIENTRY wrap_glGetBufferParameteri64v(unsigned int target, unsigned int pname, long long* params) {  (*pfn_glGetBufferParameteri64v)((GLenum)target, (GLenum)pname, (GLint64 *)params); }
//
// void (APIENTRYP pfn_glGetBufferParameteriv)(GLenum target, GLenum pname, GLint *params);
// GLAPI void APIENTRY wrap_glGetBufferParameteriv(unsigned int target, unsigned int pname, int* params) {  (*pfn_glGetBufferParameteriv)((GLenum)target, (GLenum)pname, (GLint *)params); }
//
// void (APIENTRYP pfn_glGetBufferPointerv)(GLenum target, GLenum pname, GLvoid* *params);
// GLAPI void APIENTRY wrap_glGetBufferPointerv(unsigned int target, unsigned int pname, void** params) {  (*pfn_glGetBufferPointerv)((GLenum)target, (GLenum)pname, (GLvoid* *)params); }
//
// void (APIENTRYP pfn_glGetBufferSubData)(GLenum target, GLintptr offset, GLsizeiptr size, GLvoid *data);
// GLAPI void APIENTRY wrap_glGetBufferSubData(unsigned int target, long long offset, long long size, void* data) {  (*pfn_glGetBufferSubData)((GLenum)target, (GLintptr)offset, (GLsizeiptr)size, (GLvoid *)data); }
//
// void (APIENTRYP pfn_glGetCompressedTexImage)(GLenum target, GLint level, GLvoid *img);
// GLAPI void APIENTRY wrap_glGetCompressedTexImage(unsigned int target, int level, void* img) {  (*pfn_glGetCompressedTexImage)((GLenum)target, (GLint)level, (GLvoid *)img); }
//
// void (APIENTRYP pfn_glGetDoublei_v)(GLenum target, GLuint index, GLdouble *data);
// GLAPI void APIENTRY wrap_glGetDoublei_v(unsigned int target, unsigned int index, double* data) {  (*pfn_glGetDoublei_v)((GLenum)target, (GLuint)index, (GLdouble *)data); }
//
// void (APIENTRYP pfn_glGetDoublev)(GLenum pname, GLdouble *params);
// GLAPI void APIENTRY wrap_glGetDoublev(unsigned int pname, double* params) {  (*pfn_glGetDoublev)((GLenum)pname, (GLdouble *)params); }
//
// void (APIENTRYP pfn_glGetFloati_v)(GLenum target, GLuint index, GLfloat *data);
// GLAPI void APIENTRY wrap_glGetFloati_v(unsigned int target, unsigned int index, float* data) {  (*pfn_glGetFloati_v)((GLenum)target, (GLuint)index, (GLfloat *)data); }
//
// void (APIENTRYP pfn_glGetFloatv)(GLenum pname, GLfloat *params);
// GLAPI void APIENTRY wrap_glGetFloatv(unsigned int pname, float* params) {  (*pfn_glGetFloatv)((GLenum)pname, (GLfloat *)params); }
//
// void (APIENTRYP pfn_glGetFramebufferAttachmentParameteriv)(GLenum target, GLenum attachment, GLenum pname, GLint *params);
// GLAPI void APIENTRY wrap_glGetFramebufferAttachmentParameteriv(unsigned int target, unsigned int attachment, unsigned int pname, int* params) {  (*pfn_glGetFramebufferAttachmentParameteriv)((GLenum)target, (GLenum)attachment, (GLenum)pname, (GLint *)params); }
//
// void (APIENTRYP pfn_glGetFramebufferParameteriv)(GLenum target, GLenum pname, GLint *params);
// GLAPI void APIENTRY wrap_glGetFramebufferParameteriv(unsigned int target, unsigned int pname, int* params) {  (*pfn_glGetFramebufferParameteriv)((GLenum)target, (GLenum)pname, (GLint *)params); }
//
// void (APIENTRYP pfn_glGetInteger64i_v)(GLenum target, GLuint index, GLint64 *data);
// GLAPI void APIENTRY wrap_glGetInteger64i_v(unsigned int target, unsigned int index, long long* data) {  (*pfn_glGetInteger64i_v)((GLenum)target, (GLuint)index, (GLint64 *)data); }
//
// void (APIENTRYP pfn_glGetInteger64v)(GLenum pname, GLint64 *params);
// GLAPI void APIENTRY wrap_glGetInteger64v(unsigned int pname, long long* params) {  (*pfn_glGetInteger64v)((GLenum)pname, (GLint64 *)params); }
//
// void (APIENTRYP pfn_glGetIntegeri_v)(GLenum target, GLuint index, GLint *data);
// GLAPI void APIENTRY wrap_glGetIntegeri_v(unsigned int target, unsigned int index, int* data) {  (*pfn_glGetIntegeri_v)((GLenum)target, (GLuint)index, (GLint *)data); }
//
// void (APIENTRYP pfn_glGetIntegerv)(GLenum pname, GLint *params);
// GLAPI void APIENTRY wrap_glGetIntegerv(unsigned int pname, int* params) {  (*pfn_glGetIntegerv)((GLenum)pname, (GLint *)params); }
//
// void (APIENTRYP pfn_glGetInternalformati64v)(GLenum target, GLenum internalformat, GLenum pname, GLsizei bufSize, GLint64 *params);
// GLAPI void APIENTRY wrap_glGetInternalformati64v(unsigned int target, unsigned int internalformat, unsigned int pname, int bufSize, long long* params) {  (*pfn_glGetInternalformati64v)((GLenum)target, (GLenum)internalformat, (GLenum)pname, (GLsizei)bufSize, (GLint64 *)params); }
//
// void (APIENTRYP pfn_glGetInternalformativ)(GLenum target, GLenum internalformat, GLenum pname, GLsizei bufSize, GLint *params);
// GLAPI void APIENTRY wrap_glGetInternalformativ(unsigned int target, unsigned int internalformat, unsigned int pname, int bufSize, int* params) {  (*pfn_glGetInternalformativ)((GLenum)target, (GLenum)internalformat, (GLenum)pname, (GLsizei)bufSize, (GLint *)params); }
//
// void (APIENTRYP pfn_glGetMultisamplefv)(GLenum pname, GLuint index, GLfloat *val);
// GLAPI void APIENTRY wrap_glGetMultisamplefv(unsigned int pname, unsigned int index, float* val) {  (*pfn_glGetMultisamplefv)((GLenum)pname, (GLuint)index, (GLfloat *)val); }
//
// void (APIENTRYP pfn_glGetNamedFramebufferParameterivEXT)(GLuint framebuffer, GLenum pname, GLint *params);
// GLAPI void APIENTRY wrap_glGetNamedFramebufferParameterivEXT(unsigned int framebuffer, unsigned int pname, int* params) {  (*pfn_glGetNamedFramebufferParameterivEXT)((GLuint)framebuffer, (GLenum)pname, (GLint *)params); }
//
// void (APIENTRYP pfn_glGetNamedStringARB)(GLint namelen, const GLchar *name, GLsizei bufSize, GLint *stringlen, GLchar *string);
// GLAPI void APIENTRY wrap_glGetNamedStringARB(int namelen, const char* name, int bufSize, int* stringlen, char* s_tring) {  (*pfn_glGetNamedStringARB)((GLint)namelen, (const GLchar *)name, (GLsizei)bufSize, (GLint *)stringlen, (GLchar *)s_tring); }
//
// void (APIENTRYP pfn_glGetNamedStringivARB)(GLint namelen, const GLchar *name, GLenum pname, GLint *params);
// GLAPI void APIENTRY wrap_glGetNamedStringivARB(int namelen, const char* name, unsigned int pname, int* params) {  (*pfn_glGetNamedStringivARB)((GLint)namelen, (const GLchar *)name, (GLenum)pname, (GLint *)params); }
//
// void (APIENTRYP pfn_glGetObjectLabel)(GLenum identifier, GLuint name, GLsizei bufSize, GLsizei *length, GLchar *label);
// GLAPI void APIENTRY wrap_glGetObjectLabel(unsigned int identifier, unsigned int name, int bufSize, int* length, unsigned char* label) {  (*pfn_glGetObjectLabel)((GLenum)identifier, (GLuint)name, (GLsizei)bufSize, (GLsizei *)length, (GLchar *)label); }
//
// void (APIENTRYP pfn_glGetObjectPtrLabel)(const void *ptr, GLsizei bufSize, GLsizei *length, GLchar *label);
// GLAPI void APIENTRY wrap_glGetObjectPtrLabel(const void* ptr, int bufSize, int* length, unsigned char* label) {  (*pfn_glGetObjectPtrLabel)((const void *)ptr, (GLsizei)bufSize, (GLsizei *)length, (GLchar *)label); }
//
// void (APIENTRYP pfn_glGetPointerv)(GLenum pname, GLvoid* *params);
// GLAPI void APIENTRY wrap_glGetPointerv(unsigned int pname, void** params) {  (*pfn_glGetPointerv)((GLenum)pname, (GLvoid* *)params); }
//
// void (APIENTRYP pfn_glGetProgramBinary)(GLuint program, GLsizei bufSize, GLsizei *length, GLenum *binaryFormat, GLvoid *binary);
// GLAPI void APIENTRY wrap_glGetProgramBinary(unsigned int program, int bufSize, int* length, unsigned int* binaryFormat, void* binary) {  (*pfn_glGetProgramBinary)((GLuint)program, (GLsizei)bufSize, (GLsizei *)length, (GLenum *)binaryFormat, (GLvoid *)binary); }
//
// void (APIENTRYP pfn_glGetProgramInfoLog)(GLuint program, GLsizei bufSize, GLsizei *length, GLchar *infoLog);
// GLAPI void APIENTRY wrap_glGetProgramInfoLog(unsigned int program, int bufSize, int* length, unsigned char* infoLog) {  (*pfn_glGetProgramInfoLog)((GLuint)program, (GLsizei)bufSize, (GLsizei *)length, (GLchar *)infoLog); }
//
// void (APIENTRYP pfn_glGetProgramInterfaceiv)(GLuint program, GLenum programInterface, GLenum pname, GLint *params);
// GLAPI void APIENTRY wrap_glGetProgramInterfaceiv(unsigned int program, unsigned int programInterface, unsigned int pname, int* params) {  (*pfn_glGetProgramInterfaceiv)((GLuint)program, (GLenum)programInterface, (GLenum)pname, (GLint *)params); }
//
// void (APIENTRYP pfn_glGetProgramPipelineInfoLog)(GLuint pipeline, GLsizei bufSize, GLsizei *length, GLchar *infoLog);
// GLAPI void APIENTRY wrap_glGetProgramPipelineInfoLog(unsigned int pipeline, int bufSize, int* length, unsigned char* infoLog) {  (*pfn_glGetProgramPipelineInfoLog)((GLuint)pipeline, (GLsizei)bufSize, (GLsizei *)length, (GLchar *)infoLog); }
//
// void (APIENTRYP pfn_glGetProgramPipelineiv)(GLuint pipeline, GLenum pname, GLint *params);
// GLAPI void APIENTRY wrap_glGetProgramPipelineiv(unsigned int pipeline, unsigned int pname, int* params) {  (*pfn_glGetProgramPipelineiv)((GLuint)pipeline, (GLenum)pname, (GLint *)params); }
//
// void (APIENTRYP pfn_glGetProgramResourceName)(GLuint program, GLenum programInterface, GLuint index, GLsizei bufSize, GLsizei *length, GLchar *name);
// GLAPI void APIENTRY wrap_glGetProgramResourceName(unsigned int program, unsigned int programInterface, unsigned int index, int bufSize, int* length, char* name) {  (*pfn_glGetProgramResourceName)((GLuint)program, (GLenum)programInterface, (GLuint)index, (GLsizei)bufSize, (GLsizei *)length, (GLchar *)name); }
//
// void (APIENTRYP pfn_glGetProgramResourceiv)(GLuint program, GLenum programInterface, GLuint index, GLsizei propCount, const GLenum *props, GLsizei bufSize, GLsizei *length, GLint *params);
// GLAPI void APIENTRY wrap_glGetProgramResourceiv(unsigned int program, unsigned int programInterface, unsigned int index, int propCount, const unsigned int* props, int bufSize, int* length, int* params) {  (*pfn_glGetProgramResourceiv)((GLuint)program, (GLenum)programInterface, (GLuint)index, (GLsizei)propCount, (const GLenum *)props, (GLsizei)bufSize, (GLsizei *)length, (GLint *)params); }
//
// void (APIENTRYP pfn_glGetProgramStageiv)(GLuint program, GLenum shadertype, GLenum pname, GLint *values);
// GLAPI void APIENTRY wrap_glGetProgramStageiv(unsigned int program, unsigned int shadertype, unsigned int pname, int* values) {  (*pfn_glGetProgramStageiv)((GLuint)program, (GLenum)shadertype, (GLenum)pname, (GLint *)values); }
//
// void (APIENTRYP pfn_glGetProgramiv)(GLuint program, GLenum pname, GLint *params);
// GLAPI void APIENTRY wrap_glGetProgramiv(unsigned int program, unsigned int pname, int* params) {  (*pfn_glGetProgramiv)((GLuint)program, (GLenum)pname, (GLint *)params); }
//
// void (APIENTRYP pfn_glGetQueryIndexediv)(GLenum target, GLuint index, GLenum pname, GLint *params);
// GLAPI void APIENTRY wrap_glGetQueryIndexediv(unsigned int target, unsigned int index, unsigned int pname, int* params) {  (*pfn_glGetQueryIndexediv)((GLenum)target, (GLuint)index, (GLenum)pname, (GLint *)params); }
//
// void (APIENTRYP pfn_glGetQueryObjecti64v)(GLuint id, GLenum pname, GLint64 *params);
// GLAPI void APIENTRY wrap_glGetQueryObjecti64v(unsigned int id, unsigned int pname, long long* params) {  (*pfn_glGetQueryObjecti64v)((GLuint)id, (GLenum)pname, (GLint64 *)params); }
//
// void (APIENTRYP pfn_glGetQueryObjectiv)(GLuint id, GLenum pname, GLint *params);
// GLAPI void APIENTRY wrap_glGetQueryObjectiv(unsigned int id, unsigned int pname, int* params) {  (*pfn_glGetQueryObjectiv)((GLuint)id, (GLenum)pname, (GLint *)params); }
//
// void (APIENTRYP pfn_glGetQueryObjectui64v)(GLuint id, GLenum pname, GLuint64 *params);
// GLAPI void APIENTRY wrap_glGetQueryObjectui64v(unsigned int id, unsigned int pname, unsigned long long* params) {  (*pfn_glGetQueryObjectui64v)((GLuint)id, (GLenum)pname, (GLuint64 *)params); }
//
// void (APIENTRYP pfn_glGetQueryObjectuiv)(GLuint id, GLenum pname, GLuint *params);
// GLAPI void APIENTRY wrap_glGetQueryObjectuiv(unsigned int id, unsigned int pname, unsigned int* params) {  (*pfn_glGetQueryObjectuiv)((GLuint)id, (GLenum)pname, (GLuint *)params); }
//
// void (APIENTRYP pfn_glGetQueryiv)(GLenum target, GLenum pname, GLint *params);
// GLAPI void APIENTRY wrap_glGetQueryiv(unsigned int target, unsigned int pname, int* params) {  (*pfn_glGetQueryiv)((GLenum)target, (GLenum)pname, (GLint *)params); }
//
// void (APIENTRYP pfn_glGetRenderbufferParameteriv)(GLenum target, GLenum pname, GLint *params);
// GLAPI void APIENTRY wrap_glGetRenderbufferParameteriv(unsigned int target, unsigned int pname, int* params) {  (*pfn_glGetRenderbufferParameteriv)((GLenum)target, (GLenum)pname, (GLint *)params); }
//
// void (APIENTRYP pfn_glGetSamplerParameterIiv)(GLuint sampler, GLenum pname, GLint *params);
// GLAPI void APIENTRY wrap_glGetSamplerParameterIiv(unsigned int sampler, unsigned int pname, int* params) {  (*pfn_glGetSamplerParameterIiv)((GLuint)sampler, (GLenum)pname, (GLint *)params); }
//
// void (APIENTRYP pfn_glGetSamplerParameterIuiv)(GLuint sampler, GLenum pname, GLuint *params);
// GLAPI void APIENTRY wrap_glGetSamplerParameterIuiv(unsigned int sampler, unsigned int pname, unsigned int* params) {  (*pfn_glGetSamplerParameterIuiv)((GLuint)sampler, (GLenum)pname, (GLuint *)params); }
//
// void (APIENTRYP pfn_glGetSamplerParameterfv)(GLuint sampler, GLenum pname, GLfloat *params);
// GLAPI void APIENTRY wrap_glGetSamplerParameterfv(unsigned int sampler, unsigned int pname, float* params) {  (*pfn_glGetSamplerParameterfv)((GLuint)sampler, (GLenum)pname, (GLfloat *)params); }
//
// void (APIENTRYP pfn_glGetSamplerParameteriv)(GLuint sampler, GLenum pname, GLint *params);
// GLAPI void APIENTRY wrap_glGetSamplerParameteriv(unsigned int sampler, unsigned int pname, int* params) {  (*pfn_glGetSamplerParameteriv)((GLuint)sampler, (GLenum)pname, (GLint *)params); }
//
// void (APIENTRYP pfn_glGetShaderInfoLog)(GLuint shader, GLsizei bufSize, GLsizei *length, GLchar *infoLog);
// GLAPI void APIENTRY wrap_glGetShaderInfoLog(unsigned int shader, int bufSize, int* length, unsigned char* infoLog) {  (*pfn_glGetShaderInfoLog)((GLuint)shader, (GLsizei)bufSize, (GLsizei *)length, (GLchar *)infoLog); }
//
// void (APIENTRYP pfn_glGetShaderPrecisionFormat)(GLenum shadertype, GLenum precisiontype, GLint *range, GLint *precision);
// GLAPI void APIENTRY wrap_glGetShaderPrecisionFormat(unsigned int shadertype, unsigned int precisiontype, int* r_ange, int* precision) {  (*pfn_glGetShaderPrecisionFormat)((GLenum)shadertype, (GLenum)precisiontype, (GLint *)r_ange, (GLint *)precision); }
//
// void (APIENTRYP pfn_glGetShaderSource)(GLuint shader, GLsizei bufSize, GLsizei *length, GLchar *source);
// GLAPI void APIENTRY wrap_glGetShaderSource(unsigned int shader, int bufSize, int* length, unsigned char* source) {  (*pfn_glGetShaderSource)((GLuint)shader, (GLsizei)bufSize, (GLsizei *)length, (GLchar *)source); }
//
// void (APIENTRYP pfn_glGetShaderiv)(GLuint shader, GLenum pname, GLint *params);
// GLAPI void APIENTRY wrap_glGetShaderiv(unsigned int shader, unsigned int pname, int* params) {  (*pfn_glGetShaderiv)((GLuint)shader, (GLenum)pname, (GLint *)params); }
//
// void (APIENTRYP pfn_glGetSynciv)(GLsync sync, GLenum pname, GLsizei bufSize, GLsizei *length, GLint *values);
// GLAPI void APIENTRY wrap_glGetSynciv(GLsync sync, unsigned int pname, int bufSize, int* length, int* values) {  (*pfn_glGetSynciv)((GLsync)sync, (GLenum)pname, (GLsizei)bufSize, (GLsizei *)length, (GLint *)values); }
//
// void (APIENTRYP pfn_glGetTexImage)(GLenum target, GLint level, GLenum format, GLenum type, GLvoid *pixels);
// GLAPI void APIENTRY wrap_glGetTexImage(unsigned int target, int level, unsigned int format, unsigned int t_ype, void* pixels) {  (*pfn_glGetTexImage)((GLenum)target, (GLint)level, (GLenum)format, (GLenum)t_ype, (GLvoid *)pixels); }
//
// void (APIENTRYP pfn_glGetTexLevelParameterfv)(GLenum target, GLint level, GLenum pname, GLfloat *params);
// GLAPI void APIENTRY wrap_glGetTexLevelParameterfv(unsigned int target, int level, unsigned int pname, float* params) {  (*pfn_glGetTexLevelParameterfv)((GLenum)target, (GLint)level, (GLenum)pname, (GLfloat *)params); }
//
// void (APIENTRYP pfn_glGetTexLevelParameteriv)(GLenum target, GLint level, GLenum pname, GLint *params);
// GLAPI void APIENTRY wrap_glGetTexLevelParameteriv(unsigned int target, int level, unsigned int pname, int* params) {  (*pfn_glGetTexLevelParameteriv)((GLenum)target, (GLint)level, (GLenum)pname, (GLint *)params); }
//
// void (APIENTRYP pfn_glGetTexParameterIiv)(GLenum target, GLenum pname, GLint *params);
// GLAPI void APIENTRY wrap_glGetTexParameterIiv(unsigned int target, unsigned int pname, int* params) {  (*pfn_glGetTexParameterIiv)((GLenum)target, (GLenum)pname, (GLint *)params); }
//
// void (APIENTRYP pfn_glGetTexParameterIuiv)(GLenum target, GLenum pname, GLuint *params);
// GLAPI void APIENTRY wrap_glGetTexParameterIuiv(unsigned int target, unsigned int pname, unsigned int* params) {  (*pfn_glGetTexParameterIuiv)((GLenum)target, (GLenum)pname, (GLuint *)params); }
//
// void (APIENTRYP pfn_glGetTexParameterfv)(GLenum target, GLenum pname, GLfloat *params);
// GLAPI void APIENTRY wrap_glGetTexParameterfv(unsigned int target, unsigned int pname, float* params) {  (*pfn_glGetTexParameterfv)((GLenum)target, (GLenum)pname, (GLfloat *)params); }
//
// void (APIENTRYP pfn_glGetTexParameteriv)(GLenum target, GLenum pname, GLint *params);
// GLAPI void APIENTRY wrap_glGetTexParameteriv(unsigned int target, unsigned int pname, int* params) {  (*pfn_glGetTexParameteriv)((GLenum)target, (GLenum)pname, (GLint *)params); }
//
// void (APIENTRYP pfn_glGetTransformFeedbackVarying)(GLuint program, GLuint index, GLsizei bufSize, GLsizei *length, GLsizei *size, GLenum *type, GLchar *name);
// GLAPI void APIENTRY wrap_glGetTransformFeedbackVarying(unsigned int program, unsigned int index, int bufSize, int* length, int* size, unsigned int* t_ype, char* name) {  (*pfn_glGetTransformFeedbackVarying)((GLuint)program, (GLuint)index, (GLsizei)bufSize, (GLsizei *)length, (GLsizei *)size, (GLenum *)t_ype, (GLchar *)name); }
//
// void (APIENTRYP pfn_glGetUniformIndices)(GLuint program, GLsizei uniformCount, const GLchar* const *uniformNames, GLuint *uniformIndices);
// GLAPI void APIENTRY wrap_glGetUniformIndices(unsigned int program, int uniformCount, const char* const* uniformNames, unsigned int* uniformIndices) {  (*pfn_glGetUniformIndices)((GLuint)program, (GLsizei)uniformCount, (const GLchar* const *)uniformNames, (GLuint *)uniformIndices); }
//
// void (APIENTRYP pfn_glGetUniformSubroutineuiv)(GLenum shadertype, GLint location, GLuint *params);
// GLAPI void APIENTRY wrap_glGetUniformSubroutineuiv(unsigned int shadertype, int location, unsigned int* params) {  (*pfn_glGetUniformSubroutineuiv)((GLenum)shadertype, (GLint)location, (GLuint *)params); }
//
// void (APIENTRYP pfn_glGetUniformdv)(GLuint program, GLint location, GLdouble *params);
// GLAPI void APIENTRY wrap_glGetUniformdv(unsigned int program, int location, double* params) {  (*pfn_glGetUniformdv)((GLuint)program, (GLint)location, (GLdouble *)params); }
//
// void (APIENTRYP pfn_glGetUniformfv)(GLuint program, GLint location, GLfloat *params);
// GLAPI void APIENTRY wrap_glGetUniformfv(unsigned int program, int location, float* params) {  (*pfn_glGetUniformfv)((GLuint)program, (GLint)location, (GLfloat *)params); }
//
// void (APIENTRYP pfn_glGetUniformiv)(GLuint program, GLint location, GLint *params);
// GLAPI void APIENTRY wrap_glGetUniformiv(unsigned int program, int location, int* params) {  (*pfn_glGetUniformiv)((GLuint)program, (GLint)location, (GLint *)params); }
//
// void (APIENTRYP pfn_glGetUniformuiv)(GLuint program, GLint location, GLuint *params);
// GLAPI void APIENTRY wrap_glGetUniformuiv(unsigned int program, int location, unsigned int* params) {  (*pfn_glGetUniformuiv)((GLuint)program, (GLint)location, (GLuint *)params); }
//
// void (APIENTRYP pfn_glGetVertexAttribIiv)(GLuint index, GLenum pname, GLint *params);
// GLAPI void APIENTRY wrap_glGetVertexAttribIiv(unsigned int index, unsigned int pname, int* params) {  (*pfn_glGetVertexAttribIiv)((GLuint)index, (GLenum)pname, (GLint *)params); }
//
// void (APIENTRYP pfn_glGetVertexAttribIuiv)(GLuint index, GLenum pname, GLuint *params);
// GLAPI void APIENTRY wrap_glGetVertexAttribIuiv(unsigned int index, unsigned int pname, unsigned int* params) {  (*pfn_glGetVertexAttribIuiv)((GLuint)index, (GLenum)pname, (GLuint *)params); }
//
// void (APIENTRYP pfn_glGetVertexAttribLdv)(GLuint index, GLenum pname, GLdouble *params);
// GLAPI void APIENTRY wrap_glGetVertexAttribLdv(unsigned int index, unsigned int pname, double* params) {  (*pfn_glGetVertexAttribLdv)((GLuint)index, (GLenum)pname, (GLdouble *)params); }
//
// void (APIENTRYP pfn_glGetVertexAttribPointerv)(GLuint index, GLenum pname, GLvoid* *pointer);
// GLAPI void APIENTRY wrap_glGetVertexAttribPointerv(unsigned int index, unsigned int pname, long long* pointer) {  (*pfn_glGetVertexAttribPointerv)((GLuint)index, (GLenum)pname, (GLvoid* *)pointer); }
//
// void (APIENTRYP pfn_glGetVertexAttribdv)(GLuint index, GLenum pname, GLdouble *params);
// GLAPI void APIENTRY wrap_glGetVertexAttribdv(unsigned int index, unsigned int pname, double* params) {  (*pfn_glGetVertexAttribdv)((GLuint)index, (GLenum)pname, (GLdouble *)params); }
//
// void (APIENTRYP pfn_glGetVertexAttribfv)(GLuint index, GLenum pname, GLfloat *params);
// GLAPI void APIENTRY wrap_glGetVertexAttribfv(unsigned int index, unsigned int pname, float* params) {  (*pfn_glGetVertexAttribfv)((GLuint)index, (GLenum)pname, (GLfloat *)params); }
//
// void (APIENTRYP pfn_glGetVertexAttribiv)(GLuint index, GLenum pname, GLint *params);
// GLAPI void APIENTRY wrap_glGetVertexAttribiv(unsigned int index, unsigned int pname, int* params) {  (*pfn_glGetVertexAttribiv)((GLuint)index, (GLenum)pname, (GLint *)params); }
//
// void (APIENTRYP pfn_glGetnColorTableARB)(GLenum target, GLenum format, GLenum type, GLsizei bufSize, GLvoid *table);
// GLAPI void APIENTRY wrap_glGetnColorTableARB(unsigned int target, unsigned int format, unsigned int t_ype, int bufSize, void* table) {  (*pfn_glGetnColorTableARB)((GLenum)target, (GLenum)format, (GLenum)t_ype, (GLsizei)bufSize, (GLvoid *)table); }
//
// void (APIENTRYP pfn_glGetnCompressedTexImageARB)(GLenum target, GLint lod, GLsizei bufSize, GLvoid *img);
// GLAPI void APIENTRY wrap_glGetnCompressedTexImageARB(unsigned int target, int lod, int bufSize, void* img) {  (*pfn_glGetnCompressedTexImageARB)((GLenum)target, (GLint)lod, (GLsizei)bufSize, (GLvoid *)img); }
//
// void (APIENTRYP pfn_glGetnConvolutionFilterARB)(GLenum target, GLenum format, GLenum type, GLsizei bufSize, GLvoid *image);
// GLAPI void APIENTRY wrap_glGetnConvolutionFilterARB(unsigned int target, unsigned int format, unsigned int t_ype, int bufSize, void* image) {  (*pfn_glGetnConvolutionFilterARB)((GLenum)target, (GLenum)format, (GLenum)t_ype, (GLsizei)bufSize, (GLvoid *)image); }
//
// void (APIENTRYP pfn_glGetnHistogramARB)(GLenum target, GLboolean reset, GLenum format, GLenum type, GLsizei bufSize, GLvoid *values);
// GLAPI void APIENTRY wrap_glGetnHistogramARB(unsigned int target, unsigned char reset, unsigned int format, unsigned int t_ype, int bufSize, void* values) {  (*pfn_glGetnHistogramARB)((GLenum)target, (GLboolean)reset, (GLenum)format, (GLenum)t_ype, (GLsizei)bufSize, (GLvoid *)values); }
//
// void (APIENTRYP pfn_glGetnMapdvARB)(GLenum target, GLenum query, GLsizei bufSize, GLdouble *v);
// GLAPI void APIENTRY wrap_glGetnMapdvARB(unsigned int target, unsigned int query, int bufSize, double* v) {  (*pfn_glGetnMapdvARB)((GLenum)target, (GLenum)query, (GLsizei)bufSize, (GLdouble *)v); }
//
// void (APIENTRYP pfn_glGetnMapfvARB)(GLenum target, GLenum query, GLsizei bufSize, GLfloat *v);
// GLAPI void APIENTRY wrap_glGetnMapfvARB(unsigned int target, unsigned int query, int bufSize, float* v) {  (*pfn_glGetnMapfvARB)((GLenum)target, (GLenum)query, (GLsizei)bufSize, (GLfloat *)v); }
//
// void (APIENTRYP pfn_glGetnMapivARB)(GLenum target, GLenum query, GLsizei bufSize, GLint *v);
// GLAPI void APIENTRY wrap_glGetnMapivARB(unsigned int target, unsigned int query, int bufSize, int* v) {  (*pfn_glGetnMapivARB)((GLenum)target, (GLenum)query, (GLsizei)bufSize, (GLint *)v); }
//
// void (APIENTRYP pfn_glGetnMinmaxARB)(GLenum target, GLboolean reset, GLenum format, GLenum type, GLsizei bufSize, GLvoid *values);
// GLAPI void APIENTRY wrap_glGetnMinmaxARB(unsigned int target, unsigned char reset, unsigned int format, unsigned int t_ype, int bufSize, void* values) {  (*pfn_glGetnMinmaxARB)((GLenum)target, (GLboolean)reset, (GLenum)format, (GLenum)t_ype, (GLsizei)bufSize, (GLvoid *)values); }
//
// void (APIENTRYP pfn_glGetnPixelMapfvARB)(GLenum map, GLsizei bufSize, GLfloat *values);
// GLAPI void APIENTRY wrap_glGetnPixelMapfvARB(unsigned int m_ap, int bufSize, float* values) {  (*pfn_glGetnPixelMapfvARB)((GLenum)m_ap, (GLsizei)bufSize, (GLfloat *)values); }
//
// void (APIENTRYP pfn_glGetnPixelMapuivARB)(GLenum map, GLsizei bufSize, GLuint *values);
// GLAPI void APIENTRY wrap_glGetnPixelMapuivARB(unsigned int m_ap, int bufSize, unsigned int* values) {  (*pfn_glGetnPixelMapuivARB)((GLenum)m_ap, (GLsizei)bufSize, (GLuint *)values); }
//
// void (APIENTRYP pfn_glGetnPixelMapusvARB)(GLenum map, GLsizei bufSize, GLushort *values);
// GLAPI void APIENTRY wrap_glGetnPixelMapusvARB(unsigned int m_ap, int bufSize, unsigned short* values) {  (*pfn_glGetnPixelMapusvARB)((GLenum)m_ap, (GLsizei)bufSize, (GLushort *)values); }
//
// void (APIENTRYP pfn_glGetnPolygonStippleARB)(GLsizei bufSize, GLubyte *pattern);
// GLAPI void APIENTRY wrap_glGetnPolygonStippleARB(int bufSize, unsigned char* pattern) {  (*pfn_glGetnPolygonStippleARB)((GLsizei)bufSize, (GLubyte *)pattern); }
//
// void (APIENTRYP pfn_glGetnSeparableFilterARB)(GLenum target, GLenum format, GLenum type, GLsizei rowBufSize, GLvoid *row, GLsizei columnBufSize, GLvoid *column, GLvoid *span);
// GLAPI void APIENTRY wrap_glGetnSeparableFilterARB(unsigned int target, unsigned int format, unsigned int t_ype, int rowBufSize, void* row, int columnBufSize, void* column, void* span) {  (*pfn_glGetnSeparableFilterARB)((GLenum)target, (GLenum)format, (GLenum)t_ype, (GLsizei)rowBufSize, (GLvoid *)row, (GLsizei)columnBufSize, (GLvoid *)column, (GLvoid *)span); }
//
// void (APIENTRYP pfn_glGetnTexImageARB)(GLenum target, GLint level, GLenum format, GLenum type, GLsizei bufSize, GLvoid *img);
// GLAPI void APIENTRY wrap_glGetnTexImageARB(unsigned int target, int level, unsigned int format, unsigned int t_ype, int bufSize, void* img) {  (*pfn_glGetnTexImageARB)((GLenum)target, (GLint)level, (GLenum)format, (GLenum)t_ype, (GLsizei)bufSize, (GLvoid *)img); }
//
// void (APIENTRYP pfn_glGetnUniformdvARB)(GLuint program, GLint location, GLsizei bufSize, GLdouble *params);
// GLAPI void APIENTRY wrap_glGetnUniformdvARB(unsigned int program, int location, int bufSize, double* params) {  (*pfn_glGetnUniformdvARB)((GLuint)program, (GLint)location, (GLsizei)bufSize, (GLdouble *)params); }
//
// void (APIENTRYP pfn_glGetnUniformfvARB)(GLuint program, GLint location, GLsizei bufSize, GLfloat *params);
// GLAPI void APIENTRY wrap_glGetnUniformfvARB(unsigned int program, int location, int bufSize, float* params) {  (*pfn_glGetnUniformfvARB)((GLuint)program, (GLint)location, (GLsizei)bufSize, (GLfloat *)params); }
//
// void (APIENTRYP pfn_glGetnUniformivARB)(GLuint program, GLint location, GLsizei bufSize, GLint *params);
// GLAPI void APIENTRY wrap_glGetnUniformivARB(unsigned int program, int location, int bufSize, int* params) {  (*pfn_glGetnUniformivARB)((GLuint)program, (GLint)location, (GLsizei)bufSize, (GLint *)params); }
//
// void (APIENTRYP pfn_glGetnUniformuivARB)(GLuint program, GLint location, GLsizei bufSize, GLuint *params);
// GLAPI void APIENTRY wrap_glGetnUniformuivARB(unsigned int program, int location, int bufSize, unsigned int* params) {  (*pfn_glGetnUniformuivARB)((GLuint)program, (GLint)location, (GLsizei)bufSize, (GLuint *)params); }
//
// void (APIENTRYP pfn_glHint)(GLenum target, GLenum mode);
// GLAPI void APIENTRY wrap_glHint(unsigned int target, unsigned int mode) {  (*pfn_glHint)((GLenum)target, (GLenum)mode); }
//
// void (APIENTRYP pfn_glInvalidateBufferData)(GLuint buffer);
// GLAPI void APIENTRY wrap_glInvalidateBufferData(unsigned int buffer) {  (*pfn_glInvalidateBufferData)((GLuint)buffer); }
//
// void (APIENTRYP pfn_glInvalidateBufferSubData)(GLuint buffer, GLintptr offset, GLsizeiptr length);
// GLAPI void APIENTRY wrap_glInvalidateBufferSubData(unsigned int buffer, long long offset, long long length) {  (*pfn_glInvalidateBufferSubData)((GLuint)buffer, (GLintptr)offset, (GLsizeiptr)length); }
//
// void (APIENTRYP pfn_glInvalidateFramebuffer)(GLenum target, GLsizei numAttachments, const GLenum *attachments);
// GLAPI void APIENTRY wrap_glInvalidateFramebuffer(unsigned int target, int numAttachments, const unsigned int* attachments) {  (*pfn_glInvalidateFramebuffer)((GLenum)target, (GLsizei)numAttachments, (const GLenum *)attachments); }
//
// void (APIENTRYP pfn_glInvalidateSubFramebuffer)(GLenum target, GLsizei numAttachments, const GLenum *attachments, GLint x, GLint y, GLsizei width, GLsizei height);
// GLAPI void APIENTRY wrap_glInvalidateSubFramebuffer(unsigned int target, int numAttachments, const unsigned int* attachments, int x, int y, int width, int height) {  (*pfn_glInvalidateSubFramebuffer)((GLenum)target, (GLsizei)numAttachments, (const GLenum *)attachments, (GLint)x, (GLint)y, (GLsizei)width, (GLsizei)height); }
//
// void (APIENTRYP pfn_glInvalidateTexImage)(GLuint texture, GLint level);
// GLAPI void APIENTRY wrap_glInvalidateTexImage(unsigned int texture, int level) {  (*pfn_glInvalidateTexImage)((GLuint)texture, (GLint)level); }
//
// void (APIENTRYP pfn_glInvalidateTexSubImage)(GLuint texture, GLint level, GLint xoffset, GLint yoffset, GLint zoffset, GLsizei width, GLsizei height, GLsizei depth);
// GLAPI void APIENTRY wrap_glInvalidateTexSubImage(unsigned int texture, int level, int xoffset, int yoffset, int zoffset, int width, int height, int depth) {  (*pfn_glInvalidateTexSubImage)((GLuint)texture, (GLint)level, (GLint)xoffset, (GLint)yoffset, (GLint)zoffset, (GLsizei)width, (GLsizei)height, (GLsizei)depth); }
//
// void (APIENTRYP pfn_glLineWidth)(GLfloat width);
// GLAPI void APIENTRY wrap_glLineWidth(float width) {  (*pfn_glLineWidth)((GLfloat)width); }
//
// void (APIENTRYP pfn_glLinkProgram)(GLuint program);
// GLAPI void APIENTRY wrap_glLinkProgram(unsigned int program) {  (*pfn_glLinkProgram)((GLuint)program); }
//
// void (APIENTRYP pfn_glLogicOp)(GLenum opcode);
// GLAPI void APIENTRY wrap_glLogicOp(unsigned int opcode) {  (*pfn_glLogicOp)((GLenum)opcode); }
//
// void (APIENTRYP pfn_glMemoryBarrier)(GLbitfield barriers);
// GLAPI void APIENTRY wrap_glMemoryBarrier(unsigned int barriers) {  (*pfn_glMemoryBarrier)((GLbitfield)barriers); }
//
// void (APIENTRYP pfn_glMinSampleShading)(GLfloat value);
// GLAPI void APIENTRY wrap_glMinSampleShading(float value) {  (*pfn_glMinSampleShading)((GLfloat)value); }
//
// void (APIENTRYP pfn_glMinSampleShadingARB)(GLfloat value);
// GLAPI void APIENTRY wrap_glMinSampleShadingARB(float value) {  (*pfn_glMinSampleShadingARB)((GLfloat)value); }
//
// void (APIENTRYP pfn_glMultiDrawArrays)(GLenum mode, const GLint *first, const GLsizei *count, GLsizei drawcount);
// GLAPI void APIENTRY wrap_glMultiDrawArrays(unsigned int mode, const int* first, const int* count, int drawcount) {  (*pfn_glMultiDrawArrays)((GLenum)mode, (const GLint *)first, (const GLsizei *)count, (GLsizei)drawcount); }
//
// void (APIENTRYP pfn_glMultiDrawArraysIndirect)(GLenum mode, const void *indirect, GLsizei drawcount, GLsizei stride);
// GLAPI void APIENTRY wrap_glMultiDrawArraysIndirect(unsigned int mode, const void* indirect, int drawcount, int stride) {  (*pfn_glMultiDrawArraysIndirect)((GLenum)mode, (const void *)indirect, (GLsizei)drawcount, (GLsizei)stride); }
//
// void (APIENTRYP pfn_glMultiDrawElements)(GLenum mode, const GLsizei *count, GLenum type, const GLvoid* const *indices, GLsizei drawcount);
// GLAPI void APIENTRY wrap_glMultiDrawElements(unsigned int mode, const int* count, unsigned int t_ype, const void* const* indices, int drawcount) {  (*pfn_glMultiDrawElements)((GLenum)mode, (const GLsizei *)count, (GLenum)t_ype, (const GLvoid* const *)indices, (GLsizei)drawcount); }
//
// void (APIENTRYP pfn_glMultiDrawElementsBaseVertex)(GLenum mode, const GLsizei *count, GLenum type, const GLvoid* const *indices, GLsizei drawcount, const GLint *basevertex);
// GLAPI void APIENTRY wrap_glMultiDrawElementsBaseVertex(unsigned int mode, const int* count, unsigned int t_ype, const void* const* indices, int drawcount, const int* basevertex) {  (*pfn_glMultiDrawElementsBaseVertex)((GLenum)mode, (const GLsizei *)count, (GLenum)t_ype, (const GLvoid* const *)indices, (GLsizei)drawcount, (const GLint *)basevertex); }
//
// void (APIENTRYP pfn_glMultiDrawElementsIndirect)(GLenum mode, GLenum type, const void *indirect, GLsizei drawcount, GLsizei stride);
// GLAPI void APIENTRY wrap_glMultiDrawElementsIndirect(unsigned int mode, unsigned int t_ype, const void* indirect, int drawcount, int stride) {  (*pfn_glMultiDrawElementsIndirect)((GLenum)mode, (GLenum)t_ype, (const void *)indirect, (GLsizei)drawcount, (GLsizei)stride); }
//
// void (APIENTRYP pfn_glMultiTexCoordP1ui)(GLenum texture, GLenum type, GLuint coords);
// GLAPI void APIENTRY wrap_glMultiTexCoordP1ui(unsigned int texture, unsigned int t_ype, unsigned int coords) {  (*pfn_glMultiTexCoordP1ui)((GLenum)texture, (GLenum)t_ype, (GLuint)coords); }
//
// void (APIENTRYP pfn_glMultiTexCoordP1uiv)(GLenum texture, GLenum type, const GLuint *coords);
// GLAPI void APIENTRY wrap_glMultiTexCoordP1uiv(unsigned int texture, unsigned int t_ype, const unsigned int* coords) {  (*pfn_glMultiTexCoordP1uiv)((GLenum)texture, (GLenum)t_ype, (const GLuint *)coords); }
//
// void (APIENTRYP pfn_glMultiTexCoordP2ui)(GLenum texture, GLenum type, GLuint coords);
// GLAPI void APIENTRY wrap_glMultiTexCoordP2ui(unsigned int texture, unsigned int t_ype, unsigned int coords) {  (*pfn_glMultiTexCoordP2ui)((GLenum)texture, (GLenum)t_ype, (GLuint)coords); }
//
// void (APIENTRYP pfn_glMultiTexCoordP2uiv)(GLenum texture, GLenum type, const GLuint *coords);
// GLAPI void APIENTRY wrap_glMultiTexCoordP2uiv(unsigned int texture, unsigned int t_ype, const unsigned int* coords) {  (*pfn_glMultiTexCoordP2uiv)((GLenum)texture, (GLenum)t_ype, (const GLuint *)coords); }
//
// void (APIENTRYP pfn_glMultiTexCoordP3ui)(GLenum texture, GLenum type, GLuint coords);
// GLAPI void APIENTRY wrap_glMultiTexCoordP3ui(unsigned int texture, unsigned int t_ype, unsigned int coords) {  (*pfn_glMultiTexCoordP3ui)((GLenum)texture, (GLenum)t_ype, (GLuint)coords); }
//
// void (APIENTRYP pfn_glMultiTexCoordP3uiv)(GLenum texture, GLenum type, const GLuint *coords);
// GLAPI void APIENTRY wrap_glMultiTexCoordP3uiv(unsigned int texture, unsigned int t_ype, const unsigned int* coords) {  (*pfn_glMultiTexCoordP3uiv)((GLenum)texture, (GLenum)t_ype, (const GLuint *)coords); }
//
// void (APIENTRYP pfn_glMultiTexCoordP4ui)(GLenum texture, GLenum type, GLuint coords);
// GLAPI void APIENTRY wrap_glMultiTexCoordP4ui(unsigned int texture, unsigned int t_ype, unsigned int coords) {  (*pfn_glMultiTexCoordP4ui)((GLenum)texture, (GLenum)t_ype, (GLuint)coords); }
//
// void (APIENTRYP pfn_glMultiTexCoordP4uiv)(GLenum texture, GLenum type, const GLuint *coords);
// GLAPI void APIENTRY wrap_glMultiTexCoordP4uiv(unsigned int texture, unsigned int t_ype, const unsigned int* coords) {  (*pfn_glMultiTexCoordP4uiv)((GLenum)texture, (GLenum)t_ype, (const GLuint *)coords); }
//
// void (APIENTRYP pfn_glNamedFramebufferParameteriEXT)(GLuint framebuffer, GLenum pname, GLint param);
// GLAPI void APIENTRY wrap_glNamedFramebufferParameteriEXT(unsigned int framebuffer, unsigned int pname, int param) {  (*pfn_glNamedFramebufferParameteriEXT)((GLuint)framebuffer, (GLenum)pname, (GLint)param); }
//
// void (APIENTRYP pfn_glNamedStringARB)(GLenum type, GLint namelen, const GLchar *name, GLint stringlen, const GLchar *string);
// GLAPI void APIENTRY wrap_glNamedStringARB(unsigned int t_ype, int namelen, const char* name, int stringlen, const char* s_tring) {  (*pfn_glNamedStringARB)((GLenum)t_ype, (GLint)namelen, (const GLchar *)name, (GLint)stringlen, (const GLchar *)s_tring); }
//
// void (APIENTRYP pfn_glNormalP3ui)(GLenum type, GLuint coords);
// GLAPI void APIENTRY wrap_glNormalP3ui(unsigned int t_ype, unsigned int coords) {  (*pfn_glNormalP3ui)((GLenum)t_ype, (GLuint)coords); }
//
// void (APIENTRYP pfn_glNormalP3uiv)(GLenum type, const GLuint *coords);
// GLAPI void APIENTRY wrap_glNormalP3uiv(unsigned int t_ype, const unsigned int* coords) {  (*pfn_glNormalP3uiv)((GLenum)t_ype, (const GLuint *)coords); }
//
// void (APIENTRYP pfn_glObjectLabel)(GLenum identifier, GLuint name, GLsizei length, const GLchar *label);
// GLAPI void APIENTRY wrap_glObjectLabel(unsigned int identifier, unsigned int name, int length, const char* label) {  (*pfn_glObjectLabel)((GLenum)identifier, (GLuint)name, (GLsizei)length, (const GLchar *)label); }
//
// void (APIENTRYP pfn_glObjectPtrLabel)(const void *ptr, GLsizei length, const GLchar *label);
// GLAPI void APIENTRY wrap_glObjectPtrLabel(const void* ptr, int length, const char* label) {  (*pfn_glObjectPtrLabel)((const void *)ptr, (GLsizei)length, (const GLchar *)label); }
//
// void (APIENTRYP pfn_glPatchParameterfv)(GLenum pname, const GLfloat *values);
// GLAPI void APIENTRY wrap_glPatchParameterfv(unsigned int pname, const float* values) {  (*pfn_glPatchParameterfv)((GLenum)pname, (const GLfloat *)values); }
//
// void (APIENTRYP pfn_glPatchParameteri)(GLenum pname, GLint value);
// GLAPI void APIENTRY wrap_glPatchParameteri(unsigned int pname, int value) {  (*pfn_glPatchParameteri)((GLenum)pname, (GLint)value); }
//
// void (APIENTRYP pfn_glPauseTransformFeedback)(void);
// GLAPI void APIENTRY wrap_glPauseTransformFeedback() {  (*pfn_glPauseTransformFeedback)(); }
//
// void (APIENTRYP pfn_glPixelStoref)(GLenum pname, GLfloat param);
// GLAPI void APIENTRY wrap_glPixelStoref(unsigned int pname, float param) {  (*pfn_glPixelStoref)((GLenum)pname, (GLfloat)param); }
//
// void (APIENTRYP pfn_glPixelStorei)(GLenum pname, GLint param);
// GLAPI void APIENTRY wrap_glPixelStorei(unsigned int pname, int param) {  (*pfn_glPixelStorei)((GLenum)pname, (GLint)param); }
//
// void (APIENTRYP pfn_glPointParameterf)(GLenum pname, GLfloat param);
// GLAPI void APIENTRY wrap_glPointParameterf(unsigned int pname, float param) {  (*pfn_glPointParameterf)((GLenum)pname, (GLfloat)param); }
//
// void (APIENTRYP pfn_glPointParameterfv)(GLenum pname, const GLfloat *params);
// GLAPI void APIENTRY wrap_glPointParameterfv(unsigned int pname, const float* params) {  (*pfn_glPointParameterfv)((GLenum)pname, (const GLfloat *)params); }
//
// void (APIENTRYP pfn_glPointParameteri)(GLenum pname, GLint param);
// GLAPI void APIENTRY wrap_glPointParameteri(unsigned int pname, int param) {  (*pfn_glPointParameteri)((GLenum)pname, (GLint)param); }
//
// void (APIENTRYP pfn_glPointParameteriv)(GLenum pname, const GLint *params);
// GLAPI void APIENTRY wrap_glPointParameteriv(unsigned int pname, const int* params) {  (*pfn_glPointParameteriv)((GLenum)pname, (const GLint *)params); }
//
// void (APIENTRYP pfn_glPointSize)(GLfloat size);
// GLAPI void APIENTRY wrap_glPointSize(float size) {  (*pfn_glPointSize)((GLfloat)size); }
//
// void (APIENTRYP pfn_glPolygonMode)(GLenum face, GLenum mode);
// GLAPI void APIENTRY wrap_glPolygonMode(unsigned int face, unsigned int mode) {  (*pfn_glPolygonMode)((GLenum)face, (GLenum)mode); }
//
// void (APIENTRYP pfn_glPolygonOffset)(GLfloat factor, GLfloat units);
// GLAPI void APIENTRY wrap_glPolygonOffset(float factor, float units) {  (*pfn_glPolygonOffset)((GLfloat)factor, (GLfloat)units); }
//
// void (APIENTRYP pfn_glPopDebugGroup)(void);
// GLAPI void APIENTRY wrap_glPopDebugGroup() {  (*pfn_glPopDebugGroup)(); }
//
// void (APIENTRYP pfn_glPrimitiveRestartIndex)(GLuint index);
// GLAPI void APIENTRY wrap_glPrimitiveRestartIndex(unsigned int index) {  (*pfn_glPrimitiveRestartIndex)((GLuint)index); }
//
// void (APIENTRYP pfn_glProgramBinary)(GLuint program, GLenum binaryFormat, const GLvoid *binary, GLsizei length);
// GLAPI void APIENTRY wrap_glProgramBinary(unsigned int program, unsigned int binaryFormat, const void* binary, int length) {  (*pfn_glProgramBinary)((GLuint)program, (GLenum)binaryFormat, (const GLvoid *)binary, (GLsizei)length); }
//
// void (APIENTRYP pfn_glProgramParameteri)(GLuint program, GLenum pname, GLint value);
// GLAPI void APIENTRY wrap_glProgramParameteri(unsigned int program, unsigned int pname, int value) {  (*pfn_glProgramParameteri)((GLuint)program, (GLenum)pname, (GLint)value); }
//
// void (APIENTRYP pfn_glProgramUniform1d)(GLuint program, GLint location, GLdouble v0);
// GLAPI void APIENTRY wrap_glProgramUniform1d(unsigned int program, int location, double v0) {  (*pfn_glProgramUniform1d)((GLuint)program, (GLint)location, (GLdouble)v0); }
//
// void (APIENTRYP pfn_glProgramUniform1dv)(GLuint program, GLint location, GLsizei count, const GLdouble *value);
// GLAPI void APIENTRY wrap_glProgramUniform1dv(unsigned int program, int location, int count, const double* value) {  (*pfn_glProgramUniform1dv)((GLuint)program, (GLint)location, (GLsizei)count, (const GLdouble *)value); }
//
// void (APIENTRYP pfn_glProgramUniform1f)(GLuint program, GLint location, GLfloat v0);
// GLAPI void APIENTRY wrap_glProgramUniform1f(unsigned int program, int location, float v0) {  (*pfn_glProgramUniform1f)((GLuint)program, (GLint)location, (GLfloat)v0); }
//
// void (APIENTRYP pfn_glProgramUniform1fv)(GLuint program, GLint location, GLsizei count, const GLfloat *value);
// GLAPI void APIENTRY wrap_glProgramUniform1fv(unsigned int program, int location, int count, const float* value) {  (*pfn_glProgramUniform1fv)((GLuint)program, (GLint)location, (GLsizei)count, (const GLfloat *)value); }
//
// void (APIENTRYP pfn_glProgramUniform1i)(GLuint program, GLint location, GLint v0);
// GLAPI void APIENTRY wrap_glProgramUniform1i(unsigned int program, int location, int v0) {  (*pfn_glProgramUniform1i)((GLuint)program, (GLint)location, (GLint)v0); }
//
// void (APIENTRYP pfn_glProgramUniform1iv)(GLuint program, GLint location, GLsizei count, const GLint *value);
// GLAPI void APIENTRY wrap_glProgramUniform1iv(unsigned int program, int location, int count, const int* value) {  (*pfn_glProgramUniform1iv)((GLuint)program, (GLint)location, (GLsizei)count, (const GLint *)value); }
//
// void (APIENTRYP pfn_glProgramUniform1ui)(GLuint program, GLint location, GLuint v0);
// GLAPI void APIENTRY wrap_glProgramUniform1ui(unsigned int program, int location, unsigned int v0) {  (*pfn_glProgramUniform1ui)((GLuint)program, (GLint)location, (GLuint)v0); }
//
// void (APIENTRYP pfn_glProgramUniform1uiv)(GLuint program, GLint location, GLsizei count, const GLuint *value);
// GLAPI void APIENTRY wrap_glProgramUniform1uiv(unsigned int program, int location, int count, const unsigned int* value) {  (*pfn_glProgramUniform1uiv)((GLuint)program, (GLint)location, (GLsizei)count, (const GLuint *)value); }
//
// void (APIENTRYP pfn_glProgramUniform2d)(GLuint program, GLint location, GLdouble v0, GLdouble v1);
// GLAPI void APIENTRY wrap_glProgramUniform2d(unsigned int program, int location, double v0, double v1) {  (*pfn_glProgramUniform2d)((GLuint)program, (GLint)location, (GLdouble)v0, (GLdouble)v1); }
//
// void (APIENTRYP pfn_glProgramUniform2dv)(GLuint program, GLint location, GLsizei count, const GLdouble *value);
// GLAPI void APIENTRY wrap_glProgramUniform2dv(unsigned int program, int location, int count, const double* value) {  (*pfn_glProgramUniform2dv)((GLuint)program, (GLint)location, (GLsizei)count, (const GLdouble *)value); }
//
// void (APIENTRYP pfn_glProgramUniform2f)(GLuint program, GLint location, GLfloat v0, GLfloat v1);
// GLAPI void APIENTRY wrap_glProgramUniform2f(unsigned int program, int location, float v0, float v1) {  (*pfn_glProgramUniform2f)((GLuint)program, (GLint)location, (GLfloat)v0, (GLfloat)v1); }
//
// void (APIENTRYP pfn_glProgramUniform2fv)(GLuint program, GLint location, GLsizei count, const GLfloat *value);
// GLAPI void APIENTRY wrap_glProgramUniform2fv(unsigned int program, int location, int count, const float* value) {  (*pfn_glProgramUniform2fv)((GLuint)program, (GLint)location, (GLsizei)count, (const GLfloat *)value); }
//
// void (APIENTRYP pfn_glProgramUniform2i)(GLuint program, GLint location, GLint v0, GLint v1);
// GLAPI void APIENTRY wrap_glProgramUniform2i(unsigned int program, int location, int v0, int v1) {  (*pfn_glProgramUniform2i)((GLuint)program, (GLint)location, (GLint)v0, (GLint)v1); }
//
// void (APIENTRYP pfn_glProgramUniform2iv)(GLuint program, GLint location, GLsizei count, const GLint *value);
// GLAPI void APIENTRY wrap_glProgramUniform2iv(unsigned int program, int location, int count, const int* value) {  (*pfn_glProgramUniform2iv)((GLuint)program, (GLint)location, (GLsizei)count, (const GLint *)value); }
//
// void (APIENTRYP pfn_glProgramUniform2ui)(GLuint program, GLint location, GLuint v0, GLuint v1);
// GLAPI void APIENTRY wrap_glProgramUniform2ui(unsigned int program, int location, unsigned int v0, unsigned int v1) {  (*pfn_glProgramUniform2ui)((GLuint)program, (GLint)location, (GLuint)v0, (GLuint)v1); }
//
// void (APIENTRYP pfn_glProgramUniform2uiv)(GLuint program, GLint location, GLsizei count, const GLuint *value);
// GLAPI void APIENTRY wrap_glProgramUniform2uiv(unsigned int program, int location, int count, const unsigned int* value) {  (*pfn_glProgramUniform2uiv)((GLuint)program, (GLint)location, (GLsizei)count, (const GLuint *)value); }
//
// void (APIENTRYP pfn_glProgramUniform3d)(GLuint program, GLint location, GLdouble v0, GLdouble v1, GLdouble v2);
// GLAPI void APIENTRY wrap_glProgramUniform3d(unsigned int program, int location, double v0, double v1, double v2) {  (*pfn_glProgramUniform3d)((GLuint)program, (GLint)location, (GLdouble)v0, (GLdouble)v1, (GLdouble)v2); }
//
// void (APIENTRYP pfn_glProgramUniform3dv)(GLuint program, GLint location, GLsizei count, const GLdouble *value);
// GLAPI void APIENTRY wrap_glProgramUniform3dv(unsigned int program, int location, int count, const double* value) {  (*pfn_glProgramUniform3dv)((GLuint)program, (GLint)location, (GLsizei)count, (const GLdouble *)value); }
//
// void (APIENTRYP pfn_glProgramUniform3f)(GLuint program, GLint location, GLfloat v0, GLfloat v1, GLfloat v2);
// GLAPI void APIENTRY wrap_glProgramUniform3f(unsigned int program, int location, float v0, float v1, float v2) {  (*pfn_glProgramUniform3f)((GLuint)program, (GLint)location, (GLfloat)v0, (GLfloat)v1, (GLfloat)v2); }
//
// void (APIENTRYP pfn_glProgramUniform3fv)(GLuint program, GLint location, GLsizei count, const GLfloat *value);
// GLAPI void APIENTRY wrap_glProgramUniform3fv(unsigned int program, int location, int count, const float* value) {  (*pfn_glProgramUniform3fv)((GLuint)program, (GLint)location, (GLsizei)count, (const GLfloat *)value); }
//
// void (APIENTRYP pfn_glProgramUniform3i)(GLuint program, GLint location, GLint v0, GLint v1, GLint v2);
// GLAPI void APIENTRY wrap_glProgramUniform3i(unsigned int program, int location, int v0, int v1, int v2) {  (*pfn_glProgramUniform3i)((GLuint)program, (GLint)location, (GLint)v0, (GLint)v1, (GLint)v2); }
//
// void (APIENTRYP pfn_glProgramUniform3iv)(GLuint program, GLint location, GLsizei count, const GLint *value);
// GLAPI void APIENTRY wrap_glProgramUniform3iv(unsigned int program, int location, int count, const int* value) {  (*pfn_glProgramUniform3iv)((GLuint)program, (GLint)location, (GLsizei)count, (const GLint *)value); }
//
// void (APIENTRYP pfn_glProgramUniform3ui)(GLuint program, GLint location, GLuint v0, GLuint v1, GLuint v2);
// GLAPI void APIENTRY wrap_glProgramUniform3ui(unsigned int program, int location, unsigned int v0, unsigned int v1, unsigned int v2) {  (*pfn_glProgramUniform3ui)((GLuint)program, (GLint)location, (GLuint)v0, (GLuint)v1, (GLuint)v2); }
//
// void (APIENTRYP pfn_glProgramUniform3uiv)(GLuint program, GLint location, GLsizei count, const GLuint *value);
// GLAPI void APIENTRY wrap_glProgramUniform3uiv(unsigned int program, int location, int count, const unsigned int* value) {  (*pfn_glProgramUniform3uiv)((GLuint)program, (GLint)location, (GLsizei)count, (const GLuint *)value); }
//
// void (APIENTRYP pfn_glProgramUniform4d)(GLuint program, GLint location, GLdouble v0, GLdouble v1, GLdouble v2, GLdouble v3);
// GLAPI void APIENTRY wrap_glProgramUniform4d(unsigned int program, int location, double v0, double v1, double v2, double v3) {  (*pfn_glProgramUniform4d)((GLuint)program, (GLint)location, (GLdouble)v0, (GLdouble)v1, (GLdouble)v2, (GLdouble)v3); }
//
// void (APIENTRYP pfn_glProgramUniform4dv)(GLuint program, GLint location, GLsizei count, const GLdouble *value);
// GLAPI void APIENTRY wrap_glProgramUniform4dv(unsigned int program, int location, int count, const double* value) {  (*pfn_glProgramUniform4dv)((GLuint)program, (GLint)location, (GLsizei)count, (const GLdouble *)value); }
//
// void (APIENTRYP pfn_glProgramUniform4f)(GLuint program, GLint location, GLfloat v0, GLfloat v1, GLfloat v2, GLfloat v3);
// GLAPI void APIENTRY wrap_glProgramUniform4f(unsigned int program, int location, float v0, float v1, float v2, float v3) {  (*pfn_glProgramUniform4f)((GLuint)program, (GLint)location, (GLfloat)v0, (GLfloat)v1, (GLfloat)v2, (GLfloat)v3); }
//
// void (APIENTRYP pfn_glProgramUniform4fv)(GLuint program, GLint location, GLsizei count, const GLfloat *value);
// GLAPI void APIENTRY wrap_glProgramUniform4fv(unsigned int program, int location, int count, const float* value) {  (*pfn_glProgramUniform4fv)((GLuint)program, (GLint)location, (GLsizei)count, (const GLfloat *)value); }
//
// void (APIENTRYP pfn_glProgramUniform4i)(GLuint program, GLint location, GLint v0, GLint v1, GLint v2, GLint v3);
// GLAPI void APIENTRY wrap_glProgramUniform4i(unsigned int program, int location, int v0, int v1, int v2, int v3) {  (*pfn_glProgramUniform4i)((GLuint)program, (GLint)location, (GLint)v0, (GLint)v1, (GLint)v2, (GLint)v3); }
//
// void (APIENTRYP pfn_glProgramUniform4iv)(GLuint program, GLint location, GLsizei count, const GLint *value);
// GLAPI void APIENTRY wrap_glProgramUniform4iv(unsigned int program, int location, int count, const int* value) {  (*pfn_glProgramUniform4iv)((GLuint)program, (GLint)location, (GLsizei)count, (const GLint *)value); }
//
// void (APIENTRYP pfn_glProgramUniform4ui)(GLuint program, GLint location, GLuint v0, GLuint v1, GLuint v2, GLuint v3);
// GLAPI void APIENTRY wrap_glProgramUniform4ui(unsigned int program, int location, unsigned int v0, unsigned int v1, unsigned int v2, unsigned int v3) {  (*pfn_glProgramUniform4ui)((GLuint)program, (GLint)location, (GLuint)v0, (GLuint)v1, (GLuint)v2, (GLuint)v3); }
//
// void (APIENTRYP pfn_glProgramUniform4uiv)(GLuint program, GLint location, GLsizei count, const GLuint *value);
// GLAPI void APIENTRY wrap_glProgramUniform4uiv(unsigned int program, int location, int count, const unsigned int* value) {  (*pfn_glProgramUniform4uiv)((GLuint)program, (GLint)location, (GLsizei)count, (const GLuint *)value); }
//
// void (APIENTRYP pfn_glProgramUniformMatrix2dv)(GLuint program, GLint location, GLsizei count, GLboolean transpose, const GLdouble *value);
// GLAPI void APIENTRY wrap_glProgramUniformMatrix2dv(unsigned int program, int location, int count, unsigned char transpose, const double* value) {  (*pfn_glProgramUniformMatrix2dv)((GLuint)program, (GLint)location, (GLsizei)count, (GLboolean)transpose, (const GLdouble *)value); }
//
// void (APIENTRYP pfn_glProgramUniformMatrix2fv)(GLuint program, GLint location, GLsizei count, GLboolean transpose, const GLfloat *value);
// GLAPI void APIENTRY wrap_glProgramUniformMatrix2fv(unsigned int program, int location, int count, unsigned char transpose, const float* value) {  (*pfn_glProgramUniformMatrix2fv)((GLuint)program, (GLint)location, (GLsizei)count, (GLboolean)transpose, (const GLfloat *)value); }
//
// void (APIENTRYP pfn_glProgramUniformMatrix2x3dv)(GLuint program, GLint location, GLsizei count, GLboolean transpose, const GLdouble *value);
// GLAPI void APIENTRY wrap_glProgramUniformMatrix2x3dv(unsigned int program, int location, int count, unsigned char transpose, const double* value) {  (*pfn_glProgramUniformMatrix2x3dv)((GLuint)program, (GLint)location, (GLsizei)count, (GLboolean)transpose, (const GLdouble *)value); }
//
// void (APIENTRYP pfn_glProgramUniformMatrix2x3fv)(GLuint program, GLint location, GLsizei count, GLboolean transpose, const GLfloat *value);
// GLAPI void APIENTRY wrap_glProgramUniformMatrix2x3fv(unsigned int program, int location, int count, unsigned char transpose, const float* value) {  (*pfn_glProgramUniformMatrix2x3fv)((GLuint)program, (GLint)location, (GLsizei)count, (GLboolean)transpose, (const GLfloat *)value); }
//
// void (APIENTRYP pfn_glProgramUniformMatrix2x4dv)(GLuint program, GLint location, GLsizei count, GLboolean transpose, const GLdouble *value);
// GLAPI void APIENTRY wrap_glProgramUniformMatrix2x4dv(unsigned int program, int location, int count, unsigned char transpose, const double* value) {  (*pfn_glProgramUniformMatrix2x4dv)((GLuint)program, (GLint)location, (GLsizei)count, (GLboolean)transpose, (const GLdouble *)value); }
//
// void (APIENTRYP pfn_glProgramUniformMatrix2x4fv)(GLuint program, GLint location, GLsizei count, GLboolean transpose, const GLfloat *value);
// GLAPI void APIENTRY wrap_glProgramUniformMatrix2x4fv(unsigned int program, int location, int count, unsigned char transpose, const float* value) {  (*pfn_glProgramUniformMatrix2x4fv)((GLuint)program, (GLint)location, (GLsizei)count, (GLboolean)transpose, (const GLfloat *)value); }
//
// void (APIENTRYP pfn_glProgramUniformMatrix3dv)(GLuint program, GLint location, GLsizei count, GLboolean transpose, const GLdouble *value);
// GLAPI void APIENTRY wrap_glProgramUniformMatrix3dv(unsigned int program, int location, int count, unsigned char transpose, const double* value) {  (*pfn_glProgramUniformMatrix3dv)((GLuint)program, (GLint)location, (GLsizei)count, (GLboolean)transpose, (const GLdouble *)value); }
//
// void (APIENTRYP pfn_glProgramUniformMatrix3fv)(GLuint program, GLint location, GLsizei count, GLboolean transpose, const GLfloat *value);
// GLAPI void APIENTRY wrap_glProgramUniformMatrix3fv(unsigned int program, int location, int count, unsigned char transpose, const float* value) {  (*pfn_glProgramUniformMatrix3fv)((GLuint)program, (GLint)location, (GLsizei)count, (GLboolean)transpose, (const GLfloat *)value); }
//
// void (APIENTRYP pfn_glProgramUniformMatrix3x2dv)(GLuint program, GLint location, GLsizei count, GLboolean transpose, const GLdouble *value);
// GLAPI void APIENTRY wrap_glProgramUniformMatrix3x2dv(unsigned int program, int location, int count, unsigned char transpose, const double* value) {  (*pfn_glProgramUniformMatrix3x2dv)((GLuint)program, (GLint)location, (GLsizei)count, (GLboolean)transpose, (const GLdouble *)value); }
//
// void (APIENTRYP pfn_glProgramUniformMatrix3x2fv)(GLuint program, GLint location, GLsizei count, GLboolean transpose, const GLfloat *value);
// GLAPI void APIENTRY wrap_glProgramUniformMatrix3x2fv(unsigned int program, int location, int count, unsigned char transpose, const float* value) {  (*pfn_glProgramUniformMatrix3x2fv)((GLuint)program, (GLint)location, (GLsizei)count, (GLboolean)transpose, (const GLfloat *)value); }
//
// void (APIENTRYP pfn_glProgramUniformMatrix3x4dv)(GLuint program, GLint location, GLsizei count, GLboolean transpose, const GLdouble *value);
// GLAPI void APIENTRY wrap_glProgramUniformMatrix3x4dv(unsigned int program, int location, int count, unsigned char transpose, const double* value) {  (*pfn_glProgramUniformMatrix3x4dv)((GLuint)program, (GLint)location, (GLsizei)count, (GLboolean)transpose, (const GLdouble *)value); }
//
// void (APIENTRYP pfn_glProgramUniformMatrix3x4fv)(GLuint program, GLint location, GLsizei count, GLboolean transpose, const GLfloat *value);
// GLAPI void APIENTRY wrap_glProgramUniformMatrix3x4fv(unsigned int program, int location, int count, unsigned char transpose, const float* value) {  (*pfn_glProgramUniformMatrix3x4fv)((GLuint)program, (GLint)location, (GLsizei)count, (GLboolean)transpose, (const GLfloat *)value); }
//
// void (APIENTRYP pfn_glProgramUniformMatrix4dv)(GLuint program, GLint location, GLsizei count, GLboolean transpose, const GLdouble *value);
// GLAPI void APIENTRY wrap_glProgramUniformMatrix4dv(unsigned int program, int location, int count, unsigned char transpose, const double* value) {  (*pfn_glProgramUniformMatrix4dv)((GLuint)program, (GLint)location, (GLsizei)count, (GLboolean)transpose, (const GLdouble *)value); }
//
// void (APIENTRYP pfn_glProgramUniformMatrix4fv)(GLuint program, GLint location, GLsizei count, GLboolean transpose, const GLfloat *value);
// GLAPI void APIENTRY wrap_glProgramUniformMatrix4fv(unsigned int program, int location, int count, unsigned char transpose, const float* value) {  (*pfn_glProgramUniformMatrix4fv)((GLuint)program, (GLint)location, (GLsizei)count, (GLboolean)transpose, (const GLfloat *)value); }
//
// void (APIENTRYP pfn_glProgramUniformMatrix4x2dv)(GLuint program, GLint location, GLsizei count, GLboolean transpose, const GLdouble *value);
// GLAPI void APIENTRY wrap_glProgramUniformMatrix4x2dv(unsigned int program, int location, int count, unsigned char transpose, const double* value) {  (*pfn_glProgramUniformMatrix4x2dv)((GLuint)program, (GLint)location, (GLsizei)count, (GLboolean)transpose, (const GLdouble *)value); }
//
// void (APIENTRYP pfn_glProgramUniformMatrix4x2fv)(GLuint program, GLint location, GLsizei count, GLboolean transpose, const GLfloat *value);
// GLAPI void APIENTRY wrap_glProgramUniformMatrix4x2fv(unsigned int program, int location, int count, unsigned char transpose, const float* value) {  (*pfn_glProgramUniformMatrix4x2fv)((GLuint)program, (GLint)location, (GLsizei)count, (GLboolean)transpose, (const GLfloat *)value); }
//
// void (APIENTRYP pfn_glProgramUniformMatrix4x3dv)(GLuint program, GLint location, GLsizei count, GLboolean transpose, const GLdouble *value);
// GLAPI void APIENTRY wrap_glProgramUniformMatrix4x3dv(unsigned int program, int location, int count, unsigned char transpose, const double* value) {  (*pfn_glProgramUniformMatrix4x3dv)((GLuint)program, (GLint)location, (GLsizei)count, (GLboolean)transpose, (const GLdouble *)value); }
//
// void (APIENTRYP pfn_glProgramUniformMatrix4x3fv)(GLuint program, GLint location, GLsizei count, GLboolean transpose, const GLfloat *value);
// GLAPI void APIENTRY wrap_glProgramUniformMatrix4x3fv(unsigned int program, int location, int count, unsigned char transpose, const float* value) {  (*pfn_glProgramUniformMatrix4x3fv)((GLuint)program, (GLint)location, (GLsizei)count, (GLboolean)transpose, (const GLfloat *)value); }
//
// void (APIENTRYP pfn_glProvokingVertex)(GLenum mode);
// GLAPI void APIENTRY wrap_glProvokingVertex(unsigned int mode) {  (*pfn_glProvokingVertex)((GLenum)mode); }
//
// void (APIENTRYP pfn_glPushDebugGroup)(GLenum source, GLuint id, GLsizei length, const GLchar *message);
// GLAPI void APIENTRY wrap_glPushDebugGroup(unsigned int source, unsigned int id, int length, const char* message) {  (*pfn_glPushDebugGroup)((GLenum)source, (GLuint)id, (GLsizei)length, (const GLchar *)message); }
//
// void (APIENTRYP pfn_glQueryCounter)(GLuint id, GLenum target);
// GLAPI void APIENTRY wrap_glQueryCounter(unsigned int id, unsigned int target) {  (*pfn_glQueryCounter)((GLuint)id, (GLenum)target); }
//
// void (APIENTRYP pfn_glReadBuffer)(GLenum mode);
// GLAPI void APIENTRY wrap_glReadBuffer(unsigned int mode) {  (*pfn_glReadBuffer)((GLenum)mode); }
//
// void (APIENTRYP pfn_glReadPixels)(GLint x, GLint y, GLsizei width, GLsizei height, GLenum format, GLenum type, GLvoid *pixels);
// GLAPI void APIENTRY wrap_glReadPixels(int x, int y, int width, int height, unsigned int format, unsigned int t_ype, void* pixels) {  (*pfn_glReadPixels)((GLint)x, (GLint)y, (GLsizei)width, (GLsizei)height, (GLenum)format, (GLenum)t_ype, (GLvoid *)pixels); }
//
// void (APIENTRYP pfn_glReadnPixelsARB)(GLint x, GLint y, GLsizei width, GLsizei height, GLenum format, GLenum type, GLsizei bufSize, GLvoid *data);
// GLAPI void APIENTRY wrap_glReadnPixelsARB(int x, int y, int width, int height, unsigned int format, unsigned int t_ype, int bufSize, void* data) {  (*pfn_glReadnPixelsARB)((GLint)x, (GLint)y, (GLsizei)width, (GLsizei)height, (GLenum)format, (GLenum)t_ype, (GLsizei)bufSize, (GLvoid *)data); }
//
// void (APIENTRYP pfn_glReleaseShaderCompiler)(void);
// GLAPI void APIENTRY wrap_glReleaseShaderCompiler() {  (*pfn_glReleaseShaderCompiler)(); }
//
// void (APIENTRYP pfn_glRenderbufferStorage)(GLenum target, GLenum internalformat, GLsizei width, GLsizei height);
// GLAPI void APIENTRY wrap_glRenderbufferStorage(unsigned int target, unsigned int internalformat, int width, int height) {  (*pfn_glRenderbufferStorage)((GLenum)target, (GLenum)internalformat, (GLsizei)width, (GLsizei)height); }
//
// void (APIENTRYP pfn_glRenderbufferStorageMultisample)(GLenum target, GLsizei samples, GLenum internalformat, GLsizei width, GLsizei height);
// GLAPI void APIENTRY wrap_glRenderbufferStorageMultisample(unsigned int target, int samples, unsigned int internalformat, int width, int height) {  (*pfn_glRenderbufferStorageMultisample)((GLenum)target, (GLsizei)samples, (GLenum)internalformat, (GLsizei)width, (GLsizei)height); }
//
// void (APIENTRYP pfn_glResumeTransformFeedback)(void);
// GLAPI void APIENTRY wrap_glResumeTransformFeedback() {  (*pfn_glResumeTransformFeedback)(); }
//
// void (APIENTRYP pfn_glSampleCoverage)(GLfloat value, GLboolean invert);
// GLAPI void APIENTRY wrap_glSampleCoverage(float value, unsigned char invert) {  (*pfn_glSampleCoverage)((GLfloat)value, (GLboolean)invert); }
//
// void (APIENTRYP pfn_glSampleMaski)(GLuint index, GLbitfield mask);
// GLAPI void APIENTRY wrap_glSampleMaski(unsigned int index, unsigned int mask) {  (*pfn_glSampleMaski)((GLuint)index, (GLbitfield)mask); }
//
// void (APIENTRYP pfn_glSamplerParameterIiv)(GLuint sampler, GLenum pname, const GLint *param);
// GLAPI void APIENTRY wrap_glSamplerParameterIiv(unsigned int sampler, unsigned int pname, const int* param) {  (*pfn_glSamplerParameterIiv)((GLuint)sampler, (GLenum)pname, (const GLint *)param); }
//
// void (APIENTRYP pfn_glSamplerParameterIuiv)(GLuint sampler, GLenum pname, const GLuint *param);
// GLAPI void APIENTRY wrap_glSamplerParameterIuiv(unsigned int sampler, unsigned int pname, const unsigned int* param) {  (*pfn_glSamplerParameterIuiv)((GLuint)sampler, (GLenum)pname, (const GLuint *)param); }
//
// void (APIENTRYP pfn_glSamplerParameterf)(GLuint sampler, GLenum pname, GLfloat param);
// GLAPI void APIENTRY wrap_glSamplerParameterf(unsigned int sampler, unsigned int pname, float param) {  (*pfn_glSamplerParameterf)((GLuint)sampler, (GLenum)pname, (GLfloat)param); }
//
// void (APIENTRYP pfn_glSamplerParameterfv)(GLuint sampler, GLenum pname, const GLfloat *param);
// GLAPI void APIENTRY wrap_glSamplerParameterfv(unsigned int sampler, unsigned int pname, const float* param) {  (*pfn_glSamplerParameterfv)((GLuint)sampler, (GLenum)pname, (const GLfloat *)param); }
//
// void (APIENTRYP pfn_glSamplerParameteri)(GLuint sampler, GLenum pname, GLint param);
// GLAPI void APIENTRY wrap_glSamplerParameteri(unsigned int sampler, unsigned int pname, int param) {  (*pfn_glSamplerParameteri)((GLuint)sampler, (GLenum)pname, (GLint)param); }
//
// void (APIENTRYP pfn_glSamplerParameteriv)(GLuint sampler, GLenum pname, const GLint *param);
// GLAPI void APIENTRY wrap_glSamplerParameteriv(unsigned int sampler, unsigned int pname, const int* param) {  (*pfn_glSamplerParameteriv)((GLuint)sampler, (GLenum)pname, (const GLint *)param); }
//
// void (APIENTRYP pfn_glScissor)(GLint x, GLint y, GLsizei width, GLsizei height);
// GLAPI void APIENTRY wrap_glScissor(int x, int y, int width, int height) {  (*pfn_glScissor)((GLint)x, (GLint)y, (GLsizei)width, (GLsizei)height); }
//
// void (APIENTRYP pfn_glScissorArrayv)(GLuint first, GLsizei count, const GLint *v);
// GLAPI void APIENTRY wrap_glScissorArrayv(unsigned int first, int count, const int* v) {  (*pfn_glScissorArrayv)((GLuint)first, (GLsizei)count, (const GLint *)v); }
//
// void (APIENTRYP pfn_glScissorIndexed)(GLuint index, GLint left, GLint bottom, GLsizei width, GLsizei height);
// GLAPI void APIENTRY wrap_glScissorIndexed(unsigned int index, int left, int bottom, int width, int height) {  (*pfn_glScissorIndexed)((GLuint)index, (GLint)left, (GLint)bottom, (GLsizei)width, (GLsizei)height); }
//
// void (APIENTRYP pfn_glScissorIndexedv)(GLuint index, const GLint *v);
// GLAPI void APIENTRY wrap_glScissorIndexedv(unsigned int index, const int* v) {  (*pfn_glScissorIndexedv)((GLuint)index, (const GLint *)v); }
//
// void (APIENTRYP pfn_glSecondaryColorP3ui)(GLenum type, GLuint color);
// GLAPI void APIENTRY wrap_glSecondaryColorP3ui(unsigned int t_ype, unsigned int color) {  (*pfn_glSecondaryColorP3ui)((GLenum)t_ype, (GLuint)color); }
//
// void (APIENTRYP pfn_glSecondaryColorP3uiv)(GLenum type, const GLuint *color);
// GLAPI void APIENTRY wrap_glSecondaryColorP3uiv(unsigned int t_ype, const unsigned int* color) {  (*pfn_glSecondaryColorP3uiv)((GLenum)t_ype, (const GLuint *)color); }
//
// void (APIENTRYP pfn_glShaderBinary)(GLsizei count, const GLuint *shaders, GLenum binaryformat, const GLvoid *binary, GLsizei length);
// GLAPI void APIENTRY wrap_glShaderBinary(int count, const unsigned int* shaders, unsigned int binaryformat, const void* binary, int length) {  (*pfn_glShaderBinary)((GLsizei)count, (const GLuint *)shaders, (GLenum)binaryformat, (const GLvoid *)binary, (GLsizei)length); }
//
// void (APIENTRYP pfn_glShaderSource)(GLuint shader, GLsizei count, const GLchar* const *string, const GLint *length);
// GLAPI void APIENTRY wrap_glShaderSource(unsigned int shader, int count, const char* const* s_tring, const int* length) {  (*pfn_glShaderSource)((GLuint)shader, (GLsizei)count, (const GLchar* const *)s_tring, (const GLint *)length); }
//
// void (APIENTRYP pfn_glShaderStorageBlockBinding)(GLuint program, GLuint storageBlockIndex, GLuint storageBlockBinding);
// GLAPI void APIENTRY wrap_glShaderStorageBlockBinding(unsigned int program, unsigned int storageBlockIndex, unsigned int storageBlockBinding) {  (*pfn_glShaderStorageBlockBinding)((GLuint)program, (GLuint)storageBlockIndex, (GLuint)storageBlockBinding); }
//
// void (APIENTRYP pfn_glStencilFunc)(GLenum func, GLint ref, GLuint mask);
// GLAPI void APIENTRY wrap_glStencilFunc(unsigned int f_unc, int ref, unsigned int mask) {  (*pfn_glStencilFunc)((GLenum)f_unc, (GLint)ref, (GLuint)mask); }
//
// void (APIENTRYP pfn_glStencilFuncSeparate)(GLenum face, GLenum func, GLint ref, GLuint mask);
// GLAPI void APIENTRY wrap_glStencilFuncSeparate(unsigned int face, unsigned int f_unc, int ref, unsigned int mask) {  (*pfn_glStencilFuncSeparate)((GLenum)face, (GLenum)f_unc, (GLint)ref, (GLuint)mask); }
//
// void (APIENTRYP pfn_glStencilMask)(GLuint mask);
// GLAPI void APIENTRY wrap_glStencilMask(unsigned int mask) {  (*pfn_glStencilMask)((GLuint)mask); }
//
// void (APIENTRYP pfn_glStencilMaskSeparate)(GLenum face, GLuint mask);
// GLAPI void APIENTRY wrap_glStencilMaskSeparate(unsigned int face, unsigned int mask) {  (*pfn_glStencilMaskSeparate)((GLenum)face, (GLuint)mask); }
//
// void (APIENTRYP pfn_glStencilOp)(GLenum fail, GLenum zfail, GLenum zpass);
// GLAPI void APIENTRY wrap_glStencilOp(unsigned int fail, unsigned int zfail, unsigned int zpass) {  (*pfn_glStencilOp)((GLenum)fail, (GLenum)zfail, (GLenum)zpass); }
//
// void (APIENTRYP pfn_glStencilOpSeparate)(GLenum face, GLenum sfail, GLenum dpfail, GLenum dppass);
// GLAPI void APIENTRY wrap_glStencilOpSeparate(unsigned int face, unsigned int sfail, unsigned int dpfail, unsigned int dppass) {  (*pfn_glStencilOpSeparate)((GLenum)face, (GLenum)sfail, (GLenum)dpfail, (GLenum)dppass); }
//
// void (APIENTRYP pfn_glTexBuffer)(GLenum target, GLenum internalformat, GLuint buffer);
// GLAPI void APIENTRY wrap_glTexBuffer(unsigned int target, unsigned int internalformat, unsigned int buffer) {  (*pfn_glTexBuffer)((GLenum)target, (GLenum)internalformat, (GLuint)buffer); }
//
// void (APIENTRYP pfn_glTexBufferRange)(GLenum target, GLenum internalformat, GLuint buffer, GLintptr offset, GLsizeiptr size);
// GLAPI void APIENTRY wrap_glTexBufferRange(unsigned int target, unsigned int internalformat, unsigned int buffer, long long offset, long long size) {  (*pfn_glTexBufferRange)((GLenum)target, (GLenum)internalformat, (GLuint)buffer, (GLintptr)offset, (GLsizeiptr)size); }
//
// void (APIENTRYP pfn_glTexCoordP1ui)(GLenum type, GLuint coords);
// GLAPI void APIENTRY wrap_glTexCoordP1ui(unsigned int t_ype, unsigned int coords) {  (*pfn_glTexCoordP1ui)((GLenum)t_ype, (GLuint)coords); }
//
// void (APIENTRYP pfn_glTexCoordP1uiv)(GLenum type, const GLuint *coords);
// GLAPI void APIENTRY wrap_glTexCoordP1uiv(unsigned int t_ype, const unsigned int* coords) {  (*pfn_glTexCoordP1uiv)((GLenum)t_ype, (const GLuint *)coords); }
//
// void (APIENTRYP pfn_glTexCoordP2ui)(GLenum type, GLuint coords);
// GLAPI void APIENTRY wrap_glTexCoordP2ui(unsigned int t_ype, unsigned int coords) {  (*pfn_glTexCoordP2ui)((GLenum)t_ype, (GLuint)coords); }
//
// void (APIENTRYP pfn_glTexCoordP2uiv)(GLenum type, const GLuint *coords);
// GLAPI void APIENTRY wrap_glTexCoordP2uiv(unsigned int t_ype, const unsigned int* coords) {  (*pfn_glTexCoordP2uiv)((GLenum)t_ype, (const GLuint *)coords); }
//
// void (APIENTRYP pfn_glTexCoordP3ui)(GLenum type, GLuint coords);
// GLAPI void APIENTRY wrap_glTexCoordP3ui(unsigned int t_ype, unsigned int coords) {  (*pfn_glTexCoordP3ui)((GLenum)t_ype, (GLuint)coords); }
//
// void (APIENTRYP pfn_glTexCoordP3uiv)(GLenum type, const GLuint *coords);
// GLAPI void APIENTRY wrap_glTexCoordP3uiv(unsigned int t_ype, const unsigned int* coords) {  (*pfn_glTexCoordP3uiv)((GLenum)t_ype, (const GLuint *)coords); }
//
// void (APIENTRYP pfn_glTexCoordP4ui)(GLenum type, GLuint coords);
// GLAPI void APIENTRY wrap_glTexCoordP4ui(unsigned int t_ype, unsigned int coords) {  (*pfn_glTexCoordP4ui)((GLenum)t_ype, (GLuint)coords); }
//
// void (APIENTRYP pfn_glTexCoordP4uiv)(GLenum type, const GLuint *coords);
// GLAPI void APIENTRY wrap_glTexCoordP4uiv(unsigned int t_ype, const unsigned int* coords) {  (*pfn_glTexCoordP4uiv)((GLenum)t_ype, (const GLuint *)coords); }
//
// void (APIENTRYP pfn_glTexImage1D)(GLenum target, GLint level, GLint internalformat, GLsizei width, GLint border, GLenum format, GLenum type, const GLvoid *pixels);
// GLAPI void APIENTRY wrap_glTexImage1D(unsigned int target, int level, int internalformat, int width, int border, unsigned int format, unsigned int t_ype, const void* pixels) {  (*pfn_glTexImage1D)((GLenum)target, (GLint)level, (GLint)internalformat, (GLsizei)width, (GLint)border, (GLenum)format, (GLenum)t_ype, (const GLvoid *)pixels); }
//
// void (APIENTRYP pfn_glTexImage2D)(GLenum target, GLint level, GLint internalformat, GLsizei width, GLsizei height, GLint border, GLenum format, GLenum type, const GLvoid *pixels);
// GLAPI void APIENTRY wrap_glTexImage2D(unsigned int target, int level, int internalformat, int width, int height, int border, unsigned int format, unsigned int t_ype, const void* pixels) {  (*pfn_glTexImage2D)((GLenum)target, (GLint)level, (GLint)internalformat, (GLsizei)width, (GLsizei)height, (GLint)border, (GLenum)format, (GLenum)t_ype, (const GLvoid *)pixels); }
//
// void (APIENTRYP pfn_glTexImage2DMultisample)(GLenum target, GLsizei samples, GLint internalformat, GLsizei width, GLsizei height, GLboolean fixedsamplelocations);
// GLAPI void APIENTRY wrap_glTexImage2DMultisample(unsigned int target, int samples, int internalformat, int width, int height, unsigned char fixedsamplelocations) {  (*pfn_glTexImage2DMultisample)((GLenum)target, (GLsizei)samples, (GLint)internalformat, (GLsizei)width, (GLsizei)height, (GLboolean)fixedsamplelocations); }
//
// void (APIENTRYP pfn_glTexImage3D)(GLenum target, GLint level, GLint internalformat, GLsizei width, GLsizei height, GLsizei depth, GLint border, GLenum format, GLenum type, const GLvoid *pixels);
// GLAPI void APIENTRY wrap_glTexImage3D(unsigned int target, int level, int internalformat, int width, int height, int depth, int border, unsigned int format, unsigned int t_ype, const void* pixels) {  (*pfn_glTexImage3D)((GLenum)target, (GLint)level, (GLint)internalformat, (GLsizei)width, (GLsizei)height, (GLsizei)depth, (GLint)border, (GLenum)format, (GLenum)t_ype, (const GLvoid *)pixels); }
//
// void (APIENTRYP pfn_glTexImage3DMultisample)(GLenum target, GLsizei samples, GLint internalformat, GLsizei width, GLsizei height, GLsizei depth, GLboolean fixedsamplelocations);
// GLAPI void APIENTRY wrap_glTexImage3DMultisample(unsigned int target, int samples, int internalformat, int width, int height, int depth, unsigned char fixedsamplelocations) {  (*pfn_glTexImage3DMultisample)((GLenum)target, (GLsizei)samples, (GLint)internalformat, (GLsizei)width, (GLsizei)height, (GLsizei)depth, (GLboolean)fixedsamplelocations); }
//
// void (APIENTRYP pfn_glTexParameterIiv)(GLenum target, GLenum pname, const GLint *params);
// GLAPI void APIENTRY wrap_glTexParameterIiv(unsigned int target, unsigned int pname, const int* params) {  (*pfn_glTexParameterIiv)((GLenum)target, (GLenum)pname, (const GLint *)params); }
//
// void (APIENTRYP pfn_glTexParameterIuiv)(GLenum target, GLenum pname, const GLuint *params);
// GLAPI void APIENTRY wrap_glTexParameterIuiv(unsigned int target, unsigned int pname, const unsigned int* params) {  (*pfn_glTexParameterIuiv)((GLenum)target, (GLenum)pname, (const GLuint *)params); }
//
// void (APIENTRYP pfn_glTexParameterf)(GLenum target, GLenum pname, GLfloat param);
// GLAPI void APIENTRY wrap_glTexParameterf(unsigned int target, unsigned int pname, float param) {  (*pfn_glTexParameterf)((GLenum)target, (GLenum)pname, (GLfloat)param); }
//
// void (APIENTRYP pfn_glTexParameterfv)(GLenum target, GLenum pname, const GLfloat *params);
// GLAPI void APIENTRY wrap_glTexParameterfv(unsigned int target, unsigned int pname, const float* params) {  (*pfn_glTexParameterfv)((GLenum)target, (GLenum)pname, (const GLfloat *)params); }
//
// void (APIENTRYP pfn_glTexParameteri)(GLenum target, GLenum pname, GLint param);
// GLAPI void APIENTRY wrap_glTexParameteri(unsigned int target, unsigned int pname, int param) {  (*pfn_glTexParameteri)((GLenum)target, (GLenum)pname, (GLint)param); }
//
// void (APIENTRYP pfn_glTexParameteriv)(GLenum target, GLenum pname, const GLint *params);
// GLAPI void APIENTRY wrap_glTexParameteriv(unsigned int target, unsigned int pname, const int* params) {  (*pfn_glTexParameteriv)((GLenum)target, (GLenum)pname, (const GLint *)params); }
//
// void (APIENTRYP pfn_glTexStorage1D)(GLenum target, GLsizei levels, GLenum internalformat, GLsizei width);
// GLAPI void APIENTRY wrap_glTexStorage1D(unsigned int target, int levels, unsigned int internalformat, int width) {  (*pfn_glTexStorage1D)((GLenum)target, (GLsizei)levels, (GLenum)internalformat, (GLsizei)width); }
//
// void (APIENTRYP pfn_glTexStorage2D)(GLenum target, GLsizei levels, GLenum internalformat, GLsizei width, GLsizei height);
// GLAPI void APIENTRY wrap_glTexStorage2D(unsigned int target, int levels, unsigned int internalformat, int width, int height) {  (*pfn_glTexStorage2D)((GLenum)target, (GLsizei)levels, (GLenum)internalformat, (GLsizei)width, (GLsizei)height); }
//
// void (APIENTRYP pfn_glTexStorage2DMultisample)(GLenum target, GLsizei samples, GLenum internalformat, GLsizei width, GLsizei height, GLboolean fixedsamplelocations);
// GLAPI void APIENTRY wrap_glTexStorage2DMultisample(unsigned int target, int samples, unsigned int internalformat, int width, int height, unsigned char fixedsamplelocations) {  (*pfn_glTexStorage2DMultisample)((GLenum)target, (GLsizei)samples, (GLenum)internalformat, (GLsizei)width, (GLsizei)height, (GLboolean)fixedsamplelocations); }
//
// void (APIENTRYP pfn_glTexStorage3D)(GLenum target, GLsizei levels, GLenum internalformat, GLsizei width, GLsizei height, GLsizei depth);
// GLAPI void APIENTRY wrap_glTexStorage3D(unsigned int target, int levels, unsigned int internalformat, int width, int height, int depth) {  (*pfn_glTexStorage3D)((GLenum)target, (GLsizei)levels, (GLenum)internalformat, (GLsizei)width, (GLsizei)height, (GLsizei)depth); }
//
// void (APIENTRYP pfn_glTexStorage3DMultisample)(GLenum target, GLsizei samples, GLenum internalformat, GLsizei width, GLsizei height, GLsizei depth, GLboolean fixedsamplelocations);
// GLAPI void APIENTRY wrap_glTexStorage3DMultisample(unsigned int target, int samples, unsigned int internalformat, int width, int height, int depth, unsigned char fixedsamplelocations) {  (*pfn_glTexStorage3DMultisample)((GLenum)target, (GLsizei)samples, (GLenum)internalformat, (GLsizei)width, (GLsizei)height, (GLsizei)depth, (GLboolean)fixedsamplelocations); }
//
// void (APIENTRYP pfn_glTexSubImage1D)(GLenum target, GLint level, GLint xoffset, GLsizei width, GLenum format, GLenum type, const GLvoid *pixels);
// GLAPI void APIENTRY wrap_glTexSubImage1D(unsigned int target, int level, int xoffset, int width, unsigned int format, unsigned int t_ype, const void* pixels) {  (*pfn_glTexSubImage1D)((GLenum)target, (GLint)level, (GLint)xoffset, (GLsizei)width, (GLenum)format, (GLenum)t_ype, (const GLvoid *)pixels); }
//
// void (APIENTRYP pfn_glTexSubImage2D)(GLenum target, GLint level, GLint xoffset, GLint yoffset, GLsizei width, GLsizei height, GLenum format, GLenum type, const GLvoid *pixels);
// GLAPI void APIENTRY wrap_glTexSubImage2D(unsigned int target, int level, int xoffset, int yoffset, int width, int height, unsigned int format, unsigned int t_ype, const void* pixels) {  (*pfn_glTexSubImage2D)((GLenum)target, (GLint)level, (GLint)xoffset, (GLint)yoffset, (GLsizei)width, (GLsizei)height, (GLenum)format, (GLenum)t_ype, (const GLvoid *)pixels); }
//
// void (APIENTRYP pfn_glTexSubImage3D)(GLenum target, GLint level, GLint xoffset, GLint yoffset, GLint zoffset, GLsizei width, GLsizei height, GLsizei depth, GLenum format, GLenum type, const GLvoid *pixels);
// GLAPI void APIENTRY wrap_glTexSubImage3D(unsigned int target, int level, int xoffset, int yoffset, int zoffset, int width, int height, int depth, unsigned int format, unsigned int t_ype, const void* pixels) {  (*pfn_glTexSubImage3D)((GLenum)target, (GLint)level, (GLint)xoffset, (GLint)yoffset, (GLint)zoffset, (GLsizei)width, (GLsizei)height, (GLsizei)depth, (GLenum)format, (GLenum)t_ype, (const GLvoid *)pixels); }
//
// void (APIENTRYP pfn_glTextureBufferRangeEXT)(GLuint texture, GLenum target, GLenum internalformat, GLuint buffer, GLintptr offset, GLsizeiptr size);
// GLAPI void APIENTRY wrap_glTextureBufferRangeEXT(unsigned int texture, unsigned int target, unsigned int internalformat, unsigned int buffer, long long offset, long long size) {  (*pfn_glTextureBufferRangeEXT)((GLuint)texture, (GLenum)target, (GLenum)internalformat, (GLuint)buffer, (GLintptr)offset, (GLsizeiptr)size); }
//
// void (APIENTRYP pfn_glTextureStorage1DEXT)(GLuint texture, GLenum target, GLsizei levels, GLenum internalformat, GLsizei width);
// GLAPI void APIENTRY wrap_glTextureStorage1DEXT(unsigned int texture, unsigned int target, int levels, unsigned int internalformat, int width) {  (*pfn_glTextureStorage1DEXT)((GLuint)texture, (GLenum)target, (GLsizei)levels, (GLenum)internalformat, (GLsizei)width); }
//
// void (APIENTRYP pfn_glTextureStorage2DEXT)(GLuint texture, GLenum target, GLsizei levels, GLenum internalformat, GLsizei width, GLsizei height);
// GLAPI void APIENTRY wrap_glTextureStorage2DEXT(unsigned int texture, unsigned int target, int levels, unsigned int internalformat, int width, int height) {  (*pfn_glTextureStorage2DEXT)((GLuint)texture, (GLenum)target, (GLsizei)levels, (GLenum)internalformat, (GLsizei)width, (GLsizei)height); }
//
// void (APIENTRYP pfn_glTextureStorage2DMultisampleEXT)(GLuint texture, GLenum target, GLsizei samples, GLenum internalformat, GLsizei width, GLsizei height, GLboolean fixedsamplelocations);
// GLAPI void APIENTRY wrap_glTextureStorage2DMultisampleEXT(unsigned int texture, unsigned int target, int samples, unsigned int internalformat, int width, int height, unsigned char fixedsamplelocations) {  (*pfn_glTextureStorage2DMultisampleEXT)((GLuint)texture, (GLenum)target, (GLsizei)samples, (GLenum)internalformat, (GLsizei)width, (GLsizei)height, (GLboolean)fixedsamplelocations); }
//
// void (APIENTRYP pfn_glTextureStorage3DEXT)(GLuint texture, GLenum target, GLsizei levels, GLenum internalformat, GLsizei width, GLsizei height, GLsizei depth);
// GLAPI void APIENTRY wrap_glTextureStorage3DEXT(unsigned int texture, unsigned int target, int levels, unsigned int internalformat, int width, int height, int depth) {  (*pfn_glTextureStorage3DEXT)((GLuint)texture, (GLenum)target, (GLsizei)levels, (GLenum)internalformat, (GLsizei)width, (GLsizei)height, (GLsizei)depth); }
//
// void (APIENTRYP pfn_glTextureStorage3DMultisampleEXT)(GLuint texture, GLenum target, GLsizei samples, GLenum internalformat, GLsizei width, GLsizei height, GLsizei depth, GLboolean fixedsamplelocations);
// GLAPI void APIENTRY wrap_glTextureStorage3DMultisampleEXT(unsigned int texture, unsigned int target, int samples, unsigned int internalformat, int width, int height, int depth, unsigned char fixedsamplelocations) {  (*pfn_glTextureStorage3DMultisampleEXT)((GLuint)texture, (GLenum)target, (GLsizei)samples, (GLenum)internalformat, (GLsizei)width, (GLsizei)height, (GLsizei)depth, (GLboolean)fixedsamplelocations); }
//
// void (APIENTRYP pfn_glTextureView)(GLuint texture, GLenum target, GLuint origtexture, GLenum internalformat, GLuint minlevel, GLuint numlevels, GLuint minlayer, GLuint numlayers);
// GLAPI void APIENTRY wrap_glTextureView(unsigned int texture, unsigned int target, unsigned int origtexture, unsigned int internalformat, unsigned int minlevel, unsigned int numlevels, unsigned int minlayer, unsigned int numlayers) {  (*pfn_glTextureView)((GLuint)texture, (GLenum)target, (GLuint)origtexture, (GLenum)internalformat, (GLuint)minlevel, (GLuint)numlevels, (GLuint)minlayer, (GLuint)numlayers); }
//
// void (APIENTRYP pfn_glTransformFeedbackVaryings)(GLuint program, GLsizei count, const GLchar* const *varyings, GLenum bufferMode);
// GLAPI void APIENTRY wrap_glTransformFeedbackVaryings(unsigned int program, int count, const char* const* varyings, unsigned int bufferMode) {  (*pfn_glTransformFeedbackVaryings)((GLuint)program, (GLsizei)count, (const GLchar* const *)varyings, (GLenum)bufferMode); }
//
// void (APIENTRYP pfn_glUniform1d)(GLint location, GLdouble x);
// GLAPI void APIENTRY wrap_glUniform1d(int location, double x) {  (*pfn_glUniform1d)((GLint)location, (GLdouble)x); }
//
// void (APIENTRYP pfn_glUniform1dv)(GLint location, GLsizei count, const GLdouble *value);
// GLAPI void APIENTRY wrap_glUniform1dv(int location, int count, const double* value) {  (*pfn_glUniform1dv)((GLint)location, (GLsizei)count, (const GLdouble *)value); }
//
// void (APIENTRYP pfn_glUniform1f)(GLint location, GLfloat v0);
// GLAPI void APIENTRY wrap_glUniform1f(int location, float v0) {  (*pfn_glUniform1f)((GLint)location, (GLfloat)v0); }
//
// void (APIENTRYP pfn_glUniform1fv)(GLint location, GLsizei count, const GLfloat *value);
// GLAPI void APIENTRY wrap_glUniform1fv(int location, int count, const float* value) {  (*pfn_glUniform1fv)((GLint)location, (GLsizei)count, (const GLfloat *)value); }
//
// void (APIENTRYP pfn_glUniform1i)(GLint location, GLint v0);
// GLAPI void APIENTRY wrap_glUniform1i(int location, int v0) {  (*pfn_glUniform1i)((GLint)location, (GLint)v0); }
//
// void (APIENTRYP pfn_glUniform1iv)(GLint location, GLsizei count, const GLint *value);
// GLAPI void APIENTRY wrap_glUniform1iv(int location, int count, const int* value) {  (*pfn_glUniform1iv)((GLint)location, (GLsizei)count, (const GLint *)value); }
//
// void (APIENTRYP pfn_glUniform1ui)(GLint location, GLuint v0);
// GLAPI void APIENTRY wrap_glUniform1ui(int location, unsigned int v0) {  (*pfn_glUniform1ui)((GLint)location, (GLuint)v0); }
//
// void (APIENTRYP pfn_glUniform1uiv)(GLint location, GLsizei count, const GLuint *value);
// GLAPI void APIENTRY wrap_glUniform1uiv(int location, int count, const unsigned int* value) {  (*pfn_glUniform1uiv)((GLint)location, (GLsizei)count, (const GLuint *)value); }
//
// void (APIENTRYP pfn_glUniform2d)(GLint location, GLdouble x, GLdouble y);
// GLAPI void APIENTRY wrap_glUniform2d(int location, double x, double y) {  (*pfn_glUniform2d)((GLint)location, (GLdouble)x, (GLdouble)y); }
//
// void (APIENTRYP pfn_glUniform2dv)(GLint location, GLsizei count, const GLdouble *value);
// GLAPI void APIENTRY wrap_glUniform2dv(int location, int count, const double* value) {  (*pfn_glUniform2dv)((GLint)location, (GLsizei)count, (const GLdouble *)value); }
//
// void (APIENTRYP pfn_glUniform2f)(GLint location, GLfloat v0, GLfloat v1);
// GLAPI void APIENTRY wrap_glUniform2f(int location, float v0, float v1) {  (*pfn_glUniform2f)((GLint)location, (GLfloat)v0, (GLfloat)v1); }
//
// void (APIENTRYP pfn_glUniform2fv)(GLint location, GLsizei count, const GLfloat *value);
// GLAPI void APIENTRY wrap_glUniform2fv(int location, int count, const float* value) {  (*pfn_glUniform2fv)((GLint)location, (GLsizei)count, (const GLfloat *)value); }
//
// void (APIENTRYP pfn_glUniform2i)(GLint location, GLint v0, GLint v1);
// GLAPI void APIENTRY wrap_glUniform2i(int location, int v0, int v1) {  (*pfn_glUniform2i)((GLint)location, (GLint)v0, (GLint)v1); }
//
// void (APIENTRYP pfn_glUniform2iv)(GLint location, GLsizei count, const GLint *value);
// GLAPI void APIENTRY wrap_glUniform2iv(int location, int count, const int* value) {  (*pfn_glUniform2iv)((GLint)location, (GLsizei)count, (const GLint *)value); }
//
// void (APIENTRYP pfn_glUniform2ui)(GLint location, GLuint v0, GLuint v1);
// GLAPI void APIENTRY wrap_glUniform2ui(int location, unsigned int v0, unsigned int v1) {  (*pfn_glUniform2ui)((GLint)location, (GLuint)v0, (GLuint)v1); }
//
// void (APIENTRYP pfn_glUniform2uiv)(GLint location, GLsizei count, const GLuint *value);
// GLAPI void APIENTRY wrap_glUniform2uiv(int location, int count, const unsigned int* value) {  (*pfn_glUniform2uiv)((GLint)location, (GLsizei)count, (const GLuint *)value); }
//
// void (APIENTRYP pfn_glUniform3d)(GLint location, GLdouble x, GLdouble y, GLdouble z);
// GLAPI void APIENTRY wrap_glUniform3d(int location, double x, double y, double z) {  (*pfn_glUniform3d)((GLint)location, (GLdouble)x, (GLdouble)y, (GLdouble)z); }
//
// void (APIENTRYP pfn_glUniform3dv)(GLint location, GLsizei count, const GLdouble *value);
// GLAPI void APIENTRY wrap_glUniform3dv(int location, int count, const double* value) {  (*pfn_glUniform3dv)((GLint)location, (GLsizei)count, (const GLdouble *)value); }
//
// void (APIENTRYP pfn_glUniform3f)(GLint location, GLfloat v0, GLfloat v1, GLfloat v2);
// GLAPI void APIENTRY wrap_glUniform3f(int location, float v0, float v1, float v2) {  (*pfn_glUniform3f)((GLint)location, (GLfloat)v0, (GLfloat)v1, (GLfloat)v2); }
//
// void (APIENTRYP pfn_glUniform3fv)(GLint location, GLsizei count, const GLfloat *value);
// GLAPI void APIENTRY wrap_glUniform3fv(int location, int count, const float* value) {  (*pfn_glUniform3fv)((GLint)location, (GLsizei)count, (const GLfloat *)value); }
//
// void (APIENTRYP pfn_glUniform3i)(GLint location, GLint v0, GLint v1, GLint v2);
// GLAPI void APIENTRY wrap_glUniform3i(int location, int v0, int v1, int v2) {  (*pfn_glUniform3i)((GLint)location, (GLint)v0, (GLint)v1, (GLint)v2); }
//
// void (APIENTRYP pfn_glUniform3iv)(GLint location, GLsizei count, const GLint *value);
// GLAPI void APIENTRY wrap_glUniform3iv(int location, int count, const int* value) {  (*pfn_glUniform3iv)((GLint)location, (GLsizei)count, (const GLint *)value); }
//
// void (APIENTRYP pfn_glUniform3ui)(GLint location, GLuint v0, GLuint v1, GLuint v2);
// GLAPI void APIENTRY wrap_glUniform3ui(int location, unsigned int v0, unsigned int v1, unsigned int v2) {  (*pfn_glUniform3ui)((GLint)location, (GLuint)v0, (GLuint)v1, (GLuint)v2); }
//
// void (APIENTRYP pfn_glUniform3uiv)(GLint location, GLsizei count, const GLuint *value);
// GLAPI void APIENTRY wrap_glUniform3uiv(int location, int count, const unsigned int* value) {  (*pfn_glUniform3uiv)((GLint)location, (GLsizei)count, (const GLuint *)value); }
//
// void (APIENTRYP pfn_glUniform4d)(GLint location, GLdouble x, GLdouble y, GLdouble z, GLdouble w);
// GLAPI void APIENTRY wrap_glUniform4d(int location, double x, double y, double z, double w) {  (*pfn_glUniform4d)((GLint)location, (GLdouble)x, (GLdouble)y, (GLdouble)z, (GLdouble)w); }
//
// void (APIENTRYP pfn_glUniform4dv)(GLint location, GLsizei count, const GLdouble *value);
// GLAPI void APIENTRY wrap_glUniform4dv(int location, int count, const double* value) {  (*pfn_glUniform4dv)((GLint)location, (GLsizei)count, (const GLdouble *)value); }
//
// void (APIENTRYP pfn_glUniform4f)(GLint location, GLfloat v0, GLfloat v1, GLfloat v2, GLfloat v3);
// GLAPI void APIENTRY wrap_glUniform4f(int location, float v0, float v1, float v2, float v3) {  (*pfn_glUniform4f)((GLint)location, (GLfloat)v0, (GLfloat)v1, (GLfloat)v2, (GLfloat)v3); }
//
// void (APIENTRYP pfn_glUniform4fv)(GLint location, GLsizei count, const GLfloat *value);
// GLAPI void APIENTRY wrap_glUniform4fv(int location, int count, const float* value) {  (*pfn_glUniform4fv)((GLint)location, (GLsizei)count, (const GLfloat *)value); }
//
// void (APIENTRYP pfn_glUniform4i)(GLint location, GLint v0, GLint v1, GLint v2, GLint v3);
// GLAPI void APIENTRY wrap_glUniform4i(int location, int v0, int v1, int v2, int v3) {  (*pfn_glUniform4i)((GLint)location, (GLint)v0, (GLint)v1, (GLint)v2, (GLint)v3); }
//
// void (APIENTRYP pfn_glUniform4iv)(GLint location, GLsizei count, const GLint *value);
// GLAPI void APIENTRY wrap_glUniform4iv(int location, int count, const int* value) {  (*pfn_glUniform4iv)((GLint)location, (GLsizei)count, (const GLint *)value); }
//
// void (APIENTRYP pfn_glUniform4ui)(GLint location, GLuint v0, GLuint v1, GLuint v2, GLuint v3);
// GLAPI void APIENTRY wrap_glUniform4ui(int location, unsigned int v0, unsigned int v1, unsigned int v2, unsigned int v3) {  (*pfn_glUniform4ui)((GLint)location, (GLuint)v0, (GLuint)v1, (GLuint)v2, (GLuint)v3); }
//
// void (APIENTRYP pfn_glUniform4uiv)(GLint location, GLsizei count, const GLuint *value);
// GLAPI void APIENTRY wrap_glUniform4uiv(int location, int count, const unsigned int* value) {  (*pfn_glUniform4uiv)((GLint)location, (GLsizei)count, (const GLuint *)value); }
//
// void (APIENTRYP pfn_glUniformBlockBinding)(GLuint program, GLuint uniformBlockIndex, GLuint uniformBlockBinding);
// GLAPI void APIENTRY wrap_glUniformBlockBinding(unsigned int program, unsigned int uniformBlockIndex, unsigned int uniformBlockBinding) {  (*pfn_glUniformBlockBinding)((GLuint)program, (GLuint)uniformBlockIndex, (GLuint)uniformBlockBinding); }
//
// void (APIENTRYP pfn_glUniformMatrix2dv)(GLint location, GLsizei count, GLboolean transpose, const GLdouble *value);
// GLAPI void APIENTRY wrap_glUniformMatrix2dv(int location, int count, unsigned char transpose, const double* value) {  (*pfn_glUniformMatrix2dv)((GLint)location, (GLsizei)count, (GLboolean)transpose, (const GLdouble *)value); }
//
// void (APIENTRYP pfn_glUniformMatrix2fv)(GLint location, GLsizei count, GLboolean transpose, const GLfloat *value);
// GLAPI void APIENTRY wrap_glUniformMatrix2fv(int location, int count, unsigned char transpose, const float* value) {  (*pfn_glUniformMatrix2fv)((GLint)location, (GLsizei)count, (GLboolean)transpose, (const GLfloat *)value); }
//
// void (APIENTRYP pfn_glUniformMatrix2x3dv)(GLint location, GLsizei count, GLboolean transpose, const GLdouble *value);
// GLAPI void APIENTRY wrap_glUniformMatrix2x3dv(int location, int count, unsigned char transpose, const double* value) {  (*pfn_glUniformMatrix2x3dv)((GLint)location, (GLsizei)count, (GLboolean)transpose, (const GLdouble *)value); }
//
// void (APIENTRYP pfn_glUniformMatrix2x3fv)(GLint location, GLsizei count, GLboolean transpose, const GLfloat *value);
// GLAPI void APIENTRY wrap_glUniformMatrix2x3fv(int location, int count, unsigned char transpose, const float* value) {  (*pfn_glUniformMatrix2x3fv)((GLint)location, (GLsizei)count, (GLboolean)transpose, (const GLfloat *)value); }
//
// void (APIENTRYP pfn_glUniformMatrix2x4dv)(GLint location, GLsizei count, GLboolean transpose, const GLdouble *value);
// GLAPI void APIENTRY wrap_glUniformMatrix2x4dv(int location, int count, unsigned char transpose, const double* value) {  (*pfn_glUniformMatrix2x4dv)((GLint)location, (GLsizei)count, (GLboolean)transpose, (const GLdouble *)value); }
//
// void (APIENTRYP pfn_glUniformMatrix2x4fv)(GLint location, GLsizei count, GLboolean transpose, const GLfloat *value);
// GLAPI void APIENTRY wrap_glUniformMatrix2x4fv(int location, int count, unsigned char transpose, const float* value) {  (*pfn_glUniformMatrix2x4fv)((GLint)location, (GLsizei)count, (GLboolean)transpose, (const GLfloat *)value); }
//
// void (APIENTRYP pfn_glUniformMatrix3dv)(GLint location, GLsizei count, GLboolean transpose, const GLdouble *value);
// GLAPI void APIENTRY wrap_glUniformMatrix3dv(int location, int count, unsigned char transpose, const double* value) {  (*pfn_glUniformMatrix3dv)((GLint)location, (GLsizei)count, (GLboolean)transpose, (const GLdouble *)value); }
//
// void (APIENTRYP pfn_glUniformMatrix3fv)(GLint location, GLsizei count, GLboolean transpose, const GLfloat *value);
// GLAPI void APIENTRY wrap_glUniformMatrix3fv(int location, int count, unsigned char transpose, const float* value) {  (*pfn_glUniformMatrix3fv)((GLint)location, (GLsizei)count, (GLboolean)transpose, (const GLfloat *)value); }
//
// void (APIENTRYP pfn_glUniformMatrix3x2dv)(GLint location, GLsizei count, GLboolean transpose, const GLdouble *value);
// GLAPI void APIENTRY wrap_glUniformMatrix3x2dv(int location, int count, unsigned char transpose, const double* value) {  (*pfn_glUniformMatrix3x2dv)((GLint)location, (GLsizei)count, (GLboolean)transpose, (const GLdouble *)value); }
//
// void (APIENTRYP pfn_glUniformMatrix3x2fv)(GLint location, GLsizei count, GLboolean transpose, const GLfloat *value);
// GLAPI void APIENTRY wrap_glUniformMatrix3x2fv(int location, int count, unsigned char transpose, const float* value) {  (*pfn_glUniformMatrix3x2fv)((GLint)location, (GLsizei)count, (GLboolean)transpose, (const GLfloat *)value); }
//
// void (APIENTRYP pfn_glUniformMatrix3x4dv)(GLint location, GLsizei count, GLboolean transpose, const GLdouble *value);
// GLAPI void APIENTRY wrap_glUniformMatrix3x4dv(int location, int count, unsigned char transpose, const double* value) {  (*pfn_glUniformMatrix3x4dv)((GLint)location, (GLsizei)count, (GLboolean)transpose, (const GLdouble *)value); }
//
// void (APIENTRYP pfn_glUniformMatrix3x4fv)(GLint location, GLsizei count, GLboolean transpose, const GLfloat *value);
// GLAPI void APIENTRY wrap_glUniformMatrix3x4fv(int location, int count, unsigned char transpose, const float* value) {  (*pfn_glUniformMatrix3x4fv)((GLint)location, (GLsizei)count, (GLboolean)transpose, (const GLfloat *)value); }
//
// void (APIENTRYP pfn_glUniformMatrix4dv)(GLint location, GLsizei count, GLboolean transpose, const GLdouble *value);
// GLAPI void APIENTRY wrap_glUniformMatrix4dv(int location, int count, unsigned char transpose, const double* value) {  (*pfn_glUniformMatrix4dv)((GLint)location, (GLsizei)count, (GLboolean)transpose, (const GLdouble *)value); }
//
// void (APIENTRYP pfn_glUniformMatrix4fv)(GLint location, GLsizei count, GLboolean transpose, const GLfloat *value);
// GLAPI void APIENTRY wrap_glUniformMatrix4fv(int location, int count, unsigned char transpose, const float* value) {  (*pfn_glUniformMatrix4fv)((GLint)location, (GLsizei)count, (GLboolean)transpose, (const GLfloat *)value); }
//
// void (APIENTRYP pfn_glUniformMatrix4x2dv)(GLint location, GLsizei count, GLboolean transpose, const GLdouble *value);
// GLAPI void APIENTRY wrap_glUniformMatrix4x2dv(int location, int count, unsigned char transpose, const double* value) {  (*pfn_glUniformMatrix4x2dv)((GLint)location, (GLsizei)count, (GLboolean)transpose, (const GLdouble *)value); }
//
// void (APIENTRYP pfn_glUniformMatrix4x2fv)(GLint location, GLsizei count, GLboolean transpose, const GLfloat *value);
// GLAPI void APIENTRY wrap_glUniformMatrix4x2fv(int location, int count, unsigned char transpose, const float* value) {  (*pfn_glUniformMatrix4x2fv)((GLint)location, (GLsizei)count, (GLboolean)transpose, (const GLfloat *)value); }
//
// void (APIENTRYP pfn_glUniformMatrix4x3dv)(GLint location, GLsizei count, GLboolean transpose, const GLdouble *value);
// GLAPI void APIENTRY wrap_glUniformMatrix4x3dv(int location, int count, unsigned char transpose, const double* value) {  (*pfn_glUniformMatrix4x3dv)((GLint)location, (GLsizei)count, (GLboolean)transpose, (const GLdouble *)value); }
//
// void (APIENTRYP pfn_glUniformMatrix4x3fv)(GLint location, GLsizei count, GLboolean transpose, const GLfloat *value);
// GLAPI void APIENTRY wrap_glUniformMatrix4x3fv(int location, int count, unsigned char transpose, const float* value) {  (*pfn_glUniformMatrix4x3fv)((GLint)location, (GLsizei)count, (GLboolean)transpose, (const GLfloat *)value); }
//
// void (APIENTRYP pfn_glUniformSubroutinesuiv)(GLenum shadertype, GLsizei count, const GLuint *indices);
// GLAPI void APIENTRY wrap_glUniformSubroutinesuiv(unsigned int shadertype, int count, const unsigned int* indices) {  (*pfn_glUniformSubroutinesuiv)((GLenum)shadertype, (GLsizei)count, (const GLuint *)indices); }
//
// void (APIENTRYP pfn_glUseProgram)(GLuint program);
// GLAPI void APIENTRY wrap_glUseProgram(unsigned int program) {  (*pfn_glUseProgram)((GLuint)program); }
//
// void (APIENTRYP pfn_glUseProgramStages)(GLuint pipeline, GLbitfield stages, GLuint program);
// GLAPI void APIENTRY wrap_glUseProgramStages(unsigned int pipeline, unsigned int stages, unsigned int program) {  (*pfn_glUseProgramStages)((GLuint)pipeline, (GLbitfield)stages, (GLuint)program); }
//
// void (APIENTRYP pfn_glValidateProgram)(GLuint program);
// GLAPI void APIENTRY wrap_glValidateProgram(unsigned int program) {  (*pfn_glValidateProgram)((GLuint)program); }
//
// void (APIENTRYP pfn_glValidateProgramPipeline)(GLuint pipeline);
// GLAPI void APIENTRY wrap_glValidateProgramPipeline(unsigned int pipeline) {  (*pfn_glValidateProgramPipeline)((GLuint)pipeline); }
//
// void (APIENTRYP pfn_glVertexArrayBindVertexBufferEXT)(GLuint vaobj, GLuint bindingindex, GLuint buffer, GLintptr offset, GLsizei stride);
// GLAPI void APIENTRY wrap_glVertexArrayBindVertexBufferEXT(unsigned int vaobj, unsigned int bindingindex, unsigned int buffer, long long offset, int stride) {  (*pfn_glVertexArrayBindVertexBufferEXT)((GLuint)vaobj, (GLuint)bindingindex, (GLuint)buffer, (GLintptr)offset, (GLsizei)stride); }
//
// void (APIENTRYP pfn_glVertexArrayVertexAttribBindingEXT)(GLuint vaobj, GLuint attribindex, GLuint bindingindex);
// GLAPI void APIENTRY wrap_glVertexArrayVertexAttribBindingEXT(unsigned int vaobj, unsigned int attribindex, unsigned int bindingindex) {  (*pfn_glVertexArrayVertexAttribBindingEXT)((GLuint)vaobj, (GLuint)attribindex, (GLuint)bindingindex); }
//
// void (APIENTRYP pfn_glVertexArrayVertexAttribFormatEXT)(GLuint vaobj, GLuint attribindex, GLint size, GLenum type, GLboolean normalized, GLuint relativeoffset);
// GLAPI void APIENTRY wrap_glVertexArrayVertexAttribFormatEXT(unsigned int vaobj, unsigned int attribindex, int size, unsigned int t_ype, unsigned char normalized, unsigned int relativeoffset) {  (*pfn_glVertexArrayVertexAttribFormatEXT)((GLuint)vaobj, (GLuint)attribindex, (GLint)size, (GLenum)t_ype, (GLboolean)normalized, (GLuint)relativeoffset); }
//
// void (APIENTRYP pfn_glVertexArrayVertexAttribIFormatEXT)(GLuint vaobj, GLuint attribindex, GLint size, GLenum type, GLuint relativeoffset);
// GLAPI void APIENTRY wrap_glVertexArrayVertexAttribIFormatEXT(unsigned int vaobj, unsigned int attribindex, int size, unsigned int t_ype, unsigned int relativeoffset) {  (*pfn_glVertexArrayVertexAttribIFormatEXT)((GLuint)vaobj, (GLuint)attribindex, (GLint)size, (GLenum)t_ype, (GLuint)relativeoffset); }
//
// void (APIENTRYP pfn_glVertexArrayVertexAttribLFormatEXT)(GLuint vaobj, GLuint attribindex, GLint size, GLenum type, GLuint relativeoffset);
// GLAPI void APIENTRY wrap_glVertexArrayVertexAttribLFormatEXT(unsigned int vaobj, unsigned int attribindex, int size, unsigned int t_ype, unsigned int relativeoffset) {  (*pfn_glVertexArrayVertexAttribLFormatEXT)((GLuint)vaobj, (GLuint)attribindex, (GLint)size, (GLenum)t_ype, (GLuint)relativeoffset); }
//
// void (APIENTRYP pfn_glVertexArrayVertexBindingDivisorEXT)(GLuint vaobj, GLuint bindingindex, GLuint divisor);
// GLAPI void APIENTRY wrap_glVertexArrayVertexBindingDivisorEXT(unsigned int vaobj, unsigned int bindingindex, unsigned int divisor) {  (*pfn_glVertexArrayVertexBindingDivisorEXT)((GLuint)vaobj, (GLuint)bindingindex, (GLuint)divisor); }
//
// void (APIENTRYP pfn_glVertexAttrib1d)(GLuint index, GLdouble x);
// GLAPI void APIENTRY wrap_glVertexAttrib1d(unsigned int index, double x) {  (*pfn_glVertexAttrib1d)((GLuint)index, (GLdouble)x); }
//
// void (APIENTRYP pfn_glVertexAttrib1dv)(GLuint index, const GLdouble *v);
// GLAPI void APIENTRY wrap_glVertexAttrib1dv(unsigned int index, const double* v) {  (*pfn_glVertexAttrib1dv)((GLuint)index, (const GLdouble *)v); }
//
// void (APIENTRYP pfn_glVertexAttrib1f)(GLuint index, GLfloat x);
// GLAPI void APIENTRY wrap_glVertexAttrib1f(unsigned int index, float x) {  (*pfn_glVertexAttrib1f)((GLuint)index, (GLfloat)x); }
//
// void (APIENTRYP pfn_glVertexAttrib1fv)(GLuint index, const GLfloat *v);
// GLAPI void APIENTRY wrap_glVertexAttrib1fv(unsigned int index, const float* v) {  (*pfn_glVertexAttrib1fv)((GLuint)index, (const GLfloat *)v); }
//
// void (APIENTRYP pfn_glVertexAttrib1s)(GLuint index, GLshort x);
// GLAPI void APIENTRY wrap_glVertexAttrib1s(unsigned int index, short  x) {  (*pfn_glVertexAttrib1s)((GLuint)index, (GLshort)x); }
//
// void (APIENTRYP pfn_glVertexAttrib1sv)(GLuint index, const GLshort *v);
// GLAPI void APIENTRY wrap_glVertexAttrib1sv(unsigned int index, const short * v) {  (*pfn_glVertexAttrib1sv)((GLuint)index, (const GLshort *)v); }
//
// void (APIENTRYP pfn_glVertexAttrib2d)(GLuint index, GLdouble x, GLdouble y);
// GLAPI void APIENTRY wrap_glVertexAttrib2d(unsigned int index, double x, double y) {  (*pfn_glVertexAttrib2d)((GLuint)index, (GLdouble)x, (GLdouble)y); }
//
// void (APIENTRYP pfn_glVertexAttrib2dv)(GLuint index, const GLdouble *v);
// GLAPI void APIENTRY wrap_glVertexAttrib2dv(unsigned int index, const double* v) {  (*pfn_glVertexAttrib2dv)((GLuint)index, (const GLdouble *)v); }
//
// void (APIENTRYP pfn_glVertexAttrib2f)(GLuint index, GLfloat x, GLfloat y);
// GLAPI void APIENTRY wrap_glVertexAttrib2f(unsigned int index, float x, float y) {  (*pfn_glVertexAttrib2f)((GLuint)index, (GLfloat)x, (GLfloat)y); }
//
// void (APIENTRYP pfn_glVertexAttrib2fv)(GLuint index, const GLfloat *v);
// GLAPI void APIENTRY wrap_glVertexAttrib2fv(unsigned int index, const float* v) {  (*pfn_glVertexAttrib2fv)((GLuint)index, (const GLfloat *)v); }
//
// void (APIENTRYP pfn_glVertexAttrib2s)(GLuint index, GLshort x, GLshort y);
// GLAPI void APIENTRY wrap_glVertexAttrib2s(unsigned int index, short  x, short  y) {  (*pfn_glVertexAttrib2s)((GLuint)index, (GLshort)x, (GLshort)y); }
//
// void (APIENTRYP pfn_glVertexAttrib2sv)(GLuint index, const GLshort *v);
// GLAPI void APIENTRY wrap_glVertexAttrib2sv(unsigned int index, const short * v) {  (*pfn_glVertexAttrib2sv)((GLuint)index, (const GLshort *)v); }
//
// void (APIENTRYP pfn_glVertexAttrib3d)(GLuint index, GLdouble x, GLdouble y, GLdouble z);
// GLAPI void APIENTRY wrap_glVertexAttrib3d(unsigned int index, double x, double y, double z) {  (*pfn_glVertexAttrib3d)((GLuint)index, (GLdouble)x, (GLdouble)y, (GLdouble)z); }
//
// void (APIENTRYP pfn_glVertexAttrib3dv)(GLuint index, const GLdouble *v);
// GLAPI void APIENTRY wrap_glVertexAttrib3dv(unsigned int index, const double* v) {  (*pfn_glVertexAttrib3dv)((GLuint)index, (const GLdouble *)v); }
//
// void (APIENTRYP pfn_glVertexAttrib3f)(GLuint index, GLfloat x, GLfloat y, GLfloat z);
// GLAPI void APIENTRY wrap_glVertexAttrib3f(unsigned int index, float x, float y, float z) {  (*pfn_glVertexAttrib3f)((GLuint)index, (GLfloat)x, (GLfloat)y, (GLfloat)z); }
//
// void (APIENTRYP pfn_glVertexAttrib3fv)(GLuint index, const GLfloat *v);
// GLAPI void APIENTRY wrap_glVertexAttrib3fv(unsigned int index, const float* v) {  (*pfn_glVertexAttrib3fv)((GLuint)index, (const GLfloat *)v); }
//
// void (APIENTRYP pfn_glVertexAttrib3s)(GLuint index, GLshort x, GLshort y, GLshort z);
// GLAPI void APIENTRY wrap_glVertexAttrib3s(unsigned int index, short  x, short  y, short  z) {  (*pfn_glVertexAttrib3s)((GLuint)index, (GLshort)x, (GLshort)y, (GLshort)z); }
//
// void (APIENTRYP pfn_glVertexAttrib3sv)(GLuint index, const GLshort *v);
// GLAPI void APIENTRY wrap_glVertexAttrib3sv(unsigned int index, const short * v) {  (*pfn_glVertexAttrib3sv)((GLuint)index, (const GLshort *)v); }
//
// void (APIENTRYP pfn_glVertexAttrib4Nbv)(GLuint index, const GLbyte *v);
// GLAPI void APIENTRY wrap_glVertexAttrib4Nbv(unsigned int index, const signed char* v) {  (*pfn_glVertexAttrib4Nbv)((GLuint)index, (const GLbyte *)v); }
//
// void (APIENTRYP pfn_glVertexAttrib4Niv)(GLuint index, const GLint *v);
// GLAPI void APIENTRY wrap_glVertexAttrib4Niv(unsigned int index, const int* v) {  (*pfn_glVertexAttrib4Niv)((GLuint)index, (const GLint *)v); }
//
// void (APIENTRYP pfn_glVertexAttrib4Nsv)(GLuint index, const GLshort *v);
// GLAPI void APIENTRY wrap_glVertexAttrib4Nsv(unsigned int index, const short * v) {  (*pfn_glVertexAttrib4Nsv)((GLuint)index, (const GLshort *)v); }
//
// void (APIENTRYP pfn_glVertexAttrib4Nub)(GLuint index, GLubyte x, GLubyte y, GLubyte z, GLubyte w);
// GLAPI void APIENTRY wrap_glVertexAttrib4Nub(unsigned int index, unsigned char x, unsigned char y, unsigned char z, unsigned char w) {  (*pfn_glVertexAttrib4Nub)((GLuint)index, (GLubyte)x, (GLubyte)y, (GLubyte)z, (GLubyte)w); }
//
// void (APIENTRYP pfn_glVertexAttrib4Nubv)(GLuint index, const GLubyte *v);
// GLAPI void APIENTRY wrap_glVertexAttrib4Nubv(unsigned int index, const unsigned char* v) {  (*pfn_glVertexAttrib4Nubv)((GLuint)index, (const GLubyte *)v); }
//
// void (APIENTRYP pfn_glVertexAttrib4Nuiv)(GLuint index, const GLuint *v);
// GLAPI void APIENTRY wrap_glVertexAttrib4Nuiv(unsigned int index, const unsigned int* v) {  (*pfn_glVertexAttrib4Nuiv)((GLuint)index, (const GLuint *)v); }
//
// void (APIENTRYP pfn_glVertexAttrib4Nusv)(GLuint index, const GLushort *v);
// GLAPI void APIENTRY wrap_glVertexAttrib4Nusv(unsigned int index, const unsigned short* v) {  (*pfn_glVertexAttrib4Nusv)((GLuint)index, (const GLushort *)v); }
//
// void (APIENTRYP pfn_glVertexAttrib4bv)(GLuint index, const GLbyte *v);
// GLAPI void APIENTRY wrap_glVertexAttrib4bv(unsigned int index, const signed char* v) {  (*pfn_glVertexAttrib4bv)((GLuint)index, (const GLbyte *)v); }
//
// void (APIENTRYP pfn_glVertexAttrib4d)(GLuint index, GLdouble x, GLdouble y, GLdouble z, GLdouble w);
// GLAPI void APIENTRY wrap_glVertexAttrib4d(unsigned int index, double x, double y, double z, double w) {  (*pfn_glVertexAttrib4d)((GLuint)index, (GLdouble)x, (GLdouble)y, (GLdouble)z, (GLdouble)w); }
//
// void (APIENTRYP pfn_glVertexAttrib4dv)(GLuint index, const GLdouble *v);
// GLAPI void APIENTRY wrap_glVertexAttrib4dv(unsigned int index, const double* v) {  (*pfn_glVertexAttrib4dv)((GLuint)index, (const GLdouble *)v); }
//
// void (APIENTRYP pfn_glVertexAttrib4f)(GLuint index, GLfloat x, GLfloat y, GLfloat z, GLfloat w);
// GLAPI void APIENTRY wrap_glVertexAttrib4f(unsigned int index, float x, float y, float z, float w) {  (*pfn_glVertexAttrib4f)((GLuint)index, (GLfloat)x, (GLfloat)y, (GLfloat)z, (GLfloat)w); }
//
// void (APIENTRYP pfn_glVertexAttrib4fv)(GLuint index, const GLfloat *v);
// GLAPI void APIENTRY wrap_glVertexAttrib4fv(unsigned int index, const float* v) {  (*pfn_glVertexAttrib4fv)((GLuint)index, (const GLfloat *)v); }
//
// void (APIENTRYP pfn_glVertexAttrib4iv)(GLuint index, const GLint *v);
// GLAPI void APIENTRY wrap_glVertexAttrib4iv(unsigned int index, const int* v) {  (*pfn_glVertexAttrib4iv)((GLuint)index, (const GLint *)v); }
//
// void (APIENTRYP pfn_glVertexAttrib4s)(GLuint index, GLshort x, GLshort y, GLshort z, GLshort w);
// GLAPI void APIENTRY wrap_glVertexAttrib4s(unsigned int index, short  x, short  y, short  z, short  w) {  (*pfn_glVertexAttrib4s)((GLuint)index, (GLshort)x, (GLshort)y, (GLshort)z, (GLshort)w); }
//
// void (APIENTRYP pfn_glVertexAttrib4sv)(GLuint index, const GLshort *v);
// GLAPI void APIENTRY wrap_glVertexAttrib4sv(unsigned int index, const short * v) {  (*pfn_glVertexAttrib4sv)((GLuint)index, (const GLshort *)v); }
//
// void (APIENTRYP pfn_glVertexAttrib4ubv)(GLuint index, const GLubyte *v);
// GLAPI void APIENTRY wrap_glVertexAttrib4ubv(unsigned int index, const unsigned char* v) {  (*pfn_glVertexAttrib4ubv)((GLuint)index, (const GLubyte *)v); }
//
// void (APIENTRYP pfn_glVertexAttrib4uiv)(GLuint index, const GLuint *v);
// GLAPI void APIENTRY wrap_glVertexAttrib4uiv(unsigned int index, const unsigned int* v) {  (*pfn_glVertexAttrib4uiv)((GLuint)index, (const GLuint *)v); }
//
// void (APIENTRYP pfn_glVertexAttrib4usv)(GLuint index, const GLushort *v);
// GLAPI void APIENTRY wrap_glVertexAttrib4usv(unsigned int index, const unsigned short* v) {  (*pfn_glVertexAttrib4usv)((GLuint)index, (const GLushort *)v); }
//
// void (APIENTRYP pfn_glVertexAttribBinding)(GLuint attribindex, GLuint bindingindex);
// GLAPI void APIENTRY wrap_glVertexAttribBinding(unsigned int attribindex, unsigned int bindingindex) {  (*pfn_glVertexAttribBinding)((GLuint)attribindex, (GLuint)bindingindex); }
//
// void (APIENTRYP pfn_glVertexAttribDivisor)(GLuint index, GLuint divisor);
// GLAPI void APIENTRY wrap_glVertexAttribDivisor(unsigned int index, unsigned int divisor) {  (*pfn_glVertexAttribDivisor)((GLuint)index, (GLuint)divisor); }
//
// void (APIENTRYP pfn_glVertexAttribFormat)(GLuint attribindex, GLint size, GLenum type, GLboolean normalized, GLuint relativeoffset);
// GLAPI void APIENTRY wrap_glVertexAttribFormat(unsigned int attribindex, int size, unsigned int t_ype, unsigned char normalized, unsigned int relativeoffset) {  (*pfn_glVertexAttribFormat)((GLuint)attribindex, (GLint)size, (GLenum)t_ype, (GLboolean)normalized, (GLuint)relativeoffset); }
//
// void (APIENTRYP pfn_glVertexAttribI1i)(GLuint index, GLint x);
// GLAPI void APIENTRY wrap_glVertexAttribI1i(unsigned int index, int x) {  (*pfn_glVertexAttribI1i)((GLuint)index, (GLint)x); }
//
// void (APIENTRYP pfn_glVertexAttribI1iv)(GLuint index, const GLint *v);
// GLAPI void APIENTRY wrap_glVertexAttribI1iv(unsigned int index, const int* v) {  (*pfn_glVertexAttribI1iv)((GLuint)index, (const GLint *)v); }
//
// void (APIENTRYP pfn_glVertexAttribI1ui)(GLuint index, GLuint x);
// GLAPI void APIENTRY wrap_glVertexAttribI1ui(unsigned int index, unsigned int x) {  (*pfn_glVertexAttribI1ui)((GLuint)index, (GLuint)x); }
//
// void (APIENTRYP pfn_glVertexAttribI1uiv)(GLuint index, const GLuint *v);
// GLAPI void APIENTRY wrap_glVertexAttribI1uiv(unsigned int index, const unsigned int* v) {  (*pfn_glVertexAttribI1uiv)((GLuint)index, (const GLuint *)v); }
//
// void (APIENTRYP pfn_glVertexAttribI2i)(GLuint index, GLint x, GLint y);
// GLAPI void APIENTRY wrap_glVertexAttribI2i(unsigned int index, int x, int y) {  (*pfn_glVertexAttribI2i)((GLuint)index, (GLint)x, (GLint)y); }
//
// void (APIENTRYP pfn_glVertexAttribI2iv)(GLuint index, const GLint *v);
// GLAPI void APIENTRY wrap_glVertexAttribI2iv(unsigned int index, const int* v) {  (*pfn_glVertexAttribI2iv)((GLuint)index, (const GLint *)v); }
//
// void (APIENTRYP pfn_glVertexAttribI2ui)(GLuint index, GLuint x, GLuint y);
// GLAPI void APIENTRY wrap_glVertexAttribI2ui(unsigned int index, unsigned int x, unsigned int y) {  (*pfn_glVertexAttribI2ui)((GLuint)index, (GLuint)x, (GLuint)y); }
//
// void (APIENTRYP pfn_glVertexAttribI2uiv)(GLuint index, const GLuint *v);
// GLAPI void APIENTRY wrap_glVertexAttribI2uiv(unsigned int index, const unsigned int* v) {  (*pfn_glVertexAttribI2uiv)((GLuint)index, (const GLuint *)v); }
//
// void (APIENTRYP pfn_glVertexAttribI3i)(GLuint index, GLint x, GLint y, GLint z);
// GLAPI void APIENTRY wrap_glVertexAttribI3i(unsigned int index, int x, int y, int z) {  (*pfn_glVertexAttribI3i)((GLuint)index, (GLint)x, (GLint)y, (GLint)z); }
//
// void (APIENTRYP pfn_glVertexAttribI3iv)(GLuint index, const GLint *v);
// GLAPI void APIENTRY wrap_glVertexAttribI3iv(unsigned int index, const int* v) {  (*pfn_glVertexAttribI3iv)((GLuint)index, (const GLint *)v); }
//
// void (APIENTRYP pfn_glVertexAttribI3ui)(GLuint index, GLuint x, GLuint y, GLuint z);
// GLAPI void APIENTRY wrap_glVertexAttribI3ui(unsigned int index, unsigned int x, unsigned int y, unsigned int z) {  (*pfn_glVertexAttribI3ui)((GLuint)index, (GLuint)x, (GLuint)y, (GLuint)z); }
//
// void (APIENTRYP pfn_glVertexAttribI3uiv)(GLuint index, const GLuint *v);
// GLAPI void APIENTRY wrap_glVertexAttribI3uiv(unsigned int index, const unsigned int* v) {  (*pfn_glVertexAttribI3uiv)((GLuint)index, (const GLuint *)v); }
//
// void (APIENTRYP pfn_glVertexAttribI4bv)(GLuint index, const GLbyte *v);
// GLAPI void APIENTRY wrap_glVertexAttribI4bv(unsigned int index, const signed char* v) {  (*pfn_glVertexAttribI4bv)((GLuint)index, (const GLbyte *)v); }
//
// void (APIENTRYP pfn_glVertexAttribI4i)(GLuint index, GLint x, GLint y, GLint z, GLint w);
// GLAPI void APIENTRY wrap_glVertexAttribI4i(unsigned int index, int x, int y, int z, int w) {  (*pfn_glVertexAttribI4i)((GLuint)index, (GLint)x, (GLint)y, (GLint)z, (GLint)w); }
//
// void (APIENTRYP pfn_glVertexAttribI4iv)(GLuint index, const GLint *v);
// GLAPI void APIENTRY wrap_glVertexAttribI4iv(unsigned int index, const int* v) {  (*pfn_glVertexAttribI4iv)((GLuint)index, (const GLint *)v); }
//
// void (APIENTRYP pfn_glVertexAttribI4sv)(GLuint index, const GLshort *v);
// GLAPI void APIENTRY wrap_glVertexAttribI4sv(unsigned int index, const short * v) {  (*pfn_glVertexAttribI4sv)((GLuint)index, (const GLshort *)v); }
//
// void (APIENTRYP pfn_glVertexAttribI4ubv)(GLuint index, const GLubyte *v);
// GLAPI void APIENTRY wrap_glVertexAttribI4ubv(unsigned int index, const unsigned char* v) {  (*pfn_glVertexAttribI4ubv)((GLuint)index, (const GLubyte *)v); }
//
// void (APIENTRYP pfn_glVertexAttribI4ui)(GLuint index, GLuint x, GLuint y, GLuint z, GLuint w);
// GLAPI void APIENTRY wrap_glVertexAttribI4ui(unsigned int index, unsigned int x, unsigned int y, unsigned int z, unsigned int w) {  (*pfn_glVertexAttribI4ui)((GLuint)index, (GLuint)x, (GLuint)y, (GLuint)z, (GLuint)w); }
//
// void (APIENTRYP pfn_glVertexAttribI4uiv)(GLuint index, const GLuint *v);
// GLAPI void APIENTRY wrap_glVertexAttribI4uiv(unsigned int index, const unsigned int* v) {  (*pfn_glVertexAttribI4uiv)((GLuint)index, (const GLuint *)v); }
//
// void (APIENTRYP pfn_glVertexAttribI4usv)(GLuint index, const GLushort *v);
// GLAPI void APIENTRY wrap_glVertexAttribI4usv(unsigned int index, const unsigned short* v) {  (*pfn_glVertexAttribI4usv)((GLuint)index, (const GLushort *)v); }
//
// void (APIENTRYP pfn_glVertexAttribIFormat)(GLuint attribindex, GLint size, GLenum type, GLuint relativeoffset);
// GLAPI void APIENTRY wrap_glVertexAttribIFormat(unsigned int attribindex, int size, unsigned int t_ype, unsigned int relativeoffset) {  (*pfn_glVertexAttribIFormat)((GLuint)attribindex, (GLint)size, (GLenum)t_ype, (GLuint)relativeoffset); }
//
// void (APIENTRYP pfn_glVertexAttribIPointer)(GLuint index, GLint size, GLenum type, GLsizei stride, const GLvoid *pointer);
// GLAPI void APIENTRY wrap_glVertexAttribIPointer(unsigned int index, int size, unsigned int t_ype, int stride, const void* pointer) {  (*pfn_glVertexAttribIPointer)((GLuint)index, (GLint)size, (GLenum)t_ype, (GLsizei)stride, (const GLvoid *)pointer); }
//
// void (APIENTRYP pfn_glVertexAttribL1d)(GLuint index, GLdouble x);
// GLAPI void APIENTRY wrap_glVertexAttribL1d(unsigned int index, double x) {  (*pfn_glVertexAttribL1d)((GLuint)index, (GLdouble)x); }
//
// void (APIENTRYP pfn_glVertexAttribL1dv)(GLuint index, const GLdouble *v);
// GLAPI void APIENTRY wrap_glVertexAttribL1dv(unsigned int index, const double* v) {  (*pfn_glVertexAttribL1dv)((GLuint)index, (const GLdouble *)v); }
//
// void (APIENTRYP pfn_glVertexAttribL2d)(GLuint index, GLdouble x, GLdouble y);
// GLAPI void APIENTRY wrap_glVertexAttribL2d(unsigned int index, double x, double y) {  (*pfn_glVertexAttribL2d)((GLuint)index, (GLdouble)x, (GLdouble)y); }
//
// void (APIENTRYP pfn_glVertexAttribL2dv)(GLuint index, const GLdouble *v);
// GLAPI void APIENTRY wrap_glVertexAttribL2dv(unsigned int index, const double* v) {  (*pfn_glVertexAttribL2dv)((GLuint)index, (const GLdouble *)v); }
//
// void (APIENTRYP pfn_glVertexAttribL3d)(GLuint index, GLdouble x, GLdouble y, GLdouble z);
// GLAPI void APIENTRY wrap_glVertexAttribL3d(unsigned int index, double x, double y, double z) {  (*pfn_glVertexAttribL3d)((GLuint)index, (GLdouble)x, (GLdouble)y, (GLdouble)z); }
//
// void (APIENTRYP pfn_glVertexAttribL3dv)(GLuint index, const GLdouble *v);
// GLAPI void APIENTRY wrap_glVertexAttribL3dv(unsigned int index, const double* v) {  (*pfn_glVertexAttribL3dv)((GLuint)index, (const GLdouble *)v); }
//
// void (APIENTRYP pfn_glVertexAttribL4d)(GLuint index, GLdouble x, GLdouble y, GLdouble z, GLdouble w);
// GLAPI void APIENTRY wrap_glVertexAttribL4d(unsigned int index, double x, double y, double z, double w) {  (*pfn_glVertexAttribL4d)((GLuint)index, (GLdouble)x, (GLdouble)y, (GLdouble)z, (GLdouble)w); }
//
// void (APIENTRYP pfn_glVertexAttribL4dv)(GLuint index, const GLdouble *v);
// GLAPI void APIENTRY wrap_glVertexAttribL4dv(unsigned int index, const double* v) {  (*pfn_glVertexAttribL4dv)((GLuint)index, (const GLdouble *)v); }
//
// void (APIENTRYP pfn_glVertexAttribLFormat)(GLuint attribindex, GLint size, GLenum type, GLuint relativeoffset);
// GLAPI void APIENTRY wrap_glVertexAttribLFormat(unsigned int attribindex, int size, unsigned int t_ype, unsigned int relativeoffset) {  (*pfn_glVertexAttribLFormat)((GLuint)attribindex, (GLint)size, (GLenum)t_ype, (GLuint)relativeoffset); }
//
// void (APIENTRYP pfn_glVertexAttribLPointer)(GLuint index, GLint size, GLenum type, GLsizei stride, const GLvoid *pointer);
// GLAPI void APIENTRY wrap_glVertexAttribLPointer(unsigned int index, int size, unsigned int t_ype, int stride, const void* pointer) {  (*pfn_glVertexAttribLPointer)((GLuint)index, (GLint)size, (GLenum)t_ype, (GLsizei)stride, (const GLvoid *)pointer); }
//
// void (APIENTRYP pfn_glVertexAttribP1ui)(GLuint index, GLenum type, GLboolean normalized, GLuint value);
// GLAPI void APIENTRY wrap_glVertexAttribP1ui(unsigned int index, unsigned int t_ype, unsigned char normalized, unsigned int value) {  (*pfn_glVertexAttribP1ui)((GLuint)index, (GLenum)t_ype, (GLboolean)normalized, (GLuint)value); }
//
// void (APIENTRYP pfn_glVertexAttribP1uiv)(GLuint index, GLenum type, GLboolean normalized, const GLuint *value);
// GLAPI void APIENTRY wrap_glVertexAttribP1uiv(unsigned int index, unsigned int t_ype, unsigned char normalized, const unsigned int* value) {  (*pfn_glVertexAttribP1uiv)((GLuint)index, (GLenum)t_ype, (GLboolean)normalized, (const GLuint *)value); }
//
// void (APIENTRYP pfn_glVertexAttribP2ui)(GLuint index, GLenum type, GLboolean normalized, GLuint value);
// GLAPI void APIENTRY wrap_glVertexAttribP2ui(unsigned int index, unsigned int t_ype, unsigned char normalized, unsigned int value) {  (*pfn_glVertexAttribP2ui)((GLuint)index, (GLenum)t_ype, (GLboolean)normalized, (GLuint)value); }
//
// void (APIENTRYP pfn_glVertexAttribP2uiv)(GLuint index, GLenum type, GLboolean normalized, const GLuint *value);
// GLAPI void APIENTRY wrap_glVertexAttribP2uiv(unsigned int index, unsigned int t_ype, unsigned char normalized, const unsigned int* value) {  (*pfn_glVertexAttribP2uiv)((GLuint)index, (GLenum)t_ype, (GLboolean)normalized, (const GLuint *)value); }
//
// void (APIENTRYP pfn_glVertexAttribP3ui)(GLuint index, GLenum type, GLboolean normalized, GLuint value);
// GLAPI void APIENTRY wrap_glVertexAttribP3ui(unsigned int index, unsigned int t_ype, unsigned char normalized, unsigned int value) {  (*pfn_glVertexAttribP3ui)((GLuint)index, (GLenum)t_ype, (GLboolean)normalized, (GLuint)value); }
//
// void (APIENTRYP pfn_glVertexAttribP3uiv)(GLuint index, GLenum type, GLboolean normalized, const GLuint *value);
// GLAPI void APIENTRY wrap_glVertexAttribP3uiv(unsigned int index, unsigned int t_ype, unsigned char normalized, const unsigned int* value) {  (*pfn_glVertexAttribP3uiv)((GLuint)index, (GLenum)t_ype, (GLboolean)normalized, (const GLuint *)value); }
//
// void (APIENTRYP pfn_glVertexAttribP4ui)(GLuint index, GLenum type, GLboolean normalized, GLuint value);
// GLAPI void APIENTRY wrap_glVertexAttribP4ui(unsigned int index, unsigned int t_ype, unsigned char normalized, unsigned int value) {  (*pfn_glVertexAttribP4ui)((GLuint)index, (GLenum)t_ype, (GLboolean)normalized, (GLuint)value); }
//
// void (APIENTRYP pfn_glVertexAttribP4uiv)(GLuint index, GLenum type, GLboolean normalized, const GLuint *value);
// GLAPI void APIENTRY wrap_glVertexAttribP4uiv(unsigned int index, unsigned int t_ype, unsigned char normalized, const unsigned int* value) {  (*pfn_glVertexAttribP4uiv)((GLuint)index, (GLenum)t_ype, (GLboolean)normalized, (const GLuint *)value); }
//
// void (APIENTRYP pfn_glVertexAttribPointer)(GLuint index, GLint size, GLenum type, GLboolean normalized, GLsizei stride, const GLvoid *pointer);
// GLAPI void APIENTRY wrap_glVertexAttribPointer(unsigned int index, int size, unsigned int t_ype, unsigned char normalized, int stride, long long pointer) {  (*pfn_glVertexAttribPointer)((GLuint)index, (GLint)size, (GLenum)t_ype, (GLboolean)normalized, (GLsizei)stride, (const GLvoid *)pointer); }
//
// void (APIENTRYP pfn_glVertexBindingDivisor)(GLuint bindingindex, GLuint divisor);
// GLAPI void APIENTRY wrap_glVertexBindingDivisor(unsigned int bindingindex, unsigned int divisor) {  (*pfn_glVertexBindingDivisor)((GLuint)bindingindex, (GLuint)divisor); }
//
// void (APIENTRYP pfn_glVertexP2ui)(GLenum type, GLuint value);
// GLAPI void APIENTRY wrap_glVertexP2ui(unsigned int t_ype, unsigned int value) {  (*pfn_glVertexP2ui)((GLenum)t_ype, (GLuint)value); }
//
// void (APIENTRYP pfn_glVertexP2uiv)(GLenum type, const GLuint *value);
// GLAPI void APIENTRY wrap_glVertexP2uiv(unsigned int t_ype, const unsigned int* value) {  (*pfn_glVertexP2uiv)((GLenum)t_ype, (const GLuint *)value); }
//
// void (APIENTRYP pfn_glVertexP3ui)(GLenum type, GLuint value);
// GLAPI void APIENTRY wrap_glVertexP3ui(unsigned int t_ype, unsigned int value) {  (*pfn_glVertexP3ui)((GLenum)t_ype, (GLuint)value); }
//
// void (APIENTRYP pfn_glVertexP3uiv)(GLenum type, const GLuint *value);
// GLAPI void APIENTRY wrap_glVertexP3uiv(unsigned int t_ype, const unsigned int* value) {  (*pfn_glVertexP3uiv)((GLenum)t_ype, (const GLuint *)value); }
//
// void (APIENTRYP pfn_glVertexP4ui)(GLenum type, GLuint value);
// GLAPI void APIENTRY wrap_glVertexP4ui(unsigned int t_ype, unsigned int value) {  (*pfn_glVertexP4ui)((GLenum)t_ype, (GLuint)value); }
//
// void (APIENTRYP pfn_glVertexP4uiv)(GLenum type, const GLuint *value);
// GLAPI void APIENTRY wrap_glVertexP4uiv(unsigned int t_ype, const unsigned int* value) {  (*pfn_glVertexP4uiv)((GLenum)t_ype, (const GLuint *)value); }
//
// void (APIENTRYP pfn_glViewport)(GLint x, GLint y, GLsizei width, GLsizei height);
// GLAPI void APIENTRY wrap_glViewport(int x, int y, int width, int height) {  (*pfn_glViewport)((GLint)x, (GLint)y, (GLsizei)width, (GLsizei)height); }
//
// void (APIENTRYP pfn_glViewportArrayv)(GLuint first, GLsizei count, const GLfloat *v);
// GLAPI void APIENTRY wrap_glViewportArrayv(unsigned int first, int count, const float* v) {  (*pfn_glViewportArrayv)((GLuint)first, (GLsizei)count, (const GLfloat *)v); }
//
// void (APIENTRYP pfn_glViewportIndexedf)(GLuint index, GLfloat x, GLfloat y, GLfloat w, GLfloat h);
// GLAPI void APIENTRY wrap_glViewportIndexedf(unsigned int index, float x, float y, float w, float h) {  (*pfn_glViewportIndexedf)((GLuint)index, (GLfloat)x, (GLfloat)y, (GLfloat)w, (GLfloat)h); }
//
// void (APIENTRYP pfn_glViewportIndexedfv)(GLuint index, const GLfloat *v);
// GLAPI void APIENTRY wrap_glViewportIndexedfv(unsigned int index, const float* v) {  (*pfn_glViewportIndexedfv)((GLuint)index, (const GLfloat *)v); }
//
// void (APIENTRYP pfn_glWaitSync)(GLsync sync, GLbitfield flags, GLuint64 timeout);
// GLAPI void APIENTRY wrap_glWaitSync(GLsync sync, unsigned int flags, unsigned long long timeout) {  (*pfn_glWaitSync)((GLsync)sync, (GLbitfield)flags, (GLuint64)timeout); }
//
//
// void init() {
//   pfn_glIsBuffer = bindMethod("glIsBuffer");
//   pfn_glIsEnabled = bindMethod("glIsEnabled");
//   pfn_glIsEnabledi = bindMethod("glIsEnabledi");
//   pfn_glIsFramebuffer = bindMethod("glIsFramebuffer");
//   pfn_glIsNamedStringARB = bindMethod("glIsNamedStringARB");
//   pfn_glIsProgram = bindMethod("glIsProgram");
//   pfn_glIsProgramPipeline = bindMethod("glIsProgramPipeline");
//   pfn_glIsQuery = bindMethod("glIsQuery");
//   pfn_glIsRenderbuffer = bindMethod("glIsRenderbuffer");
//   pfn_glIsSampler = bindMethod("glIsSampler");
//   pfn_glIsShader = bindMethod("glIsShader");
//   pfn_glIsSync = bindMethod("glIsSync");
//   pfn_glIsTexture = bindMethod("glIsTexture");
//   pfn_glIsTransformFeedback = bindMethod("glIsTransformFeedback");
//   pfn_glIsVertexArray = bindMethod("glIsVertexArray");
//   pfn_glUnmapBuffer = bindMethod("glUnmapBuffer");
//   pfn_glCheckFramebufferStatus = bindMethod("glCheckFramebufferStatus");
//   pfn_glClientWaitSync = bindMethod("glClientWaitSync");
//   pfn_glGetError = bindMethod("glGetError");
//   pfn_glGetGraphicsResetStatusARB = bindMethod("glGetGraphicsResetStatusARB");
//   pfn_glGetAttribLocation = bindMethod("glGetAttribLocation");
//   pfn_glGetFragDataIndex = bindMethod("glGetFragDataIndex");
//   pfn_glGetFragDataLocation = bindMethod("glGetFragDataLocation");
//   pfn_glGetProgramResourceLocation = bindMethod("glGetProgramResourceLocation");
//   pfn_glGetProgramResourceLocationIndex = bindMethod("glGetProgramResourceLocationIndex");
//   pfn_glGetSubroutineUniformLocation = bindMethod("glGetSubroutineUniformLocation");
//   pfn_glGetUniformLocation = bindMethod("glGetUniformLocation");
//   pfn_glCreateSyncFromCLeventARB = bindMethod("glCreateSyncFromCLeventARB");
//   pfn_glFenceSync = bindMethod("glFenceSync");
//   pfn_glCreateProgram = bindMethod("glCreateProgram");
//   pfn_glCreateShader = bindMethod("glCreateShader");
//   pfn_glCreateShaderProgramv = bindMethod("glCreateShaderProgramv");
//   pfn_glGetDebugMessageLog = bindMethod("glGetDebugMessageLog");
//   pfn_glGetDebugMessageLogARB = bindMethod("glGetDebugMessageLogARB");
//   pfn_glGetProgramResourceIndex = bindMethod("glGetProgramResourceIndex");
//   pfn_glGetSubroutineIndex = bindMethod("glGetSubroutineIndex");
//   pfn_glGetUniformBlockIndex = bindMethod("glGetUniformBlockIndex");
//   pfn_glMapBuffer = bindMethod("glMapBuffer");
//   pfn_glMapBufferRange = bindMethod("glMapBufferRange");
//   pfn_glGetString = bindMethod("glGetString");
//   pfn_glGetStringi = bindMethod("glGetStringi");
//   pfn_glActiveShaderProgram = bindMethod("glActiveShaderProgram");
//   pfn_glActiveTexture = bindMethod("glActiveTexture");
//   pfn_glAttachShader = bindMethod("glAttachShader");
//   pfn_glBeginConditionalRender = bindMethod("glBeginConditionalRender");
//   pfn_glBeginQuery = bindMethod("glBeginQuery");
//   pfn_glBeginQueryIndexed = bindMethod("glBeginQueryIndexed");
//   pfn_glBeginTransformFeedback = bindMethod("glBeginTransformFeedback");
//   pfn_glBindAttribLocation = bindMethod("glBindAttribLocation");
//   pfn_glBindBuffer = bindMethod("glBindBuffer");
//   pfn_glBindBufferBase = bindMethod("glBindBufferBase");
//   pfn_glBindBufferRange = bindMethod("glBindBufferRange");
//   pfn_glBindFragDataLocation = bindMethod("glBindFragDataLocation");
//   pfn_glBindFragDataLocationIndexed = bindMethod("glBindFragDataLocationIndexed");
//   pfn_glBindFramebuffer = bindMethod("glBindFramebuffer");
//   pfn_glBindImageTexture = bindMethod("glBindImageTexture");
//   pfn_glBindProgramPipeline = bindMethod("glBindProgramPipeline");
//   pfn_glBindRenderbuffer = bindMethod("glBindRenderbuffer");
//   pfn_glBindSampler = bindMethod("glBindSampler");
//   pfn_glBindTexture = bindMethod("glBindTexture");
//   pfn_glBindTransformFeedback = bindMethod("glBindTransformFeedback");
//   pfn_glBindVertexArray = bindMethod("glBindVertexArray");
//   pfn_glBindVertexBuffer = bindMethod("glBindVertexBuffer");
//   pfn_glBlendColor = bindMethod("glBlendColor");
//   pfn_glBlendEquation = bindMethod("glBlendEquation");
//   pfn_glBlendEquationSeparate = bindMethod("glBlendEquationSeparate");
//   pfn_glBlendEquationSeparatei = bindMethod("glBlendEquationSeparatei");
//   pfn_glBlendEquationSeparateiARB = bindMethod("glBlendEquationSeparateiARB");
//   pfn_glBlendEquationi = bindMethod("glBlendEquationi");
//   pfn_glBlendEquationiARB = bindMethod("glBlendEquationiARB");
//   pfn_glBlendFunc = bindMethod("glBlendFunc");
//   pfn_glBlendFuncSeparate = bindMethod("glBlendFuncSeparate");
//   pfn_glBlendFuncSeparatei = bindMethod("glBlendFuncSeparatei");
//   pfn_glBlendFuncSeparateiARB = bindMethod("glBlendFuncSeparateiARB");
//   pfn_glBlendFunci = bindMethod("glBlendFunci");
//   pfn_glBlendFunciARB = bindMethod("glBlendFunciARB");
//   pfn_glBlitFramebuffer = bindMethod("glBlitFramebuffer");
//   pfn_glBufferData = bindMethod("glBufferData");
//   pfn_glBufferSubData = bindMethod("glBufferSubData");
//   pfn_glClampColor = bindMethod("glClampColor");
//   pfn_glClear = bindMethod("glClear");
//   pfn_glClearBufferData = bindMethod("glClearBufferData");
//   pfn_glClearBufferSubData = bindMethod("glClearBufferSubData");
//   pfn_glClearBufferfi = bindMethod("glClearBufferfi");
//   pfn_glClearBufferfv = bindMethod("glClearBufferfv");
//   pfn_glClearBufferiv = bindMethod("glClearBufferiv");
//   pfn_glClearBufferuiv = bindMethod("glClearBufferuiv");
//   pfn_glClearColor = bindMethod("glClearColor");
//   pfn_glClearDepth = bindMethod("glClearDepth");
//   pfn_glClearDepthf = bindMethod("glClearDepthf");
//   pfn_glClearNamedBufferDataEXT = bindMethod("glClearNamedBufferDataEXT");
//   pfn_glClearNamedBufferSubDataEXT = bindMethod("glClearNamedBufferSubDataEXT");
//   pfn_glClearStencil = bindMethod("glClearStencil");
//   pfn_glColorMask = bindMethod("glColorMask");
//   pfn_glColorMaski = bindMethod("glColorMaski");
//   pfn_glColorP3ui = bindMethod("glColorP3ui");
//   pfn_glColorP3uiv = bindMethod("glColorP3uiv");
//   pfn_glColorP4ui = bindMethod("glColorP4ui");
//   pfn_glColorP4uiv = bindMethod("glColorP4uiv");
//   pfn_glCompileShader = bindMethod("glCompileShader");
//   pfn_glCompileShaderIncludeARB = bindMethod("glCompileShaderIncludeARB");
//   pfn_glCompressedTexImage1D = bindMethod("glCompressedTexImage1D");
//   pfn_glCompressedTexImage2D = bindMethod("glCompressedTexImage2D");
//   pfn_glCompressedTexImage3D = bindMethod("glCompressedTexImage3D");
//   pfn_glCompressedTexSubImage1D = bindMethod("glCompressedTexSubImage1D");
//   pfn_glCompressedTexSubImage2D = bindMethod("glCompressedTexSubImage2D");
//   pfn_glCompressedTexSubImage3D = bindMethod("glCompressedTexSubImage3D");
//   pfn_glCopyBufferSubData = bindMethod("glCopyBufferSubData");
//   pfn_glCopyImageSubData = bindMethod("glCopyImageSubData");
//   pfn_glCopyTexImage1D = bindMethod("glCopyTexImage1D");
//   pfn_glCopyTexImage2D = bindMethod("glCopyTexImage2D");
//   pfn_glCopyTexSubImage1D = bindMethod("glCopyTexSubImage1D");
//   pfn_glCopyTexSubImage2D = bindMethod("glCopyTexSubImage2D");
//   pfn_glCopyTexSubImage3D = bindMethod("glCopyTexSubImage3D");
//   pfn_glCullFace = bindMethod("glCullFace");
//   pfn_glDebugMessageCallback = bindMethod("glDebugMessageCallback");
//   pfn_glDebugMessageCallbackARB = bindMethod("glDebugMessageCallbackARB");
//   pfn_glDebugMessageControl = bindMethod("glDebugMessageControl");
//   pfn_glDebugMessageControlARB = bindMethod("glDebugMessageControlARB");
//   pfn_glDebugMessageInsert = bindMethod("glDebugMessageInsert");
//   pfn_glDebugMessageInsertARB = bindMethod("glDebugMessageInsertARB");
//   pfn_glDeleteBuffers = bindMethod("glDeleteBuffers");
//   pfn_glDeleteFramebuffers = bindMethod("glDeleteFramebuffers");
//   pfn_glDeleteNamedStringARB = bindMethod("glDeleteNamedStringARB");
//   pfn_glDeleteProgram = bindMethod("glDeleteProgram");
//   pfn_glDeleteProgramPipelines = bindMethod("glDeleteProgramPipelines");
//   pfn_glDeleteQueries = bindMethod("glDeleteQueries");
//   pfn_glDeleteRenderbuffers = bindMethod("glDeleteRenderbuffers");
//   pfn_glDeleteSamplers = bindMethod("glDeleteSamplers");
//   pfn_glDeleteShader = bindMethod("glDeleteShader");
//   pfn_glDeleteSync = bindMethod("glDeleteSync");
//   pfn_glDeleteTextures = bindMethod("glDeleteTextures");
//   pfn_glDeleteTransformFeedbacks = bindMethod("glDeleteTransformFeedbacks");
//   pfn_glDeleteVertexArrays = bindMethod("glDeleteVertexArrays");
//   pfn_glDepthFunc = bindMethod("glDepthFunc");
//   pfn_glDepthMask = bindMethod("glDepthMask");
//   pfn_glDepthRange = bindMethod("glDepthRange");
//   pfn_glDepthRangeArrayv = bindMethod("glDepthRangeArrayv");
//   pfn_glDepthRangeIndexed = bindMethod("glDepthRangeIndexed");
//   pfn_glDepthRangef = bindMethod("glDepthRangef");
//   pfn_glDetachShader = bindMethod("glDetachShader");
//   pfn_glDisable = bindMethod("glDisable");
//   pfn_glDisableVertexAttribArray = bindMethod("glDisableVertexAttribArray");
//   pfn_glDisablei = bindMethod("glDisablei");
//   pfn_glDispatchCompute = bindMethod("glDispatchCompute");
//   pfn_glDispatchComputeIndirect = bindMethod("glDispatchComputeIndirect");
//   pfn_glDrawArrays = bindMethod("glDrawArrays");
//   pfn_glDrawArraysIndirect = bindMethod("glDrawArraysIndirect");
//   pfn_glDrawArraysInstanced = bindMethod("glDrawArraysInstanced");
//   pfn_glDrawArraysInstancedBaseInstance = bindMethod("glDrawArraysInstancedBaseInstance");
//   pfn_glDrawBuffer = bindMethod("glDrawBuffer");
//   pfn_glDrawBuffers = bindMethod("glDrawBuffers");
//   pfn_glDrawElements = bindMethod("glDrawElements");
//   pfn_glDrawElementsBaseVertex = bindMethod("glDrawElementsBaseVertex");
//   pfn_glDrawElementsIndirect = bindMethod("glDrawElementsIndirect");
//   pfn_glDrawElementsInstanced = bindMethod("glDrawElementsInstanced");
//   pfn_glDrawElementsInstancedBaseInstance = bindMethod("glDrawElementsInstancedBaseInstance");
//   pfn_glDrawElementsInstancedBaseVertex = bindMethod("glDrawElementsInstancedBaseVertex");
//   pfn_glDrawElementsInstancedBaseVertexBaseInstance = bindMethod("glDrawElementsInstancedBaseVertexBaseInstance");
//   pfn_glDrawRangeElements = bindMethod("glDrawRangeElements");
//   pfn_glDrawRangeElementsBaseVertex = bindMethod("glDrawRangeElementsBaseVertex");
//   pfn_glDrawTransformFeedback = bindMethod("glDrawTransformFeedback");
//   pfn_glDrawTransformFeedbackInstanced = bindMethod("glDrawTransformFeedbackInstanced");
//   pfn_glDrawTransformFeedbackStream = bindMethod("glDrawTransformFeedbackStream");
//   pfn_glDrawTransformFeedbackStreamInstanced = bindMethod("glDrawTransformFeedbackStreamInstanced");
//   pfn_glEnable = bindMethod("glEnable");
//   pfn_glEnableVertexAttribArray = bindMethod("glEnableVertexAttribArray");
//   pfn_glEnablei = bindMethod("glEnablei");
//   pfn_glEndConditionalRender = bindMethod("glEndConditionalRender");
//   pfn_glEndQuery = bindMethod("glEndQuery");
//   pfn_glEndQueryIndexed = bindMethod("glEndQueryIndexed");
//   pfn_glEndTransformFeedback = bindMethod("glEndTransformFeedback");
//   pfn_glFinish = bindMethod("glFinish");
//   pfn_glFlush = bindMethod("glFlush");
//   pfn_glFlushMappedBufferRange = bindMethod("glFlushMappedBufferRange");
//   pfn_glFramebufferParameteri = bindMethod("glFramebufferParameteri");
//   pfn_glFramebufferRenderbuffer = bindMethod("glFramebufferRenderbuffer");
//   pfn_glFramebufferTexture = bindMethod("glFramebufferTexture");
//   pfn_glFramebufferTexture1D = bindMethod("glFramebufferTexture1D");
//   pfn_glFramebufferTexture2D = bindMethod("glFramebufferTexture2D");
//   pfn_glFramebufferTexture3D = bindMethod("glFramebufferTexture3D");
//   pfn_glFramebufferTextureLayer = bindMethod("glFramebufferTextureLayer");
//   pfn_glFrontFace = bindMethod("glFrontFace");
//   pfn_glGenBuffers = bindMethod("glGenBuffers");
//   pfn_glGenFramebuffers = bindMethod("glGenFramebuffers");
//   pfn_glGenProgramPipelines = bindMethod("glGenProgramPipelines");
//   pfn_glGenQueries = bindMethod("glGenQueries");
//   pfn_glGenRenderbuffers = bindMethod("glGenRenderbuffers");
//   pfn_glGenSamplers = bindMethod("glGenSamplers");
//   pfn_glGenTextures = bindMethod("glGenTextures");
//   pfn_glGenTransformFeedbacks = bindMethod("glGenTransformFeedbacks");
//   pfn_glGenVertexArrays = bindMethod("glGenVertexArrays");
//   pfn_glGenerateMipmap = bindMethod("glGenerateMipmap");
//   pfn_glGetActiveAtomicCounterBufferiv = bindMethod("glGetActiveAtomicCounterBufferiv");
//   pfn_glGetActiveAttrib = bindMethod("glGetActiveAttrib");
//   pfn_glGetActiveSubroutineName = bindMethod("glGetActiveSubroutineName");
//   pfn_glGetActiveSubroutineUniformName = bindMethod("glGetActiveSubroutineUniformName");
//   pfn_glGetActiveSubroutineUniformiv = bindMethod("glGetActiveSubroutineUniformiv");
//   pfn_glGetActiveUniform = bindMethod("glGetActiveUniform");
//   pfn_glGetActiveUniformBlockName = bindMethod("glGetActiveUniformBlockName");
//   pfn_glGetActiveUniformBlockiv = bindMethod("glGetActiveUniformBlockiv");
//   pfn_glGetActiveUniformName = bindMethod("glGetActiveUniformName");
//   pfn_glGetActiveUniformsiv = bindMethod("glGetActiveUniformsiv");
//   pfn_glGetAttachedShaders = bindMethod("glGetAttachedShaders");
//   pfn_glGetBooleani_v = bindMethod("glGetBooleani_v");
//   pfn_glGetBooleanv = bindMethod("glGetBooleanv");
//   pfn_glGetBufferParameteri64v = bindMethod("glGetBufferParameteri64v");
//   pfn_glGetBufferParameteriv = bindMethod("glGetBufferParameteriv");
//   pfn_glGetBufferPointerv = bindMethod("glGetBufferPointerv");
//   pfn_glGetBufferSubData = bindMethod("glGetBufferSubData");
//   pfn_glGetCompressedTexImage = bindMethod("glGetCompressedTexImage");
//   pfn_glGetDoublei_v = bindMethod("glGetDoublei_v");
//   pfn_glGetDoublev = bindMethod("glGetDoublev");
//   pfn_glGetFloati_v = bindMethod("glGetFloati_v");
//   pfn_glGetFloatv = bindMethod("glGetFloatv");
//   pfn_glGetFramebufferAttachmentParameteriv = bindMethod("glGetFramebufferAttachmentParameteriv");
//   pfn_glGetFramebufferParameteriv = bindMethod("glGetFramebufferParameteriv");
//   pfn_glGetInteger64i_v = bindMethod("glGetInteger64i_v");
//   pfn_glGetInteger64v = bindMethod("glGetInteger64v");
//   pfn_glGetIntegeri_v = bindMethod("glGetIntegeri_v");
//   pfn_glGetIntegerv = bindMethod("glGetIntegerv");
//   pfn_glGetInternalformati64v = bindMethod("glGetInternalformati64v");
//   pfn_glGetInternalformativ = bindMethod("glGetInternalformativ");
//   pfn_glGetMultisamplefv = bindMethod("glGetMultisamplefv");
//   pfn_glGetNamedFramebufferParameterivEXT = bindMethod("glGetNamedFramebufferParameterivEXT");
//   pfn_glGetNamedStringARB = bindMethod("glGetNamedStringARB");
//   pfn_glGetNamedStringivARB = bindMethod("glGetNamedStringivARB");
//   pfn_glGetObjectLabel = bindMethod("glGetObjectLabel");
//   pfn_glGetObjectPtrLabel = bindMethod("glGetObjectPtrLabel");
//   pfn_glGetPointerv = bindMethod("glGetPointerv");
//   pfn_glGetProgramBinary = bindMethod("glGetProgramBinary");
//   pfn_glGetProgramInfoLog = bindMethod("glGetProgramInfoLog");
//   pfn_glGetProgramInterfaceiv = bindMethod("glGetProgramInterfaceiv");
//   pfn_glGetProgramPipelineInfoLog = bindMethod("glGetProgramPipelineInfoLog");
//   pfn_glGetProgramPipelineiv = bindMethod("glGetProgramPipelineiv");
//   pfn_glGetProgramResourceName = bindMethod("glGetProgramResourceName");
//   pfn_glGetProgramResourceiv = bindMethod("glGetProgramResourceiv");
//   pfn_glGetProgramStageiv = bindMethod("glGetProgramStageiv");
//   pfn_glGetProgramiv = bindMethod("glGetProgramiv");
//   pfn_glGetQueryIndexediv = bindMethod("glGetQueryIndexediv");
//   pfn_glGetQueryObjecti64v = bindMethod("glGetQueryObjecti64v");
//   pfn_glGetQueryObjectiv = bindMethod("glGetQueryObjectiv");
//   pfn_glGetQueryObjectui64v = bindMethod("glGetQueryObjectui64v");
//   pfn_glGetQueryObjectuiv = bindMethod("glGetQueryObjectuiv");
//   pfn_glGetQueryiv = bindMethod("glGetQueryiv");
//   pfn_glGetRenderbufferParameteriv = bindMethod("glGetRenderbufferParameteriv");
//   pfn_glGetSamplerParameterIiv = bindMethod("glGetSamplerParameterIiv");
//   pfn_glGetSamplerParameterIuiv = bindMethod("glGetSamplerParameterIuiv");
//   pfn_glGetSamplerParameterfv = bindMethod("glGetSamplerParameterfv");
//   pfn_glGetSamplerParameteriv = bindMethod("glGetSamplerParameteriv");
//   pfn_glGetShaderInfoLog = bindMethod("glGetShaderInfoLog");
//   pfn_glGetShaderPrecisionFormat = bindMethod("glGetShaderPrecisionFormat");
//   pfn_glGetShaderSource = bindMethod("glGetShaderSource");
//   pfn_glGetShaderiv = bindMethod("glGetShaderiv");
//   pfn_glGetSynciv = bindMethod("glGetSynciv");
//   pfn_glGetTexImage = bindMethod("glGetTexImage");
//   pfn_glGetTexLevelParameterfv = bindMethod("glGetTexLevelParameterfv");
//   pfn_glGetTexLevelParameteriv = bindMethod("glGetTexLevelParameteriv");
//   pfn_glGetTexParameterIiv = bindMethod("glGetTexParameterIiv");
//   pfn_glGetTexParameterIuiv = bindMethod("glGetTexParameterIuiv");
//   pfn_glGetTexParameterfv = bindMethod("glGetTexParameterfv");
//   pfn_glGetTexParameteriv = bindMethod("glGetTexParameteriv");
//   pfn_glGetTransformFeedbackVarying = bindMethod("glGetTransformFeedbackVarying");
//   pfn_glGetUniformIndices = bindMethod("glGetUniformIndices");
//   pfn_glGetUniformSubroutineuiv = bindMethod("glGetUniformSubroutineuiv");
//   pfn_glGetUniformdv = bindMethod("glGetUniformdv");
//   pfn_glGetUniformfv = bindMethod("glGetUniformfv");
//   pfn_glGetUniformiv = bindMethod("glGetUniformiv");
//   pfn_glGetUniformuiv = bindMethod("glGetUniformuiv");
//   pfn_glGetVertexAttribIiv = bindMethod("glGetVertexAttribIiv");
//   pfn_glGetVertexAttribIuiv = bindMethod("glGetVertexAttribIuiv");
//   pfn_glGetVertexAttribLdv = bindMethod("glGetVertexAttribLdv");
//   pfn_glGetVertexAttribPointerv = bindMethod("glGetVertexAttribPointerv");
//   pfn_glGetVertexAttribdv = bindMethod("glGetVertexAttribdv");
//   pfn_glGetVertexAttribfv = bindMethod("glGetVertexAttribfv");
//   pfn_glGetVertexAttribiv = bindMethod("glGetVertexAttribiv");
//   pfn_glGetnColorTableARB = bindMethod("glGetnColorTableARB");
//   pfn_glGetnCompressedTexImageARB = bindMethod("glGetnCompressedTexImageARB");
//   pfn_glGetnConvolutionFilterARB = bindMethod("glGetnConvolutionFilterARB");
//   pfn_glGetnHistogramARB = bindMethod("glGetnHistogramARB");
//   pfn_glGetnMapdvARB = bindMethod("glGetnMapdvARB");
//   pfn_glGetnMapfvARB = bindMethod("glGetnMapfvARB");
//   pfn_glGetnMapivARB = bindMethod("glGetnMapivARB");
//   pfn_glGetnMinmaxARB = bindMethod("glGetnMinmaxARB");
//   pfn_glGetnPixelMapfvARB = bindMethod("glGetnPixelMapfvARB");
//   pfn_glGetnPixelMapuivARB = bindMethod("glGetnPixelMapuivARB");
//   pfn_glGetnPixelMapusvARB = bindMethod("glGetnPixelMapusvARB");
//   pfn_glGetnPolygonStippleARB = bindMethod("glGetnPolygonStippleARB");
//   pfn_glGetnSeparableFilterARB = bindMethod("glGetnSeparableFilterARB");
//   pfn_glGetnTexImageARB = bindMethod("glGetnTexImageARB");
//   pfn_glGetnUniformdvARB = bindMethod("glGetnUniformdvARB");
//   pfn_glGetnUniformfvARB = bindMethod("glGetnUniformfvARB");
//   pfn_glGetnUniformivARB = bindMethod("glGetnUniformivARB");
//   pfn_glGetnUniformuivARB = bindMethod("glGetnUniformuivARB");
//   pfn_glHint = bindMethod("glHint");
//   pfn_glInvalidateBufferData = bindMethod("glInvalidateBufferData");
//   pfn_glInvalidateBufferSubData = bindMethod("glInvalidateBufferSubData");
//   pfn_glInvalidateFramebuffer = bindMethod("glInvalidateFramebuffer");
//   pfn_glInvalidateSubFramebuffer = bindMethod("glInvalidateSubFramebuffer");
//   pfn_glInvalidateTexImage = bindMethod("glInvalidateTexImage");
//   pfn_glInvalidateTexSubImage = bindMethod("glInvalidateTexSubImage");
//   pfn_glLineWidth = bindMethod("glLineWidth");
//   pfn_glLinkProgram = bindMethod("glLinkProgram");
//   pfn_glLogicOp = bindMethod("glLogicOp");
//   pfn_glMemoryBarrier = bindMethod("glMemoryBarrier");
//   pfn_glMinSampleShading = bindMethod("glMinSampleShading");
//   pfn_glMinSampleShadingARB = bindMethod("glMinSampleShadingARB");
//   pfn_glMultiDrawArrays = bindMethod("glMultiDrawArrays");
//   pfn_glMultiDrawArraysIndirect = bindMethod("glMultiDrawArraysIndirect");
//   pfn_glMultiDrawElements = bindMethod("glMultiDrawElements");
//   pfn_glMultiDrawElementsBaseVertex = bindMethod("glMultiDrawElementsBaseVertex");
//   pfn_glMultiDrawElementsIndirect = bindMethod("glMultiDrawElementsIndirect");
//   pfn_glMultiTexCoordP1ui = bindMethod("glMultiTexCoordP1ui");
//   pfn_glMultiTexCoordP1uiv = bindMethod("glMultiTexCoordP1uiv");
//   pfn_glMultiTexCoordP2ui = bindMethod("glMultiTexCoordP2ui");
//   pfn_glMultiTexCoordP2uiv = bindMethod("glMultiTexCoordP2uiv");
//   pfn_glMultiTexCoordP3ui = bindMethod("glMultiTexCoordP3ui");
//   pfn_glMultiTexCoordP3uiv = bindMethod("glMultiTexCoordP3uiv");
//   pfn_glMultiTexCoordP4ui = bindMethod("glMultiTexCoordP4ui");
//   pfn_glMultiTexCoordP4uiv = bindMethod("glMultiTexCoordP4uiv");
//   pfn_glNamedFramebufferParameteriEXT = bindMethod("glNamedFramebufferParameteriEXT");
//   pfn_glNamedStringARB = bindMethod("glNamedStringARB");
//   pfn_glNormalP3ui = bindMethod("glNormalP3ui");
//   pfn_glNormalP3uiv = bindMethod("glNormalP3uiv");
//   pfn_glObjectLabel = bindMethod("glObjectLabel");
//   pfn_glObjectPtrLabel = bindMethod("glObjectPtrLabel");
//   pfn_glPatchParameterfv = bindMethod("glPatchParameterfv");
//   pfn_glPatchParameteri = bindMethod("glPatchParameteri");
//   pfn_glPauseTransformFeedback = bindMethod("glPauseTransformFeedback");
//   pfn_glPixelStoref = bindMethod("glPixelStoref");
//   pfn_glPixelStorei = bindMethod("glPixelStorei");
//   pfn_glPointParameterf = bindMethod("glPointParameterf");
//   pfn_glPointParameterfv = bindMethod("glPointParameterfv");
//   pfn_glPointParameteri = bindMethod("glPointParameteri");
//   pfn_glPointParameteriv = bindMethod("glPointParameteriv");
//   pfn_glPointSize = bindMethod("glPointSize");
//   pfn_glPolygonMode = bindMethod("glPolygonMode");
//   pfn_glPolygonOffset = bindMethod("glPolygonOffset");
//   pfn_glPopDebugGroup = bindMethod("glPopDebugGroup");
//   pfn_glPrimitiveRestartIndex = bindMethod("glPrimitiveRestartIndex");
//   pfn_glProgramBinary = bindMethod("glProgramBinary");
//   pfn_glProgramParameteri = bindMethod("glProgramParameteri");
//   pfn_glProgramUniform1d = bindMethod("glProgramUniform1d");
//   pfn_glProgramUniform1dv = bindMethod("glProgramUniform1dv");
//   pfn_glProgramUniform1f = bindMethod("glProgramUniform1f");
//   pfn_glProgramUniform1fv = bindMethod("glProgramUniform1fv");
//   pfn_glProgramUniform1i = bindMethod("glProgramUniform1i");
//   pfn_glProgramUniform1iv = bindMethod("glProgramUniform1iv");
//   pfn_glProgramUniform1ui = bindMethod("glProgramUniform1ui");
//   pfn_glProgramUniform1uiv = bindMethod("glProgramUniform1uiv");
//   pfn_glProgramUniform2d = bindMethod("glProgramUniform2d");
//   pfn_glProgramUniform2dv = bindMethod("glProgramUniform2dv");
//   pfn_glProgramUniform2f = bindMethod("glProgramUniform2f");
//   pfn_glProgramUniform2fv = bindMethod("glProgramUniform2fv");
//   pfn_glProgramUniform2i = bindMethod("glProgramUniform2i");
//   pfn_glProgramUniform2iv = bindMethod("glProgramUniform2iv");
//   pfn_glProgramUniform2ui = bindMethod("glProgramUniform2ui");
//   pfn_glProgramUniform2uiv = bindMethod("glProgramUniform2uiv");
//   pfn_glProgramUniform3d = bindMethod("glProgramUniform3d");
//   pfn_glProgramUniform3dv = bindMethod("glProgramUniform3dv");
//   pfn_glProgramUniform3f = bindMethod("glProgramUniform3f");
//   pfn_glProgramUniform3fv = bindMethod("glProgramUniform3fv");
//   pfn_glProgramUniform3i = bindMethod("glProgramUniform3i");
//   pfn_glProgramUniform3iv = bindMethod("glProgramUniform3iv");
//   pfn_glProgramUniform3ui = bindMethod("glProgramUniform3ui");
//   pfn_glProgramUniform3uiv = bindMethod("glProgramUniform3uiv");
//   pfn_glProgramUniform4d = bindMethod("glProgramUniform4d");
//   pfn_glProgramUniform4dv = bindMethod("glProgramUniform4dv");
//   pfn_glProgramUniform4f = bindMethod("glProgramUniform4f");
//   pfn_glProgramUniform4fv = bindMethod("glProgramUniform4fv");
//   pfn_glProgramUniform4i = bindMethod("glProgramUniform4i");
//   pfn_glProgramUniform4iv = bindMethod("glProgramUniform4iv");
//   pfn_glProgramUniform4ui = bindMethod("glProgramUniform4ui");
//   pfn_glProgramUniform4uiv = bindMethod("glProgramUniform4uiv");
//   pfn_glProgramUniformMatrix2dv = bindMethod("glProgramUniformMatrix2dv");
//   pfn_glProgramUniformMatrix2fv = bindMethod("glProgramUniformMatrix2fv");
//   pfn_glProgramUniformMatrix2x3dv = bindMethod("glProgramUniformMatrix2x3dv");
//   pfn_glProgramUniformMatrix2x3fv = bindMethod("glProgramUniformMatrix2x3fv");
//   pfn_glProgramUniformMatrix2x4dv = bindMethod("glProgramUniformMatrix2x4dv");
//   pfn_glProgramUniformMatrix2x4fv = bindMethod("glProgramUniformMatrix2x4fv");
//   pfn_glProgramUniformMatrix3dv = bindMethod("glProgramUniformMatrix3dv");
//   pfn_glProgramUniformMatrix3fv = bindMethod("glProgramUniformMatrix3fv");
//   pfn_glProgramUniformMatrix3x2dv = bindMethod("glProgramUniformMatrix3x2dv");
//   pfn_glProgramUniformMatrix3x2fv = bindMethod("glProgramUniformMatrix3x2fv");
//   pfn_glProgramUniformMatrix3x4dv = bindMethod("glProgramUniformMatrix3x4dv");
//   pfn_glProgramUniformMatrix3x4fv = bindMethod("glProgramUniformMatrix3x4fv");
//   pfn_glProgramUniformMatrix4dv = bindMethod("glProgramUniformMatrix4dv");
//   pfn_glProgramUniformMatrix4fv = bindMethod("glProgramUniformMatrix4fv");
//   pfn_glProgramUniformMatrix4x2dv = bindMethod("glProgramUniformMatrix4x2dv");
//   pfn_glProgramUniformMatrix4x2fv = bindMethod("glProgramUniformMatrix4x2fv");
//   pfn_glProgramUniformMatrix4x3dv = bindMethod("glProgramUniformMatrix4x3dv");
//   pfn_glProgramUniformMatrix4x3fv = bindMethod("glProgramUniformMatrix4x3fv");
//   pfn_glProvokingVertex = bindMethod("glProvokingVertex");
//   pfn_glPushDebugGroup = bindMethod("glPushDebugGroup");
//   pfn_glQueryCounter = bindMethod("glQueryCounter");
//   pfn_glReadBuffer = bindMethod("glReadBuffer");
//   pfn_glReadPixels = bindMethod("glReadPixels");
//   pfn_glReadnPixelsARB = bindMethod("glReadnPixelsARB");
//   pfn_glReleaseShaderCompiler = bindMethod("glReleaseShaderCompiler");
//   pfn_glRenderbufferStorage = bindMethod("glRenderbufferStorage");
//   pfn_glRenderbufferStorageMultisample = bindMethod("glRenderbufferStorageMultisample");
//   pfn_glResumeTransformFeedback = bindMethod("glResumeTransformFeedback");
//   pfn_glSampleCoverage = bindMethod("glSampleCoverage");
//   pfn_glSampleMaski = bindMethod("glSampleMaski");
//   pfn_glSamplerParameterIiv = bindMethod("glSamplerParameterIiv");
//   pfn_glSamplerParameterIuiv = bindMethod("glSamplerParameterIuiv");
//   pfn_glSamplerParameterf = bindMethod("glSamplerParameterf");
//   pfn_glSamplerParameterfv = bindMethod("glSamplerParameterfv");
//   pfn_glSamplerParameteri = bindMethod("glSamplerParameteri");
//   pfn_glSamplerParameteriv = bindMethod("glSamplerParameteriv");
//   pfn_glScissor = bindMethod("glScissor");
//   pfn_glScissorArrayv = bindMethod("glScissorArrayv");
//   pfn_glScissorIndexed = bindMethod("glScissorIndexed");
//   pfn_glScissorIndexedv = bindMethod("glScissorIndexedv");
//   pfn_glSecondaryColorP3ui = bindMethod("glSecondaryColorP3ui");
//   pfn_glSecondaryColorP3uiv = bindMethod("glSecondaryColorP3uiv");
//   pfn_glShaderBinary = bindMethod("glShaderBinary");
//   pfn_glShaderSource = bindMethod("glShaderSource");
//   pfn_glShaderStorageBlockBinding = bindMethod("glShaderStorageBlockBinding");
//   pfn_glStencilFunc = bindMethod("glStencilFunc");
//   pfn_glStencilFuncSeparate = bindMethod("glStencilFuncSeparate");
//   pfn_glStencilMask = bindMethod("glStencilMask");
//   pfn_glStencilMaskSeparate = bindMethod("glStencilMaskSeparate");
//   pfn_glStencilOp = bindMethod("glStencilOp");
//   pfn_glStencilOpSeparate = bindMethod("glStencilOpSeparate");
//   pfn_glTexBuffer = bindMethod("glTexBuffer");
//   pfn_glTexBufferRange = bindMethod("glTexBufferRange");
//   pfn_glTexCoordP1ui = bindMethod("glTexCoordP1ui");
//   pfn_glTexCoordP1uiv = bindMethod("glTexCoordP1uiv");
//   pfn_glTexCoordP2ui = bindMethod("glTexCoordP2ui");
//   pfn_glTexCoordP2uiv = bindMethod("glTexCoordP2uiv");
//   pfn_glTexCoordP3ui = bindMethod("glTexCoordP3ui");
//   pfn_glTexCoordP3uiv = bindMethod("glTexCoordP3uiv");
//   pfn_glTexCoordP4ui = bindMethod("glTexCoordP4ui");
//   pfn_glTexCoordP4uiv = bindMethod("glTexCoordP4uiv");
//   pfn_glTexImage1D = bindMethod("glTexImage1D");
//   pfn_glTexImage2D = bindMethod("glTexImage2D");
//   pfn_glTexImage2DMultisample = bindMethod("glTexImage2DMultisample");
//   pfn_glTexImage3D = bindMethod("glTexImage3D");
//   pfn_glTexImage3DMultisample = bindMethod("glTexImage3DMultisample");
//   pfn_glTexParameterIiv = bindMethod("glTexParameterIiv");
//   pfn_glTexParameterIuiv = bindMethod("glTexParameterIuiv");
//   pfn_glTexParameterf = bindMethod("glTexParameterf");
//   pfn_glTexParameterfv = bindMethod("glTexParameterfv");
//   pfn_glTexParameteri = bindMethod("glTexParameteri");
//   pfn_glTexParameteriv = bindMethod("glTexParameteriv");
//   pfn_glTexStorage1D = bindMethod("glTexStorage1D");
//   pfn_glTexStorage2D = bindMethod("glTexStorage2D");
//   pfn_glTexStorage2DMultisample = bindMethod("glTexStorage2DMultisample");
//   pfn_glTexStorage3D = bindMethod("glTexStorage3D");
//   pfn_glTexStorage3DMultisample = bindMethod("glTexStorage3DMultisample");
//   pfn_glTexSubImage1D = bindMethod("glTexSubImage1D");
//   pfn_glTexSubImage2D = bindMethod("glTexSubImage2D");
//   pfn_glTexSubImage3D = bindMethod("glTexSubImage3D");
//   pfn_glTextureBufferRangeEXT = bindMethod("glTextureBufferRangeEXT");
//   pfn_glTextureStorage1DEXT = bindMethod("glTextureStorage1DEXT");
//   pfn_glTextureStorage2DEXT = bindMethod("glTextureStorage2DEXT");
//   pfn_glTextureStorage2DMultisampleEXT = bindMethod("glTextureStorage2DMultisampleEXT");
//   pfn_glTextureStorage3DEXT = bindMethod("glTextureStorage3DEXT");
//   pfn_glTextureStorage3DMultisampleEXT = bindMethod("glTextureStorage3DMultisampleEXT");
//   pfn_glTextureView = bindMethod("glTextureView");
//   pfn_glTransformFeedbackVaryings = bindMethod("glTransformFeedbackVaryings");
//   pfn_glUniform1d = bindMethod("glUniform1d");
//   pfn_glUniform1dv = bindMethod("glUniform1dv");
//   pfn_glUniform1f = bindMethod("glUniform1f");
//   pfn_glUniform1fv = bindMethod("glUniform1fv");
//   pfn_glUniform1i = bindMethod("glUniform1i");
//   pfn_glUniform1iv = bindMethod("glUniform1iv");
//   pfn_glUniform1ui = bindMethod("glUniform1ui");
//   pfn_glUniform1uiv = bindMethod("glUniform1uiv");
//   pfn_glUniform2d = bindMethod("glUniform2d");
//   pfn_glUniform2dv = bindMethod("glUniform2dv");
//   pfn_glUniform2f = bindMethod("glUniform2f");
//   pfn_glUniform2fv = bindMethod("glUniform2fv");
//   pfn_glUniform2i = bindMethod("glUniform2i");
//   pfn_glUniform2iv = bindMethod("glUniform2iv");
//   pfn_glUniform2ui = bindMethod("glUniform2ui");
//   pfn_glUniform2uiv = bindMethod("glUniform2uiv");
//   pfn_glUniform3d = bindMethod("glUniform3d");
//   pfn_glUniform3dv = bindMethod("glUniform3dv");
//   pfn_glUniform3f = bindMethod("glUniform3f");
//   pfn_glUniform3fv = bindMethod("glUniform3fv");
//   pfn_glUniform3i = bindMethod("glUniform3i");
//   pfn_glUniform3iv = bindMethod("glUniform3iv");
//   pfn_glUniform3ui = bindMethod("glUniform3ui");
//   pfn_glUniform3uiv = bindMethod("glUniform3uiv");
//   pfn_glUniform4d = bindMethod("glUniform4d");
//   pfn_glUniform4dv = bindMethod("glUniform4dv");
//   pfn_glUniform4f = bindMethod("glUniform4f");
//   pfn_glUniform4fv = bindMethod("glUniform4fv");
//   pfn_glUniform4i = bindMethod("glUniform4i");
//   pfn_glUniform4iv = bindMethod("glUniform4iv");
//   pfn_glUniform4ui = bindMethod("glUniform4ui");
//   pfn_glUniform4uiv = bindMethod("glUniform4uiv");
//   pfn_glUniformBlockBinding = bindMethod("glUniformBlockBinding");
//   pfn_glUniformMatrix2dv = bindMethod("glUniformMatrix2dv");
//   pfn_glUniformMatrix2fv = bindMethod("glUniformMatrix2fv");
//   pfn_glUniformMatrix2x3dv = bindMethod("glUniformMatrix2x3dv");
//   pfn_glUniformMatrix2x3fv = bindMethod("glUniformMatrix2x3fv");
//   pfn_glUniformMatrix2x4dv = bindMethod("glUniformMatrix2x4dv");
//   pfn_glUniformMatrix2x4fv = bindMethod("glUniformMatrix2x4fv");
//   pfn_glUniformMatrix3dv = bindMethod("glUniformMatrix3dv");
//   pfn_glUniformMatrix3fv = bindMethod("glUniformMatrix3fv");
//   pfn_glUniformMatrix3x2dv = bindMethod("glUniformMatrix3x2dv");
//   pfn_glUniformMatrix3x2fv = bindMethod("glUniformMatrix3x2fv");
//   pfn_glUniformMatrix3x4dv = bindMethod("glUniformMatrix3x4dv");
//   pfn_glUniformMatrix3x4fv = bindMethod("glUniformMatrix3x4fv");
//   pfn_glUniformMatrix4dv = bindMethod("glUniformMatrix4dv");
//   pfn_glUniformMatrix4fv = bindMethod("glUniformMatrix4fv");
//   pfn_glUniformMatrix4x2dv = bindMethod("glUniformMatrix4x2dv");
//   pfn_glUniformMatrix4x2fv = bindMethod("glUniformMatrix4x2fv");
//   pfn_glUniformMatrix4x3dv = bindMethod("glUniformMatrix4x3dv");
//   pfn_glUniformMatrix4x3fv = bindMethod("glUniformMatrix4x3fv");
//   pfn_glUniformSubroutinesuiv = bindMethod("glUniformSubroutinesuiv");
//   pfn_glUseProgram = bindMethod("glUseProgram");
//   pfn_glUseProgramStages = bindMethod("glUseProgramStages");
//   pfn_glValidateProgram = bindMethod("glValidateProgram");
//   pfn_glValidateProgramPipeline = bindMethod("glValidateProgramPipeline");
//   pfn_glVertexArrayBindVertexBufferEXT = bindMethod("glVertexArrayBindVertexBufferEXT");
//   pfn_glVertexArrayVertexAttribBindingEXT = bindMethod("glVertexArrayVertexAttribBindingEXT");
//   pfn_glVertexArrayVertexAttribFormatEXT = bindMethod("glVertexArrayVertexAttribFormatEXT");
//   pfn_glVertexArrayVertexAttribIFormatEXT = bindMethod("glVertexArrayVertexAttribIFormatEXT");
//   pfn_glVertexArrayVertexAttribLFormatEXT = bindMethod("glVertexArrayVertexAttribLFormatEXT");
//   pfn_glVertexArrayVertexBindingDivisorEXT = bindMethod("glVertexArrayVertexBindingDivisorEXT");
//   pfn_glVertexAttrib1d = bindMethod("glVertexAttrib1d");
//   pfn_glVertexAttrib1dv = bindMethod("glVertexAttrib1dv");
//   pfn_glVertexAttrib1f = bindMethod("glVertexAttrib1f");
//   pfn_glVertexAttrib1fv = bindMethod("glVertexAttrib1fv");
//   pfn_glVertexAttrib1s = bindMethod("glVertexAttrib1s");
//   pfn_glVertexAttrib1sv = bindMethod("glVertexAttrib1sv");
//   pfn_glVertexAttrib2d = bindMethod("glVertexAttrib2d");
//   pfn_glVertexAttrib2dv = bindMethod("glVertexAttrib2dv");
//   pfn_glVertexAttrib2f = bindMethod("glVertexAttrib2f");
//   pfn_glVertexAttrib2fv = bindMethod("glVertexAttrib2fv");
//   pfn_glVertexAttrib2s = bindMethod("glVertexAttrib2s");
//   pfn_glVertexAttrib2sv = bindMethod("glVertexAttrib2sv");
//   pfn_glVertexAttrib3d = bindMethod("glVertexAttrib3d");
//   pfn_glVertexAttrib3dv = bindMethod("glVertexAttrib3dv");
//   pfn_glVertexAttrib3f = bindMethod("glVertexAttrib3f");
//   pfn_glVertexAttrib3fv = bindMethod("glVertexAttrib3fv");
//   pfn_glVertexAttrib3s = bindMethod("glVertexAttrib3s");
//   pfn_glVertexAttrib3sv = bindMethod("glVertexAttrib3sv");
//   pfn_glVertexAttrib4Nbv = bindMethod("glVertexAttrib4Nbv");
//   pfn_glVertexAttrib4Niv = bindMethod("glVertexAttrib4Niv");
//   pfn_glVertexAttrib4Nsv = bindMethod("glVertexAttrib4Nsv");
//   pfn_glVertexAttrib4Nub = bindMethod("glVertexAttrib4Nub");
//   pfn_glVertexAttrib4Nubv = bindMethod("glVertexAttrib4Nubv");
//   pfn_glVertexAttrib4Nuiv = bindMethod("glVertexAttrib4Nuiv");
//   pfn_glVertexAttrib4Nusv = bindMethod("glVertexAttrib4Nusv");
//   pfn_glVertexAttrib4bv = bindMethod("glVertexAttrib4bv");
//   pfn_glVertexAttrib4d = bindMethod("glVertexAttrib4d");
//   pfn_glVertexAttrib4dv = bindMethod("glVertexAttrib4dv");
//   pfn_glVertexAttrib4f = bindMethod("glVertexAttrib4f");
//   pfn_glVertexAttrib4fv = bindMethod("glVertexAttrib4fv");
//   pfn_glVertexAttrib4iv = bindMethod("glVertexAttrib4iv");
//   pfn_glVertexAttrib4s = bindMethod("glVertexAttrib4s");
//   pfn_glVertexAttrib4sv = bindMethod("glVertexAttrib4sv");
//   pfn_glVertexAttrib4ubv = bindMethod("glVertexAttrib4ubv");
//   pfn_glVertexAttrib4uiv = bindMethod("glVertexAttrib4uiv");
//   pfn_glVertexAttrib4usv = bindMethod("glVertexAttrib4usv");
//   pfn_glVertexAttribBinding = bindMethod("glVertexAttribBinding");
//   pfn_glVertexAttribDivisor = bindMethod("glVertexAttribDivisor");
//   pfn_glVertexAttribFormat = bindMethod("glVertexAttribFormat");
//   pfn_glVertexAttribI1i = bindMethod("glVertexAttribI1i");
//   pfn_glVertexAttribI1iv = bindMethod("glVertexAttribI1iv");
//   pfn_glVertexAttribI1ui = bindMethod("glVertexAttribI1ui");
//   pfn_glVertexAttribI1uiv = bindMethod("glVertexAttribI1uiv");
//   pfn_glVertexAttribI2i = bindMethod("glVertexAttribI2i");
//   pfn_glVertexAttribI2iv = bindMethod("glVertexAttribI2iv");
//   pfn_glVertexAttribI2ui = bindMethod("glVertexAttribI2ui");
//   pfn_glVertexAttribI2uiv = bindMethod("glVertexAttribI2uiv");
//   pfn_glVertexAttribI3i = bindMethod("glVertexAttribI3i");
//   pfn_glVertexAttribI3iv = bindMethod("glVertexAttribI3iv");
//   pfn_glVertexAttribI3ui = bindMethod("glVertexAttribI3ui");
//   pfn_glVertexAttribI3uiv = bindMethod("glVertexAttribI3uiv");
//   pfn_glVertexAttribI4bv = bindMethod("glVertexAttribI4bv");
//   pfn_glVertexAttribI4i = bindMethod("glVertexAttribI4i");
//   pfn_glVertexAttribI4iv = bindMethod("glVertexAttribI4iv");
//   pfn_glVertexAttribI4sv = bindMethod("glVertexAttribI4sv");
//   pfn_glVertexAttribI4ubv = bindMethod("glVertexAttribI4ubv");
//   pfn_glVertexAttribI4ui = bindMethod("glVertexAttribI4ui");
//   pfn_glVertexAttribI4uiv = bindMethod("glVertexAttribI4uiv");
//   pfn_glVertexAttribI4usv = bindMethod("glVertexAttribI4usv");
//   pfn_glVertexAttribIFormat = bindMethod("glVertexAttribIFormat");
//   pfn_glVertexAttribIPointer = bindMethod("glVertexAttribIPointer");
//   pfn_glVertexAttribL1d = bindMethod("glVertexAttribL1d");
//   pfn_glVertexAttribL1dv = bindMethod("glVertexAttribL1dv");
//   pfn_glVertexAttribL2d = bindMethod("glVertexAttribL2d");
//   pfn_glVertexAttribL2dv = bindMethod("glVertexAttribL2dv");
//   pfn_glVertexAttribL3d = bindMethod("glVertexAttribL3d");
//   pfn_glVertexAttribL3dv = bindMethod("glVertexAttribL3dv");
//   pfn_glVertexAttribL4d = bindMethod("glVertexAttribL4d");
//   pfn_glVertexAttribL4dv = bindMethod("glVertexAttribL4dv");
//   pfn_glVertexAttribLFormat = bindMethod("glVertexAttribLFormat");
//   pfn_glVertexAttribLPointer = bindMethod("glVertexAttribLPointer");
//   pfn_glVertexAttribP1ui = bindMethod("glVertexAttribP1ui");
//   pfn_glVertexAttribP1uiv = bindMethod("glVertexAttribP1uiv");
//   pfn_glVertexAttribP2ui = bindMethod("glVertexAttribP2ui");
//   pfn_glVertexAttribP2uiv = bindMethod("glVertexAttribP2uiv");
//   pfn_glVertexAttribP3ui = bindMethod("glVertexAttribP3ui");
//   pfn_glVertexAttribP3uiv = bindMethod("glVertexAttribP3uiv");
//   pfn_glVertexAttribP4ui = bindMethod("glVertexAttribP4ui");
//   pfn_glVertexAttribP4uiv = bindMethod("glVertexAttribP4uiv");
//   pfn_glVertexAttribPointer = bindMethod("glVertexAttribPointer");
//   pfn_glVertexBindingDivisor = bindMethod("glVertexBindingDivisor");
//   pfn_glVertexP2ui = bindMethod("glVertexP2ui");
//   pfn_glVertexP2uiv = bindMethod("glVertexP2uiv");
//   pfn_glVertexP3ui = bindMethod("glVertexP3ui");
//   pfn_glVertexP3uiv = bindMethod("glVertexP3uiv");
//   pfn_glVertexP4ui = bindMethod("glVertexP4ui");
//   pfn_glVertexP4uiv = bindMethod("glVertexP4uiv");
//   pfn_glViewport = bindMethod("glViewport");
//   pfn_glViewportArrayv = bindMethod("glViewportArrayv");
//   pfn_glViewportIndexedf = bindMethod("glViewportIndexedf");
//   pfn_glViewportIndexedfv = bindMethod("glViewportIndexedfv");
//   pfn_glWaitSync = bindMethod("glWaitSync");
// }
//
import "C"
import "unsafe"

import "fmt"

// convert a GL uint boolean to a go bool
func cbool(glbool uint) bool {
	return glbool == TRUE
}

// Special type mappings
type (
	Pointer      unsafe.Pointer
	Sync         C.GLsync
	clContext    C.struct_Cl_context
	clEvent      C.struct_Cl_event
	DEBUGPROCARB C.GLDEBUGPROCARB
	DEBUGPROC    C.GLDEBUGPROC
)

// bind the methods to the function pointers
func Init() {
	C.init()
}

const (
	DEPTH_BUFFER_BIT                                           = 0x00000100
	STENCIL_BUFFER_BIT                                         = 0x00000400
	COLOR_BUFFER_BIT                                           = 0x00004000
	FALSE                                                      = 0
	TRUE                                                       = 1
	POINTS                                                     = 0x0000
	LINES                                                      = 0x0001
	LINE_LOOP                                                  = 0x0002
	LINE_STRIP                                                 = 0x0003
	TRIANGLES                                                  = 0x0004
	TRIANGLE_STRIP                                             = 0x0005
	TRIANGLE_FAN                                               = 0x0006
	NEVER                                                      = 0x0200
	LESS                                                       = 0x0201
	EQUAL                                                      = 0x0202
	LEQUAL                                                     = 0x0203
	GREATER                                                    = 0x0204
	NOTEQUAL                                                   = 0x0205
	GEQUAL                                                     = 0x0206
	ALWAYS                                                     = 0x0207
	ZERO                                                       = 0
	ONE                                                        = 1
	SRC_COLOR                                                  = 0x0300
	ONE_MINUS_SRC_COLOR                                        = 0x0301
	SRC_ALPHA                                                  = 0x0302
	ONE_MINUS_SRC_ALPHA                                        = 0x0303
	DST_ALPHA                                                  = 0x0304
	ONE_MINUS_DST_ALPHA                                        = 0x0305
	DST_COLOR                                                  = 0x0306
	ONE_MINUS_DST_COLOR                                        = 0x0307
	SRC_ALPHA_SATURATE                                         = 0x0308
	NONE                                                       = 0
	FRONT_LEFT                                                 = 0x0400
	FRONT_RIGHT                                                = 0x0401
	BACK_LEFT                                                  = 0x0402
	BACK_RIGHT                                                 = 0x0403
	FRONT                                                      = 0x0404
	BACK                                                       = 0x0405
	LEFT                                                       = 0x0406
	RIGHT                                                      = 0x0407
	FRONT_AND_BACK                                             = 0x0408
	NO_ERROR                                                   = 0
	INVALID_ENUM                                               = 0x0500
	INVALID_VALUE                                              = 0x0501
	INVALID_OPERATION                                          = 0x0502
	OUT_OF_MEMORY                                              = 0x0505
	CW                                                         = 0x0900
	CCW                                                        = 0x0901
	POINT_SIZE                                                 = 0x0B11
	POINT_SIZE_RANGE                                           = 0x0B12
	POINT_SIZE_GRANULARITY                                     = 0x0B13
	LINE_SMOOTH                                                = 0x0B20
	LINE_WIDTH                                                 = 0x0B21
	LINE_WIDTH_RANGE                                           = 0x0B22
	LINE_WIDTH_GRANULARITY                                     = 0x0B23
	POLYGON_SMOOTH                                             = 0x0B41
	CULL_FACE                                                  = 0x0B44
	CULL_FACE_MODE                                             = 0x0B45
	FRONT_FACE                                                 = 0x0B46
	DEPTH_RANGE                                                = 0x0B70
	DEPTH_TEST                                                 = 0x0B71
	DEPTH_WRITEMASK                                            = 0x0B72
	DEPTH_CLEAR_VALUE                                          = 0x0B73
	DEPTH_FUNC                                                 = 0x0B74
	STENCIL_TEST                                               = 0x0B90
	STENCIL_CLEAR_VALUE                                        = 0x0B91
	STENCIL_FUNC                                               = 0x0B92
	STENCIL_VALUE_MASK                                         = 0x0B93
	STENCIL_FAIL                                               = 0x0B94
	STENCIL_PASS_DEPTH_FAIL                                    = 0x0B95
	STENCIL_PASS_DEPTH_PASS                                    = 0x0B96
	STENCIL_REF                                                = 0x0B97
	STENCIL_WRITEMASK                                          = 0x0B98
	VIEWPORT                                                   = 0x0BA2
	DITHER                                                     = 0x0BD0
	BLEND_DST                                                  = 0x0BE0
	BLEND_SRC                                                  = 0x0BE1
	BLEND                                                      = 0x0BE2
	LOGIC_OP_MODE                                              = 0x0BF0
	COLOR_LOGIC_OP                                             = 0x0BF2
	DRAW_BUFFER                                                = 0x0C01
	READ_BUFFER                                                = 0x0C02
	SCISSOR_BOX                                                = 0x0C10
	SCISSOR_TEST                                               = 0x0C11
	COLOR_CLEAR_VALUE                                          = 0x0C22
	COLOR_WRITEMASK                                            = 0x0C23
	DOUBLEBUFFER                                               = 0x0C32
	STEREO                                                     = 0x0C33
	LINE_SMOOTH_HINT                                           = 0x0C52
	POLYGON_SMOOTH_HINT                                        = 0x0C53
	UNPACK_SWAP_BYTES                                          = 0x0CF0
	UNPACK_LSB_FIRST                                           = 0x0CF1
	UNPACK_ROW_LENGTH                                          = 0x0CF2
	UNPACK_SKIP_ROWS                                           = 0x0CF3
	UNPACK_SKIP_PIXELS                                         = 0x0CF4
	UNPACK_ALIGNMENT                                           = 0x0CF5
	PACK_SWAP_BYTES                                            = 0x0D00
	PACK_LSB_FIRST                                             = 0x0D01
	PACK_ROW_LENGTH                                            = 0x0D02
	PACK_SKIP_ROWS                                             = 0x0D03
	PACK_SKIP_PIXELS                                           = 0x0D04
	PACK_ALIGNMENT                                             = 0x0D05
	MAX_TEXTURE_SIZE                                           = 0x0D33
	MAX_VIEWPORT_DIMS                                          = 0x0D3A
	SUBPIXEL_BITS                                              = 0x0D50
	TEXTURE_1D                                                 = 0x0DE0
	TEXTURE_2D                                                 = 0x0DE1
	POLYGON_OFFSET_UNITS                                       = 0x2A00
	POLYGON_OFFSET_POINT                                       = 0x2A01
	POLYGON_OFFSET_LINE                                        = 0x2A02
	POLYGON_OFFSET_FILL                                        = 0x8037
	POLYGON_OFFSET_FACTOR                                      = 0x8038
	TEXTURE_BINDING_1D                                         = 0x8068
	TEXTURE_BINDING_2D                                         = 0x8069
	TEXTURE_WIDTH                                              = 0x1000
	TEXTURE_HEIGHT                                             = 0x1001
	TEXTURE_INTERNAL_FORMAT                                    = 0x1003
	TEXTURE_BORDER_COLOR                                       = 0x1004
	TEXTURE_RED_SIZE                                           = 0x805C
	TEXTURE_GREEN_SIZE                                         = 0x805D
	TEXTURE_BLUE_SIZE                                          = 0x805E
	TEXTURE_ALPHA_SIZE                                         = 0x805F
	DONT_CARE                                                  = 0x1100
	FASTEST                                                    = 0x1101
	NICEST                                                     = 0x1102
	BYTE                                                       = 0x1400
	UNSIGNED_BYTE                                              = 0x1401
	SHORT                                                      = 0x1402
	UNSIGNED_SHORT                                             = 0x1403
	INT                                                        = 0x1404
	UNSIGNED_INT                                               = 0x1405
	FLOAT                                                      = 0x1406
	DOUBLE                                                     = 0x140A
	STACK_OVERFLOW                                             = 0x0503
	STACK_UNDERFLOW                                            = 0x0504
	CLEAR                                                      = 0x1500
	AND                                                        = 0x1501
	AND_REVERSE                                                = 0x1502
	COPY                                                       = 0x1503
	AND_INVERTED                                               = 0x1504
	NOOP                                                       = 0x1505
	XOR                                                        = 0x1506
	OR                                                         = 0x1507
	NOR                                                        = 0x1508
	EQUIV                                                      = 0x1509
	INVERT                                                     = 0x150A
	OR_REVERSE                                                 = 0x150B
	COPY_INVERTED                                              = 0x150C
	OR_INVERTED                                                = 0x150D
	NAND                                                       = 0x150E
	SET                                                        = 0x150F
	TEXTURE                                                    = 0x1702
	COLOR                                                      = 0x1800
	DEPTH                                                      = 0x1801
	STENCIL                                                    = 0x1802
	STENCIL_INDEX                                              = 0x1901
	DEPTH_COMPONENT                                            = 0x1902
	RED                                                        = 0x1903
	GREEN                                                      = 0x1904
	BLUE                                                       = 0x1905
	ALPHA                                                      = 0x1906
	RGB                                                        = 0x1907
	RGBA                                                       = 0x1908
	POINT                                                      = 0x1B00
	LINE                                                       = 0x1B01
	FILL                                                       = 0x1B02
	KEEP                                                       = 0x1E00
	REPLACE                                                    = 0x1E01
	INCR                                                       = 0x1E02
	DECR                                                       = 0x1E03
	VENDOR                                                     = 0x1F00
	RENDERER                                                   = 0x1F01
	VERSION                                                    = 0x1F02
	EXTENSIONS                                                 = 0x1F03
	NEAREST                                                    = 0x2600
	LINEAR                                                     = 0x2601
	NEAREST_MIPMAP_NEAREST                                     = 0x2700
	LINEAR_MIPMAP_NEAREST                                      = 0x2701
	NEAREST_MIPMAP_LINEAR                                      = 0x2702
	LINEAR_MIPMAP_LINEAR                                       = 0x2703
	TEXTURE_MAG_FILTER                                         = 0x2800
	TEXTURE_MIN_FILTER                                         = 0x2801
	TEXTURE_WRAP_S                                             = 0x2802
	TEXTURE_WRAP_T                                             = 0x2803
	PROXY_TEXTURE_1D                                           = 0x8063
	PROXY_TEXTURE_2D                                           = 0x8064
	REPEAT                                                     = 0x2901
	R3_G3_B2                                                   = 0x2A10
	RGB4                                                       = 0x804F
	RGB5                                                       = 0x8050
	RGB8                                                       = 0x8051
	RGB10                                                      = 0x8052
	RGB12                                                      = 0x8053
	RGB16                                                      = 0x8054
	RGBA2                                                      = 0x8055
	RGBA4                                                      = 0x8056
	RGB5_A1                                                    = 0x8057
	RGBA8                                                      = 0x8058
	RGB10_A2                                                   = 0x8059
	RGBA12                                                     = 0x805A
	RGBA16                                                     = 0x805B
	UNSIGNED_BYTE_3_3_2                                        = 0x8032
	UNSIGNED_SHORT_4_4_4_4                                     = 0x8033
	UNSIGNED_SHORT_5_5_5_1                                     = 0x8034
	UNSIGNED_INT_8_8_8_8                                       = 0x8035
	UNSIGNED_INT_10_10_10_2                                    = 0x8036
	TEXTURE_BINDING_3D                                         = 0x806A
	PACK_SKIP_IMAGES                                           = 0x806B
	PACK_IMAGE_HEIGHT                                          = 0x806C
	UNPACK_SKIP_IMAGES                                         = 0x806D
	UNPACK_IMAGE_HEIGHT                                        = 0x806E
	TEXTURE_3D                                                 = 0x806F
	PROXY_TEXTURE_3D                                           = 0x8070
	TEXTURE_DEPTH                                              = 0x8071
	TEXTURE_WRAP_R                                             = 0x8072
	MAX_3D_TEXTURE_SIZE                                        = 0x8073
	UNSIGNED_BYTE_2_3_3_REV                                    = 0x8362
	UNSIGNED_SHORT_5_6_5                                       = 0x8363
	UNSIGNED_SHORT_5_6_5_REV                                   = 0x8364
	UNSIGNED_SHORT_4_4_4_4_REV                                 = 0x8365
	UNSIGNED_SHORT_1_5_5_5_REV                                 = 0x8366
	UNSIGNED_INT_8_8_8_8_REV                                   = 0x8367
	UNSIGNED_INT_2_10_10_10_REV                                = 0x8368
	BGR                                                        = 0x80E0
	BGRA                                                       = 0x80E1
	MAX_ELEMENTS_VERTICES                                      = 0x80E8
	MAX_ELEMENTS_INDICES                                       = 0x80E9
	CLAMP_TO_EDGE                                              = 0x812F
	TEXTURE_MIN_LOD                                            = 0x813A
	TEXTURE_MAX_LOD                                            = 0x813B
	TEXTURE_BASE_LEVEL                                         = 0x813C
	TEXTURE_MAX_LEVEL                                          = 0x813D
	SMOOTH_POINT_SIZE_RANGE                                    = 0x0B12
	SMOOTH_POINT_SIZE_GRANULARITY                              = 0x0B13
	SMOOTH_LINE_WIDTH_RANGE                                    = 0x0B22
	SMOOTH_LINE_WIDTH_GRANULARITY                              = 0x0B23
	ALIASED_LINE_WIDTH_RANGE                                   = 0x846E
	CONSTANT_COLOR                                             = 0x8001
	ONE_MINUS_CONSTANT_COLOR                                   = 0x8002
	CONSTANT_ALPHA                                             = 0x8003
	ONE_MINUS_CONSTANT_ALPHA                                   = 0x8004
	BLEND_COLOR                                                = 0x8005
	FUNC_ADD                                                   = 0x8006
	MIN                                                        = 0x8007
	MAX                                                        = 0x8008
	BLEND_EQUATION                                             = 0x8009
	FUNC_SUBTRACT                                              = 0x800A
	FUNC_REVERSE_SUBTRACT                                      = 0x800B
	TEXTURE0                                                   = 0x84C0
	TEXTURE1                                                   = 0x84C1
	TEXTURE2                                                   = 0x84C2
	TEXTURE3                                                   = 0x84C3
	TEXTURE4                                                   = 0x84C4
	TEXTURE5                                                   = 0x84C5
	TEXTURE6                                                   = 0x84C6
	TEXTURE7                                                   = 0x84C7
	TEXTURE8                                                   = 0x84C8
	TEXTURE9                                                   = 0x84C9
	TEXTURE10                                                  = 0x84CA
	TEXTURE11                                                  = 0x84CB
	TEXTURE12                                                  = 0x84CC
	TEXTURE13                                                  = 0x84CD
	TEXTURE14                                                  = 0x84CE
	TEXTURE15                                                  = 0x84CF
	TEXTURE16                                                  = 0x84D0
	TEXTURE17                                                  = 0x84D1
	TEXTURE18                                                  = 0x84D2
	TEXTURE19                                                  = 0x84D3
	TEXTURE20                                                  = 0x84D4
	TEXTURE21                                                  = 0x84D5
	TEXTURE22                                                  = 0x84D6
	TEXTURE23                                                  = 0x84D7
	TEXTURE24                                                  = 0x84D8
	TEXTURE25                                                  = 0x84D9
	TEXTURE26                                                  = 0x84DA
	TEXTURE27                                                  = 0x84DB
	TEXTURE28                                                  = 0x84DC
	TEXTURE29                                                  = 0x84DD
	TEXTURE30                                                  = 0x84DE
	TEXTURE31                                                  = 0x84DF
	ACTIVE_TEXTURE                                             = 0x84E0
	MULTISAMPLE                                                = 0x809D
	SAMPLE_ALPHA_TO_COVERAGE                                   = 0x809E
	SAMPLE_ALPHA_TO_ONE                                        = 0x809F
	SAMPLE_COVERAGE                                            = 0x80A0
	SAMPLE_BUFFERS                                             = 0x80A8
	SAMPLES                                                    = 0x80A9
	SAMPLE_COVERAGE_VALUE                                      = 0x80AA
	SAMPLE_COVERAGE_INVERT                                     = 0x80AB
	TEXTURE_CUBE_MAP                                           = 0x8513
	TEXTURE_BINDING_CUBE_MAP                                   = 0x8514
	TEXTURE_CUBE_MAP_POSITIVE_X                                = 0x8515
	TEXTURE_CUBE_MAP_NEGATIVE_X                                = 0x8516
	TEXTURE_CUBE_MAP_POSITIVE_Y                                = 0x8517
	TEXTURE_CUBE_MAP_NEGATIVE_Y                                = 0x8518
	TEXTURE_CUBE_MAP_POSITIVE_Z                                = 0x8519
	TEXTURE_CUBE_MAP_NEGATIVE_Z                                = 0x851A
	PROXY_TEXTURE_CUBE_MAP                                     = 0x851B
	MAX_CUBE_MAP_TEXTURE_SIZE                                  = 0x851C
	COMPRESSED_RGB                                             = 0x84ED
	COMPRESSED_RGBA                                            = 0x84EE
	TEXTURE_COMPRESSION_HINT                                   = 0x84EF
	TEXTURE_COMPRESSED_IMAGE_SIZE                              = 0x86A0
	TEXTURE_COMPRESSED                                         = 0x86A1
	NUM_COMPRESSED_TEXTURE_FORMATS                             = 0x86A2
	COMPRESSED_TEXTURE_FORMATS                                 = 0x86A3
	CLAMP_TO_BORDER                                            = 0x812D
	BLEND_DST_RGB                                              = 0x80C8
	BLEND_SRC_RGB                                              = 0x80C9
	BLEND_DST_ALPHA                                            = 0x80CA
	BLEND_SRC_ALPHA                                            = 0x80CB
	POINT_FADE_THRESHOLD_SIZE                                  = 0x8128
	DEPTH_COMPONENT16                                          = 0x81A5
	DEPTH_COMPONENT24                                          = 0x81A6
	DEPTH_COMPONENT32                                          = 0x81A7
	MIRRORED_REPEAT                                            = 0x8370
	MAX_TEXTURE_LOD_BIAS                                       = 0x84FD
	TEXTURE_LOD_BIAS                                           = 0x8501
	INCR_WRAP                                                  = 0x8507
	DECR_WRAP                                                  = 0x8508
	TEXTURE_DEPTH_SIZE                                         = 0x884A
	TEXTURE_COMPARE_MODE                                       = 0x884C
	TEXTURE_COMPARE_FUNC                                       = 0x884D
	BUFFER_SIZE                                                = 0x8764
	BUFFER_USAGE                                               = 0x8765
	QUERY_COUNTER_BITS                                         = 0x8864
	CURRENT_QUERY                                              = 0x8865
	QUERY_RESULT                                               = 0x8866
	QUERY_RESULT_AVAILABLE                                     = 0x8867
	ARRAY_BUFFER                                               = 0x8892
	ELEMENT_ARRAY_BUFFER                                       = 0x8893
	ARRAY_BUFFER_BINDING                                       = 0x8894
	ELEMENT_ARRAY_BUFFER_BINDING                               = 0x8895
	VERTEX_ATTRIB_ARRAY_BUFFER_BINDING                         = 0x889F
	READ_ONLY                                                  = 0x88B8
	WRITE_ONLY                                                 = 0x88B9
	READ_WRITE                                                 = 0x88BA
	BUFFER_ACCESS                                              = 0x88BB
	BUFFER_MAPPED                                              = 0x88BC
	BUFFER_MAP_POINTER                                         = 0x88BD
	STREAM_DRAW                                                = 0x88E0
	STREAM_READ                                                = 0x88E1
	STREAM_COPY                                                = 0x88E2
	STATIC_DRAW                                                = 0x88E4
	STATIC_READ                                                = 0x88E5
	STATIC_COPY                                                = 0x88E6
	DYNAMIC_DRAW                                               = 0x88E8
	DYNAMIC_READ                                               = 0x88E9
	DYNAMIC_COPY                                               = 0x88EA
	SAMPLES_PASSED                                             = 0x8914
	BLEND_EQUATION_RGB                                         = 0x8009
	VERTEX_ATTRIB_ARRAY_ENABLED                                = 0x8622
	VERTEX_ATTRIB_ARRAY_SIZE                                   = 0x8623
	VERTEX_ATTRIB_ARRAY_STRIDE                                 = 0x8624
	VERTEX_ATTRIB_ARRAY_TYPE                                   = 0x8625
	CURRENT_VERTEX_ATTRIB                                      = 0x8626
	VERTEX_PROGRAM_POINT_SIZE                                  = 0x8642
	VERTEX_ATTRIB_ARRAY_POINTER                                = 0x8645
	STENCIL_BACK_FUNC                                          = 0x8800
	STENCIL_BACK_FAIL                                          = 0x8801
	STENCIL_BACK_PASS_DEPTH_FAIL                               = 0x8802
	STENCIL_BACK_PASS_DEPTH_PASS                               = 0x8803
	MAX_DRAW_BUFFERS                                           = 0x8824
	DRAW_BUFFER0                                               = 0x8825
	DRAW_BUFFER1                                               = 0x8826
	DRAW_BUFFER2                                               = 0x8827
	DRAW_BUFFER3                                               = 0x8828
	DRAW_BUFFER4                                               = 0x8829
	DRAW_BUFFER5                                               = 0x882A
	DRAW_BUFFER6                                               = 0x882B
	DRAW_BUFFER7                                               = 0x882C
	DRAW_BUFFER8                                               = 0x882D
	DRAW_BUFFER9                                               = 0x882E
	DRAW_BUFFER10                                              = 0x882F
	DRAW_BUFFER11                                              = 0x8830
	DRAW_BUFFER12                                              = 0x8831
	DRAW_BUFFER13                                              = 0x8832
	DRAW_BUFFER14                                              = 0x8833
	DRAW_BUFFER15                                              = 0x8834
	BLEND_EQUATION_ALPHA                                       = 0x883D
	MAX_VERTEX_ATTRIBS                                         = 0x8869
	VERTEX_ATTRIB_ARRAY_NORMALIZED                             = 0x886A
	MAX_TEXTURE_IMAGE_UNITS                                    = 0x8872
	FRAGMENT_SHADER                                            = 0x8B30
	VERTEX_SHADER                                              = 0x8B31
	MAX_FRAGMENT_UNIFORM_COMPONENTS                            = 0x8B49
	MAX_VERTEX_UNIFORM_COMPONENTS                              = 0x8B4A
	MAX_VARYING_FLOATS                                         = 0x8B4B
	MAX_VERTEX_TEXTURE_IMAGE_UNITS                             = 0x8B4C
	MAX_COMBINED_TEXTURE_IMAGE_UNITS                           = 0x8B4D
	SHADER_TYPE                                                = 0x8B4F
	FLOAT_VEC2                                                 = 0x8B50
	FLOAT_VEC3                                                 = 0x8B51
	FLOAT_VEC4                                                 = 0x8B52
	INT_VEC2                                                   = 0x8B53
	INT_VEC3                                                   = 0x8B54
	INT_VEC4                                                   = 0x8B55
	BOOL                                                       = 0x8B56
	BOOL_VEC2                                                  = 0x8B57
	BOOL_VEC3                                                  = 0x8B58
	BOOL_VEC4                                                  = 0x8B59
	FLOAT_MAT2                                                 = 0x8B5A
	FLOAT_MAT3                                                 = 0x8B5B
	FLOAT_MAT4                                                 = 0x8B5C
	SAMPLER_1D                                                 = 0x8B5D
	SAMPLER_2D                                                 = 0x8B5E
	SAMPLER_3D                                                 = 0x8B5F
	SAMPLER_CUBE                                               = 0x8B60
	SAMPLER_1D_SHADOW                                          = 0x8B61
	SAMPLER_2D_SHADOW                                          = 0x8B62
	DELETE_STATUS                                              = 0x8B80
	COMPILE_STATUS                                             = 0x8B81
	LINK_STATUS                                                = 0x8B82
	VALIDATE_STATUS                                            = 0x8B83
	INFO_LOG_LENGTH                                            = 0x8B84
	ATTACHED_SHADERS                                           = 0x8B85
	ACTIVE_UNIFORMS                                            = 0x8B86
	ACTIVE_UNIFORM_MAX_LENGTH                                  = 0x8B87
	SHADER_SOURCE_LENGTH                                       = 0x8B88
	ACTIVE_ATTRIBUTES                                          = 0x8B89
	ACTIVE_ATTRIBUTE_MAX_LENGTH                                = 0x8B8A
	FRAGMENT_SHADER_DERIVATIVE_HINT                            = 0x8B8B
	SHADING_LANGUAGE_VERSION                                   = 0x8B8C
	CURRENT_PROGRAM                                            = 0x8B8D
	POINT_SPRITE_COORD_ORIGIN                                  = 0x8CA0
	LOWER_LEFT                                                 = 0x8CA1
	UPPER_LEFT                                                 = 0x8CA2
	STENCIL_BACK_REF                                           = 0x8CA3
	STENCIL_BACK_VALUE_MASK                                    = 0x8CA4
	STENCIL_BACK_WRITEMASK                                     = 0x8CA5
	PIXEL_PACK_BUFFER                                          = 0x88EB
	PIXEL_UNPACK_BUFFER                                        = 0x88EC
	PIXEL_PACK_BUFFER_BINDING                                  = 0x88ED
	PIXEL_UNPACK_BUFFER_BINDING                                = 0x88EF
	FLOAT_MAT2x3                                               = 0x8B65
	FLOAT_MAT2x4                                               = 0x8B66
	FLOAT_MAT3x2                                               = 0x8B67
	FLOAT_MAT3x4                                               = 0x8B68
	FLOAT_MAT4x2                                               = 0x8B69
	FLOAT_MAT4x3                                               = 0x8B6A
	SRGB                                                       = 0x8C40
	SRGB8                                                      = 0x8C41
	SRGB_ALPHA                                                 = 0x8C42
	SRGB8_ALPHA8                                               = 0x8C43
	COMPRESSED_SRGB                                            = 0x8C48
	COMPRESSED_SRGB_ALPHA                                      = 0x8C49
	COMPARE_REF_TO_TEXTURE                                     = 0x884E
	CLIP_DISTANCE0                                             = 0x3000
	CLIP_DISTANCE1                                             = 0x3001
	CLIP_DISTANCE2                                             = 0x3002
	CLIP_DISTANCE3                                             = 0x3003
	CLIP_DISTANCE4                                             = 0x3004
	CLIP_DISTANCE5                                             = 0x3005
	CLIP_DISTANCE6                                             = 0x3006
	CLIP_DISTANCE7                                             = 0x3007
	MAX_CLIP_DISTANCES                                         = 0x0D32
	MAJOR_VERSION                                              = 0x821B
	MINOR_VERSION                                              = 0x821C
	NUM_EXTENSIONS                                             = 0x821D
	CONTEXT_FLAGS                                              = 0x821E
	COMPRESSED_RED                                             = 0x8225
	COMPRESSED_RG                                              = 0x8226
	CONTEXT_FLAG_FORWARD_COMPATIBLE_BIT                        = 0x0001
	RGBA32F                                                    = 0x8814
	RGB32F                                                     = 0x8815
	RGBA16F                                                    = 0x881A
	RGB16F                                                     = 0x881B
	VERTEX_ATTRIB_ARRAY_INTEGER                                = 0x88FD
	MAX_ARRAY_TEXTURE_LAYERS                                   = 0x88FF
	MIN_PROGRAM_TEXEL_OFFSET                                   = 0x8904
	MAX_PROGRAM_TEXEL_OFFSET                                   = 0x8905
	CLAMP_READ_COLOR                                           = 0x891C
	FIXED_ONLY                                                 = 0x891D
	MAX_VARYING_COMPONENTS                                     = 0x8B4B
	TEXTURE_1D_ARRAY                                           = 0x8C18
	PROXY_TEXTURE_1D_ARRAY                                     = 0x8C19
	TEXTURE_2D_ARRAY                                           = 0x8C1A
	PROXY_TEXTURE_2D_ARRAY                                     = 0x8C1B
	TEXTURE_BINDING_1D_ARRAY                                   = 0x8C1C
	TEXTURE_BINDING_2D_ARRAY                                   = 0x8C1D
	R11F_G11F_B10F                                             = 0x8C3A
	UNSIGNED_INT_10F_11F_11F_REV                               = 0x8C3B
	RGB9_E5                                                    = 0x8C3D
	UNSIGNED_INT_5_9_9_9_REV                                   = 0x8C3E
	TEXTURE_SHARED_SIZE                                        = 0x8C3F
	TRANSFORM_FEEDBACK_VARYING_MAX_LENGTH                      = 0x8C76
	TRANSFORM_FEEDBACK_BUFFER_MODE                             = 0x8C7F
	MAX_TRANSFORM_FEEDBACK_SEPARATE_COMPONENTS                 = 0x8C80
	TRANSFORM_FEEDBACK_VARYINGS                                = 0x8C83
	TRANSFORM_FEEDBACK_BUFFER_START                            = 0x8C84
	TRANSFORM_FEEDBACK_BUFFER_SIZE                             = 0x8C85
	PRIMITIVES_GENERATED                                       = 0x8C87
	TRANSFORM_FEEDBACK_PRIMITIVES_WRITTEN                      = 0x8C88
	RASTERIZER_DISCARD                                         = 0x8C89
	MAX_TRANSFORM_FEEDBACK_INTERLEAVED_COMPONENTS              = 0x8C8A
	MAX_TRANSFORM_FEEDBACK_SEPARATE_ATTRIBS                    = 0x8C8B
	INTERLEAVED_ATTRIBS                                        = 0x8C8C
	SEPARATE_ATTRIBS                                           = 0x8C8D
	TRANSFORM_FEEDBACK_BUFFER                                  = 0x8C8E
	TRANSFORM_FEEDBACK_BUFFER_BINDING                          = 0x8C8F
	RGBA32UI                                                   = 0x8D70
	RGB32UI                                                    = 0x8D71
	RGBA16UI                                                   = 0x8D76
	RGB16UI                                                    = 0x8D77
	RGBA8UI                                                    = 0x8D7C
	RGB8UI                                                     = 0x8D7D
	RGBA32I                                                    = 0x8D82
	RGB32I                                                     = 0x8D83
	RGBA16I                                                    = 0x8D88
	RGB16I                                                     = 0x8D89
	RGBA8I                                                     = 0x8D8E
	RGB8I                                                      = 0x8D8F
	RED_INTEGER                                                = 0x8D94
	GREEN_INTEGER                                              = 0x8D95
	BLUE_INTEGER                                               = 0x8D96
	RGB_INTEGER                                                = 0x8D98
	RGBA_INTEGER                                               = 0x8D99
	BGR_INTEGER                                                = 0x8D9A
	BGRA_INTEGER                                               = 0x8D9B
	SAMPLER_1D_ARRAY                                           = 0x8DC0
	SAMPLER_2D_ARRAY                                           = 0x8DC1
	SAMPLER_1D_ARRAY_SHADOW                                    = 0x8DC3
	SAMPLER_2D_ARRAY_SHADOW                                    = 0x8DC4
	SAMPLER_CUBE_SHADOW                                        = 0x8DC5
	UNSIGNED_INT_VEC2                                          = 0x8DC6
	UNSIGNED_INT_VEC3                                          = 0x8DC7
	UNSIGNED_INT_VEC4                                          = 0x8DC8
	INT_SAMPLER_1D                                             = 0x8DC9
	INT_SAMPLER_2D                                             = 0x8DCA
	INT_SAMPLER_3D                                             = 0x8DCB
	INT_SAMPLER_CUBE                                           = 0x8DCC
	INT_SAMPLER_1D_ARRAY                                       = 0x8DCE
	INT_SAMPLER_2D_ARRAY                                       = 0x8DCF
	UNSIGNED_INT_SAMPLER_1D                                    = 0x8DD1
	UNSIGNED_INT_SAMPLER_2D                                    = 0x8DD2
	UNSIGNED_INT_SAMPLER_3D                                    = 0x8DD3
	UNSIGNED_INT_SAMPLER_CUBE                                  = 0x8DD4
	UNSIGNED_INT_SAMPLER_1D_ARRAY                              = 0x8DD6
	UNSIGNED_INT_SAMPLER_2D_ARRAY                              = 0x8DD7
	QUERY_WAIT                                                 = 0x8E13
	QUERY_NO_WAIT                                              = 0x8E14
	QUERY_BY_REGION_WAIT                                       = 0x8E15
	QUERY_BY_REGION_NO_WAIT                                    = 0x8E16
	BUFFER_ACCESS_FLAGS                                        = 0x911F
	BUFFER_MAP_LENGTH                                          = 0x9120
	BUFFER_MAP_OFFSET                                          = 0x9121
	SAMPLER_2D_RECT                                            = 0x8B63
	SAMPLER_2D_RECT_SHADOW                                     = 0x8B64
	SAMPLER_BUFFER                                             = 0x8DC2
	INT_SAMPLER_2D_RECT                                        = 0x8DCD
	INT_SAMPLER_BUFFER                                         = 0x8DD0
	UNSIGNED_INT_SAMPLER_2D_RECT                               = 0x8DD5
	UNSIGNED_INT_SAMPLER_BUFFER                                = 0x8DD8
	TEXTURE_BUFFER                                             = 0x8C2A
	MAX_TEXTURE_BUFFER_SIZE                                    = 0x8C2B
	TEXTURE_BINDING_BUFFER                                     = 0x8C2C
	TEXTURE_BUFFER_DATA_STORE_BINDING                          = 0x8C2D
	TEXTURE_BUFFER_FORMAT                                      = 0x8C2E
	TEXTURE_RECTANGLE                                          = 0x84F5
	TEXTURE_BINDING_RECTANGLE                                  = 0x84F6
	PROXY_TEXTURE_RECTANGLE                                    = 0x84F7
	MAX_RECTANGLE_TEXTURE_SIZE                                 = 0x84F8
	RED_SNORM                                                  = 0x8F90
	RG_SNORM                                                   = 0x8F91
	RGB_SNORM                                                  = 0x8F92
	RGBA_SNORM                                                 = 0x8F93
	R8_SNORM                                                   = 0x8F94
	RG8_SNORM                                                  = 0x8F95
	RGB8_SNORM                                                 = 0x8F96
	RGBA8_SNORM                                                = 0x8F97
	R16_SNORM                                                  = 0x8F98
	RG16_SNORM                                                 = 0x8F99
	RGB16_SNORM                                                = 0x8F9A
	RGBA16_SNORM                                               = 0x8F9B
	SIGNED_NORMALIZED                                          = 0x8F9C
	PRIMITIVE_RESTART                                          = 0x8F9D
	PRIMITIVE_RESTART_INDEX                                    = 0x8F9E
	CONTEXT_CORE_PROFILE_BIT                                   = 0x00000001
	CONTEXT_COMPATIBILITY_PROFILE_BIT                          = 0x00000002
	LINES_ADJACENCY                                            = 0x000A
	LINE_STRIP_ADJACENCY                                       = 0x000B
	TRIANGLES_ADJACENCY                                        = 0x000C
	TRIANGLE_STRIP_ADJACENCY                                   = 0x000D
	PROGRAM_POINT_SIZE                                         = 0x8642
	MAX_GEOMETRY_TEXTURE_IMAGE_UNITS                           = 0x8C29
	FRAMEBUFFER_ATTACHMENT_LAYERED                             = 0x8DA7
	FRAMEBUFFER_INCOMPLETE_LAYER_TARGETS                       = 0x8DA8
	GEOMETRY_SHADER                                            = 0x8DD9
	GEOMETRY_VERTICES_OUT                                      = 0x8916
	GEOMETRY_INPUT_TYPE                                        = 0x8917
	GEOMETRY_OUTPUT_TYPE                                       = 0x8918
	MAX_GEOMETRY_UNIFORM_COMPONENTS                            = 0x8DDF
	MAX_GEOMETRY_OUTPUT_VERTICES                               = 0x8DE0
	MAX_GEOMETRY_TOTAL_OUTPUT_COMPONENTS                       = 0x8DE1
	MAX_VERTEX_OUTPUT_COMPONENTS                               = 0x9122
	MAX_GEOMETRY_INPUT_COMPONENTS                              = 0x9123
	MAX_GEOMETRY_OUTPUT_COMPONENTS                             = 0x9124
	MAX_FRAGMENT_INPUT_COMPONENTS                              = 0x9125
	CONTEXT_PROFILE_MASK                                       = 0x9126
	VERTEX_ATTRIB_ARRAY_DIVISOR                                = 0x88FE
	SAMPLE_SHADING                                             = 0x8C36
	MIN_SAMPLE_SHADING_VALUE                                   = 0x8C37
	MIN_PROGRAM_TEXTURE_GATHER_OFFSET                          = 0x8E5E
	MAX_PROGRAM_TEXTURE_GATHER_OFFSET                          = 0x8E5F
	TEXTURE_CUBE_MAP_ARRAY                                     = 0x9009
	TEXTURE_BINDING_CUBE_MAP_ARRAY                             = 0x900A
	PROXY_TEXTURE_CUBE_MAP_ARRAY                               = 0x900B
	SAMPLER_CUBE_MAP_ARRAY                                     = 0x900C
	SAMPLER_CUBE_MAP_ARRAY_SHADOW                              = 0x900D
	INT_SAMPLER_CUBE_MAP_ARRAY                                 = 0x900E
	UNSIGNED_INT_SAMPLER_CUBE_MAP_ARRAY                        = 0x900F
	NUM_SHADING_LANGUAGE_VERSIONS                              = 0x82E9
	VERTEX_ATTRIB_ARRAY_LONG                                   = 0x874E
	DEPTH_COMPONENT32F                                         = 0x8CAC
	DEPTH32F_STENCIL8                                          = 0x8CAD
	FLOAT_32_UNSIGNED_INT_24_8_REV                             = 0x8DAD
	INVALID_FRAMEBUFFER_OPERATION                              = 0x0506
	FRAMEBUFFER_ATTACHMENT_COLOR_ENCODING                      = 0x8210
	FRAMEBUFFER_ATTACHMENT_COMPONENT_TYPE                      = 0x8211
	FRAMEBUFFER_ATTACHMENT_RED_SIZE                            = 0x8212
	FRAMEBUFFER_ATTACHMENT_GREEN_SIZE                          = 0x8213
	FRAMEBUFFER_ATTACHMENT_BLUE_SIZE                           = 0x8214
	FRAMEBUFFER_ATTACHMENT_ALPHA_SIZE                          = 0x8215
	FRAMEBUFFER_ATTACHMENT_DEPTH_SIZE                          = 0x8216
	FRAMEBUFFER_ATTACHMENT_STENCIL_SIZE                        = 0x8217
	FRAMEBUFFER_DEFAULT                                        = 0x8218
	FRAMEBUFFER_UNDEFINED                                      = 0x8219
	DEPTH_STENCIL_ATTACHMENT                                   = 0x821A
	MAX_RENDERBUFFER_SIZE                                      = 0x84E8
	DEPTH_STENCIL                                              = 0x84F9
	UNSIGNED_INT_24_8                                          = 0x84FA
	DEPTH24_STENCIL8                                           = 0x88F0
	TEXTURE_STENCIL_SIZE                                       = 0x88F1
	TEXTURE_RED_TYPE                                           = 0x8C10
	TEXTURE_GREEN_TYPE                                         = 0x8C11
	TEXTURE_BLUE_TYPE                                          = 0x8C12
	TEXTURE_ALPHA_TYPE                                         = 0x8C13
	TEXTURE_DEPTH_TYPE                                         = 0x8C16
	UNSIGNED_NORMALIZED                                        = 0x8C17
	FRAMEBUFFER_BINDING                                        = 0x8CA6
	DRAW_FRAMEBUFFER_BINDING                                   = FRAMEBUFFER_BINDING
	RENDERBUFFER_BINDING                                       = 0x8CA7
	READ_FRAMEBUFFER                                           = 0x8CA8
	DRAW_FRAMEBUFFER                                           = 0x8CA9
	READ_FRAMEBUFFER_BINDING                                   = 0x8CAA
	RENDERBUFFER_SAMPLES                                       = 0x8CAB
	FRAMEBUFFER_ATTACHMENT_OBJECT_TYPE                         = 0x8CD0
	FRAMEBUFFER_ATTACHMENT_OBJECT_NAME                         = 0x8CD1
	FRAMEBUFFER_ATTACHMENT_TEXTURE_LEVEL                       = 0x8CD2
	FRAMEBUFFER_ATTACHMENT_TEXTURE_CUBE_MAP_FACE               = 0x8CD3
	FRAMEBUFFER_ATTACHMENT_TEXTURE_LAYER                       = 0x8CD4
	FRAMEBUFFER_COMPLETE                                       = 0x8CD5
	FRAMEBUFFER_INCOMPLETE_ATTACHMENT                          = 0x8CD6
	FRAMEBUFFER_INCOMPLETE_MISSING_ATTACHMENT                  = 0x8CD7
	FRAMEBUFFER_INCOMPLETE_DRAW_BUFFER                         = 0x8CDB
	FRAMEBUFFER_INCOMPLETE_READ_BUFFER                         = 0x8CDC
	FRAMEBUFFER_UNSUPPORTED                                    = 0x8CDD
	MAX_COLOR_ATTACHMENTS                                      = 0x8CDF
	COLOR_ATTACHMENT0                                          = 0x8CE0
	COLOR_ATTACHMENT1                                          = 0x8CE1
	COLOR_ATTACHMENT2                                          = 0x8CE2
	COLOR_ATTACHMENT3                                          = 0x8CE3
	COLOR_ATTACHMENT4                                          = 0x8CE4
	COLOR_ATTACHMENT5                                          = 0x8CE5
	COLOR_ATTACHMENT6                                          = 0x8CE6
	COLOR_ATTACHMENT7                                          = 0x8CE7
	COLOR_ATTACHMENT8                                          = 0x8CE8
	COLOR_ATTACHMENT9                                          = 0x8CE9
	COLOR_ATTACHMENT10                                         = 0x8CEA
	COLOR_ATTACHMENT11                                         = 0x8CEB
	COLOR_ATTACHMENT12                                         = 0x8CEC
	COLOR_ATTACHMENT13                                         = 0x8CED
	COLOR_ATTACHMENT14                                         = 0x8CEE
	COLOR_ATTACHMENT15                                         = 0x8CEF
	DEPTH_ATTACHMENT                                           = 0x8D00
	STENCIL_ATTACHMENT                                         = 0x8D20
	FRAMEBUFFER                                                = 0x8D40
	RENDERBUFFER                                               = 0x8D41
	RENDERBUFFER_WIDTH                                         = 0x8D42
	RENDERBUFFER_HEIGHT                                        = 0x8D43
	RENDERBUFFER_INTERNAL_FORMAT                               = 0x8D44
	STENCIL_INDEX1                                             = 0x8D46
	STENCIL_INDEX4                                             = 0x8D47
	STENCIL_INDEX8                                             = 0x8D48
	STENCIL_INDEX16                                            = 0x8D49
	RENDERBUFFER_RED_SIZE                                      = 0x8D50
	RENDERBUFFER_GREEN_SIZE                                    = 0x8D51
	RENDERBUFFER_BLUE_SIZE                                     = 0x8D52
	RENDERBUFFER_ALPHA_SIZE                                    = 0x8D53
	RENDERBUFFER_DEPTH_SIZE                                    = 0x8D54
	RENDERBUFFER_STENCIL_SIZE                                  = 0x8D55
	FRAMEBUFFER_INCOMPLETE_MULTISAMPLE                         = 0x8D56
	MAX_SAMPLES                                                = 0x8D57
	FRAMEBUFFER_SRGB                                           = 0x8DB9
	HALF_FLOAT                                                 = 0x140B
	MAP_READ_BIT                                               = 0x0001
	MAP_WRITE_BIT                                              = 0x0002
	MAP_INVALIDATE_RANGE_BIT                                   = 0x0004
	MAP_INVALIDATE_BUFFER_BIT                                  = 0x0008
	MAP_FLUSH_EXPLICIT_BIT                                     = 0x0010
	MAP_UNSYNCHRONIZED_BIT                                     = 0x0020
	COMPRESSED_RED_RGTC1                                       = 0x8DBB
	COMPRESSED_SIGNED_RED_RGTC1                                = 0x8DBC
	COMPRESSED_RG_RGTC2                                        = 0x8DBD
	COMPRESSED_SIGNED_RG_RGTC2                                 = 0x8DBE
	RG                                                         = 0x8227
	RG_INTEGER                                                 = 0x8228
	R8                                                         = 0x8229
	R16                                                        = 0x822A
	RG8                                                        = 0x822B
	RG16                                                       = 0x822C
	R16F                                                       = 0x822D
	R32F                                                       = 0x822E
	RG16F                                                      = 0x822F
	RG32F                                                      = 0x8230
	R8I                                                        = 0x8231
	R8UI                                                       = 0x8232
	R16I                                                       = 0x8233
	R16UI                                                      = 0x8234
	R32I                                                       = 0x8235
	R32UI                                                      = 0x8236
	RG8I                                                       = 0x8237
	RG8UI                                                      = 0x8238
	RG16I                                                      = 0x8239
	RG16UI                                                     = 0x823A
	RG32I                                                      = 0x823B
	RG32UI                                                     = 0x823C
	VERTEX_ARRAY_BINDING                                       = 0x85B5
	UNIFORM_BUFFER                                             = 0x8A11
	UNIFORM_BUFFER_BINDING                                     = 0x8A28
	UNIFORM_BUFFER_START                                       = 0x8A29
	UNIFORM_BUFFER_SIZE                                        = 0x8A2A
	MAX_VERTEX_UNIFORM_BLOCKS                                  = 0x8A2B
	MAX_GEOMETRY_UNIFORM_BLOCKS                                = 0x8A2C
	MAX_FRAGMENT_UNIFORM_BLOCKS                                = 0x8A2D
	MAX_COMBINED_UNIFORM_BLOCKS                                = 0x8A2E
	MAX_UNIFORM_BUFFER_BINDINGS                                = 0x8A2F
	MAX_UNIFORM_BLOCK_SIZE                                     = 0x8A30
	MAX_COMBINED_VERTEX_UNIFORM_COMPONENTS                     = 0x8A31
	MAX_COMBINED_GEOMETRY_UNIFORM_COMPONENTS                   = 0x8A32
	MAX_COMBINED_FRAGMENT_UNIFORM_COMPONENTS                   = 0x8A33
	UNIFORM_BUFFER_OFFSET_ALIGNMENT                            = 0x8A34
	ACTIVE_UNIFORM_BLOCK_MAX_NAME_LENGTH                       = 0x8A35
	ACTIVE_UNIFORM_BLOCKS                                      = 0x8A36
	UNIFORM_TYPE                                               = 0x8A37
	UNIFORM_SIZE                                               = 0x8A38
	UNIFORM_NAME_LENGTH                                        = 0x8A39
	UNIFORM_BLOCK_INDEX                                        = 0x8A3A
	UNIFORM_OFFSET                                             = 0x8A3B
	UNIFORM_ARRAY_STRIDE                                       = 0x8A3C
	UNIFORM_MATRIX_STRIDE                                      = 0x8A3D
	UNIFORM_IS_ROW_MAJOR                                       = 0x8A3E
	UNIFORM_BLOCK_BINDING                                      = 0x8A3F
	UNIFORM_BLOCK_DATA_SIZE                                    = 0x8A40
	UNIFORM_BLOCK_NAME_LENGTH                                  = 0x8A41
	UNIFORM_BLOCK_ACTIVE_UNIFORMS                              = 0x8A42
	UNIFORM_BLOCK_ACTIVE_UNIFORM_INDICES                       = 0x8A43
	UNIFORM_BLOCK_REFERENCED_BY_VERTEX_SHADER                  = 0x8A44
	UNIFORM_BLOCK_REFERENCED_BY_GEOMETRY_SHADER                = 0x8A45
	UNIFORM_BLOCK_REFERENCED_BY_FRAGMENT_SHADER                = 0x8A46
	INVALID_INDEX                                              = 0xFFFFFFFF
	COPY_READ_BUFFER_BINDING                                   = 0x8F36
	COPY_READ_BUFFER                                           = COPY_READ_BUFFER_BINDING
	COPY_WRITE_BUFFER_BINDING                                  = 0x8F37
	COPY_WRITE_BUFFER                                          = COPY_WRITE_BUFFER_BINDING
	DEPTH_CLAMP                                                = 0x864F
	QUADS_FOLLOW_PROVOKING_VERTEX_CONVENTION                   = 0x8E4C
	FIRST_VERTEX_CONVENTION                                    = 0x8E4D
	LAST_VERTEX_CONVENTION                                     = 0x8E4E
	PROVOKING_VERTEX                                           = 0x8E4F
	TEXTURE_CUBE_MAP_SEAMLESS                                  = 0x884F
	MAX_SERVER_WAIT_TIMEOUT                                    = 0x9111
	OBJECT_TYPE                                                = 0x9112
	SYNC_CONDITION                                             = 0x9113
	SYNC_STATUS                                                = 0x9114
	SYNC_FLAGS                                                 = 0x9115
	SYNC_FENCE                                                 = 0x9116
	SYNC_GPU_COMMANDS_COMPLETE                                 = 0x9117
	UNSIGNALED                                                 = 0x9118
	SIGNALED                                                   = 0x9119
	ALREADY_SIGNALED                                           = 0x911A
	TIMEOUT_EXPIRED                                            = 0x911B
	CONDITION_SATISFIED                                        = 0x911C
	WAIT_FAILED                                                = 0x911D
	SYNC_FLUSH_COMMANDS_BIT                                    = 0x00000001
	TIMEOUT_IGNORED                                            = 0xFFFFFFFFFFFFFFFF
	SAMPLE_POSITION                                            = 0x8E50
	SAMPLE_MASK                                                = 0x8E51
	SAMPLE_MASK_VALUE                                          = 0x8E52
	MAX_SAMPLE_MASK_WORDS                                      = 0x8E59
	TEXTURE_2D_MULTISAMPLE                                     = 0x9100
	PROXY_TEXTURE_2D_MULTISAMPLE                               = 0x9101
	TEXTURE_2D_MULTISAMPLE_ARRAY                               = 0x9102
	PROXY_TEXTURE_2D_MULTISAMPLE_ARRAY                         = 0x9103
	TEXTURE_BINDING_2D_MULTISAMPLE                             = 0x9104
	TEXTURE_BINDING_2D_MULTISAMPLE_ARRAY                       = 0x9105
	TEXTURE_SAMPLES                                            = 0x9106
	TEXTURE_FIXED_SAMPLE_LOCATIONS                             = 0x9107
	SAMPLER_2D_MULTISAMPLE                                     = 0x9108
	INT_SAMPLER_2D_MULTISAMPLE                                 = 0x9109
	UNSIGNED_INT_SAMPLER_2D_MULTISAMPLE                        = 0x910A
	SAMPLER_2D_MULTISAMPLE_ARRAY                               = 0x910B
	INT_SAMPLER_2D_MULTISAMPLE_ARRAY                           = 0x910C
	UNSIGNED_INT_SAMPLER_2D_MULTISAMPLE_ARRAY                  = 0x910D
	MAX_COLOR_TEXTURE_SAMPLES                                  = 0x910E
	MAX_DEPTH_TEXTURE_SAMPLES                                  = 0x910F
	MAX_INTEGER_SAMPLES                                        = 0x9110
	SAMPLE_SHADING_ARB                                         = 0x8C36
	MIN_SAMPLE_SHADING_VALUE_ARB                               = 0x8C37
	TEXTURE_CUBE_MAP_ARRAY_ARB                                 = 0x9009
	TEXTURE_BINDING_CUBE_MAP_ARRAY_ARB                         = 0x900A
	PROXY_TEXTURE_CUBE_MAP_ARRAY_ARB                           = 0x900B
	SAMPLER_CUBE_MAP_ARRAY_ARB                                 = 0x900C
	SAMPLER_CUBE_MAP_ARRAY_SHADOW_ARB                          = 0x900D
	INT_SAMPLER_CUBE_MAP_ARRAY_ARB                             = 0x900E
	UNSIGNED_INT_SAMPLER_CUBE_MAP_ARRAY_ARB                    = 0x900F
	MIN_PROGRAM_TEXTURE_GATHER_OFFSET_ARB                      = 0x8E5E
	MAX_PROGRAM_TEXTURE_GATHER_OFFSET_ARB                      = 0x8E5F
	SHADER_INCLUDE_ARB                                         = 0x8DAE
	NAMED_STRING_LENGTH_ARB                                    = 0x8DE9
	NAMED_STRING_TYPE_ARB                                      = 0x8DEA
	COMPRESSED_RGBA_BPTC_UNORM_ARB                             = 0x8E8C
	COMPRESSED_SRGB_ALPHA_BPTC_UNORM_ARB                       = 0x8E8D
	COMPRESSED_RGB_BPTC_SIGNED_FLOAT_ARB                       = 0x8E8E
	COMPRESSED_RGB_BPTC_UNSIGNED_FLOAT_ARB                     = 0x8E8F
	SRC1_COLOR                                                 = 0x88F9
	ONE_MINUS_SRC1_COLOR                                       = 0x88FA
	ONE_MINUS_SRC1_ALPHA                                       = 0x88FB
	MAX_DUAL_SOURCE_DRAW_BUFFERS                               = 0x88FC
	ANY_SAMPLES_PASSED                                         = 0x8C2F
	SAMPLER_BINDING                                            = 0x8919
	RGB10_A2UI                                                 = 0x906F
	TEXTURE_SWIZZLE_R                                          = 0x8E42
	TEXTURE_SWIZZLE_G                                          = 0x8E43
	TEXTURE_SWIZZLE_B                                          = 0x8E44
	TEXTURE_SWIZZLE_A                                          = 0x8E45
	TEXTURE_SWIZZLE_RGBA                                       = 0x8E46
	TIME_ELAPSED                                               = 0x88BF
	TIMESTAMP                                                  = 0x8E28
	INT_2_10_10_10_REV                                         = 0x8D9F
	DRAW_INDIRECT_BUFFER                                       = 0x8F3F
	DRAW_INDIRECT_BUFFER_BINDING                               = 0x8F43
	GEOMETRY_SHADER_INVOCATIONS                                = 0x887F
	MAX_GEOMETRY_SHADER_INVOCATIONS                            = 0x8E5A
	MIN_FRAGMENT_INTERPOLATION_OFFSET                          = 0x8E5B
	MAX_FRAGMENT_INTERPOLATION_OFFSET                          = 0x8E5C
	FRAGMENT_INTERPOLATION_OFFSET_BITS                         = 0x8E5D
	DOUBLE_VEC2                                                = 0x8FFC
	DOUBLE_VEC3                                                = 0x8FFD
	DOUBLE_VEC4                                                = 0x8FFE
	DOUBLE_MAT2                                                = 0x8F46
	DOUBLE_MAT3                                                = 0x8F47
	DOUBLE_MAT4                                                = 0x8F48
	DOUBLE_MAT2x3                                              = 0x8F49
	DOUBLE_MAT2x4                                              = 0x8F4A
	DOUBLE_MAT3x2                                              = 0x8F4B
	DOUBLE_MAT3x4                                              = 0x8F4C
	DOUBLE_MAT4x2                                              = 0x8F4D
	DOUBLE_MAT4x3                                              = 0x8F4E
	ACTIVE_SUBROUTINES                                         = 0x8DE5
	ACTIVE_SUBROUTINE_UNIFORMS                                 = 0x8DE6
	ACTIVE_SUBROUTINE_UNIFORM_LOCATIONS                        = 0x8E47
	ACTIVE_SUBROUTINE_MAX_LENGTH                               = 0x8E48
	ACTIVE_SUBROUTINE_UNIFORM_MAX_LENGTH                       = 0x8E49
	MAX_SUBROUTINES                                            = 0x8DE7
	MAX_SUBROUTINE_UNIFORM_LOCATIONS                           = 0x8DE8
	NUM_COMPATIBLE_SUBROUTINES                                 = 0x8E4A
	COMPATIBLE_SUBROUTINES                                     = 0x8E4B
	PATCHES                                                    = 0x000E
	PATCH_VERTICES                                             = 0x8E72
	PATCH_DEFAULT_INNER_LEVEL                                  = 0x8E73
	PATCH_DEFAULT_OUTER_LEVEL                                  = 0x8E74
	TESS_CONTROL_OUTPUT_VERTICES                               = 0x8E75
	TESS_GEN_MODE                                              = 0x8E76
	TESS_GEN_SPACING                                           = 0x8E77
	TESS_GEN_VERTEX_ORDER                                      = 0x8E78
	TESS_GEN_POINT_MODE                                        = 0x8E79
	ISOLINES                                                   = 0x8E7A
	FRACTIONAL_ODD                                             = 0x8E7B
	FRACTIONAL_EVEN                                            = 0x8E7C
	MAX_PATCH_VERTICES                                         = 0x8E7D
	MAX_TESS_GEN_LEVEL                                         = 0x8E7E
	MAX_TESS_CONTROL_UNIFORM_COMPONENTS                        = 0x8E7F
	MAX_TESS_EVALUATION_UNIFORM_COMPONENTS                     = 0x8E80
	MAX_TESS_CONTROL_TEXTURE_IMAGE_UNITS                       = 0x8E81
	MAX_TESS_EVALUATION_TEXTURE_IMAGE_UNITS                    = 0x8E82
	MAX_TESS_CONTROL_OUTPUT_COMPONENTS                         = 0x8E83
	MAX_TESS_PATCH_COMPONENTS                                  = 0x8E84
	MAX_TESS_CONTROL_TOTAL_OUTPUT_COMPONENTS                   = 0x8E85
	MAX_TESS_EVALUATION_OUTPUT_COMPONENTS                      = 0x8E86
	MAX_TESS_CONTROL_UNIFORM_BLOCKS                            = 0x8E89
	MAX_TESS_EVALUATION_UNIFORM_BLOCKS                         = 0x8E8A
	MAX_TESS_CONTROL_INPUT_COMPONENTS                          = 0x886C
	MAX_TESS_EVALUATION_INPUT_COMPONENTS                       = 0x886D
	MAX_COMBINED_TESS_CONTROL_UNIFORM_COMPONENTS               = 0x8E1E
	MAX_COMBINED_TESS_EVALUATION_UNIFORM_COMPONENTS            = 0x8E1F
	UNIFORM_BLOCK_REFERENCED_BY_TESS_CONTROL_SHADER            = 0x84F0
	UNIFORM_BLOCK_REFERENCED_BY_TESS_EVALUATION_SHADER         = 0x84F1
	TESS_EVALUATION_SHADER                                     = 0x8E87
	TESS_CONTROL_SHADER                                        = 0x8E88
	TRANSFORM_FEEDBACK                                         = 0x8E22
	TRANSFORM_FEEDBACK_PAUSED                                  = 0x8E23
	TRANSFORM_FEEDBACK_BUFFER_PAUSED                           = TRANSFORM_FEEDBACK_PAUSED
	TRANSFORM_FEEDBACK_ACTIVE                                  = 0x8E24
	TRANSFORM_FEEDBACK_BUFFER_ACTIVE                           = TRANSFORM_FEEDBACK_ACTIVE
	TRANSFORM_FEEDBACK_BINDING                                 = 0x8E25
	MAX_TRANSFORM_FEEDBACK_BUFFERS                             = 0x8E70
	MAX_VERTEX_STREAMS                                         = 0x8E71
	FIXED                                                      = 0x140C
	IMPLEMENTATION_COLOR_READ_TYPE                             = 0x8B9A
	IMPLEMENTATION_COLOR_READ_FORMAT                           = 0x8B9B
	LOW_FLOAT                                                  = 0x8DF0
	MEDIUM_FLOAT                                               = 0x8DF1
	HIGH_FLOAT                                                 = 0x8DF2
	LOW_INT                                                    = 0x8DF3
	MEDIUM_INT                                                 = 0x8DF4
	HIGH_INT                                                   = 0x8DF5
	SHADER_COMPILER                                            = 0x8DFA
	NUM_SHADER_BINARY_FORMATS                                  = 0x8DF9
	MAX_VERTEX_UNIFORM_VECTORS                                 = 0x8DFB
	MAX_VARYING_VECTORS                                        = 0x8DFC
	MAX_FRAGMENT_UNIFORM_VECTORS                               = 0x8DFD
	RGB565                                                     = 0x8D62
	PROGRAM_BINARY_RETRIEVABLE_HINT                            = 0x8257
	PROGRAM_BINARY_LENGTH                                      = 0x8741
	NUM_PROGRAM_BINARY_FORMATS                                 = 0x87FE
	PROGRAM_BINARY_FORMATS                                     = 0x87FF
	VERTEX_SHADER_BIT                                          = 0x00000001
	FRAGMENT_SHADER_BIT                                        = 0x00000002
	GEOMETRY_SHADER_BIT                                        = 0x00000004
	TESS_CONTROL_SHADER_BIT                                    = 0x00000008
	TESS_EVALUATION_SHADER_BIT                                 = 0x00000010
	ALL_SHADER_BITS                                            = 0xFFFFFFFF
	PROGRAM_SEPARABLE                                          = 0x8258
	ACTIVE_PROGRAM                                             = 0x8259
	PROGRAM_PIPELINE_BINDING                                   = 0x825A
	MAX_VIEWPORTS                                              = 0x825B
	VIEWPORT_SUBPIXEL_BITS                                     = 0x825C
	VIEWPORT_BOUNDS_RANGE                                      = 0x825D
	LAYER_PROVOKING_VERTEX                                     = 0x825E
	VIEWPORT_INDEX_PROVOKING_VERTEX                            = 0x825F
	UNDEFINED_VERTEX                                           = 0x8260
	SYNC_CL_EVENT_ARB                                          = 0x8240
	SYNC_CL_EVENT_COMPLETE_ARB                                 = 0x8241
	DEBUG_OUTPUT_SYNCHRONOUS_ARB                               = 0x8242
	DEBUG_NEXT_LOGGED_MESSAGE_LENGTH_ARB                       = 0x8243
	DEBUG_CALLBACK_FUNCTION_ARB                                = 0x8244
	DEBUG_CALLBACK_USER_PARAM_ARB                              = 0x8245
	DEBUG_SOURCE_API_ARB                                       = 0x8246
	DEBUG_SOURCE_WINDOW_SYSTEM_ARB                             = 0x8247
	DEBUG_SOURCE_SHADER_COMPILER_ARB                           = 0x8248
	DEBUG_SOURCE_THIRD_PARTY_ARB                               = 0x8249
	DEBUG_SOURCE_APPLICATION_ARB                               = 0x824A
	DEBUG_SOURCE_OTHER_ARB                                     = 0x824B
	DEBUG_TYPE_ERROR_ARB                                       = 0x824C
	DEBUG_TYPE_DEPRECATED_BEHAVIOR_ARB                         = 0x824D
	DEBUG_TYPE_UNDEFINED_BEHAVIOR_ARB                          = 0x824E
	DEBUG_TYPE_PORTABILITY_ARB                                 = 0x824F
	DEBUG_TYPE_PERFORMANCE_ARB                                 = 0x8250
	DEBUG_TYPE_OTHER_ARB                                       = 0x8251
	MAX_DEBUG_MESSAGE_LENGTH_ARB                               = 0x9143
	MAX_DEBUG_LOGGED_MESSAGES_ARB                              = 0x9144
	DEBUG_LOGGED_MESSAGES_ARB                                  = 0x9145
	DEBUG_SEVERITY_HIGH_ARB                                    = 0x9146
	DEBUG_SEVERITY_MEDIUM_ARB                                  = 0x9147
	DEBUG_SEVERITY_LOW_ARB                                     = 0x9148
	CONTEXT_FLAG_ROBUST_ACCESS_BIT_ARB                         = 0x00000004
	LOSE_CONTEXT_ON_RESET_ARB                                  = 0x8252
	GUILTY_CONTEXT_RESET_ARB                                   = 0x8253
	INNOCENT_CONTEXT_RESET_ARB                                 = 0x8254
	UNKNOWN_CONTEXT_RESET_ARB                                  = 0x8255
	RESET_NOTIFICATION_STRATEGY_ARB                            = 0x8256
	NO_RESET_NOTIFICATION_ARB                                  = 0x8261
	UNPACK_COMPRESSED_BLOCK_WIDTH                              = 0x9127
	UNPACK_COMPRESSED_BLOCK_HEIGHT                             = 0x9128
	UNPACK_COMPRESSED_BLOCK_DEPTH                              = 0x9129
	UNPACK_COMPRESSED_BLOCK_SIZE                               = 0x912A
	PACK_COMPRESSED_BLOCK_WIDTH                                = 0x912B
	PACK_COMPRESSED_BLOCK_HEIGHT                               = 0x912C
	PACK_COMPRESSED_BLOCK_DEPTH                                = 0x912D
	PACK_COMPRESSED_BLOCK_SIZE                                 = 0x912E
	NUM_SAMPLE_COUNTS                                          = 0x9380
	MIN_MAP_BUFFER_ALIGNMENT                                   = 0x90BC
	ATOMIC_COUNTER_BUFFER                                      = 0x92C0
	ATOMIC_COUNTER_BUFFER_BINDING                              = 0x92C1
	ATOMIC_COUNTER_BUFFER_START                                = 0x92C2
	ATOMIC_COUNTER_BUFFER_SIZE                                 = 0x92C3
	ATOMIC_COUNTER_BUFFER_DATA_SIZE                            = 0x92C4
	ATOMIC_COUNTER_BUFFER_ACTIVE_ATOMIC_COUNTERS               = 0x92C5
	ATOMIC_COUNTER_BUFFER_ACTIVE_ATOMIC_COUNTER_INDICES        = 0x92C6
	ATOMIC_COUNTER_BUFFER_REFERENCED_BY_VERTEX_SHADER          = 0x92C7
	ATOMIC_COUNTER_BUFFER_REFERENCED_BY_TESS_CONTROL_SHADER    = 0x92C8
	ATOMIC_COUNTER_BUFFER_REFERENCED_BY_TESS_EVALUATION_SHADER = 0x92C9
	ATOMIC_COUNTER_BUFFER_REFERENCED_BY_GEOMETRY_SHADER        = 0x92CA
	ATOMIC_COUNTER_BUFFER_REFERENCED_BY_FRAGMENT_SHADER        = 0x92CB
	MAX_VERTEX_ATOMIC_COUNTER_BUFFERS                          = 0x92CC
	MAX_TESS_CONTROL_ATOMIC_COUNTER_BUFFERS                    = 0x92CD
	MAX_TESS_EVALUATION_ATOMIC_COUNTER_BUFFERS                 = 0x92CE
	MAX_GEOMETRY_ATOMIC_COUNTER_BUFFERS                        = 0x92CF
	MAX_FRAGMENT_ATOMIC_COUNTER_BUFFERS                        = 0x92D0
	MAX_COMBINED_ATOMIC_COUNTER_BUFFERS                        = 0x92D1
	MAX_VERTEX_ATOMIC_COUNTERS                                 = 0x92D2
	MAX_TESS_CONTROL_ATOMIC_COUNTERS                           = 0x92D3
	MAX_TESS_EVALUATION_ATOMIC_COUNTERS                        = 0x92D4
	MAX_GEOMETRY_ATOMIC_COUNTERS                               = 0x92D5
	MAX_FRAGMENT_ATOMIC_COUNTERS                               = 0x92D6
	MAX_COMBINED_ATOMIC_COUNTERS                               = 0x92D7
	MAX_ATOMIC_COUNTER_BUFFER_SIZE                             = 0x92D8
	MAX_ATOMIC_COUNTER_BUFFER_BINDINGS                         = 0x92DC
	ACTIVE_ATOMIC_COUNTER_BUFFERS                              = 0x92D9
	UNIFORM_ATOMIC_COUNTER_BUFFER_INDEX                        = 0x92DA
	UNSIGNED_INT_ATOMIC_COUNTER                                = 0x92DB
	VERTEX_ATTRIB_ARRAY_BARRIER_BIT                            = 0x00000001
	ELEMENT_ARRAY_BARRIER_BIT                                  = 0x00000002
	UNIFORM_BARRIER_BIT                                        = 0x00000004
	TEXTURE_FETCH_BARRIER_BIT                                  = 0x00000008
	SHADER_IMAGE_ACCESS_BARRIER_BIT                            = 0x00000020
	COMMAND_BARRIER_BIT                                        = 0x00000040
	PIXEL_BUFFER_BARRIER_BIT                                   = 0x00000080
	TEXTURE_UPDATE_BARRIER_BIT                                 = 0x00000100
	BUFFER_UPDATE_BARRIER_BIT                                  = 0x00000200
	FRAMEBUFFER_BARRIER_BIT                                    = 0x00000400
	TRANSFORM_FEEDBACK_BARRIER_BIT                             = 0x00000800
	ATOMIC_COUNTER_BARRIER_BIT                                 = 0x00001000
	ALL_BARRIER_BITS                                           = 0xFFFFFFFF
	MAX_IMAGE_UNITS                                            = 0x8F38
	MAX_COMBINED_IMAGE_UNITS_AND_FRAGMENT_OUTPUTS              = 0x8F39
	IMAGE_BINDING_NAME                                         = 0x8F3A
	IMAGE_BINDING_LEVEL                                        = 0x8F3B
	IMAGE_BINDING_LAYERED                                      = 0x8F3C
	IMAGE_BINDING_LAYER                                        = 0x8F3D
	IMAGE_BINDING_ACCESS                                       = 0x8F3E
	IMAGE_1D                                                   = 0x904C
	IMAGE_2D                                                   = 0x904D
	IMAGE_3D                                                   = 0x904E
	IMAGE_2D_RECT                                              = 0x904F
	IMAGE_CUBE                                                 = 0x9050
	IMAGE_BUFFER                                               = 0x9051
	IMAGE_1D_ARRAY                                             = 0x9052
	IMAGE_2D_ARRAY                                             = 0x9053
	IMAGE_CUBE_MAP_ARRAY                                       = 0x9054
	IMAGE_2D_MULTISAMPLE                                       = 0x9055
	IMAGE_2D_MULTISAMPLE_ARRAY                                 = 0x9056
	INT_IMAGE_1D                                               = 0x9057
	INT_IMAGE_2D                                               = 0x9058
	INT_IMAGE_3D                                               = 0x9059
	INT_IMAGE_2D_RECT                                          = 0x905A
	INT_IMAGE_CUBE                                             = 0x905B
	INT_IMAGE_BUFFER                                           = 0x905C
	INT_IMAGE_1D_ARRAY                                         = 0x905D
	INT_IMAGE_2D_ARRAY                                         = 0x905E
	INT_IMAGE_CUBE_MAP_ARRAY                                   = 0x905F
	INT_IMAGE_2D_MULTISAMPLE                                   = 0x9060
	INT_IMAGE_2D_MULTISAMPLE_ARRAY                             = 0x9061
	UNSIGNED_INT_IMAGE_1D                                      = 0x9062
	UNSIGNED_INT_IMAGE_2D                                      = 0x9063
	UNSIGNED_INT_IMAGE_3D                                      = 0x9064
	UNSIGNED_INT_IMAGE_2D_RECT                                 = 0x9065
	UNSIGNED_INT_IMAGE_CUBE                                    = 0x9066
	UNSIGNED_INT_IMAGE_BUFFER                                  = 0x9067
	UNSIGNED_INT_IMAGE_1D_ARRAY                                = 0x9068
	UNSIGNED_INT_IMAGE_2D_ARRAY                                = 0x9069
	UNSIGNED_INT_IMAGE_CUBE_MAP_ARRAY                          = 0x906A
	UNSIGNED_INT_IMAGE_2D_MULTISAMPLE                          = 0x906B
	UNSIGNED_INT_IMAGE_2D_MULTISAMPLE_ARRAY                    = 0x906C
	MAX_IMAGE_SAMPLES                                          = 0x906D
	IMAGE_BINDING_FORMAT                                       = 0x906E
	IMAGE_FORMAT_COMPATIBILITY_TYPE                            = 0x90C7
	IMAGE_FORMAT_COMPATIBILITY_BY_SIZE                         = 0x90C8
	IMAGE_FORMAT_COMPATIBILITY_BY_CLASS                        = 0x90C9
	MAX_VERTEX_IMAGE_UNIFORMS                                  = 0x90CA
	MAX_TESS_CONTROL_IMAGE_UNIFORMS                            = 0x90CB
	MAX_TESS_EVALUATION_IMAGE_UNIFORMS                         = 0x90CC
	MAX_GEOMETRY_IMAGE_UNIFORMS                                = 0x90CD
	MAX_FRAGMENT_IMAGE_UNIFORMS                                = 0x90CE
	MAX_COMBINED_IMAGE_UNIFORMS                                = 0x90CF
	TEXTURE_IMMUTABLE_FORMAT                                   = 0x912F
	COMPRESSED_RGBA_ASTC_4x4_KHR                               = 0x93B0
	COMPRESSED_RGBA_ASTC_5x4_KHR                               = 0x93B1
	COMPRESSED_RGBA_ASTC_5x5_KHR                               = 0x93B2
	COMPRESSED_RGBA_ASTC_6x5_KHR                               = 0x93B3
	COMPRESSED_RGBA_ASTC_6x6_KHR                               = 0x93B4
	COMPRESSED_RGBA_ASTC_8x5_KHR                               = 0x93B5
	COMPRESSED_RGBA_ASTC_8x6_KHR                               = 0x93B6
	COMPRESSED_RGBA_ASTC_8x8_KHR                               = 0x93B7
	COMPRESSED_RGBA_ASTC_10x5_KHR                              = 0x93B8
	COMPRESSED_RGBA_ASTC_10x6_KHR                              = 0x93B9
	COMPRESSED_RGBA_ASTC_10x8_KHR                              = 0x93BA
	COMPRESSED_RGBA_ASTC_10x10_KHR                             = 0x93BB
	COMPRESSED_RGBA_ASTC_12x10_KHR                             = 0x93BC
	COMPRESSED_RGBA_ASTC_12x12_KHR                             = 0x93BD
	COMPRESSED_SRGB8_ALPHA8_ASTC_4x4_KHR                       = 0x93D0
	COMPRESSED_SRGB8_ALPHA8_ASTC_5x4_KHR                       = 0x93D1
	COMPRESSED_SRGB8_ALPHA8_ASTC_5x5_KHR                       = 0x93D2
	COMPRESSED_SRGB8_ALPHA8_ASTC_6x5_KHR                       = 0x93D3
	COMPRESSED_SRGB8_ALPHA8_ASTC_6x6_KHR                       = 0x93D4
	COMPRESSED_SRGB8_ALPHA8_ASTC_8x5_KHR                       = 0x93D5
	COMPRESSED_SRGB8_ALPHA8_ASTC_8x6_KHR                       = 0x93D6
	COMPRESSED_SRGB8_ALPHA8_ASTC_8x8_KHR                       = 0x93D7
	COMPRESSED_SRGB8_ALPHA8_ASTC_10x5_KHR                      = 0x93D8
	COMPRESSED_SRGB8_ALPHA8_ASTC_10x6_KHR                      = 0x93D9
	COMPRESSED_SRGB8_ALPHA8_ASTC_10x8_KHR                      = 0x93DA
	COMPRESSED_SRGB8_ALPHA8_ASTC_10x10_KHR                     = 0x93DB
	COMPRESSED_SRGB8_ALPHA8_ASTC_12x10_KHR                     = 0x93DC
	COMPRESSED_SRGB8_ALPHA8_ASTC_12x12_KHR                     = 0x93DD
	DEBUG_OUTPUT_SYNCHRONOUS                                   = 0x8242
	DEBUG_NEXT_LOGGED_MESSAGE_LENGTH                           = 0x8243
	DEBUG_CALLBACK_FUNCTION                                    = 0x8244
	DEBUG_CALLBACK_USER_PARAM                                  = 0x8245
	DEBUG_SOURCE_API                                           = 0x8246
	DEBUG_SOURCE_WINDOW_SYSTEM                                 = 0x8247
	DEBUG_SOURCE_SHADER_COMPILER                               = 0x8248
	DEBUG_SOURCE_THIRD_PARTY                                   = 0x8249
	DEBUG_SOURCE_APPLICATION                                   = 0x824A
	DEBUG_SOURCE_OTHER                                         = 0x824B
	DEBUG_TYPE_ERROR                                           = 0x824C
	DEBUG_TYPE_DEPRECATED_BEHAVIOR                             = 0x824D
	DEBUG_TYPE_UNDEFINED_BEHAVIOR                              = 0x824E
	DEBUG_TYPE_PORTABILITY                                     = 0x824F
	DEBUG_TYPE_PERFORMANCE                                     = 0x8250
	DEBUG_TYPE_OTHER                                           = 0x8251
	DEBUG_TYPE_MARKER                                          = 0x8268
	DEBUG_TYPE_PUSH_GROUP                                      = 0x8269
	DEBUG_TYPE_POP_GROUP                                       = 0x826A
	DEBUG_SEVERITY_NOTIFICATION                                = 0x826B
	MAX_DEBUG_GROUP_STACK_DEPTH                                = 0x826C
	DEBUG_GROUP_STACK_DEPTH                                    = 0x826D
	BUFFER                                                     = 0x82E0
	SHADER                                                     = 0x82E1
	PROGRAM                                                    = 0x82E2
	QUERY                                                      = 0x82E3
	PROGRAM_PIPELINE                                           = 0x82E4
	SAMPLER                                                    = 0x82E6
	DISPLAY_LIST                                               = 0x82E7
	MAX_LABEL_LENGTH                                           = 0x82E8
	MAX_DEBUG_MESSAGE_LENGTH                                   = 0x9143
	MAX_DEBUG_LOGGED_MESSAGES                                  = 0x9144
	DEBUG_LOGGED_MESSAGES                                      = 0x9145
	DEBUG_SEVERITY_HIGH                                        = 0x9146
	DEBUG_SEVERITY_MEDIUM                                      = 0x9147
	DEBUG_SEVERITY_LOW                                         = 0x9148
	DEBUG_OUTPUT                                               = 0x92E0
	CONTEXT_FLAG_DEBUG_BIT                                     = 0x00000002
	COMPUTE_SHADER                                             = 0x91B9
	MAX_COMPUTE_UNIFORM_BLOCKS                                 = 0x91BB
	MAX_COMPUTE_TEXTURE_IMAGE_UNITS                            = 0x91BC
	MAX_COMPUTE_IMAGE_UNIFORMS                                 = 0x91BD
	MAX_COMPUTE_SHARED_MEMORY_SIZE                             = 0x8262
	MAX_COMPUTE_UNIFORM_COMPONENTS                             = 0x8263
	MAX_COMPUTE_ATOMIC_COUNTER_BUFFERS                         = 0x8264
	MAX_COMPUTE_ATOMIC_COUNTERS                                = 0x8265
	MAX_COMBINED_COMPUTE_UNIFORM_COMPONENTS                    = 0x8266
	MAX_COMPUTE_LOCAL_INVOCATIONS                              = 0x90EB
	MAX_COMPUTE_WORK_GROUP_COUNT                               = 0x91BE
	MAX_COMPUTE_WORK_GROUP_SIZE                                = 0x91BF
	COMPUTE_LOCAL_WORK_SIZE                                    = 0x8267
	UNIFORM_BLOCK_REFERENCED_BY_COMPUTE_SHADER                 = 0x90EC
	ATOMIC_COUNTER_BUFFER_REFERENCED_BY_COMPUTE_SHADER         = 0x90ED
	DISPATCH_INDIRECT_BUFFER                                   = 0x90EE
	DISPATCH_INDIRECT_BUFFER_BINDING                           = 0x90EF
	COMPUTE_SHADER_BIT                                         = 0x00000020
	TEXTURE_VIEW_MIN_LEVEL                                     = 0x82DB
	TEXTURE_VIEW_NUM_LEVELS                                    = 0x82DC
	TEXTURE_VIEW_MIN_LAYER                                     = 0x82DD
	TEXTURE_VIEW_NUM_LAYERS                                    = 0x82DE
	TEXTURE_IMMUTABLE_LEVELS                                   = 0x82DF
	VERTEX_ATTRIB_BINDING                                      = 0x82D4
	VERTEX_ATTRIB_RELATIVE_OFFSET                              = 0x82D5
	VERTEX_BINDING_DIVISOR                                     = 0x82D6
	VERTEX_BINDING_OFFSET                                      = 0x82D7
	VERTEX_BINDING_STRIDE                                      = 0x82D8
	MAX_VERTEX_ATTRIB_RELATIVE_OFFSET                          = 0x82D9
	MAX_VERTEX_ATTRIB_BINDINGS                                 = 0x82DA
	COMPRESSED_RGB8_ETC2                                       = 0x9274
	COMPRESSED_SRGB8_ETC2                                      = 0x9275
	COMPRESSED_RGB8_PUNCHTHROUGH_ALPHA1_ETC2                   = 0x9276
	COMPRESSED_SRGB8_PUNCHTHROUGH_ALPHA1_ETC2                  = 0x9277
	COMPRESSED_RGBA8_ETC2_EAC                                  = 0x9278
	COMPRESSED_SRGB8_ALPHA8_ETC2_EAC                           = 0x9279
	COMPRESSED_R11_EAC                                         = 0x9270
	COMPRESSED_SIGNED_R11_EAC                                  = 0x9271
	COMPRESSED_RG11_EAC                                        = 0x9272
	COMPRESSED_SIGNED_RG11_EAC                                 = 0x9273
	PRIMITIVE_RESTART_FIXED_INDEX                              = 0x8D69
	ANY_SAMPLES_PASSED_CONSERVATIVE                            = 0x8D6A
	MAX_ELEMENT_INDEX                                          = 0x8D6B
	MAX_UNIFORM_LOCATIONS                                      = 0x826E
	FRAMEBUFFER_DEFAULT_WIDTH                                  = 0x9310
	FRAMEBUFFER_DEFAULT_HEIGHT                                 = 0x9311
	FRAMEBUFFER_DEFAULT_LAYERS                                 = 0x9312
	FRAMEBUFFER_DEFAULT_SAMPLES                                = 0x9313
	FRAMEBUFFER_DEFAULT_FIXED_SAMPLE_LOCATIONS                 = 0x9314
	MAX_FRAMEBUFFER_WIDTH                                      = 0x9315
	MAX_FRAMEBUFFER_HEIGHT                                     = 0x9316
	MAX_FRAMEBUFFER_LAYERS                                     = 0x9317
	MAX_FRAMEBUFFER_SAMPLES                                    = 0x9318
	INTERNALFORMAT_SUPPORTED                                   = 0x826F
	INTERNALFORMAT_PREFERRED                                   = 0x8270
	INTERNALFORMAT_RED_SIZE                                    = 0x8271
	INTERNALFORMAT_GREEN_SIZE                                  = 0x8272
	INTERNALFORMAT_BLUE_SIZE                                   = 0x8273
	INTERNALFORMAT_ALPHA_SIZE                                  = 0x8274
	INTERNALFORMAT_DEPTH_SIZE                                  = 0x8275
	INTERNALFORMAT_STENCIL_SIZE                                = 0x8276
	INTERNALFORMAT_SHARED_SIZE                                 = 0x8277
	INTERNALFORMAT_RED_TYPE                                    = 0x8278
	INTERNALFORMAT_GREEN_TYPE                                  = 0x8279
	INTERNALFORMAT_BLUE_TYPE                                   = 0x827A
	INTERNALFORMAT_ALPHA_TYPE                                  = 0x827B
	INTERNALFORMAT_DEPTH_TYPE                                  = 0x827C
	INTERNALFORMAT_STENCIL_TYPE                                = 0x827D
	MAX_WIDTH                                                  = 0x827E
	MAX_HEIGHT                                                 = 0x827F
	MAX_DEPTH                                                  = 0x8280
	MAX_LAYERS                                                 = 0x8281
	MAX_COMBINED_DIMENSIONS                                    = 0x8282
	COLOR_COMPONENTS                                           = 0x8283
	DEPTH_COMPONENTS                                           = 0x8284
	STENCIL_COMPONENTS                                         = 0x8285
	COLOR_RENDERABLE                                           = 0x8286
	DEPTH_RENDERABLE                                           = 0x8287
	STENCIL_RENDERABLE                                         = 0x8288
	FRAMEBUFFER_RENDERABLE                                     = 0x8289
	FRAMEBUFFER_RENDERABLE_LAYERED                             = 0x828A
	FRAMEBUFFER_BLEND                                          = 0x828B
	READ_PIXELS                                                = 0x828C
	READ_PIXELS_FORMAT                                         = 0x828D
	READ_PIXELS_TYPE                                           = 0x828E
	TEXTURE_IMAGE_FORMAT                                       = 0x828F
	TEXTURE_IMAGE_TYPE                                         = 0x8290
	GET_TEXTURE_IMAGE_FORMAT                                   = 0x8291
	GET_TEXTURE_IMAGE_TYPE                                     = 0x8292
	MIPMAP                                                     = 0x8293
	MANUAL_GENERATE_MIPMAP                                     = 0x8294
	AUTO_GENERATE_MIPMAP                                       = 0x8295
	COLOR_ENCODING                                             = 0x8296
	SRGB_READ                                                  = 0x8297
	SRGB_WRITE                                                 = 0x8298
	SRGB_DECODE_ARB                                            = 0x8299
	FILTER                                                     = 0x829A
	VERTEX_TEXTURE                                             = 0x829B
	TESS_CONTROL_TEXTURE                                       = 0x829C
	TESS_EVALUATION_TEXTURE                                    = 0x829D
	GEOMETRY_TEXTURE                                           = 0x829E
	FRAGMENT_TEXTURE                                           = 0x829F
	COMPUTE_TEXTURE                                            = 0x82A0
	TEXTURE_SHADOW                                             = 0x82A1
	TEXTURE_GATHER                                             = 0x82A2
	TEXTURE_GATHER_SHADOW                                      = 0x82A3
	SHADER_IMAGE_LOAD                                          = 0x82A4
	SHADER_IMAGE_STORE                                         = 0x82A5
	SHADER_IMAGE_ATOMIC                                        = 0x82A6
	IMAGE_TEXEL_SIZE                                           = 0x82A7
	IMAGE_COMPATIBILITY_CLASS                                  = 0x82A8
	IMAGE_PIXEL_FORMAT                                         = 0x82A9
	IMAGE_PIXEL_TYPE                                           = 0x82AA
	SIMULTANEOUS_TEXTURE_AND_DEPTH_TEST                        = 0x82AC
	SIMULTANEOUS_TEXTURE_AND_STENCIL_TEST                      = 0x82AD
	SIMULTANEOUS_TEXTURE_AND_DEPTH_WRITE                       = 0x82AE
	SIMULTANEOUS_TEXTURE_AND_STENCIL_WRITE                     = 0x82AF
	TEXTURE_COMPRESSED_BLOCK_WIDTH                             = 0x82B1
	TEXTURE_COMPRESSED_BLOCK_HEIGHT                            = 0x82B2
	TEXTURE_COMPRESSED_BLOCK_SIZE                              = 0x82B3
	CLEAR_BUFFER                                               = 0x82B4
	TEXTURE_VIEW                                               = 0x82B5
	VIEW_COMPATIBILITY_CLASS                                   = 0x82B6
	FULL_SUPPORT                                               = 0x82B7
	CAVEAT_SUPPORT                                             = 0x82B8
	IMAGE_CLASS_4_X_32                                         = 0x82B9
	IMAGE_CLASS_2_X_32                                         = 0x82BA
	IMAGE_CLASS_1_X_32                                         = 0x82BB
	IMAGE_CLASS_4_X_16                                         = 0x82BC
	IMAGE_CLASS_2_X_16                                         = 0x82BD
	IMAGE_CLASS_1_X_16                                         = 0x82BE
	IMAGE_CLASS_4_X_8                                          = 0x82BF
	IMAGE_CLASS_2_X_8                                          = 0x82C0
	IMAGE_CLASS_1_X_8                                          = 0x82C1
	IMAGE_CLASS_11_11_10                                       = 0x82C2
	IMAGE_CLASS_10_10_10_2                                     = 0x82C3
	VIEW_CLASS_128_BITS                                        = 0x82C4
	VIEW_CLASS_96_BITS                                         = 0x82C5
	VIEW_CLASS_64_BITS                                         = 0x82C6
	VIEW_CLASS_48_BITS                                         = 0x82C7
	VIEW_CLASS_32_BITS                                         = 0x82C8
	VIEW_CLASS_24_BITS                                         = 0x82C9
	VIEW_CLASS_16_BITS                                         = 0x82CA
	VIEW_CLASS_8_BITS                                          = 0x82CB
	VIEW_CLASS_S3TC_DXT1_RGB                                   = 0x82CC
	VIEW_CLASS_S3TC_DXT1_RGBA                                  = 0x82CD
	VIEW_CLASS_S3TC_DXT3_RGBA                                  = 0x82CE
	VIEW_CLASS_S3TC_DXT5_RGBA                                  = 0x82CF
	VIEW_CLASS_RGTC1_RED                                       = 0x82D0
	VIEW_CLASS_RGTC2_RG                                        = 0x82D1
	VIEW_CLASS_BPTC_UNORM                                      = 0x82D2
	VIEW_CLASS_BPTC_FLOAT                                      = 0x82D3
	UNIFORM                                                    = 0x92E1
	UNIFORM_BLOCK                                              = 0x92E2
	PROGRAM_INPUT                                              = 0x92E3
	PROGRAM_OUTPUT                                             = 0x92E4
	BUFFER_VARIABLE                                            = 0x92E5
	SHADER_STORAGE_BLOCK                                       = 0x92E6
	VERTEX_SUBROUTINE                                          = 0x92E8
	TESS_CONTROL_SUBROUTINE                                    = 0x92E9
	TESS_EVALUATION_SUBROUTINE                                 = 0x92EA
	GEOMETRY_SUBROUTINE                                        = 0x92EB
	FRAGMENT_SUBROUTINE                                        = 0x92EC
	COMPUTE_SUBROUTINE                                         = 0x92ED
	VERTEX_SUBROUTINE_UNIFORM                                  = 0x92EE
	TESS_CONTROL_SUBROUTINE_UNIFORM                            = 0x92EF
	TESS_EVALUATION_SUBROUTINE_UNIFORM                         = 0x92F0
	GEOMETRY_SUBROUTINE_UNIFORM                                = 0x92F1
	FRAGMENT_SUBROUTINE_UNIFORM                                = 0x92F2
	COMPUTE_SUBROUTINE_UNIFORM                                 = 0x92F3
	TRANSFORM_FEEDBACK_VARYING                                 = 0x92F4
	ACTIVE_RESOURCES                                           = 0x92F5
	MAX_NAME_LENGTH                                            = 0x92F6
	MAX_NUM_ACTIVE_VARIABLES                                   = 0x92F7
	MAX_NUM_COMPATIBLE_SUBROUTINES                             = 0x92F8
	NAME_LENGTH                                                = 0x92F9
	TYPE                                                       = 0x92FA
	ARRAY_SIZE                                                 = 0x92FB
	OFFSET                                                     = 0x92FC
	BLOCK_INDEX                                                = 0x92FD
	ARRAY_STRIDE                                               = 0x92FE
	MATRIX_STRIDE                                              = 0x92FF
	IS_ROW_MAJOR                                               = 0x9300
	ATOMIC_COUNTER_BUFFER_INDEX                                = 0x9301
	BUFFER_BINDING                                             = 0x9302
	BUFFER_DATA_SIZE                                           = 0x9303
	NUM_ACTIVE_VARIABLES                                       = 0x9304
	ACTIVE_VARIABLES                                           = 0x9305
	REFERENCED_BY_VERTEX_SHADER                                = 0x9306
	REFERENCED_BY_TESS_CONTROL_SHADER                          = 0x9307
	REFERENCED_BY_TESS_EVALUATION_SHADER                       = 0x9308
	REFERENCED_BY_GEOMETRY_SHADER                              = 0x9309
	REFERENCED_BY_FRAGMENT_SHADER                              = 0x930A
	REFERENCED_BY_COMPUTE_SHADER                               = 0x930B
	TOP_LEVEL_ARRAY_SIZE                                       = 0x930C
	TOP_LEVEL_ARRAY_STRIDE                                     = 0x930D
	LOCATION                                                   = 0x930E
	LOCATION_INDEX                                             = 0x930F
	IS_PER_PATCH                                               = 0x92E7
	SHADER_STORAGE_BUFFER                                      = 0x90D2
	SHADER_STORAGE_BUFFER_BINDING                              = 0x90D3
	SHADER_STORAGE_BUFFER_START                                = 0x90D4
	SHADER_STORAGE_BUFFER_SIZE                                 = 0x90D5
	MAX_VERTEX_SHADER_STORAGE_BLOCKS                           = 0x90D6
	MAX_GEOMETRY_SHADER_STORAGE_BLOCKS                         = 0x90D7
	MAX_TESS_CONTROL_SHADER_STORAGE_BLOCKS                     = 0x90D8
	MAX_TESS_EVALUATION_SHADER_STORAGE_BLOCKS                  = 0x90D9
	MAX_FRAGMENT_SHADER_STORAGE_BLOCKS                         = 0x90DA
	MAX_COMPUTE_SHADER_STORAGE_BLOCKS                          = 0x90DB
	MAX_COMBINED_SHADER_STORAGE_BLOCKS                         = 0x90DC
	MAX_SHADER_STORAGE_BUFFER_BINDINGS                         = 0x90DD
	MAX_SHADER_STORAGE_BLOCK_SIZE                              = 0x90DE
	SHADER_STORAGE_BUFFER_OFFSET_ALIGNMENT                     = 0x90DF
	SHADER_STORAGE_BARRIER_BIT                                 = 0x2000
	MAX_COMBINED_SHADER_OUTPUT_RESOURCES                       = MAX_COMBINED_IMAGE_UNITS_AND_FRAGMENT_OUTPUTS
	DEPTH_STENCIL_TEXTURE_MODE                                 = 0x90EA
	TEXTURE_BUFFER_OFFSET                                      = 0x919D
	TEXTURE_BUFFER_SIZE                                        = 0x919E
	TEXTURE_BUFFER_OFFSET_ALIGNMENT                            = 0x919F
	KHR_texture_compression_astc_ldr                           = 1
	KHR_debug                                                  = 1
)

func IsBuffer(buffer uint32) bool {
	return cbool(uint(C.wrap_glIsBuffer(C.uint(buffer))))
}
func IsEnabled(cap uint32) bool {
	return cbool(uint(C.wrap_glIsEnabled(C.uint(cap))))
}
func IsEnabledi(target uint32, index uint32) bool {
	return cbool(uint(C.wrap_glIsEnabledi(C.uint(target), C.uint(index))))
}
func IsFramebuffer(framebuffer uint32) bool {
	return cbool(uint(C.wrap_glIsFramebuffer(C.uint(framebuffer))))
}
func IsNamedStringARB(namelen int32, name string) bool {
	cstr1 := C.CString(name)
	defer C.free(unsafe.Pointer(cstr1))
	return cbool(uint(C.wrap_glIsNamedStringARB(C.int(namelen), cstr1)))
}
func IsProgram(program uint32) bool {
	return cbool(uint(C.wrap_glIsProgram(C.uint(program))))
}
func IsProgramPipeline(pipeline uint32) bool {
	return cbool(uint(C.wrap_glIsProgramPipeline(C.uint(pipeline))))
}
func IsQuery(id uint32) bool {
	return cbool(uint(C.wrap_glIsQuery(C.uint(id))))
}
func IsRenderbuffer(renderbuffer uint32) bool {
	return cbool(uint(C.wrap_glIsRenderbuffer(C.uint(renderbuffer))))
}
func IsSampler(sampler uint32) bool {
	return cbool(uint(C.wrap_glIsSampler(C.uint(sampler))))
}
func IsShader(shader uint32) bool {
	return cbool(uint(C.wrap_glIsShader(C.uint(shader))))
}
func IsSync(sync Sync) bool {
	return cbool(uint(C.wrap_glIsSync(C.GLsync(sync))))
}
func IsTexture(texture uint32) bool {
	return cbool(uint(C.wrap_glIsTexture(C.uint(texture))))
}
func IsTransformFeedback(id uint32) bool {
	return cbool(uint(C.wrap_glIsTransformFeedback(C.uint(id))))
}
func IsVertexArray(array uint32) bool {
	return cbool(uint(C.wrap_glIsVertexArray(C.uint(array))))
}
func UnmapBuffer(target uint32) bool {
	return cbool(uint(C.wrap_glUnmapBuffer(C.uint(target))))
}
func CheckFramebufferStatus(target uint32) uint32 {
	return uint32(C.wrap_glCheckFramebufferStatus(C.uint(target)))
}
func ClientWaitSync(sync Sync, flags uint32, timeout uint64) uint32 {
	return uint32(C.wrap_glClientWaitSync(C.GLsync(sync), C.uint(flags), C.ulonglong(timeout)))
}
func GetError() uint32 {
	return uint32(C.wrap_glGetError())
}
func GetGraphicsResetStatusARB() uint32 {
	return uint32(C.wrap_glGetGraphicsResetStatusARB())
}
func GetAttribLocation(program uint32, name string) int32 {
	cstr1 := C.CString(name)
	defer C.free(unsafe.Pointer(cstr1))
	return int32(C.wrap_glGetAttribLocation(C.uint(program), cstr1))
}
func GetFragDataIndex(program uint32, name string) int32 {
	cstr1 := C.CString(name)
	defer C.free(unsafe.Pointer(cstr1))
	return int32(C.wrap_glGetFragDataIndex(C.uint(program), cstr1))
}
func GetFragDataLocation(program uint32, name string) int32 {
	cstr1 := C.CString(name)
	defer C.free(unsafe.Pointer(cstr1))
	return int32(C.wrap_glGetFragDataLocation(C.uint(program), cstr1))
}
func GetProgramResourceLocation(program uint32, programInterface uint32, name string) int32 {
	cstr1 := C.CString(name)
	defer C.free(unsafe.Pointer(cstr1))
	return int32(C.wrap_glGetProgramResourceLocation(C.uint(program), C.uint(programInterface), cstr1))
}
func GetProgramResourceLocationIndex(program uint32, programInterface uint32, name string) int32 {
	cstr1 := C.CString(name)
	defer C.free(unsafe.Pointer(cstr1))
	return int32(C.wrap_glGetProgramResourceLocationIndex(C.uint(program), C.uint(programInterface), cstr1))
}
func GetSubroutineUniformLocation(program uint32, shadertype uint32, name string) int32 {
	cstr1 := C.CString(name)
	defer C.free(unsafe.Pointer(cstr1))
	return int32(C.wrap_glGetSubroutineUniformLocation(C.uint(program), C.uint(shadertype), cstr1))
}
func GetUniformLocation(program uint32, name string) int32 {
	cstr1 := C.CString(name)
	defer C.free(unsafe.Pointer(cstr1))
	return int32(C.wrap_glGetUniformLocation(C.uint(program), cstr1))
}
func CreateSyncFromCLeventARB(context *clContext, event *clEvent, flags uint32) Sync {
	return Sync(C.wrap_glCreateSyncFromCLeventARB((*C.struct__cl_context)(context), (*C.struct__cl_event)(event), C.uint(flags)))
}
func FenceSync(condition uint32, flags uint32) Sync {
	return Sync(C.wrap_glFenceSync(C.uint(condition), C.uint(flags)))
}
func CreateProgram() uint32 {
	return uint32(C.wrap_glCreateProgram())
}
func CreateShader(t_ype uint32) uint32 {
	return uint32(C.wrap_glCreateShader(C.uint(t_ype)))
}
func CreateShaderProgramv(t_ype uint32, count int32, strings []string) uint32 {
	cstrings := C.newStringArray(C.int(len(strings)))
	defer C.freeStringArray(cstrings, C.int(len(strings)))
	for cnt, str := range strings {
		C.assignString(cstrings, C.CString(str), C.int(cnt))
	}
	return uint32(C.wrap_glCreateShaderProgramv(C.uint(t_ype), C.int(count), cstrings))
}
func GetDebugMessageLog(count uint32, bufsize int32, sources *uint32, types *uint32, ids *uint32, severities *uint32, lengths *int32, messageLog *uint8) uint32 {
	return uint32(C.wrap_glGetDebugMessageLog(C.uint(count), C.int(bufsize), (*C.uint)(sources), (*C.uint)(types), (*C.uint)(ids), (*C.uint)(severities), (*C.int)(lengths), (*C.uchar)(messageLog)))
}
func GetDebugMessageLogARB(count uint32, bufsize int32, sources *uint32, types *uint32, ids *uint32, severities *uint32, lengths *int32, messageLog *uint8) uint32 {
	return uint32(C.wrap_glGetDebugMessageLogARB(C.uint(count), C.int(bufsize), (*C.uint)(sources), (*C.uint)(types), (*C.uint)(ids), (*C.uint)(severities), (*C.int)(lengths), (*C.uchar)(messageLog)))
}
func GetProgramResourceIndex(program uint32, programInterface uint32, name string) uint32 {
	cstr1 := C.CString(name)
	defer C.free(unsafe.Pointer(cstr1))
	return uint32(C.wrap_glGetProgramResourceIndex(C.uint(program), C.uint(programInterface), cstr1))
}
func GetSubroutineIndex(program uint32, shadertype uint32, name string) uint32 {
	cstr1 := C.CString(name)
	defer C.free(unsafe.Pointer(cstr1))
	return uint32(C.wrap_glGetSubroutineIndex(C.uint(program), C.uint(shadertype), cstr1))
}
func GetUniformBlockIndex(program uint32, uniformBlockName string) uint32 {
	cstr1 := C.CString(uniformBlockName)
	defer C.free(unsafe.Pointer(cstr1))
	return uint32(C.wrap_glGetUniformBlockIndex(C.uint(program), cstr1))
}
func MapBuffer(target uint32, access uint32) Pointer {
	return Pointer(C.wrap_glMapBuffer(C.uint(target), C.uint(access)))
}
func MapBufferRange(target uint32, offset int64, length int64, access uint32) Pointer {
	return Pointer(C.wrap_glMapBufferRange(C.uint(target), C.longlong(offset), C.longlong(length), C.uint(access)))
}
func GetString(name uint32) string {
	return C.GoString(C.wrap_glGetString(C.uint(name)))
}
func GetStringi(name uint32, index uint32) string {
	return C.GoString(C.wrap_glGetStringi(C.uint(name), C.uint(index)))
}
func ActiveShaderProgram(pipeline uint32, program uint32) {
	C.wrap_glActiveShaderProgram(C.uint(pipeline), C.uint(program))
}
func ActiveTexture(texture uint32) {
	C.wrap_glActiveTexture(C.uint(texture))
}
func AttachShader(program uint32, shader uint32) {
	C.wrap_glAttachShader(C.uint(program), C.uint(shader))
}
func BeginConditionalRender(id uint32, mode uint32) {
	C.wrap_glBeginConditionalRender(C.uint(id), C.uint(mode))
}
func BeginQuery(target uint32, id uint32) {
	C.wrap_glBeginQuery(C.uint(target), C.uint(id))
}
func BeginQueryIndexed(target uint32, index uint32, id uint32) {
	C.wrap_glBeginQueryIndexed(C.uint(target), C.uint(index), C.uint(id))
}
func BeginTransformFeedback(primitiveMode uint32) {
	C.wrap_glBeginTransformFeedback(C.uint(primitiveMode))
}
func BindAttribLocation(program uint32, index uint32, name string) {
	cstr1 := C.CString(name)
	defer C.free(unsafe.Pointer(cstr1))
	C.wrap_glBindAttribLocation(C.uint(program), C.uint(index), cstr1)
}
func BindBuffer(target uint32, buffer uint32) {
	C.wrap_glBindBuffer(C.uint(target), C.uint(buffer))
}
func BindBufferBase(target uint32, index uint32, buffer uint32) {
	C.wrap_glBindBufferBase(C.uint(target), C.uint(index), C.uint(buffer))
}
func BindBufferRange(target uint32, index uint32, buffer uint32, offset int64, size int64) {
	C.wrap_glBindBufferRange(C.uint(target), C.uint(index), C.uint(buffer), C.longlong(offset), C.longlong(size))
}
func BindFragDataLocation(program uint32, color uint32, name string) {
	cstr1 := C.CString(name)
	defer C.free(unsafe.Pointer(cstr1))
	C.wrap_glBindFragDataLocation(C.uint(program), C.uint(color), cstr1)
}
func BindFragDataLocationIndexed(program uint32, colorNumber uint32, index uint32, name string) {
	cstr1 := C.CString(name)
	defer C.free(unsafe.Pointer(cstr1))
	C.wrap_glBindFragDataLocationIndexed(C.uint(program), C.uint(colorNumber), C.uint(index), cstr1)
}
func BindFramebuffer(target uint32, framebuffer uint32) {
	C.wrap_glBindFramebuffer(C.uint(target), C.uint(framebuffer))
}
func BindImageTexture(unit uint32, texture uint32, level int32, layered bool, layer int32, access uint32, format uint32) {
	tf1 := FALSE
	if layered {
		tf1 = TRUE
	}
	C.wrap_glBindImageTexture(C.uint(unit), C.uint(texture), C.int(level), C.uchar(tf1), C.int(layer), C.uint(access), C.uint(format))
}
func BindProgramPipeline(pipeline uint32) {
	C.wrap_glBindProgramPipeline(C.uint(pipeline))
}
func BindRenderbuffer(target uint32, renderbuffer uint32) {
	C.wrap_glBindRenderbuffer(C.uint(target), C.uint(renderbuffer))
}
func BindSampler(unit uint32, sampler uint32) {
	C.wrap_glBindSampler(C.uint(unit), C.uint(sampler))
}
func BindTexture(target uint32, texture uint32) {
	C.wrap_glBindTexture(C.uint(target), C.uint(texture))
}
func BindTransformFeedback(target uint32, id uint32) {
	C.wrap_glBindTransformFeedback(C.uint(target), C.uint(id))
}
func BindVertexArray(array uint32) {
	C.wrap_glBindVertexArray(C.uint(array))
}
func BindVertexBuffer(bindingindex uint32, buffer uint32, offset int64, stride int32) {
	C.wrap_glBindVertexBuffer(C.uint(bindingindex), C.uint(buffer), C.longlong(offset), C.int(stride))
}
func BlendColor(red float32, green float32, blue float32, alpha float32) {
	C.wrap_glBlendColor(C.float(red), C.float(green), C.float(blue), C.float(alpha))
}
func BlendEquation(mode uint32) {
	C.wrap_glBlendEquation(C.uint(mode))
}
func BlendEquationSeparate(modeRGB uint32, modeAlpha uint32) {
	C.wrap_glBlendEquationSeparate(C.uint(modeRGB), C.uint(modeAlpha))
}
func BlendEquationSeparatei(buf uint32, modeRGB uint32, modeAlpha uint32) {
	C.wrap_glBlendEquationSeparatei(C.uint(buf), C.uint(modeRGB), C.uint(modeAlpha))
}
func BlendEquationSeparateiARB(buf uint32, modeRGB uint32, modeAlpha uint32) {
	C.wrap_glBlendEquationSeparateiARB(C.uint(buf), C.uint(modeRGB), C.uint(modeAlpha))
}
func BlendEquationi(buf uint32, mode uint32) {
	C.wrap_glBlendEquationi(C.uint(buf), C.uint(mode))
}
func BlendEquationiARB(buf uint32, mode uint32) {
	C.wrap_glBlendEquationiARB(C.uint(buf), C.uint(mode))
}
func BlendFunc(sfactor uint32, dfactor uint32) {
	C.wrap_glBlendFunc(C.uint(sfactor), C.uint(dfactor))
}
func BlendFuncSeparate(sfactorRGB uint32, dfactorRGB uint32, sfactorAlpha uint32, dfactorAlpha uint32) {
	C.wrap_glBlendFuncSeparate(C.uint(sfactorRGB), C.uint(dfactorRGB), C.uint(sfactorAlpha), C.uint(dfactorAlpha))
}
func BlendFuncSeparatei(buf uint32, srcRGB uint32, dstRGB uint32, srcAlpha uint32, dstAlpha uint32) {
	C.wrap_glBlendFuncSeparatei(C.uint(buf), C.uint(srcRGB), C.uint(dstRGB), C.uint(srcAlpha), C.uint(dstAlpha))
}
func BlendFuncSeparateiARB(buf uint32, srcRGB uint32, dstRGB uint32, srcAlpha uint32, dstAlpha uint32) {
	C.wrap_glBlendFuncSeparateiARB(C.uint(buf), C.uint(srcRGB), C.uint(dstRGB), C.uint(srcAlpha), C.uint(dstAlpha))
}
func BlendFunci(buf uint32, src uint32, dst uint32) {
	C.wrap_glBlendFunci(C.uint(buf), C.uint(src), C.uint(dst))
}
func BlendFunciARB(buf uint32, src uint32, dst uint32) {
	C.wrap_glBlendFunciARB(C.uint(buf), C.uint(src), C.uint(dst))
}
func BlitFramebuffer(srcX0 int32, srcY0 int32, srcX1 int32, srcY1 int32, dstX0 int32, dstY0 int32, dstX1 int32, dstY1 int32, mask uint32, filter uint32) {
	C.wrap_glBlitFramebuffer(C.int(srcX0), C.int(srcY0), C.int(srcX1), C.int(srcY1), C.int(dstX0), C.int(dstY0), C.int(dstX1), C.int(dstY1), C.uint(mask), C.uint(filter))
}
func BufferData(target uint32, size int64, data Pointer, usage uint32) {
	C.wrap_glBufferData(C.uint(target), C.longlong(size), unsafe.Pointer(data), C.uint(usage))
}
func BufferSubData(target uint32, offset int64, size int64, data Pointer) {
	C.wrap_glBufferSubData(C.uint(target), C.longlong(offset), C.longlong(size), unsafe.Pointer(data))
}
func ClampColor(target uint32, clamp uint32) {
	C.wrap_glClampColor(C.uint(target), C.uint(clamp))
}
func Clear(mask uint32) {
	C.wrap_glClear(C.uint(mask))
}
func ClearBufferData(target uint32, internalformat uint32, format uint32, t_ype uint32, data Pointer) {
	C.wrap_glClearBufferData(C.uint(target), C.uint(internalformat), C.uint(format), C.uint(t_ype), unsafe.Pointer(data))
}
func ClearBufferSubData(target uint32, internalformat uint32, offset int64, size int64, format uint32, t_ype uint32, data Pointer) {
	C.wrap_glClearBufferSubData(C.uint(target), C.uint(internalformat), C.longlong(offset), C.longlong(size), C.uint(format), C.uint(t_ype), unsafe.Pointer(data))
}
func ClearBufferfi(buffer uint32, drawbuffer int32, depth float32, stencil int32) {
	C.wrap_glClearBufferfi(C.uint(buffer), C.int(drawbuffer), C.float(depth), C.int(stencil))
}
func ClearBufferfv(buffer uint32, drawbuffer int32, value *float32) {
	C.wrap_glClearBufferfv(C.uint(buffer), C.int(drawbuffer), (*C.float)(value))
}
func ClearBufferiv(buffer uint32, drawbuffer int32, value *int32) {
	C.wrap_glClearBufferiv(C.uint(buffer), C.int(drawbuffer), (*C.int)(value))
}
func ClearBufferuiv(buffer uint32, drawbuffer int32, value *uint32) {
	C.wrap_glClearBufferuiv(C.uint(buffer), C.int(drawbuffer), (*C.uint)(value))
}
func ClearColor(red float32, green float32, blue float32, alpha float32) {
	C.wrap_glClearColor(C.float(red), C.float(green), C.float(blue), C.float(alpha))
}
func ClearDepth(depth float64) {
	C.wrap_glClearDepth(C.double(depth))
}
func ClearDepthf(d float32) {
	C.wrap_glClearDepthf(C.float(d))
}
func ClearNamedBufferDataEXT(buffer uint32, internalformat uint32, format uint32, t_ype uint32, data Pointer) {
	C.wrap_glClearNamedBufferDataEXT(C.uint(buffer), C.uint(internalformat), C.uint(format), C.uint(t_ype), unsafe.Pointer(data))
}
func ClearNamedBufferSubDataEXT(buffer uint32, internalformat uint32, format uint32, t_ype uint32, offset int64, size int64, data Pointer) {
	C.wrap_glClearNamedBufferSubDataEXT(C.uint(buffer), C.uint(internalformat), C.uint(format), C.uint(t_ype), C.longlong(offset), C.longlong(size), unsafe.Pointer(data))
}
func ClearStencil(s int32) {
	C.wrap_glClearStencil(C.int(s))
}
func ColorMask(red bool, green bool, blue bool, alpha bool) {
	tf4 := FALSE
	if alpha {
		tf4 = TRUE
	}
	tf3 := FALSE
	if blue {
		tf3 = TRUE
	}
	tf2 := FALSE
	if green {
		tf2 = TRUE
	}
	tf1 := FALSE
	if red {
		tf1 = TRUE
	}
	C.wrap_glColorMask(C.uchar(tf1), C.uchar(tf2), C.uchar(tf3), C.uchar(tf4))
}
func ColorMaski(index uint32, r bool, g bool, b bool, a bool) {
	tf4 := FALSE
	if a {
		tf4 = TRUE
	}
	tf3 := FALSE
	if b {
		tf3 = TRUE
	}
	tf2 := FALSE
	if g {
		tf2 = TRUE
	}
	tf1 := FALSE
	if r {
		tf1 = TRUE
	}
	C.wrap_glColorMaski(C.uint(index), C.uchar(tf1), C.uchar(tf2), C.uchar(tf3), C.uchar(tf4))
}
func ColorP3ui(t_ype uint32, color uint32) {
	C.wrap_glColorP3ui(C.uint(t_ype), C.uint(color))
}
func ColorP3uiv(t_ype uint32, color *uint32) {
	C.wrap_glColorP3uiv(C.uint(t_ype), (*C.uint)(color))
}
func ColorP4ui(t_ype uint32, color uint32) {
	C.wrap_glColorP4ui(C.uint(t_ype), C.uint(color))
}
func ColorP4uiv(t_ype uint32, color *uint32) {
	C.wrap_glColorP4uiv(C.uint(t_ype), (*C.uint)(color))
}
func CompileShader(shader uint32) {
	C.wrap_glCompileShader(C.uint(shader))
}
func CompileShaderIncludeARB(shader uint32, count int32, path []string, length *int32) {
	cstrings := C.newStringArray(C.int(len(path)))
	defer C.freeStringArray(cstrings, C.int(len(path)))
	for cnt, str := range path {
		C.assignString(cstrings, C.CString(str), C.int(cnt))
	}
	C.wrap_glCompileShaderIncludeARB(C.uint(shader), C.int(count), cstrings, (*C.int)(length))
}
func CompressedTexImage1D(target uint32, level int32, internalformat uint32, width int32, border int32, imageSize int32, data Pointer) {
	C.wrap_glCompressedTexImage1D(C.uint(target), C.int(level), C.uint(internalformat), C.int(width), C.int(border), C.int(imageSize), unsafe.Pointer(data))
}
func CompressedTexImage2D(target uint32, level int32, internalformat uint32, width int32, height int32, border int32, imageSize int32, data Pointer) {
	C.wrap_glCompressedTexImage2D(C.uint(target), C.int(level), C.uint(internalformat), C.int(width), C.int(height), C.int(border), C.int(imageSize), unsafe.Pointer(data))
}
func CompressedTexImage3D(target uint32, level int32, internalformat uint32, width int32, height int32, depth int32, border int32, imageSize int32, data Pointer) {
	C.wrap_glCompressedTexImage3D(C.uint(target), C.int(level), C.uint(internalformat), C.int(width), C.int(height), C.int(depth), C.int(border), C.int(imageSize), unsafe.Pointer(data))
}
func CompressedTexSubImage1D(target uint32, level int32, xoffset int32, width int32, format uint32, imageSize int32, data Pointer) {
	C.wrap_glCompressedTexSubImage1D(C.uint(target), C.int(level), C.int(xoffset), C.int(width), C.uint(format), C.int(imageSize), unsafe.Pointer(data))
}
func CompressedTexSubImage2D(target uint32, level int32, xoffset int32, yoffset int32, width int32, height int32, format uint32, imageSize int32, data Pointer) {
	C.wrap_glCompressedTexSubImage2D(C.uint(target), C.int(level), C.int(xoffset), C.int(yoffset), C.int(width), C.int(height), C.uint(format), C.int(imageSize), unsafe.Pointer(data))
}
func CompressedTexSubImage3D(target uint32, level int32, xoffset int32, yoffset int32, zoffset int32, width int32, height int32, depth int32, format uint32, imageSize int32, data Pointer) {
	C.wrap_glCompressedTexSubImage3D(C.uint(target), C.int(level), C.int(xoffset), C.int(yoffset), C.int(zoffset), C.int(width), C.int(height), C.int(depth), C.uint(format), C.int(imageSize), unsafe.Pointer(data))
}
func CopyBufferSubData(readTarget uint32, writeTarget uint32, readOffset int64, writeOffset int64, size int64) {
	C.wrap_glCopyBufferSubData(C.uint(readTarget), C.uint(writeTarget), C.longlong(readOffset), C.longlong(writeOffset), C.longlong(size))
}
func CopyImageSubData(srcName uint32, srcTarget uint32, srcLevel int32, srcX int32, srcY int32, srcZ int32, dstName uint32, dstTarget uint32, dstLevel int32, dstX int32, dstY int32, dstZ int32, srcWidth int32, srcHeight int32, srcDepth int32) {
	C.wrap_glCopyImageSubData(C.uint(srcName), C.uint(srcTarget), C.int(srcLevel), C.int(srcX), C.int(srcY), C.int(srcZ), C.uint(dstName), C.uint(dstTarget), C.int(dstLevel), C.int(dstX), C.int(dstY), C.int(dstZ), C.int(srcWidth), C.int(srcHeight), C.int(srcDepth))
}
func CopyTexImage1D(target uint32, level int32, internalformat uint32, x int32, y int32, width int32, border int32) {
	C.wrap_glCopyTexImage1D(C.uint(target), C.int(level), C.uint(internalformat), C.int(x), C.int(y), C.int(width), C.int(border))
}
func CopyTexImage2D(target uint32, level int32, internalformat uint32, x int32, y int32, width int32, height int32, border int32) {
	C.wrap_glCopyTexImage2D(C.uint(target), C.int(level), C.uint(internalformat), C.int(x), C.int(y), C.int(width), C.int(height), C.int(border))
}
func CopyTexSubImage1D(target uint32, level int32, xoffset int32, x int32, y int32, width int32) {
	C.wrap_glCopyTexSubImage1D(C.uint(target), C.int(level), C.int(xoffset), C.int(x), C.int(y), C.int(width))
}
func CopyTexSubImage2D(target uint32, level int32, xoffset int32, yoffset int32, x int32, y int32, width int32, height int32) {
	C.wrap_glCopyTexSubImage2D(C.uint(target), C.int(level), C.int(xoffset), C.int(yoffset), C.int(x), C.int(y), C.int(width), C.int(height))
}
func CopyTexSubImage3D(target uint32, level int32, xoffset int32, yoffset int32, zoffset int32, x int32, y int32, width int32, height int32) {
	C.wrap_glCopyTexSubImage3D(C.uint(target), C.int(level), C.int(xoffset), C.int(yoffset), C.int(zoffset), C.int(x), C.int(y), C.int(width), C.int(height))
}
func CullFace(mode uint32) {
	C.wrap_glCullFace(C.uint(mode))
}
func DebugMessageCallback(callback DEBUGPROC, userParam Pointer) {
	C.wrap_glDebugMessageCallback(C.GLDEBUGPROC(callback), unsafe.Pointer(userParam))
}
func DebugMessageCallbackARB(callback DEBUGPROCARB, userParam Pointer) {
	C.wrap_glDebugMessageCallbackARB(C.GLDEBUGPROCARB(callback), unsafe.Pointer(userParam))
}
func DebugMessageControl(source uint32, t_ype uint32, severity uint32, count int32, ids *uint32, enabled bool) {
	tf1 := FALSE
	if enabled {
		tf1 = TRUE
	}
	C.wrap_glDebugMessageControl(C.uint(source), C.uint(t_ype), C.uint(severity), C.int(count), (*C.uint)(ids), C.uchar(tf1))
}
func DebugMessageControlARB(source uint32, t_ype uint32, severity uint32, count int32, ids *uint32, enabled bool) {
	tf1 := FALSE
	if enabled {
		tf1 = TRUE
	}
	C.wrap_glDebugMessageControlARB(C.uint(source), C.uint(t_ype), C.uint(severity), C.int(count), (*C.uint)(ids), C.uchar(tf1))
}
func DebugMessageInsert(source uint32, t_ype uint32, id uint32, severity uint32, length int32, buf string) {
	cstr1 := C.CString(buf)
	defer C.free(unsafe.Pointer(cstr1))
	C.wrap_glDebugMessageInsert(C.uint(source), C.uint(t_ype), C.uint(id), C.uint(severity), C.int(length), cstr1)
}
func DebugMessageInsertARB(source uint32, t_ype uint32, id uint32, severity uint32, length int32, buf string) {
	cstr1 := C.CString(buf)
	defer C.free(unsafe.Pointer(cstr1))
	C.wrap_glDebugMessageInsertARB(C.uint(source), C.uint(t_ype), C.uint(id), C.uint(severity), C.int(length), cstr1)
}
func DeleteBuffers(n int32, buffers *uint32) {
	C.wrap_glDeleteBuffers(C.int(n), (*C.uint)(buffers))
}
func DeleteFramebuffers(n int32, framebuffers *uint32) {
	C.wrap_glDeleteFramebuffers(C.int(n), (*C.uint)(framebuffers))
}
func DeleteNamedStringARB(namelen int32, name string) {
	cstr1 := C.CString(name)
	defer C.free(unsafe.Pointer(cstr1))
	C.wrap_glDeleteNamedStringARB(C.int(namelen), cstr1)
}
func DeleteProgram(program uint32) {
	C.wrap_glDeleteProgram(C.uint(program))
}
func DeleteProgramPipelines(n int32, pipelines *uint32) {
	C.wrap_glDeleteProgramPipelines(C.int(n), (*C.uint)(pipelines))
}
func DeleteQueries(n int32, ids *uint32) {
	C.wrap_glDeleteQueries(C.int(n), (*C.uint)(ids))
}
func DeleteRenderbuffers(n int32, renderbuffers *uint32) {
	C.wrap_glDeleteRenderbuffers(C.int(n), (*C.uint)(renderbuffers))
}
func DeleteSamplers(count int32, samplers *uint32) {
	C.wrap_glDeleteSamplers(C.int(count), (*C.uint)(samplers))
}
func DeleteShader(shader uint32) {
	C.wrap_glDeleteShader(C.uint(shader))
}
func DeleteSync(sync Sync) {
	C.wrap_glDeleteSync(C.GLsync(sync))
}
func DeleteTextures(n int32, textures *uint32) {
	C.wrap_glDeleteTextures(C.int(n), (*C.uint)(textures))
}
func DeleteTransformFeedbacks(n int32, ids *uint32) {
	C.wrap_glDeleteTransformFeedbacks(C.int(n), (*C.uint)(ids))
}
func DeleteVertexArrays(n int32, arrays *uint32) {
	C.wrap_glDeleteVertexArrays(C.int(n), (*C.uint)(arrays))
}
func DepthFunc(f_unc uint32) {
	C.wrap_glDepthFunc(C.uint(f_unc))
}
func DepthMask(flag bool) {
	tf1 := FALSE
	if flag {
		tf1 = TRUE
	}
	C.wrap_glDepthMask(C.uchar(tf1))
}
func DepthRange(n_ear float64, f_ar float64) {
	C.wrap_glDepthRange(C.double(n_ear), C.double(f_ar))
}
func DepthRangeArrayv(first uint32, count int32, v *float64) {
	C.wrap_glDepthRangeArrayv(C.uint(first), C.int(count), (*C.double)(v))
}
func DepthRangeIndexed(index uint32, n float64, f float64) {
	C.wrap_glDepthRangeIndexed(C.uint(index), C.double(n), C.double(f))
}
func DepthRangef(n float32, f float32) {
	C.wrap_glDepthRangef(C.float(n), C.float(f))
}
func DetachShader(program uint32, shader uint32) {
	C.wrap_glDetachShader(C.uint(program), C.uint(shader))
}
func Disable(cap uint32) {
	C.wrap_glDisable(C.uint(cap))
}
func DisableVertexAttribArray(index uint32) {
	C.wrap_glDisableVertexAttribArray(C.uint(index))
}
func Disablei(target uint32, index uint32) {
	C.wrap_glDisablei(C.uint(target), C.uint(index))
}
func DispatchCompute(num_groups_x uint32, num_groups_y uint32, num_groups_z uint32) {
	C.wrap_glDispatchCompute(C.uint(num_groups_x), C.uint(num_groups_y), C.uint(num_groups_z))
}
func DispatchComputeIndirect(indirect int64) {
	C.wrap_glDispatchComputeIndirect(C.longlong(indirect))
}
func DrawArrays(mode uint32, first int32, count int32) {
	C.wrap_glDrawArrays(C.uint(mode), C.int(first), C.int(count))
}
func DrawArraysIndirect(mode uint32, indirect Pointer) {
	C.wrap_glDrawArraysIndirect(C.uint(mode), unsafe.Pointer(indirect))
}
func DrawArraysInstanced(mode uint32, first int32, count int32, instancecount int32) {
	C.wrap_glDrawArraysInstanced(C.uint(mode), C.int(first), C.int(count), C.int(instancecount))
}
func DrawArraysInstancedBaseInstance(mode uint32, first int32, count int32, instancecount int32, baseinstance uint32) {
	C.wrap_glDrawArraysInstancedBaseInstance(C.uint(mode), C.int(first), C.int(count), C.int(instancecount), C.uint(baseinstance))
}
func DrawBuffer(mode uint32) {
	C.wrap_glDrawBuffer(C.uint(mode))
}
func DrawBuffers(n int32, bufs *uint32) {
	C.wrap_glDrawBuffers(C.int(n), (*C.uint)(bufs))
}
func DrawElements(mode uint32, count int32, t_ype uint32, indicies int64) {
	C.wrap_glDrawElements(C.uint(mode), C.int(count), C.uint(t_ype), C.longlong(indicies))
}
func DrawElementsBaseVertex(mode uint32, count int32, t_ype uint32, indicies int64, basevertex int32) {
	C.wrap_glDrawElementsBaseVertex(C.uint(mode), C.int(count), C.uint(t_ype), C.longlong(indicies), C.int(basevertex))
}
func DrawElementsIndirect(mode uint32, t_ype uint32, indirect Pointer) {
	C.wrap_glDrawElementsIndirect(C.uint(mode), C.uint(t_ype), unsafe.Pointer(indirect))
}
func DrawElementsInstanced(mode uint32, count int32, t_ype uint32, indicies int64, instancecount int32) {
	C.wrap_glDrawElementsInstanced(C.uint(mode), C.int(count), C.uint(t_ype), C.longlong(indicies), C.int(instancecount))
}
func DrawElementsInstancedBaseInstance(mode uint32, count int32, t_ype uint32, indices Pointer, instancecount int32, baseinstance uint32) {
	C.wrap_glDrawElementsInstancedBaseInstance(C.uint(mode), C.int(count), C.uint(t_ype), unsafe.Pointer(indices), C.int(instancecount), C.uint(baseinstance))
}
func DrawElementsInstancedBaseVertex(mode uint32, count int32, t_ype uint32, indicies int64, instancecount int32, basevertex int32) {
	C.wrap_glDrawElementsInstancedBaseVertex(C.uint(mode), C.int(count), C.uint(t_ype), C.longlong(indicies), C.int(instancecount), C.int(basevertex))
}
func DrawElementsInstancedBaseVertexBaseInstance(mode uint32, count int32, t_ype uint32, indices Pointer, instancecount int32, basevertex int32, baseinstance uint32) {
	C.wrap_glDrawElementsInstancedBaseVertexBaseInstance(C.uint(mode), C.int(count), C.uint(t_ype), unsafe.Pointer(indices), C.int(instancecount), C.int(basevertex), C.uint(baseinstance))
}
func DrawRangeElements(mode uint32, start uint32, end uint32, count int32, t_ype uint32, indicies int64) {
	C.wrap_glDrawRangeElements(C.uint(mode), C.uint(start), C.uint(end), C.int(count), C.uint(t_ype), C.longlong(indicies))
}
func DrawRangeElementsBaseVertex(mode uint32, start uint32, end uint32, count int32, t_ype uint32, indicies int64, basevertex int32) {
	C.wrap_glDrawRangeElementsBaseVertex(C.uint(mode), C.uint(start), C.uint(end), C.int(count), C.uint(t_ype), C.longlong(indicies), C.int(basevertex))
}
func DrawTransformFeedback(mode uint32, id uint32) {
	C.wrap_glDrawTransformFeedback(C.uint(mode), C.uint(id))
}
func DrawTransformFeedbackInstanced(mode uint32, id uint32, instancecount int32) {
	C.wrap_glDrawTransformFeedbackInstanced(C.uint(mode), C.uint(id), C.int(instancecount))
}
func DrawTransformFeedbackStream(mode uint32, id uint32, stream uint32) {
	C.wrap_glDrawTransformFeedbackStream(C.uint(mode), C.uint(id), C.uint(stream))
}
func DrawTransformFeedbackStreamInstanced(mode uint32, id uint32, stream uint32, instancecount int32) {
	C.wrap_glDrawTransformFeedbackStreamInstanced(C.uint(mode), C.uint(id), C.uint(stream), C.int(instancecount))
}
func Enable(cap uint32) {
	C.wrap_glEnable(C.uint(cap))
}
func EnableVertexAttribArray(index uint32) {
	C.wrap_glEnableVertexAttribArray(C.uint(index))
}
func Enablei(target uint32, index uint32) {
	C.wrap_glEnablei(C.uint(target), C.uint(index))
}
func EndConditionalRender() {
	C.wrap_glEndConditionalRender()
}
func EndQuery(target uint32) {
	C.wrap_glEndQuery(C.uint(target))
}
func EndQueryIndexed(target uint32, index uint32) {
	C.wrap_glEndQueryIndexed(C.uint(target), C.uint(index))
}
func EndTransformFeedback() {
	C.wrap_glEndTransformFeedback()
}
func Finish() {
	C.wrap_glFinish()
}
func Flush() {
	C.wrap_glFlush()
}
func FlushMappedBufferRange(target uint32, offset int64, length int64) {
	C.wrap_glFlushMappedBufferRange(C.uint(target), C.longlong(offset), C.longlong(length))
}
func FramebufferParameteri(target uint32, pname uint32, param int32) {
	C.wrap_glFramebufferParameteri(C.uint(target), C.uint(pname), C.int(param))
}
func FramebufferRenderbuffer(target uint32, attachment uint32, renderbuffertarget uint32, renderbuffer uint32) {
	C.wrap_glFramebufferRenderbuffer(C.uint(target), C.uint(attachment), C.uint(renderbuffertarget), C.uint(renderbuffer))
}
func FramebufferTexture(target uint32, attachment uint32, texture uint32, level int32) {
	C.wrap_glFramebufferTexture(C.uint(target), C.uint(attachment), C.uint(texture), C.int(level))
}
func FramebufferTexture1D(target uint32, attachment uint32, textarget uint32, texture uint32, level int32) {
	C.wrap_glFramebufferTexture1D(C.uint(target), C.uint(attachment), C.uint(textarget), C.uint(texture), C.int(level))
}
func FramebufferTexture2D(target uint32, attachment uint32, textarget uint32, texture uint32, level int32) {
	C.wrap_glFramebufferTexture2D(C.uint(target), C.uint(attachment), C.uint(textarget), C.uint(texture), C.int(level))
}
func FramebufferTexture3D(target uint32, attachment uint32, textarget uint32, texture uint32, level int32, zoffset int32) {
	C.wrap_glFramebufferTexture3D(C.uint(target), C.uint(attachment), C.uint(textarget), C.uint(texture), C.int(level), C.int(zoffset))
}
func FramebufferTextureLayer(target uint32, attachment uint32, texture uint32, level int32, layer int32) {
	C.wrap_glFramebufferTextureLayer(C.uint(target), C.uint(attachment), C.uint(texture), C.int(level), C.int(layer))
}
func FrontFace(mode uint32) {
	C.wrap_glFrontFace(C.uint(mode))
}
func GenBuffers(n int32, buffers *uint32) {
	C.wrap_glGenBuffers(C.int(n), (*C.uint)(buffers))
}
func GenFramebuffers(n int32, framebuffers *uint32) {
	C.wrap_glGenFramebuffers(C.int(n), (*C.uint)(framebuffers))
}
func GenProgramPipelines(n int32, pipelines *uint32) {
	C.wrap_glGenProgramPipelines(C.int(n), (*C.uint)(pipelines))
}
func GenQueries(n int32, ids *uint32) {
	C.wrap_glGenQueries(C.int(n), (*C.uint)(ids))
}
func GenRenderbuffers(n int32, renderbuffers *uint32) {
	C.wrap_glGenRenderbuffers(C.int(n), (*C.uint)(renderbuffers))
}
func GenSamplers(count int32, samplers *uint32) {
	C.wrap_glGenSamplers(C.int(count), (*C.uint)(samplers))
}
func GenTextures(n int32, textures *uint32) {
	C.wrap_glGenTextures(C.int(n), (*C.uint)(textures))
}
func GenTransformFeedbacks(n int32, ids *uint32) {
	C.wrap_glGenTransformFeedbacks(C.int(n), (*C.uint)(ids))
}
func GenVertexArrays(n int32, arrays *uint32) {
	C.wrap_glGenVertexArrays(C.int(n), (*C.uint)(arrays))
}
func GenerateMipmap(target uint32) {
	C.wrap_glGenerateMipmap(C.uint(target))
}
func GetActiveAtomicCounterBufferiv(program uint32, bufferIndex uint32, pname uint32, params *int32) {
	C.wrap_glGetActiveAtomicCounterBufferiv(C.uint(program), C.uint(bufferIndex), C.uint(pname), (*C.int)(params))
}
func GetActiveAttrib(program uint32, index uint32, bufSize int32, length *int32, size *int32, t_ype *uint32, name string) {
	cstr1 := C.CString(name)
	defer C.free(unsafe.Pointer(cstr1))
	C.wrap_glGetActiveAttrib(C.uint(program), C.uint(index), C.int(bufSize), (*C.int)(length), (*C.int)(size), (*C.uint)(t_ype), cstr1)
}
func GetActiveSubroutineName(program uint32, shadertype uint32, index uint32, bufsize int32, length *int32, name *uint8) {
	C.wrap_glGetActiveSubroutineName(C.uint(program), C.uint(shadertype), C.uint(index), C.int(bufsize), (*C.int)(length), (*C.uchar)(name))
}
func GetActiveSubroutineUniformName(program uint32, shadertype uint32, index uint32, bufsize int32, length *int32, name string) {
	cstr1 := C.CString(name)
	defer C.free(unsafe.Pointer(cstr1))
	C.wrap_glGetActiveSubroutineUniformName(C.uint(program), C.uint(shadertype), C.uint(index), C.int(bufsize), (*C.int)(length), cstr1)
}
func GetActiveSubroutineUniformiv(program uint32, shadertype uint32, index uint32, pname uint32, values *int32) {
	C.wrap_glGetActiveSubroutineUniformiv(C.uint(program), C.uint(shadertype), C.uint(index), C.uint(pname), (*C.int)(values))
}
func GetActiveUniform(program uint32, index uint32, bufSize int32, length *int32, size *int32, t_ype *uint32, name *uint8) {
	C.wrap_glGetActiveUniform(C.uint(program), C.uint(index), C.int(bufSize), (*C.int)(length), (*C.int)(size), (*C.uint)(t_ype), (*C.uchar)(name))
}
func GetActiveUniformBlockName(program uint32, uniformBlockIndex uint32, bufSize int32, length *int32, uniformBlockName *uint8) {
	C.wrap_glGetActiveUniformBlockName(C.uint(program), C.uint(uniformBlockIndex), C.int(bufSize), (*C.int)(length), (*C.uchar)(uniformBlockName))
}
func GetActiveUniformBlockiv(program uint32, uniformBlockIndex uint32, pname uint32, params *int32) {
	C.wrap_glGetActiveUniformBlockiv(C.uint(program), C.uint(uniformBlockIndex), C.uint(pname), (*C.int)(params))
}
func GetActiveUniformName(program uint32, uniformIndex uint32, bufSize int32, length *int32, uniformName *uint8) {
	C.wrap_glGetActiveUniformName(C.uint(program), C.uint(uniformIndex), C.int(bufSize), (*C.int)(length), (*C.uchar)(uniformName))
}
func GetActiveUniformsiv(program uint32, uniformCount int32, uniformIndices *uint32, pname uint32, params *int32) {
	C.wrap_glGetActiveUniformsiv(C.uint(program), C.int(uniformCount), (*C.uint)(uniformIndices), C.uint(pname), (*C.int)(params))
}
func GetAttachedShaders(program uint32, maxCount int32, count *int32, obj *uint32) {
	C.wrap_glGetAttachedShaders(C.uint(program), C.int(maxCount), (*C.int)(count), (*C.uint)(obj))
}
func GetBooleani_v(target uint32, index uint32, data *uint8) {
	C.wrap_glGetBooleani_v(C.uint(target), C.uint(index), (*C.uchar)(data))
}
func GetBooleanv(pname uint32, params *uint8) {
	C.wrap_glGetBooleanv(C.uint(pname), (*C.uchar)(params))
}
func GetBufferParameteri64v(target uint32, pname uint32, params *int64) {
	C.wrap_glGetBufferParameteri64v(C.uint(target), C.uint(pname), (*C.longlong)(params))
}
func GetBufferParameteriv(target uint32, pname uint32, params *int32) {
	C.wrap_glGetBufferParameteriv(C.uint(target), C.uint(pname), (*C.int)(params))
}
func GetBufferPointerv(target uint32, pname uint32, params *Pointer) {
	C.wrap_glGetBufferPointerv(C.uint(target), C.uint(pname), (*unsafe.Pointer)(params))
}
func GetBufferSubData(target uint32, offset int64, size int64, data Pointer) {
	C.wrap_glGetBufferSubData(C.uint(target), C.longlong(offset), C.longlong(size), unsafe.Pointer(data))
}
func GetCompressedTexImage(target uint32, level int32, img Pointer) {
	C.wrap_glGetCompressedTexImage(C.uint(target), C.int(level), unsafe.Pointer(img))
}
func GetDoublei_v(target uint32, index uint32, data *float64) {
	C.wrap_glGetDoublei_v(C.uint(target), C.uint(index), (*C.double)(data))
}
func GetDoublev(pname uint32, params *float64) {
	C.wrap_glGetDoublev(C.uint(pname), (*C.double)(params))
}
func GetFloati_v(target uint32, index uint32, data *float32) {
	C.wrap_glGetFloati_v(C.uint(target), C.uint(index), (*C.float)(data))
}
func GetFloatv(pname uint32, params *float32) {
	C.wrap_glGetFloatv(C.uint(pname), (*C.float)(params))
}
func GetFramebufferAttachmentParameteriv(target uint32, attachment uint32, pname uint32, params *int32) {
	C.wrap_glGetFramebufferAttachmentParameteriv(C.uint(target), C.uint(attachment), C.uint(pname), (*C.int)(params))
}
func GetFramebufferParameteriv(target uint32, pname uint32, params *int32) {
	C.wrap_glGetFramebufferParameteriv(C.uint(target), C.uint(pname), (*C.int)(params))
}
func GetInteger64i_v(target uint32, index uint32, data *int64) {
	C.wrap_glGetInteger64i_v(C.uint(target), C.uint(index), (*C.longlong)(data))
}
func GetInteger64v(pname uint32, params *int64) {
	C.wrap_glGetInteger64v(C.uint(pname), (*C.longlong)(params))
}
func GetIntegeri_v(target uint32, index uint32, data *int32) {
	C.wrap_glGetIntegeri_v(C.uint(target), C.uint(index), (*C.int)(data))
}
func GetIntegerv(pname uint32, params *int32) {
	C.wrap_glGetIntegerv(C.uint(pname), (*C.int)(params))
}
func GetInternalformati64v(target uint32, internalformat uint32, pname uint32, bufSize int32, params *int64) {
	C.wrap_glGetInternalformati64v(C.uint(target), C.uint(internalformat), C.uint(pname), C.int(bufSize), (*C.longlong)(params))
}
func GetInternalformativ(target uint32, internalformat uint32, pname uint32, bufSize int32, params *int32) {
	C.wrap_glGetInternalformativ(C.uint(target), C.uint(internalformat), C.uint(pname), C.int(bufSize), (*C.int)(params))
}
func GetMultisamplefv(pname uint32, index uint32, val *float32) {
	C.wrap_glGetMultisamplefv(C.uint(pname), C.uint(index), (*C.float)(val))
}
func GetNamedFramebufferParameterivEXT(framebuffer uint32, pname uint32, params *int32) {
	C.wrap_glGetNamedFramebufferParameterivEXT(C.uint(framebuffer), C.uint(pname), (*C.int)(params))
}
func GetNamedStringARB(namelen int32, name string, bufSize int32, stringlen *int32, s_tring string) {
	cstr2 := C.CString(s_tring)
	defer C.free(unsafe.Pointer(cstr2))
	cstr1 := C.CString(name)
	defer C.free(unsafe.Pointer(cstr1))
	C.wrap_glGetNamedStringARB(C.int(namelen), cstr1, C.int(bufSize), (*C.int)(stringlen), cstr2)
}
func GetNamedStringivARB(namelen int32, name string, pname uint32, params *int32) {
	cstr1 := C.CString(name)
	defer C.free(unsafe.Pointer(cstr1))
	C.wrap_glGetNamedStringivARB(C.int(namelen), cstr1, C.uint(pname), (*C.int)(params))
}
func GetObjectLabel(identifier uint32, name uint32, bufSize int32, length *int32, label *uint8) {
	C.wrap_glGetObjectLabel(C.uint(identifier), C.uint(name), C.int(bufSize), (*C.int)(length), (*C.uchar)(label))
}
func GetObjectPtrLabel(ptr Pointer, bufSize int32, length *int32, label *uint8) {
	C.wrap_glGetObjectPtrLabel(unsafe.Pointer(ptr), C.int(bufSize), (*C.int)(length), (*C.uchar)(label))
}
func GetPointerv(pname uint32, params *Pointer) {
	C.wrap_glGetPointerv(C.uint(pname), (*unsafe.Pointer)(params))
}
func GetProgramBinary(program uint32, bufSize int32, length *int32, binaryFormat *uint32, binary Pointer) {
	C.wrap_glGetProgramBinary(C.uint(program), C.int(bufSize), (*C.int)(length), (*C.uint)(binaryFormat), unsafe.Pointer(binary))
}
func GetProgramInfoLog(program uint32, bufSize int32, length *int32, infoLog *uint8) {
	C.wrap_glGetProgramInfoLog(C.uint(program), C.int(bufSize), (*C.int)(length), (*C.uchar)(infoLog))
}
func GetProgramInterfaceiv(program uint32, programInterface uint32, pname uint32, params *int32) {
	C.wrap_glGetProgramInterfaceiv(C.uint(program), C.uint(programInterface), C.uint(pname), (*C.int)(params))
}
func GetProgramPipelineInfoLog(pipeline uint32, bufSize int32, length *int32, infoLog *uint8) {
	C.wrap_glGetProgramPipelineInfoLog(C.uint(pipeline), C.int(bufSize), (*C.int)(length), (*C.uchar)(infoLog))
}
func GetProgramPipelineiv(pipeline uint32, pname uint32, params *int32) {
	C.wrap_glGetProgramPipelineiv(C.uint(pipeline), C.uint(pname), (*C.int)(params))
}
func GetProgramResourceName(program uint32, programInterface uint32, index uint32, bufSize int32, length *int32, name string) {
	cstr1 := C.CString(name)
	defer C.free(unsafe.Pointer(cstr1))
	C.wrap_glGetProgramResourceName(C.uint(program), C.uint(programInterface), C.uint(index), C.int(bufSize), (*C.int)(length), cstr1)
}
func GetProgramResourceiv(program uint32, programInterface uint32, index uint32, propCount int32, props *uint32, bufSize int32, length *int32, params *int32) {
	C.wrap_glGetProgramResourceiv(C.uint(program), C.uint(programInterface), C.uint(index), C.int(propCount), (*C.uint)(props), C.int(bufSize), (*C.int)(length), (*C.int)(params))
}
func GetProgramStageiv(program uint32, shadertype uint32, pname uint32, values *int32) {
	C.wrap_glGetProgramStageiv(C.uint(program), C.uint(shadertype), C.uint(pname), (*C.int)(values))
}
func GetProgramiv(program uint32, pname uint32, params *int32) {
	C.wrap_glGetProgramiv(C.uint(program), C.uint(pname), (*C.int)(params))
}
func GetQueryIndexediv(target uint32, index uint32, pname uint32, params *int32) {
	C.wrap_glGetQueryIndexediv(C.uint(target), C.uint(index), C.uint(pname), (*C.int)(params))
}
func GetQueryObjecti64v(id uint32, pname uint32, params *int64) {
	C.wrap_glGetQueryObjecti64v(C.uint(id), C.uint(pname), (*C.longlong)(params))
}
func GetQueryObjectiv(id uint32, pname uint32, params *int32) {
	C.wrap_glGetQueryObjectiv(C.uint(id), C.uint(pname), (*C.int)(params))
}
func GetQueryObjectui64v(id uint32, pname uint32, params *uint64) {
	C.wrap_glGetQueryObjectui64v(C.uint(id), C.uint(pname), (*C.ulonglong)(params))
}
func GetQueryObjectuiv(id uint32, pname uint32, params *uint32) {
	C.wrap_glGetQueryObjectuiv(C.uint(id), C.uint(pname), (*C.uint)(params))
}
func GetQueryiv(target uint32, pname uint32, params *int32) {
	C.wrap_glGetQueryiv(C.uint(target), C.uint(pname), (*C.int)(params))
}
func GetRenderbufferParameteriv(target uint32, pname uint32, params *int32) {
	C.wrap_glGetRenderbufferParameteriv(C.uint(target), C.uint(pname), (*C.int)(params))
}
func GetSamplerParameterIiv(sampler uint32, pname uint32, params *int32) {
	C.wrap_glGetSamplerParameterIiv(C.uint(sampler), C.uint(pname), (*C.int)(params))
}
func GetSamplerParameterIuiv(sampler uint32, pname uint32, params *uint32) {
	C.wrap_glGetSamplerParameterIuiv(C.uint(sampler), C.uint(pname), (*C.uint)(params))
}
func GetSamplerParameterfv(sampler uint32, pname uint32, params *float32) {
	C.wrap_glGetSamplerParameterfv(C.uint(sampler), C.uint(pname), (*C.float)(params))
}
func GetSamplerParameteriv(sampler uint32, pname uint32, params *int32) {
	C.wrap_glGetSamplerParameteriv(C.uint(sampler), C.uint(pname), (*C.int)(params))
}
func GetShaderInfoLog(shader uint32, bufSize int32, length *int32, infoLog *uint8) {
	C.wrap_glGetShaderInfoLog(C.uint(shader), C.int(bufSize), (*C.int)(length), (*C.uchar)(infoLog))
}
func GetShaderPrecisionFormat(shadertype uint32, precisiontype uint32, r_ange *int32, precision *int32) {
	C.wrap_glGetShaderPrecisionFormat(C.uint(shadertype), C.uint(precisiontype), (*C.int)(r_ange), (*C.int)(precision))
}
func GetShaderSource(shader uint32, bufSize int32, length *int32, source *uint8) {
	C.wrap_glGetShaderSource(C.uint(shader), C.int(bufSize), (*C.int)(length), (*C.uchar)(source))
}
func GetShaderiv(shader uint32, pname uint32, params *int32) {
	C.wrap_glGetShaderiv(C.uint(shader), C.uint(pname), (*C.int)(params))
}
func GetSynciv(sync Sync, pname uint32, bufSize int32, length *int32, values *int32) {
	C.wrap_glGetSynciv(C.GLsync(sync), C.uint(pname), C.int(bufSize), (*C.int)(length), (*C.int)(values))
}
func GetTexImage(target uint32, level int32, format uint32, t_ype uint32, pixels Pointer) {
	C.wrap_glGetTexImage(C.uint(target), C.int(level), C.uint(format), C.uint(t_ype), unsafe.Pointer(pixels))
}
func GetTexLevelParameterfv(target uint32, level int32, pname uint32, params *float32) {
	C.wrap_glGetTexLevelParameterfv(C.uint(target), C.int(level), C.uint(pname), (*C.float)(params))
}
func GetTexLevelParameteriv(target uint32, level int32, pname uint32, params *int32) {
	C.wrap_glGetTexLevelParameteriv(C.uint(target), C.int(level), C.uint(pname), (*C.int)(params))
}
func GetTexParameterIiv(target uint32, pname uint32, params *int32) {
	C.wrap_glGetTexParameterIiv(C.uint(target), C.uint(pname), (*C.int)(params))
}
func GetTexParameterIuiv(target uint32, pname uint32, params *uint32) {
	C.wrap_glGetTexParameterIuiv(C.uint(target), C.uint(pname), (*C.uint)(params))
}
func GetTexParameterfv(target uint32, pname uint32, params *float32) {
	C.wrap_glGetTexParameterfv(C.uint(target), C.uint(pname), (*C.float)(params))
}
func GetTexParameteriv(target uint32, pname uint32, params *int32) {
	C.wrap_glGetTexParameteriv(C.uint(target), C.uint(pname), (*C.int)(params))
}
func GetTransformFeedbackVarying(program uint32, index uint32, bufSize int32, length *int32, size *int32, t_ype *uint32, name string) {
	cstr1 := C.CString(name)
	defer C.free(unsafe.Pointer(cstr1))
	C.wrap_glGetTransformFeedbackVarying(C.uint(program), C.uint(index), C.int(bufSize), (*C.int)(length), (*C.int)(size), (*C.uint)(t_ype), cstr1)
}
func GetUniformIndices(program uint32, uniformCount int32, uniformNames []string, uniformIndices *uint32) {
	cstrings := C.newStringArray(C.int(len(uniformNames)))
	defer C.freeStringArray(cstrings, C.int(len(uniformNames)))
	for cnt, str := range uniformNames {
		C.assignString(cstrings, C.CString(str), C.int(cnt))
	}
	C.wrap_glGetUniformIndices(C.uint(program), C.int(uniformCount), cstrings, (*C.uint)(uniformIndices))
}
func GetUniformSubroutineuiv(shadertype uint32, location int32, params *uint32) {
	C.wrap_glGetUniformSubroutineuiv(C.uint(shadertype), C.int(location), (*C.uint)(params))
}
func GetUniformdv(program uint32, location int32, params *float64) {
	C.wrap_glGetUniformdv(C.uint(program), C.int(location), (*C.double)(params))
}
func GetUniformfv(program uint32, location int32, params *float32) {
	C.wrap_glGetUniformfv(C.uint(program), C.int(location), (*C.float)(params))
}
func GetUniformiv(program uint32, location int32, params *int32) {
	C.wrap_glGetUniformiv(C.uint(program), C.int(location), (*C.int)(params))
}
func GetUniformuiv(program uint32, location int32, params *uint32) {
	C.wrap_glGetUniformuiv(C.uint(program), C.int(location), (*C.uint)(params))
}
func GetVertexAttribIiv(index uint32, pname uint32, params *int32) {
	C.wrap_glGetVertexAttribIiv(C.uint(index), C.uint(pname), (*C.int)(params))
}
func GetVertexAttribIuiv(index uint32, pname uint32, params *uint32) {
	C.wrap_glGetVertexAttribIuiv(C.uint(index), C.uint(pname), (*C.uint)(params))
}
func GetVertexAttribLdv(index uint32, pname uint32, params *float64) {
	C.wrap_glGetVertexAttribLdv(C.uint(index), C.uint(pname), (*C.double)(params))
}
func GetVertexAttribPointerv(index uint32, pname uint32, pointer *int64) {
	C.wrap_glGetVertexAttribPointerv(C.uint(index), C.uint(pname), (*C.longlong)(pointer))
}
func GetVertexAttribdv(index uint32, pname uint32, params *float64) {
	C.wrap_glGetVertexAttribdv(C.uint(index), C.uint(pname), (*C.double)(params))
}
func GetVertexAttribfv(index uint32, pname uint32, params *float32) {
	C.wrap_glGetVertexAttribfv(C.uint(index), C.uint(pname), (*C.float)(params))
}
func GetVertexAttribiv(index uint32, pname uint32, params *int32) {
	C.wrap_glGetVertexAttribiv(C.uint(index), C.uint(pname), (*C.int)(params))
}
func GetnColorTableARB(target uint32, format uint32, t_ype uint32, bufSize int32, table Pointer) {
	C.wrap_glGetnColorTableARB(C.uint(target), C.uint(format), C.uint(t_ype), C.int(bufSize), unsafe.Pointer(table))
}
func GetnCompressedTexImageARB(target uint32, lod int32, bufSize int32, img Pointer) {
	C.wrap_glGetnCompressedTexImageARB(C.uint(target), C.int(lod), C.int(bufSize), unsafe.Pointer(img))
}
func GetnConvolutionFilterARB(target uint32, format uint32, t_ype uint32, bufSize int32, image Pointer) {
	C.wrap_glGetnConvolutionFilterARB(C.uint(target), C.uint(format), C.uint(t_ype), C.int(bufSize), unsafe.Pointer(image))
}
func GetnHistogramARB(target uint32, reset bool, format uint32, t_ype uint32, bufSize int32, values Pointer) {
	tf1 := FALSE
	if reset {
		tf1 = TRUE
	}
	C.wrap_glGetnHistogramARB(C.uint(target), C.uchar(tf1), C.uint(format), C.uint(t_ype), C.int(bufSize), unsafe.Pointer(values))
}
func GetnMapdvARB(target uint32, query uint32, bufSize int32, v *float64) {
	C.wrap_glGetnMapdvARB(C.uint(target), C.uint(query), C.int(bufSize), (*C.double)(v))
}
func GetnMapfvARB(target uint32, query uint32, bufSize int32, v *float32) {
	C.wrap_glGetnMapfvARB(C.uint(target), C.uint(query), C.int(bufSize), (*C.float)(v))
}
func GetnMapivARB(target uint32, query uint32, bufSize int32, v *int32) {
	C.wrap_glGetnMapivARB(C.uint(target), C.uint(query), C.int(bufSize), (*C.int)(v))
}
func GetnMinmaxARB(target uint32, reset bool, format uint32, t_ype uint32, bufSize int32, values Pointer) {
	tf1 := FALSE
	if reset {
		tf1 = TRUE
	}
	C.wrap_glGetnMinmaxARB(C.uint(target), C.uchar(tf1), C.uint(format), C.uint(t_ype), C.int(bufSize), unsafe.Pointer(values))
}
func GetnPixelMapfvARB(m_ap uint32, bufSize int32, values *float32) {
	C.wrap_glGetnPixelMapfvARB(C.uint(m_ap), C.int(bufSize), (*C.float)(values))
}
func GetnPixelMapuivARB(m_ap uint32, bufSize int32, values *uint32) {
	C.wrap_glGetnPixelMapuivARB(C.uint(m_ap), C.int(bufSize), (*C.uint)(values))
}
func GetnPixelMapusvARB(m_ap uint32, bufSize int32, values *uint16) {
	C.wrap_glGetnPixelMapusvARB(C.uint(m_ap), C.int(bufSize), (*C.ushort)(values))
}
func GetnPolygonStippleARB(bufSize int32, pattern *uint8) {
	C.wrap_glGetnPolygonStippleARB(C.int(bufSize), (*C.uchar)(pattern))
}
func GetnSeparableFilterARB(target uint32, format uint32, t_ype uint32, rowBufSize int32, row Pointer, columnBufSize int32, column Pointer, span Pointer) {
	C.wrap_glGetnSeparableFilterARB(C.uint(target), C.uint(format), C.uint(t_ype), C.int(rowBufSize), unsafe.Pointer(row), C.int(columnBufSize), unsafe.Pointer(column), unsafe.Pointer(span))
}
func GetnTexImageARB(target uint32, level int32, format uint32, t_ype uint32, bufSize int32, img Pointer) {
	C.wrap_glGetnTexImageARB(C.uint(target), C.int(level), C.uint(format), C.uint(t_ype), C.int(bufSize), unsafe.Pointer(img))
}
func GetnUniformdvARB(program uint32, location int32, bufSize int32, params *float64) {
	C.wrap_glGetnUniformdvARB(C.uint(program), C.int(location), C.int(bufSize), (*C.double)(params))
}
func GetnUniformfvARB(program uint32, location int32, bufSize int32, params *float32) {
	C.wrap_glGetnUniformfvARB(C.uint(program), C.int(location), C.int(bufSize), (*C.float)(params))
}
func GetnUniformivARB(program uint32, location int32, bufSize int32, params *int32) {
	C.wrap_glGetnUniformivARB(C.uint(program), C.int(location), C.int(bufSize), (*C.int)(params))
}
func GetnUniformuivARB(program uint32, location int32, bufSize int32, params *uint32) {
	C.wrap_glGetnUniformuivARB(C.uint(program), C.int(location), C.int(bufSize), (*C.uint)(params))
}
func Hint(target uint32, mode uint32) {
	C.wrap_glHint(C.uint(target), C.uint(mode))
}
func InvalidateBufferData(buffer uint32) {
	C.wrap_glInvalidateBufferData(C.uint(buffer))
}
func InvalidateBufferSubData(buffer uint32, offset int64, length int64) {
	C.wrap_glInvalidateBufferSubData(C.uint(buffer), C.longlong(offset), C.longlong(length))
}
func InvalidateFramebuffer(target uint32, numAttachments int32, attachments *uint32) {
	C.wrap_glInvalidateFramebuffer(C.uint(target), C.int(numAttachments), (*C.uint)(attachments))
}
func InvalidateSubFramebuffer(target uint32, numAttachments int32, attachments *uint32, x int32, y int32, width int32, height int32) {
	C.wrap_glInvalidateSubFramebuffer(C.uint(target), C.int(numAttachments), (*C.uint)(attachments), C.int(x), C.int(y), C.int(width), C.int(height))
}
func InvalidateTexImage(texture uint32, level int32) {
	C.wrap_glInvalidateTexImage(C.uint(texture), C.int(level))
}
func InvalidateTexSubImage(texture uint32, level int32, xoffset int32, yoffset int32, zoffset int32, width int32, height int32, depth int32) {
	C.wrap_glInvalidateTexSubImage(C.uint(texture), C.int(level), C.int(xoffset), C.int(yoffset), C.int(zoffset), C.int(width), C.int(height), C.int(depth))
}
func LineWidth(width float32) {
	C.wrap_glLineWidth(C.float(width))
}
func LinkProgram(program uint32) {
	C.wrap_glLinkProgram(C.uint(program))
}
func LogicOp(opcode uint32) {
	C.wrap_glLogicOp(C.uint(opcode))
}
func MemoryBarrier(barriers uint32) {
	C.wrap_glMemoryBarrier(C.uint(barriers))
}
func MinSampleShading(value float32) {
	C.wrap_glMinSampleShading(C.float(value))
}
func MinSampleShadingARB(value float32) {
	C.wrap_glMinSampleShadingARB(C.float(value))
}
func MultiDrawArrays(mode uint32, first *int32, count *int32, drawcount int32) {
	C.wrap_glMultiDrawArrays(C.uint(mode), (*C.int)(first), (*C.int)(count), C.int(drawcount))
}
func MultiDrawArraysIndirect(mode uint32, indirect Pointer, drawcount int32, stride int32) {
	C.wrap_glMultiDrawArraysIndirect(C.uint(mode), unsafe.Pointer(indirect), C.int(drawcount), C.int(stride))
}
func MultiDrawElements(mode uint32, count *int32, t_ype uint32, indices *Pointer, drawcount int32) {
	C.wrap_glMultiDrawElements(C.uint(mode), (*C.int)(count), C.uint(t_ype), (*unsafe.Pointer)(indices), C.int(drawcount))
}
func MultiDrawElementsBaseVertex(mode uint32, count *int32, t_ype uint32, indices *Pointer, drawcount int32, basevertex *int32) {
	C.wrap_glMultiDrawElementsBaseVertex(C.uint(mode), (*C.int)(count), C.uint(t_ype), (*unsafe.Pointer)(indices), C.int(drawcount), (*C.int)(basevertex))
}
func MultiDrawElementsIndirect(mode uint32, t_ype uint32, indirect Pointer, drawcount int32, stride int32) {
	C.wrap_glMultiDrawElementsIndirect(C.uint(mode), C.uint(t_ype), unsafe.Pointer(indirect), C.int(drawcount), C.int(stride))
}
func MultiTexCoordP1ui(texture uint32, t_ype uint32, coords uint32) {
	C.wrap_glMultiTexCoordP1ui(C.uint(texture), C.uint(t_ype), C.uint(coords))
}
func MultiTexCoordP1uiv(texture uint32, t_ype uint32, coords *uint32) {
	C.wrap_glMultiTexCoordP1uiv(C.uint(texture), C.uint(t_ype), (*C.uint)(coords))
}
func MultiTexCoordP2ui(texture uint32, t_ype uint32, coords uint32) {
	C.wrap_glMultiTexCoordP2ui(C.uint(texture), C.uint(t_ype), C.uint(coords))
}
func MultiTexCoordP2uiv(texture uint32, t_ype uint32, coords *uint32) {
	C.wrap_glMultiTexCoordP2uiv(C.uint(texture), C.uint(t_ype), (*C.uint)(coords))
}
func MultiTexCoordP3ui(texture uint32, t_ype uint32, coords uint32) {
	C.wrap_glMultiTexCoordP3ui(C.uint(texture), C.uint(t_ype), C.uint(coords))
}
func MultiTexCoordP3uiv(texture uint32, t_ype uint32, coords *uint32) {
	C.wrap_glMultiTexCoordP3uiv(C.uint(texture), C.uint(t_ype), (*C.uint)(coords))
}
func MultiTexCoordP4ui(texture uint32, t_ype uint32, coords uint32) {
	C.wrap_glMultiTexCoordP4ui(C.uint(texture), C.uint(t_ype), C.uint(coords))
}
func MultiTexCoordP4uiv(texture uint32, t_ype uint32, coords *uint32) {
	C.wrap_glMultiTexCoordP4uiv(C.uint(texture), C.uint(t_ype), (*C.uint)(coords))
}
func NamedFramebufferParameteriEXT(framebuffer uint32, pname uint32, param int32) {
	C.wrap_glNamedFramebufferParameteriEXT(C.uint(framebuffer), C.uint(pname), C.int(param))
}
func NamedStringARB(t_ype uint32, namelen int32, name string, stringlen int32, s_tring string) {
	cstr2 := C.CString(s_tring)
	defer C.free(unsafe.Pointer(cstr2))
	cstr1 := C.CString(name)
	defer C.free(unsafe.Pointer(cstr1))
	C.wrap_glNamedStringARB(C.uint(t_ype), C.int(namelen), cstr1, C.int(stringlen), cstr2)
}
func NormalP3ui(t_ype uint32, coords uint32) {
	C.wrap_glNormalP3ui(C.uint(t_ype), C.uint(coords))
}
func NormalP3uiv(t_ype uint32, coords *uint32) {
	C.wrap_glNormalP3uiv(C.uint(t_ype), (*C.uint)(coords))
}
func ObjectLabel(identifier uint32, name uint32, length int32, label string) {
	cstr1 := C.CString(label)
	defer C.free(unsafe.Pointer(cstr1))
	C.wrap_glObjectLabel(C.uint(identifier), C.uint(name), C.int(length), cstr1)
}
func ObjectPtrLabel(ptr Pointer, length int32, label string) {
	cstr1 := C.CString(label)
	defer C.free(unsafe.Pointer(cstr1))
	C.wrap_glObjectPtrLabel(unsafe.Pointer(ptr), C.int(length), cstr1)
}
func PatchParameterfv(pname uint32, values *float32) {
	C.wrap_glPatchParameterfv(C.uint(pname), (*C.float)(values))
}
func PatchParameteri(pname uint32, value int32) {
	C.wrap_glPatchParameteri(C.uint(pname), C.int(value))
}
func PauseTransformFeedback() {
	C.wrap_glPauseTransformFeedback()
}
func PixelStoref(pname uint32, param float32) {
	C.wrap_glPixelStoref(C.uint(pname), C.float(param))
}
func PixelStorei(pname uint32, param int32) {
	C.wrap_glPixelStorei(C.uint(pname), C.int(param))
}
func PointParameterf(pname uint32, param float32) {
	C.wrap_glPointParameterf(C.uint(pname), C.float(param))
}
func PointParameterfv(pname uint32, params *float32) {
	C.wrap_glPointParameterfv(C.uint(pname), (*C.float)(params))
}
func PointParameteri(pname uint32, param int32) {
	C.wrap_glPointParameteri(C.uint(pname), C.int(param))
}
func PointParameteriv(pname uint32, params *int32) {
	C.wrap_glPointParameteriv(C.uint(pname), (*C.int)(params))
}
func PointSize(size float32) {
	C.wrap_glPointSize(C.float(size))
}
func PolygonMode(face uint32, mode uint32) {
	C.wrap_glPolygonMode(C.uint(face), C.uint(mode))
}
func PolygonOffset(factor float32, units float32) {
	C.wrap_glPolygonOffset(C.float(factor), C.float(units))
}
func PopDebugGroup() {
	C.wrap_glPopDebugGroup()
}
func PrimitiveRestartIndex(index uint32) {
	C.wrap_glPrimitiveRestartIndex(C.uint(index))
}
func ProgramBinary(program uint32, binaryFormat uint32, binary Pointer, length int32) {
	C.wrap_glProgramBinary(C.uint(program), C.uint(binaryFormat), unsafe.Pointer(binary), C.int(length))
}
func ProgramParameteri(program uint32, pname uint32, value int32) {
	C.wrap_glProgramParameteri(C.uint(program), C.uint(pname), C.int(value))
}
func ProgramUniform1d(program uint32, location int32, v0 float64) {
	C.wrap_glProgramUniform1d(C.uint(program), C.int(location), C.double(v0))
}
func ProgramUniform1dv(program uint32, location int32, count int32, value *float64) {
	C.wrap_glProgramUniform1dv(C.uint(program), C.int(location), C.int(count), (*C.double)(value))
}
func ProgramUniform1f(program uint32, location int32, v0 float32) {
	C.wrap_glProgramUniform1f(C.uint(program), C.int(location), C.float(v0))
}
func ProgramUniform1fv(program uint32, location int32, count int32, value *float32) {
	C.wrap_glProgramUniform1fv(C.uint(program), C.int(location), C.int(count), (*C.float)(value))
}
func ProgramUniform1i(program uint32, location int32, v0 int32) {
	C.wrap_glProgramUniform1i(C.uint(program), C.int(location), C.int(v0))
}
func ProgramUniform1iv(program uint32, location int32, count int32, value *int32) {
	C.wrap_glProgramUniform1iv(C.uint(program), C.int(location), C.int(count), (*C.int)(value))
}
func ProgramUniform1ui(program uint32, location int32, v0 uint32) {
	C.wrap_glProgramUniform1ui(C.uint(program), C.int(location), C.uint(v0))
}
func ProgramUniform1uiv(program uint32, location int32, count int32, value *uint32) {
	C.wrap_glProgramUniform1uiv(C.uint(program), C.int(location), C.int(count), (*C.uint)(value))
}
func ProgramUniform2d(program uint32, location int32, v0 float64, v1 float64) {
	C.wrap_glProgramUniform2d(C.uint(program), C.int(location), C.double(v0), C.double(v1))
}
func ProgramUniform2dv(program uint32, location int32, count int32, value *float64) {
	C.wrap_glProgramUniform2dv(C.uint(program), C.int(location), C.int(count), (*C.double)(value))
}
func ProgramUniform2f(program uint32, location int32, v0 float32, v1 float32) {
	C.wrap_glProgramUniform2f(C.uint(program), C.int(location), C.float(v0), C.float(v1))
}
func ProgramUniform2fv(program uint32, location int32, count int32, value *float32) {
	C.wrap_glProgramUniform2fv(C.uint(program), C.int(location), C.int(count), (*C.float)(value))
}
func ProgramUniform2i(program uint32, location int32, v0 int32, v1 int32) {
	C.wrap_glProgramUniform2i(C.uint(program), C.int(location), C.int(v0), C.int(v1))
}
func ProgramUniform2iv(program uint32, location int32, count int32, value *int32) {
	C.wrap_glProgramUniform2iv(C.uint(program), C.int(location), C.int(count), (*C.int)(value))
}
func ProgramUniform2ui(program uint32, location int32, v0 uint32, v1 uint32) {
	C.wrap_glProgramUniform2ui(C.uint(program), C.int(location), C.uint(v0), C.uint(v1))
}
func ProgramUniform2uiv(program uint32, location int32, count int32, value *uint32) {
	C.wrap_glProgramUniform2uiv(C.uint(program), C.int(location), C.int(count), (*C.uint)(value))
}
func ProgramUniform3d(program uint32, location int32, v0 float64, v1 float64, v2 float64) {
	C.wrap_glProgramUniform3d(C.uint(program), C.int(location), C.double(v0), C.double(v1), C.double(v2))
}
func ProgramUniform3dv(program uint32, location int32, count int32, value *float64) {
	C.wrap_glProgramUniform3dv(C.uint(program), C.int(location), C.int(count), (*C.double)(value))
}
func ProgramUniform3f(program uint32, location int32, v0 float32, v1 float32, v2 float32) {
	C.wrap_glProgramUniform3f(C.uint(program), C.int(location), C.float(v0), C.float(v1), C.float(v2))
}
func ProgramUniform3fv(program uint32, location int32, count int32, value *float32) {
	C.wrap_glProgramUniform3fv(C.uint(program), C.int(location), C.int(count), (*C.float)(value))
}
func ProgramUniform3i(program uint32, location int32, v0 int32, v1 int32, v2 int32) {
	C.wrap_glProgramUniform3i(C.uint(program), C.int(location), C.int(v0), C.int(v1), C.int(v2))
}
func ProgramUniform3iv(program uint32, location int32, count int32, value *int32) {
	C.wrap_glProgramUniform3iv(C.uint(program), C.int(location), C.int(count), (*C.int)(value))
}
func ProgramUniform3ui(program uint32, location int32, v0 uint32, v1 uint32, v2 uint32) {
	C.wrap_glProgramUniform3ui(C.uint(program), C.int(location), C.uint(v0), C.uint(v1), C.uint(v2))
}
func ProgramUniform3uiv(program uint32, location int32, count int32, value *uint32) {
	C.wrap_glProgramUniform3uiv(C.uint(program), C.int(location), C.int(count), (*C.uint)(value))
}
func ProgramUniform4d(program uint32, location int32, v0 float64, v1 float64, v2 float64, v3 float64) {
	C.wrap_glProgramUniform4d(C.uint(program), C.int(location), C.double(v0), C.double(v1), C.double(v2), C.double(v3))
}
func ProgramUniform4dv(program uint32, location int32, count int32, value *float64) {
	C.wrap_glProgramUniform4dv(C.uint(program), C.int(location), C.int(count), (*C.double)(value))
}
func ProgramUniform4f(program uint32, location int32, v0 float32, v1 float32, v2 float32, v3 float32) {
	C.wrap_glProgramUniform4f(C.uint(program), C.int(location), C.float(v0), C.float(v1), C.float(v2), C.float(v3))
}
func ProgramUniform4fv(program uint32, location int32, count int32, value *float32) {
	C.wrap_glProgramUniform4fv(C.uint(program), C.int(location), C.int(count), (*C.float)(value))
}
func ProgramUniform4i(program uint32, location int32, v0 int32, v1 int32, v2 int32, v3 int32) {
	C.wrap_glProgramUniform4i(C.uint(program), C.int(location), C.int(v0), C.int(v1), C.int(v2), C.int(v3))
}
func ProgramUniform4iv(program uint32, location int32, count int32, value *int32) {
	C.wrap_glProgramUniform4iv(C.uint(program), C.int(location), C.int(count), (*C.int)(value))
}
func ProgramUniform4ui(program uint32, location int32, v0 uint32, v1 uint32, v2 uint32, v3 uint32) {
	C.wrap_glProgramUniform4ui(C.uint(program), C.int(location), C.uint(v0), C.uint(v1), C.uint(v2), C.uint(v3))
}
func ProgramUniform4uiv(program uint32, location int32, count int32, value *uint32) {
	C.wrap_glProgramUniform4uiv(C.uint(program), C.int(location), C.int(count), (*C.uint)(value))
}
func ProgramUniformMatrix2dv(program uint32, location int32, count int32, transpose bool, value *float64) {
	tf1 := FALSE
	if transpose {
		tf1 = TRUE
	}
	C.wrap_glProgramUniformMatrix2dv(C.uint(program), C.int(location), C.int(count), C.uchar(tf1), (*C.double)(value))
}
func ProgramUniformMatrix2fv(program uint32, location int32, count int32, transpose bool, value *float32) {
	tf1 := FALSE
	if transpose {
		tf1 = TRUE
	}
	C.wrap_glProgramUniformMatrix2fv(C.uint(program), C.int(location), C.int(count), C.uchar(tf1), (*C.float)(value))
}
func ProgramUniformMatrix2x3dv(program uint32, location int32, count int32, transpose bool, value *float64) {
	tf1 := FALSE
	if transpose {
		tf1 = TRUE
	}
	C.wrap_glProgramUniformMatrix2x3dv(C.uint(program), C.int(location), C.int(count), C.uchar(tf1), (*C.double)(value))
}
func ProgramUniformMatrix2x3fv(program uint32, location int32, count int32, transpose bool, value *float32) {
	tf1 := FALSE
	if transpose {
		tf1 = TRUE
	}
	C.wrap_glProgramUniformMatrix2x3fv(C.uint(program), C.int(location), C.int(count), C.uchar(tf1), (*C.float)(value))
}
func ProgramUniformMatrix2x4dv(program uint32, location int32, count int32, transpose bool, value *float64) {
	tf1 := FALSE
	if transpose {
		tf1 = TRUE
	}
	C.wrap_glProgramUniformMatrix2x4dv(C.uint(program), C.int(location), C.int(count), C.uchar(tf1), (*C.double)(value))
}
func ProgramUniformMatrix2x4fv(program uint32, location int32, count int32, transpose bool, value *float32) {
	tf1 := FALSE
	if transpose {
		tf1 = TRUE
	}
	C.wrap_glProgramUniformMatrix2x4fv(C.uint(program), C.int(location), C.int(count), C.uchar(tf1), (*C.float)(value))
}
func ProgramUniformMatrix3dv(program uint32, location int32, count int32, transpose bool, value *float64) {
	tf1 := FALSE
	if transpose {
		tf1 = TRUE
	}
	C.wrap_glProgramUniformMatrix3dv(C.uint(program), C.int(location), C.int(count), C.uchar(tf1), (*C.double)(value))
}
func ProgramUniformMatrix3fv(program uint32, location int32, count int32, transpose bool, value *float32) {
	tf1 := FALSE
	if transpose {
		tf1 = TRUE
	}
	C.wrap_glProgramUniformMatrix3fv(C.uint(program), C.int(location), C.int(count), C.uchar(tf1), (*C.float)(value))
}
func ProgramUniformMatrix3x2dv(program uint32, location int32, count int32, transpose bool, value *float64) {
	tf1 := FALSE
	if transpose {
		tf1 = TRUE
	}
	C.wrap_glProgramUniformMatrix3x2dv(C.uint(program), C.int(location), C.int(count), C.uchar(tf1), (*C.double)(value))
}
func ProgramUniformMatrix3x2fv(program uint32, location int32, count int32, transpose bool, value *float32) {
	tf1 := FALSE
	if transpose {
		tf1 = TRUE
	}
	C.wrap_glProgramUniformMatrix3x2fv(C.uint(program), C.int(location), C.int(count), C.uchar(tf1), (*C.float)(value))
}
func ProgramUniformMatrix3x4dv(program uint32, location int32, count int32, transpose bool, value *float64) {
	tf1 := FALSE
	if transpose {
		tf1 = TRUE
	}
	C.wrap_glProgramUniformMatrix3x4dv(C.uint(program), C.int(location), C.int(count), C.uchar(tf1), (*C.double)(value))
}
func ProgramUniformMatrix3x4fv(program uint32, location int32, count int32, transpose bool, value *float32) {
	tf1 := FALSE
	if transpose {
		tf1 = TRUE
	}
	C.wrap_glProgramUniformMatrix3x4fv(C.uint(program), C.int(location), C.int(count), C.uchar(tf1), (*C.float)(value))
}
func ProgramUniformMatrix4dv(program uint32, location int32, count int32, transpose bool, value *float64) {
	tf1 := FALSE
	if transpose {
		tf1 = TRUE
	}
	C.wrap_glProgramUniformMatrix4dv(C.uint(program), C.int(location), C.int(count), C.uchar(tf1), (*C.double)(value))
}
func ProgramUniformMatrix4fv(program uint32, location int32, count int32, transpose bool, value *float32) {
	tf1 := FALSE
	if transpose {
		tf1 = TRUE
	}
	C.wrap_glProgramUniformMatrix4fv(C.uint(program), C.int(location), C.int(count), C.uchar(tf1), (*C.float)(value))
}
func ProgramUniformMatrix4x2dv(program uint32, location int32, count int32, transpose bool, value *float64) {
	tf1 := FALSE
	if transpose {
		tf1 = TRUE
	}
	C.wrap_glProgramUniformMatrix4x2dv(C.uint(program), C.int(location), C.int(count), C.uchar(tf1), (*C.double)(value))
}
func ProgramUniformMatrix4x2fv(program uint32, location int32, count int32, transpose bool, value *float32) {
	tf1 := FALSE
	if transpose {
		tf1 = TRUE
	}
	C.wrap_glProgramUniformMatrix4x2fv(C.uint(program), C.int(location), C.int(count), C.uchar(tf1), (*C.float)(value))
}
func ProgramUniformMatrix4x3dv(program uint32, location int32, count int32, transpose bool, value *float64) {
	tf1 := FALSE
	if transpose {
		tf1 = TRUE
	}
	C.wrap_glProgramUniformMatrix4x3dv(C.uint(program), C.int(location), C.int(count), C.uchar(tf1), (*C.double)(value))
}
func ProgramUniformMatrix4x3fv(program uint32, location int32, count int32, transpose bool, value *float32) {
	tf1 := FALSE
	if transpose {
		tf1 = TRUE
	}
	C.wrap_glProgramUniformMatrix4x3fv(C.uint(program), C.int(location), C.int(count), C.uchar(tf1), (*C.float)(value))
}
func ProvokingVertex(mode uint32) {
	C.wrap_glProvokingVertex(C.uint(mode))
}
func PushDebugGroup(source uint32, id uint32, length int32, message string) {
	cstr1 := C.CString(message)
	defer C.free(unsafe.Pointer(cstr1))
	C.wrap_glPushDebugGroup(C.uint(source), C.uint(id), C.int(length), cstr1)
}
func QueryCounter(id uint32, target uint32) {
	C.wrap_glQueryCounter(C.uint(id), C.uint(target))
}
func ReadBuffer(mode uint32) {
	C.wrap_glReadBuffer(C.uint(mode))
}
func ReadPixels(x int32, y int32, width int32, height int32, format uint32, t_ype uint32, pixels Pointer) {
	C.wrap_glReadPixels(C.int(x), C.int(y), C.int(width), C.int(height), C.uint(format), C.uint(t_ype), unsafe.Pointer(pixels))
}
func ReadnPixelsARB(x int32, y int32, width int32, height int32, format uint32, t_ype uint32, bufSize int32, data Pointer) {
	C.wrap_glReadnPixelsARB(C.int(x), C.int(y), C.int(width), C.int(height), C.uint(format), C.uint(t_ype), C.int(bufSize), unsafe.Pointer(data))
}
func ReleaseShaderCompiler() {
	C.wrap_glReleaseShaderCompiler()
}
func RenderbufferStorage(target uint32, internalformat uint32, width int32, height int32) {
	C.wrap_glRenderbufferStorage(C.uint(target), C.uint(internalformat), C.int(width), C.int(height))
}
func RenderbufferStorageMultisample(target uint32, samples int32, internalformat uint32, width int32, height int32) {
	C.wrap_glRenderbufferStorageMultisample(C.uint(target), C.int(samples), C.uint(internalformat), C.int(width), C.int(height))
}
func ResumeTransformFeedback() {
	C.wrap_glResumeTransformFeedback()
}
func SampleCoverage(value float32, invert bool) {
	tf1 := FALSE
	if invert {
		tf1 = TRUE
	}
	C.wrap_glSampleCoverage(C.float(value), C.uchar(tf1))
}
func SampleMaski(index uint32, mask uint32) {
	C.wrap_glSampleMaski(C.uint(index), C.uint(mask))
}
func SamplerParameterIiv(sampler uint32, pname uint32, param *int32) {
	C.wrap_glSamplerParameterIiv(C.uint(sampler), C.uint(pname), (*C.int)(param))
}
func SamplerParameterIuiv(sampler uint32, pname uint32, param *uint32) {
	C.wrap_glSamplerParameterIuiv(C.uint(sampler), C.uint(pname), (*C.uint)(param))
}
func SamplerParameterf(sampler uint32, pname uint32, param float32) {
	C.wrap_glSamplerParameterf(C.uint(sampler), C.uint(pname), C.float(param))
}
func SamplerParameterfv(sampler uint32, pname uint32, param *float32) {
	C.wrap_glSamplerParameterfv(C.uint(sampler), C.uint(pname), (*C.float)(param))
}
func SamplerParameteri(sampler uint32, pname uint32, param int32) {
	C.wrap_glSamplerParameteri(C.uint(sampler), C.uint(pname), C.int(param))
}
func SamplerParameteriv(sampler uint32, pname uint32, param *int32) {
	C.wrap_glSamplerParameteriv(C.uint(sampler), C.uint(pname), (*C.int)(param))
}
func Scissor(x int32, y int32, width int32, height int32) {
	C.wrap_glScissor(C.int(x), C.int(y), C.int(width), C.int(height))
}
func ScissorArrayv(first uint32, count int32, v *int32) {
	C.wrap_glScissorArrayv(C.uint(first), C.int(count), (*C.int)(v))
}
func ScissorIndexed(index uint32, left int32, bottom int32, width int32, height int32) {
	C.wrap_glScissorIndexed(C.uint(index), C.int(left), C.int(bottom), C.int(width), C.int(height))
}
func ScissorIndexedv(index uint32, v *int32) {
	C.wrap_glScissorIndexedv(C.uint(index), (*C.int)(v))
}
func SecondaryColorP3ui(t_ype uint32, color uint32) {
	C.wrap_glSecondaryColorP3ui(C.uint(t_ype), C.uint(color))
}
func SecondaryColorP3uiv(t_ype uint32, color *uint32) {
	C.wrap_glSecondaryColorP3uiv(C.uint(t_ype), (*C.uint)(color))
}
func ShaderBinary(count int32, shaders *uint32, binaryformat uint32, binary Pointer, length int32) {
	C.wrap_glShaderBinary(C.int(count), (*C.uint)(shaders), C.uint(binaryformat), unsafe.Pointer(binary), C.int(length))
}
func ShaderSource(shader uint32, count int32, s_tring []string, length *int32) {
	cstrings := C.newStringArray(C.int(len(s_tring)))
	defer C.freeStringArray(cstrings, C.int(len(s_tring)))
	for cnt, str := range s_tring {
		C.assignString(cstrings, C.CString(str), C.int(cnt))
	}
	C.wrap_glShaderSource(C.uint(shader), C.int(count), cstrings, (*C.int)(length))
}
func ShaderStorageBlockBinding(program uint32, storageBlockIndex uint32, storageBlockBinding uint32) {
	C.wrap_glShaderStorageBlockBinding(C.uint(program), C.uint(storageBlockIndex), C.uint(storageBlockBinding))
}
func StencilFunc(f_unc uint32, ref int32, mask uint32) {
	C.wrap_glStencilFunc(C.uint(f_unc), C.int(ref), C.uint(mask))
}
func StencilFuncSeparate(face uint32, f_unc uint32, ref int32, mask uint32) {
	C.wrap_glStencilFuncSeparate(C.uint(face), C.uint(f_unc), C.int(ref), C.uint(mask))
}
func StencilMask(mask uint32) {
	C.wrap_glStencilMask(C.uint(mask))
}
func StencilMaskSeparate(face uint32, mask uint32) {
	C.wrap_glStencilMaskSeparate(C.uint(face), C.uint(mask))
}
func StencilOp(fail uint32, zfail uint32, zpass uint32) {
	C.wrap_glStencilOp(C.uint(fail), C.uint(zfail), C.uint(zpass))
}
func StencilOpSeparate(face uint32, sfail uint32, dpfail uint32, dppass uint32) {
	C.wrap_glStencilOpSeparate(C.uint(face), C.uint(sfail), C.uint(dpfail), C.uint(dppass))
}
func TexBuffer(target uint32, internalformat uint32, buffer uint32) {
	C.wrap_glTexBuffer(C.uint(target), C.uint(internalformat), C.uint(buffer))
}
func TexBufferRange(target uint32, internalformat uint32, buffer uint32, offset int64, size int64) {
	C.wrap_glTexBufferRange(C.uint(target), C.uint(internalformat), C.uint(buffer), C.longlong(offset), C.longlong(size))
}
func TexCoordP1ui(t_ype uint32, coords uint32) {
	C.wrap_glTexCoordP1ui(C.uint(t_ype), C.uint(coords))
}
func TexCoordP1uiv(t_ype uint32, coords *uint32) {
	C.wrap_glTexCoordP1uiv(C.uint(t_ype), (*C.uint)(coords))
}
func TexCoordP2ui(t_ype uint32, coords uint32) {
	C.wrap_glTexCoordP2ui(C.uint(t_ype), C.uint(coords))
}
func TexCoordP2uiv(t_ype uint32, coords *uint32) {
	C.wrap_glTexCoordP2uiv(C.uint(t_ype), (*C.uint)(coords))
}
func TexCoordP3ui(t_ype uint32, coords uint32) {
	C.wrap_glTexCoordP3ui(C.uint(t_ype), C.uint(coords))
}
func TexCoordP3uiv(t_ype uint32, coords *uint32) {
	C.wrap_glTexCoordP3uiv(C.uint(t_ype), (*C.uint)(coords))
}
func TexCoordP4ui(t_ype uint32, coords uint32) {
	C.wrap_glTexCoordP4ui(C.uint(t_ype), C.uint(coords))
}
func TexCoordP4uiv(t_ype uint32, coords *uint32) {
	C.wrap_glTexCoordP4uiv(C.uint(t_ype), (*C.uint)(coords))
}
func TexImage1D(target uint32, level int32, internalformat int32, width int32, border int32, format uint32, t_ype uint32, pixels Pointer) {
	C.wrap_glTexImage1D(C.uint(target), C.int(level), C.int(internalformat), C.int(width), C.int(border), C.uint(format), C.uint(t_ype), unsafe.Pointer(pixels))
}
func TexImage2D(target uint32, level int32, internalformat int32, width int32, height int32, border int32, format uint32, t_ype uint32, pixels Pointer) {
	C.wrap_glTexImage2D(C.uint(target), C.int(level), C.int(internalformat), C.int(width), C.int(height), C.int(border), C.uint(format), C.uint(t_ype), unsafe.Pointer(pixels))
}
func TexImage2DMultisample(target uint32, samples int32, internalformat int32, width int32, height int32, fixedsamplelocations bool) {
	tf1 := FALSE
	if fixedsamplelocations {
		tf1 = TRUE
	}
	C.wrap_glTexImage2DMultisample(C.uint(target), C.int(samples), C.int(internalformat), C.int(width), C.int(height), C.uchar(tf1))
}
func TexImage3D(target uint32, level int32, internalformat int32, width int32, height int32, depth int32, border int32, format uint32, t_ype uint32, pixels Pointer) {
	C.wrap_glTexImage3D(C.uint(target), C.int(level), C.int(internalformat), C.int(width), C.int(height), C.int(depth), C.int(border), C.uint(format), C.uint(t_ype), unsafe.Pointer(pixels))
}
func TexImage3DMultisample(target uint32, samples int32, internalformat int32, width int32, height int32, depth int32, fixedsamplelocations bool) {
	tf1 := FALSE
	if fixedsamplelocations {
		tf1 = TRUE
	}
	C.wrap_glTexImage3DMultisample(C.uint(target), C.int(samples), C.int(internalformat), C.int(width), C.int(height), C.int(depth), C.uchar(tf1))
}
func TexParameterIiv(target uint32, pname uint32, params *int32) {
	C.wrap_glTexParameterIiv(C.uint(target), C.uint(pname), (*C.int)(params))
}
func TexParameterIuiv(target uint32, pname uint32, params *uint32) {
	C.wrap_glTexParameterIuiv(C.uint(target), C.uint(pname), (*C.uint)(params))
}
func TexParameterf(target uint32, pname uint32, param float32) {
	C.wrap_glTexParameterf(C.uint(target), C.uint(pname), C.float(param))
}
func TexParameterfv(target uint32, pname uint32, params *float32) {
	C.wrap_glTexParameterfv(C.uint(target), C.uint(pname), (*C.float)(params))
}
func TexParameteri(target uint32, pname uint32, param int32) {
	C.wrap_glTexParameteri(C.uint(target), C.uint(pname), C.int(param))
}
func TexParameteriv(target uint32, pname uint32, params *int32) {
	C.wrap_glTexParameteriv(C.uint(target), C.uint(pname), (*C.int)(params))
}
func TexStorage1D(target uint32, levels int32, internalformat uint32, width int32) {
	C.wrap_glTexStorage1D(C.uint(target), C.int(levels), C.uint(internalformat), C.int(width))
}
func TexStorage2D(target uint32, levels int32, internalformat uint32, width int32, height int32) {
	C.wrap_glTexStorage2D(C.uint(target), C.int(levels), C.uint(internalformat), C.int(width), C.int(height))
}
func TexStorage2DMultisample(target uint32, samples int32, internalformat uint32, width int32, height int32, fixedsamplelocations bool) {
	tf1 := FALSE
	if fixedsamplelocations {
		tf1 = TRUE
	}
	C.wrap_glTexStorage2DMultisample(C.uint(target), C.int(samples), C.uint(internalformat), C.int(width), C.int(height), C.uchar(tf1))
}
func TexStorage3D(target uint32, levels int32, internalformat uint32, width int32, height int32, depth int32) {
	C.wrap_glTexStorage3D(C.uint(target), C.int(levels), C.uint(internalformat), C.int(width), C.int(height), C.int(depth))
}
func TexStorage3DMultisample(target uint32, samples int32, internalformat uint32, width int32, height int32, depth int32, fixedsamplelocations bool) {
	tf1 := FALSE
	if fixedsamplelocations {
		tf1 = TRUE
	}
	C.wrap_glTexStorage3DMultisample(C.uint(target), C.int(samples), C.uint(internalformat), C.int(width), C.int(height), C.int(depth), C.uchar(tf1))
}
func TexSubImage1D(target uint32, level int32, xoffset int32, width int32, format uint32, t_ype uint32, pixels Pointer) {
	C.wrap_glTexSubImage1D(C.uint(target), C.int(level), C.int(xoffset), C.int(width), C.uint(format), C.uint(t_ype), unsafe.Pointer(pixels))
}
func TexSubImage2D(target uint32, level int32, xoffset int32, yoffset int32, width int32, height int32, format uint32, t_ype uint32, pixels Pointer) {
	C.wrap_glTexSubImage2D(C.uint(target), C.int(level), C.int(xoffset), C.int(yoffset), C.int(width), C.int(height), C.uint(format), C.uint(t_ype), unsafe.Pointer(pixels))
}
func TexSubImage3D(target uint32, level int32, xoffset int32, yoffset int32, zoffset int32, width int32, height int32, depth int32, format uint32, t_ype uint32, pixels Pointer) {
	C.wrap_glTexSubImage3D(C.uint(target), C.int(level), C.int(xoffset), C.int(yoffset), C.int(zoffset), C.int(width), C.int(height), C.int(depth), C.uint(format), C.uint(t_ype), unsafe.Pointer(pixels))
}
func TextureBufferRangeEXT(texture uint32, target uint32, internalformat uint32, buffer uint32, offset int64, size int64) {
	C.wrap_glTextureBufferRangeEXT(C.uint(texture), C.uint(target), C.uint(internalformat), C.uint(buffer), C.longlong(offset), C.longlong(size))
}
func TextureStorage1DEXT(texture uint32, target uint32, levels int32, internalformat uint32, width int32) {
	C.wrap_glTextureStorage1DEXT(C.uint(texture), C.uint(target), C.int(levels), C.uint(internalformat), C.int(width))
}
func TextureStorage2DEXT(texture uint32, target uint32, levels int32, internalformat uint32, width int32, height int32) {
	C.wrap_glTextureStorage2DEXT(C.uint(texture), C.uint(target), C.int(levels), C.uint(internalformat), C.int(width), C.int(height))
}
func TextureStorage2DMultisampleEXT(texture uint32, target uint32, samples int32, internalformat uint32, width int32, height int32, fixedsamplelocations bool) {
	tf1 := FALSE
	if fixedsamplelocations {
		tf1 = TRUE
	}
	C.wrap_glTextureStorage2DMultisampleEXT(C.uint(texture), C.uint(target), C.int(samples), C.uint(internalformat), C.int(width), C.int(height), C.uchar(tf1))
}
func TextureStorage3DEXT(texture uint32, target uint32, levels int32, internalformat uint32, width int32, height int32, depth int32) {
	C.wrap_glTextureStorage3DEXT(C.uint(texture), C.uint(target), C.int(levels), C.uint(internalformat), C.int(width), C.int(height), C.int(depth))
}
func TextureStorage3DMultisampleEXT(texture uint32, target uint32, samples int32, internalformat uint32, width int32, height int32, depth int32, fixedsamplelocations bool) {
	tf1 := FALSE
	if fixedsamplelocations {
		tf1 = TRUE
	}
	C.wrap_glTextureStorage3DMultisampleEXT(C.uint(texture), C.uint(target), C.int(samples), C.uint(internalformat), C.int(width), C.int(height), C.int(depth), C.uchar(tf1))
}
func TextureView(texture uint32, target uint32, origtexture uint32, internalformat uint32, minlevel uint32, numlevels uint32, minlayer uint32, numlayers uint32) {
	C.wrap_glTextureView(C.uint(texture), C.uint(target), C.uint(origtexture), C.uint(internalformat), C.uint(minlevel), C.uint(numlevels), C.uint(minlayer), C.uint(numlayers))
}
func TransformFeedbackVaryings(program uint32, count int32, varyings []string, bufferMode uint32) {
	cstrings := C.newStringArray(C.int(len(varyings)))
	defer C.freeStringArray(cstrings, C.int(len(varyings)))
	for cnt, str := range varyings {
		C.assignString(cstrings, C.CString(str), C.int(cnt))
	}
	C.wrap_glTransformFeedbackVaryings(C.uint(program), C.int(count), cstrings, C.uint(bufferMode))
}
func Uniform1d(location int32, x float64) {
	C.wrap_glUniform1d(C.int(location), C.double(x))
}
func Uniform1dv(location int32, count int32, value *float64) {
	C.wrap_glUniform1dv(C.int(location), C.int(count), (*C.double)(value))
}
func Uniform1f(location int32, v0 float32) {
	C.wrap_glUniform1f(C.int(location), C.float(v0))
}
func Uniform1fv(location int32, count int32, value *float32) {
	C.wrap_glUniform1fv(C.int(location), C.int(count), (*C.float)(value))
}
func Uniform1i(location int32, v0 int32) {
	C.wrap_glUniform1i(C.int(location), C.int(v0))
}
func Uniform1iv(location int32, count int32, value *int32) {
	C.wrap_glUniform1iv(C.int(location), C.int(count), (*C.int)(value))
}
func Uniform1ui(location int32, v0 uint32) {
	C.wrap_glUniform1ui(C.int(location), C.uint(v0))
}
func Uniform1uiv(location int32, count int32, value *uint32) {
	C.wrap_glUniform1uiv(C.int(location), C.int(count), (*C.uint)(value))
}
func Uniform2d(location int32, x float64, y float64) {
	C.wrap_glUniform2d(C.int(location), C.double(x), C.double(y))
}
func Uniform2dv(location int32, count int32, value *float64) {
	C.wrap_glUniform2dv(C.int(location), C.int(count), (*C.double)(value))
}
func Uniform2f(location int32, v0 float32, v1 float32) {
	C.wrap_glUniform2f(C.int(location), C.float(v0), C.float(v1))
}
func Uniform2fv(location int32, count int32, value *float32) {
	C.wrap_glUniform2fv(C.int(location), C.int(count), (*C.float)(value))
}
func Uniform2i(location int32, v0 int32, v1 int32) {
	C.wrap_glUniform2i(C.int(location), C.int(v0), C.int(v1))
}
func Uniform2iv(location int32, count int32, value *int32) {
	C.wrap_glUniform2iv(C.int(location), C.int(count), (*C.int)(value))
}
func Uniform2ui(location int32, v0 uint32, v1 uint32) {
	C.wrap_glUniform2ui(C.int(location), C.uint(v0), C.uint(v1))
}
func Uniform2uiv(location int32, count int32, value *uint32) {
	C.wrap_glUniform2uiv(C.int(location), C.int(count), (*C.uint)(value))
}
func Uniform3d(location int32, x float64, y float64, z float64) {
	C.wrap_glUniform3d(C.int(location), C.double(x), C.double(y), C.double(z))
}
func Uniform3dv(location int32, count int32, value *float64) {
	C.wrap_glUniform3dv(C.int(location), C.int(count), (*C.double)(value))
}
func Uniform3f(location int32, v0 float32, v1 float32, v2 float32) {
	C.wrap_glUniform3f(C.int(location), C.float(v0), C.float(v1), C.float(v2))
}
func Uniform3fv(location int32, count int32, value *float32) {
	C.wrap_glUniform3fv(C.int(location), C.int(count), (*C.float)(value))
}
func Uniform3i(location int32, v0 int32, v1 int32, v2 int32) {
	C.wrap_glUniform3i(C.int(location), C.int(v0), C.int(v1), C.int(v2))
}
func Uniform3iv(location int32, count int32, value *int32) {
	C.wrap_glUniform3iv(C.int(location), C.int(count), (*C.int)(value))
}
func Uniform3ui(location int32, v0 uint32, v1 uint32, v2 uint32) {
	C.wrap_glUniform3ui(C.int(location), C.uint(v0), C.uint(v1), C.uint(v2))
}
func Uniform3uiv(location int32, count int32, value *uint32) {
	C.wrap_glUniform3uiv(C.int(location), C.int(count), (*C.uint)(value))
}
func Uniform4d(location int32, x float64, y float64, z float64, w float64) {
	C.wrap_glUniform4d(C.int(location), C.double(x), C.double(y), C.double(z), C.double(w))
}
func Uniform4dv(location int32, count int32, value *float64) {
	C.wrap_glUniform4dv(C.int(location), C.int(count), (*C.double)(value))
}
func Uniform4f(location int32, v0 float32, v1 float32, v2 float32, v3 float32) {
	C.wrap_glUniform4f(C.int(location), C.float(v0), C.float(v1), C.float(v2), C.float(v3))
}
func Uniform4fv(location int32, count int32, value *float32) {
	C.wrap_glUniform4fv(C.int(location), C.int(count), (*C.float)(value))
}
func Uniform4i(location int32, v0 int32, v1 int32, v2 int32, v3 int32) {
	C.wrap_glUniform4i(C.int(location), C.int(v0), C.int(v1), C.int(v2), C.int(v3))
}
func Uniform4iv(location int32, count int32, value *int32) {
	C.wrap_glUniform4iv(C.int(location), C.int(count), (*C.int)(value))
}
func Uniform4ui(location int32, v0 uint32, v1 uint32, v2 uint32, v3 uint32) {
	C.wrap_glUniform4ui(C.int(location), C.uint(v0), C.uint(v1), C.uint(v2), C.uint(v3))
}
func Uniform4uiv(location int32, count int32, value *uint32) {
	C.wrap_glUniform4uiv(C.int(location), C.int(count), (*C.uint)(value))
}
func UniformBlockBinding(program uint32, uniformBlockIndex uint32, uniformBlockBinding uint32) {
	C.wrap_glUniformBlockBinding(C.uint(program), C.uint(uniformBlockIndex), C.uint(uniformBlockBinding))
}
func UniformMatrix2dv(location int32, count int32, transpose bool, value *float64) {
	tf1 := FALSE
	if transpose {
		tf1 = TRUE
	}
	C.wrap_glUniformMatrix2dv(C.int(location), C.int(count), C.uchar(tf1), (*C.double)(value))
}
func UniformMatrix2fv(location int32, count int32, transpose bool, value *float32) {
	tf1 := FALSE
	if transpose {
		tf1 = TRUE
	}
	C.wrap_glUniformMatrix2fv(C.int(location), C.int(count), C.uchar(tf1), (*C.float)(value))
}
func UniformMatrix2x3dv(location int32, count int32, transpose bool, value *float64) {
	tf1 := FALSE
	if transpose {
		tf1 = TRUE
	}
	C.wrap_glUniformMatrix2x3dv(C.int(location), C.int(count), C.uchar(tf1), (*C.double)(value))
}
func UniformMatrix2x3fv(location int32, count int32, transpose bool, value *float32) {
	tf1 := FALSE
	if transpose {
		tf1 = TRUE
	}
	C.wrap_glUniformMatrix2x3fv(C.int(location), C.int(count), C.uchar(tf1), (*C.float)(value))
}
func UniformMatrix2x4dv(location int32, count int32, transpose bool, value *float64) {
	tf1 := FALSE
	if transpose {
		tf1 = TRUE
	}
	C.wrap_glUniformMatrix2x4dv(C.int(location), C.int(count), C.uchar(tf1), (*C.double)(value))
}
func UniformMatrix2x4fv(location int32, count int32, transpose bool, value *float32) {
	tf1 := FALSE
	if transpose {
		tf1 = TRUE
	}
	C.wrap_glUniformMatrix2x4fv(C.int(location), C.int(count), C.uchar(tf1), (*C.float)(value))
}
func UniformMatrix3dv(location int32, count int32, transpose bool, value *float64) {
	tf1 := FALSE
	if transpose {
		tf1 = TRUE
	}
	C.wrap_glUniformMatrix3dv(C.int(location), C.int(count), C.uchar(tf1), (*C.double)(value))
}
func UniformMatrix3fv(location int32, count int32, transpose bool, value *float32) {
	tf1 := FALSE
	if transpose {
		tf1 = TRUE
	}
	C.wrap_glUniformMatrix3fv(C.int(location), C.int(count), C.uchar(tf1), (*C.float)(value))
}
func UniformMatrix3x2dv(location int32, count int32, transpose bool, value *float64) {
	tf1 := FALSE
	if transpose {
		tf1 = TRUE
	}
	C.wrap_glUniformMatrix3x2dv(C.int(location), C.int(count), C.uchar(tf1), (*C.double)(value))
}
func UniformMatrix3x2fv(location int32, count int32, transpose bool, value *float32) {
	tf1 := FALSE
	if transpose {
		tf1 = TRUE
	}
	C.wrap_glUniformMatrix3x2fv(C.int(location), C.int(count), C.uchar(tf1), (*C.float)(value))
}
func UniformMatrix3x4dv(location int32, count int32, transpose bool, value *float64) {
	tf1 := FALSE
	if transpose {
		tf1 = TRUE
	}
	C.wrap_glUniformMatrix3x4dv(C.int(location), C.int(count), C.uchar(tf1), (*C.double)(value))
}
func UniformMatrix3x4fv(location int32, count int32, transpose bool, value *float32) {
	tf1 := FALSE
	if transpose {
		tf1 = TRUE
	}
	C.wrap_glUniformMatrix3x4fv(C.int(location), C.int(count), C.uchar(tf1), (*C.float)(value))
}
func UniformMatrix4dv(location int32, count int32, transpose bool, value *float64) {
	tf1 := FALSE
	if transpose {
		tf1 = TRUE
	}
	C.wrap_glUniformMatrix4dv(C.int(location), C.int(count), C.uchar(tf1), (*C.double)(value))
}
func UniformMatrix4fv(location int32, count int32, transpose bool, value *float32) {
	tf1 := FALSE
	if transpose {
		tf1 = TRUE
	}
	C.wrap_glUniformMatrix4fv(C.int(location), C.int(count), C.uchar(tf1), (*C.float)(value))
}
func UniformMatrix4x2dv(location int32, count int32, transpose bool, value *float64) {
	tf1 := FALSE
	if transpose {
		tf1 = TRUE
	}
	C.wrap_glUniformMatrix4x2dv(C.int(location), C.int(count), C.uchar(tf1), (*C.double)(value))
}
func UniformMatrix4x2fv(location int32, count int32, transpose bool, value *float32) {
	tf1 := FALSE
	if transpose {
		tf1 = TRUE
	}
	C.wrap_glUniformMatrix4x2fv(C.int(location), C.int(count), C.uchar(tf1), (*C.float)(value))
}
func UniformMatrix4x3dv(location int32, count int32, transpose bool, value *float64) {
	tf1 := FALSE
	if transpose {
		tf1 = TRUE
	}
	C.wrap_glUniformMatrix4x3dv(C.int(location), C.int(count), C.uchar(tf1), (*C.double)(value))
}
func UniformMatrix4x3fv(location int32, count int32, transpose bool, value *float32) {
	tf1 := FALSE
	if transpose {
		tf1 = TRUE
	}
	C.wrap_glUniformMatrix4x3fv(C.int(location), C.int(count), C.uchar(tf1), (*C.float)(value))
}
func UniformSubroutinesuiv(shadertype uint32, count int32, indices *uint32) {
	C.wrap_glUniformSubroutinesuiv(C.uint(shadertype), C.int(count), (*C.uint)(indices))
}
func UseProgram(program uint32) {
	C.wrap_glUseProgram(C.uint(program))
}
func UseProgramStages(pipeline uint32, stages uint32, program uint32) {
	C.wrap_glUseProgramStages(C.uint(pipeline), C.uint(stages), C.uint(program))
}
func ValidateProgram(program uint32) {
	C.wrap_glValidateProgram(C.uint(program))
}
func ValidateProgramPipeline(pipeline uint32) {
	C.wrap_glValidateProgramPipeline(C.uint(pipeline))
}
func VertexArrayBindVertexBufferEXT(vaobj uint32, bindingindex uint32, buffer uint32, offset int64, stride int32) {
	C.wrap_glVertexArrayBindVertexBufferEXT(C.uint(vaobj), C.uint(bindingindex), C.uint(buffer), C.longlong(offset), C.int(stride))
}
func VertexArrayVertexAttribBindingEXT(vaobj uint32, attribindex uint32, bindingindex uint32) {
	C.wrap_glVertexArrayVertexAttribBindingEXT(C.uint(vaobj), C.uint(attribindex), C.uint(bindingindex))
}
func VertexArrayVertexAttribFormatEXT(vaobj uint32, attribindex uint32, size int32, t_ype uint32, normalized bool, relativeoffset uint32) {
	tf1 := FALSE
	if normalized {
		tf1 = TRUE
	}
	C.wrap_glVertexArrayVertexAttribFormatEXT(C.uint(vaobj), C.uint(attribindex), C.int(size), C.uint(t_ype), C.uchar(tf1), C.uint(relativeoffset))
}
func VertexArrayVertexAttribIFormatEXT(vaobj uint32, attribindex uint32, size int32, t_ype uint32, relativeoffset uint32) {
	C.wrap_glVertexArrayVertexAttribIFormatEXT(C.uint(vaobj), C.uint(attribindex), C.int(size), C.uint(t_ype), C.uint(relativeoffset))
}
func VertexArrayVertexAttribLFormatEXT(vaobj uint32, attribindex uint32, size int32, t_ype uint32, relativeoffset uint32) {
	C.wrap_glVertexArrayVertexAttribLFormatEXT(C.uint(vaobj), C.uint(attribindex), C.int(size), C.uint(t_ype), C.uint(relativeoffset))
}
func VertexArrayVertexBindingDivisorEXT(vaobj uint32, bindingindex uint32, divisor uint32) {
	C.wrap_glVertexArrayVertexBindingDivisorEXT(C.uint(vaobj), C.uint(bindingindex), C.uint(divisor))
}
func VertexAttrib1d(index uint32, x float64) {
	C.wrap_glVertexAttrib1d(C.uint(index), C.double(x))
}
func VertexAttrib1dv(index uint32, v *float64) {
	C.wrap_glVertexAttrib1dv(C.uint(index), (*C.double)(v))
}
func VertexAttrib1f(index uint32, x float32) {
	C.wrap_glVertexAttrib1f(C.uint(index), C.float(x))
}
func VertexAttrib1fv(index uint32, v *float32) {
	C.wrap_glVertexAttrib1fv(C.uint(index), (*C.float)(v))
}
func VertexAttrib1s(index uint32, x int16) {
	C.wrap_glVertexAttrib1s(C.uint(index), C.short(x))
}
func VertexAttrib1sv(index uint32, v *int16) {
	C.wrap_glVertexAttrib1sv(C.uint(index), (*C.short)(v))
}
func VertexAttrib2d(index uint32, x float64, y float64) {
	C.wrap_glVertexAttrib2d(C.uint(index), C.double(x), C.double(y))
}
func VertexAttrib2dv(index uint32, v *float64) {
	C.wrap_glVertexAttrib2dv(C.uint(index), (*C.double)(v))
}
func VertexAttrib2f(index uint32, x float32, y float32) {
	C.wrap_glVertexAttrib2f(C.uint(index), C.float(x), C.float(y))
}
func VertexAttrib2fv(index uint32, v *float32) {
	C.wrap_glVertexAttrib2fv(C.uint(index), (*C.float)(v))
}
func VertexAttrib2s(index uint32, x int16, y int16) {
	C.wrap_glVertexAttrib2s(C.uint(index), C.short(x), C.short(y))
}
func VertexAttrib2sv(index uint32, v *int16) {
	C.wrap_glVertexAttrib2sv(C.uint(index), (*C.short)(v))
}
func VertexAttrib3d(index uint32, x float64, y float64, z float64) {
	C.wrap_glVertexAttrib3d(C.uint(index), C.double(x), C.double(y), C.double(z))
}
func VertexAttrib3dv(index uint32, v *float64) {
	C.wrap_glVertexAttrib3dv(C.uint(index), (*C.double)(v))
}
func VertexAttrib3f(index uint32, x float32, y float32, z float32) {
	C.wrap_glVertexAttrib3f(C.uint(index), C.float(x), C.float(y), C.float(z))
}
func VertexAttrib3fv(index uint32, v *float32) {
	C.wrap_glVertexAttrib3fv(C.uint(index), (*C.float)(v))
}
func VertexAttrib3s(index uint32, x int16, y int16, z int16) {
	C.wrap_glVertexAttrib3s(C.uint(index), C.short(x), C.short(y), C.short(z))
}
func VertexAttrib3sv(index uint32, v *int16) {
	C.wrap_glVertexAttrib3sv(C.uint(index), (*C.short)(v))
}
func VertexAttrib4Nbv(index uint32, v *int8) {
	C.wrap_glVertexAttrib4Nbv(C.uint(index), (*C.schar)(v))
}
func VertexAttrib4Niv(index uint32, v *int32) {
	C.wrap_glVertexAttrib4Niv(C.uint(index), (*C.int)(v))
}
func VertexAttrib4Nsv(index uint32, v *int16) {
	C.wrap_glVertexAttrib4Nsv(C.uint(index), (*C.short)(v))
}
func VertexAttrib4Nub(index uint32, x uint8, y uint8, z uint8, w uint8) {
	C.wrap_glVertexAttrib4Nub(C.uint(index), C.uchar(x), C.uchar(y), C.uchar(z), C.uchar(w))
}
func VertexAttrib4Nubv(index uint32, v *uint8) {
	C.wrap_glVertexAttrib4Nubv(C.uint(index), (*C.uchar)(v))
}
func VertexAttrib4Nuiv(index uint32, v *uint32) {
	C.wrap_glVertexAttrib4Nuiv(C.uint(index), (*C.uint)(v))
}
func VertexAttrib4Nusv(index uint32, v *uint16) {
	C.wrap_glVertexAttrib4Nusv(C.uint(index), (*C.ushort)(v))
}
func VertexAttrib4bv(index uint32, v *int8) {
	C.wrap_glVertexAttrib4bv(C.uint(index), (*C.schar)(v))
}
func VertexAttrib4d(index uint32, x float64, y float64, z float64, w float64) {
	C.wrap_glVertexAttrib4d(C.uint(index), C.double(x), C.double(y), C.double(z), C.double(w))
}
func VertexAttrib4dv(index uint32, v *float64) {
	C.wrap_glVertexAttrib4dv(C.uint(index), (*C.double)(v))
}
func VertexAttrib4f(index uint32, x float32, y float32, z float32, w float32) {
	C.wrap_glVertexAttrib4f(C.uint(index), C.float(x), C.float(y), C.float(z), C.float(w))
}
func VertexAttrib4fv(index uint32, v *float32) {
	C.wrap_glVertexAttrib4fv(C.uint(index), (*C.float)(v))
}
func VertexAttrib4iv(index uint32, v *int32) {
	C.wrap_glVertexAttrib4iv(C.uint(index), (*C.int)(v))
}
func VertexAttrib4s(index uint32, x int16, y int16, z int16, w int16) {
	C.wrap_glVertexAttrib4s(C.uint(index), C.short(x), C.short(y), C.short(z), C.short(w))
}
func VertexAttrib4sv(index uint32, v *int16) {
	C.wrap_glVertexAttrib4sv(C.uint(index), (*C.short)(v))
}
func VertexAttrib4ubv(index uint32, v *uint8) {
	C.wrap_glVertexAttrib4ubv(C.uint(index), (*C.uchar)(v))
}
func VertexAttrib4uiv(index uint32, v *uint32) {
	C.wrap_glVertexAttrib4uiv(C.uint(index), (*C.uint)(v))
}
func VertexAttrib4usv(index uint32, v *uint16) {
	C.wrap_glVertexAttrib4usv(C.uint(index), (*C.ushort)(v))
}
func VertexAttribBinding(attribindex uint32, bindingindex uint32) {
	C.wrap_glVertexAttribBinding(C.uint(attribindex), C.uint(bindingindex))
}
func VertexAttribDivisor(index uint32, divisor uint32) {
	C.wrap_glVertexAttribDivisor(C.uint(index), C.uint(divisor))
}
func VertexAttribFormat(attribindex uint32, size int32, t_ype uint32, normalized bool, relativeoffset uint32) {
	tf1 := FALSE
	if normalized {
		tf1 = TRUE
	}
	C.wrap_glVertexAttribFormat(C.uint(attribindex), C.int(size), C.uint(t_ype), C.uchar(tf1), C.uint(relativeoffset))
}
func VertexAttribI1i(index uint32, x int32) {
	C.wrap_glVertexAttribI1i(C.uint(index), C.int(x))
}
func VertexAttribI1iv(index uint32, v *int32) {
	C.wrap_glVertexAttribI1iv(C.uint(index), (*C.int)(v))
}
func VertexAttribI1ui(index uint32, x uint32) {
	C.wrap_glVertexAttribI1ui(C.uint(index), C.uint(x))
}
func VertexAttribI1uiv(index uint32, v *uint32) {
	C.wrap_glVertexAttribI1uiv(C.uint(index), (*C.uint)(v))
}
func VertexAttribI2i(index uint32, x int32, y int32) {
	C.wrap_glVertexAttribI2i(C.uint(index), C.int(x), C.int(y))
}
func VertexAttribI2iv(index uint32, v *int32) {
	C.wrap_glVertexAttribI2iv(C.uint(index), (*C.int)(v))
}
func VertexAttribI2ui(index uint32, x uint32, y uint32) {
	C.wrap_glVertexAttribI2ui(C.uint(index), C.uint(x), C.uint(y))
}
func VertexAttribI2uiv(index uint32, v *uint32) {
	C.wrap_glVertexAttribI2uiv(C.uint(index), (*C.uint)(v))
}
func VertexAttribI3i(index uint32, x int32, y int32, z int32) {
	C.wrap_glVertexAttribI3i(C.uint(index), C.int(x), C.int(y), C.int(z))
}
func VertexAttribI3iv(index uint32, v *int32) {
	C.wrap_glVertexAttribI3iv(C.uint(index), (*C.int)(v))
}
func VertexAttribI3ui(index uint32, x uint32, y uint32, z uint32) {
	C.wrap_glVertexAttribI3ui(C.uint(index), C.uint(x), C.uint(y), C.uint(z))
}
func VertexAttribI3uiv(index uint32, v *uint32) {
	C.wrap_glVertexAttribI3uiv(C.uint(index), (*C.uint)(v))
}
func VertexAttribI4bv(index uint32, v *int8) {
	C.wrap_glVertexAttribI4bv(C.uint(index), (*C.schar)(v))
}
func VertexAttribI4i(index uint32, x int32, y int32, z int32, w int32) {
	C.wrap_glVertexAttribI4i(C.uint(index), C.int(x), C.int(y), C.int(z), C.int(w))
}
func VertexAttribI4iv(index uint32, v *int32) {
	C.wrap_glVertexAttribI4iv(C.uint(index), (*C.int)(v))
}
func VertexAttribI4sv(index uint32, v *int16) {
	C.wrap_glVertexAttribI4sv(C.uint(index), (*C.short)(v))
}
func VertexAttribI4ubv(index uint32, v *uint8) {
	C.wrap_glVertexAttribI4ubv(C.uint(index), (*C.uchar)(v))
}
func VertexAttribI4ui(index uint32, x uint32, y uint32, z uint32, w uint32) {
	C.wrap_glVertexAttribI4ui(C.uint(index), C.uint(x), C.uint(y), C.uint(z), C.uint(w))
}
func VertexAttribI4uiv(index uint32, v *uint32) {
	C.wrap_glVertexAttribI4uiv(C.uint(index), (*C.uint)(v))
}
func VertexAttribI4usv(index uint32, v *uint16) {
	C.wrap_glVertexAttribI4usv(C.uint(index), (*C.ushort)(v))
}
func VertexAttribIFormat(attribindex uint32, size int32, t_ype uint32, relativeoffset uint32) {
	C.wrap_glVertexAttribIFormat(C.uint(attribindex), C.int(size), C.uint(t_ype), C.uint(relativeoffset))
}
func VertexAttribIPointer(index uint32, size int32, t_ype uint32, stride int32, pointer Pointer) {
	C.wrap_glVertexAttribIPointer(C.uint(index), C.int(size), C.uint(t_ype), C.int(stride), unsafe.Pointer(pointer))
}
func VertexAttribL1d(index uint32, x float64) {
	C.wrap_glVertexAttribL1d(C.uint(index), C.double(x))
}
func VertexAttribL1dv(index uint32, v *float64) {
	C.wrap_glVertexAttribL1dv(C.uint(index), (*C.double)(v))
}
func VertexAttribL2d(index uint32, x float64, y float64) {
	C.wrap_glVertexAttribL2d(C.uint(index), C.double(x), C.double(y))
}
func VertexAttribL2dv(index uint32, v *float64) {
	C.wrap_glVertexAttribL2dv(C.uint(index), (*C.double)(v))
}
func VertexAttribL3d(index uint32, x float64, y float64, z float64) {
	C.wrap_glVertexAttribL3d(C.uint(index), C.double(x), C.double(y), C.double(z))
}
func VertexAttribL3dv(index uint32, v *float64) {
	C.wrap_glVertexAttribL3dv(C.uint(index), (*C.double)(v))
}
func VertexAttribL4d(index uint32, x float64, y float64, z float64, w float64) {
	C.wrap_glVertexAttribL4d(C.uint(index), C.double(x), C.double(y), C.double(z), C.double(w))
}
func VertexAttribL4dv(index uint32, v *float64) {
	C.wrap_glVertexAttribL4dv(C.uint(index), (*C.double)(v))
}
func VertexAttribLFormat(attribindex uint32, size int32, t_ype uint32, relativeoffset uint32) {
	C.wrap_glVertexAttribLFormat(C.uint(attribindex), C.int(size), C.uint(t_ype), C.uint(relativeoffset))
}
func VertexAttribLPointer(index uint32, size int32, t_ype uint32, stride int32, pointer Pointer) {
	C.wrap_glVertexAttribLPointer(C.uint(index), C.int(size), C.uint(t_ype), C.int(stride), unsafe.Pointer(pointer))
}
func VertexAttribP1ui(index uint32, t_ype uint32, normalized bool, value uint32) {
	tf1 := FALSE
	if normalized {
		tf1 = TRUE
	}
	C.wrap_glVertexAttribP1ui(C.uint(index), C.uint(t_ype), C.uchar(tf1), C.uint(value))
}
func VertexAttribP1uiv(index uint32, t_ype uint32, normalized bool, value *uint32) {
	tf1 := FALSE
	if normalized {
		tf1 = TRUE
	}
	C.wrap_glVertexAttribP1uiv(C.uint(index), C.uint(t_ype), C.uchar(tf1), (*C.uint)(value))
}
func VertexAttribP2ui(index uint32, t_ype uint32, normalized bool, value uint32) {
	tf1 := FALSE
	if normalized {
		tf1 = TRUE
	}
	C.wrap_glVertexAttribP2ui(C.uint(index), C.uint(t_ype), C.uchar(tf1), C.uint(value))
}
func VertexAttribP2uiv(index uint32, t_ype uint32, normalized bool, value *uint32) {
	tf1 := FALSE
	if normalized {
		tf1 = TRUE
	}
	C.wrap_glVertexAttribP2uiv(C.uint(index), C.uint(t_ype), C.uchar(tf1), (*C.uint)(value))
}
func VertexAttribP3ui(index uint32, t_ype uint32, normalized bool, value uint32) {
	tf1 := FALSE
	if normalized {
		tf1 = TRUE
	}
	C.wrap_glVertexAttribP3ui(C.uint(index), C.uint(t_ype), C.uchar(tf1), C.uint(value))
}
func VertexAttribP3uiv(index uint32, t_ype uint32, normalized bool, value *uint32) {
	tf1 := FALSE
	if normalized {
		tf1 = TRUE
	}
	C.wrap_glVertexAttribP3uiv(C.uint(index), C.uint(t_ype), C.uchar(tf1), (*C.uint)(value))
}
func VertexAttribP4ui(index uint32, t_ype uint32, normalized bool, value uint32) {
	tf1 := FALSE
	if normalized {
		tf1 = TRUE
	}
	C.wrap_glVertexAttribP4ui(C.uint(index), C.uint(t_ype), C.uchar(tf1), C.uint(value))
}
func VertexAttribP4uiv(index uint32, t_ype uint32, normalized bool, value *uint32) {
	tf1 := FALSE
	if normalized {
		tf1 = TRUE
	}
	C.wrap_glVertexAttribP4uiv(C.uint(index), C.uint(t_ype), C.uchar(tf1), (*C.uint)(value))
}
func VertexAttribPointer(index uint32, size int32, t_ype uint32, normalized bool, stride int32, pointer int64) {
	tf1 := FALSE
	if normalized {
		tf1 = TRUE
	}
	C.wrap_glVertexAttribPointer(C.uint(index), C.int(size), C.uint(t_ype), C.uchar(tf1), C.int(stride), C.longlong(pointer))
}
func VertexBindingDivisor(bindingindex uint32, divisor uint32) {
	C.wrap_glVertexBindingDivisor(C.uint(bindingindex), C.uint(divisor))
}
func VertexP2ui(t_ype uint32, value uint32) {
	C.wrap_glVertexP2ui(C.uint(t_ype), C.uint(value))
}
func VertexP2uiv(t_ype uint32, value *uint32) {
	C.wrap_glVertexP2uiv(C.uint(t_ype), (*C.uint)(value))
}
func VertexP3ui(t_ype uint32, value uint32) {
	C.wrap_glVertexP3ui(C.uint(t_ype), C.uint(value))
}
func VertexP3uiv(t_ype uint32, value *uint32) {
	C.wrap_glVertexP3uiv(C.uint(t_ype), (*C.uint)(value))
}
func VertexP4ui(t_ype uint32, value uint32) {
	C.wrap_glVertexP4ui(C.uint(t_ype), C.uint(value))
}
func VertexP4uiv(t_ype uint32, value *uint32) {
	C.wrap_glVertexP4uiv(C.uint(t_ype), (*C.uint)(value))
}
func Viewport(x int32, y int32, width int32, height int32) {
	C.wrap_glViewport(C.int(x), C.int(y), C.int(width), C.int(height))
}
func ViewportArrayv(first uint32, count int32, v *float32) {
	C.wrap_glViewportArrayv(C.uint(first), C.int(count), (*C.float)(v))
}
func ViewportIndexedf(index uint32, x float32, y float32, w float32, h float32) {
	C.wrap_glViewportIndexedf(C.uint(index), C.float(x), C.float(y), C.float(w), C.float(h))
}
func ViewportIndexedfv(index uint32, v *float32) {
	C.wrap_glViewportIndexedfv(C.uint(index), (*C.float)(v))
}
func WaitSync(sync Sync, flags uint32, timeout uint64) {
	C.wrap_glWaitSync(C.GLsync(sync), C.uint(flags), C.ulonglong(timeout))
}

// Show which function pointers are bound
func BindingReport() (report []string) {
	report = []string{}
	report = append(report, isBound(unsafe.Pointer(C.pfn_glIsBuffer), "glIsBuffer"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glIsEnabled), "glIsEnabled"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glIsEnabledi), "glIsEnabledi"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glIsFramebuffer), "glIsFramebuffer"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glIsNamedStringARB), "glIsNamedStringARB"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glIsProgram), "glIsProgram"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glIsProgramPipeline), "glIsProgramPipeline"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glIsQuery), "glIsQuery"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glIsRenderbuffer), "glIsRenderbuffer"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glIsSampler), "glIsSampler"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glIsShader), "glIsShader"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glIsSync), "glIsSync"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glIsTexture), "glIsTexture"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glIsTransformFeedback), "glIsTransformFeedback"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glIsVertexArray), "glIsVertexArray"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUnmapBuffer), "glUnmapBuffer"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glCheckFramebufferStatus), "glCheckFramebufferStatus"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glClientWaitSync), "glClientWaitSync"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetError), "glGetError"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetGraphicsResetStatusARB), "glGetGraphicsResetStatusARB"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetAttribLocation), "glGetAttribLocation"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetFragDataIndex), "glGetFragDataIndex"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetFragDataLocation), "glGetFragDataLocation"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetProgramResourceLocation), "glGetProgramResourceLocation"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetProgramResourceLocationIndex), "glGetProgramResourceLocationIndex"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetSubroutineUniformLocation), "glGetSubroutineUniformLocation"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetUniformLocation), "glGetUniformLocation"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glCreateSyncFromCLeventARB), "glCreateSyncFromCLeventARB"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glFenceSync), "glFenceSync"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glCreateProgram), "glCreateProgram"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glCreateShader), "glCreateShader"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glCreateShaderProgramv), "glCreateShaderProgramv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetDebugMessageLog), "glGetDebugMessageLog"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetDebugMessageLogARB), "glGetDebugMessageLogARB"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetProgramResourceIndex), "glGetProgramResourceIndex"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetSubroutineIndex), "glGetSubroutineIndex"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetUniformBlockIndex), "glGetUniformBlockIndex"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glMapBuffer), "glMapBuffer"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glMapBufferRange), "glMapBufferRange"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetString), "glGetString"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetStringi), "glGetStringi"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glActiveShaderProgram), "glActiveShaderProgram"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glActiveTexture), "glActiveTexture"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glAttachShader), "glAttachShader"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glBeginConditionalRender), "glBeginConditionalRender"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glBeginQuery), "glBeginQuery"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glBeginQueryIndexed), "glBeginQueryIndexed"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glBeginTransformFeedback), "glBeginTransformFeedback"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glBindAttribLocation), "glBindAttribLocation"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glBindBuffer), "glBindBuffer"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glBindBufferBase), "glBindBufferBase"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glBindBufferRange), "glBindBufferRange"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glBindFragDataLocation), "glBindFragDataLocation"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glBindFragDataLocationIndexed), "glBindFragDataLocationIndexed"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glBindFramebuffer), "glBindFramebuffer"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glBindImageTexture), "glBindImageTexture"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glBindProgramPipeline), "glBindProgramPipeline"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glBindRenderbuffer), "glBindRenderbuffer"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glBindSampler), "glBindSampler"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glBindTexture), "glBindTexture"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glBindTransformFeedback), "glBindTransformFeedback"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glBindVertexArray), "glBindVertexArray"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glBindVertexBuffer), "glBindVertexBuffer"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glBlendColor), "glBlendColor"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glBlendEquation), "glBlendEquation"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glBlendEquationSeparate), "glBlendEquationSeparate"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glBlendEquationSeparatei), "glBlendEquationSeparatei"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glBlendEquationSeparateiARB), "glBlendEquationSeparateiARB"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glBlendEquationi), "glBlendEquationi"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glBlendEquationiARB), "glBlendEquationiARB"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glBlendFunc), "glBlendFunc"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glBlendFuncSeparate), "glBlendFuncSeparate"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glBlendFuncSeparatei), "glBlendFuncSeparatei"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glBlendFuncSeparateiARB), "glBlendFuncSeparateiARB"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glBlendFunci), "glBlendFunci"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glBlendFunciARB), "glBlendFunciARB"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glBlitFramebuffer), "glBlitFramebuffer"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glBufferData), "glBufferData"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glBufferSubData), "glBufferSubData"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glClampColor), "glClampColor"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glClear), "glClear"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glClearBufferData), "glClearBufferData"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glClearBufferSubData), "glClearBufferSubData"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glClearBufferfi), "glClearBufferfi"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glClearBufferfv), "glClearBufferfv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glClearBufferiv), "glClearBufferiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glClearBufferuiv), "glClearBufferuiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glClearColor), "glClearColor"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glClearDepth), "glClearDepth"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glClearDepthf), "glClearDepthf"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glClearNamedBufferDataEXT), "glClearNamedBufferDataEXT"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glClearNamedBufferSubDataEXT), "glClearNamedBufferSubDataEXT"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glClearStencil), "glClearStencil"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glColorMask), "glColorMask"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glColorMaski), "glColorMaski"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glColorP3ui), "glColorP3ui"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glColorP3uiv), "glColorP3uiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glColorP4ui), "glColorP4ui"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glColorP4uiv), "glColorP4uiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glCompileShader), "glCompileShader"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glCompileShaderIncludeARB), "glCompileShaderIncludeARB"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glCompressedTexImage1D), "glCompressedTexImage1D"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glCompressedTexImage2D), "glCompressedTexImage2D"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glCompressedTexImage3D), "glCompressedTexImage3D"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glCompressedTexSubImage1D), "glCompressedTexSubImage1D"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glCompressedTexSubImage2D), "glCompressedTexSubImage2D"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glCompressedTexSubImage3D), "glCompressedTexSubImage3D"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glCopyBufferSubData), "glCopyBufferSubData"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glCopyImageSubData), "glCopyImageSubData"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glCopyTexImage1D), "glCopyTexImage1D"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glCopyTexImage2D), "glCopyTexImage2D"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glCopyTexSubImage1D), "glCopyTexSubImage1D"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glCopyTexSubImage2D), "glCopyTexSubImage2D"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glCopyTexSubImage3D), "glCopyTexSubImage3D"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glCullFace), "glCullFace"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDebugMessageCallback), "glDebugMessageCallback"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDebugMessageCallbackARB), "glDebugMessageCallbackARB"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDebugMessageControl), "glDebugMessageControl"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDebugMessageControlARB), "glDebugMessageControlARB"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDebugMessageInsert), "glDebugMessageInsert"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDebugMessageInsertARB), "glDebugMessageInsertARB"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDeleteBuffers), "glDeleteBuffers"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDeleteFramebuffers), "glDeleteFramebuffers"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDeleteNamedStringARB), "glDeleteNamedStringARB"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDeleteProgram), "glDeleteProgram"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDeleteProgramPipelines), "glDeleteProgramPipelines"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDeleteQueries), "glDeleteQueries"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDeleteRenderbuffers), "glDeleteRenderbuffers"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDeleteSamplers), "glDeleteSamplers"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDeleteShader), "glDeleteShader"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDeleteSync), "glDeleteSync"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDeleteTextures), "glDeleteTextures"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDeleteTransformFeedbacks), "glDeleteTransformFeedbacks"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDeleteVertexArrays), "glDeleteVertexArrays"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDepthFunc), "glDepthFunc"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDepthMask), "glDepthMask"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDepthRange), "glDepthRange"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDepthRangeArrayv), "glDepthRangeArrayv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDepthRangeIndexed), "glDepthRangeIndexed"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDepthRangef), "glDepthRangef"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDetachShader), "glDetachShader"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDisable), "glDisable"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDisableVertexAttribArray), "glDisableVertexAttribArray"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDisablei), "glDisablei"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDispatchCompute), "glDispatchCompute"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDispatchComputeIndirect), "glDispatchComputeIndirect"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDrawArrays), "glDrawArrays"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDrawArraysIndirect), "glDrawArraysIndirect"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDrawArraysInstanced), "glDrawArraysInstanced"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDrawArraysInstancedBaseInstance), "glDrawArraysInstancedBaseInstance"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDrawBuffer), "glDrawBuffer"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDrawBuffers), "glDrawBuffers"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDrawElements), "glDrawElements"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDrawElementsBaseVertex), "glDrawElementsBaseVertex"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDrawElementsIndirect), "glDrawElementsIndirect"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDrawElementsInstanced), "glDrawElementsInstanced"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDrawElementsInstancedBaseInstance), "glDrawElementsInstancedBaseInstance"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDrawElementsInstancedBaseVertex), "glDrawElementsInstancedBaseVertex"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDrawElementsInstancedBaseVertexBaseInstance), "glDrawElementsInstancedBaseVertexBaseInstance"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDrawRangeElements), "glDrawRangeElements"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDrawRangeElementsBaseVertex), "glDrawRangeElementsBaseVertex"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDrawTransformFeedback), "glDrawTransformFeedback"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDrawTransformFeedbackInstanced), "glDrawTransformFeedbackInstanced"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDrawTransformFeedbackStream), "glDrawTransformFeedbackStream"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glDrawTransformFeedbackStreamInstanced), "glDrawTransformFeedbackStreamInstanced"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glEnable), "glEnable"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glEnableVertexAttribArray), "glEnableVertexAttribArray"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glEnablei), "glEnablei"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glEndConditionalRender), "glEndConditionalRender"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glEndQuery), "glEndQuery"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glEndQueryIndexed), "glEndQueryIndexed"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glEndTransformFeedback), "glEndTransformFeedback"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glFinish), "glFinish"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glFlush), "glFlush"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glFlushMappedBufferRange), "glFlushMappedBufferRange"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glFramebufferParameteri), "glFramebufferParameteri"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glFramebufferRenderbuffer), "glFramebufferRenderbuffer"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glFramebufferTexture), "glFramebufferTexture"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glFramebufferTexture1D), "glFramebufferTexture1D"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glFramebufferTexture2D), "glFramebufferTexture2D"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glFramebufferTexture3D), "glFramebufferTexture3D"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glFramebufferTextureLayer), "glFramebufferTextureLayer"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glFrontFace), "glFrontFace"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGenBuffers), "glGenBuffers"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGenFramebuffers), "glGenFramebuffers"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGenProgramPipelines), "glGenProgramPipelines"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGenQueries), "glGenQueries"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGenRenderbuffers), "glGenRenderbuffers"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGenSamplers), "glGenSamplers"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGenTextures), "glGenTextures"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGenTransformFeedbacks), "glGenTransformFeedbacks"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGenVertexArrays), "glGenVertexArrays"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGenerateMipmap), "glGenerateMipmap"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetActiveAtomicCounterBufferiv), "glGetActiveAtomicCounterBufferiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetActiveAttrib), "glGetActiveAttrib"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetActiveSubroutineName), "glGetActiveSubroutineName"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetActiveSubroutineUniformName), "glGetActiveSubroutineUniformName"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetActiveSubroutineUniformiv), "glGetActiveSubroutineUniformiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetActiveUniform), "glGetActiveUniform"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetActiveUniformBlockName), "glGetActiveUniformBlockName"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetActiveUniformBlockiv), "glGetActiveUniformBlockiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetActiveUniformName), "glGetActiveUniformName"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetActiveUniformsiv), "glGetActiveUniformsiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetAttachedShaders), "glGetAttachedShaders"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetBooleani_v), "glGetBooleani_v"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetBooleanv), "glGetBooleanv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetBufferParameteri64v), "glGetBufferParameteri64v"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetBufferParameteriv), "glGetBufferParameteriv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetBufferPointerv), "glGetBufferPointerv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetBufferSubData), "glGetBufferSubData"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetCompressedTexImage), "glGetCompressedTexImage"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetDoublei_v), "glGetDoublei_v"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetDoublev), "glGetDoublev"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetFloati_v), "glGetFloati_v"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetFloatv), "glGetFloatv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetFramebufferAttachmentParameteriv), "glGetFramebufferAttachmentParameteriv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetFramebufferParameteriv), "glGetFramebufferParameteriv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetInteger64i_v), "glGetInteger64i_v"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetInteger64v), "glGetInteger64v"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetIntegeri_v), "glGetIntegeri_v"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetIntegerv), "glGetIntegerv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetInternalformati64v), "glGetInternalformati64v"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetInternalformativ), "glGetInternalformativ"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetMultisamplefv), "glGetMultisamplefv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetNamedFramebufferParameterivEXT), "glGetNamedFramebufferParameterivEXT"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetNamedStringARB), "glGetNamedStringARB"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetNamedStringivARB), "glGetNamedStringivARB"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetObjectLabel), "glGetObjectLabel"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetObjectPtrLabel), "glGetObjectPtrLabel"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetPointerv), "glGetPointerv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetProgramBinary), "glGetProgramBinary"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetProgramInfoLog), "glGetProgramInfoLog"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetProgramInterfaceiv), "glGetProgramInterfaceiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetProgramPipelineInfoLog), "glGetProgramPipelineInfoLog"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetProgramPipelineiv), "glGetProgramPipelineiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetProgramResourceName), "glGetProgramResourceName"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetProgramResourceiv), "glGetProgramResourceiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetProgramStageiv), "glGetProgramStageiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetProgramiv), "glGetProgramiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetQueryIndexediv), "glGetQueryIndexediv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetQueryObjecti64v), "glGetQueryObjecti64v"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetQueryObjectiv), "glGetQueryObjectiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetQueryObjectui64v), "glGetQueryObjectui64v"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetQueryObjectuiv), "glGetQueryObjectuiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetQueryiv), "glGetQueryiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetRenderbufferParameteriv), "glGetRenderbufferParameteriv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetSamplerParameterIiv), "glGetSamplerParameterIiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetSamplerParameterIuiv), "glGetSamplerParameterIuiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetSamplerParameterfv), "glGetSamplerParameterfv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetSamplerParameteriv), "glGetSamplerParameteriv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetShaderInfoLog), "glGetShaderInfoLog"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetShaderPrecisionFormat), "glGetShaderPrecisionFormat"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetShaderSource), "glGetShaderSource"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetShaderiv), "glGetShaderiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetSynciv), "glGetSynciv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetTexImage), "glGetTexImage"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetTexLevelParameterfv), "glGetTexLevelParameterfv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetTexLevelParameteriv), "glGetTexLevelParameteriv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetTexParameterIiv), "glGetTexParameterIiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetTexParameterIuiv), "glGetTexParameterIuiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetTexParameterfv), "glGetTexParameterfv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetTexParameteriv), "glGetTexParameteriv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetTransformFeedbackVarying), "glGetTransformFeedbackVarying"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetUniformIndices), "glGetUniformIndices"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetUniformSubroutineuiv), "glGetUniformSubroutineuiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetUniformdv), "glGetUniformdv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetUniformfv), "glGetUniformfv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetUniformiv), "glGetUniformiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetUniformuiv), "glGetUniformuiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetVertexAttribIiv), "glGetVertexAttribIiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetVertexAttribIuiv), "glGetVertexAttribIuiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetVertexAttribLdv), "glGetVertexAttribLdv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetVertexAttribPointerv), "glGetVertexAttribPointerv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetVertexAttribdv), "glGetVertexAttribdv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetVertexAttribfv), "glGetVertexAttribfv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetVertexAttribiv), "glGetVertexAttribiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetnColorTableARB), "glGetnColorTableARB"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetnCompressedTexImageARB), "glGetnCompressedTexImageARB"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetnConvolutionFilterARB), "glGetnConvolutionFilterARB"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetnHistogramARB), "glGetnHistogramARB"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetnMapdvARB), "glGetnMapdvARB"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetnMapfvARB), "glGetnMapfvARB"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetnMapivARB), "glGetnMapivARB"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetnMinmaxARB), "glGetnMinmaxARB"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetnPixelMapfvARB), "glGetnPixelMapfvARB"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetnPixelMapuivARB), "glGetnPixelMapuivARB"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetnPixelMapusvARB), "glGetnPixelMapusvARB"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetnPolygonStippleARB), "glGetnPolygonStippleARB"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetnSeparableFilterARB), "glGetnSeparableFilterARB"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetnTexImageARB), "glGetnTexImageARB"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetnUniformdvARB), "glGetnUniformdvARB"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetnUniformfvARB), "glGetnUniformfvARB"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetnUniformivARB), "glGetnUniformivARB"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glGetnUniformuivARB), "glGetnUniformuivARB"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glHint), "glHint"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glInvalidateBufferData), "glInvalidateBufferData"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glInvalidateBufferSubData), "glInvalidateBufferSubData"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glInvalidateFramebuffer), "glInvalidateFramebuffer"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glInvalidateSubFramebuffer), "glInvalidateSubFramebuffer"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glInvalidateTexImage), "glInvalidateTexImage"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glInvalidateTexSubImage), "glInvalidateTexSubImage"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glLineWidth), "glLineWidth"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glLinkProgram), "glLinkProgram"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glLogicOp), "glLogicOp"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glMemoryBarrier), "glMemoryBarrier"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glMinSampleShading), "glMinSampleShading"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glMinSampleShadingARB), "glMinSampleShadingARB"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glMultiDrawArrays), "glMultiDrawArrays"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glMultiDrawArraysIndirect), "glMultiDrawArraysIndirect"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glMultiDrawElements), "glMultiDrawElements"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glMultiDrawElementsBaseVertex), "glMultiDrawElementsBaseVertex"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glMultiDrawElementsIndirect), "glMultiDrawElementsIndirect"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glMultiTexCoordP1ui), "glMultiTexCoordP1ui"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glMultiTexCoordP1uiv), "glMultiTexCoordP1uiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glMultiTexCoordP2ui), "glMultiTexCoordP2ui"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glMultiTexCoordP2uiv), "glMultiTexCoordP2uiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glMultiTexCoordP3ui), "glMultiTexCoordP3ui"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glMultiTexCoordP3uiv), "glMultiTexCoordP3uiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glMultiTexCoordP4ui), "glMultiTexCoordP4ui"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glMultiTexCoordP4uiv), "glMultiTexCoordP4uiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glNamedFramebufferParameteriEXT), "glNamedFramebufferParameteriEXT"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glNamedStringARB), "glNamedStringARB"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glNormalP3ui), "glNormalP3ui"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glNormalP3uiv), "glNormalP3uiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glObjectLabel), "glObjectLabel"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glObjectPtrLabel), "glObjectPtrLabel"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glPatchParameterfv), "glPatchParameterfv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glPatchParameteri), "glPatchParameteri"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glPauseTransformFeedback), "glPauseTransformFeedback"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glPixelStoref), "glPixelStoref"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glPixelStorei), "glPixelStorei"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glPointParameterf), "glPointParameterf"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glPointParameterfv), "glPointParameterfv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glPointParameteri), "glPointParameteri"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glPointParameteriv), "glPointParameteriv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glPointSize), "glPointSize"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glPolygonMode), "glPolygonMode"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glPolygonOffset), "glPolygonOffset"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glPopDebugGroup), "glPopDebugGroup"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glPrimitiveRestartIndex), "glPrimitiveRestartIndex"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramBinary), "glProgramBinary"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramParameteri), "glProgramParameteri"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniform1d), "glProgramUniform1d"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniform1dv), "glProgramUniform1dv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniform1f), "glProgramUniform1f"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniform1fv), "glProgramUniform1fv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniform1i), "glProgramUniform1i"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniform1iv), "glProgramUniform1iv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniform1ui), "glProgramUniform1ui"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniform1uiv), "glProgramUniform1uiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniform2d), "glProgramUniform2d"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniform2dv), "glProgramUniform2dv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniform2f), "glProgramUniform2f"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniform2fv), "glProgramUniform2fv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniform2i), "glProgramUniform2i"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniform2iv), "glProgramUniform2iv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniform2ui), "glProgramUniform2ui"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniform2uiv), "glProgramUniform2uiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniform3d), "glProgramUniform3d"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniform3dv), "glProgramUniform3dv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniform3f), "glProgramUniform3f"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniform3fv), "glProgramUniform3fv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniform3i), "glProgramUniform3i"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniform3iv), "glProgramUniform3iv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniform3ui), "glProgramUniform3ui"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniform3uiv), "glProgramUniform3uiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniform4d), "glProgramUniform4d"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniform4dv), "glProgramUniform4dv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniform4f), "glProgramUniform4f"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniform4fv), "glProgramUniform4fv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniform4i), "glProgramUniform4i"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniform4iv), "glProgramUniform4iv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniform4ui), "glProgramUniform4ui"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniform4uiv), "glProgramUniform4uiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniformMatrix2dv), "glProgramUniformMatrix2dv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniformMatrix2fv), "glProgramUniformMatrix2fv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniformMatrix2x3dv), "glProgramUniformMatrix2x3dv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniformMatrix2x3fv), "glProgramUniformMatrix2x3fv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniformMatrix2x4dv), "glProgramUniformMatrix2x4dv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniformMatrix2x4fv), "glProgramUniformMatrix2x4fv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniformMatrix3dv), "glProgramUniformMatrix3dv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniformMatrix3fv), "glProgramUniformMatrix3fv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniformMatrix3x2dv), "glProgramUniformMatrix3x2dv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniformMatrix3x2fv), "glProgramUniformMatrix3x2fv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniformMatrix3x4dv), "glProgramUniformMatrix3x4dv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniformMatrix3x4fv), "glProgramUniformMatrix3x4fv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniformMatrix4dv), "glProgramUniformMatrix4dv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniformMatrix4fv), "glProgramUniformMatrix4fv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniformMatrix4x2dv), "glProgramUniformMatrix4x2dv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniformMatrix4x2fv), "glProgramUniformMatrix4x2fv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniformMatrix4x3dv), "glProgramUniformMatrix4x3dv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProgramUniformMatrix4x3fv), "glProgramUniformMatrix4x3fv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glProvokingVertex), "glProvokingVertex"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glPushDebugGroup), "glPushDebugGroup"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glQueryCounter), "glQueryCounter"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glReadBuffer), "glReadBuffer"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glReadPixels), "glReadPixels"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glReadnPixelsARB), "glReadnPixelsARB"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glReleaseShaderCompiler), "glReleaseShaderCompiler"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glRenderbufferStorage), "glRenderbufferStorage"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glRenderbufferStorageMultisample), "glRenderbufferStorageMultisample"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glResumeTransformFeedback), "glResumeTransformFeedback"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glSampleCoverage), "glSampleCoverage"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glSampleMaski), "glSampleMaski"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glSamplerParameterIiv), "glSamplerParameterIiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glSamplerParameterIuiv), "glSamplerParameterIuiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glSamplerParameterf), "glSamplerParameterf"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glSamplerParameterfv), "glSamplerParameterfv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glSamplerParameteri), "glSamplerParameteri"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glSamplerParameteriv), "glSamplerParameteriv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glScissor), "glScissor"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glScissorArrayv), "glScissorArrayv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glScissorIndexed), "glScissorIndexed"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glScissorIndexedv), "glScissorIndexedv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glSecondaryColorP3ui), "glSecondaryColorP3ui"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glSecondaryColorP3uiv), "glSecondaryColorP3uiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glShaderBinary), "glShaderBinary"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glShaderSource), "glShaderSource"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glShaderStorageBlockBinding), "glShaderStorageBlockBinding"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glStencilFunc), "glStencilFunc"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glStencilFuncSeparate), "glStencilFuncSeparate"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glStencilMask), "glStencilMask"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glStencilMaskSeparate), "glStencilMaskSeparate"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glStencilOp), "glStencilOp"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glStencilOpSeparate), "glStencilOpSeparate"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glTexBuffer), "glTexBuffer"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glTexBufferRange), "glTexBufferRange"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glTexCoordP1ui), "glTexCoordP1ui"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glTexCoordP1uiv), "glTexCoordP1uiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glTexCoordP2ui), "glTexCoordP2ui"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glTexCoordP2uiv), "glTexCoordP2uiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glTexCoordP3ui), "glTexCoordP3ui"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glTexCoordP3uiv), "glTexCoordP3uiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glTexCoordP4ui), "glTexCoordP4ui"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glTexCoordP4uiv), "glTexCoordP4uiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glTexImage1D), "glTexImage1D"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glTexImage2D), "glTexImage2D"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glTexImage2DMultisample), "glTexImage2DMultisample"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glTexImage3D), "glTexImage3D"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glTexImage3DMultisample), "glTexImage3DMultisample"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glTexParameterIiv), "glTexParameterIiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glTexParameterIuiv), "glTexParameterIuiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glTexParameterf), "glTexParameterf"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glTexParameterfv), "glTexParameterfv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glTexParameteri), "glTexParameteri"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glTexParameteriv), "glTexParameteriv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glTexStorage1D), "glTexStorage1D"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glTexStorage2D), "glTexStorage2D"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glTexStorage2DMultisample), "glTexStorage2DMultisample"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glTexStorage3D), "glTexStorage3D"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glTexStorage3DMultisample), "glTexStorage3DMultisample"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glTexSubImage1D), "glTexSubImage1D"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glTexSubImage2D), "glTexSubImage2D"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glTexSubImage3D), "glTexSubImage3D"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glTextureBufferRangeEXT), "glTextureBufferRangeEXT"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glTextureStorage1DEXT), "glTextureStorage1DEXT"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glTextureStorage2DEXT), "glTextureStorage2DEXT"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glTextureStorage2DMultisampleEXT), "glTextureStorage2DMultisampleEXT"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glTextureStorage3DEXT), "glTextureStorage3DEXT"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glTextureStorage3DMultisampleEXT), "glTextureStorage3DMultisampleEXT"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glTextureView), "glTextureView"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glTransformFeedbackVaryings), "glTransformFeedbackVaryings"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniform1d), "glUniform1d"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniform1dv), "glUniform1dv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniform1f), "glUniform1f"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniform1fv), "glUniform1fv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniform1i), "glUniform1i"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniform1iv), "glUniform1iv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniform1ui), "glUniform1ui"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniform1uiv), "glUniform1uiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniform2d), "glUniform2d"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniform2dv), "glUniform2dv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniform2f), "glUniform2f"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniform2fv), "glUniform2fv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniform2i), "glUniform2i"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniform2iv), "glUniform2iv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniform2ui), "glUniform2ui"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniform2uiv), "glUniform2uiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniform3d), "glUniform3d"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniform3dv), "glUniform3dv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniform3f), "glUniform3f"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniform3fv), "glUniform3fv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniform3i), "glUniform3i"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniform3iv), "glUniform3iv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniform3ui), "glUniform3ui"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniform3uiv), "glUniform3uiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniform4d), "glUniform4d"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniform4dv), "glUniform4dv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniform4f), "glUniform4f"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniform4fv), "glUniform4fv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniform4i), "glUniform4i"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniform4iv), "glUniform4iv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniform4ui), "glUniform4ui"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniform4uiv), "glUniform4uiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniformBlockBinding), "glUniformBlockBinding"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniformMatrix2dv), "glUniformMatrix2dv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniformMatrix2fv), "glUniformMatrix2fv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniformMatrix2x3dv), "glUniformMatrix2x3dv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniformMatrix2x3fv), "glUniformMatrix2x3fv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniformMatrix2x4dv), "glUniformMatrix2x4dv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniformMatrix2x4fv), "glUniformMatrix2x4fv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniformMatrix3dv), "glUniformMatrix3dv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniformMatrix3fv), "glUniformMatrix3fv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniformMatrix3x2dv), "glUniformMatrix3x2dv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniformMatrix3x2fv), "glUniformMatrix3x2fv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniformMatrix3x4dv), "glUniformMatrix3x4dv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniformMatrix3x4fv), "glUniformMatrix3x4fv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniformMatrix4dv), "glUniformMatrix4dv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniformMatrix4fv), "glUniformMatrix4fv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniformMatrix4x2dv), "glUniformMatrix4x2dv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniformMatrix4x2fv), "glUniformMatrix4x2fv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniformMatrix4x3dv), "glUniformMatrix4x3dv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniformMatrix4x3fv), "glUniformMatrix4x3fv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUniformSubroutinesuiv), "glUniformSubroutinesuiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUseProgram), "glUseProgram"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glUseProgramStages), "glUseProgramStages"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glValidateProgram), "glValidateProgram"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glValidateProgramPipeline), "glValidateProgramPipeline"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexArrayBindVertexBufferEXT), "glVertexArrayBindVertexBufferEXT"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexArrayVertexAttribBindingEXT), "glVertexArrayVertexAttribBindingEXT"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexArrayVertexAttribFormatEXT), "glVertexArrayVertexAttribFormatEXT"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexArrayVertexAttribIFormatEXT), "glVertexArrayVertexAttribIFormatEXT"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexArrayVertexAttribLFormatEXT), "glVertexArrayVertexAttribLFormatEXT"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexArrayVertexBindingDivisorEXT), "glVertexArrayVertexBindingDivisorEXT"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttrib1d), "glVertexAttrib1d"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttrib1dv), "glVertexAttrib1dv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttrib1f), "glVertexAttrib1f"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttrib1fv), "glVertexAttrib1fv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttrib1s), "glVertexAttrib1s"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttrib1sv), "glVertexAttrib1sv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttrib2d), "glVertexAttrib2d"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttrib2dv), "glVertexAttrib2dv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttrib2f), "glVertexAttrib2f"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttrib2fv), "glVertexAttrib2fv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttrib2s), "glVertexAttrib2s"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttrib2sv), "glVertexAttrib2sv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttrib3d), "glVertexAttrib3d"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttrib3dv), "glVertexAttrib3dv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttrib3f), "glVertexAttrib3f"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttrib3fv), "glVertexAttrib3fv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttrib3s), "glVertexAttrib3s"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttrib3sv), "glVertexAttrib3sv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttrib4Nbv), "glVertexAttrib4Nbv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttrib4Niv), "glVertexAttrib4Niv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttrib4Nsv), "glVertexAttrib4Nsv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttrib4Nub), "glVertexAttrib4Nub"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttrib4Nubv), "glVertexAttrib4Nubv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttrib4Nuiv), "glVertexAttrib4Nuiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttrib4Nusv), "glVertexAttrib4Nusv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttrib4bv), "glVertexAttrib4bv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttrib4d), "glVertexAttrib4d"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttrib4dv), "glVertexAttrib4dv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttrib4f), "glVertexAttrib4f"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttrib4fv), "glVertexAttrib4fv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttrib4iv), "glVertexAttrib4iv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttrib4s), "glVertexAttrib4s"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttrib4sv), "glVertexAttrib4sv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttrib4ubv), "glVertexAttrib4ubv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttrib4uiv), "glVertexAttrib4uiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttrib4usv), "glVertexAttrib4usv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribBinding), "glVertexAttribBinding"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribDivisor), "glVertexAttribDivisor"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribFormat), "glVertexAttribFormat"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribI1i), "glVertexAttribI1i"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribI1iv), "glVertexAttribI1iv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribI1ui), "glVertexAttribI1ui"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribI1uiv), "glVertexAttribI1uiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribI2i), "glVertexAttribI2i"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribI2iv), "glVertexAttribI2iv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribI2ui), "glVertexAttribI2ui"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribI2uiv), "glVertexAttribI2uiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribI3i), "glVertexAttribI3i"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribI3iv), "glVertexAttribI3iv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribI3ui), "glVertexAttribI3ui"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribI3uiv), "glVertexAttribI3uiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribI4bv), "glVertexAttribI4bv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribI4i), "glVertexAttribI4i"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribI4iv), "glVertexAttribI4iv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribI4sv), "glVertexAttribI4sv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribI4ubv), "glVertexAttribI4ubv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribI4ui), "glVertexAttribI4ui"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribI4uiv), "glVertexAttribI4uiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribI4usv), "glVertexAttribI4usv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribIFormat), "glVertexAttribIFormat"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribIPointer), "glVertexAttribIPointer"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribL1d), "glVertexAttribL1d"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribL1dv), "glVertexAttribL1dv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribL2d), "glVertexAttribL2d"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribL2dv), "glVertexAttribL2dv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribL3d), "glVertexAttribL3d"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribL3dv), "glVertexAttribL3dv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribL4d), "glVertexAttribL4d"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribL4dv), "glVertexAttribL4dv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribLFormat), "glVertexAttribLFormat"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribLPointer), "glVertexAttribLPointer"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribP1ui), "glVertexAttribP1ui"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribP1uiv), "glVertexAttribP1uiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribP2ui), "glVertexAttribP2ui"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribP2uiv), "glVertexAttribP2uiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribP3ui), "glVertexAttribP3ui"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribP3uiv), "glVertexAttribP3uiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribP4ui), "glVertexAttribP4ui"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribP4uiv), "glVertexAttribP4uiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexAttribPointer), "glVertexAttribPointer"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexBindingDivisor), "glVertexBindingDivisor"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexP2ui), "glVertexP2ui"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexP2uiv), "glVertexP2uiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexP3ui), "glVertexP3ui"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexP3uiv), "glVertexP3uiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexP4ui), "glVertexP4ui"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glVertexP4uiv), "glVertexP4uiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glViewport), "glViewport"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glViewportArrayv), "glViewportArrayv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glViewportIndexedf), "glViewportIndexedf"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glViewportIndexedfv), "glViewportIndexedfv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_glWaitSync), "glWaitSync"))
	return
}
func isBound(pfn unsafe.Pointer, fn string) string {
	inc := " "
	if pfn != nil {
		inc = "+"
	}
	return fmt.Sprintf("   [%s] %s", inc, fn)
}
