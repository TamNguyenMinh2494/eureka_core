package business

import (
	"main/core/models"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuizSelectorBusiness struct {
	Quiz *QuizBusiness
}

func (b *QuizSelectorBusiness) Select(email string, exam *models.Exams) (models.TakenExams, error) {
	takenExams := new(models.TakenExams)
	quizzes := make([]models.Quiz, 0)
	rand.Seed(time.Now().UnixNano())
	for _, selector := range exam.QuizSelector {
		list := make([]models.Quiz, 0)
		for _, quizId := range selector.List {
			q, err := b.Quiz.GetById(quizId)
			if err != nil {
				return *takenExams, err
			}
			list = append(list, q)
		}
		if selector.Shuffle {
			rand.Shuffle(len(list), func(i, j int) { list[i], list[j] = list[j], list[i] })
		}
		quizzes = append(quizzes, list...)
	}
	takenExams.Quizzes = quizzes
	takenExams.Id = primitive.NewObjectID()
	takenExams.CreatedDate = time.Now().Unix()
	takenExams.Examinee = email
	return *takenExams, nil
}
