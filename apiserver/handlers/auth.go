package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/info344-s17/challenges-leedann/apiserver/models/users"
	"github.com/info344-s17/challenges-leedann/apiserver/sessions"
)

const (
	charsetUTF8         = "charset=utf-8"
	contentTypeJSON     = "application/json"
	contentTypeJSONUTF8 = contentTypeJSON + "; " + charsetUTF8
)

//UserHandler allows users to sign up or gets all users
func (ctx *Context) UserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		decoder := json.NewDecoder(r.Body)
		newuser := &users.NewUser{}
		if err := decoder.Decode(newuser); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		err := newuser.Validate()
		if err != nil {
			http.Error(w, "User not valid", http.StatusBadRequest)
			return
		}
		_, err = ctx.UserStore.GetByEmail(newuser.Email)
		if err != nil {
			http.Error(w, "Email Already Exists", http.StatusBadRequest)
			return
		}
		_, err = ctx.UserStore.GetByUserName(newuser.UserName)
		if err != nil {
			http.Error(w, "Username Already Exists", http.StatusBadRequest)
			return
		}
		user, err := ctx.UserStore.Insert(newuser)
		state := &SessionState{}
		_, err = sessions.BeginSession(ctx.SessionKey, ctx.SessionStore, state, w)
		w.Header().Add("Content-Type", contentTypeJSONUTF8)
		encoder := json.NewEncoder(w)
		encoder.Encode(user)
	case "GET":
		users, err := ctx.UserStore.GetAll()
		if err != nil {
			http.Error(w, "Error fetching users", http.StatusInternalServerError)
		}
		w.Header().Add("Content-Type", contentTypeJSONUTF8)
		encoder := json.NewEncoder(w)
		encoder.Encode(users)
	}
}

//SessionsHandler allows existing users to sign in
func (ctx *Context) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		creds := &users.Credentials{}
		if err := decoder.Decode(creds); err != nil {
			http.Error(w, "Error in Credentials", http.StatusBadRequest)
			return
		}
		u, err := ctx.UserStore.GetByEmail(creds.Email)
		if err != nil {
			http.Error(w, "Email not found", http.StatusUnauthorized)
		}
		err = u.Authenticate(creds.Password)
		if err != nil {
			http.Error(w, "Error authenticating user", http.StatusUnauthorized)
		}
		state := &SessionState{}
		_, err = sessions.BeginSession(ctx.SessionKey, ctx.SessionStore, state, w)
		w.Header().Add("Content-Type", contentTypeJSONUTF8)
		encoder := json.NewEncoder(w)
		encoder.Encode(u)
	} else {
		http.Error(w, "Error with request", http.StatusBadRequest)
		return
	}
}

//SessionsMineHandler allows authenticated users to sign out
func (ctx *Context) SessionsMineHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		s := sessions.SessionID(ctx.SessionKey)
		err := ctx.SessionStore.Delete(s)
		if err != nil {
			http.Error(w, "Error ending session", http.StatusInternalServerError)
		}
		w.Header().Add("Content-Type", "text/plain")
		w.Write([]byte("User has been signed out"))
	} else {
		http.Error(w, "Error with request", http.StatusBadRequest)
		return
	}
}

//UsersMeHandler Get the session state
func (ctx *Context) UsersMeHandler(w http.ResponseWriter, r *http.Request) {
	state := &SessionState{}
	_, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, state)
	if err != nil {
		http.Error(w, "Error getting sessionID", http.StatusInternalServerError)
		return
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(state.User)
}
