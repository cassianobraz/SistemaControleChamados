package domain

type Status string

const (
	StatusOpen       Status = "aberto"
	StatusInProgress Status = "em_andamento"
	StatusResolved   Status = "resolvido"
	StatusClosed     Status = "fechado"
)

func (s Status) Valid() bool {
	switch s {
	case StatusOpen, StatusInProgress, StatusResolved, StatusClosed:
		return true
	default:
		return false
	}
}

func (s Status) IsOpen() bool {
	return s == StatusOpen || s == StatusInProgress
}

func (s Status) String() string {
	return string(s)
}
