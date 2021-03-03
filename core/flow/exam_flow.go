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
	Enrollment    *business.EnrollmentBusiness
}

func (f *Exam) GetExams(email string, courseId string, sectionId string) ([]models.Exams, error) {

	isEnrolled := false

	if !f.Course.IsAuthor(courseId, email) {
		if !f.Enrollment.IsEnroll(courseId, email) {
			return nil, errors.New("Cannot create exam for nonpossession and unenrolled course")
		}
		isEnrolled = true

	}

	sections, err := f.CourseSection.GetSectionsByCourse(courseId)
	if err != nil {
		return nil, err
	}

	ind, _ := utils.Find(sections, func(index int, elem interface{}) bool {
		return elem.(models.CourseSection).Id.Hex() == sectionId
	})
	if ind == -1 {
		return nil, errors.New("Cannot get exams of an inexistent section")
	}

	r, err := f.Exam.GetExamsBySectionId(sectionId)
	if isEnrolled {
		for i := range r {
			r[i].QuizSelector = nil
		}
	}
	return r, err
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
		return errors.New("Cannot delete exam for nonpossession course")
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

func (f *Exam) PreviewExam(email string, courseId string, examId string) (models.TakenExams, error) {
	takenExam := new(models.TakenExams)
	if !f.Course.IsAuthor(courseId, email) {
		return *takenExam, errors.New("Cannot preview exam for nonpossession course")
	}

	sections, err := f.CourseSection.GetSectionsByCourse(courseId)
	if err != nil {
		return *takenExam, err
	}

	exam, err := f.Exam.GetById(examId)
	if err != nil {
		return *takenExam, err
	}
	ind, _ := utils.Find(sections, func(index int, elem interface{}) bool {
		return elem.(models.CourseSection).Id.Hex() == exam.SectionId
	})
	if ind == -1 {
		return *takenExam, errors.New("Cannot preview exam")
	}
	return f.Exam.Preview(examId)
}

func (f *Exam) TakeExam(examinee string, courseId string, examId string) (models.TakenExams, error) {
	takenExam := new(models.TakenExams)
	if !f.Enrollment.IsEnroll(courseId, examinee) {
		return *takenExam, errors.New("Cannot take exam for unenrolled course")
	}

	sections, err := f.CourseSection.GetSectionsByCourse(courseId)
	if err != nil {
		return *takenExam, err
	}

	exam, err := f.Exam.GetById(examId)
	if err != nil {
		return *takenExam, err
	}
	ind, _ := utils.Find(sections, func(index int, elem interface{}) bool {
		return elem.(models.CourseSection).Id.Hex() == exam.SectionId
	})
	if ind == -1 {
		return *takenExam, errors.New("Cannot take exam")
	}

	return f.Exam.Take(examinee, examId)
}

func (f *Exam) SubmitExam(courseId string, submission *models.SubmittedExams) error {
	takenExam, err := f.Exam.GetTakenExamById(submission.Id.Hex())
	if err != nil {
		return err
	}

	if !f.Enrollment.IsEnroll(courseId, takenExam.Examinee) {
		return errors.New("Cannot submit an exam for unenrolled course")
	}

	// WARNING !!! Skipping some illogical cases

	return f.Exam.Submit(submission)

}

func (f *Exam) GetSubmittedExam(examinee string, takenExamId string) (*models.SubmittedExams, error) {
	nilSubmit := new(models.SubmittedExams)
	submission, err := f.Exam.GetSubmit(takenExamId)
	if err != nil {
		return submission, err
	}
	if examinee != submission.Examinee {
		return nilSubmit, errors.New("Cannot get nonpossession submission")
	}
	return submission, err
}
