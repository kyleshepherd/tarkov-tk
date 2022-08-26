package firestore

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/kyleshepherd/discord-tk-bot/internal/storage"
)

type (
	KillStore struct {
		client *firestore.Client
		now    func() time.Time
	}
)

func NewKillStore(ctx context.Context, projectID string) (*KillStore, error) {
	conf := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	return &KillStore{
		client: client,
		now:    time.Now,
	}, nil
}

func (s *KillStore) GetKillByID(ctx context.Context, id string) (*storage.Kill, error) {
	return nil, nil
}

func (s *KillStore) CreateKill(ctx context.Context, kill *storage.Kill) (*storage.Kill, error) {
	return nil, nil
}

func (s *KillStore) DeleteKill(ctx context.Context, id string) error {
	return nil
}

func (s *KillStore) DeleteKillsForServer(ctx context.Context, serverId string) error {
	return nil
}

func (s *KillStore) ListKillsForServer(ctx context.Context, serverId string) ([]*storage.Kill, error) {
	return nil, nil
}

func (s *KillStore) ListPlayerKillsForServer(ctx context.Context, serverId string, killerId string) ([]*storage.Kill, error) {
	return nil, nil
}
