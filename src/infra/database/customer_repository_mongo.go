package database

import (
	"context"
	"github.com/tonycarvalho1994/rinha_backend_2024_q1/src/core/entity"
	"go.mongodb.org/mongo-driver/mongo/options"
	_ "time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomerRepositoryMongo struct {
	Collection *mongo.Collection
	Ctx        context.Context
}

func (c *CustomerRepositoryMongo) CreateCustomer(id string, limit int) error {
	_, err := c.Collection.InsertOne(c.Ctx, entity.Customer{
		ID:           id,
		Limit:        limit,
		Transactions: []entity.Transaction{},
	})
	return err
}

func (c *CustomerRepositoryMongo) FindById(id string) (*entity.Customer, error) {
	var result entity.Customer
	err := c.Collection.FindOne(context.Background(), bson.M{"id": id}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *CustomerRepositoryMongo) Update(customer *entity.Customer) error {
	filter := bson.M{"id": customer.ID}
	update := bson.M{"$set": bson.M{"transactions": customer.Transactions}}
	_, err := c.Collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (c *CustomerRepositoryMongo) AddTransaction(customerId string, transaction *entity.Transaction) (int, error) {
	return 0, nil
}

func CreateConnection(uri, databaseName, collectionName string, ctx context.Context) (*mongo.Collection, *mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, nil, err
	}
	collection := client.Database(databaseName).Collection(collectionName)

	return collection, client, nil
}
