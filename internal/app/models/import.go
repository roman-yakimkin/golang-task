package models

type ImportTask struct {
	AuthorID int      `json:"author_id"`
	Title    string   `json:"title"`
	Body     string   `json:"body"`
	Emails   []string `json:"emails"`
}

type UpdateTask struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}
