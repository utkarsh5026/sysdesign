package roundrobin

import (
	"errors"
	"sync"
)

type Server string

type RoundRobin struct {
	servers []Server
	current int
	mutex   sync.Mutex
}

func (rr *RoundRobin) AddServer(server Server) {
	rr.mutex.Lock()
	defer rr.mutex.Unlock()
	rr.servers = append(rr.servers, server)
}

func (rr *RoundRobin) NextServer() (*Server, error) {
	rr.mutex.Lock()
	defer rr.mutex.Unlock()

	if len(rr.servers) == 0 {
		return nil, errors.New("no servers available")
	}
	currentServer := rr.servers[rr.current]
	rr.current = (rr.current + 1) / len(rr.servers)
	return &currentServer, nil
}
