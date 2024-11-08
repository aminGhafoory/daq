package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	ctx "github.com/aminGhafoory/daq/context"
	"github.com/aminGhafoory/daq/models"
	"github.com/aminGhafoory/daq/views/signIn"
	"github.com/aminGhafoory/daq/views/signUp"
)

type Users struct {
	UserService    *models.UserService
	SessionService *models.SessionService
}

func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	c := signUp.NewUser("SignIn", r, []string{})
	c.Render(context.Background(), w)
}

func (u *Users) Create(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var data struct {
		Email    string
		Password string
	}

	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")

	fmt.Println(data)

	user, err := u.UserService.CreateUser(data.Email, data.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("err create user : %v", err), http.StatusInternalServerError)
		return
	}

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)

}

func (u *Users) NewSignInPage(w http.ResponseWriter, r *http.Request) {
	user := ctx.User(r.Context())
	if user != nil {
		http.Redirect(w, r, "/users/me", http.StatusFound)
	}
	c := signIn.SignIn("signIn", r, []string{})
	c.Render(context.Background(), w)
}

func (u *Users) ProccessSignIn(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	user, err := u.UserService.Auth(email, password)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	fmt.Println(session)
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)

}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {

	user := ctx.User(r.Context())
	w.Write([]byte(fmt.Sprintf("UserID: %s UserEmail: %s PasswordHash: %s ",
		user.ID, user.Email, user.PasswordHash)))

}

func (u Users) SignOutUser(w http.ResponseWriter, r *http.Request) {
	tokenHash, err := readCookie(r, CookieSession)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
	}
	err = u.SessionService.DB.DeleteUserSession(context.Background(), tokenHash)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
	}
	deleteCookie(w, CookieSession)
	http.Redirect(w, r, "/", http.StatusFound)

}

type UserMiddleware struct {
	SessionService *models.SessionService
}

func (userMw UserMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token, err := readCookie(r, CookieSession)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		user, err := userMw.SessionService.User(token)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		contx := r.Context()

		contx = ctx.WithUser(contx, user)
		r = r.WithContext(contx)
		next.ServeHTTP(w, r)

	})

}

func (UserMw UserMiddleware) ReqireUser(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := ctx.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/sign-in", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})

}
