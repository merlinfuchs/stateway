package batcher

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/merlinfuchs/stateway/stateway-audit/model"
	"github.com/merlinfuchs/stateway/stateway-audit/store"
	"github.com/nats-io/nats.go/jetstream"
)

var _ Batcher = (*JetStreamBatcher)(nil)

const (
	// EntityChangesStreamName is the name of the JetStream stream for entity changes
	EntityChangesStreamName = "AUDIT_ENTITY_CHANGES"
	// EntityChangesStreamSubject is the subject pattern for entity changes
	EntityChangesStreamSubject = "audit.entity.changes"
	// EntityChangesConsumerName is the name of the durable consumer for batching
	EntityChangesConsumerName = "AUDIT_ENTITY_CHANGES_BATCHER"
)

// JetStreamBatcher uses NATS JetStream as a stateful batcher for entity changes.
// It consumes messages from a JetStream stream in batches and commits them to the database.
type JetStreamBatcher struct {
	js     jetstream.JetStream
	store  store.EntityChangeStore
	config JetStreamBatcherConfig

	streamName   string
	consumerName string
}

// JetStreamBatcherConfig holds the configuration for the JetStream batcher.
type JetStreamBatcherConfig struct {
	// NamePrefix is the prefix of the JetStream stream and consumer names. Defaults to "AUDIT_".
	NamePrefix string

	// BatchSize is the maximum number of messages to fetch in a single batch.
	// Defaults to 1000 if not set.
	BatchSize int

	// BatchTimeout is the maximum time to wait before processing a batch, even if
	// it hasn't reached BatchSize. Defaults to 5 seconds if not set.
	BatchTimeout time.Duration

	// MaxAckPending is the maximum number of unacknowledged messages that can be
	// in flight. This helps control memory usage. Defaults to 10000 if not set.
	MaxAckPending int
}

// NewJetStreamBatcher creates a new JetStream batcher with the given JetStream context,
// store, and configuration.
func NewJetStreamBatcher(js jetstream.JetStream, store store.EntityChangeStore, config JetStreamBatcherConfig) *JetStreamBatcher {
	// Set defaults
	if config.BatchSize == 0 {
		config.BatchSize = 1000
	}
	if config.BatchTimeout == 0 {
		config.BatchTimeout = 5 * time.Second
	}
	if config.MaxAckPending == 0 {
		config.MaxAckPending = 10000
	}

	return &JetStreamBatcher{
		js:     js,
		store:  store,
		config: config,

		streamName:   config.NamePrefix + EntityChangesStreamName,
		consumerName: config.NamePrefix + EntityChangesConsumerName,
	}
}

// CreateStream creates or updates the JetStream stream for entity changes.
func (b *JetStreamBatcher) CreateStream(ctx context.Context) error {
	_, err := b.js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name: b.streamName,
		Subjects: []string{
			fmt.Sprintf("%s.>", EntityChangesStreamSubject),
		},
		Retention: jetstream.InterestPolicy, // Interest policy ensures messages are discarded when all consumers have acknowledged them
		MaxAge:    24 * time.Hour,           // Keep messages for 24 hours
		MaxBytes:  32 * 1024 * 1024 * 1024,  // 32GB max
		MaxMsgs:   -1,                       // No limit on message count
		Discard:   jetstream.DiscardOld,
		Storage:   jetstream.FileStorage,
		Replicas:  1,
	})
	if err != nil {
		return fmt.Errorf("failed to create or update stream %s: %w", b.streamName, err)
	}

	return nil
}

// Start starts the batcher. It will consume messages from the JetStream stream
// in batches and commit them to the database. The batcher will stop when the
// context is cancelled.
func (b *JetStreamBatcher) Start(ctx context.Context) error {
	// Create or get the consumer
	consumer, err := b.js.CreateOrUpdateConsumer(ctx, b.streamName, jetstream.ConsumerConfig{
		Name:          b.consumerName,
		Durable:       b.consumerName,
		AckPolicy:     jetstream.AckAllPolicy, // Ack all messages up to and including the one being acked
		MaxAckPending: b.config.MaxAckPending,
	})
	if err != nil {
		return fmt.Errorf("failed to create or update consumer: %w", err)
	}

	// Start the batch processing loop
	go b.processBatches(ctx, consumer)

	return nil
}

// processBatches continuously fetches batches of messages and processes them.
func (b *JetStreamBatcher) processBatches(ctx context.Context, consumer jetstream.Consumer) {
	ticker := time.NewTicker(b.config.BatchTimeout)
	defer ticker.Stop()

	var pendingMessages []jetstream.Msg

	for {
		select {
		case <-ctx.Done():
			// Process any remaining messages before shutdown
			if len(pendingMessages) > 0 {
				b.processBatch(ctx, pendingMessages)
			}
			return

		case <-ticker.C:
			// Timeout reached, process current batch if any
			if len(pendingMessages) > 0 {
				b.processBatch(ctx, pendingMessages)
				pendingMessages = nil
			}

		default:
			// Fetch messages up to BatchSize
			msgs, err := consumer.Fetch(b.config.BatchSize, jetstream.FetchMaxWait(b.config.BatchTimeout))
			if err != nil {
				if err == jetstream.ErrNoMessages {
					// No messages available, continue
					time.Sleep(100 * time.Millisecond)
					continue
				}

				slog.Error("Failed to fetch messages", slog.Any("error", err))
				time.Sleep(1 * time.Second)
				continue
			}

			// Collect messages
			for msg := range msgs.Messages() {
				pendingMessages = append(pendingMessages, msg)

				// If we've reached the batch size, process immediately
				if len(pendingMessages) >= b.config.BatchSize {
					b.processBatch(ctx, pendingMessages)
					pendingMessages = nil
					ticker.Reset(b.config.BatchTimeout) // Reset timer
				}
			}
		}
	}
}

// processBatch processes a batch of messages and acknowledges them.
func (b *JetStreamBatcher) processBatch(ctx context.Context, msgs []jetstream.Msg) {
	if len(msgs) == 0 {
		return
	}

	slog.Debug("Processing batch of changes", slog.Int("count", len(msgs)))

	// Unmarshal messages into entity changes
	entityChanges := make([]model.EntityChange, 0, len(msgs))
	validMessages := make([]jetstream.Msg, 0, len(msgs))

	for _, msg := range msgs {
		var change model.EntityChange
		if err := json.Unmarshal(msg.Data(), &change); err != nil {
			// Log error and nack the message individually
			slog.Error("Failed to unmarshal entity change", slog.Any("error", err))
			// Nack with delay to retry later
			if err := msg.NakWithDelay(5 * time.Second); err != nil {
				slog.Error("Failed to nack message", slog.Any("error", err))
			}
			continue
		}
		change.IngestedAt = time.Now().UTC()
		entityChanges = append(entityChanges, change)
		validMessages = append(validMessages, msg)
	}

	if len(entityChanges) == 0 {
		// All messages failed to unmarshal, nothing to ack
		return
	}

	// Insert batch into database
	flushCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 30*time.Second)
	defer cancel()

	err := b.store.InsertEntityChanges(flushCtx, entityChanges...)
	if err != nil {
		// Log error and nack all valid messages to retry
		slog.Error("Failed to insert entity changes batch", slog.Any("error", err))

		// Nack all valid messages with delay
		for _, msg := range validMessages {
			if err := msg.NakWithDelay(5 * time.Second); err != nil {
				slog.Error("Failed to nack message", slog.Any("error", err))
			}
		}
		return
	}

	// Successfully inserted, ack only the last message (AckAllPolicy will ack all previous messages)
	if len(validMessages) > 0 {
		lastMsg := validMessages[len(validMessages)-1]
		if err := lastMsg.Ack(); err != nil {
			slog.Error("Failed to ack last message", slog.Any("error", err))
		}
	}
}

// PublishEntityChange publishes an entity change to the JetStream stream.
// This can be called from anywhere in your application to queue entity changes for batching.
func (b *JetStreamBatcher) Push(ctx context.Context, change model.EntityChange) error {
	data, err := json.Marshal(change)
	if err != nil {
		return fmt.Errorf("failed to marshal entity change: %w", err)
	}

	subject := fmt.Sprintf("%s.%d.%d.%s.%s",
		EntityChangesStreamSubject,
		change.AppID,
		change.GuildID,
		change.EntityType,
		change.EntityID,
	)

	_, err = b.js.PublishAsync(subject, data)
	if err != nil {
		return fmt.Errorf("failed to publish entity change: %w", err)
	}

	return nil
}
