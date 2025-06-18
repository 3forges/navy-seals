package client

import (
	"errors"
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
	log.Println("Initializing new vault client pool")
	pool.expandPool(maxVaultClients / 2) // Start with half of the max capacity
	return pool
}

func (p *VaultClientPool) BorrowVaultClient() (*VaultClient, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	timeoutChan := time.After(p.timeout)
	for {
		if len(p.clients) > 0 {
			conn := p.clients[0]
			p.clients = p.clients[1:]
			if conn.HealthCheck() {
				log.Printf("Borrowed vault client ID %d", conn.ID)
				return conn, nil
			}
		}

		if p.currentVaultClients < p.maxVaultClients {
			log.Println("Expanding pool due to high demand")
			p.expandPool(1)
		}

		select {
		case <-timeoutChan:
			log.Println("Failed to borrow vault client: timeout exceeded")
			return nil, errors.New("timeout exceeded, no healthy clients available")
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (p *VaultClientPool) ReleaseVaultClient(conn *VaultClient) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if conn.HealthCheck() {
		p.clients = append([]*VaultClient{conn}, p.clients...)
		log.Printf("Returned vault client ID %d to pool", conn.ID)
	} else {
		log.Printf("VaultClient ID %d failed health check, not returned to pool", conn.ID)
		if p.currentVaultClients < p.maxVaultClients {
			p.expandPool(1)
		}
	}
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
			log.Printf("[expandPool(num int)] - ERROR creating new vault client %v to pool", err)
			// os.Exit(7)
		} else {
			p.clients = append(p.clients, newclient)
		}

		p.currentVaultClients++
		log.Printf("Added new vault client ID %d to pool", p.currentVaultClients)
	}
}
