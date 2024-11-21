package vault

import (
	"context"
	"fmt"
	"os"

	"github.com/abhissng/internal/utils"
	infisical "github.com/infisical/go-sdk"
	"github.com/infisical/go-sdk/packages/models"
)

type Vault struct {
	client    infisical.InfisicalClientInterface
	env       string
	projectId string
	path      string
}

func NewVault(env, projectId, path string) *Vault {
	if utils.IsEmpty(path) {
		path = "/"
	}
	return &Vault{
		env:       env,
		projectId: projectId,
		path:      path,
	}
}

// Vault client initialization
func (v *Vault) InitVaultClient() {

	v.client = infisical.NewInfisicalClient(context.Background(), infisical.Config{
		SiteUrl:          "https://app.infisical.com", // Optional, default is https://app.infisical.com
		AutoTokenRefresh: true,                        // Wether or not to let the SDK handle the access token lifecycle. Defaults to true if not specified.
	})

	_, err := v.client.Auth().UniversalAuthLogin("", "")
	if err != nil {
		fmt.Printf("Authentication failed with the vault: %v", err)
		os.Exit(1)
	}

}

func (v *Vault) retreiveSecret(key string) (models.Secret, error) {
	apiKeySecret, err := v.client.Secrets().Retrieve(infisical.RetrieveSecretOptions{
		SecretKey:   key,
		Environment: v.env,
		ProjectID:   v.projectId,
		SecretPath:  v.path,
	})
	if err != nil {
		fmt.Printf("Error retreiving secret %s from vault: %v", key, err)
		return models.Secret{}, err
	}

	return apiKeySecret, nil
}

func (v *Vault) FetchVaultValue(key string) (string, error) {

	secret, err := v.retreiveSecret(key)
	if err != nil {
		fmt.Printf("Error fetching %s values from vault: %v", key, err)
		return "", err
	}

	return secret.SecretValue, nil
}
