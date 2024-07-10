// Package leastconnection implements a least connection load balancing algorithm.
// It provides data structures and methods to manage servers and distribute
// connections based on the number of active connections each server has.
package leastconnection

// Server represents a server in the load balancing pool.
type Server struct {
	ID          string // Unique identifier for the server
	Connections int    // Number of active connections to this server
	index       int    // Index of the server in the heap (used internally)
}

// FreeConnection decrements the number of active connections for the server.
// This method should be called when a connection to the server is closed.
func (s *Server) FreeConnection() {
	s.Connections--
}

// AddConnection increments the number of active connections for the server.
// This method should be called when a new connection is established with the server.
func (s *Server) AddConnection() {
	s.Connections++
}

// ServerQueue is a priority queue of servers, implemented as a min-heap.
// The server with the least number of connections is always at the top of the heap.
type ServerQueue []*Server

// Len returns the number of servers in the queue.
// This method is part of the sort.Interface implementation.
func (sq ServerQueue) Len() int {
	return len(sq)
}

// Less compares two servers based on their number of connections.
// This method is part of the sort.Interface implementation.
func (sq ServerQueue) Less(i int, j int) bool {
	return sq[i].Connections < sq[j].Connections
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

// Pop removes and returns the server with the least connections from the queue.
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
