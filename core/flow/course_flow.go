package flow

import (
	"errors"
	"main/core/business"
	"main/core/models"
)

type Course struct {
	Course        *business.CourseBusiness
	CourseSection *business.CourseSectionBusiness
	User          *business.UserBusiness
	Quiz          *business.QuizBusiness
}

func (f *Course) GetQuiz(email string, courseId string) ([]models.Quiz, error) {
	if !f.Course.IsAuthor(courseId, email) {
		return nil, errors.New("Cannot get quiz of nonpossession course")
	}
	return f.Quiz.GetByCourse(courseId)
}

func (f *Course) CreateQuiz(email string, courseId string, quiz *models.Quiz) error {
	if !f.Course.IsAuthor(courseId, email) {
		return errors.New("Cannot create quiz for nonpossession course")
	}
	_, err := f.Quiz.Parse(quiz)
	if err != nil {
		return err
	}
	return f.Quiz.Create(quiz)
}

func (f *Course) UpdateQuiz(email string, courseId string, quiz *models.Quiz) error {
	if !f.Course.IsAuthor(courseId, email) {
		return errors.New("Cannot update quiz for nonpossession course")
	}
	quizFromDB, err := f.Quiz.GetById(quiz.Id.String())
	if err != nil {
		return err
	}
	if quizFromDB.CourseId != courseId {
		return errors.New("This quiz does not exist in the course")
	}

	_, err = f.Quiz.Parse(quiz)
	if err != nil {
		return err
	}

	return f.Quiz.Update(quiz)
}

func (f *Course) DeleteQuiz(email string, courseId string, quizId string) error {
	if !f.Course.IsAuthor(courseId, email) {
		return errors.New("Cannot delete quiz for nonpossession course")
	}
	quizFromDB, err := f.Quiz.GetById(quizId)
	if err != nil {
		return err
	}
	if quizFromDB.CourseId != courseId {
		return errors.New("This quiz does not exist in the course")
	}
	return f.Quiz.Delete(quizId)
}
