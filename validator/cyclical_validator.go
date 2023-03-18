package validator

import (
	"math"
)

type Graph struct {
	edges    []*Edge
	vertices []uint
}

type Edge struct {
	From, To uint
	Weight   float64
}

func new_edge(from, to uint, weight float64) *Edge {
	return &Edge{From: from, To: to, Weight: weight}
}

func new_graph(edges []*Edge, vertices []uint) *Graph {
	return &Graph{edges: edges, vertices: vertices}
}

func (g *Graph) find_arbitrage_loop(source uint) []uint {
	predecessors, distances := g.bellman_ford(source)
	return g.find_negative_weight_cycle(predecessors, distances, source)
}

func (g *Graph) bellman_ford(source uint) ([]uint, []float64) {
	size := len(g.vertices)
	distances := make([]float64, size)
	predecessors := make([]uint, size)
	for _, v := range g.vertices {
		distances[v] = math.MaxFloat64
	}
	distances[source] = 0

	for i, changes := 0, 0; i < size-1; i, changes = i+1, 0 {
		for _, edge := range g.edges {
			if newDist := distances[edge.From] + edge.Weight; newDist < distances[edge.To] {
				distances[edge.To] = newDist
				predecessors[edge.To] = edge.From
				changes++
			}
		}
		if changes == 0 {
			break
		}
	}
	return predecessors, distances
}

func (g *Graph) find_negative_weight_cycle(predecessors []uint, distances []float64, source uint) []uint {
	for _, edge := range g.edges {
		if distances[edge.From]+edge.Weight < distances[edge.To] {
			return arbitrage_loop(predecessors, source)
		}
	}
	return nil
}

func arbitrage_loop(predecessors []uint, source uint) []uint {
	size := len(predecessors)
	loop := make([]uint, size)
	loop[0] = source

	exists := make([]bool, size)
	exists[source] = true

	indices := make([]uint, size)

	var index, next uint
	for index, next = 1, source; ; index++ {
		next = predecessors[next]
		loop[index] = next
		if exists[next] {
			return loop[indices[next] : index+1]
		}
		indices[next] = index
		exists[next] = true
	}
}
