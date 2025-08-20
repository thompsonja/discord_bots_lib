package secrets

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/cenkalti/backoff/v4"
	"google.golang.org/api/iterator"
)

func GetLatestSecretValue(ctx context.Context, gcpSecretName, projectID string) (string, error) {
	if gcpSecretName == "" && projectID == "" {
		return "", fmt.Errorf("gcpSecretName and projectID must be provided")
	}

	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("secretmanager.NewClient: %v", err)
	}
	defer client.Close()

	var version *secretmanagerpb.SecretVersion
	listVersionsOp := func() error {
		it := client.ListSecretVersions(ctx, &secretmanagerpb.ListSecretVersionsRequest{
			Parent: fmt.Sprintf("projects/%s/secrets/%s", projectID, gcpSecretName),
		})

		v, err := it.Next()
		if err == iterator.Done {
			return backoff.Permanent(fmt.Errorf("no secret versions for secret %s", gcpSecretName))
		}
		if err != nil {
			return fmt.Errorf("ListSecretVersions failed: %w", err)
		}

		version = v
		return nil
	}

	// Retry with exponential backoff
	err = backoff.Retry(listVersionsOp, backoff.WithContext(backoff.NewExponentialBackOff(), ctx))
	if err != nil {
		return "", fmt.Errorf("backoff.Retry for ListSecretVersions: %w", err)
	}

	var result *secretmanagerpb.AccessSecretVersionResponse
	accessSecretOp := func() error {
		res, err := client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
			Name: version.Name,
		})
		if err != nil {
			return fmt.Errorf("client.AccessSecretVersion: %v", err)
		}

		result = res
		return nil
	}

	// Retry with exponential backoff
	err = backoff.Retry(accessSecretOp, backoff.WithContext(backoff.NewExponentialBackOff(), ctx))
	if err != nil {
		return "", fmt.Errorf("backoff.Retry for AccessSecretVersion: %w", err)
	}

	return string(result.Payload.Data), nil
}
