package entity

import (
	"go/ast"

	"github.com/keisuke-m123/gomoduler/analyzer/astutil"
	"github.com/keisuke-m123/gomoduler/analyzer/domain/checker"
	"github.com/keisuke-m123/gomoduler/annotation"
)

type (
	InitializationChecker struct {
		passInfo    *checker.PassInfo
		annotations *annotation.Annotations
	}
)

func NewInitializationChecker(passInfo *checker.PassInfo, annotations *annotation.Annotations) *InitializationChecker {
	return &InitializationChecker{
		passInfo:    passInfo,
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
		e.check(nodeInfo)
		return true
	}
}

func (e *InitializationChecker) check(nodeInfo *astutil.NodeInfo) {
	if !nodeInfo.Push() {
		return
	}

	switch n := nodeInfo.Current().(type) {
	case *ast.CompositeLit:
		if e.entityCompositeLiteral(n) {
			e.passInfo.Pass().Reportf(n.Pos(), "Entityを実装した構造体はEntityが存在するパッケージ以外からcomposite literalで生成することはできません。")
		}
	}
}

func (e *InitializationChecker) entityCompositeLiteral(cl *ast.CompositeLit) bool {
	t := e.passInfo.Pass().TypesInfo.TypeOf(cl)
	if t == nil {
		return false
	}

	s, ok := e.annotations.GetEntity(t)
	if !ok {
		return false
	}

	return s.PackageSummary().Path().String() != e.passInfo.Pass().Pkg.Path()
}
