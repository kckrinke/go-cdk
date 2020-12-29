package cdk

// Imaginary Type Tag

type TypeTag interface {
	Tag() CTypeTag
	String() string
}

// Used to denote a concrete type identity
type CTypeTag string

func NewTypeTag(tag string) TypeTag {
	return CTypeTag(tag)
}

func (tag CTypeTag) Tag() CTypeTag {
	return tag
}

// Stringer interface implementation
func (tag CTypeTag) String() string {
	return string(tag)
}
