package domain

// Account
type Account interface {
	Account() string
}

func NewAccount(v string) (Account, error) {
	// TODO: format account

	return dpAccount(v), nil
}

type dpAccount string

func (r dpAccount) Account() string {
	return string(r)
}

// Nickname
type Nickname interface {
	Nickname() string
}

func NewNickname(v string) (Nickname, error) {
	// TODO: format nickname

	return dpNickname(v), nil
}

type dpNickname string

func (r dpNickname) Nickname() string {
	return string(r)
}

// Bio
type Bio interface {
	Bio() string
}

func NewBio(v string) (Bio, error) {
	// TODO: limited length for bio

	return dpBio(v), nil
}

type dpBio string

func (r dpBio) Bio() string {
	return string(r)
}

// Email
type Email interface {
	Email() string
}

func NewEmail(v string) (Email, error) {
	// TODO: check format

	return dpEmail(v), nil
}

type dpEmail string

func (r dpEmail) Email() string {
	return string(r)
}

// Phone Number
type PhoneNumber interface {
	PhoneNumber() string
}

func NewPhoneNumber(v string) (PhoneNumber, error) {
	// TODO: check format

	return dpPhoneNumber(v), nil
}

type dpPhoneNumber string

func (r dpPhoneNumber) PhoneNumber() string {
	return string(r)
}

// AvatarId
type AvatarId interface {
	AvatarId() string
}

func NewAvatarId(v string) (AvatarId, error) {
	// TODO: check the range of v

	return dpAvatarId(v), nil
}

type dpAvatarId string

func (r dpAvatarId) AvatarId() string {
	return string(r)
}
