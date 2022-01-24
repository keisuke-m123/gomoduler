package domain_test

import (
	"path"
	"testing"

	"github.com/keisuke-m123/gomoduler/analyzer/domain"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestDomainAnalyzer(t *testing.T) {
	analysistest.Run(
		t,
		path.Join(analysistest.TestData(), "analyzerproject"),
		domain.NewDomainAnalyzer([]string{"./testdata/analyzerproject/src/analyzerproject/"}),
		"./...",
	)
}
