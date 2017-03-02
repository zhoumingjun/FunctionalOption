package fo

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/format"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}

// Generator holds the state of the analysis. Primarily used to buffer
// the output for format.Source.
type Generator struct {
	buf bytes.Buffer // Accumulated output.
	pkg *Package     // Package we are scanning.
}

func (g *Generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

type Option struct {
	FieldName  string
	Type       string
	OptionName string
}

// File holds a single parsed file and associated data.
type File struct {
	pkg  *Package  // Package to which this file belongs.
	file *ast.File // Parsed AST.
	// These fields are reset for each type being generated.
	typeName string   // Name of the constant type.
	options  []Option // Accumulator for constant values of that type.
}

type Package struct {
	dir      string
	name     string
	defs     map[*ast.Ident]types.Object
	files    []*File
	typesPkg *types.Package
}

// parsePackageDir parses the package residing in the directory.
func (g *Generator) ParsePackageDir(directory string) {
	pkg, err := build.Default.ImportDir(directory, 0)
	if err != nil {
		log.Fatalf("cannot process directory %s: %s", directory, err)
	}
	var names []string
	names = append(names, pkg.GoFiles...)
	names = append(names, pkg.CgoFiles...)
	// TODO: Need to think about constants in test files. Maybe write type_string_test.go
	// in a separate pass? For later.
	// names = append(names, pkg.TestGoFiles...) // These are also in the "foo" package.
	names = append(names, pkg.SFiles...)
	names = prefixDirectory(directory, names)
	g.parsePackage(directory, names, nil)
}

// parsePackageFiles parses the package occupying the named files.
func (g *Generator) parsePackageFiles(names []string) {
	g.parsePackage(".", names, nil)
}

// prefixDirectory places the directory name on the beginning of each name in the list.
func prefixDirectory(directory string, names []string) []string {
	if directory == "." {
		return names
	}
	ret := make([]string, len(names))
	for i, name := range names {
		ret[i] = filepath.Join(directory, name)
	}
	return ret
}

// parsePackage analyzes the single package constructed from the named files.
// If text is non-nil, it is a string to be used instead of the content of the file,
// to be used for testing. parsePackage exits if there is an error.
func (g *Generator) parsePackage(directory string, names []string, text interface{}) {
	var files []*File
	var astFiles []*ast.File
	g.pkg = new(Package)
	fs := token.NewFileSet()
	for _, name := range names {
		if !strings.HasSuffix(name, ".go") {
			continue
		}
		parsedFile, err := parser.ParseFile(fs, name, text, 0)
		if err != nil {
			log.Fatalf("parsing package: %s: %s", name, err)
		}
		astFiles = append(astFiles, parsedFile)
		files = append(files, &File{
			file: parsedFile,
			pkg:  g.pkg,
		})
	}
	if len(astFiles) == 0 {
		log.Fatalf("%s: no buildable Go files", directory)
	}
	g.pkg.name = astFiles[0].Name.Name
	g.pkg.files = files
	g.pkg.dir = directory
	// Type check the package.
	g.pkg.check(fs, astFiles)
}

// check type-checks the package. The package must be OK to proceed.
func (pkg *Package) check(fs *token.FileSet, astFiles []*ast.File) {
	pkg.defs = make(map[*ast.Ident]types.Object)
	config := types.Config{Importer: importer.Default(), FakeImportC: true}
	info := &types.Info{
		Defs: pkg.defs,
	}
	typesPkg, err := config.Check(pkg.dir, fs, astFiles, info)
	if err != nil {
		log.Fatalf("checking package: %s", err)
	}
	pkg.typesPkg = typesPkg
}

// generate produces the String method for the named type.
func (g *Generator) Generate(typeName string) {
	options := make([]Option, 0, 100)
	for _, file := range g.pkg.files {
		// Set the state for this run of the walker.
		file.typeName = typeName
		file.options = nil

		if file.file != nil {
			ast.Inspect(file.file, file.genDecl)
			options = append(options, file.options...)
		}
	}

	if len(options) == 0 {
		log.Fatalf("no values defined for type %s", typeName)
		return
	}

	tmpl, err := template.New("model").Parse(tmplOptions)
	if err != nil {
		fmt.Println(err)
		return
	}

	var b bytes.Buffer

	err = tmpl.Execute(&b, map[string]interface{}{
		"PackageName": g.pkg.name,
		"TypeName":    typeName,
		"Options":     options,
	})

	out, err := format.Source(b.Bytes())
	if err != nil {
		fmt.Println(" format", err)
		return
	}

	output := strings.ToLower(typeName + "_option.go")
	outputPath := filepath.Join(g.pkg.dir, output)
	if err := ioutil.WriteFile(outputPath, out, 0644); err != nil {
		log.Fatalf("writing output: %s", err)
	}
}

// genDecl processes one declaration clause.
func (f *File) genDecl(node ast.Node) bool {

	spec, ok := node.(*ast.TypeSpec)
	if !ok || spec.Name.Name != f.typeName {
		return true
	}

	t, ok := spec.Type.(*ast.StructType)
	for _, field := range t.Fields.List {

		if field.Tag != nil {
			strTag := field.Tag.Value
			strTag = strTag[1 : len(strTag)-1]
			tags := strings.Split(strTag, ",")
			for _, tag := range tags {
				if strings.HasPrefix(tag, "option") {
					fieldName := field.Names[0].Name
					fieldType := field.Type.(*ast.Ident).Name
					option := Option{
						FieldName:  fieldName,
						Type:       fieldType,
						OptionName: fieldName,
					}

					eles := strings.Split(tag, ":")
					if len(eles) == 2 {
						option.OptionName = eles[1][1 : len(eles[1])-1]
					}

					f.options = append(f.options, option)
				}
			}

		}

	}
	return false
}
