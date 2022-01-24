package valueobject

import (
	"github.com/keisuke-m123/gomoduler/analyzer/domain/checker"
	"github.com/keisuke-m123/gomoduler/annotation"
)

type (
	ExportedFieldChecker struct {
		passInfo *checker.PassInfo
	}
)

func NewExportedFieldChecker(passInfo *checker.PassInfo) *ExportedFieldChecker {
	return &ExportedFieldChecker{passInfo: passInfo}
}

func (e *ExportedFieldChecker) Check() {
	for _, vo := range annotation.FindValueObjectStructs(e.passInfo.Relations()) {
		for _, f := range vo.Fields() {
			if f.Embedded() && annotation.DomainStruct(f.Type()) {
				continue
			}
			if f.Exported() {
				e.passInfo.Pass().Reportf(f.DefinedPos(), "ValueObjectはExportedなフィールドを定義することはできません。")
			}
		}
	}
}
