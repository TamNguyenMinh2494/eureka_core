package business

import (
	"context"
	"main/core/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CourseBusiness struct {
	DB *mongo.Database
}

func (b *CourseBusiness) CreateIndexes() {
	b.DB.Collection("courses").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.M{"id": 1},
		Options: options.Index().SetUnique(true),
	})
}

func (b *CourseBusiness) GetOneById(id string) (models.Course, error) {
	course := new(models.Course)
	result := b.DB.Collection("courses").FindOne(context.TODO(), bson.M{"id": id})
	if result.Err() != nil {
		return *course, result.Err()
	}
	err := result.Decode(course)
	if err != nil {
		return *course, err
	}
	return *course, nil
}

func (b *CourseBusiness) IsAuthor(courseId string, email string) bool {
	r := b.DB.Collection("courses").FindOne(context.TODO(), bson.M{"id": courseId, "authoremail": email})
	return r.Err() == nil
}

func (b *CourseBusiness) GetByAuthor(email string) ([]models.Course, error) {
	cursor, err := b.DB.Collection("courses").Find(context.TODO(), bson.M{"authoremail": email})
	if err != nil {
		return nil, err
	}
	courses := make([]models.Course, 0)
	//utils.CursorToList(context.TODO(), cursor, courses) // !!! Wait for testing
	for cursor.Next(context.TODO()) {
		course := new(models.Course)
		err = cursor.Decode(course)
		if err != nil {
			return nil, err
		}
		courses = append(courses, *course)
	}
	return courses, nil
}

func (b *CourseBusiness) GetPublic() ([]models.Course, error) {
	cursor, err := b.DB.Collection("courses").Find(context.TODO(), bson.M{"ispublic": true})
	if err != nil {
		return nil, err
	}
	courses := make([]models.Course, 0)
	//utils.CursorToList(context.TODO(), cursor, courses) // !!! Wait for testing
	for cursor.Next(context.TODO()) {
		course := new(models.Course)
		err = cursor.Decode(course)
		if err != nil {
			return nil, err
		}
		courses = append(courses, *course)
	}
	return courses, nil
}

func (b *CourseBusiness) Create(course models.Course) error {
	_, err := b.DB.Collection("courses").InsertOne(context.TODO(), course)
	if err != nil {
		return err
	}
	return nil
}

func (b *CourseBusiness) Update(id string, email string, course models.Course) error {
	course.AuthorEmail = email
	updatedResult := b.DB.Collection("courses").FindOneAndUpdate(context.TODO(),
		bson.M{"id": id},
		bson.M{"$set": course})
	if updatedResult.Err() != nil {
		return updatedResult.Err()
	}
	return nil
}

func (b *CourseBusiness) Delete(id string) error {
	_, err := b.DB.Collection("courses").DeleteOne(context.TODO(), bson.M{"id": id})
	if err != nil {
		return err
	}
	return nil
}
