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
	_ "github.com/lib/pq"
)

const (
	USR    = "users"
	SESS   = "sessions"
	SESSME = "sessions/me"
	USRME  = "users/me"
)

//calls all the tests
func TestCases(t *testing.T) {
	ctx := &Context{
		SessionKey:   "8675309",
		SessionStore: sessions.NewMemStore(time.Hour),
		UserStore:    users.NewMemStore(),
	}
	testUser(t, ctx)
	testSession(t, ctx)
	testSessionsMine(t, ctx)
	testUsersMe(t, ctx)

}

//creating a user and logging in
func testUser(t *testing.T, ctx *Context) {
	user := &users.NewUser{
		Email:        "test@test.com",
		Password:     "password",
		PasswordConf: "password",
		UserName:     "mrtester",
		FirstName:    "test",
		LastName:     "tester",
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
	req, err = http.NewRequest("GET", USR, bytes.NewBuffer(jsonUsr))

	handler.ServeHTTP(resRec, req)
	if resRec.Code == http.StatusInternalServerError {
		t.Errorf("handler returned internal error: %d ", resRec.Code)
	}

	if contentType != expectedContentType {
		t.Errorf("incorrect Content-Type response header: expected %s; got %s", expectedContentType, contentType)
	}
}

//tests the session-- logging in
func testSession(t *testing.T, ctx *Context) {
	creds := &users.Credentials{
		Email:    "test@test.com",
		Password: "password",
	}

	jsonCreds, err := json.Marshal(creds)
	if err != nil {
		t.Fatalf("error encoding test credentials")
	}

	handler := http.HandlerFunc(ctx.SessionsHandler)
	resRec := httptest.NewRecorder()
	req, err := http.NewRequest("POST", SESS, bytes.NewBuffer(jsonCreds))
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

//signing out
func testSessionsMine(t *testing.T, ctx *Context) {
	handler := http.HandlerFunc(ctx.SessionsMineHandler)
	resRec := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", SESSME, nil)
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(resRec, req)
	if resRec.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: expected `%d` but got `%d`\n", http.StatusOK, resRec.Code)
	}
	contentType := resRec.Header().Get("Content-Type")
	expectedContentType := "text/plain; charset=utf-8"
	if contentType != expectedContentType {
		t.Errorf("incorrect Content-Type response header: expected %s; got %s", expectedContentType, contentType)
	}
}

//getting user
func testUsersMe(t *testing.T, ctx *Context) {
	handler := http.HandlerFunc(ctx.UsersMeHandler)
	resRec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", USRME, nil)
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(resRec, req)
	if resRec.Code != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: expected `%d` but got `%d`\n", http.StatusInternalServerError, resRec.Code)
	}
	contentType := resRec.Header().Get("Content-Type")
	expectedContentType := "text/plain; charset=utf-8"
	if contentType != expectedContentType {
		t.Errorf("incorrect Content-Type response header: expected %s; got %s", expectedContentType, contentType)
	}

}
