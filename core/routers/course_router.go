package routers

import (
	"main/core"
	"main/core/business"
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

	r.g.GET("listing", func(c echo.Context) error {
		courses, err := course.GetPublic()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, courses)
	})
}
