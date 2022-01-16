package domain

import (
	"github.com/keisuke-m123/gomoduler/analyzer/astutil"
	"github.com/keisuke-m123/gomoduler/analyzer/domain/entity"
	"github.com/keisuke-m123/gomoduler/analyzer/domain/valueobject"
	"github.com/keisuke-m123/gomoduler/annotation"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

type (
	analyzer struct {
		analyzer    *analysis.Analyzer
		annotations *annotation.Annotations
	}
)

func NewEntityAnalyzer(domainPath string) *analysis.Analyzer {
	ea := &analyzer{
		analyzer: &analysis.Analyzer{
			Name:     "gomoduler-domain",
			Doc:      "gomoduler domain analyzer",
			Requires: []*analysis.Analyzer{inspect.Analyzer},
		},
		annotations: annotation.FindAnnotations(domainPath),
	}
	ea.analyzer.Run = ea.analyze
	return ea.analyzer
}

func (a *analyzer) analyze(pass *analysis.Pass) (interface{}, error) {
	astutil.NewInspector(pass.Files, []astutil.InspectionProcessor{
		entity.NewInitializationChecker(pass, a.annotations),
		valueobject.NewInitializationChecker(pass, a.annotations),
		valueobject.NewImmutableChecker(pass, a.annotations),
	}).WithStack()

	return nil, nil
}
