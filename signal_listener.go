package cdk

type SignalCallbackFn func(data []interface{}, argv ...interface{}) EventFlag

type SignalCallbackData []interface{}

type CSignalListener struct {
	n Signal
	c SignalCallbackFn
	d SignalCallbackData
}
