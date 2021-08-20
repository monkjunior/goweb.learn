package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/monkjunior/goweb.learn/context"
	"github.com/monkjunior/goweb.learn/models"
	"github.com/monkjunior/goweb.learn/rand"
	"github.com/monkjunior/goweb.learn/views"
)

func NewUsers(us models.UserService) *Users {
	return &Users{
		NewView:   views.NewView("bootstrap", "users/new"),
		LoginView: views.NewView("bootstrap", "users/login"),
		us:        us,
	}
}

type Users struct {
	NewView   *views.View
	LoginView *views.View
	us        models.UserService
}

// New is used to render the form where a user can create
// a new user account
//
// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, r, nil)
}

type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// Create is used to process signup form when a user tries to
// create a new user account
//
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, r, vd)
		return
	}

	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	if err := u.us.Create(&user); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, r, vd)
		return
	}

	err := u.signIn(w, &user)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/galleries", http.StatusFound)
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// GetLogin is used to render the form where a user can login
// a registered account
//
// GET /login
func (u *Users) GetLogin(w http.ResponseWriter, r *http.Request) {
	u.LoginView.Render(w, r, nil)
}

// PostLogin is used to process login form when a user tries to
// login by a registered account
//
// POST /login
func (u *Users) PostLogin(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form LoginForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)
		return
	}

	user, err := u.us.Authenticate(form.Email, form.Password)

	if err != nil {
		switch err {
		case models.ErrNotFound:
			vd.AlertError("Invalid email address")
		default:
			vd.SetAlert(err)
		}
		u.LoginView.Render(w, r, vd)
		return
	}

	err = u.signIn(w, user)
	if err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)
	}

	http.Redirect(w, r, "/galleries", http.StatusFound)
}

// signIn is used to sign the given user in via cookie
func (u *Users) signIn(w http.ResponseWriter, user *models.User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}

		user.Remember = token
		err = u.us.Update(user)
		if err != nil {
			return nil
		}
	}

	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	return nil
}

// Logout is used to delete the user's cookies (remember_token)
// and then will update the user resource with a new remember_token
//
// POST /logout
func (u *Users) Logout(w http.ResponseWriter, r *http.Request) {
	// Invalidate the user's cookie
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	// Change remember_token
	user := context.User(r.Context())
	token, _ := rand.RememberToken()
	user.Remember = token
	_ = u.us.Update(user)
	http.Redirect(w, r, "/", http.StatusFound)
}

// CookieTest is used to display cookies set on the current user
func (u *Users) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("remember_token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user, err := u.us.ByRemember(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%+v\n", user)
}
