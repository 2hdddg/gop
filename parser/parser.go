package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"strconv"
)

type Type int

const (
	Undefined Type = iota
	Function
	Method
	Struct
	Interface
	Field
)

type Symbol struct {
	Name        string
	Type        Type
	Line        int
	ContextName string // For Method or Field
	ContextType Type   // For Field
}

type Symbols struct {
	List []Symbol
}

func NewSymbols() *Symbols {
	return &Symbols{
		List: make([]Symbol, 0, 10),
	}
}

func typeToString(t Type) string {
	switch t {
	case Function:
		return "func"
	case Method:
		return "func"
	case Struct:
		return "struct"
	case Interface:
		return "interface"
	case Field:
		return "field"
	}
	return fmt.Sprintf("Unknown: %v", t)
}

func (o *Symbol) ToString() string {
	s := typeToString(o.Type)

	fmt.Println(*o)

	if o.ContextName != "" {
		switch o.Type {
		case Method:
			return s + fmt.Sprintf(" on %v", o.ContextName)
		case Field:
			return s + fmt.Sprintf(" of %v(%v)", o.ContextName, typeToString(o.ContextType))
		}
	}
	return s
}

func linenumber(fset *token.FileSet, p token.Pos) int {
	return fset.Position(p).Line
}

func (o *Symbols) add(sym Symbol) {
	o.List = append(o.List, sym)
}

func (o *Symbols) fun(fs *token.FileSet, f *ast.FuncDecl) {
	if f.Recv != nil {
		// Method
		if len(f.Recv.List) != 1 {
			log.Println("Unexpected")
		}
		field := f.Recv.List[0]
		context := ""
		switch x := field.Type.(type) {
		case *ast.Ident:
			context = x.Name
		case *ast.StarExpr:
			switch y := x.X.(type) {
			case *ast.Ident:
				context = "*" + y.Name
			default:
				//log.Printf("Unexpected *: %T", y)
			}
		default:
			log.Printf("Unexpected: %T", x)
		}
		o.add(Symbol{
			Name:        f.Name.Name,
			Line:        linenumber(fs, f.Pos()),
			Type:        Method,
			ContextName: context,
		})
		return
	}

	// Function
	o.add(Symbol{
		Name: f.Name.Name,
		Line: linenumber(fs, f.Pos()),
		Type: Function,
	})
}

func (o *Symbols) typ(fs *token.FileSet, s *ast.TypeSpec) {
	addFields := func(fields *ast.FieldList, c string, t, ct Type) {
		for _, f := range fields.List {
			if len(f.Names) > 0 {
				o.add(Symbol{
					Name:        f.Names[0].Name,
					Line:        linenumber(fs, f.Pos()),
					Type:        t,
					ContextName: c,
					ContextType: ct,
				})
			}
		}
	}

	switch t := s.Type.(type) {
	case *ast.StructType:
		o.add(Symbol{
			Name: s.Name.Name,
			Line: linenumber(fs, s.Pos()),
			Type: Struct,
		})
		// Add members of struct
		addFields(t.Fields, s.Name.Name, Field, Struct)
	case *ast.InterfaceType:
		o.add(Symbol{
			Name: s.Name.Name,
			Line: linenumber(fs, s.Pos()),
			Type: Interface,
		})
		// Add methods of interface as fields
		addFields(t.Methods, s.Name.Name, Method, Interface)
	default:
		//log.Printf("Unknown type: %T", s.Type)
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
			//log.Printf("Unknown decl: %T", d)
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
