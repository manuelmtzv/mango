package store

import "go.mongodb.org/mongo-driver/mongo"

type Storage struct {
	Users UserStorage
	Notes NoteStorage
	Tags  TagStorage
}

func NewStorage(db *mongo.Database) *Storage {
	return &Storage{
		Users: NewUserStore(db),
		Notes: NewNoteStore(db),
		Tags:  NewTagStore(db),
	}
}
