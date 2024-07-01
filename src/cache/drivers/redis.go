package drivers

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type Redis struct {
	Host     string
	Port     string
	Password string
	Database int
}

var (
	client *redis.Client
)

// Connect establishes a connection to Redis
func (r *Redis) Connect() (err error) {
	//connStr := fmt.Sprintf("redis://%s:%s@%s:%s/%s")

	// Create Redis client
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", r.Host, r.Port),
		Password: r.Password,
		DB:       0, //todo: should reconsider.
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS13,
		},
	})

	// Ping Redis server to ensure connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err = client.Ping(ctx).Result(); err != nil {
		log.Fatalln(err)
	}

	return
}

// Close closes the connection to Redis
func (r *Redis) Close() (err error) {
	if err = client.Close(); err != nil {
		log.Fatalln(err)
	}
	return
}

func (r *Redis) GetClient() *redis.Client {
	return client
}
