package domain_test

import (
	"path"
	"testing"

	"github.com/keisuke-m123/gomoduler/analyzer/domain"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestEntityFinder(t *testing.T) {
	analysistest.Run(
		t,
		path.Join(analysistest.TestData(), "analyzerproject"),
		domain.NewEntityAnalyzer("./testdata/analyzerproject/src/analyzerproject/"),
		"./...",
	)
}
