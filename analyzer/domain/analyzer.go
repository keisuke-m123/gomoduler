package domain

import (
	"go/ast"
	"go/types"
	"strings"

	"github.com/keisuke-m123/goanalyzer/gocode"
	"github.com/keisuke-m123/gomoduler/internal/relations"
	"github.com/spf13/afero"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

type (
	analyzer struct {
		nodeStack                 []ast.Node
		analyzer                  *analysis.Analyzer
		foundEntities             *structMap
		foundValueObjects         *structAndDefinedTypeMap
		foundValueObjectGenerator *structMap
	}

	structMap struct {
		m map[string]*gocode.Struct
	}

	structAndDefinedTypeMap struct {
		sm map[string]*gocode.Struct
		dm map[string]*gocode.DefinedType
	}
)

func newStructMap() *structMap {
	return &structMap{
		m: make(map[string]*gocode.Struct),
	}
}

func (e *structMap) add(s *gocode.Struct) {
	e.m[strings.Join([]string{s.PackageSummary().Path().String(), s.Name().String()}, ".")] = s
}

func (e *structMap) get(t types.Type) (*gocode.Struct, bool) {
	s, ok := e.m[t.String()]
	return s, ok
}

func newValueObject() *structAndDefinedTypeMap {
	return &structAndDefinedTypeMap{
		sm: make(map[string]*gocode.Struct),
		dm: make(map[string]*gocode.DefinedType),
	}
}

func (v *structAndDefinedTypeMap) addStruct(s *gocode.Struct) {
	v.sm[strings.Join([]string{s.PackageSummary().Path().String(), s.Name().String()}, ".")] = s
}

func (v *structAndDefinedTypeMap) addDefinedType(d *gocode.DefinedType) {
	v.dm[strings.Join([]string{d.PackageSummary().Path().String(), d.Name().String()}, ".")] = d
}

func (v *structAndDefinedTypeMap) getStruct(t types.Type) (*gocode.Struct, bool) {
	s, ok := v.sm[t.String()]
	return s, ok
}

func (v *structAndDefinedTypeMap) getValueObject(t types.Type) (*gocode.DefinedType, bool) {
	d, ok := v.dm[t.String()]
	return d, ok
}

func NewEntityAnalyzer(domainPath string) *analysis.Analyzer {
	r, err := gocode.LoadRelations(&gocode.LoadOptions{
		FileSystem:  afero.NewOsFs(),
		Directories: []string{domainPath},
		Recursive:   true,
	})
	if err != nil {
		panic(err)
	}

	ea := &analyzer{
		analyzer: &analysis.Analyzer{
			Name:     "entity",
			Doc:      "entity analyzer",
			Requires: []*analysis.Analyzer{inspect.Analyzer},
		},
		foundEntities:             newStructMap(),
		foundValueObjects:         newValueObject(),
		foundValueObjectGenerator: newStructMap(),
	}
	ea.analyzer.Run = ea.analyze

	ea.findEntities(r)
	ea.findValueObjects(r)
	ea.findValueObjectGenerator(r)

	return ea.analyzer
}

func (a *analyzer) findEntities(r *gocode.Relations) {
	for _, s := range r.Structs().StructAll() {
		if s.Implements(relations.GoCodeInterfaceEntity()) {
			a.foundEntities.add(s)
		}
	}
}

func (a *analyzer) findValueObjects(r *gocode.Relations) {
	for _, s := range r.Structs().StructAll() {
		if s.Implements(relations.GoCodeInterfaceValueObject()) {
			a.foundValueObjects.addStruct(s)
		}
	}
	for _, d := range r.DefinedTypes().DefinedTypeAll() {
		if d.Implements(relations.GoCodeInterfaceValueObject()) {
			a.foundValueObjects.addDefinedType(d)
		}
	}
}

func (a *analyzer) findValueObjectGenerator(r *gocode.Relations) {
	for _, s := range r.Structs().StructAll() {
		if s.Implements(relations.GoCodeInterfaceValueObjectGenerator()) {
			a.foundValueObjectGenerator.add(s)
		}
	}
}

func (a *analyzer) analyze(pass *analysis.Pass) (interface{}, error) {
	for _, f := range pass.Files {
		ast.Inspect(f, func(node ast.Node) bool {
			a.checkEntity(pass, node)
			a.checkValueObject(pass, node)
			a.updateNodeStack(node)
			return true
		})
	}
	return nil, nil
}

func (a *analyzer) updateNodeStack(node ast.Node) {
	if node == nil {
		if len(a.nodeStack) == 0 {
			return
		}
		l := len(a.nodeStack)
		a.nodeStack = a.nodeStack[:l-1]
	} else {
		a.nodeStack = append(a.nodeStack, node)
	}
}

func (a *analyzer) checkEntity(pass *analysis.Pass, node ast.Node) {
	switch n := node.(type) {
	case *ast.CompositeLit:
		t := pass.TypesInfo.TypeOf(n)
		if t == nil {
			break
		}
		s, ok := a.foundEntities.get(t)
		if !ok {
			break
		}
		if s.PackageSummary().Path().String() != pass.Pkg.Path() {
			pass.Reportf(n.Pos(), "Entityを実装した構造体はEntityが存在するパッケージ以外からcomposite literalで生成することはできません。")
		}
	}
}

func (a *analyzer) checkValueObject(pass *analysis.Pass, node ast.Node) {
	switch n := node.(type) {
	case *ast.CompositeLit:
		t := pass.TypesInfo.TypeOf(n)
		if t == nil {
			break
		}
		s, ok := a.foundValueObjects.getStruct(t)
		if ok {
			if s.PackageSummary().Path().String() != pass.Pkg.Path() {
				pass.Reportf(n.Pos(), "ValueObjectを実装した構造体はValueObjectが存在するパッケージ以外からcomposite literalで生成することはできません。")
			}
			break
		}
	case *ast.CallExpr:
		t := pass.TypesInfo.TypeOf(n.Fun)
		if t == nil {
			break
		}
		// defined typeの呼び出しのみを補足するため、prefixをみて通常の関数/メソッド呼び出しと区別
		if strings.HasPrefix(t.String(), "func(") {
			break
		}
		d, ok := a.foundValueObjects.getValueObject(t)
		if ok {
			if d.PackageSummary().Path().String() != pass.Pkg.Path() {
				pass.Reportf(n.Pos(), "ValueObjectを実装した構造体はValueObjectが存在するパッケージ以外から直接生成することはできません。")
			}
			break
		}
	}
}
