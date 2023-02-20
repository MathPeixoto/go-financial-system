package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (r *RedisDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	info, err := r.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("task_id", info.ID).
		Str("queue", info.Queue).
		Int("max retries", info.MaxRetry).
		Msg("task sent to queue")

	return nil
}

func (r *RedisTaskProcessor) ProcessSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	user, err := r.store.GetUser(ctx, payload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not found: %w", asynq.SkipRetry)
		}

		return fmt.Errorf("failed to get user: %w", err)
	}

	// TODO send email

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("email", user.Email).
		Msg("task received")

	return nil
}
