package annotation

import (
	_ "embed"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"

	"github.com/keisuke-m123/goanalyzer/gocode"
)

var (
	//go:embed domain.go
	annotationDomainFile string

	identifierInterfaceType           *types.Interface
	entityInterfaceType               *types.Interface
	valueObjectInterfaceType          *types.Interface
	valueObjectGeneratorInterfaceType *types.Interface

	identifierStructObject           types.Object
	entityStructObject               types.Object
	valueObjectStructObject          types.Object
	valueObjectGeneratorStructObject types.Object
)

func init() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "domain.go", annotationDomainFile, 0)
	if err != nil {
		panic(err)
	}

	var conf types.Config
	pkg, err := conf.Check("github.com/keisuke-m123/gomoduler/annotation", fset, []*ast.File{f}, nil)
	if err != nil {
		panic(err)
	}

	identifierInterfaceType = mustLookupInterfaceType(pkg, "identifier")
	entityInterfaceType = mustLookupInterfaceType(pkg, "entity")
	valueObjectInterfaceType = mustLookupInterfaceType(pkg, "valueObject")
	valueObjectGeneratorInterfaceType = mustLookupInterfaceType(pkg, "valueObjectGenerator")

	identifierStructObject = mustLookupStructObject(pkg, "Identifier")
	entityStructObject = mustLookupStructObject(pkg, "Entity")
	valueObjectStructObject = mustLookupStructObject(pkg, "ValueObject")
	valueObjectGeneratorStructObject = mustLookupStructObject(pkg, "ValueObjectGenerator")
}

func mustLookupInterfaceType(pkg *types.Package, interfaceName string) *types.Interface {
	obj := pkg.Scope().Lookup(interfaceName)
	if obj == nil {
		panic(fmt.Sprintf("interface %s not found", interfaceName))
	}
	typ, ok := obj.Type().Underlying().(*types.Interface)
	if !ok {
		panic(fmt.Sprintf("interface %s not found", interfaceName))
	}
	return typ
}

func mustLookupStructObject(pkg *types.Package, structName string) types.Object {
	obj := pkg.Scope().Lookup(structName)
	if obj == nil {
		panic(fmt.Sprintf("struct %s not found", structName))
	}
	return obj
}

func DomainStruct(t *gocode.Type) bool {
	objs := []types.Object{
		identifierStructObject,
		entityStructObject,
		valueObjectStructObject,
		valueObjectGeneratorStructObject,
	}
	for _, obj := range objs {
		if obj.Pkg().Path() == t.PackageSummary().Path().String() && obj.Name() == t.TypeName().String() {
			return true
		}
	}
	return false
}
