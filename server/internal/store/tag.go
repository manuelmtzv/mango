package store

import (
	"context"
	"errors"
	"time"

	"github.com/manuelmtzv/mangocatnotes-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoTagStore struct {
	coll *mongo.Collection
}

func NewTagStore(db *mongo.Database) *MongoTagStore {
	return &MongoTagStore{
		coll: db.Collection("tags"),
	}
}

func (s *MongoTagStore) Create(ctx context.Context, tag *models.Tag) error {
	count, err := s.coll.CountDocuments(ctx, bson.M{"userId": tag.UserID})
	if err != nil {
		return err
	}
	if count >= 50 {
		return errors.New("You can only have 50 tags.")
	}

	tag.CreatedAt = time.Now()
	tag.UpdatedAt = time.Now()
	if tag.NoteIDs == nil {
		tag.NoteIDs = []primitive.ObjectID{}
	}
	res, err := s.coll.InsertOne(ctx, tag)
	if err != nil {
		return err
	}
	tag.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (s *MongoTagStore) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Tag, error) {
	var tag models.Tag
	err := s.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&tag)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	if tag.NoteIDs == nil {
		tag.NoteIDs = []primitive.ObjectID{}
	}
	return &tag, nil
}

func (s *MongoTagStore) GetAll(ctx context.Context, userID primitive.ObjectID, page, limit int64) ([]models.Tag, int64, error) {
	skip := ((page - 1) * limit)
	l := (limit)
	opts := options.FindOptions{
		Skip:  &skip,
		Limit: &l,
		Sort:  bson.M{"updatedAt": -1},
	}

	filter := bson.M{"userId": userID}

	cursor, err := s.coll.Find(ctx, filter, &opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var tags []models.Tag
	if err = cursor.All(ctx, &tags); err != nil {
		return nil, 0, err
	}

	count, err := s.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	for i := range tags {
		if tags[i].NoteIDs == nil {
			tags[i].NoteIDs = []primitive.ObjectID{}
		}
	}

	return tags, count, nil
}

func (s *MongoTagStore) Update(ctx context.Context, tag *models.Tag) error {
	tag.UpdatedAt = time.Now()
	update := bson.M{
		"$set": bson.M{
			"name":      tag.Name,
			"color":     tag.Color,
			"updatedAt": tag.UpdatedAt,
		},
	}
	_, err := s.coll.UpdateOne(ctx, bson.M{"_id": tag.ID}, update)
	return err
}

func (s *MongoTagStore) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := s.coll.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (s *MongoTagStore) FindOrCreate(ctx context.Context, userID primitive.ObjectID, names []string) ([]models.Tag, error) {
	existingCount, err := s.coll.CountDocuments(ctx, bson.M{"userId": userID})
	if err != nil {
		return nil, err
	}

	var existingTags []models.Tag
	var toCreateNames []string

	for _, name := range names {
		var tag models.Tag
		err := s.coll.FindOne(ctx, bson.M{"userId": userID, "name": name}).Decode(&tag)
		if err == nil {
			existingTags = append(existingTags, tag)
		} else {
			toCreateNames = append(toCreateNames, name)
		}
	}

	if existingCount+int64(len(toCreateNames)) > 50 {
		return nil, errors.New("You can only have a maximum of 50 tags.")
	}

	var tags []models.Tag
	tags = append(tags, existingTags...)

	for _, name := range toCreateNames {
		newTag := models.Tag{
			UserID: userID,
			Name:   name,
			Color:  "#fff",
		}
		if err := s.Create(ctx, &newTag); err != nil {
			return nil, err
		}
		tags = append(tags, newTag)
	}

	return tags, nil
}
