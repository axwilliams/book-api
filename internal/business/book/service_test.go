package book_test

import (
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/axwilliams/books-api/internal/business/book"
	"github.com/axwilliams/books-api/internal/test"
)

var bookService book.Service

func TestMain(m *testing.M) {
	db, container := test.Setup()

	bookRepository := book.NewRepository(db)
	bookService = book.NewService(bookRepository)

	e := m.Run()

	test.Teardown(db, container)
	os.Exit(e)
}

func TestGetById(t *testing.T) {
	bk := &book.Book{
		ID:       "f4ac7e14-fc8e-4096-b956-34e5a33040f2",
		ISBN:     "978-0241372579",
		Title:    "The Castle",
		Author:   "Franz Kafka",
		Category: "Fiction",
	}

	res, err := bookService.GetById(bk.ID)
	if err != nil {
		t.Fatal(err)
	}

	if ok := reflect.DeepEqual(res, bk); !ok {
		t.Fatalf("\t%s\tError finding book: want %v got %v", test.Failed, bk, res)
	}
	t.Logf("\t%s\tBook found", test.Success)
}

func TestCreate(t *testing.T) {
	nb := &book.NewBook{
		ISBN:     "978-0099448792",
		Title:    "The Wind-Up Bird Chronicle",
		Author:   "Haruki Murakami",
		Category: "Fiction",
	}

	bk, err := bookService.Create(nb)
	if err != nil {
		t.Fatal(err)
	}

	res, err := bookService.GetById(bk.ID)
	if err != nil {
		log.Fatal(err)
	}

	expected := &book.Book{
		ID:       bk.ID,
		ISBN:     "978-0099448792",
		Title:    "The Wind-Up Bird Chronicle",
		Author:   "Haruki Murakami",
		Category: "Fiction",
	}

	if ok := reflect.DeepEqual(expected, res); !ok {
		t.Fatalf("\t%s\tError creating book: want %v got %v", test.Failed, expected, res)
	}
	t.Logf("\t%s\tBook created", test.Success)
}

func TestUpdate(t *testing.T) {
	ID := "562e1fe0-0dde-4717-a008-cd2a699301d2"

	c := "Physics"

	ub := book.UpdateBook{
		ISBN:     "978-0465025268",
		Title:    "Six Not-So-Easy Pieces",
		Category: &c,
	}

	if err := bookService.Update(ID, ub); err != nil {
		t.Fatal(err)
	}

	res, err := bookService.GetById(ID)
	if err != nil {
		t.Fatal(err)
	}

	expected := &book.Book{
		ID:       ID,
		ISBN:     "978-0465025268",
		Title:    "Six Not-So-Easy Pieces",
		Author:   "Richard Feynman",
		Category: "Physics",
	}

	if ok := reflect.DeepEqual(expected, res); !ok {
		t.Fatalf("\t%s\tError updating book: want %v got %v", test.Failed, expected, res)
	}
	t.Logf("\t%s\tBook updated", test.Success)
}

func TestDestroy(t *testing.T) {
	ID := "71432eb9-58da-4eae-aa20-ccc49064246f"

	if err := bookService.Destroy(ID); err != nil {
		t.Fatal(err)
	}

	if _, err := bookService.GetById(ID); err != book.ErrNoBookFound {
		t.Fatalf("\t%s\tError destroying book", test.Failed)
	}
	t.Logf("\t%s\tBook destroyed", test.Success)
}
