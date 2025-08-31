package capture


type Thing struct {
	ID int   `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Status Status `json:"status"`
}

type Status int

const (
	Pending Status = iota
	Someday
	Done
)

func (s Status) String() string {
	switch s {
	case Pending:
		return "Pending"
	case Someday:
		return "Someday"
	case Done:
		return "Done"
	default:
		return "Unknown"
	}	
}