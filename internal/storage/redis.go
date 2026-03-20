package storage

// TODO: Redis-based implementation
//
// import (
// 	"context"
// 	"fmt"
// 	"github.com/redis/go-redis/v9"
// )
//
// type RedisStore struct {
// 	client *redis.Client
// }
//
// func NewRedisStore(addr string) (*RedisStore, error) {
// 	client := redis.NewClient(&redis.Options{Addr: addr})
// 	if err := client.Ping(context.Background()).Err(); err != nil {
// 		return nil, fmt.Errorf("redis connect: %w", err)
// 	}
// 	return &RedisStore{client: client}, nil
// }
//
// Key schema:
//   notes:{userID} — Redis List of note strings
//
// SaveNote:  RPUSH notes:{userID} text
// GetNotes:  LRANGE notes:{userID} 0 -1
// CountNotes: LLEN notes:{userID}
