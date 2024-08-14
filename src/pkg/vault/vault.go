package vault

import (
	"athena/src/config"
	"fmt"
	vault "github.com/hashicorp/vault/api"
	"log"
)

var (
	configs  *config.Config
	instance *Client
)

// Client represents a client for interacting with Vault.
type Client struct {
	client *vault.Client
}

// Init creates a new Vault client.
func Init() error {
	return GetInstance().Connect()
}

// Connect Connects to vault.
func (v *Client) Connect() (err error) {
	configs = config.GetInstance()

	// Initiate new vault client.
	v.client, err = vault.NewClient(&vault.Config{
		Address: fmt.Sprintf("%s:%s", configs.Get("VAULT_HOST"), configs.Get("VAULT_PORT")),
	})

	if err != nil {
		log.Fatalln(err)
	}

	// Set Token.
	v.client.SetToken(configs.Get("VAULT_TOKEN"))

	// Check Health.
	health, err := v.client.Sys().Health()

	if err != nil {
		log.Fatalln(err)
	}

	if !health.Initialized {
		log.Fatalf("Vault Service:\n Initialization Status: Not initialized \n Please initialize vault before proceeding.")
	}

	if health.Sealed {
		log.Fatalf("Vault Service:\n Sealation Status: sealed \n Please unseal vault before proceeding.")
	}

	return nil
}

// GetVault returns the singleton instance of Vault.
func (v *Client) GetVault() *vault.Client {
	return v.client
}

// GetInstance returns the singleton instance of Vault.
func GetInstance() *Client {
	if instance == nil {
		instance = &Client{}
	}
	return instance
}
