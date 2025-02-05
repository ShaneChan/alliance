package mgo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

// InsertOne 插入一个document
func InsertOne(db string, table string, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	col, err := getCollection(db, table)
	if err != nil {
		log.Println("连接数据库报错，msg: ", err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return col.InsertOne(ctx, document, opts...)
}

// InsertMany 插入多个document数据
func InsertMany(db string, table string, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	col, err := getCollection(db, table)
	if err != nil {
		log.Println("连接数据库报错，msg: ", err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return col.InsertMany(ctx, documents, opts...)
}

// DeleteOne 删除一个document数据
func DeleteOne(db string, table string, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	col, err := getCollection(db, table)
	if err != nil {
		log.Println("连接数据库报错，msg: ", err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return col.DeleteOne(ctx, filter, opts...)
}

// DeleteMany 删除多个document
func DeleteMany(db string, table string, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	col, err := getCollection(db, table)
	if err != nil {
		log.Println("连接数据库报错，msg: ", err)

		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return col.DeleteMany(ctx, filter, opts...)
}

// UpdateOne 更新一个document
func UpdateOne(db string, table string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	col, err := getCollection(db, table)
	if err != nil {
		log.Println("连接数据库报错，msg: ", err)

		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return col.UpdateOne(ctx, filter, update, opts...)
}

// UpdateMany 更新多个document
func UpdateMany(db string, table string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	col, err := getCollection(db, table)
	if err != nil {
		log.Println("连接数据库报错，msg: ", err)

		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return col.UpdateMany(ctx, filter, update, opts...)
}

// Find 查找多个document
func Find(db string, table string, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	col, err := getCollection(db, table)
	if err != nil {
		log.Println("连接数据库报错，msg: ", err)

		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return col.Find(ctx, filter, opts...)
}

// FindOne 查找一个document
func FindOne(db string, table string, filter interface{}, opts ...*options.FindOneOptions) (*mongo.SingleResult, error) {
	col, err := getCollection(db, table)
	if err != nil {
		log.Println("连接数据库报错，msg: ", err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return col.FindOne(ctx, filter, opts...), nil
}

// FindOneAndUpdate 查找并更新一个document
func FindOneAndUpdate(db, table string, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) (*mongo.SingleResult, error) {
	col, err := getCollection(db, table)
	if err != nil {
		log.Println("连接数据库报错，msg: ", err)

		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return col.FindOneAndUpdate(ctx, filter, update, opts...), nil
}

// CountDocuments 获得document数量
func CountDocuments(db string, table string, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	col, err := getCollection(db, table)
	if err != nil {
		log.Println("连接数据库报错，msg: ", err)

		return 0, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return col.CountDocuments(ctx, filter, opts...)
}
