package ds

func NewGraph() *Graph {
	return &Graph{
		adj: make(map[int][]int),
	}
}

func (g *Graph) AddEdge(from, to int) {
	g.adj[from] = append(g.adj[from], to)
}

func (g *Graph) GetEdges(node int) []int {
	return g.adj[node]
}

// Remove directed edge from -> to
func (g *Graph) RemoveEdge(from, to int) {
	neighbors := g.adj[from]
	for i, v := range neighbors {
		if v == to {
			g.adj[from] = append(neighbors[:i], neighbors[i+1:]...)
			break
		}
	}
}
