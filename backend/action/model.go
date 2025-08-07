package action

type Action struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status 	Status `json:"status"`
}

type Status int
const  (
	Pending Status = iota
)