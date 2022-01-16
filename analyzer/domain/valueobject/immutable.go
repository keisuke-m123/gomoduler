package valueobject

import (
	"go/ast"
	"go/types"

	"github.com/keisuke-m123/gomoduler/analyzer/astutil"
	"github.com/keisuke-m123/gomoduler/annotation"
	"golang.org/x/tools/go/analysis"
)

type (
	ImmutableChecker struct {
		pass        *analysis.Pass
		annotations *annotation.Annotations
	}
)

func NewImmutableChecker(pass *analysis.Pass, annotations *annotation.Annotations) *ImmutableChecker {
	return &ImmutableChecker{
		pass:        pass,
		annotations: annotations,
	}
}

func (i *ImmutableChecker) Types() []ast.Node {
	return []ast.Node{
		new(ast.FuncDecl),
	}
}

func (i *ImmutableChecker) InspectionFunc() astutil.InspectionFunc {
	return func(nodeInfo *astutil.NodeInfo) (next bool) {
		i.checkPointerReceiver(nodeInfo)
		return true
	}
}

// 値オブジェクトのメソッドが値レシーバであることをチェックする
func (i *ImmutableChecker) checkPointerReceiver(nodeInfo *astutil.NodeInfo) {
	if !nodeInfo.Push() {
		return
	}

	fd, ok := nodeInfo.Current().(*ast.FuncDecl)
	if !ok {
		return
	}

	if fd.Recv == nil || fd.Recv.NumFields() == 0 {
		return
	}

	field := fd.Recv.List[0]
	t, ok := i.pass.TypesInfo.TypeOf(field.Type).(*types.Pointer)
	if !ok {
		return
	}

	if _, ok := i.annotations.GetValueObjectStruct(t.Elem()); ok {
		i.pass.Reportf(field.Pos(), "値オブジェクトのメソッドは値レシーバである必要があります。")
	}
}