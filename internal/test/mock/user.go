package mock

import (
	"net/http"

	"github.com/axwilliams/books-api/internal/business/user"
	"github.com/axwilliams/books-api/internal/platform/web"
	"github.com/lib/pq"
)

type MockUser interface {
	GetById(id string) (*user.User, error)
	Create(u *user.User) error
	Update(u *user.User) error
	Destroy(id string) error
	GetByUsername(username string) (*user.User, error)
	UsernameAvailable(username, cuurentID string) bool
	EmailAvailable(username, cuurentID string) bool
}

type mockUser struct{}

func NewMockUser() MockUser {
	return &mockUser{}
}

func (mu *mockUser) GetById(id string) (*user.User, error) {
	if id == "69a47775-6d89-4d38-ad38-acdb2928f6a1" {
		return &user.User{
			ID:           "69a47775-6d89-4d38-ad38-acdb2928f6a1",
			Username:     "author",
			Email:        "author@example.com",
			Roles:        pq.StringArray([]string{"AUTHOR"}),
			PasswordHash: []byte("$2a$10$ExnMCA7MuOwW.s8Ss0BvSuGNCHawMIpqMmyJ4Oa9sTCTKKw2x445e"),
		}, nil
	}

	return nil, user.ErrNoUserFound
}

func (mu *mockUser) Create(u *user.User) error {
	return nil
}

func (mu *mockUser) Update(u *user.User) error {
	return nil
}

func (mu *mockUser) Destroy(id string) error {
	if id == "bad069ce-4afa-4a53-a673-14ae7b627d06" {
		return nil
	}

	return web.NewRequestError(user.ErrNoAffect, http.StatusGone)
}

func (mu *mockUser) GetByUsername(username string) (*user.User, error) {
	if username == "author" {
		return &user.User{
			ID:           "69a47775-6d89-4d38-ad38-acdb2928f6a1",
			Username:     "author",
			Email:        "author@example.com",
			Roles:        pq.StringArray([]string{"AUTHOR"}),
			PasswordHash: []byte("$2a$10$ExnMCA7MuOwW.s8Ss0BvSuGNCHawMIpqMmyJ4Oa9sTCTKKw2x445e"),
		}, nil
	}

	return nil, user.ErrNoUserFound
}

func (mu *mockUser) UsernameAvailable(username, currentID string) bool {
	if username != "author" {
		return true
	}

	return false
}

func (mu *mockUser) EmailAvailable(email, currentID string) bool {
	if email != "author@example.com" {
		return true
	}

	return false
}
