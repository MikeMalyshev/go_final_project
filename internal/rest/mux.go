package rest

import (
	"bytes"
	"fmt"
	"net/http"

	"go_final_project/internal/app"
	"go_final_project/internal/authorization"
	"go_final_project/internal/config"
)

type Mux struct {
	cfg      *config.Handler
	app      *app.Application
	auth     *authorization.Handler
	serveMux *http.ServeMux
}

func NewMux(app *app.Application, cfg *config.Handler) *Mux {
	mux := &Mux{
		cfg:      cfg,
		app:      app,
		serveMux: http.NewServeMux(),
		auth:     authorization.Create(),
	}

	mux.serveMux.Handle("/", http.FileServer(http.Dir(cfg.WebDirPath())))
	mux.serveMux.HandleFunc("/api/nextdate", mux.NextDateHandler)
	mux.serveMux.HandleFunc("/api/task", mux.Auth(mux.TaskHandler))
	mux.serveMux.HandleFunc("/api/tasks", mux.Auth(mux.TasksHandler))
	mux.serveMux.HandleFunc("/api/task/done", mux.Auth(mux.TaskDoneHandler))
	mux.serveMux.HandleFunc("/api/signin", mux.SignupHandler) // ошибка в задании ...

	return mux
}

func (mux *Mux) ServeMux() *http.ServeMux {
	return mux.serveMux
}

func (mux Mux) Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля

		if mux.auth.PasswordSetted() {
			var jwt string // JWT-токен из куки
			// получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				jwt = cookie.Value
			}

			if !mux.auth.VerifyTocken(jwt) {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}

func (mux Mux) makeJsonResponse(jsonString string, resp http.ResponseWriter) {
	_, err := resp.Write(bytes.NewBufferString(jsonString).Bytes())
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
	}
}

func (mux Mux) makeErrorJsonResponse(error string, resp http.ResponseWriter) {
	mux.makeJsonResponse(fmt.Sprintf(`{"error":"%s"}`, error), resp)
}

func (mux Mux) makeEmptyJsonResponse(resp http.ResponseWriter) {
	mux.makeJsonResponse(fmt.Sprintf(`{}`), resp)
}
