package routers

import (
	"main/core"
	"main/core/business"
	"main/core/models"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRouter struct {
	Name string
	g    *echo.Group
}

func (r *UserRouter) Connect(s *core.Server) {
	r.g = s.Echo.Group(r.Name)

	user := business.UserBusiness{
		DB: s.DB,
	}

	account := business.AccountBusiness{
		DB: s.DB,
	}

	transaction := business.TransactionBusiness{
		DB: s.DB,
	}

	enrollment := business.EnrollmentBusiness{
		DB: s.DB,
	}

	course := business.CourseBusiness{
		DB: s.DB,
	}

	invitation := business.InvitationBusiness{
		DB: s.DB,
	}

	user.CreateIndexes()
	account.CreateIndexes()

	r.g.GET("/", func(c echo.Context) error {
		authUser := c.Get("user")
		return c.JSON(http.StatusOK, authUser)
	}, s.AuthWiddlewareJWT.Auth)

	r.g.GET("/account", func(c echo.Context) error {
		authUser := c.Get("user").(map[string]interface{})
		account, err := account.Get(authUser["email"].(string))
		if err != nil {
			return echo.ErrInternalServerError
		}
		return c.JSON(http.StatusOK, account)
	}, s.AuthWiddlewareJWT.Auth)

	r.g.GET("/transactions", func(c echo.Context) error {
		authUser := c.Get("user").(map[string]interface{})
		transactions, err := transaction.Get(authUser["email"].(string), "")
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, transactions)
	}, s.AuthWiddlewareJWT.Auth)

	r.g.GET("/courses", func(c echo.Context) error {
		authUser := c.Get("user").(map[string]interface{})
		enrollments, err := enrollment.GetByEmail(authUser["email"].(string))
		if err != nil {
			return echo.ErrInternalServerError
		}
		return c.JSON(http.StatusOK, enrollments)
	}, s.AuthWiddlewareJWT.Auth)

	r.g.POST("/", func(c echo.Context) (err error) {
		authUser := c.Get("user").(map[string]interface{})
		profile := new(models.UserProfile)

		if err = c.Bind(profile); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err = c.Validate(profile); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if profile.Email != authUser["email"] {
			return echo.NewHTTPError(http.StatusBadRequest, "Input email and authorized email did not match")
		}

		err = user.Create(s.Config.AdminEmail, *profile)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		err = account.Create(profile.Email)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, profile)
	}, s.AuthWiddlewareJWT.Auth)

	r.g.POST("/enroll", func(c echo.Context) error {
		authUser := c.Get("user").(map[string]interface{})
		courseId := c.QueryParam("course")
		enrollingCourse, err := course.GetOneById(courseId)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if !enrollingCourse.AllowEnroll {
			return echo.NewHTTPError(http.StatusBadRequest, "Cannot enroll the course")
		}
		if enrollingCourse.StartDate > time.Now().Unix() {
			return echo.NewHTTPError(http.StatusBadRequest, "The course does not start")
		}

		if !enrollingCourse.IsPublic && !invitation.IsInvited(&models.StudentInvitation{
			CourseId: enrollingCourse.Id,
			Email:    authUser["email"].(string),
		}) {
			return echo.NewHTTPError(http.StatusBadRequest, "Cannot access the course")
		}

		err = enrollment.Enroll(&models.Enrollment{
			Id:       primitive.NewObjectID(),
			CourseID: courseId,
			Email:    authUser["email"].(string),
			Date:     time.Now().Unix(),
		})

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err = transaction.Purchase(&account, &models.Transaction{
			Email:     authUser["email"].(string),
			SKU:       courseId,
			Quantity:  1,
			Amount:    -enrollingCourse.Fee,
			Timestamp: time.Now().Unix(),
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		return c.NoContent(http.StatusOK)
	}, s.AuthWiddlewareJWT.Auth)

	// [PUT] Input: Body (UserProfileUpdated)

	r.g.PUT("/", func(c echo.Context) (err error) {
		authUser := c.Get("user").(map[string]interface{})
		updatedProfile := new(models.UpdatedUserProfile)
		if err = c.Bind(updatedProfile); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err = c.Validate(updatedProfile); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		err = user.Update(authUser["email"].(string), *updatedProfile)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}, s.AuthWiddlewareJWT.Auth)

	// [DELETE] Input user's id

	r.g.DELETE("/", func(c echo.Context) (err error) {
		authUser := c.Get("user").(map[string]interface{})
		err = user.Delete(authUser["email"].(string))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}, s.AuthWiddlewareJWT.Auth)

}
