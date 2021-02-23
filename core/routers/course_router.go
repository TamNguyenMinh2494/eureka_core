package routers

import (
	"main/core"
	"main/core/business"
	"main/core/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CourseRouter struct {
	Name string
	g    *echo.Group
}

func (r *CourseRouter) Connect(s *core.Server) {
	r.g = s.Echo.Group(r.Name)
	enrollment := business.EnrollmentBusiness{
		DB: s.DB,
	}

	course := business.CourseBusiness{
		DB: s.DB,
	}

	course.CreateIndexes()

	courseSection := business.CourseSectionBusiness{
		DB: s.DB,
	}
	r.g.GET("/", func(c echo.Context) (err error) {
		authUser := c.Get("user").(map[string]interface{})
		courseId := c.QueryParam("id")
		if courseId != "" {
			courses, err := course.GetOneById(courseId)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}
			return c.JSON(http.StatusOK, courses)
		} else {
			courses, err := course.GetByAuthor(authUser["email"].(string))
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}
			return c.JSON(http.StatusOK, courses)
		}
	}, s.AuthWiddlewareJWT.Auth)

	r.g.GET("/listing", func(c echo.Context) error {
		courses, err := course.GetPublic()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, courses)
	})

	r.g.GET("/sections", func(c echo.Context) error {
		courseId := c.QueryParam("course")
		sections, err := courseSection.GetSectionsByCourse(courseId)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, sections)
	})

	r.g.GET("/enrollment", func(c echo.Context) error {
		courseId := c.QueryParam("course")
		enrols, err := enrollment.GetByCourseId(courseId)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, enrols)

	})

	r.g.POST("/", func(c echo.Context) (err error) {
		authUser := c.Get("user").(map[string]interface{})
		newCourse := new(models.Course)

		if err = c.Bind(newCourse); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err = c.Validate(newCourse); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		newCourse.AuthorEmail = authUser["email"].(string)
		err = course.Create(*newCourse)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}, s.AuthWiddlewareJWT.Auth)

	r.g.POST("/sections", func(c echo.Context) (err error) {
		authUser := c.Get("user").(map[string]interface{})
		section := new(models.CourseSection)

		if err = c.Bind(section); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err = c.Validate(section); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if !course.IsAuthor(section.CourseId, authUser["email"].(string)) {
			return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{"message": "Cannot modify nonpossession course"})
		}

		err = courseSection.Create(section)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.NoContent(http.StatusOK)
	})

	r.g.PUT("/", func(c echo.Context) (err error) {
		authUser := c.Get("user").(map[string]interface{})
		updatedCourse := new(models.Course)

		if err = c.Bind(updatedCourse); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err = c.Validate(updatedCourse); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if !course.IsAuthor(updatedCourse.Id, authUser["email"].(string)) {
			return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{"message": "Cannot modify nonpossession course"})
		}

		err = course.Update(updatedCourse.Id, *updatedCourse)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}, s.AuthWiddlewareJWT.Auth)

	r.g.PUT("/sections", func(c echo.Context) (err error) {
		authUser := c.Get("user").(map[string]interface{})
		section := new(models.CourseSection)

		if err = c.Bind(section); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err = c.Validate(section); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if !course.IsAuthor(section.CourseId, authUser["email"].(string)) {
			return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{"message": "Cannot modify nonpossession course"})
		}
		err = courseSection.Update(section.CourseId, *section)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.NoContent(http.StatusOK)
	})

	r.g.DELETE("/", func(c echo.Context) (err error) {
		authUser := c.Get("user").(map[string]interface{})
		courseId := c.QueryParam("id")

		if !course.IsAuthor(courseId, authUser["email"].(string)) {
			return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{"message": "Cannot delete nonpossession course"})
		}

		err = course.Delete(courseId)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}, s.AuthWiddlewareJWT.Auth)

	r.g.DELETE("/sections", func(c echo.Context) (err error) {
		authUser := c.Get("user").(map[string]interface{})
		courseId := c.QueryParam("course")
		sectionId := c.QueryParam("section")

		if !courseSection.HasSection(sectionId, courseId) {
			return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{"message": "The course does not contain this section"})
		}

		selectedCourse, err := course.GetOneById(courseId)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if !course.IsAuthor(selectedCourse.Id, authUser["email"].(string)) {
			return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{"message": "Cannot delete nonpossession course"})
		}

		err = courseSection.Delete(sectionId)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}, s.AuthWiddlewareJWT.Auth)

}
