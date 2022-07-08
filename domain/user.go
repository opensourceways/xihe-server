package domain

type User struct {
	Id          string
	Bio         Bio
	Email       Email
	Account     Account
	Nickname    Nickname
	AvatarId    AvatarId
	PhoneNumber PhoneNumber
}
