package mgo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Test struct {
	Name string
	Age  int
}

func InsertOne(db string, table string, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	col, err := getCollection(db, table)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return col.InsertOne(ctx, document, opts...)
}

func InsertMany(db string, table string, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	col, err := getCollection(db, table)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return col.InsertMany(ctx, documents, opts...)
}

func DeleteOne(db string, table string, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	col, err := getCollection(db, table)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return col.DeleteOne(ctx, filter, opts...)
}

func DeleteMany(db string, table string, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	col, err := getCollection(db, table)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return col.DeleteMany(ctx, filter, opts...)
}

func UpdateOne(db string, table string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	col, err := getCollection(db, table)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return col.UpdateOne(ctx, filter, update, opts...)
}

func UpdateMany(db string, table string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	col, err := getCollection(db, table)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return col.UpdateMany(ctx, filter, update, opts...)
}

func Find(db string, table string, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	col, err := getCollection(db, table)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return col.Find(ctx, filter, opts...)
}

func FindOne(db string, table string, filter interface{}, opts ...*options.FindOneOptions) (*mongo.SingleResult, error) {
	col, err := getCollection(db, table)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return col.FindOne(ctx, filter, opts...), nil
}

func FindOneAndUpdate(db, table string, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) (*mongo.SingleResult, error) {
	col, err := getCollection(db, table)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return col.FindOneAndUpdate(ctx, filter, update, opts...), nil
}

func CountDocuments(db string, table string, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	col, err := getCollection(db, table)
	if err != nil {
		return 0, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return col.CountDocuments(ctx, filter, opts...)
}
