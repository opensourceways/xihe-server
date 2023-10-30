package mongodb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoCmdAll         = "$all"
	mongoCmdSet         = "$set"
	mongoCmdInc         = "$inc"
	mongoCmdPush        = "$push"
	mongoCmdPull        = "$pull"
	mongoCmdMatch       = "$match"
	mongoCmdFilter      = "$filter"
	mongoCmdProject     = "$project"
	mongoCmdAddToSet    = "$addToSet"
	mongoCmdElemMatch   = "$elemMatch"
	mongoCmdSetOnInsert = "$setOnInsert"
)

var (
	errDocExists    = errors.New("doc exists")
	errDocNotExists = errors.New("doc doesn't exist")
)

type dbError struct {
	error
}

func isDBError(err error) bool {
	_, b := err.(dbError)

	return b
}

func isDocNotExists(err error) bool {
	return errors.Is(err, errDocNotExists)
}

func isDocExists(err error) bool {
	return errors.Is(err, errDocExists)
}

func toUID(docId interface{}) (string, error) {
	if v, ok := docId.(primitive.ObjectID); ok {
		return v.Hex(), nil
	}

	return "", errors.New("retrieve id failed")
}

func newId() string {
	return primitive.NewObjectID().Hex()
}

func objectIdFilter(s string) (bson.M, error) {
	v, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		return nil, err
	}

	return bson.M{
		"_id": v,
	}, nil
}

func genDoc(doc interface{}) (m bson.M, err error) {
	v, err := json.Marshal(doc)
	if err != nil {
		return
	}

	if err = json.Unmarshal(v, &m); err != nil {
		return
	}

	return
}

func isErrNoDocuments(err error) bool {
	return err.Error() == mongo.ErrNoDocuments.Error()
}

func appendElemMatchToFilter(array string, exists bool, cond, filter bson.M) {
	match := bson.M{mongoCmdElemMatch: cond}

	if exists {
		filter[array] = match
	} else {
		filter[array] = bson.M{"$not": match}
	}
}

func (cli *client) newDocIfNotExist(
	ctx context.Context, collection string,
	filterOfDoc, docInfo bson.M,
) (string, error) {
	upsert := true
	r, err := cli.collection(collection).UpdateOne(
		ctx, filterOfDoc,
		bson.M{mongoCmdSetOnInsert: docInfo},
		&options.UpdateOptions{Upsert: &upsert},
	)
	if err != nil {
		return "", dbError{err}
	}

	if r.UpsertedID == nil {
		return "", errDocExists
	}

	return toUID(r.UpsertedID)
}

func (cli *client) replaceDoc(
	ctx context.Context, collection string,
	filterOfDoc, docInfo bson.M,
) (string, error) {
	upsert := true

	r, err := cli.collection(collection).ReplaceOne(
		ctx, filterOfDoc, docInfo,
		&options.ReplaceOptions{Upsert: &upsert},
	)
	if err != nil {
		return "", dbError{err}
	}

	if r.UpsertedID == nil {
		return "", nil
	}

	return toUID(r.UpsertedID)
}

func (cli *client) updateDoc(
	ctx context.Context, collection string,
	filterOfDoc, update bson.M, op string, version int,
) error {
	filterOfDoc[fieldVersion] = version
	r, err := cli.collection(collection).UpdateOne(
		ctx, filterOfDoc,
		bson.M{
			op:          update,
			mongoCmdInc: bson.M{fieldVersion: 1},
		},
	)

	if err != nil {
		return err
	}

	if r.MatchedCount == 0 {
		return errDocNotExists
	}

	return nil
}

func (cli *client) updateIncDoc(
	ctx context.Context, collection string,
	filterOfDoc, update bson.M, version int,
) error {
	filterOfDoc[fieldVersion] = version
	update[fieldVersion] = 1
	r, err := cli.collection(collection).UpdateOne(
		ctx, filterOfDoc,
		bson.M{
			mongoCmdInc: update,
		},
	)

	if err != nil {
		return err
	}

	if r.MatchedCount == 0 {
		return errDocNotExists
	}

	return nil
}

func (cli *client) getDoc(
	ctx context.Context, collection string,
	filterOfDoc, project bson.M, result interface{},
) error {
	var sr *mongo.SingleResult
	col := cli.collection(collection)
	if len(project) > 0 {
		sr = col.FindOne(ctx, filterOfDoc, &options.FindOneOptions{
			Projection: project,
		})
	} else {
		sr = col.FindOne(ctx, filterOfDoc)
	}

	if err := sr.Decode(result); err != nil {
		if isErrNoDocuments(err) {
			return errDocNotExists
		}

		return err
	}

	return nil
}

func (cli *client) getDocs(
	ctx context.Context, collection string,
	filterOfDoc, project bson.M, result interface{},
) error {
	col := cli.collection(collection)

	var cursor *mongo.Cursor
	var err error
	if len(project) > 0 {
		cursor, err = col.Find(ctx, filterOfDoc, &options.FindOptions{
			Projection: project,
		})
	} else {
		cursor, err = col.Find(ctx, filterOfDoc)
	}

	if err != nil {
		return dbError{err}
	}

	return cursor.All(ctx, result)
}

func (cli *client) addToSimpleArray(
	ctx context.Context, collection, array string,
	filterOfDoc, value interface{},
) error {
	r, err := cli.collection(collection).UpdateOne(
		ctx, filterOfDoc,
		bson.M{mongoCmdAddToSet: bson.M{array: value}},
	)
	if err != nil {
		return dbError{err}
	}

	if r.MatchedCount == 0 {
		return errDocNotExists
	}

	if r.ModifiedCount == 0 {
		return errDocExists
	}

	return nil
}

func (cli *client) removeFromSimpleArray(
	ctx context.Context,
	collection, array string,
	filterOfDoc, value interface{},
) error {
	r, err := cli.collection(collection).UpdateOne(
		ctx, filterOfDoc,
		bson.M{mongoCmdPull: bson.M{array: value}},
	)
	if err != nil {
		return dbError{err}
	}

	if r.MatchedCount == 0 {
		return errDocNotExists
	}

	return nil
}

func (cli *client) pushArrayElem(
	ctx context.Context,
	collection, array string,
	filterOfDoc, value bson.M,
) error {
	r, err := cli.collection(collection).UpdateOne(
		ctx, filterOfDoc,
		bson.M{mongoCmdPush: bson.M{array: value}},
	)
	if err != nil {
		return dbError{err}
	}

	if r.MatchedCount == 0 {
		return errDocNotExists
	}

	return nil
}

func (cli *client) pushElemToLimitedArray(
	ctx context.Context,
	collection, array string, keep int,
	filterOfDoc, value bson.M,
) error {
	r, err := cli.collection(collection).UpdateOne(
		ctx, filterOfDoc,
		bson.M{mongoCmdPush: bson.M{array: bson.M{
			"$each":  bson.A{value},
			"$slice": keep,
		}}},
	)
	if err != nil {
		return dbError{err}
	}

	if r.MatchedCount == 0 {
		return errDocNotExists
	}

	return nil
}

func (cli *client) pushElemToLimitedArrayWithVersion(
	ctx context.Context,
	collection, array string, keep int,
	filterOfDoc, value bson.M, version int, otherUpdate bson.M,
) error {
	filterOfDoc[fieldVersion] = version

	updates := bson.M{
		mongoCmdPush: bson.M{array: bson.M{
			"$each":  bson.A{value},
			"$slice": keep,
		}},
		mongoCmdInc: bson.M{fieldVersion: 1},
	}

	if len(otherUpdate) > 0 {
		updates[mongoCmdSet] = otherUpdate
	}

	r, err := cli.collection(collection).UpdateOne(
		ctx, filterOfDoc, updates,
	)
	if err != nil {
		return dbError{err}
	}

	if r.MatchedCount == 0 {
		return errDocNotExists
	}

	return nil
}

func (cli *client) pullNestedArrayElem(
	ctx context.Context, collection, array string,
	filterOfDoc, filterOfArray, data bson.M,
	version int, t int64,
) (bool, error) {
	return cli.modifyArrayElem(
		ctx, collection, array,
		filterOfDoc, filterOfArray, data,
		mongoCmdPull, version, t,
	)
}

func (cli *client) pushNestedArrayElem(
	ctx context.Context, collection, array string,
	filterOfDoc, filterOfArray, data bson.M,
	version int, t int64,
) (bool, error) {
	return cli.modifyArrayElem(
		ctx, collection, array,
		filterOfDoc, filterOfArray, data,
		mongoCmdPush, version, t,
	)
}

func (cli *client) updateArrayElem(
	ctx context.Context, collection, array string,
	filterOfDoc, filterOfArray, updateCmd bson.M,
	version int, t int64,
) (bool, error) {
	return cli.modifyArrayElem(
		ctx, collection, array,
		filterOfDoc, filterOfArray, updateCmd,
		mongoCmdSet, version, t,
	)
}

func (cli *client) modifyArrayElem(
	ctx context.Context, collection, array string,
	filterOfDoc, filterOfArray, updateCmd bson.M,
	op string, version int, t int64,
) (bool, error) {
	key := func(k string) string {
		return fmt.Sprintf("%s.$[i].%s", array, k)
	}

	cmd := bson.M{}
	for k, v := range updateCmd {
		cmd[key(k)] = v
	}

	arrayFilter := bson.M{}
	for k, v := range filterOfArray {
		arrayFilter["i."+k] = v
	}
	arrayFilter["i."+fieldVersion] = version

	updates := bson.M{
		mongoCmdInc: bson.M{key(fieldVersion): 1},
	}

	if op == mongoCmdSet {
		cmd[key(fieldUpdatedAt)] = t
	} else {
		updates[mongoCmdSet] = bson.M{key(fieldUpdatedAt): t}
	}

	updates[op] = cmd

	col := cli.collection(collection)
	r, err := col.UpdateOne(
		ctx, filterOfDoc, updates,
		&options.UpdateOptions{
			ArrayFilters: &options.ArrayFilters{
				Filters: bson.A{
					arrayFilter,
				},
			},
		},
	)
	if err != nil {
		return false, err
	}

	if r.MatchedCount == 0 {
		return false, errDocNotExists
	}

	return r.ModifiedCount > 0, nil
}

func (cli *client) pushNestedArrayElemAndUpdate(
	ctx context.Context, collection, array string,
	filterOfDoc, filterOfArray, updateCmd bson.M,
	version int, otherUpdate bson.M,
) (bool, error) {
	key := func(k string) string {
		return fmt.Sprintf("%s.$[i].%s", array, k)
	}

	cmd := bson.M{}
	for k, v := range updateCmd {
		cmd[key(k)] = v
	}

	updates := bson.M{
		mongoCmdInc:  bson.M{fieldVersion: 1},
		mongoCmdPush: cmd,
	}
	if len(otherUpdate) > 0 {
		updates[mongoCmdSet] = otherUpdate
	}

	arrayFilter := bson.M{}
	for k, v := range filterOfArray {
		arrayFilter["i."+k] = v
	}

	filterOfDoc[fieldVersion] = version

	col := cli.collection(collection)
	r, err := col.UpdateOne(
		ctx, filterOfDoc, updates,
		&options.UpdateOptions{
			ArrayFilters: &options.ArrayFilters{
				Filters: bson.A{
					arrayFilter,
				},
			},
		},
	)
	if err != nil {
		return false, err
	}

	if r.MatchedCount == 0 {
		return false, errDocNotExists
	}

	return r.ModifiedCount > 0, nil
}

func (cli *client) modifyArrayElemWithoutVersion(
	ctx context.Context, collection, array string,
	filterOfDoc, filterOfArray, updateCmd bson.M,
	op string,
) (bool, error) {
	key := func(k string) string {
		return fmt.Sprintf("%s.$[i].%s", array, k)
	}

	cmd := bson.M{}
	for k, v := range updateCmd {
		cmd[key(k)] = v
	}

	arrayFilter := bson.M{}
	for k, v := range filterOfArray {
		arrayFilter["i."+k] = v
	}

	col := cli.collection(collection)
	r, err := col.UpdateOne(
		ctx, filterOfDoc, bson.M{op: cmd},
		&options.UpdateOptions{
			ArrayFilters: &options.ArrayFilters{
				Filters: bson.A{
					arrayFilter,
				},
			},
		},
	)
	if err != nil {
		return false, err
	}

	if r.MatchedCount == 0 {
		return false, errDocNotExists
	}

	return r.ModifiedCount > 0, nil
}

func (cli *client) updateArrayElemCount(
	ctx context.Context,
	collection, array, field string, num int,
	filterOfDoc, filterOfArray bson.M,
) (bool, error) {
	arrayFilter := bson.M{}
	for k, v := range filterOfArray {
		arrayFilter["i."+k] = v
	}

	col := cli.collection(collection)
	r, err := col.UpdateOne(
		ctx, filterOfDoc,
		bson.M{
			mongoCmdInc: bson.M{fmt.Sprintf("%s.$[i].%s", array, field): num},
		},
		&options.UpdateOptions{
			ArrayFilters: &options.ArrayFilters{
				Filters: bson.A{
					arrayFilter,
				},
			},
		},
	)
	if err != nil {
		return false, err
	}

	if r.MatchedCount == 0 {
		return false, errDocNotExists
	}

	return r.ModifiedCount > 0, nil
}

func (cli *client) isArrayDocExists(
	ctx context.Context, collection string,
	filterOfDoc bson.M, array string, filterOfArray bson.M,
) (bool, error) {
	filterOfDoc[array] = bson.M{mongoCmdElemMatch: filterOfArray}

	return cli.containsArrayElem(ctx, collection, filterOfDoc)
}

func (cli *client) containsArrayElem(
	ctx context.Context, collection string, filterOfDoc bson.M,
) (bool, error) {
	var v struct {
		ID primitive.ObjectID `bson:"_id"`
	}

	err := cli.getDoc(ctx, collection, filterOfDoc, bson.M{"_id": 1}, &v)
	if err != nil {
		if isDocNotExists(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (cli *client) pullArrayElem(
	ctx context.Context, collection, array string,
	filterOfDoc, filterOfArray bson.M,
) error {
	update := bson.M{mongoCmdPull: bson.M{array: filterOfArray}}

	col := cli.collection(collection)

	if _, err := col.UpdateOne(ctx, filterOfDoc, update); err != nil {
		return dbError{err}
	}

	return nil
}

func (cli *client) getArrayElem(
	ctx context.Context, collection, array string,
	filterOfDoc, filterOfArray bson.M,
	project bson.M, result interface{},
) error {
	m := map[string]bson.M{}
	if len(filterOfArray) > 0 {
		m[array] = conditionTofilterArray(filterOfArray)
	}

	return cli.getArraysElemsHelper(
		ctx, collection, filterOfDoc, m,
		project, result,
	)
}

func (cli *client) getArraysElem(
	ctx context.Context, collection string,
	filterOfDoc bson.M, filterOfArrays map[string]bson.M,
	project bson.M, result interface{},
) error {
	m := map[string]bson.M{}
	for k, v := range filterOfArrays {
		m[k] = conditionTofilterArray(v)
	}

	return cli.getArraysElemsHelper(
		ctx, collection, filterOfDoc, m,
		project, result,
	)
}

func (cli *client) getArraysElemsByCustomizedCond(
	ctx context.Context, collection string,
	filterOfDoc bson.M, filterOfArrays map[string]func() bson.M,
	project bson.M, result interface{},
) error {
	m := map[string]bson.M{}
	for k, cond := range filterOfArrays {
		m[k] = cond()
	}

	return cli.getArraysElemsHelper(
		ctx, collection, filterOfDoc, m,
		project, result,
	)
}

func (cli *client) getArraysElemsHelper(
	ctx context.Context, collection string,
	filterOfDoc bson.M, filterOfArrays map[string]bson.M,
	project bson.M, result interface{},
) error {

	pipeline := bson.A{bson.M{mongoCmdMatch: filterOfDoc}}

	if len(filterOfArrays) > 0 {
		project1 := bson.M{}

		for array, cond := range filterOfArrays {
			project1[array] = bson.M{mongoCmdFilter: bson.M{
				"input": fmt.Sprintf("$%s", array),
				"cond":  cond,
			}}
		}

		for k, v := range project {
			s := k
			if i := strings.Index(k, "."); i >= 0 {
				s = k[:i]
			}
			if _, ok := filterOfArrays[s]; !ok {
				project1[k] = v
			}
		}

		pipeline = append(pipeline, bson.M{mongoCmdProject: project1})
	}

	if len(project) > 0 {
		pipeline = append(pipeline, bson.M{mongoCmdProject: project})
	}

	col := cli.collection(collection)
	cursor, err := col.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}

	return cursor.All(ctx, result)
}

func condFieldOfArrayElem(key string) string {
	return "$$this." + key
}

func eqCondForArrayElem(key string, value interface{}) bson.M {
	return bson.M{"$eq": bson.A{condFieldOfArrayElem(key), value}}
}

func inCondForArrayElem(key string, value interface{}) bson.M {
	return bson.M{"$in": bson.A{condFieldOfArrayElem(key), value}}
}

func valueInCondForArrayElem(key string, value interface{}) bson.M {
	return bson.M{"$in": bson.A{value, condFieldOfArrayElem(key)}}
}

func matchCondForArrayElem(key string, value interface{}) bson.M {
	return bson.M{
		"$regexMatch": bson.M{
			"input":   condFieldOfArrayElem(key),
			"regex":   value,
			"options": "i",
		},
	}
}

func condForArrayElem(conds bson.A) bson.M {
	n := len(conds)
	if n > 1 {
		return bson.M{"$and": conds}
	}

	if n == 1 {
		if v, ok := conds[0].(bson.M); ok {
			return v
		}
	}

	return bson.M{
		"$toBool": 1,
	}
}

func conditionTofilterArray(filterOfArray bson.M) bson.M {
	cond := make(bson.A, 0, len(filterOfArray))
	for k, v := range filterOfArray {
		cond = append(cond, eqCondForArrayElem(k, v))
	}

	return condForArrayElem(cond)
}
