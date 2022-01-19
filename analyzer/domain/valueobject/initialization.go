package valueobject

import (
	"go/ast"
	"strings"

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
		new(ast.CallExpr),
	}
}

func (v *InitializationChecker) InspectionFunc() astutil.InspectionFunc {
	return func(nodeInfo *astutil.NodeInfo) (proceed bool) {
		v.checkInitialization(nodeInfo)
		return true
	}
}

func (v *InitializationChecker) checkInitialization(nodeInfo *astutil.NodeInfo) {
	if !nodeInfo.Push() {
		return
	}

	switch n := nodeInfo.Current().(type) {
	case *ast.CompositeLit:
		if !v.compositeLiteralInitialization(n) {
			return
		}
		if v.inGeneratorFunc(nodeInfo) {
			return
		}
		v.passInfo.Pass().Reportf(n.Pos(), "ValueObjectを実装した構造体はValueObjectが存在するパッケージ以外からcomposite literalで生成することはできません。")
	case *ast.CallExpr:
		if !v.definedTypeInitialization(n) {
			return
		}
		if v.inGeneratorFunc(nodeInfo) {
			return
		}
		v.passInfo.Pass().Reportf(n.Pos(), "ValueObjectを実装した構造体はValueObjectが存在するパッケージ以外から直接生成することはできません。")
	}
}

func (v *InitializationChecker) compositeLiteralInitialization(cl *ast.CompositeLit) bool {
	t := v.passInfo.Pass().TypesInfo.TypeOf(cl)
	if t == nil {
		return false
	}

	s, ok := v.annotations.GetValueObjectStruct(t)
	if !ok {
		return false
	}

	return s.PackageSummary().Path().String() != v.passInfo.Pass().Pkg.Path()
}

func (v *InitializationChecker) definedTypeInitialization(ce *ast.CallExpr) bool {
	t := v.passInfo.Pass().TypesInfo.TypeOf(ce.Fun)
	if t == nil {
		return false
	}

	// defined typeの呼び出しのみを補足するため、prefixをみて通常の関数/メソッド呼び出しと区別
	if strings.HasPrefix(t.String(), "func(") {
		return false
	}

	d, ok := v.annotations.GetValueObjectDefinedType(t)
	if !ok {
		return false
	}

	return d.PackageSummary().Path().String() != v.passInfo.Pass().Pkg.Path()
}

func (v *InitializationChecker) inGeneratorFunc(nodeInfo *astutil.NodeInfo) bool {
	inGenerator := false
	nodeInfo.Stack().ScanFromRoot(func(n ast.Node) (next bool) {
		fd, ok := n.(*ast.FuncDecl)
		if !ok {
			return !inGenerator
		}

		if fd.Recv == nil || !v.generatorInFields(fd.Recv) {
			return !inGenerator
		}

		inGenerator = true
		return !inGenerator
	})

	return inGenerator
}

func (v *InitializationChecker) generatorInFields(fl *ast.FieldList) bool {
	for _, f := range fl.List {
		if v.isGenerator(f.Type) {
			return true
		}
	}
	return false
}

func (v *InitializationChecker) isGenerator(expr ast.Expr) bool {
	if expr == nil {
		return false
	}

	switch t := expr.(type) {
	case *ast.StarExpr:
		return v.isGenerator(t.X)
	default:
		df := v.passInfo.Pass().TypesInfo.TypeOf(t)
		if df == nil {
			return false
		}
		_, ok := v.annotations.GetValueObjectGenerator(df)
		return ok
	}
}
