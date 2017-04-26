package users

import (
	"database/sql"
	"fmt"
)

//PGStore store stucture
type PGStore struct {
	DB *sql.DB
}

//GetAll returns all users
func (ps *PGStore) GetAll() ([]*User, error) {
	var users []*User

	//Query the database to return multiple rows
	rows, err := ps.DB.Query(`SELECT ID, Email, FirstName, LastName, PassHash, PhotoURL, UserName FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	//Next refers to the first row initially
	//returns false once EOF
	for rows.Next() {
		var user = &User{}
		//scans values into User struct; error returned if scan unsuccessful
		if err := rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.PassHash, &user.PhotoURL, &user.UserName); err != nil {
			return nil, err
		}
		//adds to array
		users = append(users, user)
	}
	//error is returned if encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return users, nil
}

//GetByID returns the User with the given ID
func (ps *PGStore) GetByID(id UserID) (*User, error) {
	var user = &User{}
	//Queries and then scans; error returned if the scan unsuccessful
	err := ps.DB.QueryRow(`SELECT ID, Email, FirstName, LastName, PassHash, PhotoURL, UserName FROM users WHERE ID = $1`, id).Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.PassHash, &user.PhotoURL, &user.UserName)
	if err == sql.ErrNoRows || err != nil {
		return nil, err
	}
	return user, nil
}

//GetByEmail returns the User with the given email
func (ps *PGStore) GetByEmail(email string) (*User, error) {
	var user = &User{}
	err := ps.DB.QueryRow(`SELECT ID, Email, FirstName, LastName, PassHash, PhotoURL, UserName FROM users WHERE Email = $1`, email).Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.PassHash, &user.PhotoURL, &user.UserName)
	if err == sql.ErrNoRows || err != nil {
		return nil, err
	}
	return user, nil
}

//GetByUserName returns the User with the given user name
func (ps *PGStore) GetByUserName(name string) (*User, error) {
	var user = &User{}
	err := ps.DB.QueryRow(`SELECT ID, Email, FirstName, LastName, PassHash, PhotoURL, UserName FROM users WHERE UserName = $1`, name).Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.PassHash, &user.PhotoURL, &user.UserName)
	if err == sql.ErrNoRows || err != nil {
		return nil, err
	}
	return user, nil
}

//Insert inserts a new NewUser into the store
//and returns a User with a newly-assigned ID
func (ps *PGStore) Insert(newUser *NewUser) (*User, error) {
	u, err := newUser.ToUser()
	//Could not turn new user to user
	if err != nil {
		return nil, err
	}
	if nil == u {
		return nil, fmt.Errorf(".ToUser() returned nil")
	}

	//start a transaction
	tx, err := ps.DB.Begin()
	//err if transaction could not start
	if err != nil {
		return nil, err
	}
	sql := `INSERT INTO users (email, passhash, username, firstname, lastname, photourl, mobilephone) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	//Receives ONE row from the database
	row := tx.QueryRow(sql, u.Email, u.PassHash, u.UserName, u.FirstName, u.LastName, u.PhotoURL, u.MobilePhone)
	//scans the value of ID returned from query INTO the user
	err = row.Scan(&u.ID)
	//err if cant scan -- rollback transaction
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	//commits the transaction-- connection no longer reserved
	tx.Commit()
	return u, nil
}

//Update applies UserUpdates to the currentUser
func (ps *PGStore) Update(updates *UserUpdates, currentuser *User) error {
	//start transaction
	tx, err := ps.DB.Begin()
	if err != nil {
		return err
	}
	sql := `UPDATE users SET FirstName = $1, LastName = $2 WHERE id = $3`
	//executes the sql query
	_, err = tx.Exec(sql, updates.FirstName, updates.LastName, currentuser.ID)
	//err if could not exec, rollback transaction
	if err != nil {
		tx.Rollback()
		return err
	}
	//commits-- connection no longer reserved
	tx.Commit()
	return nil
}
