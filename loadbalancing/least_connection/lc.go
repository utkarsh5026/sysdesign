package leastconnection

import (
	"container/heap"
	"errors"
)

type LeastConnection struct {
	Servers ServerQueue
}

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

func (lc *LeastConnection) GetNextServer() (*Server, error) {
	if lc.Servers.Len() == 0 {
		return nil, errors.New("no servers available")
	}

	maxConnServer := heap.Pop(&lc.Servers).(*Server)
	maxConnServer.AddConnection()
	heap.Push(&lc.Servers, maxConnServer)
	return maxConnServer, nil
}

func (lc *LeastConnection) ReleaseServer(server *Server) {
	server.FreeConnection()
	heap.Fix(&lc.Servers, server.index)
}
