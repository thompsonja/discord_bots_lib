package secrets

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/api/iterator"
)

func GetLatestSecretValue(gcpSecretName, projectID string) (string, error) {
	if gcpSecretName == "" && projectID == "" {
		return "", fmt.Errorf("gcpSecretName and projectID must be provided")
	}

	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("secretmanager.NewClient: %v", err)
	}
	defer client.Close()

	version, err := client.ListSecretVersions(ctx, &secretmanagerpb.ListSecretVersionsRequest{
		Parent: fmt.Sprintf("projects/%s/secrets/%s", projectID, gcpSecretName),
	}).Next()
	if err == iterator.Done {
		return "", fmt.Errorf("no secret versions for secret %s in project %s", gcpSecretName, projectID)
	}
	if err != nil {
		return "", fmt.Errorf("client.ListSecretVersions: %v", err)
	}

	result, err := client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: version.Name,
	})
	if err != nil {
		return "", fmt.Errorf("client.AccessSecretVersion: %v", err)
	}

	return string(result.Payload.Data), nil
}
