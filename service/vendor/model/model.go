package model

import (
	"database"
	"entity"
	"log"

	"github.com/dchest/uniuri"
)

type needLogin struct {
	token string
}

func authenticate(token string) (*needLogin, ErrorCode) {
	b, err := database.HasToken(token)
	if err != nil {
		return nil, DatabaseFail
	}
	if b {
		return &needLogin{token}, OK
	}
	return nil, InvalidToken
}

// CreateUser ..
func CreateUser(u *entity.User) ErrorCode {
	uu, err := database.GetUser(u.Username)
	if err != nil {
		return DatabaseFail
	}
	if uu != nil {
		log.Printf(
			"Duplicate: (%s,%s,%s,%s)",
			uu.Username, uu.Password, uu.Email, uu.Phone)
		return DuplicateUser
	}
	err = database.StoreUser(u)
	if err != nil {
		return DatabaseFail
	}
	return OK
}

func getAllUsers(needLogin) ([]*entity.User, ErrorCode) {
	us, e := database.GetAllUsers()
	if e != nil {
		return nil, DatabaseFail
	}
	return us, OK
}

// GetAllUsers ..
func GetAllUsers(token string) ([]*entity.User, ErrorCode) {
	lp, ec := authenticate(token)
	if ec != OK {
		return nil, ec
	}
	return getAllUsers(*lp)
}

func removeUser(username string, l needLogin) ErrorCode {
	u, err := database.GetUsername(l.token)
	if err != nil {
		return DatabaseFail
	}
	if u != username {
		return InvalidToken
	}
	logout(l)
	err = database.RemoveUser(u)
	if err != nil {
		return DatabaseFail
	}
	return OK
}

// RemoveUser ..
func RemoveUser(username, token string) ErrorCode {
	lp, ec := authenticate(token)
	if ec != OK {
		return ec
	}
	return removeUser(username, *lp)
}

// Login ..
func Login(username, password string) (string, ErrorCode) {
	u, err := database.GetUser(username)
	if err != nil {
		return "", DatabaseFail
	}
	if u == nil || u.Password != password {
		return "", AuthenticationFail
	}
	t := uniuri.New()
	err = database.PutToken(username, t)
	if err != nil {
		return "", DatabaseFail
	}
	return t, OK
}

// Logout ..
func Logout(token string) ErrorCode {
	lp, ec := authenticate(token)
	if ec != OK {
		return ec
	}
	return logout(*lp)
}

func logout(l needLogin) ErrorCode {
	err := database.DeleteToken(l.token)
	if err != nil {
		return DatabaseFail
	}
	return OK
}

// GetUser ..
func GetUser(username, token string) (*entity.User, ErrorCode) {
	_, ec := authenticate(token)
	if ec != OK {
		return nil, ec
	}
	u, err := database.GetUser(username)
	if err != nil {
		return nil, DatabaseFail
	}
	return u, OK
}
