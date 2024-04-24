package dao

import (
	"cloud.google.com/go/datastore"
	"context"
	"errors"
	"github.com/mousybusiness/waracle-test/internal/db"
	"math"
	"math/rand/v2"
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

// NewCake creates a new Cake instance with the given parameters, validating the constraints.
func NewCake(name, comment, imageURL string, yumFactor int) (*Cake, error) {
	id := rand.N[int](math.MaxInt32)

	if len(name) > 30 {
		return nil, errors.New("name must be at most 30 characters long")
	}
	if len(comment) > 200 {
		return nil, errors.New("comment must be at most 200 characters long")
	}
	if yumFactor < 1 || yumFactor > 5 {
		return nil, errors.New("yumFactor must be between 1 and 5")
	}

	return &Cake{
		ID:        id,
		Name:      name,
		Comment:   comment,
		ImageURL:  imageURL,
		YumFactor: yumFactor,
	}, nil
}

func CreateCake(ctx context.Context, db db.Database, c *Cake) error {
	if _, err := db.CreateNamed(ctx, namespace, strconv.Itoa(c.ID), nil, c); err != nil {
		return err
	}
	return nil
}

func GetCake(ctx context.Context, db db.Database, id int) (*Cake, error) {
	var cake Cake
	if err := db.GetNamed(ctx, namespace, strconv.Itoa(id), nil, &cake); err != nil {
		return nil, err
	}
	return &cake, nil
}

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

func DeleteCake(ctx context.Context, db db.Database, id int) error {
	if err := db.DeleteNamed(ctx, namespace, "Cake", strconv.Itoa(id), nil); err != nil {
		return err
	}
	return nil
}
