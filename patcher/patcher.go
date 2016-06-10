package patcher

// Patcher is in interface or types that do patches.
type Patcher interface {
	Patch() *PatchResult
	GetSource() string
	GetDestination() string
	SetSource(source string)
	SetDestination(destination string)
}
