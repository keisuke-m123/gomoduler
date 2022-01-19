package checker

import (
	"github.com/keisuke-m123/goanalyzer/gocode"
	"golang.org/x/tools/go/analysis"
)

type (
	PassInfo struct {
		pass      *analysis.Pass
		relations *gocode.Relations
	}
)

func NewPassInfo(pass *analysis.Pass) *PassInfo {
	return &PassInfo{
		pass:      pass,
		relations: gocode.LoadRelationsFromAnalysis(pass),
	}
}

func (p *PassInfo) Pass() *analysis.Pass {
	return p.pass
}

func (p *PassInfo) Relations() *gocode.Relations {
	return p.relations
}
