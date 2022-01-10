package annotation

type (
	// identifier を実装する構造体を一意な識別子として扱う。
	identifier interface {
		ImplAsIdentifier()
	}
	// entity を実装する構造体をエンティティとして扱う。
	entity interface {
		ImplAsEntity()
	}
	// valueObject を実装する構造体を値オブジェクトとして扱う。
	valueObject interface {
		ImplAsValueObject()
	}
	// valueObjectGenerator は値オブジェクトのジェネレータとして扱う
	valueObjectGenerator interface {
		ImplAsValueObjectGenerator()
	}
	// Identifier は一意な識別子として振る舞うStructに埋め込む。
	Identifier struct{}
	// Entity はエンティティとして振る舞うStructに埋め込む。
	Entity struct{}
	// ValueObject は値オブジェクトとして振る舞うStructに埋め込む。
	ValueObject struct{}
	// ValueObjectGenerator は値オブジェクトのジェネレータとして振る舞うStructに埋め込む。
	ValueObjectGenerator struct{}
)

func (Identifier) ImplAsIdentifier()                     {}
func (Entity) ImplAsEntity()                             {}
func (ValueObject) ImplAsValueObject()                   {}
func (ValueObjectGenerator) ImplAsValueObjectGenerator() {}
