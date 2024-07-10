package leastconnection

type Server struct {
	ID          string
	Connections int
	index       int
}

func (s *Server) FreeConnection() {
	s.Connections--
}

func (s *Server) AddConnection() {
	s.Connections++
}

type ServerQueue []*Server

func (sq ServerQueue) Len() int {
	return len(sq)
}

func (sq ServerQueue) Less(i int, j int) bool {
	return sq[i].Connections < sq[j].Connections
}

func (sq ServerQueue) Swap(i int, j int) {
	sq[i], sq[j] = sq[j], sq[i]
	sq[i].index = j
	sq[j].index = i
}

func (sq *ServerQueue) Push(x any) {
	n := len(*sq)
	server := x.(*Server)
	server.index = n
	*sq = append(*sq, x.(*Server))
}

func (sq *ServerQueue) Pop() any {
	old := *sq
	n := len(old)
	x := old[n-1]
	*sq = old[0 : n-1]
	old[n-1] = nil // avoid memory leak
	x.index = -1   // for safety
	return x
}
