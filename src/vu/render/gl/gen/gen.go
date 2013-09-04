// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

// Package gen is used to generate the golang OpenGL bindings from the
// latest OpenGL specification.  The generator uses the most recent header
// file downloaded from
//     www.opengl.org/registry/api/glcorearb.h
//
// Running gl/gen will create a go package that can then be installed
// and used. Usage:
//     gen output-directory opengl-spec
//     gen .. glcorearb.h-v4.3
//
// The generated gl source package can be (is expected to be) source controlled
// since OpenGL specifications aren't updated all that often. Noted exceptions
// to the binding are:
//     glDebugMessageCallback
//     glDebugMessageCallbackARB
//
// Thanks to https://github.com/chsc/gogl for the idea of generating the bindings
// from the specification.  Chsc/gogl is likely simpler/better in that it uses
// the XML based specifications rather than the glcorearb.h.
package main

// Design Notes: Essentially straight line code to get information in one
// format: the OpenGL Specification, to another format: the Go language bindings
// for OpenGL.
//
// Maintenance Notes: The code is hand-tweaked in parts and has just enough smarts
// to process the current specification.  For example not all go keywords are
// checked for safety, only the ones that were used in the OpenGL specification.
// The code is expected to be further tweaked as new specifications come along.
// The idea is to only generalize the code where necessary without adding undo
// bulk or dependencies to the code.  Anything that would reduce the code and
// improve readability would be appreciated and should be done.

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Parse the user input, parse the spec, generate the output.
func main() {
	pkg, spec, err := parseUserInput(os.Args)
	if err != nil {
		fmt.Printf("usage: gen generated-package-directory opengl-spec.h\n")
		os.Exit(-1)
	}
	defer spec.Close()

	// parse spec and pre-generate lines into internal data structures.
	groups := parseSpec(spec)

	// generate output.
	fi, _ := spec.Stat()
	genPackage(pkg, fi.Name(), groups)
}

// Parse the input arguments and returns the opened input file and
// output package directory location.
func parseUserInput(args []string) (pkg string, spec *os.File, err error) {
	if len(args) == 3 {
		pkg = args[1]
		if os.MkdirAll(pkg, os.FileMode(0755)) == nil {
			spec, err = os.Open(args[2])
		}
	} else {
		err = errors.New("gen: need two input arguments")
	}
	return
}

// ===============================================================
// parse spec and pre-generate lines into internal data structures
// ===============================================================

// typemap turns OpenGL types into the type expressions needed for the c-wrapper
// and the go func.  Not every type is represented, only those used in the spec.
var typemap = map[string]types{
	"GLchar*":              {"char*", "string", "(*C.char)", false},
	"GLubyte":              {"unsigned char", "uint8", "C.uchar", false},
	"GLubyte*":             {"unsigned char*", "*uint8", "(*C.uchar)", false},
	"GLboolean":            {"unsigned char", "bool", "C.uchar", false},
	"GLboolean*":           {"unsigned char*", "*uint8", "(*C.uchar)", false},
	"GLenum":               {"unsigned int", "uint32", "C.uint", false},
	"GLfloat":              {"float", "float32", "C.float", false},
	"GLfloat*":             {"float*", "*float32", "(*C.float)", false},
	"GLint":                {"int", "int32", "C.int", false},
	"GLint*":               {"int*", "*int32", "(*C.int)", false},
	"GLintptr":             {"long long", "int64", "C.longlong", false},
	"GLuint":               {"unsigned int", "uint32", "C.uint", false},
	"GLuint*":              {"unsigned int*", "*uint32", "(*C.uint)", false},
	"GLsizei":              {"int", "int32", "C.int", false},
	"GLsizeiptr":           {"long long", "int64", "C.longlong", false},
	"GLbitfield":           {"unsigned int", "uint32", "C.uint", false},
	"GLdouble":             {"double", "float64", "C.double", false},
	"GLvoid*":              {"void*", "Pointer", "unsafe.Pointer", false},
	"GLvoid**":             {"void**", "*Pointer", "(*unsafe.Pointer)", false},
	"GLdouble*":            {"double*", "*float64", "(*C.double)", false},
	"GLsizei*":             {"int*", "*int32", "(*C.int)", false},
	"GLenum*":              {"unsigned int*", "*uint32", "(*C.uint)", false},
	"GLshort":              {"short ", "int16", "C.short", false},
	"GLushort*":            {"unsigned short*", "*uint16", "(*C.ushort)", false},
	"GLsync":               {"GLsync", "Sync", "C.GLsync", false},
	"GLint64*":             {"long long*", "*int64", "(*C.longlong)", false},
	"GLuint64":             {"unsigned long long", "uint64", "C.ulonglong", false},
	"GLuint64*":            {"unsigned long long*", "*uint64", "(*C.ulonglong)", false},
	"const void*":          {"void*", "Pointer", "unsafe.Pointer", false},
	"const GLfloat*":       {"float*", "*float32", "(*C.float)", false},
	"const GLint*":         {"int*", "*int32", "(*C.int)", false},
	"const GLvoid*":        {"void*", "Pointer", "unsafe.Pointer", false},
	"const GLuint*":        {"unsigned int*", "*uint32", "(*C.uint)", false},
	"const GLsizei*":       {"int*", "*int32", "(*C.int)", false},
	"const GLvoid* const*": {"const void* const*", "*Pointer", "(*unsafe.Pointer)", false},
	"const GLenum*":        {"unsigned int*", "*uint32", "(*C.uint)", false},
	"const GLchar*":        {"char*", "string", "(*C.char)", false},
	"const GLdouble*":      {"double*", "*float64", "(*C.double)", false},
	"const GLshort*":       {"short *", "*int16", "(*C.short)", false},
	"const GLbyte*":        {"signed char*", "*int8", "(*C.schar)", false},
	"const GLubyte*":       {"unsigned char*", "*uint8", "(*C.uchar)", false},
	"const GLushort*":      {"unsigned short*", "*uint16", "(*C.ushort)", false},
	"const GLint64":        {"long long", "int64", "C.longlong", false},

	// array of strings.
	"const GLchar**":       {"const char**", "[]string", strArray, false},
	"const GLchar* const*": {"const char* const*", "[]string", strArray, false},

	// return types
	"GLAPI GLenum APIENTRY":          {"unsigned int", "uint32", "C.uint", false},
	"GLAPI const GLubyte * APIENTRY": {"char *", "string", "(*C.uchar)", false},
	"GLAPI GLboolean APIENTRY":       {"unsigned char", "bool", "C.uchar", false},
	"GLAPI GLvoid* APIENTRY":         {"void*", "Pointer", "unsafe.Pointer", false},
	"GLAPI GLuint APIENTRY":          {"unsigned int", "uint32", "C.uint", false},
	"GLAPI GLint APIENTRY":           {"int", "int32", "C.int", false},
	"GLAPI GLsync APIENTRY":          {"GLsync", "Sync", "C.GLsync", false},

	// The following may not be in every spec.
	"struct _cl_context*": {"struct _cl_context*", "*clContext", "(*C.struct_cl_context)", false},
	"struct _cl_event*":   {"struct _cl_event*", "*clEvent", "(*C.struct_cl_event)", false},
	"GLDEBUGPROCARB":      {"GLDEBUGPROCARB", "DEBUGPROCARB", "C.GLDEBUGPROCARB", false},
	"GLDEBUGPROC":         {"GLDEBUGPROC", "DEBUGPROC", "C.GLDEBUGPROC", false},
}

// typemap holds the different OpenGL type mappings as well as whether or not
// a given type is used in the spec.  For example GLDEBUGPRO* is not use in
// the 3.2 spec.  The used flag is set as the spec is being parsed and is
// referenced when generating the golang special type mappings.
type types struct {
	ctype   string // used in the c-wrapper
	gotype  string // used in the go func golang wrapper
	cgotype string // used to map between golang and c
	used    bool   // marks if this type was seen in the spec
}

// Organize output according to the groupings defined in the OpenGL
// specification.  The methods and constant definitions for each grouping
// are stored in this structure.
type grouping struct {
	// A map of OpenGL method signatures for this group.
	// The key is the OpenGl method name, eg: glActiveTexture
	methods map[string]method

	// the #define constants for this group. These are the original
	// lines from the OpenGL specification.
	goconsts []string
}

// Track the various OpenGL method signatures. There is one
// instance of these for each OpenGL method signature.  The correspoding
// output for the OpenGL method is created, and stored in this structure,
// as the lines are read in.
type method struct {
	cmethod    string // the original definition from the specification.
	cmethodPtr string // definition for a function pointer to the c method.
	cwrapper   string // a generated method that calls the c function pointer.
	gofunc     string // the go method that calls the c wrapper method.
}

// Get the sorted keys for a group of methods.  This is used to ensure
// the same output is created.  Othewise accessing the hashmap for the
// list of keys will result in a different output file each time.
func (g *grouping) sortedMethodNames() []string {
	sortedKeys := make([]string, len(g.methods))
	cnt := 0
	for key := range g.methods {
		sortedKeys[cnt] = key
		cnt++
	}
	sort.Strings(sortedKeys)
	return sortedKeys
}

// Reads the OpenGL specification into internal data structures.
// Process the spec line by line and break each method line down into parts
// The corresponding c-wrapper and golang func methods are also created
// at this point. No output happens yet.
func parseSpec(spec *os.File) (groups map[string]grouping) {
	// Groups are arb or specification-release identifiers.
	groups = make(map[string]grouping)
	var groupName string

	fileReader := bufio.NewReader(spec)
	bytes, isPrefix, err := fileReader.ReadLine()
	for err == nil && !isPrefix {
		line := string(bytes)
		if strings.Index(line, "#ifndef GL_ARB") == 0 ||
			strings.Index(line, "#ifndef GL_VER") == 0 {
			// parse OpenGL specification groupings
			if currentGroupName := parseGroup(groups, line); currentGroupName != "" {
				groupName = currentGroupName
			} else {
				fmt.Fprintf(os.Stderr, "Ignoring group def %s\n", line)
			}

		} else if strings.Index(line, "GLAPI ") == 0 && groupName != "" {
			// parse and pre-gen methods
			parseMethod(groups[groupName], line)

		} else if strings.Index(line, "#define GL_") == 0 && groupName != "" {
			// map C constants to Go constants
			groups[groupName] = parseConst(groups, groupName, line)
		}

		// use parts of the spec in the C code preamble
		appendPreamble(line)
		bytes, isPrefix, err = fileReader.ReadLine()
	}
	return
}

// Create the data structure for a group if its not been seen before.
// Some groups may be empty if there were no constants or methods for
// them in the specification.
func parseGroup(groups map[string]grouping, line string) (groupName string) {
	groupName = ""
	fields := strings.Fields(line)
	if len(fields) == 2 {
		groupName = fields[1]
		if _, isKey := groups[groupName]; !isKey {
			group := grouping{}
			group.methods = map[string]method{}
			groups[groupName] = group
		}
	}
	return
}

// Generate all the various output information as each OpenGl method is
// encountered.
func parseMethod(group grouping, line string) {
	fields := strings.Fields(line)
	if len(fields) >= 4 {
		mname, m := genMethodDefs(line)
		group.methods[mname] = m
	} else {
		fmt.Fprintf(os.Stderr, "Ignoring short api %s\n", line)
	}
}

// Converts a c constant into a go constant and stores the result
// into the current specification grouping.
func parseConst(groups map[string]grouping, groupName string, cconst string) (group grouping) {
	group = groups[groupName] // return unchanged if there was a problem.
	fields := strings.Fields(cconst)
	if len(fields) == 3 {
		constName := fields[1]

		// ignore if the name is a group name.
		if _, isRelKey := groups[constName]; !isRelKey {
			goconst := genGoConst(cconst)
			group.goconsts = append(group.goconsts, goconst)
		}
	} else {
		fmt.Fprintf(os.Stderr, "Ignoring short const %s\n", cconst)
	}
	return
}

// Takes an OpenGL method and creates the c function pointer wrapper and
// the corresponding go language binding.  For example the following input:
//
//    GLAPI void APIENTRY glEnablei (GLenum target, GLuint index);
//
// results in
//
//    void (* pfn_glEnablei) (GLenum target, GLuint index);
//    void wrap_glEnablei (GLenum target, GLuint index) { (*pfn_glEnablei)(target, index); }
//    func Enablei (target uint32, index uint32) { C.wrap_glEnablei( C.GLenum(target), C.GLuint(index)) }
//
func genMethodDefs(cmethod string) (methodName string, meth method) {
	creturn, mname, cparms := splitCMethod(cmethod)
	ptype := strings.Replace(creturn, "APIENTRY", "", 1)
	ptype = strings.Replace(ptype, "GLAPI", "", 1)
	ptype = strings.TrimSpace(ptype)
	m := &method{}
	m.cmethod = cmethod
	m.cmethodPtr = ptype + " (APIENTRYP pfn_" + mname + ")" + cparms + ";"

	cparms = specialHandling(mname, cparms, true)
	ptypes, pnames := splitParms(cparms)
	m.cwrapper = genCWrapper(creturn, mname, ptypes, pnames)
	m.gofunc = genGoFunc(creturn, mname, ptypes, pnames)
	return mname, *m
}

// Break the c method definition into the return type, the method name, and
// the method parameters.
func splitCMethod(methodDef string) (creturn string, cname string, cparms string) {
	cparms = "(" + strings.Split(methodDef, "(")[1]
	nametype := strings.Replace(methodDef, cparms, "", 1)
	cparms = strings.Trim(cparms, ";")
	cparms = strings.Replace(cparms, "(void)", "()", 1)

	// method name
	nametypes := strings.Fields(nametype)
	cname = nametypes[len(nametypes)-1]

	// return types
	nametypes = nametypes[0 : len(nametypes)-1]
	creturn = strings.Join(nametypes, " ")
	return
}

// Splits a c parameter list into the parameter types and the parameter names.
func splitParms(cparms string) (ptypes []string, pnames []string) {
	// default to an empty parameter list
	ptypes = []string{}
	pnames = []string{}

	// associate pointers with the parameter type
	parms := strings.Replace(cparms, " *", "* ", -1)
	parms = strings.Trim(parms, "();")
	if parms != "" {
		args := strings.Split(parms, ",")
		for _, arg := range args {
			tokens := strings.Fields(arg)
			pname := safeName(tokens[len(tokens)-1])
			pnames = append(pnames, pname)
			ptype := arg[:len(arg)-len(pname)]
			ptypes = append(ptypes, strings.TrimSpace(ptype))
		}
	}
	return
}

// Special handling case where certain methods need to be massaged to get a nice
// binding.  The original method definition is left as is, but the wrappers will
// use slightly different types.
func specialHandling(mname string, sig string, isParameters bool) string {

	// VertexAttribPointer calls use a number, but specify a pointer due to
	// historical reasons.
	if isParameters {
		if strings.Contains(mname, "VertexAttribPointer") {
			sig = strings.Replace(sig, "GLvoid *", "GLint64 ", 1)
			sig = strings.Replace(sig, "GLvoid*", "GLint64", 1)
		}
	} else {
		if strings.Contains(mname, "glVertexAttribPointer") {
			sig = strings.Replace(sig, "pointer);", "(const GLvoid *)pointer);", 1)
		}
		if strings.Contains(mname, "glGetVertexAttribPointerv") {
			sig = strings.Replace(sig, "pointer);", "(GLvoid **)pointer);", 1)
		}
	}

	// Special handling for functions that return strings.
	returnsString := []string{
		"glGetShaderInfoLog", "glGetProgramInfoLog", "glGetProgramPipelineInfoLog",
		"glGetActiveUniform", "glGetActiveUniformBlockName", "glGetActiveUniformName",
		"glGetDebugMessageLog", "glGetDebugMessageLogARB",
		"glGetActiveSubroutineName", "glGetObjectLabel",
		"glGetObjectPtrLabel", "glGetShaderSource",
	}
	if isParameters {
		for _, function := range returnsString {
			if mname == function {
				sig = strings.Replace(sig, "GLchar *", "GLubyte *", 1)
			}
		}
	} else {
		for _, function := range returnsString {
			if mname == function {
				regx, _ := regexp.Compile(" ([A-Za-z_]+)\\);")
				if loc := regx.FindStringIndex(sig); loc != nil {
					replace := sig[loc[0]+1 : loc[1]]
					sig = strings.Replace(sig, replace, "(GLchar *)"+replace, 1)
				}
			}
		}
	}
	return sig
}

// Make c or golang reserved words safe by inserting an underscore after
// the first letter in the word.  This leaves it readable and no longer
// in conflict.
func safeName(arg string) (safearg string) {
	safearg = arg
	reserved := map[string]string{
		"func":   "f_unc",
		"range":  "r_ange",
		"type":   "t_ype",
		"string": "s_tring",
		"map":    "m_ap",
		"near":   "n_ear",
		"far":    "f_ar",
	}
	if safe, ok := reserved[arg]; ok {
		safearg = safe
	}
	return
}

// Generate the c wrapper method that calls the function pointer.
func genCWrapper(creturn string, mname string, ptypes []string, pnames []string) (cwrap string) {
	cwrap = creturn + " wrap_" + mname + "("
	args := []string{}
	for cnt := 0; cnt < len(pnames); cnt++ {
		ctype, _, _ := convertType(ptypes[cnt])
		args = append(args, ctype+" "+pnames[cnt])
	}
	cwrap += strings.Join(args, ", ")
	cwrap += ") { "
	if creturn != "GLAPI void APIENTRY" {
		cwrap += "return "
	}
	cwrap += "(*pfn_" + mname + ")("
	cwrap += strings.Join(pnames, ", ")
	cwrap += "); }"

	// special handling for functions returning strings.
	if strings.Contains(cwrap, "GLAPI const GLubyte * APIENTRY") {
		cwrap = strings.Replace(cwrap, "GLAPI const GLubyte * APIENTRY", "GLAPI const char * APIENTRY", 1)
		cwrap = strings.Replace(cwrap, "return ", "return (char *)", 1)
	}

	// the rest of the special handling.
	cwrap = specialHandling(mname, cwrap, false)
	return
}

// Create the go func call from the given c method parts
// Eg. "func Enablei (target uint32, index uint32) ..."
func genGoFunc(creturn string, cname string, ptypes []string, pnames []string) (goapi string) {
	gname := strings.Replace(cname, "gl", "", 1)
	goapi = "func " + gname + " ("
	for cnt := 0; cnt < len(pnames); cnt++ {
		pname := safeName(pnames[cnt])
		if ptypes[cnt] != "" {
			_, gotype, _ := convertType(ptypes[cnt])
			goapi += pname + " " + gotype
			if cnt < len(pnames)-1 {
				goapi += ", "
			}
		}
	}
	goapi += ")"

	// the go return type
	_, goreturn, _ := convertType(creturn)
	goapi += " " + goreturn + " { "

	// tack on a return only if necessary
	if goreturn != "" {
		goapi += "return "
		if strings.Index(goreturn, "*") == 0 {
			goapi += "(" + goreturn + ")"
		} else {
			goapi += goreturn
		}
		goapi += "("
	}

	// generate the call to the c wrapper code.
	goapi += "C.wrap_" + cname + "( "
	for cnt := 0; cnt < len(pnames); cnt++ {
		if cnt != 0 {
			goapi += ", "
		}
		_, _, cgotype := convertType(ptypes[cnt])
		goapi += cgotype + "(" + safeName(pnames[cnt]) + ")"
	}
	if goreturn != "" {
		goapi += ")"
	}
	goapi += ") }"

	// generate code blocks to get proper boolean and string handling.
	goapi = genGoStrings(goapi)
	goapi = genGoBooleans(goapi)
	goapi = genGoStringArrays(goapi)
	return
}

// Ensure that strings are created and freed safely.
// Returns goapi unchanged if there were no strings.
func genGoStrings(goapi string) string {
	index := 1

	// process each string parameter.
	for strings.Contains(goapi, "(*C.char)(") {
		cstr := fmt.Sprintf("cstr%d", index)
		tag := strings.Split(goapi, "(*C.char)(")[1]
		tag = strings.Split(tag, ")")[0]
		insert := "{ \n" +
			"   %s := C.CString(%s)\n" +
			"   defer C.free(unsafe.Pointer(%s))"
		insert = fmt.Sprintf(insert, cstr, tag, cstr)

		// insert the code block.
		parts := strings.Split(goapi, "{ ")
		goapi = fmt.Sprintf("%s %s\n  %s", parts[0], insert, parts[1])
		replace := fmt.Sprintf("(*C.char)(%s)", tag)
		goapi = strings.Replace(goapi, replace, cstr, 1)
		index += 1
	}

	// change the return code for the functions that return strings.
	if strings.Contains(goapi, "return string") {
		goapi = strings.Replace(goapi, "return string", "return C.GoString", 1)
	}
	return goapi
}

// Use golang booleans rather than GL_TRUE and GL_FALSE.
// Returns goapi unchanged if there were no booleans.
func genGoBooleans(goapi string) string {
	btoks := strings.Split(goapi, ") bool")
	if len(btoks) == 1 {
		btoks := strings.Split(goapi, " bool")
		if len(btoks) > 1 {

			//  process each boolean parameter
			for cnt := 0; cnt < len(btoks)-1; cnt++ {
				btok := strings.TrimSpace(btoks[cnt])
				startTrim := strings.LastIndexAny(btok, "( ")
				tag := btok[startTrim+1 : len(btok)]

				// create the code block
				cb := fmt.Sprintf("tf%d", cnt+1)
				insert := "{ \n" +
					"  %s := FALSE\n" +
					"  if %s {\n" +
					"    %s = TRUE\n" +
					"  }"
				insert = fmt.Sprintf(insert, cb, tag, cb)

				// insert the code block
				parts := strings.Split(goapi, "{ ")
				goapi = fmt.Sprintf("%s %s\n  %s", parts[0], insert, parts[1])
				replace := fmt.Sprintf("C.uchar(%s)", tag)
				with := fmt.Sprintf("C.uchar(%s)", cb)
				goapi = strings.Replace(goapi, replace, with, 1)
			}
		}
	}

	// convert functions that return bool.
	goapi = strings.Replace(goapi, "return bool(", "return cbool(uint(", 1)
	if strings.Contains(goapi, "return cbool(uint(") {
		goapi = strings.Replace(goapi, ")) }", "))) }", 1)
	}
	return goapi
}

// internal type to denote a string array.  This is used as function signatures
// are built, and then it is used here to generate the string array code.
const strArray = "_STR_ARRAY_"

// Handle input arrays of strings by properly allocating and freeing them.
// There is at most one string array in any OpenGL method call.
// Returns goapi unchanged if there were no string arrays.
func genGoStringArrays(goapi string) string {
	if strings.Contains(goapi, "[]string") {
		// get the tag. fail sooner than later if something doesn't match.
		regx, _ := regexp.Compile(strArray + "\\(([A-Za-z_]+)\\)")
		if loc := regx.FindStringIndex(goapi); loc != nil {
			rep := goapi[loc[0]:loc[1]]
			tag := goapi[loc[0]+len(strArray)+1 : loc[1]-1]

			// create the code block to be inserted.
			insert := "{\n" +
				"   cstrings := C.newStringArray(C.int(len(%s)))\n" +
				"   defer C.freeStringArray(cstrings, C.int(len(%s)))\n" +
				"   for cnt, str := range %s {\n" +
				"      C.assignString(cstrings, C.CString(str), C.int(cnt))\n" +
				"   }"
			insert = fmt.Sprintf(insert, tag, tag, tag)

			// insert the code block.
			parts := strings.Split(goapi, "{ ")
			goapi = fmt.Sprintf("%s%s\n   %s", parts[0], insert, parts[1])
			goapi = strings.Replace(goapi, rep, "cstrings", 1)
		}
	}
	return goapi
}

// look up the type and and keep the typemap updated with which
// types get referenced.
func convertType(ctype string) (ct string, gt string, cgt string) {
	ct, gt, cgt = "", "", ""
	if tm, ok := typemap[ctype]; ok {
		ct, gt, cgt = tm.ctype, tm.gotype, tm.cgotype
		tm.used = true
		typemap[ctype] = tm
	}
	return
}

// Convert a specification c constant into a go constant.
// The go constant is expected to be put into a const block.
// Special cases are included as appropriate for the
// specification (not all possible c-syntax cases are handled).
func genGoConst(cconst string) (goconst string) {
	fields := strings.Fields(cconst)
	name := strings.Replace(fields[1], "GL_", "", 1)
	value := fields[2]
	// handle unsigned byte strings.
	if value == "0xFFFFFFFFu" {
		value = "0xFFFFFFFF"
	}
	// unsigned 64 bit long byte strings.
	if value == "0xFFFFFFFFFFFFFFFFull" {
		value = "0xFFFFFFFFFFFFFFFF"
	}
	// handle referencing defined constants
	value = strings.Replace(value, "GL_", "", 1)
	goconst = name + " = " + value
	return
}

// Add the following C code as comments immediately preceding the
// import "C" statement.  The typedefs of the spec are appended
// to this preamble.
var cPreamble = []string{
	"// #cgo darwin  LDFLAGS: -framework OpenGL", // needed to compile on OSX
	"// #cgo linux   LDFLAGS: -lGL",              // only tested on Ubuntu
	"// #cgo windows LDFLAGS: -lopengl32",
	"// ",
	"// #include <stdlib.h>",
	"// #if defined(__APPLE__)",
	"// #include <dlfcn.h>", // for getting pointer to methods.
	"// #elif defined(_WIN32)",
	"// #define WIN32_LEAN_AND_MEAN 1",
	"// #include <windows.h>",
	"// #else",
	"// #include <X11/Xlib.h>", // linux is all X11 based.
	"// #include <GL/glx.h>",
	"// #endif",
	"// ",
	"// #ifdef _WIN32",
	"// static HMODULE hmod = NULL;",
	"// #endif",
	"// ",
	"// /* Helps bind function pointers to c functions. */",
	"// static void* bindMethod(const char* name) {",
	"// #ifdef __APPLE__",
	"// 	return dlsym(RTLD_DEFAULT, name);",
	"// #elif _WIN32",
	"// 	void* pf = wglGetProcAddress((LPCSTR)name);",
	"// 	if(pf) {",
	"// 		return pf;",
	"// 	}",
	"// 	if(hmod == NULL) {",
	"// 		hmod = LoadLibraryA(\"opengl32.dll\");",
	"// 	}",
	"// 	return GetProcAddress(hmod, (LPCSTR)name);",
	"// #else",
	"// 	return glXGetProcAddress((const GLubyte*)name);",
	"// #endif",
	"// }",
	"// ",
	"// /* Helper method for string arrays */",
	"// static char**newStringArray(int numberOfStrings) {",
	"//     return calloc(sizeof(char*), numberOfStrings);",
	"// }",
	"// ",
	"// /* Helper method for string arrays */",
	"// static void assignString(char **stringArray, char *string, int index) {",
	"//     stringArray[index] = string;",
	"// }",
	"// ",
	"// /* Helper method for string arrays */",
	"// static void freeStringArray(char **stringArray, int numberOfStrings) {",
	"//     int cnt;",
	"//     for (cnt = 0; cnt < numberOfStrings; cnt++)",
	"//         free(stringArray[cnt]);",
	"//     free(stringArray);",
	"// }",
	"// ",
	// the appendPreamble function copies specification typedefs after this line.
}

// set when lines should be appended to the preamble.
var copyLines = false

// Append specification lines to the C language preamble.
// This is needed to get the OpenGL typedefs.
//
// Expected to be called with each line in the OpenGL specification.
// This starts copying lines when an "onlines" is detected and stops copying
// when one of the "offlines" are detected..
func appendPreamble(line string) {
	var onlines = []string{"#ifndef APIENTRY", "#include <stddef.h>"}
	var offlines = []string{"/******************", "#ifndef GL_VERSION_1_0"}
	if !copyLines {
		for _, online := range onlines {
			if strings.Index(line, online) == 0 {
				copyLines = true
			}
		}
	} else {
		for _, offline := range offlines {
			if strings.Index(line, offline) == 0 {
				copyLines = false
			}
		}
	}
	if copyLines {
		cPreamble = append(cPreamble, "// "+line)
	}
}

// ================
// Generate output.
// ================

// Generates the package based on the specification data.
func genPackage(pkg, spec string, groups map[string]grouping) {
	groupNames := sortSpecGroupings(groups)
	path, _ := filepath.Abs(pkg)
	genTests(path)
	gout := genFile(path)
	defer gout.Close()
	genHeader(gout, path, spec)

	// end:: generate C code
	genCCode(gout, groups, groupNames)

	// end:: generate Go code
	genGoCode(gout, groups, groupNames)
}

// Sort the spec data group labels.
func sortSpecGroupings(groups map[string]grouping) (sortedGroupNames []string) {
	sortedGroupNames = make([]string, len(groups))
	cnt := 0
	for groupName, _ := range groups {
		sortedGroupNames[cnt] = groupName
		cnt++
	}
	sort.Strings(sortedGroupNames)
	return
}

// Create the output file.  This will be a go (cgo) file.
func genFile(path string) (gout *os.File) {
	gout, err := os.Create(path + "/gl.go")
	if err != nil {
		os.Exit(-1)
	}
	return
}

// Create the very first lines of comments and package declaration
// in the go code file.
func genHeader(gout *os.File, path, spec string) {
	pkgName := filepath.Base(path)
	fmt.Fprintf(gout, "// Package %s provides a wrapper 3D graphics library\n", pkgName)
	fmt.Fprintf(gout, "// that was auto-generated from %s.\n", spec)
	fmt.Fprintf(gout, "// The official OpenGL documentation for any of the constants\n")
	fmt.Fprintf(gout, "// or methods can be found online. Just prepend \"GL_\"\n")
	fmt.Fprintf(gout, "// to the function or constants names in this package.\n")
	fmt.Fprintf(gout, "package %s\n\n", pkgName)
}

// ===============
// Generate C code
// ===============

// Generate all the commented C code that must appear directly before the
// import "C" statement.
func genCCode(gout *os.File, groups map[string]grouping, groupNames []string) {
	genCDefs(gout)
	genCBindings(gout, groups, groupNames)
	genCInit(gout, groups, groupNames)
}

// Output the C definitions from the preamble.
func genCDefs(gout *os.File) {
	for _, line := range cPreamble {
		fmt.Fprintf(gout, "%s\n", line)
	}
}

// Put the c binding code in the preamble of the go code.
func genCBindings(gout *os.File, groups map[string]grouping, groupNames []string) {
	for _, group := range groupNames {
		spec := groups[group]
		if len(spec.methods) > 0 {
			fmt.Fprintf(gout, "//\n// // %s\n", group)
			sortedMethods := spec.sortedMethodNames()
			for _, key := range sortedMethods {
				method := spec.methods[key]
				fmt.Fprintf(gout, "//\n")
				fmt.Fprintf(gout, "// %s\n", method.cmethodPtr)
				fmt.Fprintf(gout, "// %s\n", method.cwrapper)
			}
		}
	}
}

// Generate the C Init code that performs the actual binding.
func genCInit(gout *os.File, groups map[string]grouping, groupNames []string) {
	fmt.Fprintf(gout, "//\n")
	fmt.Fprintf(gout, "// void init() {\n")
	for _, group := range groupNames {
		spec := groups[group]
		if len(spec.methods) > 0 {
			sortedMethods := spec.sortedMethodNames()
			for _, mname := range sortedMethods {
				fmt.Fprintf(gout, "//   pfn_%s = bindMethod(\"%s\");\n", mname, mname)
			}
		}
	}
	fmt.Fprintf(gout, "// }\n")
	fmt.Fprintf(gout, "//\n")
}

// ================
// Generate Go code
// ================

// Generate the Go code.  This outputs the go binding methods and the
// go constants.
func genGoCode(gout *os.File, groups map[string]grouping, groupNames []string) {
	genGoImports(gout)
	genUtilityMethods(gout)
	genGoTypedefs(gout)

	// create the go code organized by specificaiton grouping.
	for _, group := range groupNames {
		groupData := groups[group]
		genGoConstants(gout, group, groupData)
		genGoMethods(gout, group, groupData)
	}
	genGoInitMethod(gout)
	genGoBindingReport(gout, groups, groupNames)
}

// The spec has a few special types that need an extra type mapping for
// golang to be aware of them.  Just ensure that the type is used in the spec
// before including it.
func genGoTypedefs(gout *os.File) {
	specials := []string{"GLvoid*", "GLsync", "struct _cl_context*", "struct _cl_event*", "GLDEBUGPROCARB", "GLDEBUGPROC"}
	fmt.Fprintf(gout, "// Special type mappings\n")
	fmt.Fprintf(gout, "type (\n")
	for _, special := range specials {
		if typemap[special].used {
			gotype := strings.Trim(typemap[special].gotype, "()*")
			cgotype := strings.Trim(typemap[special].cgotype, "()*")
			fmt.Fprintf(gout, "  %s %s\n", gotype, cgotype)
		}
	}
	fmt.Fprintf(gout, ")\n ")
}

// Generate the go import statements.
// Import "C" must come immediately after the c code comments and preamble.
func genGoImports(gout *os.File) {
	fmt.Fprintf(gout, "import \"C\"\n")
	fmt.Fprintf(gout, "import \"unsafe\"\n")
	fmt.Fprintf(gout, "import \"fmt\"\n")
	fmt.Fprintf(gout, "\n")
}

// Generate the exposed initialition function.  This calls the c init that
// binds the available OpenGL methods to the c function pointers.
func genGoInitMethod(gout *os.File) {
	fmt.Fprintf(gout, "// bind the methods to the function pointers\n")
	fmt.Fprintf(gout, "func Init() {\n")
	fmt.Fprintf(gout, "   C.init()\n")
	fmt.Fprintf(gout, "}\n")
	fmt.Fprintf(gout, "\n")
}

// Create the go constant definitions for a group.
func genGoConstants(gout *os.File, group string, groupData grouping) {
	if len(groupData.goconsts) > 0 {
		fmt.Fprintf(gout, "\n// %s", group)
		fmt.Fprintf(gout, "\nconst (\n")
		for _, goconst := range groupData.goconsts {
			fmt.Fprintf(gout, "   %s\n", goconst)
		}
		fmt.Fprintf(gout, ")\n ")
	}
}

// Create the go methods for a group.
func genGoMethods(gout *os.File, group string, groupData grouping) {
	if len(groupData.methods) > 0 {
		fmt.Fprintf(gout, "\n// %s\n", group)
		sortedMethods := groupData.sortedMethodNames()
		for _, key := range sortedMethods {
			method := groupData.methods[key]
			fmt.Fprintf(gout, "%s\n", method.gofunc)
		}
	}
}

// Create the test methods.
func genTests(path string) {
	pkgName := filepath.Base(path)
	gout, err := os.Create(path + "/gl_test.go")
	if err != nil {
		os.Exit(-1)
	}
	fmt.Fprintf(gout, "package %s\n", pkgName)
	fmt.Fprintf(gout, "\n")
	fmt.Fprintf(gout, "import \"testing\"\n")
	fmt.Fprintf(gout, "\n")
	fmt.Fprintf(gout, "// The test passes if the binding layer can initialize without crashing.\n")
	fmt.Fprintf(gout, "func TestInit(t *testing.T) {\n")
	fmt.Fprintf(gout, "    Init()\n")
	fmt.Fprintf(gout, "}\n")
}

// Create some go utility methods.
func genUtilityMethods(gout *os.File) {
	fmt.Fprintf(gout, "// convert a GL uint boolean to a go bool\n")
	fmt.Fprintf(gout, "func cbool(glbool uint) bool {\n")
	fmt.Fprintf(gout, "    return glbool == TRUE\n")
	fmt.Fprintf(gout, "}\n")
	fmt.Fprintf(gout, "\n")
}

// Create the binding report code. The code shows which OpenGL method
// is bound on the current platform.
func genGoBindingReport(gout *os.File, groups map[string]grouping, groupNames []string) {
	fmt.Fprintf(gout, "// Show which function pointers are bound\n")
	fmt.Fprintf(gout, "func BindingReport() (report []string) {\n")
	fmt.Fprintf(gout, "   report = []string{}\n")
	for _, group := range groupNames {
		groupData := groups[group]
		if len(groupData.methods) > 0 {
			fmt.Fprintf(gout, "   report = append(report, \"%s\")\n", group)
			sortedMethods := groupData.sortedMethodNames()
			for _, mname := range sortedMethods {
				fmt.Fprintf(gout, "   report = append(report, isBound(unsafe.Pointer(C.pfn_%s), \"%s\"))\n", mname, mname)
			}
		}
	}
	fmt.Fprintf(gout, "return\n")
	fmt.Fprintf(gout, "}\n")

	// report helper function.
	fmt.Fprint(gout, "func isBound(pfn unsafe.Pointer, fn string) string {\n")
	fmt.Fprint(gout, "inc := \"-\"\n")
	fmt.Fprint(gout, "if pfn != nil {\n")
	fmt.Fprint(gout, "   inc = \"+\"\n")
	fmt.Fprint(gout, " }\n")
	fmt.Fprint(gout, "return fmt.Sprintf(\"   [%s] %s\", inc, fn)\n")
	fmt.Fprintf(gout, "}\n\n")
}
