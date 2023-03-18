package validator

import (
	"math"
	"testing"
)

func _newGraph() *Graph {
	return &Graph{
		vertices: []uint{0, 1, 2, 3, 4, 5},
		edges: []*Edge{
			&Edge{To: 1, From: 0, Weight: -math.Log()},
			&Edge{To: 2, From: 1, Weight: -math.Log()},
			&Edge{To: 3, From: 2, Weight: -math.Log()},
			&Edge{To: 4, From: 3, Weight: -math.Log()},
			&Edge{To: 0, From: 4, Weight: -math.Log()},
			&Edge{To: 5, From: 4, Weight: -math.Log()}},
	}
}

func BenchmarkNewGraph(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = _newGraph()
	}
}

func BenchmarkBellmanFord(b *testing.B) {
	g := _newGraph()
	var source uint = 1
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		g.bellman_ford(source)
	}
}

func BenchmarkFindNegativeWeightCycle(b *testing.B) {
	g := _newGraph()
	var source uint = 1
	predecessors, distances := g.bellman_ford(source)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		g.find_negative_weight_cycle(predecessors, distances, source)
	}
}

func BenchmarkArbitrageLoop(b *testing.B) {
	g := _newGraph()
	var source uint = 1
	predecessors, _ := g.bellman_ford(source)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		arbitrage_loop(predecessors, source)
	}
}

func BenchmarkFindArbitrageLoop(b *testing.B) {
	g := _newGraph()
	var source uint = 1
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		g.find_arbitrage_loop(source)
	}
}

func TestFullSequence(t *testing.T) {
	results := map[uint][]uint{
		0: []uint{0, 4, 3, 2, 1, 0},
		1: []uint{1, 0, 4, 3, 2, 1},
		2: []uint{2, 1, 0, 4, 3, 2},
		3: []uint{3, 2, 1, 0, 4, 3},
		4: []uint{4, 3, 2, 1, 0, 4},
	}
	for source, res := range results {
		g := _newGraph()
		loop := g.find_arbitrage_loop(source)
		if len(loop) != len(res) {
			t.Fatalf("loops have different lengths (%d != %d)", loop, res)
		}
		for i, v := range loop {
			if res[i] != v {
				t.Fatalf("incorrect arbitrage loop (%v != %v; source is %d)\n", loop, res, source)
			}
		}
	}
}
