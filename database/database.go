package database

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

var (
	Db               *mongo.Database
	MONGODB_USER     = "root"
	MONGODB_PASSWORD = "root"
	MONGODB_HOST     = "localhost:27017"
	MONGODB_DBNAME   = "pricesTracker"
)

// always runs
func init() {
	if len(os.Getenv("MONGODB_HOST")) < 1 {
		err := godotenv.Load()
		if err != nil {
			log.Print(err)
		}
	}

	MONGODB_USER = os.Getenv("MONGODB_USER")
	MONGODB_PASSWORD = os.Getenv("MONGODB_PASSWORD")
	MONGODB_DBNAME = os.Getenv("MONGODB_DBNAME")
	MONGODB_HOST = os.Getenv("MONGODB_HOST")

	// debug line
	//fmt.Printf("\nU:%s, P:%s, H:%s, D:%s\n\n", MONGODB_USER, MONGODB_PASSWORD, MONGODB_HOST, MONGODB_DBNAME)
	//migrateDatabase()
}

func ConnectDatabase() {
	//if envMongoDbUser := os.Getenv("MONGODB_USER"); envMongoDbUser != "" {
	//	MONGODB_USER = envMongoDbUser
	//}
	//
	//if envMongoDbPassword := os.Getenv("MONGODB_PASSWORD"); envMongoDbPassword != "" {
	//	MONGODB_PASSWORD = envMongoDbPassword
	//}
	//
	//if envMongoDbHost := os.Getenv("MONGODB_HOST"); envMongoDbHost != "" {
	//	MONGODB_HOST = envMongoDbHost
	//}

	// mongodb+srv
	// This works locally::
	//uri := "mongodb://" + MONGODB_USER + ":" + MONGODB_PASSWORD + "@" + MONGODB_HOST + "/" + MONGODB_DBNAME + "?retryWrites=true&w=majority"
	//uri := "mongodb://" + MONGODB_USER + ":" + MONGODB_PASSWORD + "@" + MONGODB_HOST
	//fmt.Println(uri)
	//uri := "mongodb://root:root@localhost:27017"

	// this works on HEROKU
	uri := "mongodb+srv://" + MONGODB_USER + ":" + MONGODB_PASSWORD + "@" + MONGODB_HOST + "/" + MONGODB_DBNAME + "?retryWrites=true&w=majority"

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	err = client.Connect(ctx)

	if err != nil {
		panic(err)
	}

	Db = client.Database(MONGODB_DBNAME)
	fmt.Println("Successfuly connected to the database.")

	migrateDatabase()
}

func migrateDatabase() {

}
