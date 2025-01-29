package Triggers

type Trigger interface {
	IsEnabled(state *TriggerState) bool
}
