package api

import (
	"AllianceServer/mgo"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
)

func Register(account string, password string) error {
	cur, err := mgo.FindOne("testing", "user", bson.M{"account": account})
	if err != nil {
		return err
	}
	var result bson.M
	err = cur.Decode(&result)
	if len(result) != 0 {
		return errors.New("exist account")
	}
	_, err = mgo.InsertOne("testing", "user", bson.M{"account": account, "password": password})
	return err
}

func Login(account string, password string) error {
	cur, err := mgo.FindOne("testing", "user", bson.M{"account": account, "password": password})
	if err != nil {
		return err
	}
	var result bson.M
	err = cur.Decode(&result)
	if len(result) == 0 {
		return errors.New("no account")
	}

	return nil
}

func Unregister(account string, password string) error {
	return nil
}
