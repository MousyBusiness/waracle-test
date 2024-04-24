package secrets

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/api/cloudresourcemanager/v1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"os"
)

type Store interface {
	GetSecret(ctx context.Context, secretID string) (string, error)
	CreateSecret(ctx context.Context, secretID string, secret []byte) error
	DeleteSecret(ctx context.Context, secretID string) error
}

type store struct{}

var Vault Store

func NewSecretManagerStore() *store {
	return &store{}
}

// GetSecret retrieves a secret from GCP Secret Manager
func (s *store) GetSecret(ctx context.Context, secretID string) (string, error) {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed to setup client")
	}
	projectNumber, err := getProjectNumber() // TODO store project number in environment variable
	if err != nil {
		return "", err
	}

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%d/secrets/%s/versions/%d", projectNumber, secretID, 1),
	}
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", errors.Wrap(err, "failed to access secret version")
	}

	if string(result.Payload.Data) == "" {
		return "", errors.New("secret empty")
	}

	return string(result.Payload.Data), nil
}

var projectNumberCache int64

func getProjectNumber() (int64, error) {
	if projectNumberCache != 0 {
		return projectNumberCache, nil
	}
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	cloudresourcemanagerService, err := cloudresourcemanager.NewService(context.Background())
	if err != nil {
		return 0, err
	}

	project, err := cloudresourcemanagerService.Projects.Get(projectID).Do()
	if err != nil {
		return 0, err
	}

	projectNumberCache = project.ProjectNumber
	return project.ProjectNumber, nil
}
