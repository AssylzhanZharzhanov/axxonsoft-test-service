package domain

// Status - specifies a particular sex type.
type Status StatusID

// Supported question types.
const (
	StatusNew       Status = 1
	StatusInProcess Status = 2
	StatusError     Status = 3
	StatusDone      Status = 4
)
