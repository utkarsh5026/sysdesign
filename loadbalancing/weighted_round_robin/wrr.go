package weightedroundrobin

import (
	"errors"
	"sync"
)

// Server represents a backend server in the load balancer.
type Server struct {
	Weight int    // The weight of the server, determining its selection frequency
	Id     string // A unique identifier for the server
}

// WeightedRoundRobin is the main struct implementing the Weighted Round Robin algorithm.
type WeightedRoundRobin struct {
	servers       []Server   // Slice of all registered servers
	current       int        // Index of the current server
	mutex         sync.Mutex // Mutex for ensuring thread-safety
	maxWeight     int        // Maximum weight among all servers
	currentWeight int        // Current weight in the selection algorithm
	gcdWeight     int        // Greatest common divisor of all server weights
}

// gcd computes the greatest common divisor of two numbers using the Euclidean algorithm.
func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// WeightedRoundRobinBalancer initializes and returns a new WeightedRoundRobin instance.
func WeightedRoundRobinBalancer(connections int) *WeightedRoundRobin {
	return &WeightedRoundRobin{
		servers: make([]Server, 0, connections), // Pre-allocate slice with capacity
		current: -1,                             // Start at -1 so first increment goes to 0
	}
}

// AddServer adds a new server to the load balancer.
// It updates MaxWeight and CurrentWeight if the new server has a higher weight.
func (wrr *WeightedRoundRobin) AddServer(server Server) {
	wrr.mutex.Lock()
	defer wrr.mutex.Unlock()

	wrr.servers = append(wrr.servers, server)

	if len(wrr.servers) == 1 {
		wrr.maxWeight = server.Weight
		wrr.currentWeight = wrr.maxWeight
		wrr.gcdWeight = server.Weight
	} else {
		if server.Weight > wrr.maxWeight {
			wrr.maxWeight = server.Weight
		}
		wrr.gcdWeight = gcd(wrr.gcdWeight, server.Weight)
	}
}

// NextServer selects the next server based on the Weighted Round Robin algorithm.
// It returns a pointer to the selected server and an error if no server is available.
func (wrr *WeightedRoundRobin) NextServer() (*Server, error) {
	wrr.mutex.Lock()
	defer wrr.mutex.Unlock()

	if len(wrr.servers) == 0 {
		return nil, errors.New("no server found")
	}

	for {
		// Move to the next server, wrapping around if necessary
		wrr.current = (wrr.current + 1) % len(wrr.servers)

		// If we've wrapped around to the start
		if wrr.current == 0 {
			// Decrease the current weight
			wrr.currentWeight -= wrr.gcdWeight
			if wrr.currentWeight <= 0 {
				// Reset the current weight to the maximum weight
				wrr.currentWeight = wrr.maxWeight
				if wrr.currentWeight == 0 {
					return nil, errors.New("bad weight")
				}
			}
		}

		// If the current server's weight is sufficient, select it
		if wrr.servers[wrr.current].Weight >= wrr.currentWeight {
			return &wrr.servers[wrr.current], nil
		}
	}
}
