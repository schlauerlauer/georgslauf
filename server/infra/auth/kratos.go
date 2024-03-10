package auth

import (
	"georgslauf/infra/persistence"

	ory "github.com/ory/client-go"
)

type AuthService struct {
	Client *ory.APIClient
	Auth   AuthInterface
}

func NewKratosClient(publicUrl string, repository *persistence.Repository) *AuthService {
	config := ory.NewConfiguration()
	config.Servers = ory.ServerConfigurations{
		{
			URL: publicUrl,
		},
	}
	oryClient := ory.NewAPIClient(config)
	return &AuthService{
		Auth:   NewAuth(oryClient, repository),
		Client: oryClient,
	}
}
