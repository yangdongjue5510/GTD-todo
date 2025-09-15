package capture


type Thing struct {
	ID int   `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Status Status `json:"status"`
}

type Status int

const (
	Active Status = iota
	Done
)

func (s *Thing) Process() {
	s.Status = Done
}

func (s Status) String() string {
	switch s {
	case Active:
		return "Active"
	case Done:
		return "Done"
	default:
		return "Unknown"
	}	
}