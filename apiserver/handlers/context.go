package handlers

import (
	"github.com/info344-s17/challenges-leedann/apiserver/models/users"
	"github.com/info344-s17/challenges-leedann/apiserver/sessions"
)

//Context struct provides context to the session context
type Context struct {
	SessionKey   string
	SessionStore sessions.Store
	UserStore    users.Store
}
