package auth_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/axwilliams/book-api/internal/platform/auth"
	"github.com/axwilliams/book-api/internal/test"
)

func TestClaims(t *testing.T) {
	claim := auth.Claims{
		UserID: "a72bec75-0a5f-49af-a844-5763d188788e",
		Roles:  []string{auth.RoleAdmin},
	}

	token, err := auth.CreateToken(claim)
	if err != nil {
		t.Fatal(err)
	}

	parsedClaim, err := auth.ParseWithClaims(token)
	if err != nil {
		t.Fatal(err)
	}

	if claim.UserID != parsedClaim.UserID {
		t.Fatalf("\t%s\tError parsing claim UserID: want %v got %v", test.Failed, claim.Roles, parsedClaim.Roles)
	}
	t.Logf("\t%s\tClaim UserID parsed", test.Success)

	if ok := reflect.DeepEqual(claim.Roles, parsedClaim.Roles); !ok {
		t.Fatalf("\t%s\tError parsing claim Roles: want %v got %v", test.Failed, claim.Roles, parsedClaim.Roles)
	}
	t.Logf("\t%s\tClaim Roles parsed", test.Success)
}

func TestUserContext(t *testing.T) {
	samples := []struct {
		claim    auth.Claims
		expected string
	}{
		{
			claim: auth.Claims{
				UserID: "c86057bd-f135-4f50-a233-41dc3963093b",
				Roles:  []string{auth.RoleAuthor, auth.RoleAuthor},
			},
			expected: ``,
		},
		{
			claim: auth.Claims{
				UserID: "a72bec75-0a5f-49af-a844-5763d188788e",
				Roles:  []string{auth.RoleAdmin},
			},
			expected: ``,
		},
		{
			claim: auth.Claims{
				UserID: "69a47775-6d89-4d38-ad38-acdb2928f6a1",
				Roles:  []string{auth.RoleAuthor},
			},
			expected: ``,
		},
		{
			claim: auth.Claims{
				UserID: "bad069ce-4afa-4a53-a673-14ae7b627d06",
				Roles:  []string{},
			},
			expected: ``,
		},
	}

	for _, sample := range samples {
		ctx := auth.ContextWithUser(context.Background(), sample.claim)

		UserID, ok := auth.UserFromContext(ctx)
		if !ok {
			t.Fatalf("\t%s\tFailed to get user from context", test.Failed)
		}

		if sample.claim.UserID != UserID {
			t.Fatalf("\t%s\tWrong UserID from context: want %v got %v", test.Failed, sample.claim.UserID, UserID)
		}
		t.Logf("\t%s\tUserID returned from context", test.Success)

		roles, ok := auth.RolesFromContext(ctx)
		if !ok {
			t.Fatalf("\t%s\tFailed to get Roles from context", test.Failed)
		}

		if ok := reflect.DeepEqual(sample.claim.Roles, roles); !ok {
			t.Fatalf("\t%s\tWrong roles returned from contect: want %v got %v", test.Failed, sample.claim.Roles, roles)
		}
		t.Logf("\t%s\tRoles returned from context", test.Success)
	}
}

func TestHasRole(t *testing.T) {
	samples := []struct {
		has      []string
		wanted   string
		expected bool
	}{
		{
			has:      []string{auth.RoleAdmin, auth.RoleAuthor},
			wanted:   auth.RoleAdmin,
			expected: true,
		},
		{
			has:      []string{auth.RoleAuthor},
			wanted:   auth.RoleAuthor,
			expected: true,
		},
		{
			has:      []string{auth.RoleAuthor},
			wanted:   auth.RoleAdmin,
			expected: false,
		},
		{
			has:      []string{},
			wanted:   auth.RoleAdmin,
			expected: false,
		},
	}

	for _, sample := range samples {
		res := auth.HasRole(sample.has, sample.wanted)

		if sample.expected != res {
			t.Fatalf("\t%s\tWrong result returned from role check: want %v got %v", test.Failed, sample.expected, res)
		}
		t.Logf("\t%s\tRole check sucess", test.Success)
	}
}
