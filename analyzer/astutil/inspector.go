package astutil

import (
	"go/ast"

	"golang.org/x/tools/go/ast/inspector"
)

type (
	NodeStack struct {
		nodes []ast.Node
	}

	NodeStackScanFunc func(n ast.Node) (next bool)

	NodeInfo struct {
		push    bool
		current ast.Node
		stack   NodeStack
	}

	Inspector struct {
		nodeInfo     *NodeInfo
		astInspector *inspector.Inspector
		processors   []InspectionProcessor
	}

	InspectionFunc func(info *NodeInfo) (proceed bool)

	InspectionProcessor interface {
		// Types Inspector 側で検知する必要のある(検査に使用する) ast.Node のリストを返す
		Types() []ast.Node
		// InspectionFunc 検査処理の関数を返す
		InspectionFunc() InspectionFunc
	}
)

func newNodeStack(nodes []ast.Node) NodeStack {
	return NodeStack{
		nodes: nodes,
	}
}

func (n NodeStack) Scan(sf NodeStackScanFunc) {
	for i := len(n.nodes) - 1; i >= 0; i-- {
		if !sf(n.nodes[i]) {
			break
		}
	}
}

func (n NodeStack) ScanFromRoot(sf NodeStackScanFunc) {
	for i := range n.nodes {
		if !sf(n.nodes[i]) {
			break
		}
	}
}

func newNodeInfo() *NodeInfo {
	return &NodeInfo{}
}

func (ni *NodeInfo) Push() bool {
	return ni.push
}

func (ni *NodeInfo) Current() ast.Node {
	return ni.current
}

func (ni *NodeInfo) Stack() NodeStack {
	return ni.stack
}

func NewInspector(files []*ast.File, processors []InspectionProcessor) *Inspector {
	return &Inspector{
		nodeInfo:     newNodeInfo(),
		processors:   processors,
		astInspector: inspector.New(files),
	}
}

func (i *Inspector) WithStack() {
	var fs []InspectionFunc
	for _, p := range i.processors {
		fs = append(fs, p.InspectionFunc())
	}

	i.astInspector.WithStack(i.types(), func(n ast.Node, push bool, stack []ast.Node) (proceed bool) {
		i.nodeInfo.current = n
		i.nodeInfo.push = push
		i.nodeInfo.stack = newNodeStack(stack)

		proceedAll := true
		for _, f := range fs {
			proceed = f(i.nodeInfo)
			proceedAll = proceedAll && proceed
		}
		return proceedAll
	})
}

// 検査処理に使用される ast.Node のリストを返す。
// inspector.Inspector のメソッド内で重複は処理されるので重複除去はここでは行わない。
func (i *Inspector) types() []ast.Node {
	var types []ast.Node
	for _, p := range i.processors {
		for _, t := range p.Types() {
			types = append(types, t)
		}
	}
	return types
}
