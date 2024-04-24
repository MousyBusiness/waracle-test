// Package ds abstracts the GCP Datastore library to improve testability
package db

import (
	"cloud.google.com/go/datastore"
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"reflect"
)

type Database interface {
	Client() *datastore.Client
	NameKey(namespace, kind, name string, parent *datastore.Key) *datastore.Key
	CreateNamed(ctx context.Context, namespace, name string, parent *datastore.Key, entity Entity) (*datastore.Key, error)
	GetNamed(ctx context.Context, namespace, name string, parent *datastore.Key, entity Entity) error
	DeleteNamed(ctx context.Context, namespace, kind, name string, parent *datastore.Key) error
	QueryAll(ctx context.Context, namespace string, query *datastore.Query, entitySlicePtr interface{}) (keys []*datastore.Key, err error)
}

type client struct {
	ds *datastore.Client
}

type Entity interface {
	GetKind() string
	GetValue() interface{} // GetValue must return a pointer to a struct
}

// ConnectToDatastore will establish a connection to Datastore and wrap the client in a Database instance
func ConnectToDatastore(ctx context.Context) (Database, error) {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	// Creates a client.
	c, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create new datastore client: %w", err)
	}

	log.Println("datastore client created")
	return client{ds: c}, nil
}

func (client client) Client() *datastore.Client {
	return client.ds
}

func (client client) NameKey(namespace, kind, name string, parent *datastore.Key) *datastore.Key {
	k := datastore.NameKey(kind, name, parent)
	k.Namespace = namespace
	return k
}

// CreateNamed will create a single Entity using its name key and parent (if exists)
func (client client) CreateNamed(ctx context.Context, namespace, name string, parent *datastore.Key, entity Entity) (*datastore.Key, error) {
	key := client.NameKey(namespace, entity.GetKind(), name, parent) // namespace is enforced by NameKey

	k, err := client.ds.Put(ctx, key, entity.GetValue())
	if err != nil {
		return nil, fmt.Errorf("failed to create named datastore entity: %w", err)
	}

	return k, nil
}

// GetNamed will get a single Entity using its name key and parent (if exists)
func (client client) GetNamed(ctx context.Context, namespace, name string, parent *datastore.Key, entity Entity) error {
	if entity.GetValue() == nil || reflect.ValueOf(entity.GetValue()).Kind() == reflect.Ptr && reflect.ValueOf(entity.GetValue()).IsNil() {
		return errors.New("entity.GetValue cannot return nil")
	}

	key := client.NameKey(namespace, entity.GetKind(), name, parent) // namespace is enforced by NameKey

	err := client.ds.Get(ctx, key, entity.GetValue())
	if err != nil {
		return fmt.Errorf("failed to get named datastore entity: %w", err)
	}
	return nil
}

// DeleteNamed will delete a single named Entity using its name key and parent (if exists)
func (client client) DeleteNamed(ctx context.Context, namespace, kind, name string, parent *datastore.Key) error {
	key := client.NameKey(namespace, kind, name, parent) // namespace is enforced by NameKey

	err := client.ds.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to delete named datastore entity: %w", err)
	}
	return nil
}

// QueryAll gets all possible results from a query
func (client client) QueryAll(ctx context.Context, namespace string, query *datastore.Query, entitySlicePtr interface{}) (keys []*datastore.Key, err error) {
	q := query.Namespace(namespace)
	return client.Client().GetAll(ctx, q, entitySlicePtr)
}
