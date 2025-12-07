package server

import (
	"context"
	"sync"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
)

// AuditLogInfo contains information from an audit log entry
type AuditLogInfo struct {
	ID     snowflake.ID
	UserID snowflake.ID
	Reason string
}

// auditLogKey is used as a key for matching audit log entries
type auditLogKey struct {
	TargetID   snowflake.ID
	ActionType discord.AuditLogEvent
}

// AuditLogMatcher matches entity change events with audit log events
type AuditLogMatcher struct {
	listenersMu sync.RWMutex
	listeners   map[auditLogKey]chan *AuditLogInfo
}

// NewAuditLogMatcher creates a new AuditLogMatcher
func NewAuditLogMatcher() *AuditLogMatcher {
	return &AuditLogMatcher{
		listeners: make(map[auditLogKey]chan *AuditLogInfo),
	}
}

// HandleAuditLog receives an audit log event and notifies waiting listeners
func (m *AuditLogMatcher) HandleAuditLog(auditLog gateway.EventGuildAuditLogEntryCreate) {
	// Give some time for the listener to be registered
	time.Sleep(150 * time.Millisecond)

	entry := auditLog.AuditLogEntry
	targetID := entry.TargetID
	actionType := entry.ActionType

	if targetID == nil || *targetID == 0 {
		// No target ID, can't match
		return
	}

	// Prepare audit log info
	reason := ""
	if entry.Reason != nil {
		reason = *entry.Reason
	}
	auditLogInfo := &AuditLogInfo{
		ID:     entry.ID,
		UserID: entry.UserID,
		Reason: reason,
	}

	// Check for listener waiting for this targetID and actionType
	key := auditLogKey{
		TargetID:   *targetID,
		ActionType: actionType,
	}

	m.listenersMu.Lock()
	ch, exists := m.listeners[key]
	if exists {
		// Remove listener immediately to prevent duplicate notifications
		delete(m.listeners, key)
	}
	m.listenersMu.Unlock()

	if exists {
		// Notify listener (channel is buffered, so this won't block)
		select {
		case ch <- auditLogInfo:
		default:
			// Channel already closed or full, ignore
		}
	}
}

// WaitForAuditLog waits up to 1 second for a matching audit log event
func (m *AuditLogMatcher) WaitForAuditLog(ctx context.Context, targetID snowflake.ID, actionType discord.AuditLogEvent) *AuditLogInfo {
	key := auditLogKey{
		TargetID:   targetID,
		ActionType: actionType,
	}

	// Create listener channel
	ch := make(chan *AuditLogInfo, 1)

	// Register listener
	m.listenersMu.Lock()
	m.listeners[key] = ch
	m.listenersMu.Unlock()

	// Clean up listener on return
	defer func() {
		m.listenersMu.Lock()
		// Only delete if it's still our channel (in case it was already removed by HandleAuditLog)
		if m.listeners[key] == ch {
			delete(m.listeners, key)
		}
		m.listenersMu.Unlock()
	}()

	// Wait for audit log with timeout
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		return nil
	case auditLogInfo := <-ch:
		return auditLogInfo
	}
}

// WaitForAuditLogAny waits up to 1 second for a matching audit log event with any of the provided actionTypes
// It tries all actionTypes concurrently and returns the first match
func (m *AuditLogMatcher) WaitForAuditLogAny(ctx context.Context, targetID snowflake.ID, actionTypes ...discord.AuditLogEvent) *AuditLogInfo {
	if len(actionTypes) == 0 {
		return nil
	}

	// If only one actionType, use the simpler method
	if len(actionTypes) == 1 {
		return m.WaitForAuditLog(ctx, targetID, actionTypes[0])
	}

	// Try all actionTypes concurrently
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	resultCh := make(chan *AuditLogInfo, 1)

	// Start a goroutine for each actionType
	for _, actionType := range actionTypes {
		go func(at discord.AuditLogEvent) {
			if result := m.WaitForAuditLog(ctx, targetID, at); result != nil {
				select {
				case resultCh <- result:
				default:
					// Already got a result, ignore
				}
			}
		}(actionType)
	}

	// Wait for first result or timeout
	select {
	case <-ctx.Done():
		return nil
	case result := <-resultCh:
		// Cancel remaining waits by letting context timeout
		return result
	}
}
