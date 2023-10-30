package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"path/filepath"
)

type Document struct {
	Package string `json:"package"`
	Doc     string `json:"doc"`

	Types  []*DocumentedType  `json:"types"`
	Consts []*DocumentedConst `json:"consts"`
	Funcs  []*DocumentedFunc  `json:"funcs"`
}

type DocumentedType struct {
	Name        string `json:"name"`
	Doc         string `json:"doc"`
	Declaration string `json:"declaration"`

	Funcs   []*DocumentedFunc `json:"funcs"`
	Methods []*DocumentedFunc `json:"methods"`

	Consts []*DocumentedValue `json:"consts"`
	Vars   []*DocumentedValue `json:"vars"`

	AstFields []*AstField `json:"astFields"`
}

type DocumentedFunc struct {
	Name        string               `json:"name"`
	Doc         string               `json:"doc"`
	Declaration string               `json:"declaration"`
	Level       int                  `json:"level"`
	Type        *AstFuncType         `json:"type"`
	Examples    []*DocumentedExample `json:"examples"`
}

type DocumentedExample struct {
	Name        string `json:"name"`
	Suffix      string `json:"suffix"`
	Doc         string `json:"doc"`
	Output      string `json:"output"`
	Unordered   bool   `json:"unordered"`
	EmptyOutput bool   `json:"emptyOutput"`
	Order       int    `json:"order"`
}

type DocumentedValue struct {
	Doc     string   `json:"doc"`
	Names   []string `json:"names"`
	Comment string   `json:"comment"`
}

type DocumentedConst struct {
	Name string `json:"name"`
	Doc  string `json:"doc"`
}

type AstFuncType struct {
	TypeParams []*AstField `json:"typeParams"`
	Params     []*AstField `json:"params"`
	Results    []*AstField `json:"results"`
}

type AstField struct {
	Doc     string   `json:"doc"`
	Comment string   `json:"comment"`
	Names   []string `json:"names"`
	Type    string   `json:"type"`
	Tag     string   `json:"tag"`
}

func main() {
	DIRECTORIES := []string{
		"./service",
		"./pangea",
		"./pangea/hash",
	}

	// Make sure we're at the root of the go-pangea sdk repo
	// and we can find the directories to pull docs from
	for _, dir := range DIRECTORIES {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			fmt.Println("::::::::::::::::::::::::::::::::::::::::::::::::::::::::")
			fmt.Println(": ERROR")
			fmt.Printf(": could not find directory: %s", dir)
			fmt.Println(":")
			fmt.Println(": make sure to run this script from")
			fmt.Println(": the root of the go-pangea repo")
			fmt.Println("::::::::::::::::::::::::::::::::::::::::::::::::::::::::")
			log.Fatal(err)
		}
	}

	serviceDirs, err := filepath.Glob("./service/*")
	if err != nil {
		log.Fatal(err)
	}

	// Walk the dir for Pangea package: "./pangea"
	// pangeaFiles := []string{}
	// err = filepath.WalkDir("./pangea", func(path string, d fs.DirEntry, err error) error {
	// 	if err != nil {
	// 		return err
	// 	}

	// 	if !d.IsDir() {
	// 		pangeaFiles = append(pangeaFiles, path)
	// 	}

	// 	return nil
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	type pkgMap struct {
		fileSet *token.FileSet
		docPkg  *doc.Package
	}
	packages := make(map[string]pkgMap)

	for _, dir := range serviceDirs {
		// Create the AST
		fset := token.NewFileSet()
		files := []*ast.File{}

		serviceFiles, err := filepath.Glob(fmt.Sprintf("%s/*.go", dir))
		if err != nil {
			log.Fatal(err)
		}

		filepaths := append([]string{}, serviceFiles...)
		// filepaths = append(filepaths, pangeaFiles...) // Maybe also add the go files in pangea/ ?

		for _, file := range filepaths {
			astFile := mustParse(fset, file)
			files = append(files, astFile)
		}

		// Compute package documentation with examples.
		p, err := doc.NewFromFiles(fset, files, fmt.Sprintf("github.com/pangeacyber/pangea-go/pangea-sdk/v3/%s", dir), doc.AllDecls)
		if err != nil {
			log.Fatal(err)
		}

		packages[p.Name] = pkgMap{
			fileSet: fset,
			docPkg:  p,
		}
	}

	documents := []*Document{}
	for _, p := range packages {
		types := gatherTypes(p.docPkg.Types, p.fileSet)
		funcs := gatherFuncs(p.docPkg.Funcs, p.fileSet)

		documents = append(documents, &Document{
			Package: p.docPkg.Name,
			Types:   types,
			Funcs:   funcs,
		})
	}

	resultJson, err := json.MarshalIndent(documents, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(resultJson))
}

func mustParse(fset *token.FileSet, filename string) *ast.File {
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func gatherTypes(docTypes []*doc.Type, fs *token.FileSet) []*DocumentedType {
	types := []*DocumentedType{}
	for _, t := range docTypes {

		vals := []*AstField{}
		ast.Inspect(t.Decl, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.Field:
				field := gatherAstField(x, fs)
				if field != nil {
					vals = append(vals, field)
				}
			}
			return true
		})

		types = append(types, &DocumentedType{
			Name:        t.Name,
			Doc:         t.Doc,
			Declaration: prettify(t.Decl, fs),
			Funcs:       gatherFuncs(t.Funcs, fs),
			Methods:     gatherFuncs(t.Methods, fs),
			Vars:        gatherValues(t.Vars),
			Consts:      gatherValues(t.Consts),
			AstFields:   vals,
		})
	}

	return types
}

func gatherValues(vals []*doc.Value) []*DocumentedValue {
	values := []*DocumentedValue{}

	for _, v := range vals {
		values = append(values, &DocumentedValue{
			Doc:   v.Doc,
			Names: v.Names,
		})
	}

	return values
}

func gatherFuncs(docFuncs []*doc.Func, fs *token.FileSet) []*DocumentedFunc {
	funcs := []*DocumentedFunc{}

	for _, f := range docFuncs {
		funcs = append(funcs, &DocumentedFunc{
			Name:        f.Name,
			Doc:         f.Doc,
			Declaration: prettify(f.Decl, fs),
			Level:       f.Level,
			Type:        gatherAstFuncType(f.Decl.Type, fs),
			Examples:    gatherExamples(f.Examples, fs),
		})
	}

	return funcs
}

// Gather examples from *_test.go files
func gatherExamples(examples []*doc.Example, fs *token.FileSet) []*DocumentedExample {
	ex := []*DocumentedExample{}

	for _, e := range examples {
		ex = append(ex, &DocumentedExample{
			Name:        e.Name,
			Suffix:      e.Suffix,
			Doc:         e.Doc,
			Output:      e.Output,
			EmptyOutput: e.EmptyOutput,
			Unordered:   e.Unordered,
			Order:       e.Order,
		})
	}

	return ex
}

func gatherAstFuncType(astType *ast.FuncType, fs *token.FileSet) *AstFuncType {
	results := []*AstField{}

	if astType.Results != nil {
		results = gatherAstFieldList(astType.Results.List, fs)
	}

	return &AstFuncType{
		Params:  gatherAstFieldList(astType.Params.List, fs),
		Results: results,
	}
}

func gatherAstField(f *ast.Field, fs *token.FileSet) *AstField {
	names := []string{}
	for _, fns := range f.Names {
		names = append(names, fns.Name)
	}

	if len(names) == 0 && f.Doc.Text() == "" {
		return nil
	}

	var tag string
	if f.Tag != nil {
		tag = f.Tag.Value
	}

	return &AstField{
		Doc:     f.Doc.Text(),
		Comment: f.Comment.Text(),
		Names:   names,
		Type:    prettify(f.Type, fs),
		Tag:     tag,
	}
}

func gatherAstFieldList(fl []*ast.Field, fs *token.FileSet) []*AstField {
	fields := []*AstField{}
	for _, field := range fl {
		fields = append(fields, gatherAstField(field, fs))
	}

	return fields
}

// Pretty print the ast Node as a string
func prettify(node interface{}, fs *token.FileSet) string {
	var stringBuffer bytes.Buffer
	printer.Fprint(&stringBuffer, fs, node)

	return string(stringBuffer.Bytes())
}
