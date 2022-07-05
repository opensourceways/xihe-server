package entity

import (
	"fmt"

	"github.com/Authing/authing-go-sdk/lib/model"
)

const (
	JwtString = "xihe-sadf43@98524"
)

type TokenItem struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	IdToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

type AuthingKey struct {
	E   string `json:"e"`
	N   string `json:"n"`
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	Kid string `json:"Kid"`
}

type User struct {
	Username          *string `form:"username,omitempty"`
	Nickname          *string `form:"nickname,omitempty"`
	Photo             *string `form:"photo,omitempty"`
	Company           *string `form:"company,omitempty"`
	Name              *string `form:"name,omitempty"`
	GivenName         *string `form:"givenName,omitempty"`
	FamilyName        *string `form:"familyName,omitempty"`
	MiddleName        *string `form:"middleName,omitempty"`
	Profile           *string `form:"profile,omitempty"`
	PreferredUsername *string `form:"preferredUsername"`
	Website           *string `form:"website,omitempty"`
	Gender            *string `form:"gender,omitempty"`
	Birthdate         *string `form:"birthdate,omitempty"`
	Zoneinfo          *string `form:"zoneinfo,omitempty"`
	Locale            *string `form:"locale,omitempty"`
	Address           *string `form:"address,omitempty"`
	Formatted         *string `form:"formatted,omitempty"`
	StreetAddress     *string `form:"streetAddress,omitempty"`
	Locality          *string `form:"locality,omitempty"`
	Region            *string `form:"region,omitempty"`
	PostalCode        *string `form:"postalCode,omitempty"`
	City              *string `form:"city,omitempty"`
	Province          *string `form:"province,omitempty"`
	Country           *string `form:"country,omitempty"`
}

func (userupdateInput *User) ExportToAuthingData() (authingUserInput *model.UpdateUserInput, err error) {
	authingUserInput = new(model.UpdateUserInput)
	if userupdateInput == nil {
		err = fmt.Errorf("UpdateUserInput 不能为空")
		return

	}
	authingUserInput.Address = userupdateInput.Address
	authingUserInput.Username = userupdateInput.Username
	authingUserInput.Nickname = userupdateInput.Nickname
	authingUserInput.Photo = userupdateInput.Photo
	authingUserInput.Company = userupdateInput.Company
	authingUserInput.Name = userupdateInput.Name
	authingUserInput.GivenName = userupdateInput.GivenName
	authingUserInput.FamilyName = userupdateInput.FamilyName
	authingUserInput.MiddleName = userupdateInput.MiddleName
	authingUserInput.Profile = userupdateInput.Profile
	authingUserInput.PreferredUsername = userupdateInput.PreferredUsername
	authingUserInput.Website = userupdateInput.Website
	authingUserInput.Gender = userupdateInput.Gender
	authingUserInput.Birthdate = userupdateInput.Birthdate
	authingUserInput.Zoneinfo = userupdateInput.Zoneinfo
	authingUserInput.Locale = userupdateInput.Locale
	authingUserInput.Formatted = userupdateInput.Formatted
	authingUserInput.StreetAddress = userupdateInput.StreetAddress
	authingUserInput.Locality = userupdateInput.Locality
	authingUserInput.Region = userupdateInput.Region
	authingUserInput.PostalCode = userupdateInput.PostalCode
	authingUserInput.City = userupdateInput.City
	authingUserInput.Province = userupdateInput.Province
	authingUserInput.Country = userupdateInput.Country
	return
}
