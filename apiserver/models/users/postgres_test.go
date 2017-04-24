package users

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

//TestPostgresStore tests the dockerized PGStore
func TestPostgresStore(t *testing.T) {
	//Preparing a Postgres data abstraction for later use
	psdb, err := sql.Open("postgres", "user=pgstest dbname=pgstest sslmode=disable")
	if err != nil {
		t.Errorf("error starting db: %v", err)
	}
	//Creates the store structure
	store := &PGStore{
		DB: psdb,
	}
	//Pings the DB-- establishes a connection to the db
	err = psdb.Ping()
	if err != nil {
		t.Errorf("error pinging db %v", err)
	}

	newUser := &NewUser{
		Email:        "test@test.com",
		Password:     "password",
		PasswordConf: "password",
		UserName:     "mrtester",
		FirstName:    "test",
		LastName:     "tester",
	}

	//clears previous test users in the DB
	_, err = psdb.Exec("DELETE FROM users WHERE Email = test@test.com")
	if err != nil {
		t.Errorf("could not delete table: %v\n", err)
	}
	//start of insert
	user, err := store.Insert(newUser)
	if err != nil {
		t.Errorf("error inserting user: %v\n", err)
	}
	//means that ToUser() probably was not implemented correctly
	if nil == user {
		t.Fatalf("Nil returned from store.Insert()\n")
	}
	//getting user from ID of previous inserted user
	user2, err := store.GetByID(user.ID)
	if err != nil {
		t.Errorf("error finding user by ID: %v\n", err)
	}
	//found something but didnt match
	if user.ID != user2.ID {
		t.Errorf("ID of user retrieved by ID does not match: expected %s but got %s\n", user.ID, user2.ID)
	}
	//getting any user with the given email
	user2, err = store.GetByEmail(newUser.Email)
	if err != nil {
		t.Errorf("error getting user by email: %v\n", err)
	}
	if user.ID != user2.ID {
		t.Errorf("ID of user retreived by Email does not match: expected %s but got %s\n", user.ID, user2.ID)
	}
	//any user with given username
	user2, err = store.GetByUserName(newUser.UserName)
	if err != nil {
		t.Errorf("error getting user by UserName: %v\n", err)
	}
	if user.ID != user2.ID {
		t.Errorf("ID of user retrieved by UserName does not match: expected %s but got %s\n", user.ID, user2.ID)
	}

	update := &UserUpdates{
		FirstName: "UPDATED Test",
		LastName:  "UPDATED Tester",
	}
	//updates the store with fields in update
	if err = store.Update(update, user); err != nil {
		t.Errorf("Error updating user: %v\n", err)
	}

	//reaquire the user -- by now user ought to have updated fields
	user, err = store.GetByID(user.ID)
	if err != nil {
		t.Errorf("error finding user by ID: %v\n", err)
	}

	if user.FirstName != update.FirstName {
		t.Errorf("FirstName field not updated: expected `%s` but got `%s`\n", update.FirstName, user.FirstName)
	}
	if user.LastName != update.LastName {
		t.Errorf("LastName field not updated: expected `%s` but got `%s`\n", update.LastName, user.LastName)
	}

	//gets all users in an array
	all, err := store.GetAll()
	if err != nil {
		t.Errorf("Error getting all users: %v\n", err)
	}
	if len(all) != 1 {
		t.Errorf("incorrect length of all users: expected %d but got %d\n", 1, len(all))
	}
	if all[0].ID != user.ID {
		t.Errorf("ID of user retrieved by all does not match: expected %s but got %s\n", user.ID, all[0].ID)
	}
}
