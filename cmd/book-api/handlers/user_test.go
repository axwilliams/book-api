package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/axwilliams/book-api/cmd/book-api/handlers"
	"github.com/axwilliams/book-api/internal/business/user"
	"github.com/axwilliams/book-api/internal/test"
	"github.com/axwilliams/book-api/internal/test/mock"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/twinj/uuid"
)

var userHandler handlers.UserHandler

func init() {
	mockUser := mock.NewMockUser()
	userService := user.NewService(mockUser)
	userHandler = handlers.NewUserHandler(userService)
}

func TestAddUser(t *testing.T) {
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
			payload:    `{"username": "","email`,
			statusCode: http.StatusBadRequest,
			expected:   `{"message":"Unable to decode JSON: unexpected EOF"}`,
		},
		// Unknown field
		{
			payload:    `{"uname":"newauthor","email": "newauthor@example.com","password": "Author#2","roles": ["AUTHOR"]}`,
			statusCode: http.StatusBadRequest,
			expected:   `{"message":"Unable to decode JSON: json: unknown field \"uname\""}`,
		},
		// Empty username
		{
			payload:    `{"username":"", "email": "newauthor@example.com","password": "Author#2","roles": ["AUTHOR"]}`,
			statusCode: http.StatusUnprocessableEntity,
			expected:   `{"message":"Validation failed","errors":["username is a required field"]}`,
		},
		// Empty email
		{
			payload:    `{"username":"newauthor", "email": "","password": "Author#2","roles": ["AUTHOR"]}`,
			statusCode: http.StatusUnprocessableEntity,
			expected:   `{"message":"Validation failed","errors":["email is a required field"]}`,
		},
		// Empty password
		{
			payload:    `{"username":"newauthor", "email": "newauthor@example.com","password": "","roles": ["AUTHOR"]}`,
			statusCode: http.StatusUnprocessableEntity,
			expected:   `{"message":"Validation failed","errors":["password is a required field"]}`,
		},
		// Invalid email
		{
			payload:    `{"username":"newauthor", "email": "newauthor.com","password": "Author#2","roles": ["AUTHOR"]}`,
			statusCode: http.StatusUnprocessableEntity,
			expected:   `{"message":"Validation failed","errors":["email must be a valid email address"]}`,
		},
		// Invalid password
		{
			payload:    `{"username":"newauthor", "email": "newauthor@example.com","password": "Author","roles": ["AUTHOR"]}`,
			statusCode: http.StatusUnprocessableEntity,
			expected:   `{"message":"Validation failed","errors":["password must greater than 5 characters and contain a capital letter, lower case letter, number, and special character"]}`,
		},
		// Username exists
		{
			payload:    `{"username":"author", "email": "newauthor@example.com","password": "Author#2","roles": ["AUTHOR"]}`,
			statusCode: http.StatusNotAcceptable,
			expected:   `{"message":"` + user.ErrUsernameExists.Error() + `"}`,
		},
		// Email exists
		{
			payload:    `{"username":"newauthor", "email": "author@example.com","password": "Author#2","roles": ["AUTHOR"]}`,
			statusCode: http.StatusNotAcceptable,
			expected:   `{"message":"` + user.ErrEmailExists.Error() + `"}`,
		},
		// Success
		{
			payload:    `{"username":"newauthor", "email": "newauthor@example.com","password": "Author#2","roles": ["AUTHOR"]}`,
			statusCode: http.StatusCreated,
			expected:   "",
		},
	}

	for _, sample := range samples {
		r, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBufferString(sample.payload))
		if err != nil {
			t.Errorf("\t%s\tRequest failed: %v\n", test.Failed, err)
		}

		rr := httptest.NewRecorder()
		h := http.HandlerFunc(userHandler.Add)
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

func TestEditUser(t *testing.T) {
	samples := []struct {
		id         string
		payload    string
		statusCode int
		expected   string
	}{
		// // Invalid ID
		{
			id:         "69a47775-6d89-4d38-ad38-acdb2928",
			payload:    `{"username":"updatedauthor", "email": "updatedauthor@example.com","password": "Author#3", "roles": ["ADMIN", "AUTHOR"]}`,
			statusCode: http.StatusBadRequest,
			expected:   `{"message":"` + user.ErrInvalidID.Error() + `"}`,
		},
		// Missing body
		{
			id:         "69a47775-6d89-4d38-ad38-acdb2928f6a1",
			payload:    ``,
			statusCode: http.StatusBadRequest,
			expected:   `{"message":"Unable to decode JSON: EOF"}`,
		},
		// Invalid JSON
		{
			id:         "69a47775-6d89-4d38-ad38-acdb2928f6a1",
			payload:    `{"username":"updatedauthor", "email":`,
			statusCode: http.StatusBadRequest,
			expected:   `{"message":"Unable to decode JSON: unexpected EOF"}`,
		},
		// Unknown field
		{
			id:         "69a47775-6d89-4d38-ad38-acdb2928f6a1",
			payload:    `{"usname":"updatedauthor", "email": "updatedauthor@example.com","password": "Author#3", "roles": ["ADMIN", "AUTHOR"]}`,
			statusCode: http.StatusBadRequest,
			expected:   `{"message":"Unable to decode JSON: json: unknown field \"usname\""}`,
		},
		// Invalid email
		{
			id:         "69a47775-6d89-4d38-ad38-acdb2928f6a1",
			payload:    `{"username":"updatedauthor", "email": "updatedauthor.com","password": "Author#3", "roles": ["ADMIN", "AUTHOR"]}`,
			statusCode: http.StatusUnprocessableEntity,
			expected:   `{"message":"Validation failed","errors":["email must be a valid email address"]}`,
		},
		// Invalid password
		{
			id:         "69a47775-6d89-4d38-ad38-acdb2928f6a1",
			payload:    `{"username":"updatedauthor", "email": "updatedauthor@example.com","password": "Author", "roles": ["ADMIN", "AUTHOR"]}`,
			statusCode: http.StatusUnprocessableEntity,
			expected:   `{"message":"Validation failed","errors":["password must greater than 5 characters and contain a capital letter, lower case letter, number, and special character"]}`,
		},
		// Username exists
		{
			id:         "69a47775-6d89-4d38-ad38-acdb2928f6a1",
			payload:    `{"username":"author", "email": "updatedauthor@example.com","password": "Author#3", "roles": ["ADMIN", "AUTHOR"]}`,
			statusCode: http.StatusNotAcceptable,
			expected:   `{"message":"` + user.ErrUsernameExists.Error() + `"}`,
		},
		// Email exists
		{
			id:         "69a47775-6d89-4d38-ad38-acdb2928f6a1",
			payload:    `{"username":"updatedauthor", "email": "author@example.com","password": "Author#3", "roles": ["ADMIN", "AUTHOR"]}`,
			statusCode: http.StatusNotAcceptable,
			expected:   `{"message":"` + user.ErrEmailExists.Error() + `"}`,
		},
		// Not found
		{
			id:         "3defcc36-9a52-4274-8b72-47cd2d0b3e5c",
			payload:    `{"username":"updatedauthor", "email": "updatedauthor@example.com","password": "Author#3", "roles": ["ADMIN", "AUTHOR"]}`,
			statusCode: http.StatusGone,
			expected:   `{"message":"` + user.ErrNoAffect.Error() + `"}`,
		},
		// Success
		{
			id:         "69a47775-6d89-4d38-ad38-acdb2928f6a1",
			payload:    `{"username":"updatedauthor", "email": "updatedauthor@example.com","password": "Author#3", "roles": ["ADMIN", "AUTHOR"]}`,
			statusCode: http.StatusOK,
			expected:   "",
		},
	}

	for _, sample := range samples {
		r, err := http.NewRequest("PATCH", "/api/v1/users", bytes.NewBufferString(sample.payload))
		if err != nil {
			t.Errorf("\t%s\tRequest failed: %v\n", test.Failed, err)
		}

		r = mux.SetURLVars(r, map[string]string{"id": sample.id})

		rr := httptest.NewRecorder()
		h := http.HandlerFunc(userHandler.Edit)
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

func TestDeleteUser(t *testing.T) {
	samples := []struct {
		id         string
		statusCode int
		expected   string
	}{
		// Invalid ID
		{
			id:         "bad069ce-4afa-4a53-a673-14ae7b62",
			statusCode: http.StatusBadRequest,
			expected:   `{"message":"` + user.ErrInvalidID.Error() + `"}`,
		},
		// Not found
		{
			id:         "bad069ce-4afa-4a53-a673-14ae7b627d09",
			statusCode: http.StatusGone,
			expected:   `{"message":"` + user.ErrNoAffect.Error() + `"}`,
		},
		// Success
		{
			id:         "bad069ce-4afa-4a53-a673-14ae7b627d06",
			statusCode: http.StatusOK,
			expected:   "",
		},
	}

	for _, sample := range samples {
		r, err := http.NewRequest("DELETE", "/api/v1/users", nil)
		if err != nil {
			t.Errorf("\t%s\tRequest failed: %v\n", test.Failed, err)
		}

		r = mux.SetURLVars(r, map[string]string{"id": sample.id})

		rr := httptest.NewRecorder()
		h := http.HandlerFunc(userHandler.Delete)
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

func TestToken(t *testing.T) {
	samples := []struct {
		statusCode int
		username   string
		password   string
		expected   string
		basicAuth  bool
	}{
		// Invalid username
		{
			statusCode: http.StatusUnauthorized,
			username:   "badname",
			password:   "Author#1",
			expected:   `{"message":"` + user.ErrBadCredentials.Error() + `"}`,
			basicAuth:  true,
		},
		// Invalid password
		{
			statusCode: http.StatusUnauthorized,
			username:   "author",
			password:   "badpassword",
			expected:   `{"message":"` + user.ErrBadCredentials.Error() + `"}`,
			basicAuth:  true,
		},
		// Missing Basic Auth
		{
			statusCode: http.StatusUnauthorized,
			username:   "author",
			password:   "Author#1",
			expected:   `{"message":"` + user.ErrBasicAuth.Error() + `"}`,
			basicAuth:  false,
		},
		// Valid request
		{
			statusCode: http.StatusOK,
			username:   "author",
			password:   "Author#1",
			expected:   "",
			basicAuth:  true,
		},
	}

	for _, sample := range samples {
		r, err := http.NewRequest("POST", "http://localhost:8080/api/v1/users/token", nil)
		if err != nil {
			t.Errorf("\t%s\tRequest failed: %v\n", test.Failed, err)
		}

		if sample.basicAuth == true {
			r.SetBasicAuth(sample.username, sample.password)
		}

		rr := httptest.NewRecorder()
		h := http.HandlerFunc(userHandler.Token)
		h.ServeHTTP(rr, r)

		if sample.statusCode != rr.Code {
			t.Fatalf("\t%s\tWrong status code: want %v got %v", test.Failed, sample.statusCode, rr.Code)
		}
		t.Logf("\t%s\tStatus code correct: %v", test.Success, rr.Code)

		res := rr.Body.String()

		if sample.statusCode == http.StatusUnauthorized {
			if res != sample.expected {
				t.Fatalf("\t%s\tWrong response: want %v got %v", test.Failed, sample.expected, res)
			}
			t.Logf("\t%s\tResponse data correct", test.Success)
		}

		if sample.statusCode == http.StatusOK {
			var resp map[string]string
			err := json.NewDecoder(rr.Body).Decode(&resp)
			if err != nil {
				t.Errorf("\t%s\tfailed to decode JSON response: %v", test.Failed, err)
			}

			token, err := jwt.Parse(resp["token"], func(*jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("API_KEY")), nil
			})

			if err != nil || !token.Valid {
				t.Fatalf("\t%s\tWrong response: returned invalid token", test.Failed)
			}
			t.Logf("\t%s\tResponse data correct", test.Success)
		}
	}
}
