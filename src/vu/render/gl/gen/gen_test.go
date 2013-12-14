// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package main

import (
	"testing"
)

func TestSplitLine(t *testing.T) {
	apiLine := "GLAPI GLboolean APIENTRY glIsEnabled (GLenum cap);"
	fname, rettype, plist, pnames, ptypes := splitLine(apiLine)
	if fname != "glIsEnabled" || rettype != "GLboolean" || plist != "GLenum cap" ||
		pnames[0] != "cap" || ptypes[0] != "GLenum" {
		t.Fail()
	}
}

// // GLAPI GLuint APIENTRY wrap_glCreateShader(unsigned int t_ype) { return (*pfn_glCreateShader)(t_ype); }
func TestCwrapper0(t *testing.T) {
	apiLine := "GLAPI GLuint APIENTRY glCreateShader (GLenum type);"
	_, rettype, _, pnames, ptypes := splitLine(apiLine)
	wret, wparms, wnames := cwrapper(rettype, pnames, ptypes, ptypes)
	if wret != "GLAPI unsigned int APIENTRY" || wparms != "unsigned int t_ype" || wnames != "(GLenum)t_ype" {
		t.Fail()
	}
}

// // GLAPI const char * APIENTRY wrap_glGetString(unsigned int name) { return (char *)(*pfn_glGetString)(name); }
func TestCwrapper1(t *testing.T) {
	apiLine := "GLAPI const GLubyte * APIENTRY glGetString (GLenum name);"
	_, rettype, _, pnames, ptypes := splitLine(apiLine)
	wret, wparms, wnames := cwrapper(rettype, pnames, ptypes, ptypes)
	if wret != "GLAPI const unsigned char* APIENTRY" || wparms != "unsigned int name" || wnames != "(GLenum)name" {
		t.Fail()
	}
}

// // GLAPI void APIENTRY wrap_glGetVertexAttribPointerv(unsigned int index, unsigned int pname, long long* pointer)
//	  { (*pfn_glGetVertexAttribPointerv)(index, pname, (GLvoid **)pointer); }"
func TestCwrapper2(t *testing.T) {
	apiLine := "GLAPI void APIENTRY glGetVertexAttribPointerv (GLuint index, GLenum pname, GLvoid* *pointer);"
	_, _, _, _, gltypes := splitLine(apiLine)
	apiLine = alterSpec(apiLine)
	_, rettype, _, pnames, ptypes := splitLine(apiLine)
	wret, wparms, wnames := cwrapper(rettype, pnames, ptypes, gltypes)
	if wret != "GLAPI void APIENTRY" || wparms != "unsigned int index, unsigned int pname, long long* pointer" ||
		wnames != "(GLuint)index, (GLenum)pname, (GLvoid* *)pointer" {
		t.Fail()
	}
}

// // GLAPI void APIENTRY wrap_glGetShaderInfoLog(unsigned int shader, int bufSize, int* length, unsigned char* infoLog)
//	  { (*pfn_glGetShaderInfoLog)(shader, bufSize, length, (GLchar *)infoLog); }"
func TestCwrapper3(t *testing.T) {
	apiLine := "GLAPI void APIENTRY glGetShaderInfoLog (GLuint shader, GLsizei bufSize, GLsizei *length, GLchar *infoLog);"
	_, rettype, _, pnames, ptypes := splitLine(apiLine)
	wret, wparms, wnames := cwrapper(rettype, pnames, ptypes, ptypes)
	if wret != "GLAPI void APIENTRY" || wparms != "unsigned int shader, int bufSize, int* length, char* infoLog" ||
		wnames != "(GLuint)shader, (GLsizei)bufSize, (GLsizei *)length, (GLchar *)infoLog" {
		t.Fail()
	}
}

// func CreateShader(t_ype uint32) uint32 {
//    return uint32(C.wrap_glCreateShader(C.uint(t_ype)))
// }
func TestGowrapper0(t *testing.T) {
	mname, goparms, goret, cgoparms := gowrapper("GLAPI GLuint APIENTRY glCreateShader (GLenum type);")
	if mname != "CreateShader" || goparms != "t_ype uint32" || goret != "uint32" ||
		cgoparms != "C.uint(t_ype)" {
		t.Fail()
	}
}

//func GetShaderInfoLog(shader uint32, bufSize int32, length *int32, infoLog *uint8) {
//    C.wrap_glGetShaderInfoLog(C.uint(shader), C.int(bufSize), (*C.int)(length), (*C.uchar)(infoLog))
//}
func TestGowrapper3(t *testing.T) {
	apiLine := "GLAPI void APIENTRY glGetShaderInfoLog (GLuint shader, GLsizei bufSize, GLsizei *length, GLchar *infoLog);"
	apiLine = alterSpec(apiLine)
	mname, goparms, goret, cgoparms := gowrapper(apiLine)
	if mname != "GetShaderInfoLog" || goparms != "shader uint32, bufSize int32, length *int32, infoLog *uint8" || goret != "" ||
		cgoparms != "C.uint(shader), C.int(bufSize), (*C.int)(length), (*C.uchar)(infoLog)" {
		t.Fail()
	}
}
