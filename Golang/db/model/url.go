package url

type Url struct {
	Id       string `json:"id" db:"id"`
	Original string `json:"original" db:"original"`
	Accesses int    `json:"accesses" db:"accesses"`
}
