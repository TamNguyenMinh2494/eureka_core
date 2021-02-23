package business

import (
	"context"
	"errors"
	"main/core/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type InvitationBusiness struct {
	DB *mongo.Database
}

func (b *InvitationBusiness) IsInvited(invitation *models.StudentInvitation) bool {
	r := b.DB.Collection("course_invitations").FindOne(context.TODO(), bson.M{"courseid": invitation.CourseId, "email": invitation.Email})
	return r.Err() == nil
}

func (b *InvitationBusiness) InviteStudent(invitation *models.StudentInvitation) error {
	if b.IsInvited(invitation) {
		return errors.New("Invited")
	}
	_, err := b.DB.Collection("course_invitations").InsertOne(context.TODO(), invitation)
	return err
}
