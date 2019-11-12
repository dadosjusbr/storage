package storage

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Client manages all iteractions with mongodb
type Client struct {
	client *mongo.Client
	dbName string
}

//NewClient returns an db connection instance that can be used for CRUD operations
func NewClient(url, dbName string) (*Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := client.Connect(ctx); err != nil {
		fmt.Printf("could not ping to mongo db service: %v\n", err)
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")
	return &Client{client, dbName}, nil
}

//SaveAgency save or update a agency
func (db *Client) SaveAgency(ag Agency) error {
	// Get Agency's Collection
	collection := db.getCollection("agency")

	// Insert a single document
	_, err := collection.ReplaceOne(context.TODO(), bson.D{{Key: "short_name", Value: ag.ShortName}}, ag, options.Replace().SetUpsert(true))
	if err != nil {
		return err
	}
	fmt.Printf("Upserted a single document: %s\n", ag.Name)
	return nil
}

//SaveMonthResults save month results
func (db *Client) SaveAgencyMonthInfo(agMonth AgencyMonthlyInfo) error {
	// Get a handle for your collection
	collection := db.getCollection("AgencyMonthlyInfo")

	// Insert a single document
	_, err := collection.ReplaceOne(context.TODO(), bson.D{{Key: "_id", Value: agMonth.AgencyID}}, agMonth, options.Replace().SetUpsert(true))
	if err != nil {
		return err
	}
	fmt.Printf("Upserted a single document: %s\n", agMonth.AgencyID)
	return nil
}

//CloseConnection closes the opened connetion to mongodb
func (db *Client) CloseConnection() error {
	err := db.client.Disconnect(context.TODO())

	if err != nil {
		return err
	}

	fmt.Println("Connection to MongoDB closed.")
	return nil
}

func (db *Client) getCollection(collName string) *mongo.Collection {
	return db.client.Database(db.dbName).Collection(collName)
}
