package main

import (
	wc "./weightedchoice"
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {

	file, _ := os.Open("./instanciasPRPP/CHRISTOFIDES/P15NoRPP")
	//file, _ := os.Open("./instanciasPRPP/RANDOM/R0NoRPP")
	lineScanner := bufio.NewScanner(file)
	line := 0
	g := NewGraph(1)
	for lineScanner.Scan() {
		contents := strings.Fields(lineScanner.Text())
		if line == 0 {
			number, _ := strconv.ParseInt(contents[len(contents)-1], 0, 0)
			g = NewGraph(int(number))
		}
		if _, err := strconv.Atoi(contents[0]); err == nil {
			startNode, _ := strconv.ParseInt(contents[0], 0, 0)
			endNode, _ := strconv.ParseInt(contents[1], 0, 0)
			cost, _ := strconv.ParseInt(contents[2], 0, 0)
			benefit, _ := strconv.ParseInt(contents[3], 0, 0)
			g.AddEdge(int(startNode), int(endNode), int(cost), int(benefit))
		}
		line++
	}

	var path []*Edge
	var initialEdge *Edge
	fmt.Println(g)
	fmt.Println("positivos", g.positiveEdges)

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
	} else {
		// Select first positive edge from depot
		for index, edge := range pEdges {
			if edge.start == 1 {
				initialEdge = edge
				pEdges = append(pEdges[:index], pEdges[index+1:]...) // Delete Edge from list
				break
			}
		}
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
			} else if pEdges[adjEdge].end == b {
				b = pEdges[adjEdge].start
			}
			path = append(path, pEdges[adjEdge])
			pEdges = append(pEdges[:adjEdge], pEdges[adjEdge+1:]...) // Delete Edge from list
		} else {
			ccm := make([][]*Edge, 0)
			for _, edge := range pEdges {
				ccm = append(ccm, g.Dijkstra(edge.start, b, path))
				ccm = append(ccm, g.Dijkstra(edge.end, b, path))
			}
			cmib := getPath(ccm)         // Probabilistic selection of the path
			path = append(path, cmib...) // Append random selected path to cycle
			// Remove cmib from pEdges
			i := 0
			for _, elem := range path {
				i = 0
				for i < len(pEdges) {
					if (elem.start == pEdges[i].start && elem.end == pEdges[i].end) ||
						(elem.end == pEdges[i].start && elem.start == pEdges[i].end) {
						pEdges = append(pEdges[:i], pEdges[i+1:]...) // Delete Edge from list
						break
					}
					i = i + 1
				}
			}
			b = path[len(path)-1].end
		}
	}
	if path[len(path)-1].end != 1 {
		minPath := g.Dijkstra(1, path[len(path)-1].end, path)
		path = append(path, minPath...)
	}
	total := 0
	for i := 0; i < len(path); i++ {
		if path[i].ocurr <= 1 {
			total = total + path[i].benefit - path[i].cost
			fmt.Println("nuevo   ", path[i])
		} else {
			total = total - path[i].cost
			fmt.Println("repetido", path[i])
		}
		g.AddOcurr(path[i].start, path[i].end)
	}
	fmt.Println("Ciclo: ", path)
	fmt.Println(total)
}

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
		weights[index] = elem.benefit - elem.cost
		if weights[index] == 0 {
			weights[index] = 1
		}
	}
	wc.Weights = weights
	random := wc.BinarySearch().(int)
	edge := positiveEdges[random]
	for edge.start != b && edge.end != b {
		rand.Seed(time.Now().UTC().UnixNano())
		random = wc.BinarySearch().(int)
		edge = positiveEdges[random]
	}
	return random
}

// Select randomly a Path with probability: BenefitPath/total(BenefitPath)
func getPath(ccm [][]*Edge) []*Edge {
	var total int
	pathCost := make([]int, len(ccm)) // Create array of path cost
	for index, path := range ccm {
		total = 0
		for _, edge := range path {
			total = total + edge.benefit - edge.cost
		}
		pathCost[index] = total
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
	for i := 0; i < len(pathCost); i++ {
		if pathCost[i] > 0 {
			pathCost[i] = pathCost[i] * int(math.Abs(float64(min)))
		} else {
			pathCost[i] = pathCost[i] + int(math.Abs(float64(min)))
		}
		if pathCost[i] == 0 {
			pathCost[i] = 1
		}
	}

	rand.Seed(time.Now().UTC().UnixNano())
	wc := new(wc.WeightedChoice)
	wc.Weights = pathCost             // Assign an array position a weight according to cost
	random := wc.BinarySearch().(int) // Random select an array position
	return ccm[random]
}
