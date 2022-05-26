package entities

type Article struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

type Articles []Article
