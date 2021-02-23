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

type EnrollmentBusiness struct {
	DB *mongo.Database
}

func (b *EnrollmentBusiness) IsEnroll(courseId string, email string) bool {
	r := b.DB.Collection("enrollments").FindOne(context.TODO(), bson.M{"courseid": courseId, "email": email})
	return r.Err() == nil
}

func (b *EnrollmentBusiness) Enroll(enrollment *models.Enrollment) error {
	enrollment.Date = time.Now().Unix()
	if b.IsEnroll(enrollment.CourseID, enrollment.Email) {
		return errors.New("User enrolled")
	}
	enrollment.Id = primitive.NewObjectID()
	_, err := b.DB.Collection("enrollments").InsertOne(context.TODO(), enrollment)
	return err
}

func (b *EnrollmentBusiness) cursorToEnrollments(cursor *mongo.Cursor) (enrollments []models.Enrollment, err error) {
	enrollments = make([]models.Enrollment, 0)
	//utils.CursorToList(context.TODO(), cursor, courses) // !!! Wait for testing
	for cursor.Next(context.TODO()) {
		enrollment := new(models.Enrollment)
		err = cursor.Decode(enrollment)
		if err != nil {
			return nil, err
		}
		enrollments = append(enrollments, *enrollment)

	}
	return enrollments, nil
}
func (b *EnrollmentBusiness) GetByCourseId(courseId string) ([]models.Enrollment, error) {
	cursor, err := b.DB.Collection("enrollments").Find(context.TODO(), bson.M{"courseid": courseId})
	if err != nil {
		return nil, err
	}
	return b.cursorToEnrollments(cursor)
}

func (b *EnrollmentBusiness) GetByEmail(email string) ([]models.Enrollment, error) {
	cursor, err := b.DB.Collection("enrollments").Find(context.TODO(), bson.M{"email": email})
	if err != nil {
		return nil, err
	}
	return b.cursorToEnrollments(cursor)
}
