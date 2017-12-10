package model

import (
	"database"
	"entity"
	"sort"
	"testing"
	"testutil"

	_ "github.com/mattn/go-sqlite3"
)

// NoError ..
func ok(t *testing.T, ec ErrorCode) bool {
	if ec != OK {
		t.Errorf("Unexpected error: %d\n", ec)
		return false
	}
	return true
}

func TestCreateUser(t *testing.T) {
	database.WithTestDB(func() {
		users := []*entity.User{
			{"foo", "fooooo", "foo@"}, {"bar", "barrrr", "bar@"},
			{"baz", "bazzzz", "baz@"}}
		now := make([]*entity.User, 0)
		for _, u := range users {
			now = append(now, u)
			CreateUser(u)
			us, _ := database.GetAllUsers()
			sort.Sort(entity.UserSlice(now))
			sort.Sort(entity.UserSlice(us))
			testutil.ExpectDeepEq(t, us, now)
		}
		t.Run("Should Duplicate", func(t *testing.T) {
			ec := CreateUser(&entity.User{Username: "foo"})
			testutil.ExpectDeepEq(t, ec, DuplicateUser)
			us, _ := database.GetAllUsers()
			sort.Sort(entity.UserSlice(us))
			sort.Sort(entity.UserSlice(users))
			testutil.ExpectDeepEq(t, us, users)
		})
	})
}

func TestGetAllUsers(t *testing.T) {
	database.WithTestDB(func() {
		users := []*entity.User{
			{"foo", "fooooo", "foo@"}, {"bar", "barrrr", "bar@"},
			{"baz", "bazzzz", "baz@"}}
		sort.Sort(entity.UserSlice(users))
		for _, u := range users {
			database.StoreUser(u)
		}
		var nilus []*entity.User
		t.Run("Without Authentication", func(t *testing.T) {
			us, ec := GetAllUsers("foo", "noSuchToken")
			testutil.ExpectDeepEq(t, us, nilus)
			testutil.ExpectDeepEq(t, ec, InvalidToken)
		})
		t.Run("With Authentication", func(t *testing.T) {
			database.PutToken("foo", "nowwehavetoken")
			us, ec := GetAllUsers("foo", "nowwehavetoken")
			if ok(t, ec) {
				sort.Sort(entity.UserSlice(users))
				testutil.ExpectDeepEq(t, us, users)
			}
		})
	})
}

func TestRemoveUser(t *testing.T) {
	database.WithTestDB(func() {
		users := []*entity.User{
			{"foo", "fooooo", "foo@"}, {"bar", "barrrr", "bar@"},
			{"baz", "bazzzz", "baz@"}}
		sort.Sort(entity.UserSlice(users))
		for _, u := range users {
			database.StoreUser(u)
		}
		t.Run("Without Authentication", func(t *testing.T) {
			ec := RemoveUser("foo", "noSuchToken")
			testutil.ExpectDeepEq(t, ec, InvalidToken)
			us, _ := database.GetAllUsers()
			sort.Sort(entity.UserSlice(us))
			testutil.ExpectDeepEq(t, us, users)
		})
		t.Run("With Authentication", func(t *testing.T) {
			database.PutToken("foo", "nowwehavetoken")
			ec := RemoveUser("foo", "notthistoken")
			testutil.ExpectDeepEq(t, ec, InvalidToken)
			ec = RemoveUser("foo", "nowwehavetoken")
			if ok(t, ec) {
				sort.Sort(entity.UserSlice(users))
				userleft := []*entity.User{
					{"bar", "barrrr", "bar@"}, {"baz", "bazzzz", "baz@"}}
				us, _ := database.GetAllUsers()
				sort.Sort(entity.UserSlice(us))
				testutil.ExpectDeepEq(t, us, userleft)
			}
		})
	})
}

func TestLogin(t *testing.T) {
	database.WithTestDB(func() {
		users := []*entity.User{
			{"foo", "fooooo", "foo@"}, {"bar", "barrrr", "bar@"},
			{"baz", "bazzzz", "baz@"}}
		for _, u := range users {
			database.StoreUser(u)
		}
		t.Run("Authentication Fail", func(t *testing.T) {
			tok, ec := Login("foo", "wrongpassword")
			testutil.ExpectDeepEq(t, tok, "")
			testutil.ExpectDeepEq(t, ec, AuthenticationFail)
			tok, ec = Login("nosuchuser", "nosuchpassword")
			testutil.ExpectDeepEq(t, tok, "")
			testutil.ExpectDeepEq(t, ec, AuthenticationFail)
		})
		t.Run("Login Success", func(t *testing.T) {
			tok, ec := Login("bar", "barrrr")
			tokStored, err := database.GetToken("bar")
			if testutil.NoError(t, err) && ok(t, ec) {
				testutil.ExpectDeepEq(t, tok, tokStored)
			}
		})
		t.Run("Duplicate Login", func(t *testing.T) {
			tok, ec := Login("bar", "barrrr")
			testutil.ExpectDeepEq(t, tok, "")
			testutil.ExpectDeepEq(t, ec, DuplicateLogin)
		})
	})
}

func TestLogout(t *testing.T) {
	database.WithTestDB(func() {
		users := []*entity.User{
			{"foo", "fooooo", "foo@"}, {"bar", "barrrr", "bar@"},
			{"baz", "bazzzz", "baz@"}}
		for _, u := range users {
			database.StoreUser(u)
		}
		t.Run("Without Authentication", func(t *testing.T) {
			ec := Logout("blabla", "dontwork")
			testutil.ExpectDeepEq(t, ec, InvalidToken)
		})
		t.Run("Logout Success", func(t *testing.T) {
			tok, ec := Login("bar", "barrrr")
			if ok(t, ec) {
				ec = Logout("bar", tok)
				ok(t, ec)
			}
		})
	})
}
