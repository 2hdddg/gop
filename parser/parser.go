package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

type Base struct {
	Name string
	Line int
}

type Function struct {
	base Base
}

type Method struct {
	base   Base
	Object string
}

type Struct struct {
	base Base
}

type Output struct {
	Functions []Function
	Methods   []Method
	Structs   []Struct
}

func NewOutput() *Output {
	return &Output{
		Functions: make([]Function, 0),
		Methods:   make([]Method, 0),
		Structs:   make([]Struct, 0),
	}
}

func linenumber(fset *token.FileSet, o ast.Node) int {
	return fset.Position(o.Pos()).Line
}

func (o *Output) fun(fs *token.FileSet, f *ast.FuncDecl) {
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
			log.Println("Unexpected")
		}
		o.Methods = append(o.Methods, Method{
			base: Base{
				Name: f.Name.Name,
				Line: linenumber(fs, f),
			},
			Object: object,
		})

		return
	}

	// Function
	o.Functions = append(o.Functions, Function{
		base: Base{
			Name: f.Name.Name,
			Line: linenumber(fs, f),
		},
	})
}

func (o *Output) typ(fs *token.FileSet, s *ast.TypeSpec) {
	switch s.Type.(type) {
	case *ast.StructType:
		// Struct
		o.Structs = append(o.Structs, Struct{
			base: Base{
				Name: s.Name.Name,
				Line: linenumber(fs, s),
			},
		})
	}
}

func (o *Output) Parse(code string) error {
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
					log.Printf("Unknown spec: %T", spec)
				}
			}
		default:
			log.Printf("Unknown decl: %T", d)
		}
	}

	return nil
}
