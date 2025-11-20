package workflow

import "time"

type Action struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      Status     `json:"status"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Context     string     `json:"context,omitempty"`
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