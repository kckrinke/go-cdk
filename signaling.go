package cdk

const (
	ITypeSignaling ITypeTag = "signaling"
)

type SignalCallbackFn func(data []interface{}, argv ...interface{}) EventFlag

type SignalCallbackData []interface{}

type CSignalListener struct {
	n Signal
	c SignalCallbackFn
	d SignalCallbackData
}

type Signal string

type Signaling interface {
	TypeItem

	Connect(signal, handle Signal, c SignalCallbackFn, data ...interface{})
	Disconnect(signal, handle Signal)
	Emit(signal Signal, argv ...interface{}) EventFlag
	StopSignal(signal Signal)
	IsSignalStopped(signal Signal) bool
	PassSignal(signal Signal)
	IsSignalPassed(signal Signal) bool
	ResumeSignal(signal Signal)
}

var (
	cdkSignalListeners = make(map[int]map[Signal][]*CSignalListener)
)

func init() {
	CursesITypeRegistry.AddType(ITypeSignaling)
}

type CSignaling struct {
	CTypeItem

	stopSignals []Signal
	passSignals []Signal
}

func (o *CSignaling) Init() (already bool) {
 	o.SetIType(ITypeSignaling)
	if o.CTypeItem.Init() {
		return true
	}
	CursesITypeRegistry.AddTypeItem(ITypeSignaling, o)
	return false
}

// Connect callback to signal, identified by handle
func (o *CSignaling) Connect(signal, handle Signal, c SignalCallbackFn, data ...interface{}) {
	oid := o.ObjectID()
	if _, ok := cdkSignalListeners[oid]; !ok {
		cdkSignalListeners[oid] = make(map[Signal][]*CSignalListener)
	}
	cdkSignalListeners[oid][signal] = append(
		cdkSignalListeners[oid][signal],
		&CSignalListener{
			handle,
			c,
			data,
		},
	)
}

// Disconnect callback from signal identified by handle
func (o *CSignaling) Disconnect(signal, handle Signal) {
	oid := o.ObjectID()
	id := 0
	for i, s := range cdkSignalListeners[oid][signal] {
		if s.n == handle {
			id = i
			break
		}
	}
	o.LogDebug("disconnecting(%v) from signal(%v)", handle, signal)
	cdkSignalListeners[oid][signal] = append(
		cdkSignalListeners[oid][signal][:id],
		cdkSignalListeners[oid][signal][id+1:]...,
	)
}

// Emit a signal event to all connected listener callbacks
func (o *CSignaling) Emit(signal Signal, argv ...interface{}) EventFlag {
	oid := o.ObjectID()
	if o.IsSignalStopped(signal) {
		return EVENT_STOP
	}
	if o.IsSignalPassed(signal) {
		return EVENT_PASS
	}
	if listeners, ok := cdkSignalListeners[oid][signal]; ok {
		for _, s := range listeners {
			r := s.c(s.d, argv...)
			if r == EVENT_STOP {
				o.LogTrace("emit(%v) stopped by listener(%v)", signal, s.n)
				return EVENT_STOP
			}
		}
	}
	return EVENT_PASS
}

// Disable propagation of the given signal
func (o *CSignaling) StopSignal(signal Signal) {
	if !o.IsSignalStopped(signal) {
		o.LogDebug("stopping signal(%v)", signal)
		o.stopSignals = append(o.stopSignals, signal)
	}
}

func (o *CSignaling) IsSignalStopped(signal Signal) bool {
	return o.getSignalStopIndex(signal) >= 0
}

func (o *CSignaling) getSignalStopIndex(signal Signal) int {
	for idx, stop := range o.stopSignals {
		if signal == stop {
			return idx
		}
	}
	return -1
}

func (o *CSignaling) PassSignal(signal Signal) {
	if !o.IsSignalPassed(signal) {
		o.LogDebug("passing signal(%v)", signal)
		o.passSignals = append(o.passSignals, signal)
	}
}

func (o *CSignaling) IsSignalPassed(signal Signal) bool {
	return o.getSignalPassIndex(signal) >= 0
}

func (o *CSignaling) getSignalPassIndex(signal Signal) int {
	for idx, stop := range o.passSignals {
		if signal == stop {
			return idx
		}
	}
	return -1
}

// Enable propagation of the given signal
func (o *CSignaling) ResumeSignal(signal Signal) {
	id := o.getSignalStopIndex(signal)
	if id >= 0 {
		o.LogDebug("resuming signal(%v) from being stopped", signal)
		o.stopSignals = append(
			o.stopSignals[:id],
			o.stopSignals[id+1:]...,
		)
		return
	}
	id = o.getSignalPassIndex(signal)
	if id >= 0 {
		o.LogDebug("resuming signal(%v) from being passed", signal)
		o.passSignals = append(
			o.passSignals[:id],
			o.passSignals[id+1:]...,
		)
		return
	}
	o.LogWarn("signal(%v) already resumed", signal)
}
