package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	wc "./weightedchoice"
)

func main() {

	file, _ := os.Open("./instanciasPRPP/CHRISTOFIDES/P01NoRPP")
	//file, _ := os.Open("./instanciasPRPP/RANDOM/R9NoRPP")
	//file, _ := os.Open("./instanciasPRPP/DEGREE/D2NoRPP")
	//file, _ := os.Open("./instanciasPRPP/GRID/G16NoRPP")
	lineScanner := bufio.NewScanner(file)
	line := 0
	g := NewGraph(1)
	maxBenefit, b, c := 0, 0, 0
	// branchG := NewGraph(1)
	for lineScanner.Scan() {
		contents := strings.Fields(lineScanner.Text())
		if line == 0 {
			number, _ := strconv.ParseInt(contents[len(contents)-1], 0, 0)
			g = NewGraph(int(number))
			// branchG = NewGraph(int(number))
		}
		if _, err := strconv.Atoi(contents[0]); err == nil {
			startNode, _ := strconv.ParseInt(contents[0], 0, 0)
			endNode, _ := strconv.ParseInt(contents[1], 0, 0)
			cost, _ := strconv.ParseInt(contents[2], 0, 0)
			benefit, _ := strconv.ParseInt(contents[3], 0, 0)
			g.AddEdge(int(startNode), int(endNode), int(cost), int(benefit))
			b = g.Benefit(int(startNode), int(endNode))
			c = g.Cost(int(startNode), int(endNode))
			if b-c >= 0 {
				maxBenefit = maxBenefit + b - c
			}
			// branchG.AddEdge(int(startNode), int(endNode), int(cost), int(benefit))
		}
		line++
	}

	var path []*Edge
	path = getCycleGRASP(g)

	fmt.Println("Ciclo Greedy: ", path)
	fmt.Println("Total: ", g.getPathBenefit(path))

	path = removeNegativeCycle(g, path)
	fmt.Println("Nuevo ciclo sin negativo: ", path)
	fmt.Println("Total: ", g.getPathBenefit(path))

	var branchSol []*Edge
	branchSol = make([]*Edge, 0, 0)
	//branchSol = g.branchAndBound(1, branchSol, path, maxBenefit)
	fmt.Println("Ciclo Branch and bound: ", branchSol)
	fmt.Println("Total: ", g.getPathBenefit(branchSol))
	// fmt.Println(g.obtenerListaSucesores(1))
}

// Remove negative cycle from a solution
func removeNegativeCycle(g *Graph, path []*Edge) []*Edge {
	start, end := 0, 0
	for i := 0; i < len(path); i++ {
		//fmt.Println("i=", i, "edge: ", path[i])
		start = path[i].start
		for j := i + 1; j < len(path); j++ {
			end = path[j].end
			if start == end {
				//fmt.Println("i,j=", i, j)
				if g.getPathBenefit(path[i:j]) <= 0 {
					fmt.Println("Removiendo ciclo negativo", path[i:j+1], "total: ", g.getPathBenefit(path[i:j+1]))
					path = append(path[:i], path[j+1:]...)
					//fmt.Println("New Path", path, "total: ", getPathBenefit(g, path))
					i = i - 1
					break
				}
			}
		}
	}
	return path
}

// Get a initial solution using GRASP based Algorithm
func getCycleGRASP(g *Graph) []*Edge {
	var path []*Edge
	var initialEdge *Edge
	fmt.Println("Imprimiendo grafo\n", g)
	fmt.Println("Lados positivos", g.positiveEdges)

	pEdges := g.positiveEdges
	// Check if deposit (1) in T
	// If no positive edge adjacent to deposit, select max benefit-cost from E
	if !inPositiveEdges(pEdges, 1) {
		max := (math.MaxInt32 - 1) * -1
		for _, node := range g.Neighbors(1) {
			if g.Benefit(1, node)-g.Cost(1, node) > max {
				max = g.Benefit(1, node) - g.Cost(1, node)
				initialEdge = g.edges[1][node]
			}
		}
		//fmt.Println("No hay lado positivo adyacente, seleccionando de E: ", initialEdge)
	} else {
		// Select first positive edge from depot
		for index, edge := range pEdges {
			if edge.start == 1 {
				initialEdge = edge
				pEdges = append(pEdges[:index], pEdges[index+1:]...) // Delete Edge from list
				break
			}
		}
		//fmt.Println("Existe un lado positivo adyacente, seleccionando: ", initialEdge)
	}

	// Set depot initial edge
	b := initialEdge.end
	path = append(path, initialEdge)
	var adjEdge int
	for len(pEdges) > 0 {
		if inPositiveEdges(pEdges, b) {
			adjEdge = getEdge(pEdges, b) // Get Edge position adjacent to node g
			if pEdges[adjEdge].start == b {
				b = pEdges[adjEdge].end
				path = append(path, g.edges[pEdges[adjEdge].start][b])
			} else if pEdges[adjEdge].end == b {
				b = pEdges[adjEdge].start
				path = append(path, g.edges[pEdges[adjEdge].end][b])
			}
			//fmt.Println("Lado positivo aleatorio seleccionado", pEdges[adjEdge])
			pEdges = append(pEdges[:adjEdge], pEdges[adjEdge+1:]...) // Delete Edge from list
			//fmt.Println("Lados positivos luego seleccion aleatorio", pEdges)
			//fmt.Println("b=", b)
		} else {
			ccm := make([][]*Edge, 0)
			for _, edge := range pEdges {
				ccm = append(ccm, g.Dijkstra(edge.start, b, path))
				ccm = append(ccm, g.Dijkstra(edge.end, b, path))
			}
			cmib := getPath(ccm)         // Probabilistic selection of the path
			path = append(path, cmib...) // Append random selected path to cycle
			// Remove cmib from pEdges
			for _, elem := range path {
				i := 0
				for i < len(pEdges) {
					if elem.equals(pEdges[i]) {
						pEdges = append(pEdges[:i], pEdges[i+1:]...) // Delete Edge from list
						break
					}
					i = i + 1
				}
			}
			if path[len(path)-1].end == path[len(path)-2].end {
				b = path[len(path)-1].start
			} else {
				b = path[len(path)-1].end
			}
			//fmt.Println("Camino aleatorio seleccionado", cmib)
			//fmt.Println("Lados positivos luego seleccion aleatorio de camino", pEdges)
		}
	}
	// Check if last is depot
	if path[len(path)-1].end != 1 {
		minPath := g.Dijkstra(1, b, path)
		path = append(path, minPath...)
	}
	return path
}

//Check if node is in positiveEdges set
func inPositiveEdges(positiveEdges []*Edge, node int) bool {
	for _, edge := range positiveEdges {
		if edge.start == node || edge.end == node {
			return true
		}
	}
	return false
}

// Select randomly an Edge with probability: benefit-cost/total(benefit-cost)
func getEdge(positiveEdges []*Edge, b int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	wc := new(wc.WeightedChoice)
	weights := make([]int, len(positiveEdges))
	// Construct weight array
	for index, elem := range positiveEdges {
		weights[index] = elem.benefit - elem.cost + 1
	}
	//fmt.Println("Lados probables a seleccionar getEdge", positiveEdges)
	//fmt.Println("Pesos getEdge:", weights)
	wc.Weights = weights
	random := wc.BinarySearch().(int) // Binary search random number with weight
	edge := positiveEdges[random]
	// Random until start or end == b
	for edge.start != b && edge.end != b {
		random = wc.BinarySearch().(int)
		edge = positiveEdges[random]
	}
	return random
}

// Select randomly a Path with probability: BenefitPath/total(BenefitPath)
func getPath(ccm [][]*Edge) []*Edge {
	var total int
	pathCost := make([]int, len(ccm)) // Create array of path cost
	//fmt.Println("caminos")
	for index, path := range ccm {
		total = 0
		for _, edge := range path {
			total = total + edge.benefit - edge.cost
		}
		pathCost[index] = total
		//fmt.Println(path, "Cost: ", total)
	}

	// Because there is a possibility to get a negative total cost
	// We find the minimum
	// For positive number we multiply minimum
	// For negative number we add abs(number) + 1
	min := math.MaxInt32
	for _, elem := range pathCost {
		if elem < min {
			min = elem
		}
	}
	if min == 0 {
		min = 1
	}
	for i := 0; i < len(pathCost); i++ {
		if pathCost[i] > 0 {
			pathCost[i] = pathCost[i]*int(math.Abs(float64(min))) + 1
		} else {
			pathCost[i] = pathCost[i] + int(math.Abs(float64(min))) + 1
		}

	}
	//fmt.Println("pathCost", pathCost)
	rand.Seed(time.Now().UTC().UnixNano())
	wc := new(wc.WeightedChoice)
	wc.Weights = pathCost             // Assign an array position a weight according to cost
	random := wc.BinarySearch().(int) // Random select an array position
	return ccm[random]
}

func beneficio(solucion []*Edge) int {
	total := 0
	for _, edge := range solucion {
		if edge.ocurr > 0 {
			total = total - edge.cost
		} else {
			total = total + edge.benefit - edge.cost
		}
	}
	return total
}
