package client

import (
	"fmt"
	"navy-seals/config"
	"time"

	vapi "github.com/openbao/openbao/api/v2"
)

const (
	// defaultCheckInterval  = 30 * time.Second
	defaultTimeout = 15 * time.Second
	// configRefreshInterval = 15 * time.Second
)

type VaultClient struct {
	ID     int
	client *vapi.Client
}

// NewService - our constructor function
func NewVaultClient(id int) (*VaultClient, error) {
	var err error

	vconfig := vapi.DefaultConfig()
	config.LoadConfig()
	fmt.Printf("Creating new vault client with VaultAddress [%v]", config.ApiConfig.VaultAddress)
	// "http://openbao.pesto.io:8200"
	vconfig.Address = config.ApiConfig.VaultAddress
	// vconfig.Address = "http://openbao.pesto.io:8200"
	vconfig.MaxRetries = 0
	vconfig.Timeout = defaultTimeout

	if err = vconfig.ConfigureTLS(&vapi.TLSConfig{Insecure: config.ApiConfig.TLSSkipVerify}); err != nil {
		// logger.WithError(err).Fatal("error initializing tls config")
		fmt.Printf("error initializing tls config %v", err)
	}

	var vault_client *vapi.Client
	if vault_client, err = vapi.NewClient(vconfig); err != nil {
		// logger.Fatalf("error creating vault client: %v", err)
		fmt.Printf("Error creating vault client: %v", err)
	}

	createdClient := &VaultClient{
		ID:     id,
		client: vault_client,
	}
	// handles other potentially more complex setup logic
	// for our component, there could be calls to downstream
	// dependencies to check connections etc that could return
	// errors
	return createdClient, err
}
func (c *VaultClient) GetClient() *vapi.Client {
	// Perform a health check about the client
	// For simulation, assume clients are always healthy
	return c.client
}

func (c *VaultClient) HealthCheck() bool {
	// Perform a health check about the client
	// For simulation, assume clients are always healthy
	return true
}
