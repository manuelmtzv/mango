package store

import (
	"context"
	"errors"
	"time"

	"github.com/manuelmtzv/mangocatnotes-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoUserStore struct {
	coll *mongo.Collection
}

func NewUserStore(db *mongo.Database) *MongoUserStore {
	return &MongoUserStore{
		coll: db.Collection("users"),
	}
}

func (s *MongoUserStore) Create(ctx context.Context, user *models.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	res, err := s.coll.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (s *MongoUserStore) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := s.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *MongoUserStore) GetByEmailOrUsername(ctx context.Context, identifier string) (*models.User, error) {
	var user models.User
	err := s.coll.FindOne(ctx, bson.M{
		"$or": []bson.M{
			{"email": identifier},
			{"username": identifier},
		},
	}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *MongoUserStore) GetByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := s.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *MongoUserStore) Update(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()
	update := bson.M{
		"$set": bson.M{
			"email":     user.Email,
			"username":  user.Username,
			"name":      user.Name,
			"updatedAt": user.UpdatedAt,
		},
	}
	_, err := s.coll.UpdateOne(ctx, bson.M{"_id": user.ID}, update)
	return err
}

func (s *MongoUserStore) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := s.coll.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
