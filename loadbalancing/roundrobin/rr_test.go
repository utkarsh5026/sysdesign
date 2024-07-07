package roundrobin

import (
	"testing"
)

func TestAddServer(t *testing.T) {
	rr := &RoundRobin{}
	server := Server("server1")

	rr.AddServer(server)

	if len(rr.servers) != 1 {
		t.Errorf("expected 1 server, got %d", len(rr.servers))
	}

	if rr.servers[0] != server {
		t.Errorf("expected server %s, got %s", server, rr.servers[0])
	}
}

func TestNextServer(t *testing.T) {
	t.Run("no servers available", func(t *testing.T) {
		rr := &RoundRobin{}
		_, err := rr.NextServer()
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("single server", func(t *testing.T) {
		rr := &RoundRobin{}
		server := Server("server1")
		rr.AddServer(server)

		srv, err := rr.NextServer()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if *srv != server {
			t.Errorf("expected server %s, got %s", server, *srv)
		}
	})

	t.Run("multiple servers", func(t *testing.T) {
		rr := &RoundRobin{}
		server1 := Server("server1")
		server2 := Server("server2")
		rr.AddServer(server1)
		rr.AddServer(server2)

		srv, err := rr.NextServer()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if *srv != server1 {
			t.Errorf("expected server %s, got %s", server1, *srv)
		}

		srv, err = rr.NextServer()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if *srv != server2 {
			t.Errorf("expected server %s, got %s", server2, *srv)
		}
	})
}
