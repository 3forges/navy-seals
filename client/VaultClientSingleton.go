package client

import (
	"log"
	"sync"
)

var once sync.Once

var (
	instance *VaultClient
)

func GetVaultClient() (*VaultClient, error) {
	var err error
	once.Do(func() { // <-- atomic, does not allow repeating

		instance, err = NewVaultClient(0) // <-- thread safe
		if err != nil {
			log.Printf("GetVaultClient() SINGLETON - Failed to instanciate Vault Client Singleton %v", err)
		}

	})

	return instance, err
}
