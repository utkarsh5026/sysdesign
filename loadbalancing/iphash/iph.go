package iphash

import (
	"errors"
	"hash/fnv"
	"net"
	"sync"
)

// Server represents a server with a unique ID.
type Server struct {
	ID string
}

// IPHash is a load balancer that distributes requests based on the client's IP address.
type IPHash struct {
	servers []Server
	mutex   sync.RWMutex
}

// IpHashLoadBalancer initializes a new IPHash load balancer with the given servers.
func IpHashLoadBalancer(servers []Server) *IPHash {
	return &IPHash{servers: servers}
}

// AddServer adds a new server to the load balancer.
func (ip *IPHash) AddServer(server *Server) {
	ip.mutex.Lock()
	defer ip.mutex.Unlock()
	ip.servers = append(ip.servers, *server)
}

// RemoveServer removes a server from the load balancer by its ID.
// Returns an error if the server is not found.
func (ip *IPHash) RemoveServer(Id string) error {
	ip.mutex.Lock()
	defer ip.mutex.Unlock()

	for i, server := range ip.servers {
		if server.ID == Id {
			ip.servers = append(ip.servers[:i], ip.servers[i+1:]...)
			return nil
		}
	}
	return errors.New("server not found")
}

// GetServer returns the server assigned to the given client IP address.
// Returns an error if no servers are available.
func (ip *IPHash) GetServer(clientIP string) (Server, error) {
	ip.mutex.Lock()
	defer ip.mutex.Unlock()

	var server Server
	if len(ip.servers) == 0 {
		return server, errors.New("no server exists")
	}

	hash := hashIp(clientIP)
	index := hash % uint32(len(ip.servers))
	return ip.servers[index], nil
}

// hashIp generates a hash value for the given IP address.
func hashIp(ip string) uint32 {
	hash := fnv.New32a()
	hash.Write(net.ParseIP(ip).To4())
	return hash.Sum32()
}
