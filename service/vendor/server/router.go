package server

import (
	"encoding/json"
	"entity"
	"model"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
)

type response struct {
	StatusCode int         `json:"-"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

type theHandler func(*http.Request) *response

func responseJSON(w http.ResponseWriter, resp *response) {
	w.WriteHeader(resp.StatusCode)
	json.NewEncoder(w).Encode(resp)
}

func newResponse(code int, message string, data interface{}) *response {
	return &response{
		StatusCode: code,
		Message:    message,
		Data:       data,
	}
}

func notOK(code int, message string) *response {
	return newResponse(code, message, struct{}{})
}

func ok(data interface{}) *response {
	return newResponse(200, "OK", data)
}

func unauthorized() *response {
	return notOK(401, "Unauthorized")
}

func unauthorizedSp() *response {
	return newResponse(401, "Unauthorized", make([]*entity.User, 0))
}
func internalServerErrorSp() *response {
	return newResponse(500, "Internal Server Error", make([]*entity.User, 0))
}

func internalServerError() *response {
	return notOK(500, "Internal Server Error")
}

func checkParameter(
	w http.ResponseWriter, form url.Values, required ...string) bool {
	missing := make([]string, 0)
	for _, s := range required {
		_, present := form[s]
		if !present {
			missing = append(missing, s)
		}
	}
	if len(missing) > 0 {
		responseJSON(w, &response{
			StatusCode: http.StatusBadRequest,
			Message:    strings.Join(missing, ","),
			Data:       struct{}{},
		})
		return false
	}
	return true
}

func withCheckParameter(
	required ...string) func(theHandler) http.HandlerFunc {
	return func(h theHandler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			if checkParameter(w, r.Form, required...) {
				resp := h(r)
				if resp != nil {
					responseJSON(w, resp)
				} else {
					w.WriteHeader(http.StatusNoContent)
				}
			}
		}
	}
}

func loginHandler() http.HandlerFunc {
	return withCheckParameter("username", "password")(
		func(r *http.Request) *response {
			username, password := r.Form["username"][0], r.Form["password"][0]
			tok, ec := model.Login(username, password)
			switch ec {
			case model.OK:
				return ok(map[string]string{"token": tok})
			case model.AuthenticationFail:
				return notOK(403, "Forbidden")
			default:
				return internalServerError()
			}
		})
}

func logoutHandler() http.HandlerFunc {
	return withCheckParameter("token")(
		func(r *http.Request) *response {
			token := r.Form["token"][0]
			switch model.Logout(token) {
			case model.OK:
				return nil
			case model.InvalidToken:
				return unauthorized()
			default:
				return internalServerError()
			}
		})
}

func getUserHandler() http.HandlerFunc {
	return withCheckParameter("token")(
		func(r *http.Request) *response {
			username := mux.Vars(r)["username"]
			token := r.Form["token"][0]
			u, ec := model.GetUser(username, token)
			switch ec {
			case model.OK:
				if u == nil {
					return notOK(404, "Not Found")
				}
				return ok(u)
			case model.InvalidToken:
				return unauthorized()
			default:
				return internalServerError()
			}
		})
}

func deleteUserHandler() http.HandlerFunc {
	return withCheckParameter("token")(
		func(r *http.Request) *response {
			username := mux.Vars(r)["username"]
			token := r.Form["token"][0]
			switch model.RemoveUser(username, token) {
			case model.OK:
				return nil
			case model.InvalidToken:
				return unauthorized()
			default:
				return internalServerError()
			}
		})
}

func addUserHandler() http.HandlerFunc {
	return withCheckParameter("username", "password", "email", "phone")(
		func(r *http.Request) *response {
			username, password, email, phone :=
				r.Form["username"][0], r.Form["password"][0],
				r.Form["email"][0], r.Form["phone"][0]
			u := &entity.User{
				Username: username,
				Password: password,
				Email:    email,
				Phone:    phone,
			}
			switch model.CreateUser(u) {
			case model.OK:
				return newResponse(201, "Created", u)
			case model.DuplicateUser:
				return notOK(409, "Conflict")
			default:
				return internalServerError()
			}
		})
}

func listAllUserHandler() http.HandlerFunc {
	return withCheckParameter("token")(
		func(r *http.Request) *response {
			token := r.Form["token"][0]
			us, ec := model.GetAllUsers(token)
			switch ec {
			case model.OK:
				return ok(us)
			case model.InvalidToken:
				return unauthorizedSp()
			default:
				return internalServerErrorSp()
			}
		})
}

func router() http.Handler {
	router := mux.NewRouter()
	v1 := router.PathPrefix("/v1").Subrouter()
	v1.HandleFunc("/auth", loginHandler()).Methods("GET")
	v1.HandleFunc("/auth", logoutHandler()).Methods("DELETE")
	v1.HandleFunc("/users/{username}", getUserHandler()).Methods("GET")
	v1.HandleFunc("/users/{username}", deleteUserHandler()).Methods("DELETE")
	v1.HandleFunc("/users", addUserHandler()).Methods("POST")
	v1.HandleFunc("/users", listAllUserHandler()).Methods("GET")
	router.Handle("/", router.NotFoundHandler)
	return router
}
