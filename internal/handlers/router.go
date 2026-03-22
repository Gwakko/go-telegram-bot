package handlers

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/gwakko/go-telegram-bot/internal/bot"
	"github.com/gwakko/go-telegram-bot/internal/storage"
)

type Router struct {
	bot   *bot.Bot
	store storage.Store
	cmds  map[string]CommandFunc
}

type CommandFunc func(ctx context.Context, msg *bot.Message) error

func NewRouter(b *bot.Bot, store storage.Store) *Router {
	r := &Router{
		bot:   b,
		store: store,
		cmds:  make(map[string]CommandFunc),
	}

	r.cmds["/start"] = r.handleStart
	r.cmds["/help"] = r.handleHelp
	r.cmds["/note"] = r.handleNote
	r.cmds["/list"] = r.handleList
	r.cmds["/stats"] = r.handleStats

	return r
}

// Handle dispatches updates to the appropriate command handler.
func (r *Router) Handle(ctx context.Context, update bot.Update) error {
	if update.Message == nil {
		return nil
	}

	msg := update.Message
	cmd := extractCommand(msg.Text)

	if handler, ok := r.cmds[cmd]; ok {
		return handler(ctx, msg)
	}

	// Echo unknown messages
	return r.bot.SendMessage(ctx, msg.Chat.ID, "Unknown command. Use /help to see available commands.")
}

func (r *Router) handleStart(ctx context.Context, msg *bot.Message) error {
	text := "👋 Welcome! I'm a note-taking bot.\n\n" +
		"Use /help to see available commands."
	return r.bot.SendMessage(ctx, msg.Chat.ID, text)
}

func (r *Router) handleHelp(ctx context.Context, msg *bot.Message) error {
	text := "<b>Available commands:</b>\n\n" +
		"/note &lt;text&gt; — Save a note\n" +
		"/list — Show your notes\n" +
		"/stats — Your usage stats\n" +
		"/help — This message"
	return r.bot.SendMessage(ctx, msg.Chat.ID, text)
}

func (r *Router) handleNote(ctx context.Context, msg *bot.Message) error {
	text := extractArgs(msg.Text)
	if text == "" {
		return r.bot.SendMessage(ctx, msg.Chat.ID, "Usage: /note <your text here>")
	}

	if err := r.store.SaveNote(ctx, msg.From.ID, text); err != nil {
		log.Printf("save note error: %v", err)
		return r.bot.SendMessage(ctx, msg.Chat.ID, "Failed to save note.")
	}

	return r.bot.SendMessage(ctx, msg.Chat.ID, "✅ Note saved!")
}

func (r *Router) handleList(ctx context.Context, msg *bot.Message) error {
	notes, err := r.store.GetNotes(ctx, msg.From.ID)
	if err != nil {
		log.Printf("get notes error: %v", err)
		return r.bot.SendMessage(ctx, msg.Chat.ID, "Failed to load notes.")
	}

	if len(notes) == 0 {
		return r.bot.SendMessage(ctx, msg.Chat.ID, "No notes yet. Use /note to add one.")
	}

	var sb strings.Builder
	sb.WriteString("<b>Your notes:</b>\n\n")
	for i, note := range notes {
		sb.WriteString(strings.Repeat(" ", 0))
		sb.WriteString(string(rune('1'+i)) + ". " + note + "\n")
	}

	return r.bot.SendMessage(ctx, msg.Chat.ID, sb.String())
}

func (r *Router) handleStats(ctx context.Context, msg *bot.Message) error {
	count, err := r.store.CountNotes(ctx, msg.From.ID)
	if err != nil {
		return r.bot.SendMessage(ctx, msg.Chat.ID, "Failed to load stats.")
	}

	return r.bot.SendMessage(ctx, msg.Chat.ID,
		fmt.Sprintf("📊 <b>Your stats:</b>\n\nNotes saved: %d", count))
}

func extractCommand(text string) string {
	parts := strings.Fields(text)
	if len(parts) == 0 {
		return ""
	}
	cmd := parts[0]
	if at := strings.Index(cmd, "@"); at != -1 {
		cmd = cmd[:at]
	}
	return strings.ToLower(cmd)
}

func extractArgs(text string) string {
	parts := strings.SplitN(text, " ", 2)
	if len(parts) < 2 {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
