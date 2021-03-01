package routers

import (
	"main/core"
	"main/core/business"
	"main/core/flow"
	"main/core/middlewares"
	"main/core/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ExamRouter struct {
	Name string
	g    *echo.Group
}

func (r *ExamRouter) Connect(s *core.Server) {
	r.g = s.Echo.Group(r.Name)
	quiz := business.QuizBusiness{
		DB: s.DB,
	}
	quizSelector := business.QuizSelectorBusiness{
		Quiz: &quiz,
	}

	exam := flow.Exam{
		Course: &business.CourseBusiness{
			DB: s.DB,
		},
		CourseSection: &business.CourseSectionBusiness{
			DB: s.DB,
		},
		Exam: &business.ExamBusiness{
			DB:           s.DB,
			Quiz:         &quiz,
			QuizSelector: &quizSelector,
		},
		Enrollment: &business.EnrollmentBusiness{
			DB: s.DB,
		},
	}
	r.g.GET("/", func(c echo.Context) error {
		userAuth := c.Get("user").(map[string]interface{})
		courseId := c.QueryParam("course")
		sectionId := c.QueryParam("section")
		exams, err := exam.GetExams(userAuth["email"].(string), courseId, sectionId)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, exams)
	}, s.AuthWiddlewareJWT.Auth)
	r.g.GET("/submission", func(c echo.Context) error {
		userAuth := c.Get("user").(map[string]interface{})
		examId := c.QueryParam("exam")
		submission, err := exam.GetSubmittedExam(userAuth["email"].(string), examId)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, submission)
	}, s.AuthWiddlewareJWT.Auth)
	r.g.GET("/taken-exams", func(c echo.Context) error {
		userAuth := c.Get("user").(map[string]interface{})
		takenExam, err := exam.Exam.GetTakenExams(userAuth["email"].(string))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, takenExam)
	}, s.AuthWiddlewareJWT.Auth)
	r.g.GET("/take", func(c echo.Context) error {
		userAuth := c.Get("user").(map[string]interface{})
		courseId := c.QueryParam("course")
		examId := c.QueryParam("exam")
		takenExam, err := exam.TakeExam(userAuth["email"].(string), courseId, examId)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, takenExam)
	}, s.AuthWiddlewareJWT.Auth)
	r.g.GET("/preview", func(c echo.Context) error {
		userAuth := c.Get("user").(map[string]interface{})
		courseId := c.QueryParam("course")
		examId := c.QueryParam("exam")
		takenExam, err := exam.PreviewExam(userAuth["email"].(string), courseId, examId)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusBadRequest, takenExam)
	}, s.AuthWiddlewareJWT.Auth)
	r.g.POST("/", func(c echo.Context) error {
		userAuth := c.Get("user").(map[string]interface{})
		courseId := c.QueryParam("course")
		exams := new(models.Exams)
		if err := c.Bind(exams); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err := c.Validate(exams); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		err := exam.CreateExam(userAuth["email"].(string), courseId, exams)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}, s.AuthWiddlewareJWT.Auth, middlewares.Deny("general"))
	r.g.POST("/submit", func(c echo.Context) error {
		userAuth := c.Get("user").(map[string]interface{})
		courseId := c.QueryParam("course")
		submission := new(models.SubmittedExams)
		if err := c.Bind(submission); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err := c.Validate(submission); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		submission.Examinee = userAuth["email"].(string)
		err := exam.SubmitExam(courseId, submission)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}, s.AuthWiddlewareJWT.Auth)
	r.g.PUT("/", func(c echo.Context) error {
		userAuth := c.Get("user").(map[string]interface{})
		courseId := c.QueryParam("course")
		exams := new(models.Exams)
		if err := c.Bind(exams); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err := c.Validate(exams); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		err := exam.UpdateExam(userAuth["email"].(string), courseId, exams)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}, s.AuthWiddlewareJWT.Auth, middlewares.Deny("general"))
	r.g.DELETE("/", func(c echo.Context) error {
		userAuth := c.Get("user").(map[string]interface{})
		courseId := c.QueryParam("course")
		examId := c.QueryParam("exam")
		err := exam.DeleteExam(userAuth["email"].(string), courseId, examId)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}, s.AuthWiddlewareJWT.Auth, middlewares.Deny("general"))
}
