package domain

import (
	"time"
)

type Process struct {
	ID          int
	RequestID   string
	Status      ProcessStatus
	DetailsID   int
	StartedAt   *time.Time
	CompletedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// master_statuses
type ProcessStatus int

const (
	ProcessStatusPending ProcessStatus = iota + 1
	ProcessStatusRunning
	ProcessStatusSucceed
	ProcessStatusFailed
)

func (p ProcessStatus) String() string {
	switch p {
	case ProcessStatusPending:
		return "pending"
	case ProcessStatusRunning:
		return "running"
	case ProcessStatusSucceed:
		return "succeed"
	case ProcessStatusFailed:
		return "failed"
	default:
		return "unknown"
	}
}

func ProcessStatusFromString(s string) (ProcessStatus, bool) {
	switch s {
	case "pending":
		return ProcessStatusPending, true
	case "running":
		return ProcessStatusRunning, true
	case "succeed":
		return ProcessStatusSucceed, true
	case "failed":
		return ProcessStatusFailed, true
	default:
		return ProcessStatusPending, false
	}
}

type ProcessDetails struct {
	ID             int
	Command        *string // raw command
	Title          string
	ScoreObjectID  int  // score object
	LogObjectID    *int // log object
	ResultObjectID *int // result object
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// master_object_types
type ObjectType int

const (
	ObjectTypeFile ObjectType = iota + 1
	ObjectTypeDir
)

type Object struct {
	ID               int
	Type             ObjectType
	Bucket           string
	Path             string
	SizeBytes        uint64
	CreatedAt        time.Time
	UpdatedAt        time.Time
	BucketPathSha256 []byte
}
