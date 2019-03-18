package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
)

type Symbol struct {
	Name   string
	Line   int
	Object string
}

type Symbols struct {
	Functions []Symbol
	Methods   []Symbol
	Structs   []Symbol
}

func NewSymbols() *Symbols {
	return &Symbols{
		Functions: make([]Symbol, 0),
		Methods:   make([]Symbol, 0),
		Structs:   make([]Symbol, 0),
	}
}

func linenumber(fset *token.FileSet, o ast.Node) int {
	return fset.Position(o.Pos()).Line
}

func (o *Symbols) fun(fs *token.FileSet, f *ast.FuncDecl) {
	if f.Recv != nil {
		// Method
		if len(f.Recv.List) != 1 {
			log.Println("Unexpected")
		}
		field := f.Recv.List[0]
		object := ""
		switch x := field.Type.(type) {
		case *ast.Ident:
			object = x.Name
		default:
			//log.Println("Unexpected")
		}
		o.Methods = append(o.Methods, Symbol{
			Name:   f.Name.Name,
			Line:   linenumber(fs, f),
			Object: object,
		})

		return
	}

	// Function
	o.Functions = append(o.Functions, Symbol{
		Name: f.Name.Name,
		Line: linenumber(fs, f),
	})
}

func (o *Symbols) typ(fs *token.FileSet, s *ast.TypeSpec) {
	switch s.Type.(type) {
	case *ast.StructType:
		// Struct
		o.Structs = append(o.Structs, Symbol{
			Name: s.Name.Name,
			Line: linenumber(fs, s),
		})
	}
}

func (o *Symbols) Parse(code string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", code, 0)
	if err != nil {
		log.Println(err)
		return err
	}

	// Loop through all declarations
	for _, d := range f.Decls {
		switch t := d.(type) {
		case *ast.FuncDecl:
			o.fun(fset, t)
		case *ast.GenDecl:
			for _, spec := range t.Specs {
				switch it := spec.(type) {
				case *ast.TypeSpec:
					o.typ(fset, it)
				default:
					//log.Printf("Unknown spec: %T", spec)
				}
			}
		default:
			log.Printf("Unknown decl: %T", d)
		}
	}

	return nil
}

func Parse(path string) (*Symbols, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	code := string(buf)
	symbols := NewSymbols()
	symbols.Parse(code)

	return symbols, nil
}
