package flow

import (
	"errors"
	"main/core/business"
	"main/core/models"
	"main/core/utils"
)

type Exam struct {
	Course        *business.CourseBusiness
	CourseSection *business.CourseSectionBusiness
	Exam          *business.ExamBusiness
}

func (f *Exam) CreateExam(email string, courseId string, exam *models.Exams) error {

	if !f.Course.IsAuthor(courseId, email) {
		return errors.New("Cannot create exam for nonpossession course")
	}

	sections, err := f.CourseSection.GetSectionsByCourse(courseId)
	if err != nil {
		return err
	}
	ind, _ := utils.Find(sections, func(index int, elem interface{}) bool {
		return elem.(models.CourseSection).Id.Hex() == exam.SectionId
	})
	if ind == -1 {
		return errors.New("Cannot create exam for an inexistent section")
	}

	return f.Exam.Create(exam)

}

func (f *Exam) UpdateExam(email string, courseId string, exam *models.Exams) error {

	if !f.Course.IsAuthor(courseId, email) {
		return errors.New("Cannot update exam for nonpossession course")
	}

	sections, err := f.CourseSection.GetSectionsByCourse(courseId)
	if err != nil {
		return err
	}
	ind, _ := utils.Find(sections, func(index int, elem interface{}) bool {
		return elem.(models.CourseSection).Id.Hex() == exam.SectionId
	})
	if ind == -1 {
		return errors.New("Cannot update exam for an inexistent section")
	}

	return f.Exam.Update(exam)

}

func (f *Exam) DeleteExam(email string, courseId string, examId string) error {

	if !f.Course.IsAuthor(courseId, email) {
		return errors.New("Cannot create exam for nonpossession course")
	}

	sections, err := f.CourseSection.GetSectionsByCourse(courseId)
	if err != nil {
		return err
	}

	exam, err := f.Exam.GetById(examId)
	if err != nil {
		return err
	}
	ind, _ := utils.Find(sections, func(index int, elem interface{}) bool {
		return elem.(models.CourseSection).Id.Hex() == exam.SectionId
	})
	if ind == -1 {
		return errors.New("Cannot delete exam")
	}

	return f.Exam.Delete(examId)

}
