package repositoryimpl

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	fieldName           = "name"
	fieldAccount        = "account"
	fieldBio            = "bio"
	fieldAvatarId       = "avatar_id"
	fieldVersion        = "version"
	fieldFollower       = "follower"
	fieldFollowing      = "following"
	fieldFollowerCount  = "follower_count"
	fieldFollowingCount = "following_count"
	fieldIsFollower     = "is_follower"
)

type DUser struct {
	Id primitive.ObjectID `bson:"_id"       json:"-"`

	Name                    string `bson:"name"       json:"name"`
	Email                   string `bson:"email"      json:"email"`
	Bio                     string `bson:"bio"        json:"bio"`
	AvatarId                string `bson:"avatar_id"  json:"avatar_id"`
	PlatformToken           string `bson:"token"      json:"token"`
	PlatformUserId          string `bson:"uid"        json:"uid"`
	PlatformUserNamespaceId string `bson:"nid"        json:"nid"`

	Follower  []string `bson:"follower"   json:"-"`
	Following []string `bson:"following"  json:"-"`

	// Version will be increased by 1 automatically.
	// So, don't marshal it to avoid setting it occasionally.
	Version int `bson:"version"    json:"-"`
}

type DUserRegInfo struct {
	Account  string            `bson:"account"        json:"account,omitempty"`
	Name     string            `bson:"name"           json:"name,omitempty"`
	City     string            `bson:"city"           json:"city,omitempty"`
	Email    string            `bson:"email"          json:"email,omitempty"`
	Phone    string            `bson:"phone"          json:"phone,omitempty"`
	Identity string            `bson:"identity"       json:"identity,omitempty"`
	Province string            `bson:"province"       json:"province,omitempty"`
	Detail   map[string]string `bson:"detail"         json:"detail,omitempty"`
	Version  int               `bson:"version"        json:"-"`
}
