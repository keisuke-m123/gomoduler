package domain

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/keisuke-m123/gomoduler/annotation"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

type (
	analyzer struct {
		analyzer    *analysis.Analyzer
		annotations *annotation.FindAnnotationsResult
	}
)

func NewEntityAnalyzer(domainPath string) *analysis.Analyzer {
	ea := &analyzer{
		analyzer: &analysis.Analyzer{
			Name:     "entity",
			Doc:      "entity analyzer",
			Requires: []*analysis.Analyzer{inspect.Analyzer},
		},
		annotations: annotation.FindAnnotations(domainPath),
	}
	ea.analyzer.Run = ea.analyze
	return ea.analyzer
}

func (a *analyzer) analyze(pass *analysis.Pass) (interface{}, error) {
	i := inspector.New(pass.Files)
	nodes := []ast.Node{
		&ast.FuncDecl{},
		&ast.CompositeLit{},
		&ast.CallExpr{},
		&ast.AssignStmt{},
	}
	i.WithStack(nodes, func(n ast.Node, push bool, stack []ast.Node) (proceed bool) {
		if push {
			a.checkEntityInitialization(pass, n, stack)
			a.checkValueObjectInitialization(pass, n, stack)
			a.checkValueObjectImmutable(pass, n, stack)
		}
		return true
	})
	return nil, nil
}

func (a *analyzer) checkEntityInitialization(pass *analysis.Pass, node ast.Node, nodeStack []ast.Node) {
	switch n := node.(type) {
	case *ast.CompositeLit:
		t := pass.TypesInfo.TypeOf(n)
		if t == nil {
			break
		}
		s, ok := a.annotations.GetEntity(t)
		if !ok {
			break
		}
		if s.PackageSummary().Path().String() != pass.Pkg.Path() {
			pass.Reportf(n.Pos(), "Entityを実装した構造体はEntityが存在するパッケージ以外からcomposite literalで生成することはできません。")
		}
	}
}

func (a *analyzer) checkValueObjectInitialization(pass *analysis.Pass, node ast.Node, nodeStack []ast.Node) {
	switch n := node.(type) {
	case *ast.CompositeLit:
		t := pass.TypesInfo.TypeOf(n)
		if t == nil {
			break
		}
		s, ok := a.annotations.GetValueObjectStruct(t)
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
		d, ok := a.annotations.GetValueObjectDefinedType(t)
		if ok {
			for i := range nodeStack {
				switch n := nodeStack[i].(type) {
				case *ast.FuncDecl:
					if n.Recv != nil {
						for _, f := range n.Recv.List {
							switch t := f.Type.(type) {
							case *ast.Ident:
								df := pass.TypesInfo.Defs[t]
								if _, ok := a.annotations.GetValueObjectGenerator(df.Type()); ok {
									return
								}
							case *ast.StarExpr:
								df := pass.TypesInfo.TypeOf(t.X)
								if df == nil {
									continue
								}
								if _, ok := a.annotations.GetValueObjectGenerator(df); ok {
									return
								}
							}
						}
					}
				}
			}
			if d.PackageSummary().Path().String() != pass.Pkg.Path() {
				pass.Reportf(n.Pos(), "ValueObjectを実装した構造体はValueObjectが存在するパッケージ以外から直接生成することはできません。")
			}
			break
		}
	}
}

func (a *analyzer) checkValueObjectImmutable(pass *analysis.Pass, node ast.Node, nodeStack []ast.Node) {
	switch n := node.(type) {
	case *ast.FuncDecl:
		fmt.Println(n.Name.Name)
	case *ast.AssignStmt:
		fmt.Println(n.Lhs[0])
	}
}
