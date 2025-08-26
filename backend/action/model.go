package action

type Action struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      Status `json:"status"`
}

type Status int

const (
	ToDo Status = iota
	InProgress
	Completed
	Delayed
	Delegated
	Planned
	Someday
	Removed
)

func (s Status) String() string {
	switch s {
	case ToDo:
		return "ToDo"
	case InProgress:
		return "InProgress"
	case Completed:
		return "Completed"
	case Delayed:
		return "Delayed"
	case Delegated:
		return "Delegated"
	case Planned:
		return "Planned"
	case Someday:
		return "Someday"
	case Removed:
		return "Removed"
	default:
		return "Unknown"
	}
}