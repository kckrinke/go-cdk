package cdk

const (
	ITypeObject       ITypeTag = "object"
	SignalDestroy     Signal   = "destroy"
	SignalSetProperty Signal   = "set-property"
)

func init() {
	CursesITypeRegistry.AddType(ITypeObject)
}

// Basic object type
type Object interface {
	Signaling

	Destroy()

	GetTheme() Theme
	SetTheme(theme Theme)

	SetProperty(name string, value interface{})
	GetProperty(name string) interface{}
	GetPropertyAsBool(name string, def bool) bool
	GetPropertyAsString(name string, def string) string
	GetPropertyAsInt(name string, def int) int
	GetPropertyAsFloat(name string, def float64) float64
}

type CObject struct {
	CSignaling

	theme Theme

	properties map[string]interface{}
}

func (o *CObject) Init() (already bool) {
	o.SetIType(ITypeObject)
	if o.CSignaling.Init() {
		return true
	}
	o.properties = make(map[string]interface{})
	return false
}

func (o *CObject) Destroy() {
	o.Emit(SignalDestroy, o)
	o.DestroyObject()
}

func (o *CObject) GetTheme() Theme {
	return o.theme
}

func (o *CObject) SetTheme(theme Theme) {
	o.theme = theme
}

// set the value for a named property
func (o *CObject) SetProperty(name string, value interface{}) {
	if f := o.Emit(SignalSetProperty, o, name, value); f == EVENT_PASS {
		o.properties[name] = value
	}
}

// return the named property value
func (o *CObject) GetProperty(name string) interface{} {
	if v, ok := o.properties[name]; ok {
		return v
	}
	return nil
}

// return the named property value as a string
func (o *CObject) GetPropertyAsBool(name string, def bool) bool {
	v := o.GetProperty(name)
	if v, ok := v.(bool); ok {
		return v
	}
	return def
}

// return the named property value as a string
func (o *CObject) GetPropertyAsString(name string, def string) string {
	v := o.GetProperty(name)
	if v, ok := v.(string); ok {
		return v
	}
	return def
}

// return the named property value as an integer
func (o *CObject) GetPropertyAsInt(name string, def int) int {
	v := o.GetProperty(name)
	if v, ok := v.(int); ok {
		return v
	}
	return def
}

// return the named property value as a float
func (o *CObject) GetPropertyAsFloat(name string, def float64) float64 {
	v := o.GetProperty(name)
	if v, ok := v.(float64); ok {
		return v
	}
	return def
}
