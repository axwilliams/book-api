package mock

import (
	"net/http"

	"github.com/axwilliams/books-api/internal/business/book"
	"github.com/axwilliams/books-api/internal/platform/web"
)

type MockBook interface {
	GetAll() ([]book.Book, error)
	GetById(id string) (*book.Book, error)
	Search(sp book.SearchParams, sortOrder string, limit, offset int) ([]book.Book, error)
	Create(bk *book.Book) error
	Update(bk *book.Book) error
	Destroy(id string) error
}

type mockBook struct{}

func NewMockBook() MockBook {
	return &mockBook{}
}

func (mb *mockBook) GetAll() ([]book.Book, error) {
	bs := make([]book.Book, 0)

	bs = append(bs, book.Book{
		ID:       "f4ac7e14-fc8e-4096-b956-34e5a33040f2",
		ISBN:     "978-0241372579",
		Title:    "The Castle",
		Author:   "Franz Kafka",
		Category: "Fiction",
	})

	bs = append(bs, book.Book{
		ID:       "71432eb9-58da-4eae-aa20-ccc49064246f",
		ISBN:     "978-1451673319",
		Title:    "Fahrenheit 451",
		Author:   "Ray Bradbury",
		Category: "Fiction",
	})

	bs = append(bs, book.Book{
		ID:       "562e1fe0-0dde-4717-a008-cd2a699301d2",
		ISBN:     "978-0465025275",
		Title:    "Six Easy Pieces",
		Author:   "Richard Feynman",
		Category: "Science",
	})

	return bs, nil
}

func (mb *mockBook) GetById(id string) (*book.Book, error) {
	if id == "f4ac7e14-fc8e-4096-b956-34e5a33040f2" {
		return &book.Book{
			ID:       "f4ac7e14-fc8e-4096-b956-34e5a33040f2",
			ISBN:     "978-0241372579",
			Title:    "The Castle",
			Author:   "Franz Kafka",
			Category: "Fiction",
		}, nil
	}

	return nil, book.ErrNoBookFound
}

func (mb *mockBook) Search(sp book.SearchParams, sortOrder string, limit, offset int) ([]book.Book, error) {
	bs := make([]book.Book, 0)

	if sp.ISBN == "978-0241372579" {
		bs = append(bs, book.Book{
			ID:       "f4ac7e14-fc8e-4096-b956-34e5a33040f2",
			ISBN:     "978-0241372579",
			Title:    "The Castle",
			Author:   "Franz Kafka",
			Category: "Fiction",
		})
	}

	if sp.Title == "The Castle" {
		bs = append(bs, book.Book{
			ID:       "f4ac7e14-fc8e-4096-b956-34e5a33040f2",
			ISBN:     "978-0241372579",
			Title:    "The Castle",
			Author:   "Franz Kafka",
			Category: "Fiction",
		})
	}

	if sp.Author == "Franz Kafka" {
		bs = append(bs, book.Book{
			ID:       "f4ac7e14-fc8e-4096-b956-34e5a33040f2",
			ISBN:     "978-0241372579",
			Title:    "The Castle",
			Author:   "Franz Kafka",
			Category: "Fiction",
		})
	}

	if sp.Category == "Fiction" {
		bs = append(bs, book.Book{
			ID:       "f4ac7e14-fc8e-4096-b956-34e5a33040f2",
			ISBN:     "978-0241372579",
			Title:    "The Castle",
			Author:   "Franz Kafka",
			Category: "Fiction",
		})

		bs = append(bs, book.Book{
			ID:       "71432eb9-58da-4eae-aa20-ccc49064246f",
			ISBN:     "978-1451673319",
			Title:    "Fahrenheit 451",
			Author:   "Ray Bradbury",
			Category: "Fiction",
		})
	}

	if sortOrder == "author desc" {
		bs = append(bs, book.Book{
			ID:       "f4ac7e14-fc8e-4096-b956-34e5a33040f2",
			ISBN:     "978-0241372579",
			Title:    "The Castle",
			Author:   "Franz Kafka",
			Category: "Fiction",
		})

		bs = append(bs, book.Book{
			ID:       "562e1fe0-0dde-4717-a008-cd2a699301d2",
			ISBN:     "978-0465025275",
			Title:    "Six Easy Pieces",
			Author:   "Richard Feynman",
			Category: "Science",
		})

		bs = append(bs, book.Book{
			ID:       "71432eb9-58da-4eae-aa20-ccc49064246f",
			ISBN:     "978-1451673319",
			Title:    "Fahrenheit 451",
			Author:   "Ray Bradbury",
			Category: "Fiction",
		})
	}

	if limit == 2 && offset == 1 {
		bs = append(bs, book.Book{
			ID:       "71432eb9-58da-4eae-aa20-ccc49064246f",
			ISBN:     "978-1451673319",
			Title:    "Fahrenheit 451",
			Author:   "Ray Bradbury",
			Category: "Fiction",
		})

		bs = append(bs, book.Book{
			ID:       "562e1fe0-0dde-4717-a008-cd2a699301d2",
			ISBN:     "978-0465025275",
			Title:    "Six Easy Pieces",
			Author:   "Richard Feynman",
			Category: "Science",
		})
	}

	return bs, nil
}

func (mb *mockBook) Create(bk *book.Book) error {
	return nil
}

func (mb *mockBook) Update(bk *book.Book) error {
	return nil
}

func (mb *mockBook) Destroy(id string) error {
	if id == "f4ac7e14-fc8e-4096-b956-34e5a33040f2" {
		return nil
	}

	return web.NewRequestError(book.ErrNoAffect, http.StatusGone)
}
