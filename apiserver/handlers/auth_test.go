package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bytes"

	"github.com/info344-s17/challenges-leedann/apiserver/models/users"
	"github.com/info344-s17/challenges-leedann/apiserver/sessions"
)

const (
	USR    = "/users"
	SESS   = "/sessions"
	SESSME = "/sessions/me"
	USRME  = "/users/me"
)

func TestUser(t *testing.T) {
	user := &users.NewUser{
		Email:        "test@test.com",
		Password:     "password",
		PasswordConf: "password",
		UserName:     "mrtester",
		FirstName:    "test",
		LastName:     "tester",
	}
	ctx := &Context{
		SessionKey:   "8675309",
		SessionStore: sessions.NewMemStore(time.Hour),
		UserStore:    users.NewMemStore(),
	}
	jsonUsr, _ := json.Marshal(user)
	handler := http.HandlerFunc(ctx.UserHandler)
	resRec := httptest.NewRecorder()
	req, err := http.NewRequest("POST", USR, bytes.NewBuffer(jsonUsr))
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(resRec, req)
	if resRec.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: expected `%d` but got `%d`\n", http.StatusOK, resRec.Code)
	}

	contentType := resRec.Header().Get("Content-Type")
	expectedContentType := "application/json; charset=utf-8"
	if contentType != expectedContentType {
		t.Errorf("incorrect Content-Type response header: expected %s; got %s", expectedContentType, contentType)
	}
}
