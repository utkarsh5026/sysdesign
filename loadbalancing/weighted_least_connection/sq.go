// Package weightedleastconnection implements a weighted least connection load balancing algorithm.
// It provides data structures and methods to manage servers and distribute
// connections based on the number of active connections each server has,
// weighted by the server's capacity.
package weightedleastconnection

// Server represents a server in the load balancing pool.
type Server struct {
	ID          string // Unique identifier for the server
	Connections int    // Number of active connections to this server
	Weight      int    // Weight representing the server's capacity
	index       int    // Index of the server in the heap (used internally)
}

// AddConnection increments the number of active connections for the server.
// This method should be called when a new connection is established with the server.
func (s *Server) AddConnection() {
	s.Connections++
}

// ReleaseConnection decrements the number of active connections for the server.
// This method should be called when a connection to the server is closed.
func (s *Server) ReleaseConnection() {
	s.Connections--
}

// WeightedConnection calculates and returns the weighted connection ratio for the server.
// This ratio is used to determine the server's load relative to its capacity.
// If the weight is 0, it returns 0 to avoid division by zero.
func (s *Server) WeightedConnection() float64 {
	if s.Weight == 0 {
		return 0
	}
	return float64(s.Connections) / float64(s.Weight)
}

// ServerQueue is a priority queue of servers, implemented as a min-heap.
// The server with the lowest weighted connection ratio is always at the top of the heap.
type ServerQueue []*Server

// Len returns the number of servers in the queue.
// This method is part of the sort.Interface implementation.
func (sq ServerQueue) Len() int {
	return len(sq)
}

// Less compares two servers based on their weighted connection ratios.
// This method is part of the sort.Interface implementation.
func (sq ServerQueue) Less(i, j int) bool {
	return sq[i].WeightedConnection() < sq[j].WeightedConnection()
}

// Swap exchanges the positions of two servers in the queue.
// This method is part of the sort.Interface implementation.
func (sq ServerQueue) Swap(i int, j int) {
	sq[i], sq[j] = sq[j], sq[i]
	sq[i].index = i
	sq[j].index = j
}

// Push adds a new server to the queue.
// This method is part of the heap.Interface implementation.
func (sq *ServerQueue) Push(x any) {
	n := len(*sq)
	server := x.(*Server)
	server.index = n
	*sq = append(*sq, server)
}

// Pop removes and returns the server with the lowest weighted connection ratio from the queue.
// This method is part of the heap.Interface implementation.
func (sq *ServerQueue) Pop() any {
	old := *sq
	n := len(old)
	x := old[n-1]
	*sq = old[0 : n-1]
	old[n-1] = nil // avoid memory leak
	x.index = -1   // for safety
	return x
}
