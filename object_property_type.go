package cdk

const (
	BoolProperty   PropertyType = "bool"
	StringProperty PropertyType = "string"
	IntProperty    PropertyType = "int"
	FloatProperty  PropertyType = "float"
	StructProperty PropertyType = "struct"
)

type PropertyType string

func (p PropertyType) String() string {
	return string(p)
}
