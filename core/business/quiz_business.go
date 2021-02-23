package business

import (
	"context"
	"main/core/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type QuizBusiness struct {
	DB *mongo.Database
}

func (b *QuizBusiness) Create(quiz *models.Quiz) error {
	_, err := b.DB.Collection("question_bank").InsertOne(context.TODO(), quiz)
	return err
}

func (b *QuizBusiness) Update(quiz *models.Quiz) error {
	r := b.DB.Collection("question_bank").FindOneAndUpdate(context.TODO(), bson.M{"_id": quiz.Id}, bson.M{"$set": quiz})
	return r.Err()
}

func (b *QuizBusiness) Delete(id string) error {
	_, err := b.DB.Collection("question_bank").DeleteOne(context.TODO(), bson.M{"_id": id})
	return err
}

func (b *QuizBusiness) Parse(quiz *models.Quiz) (models.StructuredQuiz, error) {
	structuredQuiz := new(models.StructuredQuiz)
	structuredQuiz.Raw = *quiz
	switch structuredQuiz.Raw.Type {
	case "multiple_choice":

		break
	}

	return *structuredQuiz, nil
}
