package session

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/manuelmtzv/mangocatnotes-api/internal/kvstore"
)

type SessionManager struct {
	store kvstore.KVStorage
}

func NewSessionManager(store kvstore.KVStorage) *SessionManager {
	return &SessionManager{
		store: store,
	}
}

func (sm *SessionManager) CreateSession(ctx context.Context, userID uuid.UUID, ttl time.Duration) (string, error) {
	sessionID := uuid.New().String()
	if err := sm.store.Set(ctx, sessionID, userID.String(), ttl); err != nil {
		return "", err
	}
	return sessionID, nil
}

func (sm *SessionManager) GetSession(ctx context.Context, sessionID string) (uuid.UUID, error) {
	userIDStr, err := sm.store.Get(ctx, sessionID)
	if err != nil {
		return uuid.Nil, err
	}
	return uuid.Parse(userIDStr)
}

func (sm *SessionManager) DeleteSession(ctx context.Context, sessionID string) error {
	return sm.store.Del(ctx, sessionID)
}
