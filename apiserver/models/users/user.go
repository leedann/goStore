package users

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/mail"
	"os"

	"golang.org/x/crypto/bcrypt"
)

//gravatarBasePhotoURL is the base URL for Gravatar profile photos
const gravatarBasePhotoURL = "https://www.gravatar.com/avatar/"
const cost = 10

//UserID defines the type for user IDs
type UserID interface{}

//User represents a user account in the database
type User struct {
	ID          UserID `json:"id" bson:"_id"`
	Email       string `json:"email"`
	PassHash    []byte `json:"-" bson:"passHash"` //stored in mongo, but never encoded to clients
	UserName    string `json:"userName"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	PhotoURL    string `json:"photoURL"`
	MobilePhone string `json:"mobilePhone"`
}

//Credentials represents user sign-in credentials
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//NewUser represents a new user signing up for an account
type NewUser struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	PasswordConf string `json:"passwordConf"`
	UserName     string `json:"userName"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	MobilePhone  string `json:"mobilePhone"`
}

//UserUpdates represents updates one can make to a user
type UserUpdates struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

//Validate validates the new user
func (nu *NewUser) Validate() error {
	//ensure Email field is a valid Email
	//HINT: use mail.ParseAddress()
	//https://golang.org/pkg/net/mail/#ParseAddress

	_, err := mail.ParseAddress(nu.Email)
	if err != nil {
		return err
	}

	//ensure Password is at least 6 chars

	if len(nu.Password) < 6 {
		return fmt.Errorf("password should be at least 6 characters")
	}

	//ensure Password and PasswordConf match

	if nu.Password != nu.PasswordConf {
		return fmt.Errorf("passwords do not match")
	}

	//ensure UserName has non-zero length

	if len(nu.UserName) <= 0 {
		return fmt.Errorf("missing username")
	}

	//if you made here, it's valid, so return nil
	return nil
}

//ToUser converts the NewUser to a User
func (nu *NewUser) ToUser() (*User, error) {
	//build the Gravatar photo URL by creating an MD5
	//hash of the new user's email address, converting
	//that to a hex string, and appending it to their base URL:
	//https://www.gravatar.com/avatar/ + hex-encoded md5 has of email
	hash := md5.New()
	emailByte := []byte(nu.Email)
	hash.Write(emailByte)
	md5Email := hex.EncodeToString(hash.Sum(nil))

	gravURL := gravatarBasePhotoURL + md5Email

	//construct a new User setting the various fields
	//but don't assign a new ID here--do that in your
	//concrete Store.Insert() method

	usr := &User{}
	usr.PhotoURL = gravURL
	userSetting(usr, nu)
	//call the User's SetPassword() method to set the password,
	//which will hash the plaintext password
	usr.SetPassword(nu.Password)
	//return the User and nil
	return usr, nil
}

//sets the various user fields to equal new user fields
//does not export
func userSetting(u *User, nu *NewUser) {
	u.Email = nu.Email
	u.FirstName = nu.FirstName
	u.LastName = nu.LastName
	u.UserName = nu.UserName
	u.MobilePhone = nu.MobilePhone
}

//SetPassword hashes the password and stores it in the PassHash field
func (u *User) SetPassword(password string) error {
	//hash the plaintext password using an adaptive
	//crytographic hashing algorithm like bcrypt
	//https://godoc.org/golang.org/x/crypto/bcrypt

	//converting password to byte
	bytePass := []byte(password)
	passHash, err := bcrypt.GenerateFromPassword(bytePass, cost)
	if err != nil {
		fmt.Printf("error hashing password: %v", err)
		os.Exit(1)
	}
	//set the User's PassHash field to the resulting hash
	u.PassHash = passHash
	return nil
}

//Authenticate compares the plaintext password against the stored hash
//and returns an error if they don't match, or nil if they do
func (u *User) Authenticate(password string) error {
	//compare the plaintext password with the PassHash field
	//using the same hashing algorithm you used in SetPassword
	bytePass := []byte(password)
	return bcrypt.CompareHashAndPassword(u.PassHash, bytePass)
}
