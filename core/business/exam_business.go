package business

import (
	"context"
	"errors"
	"main/core/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ExamBusiness struct {
	DB           *mongo.Database
	QuizSelector *QuizSelectorBusiness
	Quiz         *QuizBusiness
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

func (b *ExamBusiness) GetExamsBySectionId(sectionId string) ([]models.Exams, error) {
	exams := make([]models.Exams, 0)
	cursor, err := b.DB.Collection("taken_exams").Find(context.TODO(), bson.M{"sectionId": sectionId})
	if err != nil {
		return exams, nil
	}
	for cursor.Next(context.TODO()) {
		exam := new(models.Exams)
		cursor.Decode(exam)
		exams = append(exams, *exam)
	}
	return exams, nil
}

func (b *ExamBusiness) GetTakenExamById(takenExamId string) (*models.TakenExams, error) {
	takenExam := new(models.TakenExams)
	objectId, err := primitive.ObjectIDFromHex(takenExamId)
	if err != nil {
		return takenExam, err
	}
	r := b.DB.Collection("taken_exams").FindOne(context.TODO(), bson.M{"_id": objectId})
	if r.Err() != nil {
		return takenExam, r.Err()
	}
	err = r.Decode(takenExam)
	return takenExam, err
}

func (b *ExamBusiness) GetTakenExams(examinee string) ([]models.TakenExams, error) {
	takenExams := make([]models.TakenExams, 0)
	cursor, err := b.DB.Collection("taken_exams").Find(context.TODO(), bson.M{"examinee": examinee})
	if err != nil {
		return takenExams, nil
	}
	for cursor.Next(context.TODO()) {
		takenExam := new(models.TakenExams)
		cursor.Decode(takenExam)
		takenExams = append(takenExams, *takenExam)
	}
	return takenExams, nil
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
	for i := range takenExam.Quizzes {
		takenExam.Quizzes[i].CorrectAnswer = ""
	}
	takenExam.CreatedDate = time.Now().UnixNano()
	_, err = b.DB.Collection("taken_exams").InsertOne(context.TODO(), takenExam)
	takenExam.Duration = exam.Duration
	return takenExam, err
}

func (b *ExamBusiness) GetSubmit(takenExamId string) (*models.SubmittedExams, error) {
	submittedExam := new(models.SubmittedExams)
	objectId, err := primitive.ObjectIDFromHex(takenExamId)

	if err != nil {
		return submittedExam, err
	}
	r := b.DB.Collection("submissions").FindOne(context.TODO(), bson.M{"_id": objectId})
	if r.Err() != nil {
		return submittedExam, r.Err()
	}
	err = r.Decode(submittedExam)
	return submittedExam, err
}

func (b *ExamBusiness) Submit(answerSheet *models.SubmittedExams) error {
	_, err := b.GetTakenExamById(answerSheet.Id.Hex())
	if err != nil {
		return errors.New("Cannot submit to inexisted exam")
	}

	var totalScore float32 = 0
	var score float32 = 0

	for i, rawQuiz := range answerSheet.Quizzes {
		structuredQuiz, err := b.Quiz.Parse(&rawQuiz)
		if err != nil {
			continue
		}
		s := structuredQuiz.Question.CheckAnswer(rawQuiz.Answer)
		if s > 0 {
			totalScore += s
		}
		score += s
		answerSheet.Quizzes[i].Score = score
	}
	answerSheet.Score = score
	answerSheet.TotalScore = totalScore
	answerSheet.SubmitDate = time.Now().UnixNano()

	_, err = b.DB.Collection("submissions").InsertOne(context.TODO(), answerSheet)

	return err
}

//func (b *ExamBusiness)
