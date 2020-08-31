package book

type Book struct {
	ID       string `db:"id" json:"id"`
	ISBN     string `db:"isbn" json:"isbn"`
	Title    string `db:"title" json:"title"`
	Author   string `db:"author" json:"author"`
	Category string `db:"category" json:"category"`
}

type NewBook struct {
	ISBN     string `json:"isbn" validate:"required"`
	Title    string `json:"title" validate:"required"`
	Author   string `json:"author" validate:"required"`
	Category string `json:"category"`
}

type UpdateBook struct {
	ISBN     string  `json:"isbn"`
	Title    string  `json:"title"`
	Author   string  `json:"author"`
	Category *string `json:"category"`
}
