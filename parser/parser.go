package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"strconv"
)

type Symbol struct {
	Name       string
	Line       int
	Object     string // Name of struct for methods
	ObjectKind string // Struct, Interface
}

type Symbols struct {
	Functions  []Symbol // Global, on struct
	Methods    []Symbol // Remove?
	Structs    []Symbol
	Interfaces []Symbol
	//Vars	global,
	// Member of struct, interface
}

func NewSymbols() *Symbols {
	return &Symbols{
		Functions:  make([]Symbol, 0),
		Methods:    make([]Symbol, 0),
		Structs:    make([]Symbol, 0),
		Interfaces: make([]Symbol, 0),
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
		case *ast.StarExpr:
			switch y := x.X.(type) {
			case *ast.Ident:
				object = y.Name
			default:
				log.Printf("Unexpected *: %T", y)
			}
		default:
			log.Printf("Unexpected: %T", x)
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
		o.Structs = append(o.Structs, Symbol{
			Name: s.Name.Name,
			Line: linenumber(fs, s),
		})
	case *ast.InterfaceType:
		o.Interfaces = append(o.Interfaces, Symbol{
			Name: s.Name.Name,
			Line: linenumber(fs, s),
		})
	default:
		log.Printf("Unknown type: %T", s.Type)
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
					log.Printf("Unknown spec: %T", spec)
				}
			}
		default:
			log.Printf("Unknown decl: %T", d)
		}
	}

	return nil
}

func parseImports(code string) ([]string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", code, parser.ImportsOnly)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	imps := make([]string, len(f.Imports))
	for x, i := range f.Imports {
		name, _ := strconv.Unquote(i.Path.Value)
		imps[x] = name
	}
	return imps, nil
}

func getCode(path string) (string, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func Parse(path string) (*Symbols, error) {
	code, err := getCode(path)
	if err != nil {
		return nil, err
	}
	symbols := NewSymbols()
	symbols.Parse(code)

	return symbols, nil
}

func ParseImports(path string) ([]string, error) {
	code, err := getCode(path)
	if err != nil {
		return nil, err
	}

	return parseImports(code)
}
