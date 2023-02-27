package repositoryimpl

const (
	fieldAccount = "account"
)

type DUserRegInfo struct {
	Account  string            `bson:"account" json:"account,omitempty"`
	Name     string            `bson:"name" json:"name,omitempty"`
	City     string            `bson:"city" json:"city,omitempty"`
	Email    string            `bson:"email" json:"email,omitempty"`
	Phone    string            `bson:"phone" json:"phone,omitempty"`
	Identity string            `bson:"identity" json:"identity,omitempty"`
	Province string            `bson:"province" json:"province,omitempty"`
	Detail   map[string]string `bson:"detail" json:"detail,omitempty"`
}
