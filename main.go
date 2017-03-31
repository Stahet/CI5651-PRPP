package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var path []*Edge
var mejorSol []*Edge
var beneficioDisponible int
var solParcial []*Edge
var totalBenefit int

func main() {

	file, _ := os.Open("./instanciasPRPP/CHRISTOFIDES/P11NoRPP")
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
		}
		line++
	}

	path = getCycleGRASP(g)

	fmt.Println("Ciclo Greedy: ", path)
	fmt.Println("Total: ", getPathBenefit(path))

	path = removeNegativeCycle(g, path)
	fmt.Println("Nuevo ciclo sin negativo: ", path)
	fmt.Println("Total: ", getPathBenefit(path))

	var branchSol []*Edge
	branchSol = make([]*Edge, 0, 0)
	//_, branchSol, _ = g.branchAndBound(1, branchSol, path, maxBenefit)
	fmt.Println("Ciclo Branch and bound: ", branchSol)
	fmt.Println("Total: ", getPathBenefit(branchSol))
	// fmt.Println(g.obtenerListaSucesores(1))
}
