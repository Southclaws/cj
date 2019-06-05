package common

// Action specifies config for an action
type Action struct {
	Schedule string
	Chance   float64
	Call     func() error
}
