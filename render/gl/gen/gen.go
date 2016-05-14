// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// Package gen is used to generate golang OpenGL bindings using an OpenGL
// specification header file. The 4.3 specification header file is from:
//     www.opengl.org/registry/api/glcorearb.h
// Running gl/gen will create a go package. Usage:
//     gen -p output-directory -s opengl-spec
// Design notes are included in the code. Generation is expected to be
// infrequenty and the generated gl package should be source controlled.
//
// Package gen is provided as part of the vu (virtual universe) 3D engine.
package main

// Thanks to https://github.com/chsc/gogl for the idea of generating the
// bindings from a specification.
//
// Design Notes: Essentially straight line code to get information in one
// format, the OpenGL Specification, to another format, the Go language
// bindings for OpenGL.
//
// Maintenance Notes: The code has just enough smarts to process the current
// specification. The idea is to only generalize the code where necessary without
// adding undo bulk or dependencies to the code. Anything that would reduce the
// amount of code and improve readability would be appreciated and should be done.
//
// FUTURE: Consider using the XML specification instead of the header file.

// Use 'go generate' to rebuild the OpenGL bindings.
//go:generate go run gen.go -p ../../gl -s glcorearb.h-v4_3
//go:generate go fmt ../../gl

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Generate golang OpenGL bindings from a OpenGL specification header.  The
// two main calls are:
//    parseSpecification - read in the spec.
//    genBindings - write out the bindings.
func main() {
	flag.Parse()
	if os.MkdirAll(*pkg, os.FileMode(0755)) == nil {
		if specHeader, err := os.Open(*spec); err == nil {
			defer specHeader.Close()
			functions, constants, typedefs := parseSpecification(specHeader)
			specInfo, _ := specHeader.Stat()
			genBindings(*pkg, specInfo.Name(), functions, constants, typedefs)
			return
		}
	}
	flag.Usage()
}

// Input defaults that work when running from this directory.
// The build overrides as necessary.
var pkg = flag.String("p", "../../gl", "the output package name")
var spec = flag.String("s", "glcorearb.h-v4.3", "the OpenGL specification header file")

// parseSpecification reads the OpenGL specification lines into three different
// string slices.
func parseSpecification(spec *os.File) (functions, constants, typedefs []string) {
	scanner := bufio.NewScanner(spec)
	scanner.Split(bufio.ScanLines)
	linenum := 0
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		// functions are the OpenGL function API's, one per line.
		if len(fields) > 0 && fields[0] == "GLAPI" {
			functions = append(functions, line)
		}

		// constants are the OpenGL constant definitions.
		if len(fields) == 3 && fields[0] == "#define" && len(fields[1]) > 4 && fields[1][0:3] == "GL_" {
			if !strings.Contains(fields[1], "GL_ARB") && !strings.Contains(fields[1], "GL_VERSION_") {
				constants = append(constants, line)
			}
		}

		// typedefs are copied straight out of the header file and reproduced in the bindings.
		if linenum >= 62 && linenum <= 90 || linenum >= 2600 && linenum <= 2702 {
			typedefs = append(typedefs, line)
		}
		linenum++
	}
	sort.Strings(functions)
	return
}

// genBindings creates the OpenGL golang binding code.
func genBindings(pkg, spec string, functions, constants, typedefs []string) {
	gout, pkgname := genOutputFile(pkg)
	genPreamble(gout, pkgname, typedefs)
	genCwrappers(gout, functions)
	genCinit(gout, functions)
	genGoPreamble(gout)
	genGoConstants(gout, constants)
	genGoFunctions(gout, functions)
	genBindingReport(gout, functions)
	gout.Close()
}

// genOutput creates the output file in the given package.
func genOutputFile(pkg string) (gout *os.File, pkgname string) {
	path, _ := filepath.Abs(pkg)
	pkgname = filepath.Base(path)
	gout, err := os.Create(path + "/" + pkgname + ".go")
	if err != nil {
		os.Exit(-1)
	}
	return gout, pkgname
}

// genPreamble dumps the initial static code into the file. This includes
// a few comments and the start of the CGo code block.
func genPreamble(gout *os.File, pkgname string, typedefs []string) {
	fmt.Fprintf(gout, "// Package %s provides golang bindings for OpenGL\n", pkgname)
	fmt.Fprintf(gout, "// The bindings were generated from OpenGL spec %s.\n", *spec)
	fmt.Fprintf(gout, "// The official OpenGL documentation for any of the constants\n")
	fmt.Fprintf(gout, "// or methods can be found online. Just prepend \"GL_\"\n")
	fmt.Fprintf(gout, "// to the function or constant names in this package.\n//\n")
	fmt.Fprintf(gout, "// Package %s is provided as part of the vu (virtual universe) 3D engine.\n", pkgname)
	fmt.Fprintf(gout, "package %s\n\n", pkgname)
	// .
	for _, line := range cPreamble {
		fmt.Fprintf(gout, "%s\n", line)
	}
	for _, line := range typedefs {
		fmt.Fprintf(gout, "// %s\n", line)
	}
}

// cPreamble is the initial static cgo code.  The typedefs and generated cgo
// code will be placed after this block.
//
// Needed to create a link for the latest opengl library as follows:
//     sudo ln -s /usr/lib/nvidia-319-updates/libGL.so.1 /usr/lib/libGL.so
var cPreamble = []string{
	"// #cgo darwin  LDFLAGS: -framework OpenGL", // needed to compile on OSX
	"// #cgo linux   LDFLAGS: -lGL -ldl",         // only tested on Ubuntu
	"// #cgo windows LDFLAGS: -lopengl32",
	"// ",
	"// #include <stdlib.h>",
	"// #if defined(__APPLE__)",
	"// #include <dlfcn.h>", // for getting pointer to methods.
	"// #elif defined(_WIN32)",
	"// #define WIN32_LEAN_AND_MEAN 1",
	"// #include <windows.h>",
	"// #else",
	"// #include <dlfcn.h>",
	"// #endif",
	"// ",
	"// #ifdef _WIN32",
	"// static HMODULE hmod = NULL;",
	"// #elif !defined __APPLE__",
	"// static void* plib = NULL;",
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
	"// 	if(plib == NULL) {",
	"// 		plib = dlopen(\"libGL.so\", RTLD_LAZY);",
	"// 	}",
	"// 	return dlsym(plib, name);",
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
	"//     {",
	"//         free(stringArray[cnt]);",
	"//     }",
	"//     free(stringArray);",
	"// }",
	"// ", // the specification typedefs are output after this line.
}

// genCwrappers outputs cgo OpenGL function pointers and the wrapper functions.
func genCwrappers(gout *os.File, functions []string) {
	for _, apiLine := range functions {
		fname, rettype, plist, pnames, gltypes := splitLine(apiLine)

		// The cgo pointer definition of the OpenGL function.
		fmt.Fprintf(gout, "// %s (APIENTRYP pfn_%s)(%s);\n", rettype, fname, plist)

		// The cgo wrapper wraps the pointer definition.
		apiLine = alterSpec(apiLine)
		fname, rettype, plist, pnames, ptypes := splitLine(apiLine)
		wret, wparms, wnames := cwrapper(rettype, pnames, ptypes, gltypes)
		ret := "return"
		switch rettype {
		case "void":
			ret = ""
		case "const GLchar *":
			ret = "return (char *)"
		}
		fmt.Fprintf(gout, "// %s wrap_%s(%s) { %s (*pfn_%s)(%s); }\n//\n", wret, fname, wparms, ret, fname, wnames)
	}
}

// genCinit() outputs the Cgo Init() method that initializes OpenGL function pointers.
func genCinit(gout *os.File, functions []string) {
	fmt.Fprintf(gout, "//\n// void init() {\n")
	for _, apiLine := range functions {
		fname, _, _, _, _ := splitLine(apiLine)
		fmt.Fprintf(gout, "//   pfn_%s = bindMethod(\"%s\");\n", fname, fname)
	}
	fmt.Fprintf(gout, "// }\n//\n")
}

// genGoPreamble is static Go code immediately following the cgo block. The first
// line must be "import C" for cgo. A couple of small helper functions are created as well.
func genGoPreamble(gout *os.File) {
	fmt.Fprintf(gout, "import \"C\"\n")
	fmt.Fprintf(gout, "import \"unsafe\"\n\n")
	fmt.Fprintf(gout, "import \"fmt\"\n\n")
	fmt.Fprintf(gout, "// convert a GL uint boolean to a go bool\n")
	fmt.Fprintf(gout, "func cbool(glbool uint) bool {\n")
	fmt.Fprintf(gout, "	return glbool == TRUE\n")
	fmt.Fprintf(gout, "}\n\n")
	fmt.Fprintf(gout, "// Special type mappings\n")
	fmt.Fprintf(gout, "type (\n")
	fmt.Fprintf(gout, "	Pointer      unsafe.Pointer\n")
	fmt.Fprintf(gout, "	Sync         C.GLsync\n")
	fmt.Fprintf(gout, "	clContext    C.struct_Cl_context\n")
	fmt.Fprintf(gout, "	clEvent      C.struct_Cl_event\n")
	fmt.Fprintf(gout, "	DEBUGPROCARB C.GLDEBUGPROCARB\n")
	fmt.Fprintf(gout, "	DEBUGPROC    C.GLDEBUGPROC\n")
	fmt.Fprintf(gout, ")\n\n")
	fmt.Fprintf(gout, "// bind the methods to the function pointers\n")
	fmt.Fprintf(gout, "func Init() {\n")
	fmt.Fprintf(gout, "   C.init()\n")
	fmt.Fprintf(gout, "}\n\n")
}

// genGoConstants turns all the OpenGL constants into Go constants by stripping
// off the leading "GL_" (the constants are expected to be accessed as "gl.")
func genGoConstants(gout *os.File, constants []string) {
	fmt.Fprintf(gout, "const (\n")
	for _, line := range constants {
		fields := strings.Fields(line)
		fmt.Fprintf(gout, "   %s = %s\n", fields[1][3:], alterConstant(fields[2]))
	}
	fmt.Fprintf(gout, ")\n\n")
}

// genGoFunctions outputs the golang functions that call the cgo wrappers (that in turn
// call the cgo function pointers).
func genGoFunctions(gout *os.File, functions []string) {
	for _, apiLine := range functions {
		apiLine = alterSpec(apiLine)
		mname, goparms, goret, cgoparms := gowrapper(apiLine)

		// the space before the newline is significant for the mod methods.
		goapi := fmt.Sprintf("func %s(%s) %s{ \n", mname, goparms, goret)
		switch goret {
		case "":
			goapi += fmt.Sprintf("   C.wrap_gl%s(%s)\n}\n", mname, cgoparms)
		case "bool":
			goapi += fmt.Sprintf("   return cbool(uint(C.wrap_gl%s(%s)))\n}\n", mname, cgoparms)
		default:
			goapi += fmt.Sprintf("   return %s(C.wrap_gl%s(%s))\n}\n", goret, mname, cgoparms)
		}

		// hack the generated goapi, inserting translation code blocks so that
		// proper golang types are used for booleans and strings.
		goapi = modGoBooleans(goapi)
		goapi = modGoStrings(goapi)
		goapi = modGoStringArrays(goapi)
		fmt.Fprintf(gout, "%s", goapi)
	}
}

// genBindingReport outputs some golang code that can be used to check which
// functions have been bound.
func genBindingReport(gout *os.File, functions []string) {
	fmt.Fprintf(gout, "// Show which function pointers are bound\n")
	fmt.Fprintf(gout, "func BindingReport() (report []string) {\n")
	fmt.Fprintf(gout, "   report = []string{}\n")
	for _, apiLine := range functions {
		fname, _, _, _, _ := splitLine(apiLine)
		fmt.Fprintf(gout, "   report = append(report, isBound(unsafe.Pointer(C.pfn_%s), \"%s\"))\n", fname, fname)
	}
	fmt.Fprintf(gout, "   return\n}\n")

	// report helper function.
	fmt.Fprint(gout, "func isBound(pfn unsafe.Pointer, fn string) string {\n")
	fmt.Fprint(gout, "   inc := \" \"\n")
	fmt.Fprint(gout, "   if pfn != nil {\n")
	fmt.Fprint(gout, "      inc = \"+\"\n")
	fmt.Fprint(gout, "   }\n")
	fmt.Fprintf(gout, "   return fmt.Sprintf(\"   [%s] %s\", inc, fn)\n", "%s", "%s") // go vet compliant.
	fmt.Fprint(gout, "}\n\n")
}

// binding generation code
// ============================================================================
// parsing and generation helper functions and types.

// splitLine breaks an OpenGL API definition into its component parts.
func splitLine(line string) (fname, retype, plist string, pnames, ptypes []string) {
	fields := strings.Fields(line)
	retype = fields[1]
	fname = fields[3]
	if fname == "*" {
		retype = fields[1] + " " + fields[2] + " " + fields[3]
		fname = fields[5]
	}

	// split the parameters into types and names.
	plist = line[strings.Index(line, "(")+1 : strings.Index(line, ")")]
	parms := strings.Split(plist, ",")
	for _, parm := range parms {
		if split := strings.LastIndexAny(parm, " *"); split >= 0 {
			ptypes = append(ptypes, strings.TrimSpace(parm[0:split+1]))
			pnames = append(pnames, alterArg(parm[split+1:]))
		}
	}
	return
}

// cwrapper generates the string fields needed for the cwrapper definition.
func cwrapper(rettype string, pnames, ptypes, gltypes []string) (wret, wparms, wnames string) {
	wret = "GLAPI void APIENTRY"
	rettype = strings.Replace(rettype, " ", "", -1)
	if types, ok := typemap[rettype]; ok {
		wret = "GLAPI " + types[0] + " APIENTRY"
	}

	for cnt, gltype := range ptypes {
		tkey := strings.Replace(gltype, " ", "", -1)
		wparms += typemap[tkey][0] + " " + pnames[cnt] + ", "
		wnames += "(" + gltypes[cnt] + ")" + pnames[cnt] + ", "
	}
	if len(wparms) > 0 {
		wparms = wparms[0 : len(wparms)-2]
		wnames = wnames[0 : len(wnames)-2]
	}
	return
}

// gowrapper generates the string fields needed for the gowrapper definition.
func gowrapper(apiLine string) (mname, goparms, goret, cgoparms string) {
	fname, rettype, _, pnames, ptypes := splitLine(apiLine)
	mname = fname[2:]
	for cnt, gltype := range ptypes {
		tkey := strings.Replace(gltype, " ", "", -1)
		goparms += pnames[cnt] + " " + typemap[tkey][1] + ", "
		cgoparms += typemap[tkey][2] + "(" + pnames[cnt] + "), "
	}
	if len(goparms) > 0 {
		goparms = goparms[0 : len(goparms)-2]
		cgoparms = cgoparms[0 : len(cgoparms)-2]
	}
	tkey := strings.Replace(rettype, " ", "", -1)
	if types, ok := typemap[tkey]; ok {
		goret = types[1]
	}
	return
}

// alterSpec changes the OpenGL API for some OpenGL functions where the use of the parameters
// does not really align with the c-specification (making the binding translation awkward
// unless changed). Note that the pointer definition always matches the original spec with.
// the alterations being made afterwards for the other definitions (wrapper, go, cgo). Casts
// are used as appropriate to ensure everything aligns with the original spec API.
func alterSpec(apiLine string) string {
	fname, _, _, _, _ := splitLine(apiLine)
	switch fname {

	// OpenGL functions returning pointer to bytes that need to be treated as strings.
	case "glGetString", "glGetStringi":
		return strings.Replace(apiLine, "GLubyte *", "GLchar *", 1)

	// OpenGL functions with string parameters used to pass back strings
	// where the parameters need to be treated as pointer to bytes.
	case "glGetShaderInfoLog", "glGetProgramInfoLog", "glGetProgramPipelineInfoLog", "glGetActiveUniform",
		"glGetActiveUniformBlockName", "glGetActiveUniformName", "glGetDebugMessageLog", "glGetDebugMessageLogARB",
		"glGetActiveSubroutineName", "glGetObjectLabel", "glGetObjectPtrLabel", "glGetShaderSource",
		"glGetActiveAttrib":
		return strings.Replace(apiLine, "GLchar *", "GLubyte *", 1)

	// OpenGL functions where a pointer parameter is being used as a value (for historical reasons).
	case "glVertexAttribPointer":
		return strings.Replace(apiLine, "GLvoid *", "GLint64 ", 1)
	case "glGetVertexAttribPointerv":
		return strings.Replace(apiLine, "GLvoid* *", "GLint64 *", 1)
	case "glDrawElements", "glDrawElementsBaseVertex", "glDrawElementsInstanced",
		"glDrawElementsInstancedBaseInstance", "glDrawElementsInstancedBaseVertex",
		"glDrawElementsInstancedBaseVertexBaseInstance", "glDrawRangeElements",
		"glDrawRangeElementsBaseVertex":
		return strings.Replace(apiLine, "const GLvoid *indices", "GLintptr indicies", 1)
	}
	return apiLine
}

// alterConstant changes some C-constants to be valid Go-constants.
func alterConstant(constant string) string {
	switch constant {
	case "0xFFFFFFFFu":
		return "0xFFFFFFFF"
	case "0xFFFFFFFFFFFFFFFFull":
		return "0xFFFFFFFFFFFFFFFF"
	}
	if strings.Contains(constant, "GL_") {
		return constant[3:]
	}
	return constant
}

// alterArg makes parameter arguments safe in the case where OpenGL has used
// golang reserved words as parameters.
func alterArg(arg string) (safearg string) {
	safearg = strings.TrimSpace(arg)
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

// modGoBooleans modifies the basic golang function so that it uses bool rather
// than GL_TRUE and GL_FALSE.  This inserts a codeblock for each boolean parameter
// and updates the parameters lists to match. No changes are made if no bool parms.
func modGoBooleans(goapi string) string {
	btoks := strings.Split(goapi, ") bool")
	if len(btoks) == 1 {
		btoks := strings.Split(goapi, " bool")
		if len(btoks) > 1 {

			// process each boolean parameter
			for cnt := 0; cnt < len(btoks)-1; cnt++ {
				btok := strings.TrimSpace(btoks[cnt])
				startTrim := strings.LastIndexAny(btok, "( ")
				tag := btok[startTrim+1:]

				// create the code block
				cb := fmt.Sprintf("tf%d", cnt+1)
				insert := "{ \n" +
					"   %s := FALSE\n" +
					"   if %s {\n" +
					"      %s = TRUE\n" +
					"   }"
				insert = fmt.Sprintf(insert, cb, tag, cb)

				// insert the code block
				parts := strings.Split(goapi, "{ ")
				goapi = fmt.Sprintf("%s %s  %s", parts[0], insert, parts[1])
				replace := fmt.Sprintf("C.uchar(%s)", tag)
				with := fmt.Sprintf("C.uchar(%s)", cb)
				goapi = strings.Replace(goapi, replace, with, 1)
			}
		}
	}
	return goapi
}

// modGoStrings modifies the basic golang function to ensure that strings are
// created and freed safely. This inserts a codeblock for each string parameter
// and updates the parameters to match. No changes are made if no string parms.
func modGoStrings(goapi string) string {
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
		goapi = fmt.Sprintf("%s %s  %s", parts[0], insert, parts[1])
		replace := fmt.Sprintf("(*C.char)(%s)", tag)
		goapi = strings.Replace(goapi, replace, cstr, 1)
		index++
	}

	// change the return code for the functions that return strings.
	if strings.Contains(goapi, "return string") {
		goapi = strings.Replace(goapi, "return string", "return C.GoString", 1)
	}
	return goapi
}

// modGoStringArrays modifies the basic golang function to ensure that strings
// arrays are created and freed safely. This inserts a codeblock for each string
// parameter and updates the parameters to match. No changes are made if no string parms.
func modGoStringArrays(goapi string) string {
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
			goapi = fmt.Sprintf("%s%s   %s", parts[0], insert, parts[1])
			goapi = strings.Replace(goapi, rep, "cstrings", 1)
		}
	}
	return goapi
}

// internal type to denote a string array. This is used in type map and
// modGoStringArrays in order to generate string array binding code.
const strArray = "_STR_ARRAY_"

// typemap turns OpenGL types into the type expressions for the wrapper, cgo,
// and go functions. Not every type is represented, only those used in the spec.
// Map values are c-wrapper-type, go-type, cgo-type
var typemap = map[string][]string{
	"GLchar*":            {"char*", "string", "(*C.char)"},
	"GLubyte":            {"unsigned char", "uint8", "C.uchar"},
	"GLubyte*":           {"unsigned char*", "*uint8", "(*C.uchar)"},
	"GLboolean":          {"unsigned char", "bool", "C.uchar"},
	"GLboolean*":         {"unsigned char*", "*uint8", "(*C.uchar)"},
	"GLenum":             {"unsigned int", "uint32", "C.uint"},
	"GLfloat":            {"float", "float32", "C.float"},
	"GLfloat*":           {"float*", "*float32", "(*C.float)"},
	"GLint":              {"int", "int32", "C.int"},
	"GLint*":             {"int*", "*int32", "(*C.int)"},
	"GLintptr":           {"long long", "int64", "C.longlong"},
	"GLuint":             {"unsigned int", "uint32", "C.uint"},
	"GLuint*":            {"unsigned int*", "*uint32", "(*C.uint)"},
	"GLsizei":            {"int", "int32", "C.int"},
	"GLsizeiptr":         {"long long", "int64", "C.longlong"},
	"GLbitfield":         {"unsigned int", "uint32", "C.uint"},
	"GLdouble":           {"double", "float64", "C.double"},
	"GLvoid*":            {"void*", "Pointer", "unsafe.Pointer"},
	"GLvoid**":           {"void**", "*Pointer", "(*unsafe.Pointer)"},
	"GLdouble*":          {"double*", "*float64", "(*C.double)"},
	"GLsizei*":           {"int*", "*int32", "(*C.int)"},
	"GLenum*":            {"unsigned int*", "*uint32", "(*C.uint)"},
	"GLshort":            {"short ", "int16", "C.short"},
	"GLushort*":          {"unsigned short*", "*uint16", "(*C.ushort)"},
	"GLsync":             {"GLsync", "Sync", "C.GLsync"},
	"GLint64*":           {"long long*", "*int64", "(*C.longlong)"},
	"GLuint64":           {"unsigned long long", "uint64", "C.ulonglong"},
	"GLuint64*":          {"unsigned long long*", "*uint64", "(*C.ulonglong)"},
	"constvoid*":         {"const void*", "Pointer", "unsafe.Pointer"},
	"constGLfloat*":      {"const float*", "*float32", "(*C.float)"},
	"constGLint*":        {"const int*", "*int32", "(*C.int)"},
	"constGLsizei*":      {"const int*", "*int32", "(*C.int)"},
	"constGLvoid*":       {"const void*", "Pointer", "unsafe.Pointer"},
	"constGLuint*":       {"const unsigned int*", "*uint32", "(*C.uint)"},
	"constGLvoid*const*": {"const void* const*", "*Pointer", "(*unsafe.Pointer)"},
	"constGLenum*":       {"const unsigned int*", "*uint32", "(*C.uint)"},
	"constGLchar*":       {"const char*", "string", "(*C.char)"},
	"constGLubyte*":      {"const unsigned char*", "*uint8", "(*C.uchar)"},
	"constGLdouble*":     {"const double*", "*float64", "(*C.double)"},
	"constGLshort*":      {"const short *", "*int16", "(*C.short)"},
	"constGLbyte*":       {"const signed char*", "*int8", "(*C.schar)"},
	"constGLushort*":     {"const unsigned short*", "*uint16", "(*C.ushort)"},
	"constGLint64":       {"long long", "int64", "C.longlong"},

	// array of strings.
	"constGLchar**":      {"const char**", "[]string", strArray},
	"constGLchar*const*": {"const char* const*", "[]string", strArray},

	// The following may not be in every spec.
	"struct_cl_context*": {"struct _cl_context*", "*clContext", "(*C.struct__cl_context)"},
	"struct_cl_event*":   {"struct _cl_event*", "*clEvent", "(*C.struct__cl_event)"},
	"GLDEBUGPROCARB":     {"GLDEBUGPROCARB", "DEBUGPROCARB", "C.GLDEBUGPROCARB"},
	"GLDEBUGPROC":        {"GLDEBUGPROC", "DEBUGPROC", "C.GLDEBUGPROC"},
}
