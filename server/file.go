package server

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
)

type FuncKind int

const (
	FUNC FuncKind = iota
	METHOD
)

type funcDescr struct {
	line int
	kind FuncKind
	upon string // What the method acts upon
}

type file struct {
	path     string // Full path to file
	packName string // Qualified name: go/ast
	valid    bool

	funcs map[string][]funcDescr
}

func newFile(packName, path string) *file {
	return &file{
		path:     path,
		packName: packName,
		valid:    false,
	}
}

func _lineNumber(fset *token.FileSet, o ast.Node) int {
	return fset.Position(o.Pos()).Line
}

func (f *file) parse() {
	buf, err := ioutil.ReadFile(f.path)
	if err != nil {
		return
	}
	code := string(buf)
	f.funcs, err = _parse(code)
	if err != nil {
		return
	}
	f.valid = true
}

func _parse(code string) (funcs map[string][]funcDescr, err error) {
	funcs = make(map[string][]funcDescr)

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", code, 0)
	if err != nil {
		log.Println(err)
		return
	}

	// Filter out non-exported
	//ast.FileExports(f)
	for _, n := range f.Decls {
		switch x := n.(type) {
		case *ast.FuncDecl:
			if x.Recv != nil && len(x.Recv.List) > 0 {
				//log.Println("Method on type: %v", x.Recv.List[0].Names[0])
				//log.Printf("%T", x.Recv.List[0])
			} else {
				log.Println("Func")
				list := funcs[x.Name.Name]
				funcs[x.Name.Name] = append(list, funcDescr{
					line: _lineNumber(fset, n),
				})
			}
		case *ast.GenDecl:
			if len(x.Specs) > 0 {
				switch y := x.Specs[0].(type) {
				case *ast.TypeSpec:
					log.Printf("Type:%v", y.Name)
				}
				/*
					switch x.Tok {
					case token.TYPE:
						log.Printf("Type:%v", x)
					default:
						//log.Printf("%T", x)
					}
				*/
			}
		default:
			//log.Printf("%T", n)
		}
	}
	/*
		for _, o := range f.Scope.Objects {
			switch o.Kind {
			case ast.Fun:
				funcs[o.Name] = _lineNumber(fset, o)
			case ast.Typ:
				// o.Decl: TypeSpec
				log.Printf("Type:%T", o.Decl)
			default:
				log.Printf("%v", o.Kind)
			}
		}
	*/
	return
}
