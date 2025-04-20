package models

// Question represents a single assessment question
type Question struct {
	ID       string   `json:"id"`
	Text     string   `json:"text"`
	Category string   `json:"category"`
	Options  []Option `json:"options"`
	Weight   int      `json:"weight"`
}

// Option represents a possible answer to a question
type Option struct {
	ID     string `json:"id"`
	Text   string `json:"text"`
	Points int    `json:"points"`
}
