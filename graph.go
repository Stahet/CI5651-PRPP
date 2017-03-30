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

type Edges []Edge

func (slice Edges) Len() int {
	return len(slice)
}

func (slice Edges) Less(i, j int) bool {
	return (float64(slice[i].benefit) - float64(slice[i].cost)) < (float64(slice[j].benefit) - float64(slice[j].cost))
}

func (slice Edges) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// AddEdge function to ad edges to a graph
func (g *Graph) AddEdge(start, end, cost, benefit int) {

	if _, ok := g.edges[start][end]; !ok {
		g.edges[start][end] = &Edge{start, end, cost, benefit, 0}
	}

	if _, ok := g.edges[end][start]; !ok {
		g.edges[end][start] = &Edge{end, start, cost, benefit, 0}
	}

	if benefit-cost >= 0 {
		g.positiveEdges = append(g.positiveEdges, g.edges[start][end])
	}
}

// RemoveEdge function to remove edge from a graph
func (g *Graph) RemoveEdge(start, end int) {
	if _, ok := g.edges[start][end]; ok {
		delete(g.edges[start], end)
	}
	if _, ok := g.edges[end][start]; ok {
		delete(g.edges[end], start)
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

func (g *Graph) ResetOcurr() {
	for _, v := range g.edges {
		for _, edge := range v {
			edge.ocurr = 0
		}
	}
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
		// if u == to {
		// 	break
		// }

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
		if edge.equals(elem) {
			return edge.cost
		}
	}
	if edge.benefit-edge.cost < 0 {
		return (edge.benefit - edge.cost) * (-1)
	}
	return 0
}

// Check if 2 edges are equals
func (e1 *Edge) equals(e2 *Edge) bool {
	return (e1.start == e2.start && e1.end == e2.end) || (e1.start == e2.end && e1.end == e2.start)
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

func (g *Graph) estaEnSolucionParcial(e *Edge, solParcial []*Edge) bool {
	ocurr := 0
	for _, edge := range solParcial {
		if edge.equals(e) {
			ocurr = ocurr + 1
		}
		if ocurr >= 2 {
			return true
		}
	}
	return false
}

func (g *Graph) obtenerListaSucesores(v int) []*Edge {
	lista := Edges{}
	newlist := make([]*Edge, 0, len(g.edges[v])*2)
	for _, edge := range g.edges[v] {
		lista = append(lista, Edge{edge.start, edge.end, edge.cost, edge.benefit, 0})
		lista = append(lista, Edge{edge.start, edge.end, edge.cost, 0, 0})
	}
	sort.Sort(sort.Reverse(lista))
	for _, edge := range lista {
		newlist = append(newlist, g.edges[edge.start][edge.end])
	}
	return newlist
}

func (g *Graph) cumpleAcotamiento(e *Edge, solParcial []*Edge, mejorSol []*Edge, beneficioDisponible int) bool {
	beneficioE := e.benefit - e.cost
	beneficioSolParcial := g.getPathBenefit(solParcial) + beneficioE
	maxBeneficio := beneficioDisponible - int(math.Max(0, float64(beneficioE))) + beneficioSolParcial
	if maxBeneficio <= g.getPathBenefit(mejorSol) {
		return false
	}
	return true
}

func (g *Graph) branchAndBound(e int, solParcial []*Edge, mejorSol []*Edge, beneficioDisponible int) []*Edge {
	if e == 1 {
		if g.getPathBenefit(solParcial) > g.getPathBenefit(mejorSol) {
			mejorSol = solParcial
		}
	}
	fmt.Println(mejorSol)
	sucesores := g.obtenerListaSucesores(e)
	for _, edge := range sucesores {
		// if g.verifyConditions(g.edges[e][edge], solParcial, mejorSol) {
		if !(g.estaEnSolucionParcial(edge, solParcial)) && g.cumpleAcotamiento(edge, solParcial, mejorSol, beneficioDisponible) && !g.checkNegativeCycle(edge, solParcial) {
			solParcial = append(solParcial, edge)
			// g.AddEdge(g.edges[e.end][edge], solParcial)
			beneficioDisponible = beneficioDisponible - int(math.Max(0, float64(edge.benefit-edge.cost)))
			mejorSol = g.branchAndBound(edge.end, solParcial, mejorSol, beneficioDisponible)
		}
	}
	solParcial = solParcial[:len(solParcial)-1] // edge = eliminarUltimoLado(solParcial)
	ultimo := solParcial[len(solParcial)-1]
	beneficioDisponible = beneficioDisponible + int(math.Max(0, float64(ultimo.benefit-ultimo.cost)))
	return mejorSol
}

func (g *Graph) checkNegativeCycle(e *Edge, solParcial []*Edge) bool {
	path := append(solParcial, e)
	for index, edge := range solParcial {
		if edge.start == e.end {
			totalBenefit = g.getPathBenefit(path[index:])
			if totalBenefit < 0 {
				return true
			}
		}
	}
	return false
}

func (g *Graph) getPathBenefit(path []*Edge) int {
	g.ResetOcurr()
	total := 0
	for _, edge := range path {
		if edge.ocurr <= 0 {
			total = total + edge.benefit - edge.cost
			//fmt.Println("nuevo   ", edge, "total: ", total)
		} else {
			total = total - edge.cost
			//fmt.Println("repetido", edge, "total: ", total)
		}
		g.AddOcurr(edge.start, edge.end)
	}
	return total
}

// func (g *Graph) verifyConditions(e *Edge, solParcial []*Edge, mejorSol []*Edge) bool {
// return g.checkNegativeCycle(e, solParcial) && g.estaEnSolucionParcial(e, solParcial) && g.cumpleAcotamiento(e, mejorSol,)
// }

// NewGraph function to create a new graph
func NewGraph(nodes int) *Graph {
	g := &Graph{make(map[int]map[int]*Edge), make([]*Edge, 0)}
	for i := 1; i <= nodes; i++ {
		tmap := make(map[int]*Edge)
		g.edges[i] = tmap
	}
	return g
}
