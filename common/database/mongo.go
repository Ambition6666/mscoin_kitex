package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type MongoClient struct {
	cli *mongo.Client
	Db  *mongo.Database
}

func ConnectMongo(username string, password string, database string, url string) *MongoClient {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	credential := options.Credential{
		Username:   username,
		Password:   password,
		AuthSource: database, // 如果使用默认身份验证数据库
	}
	client, err := mongo.Connect(ctx, options.Client().
		ApplyURI(url).
		SetAuth(credential))
	if err != nil {
		panic(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}
	dat := client.Database(database)
	return &MongoClient{cli: client, Db: dat}
}

func (c *MongoClient) Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := c.cli.Disconnect(ctx)
	if err != nil {
		log.Println(err)
	}
	log.Println("关闭连接..")
}
