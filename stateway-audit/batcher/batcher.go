package batcher

import (
	"context"
	"fmt"
	"time"

	"github.com/friendlycaptcha/batchman"
	"github.com/merlinfuchs/stateway/stateway-audit/model"
	"github.com/merlinfuchs/stateway/stateway-audit/store"
)

// EntityChangeBatcher batches entity changes and commits them to the database
// when a certain number of changes are reached or a timeout occurs.
type EntityChangeBatcher struct {
	batcher *batchman.Batcher[model.EntityChange]
	store   store.EntityChangeStore
	config  Config
}

// Config holds the configuration for the entity change batcher.
type Config struct {
	// MaxSize is the maximum number of entity changes to batch before committing.
	// Defaults to 1000 if not set.
	MaxSize int

	// MaxDelay is the maximum time to wait before committing a batch, even if
	// it hasn't reached MaxSize. Defaults to 5 seconds if not set.
	MaxDelay time.Duration

	// BufferSize is the size of the internal buffer. If the buffer is full,
	// Push will return an error. Defaults to 10000 if not set.
	BufferSize int

	// OnError is an optional callback that will be called when an error occurs
	// during batch insertion. If not set, errors are silently ignored.
	OnError func(err error, batchSize int)
}

// New creates a new EntityChangeBatcher with the given store and configuration.
func New(store store.EntityChangeStore, config Config) *EntityChangeBatcher {
	// Set defaults
	if config.MaxSize == 0 {
		config.MaxSize = 1000
	}
	if config.MaxDelay == 0 {
		config.MaxDelay = 5 * time.Second
	}
	if config.BufferSize == 0 {
		config.BufferSize = 10000
	}

	return &EntityChangeBatcher{
		store:  store,
		config: config,
	}
}

// Start starts the batcher with the given context. The batcher will stop when
// the context is cancelled. Returns the batcher instance and any error that
// occurred during initialization.
func (b *EntityChangeBatcher) Start(ctx context.Context) (*batchman.Batcher[model.EntityChange], error) {
	// Create the flush function that will be called when a batch is ready
	flush := func(ctx context.Context, items []model.EntityChange) {
		if len(items) == 0 {
			return
		}

		// Use context.WithoutCancel to allow the flush to complete even if
		// the batcher context is cancelled during shutdown
		flushCtx := context.WithoutCancel(ctx)
		flushCtx, cancel := context.WithTimeout(flushCtx, 30*time.Second)
		defer cancel()

		err := b.store.InsertEntityChanges(flushCtx, items...)
		if err != nil {
			if b.config.OnError != nil {
				b.config.OnError(fmt.Errorf("failed to insert entity changes batch: %w", err), len(items))
			}
		}
	}

	// Create and configure the batchman batcher
	init := batchman.New[model.EntityChange]().
		MaxSize(b.config.MaxSize).
		MaxDelay(b.config.MaxDelay).
		BufferSize(b.config.BufferSize)

	// Start the batcher
	batcher, err := init.Start(ctx, flush)
	if err != nil {
		return nil, fmt.Errorf("failed to start entity change batcher: %w", err)
	}

	b.batcher = batcher
	return batcher, nil
}

// Push adds an entity change to the batcher. This is a non-blocking call.
// Returns an error if the batcher has been stopped or if the buffer is full.
func (b *EntityChangeBatcher) Push(change model.EntityChange) error {
	if b.batcher == nil {
		return fmt.Errorf("batcher has not been started")
	}
	return b.batcher.Push(change)
}

// Done returns a channel that will be closed when the batcher has finished
// flushing all remaining items. This is useful for graceful shutdown.
func (b *EntityChangeBatcher) Done() <-chan struct{} {
	if b.batcher == nil {
		// Return a closed channel if batcher hasn't been started
		done := make(chan struct{})
		close(done)
		return done
	}
	return b.batcher.Done()
}
