package weightedleastconnection

import (
	"container/heap"
	"errors"
)

// WeightedLeastConnection represents a weighted least connection load balancer.
// It maintains a priority queue of servers, where the server with the lowest
// weighted connection ratio is always at the top.
type WeightedLeastConnection struct {
	servers ServerQueue
}

// WeightedLeastConnectionLoadBalancer creates and initializes a new WeightedLeastConnection load balancer.
// It takes a slice of servers and returns a pointer to the new WeightedLeastConnection instance.
//
// Parameters:
//   - servers: A slice of Server structs representing the available servers.
//
// Returns:
//   - *WeightedLeastConnection: A pointer to the initialized WeightedLeastConnection load balancer.
func WeightedLeastConnectionLoadBalancer(servers []Server) *WeightedLeastConnection {
	wlc := &WeightedLeastConnection{
		servers: make(ServerQueue, len(servers)),
	}

	for i := range servers {
		wlc.servers[i] = &servers[i]
		wlc.servers[i].index = i
	}
	heap.Init(&wlc.servers)
	return wlc
}

// NextServer selects the server with the lowest weighted connection ratio,
// increments its connection count, and returns it.
//
// Returns:
//   - *Server: A pointer to the selected server with the lowest weighted connection ratio.
//   - error: An error if no servers are available.
func (wlc *WeightedLeastConnection) NextServer() (*Server, error) {
	if wlc.servers.Len() == 0 {
		return nil, errors.New("no servers available")
	}

	server := heap.Pop(&wlc.servers).(*Server)
	server.AddConnection()
	heap.Push(&wlc.servers, server)

	return server, nil
}

// ReleaseServer decrements the connection count for the given server
// and adjusts its position in the priority queue.
//
// Parameters:
//   - server: A pointer to the Server being released.
func (wlc *WeightedLeastConnection) ReleaseServer(server *Server) {
	server.ReleaseConnection()
	heap.Fix(&wlc.servers, server.index)
}

// UpdateServerWeight updates the weight of a server and adjusts its position in the priority queue.
// If the new weight is negative, it is set to 0.
//
// Parameters:
//   - server: A pointer to the Server whose weight is being updated.
//   - newWeight: The new weight to be assigned to the server.
func (wlc *WeightedLeastConnection) UpdateServerWeight(server *Server, newWeight int) {
	if newWeight < 0 {
		newWeight = 0
	}
	server.Weight = newWeight
	heap.Fix(&wlc.servers, server.index)
}
