package client

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

type VaultClientPool struct {
	mutex               sync.Mutex
	clients             []*VaultClient
	maxVaultClients     int
	currentVaultClients int
	timeout             time.Duration
}

func NewVaultClientPool(maxVaultClients int, timeout time.Duration) *VaultClientPool {
	pool := &VaultClientPool{
		clients:             make([]*VaultClient, 0, maxVaultClients),
		maxVaultClients:     maxVaultClients,
		currentVaultClients: 0,
		timeout:             timeout,
	}
	log.Println("\nNAVY SEAL VAULT CLIENTS POOL - Initializing new vault client pool")
	pool.expandPool(maxVaultClients / 2) // Start with half of the max capacity
	fmt.Printf("\nNAVY SEAL VAULT CLIENTS POOL - [NewVaultClientPool()] - AFTER EXPANDING POOL TO POOL SIZE IS: %v", len(pool.clients))
	return pool
}

func (p *VaultClientPool) BorrowVaultClient() (*VaultClient, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	fmt.Printf("\nNAVY SEAL VAULT CLIENTS POOL - [BorrowVaultClient()] - POOL SIZE IS: %v", len(p.clients))

	timeoutChan := time.After(p.timeout)
	for {
		if len(p.clients) > 0 {
			vaultClient := p.clients[0]
			p.clients = p.clients[1:]
			if vaultClient.HealthCheck() {
				log.Printf("\nNAVY SEAL VAULT CLIENTS POOL - [BorrowVaultClient()] - POOL - Borrowed vault client ID %d", vaultClient.ID)
				return vaultClient, nil
			}
		}

		if p.currentVaultClients < p.maxVaultClients {
			log.Println("\nNAVY SEAL VAULT CLIENTS POOL - Expanding pool due to high demand")
			p.expandPool(1)
		}

		select {
		case <-timeoutChan:
			log.Println("\nNAVY SEAL VAULT CLIENTS POOL - Failed to borrow vault client: timeout exceeded")
			return nil, errors.New("\nNAVY SEAL VAULT CLIENTS POOL - timeout exceeded, no healthy clients available")
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (p *VaultClientPool) ReleaseVaultClient(vaultClient *VaultClient) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	fmt.Printf("\nNAVY SEAL VAULT CLIENTS POOL - [ReleaseVaultClient()] - BEFORE RELEASING CLIENT, POOL SIZE IS [%v] ", len(p.clients))
	if vaultClient.HealthCheck() {
		p.clients = append([]*VaultClient{vaultClient}, p.clients...)
		log.Printf("\nNAVY SEAL VAULT CLIENTS POOL - Returned vault client ID %d to pool", vaultClient.ID)
	} else {
		log.Printf("\nNAVY SEAL VAULT CLIENTS POOL - VaultClient ID %d failed health check, not returned to pool", vaultClient.ID)
		if p.currentVaultClients < p.maxVaultClients {
			p.expandPool(1)
		}
	}
	fmt.Printf("\nNAVY SEAL VAULT CLIENTS POOL - [ReleaseVaultClient()] - AFTER RELEASING CLIENT, POOL SIZE IS [%v] ", len(p.clients))
}

func (p *VaultClientPool) expandPool(num int) {
	for i := 0; i < num; i++ {
		if p.currentVaultClients >= p.maxVaultClients {
			return
		}
		// p.clients = append(p.clients, &VaultClient{ID: p.currentVaultClients})
		var newclient *VaultClient
		var err error
		newclient, err = NewVaultClient(p.currentVaultClients)
		if err != nil {
			log.Printf("\nNAVY SEAL VAULT CLIENTS POOL - [expandPool(num int)] - ERROR creating new vault client %v to pool", err)
			// os.Exit(7)
		} else {
			p.clients = append(p.clients, newclient)
		}

		p.currentVaultClients++
		log.Printf("Added new vault client ID %d to pool", p.currentVaultClients)
	}
}
