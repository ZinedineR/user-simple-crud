package messaging

import (
	"context"
	"user-simple-crud/internal/model"
)

type ExampleProducer interface {
	GetTopic() string
	Send(ctx context.Context, order ...*model.ExampleMessage) error
}
