package main

import (
	"fmt"

	"go.uber.org/zap"
	"gopkg.in/mgo.v2"
)

// ConnectDB connects and authenticates against a MongoDB server
func (app *App) ConnectDB() {
	var err error

	app.mongo, err = mgo.Dial(fmt.Sprintf("%s:%s", app.config.MongoHost, app.config.MongoPort))
	if err != nil {
		logger.Fatal("failed to connect to mongodb",
			zap.Error(err))
	}
	logger.Info("connected to mongodb server")

	if app.config.MongoPass != "" {
		err = app.mongo.Login(&mgo.Credential{
			Source:   app.config.MongoName,
			Username: app.config.MongoUser,
			Password: app.config.MongoPass,
		})
		if err != nil {
			logger.Fatal("failed to log in to mongodb",
				zap.Error(err))
		}
		logger.Info("logged in to mongodb server")
	}
	if !app.CollectionExists(app.config.MongoCollection) {
		err = app.mongo.DB(app.config.MongoName).C(app.config.MongoCollection).Create(&mgo.CollectionInfo{})
		if err != nil {
			logger.Fatal("collection create failed",
				zap.String("collection", app.config.MongoCollection),
				zap.Error(err))
		}
	}
	app.cln = app.mongo.DB(app.config.MongoName).C(app.config.MongoCollection)

}

// CollectionExists checks if a collection exists in MongoDB
func (app *App) CollectionExists(name string) bool {
	collections, err := app.mongo.DB(app.config.MongoName).CollectionNames()
	if err != nil {
		logger.Fatal("failed to get collection names",
			zap.Error(err))
	}

	for _, collection := range collections {
		if collection == name {
			return true
		}
	}

	return false
}
