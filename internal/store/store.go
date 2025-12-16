package store

import "github.com/jackc/pgx/v5/pgxpool"

type Storage struct {
	Users UserStorage
	Notes NoteStorage
	Tags  TagStorage
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{
		Users: NewUserStore(pool),
		Notes: NewNoteStore(pool),
		Tags:  NewTagStore(pool),
	}
}
