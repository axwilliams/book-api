package user_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/axwilliams/book-api/internal/business/user"
	"github.com/axwilliams/book-api/internal/test"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var userService user.Service

func TestMain(m *testing.M) {
	db, container := test.Setup()

	userRepository := user.NewRepository(db)
	userService = user.NewService(userRepository)

	e := m.Run()

	test.Teardown(db, container)
	os.Exit(e)
}

func TestGetById(t *testing.T) {
	u := &user.User{
		ID:           "69a47775-6d89-4d38-ad38-acdb2928f6a1",
		Username:     "author",
		Email:        "author@example.com",
		Roles:        pq.StringArray([]string{"AUTHOR"}),
		PasswordHash: []byte("$2a$10$ExnMCA7MuOwW.s8Ss0BvSuGNCHawMIpqMmyJ4Oa9sTCTKKw2x445e"),
	}

	res, err := userService.GetById(u.ID)
	if err != nil {
		t.Fatal(err)
	}

	if ok := reflect.DeepEqual(res, u); !ok {
		t.Fatalf("\t%s\tError finding user: want %v got %v", test.Failed, u, res)
	}
	t.Logf("\t%s\tUser found", test.Success)
}

func TestCreate(t *testing.T) {
	nu := &user.NewUser{
		Username: "newadmin",
		Email:    "newadmin@example.com",
		Roles:    pq.StringArray([]string{"ADMIN", "AUTHOR"}),
		Password: "Admin#1",
	}

	u, err := userService.Create(nu)
	if err != nil {
		t.Fatal(err)
	}

	res, err := userService.GetById(u.ID)
	if err != nil {
		t.Fatal(err)
	}

	expected := &user.User{
		ID:           u.ID,
		Username:     "newadmin",
		Email:        "newadmin@example.com",
		Roles:        pq.StringArray([]string{"ADMIN", "AUTHOR"}),
		PasswordHash: []byte("$2a$10$IUs9j88n5g5pthZXNmU9tei2mhIX7MWTvk39AjWUx40juWOrrPOzi"),
	}

	if err := bcrypt.CompareHashAndPassword(res.PasswordHash, expected.PasswordHash); err != nil {
		res.PasswordHash = expected.PasswordHash
	}

	if ok := reflect.DeepEqual(expected, res); !ok {
		t.Fatalf("\t%s\tError creating user: want %v got %v", test.Failed, expected, res)
	}
	t.Logf("\t%s\tUser created", test.Success)
}

func TestUpdate(t *testing.T) {
	ID := "bad069ce-4afa-4a53-a673-14ae7b627d06"

	uu := user.UpdateUser{
		Username: "newauthor",
		Email:    "newauthor@example.com",
		Roles:    pq.StringArray([]string{"AUTHOR"}),
		Password: "Author#1",
	}

	if err := userService.Update(ID, uu); err != nil {
		t.Fatal(err)
	}

	res, err := userService.GetById(ID)
	if err != nil {
		t.Fatal(err)
	}

	expected := &user.User{
		ID:           ID,
		Username:     "newauthor",
		Email:        "newauthor@example.com",
		Roles:        pq.StringArray([]string{"AUTHOR"}),
		PasswordHash: []byte("$2a$10$ExnMCA7MuOwW.s8Ss0BvSuGNCHawMIpqMmyJ4Oa9sTCTKKw2x445e"),
	}

	if err := bcrypt.CompareHashAndPassword(res.PasswordHash, expected.PasswordHash); err != nil {
		res.PasswordHash = expected.PasswordHash
	}

	if ok := reflect.DeepEqual(expected, res); !ok {
		t.Fatalf("\t%s\tError updating user: want %v got %v", test.Failed, expected, res)
	}
	t.Logf("\t%s\tUser updated", test.Success)
}

func TestDestroy(t *testing.T) {
	ID := "bad069ce-4afa-4a53-a673-14ae7b627d06"

	if err := userService.Destroy(ID); err != nil {
		t.Fatal(err)
	}

	if _, err := userService.GetById(ID); err != user.ErrNoUserFound {
		t.Fatalf("\t%s\tError destroying user", test.Failed)
	}
	t.Logf("\t%s\tUser destroyed", test.Success)
}

func TestAuthenticate(t *testing.T) {
	claim, err := userService.Authenticate("admin", "Admin#1")
	if err != nil {
		t.Fatal(err)
	}

	expectedID := "a72bec75-0a5f-49af-a844-5763d188788e"
	if claim.UserID != expectedID {
		t.Fatalf("\t%s\tError authenticating claim UserID: want %v got %v", test.Failed, expectedID, claim.UserID)
	}

	expectedRoles := []string{"ADMIN"}
	if ok := reflect.DeepEqual(claim.Roles, expectedRoles); !ok {
		t.Fatalf("\t%s\tError authenticating claim Roles: want %v got %v", test.Failed, expectedRoles, claim.Roles)
	}

	t.Logf("\t%s\tClaim authenticated", test.Success)
}
