package repositories

import (
	"chat-app/config"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"sync"
)

var dbSingleton *gorm.DB
var dbOnce sync.Once

type MongoDatastore struct {
	DB      *mongo.Database
	Session *mongo.Client
}

var mongoDataStore *MongoDatastore
var connectOnce sync.Once

func InitChatAppDatabase() *gorm.DB {
	dbOnce.Do(
		func() {
			dbConfig := config.GetAppConfig().ChatAppDatabase
			dsn := fmt.Sprintf(
				"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				dbConfig.Username,
				dbConfig.Password,
				dbConfig.Host,
				dbConfig.Port,
				dbConfig.DatabaseName,
			)

			db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
			if err != nil {
				panic(fmt.Errorf("failed to connect chat app database, error: %v", err))
			}
			dbSingleton = db
		},
	)

	return dbSingleton
}

func InitChatMessageDatabase() *MongoDatastore {
	connectOnce.Do(func() {
		dbConfig := config.GetAppConfig().ChatMessageDatabase

		connectionString := fmt.Sprintf(
			"mongodb://%s:%s@%s:%d/%s",
			dbConfig.Username,
			dbConfig.Password,
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.DatabaseName,
		)

		clientOptions := options.Client().ApplyURI(connectionString)
		client, err := mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			log.Fatal(err)
		}

		db := client.Database(dbConfig.DatabaseName)
		mongoDataStore = &MongoDatastore{
			DB:      db,
			Session: client,
		}
	})
	return mongoDataStore
}
