// Package events provides event types and publishing interfaces for filesystem mount operations.
//
// This package defines a common event system for tracking and responding to mount-related
// events such as authentication failures. Events implement the MountEvent interface and can
// be published through an EventPublisher implementation.
package events

import "time"

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

// TransferDirection identifies whether mounted folder file content is moving
// to or from Files.com.
type TransferDirection string

const (
	TransferDirectionUpload   TransferDirection = "upload"
	TransferDirectionDownload TransferDirection = "download"
)

// TransferStatus identifies the lifecycle state of a mounted folder transfer.
type TransferStatus string

const (
	TransferStatusQueued       TransferStatus = "queued"
	TransferStatusTransferring TransferStatus = "transferring"
	TransferStatusComplete     TransferStatus = "complete"
	TransferStatusErrored      TransferStatus = "errored"
	TransferStatusCanceled     TransferStatus = "canceled"
)

// TransferEvent reports mounted folder upload and download lifecycle changes.
type TransferEvent struct {
	ID               string
	Direction        TransferDirection
	LocalPath        string
	RemotePath       string
	Size             int64
	TransferredBytes int64
	Status           TransferStatus
	StartedAt        time.Time
	EndedAt          time.Time
	Error            string
}

func (e TransferEvent) isMountEvent() {}
