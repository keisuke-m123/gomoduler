package annotation

import (
	"go/types"
	"strings"

	"github.com/keisuke-m123/goanalyzer/gocode"
	"github.com/keisuke-m123/gomoduler/internal/relations"
	"github.com/spf13/afero"
)

type (
	structMap struct {
		m map[string]*gocode.Struct
	}

	definedTypeMap struct {
		m map[string]*gocode.DefinedType
	}

	structAndDefinedTypeMap struct {
		sm *structMap
		dm *definedTypeMap
	}

	Annotations struct {
		foundEntities             *structMap
		foundValueObjectGenerator *structMap
		foundValueObjects         *structAndDefinedTypeMap
		foundIdentifiers          *structAndDefinedTypeMap
	}
)

func newStructMap() *structMap {
	return &structMap{
		m: make(map[string]*gocode.Struct),
	}
}

func (sm *structMap) add(s *gocode.Struct) {
	sm.m[strings.Join([]string{s.PackageSummary().Path().String(), s.Name().String()}, ".")] = s
}

func (sm *structMap) get(t types.Type) (*gocode.Struct, bool) {
	s, ok := sm.m[t.String()]
	return s, ok
}

func newDefinedType() *definedTypeMap {
	return &definedTypeMap{
		m: make(map[string]*gocode.DefinedType),
	}
}

func (dm *definedTypeMap) add(d *gocode.DefinedType) {
	dm.m[strings.Join([]string{d.PackageSummary().Path().String(), d.Name().String()}, ".")] = d
}

func (dm *definedTypeMap) get(t types.Type) (*gocode.DefinedType, bool) {
	d, ok := dm.m[t.String()]
	return d, ok
}

func newStructAndDefinedTypeMap() *structAndDefinedTypeMap {
	return &structAndDefinedTypeMap{
		sm: newStructMap(),
		dm: newDefinedType(),
	}
}

func (v *structAndDefinedTypeMap) addStruct(s *gocode.Struct) {
	v.sm.add(s)
}

func (v *structAndDefinedTypeMap) addDefinedType(d *gocode.DefinedType) {
	v.dm.add(d)
}

func (v *structAndDefinedTypeMap) getStruct(t types.Type) (*gocode.Struct, bool) {
	return v.sm.get(t)
}

func (v *structAndDefinedTypeMap) getDefinedType(t types.Type) (*gocode.DefinedType, bool) {
	return v.dm.get(t)
}

func FindAnnotations(path string) *Annotations {
	r, err := gocode.LoadRelations(&gocode.LoadOptions{
		FileSystem:  afero.NewOsFs(),
		Directories: []string{path},
		Recursive:   true,
	})
	if err != nil {
		panic(err)
	}
	return &Annotations{
		foundEntities:             findEntities(r),
		foundValueObjects:         findValueObjects(r),
		foundValueObjectGenerator: findValueObjectGenerator(r),
		foundIdentifiers:          findIdentifiers(r),
	}
}

func (f *Annotations) GetEntity(t types.Type) (*gocode.Struct, bool) {
	return f.foundEntities.get(t)
}

func (f *Annotations) GetEntitiesAll() []*gocode.Struct {
	var structs []*gocode.Struct
	for _, st := range f.foundEntities.m {
		structs = append(structs, st)
	}
	return structs
}

func (f *Annotations) GetValueObjectStruct(t types.Type) (*gocode.Struct, bool) {
	return f.foundValueObjects.getStruct(t)
}

func (f *Annotations) GetValueObjectDefinedType(t types.Type) (*gocode.DefinedType, bool) {
	return f.foundValueObjects.getDefinedType(t)
}

func (f *Annotations) GetValueObjectGenerator(t types.Type) (*gocode.Struct, bool) {
	return f.foundValueObjectGenerator.get(t)
}

func (f *Annotations) GetIdentifierStruct(t types.Type) (*gocode.Struct, bool) {
	return f.foundIdentifiers.getStruct(t)
}

func (f *Annotations) GetIdentifierDefinedType(t types.Type) (*gocode.DefinedType, bool) {
	return f.foundIdentifiers.getDefinedType(t)
}

func FindEntityStructs(r *gocode.Relations) []*gocode.Struct {
	var structs []*gocode.Struct
	for _, st := range findEntities(r).m {
		structs = append(structs, st)
	}
	return structs
}

func findEntities(r *gocode.Relations) *structMap {
	sm := newStructMap()
	for _, s := range r.Structs().StructAll() {
		if s.Implements(relations.GoCodeInterfaceEntity()) {
			sm.add(s)
		}
	}
	return sm
}

func FindValueObjectStructs(r *gocode.Relations) []*gocode.Struct {
	var structs []*gocode.Struct
	for _, st := range findValueObjects(r).sm.m {
		structs = append(structs, st)
	}
	return structs
}

func findValueObjects(r *gocode.Relations) *structAndDefinedTypeMap {
	sdm := newStructAndDefinedTypeMap()
	for _, s := range r.Structs().StructAll() {
		if s.Implements(relations.GoCodeInterfaceValueObject()) {
			sdm.addStruct(s)
		}
	}
	for _, d := range r.DefinedTypes().DefinedTypeAll() {
		if d.Implements(relations.GoCodeInterfaceValueObject()) {
			sdm.addDefinedType(d)
		}
	}
	return sdm
}

func findValueObjectGenerator(r *gocode.Relations) *structMap {
	sm := newStructMap()
	for _, s := range r.Structs().StructAll() {
		if s.Implements(relations.GoCodeInterfaceValueObjectGenerator()) {
			sm.add(s)
		}
	}
	return sm
}

func findIdentifiers(r *gocode.Relations) *structAndDefinedTypeMap {
	ids := newStructAndDefinedTypeMap()
	for _, s := range r.Structs().StructAll() {
		if s.Implements(relations.GoCodeInterfaceIdentifier()) {
			ids.addStruct(s)
		}
	}
	for _, d := range r.DefinedTypes().DefinedTypeAll() {
		if d.Implements(relations.GoCodeInterfaceIdentifier()) {
			ids.addDefinedType(d)
		}
	}
	return ids
}
