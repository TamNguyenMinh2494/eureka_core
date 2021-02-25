package business

import (
	"context"
	"main/core/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionBusiness struct {
	DB *mongo.Database
}

func (b *TransactionBusiness) Purchase(account *AccountBusiness, transaction *models.Transaction) (err error) {
	err = account.Charge(transaction.Email, -transaction.Amount*transaction.Quantity)
	if err != nil {
		return err
	}
	_, err = b.DB.Collection("transactions").InsertOne(context.TODO(), transaction)
	return err
}

func (b *TransactionBusiness) Get(email string, sku string) ([]models.Transaction, error) {
	var cursor *mongo.Cursor
	var err error

	if email != "" && sku != "" {
		cursor, err = b.DB.Collection("transactions").Find(context.TODO(), bson.M{"email": email, "sku": sku})
	} else if email != "" {
		cursor, err = b.DB.Collection("transactions").Find(context.TODO(), bson.M{"email": email})
	} else if sku != "" {
		cursor, err = b.DB.Collection("transactions").Find(context.TODO(), bson.M{"sku": sku})
	}

	if err != nil {
		return nil, err
	}
	transactions := make([]models.Transaction, 0)
	for cursor.Next(context.TODO()) {
		trans := new(models.Transaction)
		err = cursor.Decode(trans)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, *trans)
	}

	return transactions, nil
}
