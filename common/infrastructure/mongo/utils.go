package Mongo

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	MongoCmdAll         = "$all"
	MongoCmdSet         = "$set"
	MongoCmdInc         = "$inc"
	MongoCmdPush        = "$push"
	MongoCmdPull        = "$pull"
	MongoCmdMatch       = "$match"
	MongoCmdFilter      = "$filter"
	MongoCmdProject     = "$project"
	MongoCmdAddToSet    = "$addToSet"
	MongoCmdElemMatch   = "$elemMatch"
	MongoCmdSetOnInsert = "$setOnInsert"

	fieldName  = "name"
	fieldEmail = "email"
)

func ObjectIdFilter(s string) (bson.M, error) {
	v, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		return nil, err
	}

	return bson.M{
		"_id": v,
	}, nil
}

func UserDocFilter(account, email string) bson.M {
	return bson.M{
		"$or": bson.A{
			bson.M{
				fieldName: account,
			},
			bson.M{
				fieldEmail: email,
			},
		},
	}
}

func UserDocFilterByAccount(account string) bson.M {
	return bson.M{
		fieldName: account,
	}
}

func GenDoc(doc interface{}) (m bson.M, err error) {
	v, err := json.Marshal(doc)
	if err != nil {
		return
	}

	if err = json.Unmarshal(v, &m); err != nil {
		return
	}

	return
}
