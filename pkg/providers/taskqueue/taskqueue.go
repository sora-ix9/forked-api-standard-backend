package taskqueue

import (
	"fdlp-standard-api/pkg/config"
	"fmt"

	"github.com/hibiken/asynq"
)

// TaskQueue defines the interface for enqueuing tasks
type TaskQueue interface {
	Enqueue(typeName string, payload []byte, opts ...asynq.Option) (*asynq.TaskInfo, error)
	Close() error
}

type redisTaskQueue struct {
	client *asynq.Client
}

// NewTaskQueue initializes and returns a new TaskQueue client
func NewTaskQueue(cfg *config.Config) TaskQueue {
	redisAddr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)

	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     redisAddr,
		Password: cfg.RedisPassword,
		DB:       0, // Asynq recommends using a dedicated DB or the default one, but be careful with collisions
	})

	return &redisTaskQueue{
		client: client,
	}
}

func (q *redisTaskQueue) Enqueue(typeName string, payload []byte, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	task := asynq.NewTask(typeName, payload)
	return q.client.Enqueue(task, opts...)
}

func (q *redisTaskQueue) Close() error {
	return q.client.Close()
}
