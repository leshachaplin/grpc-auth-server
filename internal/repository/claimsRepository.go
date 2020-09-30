package repository

import (
	"context"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Claim struct {
	Id          bson.ObjectId `bson:"_id"`
	Login       string        `bson:"login"`
	Description string        `bson:"key"`
	Value       string        `bson:"value"`
}

type RepositoryOfClaims struct {
	db *mgo.Database
}

func NewRepositoryOfClaims(database *mgo.Database) *RepositoryOfClaims {
	return &RepositoryOfClaims{
		db: database,
	}
}

func (r *RepositoryOfClaims) GetClaims(ctx context.Context, login string) (map[string]string, error) {
	var err error
	query := bson.M{
		"login": bson.M{
			"$eq": login,
		},
	}
	claims := []Claim{}
	result := make(map[string]string, 0)

	err = r.db.C("claim").Find(query).All(&claims)
	if err != nil {
		return nil, err
	}
	for _, v := range claims {
		result[v.Description] = v.Value
	}
	return result, err
}

func (r *RepositoryOfClaims) IfExistClaim(ctx context.Context, key, login string) (bool, error) {
	var err error
	query := bson.M{
		"login": bson.M{
			"$eq": login,
		},
	}
	claims := []Claim{}

	err = r.db.C("claim").Find(query).All(&claims)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *RepositoryOfClaims) AddClaims(ctx context.Context, claims map[string]string, login string) error {
	var err error
	for k, v := range claims {
		c := &Claim{
			Id:          bson.NewObjectId(),
			Login:       login,
			Description: k,
			Value:       v,
		}
		err = r.db.C("claim").Insert(&c)
	}
	return err
}

func (r *RepositoryOfClaims) DeleteClaims(ctx context.Context, claims map[string]string, login string) error {
	var err error
	for k, v := range claims {
		_, err = r.db.C("claim").RemoveAll(bson.M{"key" : k, "value":v})
	}
	return err
}
