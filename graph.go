package main

import (
	"fmt"
	"sort"
)

type Edge struct {
	start   int
	end     int
	cost    int
	benefit int
	ocurr   int
}

type Graph struct {
	edges         map[int]map[int]*Edge
	positiveEdges []*Edge
}

func (e *Edge) String() string {
	return fmt.Sprintf("(%d,%d)", e.start, e.end)
}

func (g *Graph) String() string {
	s := ""
	keys := make([]int, len(g.edges))
	i := 0
	for k, _ := range g.edges {
		keys[i] = k
		i++
	}
	sort.Ints(keys)
	for _, elem := range keys {
		s = s + fmt.Sprintf("%d -> ", elem)
		for _, v2 := range g.edges[elem] {
			s = s + fmt.Sprintf("%s", v2)
		}
		s = s + "\n"
	}
	return s
}

func (g *Graph) AddEdge(start, end, cost, benefit int) {

	if _, ok := g.edges[start][end]; ok {
		g.edges[start][end].ocurr = g.edges[start][end].ocurr + 1
	} else {
		g.edges[start][end] = &Edge{start, end, cost, benefit, 1}
	}

	if _, ok := g.edges[end][start]; ok {
		g.edges[end][start].ocurr = g.edges[end][start].ocurr + 1
	} else {
		g.edges[end][start] = &Edge{end, start, cost, benefit, 1}
	}

	if benefit-cost >= 0 {
		g.positiveEdges = append(g.positiveEdges, g.edges[start][end])
	}
}

func (g *Graph) RemoveEdge(start, end int) {
	if _, ok := g.edges[start][end]; ok {
		if g.edges[start][end].ocurr <= 1 {
			delete(g.edges[start], end)
		} else {
			g.edges[start][end].ocurr = g.edges[start][end].ocurr - 1
		}
	}
	if _, ok := g.edges[end][start]; ok {
		if g.edges[end][start].ocurr <= 1 {
			delete(g.edges[end], start)
		} else {
			g.edges[end][start].ocurr = g.edges[end][start].ocurr - 1
		}
	}
}

func (g *Graph) Neighbors(node int) []int {
	neighbors := make([]int, 0, len(g.edges[node]))
	for k, _ := range g.edges[node] {
		neighbors = append(neighbors, k)
	}
	return neighbors
}

func (g *Graph) Degree(node int) int {
	return len(g.edges[node])
}

func (g *Graph) Cost(start, end int) int {
	return g.edges[start][end].cost
}

func (g *Graph) Benefit(start, end int) int {
	return g.edges[start][end].benefit
}

func NewGraph(nodes int) *Graph {
	g := &Graph{make(map[int]map[int]*Edge), make([]*Edge, 0)}
	for i := 1; i <= nodes; i++ {
		tmap := make(map[int]*Edge)
		g.edges[i] = tmap
	}
	return g
}
