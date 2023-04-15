package repositories

import (
	"chat-app/config"
	"fmt"
	"gopkg.in/mgo.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

var dbSingleton *gorm.DB
var dbOnce sync.Once

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

func InitChatMessageDatabase() *mgo.Session {
	dbConfig := config.GetAppConfig().ChatMessageDatabase
	info := &mgo.DialInfo{
		Addrs:    []string{fmt.Sprintf("%s:%d", dbConfig.Host, dbConfig.Port)},
		Username: dbConfig.Username,
		Password: dbConfig.Password,
		Database: dbConfig.DatabaseName,
	}

	session, err := mgo.DialWithInfo(info)
	if err != nil {
		panic(fmt.Errorf("failed to connect chat message database, error: %v", err))
	}

	session.SetMode(mgo.Monotonic, true)

	err = session.Ping()
	if err != nil {
		panic(fmt.Errorf("failed to connect chat message database, error: %v", err))
	}

	return session
}
