package cmd

import (
	"fmt"
	"model"
	"utils"
	"net/http"
	"net/url"
	"io/ioutil"
	"encoding/json"
	"bytes"
	"os"
)

type Resp struct {
	Message string
	Data    map[string]string
}

type UserList struct {
	Message string
	Data    []map[string]string
}

func printUser(user map[string]string) {
	b, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "\t")
	if err != nil {
		fmt.Println(err.Error())
	}
	out.WriteTo(os.Stdout)
	fmt.Println()
}

// Register a user
func Register(user, pass, mail, phone, baseURL string) int {
	// Send HTTP
	data := make(url.Values)
	data["username"] = []string{user}
	data["password"] = []string{pass}
	data["email"] = []string{mail}
	data["phone"] = []string{phone}
	resp, err := http.PostForm(baseURL + "/v1/users", data)

	if err != nil {
		utils.PrintError(err.Error())
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.PrintError(err.Error())
	}

	if resp.StatusCode == 201 {
		fmt.Println("Register Succeed")
		return 0
	}

	if resp.StatusCode == 409 {
		r := &Resp{}
		err = json.Unmarshal([]byte(string(body)), &r)
		if err != nil {
			utils.PrintError(err.Error())
		}
		fmt.Println(r.Message)
	}

	return 1
}

// ShowUsers print all users when logged in
func ShowUsers(baseURL string) int {
	if utils.LoginCheck() {
		// HTTP
		t := utils.GetToken()
		client := &http.Client{}
		req, _ := http.NewRequest("GET", baseURL+"/v1/users?token="+t, nil)

		resp, err := client.Do(req)
		if err != nil {
			utils.PrintError(err.Error())
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			utils.PrintError(err.Error())
		}

		r := &UserList{}
		err = json.Unmarshal([]byte(string(body)), &r)
		if err != nil {
			utils.PrintError(err.Error())
		}

		if resp.StatusCode == 200 {
			for i := 0; i < len(r.Data); i++ {
				printUser(r.Data[i])
			}
			return 0
		}

		if resp.StatusCode == 401 {
			fmt.Println(r.Message)
			return 1
		}

		return 1
	} else {
		fmt.Println("Please Login First")
		return 1
	}
}

// Login Command
func Login(user, pass, baseURL string) int {
	client := &http.Client{}
	req, err := http.NewRequest("GET", baseURL+"/v1/auth?username="+user+"&password="+pass, nil)
	resp, err := client.Do(req)

	if err != nil {
		utils.PrintError(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.PrintError(err.Error())
	}

	if resp.StatusCode == 200 {
		r := &Resp{}
		err = json.Unmarshal([]byte(string(body)), &r)
		if err != nil {
			utils.PrintError(err.Error())
		}

		utils.SaveToken(r.Data["token"])
		fmt.Println("Login Succeed")
		return 0
	}

	if resp.StatusCode == 403 {
		r := &Resp{}
		err = json.Unmarshal([]byte(string(body)), &r)
		if err != nil {
			utils.PrintError(err.Error())
		}

		fmt.Println(r.Message)
	}
	return 1
}

// Logout Command
func Logout(baseURL string) int {
	if utils.LoginCheck() {
		t := utils.GetToken()
		client := &http.Client{}
		req, _ := http.NewRequest("DELETE", baseURL+"/v1/auth?token="+t, nil)

		resp, err := client.Do(req)
		if err != nil {
			utils.PrintError(err.Error())
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			utils.PrintError(err.Error())
		}

		if resp.StatusCode == 204 {
			model.Logout()
			fmt.Println("Log Out Succeed!")
			return 0
		}

		if resp.StatusCode == 401 {
			r := &Resp{}
			err = json.Unmarshal([]byte(string(body)), &r)
			if err != nil {
				utils.PrintError(err.Error())
			}

			fmt.Println(r.Message)
		}

		return 1
	} else {
		fmt.Println("You haven't login.")
		return 1
	}
}

// DeleteUser delete current login user, and removed from its meeting
func DeleteUser(username, baseURL string) int {
	if utils.LoginCheck() {
		t := utils.GetToken()
		client := &http.Client{}
		req, _ := http.NewRequest("DELETE", baseURL+"/v1/users/" + username + "?token="+t, nil)

		resp, err := client.Do(req)
		if err != nil {
			utils.PrintError(err.Error())
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			utils.PrintError(err.Error())
		}

		if resp.StatusCode == 204 {
			fmt.Println("Delete User Succeed!")
			return 0
		}

		r := &Resp{}
		err = json.Unmarshal([]byte(string(body)), &r)
		if err != nil {
			utils.PrintError(err.Error())
		}

		if resp.StatusCode == 401 {
			fmt.Println(r.Message)
			return 1
		}

		return 0
	} else {
		fmt.Println("Please login first")
		return 1
	}
	return 0
}

func ShowInfo(username, baseURL string) int {
	if utils.LoginCheck() {
		t := utils.GetToken()
		client := http.Client{}
		req, _ := http.NewRequest("GET", baseURL+"/v1/users/" + username + "?token="+t, nil)

		resp, err := client.Do(req)

		if err != nil {
			utils.PrintError(err.Error())
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			utils.PrintError(err.Error())
		}

		r := &Resp{}
		err = json.Unmarshal([]byte(string(body)), &r)
		if err != nil {
			utils.PrintError(err.Error())
		}

		if resp.StatusCode == 200 {
			printUser(r.Data)
			return 0
		}

		if resp.StatusCode == 404 {
			fmt.Println(r.Message)
		}

		return 1
	} else {
	fmt.Println("Please Login first!")
	return 1
	}
}
