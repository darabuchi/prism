package queue

import "errors"

var (
	// ErrUnsupportedQueueType indicates an unsupported queue type
	ErrUnsupportedQueueType = errors.New("unsupported queue type")

	// ErrQueueClosed indicates the queue is closed
	ErrQueueClosed = errors.New("queue is closed")

	// ErrTopicNotFound indicates the topic does not exist
	ErrTopicNotFound = errors.New("topic not found")

	// ErrMessageNotFound indicates the message does not exist
	ErrMessageNotFound = errors.New("message not found")

	// ErrTimeout indicates an operation timeout
	ErrTimeout = errors.New("operation timeout")

	// ErrMaxRetriesExceeded indicates maximum retries exceeded
	ErrMaxRetriesExceeded = errors.New("maximum retries exceeded")

	// ErrInvalidConfig indicates invalid configuration
	ErrInvalidConfig = errors.New("invalid configuration")

	// ErrConnectionFailed indicates connection failure
	ErrConnectionFailed = errors.New("connection failed")
)
