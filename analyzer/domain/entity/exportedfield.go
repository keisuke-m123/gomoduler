package entity

import (
	"github.com/keisuke-m123/gomoduler/analyzer/domain/checker"
	"github.com/keisuke-m123/gomoduler/annotation"
)

type (
	ExportFieldChecker struct {
		passInfo *checker.PassInfo
	}
)

func NewExportFieldChecker(passInfo *checker.PassInfo) *ExportFieldChecker {
	return &ExportFieldChecker{
		passInfo: passInfo,
	}
}

func (e ExportFieldChecker) Check() {
	for _, entity := range annotation.FindEntityStructs(e.passInfo.Relations()) {
		for _, f := range entity.Fields() {
			if f.Embedded() && annotation.DomainStruct(f.Type()) {
				continue
			}
			if f.Exported() {
				e.passInfo.Pass().Reportf(f.DefinedPos(), "EntityはExportedなフィールドを定義することはできません。")
			}
		}
	}
}
