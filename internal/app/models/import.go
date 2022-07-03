package models

type ImportTask struct {
	AuthorID string   `json:"author_id"`
	Title    string   `json:"title"`
	Body     string   `json:"body"`
	Emails   []string `json:"emails"`
}

type UpdateTask struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}
