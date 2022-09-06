package storage

import (
	"context"
	"fmt"
	"time"
)

type KillStore interface {
	KillGetter
	KillPutter
	KillDeleter
	KillLister
	KillCloser
}

type KillGetter interface {
	GetKillByID(ctx context.Context, id string) (*Kill, error)
}

type KillPutter interface {
	CreateKill(ctx context.Context, kill *Kill) (*Kill, error)
}

type KillDeleter interface {
	DeleteKill(ctx context.Context, id string) error
	DeleteKillsForServer(ctx context.Context, serverId string) error
}

type KillLister interface {
	ListKillsForServer(ctx context.Context, serverId string) ([]*Kill, error)
	ListPlayerKillsForServer(ctx context.Context, serverId string, killerId string) ([]*Kill, error)
}

type KillCloser interface {
	Close()
}

type Kill struct {
	ID       string    `firestore:"id" csv:"-"`
	ServerID string    `firestore:"serverId" csv:"-"`
	Killer   string    `firestore:"killer"`
	Victim   string    `firestore:"victim"`
	Reason   string    `firestore:"reason"`
	Date     time.Time `firestore:"date"`
}

type ErrNotFound struct {
	Key string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("Kill not found: %s", e.Key)
}
