package domain

import "fmt"

type User struct {
	Id          string
	Bio         Bio
	Email       Email
	Account     Account
	Nickname    Nickname
	AvatarId    AvatarId
	PhoneNumber PhoneNumber
}

func (u User) ValidateID() error {
	if len(u.Id) == 0 {
		return fmt.Errorf("User id is inValidate")
	}
	return nil
}
