package utils

import (
	"model"
)

func GetToken() string {
	t := model.LoadLoginFile()
	if t != nil {
		return t.Token
	} else {
		return ""
	}
}

func SaveToken(token string) {
	model.WriteLoginFile(token)
}

func LoginCheck() bool {
	var t = GetToken()
	if len(t) != 0 {
		return true
	} else {
		return false
	}
}
