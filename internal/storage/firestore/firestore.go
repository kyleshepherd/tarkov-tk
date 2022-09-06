package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/kyleshepherd/discord-tk-bot/internal/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	KillStore struct {
		client *firestore.Client
	}
)

func NewKillStore(ctx context.Context, projectID string, credentialsFilePath string) (*KillStore, error) {
	conf := &firebase.Config{ProjectID: projectID}

	var app *firebase.App

	if credentialsFilePath != "" {
		var err error
		opt := option.WithCredentialsFile(credentialsFilePath)
		app, err = firebase.NewApp(ctx, conf, opt)
		if err != nil {
			return nil, err
		}
	} else {
		var err error
		app, err = firebase.NewApp(ctx, conf)
		if err != nil {
			return nil, err
		}
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	return &KillStore{
		client: client,
	}, nil
}

func (s *KillStore) GetKillByID(ctx context.Context, id string) (*storage.Kill, error) {
	dsnap, err := s.client.Collection("kills").Doc(id).Get(ctx)
	if err != nil && status.Code(err) == codes.NotFound {
		return nil, &storage.ErrNotFound{Key: id}
	}
	if err != nil {
		return nil, err
	}
	var k storage.Kill
	err = dsnap.DataTo(&k)
	if err != nil {
		return nil, err
	}
	return &k, nil
}

func (s *KillStore) CreateKill(ctx context.Context, kill *storage.Kill) (*storage.Kill, error) {
	ref := s.client.Collection("kills").NewDoc()
	kill.ID = ref.ID
	_, err := ref.Set(ctx, kill)
	if err != nil {
		return nil, err
	}
	return kill, nil
}

func (s *KillStore) DeleteKill(ctx context.Context, id string) error {
	_, err := s.client.Collection("kills").Doc(id).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *KillStore) DeleteKillsForServer(ctx context.Context, serverId string) error {
	iter := s.client.Collection("kills").Where("serverId", "==", serverId).Documents(ctx)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		_, err = doc.Ref.Delete(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *KillStore) ListKillsForServer(ctx context.Context, serverId string) ([]*storage.Kill, error) {
	iter := s.client.Collection("kills").Where("serverId", "==", serverId).OrderBy("date", firestore.Desc).Documents(ctx)
	return iterateKills(iter)
}

func (s *KillStore) ListPlayerKillsForServer(ctx context.Context, serverId string, killerId string) ([]*storage.Kill, error) {
	iter := s.client.Collection("kills").Where("serverId", "==", serverId).Where("killer", "==", killerId).OrderBy("date", firestore.Desc).Documents(ctx)
	return iterateKills(iter)
}

func iterateKills(iter *firestore.DocumentIterator) ([]*storage.Kill, error) {
	kills := []*storage.Kill{}
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var k *storage.Kill
		if err := doc.DataTo(&k); err != nil {
			return nil, err
		}
		kills = append(kills, k)
	}
	return kills, nil
}

func (s *KillStore) Close() {
	s.client.Close()
}
