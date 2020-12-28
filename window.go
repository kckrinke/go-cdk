package cdk

// basically what ctk.Window is right now

const (
	ITypeWindow ITypeTag = "window"
)

func init() {
	CursesITypeRegistry.AddType(ITypeWindow)
}

// Basic window interface
type Window interface {
	Object

	GetDisplay() Display
	SetDisplay(d Display)

	Draw(view *View) EventFlag
	ProcessEvent(evt Event) EventFlag
}

// Basic window type
type CWindow struct {
	CObject

	title string

	display Display
}

func (w *CWindow) Init() bool {
	w.SetIType(ITypeWindow)
	if w.CObject.Init() {
		return true
	}
	CursesITypeRegistry.AddTypeItem(ITypeWindow, w)
	return false
}

func (w *CWindow) GetDisplay() Display {
	return w.display
}

func (w *CWindow) SetDisplay(d Display) {
	w.display = d
}

func (w *CWindow) Draw(view *View) EventFlag {
	w.LogDebug("method not implemented")
	return EVENT_PASS
}

func (w *CWindow) ProcessEvent(evt Event) EventFlag {
	w.LogDebug("method not implemented")
	return EVENT_PASS
}
