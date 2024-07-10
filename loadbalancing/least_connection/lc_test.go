package leastconnection

import (
	"testing"
)

// TestLeastConnectionLoadBalancer tests the creation of a new load balancer
func TestLeastConnectionLoadBalancer(t *testing.T) {
	servers := []Server{
		{ID: "server1", Connections: 0},
		{ID: "server2", Connections: 0},
		{ID: "server3", Connections: 0},
	}

	lb := LeastConnectionLoadBalancer(servers)

	if lb == nil {
		t.Fatal("Expected non-nil load balancer")
	}

	if len(lb.Servers) != 3 {
		t.Errorf("Expected 3 servers, got %d", len(lb.Servers))
	}

	// Check if the servers are in the correct order (they should be, as all have 0 connections)
	for i, server := range lb.Servers {
		if server.ID != servers[i].ID {
			t.Errorf("Expected server %s at position %d, got %s", servers[i].ID, i, server.ID)
		}
	}
}

// TestGetNextServer tests the GetNextServer method
func TestGetNextServer(t *testing.T) {
	servers := []Server{
		{ID: "server1", Connections: 1},
		{ID: "server2", Connections: 0},
		{ID: "server3", Connections: 2},
	}

	lb := LeastConnectionLoadBalancer(servers)

	// First call should return server2 (0 connections)
	server, err := lb.GetNextServer()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if server.ID != "server2" {
		t.Errorf("Expected server2, got %s", server.ID)
	}
	if server.Connections != 1 {
		t.Errorf("Expected 1 connection, got %d", server.Connections)
	}

	// Second call should return server1 (1 connection)
	server, err = lb.GetNextServer()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if server.ID != "server1" {
		t.Errorf("Expected server1, got %s", server.ID)
	}
	if server.Connections != 2 {
		t.Errorf("Expected 2 connections, got %d", server.Connections)
	}
}

// TestGetNextServerNoServers tests GetNextServer when no servers are available
func TestGetNextServerNoServers(t *testing.T) {
	lb := LeastConnectionLoadBalancer([]Server{})

	_, err := lb.GetNextServer()
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err.Error() != "no servers available" {
		t.Errorf("Expected 'no servers available' error, got '%v'", err)
	}
}

// TestReleaseServer tests the ReleaseServer method
func TestReleaseServer(t *testing.T) {
	servers := []Server{
		{ID: "server1", Connections: 2},
		{ID: "server2", Connections: 1},
		{ID: "server3", Connections: 3},
	}

	lb := LeastConnectionLoadBalancer(servers)

	// Release server3
	lb.ReleaseServer(lb.Servers[2])

	// Check if server3 now has 2 connections
	if lb.Servers[2].Connections != 2 {
		t.Errorf("Expected 2 connections for server3, got %d", lb.Servers[2].Connections)
	}

	// GetNextServer should now return server2 (still 1 connection)
	server, err := lb.GetNextServer()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if server.ID != "server2" {
		t.Errorf("Expected server2, got %s", server.ID)
	}
}

// TestLoadBalancerIntegration tests the overall behavior of the load balancer
func TestLoadBalancerIntegration(t *testing.T) {
	servers := []Server{
		{ID: "server1", Connections: 2},
		{ID: "server2", Connections: 0},
		{ID: "server3", Connections: 0},
	}

	lb := LeastConnectionLoadBalancer(servers)

	// Simulate some load balancing
	expectedOrder := []string{"server1", "server2", "server3", "server1", "server2"}
	for i, expectedID := range expectedOrder {
		server, err := lb.GetNextServer()
		if err != nil {
			t.Fatalf("Unexpected error on iteration %d: %v", i, err)
		}
		if server.ID != expectedID {
			t.Errorf("Iteration %d: Expected %s, got %s", i, expectedID, server.ID)
		}
	}

	// Release server1 twice
	lb.ReleaseServer(lb.Servers[0])
	lb.ReleaseServer(lb.Servers[0])

	// Now server1 should be the next server
	server, err := lb.GetNextServer()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if server.ID != "server1" {
		t.Errorf("Expected server1, got %s", server.ID)
	}
}
