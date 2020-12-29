package cdk

type SignalListenerFn func(data []interface{}, argv ...interface{}) EventFlag

type SignalListenerData []interface{}

type CSignalListener struct {
	n Signal
	c SignalListenerFn
	d SignalListenerData
}
