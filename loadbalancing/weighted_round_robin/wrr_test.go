package weightedroundrobin

import (
	"testing"
)

func TestNoServers(t *testing.T) {
	wrr := WeightedRoundRobinBalancer(0)
	_, err := wrr.NextServer()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "no server found" {
		t.Fatalf("expected 'no server found' error, got %v", err)
	}
}

func TestBadWeight(t *testing.T) {
	wrr := WeightedRoundRobinBalancer(3)

	wrr.AddServer(Server{Weight: 0, Id: "Server 1"})
	wrr.AddServer(Server{Weight: 0, Id: "Server 2"})
	wrr.AddServer(Server{Weight: 0, Id: "Server 3"})

	_, err := wrr.NextServer()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "bad weight" {
		t.Fatalf("expected 'bad weight' error, got %v", err)
	}
}

func TestAddServers(t *testing.T) {
	wrr := WeightedRoundRobinBalancer(2)
	wrr.AddServer(Server{Weight: 2, Id: "Server 1"})
	wrr.AddServer(Server{Weight: 3, Id: "Server 2"})

	expectedMaxWeight := 3
	if wrr.maxWeight != expectedMaxWeight {
		t.Errorf("expected maxWeight %d, got %d", expectedMaxWeight, wrr.maxWeight)
	}

	expectedGCDWeight := 1
	if wrr.gcdWeight != expectedGCDWeight {
		t.Errorf("expected gcdWeight %d, got %d", expectedGCDWeight, wrr.gcdWeight)
	}
}

func TestAddServerIncreasesMaxWeight(t *testing.T) {
	wrr := WeightedRoundRobinBalancer(2)
	wrr.AddServer(Server{Weight: 2, Id: "Server 1"})
	wrr.AddServer(Server{Weight: 5, Id: "Server 2"})

	expectedMaxWeight := 5
	if wrr.maxWeight != expectedMaxWeight {
		t.Errorf("expected maxWeight %d, got %d", expectedMaxWeight, wrr.maxWeight)
	}
}
