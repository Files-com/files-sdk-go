// Package events provides event types and publishing interfaces for filesystem mount operations.
//
// This package defines a common event system for tracking and responding to mount-related
// events such as authentication failures. Events implement the MountEvent interface and can
// be published through an EventPublisher implementation.
package events

// MountEvent is the marker interface that all mount events implement.
type MountEvent interface {
	isMountEvent()
}

// EventPublisher defines the interface for publishing mount events.
type EventPublisher interface {
	Publish(MountEvent)
}

// NoOpEventPublisher is an EventPublisher that does nothing.
type NoOpEventPublisher struct{}

// Publish implements the EventPublisher interface but performs no action.
func (p *NoOpEventPublisher) Publish(ev MountEvent) {}

// AuthenticationFailedEvent represents an authentication failure during mount operations.
type AuthenticationFailedEvent struct {
	Reason string
}

func (e AuthenticationFailedEvent) isMountEvent() {}
