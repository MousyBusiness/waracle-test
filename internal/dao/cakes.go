package dao

import (
	"cloud.google.com/go/datastore"
	"context"
	"errors"
	"github.com/mousybusiness/waracle-test/internal/db"
	"strconv"
)

const namespace = "Bakery"
const kind = "Cake"

// Cake represents a cake with specific properties.
type Cake struct {
	ID        int    `json:"id"`         // Unique identifier for the cake.
	Name      string `json:"name"`       // Name of the cake, max 30 characters.
	Comment   string `json:"comment"`    // A comment about the cake, max 200 characters.
	ImageURL  string `json:"image_url"`  // URL to an image of the cake.
	YumFactor int    `json:"yum_factor"` // Rating from 1 to 5 inclusive.
}

func (c *Cake) GetKind() string {
	return kind
}

func (c *Cake) GetValue() any {
	return c
}

// Validate ensures cakes adhere to input constraints
func (c *Cake) Validate() error {

	if len(c.Name) > 30 {
		return errors.New("name must be at most 30 characters long")
	}
	if len(c.Comment) > 200 {
		return errors.New("comment must be at most 200 characters long")
	}
	if c.YumFactor < 1 || c.YumFactor > 5 {
		return errors.New("yumFactor must be between 1 and 5")
	}

	return nil
}

// CreateCake creates a cake entity in Datastore
func CreateCake(ctx context.Context, db db.Database, c *Cake) error {
	if _, err := db.CreateNamed(ctx, namespace, strconv.Itoa(c.ID), nil, c); err != nil {
		return err
	}
	return nil
}

// GetCake gets a cake entity from Datastore using the cake ID
func GetCake(ctx context.Context, db db.Database, id int) (*Cake, error) {
	var cake Cake
	if err := db.GetNamed(ctx, namespace, strconv.Itoa(id), nil, &cake); err != nil {
		return nil, err
	}
	return &cake, nil
}

// ListCakes lists all cake entities in Datastore
func ListCakes(ctx context.Context, db db.Database) ([]*Cake, error) {
	query := datastore.NewQuery(kind)
	var cakes []*Cake
	_, err := db.QueryAll(ctx, namespace, query, &cakes)
	if err != nil {
		return nil, err
	}
	return cakes, nil
}

type SearchRequest struct {
	Name      *string `json:"name"`
	YumFactor *int    `json:"yum_factor"`
}

// SearchCakes searches Datastore for matching names and yum-factors
func SearchCakes(ctx context.Context, db db.Database, request SearchRequest) ([]*Cake, error) {
	query := datastore.NewQuery(kind)
	if request.Name != nil {
		query = query.FilterField("Name", "=", *request.Name)
	}
	if request.YumFactor != nil {
		query = query.FilterField("YumFactor", "=", *request.YumFactor)
	}
	var cakes []*Cake
	_, err := db.QueryAll(ctx, namespace, query, &cakes)
	if err != nil {
		return nil, err
	}
	return cakes, nil
}

// DeleteCake deletes an entity from datastore using the cake ID
func DeleteCake(ctx context.Context, db db.Database, id int) error {
	if err := db.DeleteNamed(ctx, namespace, "Cake", strconv.Itoa(id), nil); err != nil {
		return err
	}
	return nil
}
