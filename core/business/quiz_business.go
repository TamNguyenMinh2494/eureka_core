package business

import (
	"context"
	"encoding/json"
	"main/core/models"
	"main/core/question"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type QuizBusiness struct {
	DB *mongo.Database
}

func (b *QuizBusiness) GetById(id string) (models.Quiz, error) {
	quiz := new(models.Quiz)
	r := b.DB.Collection("question_bank").FindOne(context.TODO(), bson.M{"_id": id})
	if r.Err() != nil {
		return *quiz, r.Err()
	}
	err := r.Decode(quiz)
	return *quiz, err
}

func (b *QuizBusiness) GetByCourse(courseId string) ([]models.Quiz, error) {
	cursor, err := b.DB.Collection("question_bank").Find(context.TODO(), bson.M{"course_id": courseId})
	if err != nil {
		return nil, err
	}

	quizes := make([]models.Quiz, 0)

	for cursor.Next(context.TODO()) {
		quiz := new(models.Quiz)
		err = cursor.Decode(quiz)
		if err != nil {
			return quizes, err
		}
		quizes = append(quizes, *quiz)
	}
	return quizes, nil
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
		question := new(question.MultipleChoice)
		err := json.Unmarshal([]byte(quiz.Question), question)
		if err != nil {
			return *structuredQuiz, err
		}
		structuredQuiz.Question = question
		break
	}

	return *structuredQuiz, nil
}
