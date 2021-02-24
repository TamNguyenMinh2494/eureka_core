package business

import (
	"context"
	"main/core/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ExamBusiness struct {
	DB           *mongo.Database
	QuizSelector *QuizSelectorBusiness
}

func (b *ExamBusiness) GetById(examId string) (models.Exams, error) {
	exam := new(models.Exams)
	objectId, err := primitive.ObjectIDFromHex(examId)
	if err != nil {
		return *exam, err
	}
	r := b.DB.Collection("exams").FindOne(context.TODO(), bson.M{"_id": objectId})
	if r.Err() != nil {
		return *exam, r.Err()
	}
	err = r.Decode(exam)
	return *exam, err
}

func (b *ExamBusiness) Create(exam *models.Exams) error {
	exam.Id = primitive.NewObjectID()
	_, err := b.DB.Collection("exams").InsertOne(context.TODO(), exam)
	return err
}

func (b *ExamBusiness) Update(exam *models.Exams) error {
	_, err := b.DB.Collection("exams").UpdateOne(context.TODO(), bson.M{"_id": exam.Id}, exam)
	return err
}

func (b *ExamBusiness) Delete(examId string) error {
	objectId, err := primitive.ObjectIDFromHex(examId)
	if err != nil {
		return err
	}
	_, err = b.DB.Collection("exams").DeleteOne(context.TODO(), bson.M{"_id": objectId})
	return err
}

func (b *ExamBusiness) Preview(examId string) (models.TakenExams, error) {
	takenExam := new(models.TakenExams)
	exam, err := b.GetById(examId)
	if err != nil {
		return *takenExam, nil
	}
	return b.QuizSelector.Select("", &exam)
}

func (b *ExamBusiness) Take(email string, examId string) (models.TakenExams, error) {
	takenExam := *new(models.TakenExams)
	exam, err := b.GetById(examId)
	if err != nil {
		return takenExam, nil
	}
	takenExam, err = b.QuizSelector.Select(email, &exam)
	if err != nil {
		return takenExam, nil
	}
}
