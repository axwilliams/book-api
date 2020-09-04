package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/axwilliams/books-api/cmd/books-api/handlers"
	"github.com/axwilliams/books-api/internal/business/book"
	"github.com/axwilliams/books-api/internal/test"
	"github.com/axwilliams/books-api/internal/test/mock"
	"github.com/gorilla/mux"
	"github.com/twinj/uuid"
)

var bookHandler handlers.BookHandler

func init() {
	mockBook := mock.NewMockBook()
	bookService := book.NewService(mockBook)
	bookHandler = handlers.NewBookHandler(bookService)
}

func TestFindAllBooks(t *testing.T) {
	r, err := http.NewRequest("GET", "/api/v1/books", nil)
	if err != nil {
		t.Errorf("\t%s\tRequest failed: %v\n", test.Failed, err)
	}

	rr := httptest.NewRecorder()
	h := http.HandlerFunc(bookHandler.FindAll)
	h.ServeHTTP(rr, r)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("\t%s\tWrong status code: want %v got %v", test.Failed, http.StatusOK, status)
	}
	t.Logf("\t%s\tStatus code correct: 200", test.Success)

	res := rr.Body.String()
	expected := `[{"id":"f4ac7e14-fc8e-4096-b956-34e5a33040f2","isbn":"978-0241372579","title":"The Castle","author":"Franz Kafka","category":"Fiction"},{"id":"71432eb9-58da-4eae-aa20-ccc49064246f","isbn":"978-1451673319","title":"Fahrenheit 451","author":"Ray Bradbury","category":"Fiction"},{"id":"562e1fe0-0dde-4717-a008-cd2a699301d2","isbn":"978-0465025275","title":"Six Easy Pieces","author":"Richard Feynman","category":"Science"}]`
	if res != expected {
		t.Fatalf("\t%s\tWrong response: want %v got %v", test.Failed, expected, res)
	}
	t.Logf("\t%s\tResults returned", test.Success)
}

func TestFindBookById(t *testing.T) {
	samples := []struct {
		id         string
		statusCode int
		expected   string
	}{
		// Invalid ID
		{
			id:         "f4ac7e14-fc8e-4096-b956-34e5a",
			statusCode: http.StatusBadRequest,
			expected:   `{"message":"` + book.ErrInvalidID.Error() + `"}`,
		},
		// Not found
		{
			id:         "f4ac7e14-fc8e-4096-b956-34e5a33040f5",
			statusCode: http.StatusOK,
			expected:   "",
		},
		// Success
		{
			id:         "f4ac7e14-fc8e-4096-b956-34e5a33040f2",
			statusCode: http.StatusOK,
			expected:   `{"id":"f4ac7e14-fc8e-4096-b956-34e5a33040f2","isbn":"978-0241372579","title":"The Castle","author":"Franz Kafka","category":"Fiction"}`,
		},
	}

	for _, sample := range samples {
		r, err := http.NewRequest("GET", "/api/v1/books", nil)
		if err != nil {
			t.Errorf("\t%s\tRequest failed: %v\n", test.Failed, err)
		}

		r = mux.SetURLVars(r, map[string]string{"id": sample.id})

		rr := httptest.NewRecorder()
		h := http.HandlerFunc(bookHandler.FindById)
		h.ServeHTTP(rr, r)

		if sample.statusCode != rr.Code {
			t.Fatalf("\t%s\tWrong status code: want %v got %v", test.Failed, sample.statusCode, rr.Code)
		}
		t.Logf("\t%s\tStatus code correct: %v", test.Success, rr.Code)

		res := rr.Body.String()
		if res != sample.expected {
			t.Fatalf("\t%s\tWrong response: want %v got %v", test.Failed, sample.expected, res)
		}
		t.Logf("\t%s\tResponse data correct", test.Success)
	}
}

func TestAddBook(t *testing.T) {
	samples := []struct {
		payload    string
		statusCode int
		expected   string
	}{
		// Missing body
		{
			payload:    ``,
			statusCode: http.StatusBadRequest,
			expected:   `{"message":"Unable to decode JSON: EOF"}`,
		},
		// Invalid JSON
		{
			payload:    `{"isbn": "978-0099448792","title"`,
			statusCode: http.StatusBadRequest,
			expected:   `{"message":"Unable to decode JSON: unexpected EOF"}`,
		},
		// Unknown field
		{
			payload:    `{"isbn": "978-0099448792","title": "The Wind-Up Bird Chronicle","author": "Haruki Murakami","cat": "Fiction"}`,
			statusCode: http.StatusBadRequest,
			expected:   `{"message":"Unable to decode JSON: json: unknown field \"cat\""}`,
		},
		// Missing isbn
		{
			payload:    `{"isbn": "","title": "The Wind-Up Bird Chronicle","author": "Haruki Murakami","category": "Fiction"}`,
			statusCode: http.StatusUnprocessableEntity,
			expected:   `{"message":"Validation failed","errors":["isbn is a required field"]}`,
		},
		// Missing title
		{
			payload:    `{"isbn": "978-0099448792","title": "","author": "Haruki Murakami","category": "Fiction"}`,
			statusCode: http.StatusUnprocessableEntity,
			expected:   `{"message":"Validation failed","errors":["title is a required field"]}`,
		},
		// Missing author
		{
			payload:    `{"isbn": "978-0099448792","title": "The Wind-Up Bird Chronicle","author": "","category": "Fiction"}`,
			statusCode: http.StatusUnprocessableEntity,
			expected:   `{"message":"Validation failed","errors":["author is a required field"]}`,
		},
		// Valid insert
		{
			payload:    `{"isbn":"978-0099448792","title":"The Wind-Up Bird Chronicle","author":"Haruki Murakami","category":"Fiction"}`,
			statusCode: http.StatusCreated,
			expected:   "",
		},
	}

	for _, sample := range samples {
		r, err := http.NewRequest("POST", "/api/v1/books", bytes.NewBufferString(sample.payload))
		if err != nil {
			t.Errorf("\t%s\tRequest failed: %v\n", test.Failed, err)
		}

		rr := httptest.NewRecorder()
		h := http.HandlerFunc(bookHandler.Add)
		h.ServeHTTP(rr, r)

		if sample.statusCode != rr.Code {
			t.Fatalf("\t%s\tWrong status code: want %v got %v", test.Failed, sample.statusCode, rr.Code)
		}
		t.Logf("\t%s\tStatus code correct: %v", test.Success, rr.Code)

		res := rr.Body.String()

		if sample.expected != "" && res != sample.expected {
			t.Fatalf("\t%s\tWrong response: want %v got %v", test.Failed, sample.expected, res)
		}
		t.Logf("\t%s\tResponse data correct", test.Success)

		if sample.statusCode == http.StatusCreated {
			var resp map[string]string
			err := json.NewDecoder(rr.Body).Decode(&resp)
			if err != nil {
				t.Errorf("\t%s\tFailed to decode JSON response: %v", test.Failed, err)
			}
			if _, err := uuid.Parse(resp["id"]); err != nil {
				t.Fatalf("\t%s\tResponse ID not a valid UUID", test.Failed)
			}
			t.Logf("\t%s\tResponse data correct", test.Success)
		}
	}
}

func TestEditBook(t *testing.T) {
	samples := []struct {
		id         string
		payload    string
		statusCode int
		expected   string
	}{
		// Invalid ID
		{
			id:         "f4ac7e14-fc8e-4096-b956-34e5a",
			payload:    `{"isbn": "978-0099448793","title": "Nineteen Eighty-Four","author": "George Orwell"}`,
			statusCode: http.StatusBadRequest,
			expected:   `{"message":"` + book.ErrInvalidID.Error() + `"}`,
		},
		// // Missing body
		{
			id:         "f4ac7e14-fc8e-4096-b956-34e5a33040f2",
			payload:    ``,
			statusCode: http.StatusBadRequest,
			expected:   `{"message":"Unable to decode JSON: EOF"}`,
		},
		// Invalid JSON
		{
			id:         "f4ac7e14-fc8e-4096-b956-34e5a33040f2",
			payload:    `{"isbn": "978-0099448793","title"`,
			statusCode: http.StatusBadRequest,
			expected:   `{"message":"Unable to decode JSON: unexpected EOF"}`,
		},
		// Unknown field
		{
			id:         "f4ac7e14-fc8e-4096-b956-34e5a33040f2",
			payload:    `{"isbn": "978-0099448793","title": "Nineteen Eighty-Four","author": "George Orwell","cat": "Fiction"}`,
			statusCode: http.StatusBadRequest,
			expected:   `{"message":"Unable to decode JSON: json: unknown field \"cat\""}`,
		},
		// Not found
		{
			id:         "7b6807c2-1e11-4e38-bdfd-281186885c3f",
			payload:    `{"isbn": "978-0099448793","title": "Nineteen Eighty-Four","author": "George Orwell"}`,
			statusCode: http.StatusGone,
			expected:   `{"message":"` + book.ErrNoAffect.Error() + `"}`,
		},
		// // Success
		{
			id:         "f4ac7e14-fc8e-4096-b956-34e5a33040f2",
			payload:    `{"isbn": "978-0099448793","title": "Nineteen Eighty-Four","author": "George Orwell", "category":"Fiction"}`,
			statusCode: http.StatusOK,
			expected:   "",
		},
	}

	for _, sample := range samples {
		r, err := http.NewRequest("PATCH", "/api/v1/books", bytes.NewBufferString(sample.payload))
		if err != nil {
			t.Errorf("\t%s\tRequest failed: %v\n", test.Failed, err)
		}

		r = mux.SetURLVars(r, map[string]string{"id": sample.id})

		rr := httptest.NewRecorder()
		h := http.HandlerFunc(bookHandler.Edit)
		h.ServeHTTP(rr, r)

		if sample.statusCode != rr.Code {
			t.Fatalf("\t%s\tWrong status code: want %v got %v", test.Failed, sample.statusCode, rr.Code)
		}
		t.Logf("\t%s\tStatus code correct: %v", test.Success, rr.Code)

		res := rr.Body.String()

		if sample.expected != "" && res != sample.expected {
			t.Fatalf("\t%s\tWrong response: want %v got %v", test.Failed, sample.expected, res)
		}
		t.Logf("\t%s\tResponse data correct", test.Success)
	}
}

func TestDeleteBook(t *testing.T) {
	samples := []struct {
		id         string
		statusCode int
		expected   string
	}{
		// Invalid ID
		{
			id:         "f4ac7e14-fc8e-4096-b956-34e5a",
			statusCode: http.StatusBadRequest,
			expected:   `{"message":"` + book.ErrInvalidID.Error() + `"}`,
		},
		// Not found
		{
			id:         "7b6807c2-1e11-4e38-bdfd-281186885c3f",
			statusCode: http.StatusGone,
			expected:   `{"message":"` + book.ErrNoAffect.Error() + `"}`,
		},
		// Success
		{
			id:         "f4ac7e14-fc8e-4096-b956-34e5a33040f2",
			statusCode: http.StatusOK,
			expected:   "",
		},
	}

	for _, sample := range samples {
		r, err := http.NewRequest("DELETE", "/api/v1/books", nil)
		if err != nil {
			t.Errorf("\t%s\tRequest failed: %v\n", test.Failed, err)
		}

		r = mux.SetURLVars(r, map[string]string{"id": sample.id})

		rr := httptest.NewRecorder()
		h := http.HandlerFunc(bookHandler.Delete)
		h.ServeHTTP(rr, r)

		if sample.statusCode != rr.Code {
			t.Fatalf("\t%s\tWrong status code: want %v got %v", test.Failed, sample.statusCode, rr.Code)
		}
		t.Logf("\t%s\tStatus code correct: %v", test.Success, rr.Code)

		res := rr.Body.String()
		if res != sample.expected {
			t.Fatalf("\t%s\tWrong response: want %v got %v", test.Failed, sample.expected, res)
		}
		t.Logf("\t%s\tResponse data correct", test.Success)
	}
}
