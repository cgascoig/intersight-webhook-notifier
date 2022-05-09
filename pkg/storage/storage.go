package storage

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
	log "github.com/sirupsen/logrus"
)

type Subscription struct {
	// IgnoredAlarmNames []string `firestore:"ignoredAlarmNames"`

	WebexRoomID    string `firestore:"webexRoomID"`
	subscriptionID string
}

type Store struct {
	fsClient *firestore.Client

	subscriptionCollectionName string
	projectID                  string
}

func NewStore(ctx context.Context) *Store {
	projectID := os.Getenv("PROJECT_ID")
	if projectID == "" {
		log.Error("PROJECT_ID not specified")
	}

	subscriptionCollectionName := os.Getenv("SUBSCRIPTION_COLLECTION_NAME")
	if subscriptionCollectionName == "" {
		log.Error("SUBSCRIPTION_COLLECTION_NAME not specified")
	}

	fsClient, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Errorf("Error creating firestore client: %v", err)
	}

	return &Store{
		fsClient: fsClient,

		subscriptionCollectionName: subscriptionCollectionName,
		projectID:                  projectID,
	}
}

func (s *Store) GetSubscription(ctx context.Context, subscriptionID string) (*Subscription, error) {
	ds, err := s.fsClient.Collection(s.subscriptionCollectionName).Doc(subscriptionID).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting document from firestore: %v", err)
	}

	if !ds.Exists() {
		return nil, fmt.Errorf("SubscriptionID does not exist in firestore")
	}

	var subscription Subscription
	if err := ds.DataTo(&subscription); err != nil {
		return nil, fmt.Errorf("Subscription data in firestore invalid: %v", err)
	}

	subscription.subscriptionID = ds.Ref.ID

	return &subscription, nil
}

func (s *Store) GetSubscriptionByWebexRoomID(ctx context.Context, roomID string) (*Subscription, error) {
	docs := s.fsClient.Collection(s.subscriptionCollectionName).Where("webexRoomID", "==", roomID).Documents(ctx)

	ds, err := docs.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error getting documents from firestore: %v", err)
	}

	if len(ds) != 1 {
		return nil, fmt.Errorf("firestore query did not return exactly 1 result")
	}

	var subscription Subscription
	if err := ds[0].DataTo(&subscription); err != nil {
		return nil, fmt.Errorf("Subscription data in firestore invalid: %v", err)
	}
	subscription.subscriptionID = ds[0].Ref.ID

	return &subscription, nil
}

func (s *Store) SaveSubscription(ctx context.Context, subscription *Subscription) error {
	log.Infof("Saving subscription '%s' for room ID '%s'", subscription.subscriptionID, subscription.WebexRoomID)

	var d *firestore.DocumentRef
	if subscription.subscriptionID == "" {
		log.Info("No existing subscriptionID, using new firestore document")
		d = s.fsClient.Collection(s.subscriptionCollectionName).NewDoc()
	} else {
		log.Info("Existing subscriptionID, using existing firestore document")
		d = s.fsClient.Collection(s.subscriptionCollectionName).Doc(subscription.subscriptionID)
	}
	_, err := d.Set(ctx, subscription)
	if err != nil {
		return fmt.Errorf("unable to save subscription: %v", err)
	}

	return nil
}

func NewSubscription(roomID string) *Subscription {
	return &Subscription{
		WebexRoomID: roomID,
	}
}
