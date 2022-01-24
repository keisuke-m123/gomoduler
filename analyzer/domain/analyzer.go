package domain

import (
	"github.com/keisuke-m123/gomoduler/analyzer/astutil"
	"github.com/keisuke-m123/gomoduler/analyzer/domain/checker"
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

func NewDomainAnalyzer(domainPath []string) *analysis.Analyzer {
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
	passInfo := checker.NewPassInfo(pass)
	astutil.NewInspector(pass.Files, []astutil.InspectionProcessor{
		entity.NewInitializationChecker(passInfo, a.annotations),
		valueobject.NewInitializationChecker(passInfo, a.annotations),
		valueobject.NewImmutableChecker(passInfo, a.annotations),
	}).WithStack()

	for _, c := range []checker.SimpleChecker{
		entity.NewExportFieldChecker(passInfo),
		entity.NewIdentifierChecker(passInfo, a.annotations),
		valueobject.NewExportedFieldChecker(passInfo),
	} {
		c.Check()
	}

	return nil, nil
}
