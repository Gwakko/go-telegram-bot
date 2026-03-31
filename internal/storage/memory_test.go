package storage

import (
	"context"
	"testing"
)

func TestSaveNoteAndGetNotes(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	if err := store.SaveNote(ctx, 1, "first"); err != nil {
		t.Fatalf("SaveNote: %v", err)
	}
	if err := store.SaveNote(ctx, 1, "second"); err != nil {
		t.Fatalf("SaveNote: %v", err)
	}

	notes, err := store.GetNotes(ctx, 1)
	if err != nil {
		t.Fatalf("GetNotes: %v", err)
	}
	if len(notes) != 2 {
		t.Fatalf("expected 2 notes, got %d", len(notes))
	}
	if notes[0] != "first" || notes[1] != "second" {
		t.Fatalf("unexpected notes: %v", notes)
	}
}

func TestCountNotes(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	count, err := store.CountNotes(ctx, 1)
	if err != nil {
		t.Fatalf("CountNotes: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected 0, got %d", count)
	}

	store.SaveNote(ctx, 1, "a")
	store.SaveNote(ctx, 1, "b")
	store.SaveNote(ctx, 1, "c")

	count, err = store.CountNotes(ctx, 1)
	if err != nil {
		t.Fatalf("CountNotes: %v", err)
	}
	if count != 3 {
		t.Fatalf("expected 3, got %d", count)
	}
}

func TestSeparateUsers(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	store.SaveNote(ctx, 1, "user1-note")
	store.SaveNote(ctx, 2, "user2-note")

	notes1, _ := store.GetNotes(ctx, 1)
	notes2, _ := store.GetNotes(ctx, 2)

	if len(notes1) != 1 || notes1[0] != "user1-note" {
		t.Fatalf("user 1 notes wrong: %v", notes1)
	}
	if len(notes2) != 1 || notes2[0] != "user2-note" {
		t.Fatalf("user 2 notes wrong: %v", notes2)
	}
}

func TestEmptyUserReturnsEmptySlice(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	notes, err := store.GetNotes(ctx, 999)
	if err != nil {
		t.Fatalf("GetNotes: %v", err)
	}
	if notes != nil && len(notes) != 0 {
		t.Fatalf("expected empty/nil slice, got %v", notes)
	}
}
