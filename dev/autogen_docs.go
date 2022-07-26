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
	// Make sure we're at the root of the go-pangea sdk repo
	if _, err := os.Stat("./service"); os.IsNotExist(err) {
		fmt.Println("::::::::::::::::::::::::::::::::::::::::::::::::::::::::")
		fmt.Println(": ERROR                                                :")
		fmt.Println(": run this script from the root of the go-pangea repo  :")
		fmt.Println("::::::::::::::::::::::::::::::::::::::::::::::::::::::::")
		log.Fatal(err)
	}

	fs, err := filepath.Glob("./service/**/*.go")
	if err != nil {
		log.Fatal(err)
	}

	// Create the AST
	fset := token.NewFileSet()
	files := []*ast.File{}

	for _, file := range fs {
		parsedFile := mustParse(fset, file)
		files = append(files, parsedFile)
	}

	// Compute package documentation with examples.
	p, err := doc.NewFromFiles(fset, files, "")
	if err != nil {
		log.Fatal(err)
	}

	resultJson, err := json.Marshal(Package{
		Doc:        p.Doc,
		Name:       p.Name,
		ImportPath: p.ImportPath,
		Imports:    p.Imports,
		Filenames:  p.Filenames,
		// Notes:      &p.Notes,

		// Consts: &p.Consts,
		// Types:  &p.Types,
		// Vars:   &p.Vars,
		// Funcs:  &p.Funcs,
	})

	fmt.Println(string(resultJson))

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

type Package struct {
	Doc        string             `json:"doc"`
	Name       string             `json:"name"`
	ImportPath string             `json:"importPath"`
	Imports    []string           `json:"imports"`
	Filenames  []string           `json:"filenames"`
	Notes      map[string][]*Note `json:"notes"`

	// declarations
	Consts []*Value `json:"consts"`
	Types  []*Type  `json:"types"`
	Vars   []*Value `json:"vars"`
	Funcs  []*Func  `json:"funcs"`

	// Examples is a sorted list of examples associated with
	// the package. Examples are extracted from _test.go files
	// provided to NewFromFiles.
	Examples []*Example `json:"examples"`
}

type Func struct {
	Doc  string        `json:"doc"`
	Name string        `json:"name"`
	Decl *ast.FuncDecl `json:"decl"`

	// methods
	// (for functions, these fields have the respective zero value)
	Recv  string `json:"recv"`  // actual   receiver "T" or "*T"
	Orig  string `json:"orig"`  // original receiver "T" or "*T"
	Level int    `json:"level"` // embedding level; 0 means not embedded

	// Examples is a sorted list of examples associated with this
	// function or method. Examples are extracted from _test.go files
	// provided to NewFromFiles.
	Examples []*Example `json:"examples"`
}

type Type struct {
	Doc  string       `json:"doc"`
	Name string       `json:"name"`
	Decl *ast.GenDecl `json:"decl"`

	// associated declarations
	Consts  []*Value `json:"consts"`  // sorted list of constants of (mostly) this type
	Vars    []*Value `json:"vars"`    // sorted list of variables of (mostly) this type
	Funcs   []*Func  `json:"funcs"`   // sorted list of functions returning this type
	Methods []*Func  `json:"methods"` // sorted list of methods (including embedded ones) of this type

	// Examples is a sorted list of examples associated with
	// this type. Examples are extracted from _test.go files
	// provided to NewFromFiles.
	Examples []*Example `json:"examples"`
}

type Example struct {
	Name        string              `json:"name"`   // name of the item being exemplified (including optional suffix)
	Suffix      string              `json:"suffix"` // example suffix, without leading '_' (only populated by NewFromFiles)
	Doc         string              `json:"doc"`    // example function doc string
	Code        ast.Node            `json:"code"`
	Play        *ast.File           `json:"play"` // a whole program version of the example
	Comments    []*ast.CommentGroup `json:"comments"`
	Output      string              `json:output""` // expected output
	Unordered   bool                `json:"unordered"`
	EmptyOutput bool                `json:"emptyoutput"` // expect empty output
	Order       int                 `json:"order"`       // original source code order
}

type Value struct {
	Doc   string       `json:"doc"`
	Names []string     `json:"names"` // var or const names in declaration order
	Decl  *ast.GenDecl `json:"decl"`
	// contains filtered or unexported fields
}

type Note struct {
	Pos, End token.Pos `json:"Pos"`  // position range of the comment containing the marker
	UID      string    `json:"UID"`  // uid found with the marker
	Body     string    `json:"Body"` // note body text
}
