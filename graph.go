package main

import (
	"fmt"
	"math"
	"sort"
)

// Edge Type
type Edge struct {
	start   int
	end     int
	cost    int
	benefit int
	ocurr   int
}

// Graph Type
type Graph struct {
	edges         map[int]map[int]*Edge
	positiveEdges []*Edge
}

var mejorSol []*Edge
var beneficioDisponible int
var solParcial []*Edge

var totalBenefit int

func (e *Edge) String() string {
	return fmt.Sprintf("(%d,%d)", e.start, e.end)
}

func (g *Graph) String() string {
	s := ""
	keys := make([]int, len(g.edges))
	i := 0
	for k := range g.edges {
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

// AddEdge function to ad edges to a graph
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

// RemoveEdge function to remove edge from a graph
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

// Neighbors function to return Neigbors of node
func (g *Graph) Neighbors(node int) []int {
	neighbors := make([]int, 0, len(g.edges[node]))
	for k := range g.edges[node] {
		neighbors = append(neighbors, k)
	}
	return neighbors
}

// AddOcurr function to add ocurrence to an edge
func (g *Graph) AddOcurr(start, end int) {
	g.edges[start][end].ocurr = g.edges[start][end].ocurr + 1
	g.edges[end][start].ocurr = g.edges[end][start].ocurr + 1
}

// Dijkstra algorithm
func (g *Graph) Dijkstra(source int, to int, path []*Edge) []*Edge {
	// Create map to track distances from source vertex
	var u int
	dist := make(map[int]int)
	seen := make(map[int]bool)
	prev := make([]int, len(g.edges)+1)

	// Distance from source to source is zero
	dist[source] = 0

	// Initalize all distances to maximum
	for index := range g.edges {
		if index != source {
			dist[index] = math.MaxInt32
		}
	}

	// Iterate over all vertices
	for len(seen) < len(g.edges) {
		// Find vertex with minimum distance
		min := math.MaxInt32
		for index := range g.edges {
			if _, ok := seen[index]; dist[index] < min && !ok {
				min = dist[index]
				u = index
			}
		}
		seen[u] = true

		// Calculate minimum edge distance
		for _, edge := range g.edges[u] {
			if dist[edge.end] > dist[u]+costMinimumPath(edge, path) {
				dist[edge.end] = dist[u] + costMinimumPath(edge, path)
				prev[edge.end] = u
			}
			// if dist[edge.end] > dist[u]+edge.cost {
			// 	dist[edge.end] = dist[u] + edge.cost
			// 	prev[edge.end] = u
			// }
		}
	}
	return g.reconstructPath(source, to, prev)
}

func (g *Graph) reconstructPath(from int, to int, prev []int) []*Edge {
	var path []*Edge
	next := to
	for next != from {
		path = append(path, g.edges[next][prev[next]])
		next = prev[next]
	}
	return path
}

func costMinimumPath(edge *Edge, path []*Edge) int {
	for _, elem := range path {
		if edge == elem {
			return edge.cost
		}
	}
	if edge.benefit-edge.cost < 0 {
		return (edge.benefit - edge.cost) * (-1)
	}
	return 0
}

// Degree function to return the incidence degree of a node
func (g *Graph) Degree(node int) int {
	return len(g.edges[node])
}

// Cost function to return the cost of an edge
func (g *Graph) Cost(start, end int) int {
	return g.edges[start][end].cost
}

// Benefit function to return the benefit of an edge
func (g *Graph) Benefit(start, end int) int {
	return g.edges[start][end].benefit
}

func (g *Graph) estaEnSolucionParcial(e *Edge) bool {
	if !(g.edges[e.start][e.end] == nil) {
		return false
	} else if g.edges[e.start][e.end].ocurr == 1 {
		if g.edges[e.start][e.end].benefit == 0 {
			return false
		}
		return true
	} else {
		return true
	}
}

func (g *Graph) cumpleAcotamiento(e *Edge, mejorSol []*Edge) bool {
	beneficioE := e.benefit - e.cost
	beneficioSolParcial := beneficio(mejorSol) + beneficioE
	maxBeneficio := beneficioDisponible - int(math.Max(0, float64(beneficioE))) + beneficioSolParcial
	if maxBeneficio <= beneficio(mejorSol) {
		return false
	}
	return true
}

func (g *Graph) branchAndBound(e *Edge, path []*Edge) {
	if e.end == 1 {
		if beneficio(path) > beneficio(mejorSol) {
			mejorSol = path
		}
	}
	sucesores := g.Neighbors(e.end)
	for _, edge := range sucesores {
		if g.verifyConditions(g.edges[e.end][edge], path) {
			// g.AddEdge(g.edges[e.end][edge], solParcial)
			beneficioDisponible = beneficioDisponible - int(math.Max(0, float64(g.edges[e.end][edge].benefit-g.edges[e.end][edge].cost)))
			g.branchAndBound(g.edges[e.end][edge], solParcial)
		}
		// edge = eliminarUltimoLado(solParcial)
		beneficioDisponible = beneficioDisponible - int(math.Max(0, float64(g.edges[e.end][edge].benefit-g.edges[e.end][edge].cost)))
	}
}

func (g *Graph) checkNegativeCycle(e *Edge, solParcial []*Edge) bool {
	if g.edges[e.end] == nil {
		totalBenefit = e.benefit - e.cost
		for i := len(g.edges); i > 0; i-- {
			if solParcial[i].start == e.end {
				break
			} else {
				totalBenefit = totalBenefit + solParcial[i].benefit - solParcial[i].cost
			}
		}
	}
	if totalBenefit < 0 {
		return true
	}
	return false
}

func (g *Graph) verifyConditions(e *Edge, solParcial []*Edge) bool {
	return g.checkNegativeCycle(e, solParcial) && g.estaEnSolucionParcial(e) && g.cumpleAcotamiento(e, mejorSol)
}

// NewGraph function to create a new graph
func NewGraph(nodes int) *Graph {
	g := &Graph{make(map[int]map[int]*Edge), make([]*Edge, 0)}
	for i := 1; i <= nodes; i++ {
		tmap := make(map[int]*Edge)
		g.edges[i] = tmap
	}
	return g
}
