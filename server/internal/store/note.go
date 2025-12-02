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

type MongoNoteStore struct {
	coll     *mongo.Collection
	tagsColl *mongo.Collection
}

func NewNoteStore(db *mongo.Database) *MongoNoteStore {
	return &MongoNoteStore{
		coll:     db.Collection("notes"),
		tagsColl: db.Collection("tags"),
	}
}

func (s *MongoNoteStore) Create(ctx context.Context, note *models.Note) error {
	note.CreatedAt = time.Now()
	note.UpdatedAt = time.Now()
	if note.TagIDs == nil {
		note.TagIDs = []primitive.ObjectID{}
	}
	res, err := s.coll.InsertOne(ctx, note)
	if err != nil {
		return err
	}
	note.ID = res.InsertedID.(primitive.ObjectID)

	if len(note.TagIDs) > 0 {
		_, err := s.tagsColl.UpdateMany(ctx,
			bson.M{"_id": bson.M{"$in": note.TagIDs}},
			bson.M{"$addToSet": bson.M{"noteIDs": note.ID}},
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *MongoNoteStore) AttachTags(ctx context.Context, noteID primitive.ObjectID, tagIDs []primitive.ObjectID) error {
	update := bson.M{
		"$addToSet": bson.M{
			"tagIDs": bson.M{
				"$each": tagIDs,
			},
		},
		"$set": bson.M{
			"updatedAt": time.Now(),
		},
	}
	_, err := s.coll.UpdateOne(ctx, bson.M{"_id": noteID}, update)
	if err != nil {
		return err
	}

	_, err = s.tagsColl.UpdateMany(ctx,
		bson.M{"_id": bson.M{"$in": tagIDs}},
		bson.M{"$addToSet": bson.M{"noteIDs": noteID}},
	)
	return err
}

func (s *MongoNoteStore) DetachTag(ctx context.Context, noteID primitive.ObjectID, tagID primitive.ObjectID) error {
	update := bson.M{
		"$pull": bson.M{
			"tagIDs": tagID,
		},
		"$set": bson.M{
			"updatedAt": time.Now(),
		},
	}
	_, err := s.coll.UpdateOne(ctx, bson.M{"_id": noteID}, update)
	if err != nil {
		return err
	}

	_, err = s.tagsColl.UpdateOne(ctx,
		bson.M{"_id": tagID},
		bson.M{"$pull": bson.M{"noteIDs": noteID}},
	)
	return err
}

func (s *MongoNoteStore) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Note, error) {
	var note models.Note
	err := s.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&note)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	if note.TagIDs == nil {
		note.TagIDs = []primitive.ObjectID{}
	}
	return &note, nil
}

func (s *MongoNoteStore) GetAll(ctx context.Context, userID primitive.ObjectID, page, limit int64, search string, tags []string) ([]models.Note, int64, error) {
	skip := ((page - 1) * limit)
	l := (limit)
	opts := options.FindOptions{
		Skip:  &skip,
		Limit: &l,
		Sort:  bson.M{"updatedAt": -1},
	}

	filter := bson.M{"userId": userID}

	if search != "" {
		filter["title"] = bson.M{"$regex": search, "$options": "i"}
	}

	if len(tags) > 0 {
		var tagDocs []models.Tag
		tagCursor, err := s.tagsColl.Find(ctx, bson.M{
			"userId": userID,
			"name":   bson.M{"$in": tags},
		})
		if err != nil {
			return nil, 0, err
		}
		if err = tagCursor.All(ctx, &tagDocs); err != nil {
			return nil, 0, err
		}
		tagCursor.Close(ctx)

		if len(tagDocs) > 0 {
			tagIDs := make([]primitive.ObjectID, len(tagDocs))
			for i, tag := range tagDocs {
				tagIDs[i] = tag.ID
			}
			filter["tagIDs"] = bson.M{"$in": tagIDs}
		} else {
			return []models.Note{}, 0, nil
		}
	}

	cursor, err := s.coll.Find(ctx, filter, &opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var notes []models.Note
	if err = cursor.All(ctx, &notes); err != nil {
		return nil, 0, err
	}

	count, err := s.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	for i := range notes {
		if notes[i].TagIDs == nil {
			notes[i].TagIDs = []primitive.ObjectID{}
		}
	}

	return notes, count, nil
}

func (s *MongoNoteStore) Update(ctx context.Context, note *models.Note) error {
	note.UpdatedAt = time.Now()
	update := bson.M{
		"$set": bson.M{
			"title":     note.Title,
			"content":   note.Content,
			"archived":  note.Archived,
			"tagIDs":    note.TagIDs,
			"updatedAt": note.UpdatedAt,
		},
	}
	_, err := s.coll.UpdateOne(ctx, bson.M{"_id": note.ID}, update)
	return err
}

func (s *MongoNoteStore) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := s.coll.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (s *MongoNoteStore) GetTags(ctx context.Context, noteID primitive.ObjectID) ([]models.Tag, error) {

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"_id": noteID}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "tags",
			"localField":   "tagIDs",
			"foreignField": "_id",
			"as":           "tags",
		}}},
		{{Key: "$project", Value: bson.M{"tags": 1, "_id": 0}}},
	}

	cursor, err := s.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []struct {
		Tags []models.Tag `bson:"tags"`
	}

	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return []models.Tag{}, nil
	}

	if result[0].Tags == nil {
		return []models.Tag{}, nil
	}

	return result[0].Tags, nil
}
