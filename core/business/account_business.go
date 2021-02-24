package business

import (
	"context"
	"errors"
	"main/core/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AccountBusiness struct {
	DB *mongo.Database
}

func (b *AccountBusiness) CreateIndexes() {
	b.DB.Collection("accounts").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	})
}

func (b *AccountBusiness) Create(email string) (err error) {
	account := models.Account{
		Email:   email,
		Balance: 0,
	}
	_, err = b.DB.Collection("accounts").InsertOne(context.TODO(), account)
	return err
}

func (b *AccountBusiness) Get(email string) (account *models.Account, err error) {
	result := b.DB.Collection("accounts").FindOne(context.TODO(), bson.M{"email": email})
	if result.Err() != nil {
		return nil, result.Err()
	}
	account = new(models.Account)
	err = result.Decode(account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (b *AccountBusiness) Charge(email string, amount int64) error {
	account, err := b.Get(email)
	if err != nil {
		return err
	}
	if account.Balance < amount {
		return errors.New("Balance is not enough")
	}
	result := b.DB.Collection("accounts").FindOneAndUpdate(context.TODO(),
		bson.M{"email": email},
		bson.D{{"$inc", bson.D{{"balance", -1 * amount}}}})
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}
