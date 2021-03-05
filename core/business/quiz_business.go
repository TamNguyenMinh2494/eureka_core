package business

import (
	"context"
	"encoding/json"
	"errors"
	"main/core/models"
	"main/core/question"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type QuizBusiness struct {
	DB *mongo.Database
}

func (b *QuizBusiness) GetById(id string) (models.Quiz, error) {
	quiz := new(models.Quiz)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return *quiz, err
	}
	r := b.DB.Collection("question_bank").FindOne(context.TODO(), bson.M{"_id": objectId})
	if r.Err() != nil {
		return *quiz, r.Err()
	}
	err = r.Decode(quiz)
	return *quiz, err
}

func (b *QuizBusiness) GetByCourse(courseId string) ([]models.Quiz, error) {
	cursor, err := b.DB.Collection("question_bank").Find(context.TODO(), bson.M{"courseid": courseId})
	if err != nil {
		return nil, err
	}

	quizzes := make([]models.Quiz, 0)

	for cursor.Next(context.TODO()) {
		quiz := new(models.Quiz)
		err = cursor.Decode(quiz)
		if err != nil {
			return quizzes, err
		}
		quizzes = append(quizzes, *quiz)
	}
	return quizzes, nil
}

func (b *QuizBusiness) Create(quiz *models.Quiz) error {
	quiz.LastUpdate = time.Now().Unix()
	quiz.Id = primitive.NewObjectID()
	_, err := b.DB.Collection("question_bank").InsertOne(context.TODO(), quiz)
	return err
}

func (b *QuizBusiness) Update(quiz *models.Quiz) error {
	// objectId, err := primitive.ObjectIDFromHex(quiz.id)
	// if err != nil {
	// 	return *quiz, err
	// }
	r := b.DB.Collection("question_bank").FindOneAndUpdate(context.TODO(), bson.M{"_id": quiz.Id}, bson.M{"$set": quiz})
	return r.Err()
}

func (b *QuizBusiness) Delete(id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = b.DB.Collection("question_bank").DeleteOne(context.TODO(), bson.M{"_id": objectId})
	return err
}

func (b *QuizBusiness) Parse(quiz *models.Quiz) (models.StructuredQuiz, error) {
	structuredQuiz := new(models.StructuredQuiz)
	structuredQuiz.Raw = *quiz
	switch structuredQuiz.Raw.Type {
	case "multiple_choices":
		question := new(question.MultipleChoice)
		err := json.Unmarshal([]byte(quiz.Question), question)
		if err != nil {
			return *structuredQuiz, err
		}
		structuredQuiz.Question = question
		break
	default:
		return *structuredQuiz, errors.New("Unknown question type: " + structuredQuiz.Raw.Type)
	}

	return *structuredQuiz, nil
}
