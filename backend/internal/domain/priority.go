package domain

type Priority string

const (
	PriorityLow    Priority = "baixa"
	PriorityMedium Priority = "media"
	PriorityHigh   Priority = "alta"
)

func (p Priority) Valid() bool {
	switch p {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return true
	default:
		return false
	}
}

func (p Priority) String() string {
	return string(p)
}
