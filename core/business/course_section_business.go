package business

import (
	"context"
	"main/core/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CourseSectionBusiness struct {
	DB *mongo.Database
}

func (b *CourseSectionBusiness) GetSectionsByCourse(courseId string) ([]models.CourseSection, error) {
	cursor, err := b.DB.Collection("course_sections").Find(context.TODO(), bson.M{"courseid": courseId})
	if err != nil {
		return nil, err
	}
	sections := make([]models.CourseSection, 0)
	for cursor.Next(context.TODO()) {
		section := new(models.CourseSection)
		err = cursor.Decode(section)
		if err != nil {
			return nil, err
		}
		sections = append(sections, *section)
	}
	return sections, nil
}

func (b *CourseSectionBusiness) HasSection(courseId string, sectionId string) bool {
	r := b.DB.Collection("course_sections").FindOne(context.TODO(), bson.M{"courseid": courseId, "_id": sectionId})
	return r.Err() == nil
}

func (b *CourseSectionBusiness) Create(section *models.CourseSection) error {
	_, err := b.DB.Collection("course_sections").InsertOne(context.TODO(), section)
	return err
}

func (b *CourseSectionBusiness) Update(id string, section models.CourseSection) error {
	updatedResult := b.DB.Collection("course_sections").FindOneAndUpdate(context.TODO(),
		bson.M{"_id": id},
		bson.M{"$set": section})
	if updatedResult.Err() != nil {
		return updatedResult.Err()
	}
	return nil
}

func (b *CourseSectionBusiness) Delete(id string) error {
	_, err := b.DB.Collection("course_sections").DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}
