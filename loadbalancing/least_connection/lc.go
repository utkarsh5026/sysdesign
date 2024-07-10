package leastconnection

import (
	"container/heap"
	"errors"
)

// LeastConnection represents a least connection load balancer.
// It maintains a priority queue of servers, where the server with
// the least number of active connections is always at the top.
type LeastConnection struct {
	Servers ServerQueue
}

// LeastConnectionLoadBalancer creates and initializes a new LeastConnection load balancer.
// It takes a slice of servers and returns a pointer to the new LeastConnection instance.
//
// Parameters:
//   - servers: A slice of Server structs representing the available servers.
//
// Returns:
//   - *LeastConnection: A pointer to the initialized LeastConnection load balancer.
func LeastConnectionLoadBalancer(servers []Server) *LeastConnection {
	lc := &LeastConnection{
		Servers: make(ServerQueue, len(servers)),
	}

	for i := range servers {
		lc.Servers[i] = &servers[i]
		lc.Servers[i].index = i
	}

	heap.Init(&lc.Servers)
	return lc
}

// GetNextServer selects the server with the least number of active connections,
// increments its connection count, and returns it.
//
// Returns:
//   - *Server: A pointer to the selected server with the least connections.
//   - error: An error if no servers are available.
func (lc *LeastConnection) GetNextServer() (*Server, error) {
	if lc.Servers.Len() == 0 {
		return nil, errors.New("no servers available")
	}

	maxConnServer := heap.Pop(&lc.Servers).(*Server)
	maxConnServer.AddConnection()
	heap.Push(&lc.Servers, maxConnServer)
	return maxConnServer, nil
}

// ReleaseServer decrements the connection count for the given server
// and adjusts its position in the priority queue.
//
// Parameters:
//   - server: A pointer to the Server being released.
func (lc *LeastConnection) ReleaseServer(server *Server) {
	server.FreeConnection()
	heap.Fix(&lc.Servers, server.index)
}
