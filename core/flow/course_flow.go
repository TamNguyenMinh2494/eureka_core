package flow

import (
	"errors"
	"main/core/business"
	"main/core/models"
)

type Course struct {
	course        *business.CourseBusiness
	courseSection *business.CourseSectionBusiness
	user          *business.UserBusiness
	quiz          *business.QuizBusiness
}

// func (f *Course) GetQuiz(email string, courseId string) ([]models.Quiz,error) {
// 	if !f.course.IsAuthor(courseId, email) {
// 		return errors.New("Cannot create quiz for nonpossession course")
// 	}
// }

func (f *Course) CreateQuiz(email string, courseId string, quiz *models.Quiz) error {
	if !f.course.IsAuthor(courseId, email) {
		return errors.New("Cannot create quiz for nonpossession course")
	}
	return f.quiz.Create(quiz)
}

func (f *Course) UpdateQuiz(email string, courseId string, quiz *models.Quiz) error {
	if !f.course.IsAuthor(courseId, email) {
		return errors.New("Cannot update quiz for nonpossession course")
	}
	quizFromDB, err := f.quiz.GetById(quiz.Id.String())
	if err != nil {
		return err
	}
	if quizFromDB.CourseId != courseId {
		return errors.New("This quiz does not exist in the course")
	}

	return f.quiz.Update(quiz)
}

func (f *Course) DeleteQuiz(email string, courseId string, quizId string) error {
	if !f.course.IsAuthor(courseId, email) {
		return errors.New("Cannot create quiz for nonpossession course")
	}
	quizFromDB, err := f.quiz.GetById(quizId)
	if err != nil {
		return err
	}
	if quizFromDB.CourseId != courseId {
		return errors.New("This quiz does not exist in the course")
	}
	return f.quiz.Delete(quizId)
}
