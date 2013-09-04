// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

// Try opening a non-existing file.
func TestParseBadSpec(t *testing.T) {
	args := []string{"gengl", "glpkg_test", "bad spec name"}
	if _, _, err := parseUserInput(args); err == nil {
		t.Error("Did not catch bad specification name")
	}
	os.Remove("glpkg_test")
	os.Remove(args[0])
}

// Check that valid input passes parsing.
func TestParseOk(t *testing.T) {
	args := []string{"gengl", "../gl", "glcorearb.h-v4.3"}
	if _, _, err := parseUserInput(args); err != nil {
		t.Error(err)
	}
	os.Remove("../gl")
}

// Check that c method definitions are properly split.
func TestSplitCMethod(t *testing.T) {
	in := "GLAPI void APIENTRY glEnablei (GLenum target, GLuint index);"
	want := "GLAPI void APIENTRY-glEnablei-(GLenum target, GLuint index)"
	ctype, cname, cparms := splitCMethod(in)
	out := strings.Join(append([]string{}, ctype, cname, cparms), "-")
	if out != want {
		t.Error("converting ", in, "\nwanted ", want, "\ngot    ", out)
	}
}

// Check that a parmeter list can be split into names and types
func TestSplitParms(t *testing.T) {
	in := "(GLenum target, GLenum range, const GLfloat *params);"
	want := "[GLenum GLenum const GLfloat*] [target r_ange params]"
	ptypes, pnames := splitParms(in)
	out := fmt.Sprintf("%s %s", ptypes, pnames)
	if out != want {
		t.Error("converting ", in, "\nwanted ", want, "\ngot    ", out)
	}
}

// This was an interesting test where the parameter name was a substring
// of the type name and the code assumed that would never be the case.
func TestSplitParmsGotcha(t *testing.T) {
	in := "(GLfloat n, GLfloat f);"
	want := "[GLfloat GLfloat] [n f]"
	ptypes, pnames := splitParms(in)
	out := fmt.Sprintf("%s %s", ptypes, pnames)
	if out != want {
		t.Error("converting ", in, "\nwanted ", want, "\ngot    ", out)
	}
}

// Check that the spec can be parsed.
func TestParseSpec(t *testing.T) {
	want := "DEPTH_BUFFER_BIT = 0x00000100"
	spec, _ := os.Open("glcorearb.h-v4.3")
	defer spec.Close()
	groups := parseSpec(spec)
	if groups == nil {
		t.Error("parseSpec returned nil")
	} else {
		out := groups["GL_VERSION_1_1"].goconsts[0]
		if out != want {
			t.Error("mismatch in parse", "\nwanted ", want, "\ngot    ", out)
		}
	}
}

// Check some c type conversions
func TestConvertType(t *testing.T) {
	inout := map[string]string{
		"GLenum":                         "unsigned int uint32 C.uint",
		"GLboolean*":                     "unsigned char* *uint8 (*C.uchar)",
		"const GLvoid*":                  "void* Pointer unsafe.Pointer",
		"const void*":                    "void* Pointer unsafe.Pointer",
		"const GLchar* const*":           "const char* const* []string _STR_ARRAY_",
		"GLAPI const GLubyte * APIENTRY": "char * string (*C.uchar)",
	}
	for in, want := range inout {
		ctype, gotype, cgotype := convertType(in)
		out := strings.Join(append([]string{}, ctype, gotype, cgotype), " ")
		if out != want {
			t.Error("converting ", in, "\nwanted ", want, "\ngot    ", out)
		}
	}
}

// Check that a cwrapper is properly generated.
func TestGenCWrapper(t *testing.T) {
	in := "GLAPI void APIENTRY glGetPointerv (GLenum *pname, GLvoid* *params);"
	want := "GLAPI void APIENTRY wrap_glGetPointerv(unsigned int* pname, void** params) " +
		"{ (*pfn_glGetPointerv)(pname, params); }"
	creturn, mname, cparms := splitCMethod(in)
	ptypes, pnames := splitParms(cparms)
	out := genCWrapper(creturn, mname, ptypes, pnames)
	if out != want {
		t.Error("gen for ", in, "\nwanted ", want, "\ngot    ", out)
	}
}

// Generate output strings from the API.
func TestGenGoFunc(t *testing.T) {
	in := "GLAPI void APIENTRY glEnablei (GLenum target, GLuint index);"
	want := "func Enablei (target uint32, index uint32)  { C.wrap_glEnablei( C.uint(target), C.uint(index)) }"
	creturn, mname, cparms := splitCMethod(in)
	ptype := strings.Replace(creturn, "APIENTRY", "", 1)
	ptype = strings.Replace(ptype, "GLAPI", "", 1)
	ptype = strings.TrimSpace(ptype)
	/*
		m := &method{}
		m.cmethod = cmethod
		m.cmethodPtr = ptype + " (APIENTRYP pfn_" + mname + ")" + cparms + ";"
	*/
	//cparms = specialHandling(mname, cparms, true)
	ptypes, pnames := splitParms(cparms)
	out := genGoFunc(creturn, mname, ptypes, pnames)
	if out != want {
		t.Error("gen for ", in, "\nwanted ", want, "\ngot    ", out)
	}
}

// Generate output strings with an API method that has a return value.
func TestGenGoFuncReturn(t *testing.T) {
	in := "GLAPI GLboolean APIENTRY glIsSampler (GLuint sampler);"
	want := "func IsSampler (sampler uint32) bool { return cbool(uint(C.wrap_glIsSampler( C.uint(sampler)))) }"
	creturn, mname, cparms := splitCMethod(in)
	ptypes, pnames := splitParms(cparms)
	out := genGoFunc(creturn, mname, ptypes, pnames)
	if out != want {
		t.Error("gen for ", in, "\nwanted ", want, "\ngot    ", out)
	}
}

// Generate output strings with an API method that has pointers and
// a reserved word.
func TestGenGoFuncPointers(t *testing.T) {
	in := "GLAPI void APIENTRY glMultiDrawElements (const GLsizei *count, GLenum type, const GLvoid* const *indices);"
	want := "func MultiDrawElements (count *int32, t_ype uint32, indices *Pointer)  " +
		"{ C.wrap_glMultiDrawElements( (*C.int)(count), C.uint(t_ype), (*unsafe.Pointer)(indices)) }"
	creturn, mname, cparms := splitCMethod(in)
	ptypes, pnames := splitParms(cparms)
	out := genGoFunc(creturn, mname, ptypes, pnames)
	if out != want {
		t.Error("gen for ", in, "\nwanted ", want, "\ngot    ", out)
	}
}

// Generated go func that uses string conversion.
func TestGenGoStrings(t *testing.T) {
	in := "GLAPI GLuint APIENTRY glCreateShaderProgramv (GLenum type, GLsizei count, const GLchar* const *strings);"
	want := "" +
		"func CreateShaderProgramv (t_ype uint32, count int32, strings []string) uint32 {\n" +
		"   cstrings := C.newStringArray(C.int(len(strings)))\n" +
		"   defer C.freeStringArray(cstrings, C.int(len(strings)))\n" +
		"   for cnt, str := range strings {\n" +
		"      C.assignString(cstrings, C.CString(str), C.int(cnt))\n" +
		"   }\n" +
		"   return uint32(C.wrap_glCreateShaderProgramv( C.uint(t_ype), C.int(count), cstrings)) }"
	creturn, mname, cparms := splitCMethod(in)
	ptypes, pnames := splitParms(cparms)
	out := genGoFunc(creturn, mname, ptypes, pnames)
	if out != want {
		t.Error("gen for ", in, "\nwanted\n", want, "\ngot\n", out)
	}
}
