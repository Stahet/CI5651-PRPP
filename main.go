package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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

	file, _ := os.Open("./instanciasPRPP/CHRISTOFIDES/P01NoRPP")
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
	fmt.Println(g)
}
