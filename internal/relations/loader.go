package relations

import (
	"github.com/keisuke-m123/goanalyzer/gocode"
	"github.com/spf13/afero"
)

var relations *gocode.Relations

func init() {
	r, err := gocode.LoadRelations(&gocode.LoadOptions{
		FileSystem:  afero.NewOsFs(),
		Directories: []string{"../../annotation"},
		Recursive:   false,
	})
	if err != nil {
		panic(err)
	}
	relations = r
}

const annotationPkgName = "annotation"

func GoCodeInterfaceIdentifier() *gocode.Interface {
	iface, ok := relations.Interfaces().Get(annotationPkgName, "identifier")
	if !ok {
		panic("interface identifier not found")
	}
	return iface
}

func GoCodeInterfaceEntity() *gocode.Interface {
	iface, ok := relations.Interfaces().Get(annotationPkgName, "entity")
	if !ok {
		panic("interface entity not found")
	}
	return iface
}

func GoCodeInterfaceValueObject() *gocode.Interface {
	iface, ok := relations.Interfaces().Get(annotationPkgName, "valueObject")
	if !ok {
		panic("interface valueObject not found")
	}
	return iface
}

func GoCodeInterfaceValueObjectGenerator() *gocode.Interface {
	iface, ok := relations.Interfaces().Get(annotationPkgName, "valueObjectGenerator")
	if !ok {
		panic("interface valueObjectGenerator not found")
	}
	return iface
}
