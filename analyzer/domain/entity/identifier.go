package entity

import (
	"github.com/keisuke-m123/goanalyzer/gocode"
	"github.com/keisuke-m123/gomoduler/analyzer/domain/checker"
	"github.com/keisuke-m123/gomoduler/annotation"
)

type (
	IdentifierChecker struct {
		passInfo    *checker.PassInfo
		annotations *annotation.Annotations
	}
)

func NewIdentifierChecker(passInfo *checker.PassInfo, annotations *annotation.Annotations) *IdentifierChecker {
	return &IdentifierChecker{
		passInfo:    passInfo,
		annotations: annotations,
	}
}

func (ic *IdentifierChecker) Check() {
	for _, s := range annotation.FindEntityStructs(ic.passInfo.Relations()) {
		ic.checkEntity(s)
	}
}

func (ic *IdentifierChecker) checkEntity(entity *gocode.Struct) {
	if !ic.containsIdentifierFieldOnlyOne(entity) {
		ic.passInfo.Pass().Reportf(entity.DefinedPos(), "EntityはIdentifierを実装するフィールドは1つのみ含む必要があります。")
	}
}

func (ic *IdentifierChecker) containsIdentifierFieldOnlyOne(entity *gocode.Struct) bool {
	var count int
	for _, f := range entity.Fields() {
		if _, ok := ic.annotations.GetIdentifierStruct(f.Type().GoType()); ok {
			count++
			continue
		}
		if _, ok := ic.annotations.GetIdentifierDefinedType(f.Type().GoType()); ok {
			count++
			continue
		}
	}
	return count == 1
}
