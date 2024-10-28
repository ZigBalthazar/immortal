package repositories

import (
	"context"

	"github.com/dezh-tech/immortal/management/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepositoryQ interface {
	GetUserByEmail(user *models.User, email string)
}

type UserRepository struct {
	DB *mongo.Client
}

func NewUserRepository(db *mongo.Client) *UserRepository {
	return &UserRepository{DB: db}
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, user *models.User, email string) (error){
	usersCollection := ur.DB.Database("immortal").Collection("users")

	filter := bson.D{{Key: "email", Value: email}}


	err := usersCollection.FindOne(ctx, filter).Decode(user)
	if err != nil {
		return err
	}

	return nil
}
