package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash []byte `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

// encrypt the password
func (u *User) SetPassword(plaintextPassword string) error{
	hash,err:=bcrypt.GenerateFromPassword([]byte(plaintextPassword),12)
	if err!=nil{return err}
	u.PasswordHash = hash
	return nil

}
// a method to match password 
func (u *User) MatchesPassword(plaintextPassword string)(bool,error){
	err:=bcrypt.CompareHashAndPassword(u.PasswordHash,[]byte(plaintextPassword))
	if err!=nil{
		return false,nil
	}
	return true ,nil
}