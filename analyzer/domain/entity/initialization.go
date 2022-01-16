package entity

import (
	"go/ast"

	"github.com/keisuke-m123/gomoduler/analyzer/astutil"
	"github.com/keisuke-m123/gomoduler/annotation"
	"golang.org/x/tools/go/analysis"
)

type (
	InitializationChecker struct {
		pass        *analysis.Pass
		annotations *annotation.Annotations
	}
)

func NewInitializationChecker(pass *analysis.Pass, annotations *annotation.Annotations) *InitializationChecker {
	return &InitializationChecker{
		pass:        pass,
		annotations: annotations,
	}
}

func (*InitializationChecker) Types() []ast.Node {
	return []ast.Node{
		new(ast.CompositeLit),
	}
}

func (e *InitializationChecker) InspectionFunc() astutil.InspectionFunc {
	return func(nodeInfo *astutil.NodeInfo) (proceed bool) {
		e.checkInitialization(nodeInfo)
		return true
	}
}

func (e *InitializationChecker) checkInitialization(nodeInfo *astutil.NodeInfo) {
	if !nodeInfo.Push() {
		return
	}

	switch n := nodeInfo.Current().(type) {
	case *ast.CompositeLit:
		if e.entityCompositeLiteral(n) {
			e.pass.Reportf(n.Pos(), "Entityを実装した構造体はEntityが存在するパッケージ以外からcomposite literalで生成することはできません。")
		}
	}
}

func (e *InitializationChecker) entityCompositeLiteral(cl *ast.CompositeLit) bool {
	t := e.pass.TypesInfo.TypeOf(cl)
	if t == nil {
		return false
	}

	s, ok := e.annotations.GetEntity(t)
	if !ok {
		return false
	}

	return s.PackageSummary().Path().String() != e.pass.Pkg.Path()
}
