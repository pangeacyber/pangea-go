package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
)

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

	// type pkgMap struct {
	// 	fileSet *token.FileSet
	// 	docPkg  *doc.Package
	// }
	// packages := make(map[string]pkgMap)
	packages := []*doc.Package{}

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
			// fmt.Println("LOOKING AT: ", file)
			astFile := mustParse(fset, file)
			files = append(files, astFile)

			// fmt.Println(":::::::: PACKAGE NAME:::::", astFile.Name.Name)
			// for _, comment := range astFile.Comments {
			// 	fmt.Println("COMMENT", comment.Text())
			// }
		}

		// Compute package documentation with examples.
		p, err := doc.NewFromFiles(fset, files, fmt.Sprintf("github.com/pangeacyber/go-pangea/%s", dir), doc.AllMethods)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("DOC package", p.Name, p.Filenames)
		packages = append(packages, p)
		// packages[p.Name] = pkgMap{
		// 	fileSet: fset,
		// 	docPkg:  p,
		// }
	}

	for _, p := range packages {
		fmt.Println("PACKAGE:::::", p.Name)

		types := gatherTypes(p.Types)
		funcs := gatherFuncs(p.Funcs)

		resultJson, err := json.MarshalIndent(Document{
			Package: p.Name,
			Types:   types,
			Funcs:   funcs,
		}, "", "\t")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(resultJson))

		// for _, types := range p.Types {
		// 	fmt.Println("  type: ", types.Name)
		// 	fmt.Println("    > decl: ", &types.Decl)

		// 	// ast.Inspect(types.Decl, func(node ast.Node) bool {
		// 	// 	switch n := node.(type) {
		// 	// 	case nil:
		// 	// 		return true
		// 	// 	case *ast.FuncDecl:
		// 	// 		name := n.Name.Name
		// 	// 		recv :=

		// 	// 	}
		// 	// 	// find funcs
		// 	// 	fn, ok := n.(*ast.FuncType)
		// 	// 	if ok {
		// 	// 		fmt.Printf("fn: %v\n", fn)
		// 	// 		if fn.Params != nil {
		// 	// 			for _, field := range fn.Params.List {
		// 	// 				fmt.Println("HEREEEE", field.Names)
		// 	// 			}
		// 	// 		}
		// 	// 		// fmt.Printf("func found on line %d:\n\t", fset.Position())
		// 	// 	}

		// 	// 	return true
		// 	// })

		// 	// for _, c := range types.Consts {
		// 	// 	fmt.Println("    > consts", c)
		// 	// }
		// 	// for _, v := range types.Vars {
		// 	// 	fmt.Println("    > vars", v.Names[0])
		// 	// }
		// 	// // fmt.Println("    > decl.specs", types.Decl.Specs)

		// 	// for _, methods := range types.Methods {
		// 	// 	fmt.Println("    methods: ")
		// 	// 	fmt.Println("      > name:", methods.Name)
		// 	// 	fmt.Println("      > doc:", methods.Doc)
		// 	// 	fmt.Println("      > examples:", methods.Examples)

		// 	// 	for _, params := range methods.Decl.Type.Params.List {
		// 	// 		fmt.Println("      > params:", params.Names[0], params.Type.Pos(), params.Type.End())
		// 	// 	}
		// 	// }

		// 	// for _, funcs := range types.Funcs {
		// 	// 	fmt.Println("    funcs: ", funcs.Name, funcs.Doc)
		// 	// }
		// }

		// for _, funcs := range p.Funcs {
		// 	fmt.Println("    funcs: ", funcs.Name, funcs.Doc)
		// }
	}

	// b, err := pangeautil.CanonicalizeJSONMarshall(
	// 	Package{
	// 		Doc: p.Doc,
	// 		Name: p.Name,
	// 		Types: p.Types,
	// 	},
	// )

	// fmt.Println(("result json:::::::"))
	// fmt.Println(string(resultJson))

	// fmt.Printf("package %s - %s", p.Name, p.Doc)
	// fmt.Printf("func %s - %s", p.Funcs[0].Name, p.Funcs[0].Doc)
	// fmt.Println("TYPES::::::::::")
	// // fmt.Println(len(p.Types))
	// for i := 0; i < len(p.Types); i++ {
	// 	fmt.Println(p.Types[i].Name, p.Types[i].Methods)

	// 	for k := range p.Types[i].Methods {
	// 		method := p.Types[i].Methods[k]

	// 		fmt.Println(method.Name, method.Doc)
	// 	}
	// }
	// fmt.Println("Funcs::::::::::")
	// fmt.Println(len(p.Funcs))
	// for i := range p.Funcs {
	// 	f := p.Funcs[i]
	// 	fmt.Println(f.Name, f.Doc)
	// }

	// fmt.Printf(" ⤷ example with suffix %q - %s", p.Funcs[0].Examples[0].Suffix, p.Funcs[0].Examples[0].Doc)

}

func mustParse(fset *token.FileSet, filename string) *ast.File {
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

type Document struct {
	Package string `json:"package"`
	Doc     string `json:"doc"`

	Types  []*DocumentedType  `json:"types"`
	Consts []*DocumentedConst `json:"consts"`
	Funcs  []*DocumentedFunc  `json:"funcs"`
}

type DocumentedType struct {
	Name string `json:"name"`
	Doc  string `json:"doc"`

	Funcs   []*DocumentedFunc `json:"funcs"`
	Methods []*DocumentedFunc `json:"methods"`

	Consts []*DocumentedValue `json:"consts"`
	Vars   []*DocumentedValue `json:"vars"`
}

type DocumentedFunc struct {
	Name          string       `json:"name"`
	Doc           string       `json:"doc"`
	Recv          string       `json:"recv"`
	FormattedRecv string       `json:"formattedRecv"`
	Orig          string       `json:"orig"`
	Level         int          `json:"level"`
	Type          *AstFuncType `json:"type"`
}

type DocumentedValue struct {
	Doc   string   `json:"doc"`
	Names []string `json:"names"`
}

type DocumentedConst struct {
	Name string `json:"name"`
	Doc  string `json:"doc"`
}

func gatherTypes(docTypes []*doc.Type) []*DocumentedType {
	types := []*DocumentedType{}
	for _, t := range docTypes {
		types = append(types, &DocumentedType{
			Name:    t.Name,
			Doc:     t.Doc,
			Funcs:   gatherFuncs(t.Funcs),
			Methods: gatherFuncs(t.Methods),
			Vars:    gatherValues(t.Vars),
			Consts:  gatherValues(t.Consts),
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

func gatherFuncs(docFuncs []*doc.Func) []*DocumentedFunc {
	funcs := []*DocumentedFunc{}

	for _, f := range docFuncs {

		var formattedRecv = ""
		if len(f.Recv) > 0 {
			formattedRecv = "(" + f.Recv + ")"
		}

		funcs = append(funcs, &DocumentedFunc{
			Name:          f.Name,
			Doc:           f.Doc,
			Recv:          f.Recv,
			FormattedRecv: formattedRecv,
			Orig:          f.Orig,
			Level:         f.Level,
			Type:          gatherAstFuncType(f.Decl.Type),
		})
	}

	return funcs
}

func gatherAstFuncType(astType *ast.FuncType) *AstFuncType {
	return &AstFuncType{
		Params:  gatherAstFieldList(astType.Params.List),
		Results: gatherAstFieldList(astType.Results.List),
	}
}

type AstFuncType struct {
	TypeParams []*AstField `json:"typeParams"`
	Params     []*AstField `json:"params"`
	Results    []*AstField `json:"results"`
}

type AstField struct {
	Doc        string   `json:"doc"`
	Comment    string   `json:"comment"`
	Names      []string `json:"names"`
	IsExported bool     `json:"isExported"`
}

func gatherAstFieldList(l []*ast.Field) []*AstField {
	fields := []*AstField{}
	for _, field := range l {
		names := []string{}

		for _, fn := range field.Names {
			names = append(names, fn.Obj.Name)
		}

		fields = append(fields, &AstField{
			Doc:     field.Doc.Text(),
			Comment: field.Comment.Text(),
			Names:   names,
		})
	}

	return fields
}

// Custom marshal func from
// https://blog.logrocket.com/using-json-go-guide/
// func (p *doc.Package) MarshalJSON() ([]byte, error) {
// 	type PackageAlias doc.Package
// 	return json.Marshal(&struct{
// 		*PackageAlias

// 	})
// }

// type Package struct {
// 	Doc        string            `json:"doc"`
// 	Name       string            `json:"name"`
// 	ImportPath string            `json:"importPath"`
// 	Imports    []string          `json:"imports"`
// 	Filenames  []string          `json:"filenames"`
// 	Notes      map[string][]Note `json:"notes"`

// 	// declarations
// 	Consts []Value `json:"consts"`
// 	Types  []Type  `json:"types"`
// 	Vars   []Value `json:"vars"`
// 	Funcs  []Func  `json:"funcs"`

// 	// Examples is a sorted list of examples associated with
// 	// the package. Examples are extracted from _test.go files
// 	// provided to NewFromFiles.
// 	Examples []Example `json:"examples"`
// }

// type Func struct {
// 	Doc  string        `json:"doc"`
// 	Name string        `json:"name"`
// 	Decl *ast.FuncDecl `json:"decl"`

// 	// methods
// 	// (for functions, these fields have the respective zero value)
// 	Recv  string `json:"recv"`  // actual   receiver "T" or "*T"
// 	Orig  string `json:"orig"`  // original receiver "T" or "*T"
// 	Level int    `json:"level"` // embedding level; 0 means not embedded

// 	// Examples is a sorted list of examples associated with this
// 	// function or method. Examples are extracted from _test.go files
// 	// provided to NewFromFiles.
// 	Examples []*Example `json:"examples"`
// }

// type Type struct {
// 	Doc  string       `json:"doc"`
// 	Name string       `json:"name"`
// 	Decl *ast.GenDecl `json:"decl"`

// 	// associated declarations
// 	Consts  []*Value `json:"consts"`  // sorted list of constants of (mostly) this type
// 	Vars    []*Value `json:"vars"`    // sorted list of variables of (mostly) this type
// 	Funcs   []*Func  `json:"funcs"`   // sorted list of functions returning this type
// 	Methods []*Func  `json:"methods"` // sorted list of methods (including embedded ones) of this type

// 	// Examples is a sorted list of examples associated with
// 	// this type. Examples are extracted from _test.go files
// 	// provided to NewFromFiles.
// 	Examples []*Example `json:"examples"`
// }

// type Example struct {
// 	Name        string              `json:"name"`   // name of the item being exemplified (including optional suffix)
// 	Suffix      string              `json:"suffix"` // example suffix, without leading '_' (only populated by NewFromFiles)
// 	Doc         string              `json:"doc"`    // example function doc string
// 	Code        ast.Node            `json:"code"`
// 	Play        *ast.File           `json:"play"` // a whole program version of the example
// 	Comments    []*ast.CommentGroup `json:"comments"`
// 	Output      string              `json:output""` // expected output
// 	Unordered   bool                `json:"unordered"`
// 	EmptyOutput bool                `json:"emptyoutput"` // expect empty output
// 	Order       int                 `json:"order"`       // original source code order
// }

// type Value struct {
// 	Doc   string       `json:"doc"`
// 	Names []string     `json:"names"` // var or const names in declaration order
// 	Decl  *ast.GenDecl `json:"decl"`
// 	// contains filtered or unexported fields
// }

// type Note struct {
// 	Pos, End token.Pos `json:"Pos"`  // position range of the comment containing the marker
// 	UID      string    `json:"UID"`  // uid found with the marker
// 	Body     string    `json:"Body"` // note body text
// }
