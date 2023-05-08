package transport

import (
	"fmt"
	"github.com/slavikx4/http-api-server/internal/database"
	"html/template"
	"log"
	"net/http"
	"time"
)

const (
	cookieName = "session_id"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h Handler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	sessionID, err := r.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {

			templ, err := template.ParseFiles("./web/html/login.html")
			if err != nil {
				log.Println(fmt.Sprintf("не удалось распарсить шаблон: %s", err.Error()))
			}

			if err = templ.Execute(w, nil); err != nil {
				log.Println(fmt.Sprintf("не удалось выполнить шаблон: %s", err.Error()))
			}

		} else {
			log.Println(fmt.Sprintf("не получилось прочитать cookie: %s", err.Error()))
		}

		return
	}

	session, ok := Sessions[sessionID.Value]
	if !ok {
		if _, err = w.Write([]byte("такой сессии нет")); err != nil {
			log.Println(err)
		}
	}

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			log.Println(err)
		}

		inputTask := r.Form["newTask"][0]
		if err := database.AddTaskInBase(session.Login, inputTask); err != nil {
			log.Println(err)
		}
	}

	data, err := database.UnloadTaskFromBase(session.Login)
	if err != nil {
		log.Println(err)
		return
	}

	templ, err := template.ParseFiles("./web/html/index.html")
	if err != nil {
		log.Println(fmt.Sprintf("не удалось распарсить шаблон: %s", err.Error()))
	}

	if err := templ.Execute(w, data); err != nil {
		log.Println(fmt.Sprintf("не удалось выполнить шаблон: %s", err.Error()))
	}
}

func (h Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
	}

	inputLogin := r.Form["login"][0]
	inputPassword := r.Form["password"][0]

	password, err := database.CheckUserInBase(inputLogin)
	if err != nil {
		if err := database.AddUserInBase(inputLogin, inputPassword); err != nil {
			log.Println(err)
		}
	} else {
		if inputPassword != password {
			if _, err := w.Write([]byte("неверный пароль")); err != nil {
				log.Println(err)
			}
			return
		}
	}

	sessionID := NewSession(inputLogin)
	var cookie = http.Cookie{
		Name:    cookieName,
		Value:   sessionID,
		Expires: time.Now().Add(1 * time.Minute),
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}
