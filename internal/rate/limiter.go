package rate

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Limiter struct {
	client *redis.Client
}

func NewLimiter(redisURL string) (*Limiter, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		opts = &redis.Options{Addr: redisURL}
	}

	client := redis.NewClient(opts)
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &Limiter{client: client}, nil
}

func (l *Limiter) AllowRequest(ctx context.Context, clientID string) bool {
	key := "rate_limit:" + clientID
	count, err := l.client.Incr(ctx, key).Result()
	if err != nil {
		return false
	}

	if count == 1 {
		l.client.Expire(ctx, key, 60*time.Second)
	}

	return count <= 10
}

func (l *Limiter) Close() {
	l.client.Close()
}
