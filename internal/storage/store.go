package storage

import "context"

type Store interface {
	SaveNote(ctx context.Context, userID int64, text string) error
	GetNotes(ctx context.Context, userID int64) ([]string, error)
	CountNotes(ctx context.Context, userID int64) (int, error)
	Close() error
}
