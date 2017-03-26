package main

import (
	wc "./weightedchoice"
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {

	// g := NewGraph(5)
	// g.AddEdge(1, 2, 10, 4)
	// g.AddEdge(4, 2, 10, 4)
	// g.AddEdge(3, 2, 4, 80)
	// g.AddEdge(3, 4, 7, 9)
	// g.AddEdge(3, 4, 7, 9)
	// g.AddEdge(4, 3, 7, 4)
	// fmt.Println(g)
	// // g.RemoveEdge(2, 1)
	// // g.RemoveEdge(3, 4)
	// // g.RemoveEdge(4, 3)
	// // g.RemoveEdge(3, 4)
	// // fmt.Println(g)
	// fmt.Println(g.Neighbors(2))
	// fmt.Println("4-3Cost ", g.Cost(4, 3), "Benefit ", g.Benefit(3, 4))
	// fmt.Println("1-2Cost", g.Cost(1, 2), "Benefit ", g.Benefit(2, 1))
	// fmt.Println("Node 2 Degree", g.Degree(2))
	// fmt.Println("Node 5 Degree", g.Degree(5))

	file, _ := os.Open("./instanciasPRPP/CHRISTOFIDES/P02NoRPP")
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

	var path []int
	path = append(path, 1)
	//fmt.Println(g)
	fmt.Println("positivos", g.positiveEdges)

	pEdges := g.positiveEdges
	// Check if deposit (1) in T
	// If no positive edge adjacent to deposit, select max benefit-cost from E
	if !inPositiveEdges(pEdges, 1) {
		var maxEdge *Edge
		max := -99999999
		for _, node := range g.Neighbors(1) {
			if g.Benefit(1, node)-g.Benefit(1, node) > max {
				max = g.Benefit(1, node) - g.Benefit(1, node)
				maxEdge = g.edges[1][node]
			}
		}
		pEdges = append(pEdges, maxEdge)
	}
	b := 1
	var adjEdge int
	for len(pEdges) > 0 {
		if inPositiveEdges(pEdges, b) {
			adjEdge = getEdge(pEdges, b)
			if pEdges[adjEdge].start == b {
				b = pEdges[adjEdge].end
				path = append(path, b)
			} else if pEdges[adjEdge].end == b {
				b = pEdges[adjEdge].start
				path = append(path, b)
			}
			fmt.Println("delete: ", pEdges[adjEdge], pEdges)
			pEdges = append(pEdges[:adjEdge], pEdges[adjEdge+1:]...)
			fmt.Println(path)
		}
	}
	getEdge(pEdges, 3)
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
	var adjEdges []*Edge
	for _, edge := range positiveEdges {
		if edge.start == b || edge.end == b {
			adjEdges = append(adjEdges, edge)
		}
	}

	rand.Seed(time.Now().UTC().UnixNano())
	wc := new(wc.WeightedChoice)
	weights := make([]int, len(adjEdges))
	for index, elem := range adjEdges {
		weights[index] = elem.benefit - elem.cost
	}
	wc.Weights = weights
	random := wc.BinarySearch().(int)
	// fmt.Printf("result: %s %d\n", adjEdges[random], weights[random])
	// fmt.Println(positiveEdges, "hola")
	return random
}
