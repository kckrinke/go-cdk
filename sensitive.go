package cdk

type Sensitive interface {
	ProcessEvent(evt Event) EventFlag
}
