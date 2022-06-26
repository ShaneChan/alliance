package mgo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

type Client struct {
	client *mongo.Client  // mongo连接句柄
	dbs    map[string]*db // client实例持有的db
}

type db struct {
	database    *mongo.Database              // db
	collections map[string]*mongo.Collection // 一个db拥有的collection
}

func (c *Client) getConnection(database string, table string) *mongo.Collection {
	if c.dbs[database] == nil {
		c.dbs[database] = &db{
			database:    c.client.Database(database),
			collections: map[string]*mongo.Collection{},
		}
	}

	if c.dbs[database].collections[table] == nil {
		c.dbs[database].collections[table] = c.dbs[database].database.Collection(table)
	}

	return c.dbs[database].collections[table]
}

var client *Client // 全局共用一个mongo连接实例

// 程序启动的时候初始化mongo连接实例
func init() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	option := options.Client().ApplyURI("mongodb://localhost:27017")
	cli, err := mongo.Connect(ctx, option)
	if err != nil {
		_ = cli.Disconnect(ctx)
		log.Fatalln("connect to mongo error")
	}

	if err := cli.Ping(context.Background(), readpref.Primary()); err != nil {
		_ = cli.Disconnect(ctx)
		log.Fatalln("ping mongo error")
	}

	client = &Client{
		client: cli,
		dbs:    map[string]*db{},
	}
}

func getCollection(db string, table string) (*mongo.Collection, error) {
	return client.getConnection(db, table), nil
}
